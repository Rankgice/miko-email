package smtp

import (
	"crypto/tls"
	"enco
	"fmt"
	"log"
	"net"
	"net/smtp"
	"strings"
	"time"

	"nbemail/internal/config"
)

// OutboundClient 外部SMTP客户端
type OutboundClient struct {
	config *config.Config
}

// NewOutboundClient 创建外部SMTP客户端
func NewOutboundClient(cfg *config.Config) *OutboundClient {
	return &OutboundClient{
		config: cfg,
	}
}

// SendEmail 发送邮件到外部邮箱
func (c *OutboundClient) SendEmail(from, to, subject, body string) error {
	// 根据发件人邮箱获取对应的SMTP配置
	smtpConfig := c.config.GetSMTPConfigForEmail(from)
	if smtpConfig == nil {
		// 如果没有SMTP配置，尝试直接MX发送
		log.Printf("未找到SMTP配置，尝试直接MX发送: %s -> %s", from, to)
		return c.sendDirectMX(from, to, subject, body)
	}

	// 确定实际发件人地址
	actualFrom := from
	if smtpConfig.User != "" {
		// 对于第三方SMTP服务器，使用认证用户作为发件人
		actualFrom = smtpConfig.User
	}

	// 构建邮件内容，在邮件头中保留原始发件人，但在Reply-To中设置
	message := c.buildMessage(actualFrom, to, subject, body, from)

	// 设置认证（如果配置了用户名和密码）
	var auth smtp.Auth
	if smtpConfig.User != "" && smtpConfig.Password != "" {
		auth = smtp.PlainAuth("", smtpConfig.User, smtpConfig.Password, smtpConfig.Host)
	}

	// 发送邮件 - 使用实际发件人地址
	addr := fmt.Sprintf("%s:%d", smtpConfig.Host, smtpConfig.Port)

	if smtpConfig.TLS {
		// 对于465端口，使用SSL直接连接
		if smtpConfig.Port == 465 {
			return c.sendWithSSL(addr, auth, actualFrom, []string{to}, message, smtpConfig)
		} else {
			// 对于587端口，使用STARTTLS
			return c.sendWithTLSConfig(addr, auth, actualFrom, []string{to}, message, smtpConfig)
		}
	} else {
		return smtp.SendMail(addr, auth, actualFrom, []string{to}, message)
	}
}

// sendWithTLS 使用TLS发送邮件
func (c *OutboundClient) sendWithTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	// 连接到SMTP服务器
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("连接SMTP服务器失败: %v", err)
	}
	defer client.Close()

	// 启动TLS
	tlsConfig := &tls.Config{
		ServerName: c.config.OutboundSMTPHost,
	}

	if err = client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("启动TLS失败: %v", err)
	}

	// 认证（如果提供了认证信息）
	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP认证失败: %v", err)
		}
	}

	// 设置发件人
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("设置发件人失败: %v", err)
	}

	// 设置收件人
	for _, addr := range to {
		if err = client.Rcpt(addr); err != nil {
			return fmt.Errorf("设置收件人失败: %v", err)
		}
	}

	// 发送邮件内容
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("获取数据写入器失败: %v", err)
	}

	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("写入邮件内容失败: %v", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("关闭数据写入器失败: %v", err)
	}

	return client.Quit()
}

// sendWithTLSConfig 使用指定的SMTP配置和TLS发送邮件
func (c *OutboundClient) sendWithTLSConfig(addr string, auth smtp.Auth, from string, to []string, msg []byte, smtpConfig *config.SMTPConfig) error {
	// 使用重试机制连接到SMTP服务器
	var client *smtp.Client
	var err error
	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		log.Printf("尝试连接SMTP服务器 %s (第%d次)", addr, i+1)

		// 设置连接超时
		conn, dialErr := net.DialTimeout("tcp", addr, 30*time.Second)
		if dialErr != nil {
			err = dialErr
			log.Printf("TCP连接失败 (第%d次): %v", i+1, err)
			if i < maxRetries-1 {
				time.Sleep(time.Duration(i+1) * 2 * time.Second)
			}
			continue
		}

		client, err = smtp.NewClient(conn, smtpConfig.Host)
		if err == nil {
			break
		}

		conn.Close()
		log.Printf("创建SMTP客户端失败 (第%d次): %v", i+1, err)

		// 如果不是最后一次重试，等待一段时间再重试
		if i < maxRetries-1 {
			time.Sleep(time.Duration(i+1) * 2 * time.Second)
		}
	}

	if err != nil {
		return fmt.Errorf("连接SMTP服务器失败，已重试%d次: %v", maxRetries, err)
	}
	defer client.Close()

	log.Printf("SMTP连接成功，启动TLS")

	// 启动TLS，增加更多选项以提高兼容性
	tlsConfig := &tls.Config{
		ServerName:         smtpConfig.Host,
		InsecureSkipVerify: false,            // 保持安全验证
		MinVersion:         tls.VersionTLS12, // 最低TLS版本
		MaxVersion:         tls.VersionTLS13, // 最高TLS版本
	}

	if err = client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("启动TLS失败: %v", err)
	}

	log.Printf("TLS启动成功")

	// 认证（如果提供了认证信息）
	if auth != nil {
		log.Printf("开始SMTP认证")
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP认证失败: %v", err)
		}
		log.Printf("SMTP认证成功")
	}

	// 设置发件人
	log.Printf("设置发件人: %s", from)
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("设置发件人失败: %v", err)
	}

	// 设置收件人
	for _, addr := range to {
		log.Printf("设置收件人: %s", addr)
		if err = client.Rcpt(addr); err != nil {
			return fmt.Errorf("设置收件人失败: %v", err)
		}
	}

	// 发送邮件内容
	log.Printf("开始发送邮件内容")
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("开始发送邮件内容失败: %v", err)
	}

	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("写入邮件内容失败: %v", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("关闭数据写入器失败: %v", err)
	}

	log.Printf("邮件发送完成，关闭连接")
	return client.Quit()
}

// sendWithSSL 使用SSL直接连接发送邮件（适用于465端口）
func (c *OutboundClient) sendWithSSL(addr string, auth smtp.Auth, from string, to []string, msg []byte, smtpConfig *config.SMTPConfig) error {
	// 创建TLS配置，增加更多选项以提高兼容性
	tlsConfig := &tls.Config{
		ServerName:         smtpConfig.Host,
		InsecureSkipVerify: false,            // 保持安全验证
		MinVersion:         tls.VersionTLS12, // 最低TLS版本
		MaxVersion:         tls.VersionTLS13, // 最高TLS版本
	}

	// 使用重试机制进行SSL连接
	var conn *tls.Conn
	var err error
	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		log.Printf("尝试SSL连接到 %s (第%d次)", addr, i+1)

		// 设置连接超时
		dialer := &net.Dialer{
			Timeout: 30 * time.Second,
		}

		conn, err = tls.DialWithDialer(dialer, "tcp", addr, tlsConfig)
		if err == nil {
			break
		}

		log.Printf("SSL连接失败 (第%d次): %v", i+1, err)

		// 如果不是最后一次重试，等待一段时间再重试
		if i < maxRetries-1 {
			time.Sleep(time.Duration(i+1) * 2 * time.Second)
		}
	}

	if err != nil {
		return fmt.Errorf("SSL连接失败，已重试%d次: %v", maxRetries, err)
	}
	defer conn.Close()

	log.Printf("SSL连接成功，创建SMTP客户端")

	// 创建SMTP客户端
	client, err := smtp.NewClient(conn, smtpConfig.Host)
	if err != nil {
		return fmt.Errorf("创建SMTP客户端失败: %v", err)
	}
	defer client.Close()

	// 认证（如果提供了认证信息）
	if auth != nil {
		log.Printf("开始SMTP认证")
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP认证失败: %v", err)
		}
		log.Printf("SMTP认证成功")
	}

	// 设置发件人
	log.Printf("设置发件人: %s", from)
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("设置发件人失败: %v", err)
	}

	// 设置收件人
	for _, addr := range to {
		log.Printf("设置收件人: %s", addr)
		if err = client.Rcpt(addr); err != nil {
			return fmt.Errorf("设置收件人失败: %v", err)
		}
	}

	// 发送邮件内容
	log.Printf("开始发送邮件内容")
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("开始发送邮件内容失败: %v", err)
	}

	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("写入邮件内容失败: %v", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("关闭数据写入器失败: %v", err)
	}

	log.Printf("邮件发送完成，关闭连接")
	return client.Quit()
}

// buildMessage 构建邮件消息
func (c *OutboundClient) buildMessage(from, to, subject, body string, originalFrom ...string) []byte {
	// 构建邮件头
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to

	// 如果提供了原始发件人，设置Reply-To
	if len(originalFrom) > 0 && originalFrom[0] != from {
		headers["Reply-To"] = originalFrom[0]
	}

	// 对主题进行MIME编码（如果包含非ASCII字符）
	if needsMIMEEncoding(subject) {
		headers["Subject"] = encodeMIMEHeader(subject)
	} else {
		headers["Subject"] = subject
	}

	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/plain; charset=UTF-8"

	// 对邮件正文进行编码处理
	var encodedBody string
	var transferEncoding string

	if needsMIMEEncoding(body) {
		// 使用Base64编码处理包含中文的内容
		encodedBody = base64.StdEncoding.EncodeToString([]byte(body))
		transferEncoding = "base64"

		// 将Base64编码的内容按76字符换行（RFC标准）
		var formattedBody strings.Builder
		for i := 0; i < len(encodedBody); i += 76 {
			end := i + 76
			if end > len(encodedBody) {
				end = len(encodedBody)
			}
			formattedBody.WriteString(encodedBody[i:end])
			if end < len(encodedBody) {
				formattedBody.WriteString("\r\n")
			}
		}
		encodedBody = formattedBody.String()
	} else {
		// 纯ASCII内容，直接使用
		encodedBody = body
		transferEncoding = "7bit"
	}

	headers["Content-Transfer-Encoding"] = transferEncoding

	// 添加自定义头部标识
	headers["X-Mailer"] = "NBEmail System"

	// 构建完整消息
	var message strings.Builder
	for k, v := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	message.WriteString("\r\n")
	message.WriteString(encodedBody)

	return []byte(message.String())
}

// IsExternalEmail 检查是否为外部邮箱
func (c *OutboundClient) IsExternalEmail(email string) bool {
	if !strings.Contains(email, "@") {
		return false
	}

	domain := strings.Split(email, "@")[1]
	return domain != c.config.Domain && domain != "localhost"
}

// sendDirectMX 直接通过MX记录发送邮件
func (c *OutboundClient) sendDirectMX(from, to, subject, body string) error {
	// 提取收件人域名
	parts := strings.Split(to, "@")
	if len(parts) != 2 {
		return fmt.Errorf("无效的邮箱地址: %s", to)
	}
	domain := parts[1]

	// 查询MX记录
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		return fmt.Errorf("查询MX记录失败: %v", err)
	}

	if len(mxRecords) == 0 {
		return fmt.Errorf("域名 %s 没有MX记录", domain)
	}

	// 按优先级排序，选择优先级最高的MX记录
	var bestMX *net.MX
	for _, mx := range mxRecords {
		if bestMX == nil || mx.Pref < bestMX.Pref {
			bestMX = mx
		}
	}

	mxHost := strings.TrimSuffix(bestMX.Host, ".")
	log.Printf("使用MX服务器: %s (优先级: %d)", mxHost, bestMX.Pref)

	// 构建邮件内容
	message := c.buildMessage(from, to, subject, body, from)

	// 尝试连接到MX服务器的25端口
	addr := fmt.Sprintf("%s:25", mxHost)

	// 使用重试机制
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		log.Printf("尝试连接MX服务器 %s (第%d次)", addr, i+1)

		err = c.sendDirectSMTP(addr, from, to, message, mxHost)
		if err == nil {
			log.Printf("直接MX发送成功")
			return nil
		}

		log.Printf("MX发送失败 (第%d次): %v", i+1, err)

		// 如果不是最后一次重试，等待一段时间再重试
		if i < maxRetries-1 {
			time.Sleep(time.Duration(i+1) * 2 * time.Second)
		}
	}

	return fmt.Errorf("直接MX发送失败，已重试%d次: %v", maxRetries, err)
}

// sendDirectSMTP 直接SMTP发送（无认证）
func (c *OutboundClient) sendDirectSMTP(addr, from, to string, message []byte, hostname string) error {
	// 连接到SMTP服务器
	conn, err := net.DialTimeout("tcp", addr, 30*time.Second)
	if err != nil {
		return fmt.Errorf("连接失败: %v", err)
	}
	defer conn.Close()

	// 创建SMTP客户端
	client, err := smtp.NewClient(conn, hostname)
	if err != nil {
		return fmt.Errorf("创建SMTP客户端失败: %v", err)
	}
	defer client.Close()

	// 设置发件人
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("设置发件人失败: %v", err)
	}

	// 设置收件人
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("设置收件人失败: %v", err)
	}

	// 发送邮件内容
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("开始发送邮件内容失败: %v", err)
	}

	_, err = w.Write(message)
	if err != nil {
		return fmt.Errorf("写入邮件内容失败: %v", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("关闭数据写入器失败: %v", err)
	}

	return client.Quit()
}

// LogSendAttempt 记录发送尝试
func (c *OutboundClient) LogSendAttempt(from, to, subject string, err error) {
	if err != nil {
		log.Printf("外部邮件发送失败 - From: %s, To: %s, Subject: %s, Error: %v", from, to, subject, err)
	} else {
		log.Printf("外部邮件发送成功 - From: %s, To: %s, Subject: %s", from, to, subject)
	}
}

// needsMIMEEncoding 检查字符串是否需要MIME编码
func needsMIMEEncoding(s string) bool {
	for _, r := range s {
		if r > 127 {
			return true
		}
	}
	return false
}

// encodeMIMEHeader 对邮件头部进行MIME编码
func encodeMIMEHeader(s string) string {
	if !needsMIMEEncoding(s) {
		return s
	}

	// 使用Base64编码
	encoded := base64.StdEncoding.EncodeToString([]byte(s))
	return fmt.Sprintf("=?UTF-8?B?%s?=", encoded)
}
