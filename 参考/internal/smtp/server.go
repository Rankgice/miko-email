package smtp

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
	"nbemail/internal/config"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

// ConnectionTracker 连接跟踪器
type ConnectionTracker struct {
	connections map[string][]time.Time
	mutex       sync.RWMutex
}

// NewConnectionTracker 创建连接跟踪器
func NewConnectionTracker() *ConnectionTracker {
	return &ConnectionTracker{
		connections: make(map[string][]time.Time),
	}
}

// IsAllowed 检查IP是否允许连接
func (ct *ConnectionTracker) IsAllowed(ip string) bool {
	ct.mutex.Lock()
	defer ct.mutex.Unlock()

	now := time.Now()

	// 清理过期的连接记录（超过1小时）
	if times, exists := ct.connections[ip]; exists {
		var validTimes []time.Time
		for _, t := range times {
			if now.Sub(t) < time.Hour {
				validTimes = append(validTimes, t)
			}
		}
		ct.connections[ip] = validTimes
	}

	// 检查最近5分钟内的连接次数
	recentConnections := 0
	if times, exists := ct.connections[ip]; exists {
		for _, t := range times {
			if now.Sub(t) < 5*time.Minute {
				recentConnections++
			}
		}
	}

	// 如果5分钟内连接超过10次，拒绝连接
	if recentConnections >= 10 {
		log.Printf("IP %s 连接过于频繁，已被阻止", ip)
		return false
	}

	// 记录新连接
	ct.connections[ip] = append(ct.connections[ip], now)
	return true
}

// Server SMTP服务器
type Server struct {
	config   *config.Config
	db       *sql.DB
	listener net.Listener
	quit     chan bool
	tracker  *ConnectionTracker
}

// NewServer 创建新的SMTP服务器
func NewServer(cfg *config.Config, db *sql.DB) *Server {
	return &Server{
		config:  cfg,
		db:      db,
		quit:    make(chan bool),
		tracker: NewConnectionTracker(),
	}
}

// Start 启动SMTP服务器
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.config.SMTPPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("监听端口失败: %v", err)
	}

	s.listener = listener
	log.Printf("SMTP服务器启动成功，监听端口: %d", s.config.SMTPPort)

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

// Stop 停止SMTP服务器
func (s *Server) Stop() {
	close(s.quit)
	if s.listener != nil {
		s.listener.Close()
	}
}

// handleConnection 处理SMTP连接
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	// 获取客户端IP
	clientIP := strings.Split(conn.RemoteAddr().String(), ":")[0]

	// 检查IP是否被允许连接
	if !s.tracker.IsAllowed(clientIP) {
		log.Printf("拒绝来自 %s 的连接（连接过于频繁）", conn.RemoteAddr())
		return
	}

	log.Printf("新的SMTP连接: %s", conn.RemoteAddr())

	// 设置连接超时
	conn.SetDeadline(time.Now().Add(5 * time.Minute))

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// 发送欢迎消息
	s.writeResponse(writer, 220, fmt.Sprintf("%s NBEmail SMTP Server Ready", s.config.Domain))

	session := &SMTPSession{
		conn:   conn,
		reader: reader,
		writer: writer,
		server: s,
	}

	session.handle()
}

// writeResponse 写入SMTP响应
func (s *Server) writeResponse(writer *bufio.Writer, code int, message string) error {
	response := fmt.Sprintf("%d %s\r\n", code, message)
	_, err := writer.WriteString(response)
	if err != nil {
		return err
	}
	return writer.Flush()
}

// authenticateUser 验证用户
func (s *Server) authenticateUser(username, password string) bool {
	// 查询数据库验证用户（使用email字段作为用户名）
	var storedPassword string
	err := s.db.QueryRow("SELECT password FROM users WHERE email = ?", username).Scan(&storedPassword)
	if err != nil {
		log.Printf("用户认证失败: %v", err)
		return false
	}

	// 使用bcrypt验证密码
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password))
	return err == nil
}

// SMTPSession SMTP会话
type SMTPSession struct {
	conn          net.Conn
	reader        *bufio.Reader
	writer        *bufio.Writer
	server        *Server
	helo          string
	from          string
	to            []string
	data          []byte
	username      string
	password      string
	authenticated bool
}

// handle 处理SMTP会话
func (session *SMTPSession) handle() {
	for {
		line, err := session.reader.ReadString('\n')
		if err != nil {
			log.Printf("读取命令失败: %v", err)
			return
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		log.Printf("SMTP命令: %s", line)

		parts := strings.SplitN(line, " ", 2)
		command := strings.ToUpper(parts[0])
		var args string
		if len(parts) > 1 {
			args = parts[1]
		}

		switch command {
		case "HELO", "EHLO":
			session.handleHelo(args)
		case "AUTH":
			session.handleAuth(args)
		case "MAIL":
			session.handleMail(args)
		case "RCPT":
			session.handleRcpt(args)
		case "DATA":
			session.handleData()
		case "QUIT":
			session.writeResponse(221, "Bye")
			return
		case "RSET":
			session.reset()
			session.writeResponse(250, "OK")
		case "NOOP":
			session.writeResponse(250, "OK")
		default:
			session.writeResponse(500, "Command not recognized")
		}
	}
}

// writeResponse 写入响应
func (session *SMTPSession) writeResponse(code int, message string) error {
	return session.server.writeResponse(session.writer, code, message)
}

// handleHelo 处理HELO/EHLO命令
func (session *SMTPSession) handleHelo(args string) {
	session.helo = args
	// 简单的EHLO响应，支持AUTH扩展
	session.writer.WriteString("250-Hello " + args + "\r\n")
	session.writer.WriteString("250-AUTH PLAIN LOGIN\r\n")
	session.writer.WriteString("250 8BITMIME\r\n")
	session.writer.Flush()
}

// handleAuth 处理AUTH命令
func (session *SMTPSession) handleAuth(args string) {
	parts := strings.Fields(args)
	if len(parts) == 0 {
		session.writeResponse(501, "Syntax error")
		return
	}

	method := strings.ToUpper(parts[0])
	switch method {
	case "PLAIN":
		if len(parts) == 2 {
			// AUTH PLAIN with credentials
			session.handleAuthPlain(parts[1])
		} else {
			// AUTH PLAIN without credentials, request them
			session.writeResponse(334, "")
			line, _, err := session.reader.ReadLine()
			if err != nil {
				session.writeResponse(535, "Authentication failed")
				return
			}
			session.handleAuthPlain(string(line))
		}
	case "LOGIN":
		session.handleAuthLogin()
	default:
		session.writeResponse(504, "Authentication mechanism not supported")
	}
}

// handleAuthPlain 处理PLAIN认证
func (session *SMTPSession) handleAuthPlain(credentials string) {
	// PLAIN认证格式: base64(username\0username\0password)
	decoded, err := base64.StdEncoding.DecodeString(credentials)
	if err != nil {
		session.writeResponse(535, "Authentication failed")
		return
	}

	parts := strings.Split(string(decoded), "\x00")
	if len(parts) != 3 {
		session.writeResponse(535, "Authentication failed")
		return
	}

	username := parts[1]
	password := parts[2]

	// 验证用户名和密码
	if session.server.authenticateUser(username, password) {
		session.authenticated = true
		session.username = username
		session.writeResponse(235, "Authentication successful")
	} else {
		session.writeResponse(535, "Authentication failed")
	}
}

// handleAuthLogin 处理LOGIN认证
func (session *SMTPSession) handleAuthLogin() {
	// 请求用户名
	session.writeResponse(334, base64.StdEncoding.EncodeToString([]byte("Username:")))

	line, _, err := session.reader.ReadLine()
	if err != nil {
		session.writeResponse(535, "Authentication failed")
		return
	}

	usernameBytes, err := base64.StdEncoding.DecodeString(string(line))
	if err != nil {
		session.writeResponse(535, "Authentication failed")
		return
	}
	username := string(usernameBytes)

	// 请求密码
	session.writeResponse(334, base64.StdEncoding.EncodeToString([]byte("Password:")))

	line, _, err = session.reader.ReadLine()
	if err != nil {
		session.writeResponse(535, "Authentication failed")
		return
	}

	passwordBytes, err := base64.StdEncoding.DecodeString(string(line))
	if err != nil {
		session.writeResponse(535, "Authentication failed")
		return
	}
	password := string(passwordBytes)

	// 验证用户名和密码
	if session.server.authenticateUser(username, password) {
		session.authenticated = true
		session.username = username
		session.writeResponse(235, "Authentication successful")
	} else {
		session.writeResponse(535, "Authentication failed")
	}
}

// handleMail 处理MAIL FROM命令
func (session *SMTPSession) handleMail(args string) {
	if !strings.HasPrefix(strings.ToUpper(args), "FROM:") {
		session.writeResponse(501, "Syntax error")
		return
	}

	from := strings.TrimSpace(args[5:])
	from = strings.Trim(from, "<>")
	session.from = from
	session.writeResponse(250, "OK")
}

// handleRcpt 处理RCPT TO命令
func (session *SMTPSession) handleRcpt(args string) {
	if !strings.HasPrefix(strings.ToUpper(args), "TO:") {
		session.writeResponse(501, "Syntax error")
		return
	}

	to := strings.TrimSpace(args[3:])
	to = strings.Trim(to, "<>")

	// 检查收件人是否为本域用户
	if !session.isLocalUser(to) {
		session.writeResponse(550, "User not found")
		return
	}

	session.to = append(session.to, to)
	session.writeResponse(250, "OK")
}

// handleData 处理DATA命令
func (session *SMTPSession) handleData() {
	if session.from == "" || len(session.to) == 0 {
		session.writeResponse(503, "Bad sequence of commands")
		return
	}

	session.writeResponse(354, "Start mail input; end with <CRLF>.<CRLF>")

	var data []byte
	for {
		line, err := session.reader.ReadString('\n')
		if err != nil {
			log.Printf("读取邮件数据失败: %v", err)
			return
		}

		if line == ".\r\n" || line == ".\n" {
			break
		}

		data = append(data, []byte(line)...)
	}

	session.data = data

	// 保存邮件到数据库
	if err := session.saveEmail(); err != nil {
		log.Printf("保存邮件失败: %v", err)
		session.writeResponse(550, "Failed to save email")
		return
	}

	session.writeResponse(250, "OK")
	session.reset()
}

// reset 重置会话状态
func (session *SMTPSession) reset() {
	session.from = ""
	session.to = nil
	session.data = nil
}

// isLocalUser 检查是否为本地用户
func (session *SMTPSession) isLocalUser(email string) bool {
	var count int
	err := session.server.db.QueryRow("SELECT COUNT(*) FROM mailboxes WHERE email = ? AND is_active = 1", email).Scan(&count)
	if err != nil {
		log.Printf("查询邮箱失败: %v", err)
		return false
	}
	return count > 0
}

// saveEmail 保存邮件到数据库
func (session *SMTPSession) saveEmail() error {
	// 解析邮件内容
	subject, body := session.parseEmailContent()

	// 调试日志
	log.Printf("解析后的邮件 - Subject: %s, Body: %s", subject, body)

	// 为每个收件人保存邮件
	for _, to := range session.to {
		// 获取用户ID
		var userID int
		err := session.server.db.QueryRow("SELECT user_id FROM mailboxes WHERE email = ? AND is_active = 1", to).Scan(&userID)
		if err != nil {
			log.Printf("获取用户ID失败: %v", err)
			continue
		}

		// 生成消息ID
		messageID := fmt.Sprintf("<%d.%d@%s>", time.Now().Unix(), userID, session.server.config.Domain)

		// 插入邮件记录
		log.Printf("准备插入数据库 - Body: %s", body)
		_, err = session.server.db.Exec(`
			INSERT INTO emails (message_id, from_addr, to_addr, subject, body, user_id, size, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, messageID, session.from, to, subject, body, userID, len(body), time.Now(), time.Now())

		if err != nil {
			log.Printf("插入邮件记录失败: %v", err)
			return err
		}
	}

	return nil
}

// parseEmailContent 解析邮件内容
func (session *SMTPSession) parseEmailContent() (subject, body string) {
	content := string(session.data)
	lines := strings.Split(content, "\n")

	var inHeaders = true
	var bodyLines []string
	var contentTransferEncoding string
	var contentType string
	var boundary string

	for _, line := range lines {
		line = strings.TrimRight(line, "\r")

		if inHeaders {
			if line == "" {
				inHeaders = false
				continue
			}

			if strings.HasPrefix(strings.ToLower(line), "subject:") {
				subject = strings.TrimSpace(line[8:])
				// 解码Subject中的编码内容
				subject = decodeEmailHeader(subject)
			} else if strings.HasPrefix(strings.ToLower(line), "content-transfer-encoding:") {
				contentTransferEncoding = strings.TrimSpace(strings.ToLower(line[26:]))
			} else if strings.HasPrefix(strings.ToLower(line), "content-type:") {
				contentType = strings.TrimSpace(strings.ToLower(line[13:]))
				// 提取boundary（保持原始大小写）
				if strings.Contains(strings.ToLower(line), "boundary=") {
					originalLine := strings.TrimSpace(line[13:]) // 保持原始大小写
					parts := strings.Split(originalLine, "boundary=")
					if len(parts) > 1 {
						boundary = strings.Trim(strings.Split(parts[1], ";")[0], "\"")
					}
				}
			}
		} else {
			bodyLines = append(bodyLines, line)
		}
	}

	body = strings.Join(bodyLines, "\n")

	// 处理MIME多部分邮件
	log.Printf("Content-Type: %s, Boundary: %s", contentType, boundary)
	if strings.Contains(contentType, "multipart") && boundary != "" {
		log.Printf("处理MIME多部分邮件")
		body = parseMIMEMultipart(body, boundary)
	} else {
		log.Printf("处理单部分邮件，编码: %s", contentTransferEncoding)
		// 提取字符集
		charset := extractCharsetFromContentType(contentType)

		// 根据编码方式解码body
		if contentTransferEncoding == "base64" {
			cleanContent := strings.ReplaceAll(body, "\n", "")
			cleanContent = strings.ReplaceAll(cleanContent, "\r", "")
			cleanContent = strings.TrimSpace(cleanContent)
			if decoded, err := base64.StdEncoding.DecodeString(cleanContent); err == nil {
				// 进行字符编码转换
				body = convertToUTF8(decoded, charset)
			}
		} else if contentTransferEncoding == "quoted-printable" {
			// 处理quoted-printable编码
			body = decodeQuotedPrintable(body)
			body = convertToUTF8([]byte(body), charset)
		} else if charset != "" && charset != "utf-8" {
			// 如果没有传输编码但有字符集，直接进行字符编码转换
			body = convertToUTF8([]byte(body), charset)
		}
	}

	// 确保body是有效的UTF-8
	if !utf8.ValidString(body) {
		body = strings.ToValidUTF8(body, "?")
	}

	if subject == "" {
		subject = "无主题"
	}

	return subject, body
}

// parseMIMEMultipart 解析MIME多部分邮件
func parseMIMEMultipart(body, boundary string) string {
	parts := strings.Split(body, "--"+boundary)
	var textParts []string
	var htmlParts []string

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" || part == "--" || strings.HasPrefix(part, "--") {
			continue
		}

		// 分离头部和内容
		lines := strings.Split(part, "\n")
		var inHeaders = true
		var contentLines []string
		var contentType string
		var contentTransferEncoding string
		var charset string

		for _, line := range lines {
			line = strings.TrimRight(line, "\r")

			if inHeaders {
				if line == "" {
					inHeaders = false
					continue
				}

				if strings.HasPrefix(strings.ToLower(line), "content-type:") {
					contentType = strings.TrimSpace(line[13:]) // 保持原始大小写用于提取charset
					charset = extractCharsetFromContentType(contentType)
					contentType = strings.ToLower(contentType) // 转为小写用于比较
				} else if strings.HasPrefix(strings.ToLower(line), "content-transfer-encoding:") {
					contentTransferEncoding = strings.TrimSpace(strings.ToLower(line[26:]))
				}
			} else {
				// 跳过以--开头的行（boundary分隔符）
				if !strings.HasPrefix(line, "--") {
					contentLines = append(contentLines, line)
				}
			}
		}

		content := strings.Join(contentLines, "\n")
		content = strings.TrimSpace(content)

		// 根据编码方式解码内容
		if contentTransferEncoding == "base64" {
			cleanContent := strings.ReplaceAll(content, "\n", "")
			cleanContent = strings.ReplaceAll(cleanContent, "\r", "")
			cleanContent = strings.TrimSpace(cleanContent)
			log.Printf("尝试解码base64 (charset: %s): %s", charset, cleanContent)
			if decoded, err := base64.StdEncoding.DecodeString(cleanContent); err == nil {
				// 先进行base64解码，然后进行字符编码转换
				content = convertToUTF8(decoded, charset)
				log.Printf("base64解码并转换编码成功: %s", content)
			} else {
				log.Printf("base64解码失败: %v", err)
			}
		} else if contentTransferEncoding == "quoted-printable" {
			// 处理quoted-printable编码
			content = decodeQuotedPrintable(content)
			content = convertToUTF8([]byte(content), charset)
		} else if charset != "" && charset != "utf-8" {
			// 如果没有传输编码但有字符集，直接进行字符编码转换
			content = convertToUTF8([]byte(content), charset)
		}

		// 根据内容类型分类
		if strings.Contains(contentType, "text/plain") {
			textParts = append(textParts, content)
		} else if strings.Contains(contentType, "text/html") {
			htmlParts = append(htmlParts, content)
		}
	}

	// 优先返回纯文本内容，如果没有则返回HTML内容
	if len(textParts) > 0 {
		return strings.Join(textParts, "\n\n")
	} else if len(htmlParts) > 0 {
		return strings.Join(htmlParts, "\n\n")
	}

	return body // 如果解析失败，返回原始内容
}

// decodeEmailHeader 解码邮件头部编码内容
func decodeEmailHeader(header string) string {
	// 使用我们自己的MIME头部解码函数，支持更多编码
	return decodeMIMEHeaderSMTP(header)
}

// decodeMIMEHeaderSMTP 解码MIME编码的邮件头部 (=?charset?encoding?data?=) - SMTP版本
func decodeMIMEHeaderSMTP(header string) string {
	// MIME编码格式: =?charset?encoding?encoded-text?=
	re := regexp.MustCompile(`=\?([^?]+)\?([BbQq])\?([^?]*)\?=`)

	result := header
	matches := re.FindAllStringSubmatch(header, -1)

	for _, match := range matches {
		if len(match) != 4 {
			continue
		}

		fullMatch := match[0]
		charset := strings.ToLower(match[1])
		encoding := strings.ToUpper(match[2])
		encodedText := match[3]

		var decoded string

		switch encoding {
		case "B": // Base64编码
			if decodedBytes, err := base64.StdEncoding.DecodeString(encodedText); err == nil {
				decoded = convertToUTF8(decodedBytes, charset)
			} else {
				decoded = encodedText // 解码失败，保持原样
			}
		case "Q": // Quoted-printable编码
			decoded = decodeQuotedPrintableHeaderSMTP(encodedText)
			decoded = convertToUTF8([]byte(decoded), charset)
		default:
			decoded = encodedText // 未知编码，保持原样
		}

		result = strings.Replace(result, fullMatch, decoded, 1)
	}

	return result
}

// decodeQuotedPrintableHeaderSMTP 解码quoted-printable编码的头部 - SMTP版本
func decodeQuotedPrintableHeaderSMTP(s string) string {
	// 在头部中，下划线代表空格
	s = strings.ReplaceAll(s, "_", " ")

	result := strings.Builder{}
	for i := 0; i < len(s); i++ {
		if s[i] == '=' && i+2 < len(s) {
			// 尝试解析十六进制
			hex := s[i+1 : i+3]
			if b, err := strconv.ParseUint(hex, 16, 8); err == nil {
				result.WriteByte(byte(b))
				i += 2 // 跳过已处理的字符
			} else {
				result.WriteByte(s[i])
			}
		} else {
			result.WriteByte(s[i])
		}
	}

	return result.String()
}

// getEncodingByCharset 根据字符集名称获取编码器
func getEncodingByCharset(charset string) encoding.Encoding {
	charset = strings.ToLower(strings.TrimSpace(charset))

	switch charset {
	case "gbk", "gb2312", "gb18030":
		return simplifiedchinese.GBK
	case "big5":
		return traditionalchinese.Big5
	case "shift_jis", "shift-jis", "sjis":
		return japanese.ShiftJIS
	case "euc-jp":
		return japanese.EUCJP
	case "iso-2022-jp":
		return japanese.ISO2022JP
	case "euc-kr":
		return korean.EUCKR
	case "iso-8859-1", "latin1":
		return charmap.ISO8859_1
	case "iso-8859-2", "latin2":
		return charmap.ISO8859_2
	case "iso-8859-15":
		return charmap.ISO8859_15
	case "windows-1252", "cp1252":
		return charmap.Windows1252
	case "windows-1251", "cp1251":
		return charmap.Windows1251
	case "utf-8", "utf8":
		return nil // UTF-8不需要转换
	default:
		return nil // 未知编码，不转换
	}
}

// convertToUTF8 将指定编码的字节转换为UTF-8字符串
func convertToUTF8(data []byte, charset string) string {
	encoder := getEncodingByCharset(charset)
	if encoder == nil {
		// 如果是UTF-8或未知编码，直接返回
		return string(data)
	}

	// 创建解码器
	decoder := encoder.NewDecoder()

	// 转换为UTF-8
	utf8Data, err := io.ReadAll(transform.NewReader(bytes.NewReader(data), decoder))
	if err != nil {
		log.Printf("编码转换失败 (%s): %v", charset, err)
		// 转换失败时，尝试直接返回字符串
		return string(data)
	}

	return string(utf8Data)
}

// extractCharsetFromContentType 从Content-Type中提取字符集
func extractCharsetFromContentType(contentType string) string {
	// 查找charset参数
	re := regexp.MustCompile(`charset\s*=\s*["']?([^"'\s;]+)["']?`)
	matches := re.FindStringSubmatch(contentType)
	if len(matches) > 1 {
		return matches[1]
	}
	return "utf-8" // 默认UTF-8
}

// decodeQuotedPrintable 解码quoted-printable编码
func decodeQuotedPrintable(s string) string {
	// 简单的quoted-printable解码实现
	result := strings.Builder{}
	lines := strings.Split(s, "\n")

	for i, line := range lines {
		line = strings.TrimRight(line, "\r")

		// 处理软换行（以=结尾的行）
		if strings.HasSuffix(line, "=") {
			line = line[:len(line)-1] // 移除末尾的=
		} else if i < len(lines)-1 {
			line += "\n" // 添加换行符（除了最后一行）
		}

		// 解码=XX格式的字符
		for j := 0; j < len(line); j++ {
			if line[j] == '=' && j+2 < len(line) {
				// 尝试解析十六进制
				hex := line[j+1 : j+3]
				if b, err := strconv.ParseUint(hex, 16, 8); err == nil {
					result.WriteByte(byte(b))
					j += 2 // 跳过已处理的字符
				} else {
					result.WriteByte(line[j])
				}
			} else {
				result.WriteByte(line[j])
			}
		}
	}

	return result.String()
}
