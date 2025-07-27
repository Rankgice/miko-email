package handlers

import (
	"miko-email/internal/result"
	"miko-email/internal/svc"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine, sessionStore *sessions.CookieStore, svcCtx *svc.ServiceContext) {
	// 静态文件服务 - 服务Vue构建后的文件
	r.Static("/assets", "./webvue/dist/assets")
	r.StaticFile("/", "./webvue/dist/index.html")
	r.StaticFile("/favicon.ico", "./webvue/dist/favicon.ico")

	// 处理Vue Router的路由，所有非API路径都返回index.html
	r.NoRoute(func(c *gin.Context) {
		// 如果是API请求，返回404
		if len(c.Request.URL.Path) > 4 && c.Request.URL.Path[:4] == "/api" {
			c.JSON(404, gin.H{"success": false, "message": "API endpoint not found"})
			return
		}
		// 否则返回Vue应用的index.html
		c.File("./webvue/dist/index.html")
	})

	// 创建处理器实例
	authHandler := NewAuthHandler(sessionStore, svcCtx)
	userHandler := NewUserHandler(sessionStore, svcCtx)
	mailboxHandler := NewMailboxHandler(sessionStore, svcCtx)
	emailHandler := NewEmailHandler(sessionStore, svcCtx)
	domainHandler := NewDomainHandler(sessionStore, svcCtx)
	captchaHandler := NewCaptchaHandler(sessionStore, svcCtx)
	systemHandler := NewSystemHandler(svcCtx)
	wsHandler := NewWebSocketHandler()

	// WebSocket路由
	r.GET("/ws", wsHandler.HandleWebSocket)

	// API路由组
	api := r.Group("/api")
	{
		// 认证相关 - 不需要登录
		api.POST("/login", authHandler.Login)
		api.POST("/register", authHandler.Register)
		api.POST("/logout", authHandler.Logout)
		api.POST("/admin/login", authHandler.AdminLogin)

		// 公共API - 不需要登录
		api.GET("/domains/available", domainHandler.GetAvailable)

		// 需要认证的路由
		auth := api.Group("/")
		auth.Use(AuthMiddleware(sessionStore))
		{
			// 用户相关
			auth.GET("/user/profile", userHandler.GetProfile)
			auth.PUT("/user/profile", userHandler.UpdateProfile)
			auth.GET("/user/settings", userHandler.GetSettings)
			auth.PUT("/user/password", userHandler.ChangePassword)
			auth.PUT("/user/notifications", userHandler.UpdateNotifications)
			auth.PUT("/user/theme", userHandler.UpdateTheme)
			auth.GET("/user/signature", userHandler.GetSignature)

			// 仪表板统计
			auth.GET("/dashboard/stats", userHandler.GetDashboardStats)

			// 文件夹管理
			auth.GET("/folders", emailHandler.GetFolders)

			// 邮箱管理
			auth.GET("/mailboxes", mailboxHandler.List)
			auth.POST("/mailboxes", mailboxHandler.Create)
			auth.PUT("/mailboxes/:id", mailboxHandler.Update)
			auth.PUT("/mailboxes/:id/status", mailboxHandler.UpdateStatus)
			auth.DELETE("/mailboxes/:id", mailboxHandler.Delete)

			// 邮件管理
			auth.GET("/emails", emailHandler.List)
			auth.GET("/emails/:id", emailHandler.GetDetail)
			auth.DELETE("/emails/:id", emailHandler.Delete)
			auth.GET("/emails/sent", emailHandler.GetSentEmails)
			auth.GET("/emails/recent", emailHandler.GetRecentEmails)
			auth.POST("/emails/send", emailHandler.SendEmail)
			auth.POST("/emails/draft", emailHandler.SaveDraft)
			auth.POST("/emails/search", emailHandler.SearchEmails)
			auth.PUT("/emails/:id/read", emailHandler.MarkAsRead)
			auth.PUT("/emails/:id/unread", emailHandler.MarkAsUnread)

			// 转发规则管理
			auth.GET("/forward-rules", emailHandler.GetForwardRules)
			auth.POST("/forward-rules", emailHandler.CreateForwardRule)
			auth.GET("/forward-rules/:id", emailHandler.GetForwardRule)
			auth.PUT("/forward-rules/:id", emailHandler.UpdateForwardRule)
			auth.DELETE("/forward-rules/:id", emailHandler.DeleteForwardRule)
			auth.PUT("/forward-rules/:id/toggle", emailHandler.ToggleForwardRule)
			auth.POST("/forward-rules/:id/test", emailHandler.TestForwardRule)

			// 验证码管理
			auth.POST("/captcha/send", captchaHandler.SendCaptcha)
			auth.POST("/captcha/verify", captchaHandler.VerifyCaptcha)

			// 文件上传和下载
			auth.POST("/upload/attachment", emailHandler.UploadAttachment)
			auth.GET("/attachments/:id/download", emailHandler.DownloadAttachment)

			// 统计信息
			auth.GET("/stats/emails", userHandler.GetEmailStats)
			auth.GET("/stats/storage", userHandler.GetStorageStats)
		}

		// 管理员路由
		admin := api.Group("/admin")
		admin.Use(AuthMiddleware(sessionStore), AdminMiddleware())
		{
			admin.POST("/logout", authHandler.AdminLogout)
			admin.GET("/dashboard/stats", userHandler.GetAdminDashboardStats)
			admin.GET("/system/health", userHandler.GetSystemHealth)
			admin.GET("/activities/recent", userHandler.GetRecentActivities)
			admin.GET("/notifications", userHandler.GetAdminNotifications)

			// 用户管理
			admin.GET("/users", userHandler.GetAllUsers)
			admin.POST("/users", userHandler.CreateUser)
			admin.PUT("/users/:id/status", userHandler.UpdateUserStatus)
			admin.DELETE("/users/:id", userHandler.DeleteUser)

			// 邮箱管理
			admin.GET("/mailboxes", mailboxHandler.GetAllMailboxes)
			admin.POST("/mailboxes", mailboxHandler.CreateMailboxAdmin)
			admin.PUT("/mailboxes/:id", mailboxHandler.UpdateMailboxAdmin)
			admin.DELETE("/mailboxes/:id", mailboxHandler.DeleteMailboxAdmin)

			// 域名管理
			admin.GET("/domains", domainHandler.List)
			admin.POST("/domains", domainHandler.Create)
			admin.PUT("/domains/:id", domainHandler.UpdateDomain)
			admin.DELETE("/domains/:id", domainHandler.Delete)
			admin.POST("/domains/:id/verify", domainHandler.VerifyDomain)
			admin.GET("/domains/:id/dns-records", domainHandler.GetDomainDNSRecords)
			admin.GET("/domains/dkim-record", domainHandler.GetDomainDKIMRecord)
			admin.GET("/domains/server-info", domainHandler.GetServerInfo)
			admin.POST("/domains/:id/verify-sender", domainHandler.VerifySenderConfiguration)
			admin.POST("/domains/:id/verify-receiver", domainHandler.VerifyReceiverConfiguration)
			admin.PUT("/domains/:id/status", domainHandler.UpdateDomainStatus)

			// 系统设置
			admin.GET("/settings", userHandler.GetAdminSettings)
			admin.PUT("/settings", userHandler.UpdateAdminSettings)

			// 验证码管理（管理员专用）
			admin.GET("/captcha/codes", captchaHandler.GetCodes)
			admin.GET("/captcha/stats", captchaHandler.GetStats)
			admin.POST("/captcha/generate-test", captchaHandler.GenerateTestData)
			admin.DELETE("/captcha/codes/:id", captchaHandler.DeleteCode)
			admin.DELETE("/captcha/expired", captchaHandler.ClearExpired)
			admin.POST("/captcha/resend/:id", captchaHandler.ResendCode)

			// 验证码规则管理
			admin.GET("/captcha/rules", captchaHandler.GetCaptchaRules)
			admin.PUT("/captcha/rules/:id/status", captchaHandler.UpdateCaptchaRuleStatus)

			// 系统管理
			admin.GET("/system/settings", systemHandler.GetSystemSettings)
			admin.PUT("/system/settings", systemHandler.UpdateSystemSettings)
			admin.GET("/system/status", systemHandler.GetSystemStatus)
			admin.GET("/system/logs", systemHandler.GetSystemLogs)
		}

		// 验证码管理路由（需要认证但不需要管理员权限）
		captcha := api.Group("/captcha")
		captcha.Use(AuthMiddleware(sessionStore))
		{
			captcha.GET("/codes", captchaHandler.GetCodes)
			captcha.GET("/stats", captchaHandler.GetStats)
			captcha.POST("/generate-test", captchaHandler.GenerateTestData)
			captcha.DELETE("/codes/:id", captchaHandler.DeleteCode)
			captcha.DELETE("/expired", captchaHandler.ClearExpired)
			captcha.POST("/resend/:id", captchaHandler.ResendCode)
		}

	}
}

// AuthMiddleware 认证中间件
func AuthMiddleware(sessionStore *sessions.CookieStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := sessionStore.Get(c.Request, "miko-session")
		if err != nil {
			c.JSON(http.StatusUnauthorized, result.ErrorSimpleResult("会话错误"))
			c.Abort()
			return
		}

		userID, ok := session.Values["user_id"]
		if !ok || userID == nil {
			c.JSON(http.StatusUnauthorized, result.ErrorSimpleResult("请先登录"))
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", userID)
		c.Set("is_admin", session.Values["is_admin"])
		c.Set("username", session.Values["username"])

		c.Next()
	}
}

// AdminMiddleware 管理员中间件
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get("is_admin")
		if !exists || isAdmin != true {
			c.JSON(http.StatusForbidden, result.ErrorSimpleResult("需要管理员权限"))
			c.Abort()
			return
		}
		c.Next()
	}
}
