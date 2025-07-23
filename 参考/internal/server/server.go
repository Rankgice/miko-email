package server

import (
	"database/sql"
	"embed"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"nbemail/internal/config"
)

// Server Web服务器
type Server struct {
	config        *config.Config
	db            *sql.DB
	router        *gin.Engine
	staticFiles   embed.FS
	templateFiles embed.FS
	server        *http.Server
}

// NewServer 创建新的Web服务器
func NewServer(cfg *config.Config, db *sql.DB, staticFiles, templateFiles embed.FS) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// 设置UTF-8字符集中间件
	router.Use(func(c *gin.Context) {
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.Next()
	})

	s := &Server{
		config:        cfg,
		db:            db,
		router:        router,
		staticFiles:   staticFiles,
		templateFiles: templateFiles,
	}

	s.setupRoutes()
	return s
}

// Start 启动Web服务器
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.config.WebPort)
	s.server = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}
	return s.server.ListenAndServe()
}

// Stop 停止Web服务器
func (s *Server) Stop() {
	if s.server != nil {
		s.server.Close()
	}
}

// setupRoutes 设置路由
func (s *Server) setupRoutes() {
	// 静态文件
	s.router.StaticFS("/static", http.FS(s.staticFiles))

	// 首页
	s.router.GET("/", s.handleIndex)
	s.router.GET("/login", s.handleLoginPage)
	s.router.POST("/api/login", s.handleLogin)
	s.router.POST("/api/logout", s.handleLogout)

	// API路由组
	api := s.router.Group("/api")
	api.Use(s.authMiddleware())
	{
		// 邮件相关
		api.GET("/emails", s.handleGetEmails)
		api.GET("/emails/:id", s.handleGetEmail)
		api.PUT("/emails/:id/read", s.handleMarkEmailAsRead)
		api.DELETE("/emails/:id", s.handleDeleteEmail)
		api.POST("/emails/send", s.handleSendEmail)
		api.POST("/parse-mime", s.handleParseMIME)

		// 邮箱管理相关
		api.GET("/mailboxes", s.handleGetMailboxes)
		api.GET("/mailboxes/credentials", s.handleGetMailboxCredentials)
		api.POST("/mailboxes/generate", s.handleGenerateMailboxes)
		api.POST("/mailboxes/switch", s.handleSwitchMailbox)
		api.DELETE("/mailboxes/:id", s.handleDeleteMailbox)

		// 用户相关
		api.GET("/users", s.handleGetUsers)
		api.POST("/users", s.handleCreateUser)
		api.PUT("/users/:id", s.handleUpdateUser)
		api.DELETE("/users/:id", s.handleDeleteUser)
		api.POST("/users/:id/assign-mailboxes", s.handleAssignMailboxesToUser)
		api.POST("/users/:id/assign-domains", s.handleAssignDomainsToUser)
		api.POST("/users/:id/reclaim-domains", s.handleReclaimDomainsFromUser)

		// 域名相关
		api.GET("/domains", s.handleGetDomains)
		api.GET("/user/domains", s.handleGetUserDomains)
		api.POST("/domains", s.handleCreateDomain)
		api.DELETE("/domains/:id", s.handleDeleteDomain)
		api.POST("/domains/fix-ownership", s.handleFixDomainOwnership)

		// DNS验证相关
		api.POST("/domains/:id/verify", s.handleVerifyDomain)
		api.GET("/domains/:id/dns-instructions", s.handleGetDNSInstructions)
		api.GET("/domains/:id/dns-propagation", s.handleCheckDNSPropagation)
		api.POST("/domains/batch-verify", s.handleBatchVerifyDomains)

		// SMTP配置管理
		api.GET("/smtp-configs", s.handleGetSMTPConfigs)
		api.POST("/smtp-configs", s.handleAddSMTPConfig)
		api.POST("/smtp-configs/auto-config", s.handleAutoConfigSMTP)
		api.GET("/smtp-configs/:domain/verify", s.handleVerifySMTPConfig)
		api.DELETE("/smtp-configs/:domain", s.handleDeleteSMTPConfig)
	}

	// Web页面路由
	web := s.router.Group("/")
	web.Use(s.authMiddleware())
	{
		web.GET("/inbox", s.handleInboxPage)
		web.GET("/sent", s.handleSentPage)
		web.GET("/compose", s.handleComposePage)
		web.GET("/email/:id", s.handleEmailDetailPage)
		web.GET("/users", s.handleUsersPage)
		web.GET("/domains", s.handleDomainsPage)
		web.GET("/smtp-configs", s.handleSMTPConfigsPage)
		web.GET("/guide", s.handleGuidePage)
	}
}

// handleIndex 处理首页
func (s *Server) handleIndex(c *gin.Context) {
	// 检查是否已登录
	if token, err := c.Cookie("token"); err == nil && token != "" {
		c.Redirect(http.StatusFound, "/inbox")
		return
	}

	// 显示美化的首页
	tmpl := `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>NBEmail - 临时邮箱·随心所欲</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            color: #333;
        }

        /* 导航栏 */
        .navbar {
            position: fixed;
            top: 0;
            width: 100%;
            background: rgba(255, 255, 255, 0.95);
            backdrop-filter: blur(10px);
            padding: 15px 0;
            z-index: 1000;
            box-shadow: 0 2px 20px rgba(0,0,0,0.1);
        }
        .nav-container {
            max-width: 1200px;
            margin: 0 auto;
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 0 20px;
        }
        .logo {
            font-size: 24px;
            font-weight: 600;
            color: #333;
        }
        .nav-links {
            display: flex;
            gap: 30px;
        }
        .nav-links a {
            text-decoration: none;
            color: #666;
            font-weight: 500;
            transition: color 0.3s;
        }
        .nav-links a:hover {
            color: #007bff;
        }
        .login-btn {
            background: #007bff;
            color: white;
            padding: 10px 20px;
            border-radius: 25px;
            text-decoration: none;
            font-weight: 500;
            transition: all 0.3s;
        }
        .login-btn:hover {
            background: #0056b3;
            transform: translateY(-2px);
        }

        /* 主要内容 */
        .hero {
            padding: 120px 20px 80px;
            text-align: center;
            color: white;
        }
        .hero h1 {
            font-size: 4rem;
            font-weight: 300;
            margin-bottom: 20px;
            letter-spacing: -2px;
        }
        .hero .brand {
            font-weight: 700;
            color: #fff;
        }
        .hero p {
            font-size: 1.5rem;
            margin-bottom: 40px;
            opacity: 0.9;
        }
        .hero-features {
            display: flex;
            justify-content: center;
            gap: 40px;
            margin-bottom: 50px;
            flex-wrap: wrap;
        }
        .hero-feature {
            color: rgba(255,255,255,0.8);
            font-size: 1.1rem;
        }
        .start-btn {
            background: #000;
            color: white;
            padding: 15px 40px;
            border-radius: 30px;
            text-decoration: none;
            font-size: 1.1rem;
            font-weight: 500;
            display: inline-block;
            transition: all 0.3s;
        }
        .start-btn:hover {
            background: #333;
            transform: translateY(-3px);
            box-shadow: 0 10px 30px rgba(0,0,0,0.3);
        }

        /* 特性部分 */
        .features {
            background: white;
            padding: 100px 20px;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
        }
        .section-title {
            text-align: center;
            margin-bottom: 80px;
        }
        .section-title h2 {
            font-size: 2.5rem;
            font-weight: 300;
            margin-bottom: 20px;
            color: #333;
        }
        .section-title p {
            font-size: 1.2rem;
            color: #666;
        }
        .features-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
            gap: 50px;
        }
        .feature-card {
            text-align: center;
            padding: 40px 30px;
            border-radius: 15px;
            transition: all 0.3s;
        }
        .feature-card:hover {
            transform: translateY(-10px);
            box-shadow: 0 20px 40px rgba(0,0,0,0.1);
        }
        .feature-icon {
            font-size: 3rem;
            margin-bottom: 30px;
            display: block;
        }
        .feature-card h3 {
            font-size: 1.5rem;
            margin-bottom: 20px;
            color: #333;
        }
        .feature-card p {
            color: #666;
            line-height: 1.6;
        }

        /* 对比部分 */
        .comparison {
            background: #f8f9fa;
            padding: 100px 20px;
        }
        .comparison-grid {
            display: grid;
            grid-template-columns: 1fr auto 1fr;
            gap: 50px;
            align-items: center;
            max-width: 1000px;
            margin: 0 auto;
        }
        .comparison-side {
            background: white;
            padding: 40px;
            border-radius: 15px;
            box-shadow: 0 5px 20px rgba(0,0,0,0.1);
        }
        .comparison-side h3 {
            margin-bottom: 30px;
            font-size: 1.5rem;
        }
        .comparison-side.traditional h3 {
            color: #dc3545;
        }
        .comparison-side.nbemail h3 {
            color: #28a745;
        }
        .comparison-item {
            display: flex;
            align-items: center;
            margin-bottom: 15px;
            padding: 10px 0;
        }
        .comparison-item::before {
            content: "×";
            margin-right: 15px;
            font-size: 1.2rem;
            font-weight: bold;
        }
        .traditional .comparison-item::before {
            color: #dc3545;
        }
        .nbemail .comparison-item::before {
            content: "✓";
            color: #28a745;
        }
        .vs-circle {
            width: 60px;
            height: 60px;
            background: #007bff;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            color: white;
            font-weight: bold;
            font-size: 1.2rem;
        }

        /* 响应式 */
        @media (max-width: 768px) {
            .hero h1 { font-size: 2.5rem; }
            .hero p { font-size: 1.2rem; }
            .hero-features { flex-direction: column; gap: 20px; }
            .features-grid { grid-template-columns: 1fr; }
            .comparison-grid {
                grid-template-columns: 1fr;
                gap: 30px;
            }
            .vs-circle { margin: 0 auto; }
        }
    </style>
</head>
<body>
    <!-- 导航栏 -->
    <nav class="navbar">
        <div class="nav-container">
            <div class="logo">NBEmail</div>
            <div class="nav-links">
                <a href="#features">功能特性</a>
                <a href="#comparison">产品对比</a>
                <a href="#about">关于我们</a>
            </div>
            <a href="/login" class="login-btn">开始使用</a>
        </div>
    </nav>

    <!-- 主要内容区 -->
    <section class="hero">
        <h1><span class="brand">NB</span>Email</h1>
        <p>临时邮箱 · 随心所欲</p>
        <div class="hero-features">
            <div class="hero-feature">无需注册</div>
            <div class="hero-feature">即时收发</div>
            <div class="hero-feature">隐私保护</div>
        </div>
        <a href="/login" class="start-btn">开始使用</a>
    </section>

    <!-- 特性部分 -->
    <section class="features" id="features">
        <div class="container">
            <div class="section-title">
                <h2>现代化技术 · 极简部署</h2>
                <p>基于最新技术栈构建，为您提供最佳的邮件服务体验</p>
            </div>
            <div class="features-grid">
                <div class="feature-card">
                    <span class="feature-icon">∞</span>
                    <h3>Go + Vue 技术栈</h3>
                    <p>采用 Go 语言 + Vue 前端的现代化技术栈，高性能，易维护</p>
                </div>
                <div class="feature-card">
                    <span class="feature-icon">🚀</span>
                    <h3>一键部署</h3>
                    <p>极一键部署到各种环境，无需复杂的配置，一键即可开始使用</p>
                </div>
                <div class="feature-card">
                    <span class="feature-icon">🛡️</span>
                    <h3>内置邮件服务</h3>
                    <p>内置完整的邮件服务功能，支持邮件收发及完整的邮件管理功能</p>
                </div>
            </div>
        </div>
    </section>

    <!-- 对比部分 -->
    <section class="comparison" id="comparison">
        <div class="container">
            <div class="section-title">
                <h2>为什么选择 NbEmail?</h2>
            </div>
            <div class="comparison-grid">
                <div class="comparison-side traditional">
                    <h3>🔴 传统方案</h3>
                    <div class="comparison-item">传统邮箱需要复杂注册</div>
                    <div class="comparison-item">数据隐私难以保障</div>
                    <div class="comparison-item">需要绑定第三方邮件服务</div>
                    <div class="comparison-item">在线时间长，性能受限</div>
                </div>
                <div class="vs-circle">VS</div>
                <div class="comparison-side nbemail">
                    <h3>🟢 NbEmail</h3>
                    <div class="comparison-item">基于现代技术</div>
                    <div class="comparison-item">一键部署，即时使用</div>
                    <div class="comparison-item">内置完整的邮件服务功能</div>
                    <div class="comparison-item">Go+Vue现代技术栈，高性能</div>
                </div>
            </div>
        </div>
    </section>

    <script>
        // 平滑滚动
        document.querySelectorAll('a[href^="#"]').forEach(anchor => {
            anchor.addEventListener('click', function (e) {
                e.preventDefault();
                document.querySelector(this.getAttribute('href')).scrollIntoView({
                    behavior: 'smooth'
                });
            });
        });
    </script>
</body>
</html>`
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, tmpl)
}

// handleLoginPage 处理登录页面
func (s *Server) handleLoginPage(c *gin.Context) {
	tmpl := `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>NBEmail</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            background: #f5f5f5;
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        .container {
            max-width: 400px;
            width: 90%;
            padding: 40px;
            background: white;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            text-align: center;
        }
        .logo {
            margin-bottom: 30px;
        }
        .logo h1 {
            color: #333;
            font-size: 28px;
            font-weight: 400;
            margin-bottom: 8px;
        }
        .logo p {
            color: #666;
            font-size: 14px;
            line-height: 1.4;
        }
        .form-group {
            margin-bottom: 20px;
            text-align: left;
        }
        label {
            display: block;
            margin-bottom: 6px;
            color: #333;
            font-weight: 500;
            font-size: 14px;
        }
        input[type="email"], input[type="password"] {
            width: 100%;
            padding: 12px 16px;
            border: 1px solid #ddd;
            border-radius: 4px;
            font-size: 14px;
            transition: border-color 0.2s;
            background: white;
        }
        input[type="email"]:focus, input[type="password"]:focus {
            outline: none;
            border-color: #007bff;
            box-shadow: 0 0 0 2px rgba(0,123,255,0.25);
        }
        .btn {
            width: 100%;
            padding: 12px;
            background: #007bff;
            color: white;
            border: none;
            border-radius: 4px;
            font-size: 14px;
            font-weight: 500;
            cursor: pointer;
            transition: background-color 0.2s;
            margin-top: 10px;
        }
        .btn:hover {
            background: #0056b3;
        }
        .error {
            color: #dc3545;
            margin-top: 15px;
            font-size: 14px;
            text-align: center;
            padding: 10px;
            background: rgba(220,53,69,0.1);
            border-radius: 8px;
        }
        .remember {
            margin: 15px 0;
            display: flex;
            align-items: center;
            font-size: 14px;
        }
        .remember input {
            margin-right: 8px;
        }
        .remember label {
            margin-bottom: 0;
            font-size: 14px;
            color: #666;
        }


        /* 响应式 */
        @media (max-width: 480px) {
            .container {
                padding: 40px 30px;
            }
            .logo h1 {
                font-size: 2rem;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="logo">
            <h1>NBEmail</h1>
            <p>欢迎使用邮件系统，QQ群号: 819229551</p>
        </div>
        <form id="loginForm">
            <div class="form-group">
                <input type="email" id="email" name="email" required placeholder="邮箱">
            </div>
            <div class="form-group">
                <input type="password" id="password" name="password" required placeholder="密码">
            </div>
            <div class="remember">
                <input type="checkbox" id="remember">
                <label for="remember">记住我</label>
            </div>
            <button type="submit" class="btn">登录</button>
            <div id="error" class="error" style="display: none;"></div>
        </form>
    </div>
    <script>
        document.getElementById('loginForm').addEventListener('submit', async (e) => {
            e.preventDefault();
            const email = document.getElementById('email').value;
            const password = document.getElementById('password').value;
            const errorDiv = document.getElementById('error');

            try {
                const response = await fetch('/api/login', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ email, password })
                });

                const result = await response.json();

                if (result.success) {
                    window.location.href = '/inbox';
                } else {
                    errorDiv.textContent = result.message || '登录失败，请检查邮箱和密码';
                    errorDiv.style.display = 'block';
                }
            } catch (error) {
                errorDiv.textContent = '登录失败，请重试';
                errorDiv.style.display = 'block';
            }
        });
    </script>
</body>
</html>`
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, tmpl)
}