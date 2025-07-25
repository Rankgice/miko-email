package email

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"mime/quotedprintable"
	"net"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/encoding"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
	"miko-email/internal/services/forward"
	"miko-email/internal/services/smtp"
	"miko-email/internal/svc"

	"github.com/jhillyerd/enmime/v2"
	"gorm.io/gorm"
	"miko-email/internal/model"
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

type Service struct {
	svcCtx         *svc.ServiceContext
	tracker        *ConnectionTracker
	forwardService *forward.Service
	smtpClient     *smtp.OutboundClient
}

func NewService(svcCtx *svc.ServiceContext) *Service {
	db, err := svcCtx.DB.DB()
	if err != nil {
		return nil
	}
	return &Service{
		svcCtx:         svcCtx,
		tracker:        NewConnectionTracker(),
		forwardService: forward.NewService(svcCtx),
		smtpClient:     smtp.NewOutboundClientWithDB(db), // 暂时使用固定域名
	}
}

// StartSMTPServer 启动SMTP服务器
func (s *Service) StartSMTPServer(port string) error {
	log.Printf("SMTP server starting on port %s", port)

	// 根据端口决定是否使用SSL
	if port == "465" {
		return s.startSMTPSServer(port)
	} else if port == "587" {
		return s.startSMTPWithSTARTTLS(port)
	} else {
		return s.startPlainSMTPServer(port)
	}
}

// startPlainSMTPServer 启动普通SMTP服务器（25端口）
func (s *Service) startPlainSMTPServer(port string) error {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to start SMTP server: %w", err)
	}
	defer listener.Close()

	log.Printf("SMTP server listening on port %s", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("SMTP connection error: %v", err)
			continue
		}

		go s.handleSMTPConnection(conn, false)
	}
}

// startSMTPSServer 启动SMTPS服务器（465端口，SSL）
func (s *Service) startSMTPSServer(port string) error {
	// 创建自签名证书（生产环境应使用真实证书）
	cert, err := s.generateSelfSignedCert()
	if err != nil {
		log.Printf("警告：无法生成SSL证书，465端口将使用普通连接: %v", err)
		return s.startPlainSMTPServer(port)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   "localhost",
	}

	listener, err := tls.Listen("tcp", ":"+port, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to start SMTPS server: %w", err)
	}
	defer listener.Close()

	log.Printf("SMTPS server listening on port %s (SSL)", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("SMTPS connection error: %v", err)
			continue
		}

		go s.handleSMTPConnection(conn, true)
	}
}

// startSMTPWithSTARTTLS 启动支持STARTTLS的SMTP服务器（587端口）
func (s *Service) startSMTPWithSTARTTLS(port string) error {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to start SMTP server: %w", err)
	}
	defer listener.Close()

	log.Printf("SMTP server listening on port %s (STARTTLS)", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("SMTP connection error: %v", err)
			continue
		}

		go s.handleSMTPConnection(conn, false)
	}
}

// generateSelfSignedCert 生成自签名证书
func (s *Service) generateSelfSignedCert() (tls.Certificate, error) {
	// 生成私钥
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return tls.Certificate{}, err
	}

	// 创建证书模板
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:  []string{"Miko Email System"},
			Country:       []string{"CN"},
			Province:      []string{""},
			Locality:      []string{""},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(365 * 24 * time.Hour), // 1年有效期
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		DNSNames:    []string{"localhost", "jbjj.site", "*.jbjj.site"},
	}

	// 生成证书
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return tls.Certificate{}, err
	}

	// 创建TLS证书
	cert := tls.Certificate{
		Certificate: [][]byte{certDER},
		PrivateKey:  priv,
	}

	return cert, nil
}

// StartIMAPServer 启动IMAP服务器
func (s *Service) StartIMAPServer(port string) error {
	log.Printf("IMAP server starting on port %s", port)

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to start IMAP server: %w", err)
	}
	defer listener.Close()

	log.Printf("IMAP server listening on port %s", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("IMAP connection error: %v", err)
			continue
		}

		go s.handleIMAPConnection(conn)
	}
}

// StartPOP3Server 启动POP3服务器
func (s *Service) StartPOP3Server(port string) error {
	log.Printf("POP3 server starting on port %s", port)

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to start POP3 server: %w", err)
	}
	defer listener.Close()

	log.Printf("POP3 server listening on port %s", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("POP3 connection error: %v", err)
			continue
		}

		go s.handlePOP3Connection(conn)
	}
}

// handleSMTPConnection 处理SMTP连接
func (s *Service) handleSMTPConnection(conn net.Conn, isSSL bool) {
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
	s.writeResponse(writer, 220, "jbjj.site Miko Email SMTP Server Ready")

	session := &SMTPSession{
		conn:   conn,
		reader: reader,
		writer: writer,
		server: s,
		isSSL:  isSSL,
	}

	session.handle()
}

// writeResponse 写入SMTP响应
func (s *Service) writeResponse(writer *bufio.Writer, code int, message string) error {
	response := fmt.Sprintf("%d %s\r\n", code, message)
	_, err := writer.WriteString(response)
	if err != nil {
		return err
	}
	return writer.Flush()
}

// authenticateUser 验证用户
func (s *Service) authenticateUser(username, password string) bool {
	// 首先尝试邮箱认证（mailboxes表）
	mailbox, err := s.svcCtx.MailboxModel.GetByEmailAndPassword(username, password)
	if err == nil && mailbox != nil {
		// 邮箱认证成功
		return true
	}

	// 如果邮箱认证失败，尝试用户认证（users表）
	user, err := s.svcCtx.UserModel.GetByEmail(username)
	if err != nil {
		log.Printf("用户认证失败: %v", err)
		return false
	}

	// 用户认证：使用bcrypt验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

// SMTPSession SMTP会话
type SMTPSession struct {
	conn          net.Conn
	reader        *bufio.Reader
	writer        *bufio.Writer
	server        *Service
	helo          string
	from          string
	to            []string
	data          []byte
	username      string
	password      string
	authenticated bool
	tlsEnabled    bool
	isSSL         bool
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
			session.handleHelo(args, command)
		case "STARTTLS":
			session.handleStartTLS()
		case "AUTH":
			session.handleAuth(args)
		case "MAIL":
			session.handleMail(args)
		case "RCPT":
			session.handleRcpt(args)
		case "DATA":
			session.handleData()
		case "QUIT":
			log.Printf("客户端请求断开连接: %s", session.conn.RemoteAddr())
			session.writeResponse(221, "Bye")
			return
		case "RSET":
			log.Printf("客户端重置会话: %s", session.conn.RemoteAddr())
			session.reset()
			session.writeResponse(250, "OK")
		case "NOOP":
			session.writeResponse(250, "OK")
		case ".":
			// 处理单独的点号命令（可能是DATA结束标记的误用）
			log.Printf("收到单独的点号命令，可能是协议错误: %s", session.conn.RemoteAddr())
			log.Printf("当前会话状态 - From: %s, To: %v", session.from, session.to)
			// 更宽容的处理，返回OK而不是错误，避免断开连接
			session.writeResponse(250, "OK")
		default:
			log.Printf("未识别的SMTP命令: %s (来自 %s)", command, session.conn.RemoteAddr())
			session.writeResponse(500, "Command not recognized")
		}
	}
}

// writeResponse 写入响应
func (session *SMTPSession) writeResponse(code int, message string) error {
	return session.server.writeResponse(session.writer, code, message)
}

// handleHelo 处理HELO/EHLO命令
func (session *SMTPSession) handleHelo(args string, command string) {
	session.helo = args
	log.Printf("SMTP握手: %s (来自 %s)", args, session.conn.RemoteAddr())

	if command == "EHLO" {
		// EHLO响应，支持扩展
		session.writer.WriteString("250-Hello " + args + "\r\n")

		// 如果不是SSL连接且不是已经启用TLS，则广告STARTTLS
		if !session.isSSL && !session.tlsEnabled {
			session.writer.WriteString("250-STARTTLS\r\n")
		}

		session.writer.WriteString("250-AUTH PLAIN LOGIN\r\n")
		session.writer.WriteString("250 8BITMIME\r\n")
	} else {
		// HELO响应，简单模式
		session.writer.WriteString("250 Hello " + args + "\r\n")
	}

	session.writer.Flush()
}

// handleStartTLS 处理STARTTLS命令
func (session *SMTPSession) handleStartTLS() {
	if session.isSSL || session.tlsEnabled {
		session.writeResponse(503, "TLS already active")
		return
	}

	// 生成自签名证书
	cert, err := session.server.generateSelfSignedCert()
	if err != nil {
		log.Printf("生成TLS证书失败: %v", err)
		session.writeResponse(454, "TLS not available due to temporary reason")
		return
	}

	// 发送准备开始TLS的响应
	session.writeResponse(220, "Ready to start TLS")

	// 创建TLS配置
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   "localhost",
	}

	// 将连接升级为TLS
	tlsConn := tls.Server(session.conn, tlsConfig)
	err = tlsConn.Handshake()
	if err != nil {
		log.Printf("TLS握手失败: %v", err)
		return
	}

	log.Printf("TLS连接建立成功: %s", session.conn.RemoteAddr())

	// 更新会话连接和读写器
	session.conn = tlsConn
	session.reader = bufio.NewReader(tlsConn)
	session.writer = bufio.NewWriter(tlsConn)
	session.tlsEnabled = true

	// 重置会话状态（TLS后需要重新认证）
	session.authenticated = false
	session.helo = ""
}

// IMAPSession IMAP会话
type IMAPSession struct {
	conn    net.Conn
	reader  *bufio.Reader
	writer  *bufio.Writer
	server  *Service
	state   string // NOTAUTHENTICATED, AUTHENTICATED, SELECTED
	user    string
	mailbox string
	tag     string
}

// POP3Session POP3会话
type POP3Session struct {
	conn      net.Conn
	reader    *bufio.Reader
	writer    *bufio.Writer
	server    *Service
	state     string // AUTHORIZATION, TRANSACTION, UPDATE
	user      string
	mailboxID int
	emails    []POP3Email
	deleted   map[int]bool
}

// POP3Email POP3邮件信息
type POP3Email struct {
	ID      int
	Size    int
	Subject string
	From    string
	To      string
	Date    string
	Body    string
}

// handleIMAPConnection 处理IMAP连接
func (s *Service) handleIMAPConnection(conn net.Conn) {
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
	session.writeResponse("* OK Miko Email IMAP Server Ready")

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
	session.writeResponse(session.tag + " " + response)
}

// handleCapability 处理CAPABILITY命令
func (session *IMAPSession) handleCapability() {
	session.writeResponse("* CAPABILITY IMAP4rev1 AUTH=PLAIN AUTH=LOGIN")
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

	// 支持多种认证方式
	if session.authenticateIMAPUser(username, password) {
		session.state = "AUTHENTICATED"
		session.user = username
		log.Printf("IMAP登录成功: %s", username)
		session.writeTaggedResponse("OK LOGIN completed")
	} else {
		log.Printf("IMAP登录失败: %s", username)
		session.writeTaggedResponse("NO LOGIN failed")
	}
}

// authenticateIMAPUser IMAP用户认证（支持多种认证方式）
func (session *IMAPSession) authenticateIMAPUser(username, password string) bool {
	log.Printf("IMAP认证开始: 用户名=%s", username)

	// 方式1: 直接邮箱认证 (邮箱地址 + 邮箱密码)
	log.Printf("尝试方式1: 直接邮箱认证")
	if session.authenticateByMailbox(username, password) {
		log.Printf("方式1认证成功")
		return true
	}

	// 方式2: 组合认证 (网站账号@邮箱地址 + 邮箱密码)
	// 格式: "网站用户名@邮箱地址" + "邮箱密码"
	log.Printf("尝试方式2: 组合认证")
	if strings.Contains(username, "@") {
		parts := strings.Split(username, "@")
		log.Printf("用户名分割结果: %v", parts)
		if len(parts) >= 2 {
			// 重新组合邮箱地址
			emailParts := parts[1:]
			emailAddr := strings.Join(emailParts, "@")
			websiteUser := parts[0]

			log.Printf("解析结果: 网站用户=%s, 邮箱地址=%s", websiteUser, emailAddr)

			// 验证邮箱和密码
			if session.authenticateByMailbox(emailAddr, password) {
				log.Printf("邮箱密码验证成功")
				// 同时验证网站用户是否有权限访问该邮箱
				if session.verifyUserMailboxAccess(websiteUser, emailAddr) {
					log.Printf("组合认证成功: 网站用户=%s, 邮箱=%s", websiteUser, emailAddr)
					return true
				} else {
					log.Printf("用户权限验证失败")
				}
			} else {
				log.Printf("邮箱密码验证失败")
			}
		}
	}

	// 方式3: 网站用户认证 (网站用户名 + 网站密码)
	log.Printf("尝试方式3: 网站用户认证")
	if session.authenticateByWebsiteUser(username, password) {
		log.Printf("方式3认证成功")
		return true
	}

	log.Printf("所有认证方式都失败")
	return false
}

// authenticateByMailbox 通过邮箱认证
func (session *IMAPSession) authenticateByMailbox(email, password string) bool {
	mailbox, err := session.server.svcCtx.MailboxModel.GetByEmailAndPassword(email, password)
	if err != nil {
		return false
	}

	return mailbox != nil
}

// verifyUserMailboxAccess 验证网站用户是否有权限访问指定邮箱
func (session *IMAPSession) verifyUserMailboxAccess(websiteUser, email string) bool {
	// 先获取用户
	user, err := session.server.svcCtx.UserModel.GetByUsername(websiteUser)
	if err != nil {
		log.Printf("获取用户失败: %v", err)
		return false
	}

	// 检查用户是否拥有该邮箱
	mailbox, err := session.server.svcCtx.MailboxModel.GetByEmailAndUserId(email, user.Id)
	if err != nil {
		log.Printf("验证用户邮箱权限失败: %v", err)
		return false
	} else if !mailbox.IsActive {
		log.Printf("验证用户邮箱权限失败: %v", err)
		return false
	}

	log.Printf("用户权限验证: 网站用户=%s, 邮箱=%s, 验证成功", websiteUser, email)
	return true
}

// authenticateByWebsiteUser 通过网站用户认证
func (session *IMAPSession) authenticateByWebsiteUser(username, password string) bool {
	user, err := session.server.svcCtx.UserModel.GetByEmail(username)
	if err != nil {
		return false
	}

	// 使用bcrypt验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
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

	// 获取邮件数量 - 查询收件箱邮件
	count, err := session.server.svcCtx.EmailModel.CountEmailsByUserEmail(session.user, "inbox")
	if err != nil {
		log.Printf("查询邮件数量失败: %v", err)
		count = 0
	}

	session.state = "SELECTED"
	session.mailbox = "INBOX"

	session.writeResponse(fmt.Sprintf("* %d EXISTS", count))
	session.writeResponse("* 0 RECENT")
	session.writeResponse("* OK [UIDVALIDITY 1] UIDs valid")
	session.writeTaggedResponse("OK [READ-WRITE] SELECT completed")
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

	// 查询用户的所有邮件ID
	emails, err := session.server.svcCtx.EmailModel.GetEmailsByUserEmail(session.user, "inbox", 0)
	if err != nil {
		log.Printf("SEARCH查询失败: %v", err)
		session.writeTaggedResponse("NO SEARCH failed")
		return
	}

	var emailIDs []string
	for _, email := range emails {
		emailIDs = append(emailIDs, fmt.Sprintf("%d", email.Id))
	}

	// 返回搜索结果
	if len(emailIDs) > 0 {
		session.writeResponse("* SEARCH " + strings.Join(emailIDs, " "))
	} else {
		session.writeResponse("* SEARCH")
	}
	session.writeTaggedResponse("OK SEARCH completed")
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
	dataItems := strings.Join(args[1:], " ")

	// 查询邮件数据
	emails, err := session.server.svcCtx.EmailModel.GetEmailsByUserEmail(session.user, "inbox", 10)
	if err != nil {
		session.writeTaggedResponse("NO FETCH failed")
		return
	}

	seqNum := 1
	for _, email := range emails {

		// 根据请求的数据项返回不同的信息
		if strings.Contains(dataItems, "RFC822") {
			// 构造完整的邮件内容
			emailContent := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nDate: %s\r\n\r\n%s",
				email.FromAddr, email.ToAddr, email.Subject, email.CreatedAt.Format("2006-01-02 15:04:05"), email.Body)

			// 返回RFC822格式的邮件
			session.writeResponse(fmt.Sprintf("* %d FETCH (RFC822 {%d}", seqNum, len(emailContent)))
			session.writer.WriteString(emailContent)
			session.writer.WriteString(")\r\n")
			session.writer.Flush()
		} else {
			// 返回基本信息
			session.writeResponse(fmt.Sprintf("* %d FETCH (UID %d RFC822.SIZE %d ENVELOPE (\"%s\" \"%s\" ((\"%s\" NIL \"%s\" NIL)) NIL NIL NIL NIL NIL))",
				seqNum, email.Id, len(email.Body), email.CreatedAt.Format("2006-01-02 15:04:05"), email.Subject, email.FromAddr, email.FromAddr))
		}
		seqNum++
	}

	session.writeTaggedResponse("OK FETCH completed")
}

// handleLogout 处理LOGOUT命令
func (session *IMAPSession) handleLogout() {
	session.writeResponse("* BYE Miko Email IMAP Server logging out")
	session.writeTaggedResponse("OK LOGOUT completed")
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

// handlePOP3Connection 处理POP3连接
func (s *Service) handlePOP3Connection(conn net.Conn) {
	defer conn.Close()

	log.Printf("新的POP3连接: %s", conn.RemoteAddr())

	// 设置连接超时
	conn.SetDeadline(time.Now().Add(30 * time.Minute))

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	session := &POP3Session{
		conn:    conn,
		reader:  reader,
		writer:  writer,
		server:  s,
		state:   "AUTHORIZATION",
		deleted: make(map[int]bool),
	}

	// 发送欢迎消息
	session.writeResponse("+OK Miko Email POP3 Server Ready")

	session.handle()
}

// handle 处理POP3会话
func (session *POP3Session) handle() {
	for {
		line, _, err := session.reader.ReadLine()
		if err != nil {
			log.Printf("读取POP3命令失败: %v", err)
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

		switch session.state {
		case "AUTHORIZATION":
			switch cmd {
			case "USER":
				session.handleUser(args)
			case "PASS":
				session.handlePass(args)
			case "QUIT":
				session.handleQuit()
				return
			default:
				session.writeResponse("-ERR Command not available in AUTHORIZATION state")
			}
		case "TRANSACTION":
			switch cmd {
			case "STAT":
				session.handleStat()
			case "LIST":
				session.handleList(args)
			case "RETR":
				session.handleRetr(args)
			case "DELE":
				session.handleDele(args)
			case "NOOP":
				session.handleNoop()
			case "RSET":
				session.handleRset()
			case "TOP":
				session.handleTop(args)
			case "UIDL":
				session.handleUidl(args)
			case "QUIT":
				session.handleQuit()
				return
			default:
				session.writeResponse("-ERR Command not recognized")
			}
		}
	}
}

// writeResponse 写入POP3响应
func (session *POP3Session) writeResponse(response string) {
	session.writer.WriteString(response + "\r\n")
	session.writer.Flush()
}

// handleUser 处理USER命令
func (session *POP3Session) handleUser(args []string) {
	if len(args) != 1 {
		session.writeResponse("-ERR USER requires username")
		return
	}

	session.user = args[0]
	session.writeResponse("+OK User accepted")
}

// handlePass 处理PASS命令
func (session *POP3Session) handlePass(args []string) {
	if len(args) != 1 {
		session.writeResponse("-ERR PASS requires password")
		return
	}

	if session.user == "" {
		session.writeResponse("-ERR USER required first")
		return
	}

	password := args[0]

	// 验证用户凭据 - 从mailboxes表查询
	mailbox, err := session.server.svcCtx.MailboxModel.GetByEmailAndPassword(session.user, password)
	if err != nil {
		log.Printf("POP3登录失败 - 邮箱不存在或密码错误: %s, 错误: %v", session.user, err)
		session.writeResponse("-ERR Authentication failed")
		return
	}

	session.mailboxID = int(mailbox.Id)
	session.state = "TRANSACTION"

	// 加载邮件列表
	err = session.loadEmails()
	if err != nil {
		log.Printf("POP3加载邮件失败: %v", err)
		session.writeResponse("-ERR Failed to load mailbox")
		return
	}

	log.Printf("POP3登录成功: %s", session.user)
	session.writeResponse(fmt.Sprintf("+OK Mailbox ready, %d messages", len(session.emails)))
}

// loadEmails 加载邮件列表
func (session *POP3Session) loadEmails() error {
	emails, err := session.server.svcCtx.EmailModel.GetEmailsForPOP3(int64(session.mailboxID))
	if err != nil {
		return err
	}

	session.emails = []POP3Email{}
	for _, email := range emails {
		pop3Email := POP3Email{
			ID:      int(email.Id),
			From:    email.FromAddr,
			To:      email.ToAddr,
			Subject: email.Subject,
			Body:    email.Body,
			Date:    email.CreatedAt.Format("2006-01-02 15:04:05"),
			Size:    len(email.Body) + len(email.Subject) + len(email.FromAddr) + len(email.ToAddr) + 100, // 估算大小
		}
		session.emails = append(session.emails, pop3Email)
	}

	return nil
}

// handleStat 处理STAT命令
func (session *POP3Session) handleStat() {
	totalSize := 0
	count := 0
	for i, email := range session.emails {
		if !session.deleted[i+1] {
			totalSize += email.Size
			count++
		}
	}
	session.writeResponse(fmt.Sprintf("+OK %d %d", count, totalSize))
}

// handleList 处理LIST命令
func (session *POP3Session) handleList(args []string) {
	if len(args) == 0 {
		// 列出所有邮件
		session.writeResponse("+OK")

		// 发送邮件列表
		for i, email := range session.emails {
			if !session.deleted[i+1] {
				line := fmt.Sprintf("%d %d", i+1, email.Size)
				session.writer.WriteString(line + "\r\n")
			}
		}

		// 发送结束标记
		session.writer.WriteString(".\r\n")
		session.writer.Flush()
	} else {
		// 列出指定邮件
		msgNum, err := strconv.Atoi(args[0])
		if err != nil || msgNum < 1 || msgNum > len(session.emails) {
			session.writeResponse("-ERR Invalid message number")
			return
		}

		if session.deleted[msgNum] {
			session.writeResponse("-ERR Message deleted")
			return
		}

		email := session.emails[msgNum-1]
		session.writeResponse(fmt.Sprintf("+OK %d %d", msgNum, email.Size))
	}
}

// handleRetr 处理RETR命令
func (session *POP3Session) handleRetr(args []string) {
	if len(args) != 1 {
		session.writeResponse("-ERR RETR requires message number")
		return
	}

	msgNum, err := strconv.Atoi(args[0])
	if err != nil || msgNum < 1 || msgNum > len(session.emails) {
		session.writeResponse("-ERR Invalid message number")
		return
	}

	if session.deleted[msgNum] {
		session.writeResponse("-ERR Message deleted")
		return
	}

	email := session.emails[msgNum-1]

	// 构造完整的邮件内容
	emailContent := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nDate: %s\r\n\r\n%s",
		email.From, email.To, email.Subject, email.Date, email.Body)

	session.writeResponse(fmt.Sprintf("+OK %d octets", len(emailContent)))

	// 发送邮件内容
	session.writer.WriteString(emailContent)

	// 确保邮件内容以换行结束，然后发送结束标记
	if !strings.HasSuffix(emailContent, "\r\n") {
		session.writer.WriteString("\r\n")
	}
	session.writer.WriteString(".\r\n")
	session.writer.Flush()
}

// handleDele 处理DELE命令
func (session *POP3Session) handleDele(args []string) {
	if len(args) != 1 {
		session.writeResponse("-ERR DELE requires message number")
		return
	}

	msgNum, err := strconv.Atoi(args[0])
	if err != nil || msgNum < 1 || msgNum > len(session.emails) {
		session.writeResponse("-ERR Invalid message number")
		return
	}

	if session.deleted[msgNum] {
		session.writeResponse("-ERR Message already deleted")
		return
	}

	session.deleted[msgNum] = true
	session.writeResponse("+OK Message deleted")
}

// handleNoop 处理NOOP命令
func (session *POP3Session) handleNoop() {
	session.writeResponse("+OK")
}

// handleRset 处理RSET命令
func (session *POP3Session) handleRset() {
	session.deleted = make(map[int]bool)
	session.writeResponse("+OK")
}

// handleTop 处理TOP命令
func (session *POP3Session) handleTop(args []string) {
	if len(args) != 2 {
		session.writeResponse("-ERR TOP requires message number and line count")
		return
	}

	msgNum, err := strconv.Atoi(args[0])
	if err != nil || msgNum < 1 || msgNum > len(session.emails) {
		session.writeResponse("-ERR Invalid message number")
		return
	}

	lines, err := strconv.Atoi(args[1])
	if err != nil || lines < 0 {
		session.writeResponse("-ERR Invalid line count")
		return
	}

	if session.deleted[msgNum] {
		session.writeResponse("-ERR Message deleted")
		return
	}

	email := session.emails[msgNum-1]

	// 构造邮件头
	header := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nDate: %s\r\n\r\n",
		email.From, email.To, email.Subject, email.Date)

	// 获取指定行数的正文
	bodyLines := strings.Split(email.Body, "\n")
	if lines > len(bodyLines) {
		lines = len(bodyLines)
	}

	body := strings.Join(bodyLines[:lines], "\n")
	content := header + body

	session.writeResponse("+OK")
	session.writer.WriteString(content)

	// 确保内容以换行结束，然后发送结束标记
	if !strings.HasSuffix(content, "\r\n") {
		session.writer.WriteString("\r\n")
	}
	session.writer.WriteString(".\r\n")
	session.writer.Flush()
}

// handleUidl 处理UIDL命令
func (session *POP3Session) handleUidl(args []string) {
	if len(args) == 0 {
		// 列出所有邮件的UIDL
		session.writeResponse("+OK")

		// 发送UIDL列表
		for i, email := range session.emails {
			if !session.deleted[i+1] {
				line := fmt.Sprintf("%d %d", i+1, email.ID)
				session.writer.WriteString(line + "\r\n")
			}
		}

		// 发送结束标记
		session.writer.WriteString(".\r\n")
		session.writer.Flush()
	} else {
		// 列出指定邮件的UIDL
		msgNum, err := strconv.Atoi(args[0])
		if err != nil || msgNum < 1 || msgNum > len(session.emails) {
			session.writeResponse("-ERR Invalid message number")
			return
		}

		if session.deleted[msgNum] {
			session.writeResponse("-ERR Message deleted")
			return
		}

		email := session.emails[msgNum-1]
		session.writeResponse(fmt.Sprintf("+OK %d %d", msgNum, email.ID))
	}
}

// handleQuit 处理QUIT命令
func (session *POP3Session) handleQuit() {
	if session.state == "TRANSACTION" {
		// 在UPDATE状态下删除标记为删除的邮件
		for msgNum := range session.deleted {
			if msgNum > 0 && msgNum <= len(session.emails) {
				email := session.emails[msgNum-1]
				err := session.server.svcCtx.EmailModel.DeleteEmailById(nil, int64(email.ID))
				if err != nil {
					log.Printf("删除邮件失败: %v", err)
				} else {
					log.Printf("已删除邮件 ID: %d", email.ID)
				}
			}
		}
	}

	session.writeResponse("+OK Miko Email POP3 Server signing off")
}

// handleMail 处理MAIL FROM命令
func (session *SMTPSession) handleMail(args string) {
	if !strings.HasPrefix(strings.ToUpper(args), "FROM:") {
		session.writeResponse(501, "Syntax error")
		return
	}

	// 提取FROM:后面的内容
	fromPart := strings.TrimSpace(args[5:])

	// 处理可能的参数，如 BODY=8BITMIME
	parts := strings.Fields(fromPart)
	if len(parts) > 0 {
		from := strings.Trim(parts[0], "<>")

		session.from = from
		log.Printf("设置发件人: %s (来自 %s)", from, session.conn.RemoteAddr())
	} else {
		log.Printf("MAIL FROM语法错误: %s (来自 %s)", args, session.conn.RemoteAddr())
		session.writeResponse(501, "Syntax error")
		return
	}

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
	if session.isLocalUser(to) {
		// 本域用户，允许接收
		session.to = append(session.to, to)
		log.Printf("添加本域收件人: %s (来自 %s)", to, session.conn.RemoteAddr())
		session.writeResponse(250, "OK")
		return
	}

	// 外部邮箱，需要认证才能发送
	if !session.isValidExternalEmail(to) {
		log.Printf("无效的收件人地址: %s (来自 %s)", to, session.conn.RemoteAddr())
		session.writeResponse(550, "Invalid recipient address")
		return
	}

	// 检查是否有权限发送到外部邮箱
	if !session.authenticated {
		log.Printf("未认证用户尝试发送到外部邮箱: %s (来自 %s)", to, session.conn.RemoteAddr())
		session.writeResponse(550, "Authentication required for external delivery")
		return
	}

	session.to = append(session.to, to)
	log.Printf("添加外部收件人: %s (来自 %s)", to, session.conn.RemoteAddr())
	session.writeResponse(250, "OK")
}

// handleData 处理DATA命令
func (session *SMTPSession) handleData() {
	if session.from == "" || len(session.to) == 0 {
		session.writeResponse(503, "Bad sequence of commands")
		return
	}

	// 检查发件人权限（现在有收件人信息了）
	if !session.canSendFrom(session.from) {
		log.Printf("发件人权限不足: %s (来自 %s)", session.from, session.conn.RemoteAddr())
		session.writeResponse(550, "Authentication required or sender not authorized")
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
	exists, err := session.server.svcCtx.MailboxModel.CheckEmailExists(email)
	if err != nil {
		log.Printf("查询邮箱失败: %v", err)
		return false
	}
	return exists
}

// canSendFrom 检查是否有权限从指定地址发送邮件
func (session *SMTPSession) canSendFrom(from string) bool {
	// 如果已经认证，检查是否有权限使用该发件人地址
	if session.authenticated {
		// 检查认证用户是否有权限使用该发件人地址
		return session.isAuthorizedSender(from)
	}

	// 如果未认证，检查是否是外部邮件投递到本地邮箱
	// 这是正常的邮件接收流程，应该被允许
	if session.hasLocalRecipients() {
		log.Printf("允许外部邮件投递: %s -> 本地邮箱 (来自 %s)", from, session.conn.RemoteAddr())
		return true
	}

	// 如果未认证，只允许本地连接发送本域邮件（用于内部系统）
	clientIP := strings.Split(session.conn.RemoteAddr().String(), ":")[0]
	isLocalConnection := clientIP == "127.0.0.1" || clientIP == "::1" || clientIP == "localhost"

	if isLocalConnection && session.isLocalUser(from) {
		log.Printf("允许本地连接发送本域邮件: %s", from)
		return true
	}

	log.Printf("未认证用户尝试发送邮件: %s (来自 %s)", from, session.conn.RemoteAddr())
	return false
}

// hasLocalRecipients 检查是否有本地收件人
func (session *SMTPSession) hasLocalRecipients() bool {
	for _, recipient := range session.to {
		if session.isLocalUser(recipient) {
			return true
		}
	}
	return false
}

// isAuthorizedSender 检查认证用户是否有权限使用指定发件人地址
func (session *SMTPSession) isAuthorizedSender(from string) bool {
	if session.username == "" {
		return false
	}

	// 如果认证用户就是发件人邮箱，直接允许
	if session.username == from {
		return true
	}

	// 检查用户是否拥有该邮箱（通过users表关联）
	user, err := session.server.svcCtx.UserModel.GetByEmail(session.username)
	if err != nil {
		log.Printf("获取用户失败: %v", err)
		return false
	}

	mailbox, err := session.server.svcCtx.MailboxModel.GetByEmailAndUserId(from, user.Id)
	if err != nil {
		log.Printf("检查发件人权限失败: %v", err)
		return false
	} else if !mailbox.IsActive {
		log.Printf("检查发件人权限失败: %v", err)
		return false
	}

	return true
}

// isValidExternalEmail 检查是否为有效的外部邮箱
func (session *SMTPSession) isValidExternalEmail(email string) bool {
	// 基本的邮箱格式验证
	if !strings.Contains(email, "@") {
		return false
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	domain := parts[1]
	if domain == "" {
		return false
	}

	// 检查是否为本地域名
	localDomains := []string{"localhost", "jbjj.site"}
	for _, localDomain := range localDomains {
		if domain == localDomain {
			return false // 本地域名应该通过 isLocalUser 检查
		}
	}

	// 简单的域名格式验证
	if len(domain) < 3 || !strings.Contains(domain, ".") {
		return false
	}

	return true
}

// sendToExternalEmail 发送邮件到外部邮箱
func (s *Service) sendToExternalEmail(from, to, subject, body string) error {
	// 使用SMTP客户端发送到外部邮箱
	if s.smtpClient != nil && s.smtpClient.IsExternalEmail(to) {
		err := s.smtpClient.SendEmail(from, to, subject, body)
		if err != nil {
			log.Printf("外部邮箱发送失败: %v", err)
			return fmt.Errorf("外部邮箱发送失败: %w", err)
		}
		log.Printf("✅ 外部邮箱发送成功: %s -> %s", from, to)
		return nil
	}
	return fmt.Errorf("无效的外部邮箱地址: %s", to)
}

// SaveEmail 保存邮件到数据库
func (s *Service) SaveEmail(mailboxID int64, fromAddr, toAddr, subject, body string) error {
	return s.svcCtx.EmailModel.SaveEmailToFolder(nil, mailboxID, fromAddr, toAddr, subject, body, "inbox")
}

// saveEmail 保存邮件到数据库
func (session *SMTPSession) saveEmail() error {
	// 解析邮件内容
	subject, body := session.parseEmailContent()

	// 调试日志
	log.Printf("解析后的邮件 - Subject: %s, Body: %s", subject, body)

	// 为每个收件人处理邮件
	for _, to := range session.to {
		if session.isLocalUser(to) {
			// 本地用户，保存到数据库
			mailboxID, err := session.server.svcCtx.MailboxModel.GetIdByEmail(to)
			if err != nil {
				log.Printf("获取邮箱ID失败: %v", err)
				continue
			}

			// 插入邮件记录
			log.Printf("准备插入数据库 - Body: %s", body)
			err = session.server.svcCtx.EmailModel.SaveEmailToFolder(nil, mailboxID, session.from, to, subject, body, "inbox")
			if err != nil {
				log.Printf("插入邮件记录失败: %v", err)
				return err
			}

			log.Printf("✅ 邮件保存成功 - 邮箱ID: %d, 主题: %s", mailboxID, subject)

			// 检查并执行转发规则
			session.server.processForwardRules(to, session.from, subject, body)
		} else {
			// 外部邮箱，发送到外部
			log.Printf("发送邮件到外部邮箱: %s", to)
			err := session.server.sendToExternalEmail(session.from, to, subject, body)
			if err != nil {
				log.Printf("外部邮件发送失败: %v", err)
				return err
			}
			log.Printf("✅ 外部邮件发送成功: %s -> %s", session.from, to)
		}
	}

	return nil
}

// SaveEmailToSent 保存邮件到已发送文件夹
func (s *Service) SaveEmailToSent(mailboxID int64, fromAddr, toAddr, subject, body string) error {
	return s.svcCtx.EmailModel.SaveEmailToFolder(nil, mailboxID, fromAddr, toAddr, subject, body, "sent")
}

// parseEmailContent 解析邮件内容
func (session *SMTPSession) parseEmailContent() (subject, body string) {
	content := string(session.data)

	// 首先尝试使用enmime解析
	if parsedSubject, parsedBody := session.parseEmailWithEnmime(content); parsedSubject != "" || parsedBody != "" {
		log.Printf("enmime解析成功")
		return parsedSubject, parsedBody
	}

	log.Printf("enmime解析失败，使用原始解析方法")

	// 回退到原始解析方法
	lines := strings.Split(content, "\n")

	var inHeaders = true
	var bodyLines []string
	var contentTransferEncoding string

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
			}
		} else {
			bodyLines = append(bodyLines, line)
		}
	}

	body = strings.Join(bodyLines, "\n")

	// 原始解析方法的简单处理（作为enmime的备用方案）
	log.Printf("使用原始解析方法处理邮件内容")

	// 简单的编码处理
	if contentTransferEncoding == "base64" {
		cleanContent := strings.ReplaceAll(body, "\n", "")
		cleanContent = strings.ReplaceAll(cleanContent, "\r", "")
		cleanContent = strings.TrimSpace(cleanContent)
		if decoded, err := base64.StdEncoding.DecodeString(cleanContent); err == nil {
			body = string(decoded)
		}
	} else if contentTransferEncoding == "quoted-printable" {
		body = decodeQuotedPrintable(body)
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

// processForwardRules 处理邮件转发规则
func (s *Service) processForwardRules(sourceEmail, fromAddr, subject, body string) {
	// 获取该邮箱的活跃转发规则
	rules, err := s.forwardService.GetActiveForwardRules(sourceEmail)
	if err != nil {
		log.Printf("获取转发规则失败: %v", err)
		return
	}

	if len(rules) == 0 {
		log.Printf("邮箱 %s 没有活跃的转发规则", sourceEmail)
		return
	}

	log.Printf("找到 %d 个转发规则，开始处理转发", len(rules))

	for _, rule := range rules {
		log.Printf("处理转发规则: %s -> %s", rule.SourceEmail, rule.TargetEmail)

		// 构建转发邮件的主题
		forwardSubject := subject
		if rule.SubjectPrefix != "" {
			forwardSubject = rule.SubjectPrefix + " " + subject
		}

		// 构建转发邮件的内容
		forwardBody := fmt.Sprintf(`
-------- 转发邮件 --------
发件人: %s
收件人: %s
主题: %s
时间: %s

%s
`, fromAddr, sourceEmail, subject, time.Now().Format("2006-01-02 15:04:05"), body)

		// 发送转发邮件
		err := s.sendForwardEmail(rule.SourceEmail, rule.TargetEmail, forwardSubject, forwardBody)
		if err != nil {
			log.Printf("转发邮件失败: %v", err)
			continue
		}

		// 更新转发次数
		err = s.forwardService.IncrementForwardCount(rule.ID)
		if err != nil {
			log.Printf("更新转发次数失败: %v", err)
		}

		log.Printf("✅ 邮件转发成功: %s -> %s", rule.SourceEmail, rule.TargetEmail)
	}
}

// sendForwardEmail 发送转发邮件
func (s *Service) sendForwardEmail(fromAddr, toAddr, subject, body string) error {
	// 检查目标邮箱是否是本域邮箱
	mailboxID, err := s.svcCtx.MailboxModel.GetIdByEmail(toAddr)

	if err == nil {
		// 目标是本域邮箱，直接保存到收件箱
		log.Printf("转发到本域邮箱: %s", toAddr)
		return s.SaveEmail(mailboxID, fromAddr, toAddr, subject, body)
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// 目标是外部邮箱，通过SMTP发送
		log.Printf("转发到外部邮箱: %s", toAddr)

		// 检查是否为外部邮箱
		if s.smtpClient.IsExternalEmail(toAddr) {
			// 使用SMTP客户端发送到外部邮箱
			err := s.smtpClient.SendEmail(fromAddr, toAddr, subject, body)
			if err != nil {
				log.Printf("外部邮箱转发失败: %v", err)
				return fmt.Errorf("外部邮箱转发失败: %w", err)
			}
			log.Printf("✅ 外部邮箱转发成功: %s -> %s", fromAddr, toAddr)
			return nil
		} else {
			return fmt.Errorf("无效的外部邮箱地址: %s", toAddr)
		}
	} else {
		return fmt.Errorf("查询目标邮箱失败: %w", err)
	}
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

// decodeQuotedPrintable 解码quoted-printable编码
func decodeQuotedPrintable(s string) string {
	// 使用Go标准库的quoted-printable解码器
	reader := quotedprintable.NewReader(strings.NewReader(s))
	decoded, err := io.ReadAll(reader)
	if err != nil {
		log.Printf("quoted-printable解码失败，使用备用方法: %v", err)
		// 如果标准库解码失败，使用备用方法
		return decodeQuotedPrintableFallback(s)
	}

	result := string(decoded)
	log.Printf("quoted-printable解码成功，原长度: %d, 解码后长度: %d", len(s), len(result))
	return result
}

// decodeQuotedPrintableFallback 备用的quoted-printable解码实现
func decodeQuotedPrintableFallback(s string) string {
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

// GetEmails 获取邮件列表
func (s *Service) GetEmails(mailboxID int64, folder string, page, limit int) ([]*model.Email, int64, error) {
	return s.svcCtx.EmailModel.GetEmailsByMailboxId(mailboxID, folder, page, limit)
}

// GetEmailByID 根据ID获取邮件
func (s *Service) GetEmailByID(emailID, mailboxID int64) (*model.Email, error) {
	email, err := s.svcCtx.EmailModel.GetByIdAndMailboxId(emailID, mailboxID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("邮件不存在")
		}
		return nil, err
	}

	return email, nil
}

// MarkAsRead 标记邮件为已读
func (s *Service) MarkAsRead(emailID, mailboxID int64) error {
	return s.svcCtx.EmailModel.MapUpdate(nil, emailID, map[string]interface{}{
		"is_read":    true,
		"updated_at": time.Now(),
	})
}

// DeleteEmail 删除邮件
func (s *Service) DeleteEmail(emailID, mailboxID int64) error {
	// 先验证邮件是否存在且属于指定邮箱
	email, err := s.svcCtx.EmailModel.GetByIdAndMailboxId(emailID, mailboxID)
	if err != nil {
		return err
	}

	// 删除邮件
	return s.svcCtx.EmailModel.Delete(nil, email)
}

// SendTestForwardEmail 发送测试转发邮件
func (s *Service) SendTestForwardEmail(sourceEmail, targetEmail, subject, content string, rule interface{}) error {
	// 获取源邮箱ID
	mailboxID, err := s.svcCtx.MailboxModel.GetIdByEmail(sourceEmail)
	if err != nil {
		return fmt.Errorf("获取源邮箱失败: %w", err)
	}

	// 保存测试邮件到源邮箱
	err = s.svcCtx.EmailModel.SaveEmailToFolder(nil, mailboxID, "system@test.com", sourceEmail, subject, content, "inbox")
	if err != nil {
		return fmt.Errorf("保存测试邮件失败: %w", err)
	}

	// 触发转发规则处理
	s.processForwardRules(sourceEmail, "system@test.com", subject, content)

	return nil
}

// parseEmailWithEnmime 使用enmime库解析邮件内容
func (session *SMTPSession) parseEmailWithEnmime(rawEmail string) (subject, body string) {
	log.Printf("开始使用enmime解析邮件")

	// 创建enmime解析器，禁用字符检测让库自己处理
	parser := enmime.NewParser(enmime.DisableCharacterDetection(true))

	// 解析邮件
	env, err := parser.ReadEnvelope(strings.NewReader(rawEmail))
	if err != nil {
		log.Printf("enmime解析失败: %v", err)
		return "", ""
	}

	log.Printf("enmime解析成功")

	// 获取主题
	subject = env.GetHeader("Subject")
	if subject != "" {
		log.Printf("解析到主题: %s", subject)
	}

	// 优先使用HTML内容以保持格式，如果没有则使用文本内容
	if env.HTML != "" {
		log.Printf("使用HTML内容，长度: %d", len(env.HTML))
		body = env.HTML // 保留HTML格式
	} else if env.Text != "" {
		log.Printf("使用文本内容，长度: %d", len(env.Text))
		// 将纯文本转换为HTML格式以便在前端正确显示
		body = strings.ReplaceAll(env.Text, "\n", "<br>")
		body = strings.ReplaceAll(body, "\r", "")
	} else {
		log.Printf("未找到可用的邮件内容")
		return subject, ""
	}

	return subject, body
}
