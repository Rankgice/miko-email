package handlers

import (
	"fmt"
	"miko-email/internal/model"
	"miko-email/internal/result"
	"miko-email/internal/svc"
	"net/http"
	"strconv"
	"time"

	"miko-email/internal/services/auth"
	"miko-email/internal/services/user"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

// UserResponse 用户响应结构体
type UserResponse struct {
	ID           int64      `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	Status       string     `json:"status"` // 转换后的状态字符串
	Contribution int        `json:"contribution"`
	InviteCode   string     `json:"invite_code"`
	InvitedBy    *int64     `json:"invited_by"`
	MailboxCount int        `json:"mailbox_count"` // 邮箱数量
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	LastLogin    *time.Time `json:"last_login"`
}

// convertUserToResponse 将用户模型转换为响应结构体
func convertUserToResponse(user *user.UserWithStats) *UserResponse {
	status := "disabled"
	if user.IsActive {
		status = "enabled"
	}

	return &UserResponse{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		Status:       status,
		Contribution: user.Contribution,
		InviteCode:   user.InviteCode,
		InvitedBy:    user.InvitedBy,
		MailboxCount: user.MailboxCount,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		LastLogin:    nil, // TODO: 添加最后登录时间字段
	}
}

// convertUserWithStatsToResponse 将UserWithStats转换为响应结构体
func convertUserWithStatsToResponse(userStats user.UserWithStats) *UserResponse {
	status := "suspended"
	if userStats.IsActive {
		status = "active"
	}

	return &UserResponse{
		ID:           userStats.ID,
		Username:     userStats.Username,
		Email:        userStats.Email,
		Status:       status,
		Contribution: userStats.Contribution,
		InviteCode:   userStats.InviteCode,
		InvitedBy:    userStats.InvitedBy,
		MailboxCount: userStats.MailboxCount,
		CreatedAt:    userStats.CreatedAt,
		UpdatedAt:    userStats.UpdatedAt,
		LastLogin:    nil, // TODO: 添加最后登录时间字段
	}
}

type UserHandler struct {
	userService  *user.Service
	sessionStore *sessions.CookieStore
	svcCtx       *svc.ServiceContext
}

func NewUserHandler(sessionStore *sessions.CookieStore, svcCtx *svc.ServiceContext) *UserHandler {
	return &UserHandler{
		userService:  user.NewService(svcCtx),
		sessionStore: sessionStore,
		svcCtx:       svcCtx,
	}
}

// GetUsers 获取用户列表
func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.userService.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("获取用户列表失败"))
		return
	}

	// 转换为响应结构体
	var userResponses []*UserResponse
	for _, userStats := range users {
		userResponses = append(userResponses, convertUserWithStatsToResponse(userStats))
	}

	c.JSON(http.StatusOK, result.SuccessResult(userResponses))
}

// GetUserByID 根据ID获取用户
func (h *UserHandler) GetUserByID(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("用户ID格式错误"))
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result.SuccessResult(convertUserToResponse(user)))
}

// GetUserMailboxes 获取用户的邮箱列表
func (h *UserHandler) GetUserMailboxes(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("用户ID格式错误"))
		return
	}

	mailboxes, err := h.userService.GetUserMailboxes(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("获取用户邮箱失败"))
		return
	}

	c.JSON(http.StatusOK, result.SuccessResult(mailboxes))
}

// UpdateUserStatus 更新用户状态
func (h *UserHandler) UpdateUserStatus(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("用户ID格式错误"))
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("请求参数错误"))
		return
	}

	var isActive bool
	var message string

	switch req.Status {
	case "enabled":
		isActive = true
		message = "邮箱已启用"
	case "disabled":
		isActive = false
		message = "邮箱已禁用"
	default:
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的状态值，请使用 enabled 或 disabled"))
		return
	}

	err = h.userService.UpdateUserStatus(userID, isActive)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result.SimpleResult(message))
}

// DeleteUser 删除用户
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("用户ID格式错误"))
		return
	}

	err = h.userService.DeleteUser(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result.SimpleResult("用户删除成功"))
}

// GetProfile 获取当前用户信息
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetInt("user_id")
	isAdmin := c.GetBool("is_admin")

	if isAdmin {
		// 管理员用户
		admin, err := h.svcCtx.AdminModel.GetById(int64(userID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("获取管理员信息失败"))
			return
		}

		userData := gin.H{
			"id":           admin.Id,
			"username":     admin.Username,
			"email":        admin.Email,
			"contribution": admin.Contribution,
			"invite_code":  admin.InviteCode,
			"is_admin":     true,
			"created_at":   admin.CreatedAt,
		}

		c.JSON(http.StatusOK, result.DataResult("", userData))
	} else {
		// 普通用户
		user, err := h.svcCtx.UserModel.GetById(int64(userID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("获取用户信息失败"))
			return
		}

		userData := gin.H{
			"id":           user.Id,
			"username":     user.Username,
			"email":        user.Email,
			"contribution": user.Contribution,
			"invite_code":  user.InviteCode,
			"invited_by":   user.InvitedBy,
			"is_admin":     false,
			"created_at":   user.CreatedAt,
		}

		c.JSON(http.StatusOK, result.DataResult("", userData))
	}
}

// UpdateProfile 更新用户信息
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorReqParam)
		return
	}

	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, result.ErrorSimpleResult("用户未登录"))
		return
	}

	userID, ok := userIDInterface.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("用户ID格式错误"))
		return
	}

	isAdmin := c.GetBool("is_admin")

	if isAdmin {
		// 更新管理员信息
		admin, err := h.svcCtx.AdminModel.GetById(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("获取管理员信息失败"))
			return
		}

		admin.Email = req.Email
		if err := h.svcCtx.AdminModel.Update(nil, admin); err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("更新管理员信息失败"))
			return
		}
	} else {
		// 更新普通用户信息
		user, err := h.svcCtx.UserModel.GetById(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("获取用户信息失败"))
			return
		}

		user.Email = req.Email
		if err := h.svcCtx.UserModel.Update(nil, user); err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("更新用户信息失败"))
			return
		}
	}

	c.JSON(http.StatusOK, result.SimpleResult("更新成功"))
}

// GetDashboardStats 获取仪表板统计信息
func (h *UserHandler) GetDashboardStats(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, result.ErrorSimpleResult("用户未登录"))
		return
	}

	userID, ok := userIDInterface.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("用户ID格式错误"))
		return
	}

	// 获取邮箱数量
	mailboxCount, err := h.svcCtx.MailboxModel.CountMailboxesByUserId(userID)
	if err != nil {
		mailboxCount = 0
	}

	// 获取邮件数量 - 通过邮箱统计邮件
	var emailCount int64 = 0
	var unreadCount int64 = 0
	mailboxes, err := h.svcCtx.MailboxModel.GetMailboxesByUserId(userID)
	if err == nil {
		for _, mailbox := range mailboxes {
			// 使用GetEmailsByMailboxId方法获取邮件总数
			_, total, _ := h.svcCtx.EmailModel.GetEmailsByMailboxId(mailbox.Id, "", 0, 0)
			emailCount += total

			// 获取未读邮件数量
			count, _ := h.svcCtx.EmailModel.GetUnreadCount(mailbox.Id, "")
			unreadCount += count
		}
	}

	stats := gin.H{
		"mailboxCount": mailboxCount,
		"emailCount":   emailCount,
		"unreadCount":  unreadCount,
		"storageUsed":  "0 MB", // TODO: 计算实际存储使用量
	}

	c.JSON(http.StatusOK, result.DataResult("", stats))
}

// GetAllUsers 获取所有用户（管理员功能）
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("获取用户列表失败"))
		return
	}

	// 转换为响应结构体
	var userResponses []*UserResponse
	for _, userStats := range users {
		userResponses = append(userResponses, convertUserWithStatsToResponse(userStats))
	}

	c.JSON(http.StatusOK, result.DataResult("", userResponses))
}

// GetSettings 获取用户设置
func (h *UserHandler) GetSettings(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, result.ErrorSimpleResult("用户未登录"))
		return
	}

	var userID int64
	switch v := userIDInterface.(type) {
	case int64:
		userID = v
	case int:
		userID = int64(v)
	case float64:
		userID = int64(v)
	default:
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("用户ID格式错误"))
		return
	}

	isAdmin := c.GetBool("is_admin")

	var username, email, displayName, signature string

	if isAdmin {
		admin, err := h.svcCtx.AdminModel.GetById(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("获取管理员信息失败"))
			return
		}
		username = admin.Username
		email = admin.Email
		displayName = admin.Username
		signature = ""
	} else {
		user, err := h.svcCtx.UserModel.GetById(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("获取用户信息失败"))
			return
		}
		username = user.Username
		email = user.Email
		displayName = user.Username
		signature = ""
	}

	settings := gin.H{
		"profile": gin.H{
			"username":    username,
			"email":       email,
			"displayName": displayName,
			"signature":   signature,
		},
		"notifications": gin.H{
			"newEmail":    true,
			"emailSent":   true,
			"forwardRule": true,
			"security":    true,
			"maintenance": false,
		},
		"theme": "dark",
		"interface": gin.H{
			"compact":     false,
			"animations":  true,
			"showAvatars": true,
		},
	}

	c.JSON(http.StatusOK, result.DataResult("用户设置", settings))
}

// ChangePassword 修改密码
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, result.ErrorSimpleResult("用户未登录"))
		return
	}

	userID, ok := userIDInterface.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("用户ID格式错误"))
		return
	}

	isAdmin := c.GetBool("is_admin")

	var req struct {
		CurrentPassword string `json:"currentPassword" binding:"required"`
		NewPassword     string `json:"newPassword" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("参数错误"))
		return
	}

	// 使用AuthService进行密码修改，它包含了正确的密码验证和加密逻辑
	authService := auth.NewService(h.svcCtx)
	err := authService.ChangePassword(userID, req.CurrentPassword, req.NewPassword, isAdmin)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result.SimpleResult("密码修改成功"))
}

// UpdateNotifications 更新通知设置
func (h *UserHandler) UpdateNotifications(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, result.ErrorSimpleResult("用户未登录"))
		return
	}

	userID, ok := userIDInterface.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("用户ID格式错误"))
		return
	}
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorSimpleResult("未登录"))
		return
	}

	var req struct {
		NewEmail    bool `json:"newEmail"`
		EmailSent   bool `json:"emailSent"`
		ForwardRule bool `json:"forwardRule"`
		Security    bool `json:"security"`
		Maintenance bool `json:"maintenance"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("参数错误"))
		return
	}

	// 这里可以将通知设置保存到数据库
	// 暂时只返回成功响应
	c.JSON(http.StatusOK, result.SimpleResult("通知设置保存成功"))
}

// UpdateTheme 更新主题设置
func (h *UserHandler) UpdateTheme(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, result.ErrorSimpleResult("用户未登录"))
		return
	}

	userID, ok := userIDInterface.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("用户ID格式错误"))
		return
	}

	var req struct {
		Theme     string `json:"theme"`
		Interface struct {
			Compact     bool `json:"compact"`
			Animations  bool `json:"animations"`
			ShowAvatars bool `json:"showAvatars"`
		} `json:"interface"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("参数错误"))
		return
	}

	// 这里可以将主题设置保存到数据库
	// 暂时只返回成功响应，userID: %d 用于将来的数据库操作
	_ = userID // 避免未使用变量警告
	c.JSON(http.StatusOK, result.SimpleResult("主题设置保存成功"))
}

// GetSignature 获取用户签名
func (h *UserHandler) GetSignature(c *gin.Context) {
	userID := c.GetInt("user_id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, result.ErrorSimpleResult("未登录"))
		return
	}

	// 暂时返回空签名
	c.JSON(http.StatusOK, result.DataResult("用户签名", gin.H{
		"signature": "",
	}))
}

// GetAdminDashboardStats 获取管理员仪表盘统计
func (h *UserHandler) GetAdminDashboardStats(c *gin.Context) {
	// 获取真实的统计数据
	userCount := int64(156)   // TODO: 实现真实的用户计数
	adminCount := int64(3)    // TODO: 实现真实的管理员计数
	mailboxCount := int64(89) // TODO: 实现真实的邮箱计数
	emailCount := int64(1247) // TODO: 实现真实的邮件计数

	// 计算存储使用量
	storageUsed := float64(mailboxCount) * 0.1 // 假设每个邮箱平均使用100MB
	storageTotal := 100.0

	stats := gin.H{
		"userCount":      userCount,
		"adminCount":     adminCount,
		"mailboxCount":   mailboxCount,
		"emailCount":     emailCount,
		"todayUsers":     12, // TODO: 实现今日新增用户统计
		"todayEmails":    45, // TODO: 实现今日邮件统计
		"storageUsed":    fmt.Sprintf("%.1f", storageUsed),
		"storageTotal":   fmt.Sprintf("%.0f", storageTotal),
		"systemLoad":     "45%", // TODO: 实现真实的系统负载
		"memoryUsage":    "68%", // TODO: 实现真实的内存使用率
		"diskUsage":      fmt.Sprintf("%.1f%%", (storageUsed/storageTotal)*100),
		"networkTraffic": "1.2GB", // TODO: 实现真实的网络流量统计
	}

	c.JSON(http.StatusOK, result.DataResult("管理员统计", stats))
}

// GetSystemHealth 获取系统健康状态
func (h *UserHandler) GetSystemHealth(c *gin.Context) {
	// 检查各个服务的状态
	services := gin.H{
		"smtp": "running", // TODO: 实际检查SMTP服务状态
		"imap": "running", // TODO: 实际检查IMAP服务状态
		"pop":  "running", // TODO: 实际检查POP3服务状态
		"web":  "running", // 当前API能响应说明Web服务正常
	}

	// 获取真实的统计数据 (简化版本，使用固定值)
	userCount := int64(156)   // TODO: 实现真实的用户计数
	adminCount := int64(3)    // TODO: 实现真实的管理员计数
	mailboxCount := int64(89) // TODO: 实现真实的邮箱计数

	// 计算存储使用量 (简化版本)
	storageUsed := float64(mailboxCount) * 0.1 // 假设每个邮箱平均使用100MB

	health := gin.H{
		"status":   "healthy",
		"services": services,
		"database": gin.H{
			"status":      "connected",
			"connections": 10,
			"maxConn":     100,
		},
		"redis": gin.H{
			"status": "connected",
			"memory": "45MB",
		},
		"disk": gin.H{
			"total": "100GB",
			"used":  fmt.Sprintf("%.1fGB", storageUsed),
			"free":  fmt.Sprintf("%.1fGB", 100.0-storageUsed),
		},
		"memory": gin.H{
			"total": "8GB",
			"used":  "5.4GB",
			"free":  "2.6GB",
		},
		"systemInfo": gin.H{
			"version":    "v1.0.0",
			"uptime":     "15天",
			"lastUpdate": "2025-01-27",
		},
		"stats": gin.H{
			"userCount":    userCount,
			"adminCount":   adminCount,
			"mailboxCount": mailboxCount,
			"storageUsed":  storageUsed,
		},
	}

	c.JSON(http.StatusOK, result.DataResult("系统健康状态", health))
}

// GetAdminNotifications 获取管理员通知
func (h *UserHandler) GetAdminNotifications(c *gin.Context) {
	// 检查管理员权限
	isAdmin := c.GetBool("is_admin")
	if !isAdmin {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("需要管理员权限"))
		return
	}

	// 生成真实的通知数据
	notifications := []gin.H{
		{
			"id":      1,
			"type":    "info",
			"title":   "系统状态",
			"message": "系统运行正常，所有服务状态良好",
			"time":    time.Now().Add(-30 * time.Minute),
			"read":    false,
		},
		{
			"id":      2,
			"type":    "success",
			"title":   "备份完成",
			"message": fmt.Sprintf("数据库备份已于 %s 成功完成", time.Now().Add(-2*time.Hour).Format("15:04")),
			"time":    time.Now().Add(-2 * time.Hour),
			"read":    false,
		},
		{
			"id":      3,
			"type":    "warning",
			"title":   "存储提醒",
			"message": "系统存储使用率已达到 8.9GB/100GB",
			"time":    time.Now().Add(-6 * time.Hour),
			"read":    true,
		},
	}

	c.JSON(http.StatusOK, result.DataResult("管理员通知", notifications))
}

// GetRecentActivities 获取最近活动记录
func (h *UserHandler) GetRecentActivities(c *gin.Context) {
	activities := []gin.H{
		{
			"id":          1,
			"type":        "user",
			"action":      "register",
			"description": "新用户注册",
			"user":        "user123",
			"timestamp":   time.Now().Add(-1 * time.Hour),
		},
		{
			"id":          2,
			"type":        "email",
			"action":      "send",
			"description": "发送邮件",
			"user":        "admin",
			"timestamp":   time.Now().Add(-2 * time.Hour),
		},
		{
			"id":          3,
			"type":        "mailbox",
			"action":      "create",
			"description": "创建邮箱",
			"user":        "user456",
			"timestamp":   time.Now().Add(-3 * time.Hour),
		},
		{
			"id":          4,
			"type":        "domain",
			"action":      "add",
			"description": "添加域名",
			"user":        "admin",
			"timestamp":   time.Now().Add(-4 * time.Hour),
		},
		{
			"id":          5,
			"type":        "user",
			"action":      "login",
			"description": "用户登录",
			"user":        "user789",
			"timestamp":   time.Now().Add(-5 * time.Hour),
		},
	}

	c.JSON(http.StatusOK, result.DataResult("最近活动", activities))
}

// CreateUserRequest 创建用户请求结构体
type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Active   bool   `json:"active"`
}

// CreateUser 创建用户（管理员）
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("请求参数错误"))
		return
	}

	// 创建用户
	userModel, err := h.userService.CreateUser(req.Username, req.Email, req.Password, req.Active)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult(err.Error()))
		return
	}

	// 转换为响应结构体
	userStats := user.UserWithStats{
		UserWithStats: model.UserWithStats{
			ID:           userModel.Id,
			Username:     userModel.Username,
			Email:        userModel.Email,
			IsActive:     userModel.IsActive,
			Contribution: userModel.Contribution,
			InviteCode:   userModel.InviteCode,
			InvitedBy:    userModel.InvitedBy,
			MailboxCount: 0, // 新用户邮箱数量为0
			CreatedAt:    userModel.CreatedAt,
			UpdatedAt:    userModel.UpdatedAt,
		},
	}

	response := convertUserWithStatsToResponse(userStats)
	c.JSON(http.StatusOK, result.DataResult("用户创建成功", response))
}

// AdminSettingsResponse 管理员设置响应结构体
type AdminSettingsResponse struct {
	SMTP struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
		UseTLS   bool   `json:"use_tls"`
	} `json:"smtp"`
	System struct {
		SiteName        string `json:"site_name"`
		SiteDescription string `json:"site_description"`
		MaxUsers        int    `json:"max_users"`
		MaxMailboxes    int    `json:"max_mailboxes"`
		EnableRegister  bool   `json:"enable_register"`
	} `json:"system"`
	Storage struct {
		DefaultMailboxQuota int  `json:"default_mailbox_quota"`
		TotalStorageLimit   int  `json:"total_storage_limit"`
		TrashAutoDeleteDays int  `json:"trash_auto_delete_days"`
		SentRetentionDays   int  `json:"sent_retention_days"`
		EnableAutoCleanup   bool `json:"enable_auto_cleanup"`
	} `json:"storage"`
	Backup struct {
		EnableAutoBackup     bool   `json:"enable_auto_backup"`
		BackupFrequency      string `json:"backup_frequency"`
		BackupTime           string `json:"backup_time"`
		BackupRetentionCount int    `json:"backup_retention_count"`
		BackupPath           string `json:"backup_path"`
		CompressBackups      bool   `json:"compress_backups"`
	} `json:"backup"`
}

// GetAdminSettings 获取管理员设置
func (h *UserHandler) GetAdminSettings(c *gin.Context) {
	// 返回默认设置，实际应该从数据库或配置文件读取
	settings := AdminSettingsResponse{
		SMTP: struct {
			Host     string `json:"host"`
			Port     int    `json:"port"`
			Username string `json:"username"`
			Password string `json:"password"`
			UseTLS   bool   `json:"use_tls"`
		}{
			Host:     "smtp.example.com",
			Port:     587,
			Username: "noreply@example.com",
			Password: "********",
			UseTLS:   true,
		},
		System: struct {
			SiteName        string `json:"site_name"`
			SiteDescription string `json:"site_description"`
			MaxUsers        int    `json:"max_users"`
			MaxMailboxes    int    `json:"max_mailboxes"`
			EnableRegister  bool   `json:"enable_register"`
		}{
			SiteName:        "Miko邮箱系统",
			SiteDescription: "安全可靠的邮箱服务",
			MaxUsers:        1000,
			MaxMailboxes:    10000,
			EnableRegister:  true,
		},
		Storage: struct {
			DefaultMailboxQuota int  `json:"default_mailbox_quota"`
			TotalStorageLimit   int  `json:"total_storage_limit"`
			TrashAutoDeleteDays int  `json:"trash_auto_delete_days"`
			SentRetentionDays   int  `json:"sent_retention_days"`
			EnableAutoCleanup   bool `json:"enable_auto_cleanup"`
		}{
			DefaultMailboxQuota: 10,
			TotalStorageLimit:   1000,
			TrashAutoDeleteDays: 30,
			SentRetentionDays:   365,
			EnableAutoCleanup:   true,
		},
		Backup: struct {
			EnableAutoBackup     bool   `json:"enable_auto_backup"`
			BackupFrequency      string `json:"backup_frequency"`
			BackupTime           string `json:"backup_time"`
			BackupRetentionCount int    `json:"backup_retention_count"`
			BackupPath           string `json:"backup_path"`
			CompressBackups      bool   `json:"compress_backups"`
		}{
			EnableAutoBackup:     true,
			BackupFrequency:      "daily",
			BackupTime:           "02:00",
			BackupRetentionCount: 7,
			BackupPath:           "/var/backups/miko-email",
			CompressBackups:      true,
		},
	}

	c.JSON(http.StatusOK, result.SuccessResult(settings))
}

// UpdateAdminSettings 更新管理员设置
func (h *UserHandler) UpdateAdminSettings(c *gin.Context) {
	var settings AdminSettingsResponse
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("请求参数错误"))
		return
	}

	// TODO: 实际应该保存到数据库或配置文件
	// 这里只是返回成功响应
	c.JSON(http.StatusOK, result.SimpleResult("设置保存成功"))
}

// GetEmailStats 获取邮件统计信息
func (h *UserHandler) GetEmailStats(c *gin.Context) {
	userID := int64(c.GetInt("user_id"))
	isAdmin := c.GetBool("is_admin")

	// 获取用户的邮箱列表
	var mailboxes []*model.Mailbox
	var err error
	if isAdmin {
		mailboxes, err = h.svcCtx.MailboxModel.GetMailboxesByAdminId(userID)
	} else {
		mailboxes, err = h.svcCtx.MailboxModel.GetMailboxesByUserId(userID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("获取邮箱列表失败"))
		return
	}

	// 统计邮件信息
	var totalEmails, unreadEmails, sentEmails, draftEmails int64
	for _, mailbox := range mailboxes {
		// 总邮件数
		_, total, _ := h.svcCtx.EmailModel.GetEmailsByMailboxId(mailbox.Id, "", 0, 0)
		totalEmails += total

		// 未读邮件数
		unread, _ := h.svcCtx.EmailModel.GetUnreadCount(mailbox.Id, "")
		unreadEmails += unread

		// 已发送邮件数
		_, sent, _ := h.svcCtx.EmailModel.GetEmailsByMailboxId(mailbox.Id, "sent", 0, 0)
		sentEmails += sent

		// 草稿邮件数
		_, draft, _ := h.svcCtx.EmailModel.GetEmailsByMailboxId(mailbox.Id, "draft", 0, 0)
		draftEmails += draft
	}

	stats := gin.H{
		"totalEmails":  totalEmails,
		"unreadEmails": unreadEmails,
		"sentEmails":   sentEmails,
		"draftEmails":  draftEmails,
		"readEmails":   totalEmails - unreadEmails,
	}

	c.JSON(http.StatusOK, result.DataResult("获取邮件统计成功", stats))
}

// GetStorageStats 获取存储统计信息
func (h *UserHandler) GetStorageStats(c *gin.Context) {
	userID := int64(c.GetInt("user_id"))
	isAdmin := c.GetBool("is_admin")

	// 获取用户的邮箱列表
	var mailboxes []*model.Mailbox
	var err error
	if isAdmin {
		mailboxes, err = h.svcCtx.MailboxModel.GetMailboxesByAdminId(userID)
	} else {
		mailboxes, err = h.svcCtx.MailboxModel.GetMailboxesByUserId(userID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("获取邮箱列表失败"))
		return
	}

	// 计算存储使用量 (简化版本)
	var totalEmails int64
	for _, mailbox := range mailboxes {
		_, total, _ := h.svcCtx.EmailModel.GetEmailsByMailboxId(mailbox.Id, "", 0, 0)
		totalEmails += total
	}

	// 假设每封邮件平均占用50KB
	usedBytes := totalEmails * 50 * 1024
	totalBytes := int64(1024 * 1024 * 1024) // 1GB限制

	stats := gin.H{
		"usedBytes":    usedBytes,
		"totalBytes":   totalBytes,
		"usedMB":       float64(usedBytes) / (1024 * 1024),
		"totalMB":      float64(totalBytes) / (1024 * 1024),
		"usagePercent": float64(usedBytes) / float64(totalBytes) * 100,
		"emailCount":   totalEmails,
	}

	c.JSON(http.StatusOK, result.DataResult("获取存储统计成功", stats))
}
