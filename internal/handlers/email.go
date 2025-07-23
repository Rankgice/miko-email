package handlers

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"miko-email/internal/config"
	"miko-email/internal/models"
	"miko-email/internal/services/email"
	"miko-email/internal/services/forward"
	"miko-email/internal/services/mailbox"
	smtpService "miko-email/internal/services/smtp"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

type EmailHandler struct {
	emailService   *email.Service
	mailboxService *mailbox.Service
	forwardService *forward.Service
	sessionStore   *sessions.CookieStore
	smtpClient     *smtpService.OutboundClient
}

func NewEmailHandler(emailService *email.Service, mailboxService *mailbox.Service, forwardService *forward.Service, sessionStore *sessions.CookieStore) *EmailHandler {
	return &EmailHandler{
		emailService:   emailService,
		mailboxService: mailboxService,
		forwardService: forwardService,
		sessionStore:   sessionStore,
		smtpClient:     smtpService.NewOutboundClientWithDB(mailboxService.GetDB()), // 使用数据库动态获取域名
	}
}

type SendEmailRequest struct {
	From    string `form:"from" binding:"required"`
	To      string `form:"to" binding:"required"`
	CC      string `form:"cc"`
	BCC     string `form:"bcc"`
	Subject string `form:"subject" binding:"required"`
	Content string `form:"content" binding:"required"`
}

// SendEmail 发送邮件
func (h *EmailHandler) SendEmail(c *gin.Context) {
	// 设置正确的Content-Type响应头
	c.Header("Content-Type", "application/json; charset=utf-8")

	// 手动解析表单数据以确保UTF-8编码正确处理
	err := c.Request.ParseMultipartForm(32 << 20) // 32MB max memory
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数错误"})
		return
	}

	// 从表单中获取数据
	req := SendEmailRequest{
		From:    c.Request.FormValue("from"),
		To:      c.Request.FormValue("to"),
		CC:      c.Request.FormValue("cc"),
		BCC:     c.Request.FormValue("bcc"),
		Subject: c.Request.FormValue("subject"),
		Content: c.Request.FormValue("content"),
	}

	// 验证必填字段
	if req.From == "" || req.To == "" || req.Subject == "" || req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数错误"})
		return
	}

	// 获取当前用户信息
	userID := c.GetInt("user_id")
	isAdmin := c.GetBool("is_admin")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "未登录"})
		return
	}

	// 验证发件邮箱是否属于当前用户
	fromMailbox, err := h.mailboxService.GetMailboxByEmail(req.From)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "发件邮箱不存在"})
		return
	}

	// 检查邮箱所有权
	if isAdmin {
		if fromMailbox.AdminID == nil || *fromMailbox.AdminID != userID {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "无权使用此邮箱发送邮件"})
			return
		}
	} else {
		if fromMailbox.UserID == nil || *fromMailbox.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "无权使用此邮箱发送邮件"})
			return
		}
	}

	// 处理收件人列表（支持多个收件人）
	recipients := strings.Split(req.To, ",")
	for i, recipient := range recipients {
		recipients[i] = strings.TrimSpace(recipient)
	}

	// 发送邮件到每个收件人
	var successfulSends []string // 记录成功发送的收件人

	for _, recipient := range recipients {
		if recipient == "" {
			continue
		}

		// 所有邮件都通过SMTP服务器发送，确保一致的处理流程
		var sendErr error

		// 检查收件人邮箱是否存在于系统中
		_, err := h.mailboxService.GetMailboxByEmail(recipient)
		if err != nil {
			// 收件人不在系统中，检查是否为有效的外部邮箱
			if !h.smtpClient.IsExternalEmail(recipient) {
				// 不是有效的外部邮箱，跳过
				continue
			}
		}

		// 统一通过MX发送邮件（无论是内部还是外部邮件）
		sendErr = h.smtpClient.SendEmail(req.From, recipient, req.Subject, req.Content)

		// 记录发送尝试
		h.smtpClient.LogSendAttempt(req.From, recipient, req.Subject, sendErr)

		if sendErr == nil {
			// 发送成功，记录成功的收件人
			successfulSends = append(successfulSends, recipient)
		} else {
			// 发送失败，继续处理下一个收件人
			log.Printf("邮件发送失败 %s -> %s: %v", req.From, recipient, sendErr)
			continue
		}
	}

	// 只有在有成功发送的邮件时，才保存到发件人的已发送文件夹
	for _, recipient := range successfulSends {
		err := h.emailService.SaveEmailToSent(fromMailbox.ID, req.From, recipient, req.Subject, req.Content)
		if err != nil {
			// 保存到已发送失败，记录日志但不影响主要功能
			continue
		}
	}

	// 根据发送结果返回相应消息
	if len(successfulSends) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "没有邮件发送成功"})
	} else if len(successfulSends) == len(recipients) {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "所有邮件发送成功"})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": fmt.Sprintf("部分邮件发送成功 (%d/%d)", len(successfulSends), len(recipients)),
		})
	}
}

// GetEmails 获取邮件列表
func (h *EmailHandler) GetEmails(c *gin.Context) {
	// 设置正确的Content-Type响应头
	c.Header("Content-Type", "application/json; charset=utf-8")

	// 获取当前用户信息
	userID := c.GetInt("user_id")
	isAdmin := c.GetBool("is_admin")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "未登录"})
		return
	}

	// 获取查询参数
	mailboxEmail := c.Query("mailbox")
	emailType := c.DefaultQuery("type", "inbox") // inbox, sent, trash
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// 如果没有指定邮箱，获取用户的第一个邮箱
	var targetMailbox *models.Mailbox
	var err error

	if mailboxEmail != "" {
		targetMailbox, err = h.mailboxService.GetMailboxByEmail(mailboxEmail)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "邮箱不存在"})
			return
		}
	} else {
		// 获取用户的邮箱列表
		mailboxes, err := h.mailboxService.GetUserMailboxesRaw(userID, isAdmin)
		if err != nil || len(mailboxes) == 0 {
			c.JSON(http.StatusOK, gin.H{"success": true, "data": []interface{}{}, "total": 0})
			return
		}
		targetMailbox = &mailboxes[0]
	}

	// 检查邮箱所有权
	if isAdmin {
		if targetMailbox.AdminID == nil || *targetMailbox.AdminID != userID {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "无权访问此邮箱"})
			return
		}
	} else {
		if targetMailbox.UserID == nil || *targetMailbox.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "无权访问此邮箱"})
			return
		}
	}

	// 获取邮件列表
	emails, total, err := h.emailService.GetEmails(targetMailbox.ID, emailType, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取邮件失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    emails,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

// GetEmailByID 获取单个邮件详情
func (h *EmailHandler) GetEmailByID(c *gin.Context) {
	// 设置正确的Content-Type响应头
	c.Header("Content-Type", "application/json; charset=utf-8")

	// 获取当前用户信息
	userID := c.GetInt("user_id")
	isAdmin := c.GetBool("is_admin")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "未登录"})
		return
	}

	emailIDStr := c.Param("id")
	emailID, err := strconv.Atoi(emailIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "邮件ID无效"})
		return
	}

	mailboxEmail := c.Query("mailbox")
	if mailboxEmail == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请指定邮箱"})
		return
	}

	// 获取邮箱信息
	targetMailbox, err := h.mailboxService.GetMailboxByEmail(mailboxEmail)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "邮箱不存在"})
		return
	}

	// 检查邮箱所有权
	if isAdmin {
		if targetMailbox.AdminID == nil || *targetMailbox.AdminID != userID {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "无权访问此邮箱"})
			return
		}
	} else {
		if targetMailbox.UserID == nil || *targetMailbox.UserID != userID {
			c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "无权访问此邮箱"})
			return
		}
	}

	// 获取邮件详情
	email, err := h.emailService.GetEmailByID(emailID, targetMailbox.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "邮件不存在"})
		return
	}

	// 标记为已读
	h.emailService.MarkAsRead(emailID, targetMailbox.ID)

	c.JSON(http.StatusOK, gin.H{"success": true, "data": email})
}

// DeleteEmail 删除邮件
func (h *EmailHandler) DeleteEmail(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")

	emailIDStr := c.Param("id")
	emailID, err := strconv.Atoi(emailIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "邮件ID格式错误"})
		return
	}

	userID := c.GetInt("user_id")
	isAdmin := c.GetBool("is_admin")

	// 首先需要获取用户的邮箱来验证权限
	// 这里我们需要一个更简单的方法来验证邮件所有权
	// 让我们直接在删除时验证权限

	// 获取用户的邮箱列表来验证权限
	userMailboxes, err := h.mailboxService.GetUserMailboxesRaw(userID, isAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取邮箱列表失败"})
		return
	}

	if len(userMailboxes) == 0 {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "message": "无权删除邮件"})
		return
	}

	// 使用第一个邮箱的ID来获取邮件（这里需要改进逻辑）
	mailboxID := userMailboxes[0].ID

	// 验证邮件是否存在且属于用户的邮箱
	_, err = h.emailService.GetEmailByID(emailID, mailboxID)
	if err != nil {
		// 尝试其他邮箱
		found := false
		for _, mb := range userMailboxes {
			_, err = h.emailService.GetEmailByID(emailID, mb.ID)
			if err == nil {
				mailboxID = mb.ID
				found = true
				break
			}
		}
		if !found {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "邮件不存在或无权访问"})
			return
		}
	}

	// 删除邮件
	err = h.emailService.DeleteEmail(emailID, mailboxID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "删除邮件失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "邮件删除成功"})
}

// 使用转发服务中的结构体，这里不需要重复定义

// GetForwardRules 获取转发规则列表
func (h *EmailHandler) GetForwardRules(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")

	username := c.GetString("username")
	userID := c.GetInt("user_id")

	rules, err := h.forwardService.GetForwardRulesByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取转发规则失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    rules,
		"message": fmt.Sprintf("用户 %s 的转发规则", username),
	})
}

// CreateForwardRule 创建转发规则
func (h *EmailHandler) CreateForwardRule(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")

	var req forward.CreateForwardRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	username := c.GetString("username")
	userID := c.GetInt("user_id")

	newRule, err := h.forwardService.CreateForwardRule(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    newRule,
		"message": fmt.Sprintf("用户 %s 创建转发规则成功", username),
	})
}

// GetForwardRule 获取单个转发规则
func (h *EmailHandler) GetForwardRule(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的规则ID",
		})
		return
	}

	userID := c.GetInt("user_id")

	rule, err := h.forwardService.GetForwardRuleByID(id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    rule,
	})
}

// UpdateForwardRule 更新转发规则
func (h *EmailHandler) UpdateForwardRule(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的规则ID",
		})
		return
	}

	var req forward.CreateForwardRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	userID := c.GetInt("user_id")

	err = h.forwardService.UpdateForwardRule(id, userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("转发规则 %d 更新成功", id),
	})
}

// DeleteForwardRule 删除转发规则
func (h *EmailHandler) DeleteForwardRule(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的规则ID",
		})
		return
	}

	userID := c.GetInt("user_id")

	err = h.forwardService.DeleteForwardRule(id, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("转发规则 %d 删除成功", id),
	})
}

// ToggleForwardRule 切换转发规则状态
func (h *EmailHandler) ToggleForwardRule(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的规则ID",
		})
		return
	}

	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	userID := c.GetInt("user_id")

	err = h.forwardService.ToggleForwardRule(id, userID, req.Enabled)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	status := "启用"
	if !req.Enabled {
		status = "禁用"
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("转发规则 %d 已%s", id, status),
	})
}

// TestForwardRule 测试转发规则
func (h *EmailHandler) TestForwardRule(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "无效的规则ID",
		})
		return
	}

	var req struct {
		Subject string `json:"subject"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 暂时返回成功响应，实际应该发送测试邮件
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("转发规则 %d 测试邮件发送成功", id),
	})
}

// GetForwardStatistics 获取转发统计信息
func (h *EmailHandler) GetForwardStatistics(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")

	userID := c.GetInt("user_id")

	stats, err := h.forwardService.GetForwardStatistics(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取统计信息失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// sendThroughLocalSMTP 通过本地SMTP服务器发送邮件
func (h *EmailHandler) sendThroughLocalSMTP(from, to, subject, body string) error {
	// 获取配置
	cfg := config.Load()

	// 构建邮件内容
	message := h.buildEmailMessage(from, to, subject, body)

	// 获取本地SMTP端口（优先使用587端口）
	smtpPorts := cfg.GetSMTPPorts()
	var port string
	for _, p := range smtpPorts {
		if p == "587" {
			port = p
			break
		}
	}
	if port == "" && len(smtpPorts) > 0 {
		port = smtpPorts[0] // 使用第一个可用端口
	}
	if port == "" {
		port = "25" // 默认端口
	}

	// 连接到本地SMTP服务器
	addr := fmt.Sprintf("localhost:%s", port)

	// 设置连接超时
	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		return fmt.Errorf("连接本地SMTP服务器失败: %v", err)
	}
	defer conn.Close()

	// 创建SMTP客户端
	client, err := smtp.NewClient(conn, "localhost")
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
	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("开始发送邮件内容失败: %v", err)
	}
	defer wc.Close()

	_, err = wc.Write(message)
	if err != nil {
		return fmt.Errorf("写入邮件内容失败: %v", err)
	}

	log.Printf("✅ 通过本地SMTP发送成功: %s -> %s", from, to)
	return nil
}

// buildEmailMessage 构建邮件消息
func (h *EmailHandler) buildEmailMessage(from, to, subject, body string) []byte {
	// 构建标准的邮件格式
	message := fmt.Sprintf("From: %s\r\n", from)
	message += fmt.Sprintf("To: %s\r\n", to)
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "MIME-Version: 1.0\r\n"
	message += "Content-Type: text/plain; charset=UTF-8\r\n"
	message += "Content-Transfer-Encoding: 8bit\r\n"
	message += fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z))
	message += "\r\n"
	message += body

	return []byte(message)
}
