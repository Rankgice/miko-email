package server

import (
	"database/sql"
	"net/http"
	"strings"

	"miko-email/internal/config"
	"miko-email/internal/handlers"
	"miko-email/internal/middleware"
	"miko-email/internal/services/auth"
	"miko-email/internal/services/domain"
	"miko-email/internal/services/email"
	"miko-email/internal/services/forward"
	"miko-email/internal/services/mailbox"
	"miko-email/internal/services/user"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

type Server struct {
	router         *gin.Engine
	db             *sql.DB
	config         *config.Config
	sessionStore   *sessions.CookieStore
	emailService   *email.Service
	forwardService *forward.Service
}

func New(db *sql.DB, cfg *config.Config) *Server {
	// 创建session store
	sessionStore := sessions.NewCookieStore([]byte(cfg.SessionKey))
	sessionStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7天
		HttpOnly: true,
		Secure:   false, // 在生产环境中应该设置为true
		SameSite: http.SameSiteLaxMode,
	}

	// 创建邮件服务
	emailService := email.NewService(db)

	// 创建转发服务
	forwardService := forward.NewService(db)

	server := &Server{
		router:         gin.Default(),
		db:             db,
		config:         cfg,
		sessionStore:   sessionStore,
		emailService:   emailService,
		forwardService: forwardService,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	// 设置UTF-8编码中间件
	s.router.Use(func(c *gin.Context) {
		// 对于API请求，设置JSON编码
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.Header("Content-Type", "application/json; charset=utf-8")
		} else {
			// 对于HTML页面，设置HTML编码
			c.Header("Content-Type", "text/html; charset=utf-8")
		}
		c.Next()
	})

	// 创建服务实例
	authService := auth.NewService(s.db)
	mailboxService := mailbox.NewService(s.db)
	domainService := domain.NewService(s.db)
	userService := user.NewService(s.db)

	// 创建处理器实例
	authHandler := handlers.NewAuthHandler(authService, s.sessionStore, s.db)
	mailboxHandler := handlers.NewMailboxHandler(mailboxService, s.sessionStore)
	domainHandler := handlers.NewDomainHandler(domainService, s.sessionStore)
	userHandler := handlers.NewUserHandler(userService, s.sessionStore)
	emailHandler := handlers.NewEmailHandler(s.emailService, mailboxService, s.forwardService, s.sessionStore)
	webHandler := handlers.NewWebHandler(s.sessionStore)

	// 中间件
	authMiddleware := middleware.NewAuthMiddleware(s.sessionStore)
	adminMiddleware := middleware.NewAdminMiddleware(s.sessionStore)

	// 静态文件
	s.router.Static("/static", "./web/static")
	s.router.LoadHTMLGlob("web/templates/*")

	// Web页面路由
	web := s.router.Group("/")
	{
		web.GET("/", webHandler.Home)
		web.GET("/login", webHandler.LoginPage)
		web.GET("/register", webHandler.RegisterPage)
		web.GET("/admin/login", webHandler.AdminLoginPage)
		
		// 需要登录的页面
		webAuth := web.Group("/")
		webAuth.Use(authMiddleware.RequireAuth())
		{
			webAuth.GET("/dashboard", webHandler.Dashboard)
			webAuth.GET("/compose", webHandler.ComposePage)
			webAuth.GET("/forward", webHandler.ForwardPage)
			webAuth.GET("/inbox", webHandler.InboxPage)
			webAuth.GET("/sent", webHandler.SentPage)
			webAuth.GET("/settings", webHandler.SettingsPage)
			webAuth.GET("/mailboxes", webHandler.MailboxesPage)
		}

		// 管理员页面
		webAdmin := web.Group("/admin")
		webAdmin.Use(adminMiddleware.RequireAdmin())
		{
			webAdmin.GET("/dashboard", webHandler.AdminDashboard)
			webAdmin.GET("/users", webHandler.UsersPage)
			webAdmin.GET("/mailboxes", webHandler.AdminMailboxesPage)
			webAdmin.GET("/domains", webHandler.DomainsPage)
		}
	}

	// API路由
	api := s.router.Group("/api")
	{
		// 认证相关
		api.POST("/login", authHandler.Login)
		api.POST("/register", authHandler.Register)
		api.POST("/admin/login", authHandler.AdminLogin)
		api.POST("/logout", authHandler.Logout)

		// 需要登录的API
		apiAuth := api.Group("/")
		apiAuth.Use(authMiddleware.RequireAuth())
		{
			// 用户信息
			apiAuth.GET("/profile", authHandler.GetProfile)
			apiAuth.PUT("/profile/password", authHandler.ChangePassword)

			// 邮箱管理
			apiAuth.GET("/mailboxes", mailboxHandler.GetMailboxes)
			apiAuth.POST("/mailboxes", mailboxHandler.CreateMailbox)
			apiAuth.POST("/mailboxes/batch", mailboxHandler.BatchCreateMailboxes)
			apiAuth.GET("/mailboxes/:id/password", mailboxHandler.GetMailboxPassword)
			apiAuth.DELETE("/mailboxes/:id", mailboxHandler.DeleteMailbox)

			// 邮件相关
			apiAuth.GET("/emails", emailHandler.GetEmails)
			apiAuth.GET("/emails/:id", emailHandler.GetEmailByID)
			apiAuth.POST("/emails/send", emailHandler.SendEmail)
			apiAuth.DELETE("/emails/:id", emailHandler.DeleteEmail)

			// 转发规则相关
			apiAuth.GET("/forward-rules", emailHandler.GetForwardRules)
			apiAuth.POST("/forward-rules", emailHandler.CreateForwardRule)
			apiAuth.GET("/forward-rules/:id", emailHandler.GetForwardRule)
			apiAuth.PUT("/forward-rules/:id", emailHandler.UpdateForwardRule)
			apiAuth.DELETE("/forward-rules/:id", emailHandler.DeleteForwardRule)
			apiAuth.PATCH("/forward-rules/:id/toggle", emailHandler.ToggleForwardRule)
			apiAuth.POST("/forward-rules/:id/test", emailHandler.TestForwardRule)
			apiAuth.GET("/forward-statistics", emailHandler.GetForwardStatistics)
		}

		// 管理员API
		apiAdmin := api.Group("/admin")
		apiAdmin.Use(adminMiddleware.RequireAdmin())
		{
			// 域名管理
			apiAdmin.GET("/domains", domainHandler.GetDomains)
			apiAdmin.POST("/domains", domainHandler.CreateDomain)
			apiAdmin.PUT("/domains/:id", domainHandler.UpdateDomain)
			apiAdmin.DELETE("/domains/:id", domainHandler.DeleteDomain)
			apiAdmin.POST("/domains/:id/verify", domainHandler.VerifyDomain)

			// 用户管理
			apiAdmin.GET("/users", userHandler.GetUsers)
			apiAdmin.GET("/users/:id", userHandler.GetUserByID)
			apiAdmin.GET("/users/:id/mailboxes", userHandler.GetUserMailboxes)
			apiAdmin.PUT("/users/:id/status", userHandler.UpdateUserStatus)
			apiAdmin.DELETE("/users/:id", userHandler.DeleteUser)

			// 邮箱管理
			apiAdmin.GET("/mailboxes", mailboxHandler.GetAllMailboxes)
			apiAdmin.PUT("/mailboxes/:id/status", mailboxHandler.UpdateMailboxStatus)
			apiAdmin.DELETE("/mailboxes/:id", mailboxHandler.DeleteMailboxAdmin)
			apiAdmin.GET("/mailboxes/:id/stats", mailboxHandler.GetMailboxStats)
		}

		// 公共API
		api.GET("/domains/available", domainHandler.GetAvailableDomains)
		api.GET("/domains/dns", domainHandler.GetDomainDNSRecords)
	}
}

func (s *Server) Start() error {
	return s.router.Run(":" + s.config.WebPort)
}
