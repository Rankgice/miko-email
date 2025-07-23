package imap

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"nbemail/internal/config"
)

// Server IMAP服务器
type Server struct {
	config   *config.Config
	db       *sql.DB
	listener net.Listener
	quit     chan struct{}
}

// IMAPSession IMAP会话
type IMAPSession struct {
	conn     net.Conn
	reader   *bufio.Reader
	writer   *bufio.Writer
	server   *Server
	state    string // NOTAUTHENTICATED, AUTHENTICATED, SELECTED
	user     string
	mailbox  string
	tag      string
}

// NewServer 创建新的IMAP服务器
func NewServer(cfg *config.Config, db *sql.DB) *Server {
	return &Server{
		config: cfg,
		db:     db,
		quit:   make(chan struct{}),
	}
}

// Start 启动IMAP服务器
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.config.IMAPPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("监听端口失败: %v", err)
	}

	s.listener = listener
	log.Printf("IMAP服务器启动成功，监听端口: %d", s.config.IMAPPort)

	for {
		select {
		case <-s.quit:
			return nil
		default:
			conn, err := listener.Accept()
			if err != nil {
				if strings.Contains(err.Error(), "use of closed network connection") {
					return nil
				}
				log.Printf("接受连接失败: %v", err)
				continue
			}
			go s.handleConnection(conn)
		}
	}
}

// Stop 停止IMAP服务器
func (s *Server) Stop() {
	close(s.quit)
	if s.listener != nil {
		s.listener.Close()
	}
}

// handleConnection 处理IMAP连接
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	log.Printf("新的IMAP连接: %s", conn.RemoteAddr())

	// 设置连接超时
	conn.SetDeadline(time.Now().Add(30 * time.Minute))

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	session := &IMAPSession{
		conn:   conn,
		reader: reader,
		writer: writer,
		server: s,
		state:  "NOTAUTHENTICATED",
	}

	// 发送欢迎消息
	session.writeResponse("* OK NBEmail IMAP Server Ready")

	session.handle()
}

// handle 处理IMAP会话
func (session *IMAPSession) handle() {
	for {
		line, _, err := session.reader.ReadLine()
		if err != nil {
			log.Printf("读取命令失败: %v", err)
			return
		}

		command := strings.TrimSpace(string(line))
		if command == "" {
			continue
		}

		log.Printf("IMAP命令: %s", command)

		parts := strings.Fields(command)
		if len(parts) < 2 {
			session.writeResponse("* BAD Invalid command")
			continue
		}

		tag := parts[0]
		cmd := strings.ToUpper(parts[1])
		args := parts[2:]

		session.tag = tag

		switch cmd {
		case "CAPABILITY":
			session.handleCapability()
		case "LOGIN":
			session.handleLogin(args)
		case "LIST":
			session.handleList(args)
		case "SELECT":
			session.handleSelect(args)
		case "SEARCH":
			session.handleSearch(args)
		case "FETCH":
			session.handleFetch(args)
		case "LOGOUT":
			session.handleLogout()
			return
		default:
			session.writeTaggedResponse("BAD Command not implemented")
		}
	}
}

// writeResponse 写入响应
func (session *IMAPSession) writeResponse(response string) {
	session.writer.WriteString(response + "\r\n")
	session.writer.Flush()
}

// writeTaggedResponse 写入带标签的响应
func (session *IMAPSession) writeTaggedResponse(response string) {
	session.writeResponse(fmt.Sprintf("%s %s", session.tag, response))
}

// handleCapability 处理CAPABILITY命令
func (session *IMAPSession) handleCapability() {
	session.writeResponse("* CAPABILITY IMAP4rev1 LOGIN")
	session.writeTaggedResponse("OK CAPABILITY completed")
}

// handleLogin 处理LOGIN命令
func (session *IMAPSession) handleLogin(args []string) {
	if len(args) < 2 {
		session.writeTaggedResponse("BAD LOGIN requires username and password")
		return
	}

	username := strings.Trim(args[0], "\"")
	password := strings.Trim(args[1], "\"")

	// 验证用户凭据 - 从数据库获取加密密码进行验证
	var storedPassword string
	err := session.server.db.QueryRow("SELECT password FROM users WHERE email = ?", username).Scan(&storedPassword)
	if err != nil {
		log.Printf("IMAP登录失败 - 用户不存在: %s, 错误: %v", username, err)
		session.writeTaggedResponse("NO LOGIN failed")
		return
	}

	// 验证密码 - 支持明文密码（向后兼容）和bcrypt加密密码
	var passwordValid bool
	if strings.HasPrefix(storedPassword, "$2") { // bcrypt加密密码
		// 使用bcrypt验证
		err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
		passwordValid = (err == nil)
	} else {
		// 明文密码比较（向后兼容）
		passwordValid = (storedPassword == password)
	}

	if !passwordValid {
		log.Printf("IMAP登录失败 - 密码错误: %s", username)
		session.writeTaggedResponse("NO LOGIN failed")
		return
	}

	session.state = "AUTHENTICATED"
	session.user = username
	log.Printf("IMAP登录成功: %s", username)
	session.writeTaggedResponse("OK LOGIN completed")
}

// handleList 处理LIST命令
func (session *IMAPSession) handleList(args []string) {
	if session.state != "AUTHENTICATED" && session.state != "SELECTED" {
		session.writeTaggedResponse("NO Not authenticated")
		return
	}

	session.writeResponse("* LIST () \"/\" \"INBOX\"")
	session.writeTaggedResponse("OK LIST completed")
}

// handleSelect 处理SELECT命令
func (session *IMAPSession) handleSelect(args []string) {
	if session.state != "AUTHENTICATED" && session.state != "SELECTED" {
		session.writeTaggedResponse("NO Not authenticated")
		return
	}

	if len(args) < 1 {
		session.writeTaggedResponse("BAD SELECT requires mailbox name")
		return
	}

	mailbox := strings.Trim(args[0], "\"")
	if strings.ToUpper(mailbox) != "INBOX" {
		session.writeTaggedResponse("NO Mailbox does not exist")
		return
	}

	// 获取邮件数量 - 使用精确匹配
	var count int
	err := session.server.db.QueryRow("SELECT COUNT(*) FROM emails WHERE to_addr = ?", session.user).Scan(&count)
	if err != nil {
		count = 0
	}

	session.state = "SELECTED"
	session.mailbox = "INBOX"

	session.writeResponse(fmt.Sprintf("* %d EXISTS", count))
	session.writeResponse("* 0 RECENT")
	session.writeResponse("* OK [UIDVALIDITY 1] UIDs valid")
	session.writeTaggedResponse("OK [READ-WRITE] SELECT completed")
}

// handleFetch 处理FETCH命令
func (session *IMAPSession) handleFetch(args []string) {
	if session.state != "SELECTED" {
		session.writeTaggedResponse("NO Not selected")
		return
	}

	if len(args) < 2 {
		session.writeTaggedResponse("BAD FETCH requires sequence set and data items")
		return
	}

	// 解析参数
	// sequenceSet := args[0] // 暂时不使用序列集
	dataItems := strings.Join(args[1:], " ")

	// 查询邮件数据，包含完整的邮件内容 - 使用精确匹配
	rows, err := session.server.db.Query("SELECT id, from_addr, to_addr, subject, body, created_at FROM emails WHERE to_addr = ? ORDER BY created_at DESC LIMIT 10", session.user)
	if err != nil {
		session.writeTaggedResponse("NO FETCH failed")
		return
	}
	defer rows.Close()

	seqNum := 1
	for rows.Next() {
		var id int
		var sender, recipient, subject, body, createdAt string
		if err := rows.Scan(&id, &sender, &recipient, &subject, &body, &createdAt); err != nil {
			continue
		}

		// 根据请求的数据项返回不同的信息
		if strings.Contains(dataItems, "RFC822") {
			// 构造完整的邮件内容
			emailContent := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nDate: %s\r\n\r\n%s",
				sender, recipient, subject, createdAt, body)

			// 返回RFC822格式的邮件 - 正确的IMAP格式
			// 格式: * seqnum FETCH (RFC822 {size}\r\n<content>)
			session.writeResponse(fmt.Sprintf("* %d FETCH (RFC822 {%d}", seqNum, len(emailContent)))
			session.writer.WriteString(emailContent)
			session.writer.WriteString(")\r\n")
			session.writer.Flush()
		} else {
			// 返回基本信息
			session.writeResponse(fmt.Sprintf("* %d FETCH (UID %d RFC822.SIZE %d ENVELOPE (\"%s\" \"%s\" ((\"%s\" NIL \"%s\" NIL)) NIL NIL NIL NIL NIL))",
				seqNum, id, len(body), createdAt, subject, sender, sender))
		}
		seqNum++
	}

	session.writeTaggedResponse("OK FETCH completed")
}

// handleSearch 处理SEARCH命令
func (session *IMAPSession) handleSearch(args []string) {
	if session.state != "SELECTED" {
		session.writeTaggedResponse("NO Not selected")
		return
	}

	// 简化实现，支持基本的SEARCH命令
	// 目前只支持 "ALL" 搜索
	if len(args) == 0 || strings.ToUpper(args[0]) != "ALL" {
		session.writeTaggedResponse("BAD SEARCH command requires ALL parameter")
		return
	}

	// 查询用户的所有邮件ID - 使用精确匹配
	rows, err := session.server.db.Query("SELECT id FROM emails WHERE to_addr = ? ORDER BY created_at DESC", session.user)
	if err != nil {
		log.Printf("SEARCH查询失败: %v", err)
		session.writeTaggedResponse("NO SEARCH failed")
		return
	}
	defer rows.Close()

	var emailIDs []string
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			continue
		}
		emailIDs = append(emailIDs, fmt.Sprintf("%d", id))
	}

	// 返回搜索结果
	if len(emailIDs) > 0 {
		session.writeResponse("* SEARCH " + strings.Join(emailIDs, " "))
	} else {
		session.writeResponse("* SEARCH")
	}
	session.writeTaggedResponse("OK SEARCH completed")
}

// handleLogout 处理LOGOUT命令
func (session *IMAPSession) handleLogout() {
	session.writeResponse("* BYE NBEmail IMAP Server logging out")
	session.writeTaggedResponse("OK LOGOUT completed")
}
