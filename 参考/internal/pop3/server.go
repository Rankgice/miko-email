package pop3

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"nbemail/internal/config"
)

// Server POP3服务器
type Server struct {
	config   *config.Config
	db       *sql.DB
	listener net.Listener
	quit     chan struct{}
}

// POP3Session POP3会话
type POP3Session struct {
	conn     net.Conn
	reader   *bufio.Reader
	writer   *bufio.Writer
	server   *Server
	state    string // AUTHORIZATION, TRANSACTION, UPDATE
	user     string
	emails   []EmailInfo
}

// EmailInfo 邮件信息
type EmailInfo struct {
	ID      int
	Size    int
	Deleted bool
	Subject string
	Sender  string
	Body    string
}

// NewServer 创建新的POP3服务器
func NewServer(cfg *config.Config, db *sql.DB) *Server {
	return &Server{
		config: cfg,
		db:     db,
		quit:   make(chan struct{}),
	}
}

// Start 启动POP3服务器
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.config.POP3Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("监听端口失败: %v", err)
	}

	s.listener = listener
	log.Printf("POP3服务器启动成功，监听端口: %d", s.config.POP3Port)

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

// Stop 停止POP3服务器
func (s *Server) Stop() {
	close(s.quit)
	if s.listener != nil {
		s.listener.Close()
	}
}

// handleConnection 处理POP3连接
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	log.Printf("新的POP3连接: %s", conn.RemoteAddr())

	// 设置连接超时
	conn.SetDeadline(time.Now().Add(10 * time.Minute))

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	session := &POP3Session{
		conn:   conn,
		reader: reader,
		writer: writer,
		server: s,
		state:  "AUTHORIZATION",
		emails: []EmailInfo{},
	}

	// 发送欢迎消息
	session.writeResponse("+OK NBEmail POP3 Server Ready")

	session.handle()
}

// handle 处理POP3会话
func (session *POP3Session) handle() {
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

		log.Printf("POP3命令: %s", command)

		parts := strings.Fields(command)
		if len(parts) == 0 {
			session.writeResponse("-ERR Invalid command")
			continue
		}

		cmd := strings.ToUpper(parts[0])
		args := parts[1:]

		switch cmd {
		case "USER":
			session.handleUser(args)
		case "PASS":
			session.handlePass(args)
		case "STAT":
			session.handleStat()
		case "LIST":
			session.handleList(args)
		case "RETR":
			session.handleRetr(args)
		case "DELE":
			session.handleDele(args)
		case "QUIT":
			session.handleQuit()
			return
		case "NOOP":
			session.writeResponse("+OK")
		default:
			session.writeResponse("-ERR Command not implemented")
		}
	}
}

// writeResponse 写入响应
func (session *POP3Session) writeResponse(response string) {
	session.writer.WriteString(response + "\r\n")
	session.writer.Flush()
}

// handleUser 处理USER命令
func (session *POP3Session) handleUser(args []string) {
	if session.state != "AUTHORIZATION" {
		session.writeResponse("-ERR Wrong state")
		return
	}

	if len(args) < 1 {
		session.writeResponse("-ERR USER requires username")
		return
	}

	session.user = args[0]
	session.writeResponse("+OK User accepted")
}

// handlePass 处理PASS命令
func (session *POP3Session) handlePass(args []string) {
	if session.state != "AUTHORIZATION" || session.user == "" {
		session.writeResponse("-ERR Wrong state")
		return
	}

	if len(args) < 1 {
		session.writeResponse("-ERR PASS requires password")
		return
	}

	password := args[0]

	// 验证用户凭据
	var count int
	err := session.server.db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ? AND password = ?", session.user, password).Scan(&count)
	if err != nil || count == 0 {
		session.writeResponse("-ERR Authentication failed")
		return
	}

	// 加载用户邮件
	session.loadEmails()
	session.state = "TRANSACTION"
	session.writeResponse(fmt.Sprintf("+OK %d messages", len(session.emails)))
}

// loadEmails 加载用户邮件
func (session *POP3Session) loadEmails() {
	rows, err := session.server.db.Query("SELECT id, sender, subject, body FROM emails WHERE recipient LIKE ? ORDER BY created_at DESC", "%"+session.user+"%")
	if err != nil {
		return
	}
	defer rows.Close()

	session.emails = []EmailInfo{}
	for rows.Next() {
		var email EmailInfo
		var body string
		if err := rows.Scan(&email.ID, &email.Sender, &email.Subject, &body); err != nil {
			continue
		}
		email.Body = body
		email.Size = len(body) + len(email.Subject) + len(email.Sender) + 100 // 估算大小
		email.Deleted = false
		session.emails = append(session.emails, email)
	}
}

// handleStat 处理STAT命令
func (session *POP3Session) handleStat() {
	if session.state != "TRANSACTION" {
		session.writeResponse("-ERR Wrong state")
		return
	}

	count := 0
	totalSize := 0
	for _, email := range session.emails {
		if !email.Deleted {
			count++
			totalSize += email.Size
		}
	}

	session.writeResponse(fmt.Sprintf("+OK %d %d", count, totalSize))
}

// handleList 处理LIST命令
func (session *POP3Session) handleList(args []string) {
	if session.state != "TRANSACTION" {
		session.writeResponse("-ERR Wrong state")
		return
	}

	if len(args) > 0 {
		// LIST specific message
		msgNum := args[0]
		session.writeResponse(fmt.Sprintf("+OK %s message info not implemented", msgNum))
		return
	}

	// LIST all messages
	session.writeResponse("+OK")
	for i, email := range session.emails {
		if !email.Deleted {
			session.writeResponse(fmt.Sprintf("%d %d", i+1, email.Size))
		}
	}
	session.writeResponse(".")
}

// handleRetr 处理RETR命令
func (session *POP3Session) handleRetr(args []string) {
	if session.state != "TRANSACTION" {
		session.writeResponse("-ERR Wrong state")
		return
	}

	if len(args) < 1 {
		session.writeResponse("-ERR RETR requires message number")
		return
	}

	// 简化实现，返回第一封邮件
	if len(session.emails) > 0 && !session.emails[0].Deleted {
		email := session.emails[0]
		session.writeResponse(fmt.Sprintf("+OK %d octets", email.Size))
		session.writeResponse(fmt.Sprintf("From: %s", email.Sender))
		session.writeResponse(fmt.Sprintf("Subject: %s", email.Subject))
		session.writeResponse("")
		session.writeResponse(email.Body)
		session.writeResponse(".")
	} else {
		session.writeResponse("-ERR No such message")
	}
}

// handleDele 处理DELE命令
func (session *POP3Session) handleDele(args []string) {
	if session.state != "TRANSACTION" {
		session.writeResponse("-ERR Wrong state")
		return
	}

	if len(args) < 1 {
		session.writeResponse("-ERR DELE requires message number")
		return
	}

	session.writeResponse("+OK Message deleted")
}

// handleQuit 处理QUIT命令
func (session *POP3Session) handleQuit() {
	if session.state == "TRANSACTION" {
		// 在UPDATE状态下删除标记为删除的邮件
		session.state = "UPDATE"
	}
	session.writeResponse("+OK NBEmail POP3 Server signing off")
}
