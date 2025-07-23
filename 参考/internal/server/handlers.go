package server

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"nbemail/internal/config"
	"nbemail/internal/models"
	"nbemail/internal/smtp"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

// handleLogin 处理登录
func (s *Server) handleLogin(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "请求参数错误",
		})
		return
	}

	// 查询用户
	var user models.User
	err := s.db.QueryRow(`
		SELECT id, email, password, name, is_admin, created_at, updated_at
		FROM users WHERE email = ?
	`, req.Email).Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusUnauthorized, models.Response{
			Success: false,
			Message: "邮箱或密码错误",
		})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, models.Response{
			Success: false,
			Message: "邮箱或密码错误",
		})
		return
	}

	// 设置Cookie
	c.SetCookie("token", user.Email, 3600*24*7, "/", "", false, true)
	c.SetCookie("user_id", strconv.Itoa(user.ID), 3600*24*7, "/", "", false, true)

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "登录成功",
		Data:    user,
	})
}

// handleLogout 处理登出
func (s *Server) handleLogout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.SetCookie("user_id", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "登出成功",
	})
}

// handleGetEmails 获取邮件列表
func (s *Server) handleGetEmails(c *gin.Context) {
	userID := c.GetInt("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	folder := c.DefaultQuery("folder", "inbox")
	search := c.Query("search") // 搜索关键词

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	// 获取当前用户的邮箱地址
	var userEmail string
	err := s.db.QueryRow("SELECT email FROM users WHERE id = ?", userID).Scan(&userEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "获取用户信息失败",
		})
		return
	}

	var query string
	var countQuery string
	var args []interface{}
	var countArgs []interface{}

	// 构建搜索条件
	searchCondition := ""
	if search != "" {
		searchCondition = " AND (subject LIKE ? OR from_addr LIKE ? OR to_addr LIKE ? OR body LIKE ?)"
		searchParams := []interface{}{
			"%" + search + "%",
			"%" + search + "%",
			"%" + search + "%",
			"%" + search + "%",
		}

		switch folder {
		case "sent":
			query = `SELECT id, message_id, from_addr, to_addr, subject, body, html_body, is_read, is_deleted, is_sent, user_id, size, attachments, created_at, updated_at
					FROM emails WHERE user_id = ? AND is_sent = 1 AND is_deleted = 0` + searchCondition + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
			countQuery = "SELECT COUNT(*) FROM emails WHERE user_id = ? AND is_sent = 1 AND is_deleted = 0" + searchCondition
			args = append([]interface{}{userID}, searchParams...)
			args = append(args, limit, offset)
			countArgs = append([]interface{}{userID}, searchParams...)
		default: // inbox
			query = `SELECT id, message_id, from_addr, to_addr, subject, body, html_body, is_read, is_deleted, is_sent, user_id, size, attachments, created_at, updated_at
					FROM emails WHERE user_id = ? AND is_sent = 0 AND is_deleted = 0` + searchCondition + ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
			countQuery = "SELECT COUNT(*) FROM emails WHERE user_id = ? AND is_sent = 0 AND is_deleted = 0" + searchCondition
			args = append([]interface{}{userID}, searchParams...)
			args = append(args, limit, offset)
			countArgs = append([]interface{}{userID}, searchParams...)
		}
	} else {
		switch folder {
		case "sent":
			query = `SELECT id, message_id, from_addr, to_addr, subject, body, html_body, is_read, is_deleted, is_sent, user_id, size, attachments, created_at, updated_at
					FROM emails WHERE user_id = ? AND is_sent = 1 AND is_deleted = 0 ORDER BY created_at DESC LIMIT ? OFFSET ?`
			countQuery = "SELECT COUNT(*) FROM emails WHERE user_id = ? AND is_sent = 1 AND is_deleted = 0"
			args = []interface{}{userID, limit, offset}
			countArgs = []interface{}{userID}
		default: // inbox
			query = `SELECT id, message_id, from_addr, to_addr, subject, body, html_body, is_read, is_deleted, is_sent, user_id, size, attachments, created_at, updated_at
					FROM emails WHERE user_id = ? AND is_sent = 0 AND is_deleted = 0 ORDER BY created_at DESC LIMIT ? OFFSET ?`
			countQuery = "SELECT COUNT(*) FROM emails WHERE user_id = ? AND is_sent = 0 AND is_deleted = 0"
			args = []interface{}{userID, limit, offset}
			countArgs = []interface{}{userID}
		}
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "查询邮件失败",
		})
		return
	}
	defer rows.Close()

	var emails []models.Email
	for rows.Next() {
		var email models.Email
		var htmlBody, attachments sql.NullString
		err := rows.Scan(&email.ID, &email.MessageID, &email.From, &email.To, &email.Subject,
			&email.Body, &htmlBody, &email.IsRead, &email.IsDeleted, &email.IsSent,
			&email.UserID, &email.Size, &attachments, &email.CreatedAt, &email.UpdatedAt)
		if err != nil {
			log.Printf("扫描邮件数据失败: %v", err)
			continue
		}

		// 处理可能为NULL的字段
		if htmlBody.Valid {
			email.HTMLBody = htmlBody.String
		}
		if attachments.Valid {
			email.Attachments = attachments.String
		}

		// 解析MIME内容（与详情页面使用相同的逻辑）

		// 检查是否是MIME多部分邮件
		if strings.Contains(email.Body, "boundary=") || strings.Contains(email.Body, "Content-Type:") {
			textParts, htmlParts := parseMIMEContentForAPI(email.Body)

			// 如果有解析出的内容，使用解析后的内容
			if len(textParts) > 0 {
				email.Body = strings.Join(textParts, "\n\n")
			}

			if len(htmlParts) > 0 {
				email.HTMLBody = strings.Join(htmlParts, "\n\n")
			}

			// 如果没有解析出任何内容，尝试简单的Base64解码
			if len(textParts) == 0 && len(htmlParts) == 0 {
				email.Body = decodeIfBase64(email.Body)
			}
		} else {
			// 如果没有MIME结构，尝试简单的Base64解码
			email.Body = decodeIfBase64(email.Body)
		}

		email.Subject = decodeIfBase64(email.Subject)

		emails = append(emails, email)
	}

	// 获取总数
	var total int
	s.db.QueryRow(countQuery, countArgs...).Scan(&total)

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data: models.EmailListResponse{
			Emails: emails,
			Total:  total,
			Page:   page,
			Limit:  limit,
		},
	})
}

// handleGetEmail 获取单个邮件
func (s *Server) handleGetEmail(c *gin.Context) {
	userID := c.GetInt("user_id")
	emailID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "邮件ID无效",
		})
		return
	}

	var email models.Email
	var htmlBody, attachments sql.NullString
	err = s.db.QueryRow(`
		SELECT id, message_id, from_addr, to_addr, subject, body, html_body, is_read, is_deleted, is_sent, user_id, size, attachments, created_at, updated_at
		FROM emails WHERE id = ? AND user_id = ?
	`, emailID, userID).Scan(&email.ID, &email.MessageID, &email.From, &email.To, &email.Subject,
		&email.Body, &htmlBody, &email.IsRead, &email.IsDeleted, &email.IsSent,
		&email.UserID, &email.Size, &attachments, &email.CreatedAt, &email.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, models.Response{
			Success: false,
			Message: "邮件不存在",
		})
		return
	}

	// 处理可能为NULL的字段
	if htmlBody.Valid {
		email.HTMLBody = htmlBody.String
	}
	if attachments.Valid {
		email.Attachments = attachments.String
	}

	// 解析MIME内容
	originalBody := email.Body

	// 检查是否是MIME多部分邮件
	if strings.Contains(email.Body, "boundary=") || strings.Contains(email.Body, "Content-Type:") {
		textParts, htmlParts := parseMIMEContentForAPI(email.Body)

		// 如果有解析出的内容，使用解析后的内容
		if len(textParts) > 0 {
			email.Body = strings.Join(textParts, "\n\n")
			log.Printf("Web API - MIME解析成功，文本部分: %d个", len(textParts))
		}

		if len(htmlParts) > 0 {
			email.HTMLBody = strings.Join(htmlParts, "\n\n")
			log.Printf("Web API - MIME解析成功，HTML部分: %d个", len(htmlParts))
		}

		// 如果没有解析出任何内容，尝试简单的Base64解码
		if len(textParts) == 0 && len(htmlParts) == 0 {
			email.Body = decodeIfBase64(email.Body)
		}
	} else {
		// 如果没有MIME结构，尝试简单的Base64解码
		email.Body = decodeIfBase64(email.Body)
	}

	email.Subject = decodeIfBase64(email.Subject)

	// 调试日志
	if originalBody != email.Body {
		log.Printf("Web API - 邮件内容解析成功 - 原始长度: %d, 解析后长度: %d", len(originalBody), len(email.Body))
	}

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    email,
	})
}

// handleSendEmail 发送邮件
func (s *Server) handleSendEmail(c *gin.Context) {
	// 设置正确的Content-Type
	c.Header("Content-Type", "application/json; charset=utf-8")

	// 读取原始请求体
	rawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("读取请求体失败: %v", err)
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "请求参数错误",
		})
		return
	}

	// 尝试修复可能的UTF-8编码问题
	bodyStr := string(rawBody)
	if !utf8.ValidString(bodyStr) {
		// 如果不是有效的UTF-8，尝试修复
		bodyStr = strings.ToValidUTF8(bodyStr, "")
		log.Printf("修复了请求体中的无效UTF-8字符")
	}

	// 重新创建请求体
	c.Request.Body = io.NopCloser(strings.NewReader(bodyStr))

	var req models.SendEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("JSON绑定错误: %v", err)
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "请求参数错误",
		})
		return
	}

	// 修复可能的UTF-8编码问题
	if !utf8.ValidString(req.Subject) {
		req.Subject = strings.ToValidUTF8(req.Subject, "")
		log.Printf("修复了主题中的无效UTF-8字符")
	}
	if !utf8.ValidString(req.Body) {
		req.Body = strings.ToValidUTF8(req.Body, "")
		log.Printf("修复了邮件正文中的无效UTF-8字符")
	}

	// 调试：打印接收到的数据
	log.Printf("接收到的邮件数据 - Subject: %s (len=%d), Body: %s (len=%d)",
		req.Subject, len(req.Subject), req.Body, len(req.Body))

	userID := c.GetInt("user_id")
	userEmail := c.GetString("user_email")

	// 确定发件人邮箱
	fromEmail := userEmail // 默认使用用户邮箱
	if req.From != "" {
		// 验证指定的发件人邮箱是否属于当前用户
		var count int
		err := s.db.QueryRow("SELECT COUNT(*) FROM mailboxes WHERE email = ? AND user_id = ? AND is_active = 1",
			req.From, userID).Scan(&count)
		if err != nil || count == 0 {
			c.JSON(http.StatusBadRequest, models.Response{
				Success: false,
				Message: "指定的发件人邮箱无效或不属于您",
			})
			return
		}
		fromEmail = req.From
	}

	// 生成消息ID
	messageID := fmt.Sprintf("<%d.%d@%s>", time.Now().Unix(), userID, "localhost")

	// 先检查收件人是否为本地用户
	var recipientUserID int
	err = s.db.QueryRow("SELECT user_id FROM mailboxes WHERE email = ? AND is_active = 1", req.To).Scan(&recipientUserID)

	var sendSuccess bool
	var sendError error

	if err == nil {
		// 收件人是本地用户，直接发送成功
		sendSuccess = true

		// 为收件人创建收件箱记录
		recipientMessageID := fmt.Sprintf("<%d.%d.inbox@%s>", time.Now().Unix(), recipientUserID, "localhost")
		_, err = s.db.Exec(`
			INSERT INTO emails (message_id, from_addr, to_addr, subject, body, user_id, is_sent, size, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, recipientMessageID, fromEmail, req.To, req.Subject, req.Body, recipientUserID, false, len(req.Body), time.Now(), time.Now())

		if err != nil {
			log.Printf("为收件人创建收件箱记录失败: %v", err)
		}
		log.Printf("本地邮件发送成功 - From: %s, To: %s", fromEmail, req.To)
	} else {
		// 收件人不是本地用户，尝试通过外部SMTP发送
		outboundClient := smtp.NewOutboundClient(s.config)

		if outboundClient.IsExternalEmail(req.To) {
			sendError = outboundClient.SendEmail(fromEmail, req.To, req.Subject, req.Body)
			outboundClient.LogSendAttempt(fromEmail, req.To, req.Subject, sendError)

			if sendError == nil {
				sendSuccess = true
				log.Printf("外部邮件发送成功 - From: %s, To: %s", fromEmail, req.To)
			} else {
				log.Printf("外部邮件发送失败: %v", sendError)
			}
		} else {
			sendError = fmt.Errorf("收件人邮箱格式无效或为未知本地用户: %s", req.To)
			log.Printf("发送失败: %v", sendError)
		}
	}

	// 只有发送成功才保存到发件箱
	if sendSuccess {
		_, err = s.db.Exec(`
			INSERT INTO emails (message_id, from_addr, to_addr, subject, body, user_id, is_sent, size, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, messageID, fromEmail, req.To, req.Subject, req.Body, userID, true, len(req.Body), time.Now(), time.Now())

		if err != nil {
			log.Printf("保存到发件箱失败: %v", err)
			c.JSON(http.StatusInternalServerError, models.Response{
				Success: false,
				Message: "邮件发送成功但保存失败",
			})
			return
		}

		c.JSON(http.StatusOK, models.Response{
			Success: true,
			Message: "邮件发送成功",
			Data: map[string]interface{}{
				"from": fromEmail,
				"to":   req.To,
			},
		})
	} else {
		// 发送失败，返回错误
		errorMessage := "邮件发送失败"
		if sendError != nil {
			errorMessage = fmt.Sprintf("邮件发送失败: %v", sendError)
		}

		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: errorMessage,
		})
	}
}

// SMTP配置管理相关处理函数

// handleGetSMTPConfigs 获取SMTP配置列表
func (s *Server) handleGetSMTPConfigs(c *gin.Context) {
	if !c.GetBool("is_admin") {
		c.JSON(http.StatusForbidden, models.Response{
			Success: false,
			Message: "权限不足",
		})
		return
	}

	configs := make([]map[string]interface{}, 0)

	// 添加默认配置
	if s.config.OutboundSMTPHost != "" {
		configs = append(configs, map[string]interface{}{
			"domain":     "default",
			"host":       s.config.OutboundSMTPHost,
			"port":       s.config.OutboundSMTPPort,
			"user":       s.config.OutboundSMTPUser,
			"password":   "***", // 隐藏密码
			"tls":        s.config.OutboundSMTPTLS,
			"is_default": true,
		})
	}

	// 添加域名特定配置
	for domain, smtpConfig := range s.config.DomainSMTPConfigs {
		configs = append(configs, map[string]interface{}{
			"domain":     domain,
			"host":       smtpConfig.Host,
			"port":       smtpConfig.Port,
			"user":       smtpConfig.User,
			"password":   "***", // 隐藏密码
			"tls":        smtpConfig.TLS,
			"is_default": false,
		})
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    configs,
	})
}

// SMTPConfigRequest SMTP配置请求结构
type SMTPConfigRequest struct {
	Domain   string `json:"domain" binding:"required"`
	Host     string `json:"host" binding:"required"`
	Port     int    `json:"port" binding:"required"`
	User     string `json:"user"`
	Password string `json:"password"`
	TLS      bool   `json:"tls"`
}

// handleAddSMTPConfig 添加SMTP配置
func (s *Server) handleAddSMTPConfig(c *gin.Context) {
	if !c.GetBool("is_admin") {
		c.JSON(http.StatusForbidden, models.Response{
			Success: false,
			Message: "权限不足",
		})
		return
	}

	var req SMTPConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "请求参数错误",
		})
		return
	}

	// 添加到配置中
	s.config.AddDomainSMTPConfig(req.Domain, &config.SMTPConfig{
		Host:     req.Host,
		Port:     req.Port,
		User:     req.User,
		Password: req.Password,
		TLS:      req.TLS,
	})

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "SMTP配置添加成功",
	})
}

// handleDeleteSMTPConfig 删除SMTP配置
func (s *Server) handleDeleteSMTPConfig(c *gin.Context) {
	if !c.GetBool("is_admin") {
		c.JSON(http.StatusForbidden, models.Response{
			Success: false,
			Message: "权限不足",
		})
		return
	}

	domain := c.Param("domain")
	if domain == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "域名参数不能为空",
		})
		return
	}

	if domain == "default" {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "不能删除默认配置",
		})
		return
	}

	// 从配置中删除
	s.config.RemoveDomainSMTPConfig(domain)

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "SMTP配置删除成功",
	})
}

// handleAutoConfigSMTP 自动配置SMTP
func (s *Server) handleAutoConfigSMTP(c *gin.Context) {
	if !c.GetBool("is_admin") {
		c.JSON(http.StatusForbidden, models.Response{
			Success: false,
			Message: "权限不足",
		})
		return
	}

	// 获取所有邮箱地址
	rows, err := s.db.Query("SELECT DISTINCT email FROM mailboxes WHERE is_active = 1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "获取邮箱列表失败",
		})
		return
	}
	defer rows.Close()

	var mailboxes []string
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			continue
		}
		mailboxes = append(mailboxes, email)
	}

	// 提取域名
	domains := s.config.GetDomainsFromMailboxes(mailboxes)

	// 自动配置SMTP
	s.config.AutoConfigureDomainSMTP(domains)

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: fmt.Sprintf("自动配置完成，检测到 %d 个域名", len(domains)),
		Data: map[string]interface{}{
			"domains": domains,
			"count":   len(domains),
		},
	})
}

// handleVerifySMTPConfig 验证SMTP配置的DNS记录
func (s *Server) handleVerifySMTPConfig(c *gin.Context) {
	if !c.GetBool("is_admin") {
		c.JSON(http.StatusForbidden, models.Response{
			Success: false,
			Message: "权限不足",
		})
		return
	}

	domain := c.Param("domain")
	if domain == "" {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "域名参数不能为空",
		})
		return
	}

	// 获取SMTP配置
	smtpConfig := s.config.GetSMTPConfigForDomain(domain)
	if smtpConfig == nil {
		c.JSON(http.StatusNotFound, models.Response{
			Success: false,
			Message: "未找到该域名的SMTP配置",
		})
		return
	}

	// 验证SMTP服务器DNS记录
	result := s.verifySMTPDNS(domain, smtpConfig)

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "SMTP DNS验证完成",
		Data:    result,
	})
}

// verifySMTPDNS 验证SMTP配置的DNS记录
func (s *Server) verifySMTPDNS(domain string, smtpConfig *config.SMTPConfig) map[string]interface{} {
	result := map[string]interface{}{
		"domain":    domain,
		"smtp_host": smtpConfig.Host,
		"smtp_port": smtpConfig.Port,
		"checks":    []map[string]interface{}{},
		"overall":   false,
		"message":   "",
	}

	checks := []map[string]interface{}{}

	// 1. 检查SMTP服务器主机名解析
	hostCheck := map[string]interface{}{
		"name":        "SMTP主机名解析",
		"description": fmt.Sprintf("检查 %s 是否可以解析", smtpConfig.Host),
		"success":     false,
		"message":     "",
	}

	if ips, err := net.LookupHost(smtpConfig.Host); err == nil && len(ips) > 0 {
		hostCheck["success"] = true
		hostCheck["message"] = fmt.Sprintf("解析成功，IP地址: %s", strings.Join(ips, ", "))
		hostCheck["ips"] = ips
	} else {
		hostCheck["message"] = fmt.Sprintf("解析失败: %v", err)
	}
	checks = append(checks, hostCheck)

	// 2. 检查SMTP服务器连接
	connCheck := map[string]interface{}{
		"name":        "SMTP服务器连接",
		"description": fmt.Sprintf("检查 %s:%d 是否可以连接", smtpConfig.Host, smtpConfig.Port),
		"success":     false,
		"message":     "",
	}

	if conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", smtpConfig.Host, smtpConfig.Port), 10*time.Second); err == nil {
		conn.Close()
		connCheck["success"] = true
		connCheck["message"] = "连接成功"
	} else {
		connCheck["message"] = fmt.Sprintf("连接失败: %v", err)
	}
	checks = append(checks, connCheck)

	// 3. 检查域名MX记录
	mxCheck := map[string]interface{}{
		"name":        "域名MX记录",
		"description": fmt.Sprintf("检查 %s 的MX记录", domain),
		"success":     false,
		"message":     "",
	}

	if mxRecords, err := net.LookupMX(domain); err == nil && len(mxRecords) > 0 {
		mxCheck["success"] = true
		var mxHosts []string
		for _, mx := range mxRecords {
			mxHost := strings.TrimSuffix(mx.Host, ".")
			mxHosts = append(mxHosts, fmt.Sprintf("%s (优先级: %d)", mxHost, mx.Pref))
		}
		mxCheck["message"] = fmt.Sprintf("找到 %d 条MX记录: %s", len(mxRecords), strings.Join(mxHosts, ", "))
		mxCheck["mx_records"] = mxHosts
	} else {
		mxCheck["message"] = fmt.Sprintf("未找到MX记录或查询失败: %v", err)
	}
	checks = append(checks, mxCheck)

	// 4. 检查SMTP主机是否在MX记录中
	mxMatchCheck := map[string]interface{}{
		"name":        "SMTP主机匹配",
		"description": fmt.Sprintf("检查SMTP主机 %s 是否在域名 %s 的MX记录中", smtpConfig.Host, domain),
		"success":     false,
		"message":     "",
	}

	if mxRecords, err := net.LookupMX(domain); err == nil {
		found := false
		for _, mx := range mxRecords {
			mxHost := strings.TrimSuffix(mx.Host, ".")
			if mxHost == smtpConfig.Host {
				found = true
				break
			}
		}
		if found {
			mxMatchCheck["success"] = true
			mxMatchCheck["message"] = "SMTP主机在MX记录中找到，配置正确"
		} else {
			mxMatchCheck["message"] = "SMTP主机不在MX记录中，可能需要检查配置"
		}
	} else {
		mxMatchCheck["message"] = "无法查询MX记录进行匹配检查"
	}
	checks = append(checks, mxMatchCheck)

	result["checks"] = checks

	// 计算总体结果
	successCount := 0
	for _, check := range checks {
		if check["success"].(bool) {
			successCount++
		}
	}

	if successCount >= 2 { // 至少主机解析和连接成功
		result["overall"] = true
		result["message"] = fmt.Sprintf("验证通过 (%d/%d 项检查成功)", successCount, len(checks))
	} else {
		result["message"] = fmt.Sprintf("验证失败 (%d/%d 项检查成功)", successCount, len(checks))
	}

	return result
}

// handleMarkEmailAsRead 标记邮件为已读
func (s *Server) handleMarkEmailAsRead(c *gin.Context) {
	userID := c.GetInt("user_id")
	emailID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "邮件ID无效",
		})
		return
	}

	_, err = s.db.Exec("UPDATE emails SET is_read = 1, updated_at = ? WHERE id = ? AND user_id = ?",
		time.Now(), emailID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "标记失败",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "标记成功",
	})
}

// handleDeleteEmail 删除邮件
func (s *Server) handleDeleteEmail(c *gin.Context) {
	userID := c.GetInt("user_id")
	emailID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "邮件ID无效",
		})
		return
	}

	_, err = s.db.Exec("UPDATE emails SET is_deleted = 1, updated_at = ? WHERE id = ? AND user_id = ?",
		time.Now(), emailID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "删除失败",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "删除成功",
	})
}

// handleGetUsers 获取用户列表（管理员）
func (s *Server) handleGetUsers(c *gin.Context) {
	if !c.GetBool("is_admin") {
		c.JSON(http.StatusForbidden, models.Response{
			Success: false,
			Message: "权限不足",
		})
		return
	}

	rows, err := s.db.Query(`
		SELECT id, email, name, is_admin, created_at, updated_at
		FROM users ORDER BY created_at DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "查询用户失败",
		})
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			continue
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    users,
	})
}

// handleCreateUser 创建用户（管理员）
func (s *Server) handleCreateUser(c *gin.Context) {
	if !c.GetBool("is_admin") {
		c.JSON(http.StatusForbidden, models.Response{
			Success: false,
			Message: "权限不足",
		})
		return
	}

	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "请求参数错误",
		})
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "密码加密失败",
		})
		return
	}

	// 插入用户
	_, err = s.db.Exec(`
		INSERT INTO users (email, password, name, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`, req.Email, string(hashedPassword), req.Name, time.Now(), time.Now())

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "创建用户失败",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "用户创建成功",
	})
}

// handleUpdateUser 更新用户（管理员）
func (s *Server) handleUpdateUser(c *gin.Context) {
	if !c.GetBool("is_admin") {
		c.JSON(http.StatusForbidden, models.Response{
			Success: false,
			Message: "权限不足",
		})
		return
	}

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "用户ID无效",
		})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "请求参数错误",
		})
		return
	}

	// 更新用户信息
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Response{
				Success: false,
				Message: "密码加密失败",
			})
			return
		}
		_, err = s.db.Exec("UPDATE users SET name = ?, password = ?, updated_at = ? WHERE id = ?",
			req.Name, string(hashedPassword), time.Now(), userID)
	} else {
		_, err = s.db.Exec("UPDATE users SET name = ?, updated_at = ? WHERE id = ?",
			req.Name, time.Now(), userID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "更新用户失败",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "用户更新成功",
	})
}

// handleDeleteUser 删除用户（管理员）
func (s *Server) handleDeleteUser(c *gin.Context) {
	if !c.GetBool("is_admin") {
		c.JSON(http.StatusForbidden, models.Response{
			Success: false,
			Message: "权限不足",
		})
		return
	}

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "用户ID无效",
		})
		return
	}

	_, err = s.db.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "删除用户失败",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "用户删除成功",
	})
}

// handleAssignMailboxesToUser 为用户分配邮箱
func (s *Server) handleAssignMailboxesToUser(c *gin.Context) {
	if !c.GetBool("is_admin") {
		c.JSON(http.StatusForbidden, models.Response{
			Success: false,
			Message: "权限不足",
		})
		return
	}

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "用户ID无效",
		})
		return
	}

	var req models.GenerateMailboxesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	// 验证域名是否存在且活跃
	var domainName string
	err = s.db.QueryRow("SELECT name FROM domains WHERE id = ? AND is_active = 1", req.DomainID).Scan(&domainName)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, models.Response{
				Success: false,
				Message: "域名不存在或已停用",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.Response{
				Success: false,
				Message: "验证域名失败",
			})
		}
		return
	}

	// 验证用户是否存在
	var userEmail string
	err = s.db.QueryRow("SELECT email FROM users WHERE id = ?", userID).Scan(&userEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, models.Response{
				Success: false,
				Message: "用户不存在",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.Response{
				Success: false,
				Message: "验证用户失败",
			})
		}
		return
	}

	// 生成邮箱（分配给指定用户）
	var generatedAccounts []map[string]interface{}
	for i := 0; i < req.Count; i++ {
		var email string

		// 如果前缀是完整的邮箱名（不包含@），直接使用
		if req.Prefix != "" && !strings.Contains(req.Prefix, "@") && len(req.Prefix) > 5 {
			// 检查前缀是否看起来像完整的邮箱名（包含数字和字母）
			hasDigit := false
			hasLetter := false
			for _, r := range req.Prefix {
				if r >= '0' && r <= '9' {
					hasDigit = true
				}
				if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
					hasLetter = true
				}
			}

			if hasDigit && hasLetter {
				email = req.Prefix + "@" + domainName
			} else {
				email = s.generateRandomEmail(req.Prefix, domainName)
			}
		} else {
			email = s.generateRandomEmail(req.Prefix, domainName)
		}

		// 检查邮箱是否已存在
		var count int
		err := s.db.QueryRow("SELECT COUNT(*) FROM mailboxes WHERE email = ?", email).Scan(&count)
		if err != nil || count > 0 {
			// 如果已存在，重新生成
			i--
			continue
		}

		// 生成随机密码（用于显示，实际邮箱使用用户的登录凭据）
		password := s.generateRandomPassword()

		// 插入邮箱（分配给指定用户）
		_, err = s.db.Exec(`
			INSERT INTO mailboxes (user_id, email, domain_id, password, is_active, is_current, created_at, updated_at)
			VALUES (?, ?, ?, ?, 1, 0, ?, ?)
		`, userID, email, req.DomainID, password, time.Now(), time.Now())

		if err != nil {
			log.Printf("创建邮箱失败: %v", err)
			i--
			continue
		}

		generatedAccounts = append(generatedAccounts, map[string]interface{}{
			"email":    email,
			"password": password,
		})
	}

	if len(generatedAccounts) == 0 {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "邮箱生成失败",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: fmt.Sprintf("成功为用户 %s 分配了 %d 个邮箱", userEmail, len(generatedAccounts)),
		Data: map[string]interface{}{
			"accounts": generatedAccounts,
			"count":    len(generatedAccounts),
		},
	})
}

// handleAssignDomainsToUser 为用户分配多个域名
func (s *Server) handleAssignDomainsToUser(c *gin.Context) {
	if !c.GetBool("is_admin") {
		c.JSON(http.StatusForbidden, models.Response{
			Success: false,
			Message: "权限不足",
		})
		return
	}

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "用户ID无效",
		})
		return
	}

	var req struct {
		DomainIDs []int `json:"domain_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	if len(req.DomainIDs) == 0 {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "请至少选择一个域名",
		})
		return
	}

	// 验证用户是否存在
	var userEmail string
	err = s.db.QueryRow("SELECT email FROM users WHERE id = ?", userID).Scan(&userEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, models.Response{
				Success: false,
				Message: "用户不存在",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.Response{
				Success: false,
				Message: "验证用户失败",
			})
		}
		return
	}

	// 开始事务
	tx, err := s.db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "开始事务失败",
		})
		return
	}
	defer tx.Rollback()

	var assignedDomains []string
	var failedDomains []string

	// 逐个处理域名分配
	for _, domainID := range req.DomainIDs {
		// 检查域名是否存在且未分配给其他用户
		var domainName string
		var currentUserID *int
		err := tx.QueryRow("SELECT name, user_id FROM domains WHERE id = ? AND is_active = 1", domainID).Scan(&domainName, &currentUserID)
		if err != nil {
			if err == sql.ErrNoRows {
				failedDomains = append(failedDomains, fmt.Sprintf("域名ID %d 不存在或已停用", domainID))
			} else {
				failedDomains = append(failedDomains, fmt.Sprintf("查询域名ID %d 失败", domainID))
			}
			continue
		}

		// 检查域名是否已分配给其他用户
		if currentUserID != nil {
			failedDomains = append(failedDomains, fmt.Sprintf("域名 %s 已分配给其他用户", domainName))
			continue
		}

		// 分配域名给用户
		_, err = tx.Exec("UPDATE domains SET user_id = ?, updated_at = ? WHERE id = ?", userID, time.Now(), domainID)
		if err != nil {
			failedDomains = append(failedDomains, fmt.Sprintf("分配域名 %s 失败", domainName))
			continue
		}

		assignedDomains = append(assignedDomains, domainName)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "提交事务失败",
		})
		return
	}

	// 构建响应消息
	var message string
	if len(assignedDomains) > 0 && len(failedDomains) == 0 {
		message = fmt.Sprintf("成功为用户 %s 分配了 %d 个域名", userEmail, len(assignedDomains))
	} else if len(assignedDomains) > 0 && len(failedDomains) > 0 {
		message = fmt.Sprintf("为用户 %s 分配了 %d 个域名，%d 个失败", userEmail, len(assignedDomains), len(failedDomains))
	} else {
		message = "所有域名分配都失败了"
	}

	c.JSON(http.StatusOK, models.Response{
		Success: len(assignedDomains) > 0,
		Message: message,
		Data: map[string]interface{}{
			"assigned_domains": assignedDomains,
			"failed_domains":   failedDomains,
			"user":             userEmail,
			"assigned_count":   len(assignedDomains),
			"failed_count":     len(failedDomains),
		},
	})
}

// handleReclaimDomainsFromUser 从用户回收域名
func (s *Server) handleReclaimDomainsFromUser(c *gin.Context) {
	if !c.GetBool("is_admin") {
		c.JSON(http.StatusForbidden, models.Response{
			Success: false,
			Message: "权限不足",
		})
		return
	}

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "用户ID无效",
		})
		return
	}

	var req struct {
		DomainIDs []int `json:"domain_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "请求参数错误",
		})
		return
	}

	if len(req.DomainIDs) == 0 {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "请选择要回收的域名",
		})
		return
	}

	// 获取用户邮箱
	var userEmail string
	err = s.db.QueryRow("SELECT email FROM users WHERE id = ?", userID).Scan(&userEmail)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "用户不存在",
		})
		return
	}

	// 开始事务
	tx, err := s.db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "开始事务失败",
		})
		return
	}
	defer tx.Rollback()

	var reclaimedDomains []string
	var failedDomains []string

	// 回收域名（将user_id设置为NULL并禁用域名）
	for _, domainID := range req.DomainIDs {
		// 验证域名是否属于该用户
		var currentUserID *int
		var domainName string
		err := tx.QueryRow("SELECT user_id, name FROM domains WHERE id = ?", domainID).Scan(&currentUserID, &domainName)
		if err != nil {
			failedDomains = append(failedDomains, fmt.Sprintf("域名ID %d 不存在", domainID))
			continue
		}

		if currentUserID == nil || *currentUserID != userID {
			failedDomains = append(failedDomains, fmt.Sprintf("域名 %s 不属于该用户", domainName))
			continue
		}

		// 回收域名（设置为NULL但保持活跃，可以重新分配）
		_, err = tx.Exec("UPDATE domains SET user_id = NULL, updated_at = ? WHERE id = ?", time.Now(), domainID)
		if err != nil {
			failedDomains = append(failedDomains, fmt.Sprintf("回收域名 %s 失败", domainName))
			continue
		}

		reclaimedDomains = append(reclaimedDomains, domainName)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "提交事务失败",
		})
		return
	}

	// 构建响应消息
	var message string
	if len(reclaimedDomains) > 0 && len(failedDomains) == 0 {
		message = fmt.Sprintf("成功从用户 %s 回收了 %d 个域名", userEmail, len(reclaimedDomains))
	} else if len(reclaimedDomains) > 0 && len(failedDomains) > 0 {
		message = fmt.Sprintf("从用户 %s 回收了 %d 个域名，%d 个失败", userEmail, len(reclaimedDomains), len(failedDomains))
	} else {
		message = "所有域名回收都失败了"
	}

	c.JSON(http.StatusOK, models.Response{
		Success: len(reclaimedDomains) > 0,
		Message: message,
		Data: map[string]interface{}{
			"reclaimed_domains": reclaimedDomains,
			"failed_domains":    failedDomains,
			"user":              userEmail,
			"reclaimed_count":   len(reclaimedDomains),
			"failed_count":      len(failedDomains),
		},
	})
}

// handleGetDomains 获取域名列表（管理员）
func (s *Server) handleGetDomains(c *gin.Context) {
	if !c.GetBool("is_admin") {
		c.JSON(http.StatusForbidden, models.Response{
			Success: false,
			Message: "权限不足",
		})
		return
	}

	rows, err := s.db.Query("SELECT id, name, user_id, is_active, dns_verified, mx_record, last_verified, created_at, updated_at FROM domains ORDER BY created_at DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "查询域名失败",
		})
		return
	}
	defer rows.Close()

	var domains []models.Domain
	for rows.Next() {
		var domain models.Domain
		err := rows.Scan(&domain.ID, &domain.Name, &domain.UserID, &domain.IsActive, &domain.DNSVerified,
			&domain.MXRecord, &domain.LastVerified, &domain.CreatedAt, &domain.UpdatedAt)
		if err != nil {
			continue
		}
		domains = append(domains, domain)
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    domains,
	})
}

// handleGetUserDomains 获取用户自己的域名列表
func (s *Server) handleGetUserDomains(c *gin.Context) {
	userID := c.GetInt("user_id")
	isAdmin := c.GetBool("is_admin")

	var rows *sql.Rows
	var err error

	if isAdmin {
		// 管理员可以看到所有活跃的域名
		rows, err = s.db.Query(`
			SELECT id, name, user_id, is_active, dns_verified, mx_record, last_verified, created_at, updated_at
			FROM domains
			WHERE is_active = 1
			ORDER BY created_at DESC
		`)
	} else {
		// 普通用户只能看到自己拥有的域名
		rows, err = s.db.Query(`
			SELECT id, name, user_id, is_active, dns_verified, mx_record, last_verified, created_at, updated_at
			FROM domains
			WHERE user_id = ? AND is_active = 1
			ORDER BY created_at DESC
		`, userID)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "查询域名失败",
		})
		return
	}
	defer rows.Close()

	var domains []models.Domain
	for rows.Next() {
		var domain models.Domain
		err := rows.Scan(&domain.ID, &domain.Name, &domain.UserID, &domain.IsActive, &domain.DNSVerified,
			&domain.MXRecord, &domain.LastVerified, &domain.CreatedAt, &domain.UpdatedAt)
		if err != nil {
			continue
		}
		domains = append(domains, domain)
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data:    domains,
	})
}

// handleCreateDomain 创建域名（管理员）
func (s *Server) handleCreateDomain(c *gin.Context) {
	if !c.GetBool("is_admin") {
		c.JSON(http.StatusForbidden, models.Response{
			Success: false,
			Message: "权限不足",
		})
		return
	}

	var req models.CreateDomainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "请求参数错误",
		})
		return
	}

	// 验证域名格式
	if err := validateDomainName(req.Name); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "域名格式错误: " + err.Error(),
		})
		return
	}

	now := time.Now()
	_, err := s.db.Exec(`
		INSERT INTO domains (name, is_active, dns_verified, mx_record, last_verified, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, req.Name, true, false, "", now, now, now)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "创建域名失败",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "域名创建成功",
	})
}

// handleDeleteDomain 删除域名（管理员）
func (s *Server) handleDeleteDomain(c *gin.Context) {
	if !c.GetBool("is_admin") {
		c.JSON(http.StatusForbidden, models.Response{
			Success: false,
			Message: "权限不足",
		})
		return
	}

	domainID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "域名ID无效",
		})
		return
	}

	_, err = s.db.Exec("DELETE FROM domains WHERE id = ?", domainID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "删除域名失败",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "域名删除成功",
	})
}

// decodeIfBase64 检测并解码Base64内容，支持MIME多部分邮件和MIME编码头部
func decodeIfBase64(content string) string {
	// 首先检查是否是MIME编码的头部 (=?charset?encoding?data?=)
	if strings.Contains(content, "=?") && strings.Contains(content, "?=") {
		decoded := decodeMIMEHeader(content)
		if decoded != content {
			return decoded
		}
	}

	// 检查是否是MIME多部分邮件
	if strings.Contains(content, "This is a multi-part message in MIME format") {
		log.Printf("检测到MIME多部分邮件，开始解析")
		result := parseMIMEMultipartForWeb(content)
		log.Printf("MIME解析结果: %s", result)
		return result
	}

	// 去除换行符
	cleanContent := strings.ReplaceAll(content, "\n", "")
	cleanContent = strings.ReplaceAll(cleanContent, "\r", "")
	cleanContent = strings.TrimSpace(cleanContent)

	// 检查是否看起来像Base64（只包含Base64字符且长度合理）
	base64Regex := regexp.MustCompile(`^[A-Za-z0-9+/]*={0,2}$`)
	if len(cleanContent) > 10 && len(cleanContent)%4 == 0 && base64Regex.MatchString(cleanContent) {
		// 尝试解码
		if decoded, err := base64.StdEncoding.DecodeString(cleanContent); err == nil {
			decodedStr := string(decoded)
			// 检查解码后的内容是否包含可打印字符（可能是文本）
			if isPrintableText(decodedStr) {
				return decodedStr
			}
		}
	}

	return content
}

// isPrintableText 检查字符串是否主要包含可打印字符
func isPrintableText(s string) bool {
	if len(s) == 0 {
		return false
	}

	printableCount := 0
	for _, r := range s {
		if r >= 32 && r <= 126 || r >= 0x4e00 && r <= 0x9fff || r == '\n' || r == '\r' || r == '\t' {
			printableCount++
		}
	}

	// 如果80%以上是可打印字符，认为是文本
	return float64(printableCount)/float64(len([]rune(s))) > 0.8
}

// parseMIMEMultipartForWeb 为Web界面解析MIME多部分邮件
func parseMIMEMultipartForWeb(body string) string {
	// 查找boundary - 先尝试从Content-Type头部查找
	var boundary string
	lines := strings.Split(body, "\n")
	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), "boundary=") {
			parts := strings.Split(line, "boundary=")
			if len(parts) > 1 {
				boundary = strings.Trim(strings.Split(parts[1], ";")[0], "\"")
				break
			}
		}
	}

	// 如果没有找到boundary=，尝试从MIME分隔符中提取
	if boundary == "" {
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "------") && !strings.HasSuffix(line, "--") {
				// 提取boundary（去掉前面的--）
				if len(line) > 6 {
					boundary = line[6:]
					break
				}
			}
		}
	}

	log.Printf("Web MIME解析 - 找到boundary: %s", boundary)
	if boundary == "" {
		log.Printf("Web MIME解析 - 未找到boundary，返回原始内容")
		return body
	}

	// 分割各个部分
	parts := strings.Split(body, "--"+boundary)
	var textParts []string
	var htmlParts []string

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" || part == "--" || strings.HasPrefix(part, "--") {
			continue
		}

		// 分离头部和内容
		partLines := strings.Split(part, "\n")
		var inHeaders = true
		var contentLines []string
		var contentType string
		var contentTransferEncoding string
		var charset string

		for _, line := range partLines {
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
			log.Printf("Web MIME解析 - 尝试解码base64 (charset: %s): %s", charset, cleanContent)
			if decoded, err := base64.StdEncoding.DecodeString(cleanContent); err == nil {
				// 先进行base64解码，然后进行字符编码转换
				content = convertToUTF8(decoded, charset)
				log.Printf("Web MIME解析 - base64解码并转换编码成功: %s", content)
			} else {
				log.Printf("Web MIME解析 - base64解码失败: %v", err)
			}
		} else if contentTransferEncoding == "quoted-printable" {
			// 处理quoted-printable编码
			content = decodeQuotedPrintableHeader(content)
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

// decodeMIMEHeader 解码MIME编码的邮件头部 (=?charset?encoding?data?=)
func decodeMIMEHeader(header string) string {
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
			decoded = decodeQuotedPrintableHeader(encodedText)
			decoded = convertToUTF8([]byte(decoded), charset)
		default:
			decoded = encodedText // 未知编码，保持原样
		}

		result = strings.Replace(result, fullMatch, decoded, 1)
	}

	return result
}

// decodeQuotedPrintableHeader 解码quoted-printable编码的头部
func decodeQuotedPrintableHeader(s string) string {
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

// parseMIMEContentForAPI 为API解析MIME内容，返回文本和HTML部分
func parseMIMEContentForAPI(body string) ([]string, []string) {
	var textParts []string
	var htmlParts []string

	log.Printf("Web API - 开始解析MIME内容，长度: %d", len(body))
	previewLen := 200
	if len(body) < previewLen {
		previewLen = len(body)
	}
	log.Printf("Web API - 邮件内容前%d字符: %s", previewLen, body[:previewLen])

	// 查找boundary
	var boundary string
	lines := strings.Split(body, "\n")
	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), "boundary=") {
			parts := strings.Split(line, "boundary=")
			if len(parts) > 1 {
				boundary = strings.Trim(strings.Split(parts[1], ";")[0], "\"")
				log.Printf("Web API - 从头部找到boundary: %s", boundary)
				break
			}
		}
	}

	// 如果在邮件头部没找到boundary，尝试从第一行提取
	if boundary == "" {
		firstLine := strings.TrimSpace(lines[0])
		if strings.HasPrefix(firstLine, "--") {
			// 提取boundary（去掉开头的--）
			boundary = strings.TrimPrefix(firstLine, "--")
			log.Printf("Web API - 从第一行提取boundary: %s", boundary)
		}
	}

	if boundary == "" {
		log.Printf("Web API - 未找到boundary，尝试直接解析")
		// 如果没有boundary，尝试直接解析单个部分
		if decoded := decodeIfBase64(body); decoded != body {
			textParts = append(textParts, decoded)
		}
		return textParts, htmlParts
	}

	log.Printf("Web API - 找到boundary: %s", boundary)

	// 按boundary分割内容
	parts := strings.Split(body, "--"+boundary)

	for i, part := range parts {
		if i == 0 || strings.TrimSpace(part) == "" || strings.TrimSpace(part) == "--" {
			continue
		}

		// 分离头部和内容
		headerEndIndex := strings.Index(part, "\n\n")
		if headerEndIndex == -1 {
			headerEndIndex = strings.Index(part, "\r\n\r\n")
		}

		if headerEndIndex == -1 {
			continue
		}

		headers := part[:headerEndIndex]
		content := strings.TrimSpace(part[headerEndIndex+2:])

		// 解析头部信息
		var contentType, charset, contentTransferEncoding string
		headerLines := strings.Split(headers, "\n")

		for _, line := range headerLines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(strings.ToLower(line), "content-type:") {
				contentType = strings.TrimSpace(line[13:])
				charset = extractCharsetFromContentType(contentType)
				contentType = strings.ToLower(contentType)
			} else if strings.HasPrefix(strings.ToLower(line), "content-transfer-encoding:") {
				contentTransferEncoding = strings.TrimSpace(strings.ToLower(line[26:]))
			}
		}

		log.Printf("Web API - 处理部分 %d: contentType=%s, charset=%s, encoding=%s", i, contentType, charset, contentTransferEncoding)

		// 根据编码方式解码内容
		if contentTransferEncoding == "base64" {
			cleanContent := strings.ReplaceAll(content, "\n", "")
			cleanContent = strings.ReplaceAll(cleanContent, "\r", "")
			cleanContent = strings.TrimSpace(cleanContent)
			if decoded, err := base64.StdEncoding.DecodeString(cleanContent); err == nil {
				content = convertToUTF8(decoded, charset)
				log.Printf("Web API - base64解码成功: %s", content)
			} else {
				log.Printf("Web API - base64解码失败: %v", err)
			}
		} else if contentTransferEncoding == "quoted-printable" {
			content = decodeQuotedPrintableHeader(content)
			content = convertToUTF8([]byte(content), charset)
		} else if charset != "" && charset != "utf-8" {
			content = convertToUTF8([]byte(content), charset)
		}

		// 根据内容类型分类
		if strings.Contains(contentType, "text/plain") {
			textParts = append(textParts, content)
		} else if strings.Contains(contentType, "text/html") {
			htmlParts = append(htmlParts, content)
		}
	}

	log.Printf("Web API - MIME解析完成，文本部分: %d个，HTML部分: %d个", len(textParts), len(htmlParts))
	return textParts, htmlParts
}

// handleParseMIME 处理MIME解析API请求
func (s *Server) handleParseMIME(c *gin.Context) {
	var req struct {
		Body string `json:"body"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "请求参数错误",
		})
		return
	}

	textParts, htmlParts := parseMIMEContentForAPI(req.Body)

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Data: map[string]interface{}{
			"textParts": textParts,
			"htmlParts": htmlParts,
		},
	})
}

// handleFixDomainOwnership 修复域名归属问题
func (s *Server) handleFixDomainOwnership(c *gin.Context) {
	if !c.GetBool("is_admin") {
		c.JSON(http.StatusForbidden, models.Response{
			Success: false,
			Message: "权限不足",
		})
		return
	}

	// 将jbjj.site设置为公共域名
	_, err := s.db.Exec("UPDATE domains SET user_id = NULL WHERE name = 'jbjj.site'")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "修复jbjj.site失败",
		})
		return
	}

	// 将回收的域名设置为不活跃
	_, err = s.db.Exec("UPDATE domains SET is_active = 0 WHERE user_id IS NULL AND name != 'jbjj.site'")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "修复回收域名失败",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "域名归属修复成功",
	})
}
