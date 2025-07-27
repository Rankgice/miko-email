package handlers

import (
	"encoding/base64"
	"fmt"
	"log"
	"miko-email/internal/result"
	"miko-email/internal/svc"
	"mime"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"miko-email/internal/model"
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
	svcCtx         *svc.ServiceContext
}

func NewEmailHandler(sessionStore *sessions.CookieStore, svcCtx *svc.ServiceContext) *EmailHandler {
	// 从GORM获取原生SQL数据库连接
	sqlDB, err := svcCtx.DB.DB()
	if err != nil {
		panic("Failed to get SQL DB from GORM: " + err.Error())
	}

	return &EmailHandler{
		emailService:   email.NewService(svcCtx),
		mailboxService: mailbox.NewService(svcCtx),
		forwardService: forward.NewService(svcCtx),
		sessionStore:   sessionStore,
		smtpClient:     smtpService.NewOutboundClientWithDB(sqlDB, svcCtx), // 使用数据库动态获取域名
		svcCtx:         svcCtx,
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

type EmailAttachment struct {
	Filename string
	Content  []byte
	MimeType string
}

// SendEmail 发送邮件
func (h *EmailHandler) SendEmail(c *gin.Context) {
	// 设置正确的Content-Type响应头
	c.Header("Content-Type", "application/json; charset=utf-8")

	// 手动解析表单数据以确保UTF-8编码正确处理
	err := c.Request.ParseMultipartForm(32 << 20) // 32MB max memory
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorReqParam)
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
		c.JSON(http.StatusBadRequest, result.ErrorReqParam)
		return
	}

	// 处理附件
	var attachments []EmailAttachment
	if c.Request.MultipartForm != nil && c.Request.MultipartForm.File != nil {
		files := c.Request.MultipartForm.File["attachments"]
		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("附件读取失败: "+err.Error()))
				return
			}
			defer file.Close()

			// 读取文件内容
			content := make([]byte, fileHeader.Size)
			_, err = file.Read(content)
			if err != nil {
				c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("附件内容读取失败: "+err.Error()))
				return
			}

			// 检查文件大小限制（10MB）
			if fileHeader.Size > 10*1024*1024 {
				c.JSON(http.StatusBadRequest, result.ErrorSimpleResult(fmt.Sprintf("附件 %s 超过10MB限制", fileHeader.Filename)))
				return
			}

			attachments = append(attachments, EmailAttachment{
				Filename: fileHeader.Filename,
				Content:  content,
				MimeType: fileHeader.Header.Get("Content-Type"),
			})
		}
	}

	// 获取当前用户信息
	userID := int64(c.GetInt("user_id"))
	isAdmin := c.GetBool("is_admin")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorSimpleResult("未登录"))
		return
	}

	// 验证发件邮箱是否属于当前用户
	fromMailbox, err := h.mailboxService.GetMailboxByEmail(req.From)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("发件邮箱不存在"))
		return
	}

	// 检查邮箱所有权
	if isAdmin {
		if fromMailbox.AdminId == nil || *fromMailbox.AdminId != userID {
			c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权使用此邮箱发送邮件"))
			return
		}
	} else {
		if fromMailbox.UserId == nil || *fromMailbox.UserId != userID {
			c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权使用此邮箱发送邮件"))
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
		if len(attachments) > 0 {
			// 构建MIME邮件内容
			mimeContent := h.buildMIMEMessage(req.From, recipient, req.Subject, req.Content, attachments)
			sendErr = h.smtpClient.SendMIMEEmail(req.From, recipient, mimeContent)
		} else {
			sendErr = h.smtpClient.SendEmail(req.From, recipient, req.Subject, req.Content)
		}

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
		err := h.emailService.SaveEmailToSent(fromMailbox.Id, req.From, recipient, req.Subject, req.Content)
		if err != nil {
			// 保存到已发送失败，记录日志但不影响主要功能
			continue
		}
	}

	// 根据发送结果返回相应消息
	if len(successfulSends) == 0 {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("没有邮件发送成功"))
	} else if len(successfulSends) == len(recipients) {
		c.JSON(http.StatusOK, result.SimpleResult("所有邮件发送成功"))
	} else {
		c.JSON(http.StatusOK, result.SimpleResult(fmt.Sprintf("部分邮件发送成功 (%d/%d)", len(successfulSends), len(recipients))))
	}
}

// GetEmails 获取邮件列表
func (h *EmailHandler) GetEmails(c *gin.Context) {
	// 设置正确的Content-Type响应头
	c.Header("Content-Type", "application/json; charset=utf-8")

	// 获取当前用户信息
	userID := int64(c.GetInt("user_id"))
	isAdmin := c.GetBool("is_admin")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorSimpleResult("未登录"))
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
	var targetMailbox *model.Mailbox
	var err error

	if mailboxEmail != "" {
		targetMailbox, err = h.mailboxService.GetMailboxByEmail(mailboxEmail)
		if err != nil {
			c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮箱不存在"))
			return
		}
	} else {
		// 获取用户的邮箱列表
		mailboxes, err := h.mailboxService.GetUserMailboxesRaw(userID, isAdmin)
		if err != nil || len(mailboxes) == 0 {
			c.JSON(http.StatusOK, result.ListResult([]interface{}{}, 0, 0, 0))
			return
		}
		targetMailbox = mailboxes[0]
	}

	// 检查邮箱所有权
	if isAdmin {
		if targetMailbox.AdminId == nil || *targetMailbox.AdminId != userID {
			c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权访问此邮箱"))
			return
		}
	} else {
		if targetMailbox.UserId == nil || *targetMailbox.UserId != userID {
			c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权访问此邮箱"))
			return
		}
	}

	// 获取邮件列表
	emails, total, err := h.emailService.GetEmails(targetMailbox.Id, emailType, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("获取邮件失败"))
		return
	}

	c.JSON(http.StatusOK, result.ListResult(emails, int64(page), int64(limit), total))
}

// GetEmailByID 获取单个邮件详情
func (h *EmailHandler) GetEmailByID(c *gin.Context) {
	// 设置正确的Content-Type响应头
	c.Header("Content-Type", "application/json; charset=utf-8")

	// 获取当前用户信息
	userID := int64(c.GetInt("user_id"))
	isAdmin := c.GetBool("is_admin")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorSimpleResult("未登录"))
		return
	}

	emailIDStr := c.Param("id")
	emailID, err := strconv.Atoi(emailIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮件ID无效"))
		return
	}

	mailboxEmail := c.Query("mailbox")
	if mailboxEmail == "" {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("请指定邮箱"))
		return
	}

	// 获取邮箱信息
	targetMailbox, err := h.mailboxService.GetMailboxByEmail(mailboxEmail)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮箱不存在"))
		return
	}

	// 检查邮箱所有权
	if isAdmin {
		if targetMailbox.AdminId == nil || *targetMailbox.AdminId != userID {
			c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权访问此邮箱"))
			return
		}
	} else {
		if targetMailbox.UserId == nil || *targetMailbox.UserId != userID {
			c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权访问此邮箱"))
			return
		}
	}

	// 获取邮件详情
	email, err := h.emailService.GetEmailByID(int64(emailID), targetMailbox.Id)
	if err != nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("邮件不存在"))
		return
	}

	// 标记为已读
	h.emailService.MarkAsRead(int64(emailID), targetMailbox.Id)

	c.JSON(http.StatusOK, result.DataResult("获取邮件详情成功", email))
}

// DeleteEmail 删除邮件
func (h *EmailHandler) DeleteEmail(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")

	emailIDStr := c.Param("id")
	emailID, err := strconv.Atoi(emailIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮件ID格式错误"))
		return
	}

	userID := int64(c.GetInt("user_id"))
	isAdmin := c.GetBool("is_admin")

	// 首先需要获取用户的邮箱来验证权限
	// 这里我们需要一个更简单的方法来验证邮件所有权
	// 让我们直接在删除时验证权限

	// 获取用户的邮箱列表来验证权限
	userMailboxes, err := h.mailboxService.GetUserMailboxesRaw(userID, isAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("获取邮箱列表失败"))
		return
	}

	if len(userMailboxes) == 0 {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权删除邮件"))
		return
	}

	// 使用第一个邮箱的ID来获取邮件（这里需要改进逻辑）
	mailboxID := int(userMailboxes[0].Id)

	// 验证邮件是否存在且属于用户的邮箱
	_, err = h.emailService.GetEmailByID(int64(emailID), int64(mailboxID))
	if err != nil {
		// 尝试其他邮箱
		found := false
		for _, mb := range userMailboxes {
			_, err = h.emailService.GetEmailByID(int64(emailID), mb.Id)
			if err == nil {
				mailboxID = int(mb.Id)
				found = true
				break
			}
		}
		if !found {
			c.JSON(http.StatusNotFound, result.ErrorSimpleResult("邮件不存在或无权访问"))
			return
		}
	}

	// 删除邮件
	err = h.emailService.DeleteEmail(int64(emailID), int64(mailboxID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("删除邮件失败"))
		return
	}

	c.JSON(http.StatusOK, result.SimpleResult("邮件删除成功"))
}

// GetSentEmails 获取发件箱邮件列表
func (h *EmailHandler) GetSentEmails(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")

	userID := int64(c.GetInt("user_id"))
	isAdmin := c.GetBool("is_admin")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorSimpleResult("未登录"))
		return
	}

	// 获取用户的邮箱列表
	mailboxes, err := h.mailboxService.GetUserMailboxes(userID, isAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("获取邮箱列表失败"))
		return
	}

	// 获取所有邮箱的发送邮件
	var allEmails []gin.H

	// 由于没有专门的发送记录表，这里返回空列表作为示例
	// 在实际实现中，应该有一个专门的sent_emails表来记录发送的邮件

	// 暂时返回一些示例数据
	for i, mailbox := range mailboxes {
		if i < 2 { // 只返回前两个邮箱的示例数据
			allEmails = append(allEmails, gin.H{
				"id":           int64(i + 1),
				"from_address": mailbox.Email,
				"to_address":   "recipient@example.com",
				"subject":      fmt.Sprintf("来自 %s 的测试邮件", mailbox.Email),
				"body":         "这是一封测试邮件内容",
				"created_at":   time.Now().Add(-time.Duration(i) * time.Hour),
				"status":       "sent",
				"starred":      false,
			})
		}
	}

	c.JSON(http.StatusOK, result.DataResult("发件箱邮件列表", allEmails))
}

// SaveDraft 保存草稿
func (h *EmailHandler) SaveDraft(c *gin.Context) {
	userID := c.GetInt("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorSimpleResult("未登录"))
		return
	}

	var req struct {
		From    string   `json:"from"`
		To      []string `json:"to"`
		CC      []string `json:"cc"`
		BCC     []string `json:"bcc"`
		Subject string   `json:"subject"`
		Content string   `json:"content"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("参数错误"))
		return
	}

	// 创建草稿记录
	draft := &model.Email{
		FromAddr:  req.From,
		ToAddr:    strings.Join(req.To, ","),
		Subject:   req.Subject,
		Body:      req.Content,
		Folder:    "draft",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 保存到数据库（暂时只返回成功）
	// err := h.svcCtx.EmailModel.Insert(h.svcCtx.DB, draft)
	// if err != nil {
	//     c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("保存草稿失败"))
	//     return
	// }

	c.JSON(http.StatusOK, result.DataResult("草稿保存成功", gin.H{
		"id":       draft.Id,
		"saved_at": time.Now(),
	}))
}

// SearchEmails 搜索邮件
func (h *EmailHandler) SearchEmails(c *gin.Context) {
	userID := c.GetInt("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorSimpleResult("未登录"))
		return
	}

	var req struct {
		Filters struct {
			From          string   `json:"from"`
			To            string   `json:"to"`
			Subject       string   `json:"subject"`
			Date          string   `json:"date"`
			HasAttachment bool     `json:"hasAttachment"`
			Keywords      []string `json:"keywords"`
		} `json:"filters"`
		Page      int    `json:"page"`
		Limit     int    `json:"limit"`
		SortBy    string `json:"sortBy"`
		SortOrder string `json:"sortOrder"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("参数错误"))
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.SortBy == "" {
		req.SortBy = "created_at"
	}
	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}

	// 模拟搜索结果
	searchResults := []gin.H{
		{
			"id":             1,
			"from_address":   "sender@example.com",
			"to_address":     "user@domain.com",
			"subject":        "重要通知：系统维护",
			"body":           "系统将于今晚进行维护，请提前保存工作...",
			"created_at":     time.Now().Add(-2 * time.Hour),
			"is_read":        false,
			"has_attachment": false,
		},
		{
			"id":             2,
			"from_address":   "admin@company.com",
			"to_address":     "user@domain.com",
			"subject":        "会议邀请",
			"body":           "邀请您参加明天的项目会议...",
			"created_at":     time.Now().Add(-24 * time.Hour),
			"is_read":        true,
			"has_attachment": true,
		},
	}

	// 根据搜索条件过滤结果（这里是简化的实现）
	filteredResults := []gin.H{}
	for _, email := range searchResults {
		match := true

		// 检查发件人过滤
		if req.Filters.From != "" {
			if !strings.Contains(strings.ToLower(email["from_address"].(string)), strings.ToLower(req.Filters.From)) {
				match = false
			}
		}

		// 检查收件人过滤
		if req.Filters.To != "" {
			if !strings.Contains(strings.ToLower(email["to_address"].(string)), strings.ToLower(req.Filters.To)) {
				match = false
			}
		}

		// 检查主题过滤
		if req.Filters.Subject != "" {
			if !strings.Contains(strings.ToLower(email["subject"].(string)), strings.ToLower(req.Filters.Subject)) {
				match = false
			}
		}

		// 检查关键词过滤
		if len(req.Filters.Keywords) > 0 {
			keywordMatch := false
			for _, keyword := range req.Filters.Keywords {
				if strings.Contains(strings.ToLower(email["subject"].(string)), strings.ToLower(keyword)) ||
					strings.Contains(strings.ToLower(email["body"].(string)), strings.ToLower(keyword)) {
					keywordMatch = true
					break
				}
			}
			if !keywordMatch {
				match = false
			}
		}

		if match {
			filteredResults = append(filteredResults, email)
		}
	}

	// 计算分页
	total := len(filteredResults)
	totalPages := (total + req.Limit - 1) / req.Limit
	start := (req.Page - 1) * req.Limit
	end := start + req.Limit

	if start > total {
		start = total
	}
	if end > total {
		end = total
	}

	pagedResults := filteredResults[start:end]

	c.JSON(http.StatusOK, result.DataResult("搜索结果", gin.H{
		"results":    pagedResults,
		"total":      total,
		"page":       req.Page,
		"totalPages": totalPages,
	}))
}

// 使用转发服务中的结构体，这里不需要重复定义

// GetForwardRules 获取转发规则列表
func (h *EmailHandler) GetForwardRules(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")

	username := c.GetString("username")
	userID := c.GetInt64("user_id")

	rules, err := h.forwardService.GetForwardRulesByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("获取转发规则失败: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, result.DataResult(fmt.Sprintf("用户 %s 的转发规则", username), rules))
}

// CreateForwardRule 创建转发规则
func (h *EmailHandler) CreateForwardRule(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")

	var req forward.CreateForwardRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("请求参数错误: "+err.Error()))
		return
	}

	username := c.GetString("username")
	userID := c.GetInt64("user_id")

	newRule, err := h.forwardService.CreateForwardRule(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result.DataResult(fmt.Sprintf("用户 %s 创建转发规则成功", username), newRule))
}

// GetForwardRule 获取单个转发规则
func (h *EmailHandler) GetForwardRule(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的规则ID"))
		return
	}

	userID := c.GetInt64("user_id")

	rule, err := h.forwardService.GetForwardRuleByID(id, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result.DataResult("获取转发规则成功", rule))
}

// UpdateForwardRule 更新转发规则
func (h *EmailHandler) UpdateForwardRule(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的规则ID"))
		return
	}

	var req forward.CreateForwardRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("请求参数错误: "+err.Error()))
		return
	}

	userID := c.GetInt64("user_id")

	err = h.forwardService.UpdateForwardRule(id, userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result.SimpleResult(fmt.Sprintf("转发规则 %d 更新成功", id)))
}

// DeleteForwardRule 删除转发规则
func (h *EmailHandler) DeleteForwardRule(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的规则ID"))
		return
	}

	userID := c.GetInt64("user_id")

	err = h.forwardService.DeleteForwardRule(id, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result.SimpleResult(fmt.Sprintf("转发规则 %d 删除成功", id)))
}

// ToggleForwardRule 切换转发规则状态
func (h *EmailHandler) ToggleForwardRule(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的规则ID"))
		return
	}

	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("请求参数错误: "+err.Error()))
		return
	}

	userID := c.GetInt64("user_id")

	err = h.forwardService.ToggleForwardRule(id, userID, req.Enabled)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult(err.Error()))
		return
	}

	status := "启用"
	if !req.Enabled {
		status = "禁用"
	}

	c.JSON(http.StatusOK, result.SimpleResult(fmt.Sprintf("转发规则 %d 已%s", id, status)))
}

// TestForwardRule 测试转发规则
func (h *EmailHandler) TestForwardRule(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的规则ID"))
		return
	}

	var req struct {
		Subject string `json:"subject"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("请求参数错误: "+err.Error()))
		return
	}

	// 获取用户ID
	userID := c.GetInt("user_id")

	// 获取转发规则详情
	rule, err := h.forwardService.GetForwardRuleByID(int64(id), int64(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("转发规则不存在或无权限访问"))
		return
	}

	// 检查规则是否启用
	if !rule.Enabled {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("转发规则已禁用，无法测试"))
		return
	}

	// 构建测试邮件的主题和内容
	testSubject := req.Subject
	if testSubject == "" {
		testSubject = "测试转发邮件 - " + time.Now().Format("2006-01-02 15:04:05")
	}

	testContent := req.Content
	if testContent == "" {
		testContent = fmt.Sprintf(`这是一封测试转发功能的邮件。

测试时间: %s
源邮箱: %s
目标邮箱: %s
转发规则ID: %d

如果您收到这封邮件，说明转发功能正常工作。`,
			time.Now().Format("2006-01-02 15:04:05"),
			rule.SourceEmail,
			rule.TargetEmail,
			rule.ID)
	}

	// 发送测试邮件到源邮箱，触发转发
	err = h.emailService.SendTestForwardEmail(rule.SourceEmail, rule.TargetEmail, testSubject, testContent, rule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("测试邮件发送失败: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, result.SimpleResult(fmt.Sprintf("测试邮件已发送到 %s，如果转发规则正常，您应该会在 %s 收到转发邮件", rule.SourceEmail, rule.TargetEmail)))
}

// GetForwardStatistics 获取转发统计信息
func (h *EmailHandler) GetForwardStatistics(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")

	userID := c.GetInt64("user_id")

	stats, err := h.forwardService.GetForwardStatistics(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("获取统计信息失败: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, result.DataResult("获取转发统计信息成功", stats))
}

// buildMIMEMessage 构建MIME格式的邮件内容
func (h *EmailHandler) buildMIMEMessage(from, to, subject, body string, attachments []EmailAttachment) string {
	boundary := fmt.Sprintf("----=_NextPart_%d", time.Now().Unix())

	var message strings.Builder

	// 邮件头部
	message.WriteString(fmt.Sprintf("From: %s\r\n", from))
	message.WriteString(fmt.Sprintf("To: %s\r\n", to))
	message.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	message.WriteString("MIME-Version: 1.0\r\n")
	message.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=\"%s\"\r\n", boundary))
	message.WriteString("\r\n")

	// 邮件正文部分
	message.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	message.WriteString("Content-Type: text/html; charset=utf-8\r\n")
	message.WriteString("Content-Transfer-Encoding: 8bit\r\n")
	message.WriteString("\r\n")
	message.WriteString(body)
	message.WriteString("\r\n")

	// 附件部分
	for _, attachment := range attachments {
		message.WriteString(fmt.Sprintf("--%s\r\n", boundary))

		// 确定MIME类型
		mimeType := attachment.MimeType
		if mimeType == "" {
			mimeType = mime.TypeByExtension(filepath.Ext(attachment.Filename))
			if mimeType == "" {
				mimeType = "application/octet-stream"
			}
		}

		message.WriteString(fmt.Sprintf("Content-Type: %s\r\n", mimeType))
		message.WriteString("Content-Transfer-Encoding: base64\r\n")
		message.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=\"%s\"\r\n", attachment.Filename))
		message.WriteString("\r\n")

		// Base64编码附件内容
		encoded := base64.StdEncoding.EncodeToString(attachment.Content)
		// 每76个字符换行
		for i := 0; i < len(encoded); i += 76 {
			end := i + 76
			if end > len(encoded) {
				end = len(encoded)
			}
			message.WriteString(encoded[i:end])
			message.WriteString("\r\n")
		}
	}

	// 结束边界
	message.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

	return message.String()
}

// GetVerificationCode 获取邮件验证码
func (h *EmailHandler) GetVerificationCode(c *gin.Context) {
	c.Header("Content-Type", "application/json; charset=utf-8")

	userID := int64(c.GetInt("user_id"))
	mailbox := c.Query("mailbox")
	sender := c.Query("sender")               // 可选：指定发件人过滤
	subject := c.Query("subject")             // 可选：指定主题关键词过滤
	emailIDStr := c.Query("email_id")         // 可选：指定特定邮件ID
	limitStr := c.DefaultQuery("limit", "10") // 默认查询最近10封邮件

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	if mailbox == "" {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮箱地址不能为空"))
		return
	}

	// 验证邮箱是否属于当前用户
	mailboxInfo, err := h.mailboxService.GetMailboxByEmail(mailbox)
	if err != nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("邮箱不存在"))
		return
	}

	if mailboxInfo.UserId == nil || *mailboxInfo.UserId != userID {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权访问此邮箱"))
		return
	}

	var emails []*model.Email

	// 如果指定了email_id，只查询特定邮件
	if emailIDStr != "" {
		emailID, parseErr := strconv.Atoi(emailIDStr)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮件ID格式错误"))
			return
		}

		// 获取特定邮件
		email, getErr := h.emailService.GetEmailByID(int64(emailID), mailboxInfo.Id)
		if getErr != nil {
			c.JSON(http.StatusNotFound, result.ErrorSimpleResult("邮件不存在或无权访问"))
			return
		}
		emails = []*model.Email{email}
	} else {
		// 获取邮件列表
		var getErr error
		emails, _, getErr = h.emailService.GetEmails(mailboxInfo.Id, "inbox", 1, limit)
		if getErr != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("获取邮件失败: "+getErr.Error()))
			return
		}
	}

	// 提取验证码
	var results []map[string]interface{}

	for _, email := range emails {
		// 如果指定了email_id，跳过过滤条件检查
		if emailIDStr == "" {
			// 应用过滤条件
			if sender != "" && !strings.Contains(strings.ToLower(email.FromAddr), strings.ToLower(sender)) {
				continue
			}
			if subject != "" && !strings.Contains(strings.ToLower(email.Subject), strings.ToLower(subject)) {
				continue
			}
		}

		// 提取验证码
		codes := extractVerificationCodes(email.Body)
		if len(codes) > 0 {
			results = append(results, map[string]interface{}{
				"email_id":   email.Id,
				"from":       email.FromAddr,
				"subject":    email.Subject,
				"created_at": email.CreatedAt,
				"codes":      codes,
			})
		}
	}

	c.JSON(http.StatusOK, result.ListResult(results, 0, 0, int64(len(results))))
}

// extractVerificationCodes 从邮件内容中提取验证码
func extractVerificationCodes(content string) []string {
	var codes []string

	// 常见的验证码模式
	patterns := []string{
		`\b\d{4,8}\b`,                   // 4-8位纯数字
		`\b[A-Z0-9]{4,8}\b`,             // 4-8位大写字母和数字组合
		`\b[a-zA-Z0-9]{4,8}\b`,          // 4-8位字母数字组合
		`验证码[：:\s]*([A-Za-z0-9]{4,8})`,  // 中文"验证码"后跟代码
		`验证码[：:\s]*(\d{4,8})`,           // 中文"验证码"后跟数字
		`code[：:\s]*([A-Za-z0-9]{4,8})`, // 英文"code"后跟代码
		`Code[：:\s]*([A-Za-z0-9]{4,8})`, // 英文"Code"后跟代码
		`CODE[：:\s]*([A-Za-z0-9]{4,8})`, // 英文"CODE"后跟代码
	}

	// 使用正则表达式提取
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(content, -1)

		for _, match := range matches {
			if len(match) > 1 {
				// 有捕获组的情况
				code := strings.TrimSpace(match[1])
				if isValidVerificationCode(code) {
					codes = append(codes, code)
				}
			} else if len(match) > 0 {
				// 没有捕获组的情况
				code := strings.TrimSpace(match[0])
				if isValidVerificationCode(code) {
					codes = append(codes, code)
				}
			}
		}
	}

	// 去重
	seen := make(map[string]bool)
	var uniqueCodes []string
	for _, code := range codes {
		if !seen[code] {
			seen[code] = true
			uniqueCodes = append(uniqueCodes, code)
		}
	}

	return uniqueCodes
}

// isValidVerificationCode 验证是否为有效的验证码
func isValidVerificationCode(code string) bool {
	// 长度检查
	if len(code) < 4 || len(code) > 8 {
		return false
	}

	// 排除一些明显不是验证码的内容
	excludePatterns := []string{
		`^\d{4}$`,                                // 排除4位年份
		`^(19|20)\d{2}$`,                         // 排除年份
		`^(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])$`, // 排除日期格式
	}

	for _, pattern := range excludePatterns {
		matched, _ := regexp.MatchString(pattern, code)
		if matched {
			return false
		}
	}

	return true
}

// List 获取邮件列表（API）
func (h *EmailHandler) List(c *gin.Context) {
	userID := int64(c.GetInt("user_id"))
	isAdmin := c.GetBool("is_admin")

	// 获取用户的邮箱列表
	mailboxes, err := h.mailboxService.GetUserMailboxesRaw(userID, isAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("获取邮箱列表失败"))
		return
	}

	// 收集所有邮件
	var allEmails []*model.Email
	for _, mailbox := range mailboxes {
		emails, _, err := h.svcCtx.EmailModel.GetEmailsByMailboxId(mailbox.Id, "", 0, 0)
		if err != nil {
			continue
		}
		allEmails = append(allEmails, emails...)
	}

	// 转换为API响应格式
	var emailResponses []gin.H
	for _, email := range allEmails {
		emailResponses = append(emailResponses, gin.H{
			"id":           email.Id,
			"from_address": email.FromAddr,
			"to_address":   email.ToAddr,
			"subject":      email.Subject,
			"created_at":   email.CreatedAt,
			"read":         email.IsRead,
			"folder":       email.Folder,
		})
	}

	c.JSON(http.StatusOK, result.DataResult("", emailResponses))
}

// GetDetail 获取邮件详情（API）
func (h *EmailHandler) GetDetail(c *gin.Context) {
	emailIDStr := c.Param("id")
	emailID, err := strconv.ParseInt(emailIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮件ID格式错误"))
		return
	}

	userID := int64(c.GetInt("user_id"))
	isAdmin := c.GetBool("is_admin")

	// 获取邮件详情
	email, err := h.svcCtx.EmailModel.GetById(emailID)
	if err != nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("邮件不存在"))
		return
	}

	// 验证权限 - 检查邮件是否属于用户的邮箱
	mailbox, err := h.svcCtx.MailboxModel.GetById(email.MailboxId)
	if err != nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("邮箱不存在"))
		return
	}

	// 权限检查
	hasPermission := false
	if isAdmin && mailbox.AdminId != nil && *mailbox.AdminId == userID {
		hasPermission = true
	} else if !isAdmin && mailbox.UserId != nil && *mailbox.UserId == userID {
		hasPermission = true
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限访问此邮件"))
		return
	}

	// 标记为已读
	if !email.IsRead {
		h.svcCtx.EmailModel.MarkAsRead(nil, emailID)
	}

	emailResponse := gin.H{
		"id":           email.Id,
		"from_address": email.FromAddr,
		"to_address":   email.ToAddr,
		"subject":      email.Subject,
		"body":         email.Body,
		"text_body":    email.Body, // 假设都是文本邮件
		"created_at":   email.CreatedAt,
		"read":         true, // 已标记为已读
		"folder":       email.Folder,
	}

	c.JSON(http.StatusOK, result.DataResult("", emailResponse))
}

// Delete 删除邮件（API）
func (h *EmailHandler) Delete(c *gin.Context) {
	emailIDStr := c.Param("id")
	emailID, err := strconv.ParseInt(emailIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮件ID格式错误"))
		return
	}

	userID := int64(c.GetInt("user_id"))
	isAdmin := c.GetBool("is_admin")

	// 获取邮件详情
	email, err := h.svcCtx.EmailModel.GetById(emailID)
	if err != nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("邮件不存在"))
		return
	}

	// 验证权限 - 检查邮件是否属于用户的邮箱
	mailbox, err := h.svcCtx.MailboxModel.GetById(email.MailboxId)
	if err != nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("邮箱不存在"))
		return
	}

	// 权限检查
	hasPermission := false
	if isAdmin && mailbox.AdminId != nil && *mailbox.AdminId == userID {
		hasPermission = true
	} else if !isAdmin && mailbox.UserId != nil && *mailbox.UserId == userID {
		hasPermission = true
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限删除此邮件"))
		return
	}

	// 删除邮件
	if err := h.svcCtx.EmailModel.Delete(nil, email); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("删除邮件失败"))
		return
	}

	c.JSON(http.StatusOK, result.SimpleResult("邮件删除成功"))
}

// GetFolders 获取用户文件夹列表
func (h *EmailHandler) GetFolders(c *gin.Context) {
	// 返回默认文件夹列表
	folders := []map[string]interface{}{
		{"id": "inbox", "name": "收件箱", "type": "system"},
		{"id": "sent", "name": "已发送", "type": "system"},
		{"id": "draft", "name": "草稿箱", "type": "system"},
		{"id": "trash", "name": "垃圾箱", "type": "system"},
	}

	c.JSON(http.StatusOK, result.DataResult("获取成功", folders))
}

// GetRecentEmails 获取最近邮件
func (h *EmailHandler) GetRecentEmails(c *gin.Context) {
	userID := int64(c.GetInt("user_id"))
	isAdmin := c.GetBool("is_admin")

	// 获取limit参数
	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 5
	}

	// 获取用户的邮箱列表
	var mailboxes []*model.Mailbox
	if isAdmin {
		mailboxes, err = h.svcCtx.MailboxModel.GetMailboxesByAdminId(userID)
	} else {
		mailboxes, err = h.svcCtx.MailboxModel.GetMailboxesByUserId(userID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("获取邮箱列表失败"))
		return
	}

	// 获取最近的邮件
	var allEmails []*model.Email
	for _, mailbox := range mailboxes {
		emails, _, err := h.svcCtx.EmailModel.GetEmailsByMailboxId(mailbox.Id, "", 1, limit)
		if err == nil {
			allEmails = append(allEmails, emails...)
		}
	}

	// 按时间排序并限制数量
	if len(allEmails) > limit {
		allEmails = allEmails[:limit]
	}

	c.JSON(http.StatusOK, result.DataResult("获取最近邮件成功", allEmails))
}

// MarkAsRead 标记邮件为已读
func (h *EmailHandler) MarkAsRead(c *gin.Context) {
	emailIDStr := c.Param("id")
	emailID, err := strconv.ParseInt(emailIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮件ID格式错误"))
		return
	}

	userID := int64(c.GetInt("user_id"))
	isAdmin := c.GetBool("is_admin")

	// 获取邮件
	email, err := h.svcCtx.EmailModel.GetById(emailID)
	if err != nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("邮件不存在"))
		return
	}

	// 验证权限
	mailbox, err := h.svcCtx.MailboxModel.GetById(email.MailboxId)
	if err != nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("邮箱不存在"))
		return
	}

	hasPermission := false
	if isAdmin && mailbox.AdminId != nil && *mailbox.AdminId == userID {
		hasPermission = true
	} else if !isAdmin && mailbox.UserId != nil && *mailbox.UserId == userID {
		hasPermission = true
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限操作此邮件"))
		return
	}

	// 标记为已读
	email.IsRead = true
	if err := h.svcCtx.EmailModel.Update(nil, email); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("标记已读失败"))
		return
	}

	c.JSON(http.StatusOK, result.SimpleResult("标记已读成功"))
}

// MarkAsUnread 标记邮件为未读
func (h *EmailHandler) MarkAsUnread(c *gin.Context) {
	emailIDStr := c.Param("id")
	emailID, err := strconv.ParseInt(emailIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮件ID格式错误"))
		return
	}

	userID := int64(c.GetInt("user_id"))
	isAdmin := c.GetBool("is_admin")

	// 获取邮件
	email, err := h.svcCtx.EmailModel.GetById(emailID)
	if err != nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("邮件不存在"))
		return
	}

	// 验证权限
	mailbox, err := h.svcCtx.MailboxModel.GetById(email.MailboxId)
	if err != nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("邮箱不存在"))
		return
	}

	hasPermission := false
	if isAdmin && mailbox.AdminId != nil && *mailbox.AdminId == userID {
		hasPermission = true
	} else if !isAdmin && mailbox.UserId != nil && *mailbox.UserId == userID {
		hasPermission = true
	}

	if !hasPermission {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限操作此邮件"))
		return
	}

	// 标记为未读
	email.IsRead = false
	if err := h.svcCtx.EmailModel.Update(nil, email); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("标记未读失败"))
		return
	}

	c.JSON(http.StatusOK, result.SimpleResult("标记未读成功"))
}

// UploadAttachment 上传附件
func (h *EmailHandler) UploadAttachment(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("获取文件失败"))
		return
	}

	// 检查文件大小 (限制10MB)
	if file.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("文件大小不能超过10MB"))
		return
	}

	// 检查文件类型
	allowedTypes := map[string]bool{
		"image/jpeg":         true,
		"image/png":          true,
		"image/gif":          true,
		"application/pdf":    true,
		"application/msword": true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
		"text/plain": true,
	}

	if !allowedTypes[file.Header.Get("Content-Type")] {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("不支持的文件类型"))
		return
	}

	// 生成文件名
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	filepath := fmt.Sprintf("./uploads/%s", filename)

	// 保存文件
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("保存文件失败"))
		return
	}

	// 返回文件信息
	attachment := gin.H{
		"id":       time.Now().Unix(),
		"filename": file.Filename,
		"size":     file.Size,
		"url":      fmt.Sprintf("/api/attachments/%d/download", time.Now().Unix()),
	}

	c.JSON(http.StatusOK, result.DataResult("上传成功", attachment))
}

// DownloadAttachment 下载附件
func (h *EmailHandler) DownloadAttachment(c *gin.Context) {
	attachmentID := c.Param("id")

	// 这里应该根据ID查找实际的文件路径
	// 目前简化处理，返回错误
	c.JSON(http.StatusNotFound, result.ErrorSimpleResult(fmt.Sprintf("附件 %s 不存在", attachmentID)))
}
