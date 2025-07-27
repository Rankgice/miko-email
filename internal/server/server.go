package server

import (
	"database/sql"
	"miko-email/internal/config"
	"miko-email/internal/handlers"
	"miko-email/internal/services/email"
	"miko-email/internal/services/forward"
	"miko-email/internal/svc"
	"net/http"
	"strings"

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

func New(db *sql.DB, cfg *config.Config, svcCtx *svc.ServiceContext) *Server {
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
	emailService := email.NewService(svcCtx)

	// 创建转发服务
	forwardService := forward.NewService(svcCtx)

	server := &Server{
		router:         gin.Default(),
		db:             db,
		config:         cfg,
		sessionStore:   sessionStore,
		emailService:   emailService,
		forwardService: forwardService,
	}

	server.setupRoutes(svcCtx)
	return server
}

func (s *Server) setupRoutes(svcCtx *svc.ServiceContext) {
	// 设置UTF-8编码中间件（仅对API请求设置JSON Content-Type）
	s.router.Use(func(c *gin.Context) {
		// 只对API请求设置JSON编码，让静态文件使用默认的MIME类型
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.Header("Content-Type", "application/json; charset=utf-8")
		}
		c.Next()
	})

	// 使用新的路由系统
	handlers.RegisterRoutes(s.router, s.sessionStore, svcCtx)
}

func (s *Server) Start() error {
	return s.router.Run(":" + s.config.WebPort)
}
