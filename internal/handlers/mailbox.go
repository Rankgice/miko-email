package handlers

import (
	"miko-email/internal/result"
	"miko-email/internal/svc"
	"net/http"
	"strconv"
	"strings"

	"miko-email/internal/services/mailbox"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

type MailboxHandler struct {
	mailboxService *mailbox.Service
	sessionStore   *sessions.CookieStore
	svcCtx         *svc.ServiceContext
}

func NewMailboxHandler(sessionStore *sessions.CookieStore, svcCtx *svc.ServiceContext) *MailboxHandler {
	return &MailboxHandler{
		mailboxService: mailbox.NewService(svcCtx),
		sessionStore:   sessionStore,
		svcCtx:         svcCtx,
	}
}

type CreateMailboxRequest struct {
	Prefix   string `json:"prefix" binding:"required"`
	DomainID int64  `json:"domain_id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type BatchCreateMailboxRequest struct {
	Prefixes []string `json:"prefixes" binding:"required"`
	DomainID int64    `json:"domain_id" binding:"required"`
}

// GetMailboxes 获取邮箱列表
func (h *MailboxHandler) GetMailboxes(c *gin.Context) {
	userID := int64(c.GetInt("user_id"))
	isAdmin := c.GetBool("is_admin")

	mailboxes, err := h.mailboxService.GetUserMailboxes(userID, isAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("获取邮箱列表失败"))
		return
	}

	c.JSON(http.StatusOK, result.SuccessResult(mailboxes))
}

// CreateMailbox 创建邮箱
func (h *MailboxHandler) CreateMailbox(c *gin.Context) {
	var req CreateMailboxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("请求参数错误"))
		return
	}

	userID := int64(c.GetInt("user_id"))
	isAdmin := c.GetBool("is_admin")

	// 验证前缀格式
	if !isValidEmailPrefix(req.Prefix) {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮箱前缀格式不正确"))
		return
	}

	if len(req.Password) < 6 {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("密码长度至少6位"))
		return
	}

	mailbox, err := h.mailboxService.CreateMailboxWithPassword(userID, req.Prefix, req.Password, req.DomainID, isAdmin)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result.DataResult("邮箱创建成功", mailbox))
}

// BatchCreateMailboxes 批量创建邮箱
func (h *MailboxHandler) BatchCreateMailboxes(c *gin.Context) {
	var req BatchCreateMailboxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("请求参数错误"))
		return
	}

	userID := int64(c.GetInt("user_id"))
	isAdmin := c.GetBool("is_admin")

	// 验证所有前缀格式
	for _, prefix := range req.Prefixes {
		if !isValidEmailPrefix(prefix) {
			c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮箱前缀格式不正确: "+prefix))
			return
		}
	}

	mailboxes, err := h.mailboxService.BatchCreateMailboxes(userID, req.Prefixes, req.DomainID, isAdmin)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result.DataResult("批量创建邮箱成功", mailboxes))
}

// GetMailboxPassword 获取邮箱密码
func (h *MailboxHandler) GetMailboxPassword(c *gin.Context) {
	mailboxIDStr := c.Param("id")
	mailboxID, err := strconv.ParseInt(mailboxIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮箱ID格式错误"))
		return
	}

	userID := int64(c.GetInt("user_id"))
	isAdmin := c.GetBool("is_admin")

	password, err := h.mailboxService.GetMailboxPassword(mailboxID, userID, isAdmin)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result.SuccessResult(gin.H{
		"password": password,
	}))
}

// DeleteMailbox 删除邮箱
func (h *MailboxHandler) DeleteMailbox(c *gin.Context) {
	mailboxIDStr := c.Param("id")
	mailboxID, err := strconv.ParseInt(mailboxIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮箱ID格式错误"))
		return
	}

	userID := int64(c.GetInt("user_id"))
	isAdmin := c.GetBool("is_admin")

	err = h.mailboxService.DeleteMailbox(mailboxID, userID, isAdmin)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result.SimpleResult("邮箱删除成功"))
}

// Update 更新邮箱（普通用户）
func (h *MailboxHandler) Update(c *gin.Context) {
	mailboxIDStr := c.Param("id")
	mailboxID, err := strconv.ParseInt(mailboxIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮箱ID格式错误"))
		return
	}

	var req struct {
		Password string `json:"password"`
		IsActive *bool  `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("请求参数错误"))
		return
	}

	userID := int64(c.GetInt("user_id"))
	isAdmin := c.GetBool("is_admin")

	// 验证邮箱所有权
	mailbox, err := h.svcCtx.MailboxModel.GetById(mailboxID)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮箱不存在"))
		return
	}

	// 检查权限
	if !isAdmin && (mailbox.UserId == nil || *mailbox.UserId != userID) {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限操作此邮箱"))
		return
	}

	// 更新密码
	if req.Password != "" {
		err = h.svcCtx.MailboxModel.UpdatePassword(nil, mailboxID, req.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("更新密码失败"))
			return
		}
	}

	// 更新状态
	if req.IsActive != nil {
		err = h.svcCtx.MailboxModel.UpdateStatus(nil, mailboxID, *req.IsActive)
		if err != nil {
			c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("更新状态失败"))
			return
		}
	}

	c.JSON(http.StatusOK, result.SimpleResult("邮箱更新成功"))
}

// UpdateStatus 更新邮箱状态（普通用户）
func (h *MailboxHandler) UpdateStatus(c *gin.Context) {
	mailboxIDStr := c.Param("id")
	mailboxID, err := strconv.ParseInt(mailboxIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮箱ID格式错误"))
		return
	}

	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("请求参数错误"))
		return
	}

	userID := int64(c.GetInt("user_id"))
	isAdmin := c.GetBool("is_admin")

	// 验证邮箱所有权
	mailbox, err := h.svcCtx.MailboxModel.GetById(mailboxID)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮箱不存在"))
		return
	}

	// 检查权限
	if !isAdmin && (mailbox.UserId == nil || *mailbox.UserId != userID) {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("无权限操作此邮箱"))
		return
	}

	// 更新状态
	err = h.svcCtx.MailboxModel.UpdateStatus(nil, mailboxID, req.IsActive)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("更新状态失败"))
		return
	}

	status := "启用"
	if !req.IsActive {
		status = "禁用"
	}

	c.JSON(http.StatusOK, result.SimpleResult("邮箱"+status+"成功"))
}

// isValidEmailPrefix 验证邮箱前缀格式
func isValidEmailPrefix(prefix string) bool {
	if len(prefix) == 0 || len(prefix) > 64 {
		return false
	}

	// 简单的邮箱前缀验证
	for _, char := range prefix {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '.' || char == '-' || char == '_') {
			return false
		}
	}

	// 不能以点、横线或下划线开头或结尾
	if strings.HasPrefix(prefix, ".") || strings.HasSuffix(prefix, ".") ||
		strings.HasPrefix(prefix, "-") || strings.HasSuffix(prefix, "-") ||
		strings.HasPrefix(prefix, "_") || strings.HasSuffix(prefix, "_") {
		return false
	}

	return true
}

// 管理员邮箱管理接口

// GetAllMailboxes 获取所有邮箱列表（管理员）
func (h *MailboxHandler) GetAllMailboxes(c *gin.Context) {
	mailboxes, err := h.mailboxService.GetAllMailboxes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("获取邮箱列表失败"))
		return
	}

	c.JSON(http.StatusOK, result.SuccessResult(mailboxes))
}

// UpdateMailboxStatus 更新邮箱状态（管理员）
func (h *MailboxHandler) UpdateMailboxStatus(c *gin.Context) {
	mailboxIDStr := c.Param("id")
	mailboxID, err := strconv.ParseInt(mailboxIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮箱ID格式错误"))
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("请求参数错误"))
		return
	}

	if req.Status != "active" && req.Status != "suspended" {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("状态值无效"))
		return
	}

	err = h.mailboxService.UpdateMailboxStatus(mailboxID, req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result.SimpleResult("邮箱状态更新成功"))
}

// DeleteMailboxAdmin 删除邮箱（管理员）
func (h *MailboxHandler) DeleteMailboxAdmin(c *gin.Context) {
	mailboxIDStr := c.Param("id")
	mailboxID, err := strconv.ParseInt(mailboxIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮箱ID格式错误"))
		return
	}

	err = h.mailboxService.DeleteMailboxAdmin(mailboxID)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result.SimpleResult("邮箱删除成功"))
}

// AdminCreateMailboxRequest 管理员创建邮箱请求
type AdminCreateMailboxRequest struct {
	UserID   int64  `json:"user_id" binding:"required"`
	Prefix   string `json:"prefix" binding:"required"`
	DomainID int64  `json:"domain_id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// CreateMailboxAdmin 创建邮箱（管理员）
func (h *MailboxHandler) CreateMailboxAdmin(c *gin.Context) {
	var req AdminCreateMailboxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorReqParam)
		return
	}

	// 验证前缀格式
	if !isValidEmailPrefix(req.Prefix) {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮箱前缀格式不正确"))
		return
	}

	if len(req.Password) < 6 {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("密码长度至少6位"))
		return
	}

	mailbox, err := h.mailboxService.CreateMailboxAdmin(req.UserID, req.Prefix, req.DomainID, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result.DataResult("邮箱创建成功", mailbox))
}

// GetMailboxStats 获取邮箱统计信息（管理员）
func (h *MailboxHandler) GetMailboxStats(c *gin.Context) {
	mailboxIDStr := c.Param("id")
	mailboxID, err := strconv.ParseInt(mailboxIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮箱ID格式错误"))
		return
	}

	stats, err := h.mailboxService.GetMailboxStats(mailboxID)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result.SuccessResult(stats))
}

// GetUserStats 获取用户统计信息
func (h *MailboxHandler) GetUserStats(c *gin.Context) {
	userID := int64(c.GetInt("user_id"))

	stats, err := h.mailboxService.GetUserStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("获取统计信息失败: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, result.SuccessResult(stats))
}

// List 获取邮箱列表（API）
func (h *MailboxHandler) List(c *gin.Context) {
	userID := int64(c.GetInt("user_id"))
	isAdmin := c.GetBool("is_admin")

	mailboxes, err := h.mailboxService.GetUserMailboxes(userID, isAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("获取邮箱列表失败"))
		return
	}

	c.JSON(http.StatusOK, result.DataResult("", mailboxes))
}

// Create 创建邮箱（API）
func (h *MailboxHandler) Create(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorReqParam)
		return
	}

	userID := int64(c.GetInt("user_id"))
	isAdmin := c.GetBool("is_admin")

	// 解析邮箱地址获取前缀和域名
	parts := strings.Split(req.Email, "@")
	if len(parts) != 2 {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮箱地址格式错误"))
		return
	}

	prefix := parts[0]
	domainName := parts[1]

	// 查找域名ID
	domain, err := h.svcCtx.DomainModel.GetByName(domainName)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("域名不存在"))
		return
	}

	mailbox, err := h.mailboxService.CreateMailboxWithPassword(userID, prefix, req.Password, domain.Id, isAdmin)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result.DataResult("邮箱创建成功", mailbox))
}

// Delete 删除邮箱（API）
func (h *MailboxHandler) Delete(c *gin.Context) {
	mailboxIDStr := c.Param("id")
	mailboxID, err := strconv.ParseInt(mailboxIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮箱ID格式错误"))
		return
	}

	userID := int64(c.GetInt("user_id"))
	isAdmin := c.GetBool("is_admin")

	err = h.mailboxService.DeleteMailbox(mailboxID, userID, isAdmin)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result.SimpleResult("邮箱删除成功"))
}

// UpdateMailboxAdmin 更新邮箱（管理员）
func (h *MailboxHandler) UpdateMailboxAdmin(c *gin.Context) {
	mailboxIDStr := c.Param("id")
	mailboxID, err := strconv.ParseInt(mailboxIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮箱ID格式错误"))
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Active   bool   `json:"active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("请求参数错误"))
		return
	}

	// 检查邮箱是否存在
	user, err := h.svcCtx.UserModel.GetById(mailboxID)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮箱不存在"))
		return
	}

	// 更新邮箱状态
	if user.IsActive != req.Active {
		err = h.svcCtx.UserModel.UpdateStatus(nil, mailboxID, req.Active)
		if err != nil {
			c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("更新邮箱状态失败"))
			return
		}
	}

	// 如果提供了新密码，更新密码
	if req.Password != "" {
		err = h.svcCtx.UserModel.UpdatePassword(nil, mailboxID, req.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("更新密码失败"))
			return
		}
	}

	// 如果邮箱地址有变化，更新邮箱地址（这里简化处理，实际可能需要更复杂的验证）
	if req.Email != "" && req.Email != user.Email {
		// 检查新邮箱地址是否已存在
		existingUser, _ := h.svcCtx.UserModel.GetByEmail(req.Email)
		if existingUser != nil && existingUser.Id != mailboxID {
			c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("邮箱地址已存在"))
			return
		}

		err = h.svcCtx.UserModel.UpdateEmail(nil, mailboxID, req.Email)
		if err != nil {
			c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("更新邮箱地址失败"))
			return
		}
	}

	c.JSON(http.StatusOK, result.DataResult("邮箱更新成功", nil))
}
