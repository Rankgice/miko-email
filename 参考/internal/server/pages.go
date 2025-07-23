package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// handleInboxPage 处理收件箱页面
func (s *Server) handleInboxPage(c *gin.Context) {
	userEmail := c.GetString("user_email")
	isAdmin := c.GetBool("is_admin")

	tmpl := s.generateInboxTemplate(userEmail, isAdmin, "inbox")
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, tmpl)
}

// handleSentPage 处理发件箱页面
func (s *Server) handleSentPage(c *gin.Context) {
	userEmail := c.GetString("user_email")
	isAdmin := c.GetBool("is_admin")

	tmpl := s.generateEmailPageTemplate(userEmail, isAdmin, "sent", "📤 已发送")
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, tmpl)
}

// handleComposePage 处理写邮件页面
func (s *Server) handleComposePage(c *gin.Context) {
	userEmail := c.GetString("user_email")
	isAdmin := c.GetBool("is_admin")

	tmpl := s.generateComposePageTemplate(userEmail, isAdmin)
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, tmpl)
}

// handleUsersPage 处理用户管理页面
func (s *Server) handleUsersPage(c *gin.Context) {
	if !c.GetBool("is_admin") {
		c.Redirect(http.StatusFound, "/inbox")
		return
	}

	userEmail := c.GetString("user_email")
	tmpl := s.generateUsersPageTemplate(userEmail)
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, tmpl)
}

// handleDomainsPage 处理域名管理页面
func (s *Server) handleDomainsPage(c *gin.Context) {
	if !c.GetBool("is_admin") {
		c.Redirect(http.StatusFound, "/inbox")
		return
	}

	userEmail := c.GetString("user_email")
	tmpl := s.generateDomainsPageTemplate(userEmail)
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, tmpl)
}

// handleGuidePage 处理使用指南页面
func (s *Server) handleGuidePage(c *gin.Context) {
	userEmail := c.GetString("user_email")
	isAdmin := c.GetBool("is_admin")

	tmpl := s.generateGuidePageTemplate(userEmail, isAdmin)
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, tmpl)
}

// handleSMTPConfigsPage 处理SMTP配置管理页面
func (s *Server) handleSMTPConfigsPage(c *gin.Context) {
	userEmail := c.GetString("user_email")
	isAdmin := c.GetBool("is_admin")

	if !isAdmin {
		c.Redirect(http.StatusFound, "/inbox")
		return
	}

	tmpl := s.generateSMTPConfigsPageTemplate(userEmail, isAdmin)
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, tmpl)
}

// handleEmailDetailPage 处理邮件详情页面
func (s *Server) handleEmailDetailPage(c *gin.Context) {
	userEmail := c.GetString("user_email")
	isAdmin := c.GetBool("is_admin")
	emailID := c.Param("id")

	tmpl := s.generateEmailDetailPageTemplate(userEmail, isAdmin, emailID)
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, tmpl)
}

// generateEmailPageTemplate 生成邮件页面模板
func (s *Server) generateEmailPageTemplate(userEmail string, isAdmin bool, folder, title string) string {
	// 复用收件箱的模板，但修改标题和API调用
	tmpl := s.getBasePageTemplate(userEmail, isAdmin, folder)
	// 替换标题和API调用
	if folder == "sent" {
		tmpl = strings.Replace(tmpl, "folder=inbox", "folder=sent", -1)
		tmpl = strings.Replace(tmpl, "正在加载邮件...", "正在加载已发送邮件...", -1)
		tmpl = strings.Replace(tmpl, "暂无邮件", "暂无已发送邮件", -1)
	}
	return tmpl
}

// generateComposePageTemplate 生成写邮件页面模板
func (s *Server) generateComposePageTemplate(userEmail string, isAdmin bool) string {
	return `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>NBEmail - 写邮件</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            background: #f8f9fa;
            color: #333;
        }

        /* 顶部导航栏 */
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 0 30px;
            height: 70px;
            display: flex;
            align-items: center;
            justify-content: space-between;
            box-shadow: 0 2px 20px rgba(0,0,0,0.1);
        }
        .logo {
            font-size: 28px;
            font-weight: 300;
            color: white;
            letter-spacing: -1px;
        }
        .header-nav {
            display: flex;
            align-items: center;
            gap: 30px;
        }
        .header-nav a {
            color: rgba(255,255,255,0.8);
            text-decoration: none;
            font-weight: 500;
            transition: color 0.3s;
        }
        .header-nav a:hover {
            color: white;
        }
        .user-info {
            display: flex;
            align-items: center;
            gap: 20px;
        }
        .user-email {
            color: rgba(255,255,255,0.9);
            font-weight: 500;
        }
        .logout-btn {
            background: rgba(255,255,255,0.2);
            color: white;
            border: 1px solid rgba(255,255,255,0.3);
            padding: 8px 20px;
            border-radius: 20px;
            cursor: pointer;
            font-weight: 500;
            transition: all 0.3s;
        }
        .logout-btn:hover {
            background: rgba(255,255,255,0.3);
            transform: translateY(-1px);
        }

        /* 品牌区域 */
        .brand-section {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            padding: 20px;
            display: flex;
            align-items: center;
            justify-content: space-between;
            color: white;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }

        .brand-left {
            display: flex;
            align-items: center;
            gap: 20px;
        }

        .code-mascot {
            font-size: 3rem;
            animation: float 3s ease-in-out infinite;
        }

        @keyframes float {
            0%, 100% { transform: translateY(0px); }
            50% { transform: translateY(-10px); }
        }

        .brand-info h1 {
            margin: 0;
            font-size: 2.5rem;
            font-weight: 700;
            text-shadow: 2px 2px 4px rgba(0,0,0,0.3);
        }

        .brand-info .version {
            font-size: 1rem;
            opacity: 0.9;
            margin-top: 5px;
        }

        .brand-right {
            display: flex;
            align-items: center;
            gap: 15px;
            font-size: 1.1rem;
        }

        /* 主要内容区域 */
        .container {
            display: flex;
            height: calc(100vh - 140px);
            max-width: 1600px;
            margin: 0 auto;
            gap: 20px;
            padding: 20px;
        }

        /* 侧边栏 */
        .sidebar {
            width: 250px;
            background: white;
            border-radius: 15px;
            padding: 20px 0;
            box-shadow: 0 5px 20px rgba(0,0,0,0.08);
            height: fit-content;
        }
        .nav-section {
            margin-bottom: 30px;
        }
        .nav-section-title {
            padding: 0 25px 15px;
            font-size: 12px;
            font-weight: 600;
            color: #8e9aaf;
            text-transform: uppercase;
            letter-spacing: 1px;
        }
        .nav-item {
            display: flex;
            align-items: center;
            padding: 15px 25px;
            color: #5a6c7d;
            text-decoration: none;
            transition: all 0.3s;
            font-weight: 500;
        }
        .nav-item:hover {
            background: #f8f9fa;
            color: #007bff;
            transform: translateX(5px);
        }
        .nav-item.active {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            box-shadow: 0 5px 15px rgba(102,126,234,0.3);
        }
        .nav-item .icon {
            margin-right: 12px;
            font-size: 18px;
            width: 20px;
        }
        .nav-divider {
            height: 1px;
            background: #e9ecef;
            margin: 20px 25px;
        }

        /* 主内容区 */
        .main {
            flex: 1;
            display: flex;
            flex-direction: column;
        }

        /* 写邮件表单 */
        .compose-form {
            background: white;
            border-radius: 15px;
            padding: 40px;
            box-shadow: 0 5px 20px rgba(0,0,0,0.08);
            height: 100%;
            display: flex;
            flex-direction: column;
        }
        .compose-header {
            margin-bottom: 40px;
            padding-bottom: 20px;
            border-bottom: 2px solid #f1f3f4;
        }
        .compose-header h2 {
            color: #333;
            font-weight: 600;
            font-size: 1.8rem;
            margin: 0;
        }
        .form-group {
            margin-bottom: 25px;
        }
        label {
            display: block;
            margin-bottom: 8px;
            font-weight: 600;
            color: #333;
            font-size: 14px;
        }
        input[type="email"], input[type="text"] {
            width: 100%;
            padding: 15px 20px;
            border: 2px solid #e9ecef;
            border-radius: 10px;
            font-size: 16px;
            transition: all 0.3s;
            background: #f8f9fa;
        }
        input[type="email"]:focus, input[type="text"]:focus {
            outline: none;
            border-color: #667eea;
            background: white;
            box-shadow: 0 0 0 3px rgba(102,126,234,0.1);
        }
        textarea {
            width: 100%;
            padding: 20px;
            border: 2px solid #e9ecef;
            border-radius: 10px;
            font-size: 16px;
            min-height: 300px;
            resize: vertical;
            font-family: inherit;
            line-height: 1.6;
            transition: all 0.3s;
            background: #f8f9fa;
        }
        textarea:focus {
            outline: none;
            border-color: #667eea;
            background: white;
            box-shadow: 0 0 0 3px rgba(102,126,234,0.1);
        }

        /* 按钮样式 */
        .form-actions {
            display: flex;
            gap: 15px;
            margin-top: auto;
            padding-top: 30px;
            border-top: 2px solid #f1f3f4;
        }
        .btn {
            padding: 15px 30px;
            border: none;
            border-radius: 10px;
            cursor: pointer;
            font-weight: 600;
            font-size: 16px;
            transition: all 0.3s;
            display: flex;
            align-items: center;
            gap: 8px;
        }
        .btn-primary {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
        }
        .btn-primary:hover {
            transform: translateY(-2px);
            box-shadow: 0 10px 25px rgba(102,126,234,0.3);
        }
        .btn-secondary {
            background: #f8f9fa;
            color: #5a6c7d;
            border: 2px solid #e9ecef;
        }
        .btn-secondary:hover {
            background: #e9ecef;
            transform: translateY(-2px);
        }

        /* 响应式设计 */
        @media (max-width: 1024px) {
            .brand-section {
                flex-direction: column;
                gap: 15px;
                text-align: center;
            }
            .brand-left {
                justify-content: center;
            }
            .brand-info h1 {
                font-size: 2rem;
            }
            .container {
                flex-direction: column;
                padding: 10px;
            }
            .sidebar {
                width: 100%;
                order: 2;
            }
            .main {
                order: 1;
            }
        }

        @media (max-width: 768px) {
            .header {
                padding: 0 15px;
            }
            .header-nav {
                display: none;
            }
            .brand-section {
                padding: 15px;
            }
            .code-mascot {
                font-size: 2rem;
            }
            .brand-info h1 {
                font-size: 1.8rem;
            }
            .brand-info .version {
                font-size: 0.9rem;
            }
            .brand-right {
                font-size: 1rem;
            }
            .sidebar {
                padding: 15px 0;
                margin-bottom: 20px;
            }
            .nav-section {
                margin-bottom: 20px;
            }
            .nav-section-title {
                padding: 0 20px 10px;
                font-size: 11px;
            }
            .nav-item {
                padding: 12px 20px;
                font-size: 14px;
            }
            .nav-item .icon {
                margin-right: 8px;
                font-size: 16px;
            }
            .main {
                order: 1;
                min-height: 60vh;
                width: 100%;
            }
            .compose-form {
                padding: 15px;
            }

            /* 移动端按钮样式优化 */
            .form-actions {
                flex-direction: column;
                gap: 12px;
                position: sticky;
                bottom: 0;
                background: white;
                padding: 20px 15px;
                margin: 20px -15px -15px -15px;
                border-top: 2px solid #f1f3f4;
                box-shadow: 0 -2px 10px rgba(0,0,0,0.1);
            }

            .btn {
                width: 100%;
                padding: 18px 20px;
                font-size: 16px;
                font-weight: 600;
                justify-content: center;
            }
        }
    </style>
</head>
<body>
    <div class="header">
        <div class="logo">📧 NBEmail</div>
        <div class="header-nav">
            <a href="/guide">使用指南</a>
            <a href="#about">关于</a>
        </div>
        <div class="user-info">
            <span class="user-email">` + userEmail + `</span>
            <button class="logout-btn" onclick="logout()">退出登录</button>
        </div>
    </div>

    <!-- 品牌区域 -->
    <div class="brand-section">
        <div class="brand-left">
            <div class="code-mascot">👩‍💻</div>
            <div class="brand-info">
                <h1>📧 NBEmail</h1>
                <div class="version">v1.0.0 - 专业邮件管理系统</div>
            </div>
        </div>
        <div class="brand-right">
            <span>欢迎使用，` + userEmail + `</span>
        </div>
    </div>

    <div class="container">
        <div class="sidebar">
            <div class="nav-section">
                <div class="nav-section-title">邮箱</div>
                <a href="/inbox" class="nav-item">
                    <span class="icon">📥</span>收件箱
                </a>
                <a href="/sent" class="nav-item">
                    <span class="icon">📤</span>已发送
                </a>
                <a href="/compose" class="nav-item active">
                    <span class="icon">✏️</span>写邮件
                </a>
            </div>` +
		func() string {
			if isAdmin {
				return `
            <div class="nav-divider"></div>
            <div class="nav-section">
                <div class="nav-section-title">管理</div>
                <a href="/users" class="nav-item">
                    <span class="icon">👥</span>用户管理
                </a>
                <a href="/domains" class="nav-item">
                    <span class="icon">🌐</span>域名管理
                </a>
                <a href="/smtp-configs" class="nav-item">
                    <span class="icon">📮</span>SMTP配置
                </a>
            </div>`
			}
			return ""
		}() + `
        </div>
        <div class="main">
            <div class="compose-form">
                <div class="compose-header">
                    <h2>✏️ 写邮件</h2>
                </div>
                <form id="composeForm">
                    <div class="form-group">
                        <label>发件人地址</label>
                        <div class="from-display" id="fromDisplay" style="padding: 12px 15px; background: #f8f9fa; border: 2px solid #e9ecef; border-radius: 8px; color: #666; display: flex; align-items: center; gap: 10px;">
                            <span>📧</span>
                            <span id="fromEmail">加载中...</span>
                            <button type="button" onclick="showMailboxSelector()" style="background: none; border: none; color: #007bff; cursor: pointer; padding: 2px 6px; border-radius: 4px; font-size: 12px;" title="切换发件邮箱">切换</button>
                        </div>
                    </div>
                    <div class="form-group">
                        <label for="to">* 收件人地址</label>
                        <input type="email" id="to" name="to" required placeholder="请输入收件人邮箱地址，例如：user@example.com">
                    </div>
                    <div class="form-group">
                        <label for="subject">* 邮件主题</label>
                        <input type="text" id="subject" name="subject" required placeholder="请输入邮件主题">
                    </div>
                    <div class="form-group">
                        <label for="body">* 邮件内容</label>
                        <textarea id="body" name="body" required placeholder="请输入您要发送的邮件内容..."></textarea>
                    </div>
                    <div class="form-actions">
                        <button type="submit" class="btn btn-primary">
                            <span>📤</span>发送邮件
                        </button>
                        <button type="button" class="btn btn-secondary" onclick="clearForm()">
                            <span>🗑️</span>清空内容
                        </button>
                    </div>
                </form>
            </div>
        </div>
    </div>

    <!-- 邮箱选择器模态框 -->
    <div id="mailboxSelectorModal" class="modal" style="display: none; position: fixed; z-index: 1000; left: 0; top: 0; width: 100%; height: 100%; background-color: rgba(0,0,0,0.5); backdrop-filter: blur(5px);">
        <div class="modal-content" style="background-color: white; margin: 15% auto; padding: 30px; border-radius: 15px; width: 90%; max-width: 600px; box-shadow: 0 20px 40px rgba(0,0,0,0.2);">
            <div class="modal-header" style="margin-bottom: 25px; padding-bottom: 15px; border-bottom: 2px solid #f1f3f4; display: flex; justify-content: space-between; align-items: center;">
                <h3 style="color: #333; font-weight: 600; margin: 0;">📮 选择发件邮箱</h3>
                <span class="close" onclick="closeMailboxSelector()" style="color: #aaa; font-size: 28px; font-weight: bold; cursor: pointer; line-height: 1;">&times;</span>
            </div>

            <div class="mailbox-selector-list" id="mailboxSelectorList">
                <div style="text-align: center; padding: 40px; color: #666;">正在加载邮箱列表...</div>
            </div>
        </div>
    </div>

    <script>
        let currentFromEmail = '';

        // 加载当前发件邮箱
        async function loadCurrentFromEmail() {
            try {
                const response = await fetch('/api/mailboxes');
                const result = await response.json();
                if (result.success && result.data.length > 0) {
                    const currentMailbox = result.data.find(m => m.is_current);
                    if (currentMailbox) {
                        currentFromEmail = currentMailbox.email;
                        document.getElementById('fromEmail').textContent = currentFromEmail;
                    } else {
                        document.getElementById('fromEmail').textContent = '未设置发件邮箱';
                    }
                } else {
                    document.getElementById('fromEmail').textContent = '无可用邮箱';
                }
            } catch (error) {
                document.getElementById('fromEmail').textContent = '加载失败';
            }
        }

        // 显示邮箱选择器
        async function showMailboxSelector() {
            document.getElementById('mailboxSelectorModal').style.display = 'block';
            await loadMailboxesForSelector();
        }

        // 关闭邮箱选择器
        function closeMailboxSelector() {
            document.getElementById('mailboxSelectorModal').style.display = 'none';
        }

        // 加载邮箱列表用于选择器
        async function loadMailboxesForSelector() {
            try {
                const response = await fetch('/api/mailboxes');
                const result = await response.json();
                if (result.success) {
                    const mailboxes = result.data || [];
                    renderMailboxSelector(mailboxes);
                } else {
                    document.getElementById('mailboxSelectorList').innerHTML = '<div style="text-align: center; padding: 40px; color: #666;">加载邮箱失败</div>';
                }
            } catch (error) {
                document.getElementById('mailboxSelectorList').innerHTML = '<div style="text-align: center; padding: 40px; color: #666;">加载邮箱失败</div>';
            }
        }

        // 渲染邮箱选择器
        function renderMailboxSelector(mailboxes) {
            const selectorList = document.getElementById('mailboxSelectorList');

            if (mailboxes.length === 0) {
                selectorList.innerHTML = ` + "`" + `
                    <div style="text-align: center; padding: 40px; color: #666;">
                        <div style="font-size: 3rem; margin-bottom: 15px;">📭</div>
                        <h3>暂无邮箱</h3>
                        <p>请先在邮箱管理中生成一些邮箱</p>
                    </div>
                ` + "`" + `;
                return;
            }

            selectorList.innerHTML = ` + "`" + `
                <div style="display: grid; gap: 10px;">
                    ${mailboxes.map(mailbox => ` + "`" + `
                        <div onclick="selectFromEmail('${mailbox.email}')" style="display: flex; align-items: center; justify-content: space-between; padding: 15px; border: 2px solid ${mailbox.is_current ? '#4caf50' : '#e9ecef'}; border-radius: 10px; background: ${mailbox.is_current ? '#f1f8e9' : 'white'}; cursor: pointer; transition: all 0.3s;" onmouseover="this.style.borderColor='#007bff'" onmouseout="this.style.borderColor='${mailbox.is_current ? '#4caf50' : '#e9ecef'}'">
                            <div style="display: flex; align-items: center; gap: 12px;">
                                <div style="width: 35px; height: 35px; border-radius: 50%; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); display: flex; align-items: center; justify-content: center; color: white; font-weight: 600;">
                                    ${mailbox.email.charAt(0).toUpperCase()}
                                </div>
                                <div>
                                    <div style="font-weight: 600; color: #333;">${mailbox.email}</div>
                                    <div style="font-size: 12px; color: #666;">${mailbox.domain_name}</div>
                                </div>
                                ${mailbox.is_current ? '<span style="background: #4caf50; color: white; padding: 2px 8px; border-radius: 12px; font-size: 11px; margin-left: 10px;">当前</span>' : ''}
                            </div>
                            <div style="color: #007bff; font-size: 12px;">点击选择</div>
                        </div>
                    ` + "`" + `).join('')}
                </div>
            ` + "`" + `;
        }

        // 选择发件邮箱
        function selectFromEmail(email) {
            currentFromEmail = email;
            document.getElementById('fromEmail').textContent = email;
            closeMailboxSelector();
            showNotification('已选择发件邮箱: ' + email, 'success');
        }

        // 显示通知
        function showNotification(message, type = 'info') {
            const notification = document.createElement('div');
            notification.style.cssText = ` + "`" + `
                position: fixed;
                top: 20px;
                right: 20px;
                padding: 15px 25px;
                border-radius: 10px;
                color: white;
                font-weight: 500;
                z-index: 10000;
                transform: translateX(400px);
                transition: all 0.3s ease;
                box-shadow: 0 5px 20px rgba(0,0,0,0.2);
            ` + "`" + `;

            switch(type) {
                case 'success':
                    notification.style.background = 'linear-gradient(135deg, #28a745, #20c997)';
                    break;
                case 'error':
                    notification.style.background = 'linear-gradient(135deg, #dc3545, #e74c3c)';
                    break;
                case 'info':
                default:
                    notification.style.background = 'linear-gradient(135deg, #667eea, #764ba2)';
            }

            notification.textContent = message;
            document.body.appendChild(notification);

            setTimeout(() => {
                notification.style.transform = 'translateX(0)';
            }, 100);

            setTimeout(() => {
                notification.style.transform = 'translateX(400px)';
                setTimeout(() => {
                    if (document.body.contains(notification)) {
                        document.body.removeChild(notification);
                    }
                }, 300);
            }, 3000);
        }

        // 点击模态框外部关闭
        window.onclick = function(event) {
            const modal = document.getElementById('mailboxSelectorModal');
            if (event.target == modal) {
                closeMailboxSelector();
            }
        }

        // 页面加载时获取当前发件邮箱
        loadCurrentFromEmail();

        // 表单提交处理
        document.getElementById('composeForm').addEventListener('submit', async (e) => {
            e.preventDefault();

            const submitBtn = e.target.querySelector('button[type="submit"]');
            const originalText = submitBtn.innerHTML;

            // 显示发送中状态
            submitBtn.innerHTML = '<span>⏳</span>发送中...';
            submitBtn.disabled = true;

            const to = document.getElementById('to').value;
            const subject = document.getElementById('subject').value;
            const body = document.getElementById('body').value;

            // 检查是否选择了发件邮箱
            if (!currentFromEmail) {
                showNotification('请先选择发件邮箱', 'error');
                submitBtn.innerHTML = originalText;
                submitBtn.disabled = false;
                return;
            }

            try {
                const response = await fetch('/api/emails/send', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        to,
                        subject,
                        body,
                        from: currentFromEmail  // 添加发件人邮箱
                    })
                });

                const result = await response.json();

                if (result.success) {
                    // 显示成功状态
                    submitBtn.innerHTML = '<span>✅</span>发送成功';
                    submitBtn.style.background = '#28a745';

                    // 显示成功提示
                    showNotification('邮件发送成功！', 'success');

                    // 2秒后清空表单并恢复按钮
                    setTimeout(() => {
                        clearForm();
                        submitBtn.innerHTML = originalText;
                        submitBtn.style.background = '';
                        submitBtn.disabled = false;
                    }, 2000);
                } else {
                    throw new Error(result.message || '发送失败');
                }
            } catch (error) {
                // 显示错误状态
                submitBtn.innerHTML = '<span>❌</span>发送失败';
                submitBtn.style.background = '#dc3545';

                // 根据错误类型显示不同的错误信息
                let errorMessage = '发送失败';
                if (error.message) {
                    if (error.message.includes('SSL连接失败') || error.message.includes('连接SMTP服务器失败')) {
                        errorMessage = '网络连接失败，请检查网络设置或稍后重试';
                    } else if (error.message.includes('SMTP认证失败')) {
                        errorMessage = 'SMTP认证失败，请检查邮箱配置';
                    } else if (error.message.includes('未找到发件人') || error.message.includes('SMTP配置')) {
                        errorMessage = '邮箱配置错误，请联系管理员配置SMTP';
                    } else if (error.message.includes('TLS失败') || error.message.includes('启动TLS失败')) {
                        errorMessage = 'TLS加密连接失败，请检查服务器配置';
                    } else {
                        errorMessage = '发送失败: ' + error.message;
                    }
                }

                showNotification(errorMessage, 'error');

                // 5秒后恢复按钮（给用户更多时间阅读错误信息）
                setTimeout(() => {
                    submitBtn.innerHTML = originalText;
                    submitBtn.style.background = '';
                    submitBtn.disabled = false;
                }, 5000);
            }
        });

        // 清空表单
        function clearForm() {
            if (confirm('确定要清空所有内容吗？')) {
                document.getElementById('composeForm').reset();
                showNotification('表单已清空', 'info');
            }
        }

        // 显示通知
        function showNotification(message, type = 'info') {
            // 创建通知元素
            const notification = document.createElement('div');
            notification.style.cssText = ` + "`" + `
                position: fixed;
                top: 20px;
                right: 20px;
                padding: 15px 25px;
                border-radius: 10px;
                color: white;
                font-weight: 500;
                z-index: 10000;
                transform: translateX(400px);
                transition: all 0.3s ease;
                box-shadow: 0 5px 20px rgba(0,0,0,0.2);
            ` + "`" + `;

            // 设置不同类型的样式
            switch(type) {
                case 'success':
                    notification.style.background = 'linear-gradient(135deg, #28a745, #20c997)';
                    break;
                case 'error':
                    notification.style.background = 'linear-gradient(135deg, #dc3545, #e74c3c)';
                    break;
                case 'info':
                default:
                    notification.style.background = 'linear-gradient(135deg, #667eea, #764ba2)';
            }

            notification.textContent = message;
            document.body.appendChild(notification);

            // 显示动画
            setTimeout(() => {
                notification.style.transform = 'translateX(0)';
            }, 100);

            // 自动隐藏
            setTimeout(() => {
                notification.style.transform = 'translateX(400px)';
                setTimeout(() => {
                    document.body.removeChild(notification);
                }, 300);
            }, 3000);
        }

        // 退出登录
        async function logout() {
            if (confirm('确定要退出登录吗？')) {
                try {
                    await fetch('/api/logout', { method: 'POST' });
                    window.location.href = '/login';
                } catch (error) {
                    window.location.href = '/login';
                }
            }
        }

        // 自动保存草稿（可选功能）
        let draftTimer;
        function saveDraft() {
            const to = document.getElementById('to').value;
            const subject = document.getElementById('subject').value;
            const body = document.getElementById('body').value;

            if (to || subject || body) {
                localStorage.setItem('emailDraft', JSON.stringify({ to, subject, body }));
            }
        }

        // 恢复草稿
        function loadDraft() {
            const draft = localStorage.getItem('emailDraft');
            if (draft) {
                const { to, subject, body } = JSON.parse(draft);
                if (to) document.getElementById('to').value = to;
                if (subject) document.getElementById('subject').value = subject;
                if (body) document.getElementById('body').value = body;
            }
        }

        // 监听输入变化，自动保存草稿
        ['to', 'subject', 'body'].forEach(id => {
            document.getElementById(id).addEventListener('input', () => {
                clearTimeout(draftTimer);
                draftTimer = setTimeout(saveDraft, 1000);
            });
        });

        // 页面加载时恢复草稿
        window.addEventListener('load', loadDraft);

        // 发送成功后清除草稿
        function clearDraft() {
            localStorage.removeItem('emailDraft');
        }
    </script>
</body>
</html>`
}

// generateGuidePageTemplate 生成使用指南页面模板
func (s *Server) generateGuidePageTemplate(userEmail string, isAdmin bool) string {
	return `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>NBEmail - 使用指南</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            background: #f8f9fa;
            color: #333;
            line-height: 1.6;
        }

        /* 顶部导航栏 */
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 0 30px;
            height: 70px;
            display: flex;
            align-items: center;
            justify-content: space-between;
            box-shadow: 0 2px 20px rgba(0,0,0,0.1);
        }
        .logo {
            font-size: 28px;
            font-weight: 300;
            color: white;
            letter-spacing: -1px;
        }
        .user-info {
            display: flex;
            align-items: center;
            gap: 20px;
        }
        .user-email {
            color: rgba(255,255,255,0.9);
            font-weight: 500;
        }
        .logout-btn {
            background: rgba(255,255,255,0.2);
            color: white;
            border: 1px solid rgba(255,255,255,0.3);
            padding: 8px 20px;
            border-radius: 20px;
            cursor: pointer;
            font-weight: 500;
            transition: all 0.3s;
        }
        .logout-btn:hover {
            background: rgba(255,255,255,0.3);
            transform: translateY(-1px);
        }

        /* 品牌区域 */
        .brand-section {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            padding: 20px;
            display: flex;
            align-items: center;
            justify-content: space-between;
            color: white;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }

        .brand-left {
            display: flex;
            align-items: center;
            gap: 20px;
        }

        .code-mascot {
            font-size: 3rem;
            animation: float 3s ease-in-out infinite;
        }

        @keyframes float {
            0%, 100% { transform: translateY(0px); }
            50% { transform: translateY(-10px); }
        }

        .brand-info h1 {
            margin: 0;
            font-size: 2.5rem;
            font-weight: 700;
            text-shadow: 2px 2px 4px rgba(0,0,0,0.3);
        }

        .brand-info .version {
            font-size: 1rem;
            opacity: 0.9;
            margin-top: 5px;
        }

        .brand-right {
            display: flex;
            align-items: center;
            gap: 15px;
            font-size: 1.1rem;
        }

        /* 主要内容区域 */
        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 40px 20px;
        }

        /* 指南内容 */
        .guide-content {
            background: white;
            border-radius: 15px;
            padding: 50px;
            box-shadow: 0 5px 20px rgba(0,0,0,0.08);
        }

        .guide-header {
            text-align: center;
            margin-bottom: 50px;
            padding-bottom: 30px;
            border-bottom: 2px solid #f1f3f4;
        }

        .guide-header h1 {
            font-size: 2.5rem;
            font-weight: 300;
            color: #333;
            margin-bottom: 15px;
        }

        .guide-header p {
            font-size: 1.2rem;
            color: #666;
        }

        .guide-section {
            margin-bottom: 50px;
        }

        .guide-section h2 {
            font-size: 1.8rem;
            font-weight: 600;
            color: #333;
            margin-bottom: 25px;
            padding-left: 20px;
            border-left: 4px solid #667eea;
        }

        .guide-section h3 {
            font-size: 1.3rem;
            font-weight: 600;
            color: #333;
            margin: 25px 0 15px;
        }

        .guide-section p {
            margin-bottom: 15px;
            color: #555;
            font-size: 1rem;
        }

        .guide-section ul, .guide-section ol {
            margin: 20px 0;
            padding-left: 30px;
        }

        .guide-section li {
            margin-bottom: 10px;
            color: #555;
        }

        .config-section {
            background: #f8f9fa;
            padding: 30px;
            border-radius: 10px;
            margin: 30px 0;
            border-left: 4px solid #28a745;
        }

        .config-section h4 {
            color: #28a745;
            font-weight: 600;
            margin-bottom: 15px;
        }

        .config-item {
            margin-bottom: 20px;
        }

        .config-item strong {
            color: #333;
            display: inline-block;
            min-width: 120px;
        }

        .code-block {
            background: #2d3748;
            color: #e2e8f0;
            padding: 20px;
            border-radius: 8px;
            font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
            font-size: 14px;
            overflow-x: auto;
            margin: 20px 0;
        }

        .warning-box {
            background: #fff3cd;
            border: 1px solid #ffeaa7;
            border-left: 4px solid #f39c12;
            padding: 20px;
            border-radius: 8px;
            margin: 20px 0;
        }

        .warning-box h4 {
            color: #f39c12;
            margin-bottom: 10px;
        }

        .back-btn {
            display: inline-flex;
            align-items: center;
            gap: 8px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 12px 25px;
            border-radius: 25px;
            text-decoration: none;
            font-weight: 500;
            transition: all 0.3s;
            margin-bottom: 30px;
        }

        .back-btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 10px 25px rgba(102,126,234,0.3);
            color: white;
            text-decoration: none;
        }

        /* 响应式设计 */
        @media (max-width: 768px) {
            .header {
                padding: 0 15px;
            }
            .header-nav {
                display: none;
            }
            .brand-section {
                padding: 15px;
            }
            .code-mascot {
                font-size: 2rem;
            }
            .brand-info h1 {
                font-size: 1.8rem;
            }
            .brand-info .version {
                font-size: 0.9rem;
            }
            .brand-right {
                display: none;
            }
            .sidebar {
                padding: 15px 0;
                margin-bottom: 20px;
            }
            .nav-section {
                margin-bottom: 20px;
            }
            .nav-section-title {
                padding: 0 20px 10px;
                font-size: 11px;
            }
            .nav-item {
                padding: 12px 20px;
                font-size: 14px;
            }
            .nav-item .icon {
                margin-right: 8px;
                font-size: 16px;
            }
            .guide-content {
                padding: 30px 25px;
            }
            .guide-header h1 {
                font-size: 2rem;
            }
        }
    </style>
</head>
<body>
    <div class="header">
        <div class="logo">📧 NBEmail</div>
        <div class="header-nav">
            <a href="/guide">使用指南</a>
            <a href="#about">关于</a>
        </div>
        <div class="user-info">
            <span class="user-email">` + userEmail + `</span>
            <button class="logout-btn" onclick="logout()">退出登录</button>
        </div>
    </div>

    <!-- 品牌区域 -->
    <div class="brand-section">
        <div class="brand-left">
            <div class="code-mascot">👩‍💻</div>
            <div class="brand-info">
                <h1>📧 NBEmail</h1>
                <div class="version">v1.0.0 - 专业邮件管理系统</div>
            </div>
        </div>
        <div class="brand-right">
            <span>欢迎使用，` + userEmail + `</span>
        </div>
    </div>

    <div class="container">
        <div class="guide-content">
            <a href="/inbox" class="back-btn">
                <span>←</span>返回收件箱
            </a>

            <div class="guide-header">
                <h1>📖 NBEmail 使用指南</h1>
                <p>欢迎使用 NBEmail 邮件系统，本指南将帮助您快速上手</p>
            </div>

            <div class="guide-section">
                <h2>🚀 快速开始</h2>
                <p>NBEmail 是一个现代化的邮件管理系统，基于 Go + Vue 技术栈构建，提供简洁高效的邮件收发体验。</p>

                <h3>主要功能</h3>
                <ul>
                    <li><strong>收件箱管理</strong> - 查看、阅读和管理收到的邮件</li>
                    <li><strong>邮件发送</strong> - 撰写和发送邮件给任何邮箱地址</li>
                    <li><strong>已发送邮件</strong> - 查看您发送过的所有邮件</li>
                    <li><strong>用户管理</strong> - 管理员可以创建和管理用户账户</li>
                    <li><strong>域名管理</strong> - 管理员可以配置邮件域名</li>
                </ul>
            </div>

            <div class="guide-section">
                <h2>📧 邮件操作</h2>

                <h3>发送邮件</h3>
                <ol>
                    <li>点击左侧导航栏的 "✏️ 写邮件"</li>
                    <li>填写收件人邮箱地址</li>
                    <li>输入邮件主题</li>
                    <li>撰写邮件内容</li>
                    <li>点击 "发送邮件" 按钮</li>
                </ol>

                <h3>管理收件箱</h3>
                <ul>
                    <li><strong>查看邮件</strong> - 点击邮件项目可以查看详细内容</li>
                    <li><strong>标记已读</strong> - 选择邮件后点击 "标记已读" 按钮</li>
                    <li><strong>删除邮件</strong> - 选择邮件后点击 "删除" 按钮</li>
                    <li><strong>刷新邮件</strong> - 点击 "刷新" 按钮获取最新邮件</li>
                </ul>
            </div>

            <div class="guide-section">
                <h2>⚙️ 邮件服务器配置</h2>
                <p>要使用 NBEmail 系统收发邮件，您需要正确配置邮件服务器设置。</p>

                <div class="config-section">
                    <h4>📥 接收邮件配置</h4>
                    <p>要接收外部邮件，您需要配置以下 DNS 记录：</p>

                    <div class="config-item">
                        <strong>记录类型：</strong>A<br>
                        <strong>主机记录：</strong>@<br>
                        <strong>记录值：</strong>您的服务器IP地址<br>
                        <strong>说明：</strong>将域名指向您的邮件服务器
                    </div>

                    <div class="config-item">
                        <strong>记录类型：</strong>MX<br>
                        <strong>主机记录：</strong>@<br>
                        <strong>记录值：</strong>mail.您的域名<br>
                        <strong>优先级：</strong>10<br>
                        <strong>说明：</strong>MX记录指定邮件服务器地址，用于接收邮件
                    </div>
                </div>

                <div class="config-section">
                    <h4>📤 发送邮件配置</h4>
                    <p>要发送邮件到外部邮箱，您需要配置以下设置：</p>

                    <div class="config-item">
                        <strong>记录类型：</strong>TXT<br>
                        <strong>主机记录：</strong>@<br>
                        <strong>记录值：</strong>v=spf1 ip4:您的服务器IP +all<br>
                        <strong>说明：</strong>SPF记录，用于防止邮件被标记为垃圾邮件
                    </div>

                    <div class="config-item">
                        <strong>记录类型：</strong>TXT<br>
                        <strong>主机记录：</strong>_dmarc<br>
                        <strong>记录值：</strong>v=DMARC1; p=quarantine; rua=mailto:admin@您的域名<br>
                        <strong>说明：</strong>DMARC记录，提高邮件送达率
                    </div>
                </div>

                <div class="warning-box">
                    <h4>⚠️ 重要提示</h4>
                    <p>DNS 记录修改后需要一定时间生效（通常为几分钟到几小时）。在此期间，邮件收发功能可能不稳定。</p>
                </div>
            </div>

            <div class="guide-section">
                <h2>🔧 系统管理</h2>
                <p>管理员用户可以访问额外的管理功能：</p>

                <h3>用户管理</h3>
                <ul>
                    <li>创建新用户账户</li>
                    <li>查看所有用户列表</li>
                    <li>删除用户账户</li>
                    <li>修改用户信息</li>
                </ul>

                <h3>域名管理</h3>
                <ul>
                    <li>添加新的邮件域名</li>
                    <li>查看域名状态</li>
                    <li>删除不需要的域名</li>
                </ul>
            </div>

            <div class="guide-section">
                <h2>❓ 常见问题</h2>

                <h3>Q: 为什么我发送的邮件被标记为垃圾邮件？</h3>
                <p>A: 这通常是因为没有正确配置 SPF、DKIM 或 DMARC 记录。请确保按照上述配置指南正确设置 DNS 记录。</p>

                <h3>Q: 我无法接收到外部邮件怎么办？</h3>
                <p>A: 请检查以下几点：</p>
                <ul>
                    <li>确认 MX 记录已正确配置</li>
                    <li>检查服务器防火墙是否开放了 25 端口</li>
                    <li>确认域名解析正常</li>
                </ul>

                <h3>Q: 如何备份邮件数据？</h3>
                <p>A: 系统会自动将邮件数据存储在数据库中。建议定期备份数据库文件以防数据丢失。</p>
            </div>

            <div class="guide-section">
                <h2>📞 技术支持</h2>
                <p>如果您在使用过程中遇到问题，可以通过以下方式获取帮助：</p>
                <ul>
                    <li>查看系统日志文件获取错误信息</li>
                    <li>检查服务器网络连接状态</li>
                    <li>确认配置文件设置正确</li>
                </ul>
            </div>
        </div>
    </div>

    <script>
        async function logout() {
            if (confirm('确定要退出登录吗？')) {
                try {
                    await fetch('/api/logout', { method: 'POST' });
                    window.location.href = '/login';
                } catch (error) {
                    window.location.href = '/login';
                }
            }
        }
    </script>
</body>
</html>`
}

// generateUsersPageTemplate 生成用户管理页面模板
func (s *Server) generateUsersPageTemplate(userEmail string) string {
	return `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>NBEmail - 用户管理</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            background: #f8f9fa;
            color: #333;
        }

        /* 顶部导航栏 */
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 0 30px;
            height: 70px;
            display: flex;
            align-items: center;
            justify-content: space-between;
            box-shadow: 0 2px 20px rgba(0,0,0,0.1);
        }
        .logo {
            font-size: 28px;
            font-weight: 300;
            color: white;
            letter-spacing: -1px;
        }
        .header-nav {
            display: flex;
            align-items: center;
            gap: 30px;
        }
        .header-nav a {
            color: rgba(255,255,255,0.8);
            text-decoration: none;
            font-weight: 500;
            transition: color 0.3s;
        }
        .header-nav a:hover {
            color: white;
        }
        .user-info {
            display: flex;
            align-items: center;
            gap: 20px;
        }
        .user-email {
            color: rgba(255,255,255,0.9);
            font-weight: 500;
        }
        .logout-btn {
            background: rgba(255,255,255,0.2);
            color: white;
            border: 1px solid rgba(255,255,255,0.3);
            padding: 8px 20px;
            border-radius: 20px;
            cursor: pointer;
            font-weight: 500;
            transition: all 0.3s;
        }
        .logout-btn:hover {
            background: rgba(255,255,255,0.3);
            transform: translateY(-1px);
        }

        /* 品牌区域 */
        .brand-section {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            padding: 20px;
            display: flex;
            align-items: center;
            justify-content: space-between;
            color: white;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }

        .brand-left {
            display: flex;
            align-items: center;
            gap: 20px;
        }

        .code-mascot {
            font-size: 3rem;
            animation: float 3s ease-in-out infinite;
        }

        @keyframes float {
            0%, 100% { transform: translateY(0px); }
            50% { transform: translateY(-10px); }
        }

        .brand-info h1 {
            margin: 0;
            font-size: 2.5rem;
            font-weight: 700;
            text-shadow: 2px 2px 4px rgba(0,0,0,0.3);
        }

        .brand-info .version {
            font-size: 1rem;
            opacity: 0.9;
            margin-top: 5px;
        }

        .brand-right {
            display: flex;
            align-items: center;
            gap: 15px;
            font-size: 1.1rem;
        }

        /* 主要内容区域 */
        .container {
            display: flex;
            height: calc(100vh - 140px);
            max-width: 1600px;
            margin: 0 auto;
            gap: 20px;
            padding: 20px;
        }

        /* 侧边栏 */
        .sidebar {
            width: 250px;
            background: white;
            border-radius: 15px;
            padding: 20px 0;
            box-shadow: 0 5px 20px rgba(0,0,0,0.08);
            height: fit-content;
        }
        .nav-section {
            margin-bottom: 30px;
        }
        .nav-section-title {
            padding: 0 25px 15px;
            font-size: 12px;
            font-weight: 600;
            color: #8e9aaf;
            text-transform: uppercase;
            letter-spacing: 1px;
        }
        .nav-item {
            display: flex;
            align-items: center;
            padding: 15px 25px;
            color: #5a6c7d;
            text-decoration: none;
            transition: all 0.3s;
            font-weight: 500;
        }
        .nav-item:hover {
            background: #f8f9fa;
            color: #007bff;
            transform: translateX(5px);
        }
        .nav-item.active {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            box-shadow: 0 5px 15px rgba(102,126,234,0.3);
        }
        .nav-item .icon {
            margin-right: 12px;
            font-size: 18px;
            width: 20px;
        }
        .nav-divider {
            height: 1px;
            background: #e9ecef;
            margin: 20px 25px;
        }

        /* 主内容区 */
        .main {
            flex: 1;
            display: flex;
            flex-direction: column;
            gap: 20px;
        }

        /* 工具栏 */
        .toolbar {
            background: white;
            padding: 20px 25px;
            border-radius: 15px;
            display: flex;
            align-items: center;
            justify-content: space-between;
            box-shadow: 0 5px 20px rgba(0,0,0,0.08);
        }
        .toolbar-left {
            display: flex;
            gap: 15px;
            align-items: center;
        }
        .toolbar-right {
            display: flex;
            gap: 10px;
        }
        .btn {
            padding: 10px 20px;
            border: 2px solid #e9ecef;
            background: white;
            border-radius: 10px;
            cursor: pointer;
            font-weight: 500;
            transition: all 0.3s;
            display: flex;
            align-items: center;
            gap: 8px;
        }
        .btn:hover {
            border-color: #007bff;
            color: #007bff;
            transform: translateY(-2px);
        }
        .btn-primary {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            border-color: transparent;
        }
        .btn-primary:hover {
            transform: translateY(-2px);
            box-shadow: 0 10px 25px rgba(102,126,234,0.3);
        }
        .btn-danger {
            border-color: #dc3545;
            color: #dc3545;
        }
        .btn-danger:hover {
            background: #dc3545;
            color: white;
        }

        /* 用户列表 */
        .user-list {
            background: white;
            border-radius: 15px;
            overflow: hidden;
            box-shadow: 0 5px 20px rgba(0,0,0,0.08);
            flex: 1;
        }
        .user-item {
            padding: 20px 25px;
            border-bottom: 1px solid #f1f3f4;
            display: flex;
            align-items: center;
            justify-content: space-between;
            transition: all 0.3s;
        }
        .user-item:hover {
            background: #f8f9fa;
            transform: translateX(5px);
        }
        .user-item:last-child {
            border-bottom: none;
        }
        .user-info-item {
            display: flex;
            align-items: center;
            gap: 15px;
        }
        .user-avatar {
            width: 45px;
            height: 45px;
            border-radius: 50%;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            display: flex;
            align-items: center;
            justify-content: center;
            color: white;
            font-weight: 600;
            font-size: 16px;
        }
        .user-details {
            display: flex;
            flex-direction: column;
            gap: 3px;
        }
        .user-email-item {
            font-weight: 600;
            color: #333;
            font-size: 15px;
        }
        .user-name {
            color: #8e9aaf;
            font-size: 13px;
        }
        .user-badge {
            padding: 3px 8px;
            border-radius: 12px;
            font-size: 11px;
            font-weight: 600;
            text-transform: uppercase;
            margin-left: 10px;
        }
        .badge-admin {
            background: #e3f2fd;
            color: #1976d2;
        }
        .badge-user {
            background: #f3e5f5;
            color: #7b1fa2;
        }
        .user-actions {
            display: flex;
            gap: 10px;
        }
        .btn-sm {
            padding: 6px 12px;
            font-size: 12px;
            border-radius: 8px;
        }

        /* 空状态 */
        .empty-state {
            text-align: center;
            padding: 80px 20px;
            color: #8e9aaf;
        }
        .empty-state .icon {
            font-size: 4rem;
            margin-bottom: 20px;
            opacity: 0.5;
        }
        .empty-state h3 {
            font-size: 1.2rem;
            margin-bottom: 10px;
            color: #5a6c7d;
        }

        /* 模态框样式 */
        .modal {
            display: none;
            position: fixed;
            z-index: 1000;
            left: 0;
            top: 0;
            width: 100%;
            height: 100%;
            background-color: rgba(0,0,0,0.5);
            backdrop-filter: blur(5px);
        }
        .modal-content {
            background-color: white;
            margin: 10% auto;
            padding: 30px;
            border-radius: 15px;
            width: 90%;
            max-width: 500px;
            box-shadow: 0 20px 40px rgba(0,0,0,0.2);
        }
        .modal-header {
            margin-bottom: 25px;
            padding-bottom: 15px;
            border-bottom: 2px solid #f1f3f4;
        }
        .modal-header h3 {
            color: #333;
            font-weight: 600;
        }
        .close {
            color: #aaa;
            float: right;
            font-size: 28px;
            font-weight: bold;
            cursor: pointer;
            line-height: 1;
        }
        .close:hover {
            color: #333;
        }
        .form-group {
            margin-bottom: 20px;
        }
        .form-group label {
            display: block;
            margin-bottom: 8px;
            font-weight: 600;
            color: #333;
        }
        .form-group input {
            width: 100%;
            padding: 12px 15px;
            border: 2px solid #e9ecef;
            border-radius: 8px;
            font-size: 14px;
            transition: all 0.3s;
        }
        .form-group input:focus {
            outline: none;
            border-color: #667eea;
            box-shadow: 0 0 0 3px rgba(102,126,234,0.1);
        }
        .form-group select {
            width: 100%;
            padding: 12px 15px;
            border: 2px solid #e9ecef;
            border-radius: 8px;
            font-size: 14px;
            transition: all 0.3s;
            background: white;
        }
        .form-group select:focus {
            outline: none;
            border-color: #667eea;
            box-shadow: 0 0 0 3px rgba(102,126,234,0.1);
        }

        /* 标签页样式 */
        .assign-tabs {
            display: flex;
            margin-bottom: 20px;
            border-bottom: 2px solid #f1f3f4;
        }
        .tab-btn {
            padding: 12px 24px;
            border: none;
            background: none;
            cursor: pointer;
            font-size: 14px;
            font-weight: 500;
            color: #666;
            border-bottom: 2px solid transparent;
            transition: all 0.3s;
        }
        .tab-btn.active {
            color: #667eea;
            border-bottom-color: #667eea;
        }
        .tab-btn:hover {
            color: #667eea;
        }
        .tab-content {
            display: none;
        }
        .tab-content.active {
            display: block;
        }

        /* 复选框列表样式 */
        .checkbox-list {
            max-height: 200px;
            overflow-y: auto;
            border: 2px solid #e9ecef;
            border-radius: 8px;
            padding: 10px;
            background: #f8f9fa;
        }
        .checkbox-item {
            display: flex;
            align-items: center;
            padding: 8px 0;
            border-bottom: 1px solid #e9ecef;
        }
        .checkbox-item:last-child {
            border-bottom: none;
        }
        .checkbox-item input[type="checkbox"] {
            margin-right: 10px;
            width: 16px;
            height: 16px;
        }
        .checkbox-item label {
            margin: 0;
            cursor: pointer;
            flex: 1;
            font-weight: normal;
        }
        .domain-info {
            display: flex;
            flex-direction: column;
        }
        .domain-name {
            font-weight: 600;
            color: #333;
        }
        .domain-status {
            font-size: 12px;
            color: #666;
        }

        /* 响应式设计 */
        @media (max-width: 1024px) {
            .brand-section {
                flex-direction: column;
                gap: 15px;
                text-align: center;
            }
            .brand-left {
                justify-content: center;
            }
            .brand-info h1 {
                font-size: 2rem;
            }
            .container {
                flex-direction: column;
                padding: 10px;
            }
            .sidebar {
                width: 100%;
                order: 2;
            }
            .main {
                order: 1;
                min-height: 60vh;
                width: 100%;
            }
        }

        @media (max-width: 768px) {
            .header {
                padding: 0 15px;
            }
            .header-nav {
                display: none;
            }
            .brand-section {
                padding: 15px;
            }
            .code-mascot {
                font-size: 2rem;
            }
            .brand-info h1 {
                font-size: 1.8rem;
            }
            .brand-info .version {
                font-size: 0.9rem;
            }
            .brand-right {
                font-size: 1rem;
            }
            .sidebar {
                padding: 15px 0;
                margin-bottom: 20px;
            }
            .nav-section {
                margin-bottom: 20px;
            }
            .nav-section-title {
                padding: 0 20px 10px;
                font-size: 11px;
            }
            .nav-item {
                padding: 12px 20px;
                font-size: 14px;
            }
            .nav-item .icon {
                margin-right: 8px;
                font-size: 16px;
            }
            .user-item {
                padding: 15px;
                gap: 15px;
            }
            .user-avatar {
                width: 35px;
                height: 35px;
                font-size: 14px;
            }
        }
    </style>
</head>
<body>
    <div class="header">
        <div class="logo">📧 NBEmail</div>
        <div class="header-nav">
            <a href="/guide">使用指南</a>
            <a href="#about">关于</a>
        </div>
        <div class="user-info">
            <span class="user-email">` + userEmail + `</span>
            <button class="logout-btn" onclick="logout()">退出登录</button>
        </div>
    </div>

    <!-- 品牌区域 -->
    <div class="brand-section">
        <div class="brand-left">
            <div class="code-mascot">👩‍💻</div>
            <div class="brand-info">
                <h1>📧 NBEmail</h1>
                <div class="version">v1.0.0 - 专业邮件管理系统</div>
            </div>
        </div>
        <div class="brand-right">
            <span>欢迎使用，` + userEmail + `</span>
        </div>
    </div>

    <div class="container">
        <div class="sidebar">
            <div class="nav-section">
                <div class="nav-section-title">邮箱</div>
                <a href="/inbox" class="nav-item">
                    <span class="icon">📥</span>收件箱
                </a>
                <a href="/sent" class="nav-item">
                    <span class="icon">📤</span>已发送
                </a>
                <a href="/compose" class="nav-item">
                    <span class="icon">✏️</span>写邮件
                </a>
            </div>
            <div class="nav-divider"></div>
            <div class="nav-section">
                <div class="nav-section-title">管理</div>
                <a href="/users" class="nav-item active">
                    <span class="icon">👥</span>用户管理
                </a>
                <a href="/domains" class="nav-item">
                    <span class="icon">🌐</span>域名管理
                </a>
                <a href="/smtp-configs" class="nav-item">
                    <span class="icon">📮</span>SMTP配置
                </a>
            </div>
        </div>
        <div class="main">
            <div class="toolbar">
                <div class="toolbar-left">
                    <h2 style="margin: 0; color: #333; font-weight: 600;">👥 用户管理</h2>
                </div>
                <div class="toolbar-right">
                    <button class="btn btn-primary" onclick="showCreateUserModal()">
                        <span>➕</span>创建用户
                    </button>
                    <button class="btn" onclick="refreshUsers()">
                        <span>🔄</span>刷新
                    </button>
                </div>
            </div>
            <div class="user-list" id="userList">
                <div class="empty-state">
                    <div class="icon">👥</div>
                    <h3>正在加载用户...</h3>
                    <p>请稍候，正在获取用户列表</p>
                </div>
            </div>
        </div>
    </div>

    <!-- 创建用户模态框 -->
    <div id="createUserModal" class="modal">
        <div class="modal-content">
            <div class="modal-header">
                <span class="close" onclick="closeCreateUserModal()">&times;</span>
                <h3>➕ 创建新用户</h3>
            </div>
            <form id="createUserForm">
                <div class="form-group">
                    <label for="newUserEmail">邮箱地址 *</label>
                    <input type="email" id="newUserEmail" required placeholder="请输入邮箱地址">
                </div>
                <div class="form-group">
                    <label for="newUserName">用户名 *</label>
                    <input type="text" id="newUserName" required placeholder="请输入用户名">
                </div>
                <div class="form-group">
                    <label for="newUserPassword">密码 *</label>
                    <input type="password" id="newUserPassword" required placeholder="请输入密码">
                </div>
                <div style="display: flex; gap: 10px; justify-content: flex-end; margin-top: 25px;">
                    <button type="button" class="btn" onclick="closeCreateUserModal()">取消</button>
                    <button type="submit" class="btn btn-primary">创建用户</button>
                </div>
            </form>
        </div>
    </div>

    <!-- 编辑用户模态框 -->
    <div id="editUserModal" class="modal">
        <div class="modal-content">
            <div class="modal-header">
                <span class="close" onclick="closeEditUserModal()">&times;</span>
                <h3>✏️ 编辑用户</h3>
            </div>
            <form id="editUserForm">
                <div class="form-group">
                    <label for="editUserEmail">邮箱地址</label>
                    <input type="email" id="editUserEmail" readonly style="background-color: #f5f5f5;">
                    <input type="hidden" id="editUserID">
                </div>
                <div class="form-group">
                    <label for="editUserName">用户名 *</label>
                    <input type="text" id="editUserName" required placeholder="请输入用户名">
                </div>
                <div class="form-group">
                    <label for="editUserPassword">新密码</label>
                    <input type="password" id="editUserPassword" placeholder="留空则不修改密码">
                    <small style="color: #666; font-size: 12px;">留空则不修改密码</small>
                </div>
                <div style="display: flex; gap: 10px; justify-content: flex-end; margin-top: 25px;">
                    <button type="button" class="btn" onclick="closeEditUserModal()">取消</button>
                    <button type="submit" class="btn btn-primary">保存修改</button>
                </div>
            </form>
        </div>
    </div>

    <!-- 分配邮箱模态框 -->
    <div id="assignModal" class="modal">
        <div class="modal-content">
            <div class="modal-header">
                <span class="close" onclick="closeAssignModal()">&times;</span>
                <h3>📧 分配邮箱</h3>
            </div>
            <div class="assign-tabs">
                <button class="tab-btn active" onclick="switchTab('mailbox')">分配邮箱</button>
                <button class="tab-btn" onclick="switchTab('domain')">分配域名</button>
                <button class="tab-btn" onclick="switchTab('reclaim')">回收域名</button>
            </div>

            <!-- 分配邮箱标签页 -->
            <div id="mailboxTab" class="tab-content active">
                <form id="assignMailboxForm">
                    <div class="form-group">
                        <label>用户</label>
                        <input type="text" id="assignUserEmail" readonly>
                        <input type="hidden" id="assignUserID">
                    </div>
                    <div class="form-group">
                        <label for="selectDomain">选择域名 *</label>
                        <select id="selectDomain" required>
                            <option value="">请选择域名</option>
                        </select>
                    </div>
                    <div class="form-group">
                        <label for="mailboxPrefix">邮箱前缀</label>
                        <input type="text" id="mailboxPrefix" placeholder="留空则随机生成">
                    </div>
                    <div class="form-group">
                        <label for="mailboxCount">生成数量</label>
                        <input type="number" id="mailboxCount" value="1" min="1" max="10">
                    </div>
                    <div style="display: flex; gap: 10px; justify-content: flex-end; margin-top: 25px;">
                        <button type="button" class="btn" onclick="closeAssignModal()">取消</button>
                        <button type="submit" class="btn btn-primary">分配邮箱</button>
                    </div>
                </form>
            </div>

            <!-- 分配域名标签页 -->
            <div id="domainTab" class="tab-content">
                <form id="assignDomainForm">
                    <div class="form-group">
                        <label>用户</label>
                        <input type="text" id="assignDomainUserEmail" readonly>
                        <input type="hidden" id="assignDomainUserID">
                    </div>
                    <div class="form-group">
                        <label for="selectDomains">选择域名 *</label>
                        <div id="domainCheckboxList" class="checkbox-list">
                            <!-- 域名复选框将在这里动态生成 -->
                        </div>
                    </div>
                    <div style="display: flex; gap: 10px; justify-content: flex-end; margin-top: 25px;">
                        <button type="button" class="btn" onclick="closeAssignModal()">取消</button>
                        <button type="submit" class="btn btn-primary">分配域名</button>
                    </div>
                </form>
            </div>

            <!-- 回收域名标签页 -->
            <div id="reclaimTab" class="tab-content">
                <form id="reclaimDomainForm">
                    <div class="form-group">
                        <label>用户</label>
                        <input type="text" id="reclaimDomainUserEmail" readonly>
                        <input type="hidden" id="reclaimDomainUserID">
                    </div>
                    <div class="form-group">
                        <label for="reclaimDomains">选择要回收的域名 *</label>
                        <div id="reclaimDomainCheckboxList" class="checkbox-list">
                            <!-- 用户拥有的域名复选框将在这里动态生成 -->
                        </div>
                    </div>
                    <div style="display: flex; gap: 10px; justify-content: flex-end; margin-top: 25px;">
                        <button type="button" class="btn" onclick="closeAssignModal()">取消</button>
                        <button type="submit" class="btn btn-danger">回收域名</button>
                    </div>
                </form>
            </div>
        </div>
    </div>
    <script>
        let users = [];

        async function loadUsers() {
            try {
                const response = await fetch('/api/users');
                const result = await response.json();
                if (result.success) {
                    users = result.data || [];
                    renderUsers();
                } else {
                    document.getElementById('userList').innerHTML = '<div class="empty-state">加载用户失败</div>';
                }
            } catch (error) {
                document.getElementById('userList').innerHTML = '<div class="empty-state">加载用户失败</div>';
            }
        }

        function renderUsers() {
            const userList = document.getElementById('userList');
            if (users.length === 0) {
                userList.innerHTML = ` + "`" + `
                    <div class="empty-state">
                        <div class="icon">👤</div>
                        <h3>暂无用户</h3>
                        <p>系统中还没有创建任何用户</p>
                    </div>
                ` + "`" + `;
                return;
            }

            userList.innerHTML = users.map(user => {
                const emailInitial = user.email.charAt(0).toUpperCase();
                const badgeClass = user.is_admin ? 'badge-admin' : 'badge-user';
                const badgeText = user.is_admin ? '管理员' : '普通用户';

                return ` + "`" + `
                    <div class="user-item">
                        <div class="user-info-item">
                            <div class="user-avatar">${emailInitial}</div>
                            <div class="user-details">
                                <div class="user-email-item">
                                    ${user.email}
                                    <span class="user-badge ${badgeClass}">${badgeText}</span>
                                </div>
                                <div class="user-name">${user.name}</div>
                            </div>
                        </div>
                        <div class="user-actions">
                            ${!user.is_admin ? ` + "`" + `
                                <button class="btn btn-sm btn-primary" onclick="showAssignModal(${user.id}, '${user.email}')">
                                    <span>📧</span>分配
                                </button>
                            ` + "`" + ` : ''}
                            <button class="btn btn-sm" onclick="editUser(${user.id})">
                                <span>✏️</span>编辑
                            </button>
                            <button class="btn btn-sm btn-danger" onclick="deleteUser(${user.id})">
                                <span>🗑️</span>删除
                            </button>
                        </div>
                    </div>
                ` + "`" + `;
            }).join('');
        }

        function showCreateUserModal() {
            document.getElementById('createUserModal').style.display = 'block';
            document.getElementById('newUserEmail').focus();
        }

        function closeCreateUserModal() {
            document.getElementById('createUserModal').style.display = 'none';
            document.getElementById('createUserForm').reset();
        }

        // 处理创建用户表单提交
        document.getElementById('createUserForm').addEventListener('submit', async (e) => {
            e.preventDefault();

            const email = document.getElementById('newUserEmail').value;
            const name = document.getElementById('newUserName').value;
            const password = document.getElementById('newUserPassword').value;

            const submitBtn = e.target.querySelector('button[type="submit"]');
            const originalText = submitBtn.innerHTML;

            // 显示加载状态
            submitBtn.innerHTML = '<span>⏳</span>创建中...';
            submitBtn.disabled = true;

            try {
                const response = await fetch('/api/users', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ email, name, password })
                });

                const result = await response.json();
                if (result.success) {
                    // 显示成功状态
                    submitBtn.innerHTML = '<span>✅</span>创建成功';
                    submitBtn.style.background = '#28a745';

                    showNotification('用户创建成功！', 'success');

                    // 2秒后关闭模态框并刷新列表
                    setTimeout(() => {
                        closeCreateUserModal();
                        loadUsers();
                        submitBtn.innerHTML = originalText;
                        submitBtn.style.background = '';
                        submitBtn.disabled = false;
                    }, 2000);
                } else {
                    throw new Error(result.message || '创建失败');
                }
            } catch (error) {
                // 显示错误状态
                submitBtn.innerHTML = '<span>❌</span>创建失败';
                submitBtn.style.background = '#dc3545';

                showNotification('创建失败: ' + error.message, 'error');

                // 3秒后恢复按钮
                setTimeout(() => {
                    submitBtn.innerHTML = originalText;
                    submitBtn.style.background = '';
                    submitBtn.disabled = false;
                }, 3000);
            }
        });

        // 处理分配邮箱表单提交
        document.getElementById('assignMailboxForm').addEventListener('submit', async (e) => {
            e.preventDefault();

            const userId = document.getElementById('assignUserID').value;
            const domainId = document.getElementById('selectDomain').value;
            const prefix = document.getElementById('mailboxPrefix').value;
            const count = parseInt(document.getElementById('mailboxCount').value);

            const submitBtn = e.target.querySelector('button[type="submit"]');
            const originalText = submitBtn.innerHTML;

            try {
                submitBtn.innerHTML = '<span>⏳</span>分配中...';
                submitBtn.disabled = true;

                const response = await fetch(` + "`" + `/api/users/${userId}/assign-mailboxes` + "`" + `, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        domain_id: parseInt(domainId),
                        prefix: prefix,
                        count: count
                    })
                });

                const result = await response.json();
                if (result.success) {
                    submitBtn.innerHTML = '<span>✅</span>分配成功';
                    submitBtn.style.background = '#28a745';
                    showNotification('邮箱分配成功！', 'success');

                    setTimeout(() => {
                        closeAssignModal();
                        submitBtn.innerHTML = originalText;
                        submitBtn.style.background = '';
                        submitBtn.disabled = false;
                    }, 2000);
                } else {
                    throw new Error(result.message || '分配失败');
                }
            } catch (error) {
                submitBtn.innerHTML = '<span>❌</span>分配失败';
                submitBtn.style.background = '#dc3545';
                showNotification('分配失败: ' + error.message, 'error');

                setTimeout(() => {
                    submitBtn.innerHTML = originalText;
                    submitBtn.style.background = '';
                    submitBtn.disabled = false;
                }, 3000);
            }
        });

        // 处理分配域名表单提交
        document.getElementById('assignDomainForm').addEventListener('submit', async (e) => {
            e.preventDefault();

            const userId = document.getElementById('assignDomainUserID').value;

            // 获取选中的域名ID列表
            const selectedDomains = [];
            const checkboxes = document.querySelectorAll('#domainCheckboxList input[type="checkbox"]:checked');
            checkboxes.forEach(checkbox => {
                selectedDomains.push(parseInt(checkbox.value));
            });

            if (selectedDomains.length === 0) {
                showNotification('请至少选择一个域名', 'error');
                return;
            }

            const submitBtn = e.target.querySelector('button[type="submit"]');
            const originalText = submitBtn.innerHTML;

            try {
                submitBtn.innerHTML = '<span>⏳</span>分配中...';
                submitBtn.disabled = true;

                const response = await fetch(` + "`" + `/api/users/${userId}/assign-domains` + "`" + `, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        domain_ids: selectedDomains
                    })
                });

                const result = await response.json();
                if (result.success) {
                    submitBtn.innerHTML = '<span>✅</span>分配成功';
                    submitBtn.style.background = '#28a745';
                    showNotification(` + "`" + `成功分配了 ${selectedDomains.length} 个域名！` + "`" + `, 'success');

                    setTimeout(() => {
                        closeAssignModal();
                        submitBtn.innerHTML = originalText;
                        submitBtn.style.background = '';
                        submitBtn.disabled = false;
                        loadDomainsForAssign(); // 重新加载域名列表
                    }, 2000);
                } else {
                    throw new Error(result.message || '分配失败');
                }
            } catch (error) {
                submitBtn.innerHTML = '<span>❌</span>分配失败';
                submitBtn.style.background = '#dc3545';
                showNotification('分配失败: ' + error.message, 'error');

                setTimeout(() => {
                    submitBtn.innerHTML = originalText;
                    submitBtn.style.background = '';
                    submitBtn.disabled = false;
                }, 3000);
            }
        });

        // 处理回收域名表单提交
        document.getElementById('reclaimDomainForm').addEventListener('submit', async (e) => {
            e.preventDefault();

            const userId = document.getElementById('reclaimDomainUserID').value;

            // 获取选中的域名ID列表
            const selectedDomains = [];
            const checkboxes = document.querySelectorAll('#reclaimDomainCheckboxList input[type="checkbox"]:checked');
            checkboxes.forEach(checkbox => {
                selectedDomains.push(parseInt(checkbox.value));
            });

            if (selectedDomains.length === 0) {
                showNotification('请至少选择一个要回收的域名', 'error');
                return;
            }

            if (!confirm(` + "`" + `确定要回收选中的 ${selectedDomains.length} 个域名吗？回收后域名将变为公共域名。` + "`" + `)) {
                return;
            }

            const submitBtn = e.target.querySelector('button[type="submit"]');
            const originalText = submitBtn.innerHTML;

            try {
                submitBtn.innerHTML = '<span>⏳</span>回收中...';
                submitBtn.disabled = true;

                const response = await fetch(` + "`" + `/api/users/${userId}/reclaim-domains` + "`" + `, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        domain_ids: selectedDomains
                    })
                });

                const result = await response.json();
                if (result.success) {
                    submitBtn.innerHTML = '<span>✅</span>回收成功';
                    submitBtn.style.background = '#28a745';
                    showNotification(` + "`" + `成功回收了 ${selectedDomains.length} 个域名！` + "`" + `, 'success');

                    setTimeout(() => {
                        closeAssignModal();
                        submitBtn.innerHTML = originalText;
                        submitBtn.style.background = '';
                        submitBtn.disabled = false;
                        loadDomainsForAssign(); // 重新加载域名列表
                        loadUserDomainsForReclaim(userId); // 重新加载用户域名列表
                    }, 2000);
                } else {
                    throw new Error(result.message || '回收失败');
                }
            } catch (error) {
                submitBtn.innerHTML = '<span>❌</span>回收失败';
                submitBtn.style.background = '#dc3545';
                showNotification('回收失败: ' + error.message, 'error');

                setTimeout(() => {
                    submitBtn.innerHTML = originalText;
                    submitBtn.style.background = '';
                    submitBtn.disabled = false;
                }, 3000);
            }
        });

        async function createUser(email, name, password) {
            try {
                const response = await fetch('/api/users', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ email, name, password })
                });

                const result = await response.json();
                if (result.success) {
                    alert('用户创建成功');
                    loadUsers();
                } else {
                    alert('创建失败: ' + result.message);
                }
            } catch (error) {
                alert('创建失败，请重试');
            }
        }

        async function deleteUser(userId) {
            if (!confirm('确定要删除这个用户吗？此操作不可撤销！')) return;

            try {
                const response = await fetch(` + "`" + `/api/users/${userId}` + "`" + `, { method: 'DELETE' });
                const result = await response.json();
                if (result.success) {
                    showNotification('用户删除成功', 'success');
                    loadUsers();
                } else {
                    showNotification('删除失败: ' + result.message, 'error');
                }
            } catch (error) {
                showNotification('删除失败，请重试', 'error');
            }
        }

        async function editUser(userId) {
            try {
                // 获取用户信息
                const response = await fetch('/api/users');
                const result = await response.json();
                if (result.success) {
                    const user = result.data.find(u => u.id === userId);
                    if (user) {
                        // 填充表单
                        document.getElementById('editUserID').value = user.id;
                        document.getElementById('editUserEmail').value = user.email;
                        document.getElementById('editUserName').value = user.name;
                        document.getElementById('editUserPassword').value = '';

                        // 显示模态框
                        document.getElementById('editUserModal').style.display = 'block';
                        document.getElementById('editUserName').focus();
                    } else {
                        showNotification('用户不存在', 'error');
                    }
                } else {
                    showNotification('获取用户信息失败: ' + result.message, 'error');
                }
            } catch (error) {
                showNotification('获取用户信息失败，请重试', 'error');
            }
        }

        function closeEditUserModal() {
            document.getElementById('editUserModal').style.display = 'none';
            document.getElementById('editUserForm').reset();
        }

        // 编辑用户表单提交
        document.getElementById('editUserForm').addEventListener('submit', async function(e) {
            e.preventDefault();

            const userId = document.getElementById('editUserID').value;
            const name = document.getElementById('editUserName').value.trim();
            const password = document.getElementById('editUserPassword').value.trim();

            if (!name) {
                showNotification('请输入用户名', 'error');
                return;
            }

            try {
                const requestData = { name };
                if (password) {
                    requestData.password = password;
                }

                const response = await fetch('/api/users/' + userId, {
                    method: 'PUT',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(requestData)
                });

                const result = await response.json();
                if (result.success) {
                    showNotification('用户更新成功', 'success');
                    closeEditUserModal();
                    loadUsers();
                } else {
                    showNotification('更新失败: ' + result.message, 'error');
                }
            } catch (error) {
                showNotification('更新失败，请重试', 'error');
            }
        });

        // 分配相关函数
        function showAssignModal(userId, userEmail) {
            document.getElementById('assignUserID').value = userId;
            document.getElementById('assignUserEmail').value = userEmail;
            document.getElementById('assignDomainUserID').value = userId;
            document.getElementById('assignDomainUserEmail').value = userEmail;
            document.getElementById('reclaimDomainUserID').value = userId;
            document.getElementById('reclaimDomainUserEmail').value = userEmail;
            document.getElementById('assignModal').style.display = 'block';
            loadDomainsForAssign();
            loadUserDomainsForReclaim(userId);
        }

        function closeAssignModal() {
            document.getElementById('assignModal').style.display = 'none';
            document.getElementById('assignMailboxForm').reset();
            document.getElementById('assignDomainForm').reset();
            document.getElementById('reclaimDomainForm').reset();

            // 清空域名复选框的选中状态
            const checkboxes = document.querySelectorAll('#domainCheckboxList input[type="checkbox"]');
            checkboxes.forEach(checkbox => {
                checkbox.checked = false;
            });

            // 清空回收域名复选框的选中状态
            const reclaimCheckboxes = document.querySelectorAll('#reclaimDomainCheckboxList input[type="checkbox"]');
            reclaimCheckboxes.forEach(checkbox => {
                checkbox.checked = false;
            });
        }

        function switchTab(tabName) {
            // 切换标签按钮状态
            document.querySelectorAll('.tab-btn').forEach(btn => btn.classList.remove('active'));
            event.target.classList.add('active');

            // 切换标签内容
            document.querySelectorAll('.tab-content').forEach(content => content.classList.remove('active'));
            document.getElementById(tabName + 'Tab').classList.add('active');
        }

        async function loadDomainsForAssign() {
            try {
                const response = await fetch('/api/domains');
                const result = await response.json();
                if (result.success) {
                    // 为分配邮箱加载域名下拉框
                    const select = document.getElementById('selectDomain');
                    select.innerHTML = '<option value="">请选择域名</option>';

                    // 为分配域名加载域名复选框列表
                    const checkboxList = document.getElementById('domainCheckboxList');
                    checkboxList.innerHTML = '';

                    result.data.forEach(domain => {
                        if (domain.is_active) {
                            // 添加到邮箱分配的下拉框
                            select.innerHTML += ` + "`" + `<option value="${domain.id}">${domain.name}</option>` + "`" + `;

                            // 添加到域名分配的复选框列表（只显示公共域名，即没有分配给特定用户的域名）
                            if (!domain.user_id) {
                                const checkboxItem = document.createElement('div');
                                checkboxItem.className = 'checkbox-item';
                                checkboxItem.innerHTML = ` + "`" + `
                                    <input type="checkbox" id="domain_${domain.id}" value="${domain.id}">
                                    <label for="domain_${domain.id}">
                                        <div class="domain-info">
                                            <div class="domain-name">${domain.name}</div>
                                            <div class="domain-status">${domain.dns_verified ? '✅ 已验证' : '⏳ 未验证'}</div>
                                        </div>
                                    </label>
                                ` + "`" + `;
                                checkboxList.appendChild(checkboxItem);
                            }
                        }
                    });

                    if (checkboxList.children.length === 0) {
                        checkboxList.innerHTML = '<div style="text-align: center; color: #666; padding: 20px;">暂无可分配的域名</div>';
                    }
                }
            } catch (error) {
                console.error('加载域名失败:', error);
            }
        }

        // 加载用户拥有的域名用于回收
        async function loadUserDomainsForReclaim(userId) {
            try {
                const response = await fetch('/api/domains');
                const result = await response.json();
                if (result.success) {
                    const checkboxList = document.getElementById('reclaimDomainCheckboxList');
                    checkboxList.innerHTML = '';

                    // 筛选出属于该用户的域名
                    const userDomains = result.data.filter(domain => domain.user_id === parseInt(userId) && domain.is_active);

                    if (userDomains.length === 0) {
                        checkboxList.innerHTML = '<div style="text-align: center; padding: 20px; color: #666;">该用户暂无已分配的域名</div>';
                        return;
                    }

                    userDomains.forEach(domain => {
                        const checkboxItem = document.createElement('div');
                        checkboxItem.className = 'checkbox-item';
                        checkboxItem.innerHTML = ` + "`" + `
                            <input type="checkbox" id="reclaim_domain_${domain.id}" value="${domain.id}">
                            <label for="reclaim_domain_${domain.id}">
                                <div class="domain-info">
                                    <div class="domain-name">${domain.name}</div>
                                    <div class="domain-status">${domain.dns_verified ? '✅ 已验证' : '⏳ 未验证'}</div>
                                </div>
                            </label>
                        ` + "`" + `;
                        checkboxList.appendChild(checkboxItem);
                    });
                } else {
                    document.getElementById('reclaimDomainCheckboxList').innerHTML = '<div style="text-align: center; padding: 20px; color: #dc3545;">加载域名失败</div>';
                }
            } catch (error) {
                document.getElementById('reclaimDomainCheckboxList').innerHTML = '<div style="text-align: center; padding: 20px; color: #dc3545;">网络错误</div>';
            }
        }

        function refreshUsers() {
            showNotification('正在刷新用户列表...', 'info');
            loadUsers();
        }

        // 显示通知
        function showNotification(message, type = 'info') {
            // 创建通知元素
            const notification = document.createElement('div');
            notification.style.cssText = ` + "`" + `
                position: fixed;
                top: 20px;
                right: 20px;
                padding: 15px 25px;
                border-radius: 10px;
                color: white;
                font-weight: 500;
                z-index: 10000;
                transform: translateX(400px);
                transition: all 0.3s ease;
                box-shadow: 0 5px 20px rgba(0,0,0,0.2);
            ` + "`" + `;

            // 设置不同类型的样式
            switch(type) {
                case 'success':
                    notification.style.background = 'linear-gradient(135deg, #28a745, #20c997)';
                    break;
                case 'error':
                    notification.style.background = 'linear-gradient(135deg, #dc3545, #e74c3c)';
                    break;
                case 'info':
                default:
                    notification.style.background = 'linear-gradient(135deg, #667eea, #764ba2)';
            }

            notification.textContent = message;
            document.body.appendChild(notification);

            // 显示动画
            setTimeout(() => {
                notification.style.transform = 'translateX(0)';
            }, 100);

            // 自动隐藏
            setTimeout(() => {
                notification.style.transform = 'translateX(400px)';
                setTimeout(() => {
                    if (document.body.contains(notification)) {
                        document.body.removeChild(notification);
                    }
                }, 300);
            }, 3000);
        }

        async function logout() {
            if (confirm('确定要退出登录吗？')) {
                try {
                    await fetch('/api/logout', { method: 'POST' });
                    window.location.href = '/login';
                } catch (error) {
                    window.location.href = '/login';
                }
            }
        }

        // 点击模态框外部关闭
        window.onclick = function(event) {
            const createUserModal = document.getElementById('createUserModal');
            const assignModal = document.getElementById('assignModal');

            if (event.target == createUserModal) {
                closeCreateUserModal();
            } else if (event.target == assignModal) {
                closeAssignModal();
            }
        }

        // 页面加载时获取用户列表
        loadUsers();
    </script>
</body>
</html>`
}

// generateDomainsPageTemplate 生成域名管理页面模板
func (s *Server) generateDomainsPageTemplate(userEmail string) string {
	return `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>NBEmail - 域名管理</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            background: #f8f9fa;
            color: #333;
        }

        /* 顶部导航栏 */
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 0 30px;
            height: 70px;
            display: flex;
            align-items: center;
            justify-content: space-between;
            box-shadow: 0 2px 20px rgba(0,0,0,0.1);
        }
        .logo {
            font-size: 28px;
            font-weight: 300;
            color: white;
            letter-spacing: -1px;
        }
        .header-nav {
            display: flex;
            align-items: center;
            gap: 30px;
        }
        .header-nav a {
            color: rgba(255,255,255,0.8);
            text-decoration: none;
            font-weight: 500;
            transition: color 0.3s;
        }
        .header-nav a:hover {
            color: white;
        }
        .user-info {
            display: flex;
            align-items: center;
            gap: 20px;
        }
        .user-email {
            color: rgba(255,255,255,0.9);
            font-weight: 500;
        }
        .logout-btn {
            background: rgba(255,255,255,0.2);
            color: white;
            border: 1px solid rgba(255,255,255,0.3);
            padding: 8px 20px;
            border-radius: 20px;
            cursor: pointer;
            font-weight: 500;
            transition: all 0.3s;
        }
        .logout-btn:hover {
            background: rgba(255,255,255,0.3);
            transform: translateY(-1px);
        }

        /* 品牌区域 */
        .brand-section {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            padding: 20px;
            display: flex;
            align-items: center;
            justify-content: space-between;
            color: white;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }

        .brand-left {
            display: flex;
            align-items: center;
            gap: 20px;
        }

        .code-mascot {
            font-size: 3rem;
            animation: float 3s ease-in-out infinite;
        }

        @keyframes float {
            0%, 100% { transform: translateY(0px); }
            50% { transform: translateY(-10px); }
        }

        .brand-info h1 {
            margin: 0;
            font-size: 2.5rem;
            font-weight: 700;
            text-shadow: 2px 2px 4px rgba(0,0,0,0.3);
        }

        .brand-info .version {
            font-size: 1rem;
            opacity: 0.9;
            margin-top: 5px;
        }

        .brand-right {
            display: flex;
            align-items: center;
            gap: 15px;
            font-size: 1.1rem;
        }

        /* 主要内容区域 */
        .container {
            display: flex;
            height: calc(100vh - 140px);
            max-width: 1600px;
            margin: 0 auto;
            gap: 20px;
            padding: 20px;
        }

        /* 侧边栏 */
        .sidebar {
            width: 250px;
            background: white;
            border-radius: 15px;
            padding: 20px 0;
            box-shadow: 0 5px 20px rgba(0,0,0,0.08);
            height: fit-content;
        }
        .nav-section {
            margin-bottom: 30px;
        }
        .nav-section-title {
            padding: 0 25px 15px;
            font-size: 12px;
            font-weight: 600;
            color: #8e9aaf;
            text-transform: uppercase;
            letter-spacing: 1px;
        }
        .nav-item {
            display: flex;
            align-items: center;
            padding: 15px 25px;
            color: #5a6c7d;
            text-decoration: none;
            transition: all 0.3s;
            font-weight: 500;
        }
        .nav-item:hover {
            background: #f8f9fa;
            color: #007bff;
            transform: translateX(5px);
        }
        .nav-item.active {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            box-shadow: 0 5px 15px rgba(102,126,234,0.3);
        }
        .nav-item .icon {
            margin-right: 12px;
            font-size: 18px;
            width: 20px;
        }
        .nav-divider {
            height: 1px;
            background: #e9ecef;
            margin: 20px 25px;
        }

        /* 主内容区 */
        .main {
            flex: 1;
            display: flex;
            flex-direction: column;
            gap: 20px;
        }

        /* 工具栏 */
        .toolbar {
            background: white;
            padding: 20px 25px;
            border-radius: 15px;
            display: flex;
            align-items: center;
            justify-content: space-between;
            box-shadow: 0 5px 20px rgba(0,0,0,0.08);
        }
        .toolbar-left {
            display: flex;
            gap: 15px;
            align-items: center;
        }
        .toolbar-right {
            display: flex;
            gap: 10px;
        }
        .btn {
            padding: 10px 20px;
            border: 2px solid #e9ecef;
            background: white;
            border-radius: 10px;
            cursor: pointer;
            font-weight: 500;
            transition: all 0.3s;
            display: flex;
            align-items: center;
            gap: 8px;
        }
        .btn:hover {
            border-color: #007bff;
            color: #007bff;
            transform: translateY(-2px);
        }
        .btn-primary {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            border-color: transparent;
        }
        .btn-primary:hover {
            transform: translateY(-2px);
            box-shadow: 0 10px 25px rgba(102,126,234,0.3);
        }
        .btn-danger {
            border-color: #dc3545;
            color: #dc3545;
        }
        .btn-danger:hover {
            background: #dc3545;
            color: white;
        }

        /* 域名列表 */
        .domain-list {
            background: white;
            border-radius: 15px;
            overflow: hidden;
            box-shadow: 0 5px 20px rgba(0,0,0,0.08);
            flex: 1;
        }
        .domain-item {
            padding: 20px 25px;
            border-bottom: 1px solid #f1f3f4;
            display: flex;
            align-items: center;
            justify-content: space-between;
            transition: all 0.3s;
        }
        .domain-item:hover {
            background: #f8f9fa;
            transform: translateX(5px);
        }
        .domain-item:last-child {
            border-bottom: none;
        }
        .domain-info {
            display: flex;
            align-items: center;
            gap: 15px;
        }
        .domain-icon {
            width: 45px;
            height: 45px;
            border-radius: 50%;
            background: linear-gradient(135deg, #28a745 0%, #20c997 100%);
            display: flex;
            align-items: center;
            justify-content: center;
            color: white;
            font-weight: 600;
            font-size: 16px;
        }
        .domain-details {
            display: flex;
            flex-direction: column;
            gap: 3px;
        }
        .domain-name {
            font-weight: 600;
            color: #333;
            font-size: 15px;
        }
        .domain-status {
            padding: 3px 8px;
            border-radius: 12px;
            font-size: 11px;
            font-weight: 600;
            text-transform: uppercase;
        }
        .status-active {
            background: #d4edda;
            color: #155724;
        }
        .status-inactive {
            background: #f8d7da;
            color: #721c24;
        }
        .domain-actions {
            display: flex;
            gap: 10px;
        }
        .btn-sm {
            padding: 6px 12px;
            font-size: 12px;
            border-radius: 8px;
        }

        /* 空状态 */
        .empty-state {
            text-align: center;
            padding: 80px 20px;
            color: #8e9aaf;
        }
        .empty-state .icon {
            font-size: 4rem;
            margin-bottom: 20px;
            opacity: 0.5;
        }
        .empty-state h3 {
            font-size: 1.2rem;
            margin-bottom: 10px;
            color: #5a6c7d;
        }

        /* 模态框样式 */
        .modal {
            display: none;
            position: fixed;
            z-index: 1000;
            left: 0;
            top: 0;
            width: 100%;
            height: 100%;
            background-color: rgba(0,0,0,0.5);
            backdrop-filter: blur(5px);
        }
        .modal-content {
            background-color: white;
            margin: 10% auto;
            padding: 30px;
            border-radius: 15px;
            width: 90%;
            max-width: 500px;
            box-shadow: 0 20px 40px rgba(0,0,0,0.2);
        }
        .modal-header {
            margin-bottom: 25px;
            padding-bottom: 15px;
            border-bottom: 2px solid #f1f3f4;
        }
        .modal-header h3 {
            color: #333;
            font-weight: 600;
        }
        .close {
            color: #aaa;
            float: right;
            font-size: 28px;
            font-weight: bold;
            cursor: pointer;
            line-height: 1;
        }
        .close:hover {
            color: #333;
        }
        .form-group {
            margin-bottom: 20px;
        }
        .form-group label {
            display: block;
            margin-bottom: 8px;
            font-weight: 600;
            color: #333;
        }
        .form-group input {
            width: 100%;
            padding: 12px 15px;
            border: 2px solid #e9ecef;
            border-radius: 8px;
            font-size: 14px;
            transition: all 0.3s;
        }
        .form-group input:focus {
            outline: none;
            border-color: #667eea;
            box-shadow: 0 0 0 3px rgba(102,126,234,0.1);
        }

        /* 响应式设计 */
        @media (max-width: 1024px) {
            .brand-section {
                flex-direction: column;
                gap: 15px;
                text-align: center;
            }
            .brand-left {
                justify-content: center;
            }
            .brand-info h1 {
                font-size: 2rem;
            }
            .container {
                flex-direction: column;
                padding: 10px;
            }
            .sidebar {
                width: 100%;
                order: 2;
            }
            .main {
                order: 1;
                min-height: 60vh;
                width: 100%;
            }
        }

        @media (max-width: 768px) {
            .header {
                padding: 0 15px;
            }
            .header-nav {
                display: none;
            }
            .brand-section {
                padding: 15px;
            }
            .code-mascot {
                font-size: 2rem;
            }
            .brand-info h1 {
                font-size: 1.8rem;
            }
            .brand-info .version {
                font-size: 0.9rem;
            }
            .brand-right {
                font-size: 1rem;
            }
            .sidebar {
                padding: 15px 0;
                margin-bottom: 20px;
            }
            .nav-section {
                margin-bottom: 20px;
            }
            .nav-section-title {
                padding: 0 20px 10px;
                font-size: 11px;
            }
            .nav-item {
                padding: 12px 20px;
                font-size: 14px;
            }
            .nav-item .icon {
                margin-right: 8px;
                font-size: 16px;
            }
            .domain-item {
                padding: 15px;
                gap: 15px;
            }
            .domain-icon {
                width: 35px;
                height: 35px;
                font-size: 14px;
            }
        }
    </style>
</head>
<body>
    <div class="header">
        <div class="logo">📧 NBEmail</div>
        <div class="header-nav">
            <a href="/guide">使用指南</a>
            <a href="#about">关于</a>
        </div>
        <div class="user-info">
            <span class="user-email">` + userEmail + `</span>
            <button class="logout-btn" onclick="logout()">退出登录</button>
        </div>
    </div>

    <!-- 品牌区域 -->
    <div class="brand-section">
        <div class="brand-left">
            <div class="code-mascot">👩‍💻</div>
            <div class="brand-info">
                <h1>📧 NBEmail</h1>
                <div class="version">v1.0.0 - 专业邮件管理系统</div>
            </div>
        </div>
        <div class="brand-right">
            <span>欢迎使用，` + userEmail + `</span>
        </div>
    </div>

    <div class="container">
        <div class="sidebar">
            <div class="nav-section">
                <div class="nav-section-title">邮箱</div>
                <a href="/inbox" class="nav-item">
                    <span class="icon">📥</span>收件箱
                </a>
                <a href="/sent" class="nav-item">
                    <span class="icon">📤</span>已发送
                </a>
                <a href="/compose" class="nav-item">
                    <span class="icon">✏️</span>写邮件
                </a>
            </div>
            <div class="nav-divider"></div>
            <div class="nav-section">
                <div class="nav-section-title">管理</div>
                <a href="/users" class="nav-item">
                    <span class="icon">👥</span>用户管理
                </a>
                <a href="/domains" class="nav-item active">
                    <span class="icon">🌐</span>域名管理
                </a>
                <a href="/smtp-configs" class="nav-item">
                    <span class="icon">📮</span>SMTP配置
                </a>
            </div>
        </div>
        <div class="main">
            <div class="toolbar">
                <div class="toolbar-left">
                    <h2 style="margin: 0; color: #333; font-weight: 600;">🌐 域名管理</h2>
                </div>
                <div class="toolbar-right">
                    <button class="btn" onclick="batchVerifyDomains()">
                        <span>🔍</span>批量验证
                    </button>
                    <button class="btn btn-primary" onclick="showCreateDomainModal()">
                        <span>➕</span>添加域名
                    </button>
                    <button class="btn" onclick="refreshDomains()">
                        <span>🔄</span>刷新
                    </button>
                </div>
            </div>
            <div class="domain-list" id="domainList">
                <div class="empty-state">
                    <div class="icon">🌐</div>
                    <h3>正在加载域名...</h3>
                    <p>请稍候，正在获取域名列表</p>
                </div>
            </div>
        </div>
    </div>

    <!-- 添加域名模态框 -->
    <div id="createDomainModal" class="modal">
        <div class="modal-content">
            <div class="modal-header">
                <span class="close" onclick="closeCreateDomainModal()">&times;</span>
                <h3>➕ 添加新域名</h3>
            </div>
            <form id="createDomainForm">
                <div class="form-group">
                    <label for="newDomainName">域名 *</label>
                    <input type="text" id="newDomainName" required placeholder="请输入域名，例如：example.com">
                </div>
                <div style="display: flex; gap: 10px; justify-content: flex-end; margin-top: 25px;">
                    <button type="button" class="btn" onclick="closeCreateDomainModal()">取消</button>
                    <button type="submit" class="btn btn-primary">添加域名</button>
                </div>
            </form>
        </div>
    </div>

    <!-- DNS配置说明模态框 -->
    <div id="dnsInstructionsModal" class="modal" style="display: none; position: fixed; z-index: 1000; left: 0; top: 0; width: 100%; height: 100%; background-color: rgba(0,0,0,0.5); backdrop-filter: blur(5px);">
        <div class="modal-content" style="background-color: white; margin: 3% auto; padding: 30px; border-radius: 15px; width: 95%; max-width: 900px; box-shadow: 0 20px 40px rgba(0,0,0,0.2); max-height: 90vh; overflow-y: auto;">
            <div class="modal-header" style="margin-bottom: 25px; padding-bottom: 15px; border-bottom: 2px solid #f1f3f4; display: flex; justify-content: space-between; align-items: center;">
                <h3 style="color: #333; font-weight: 600; margin: 0;">📋 DNS配置说明</h3>
                <span class="close" onclick="closeDNSInstructionsModal()" style="color: #aaa; font-size: 28px; font-weight: bold; cursor: pointer; line-height: 1;">&times;</span>
            </div>
            <div id="dnsInstructionsContent">
                <div style="text-align: center; padding: 40px; color: #666;">正在加载DNS配置说明...</div>
            </div>
        </div>
    </div>

    <!-- DNS验证结果模态框 -->
    <div id="dnsVerifyModal" class="modal" style="display: none; position: fixed; z-index: 1000; left: 0; top: 0; width: 100%; height: 100%; background-color: rgba(0,0,0,0.5); backdrop-filter: blur(5px);">
        <div class="modal-content" style="background-color: white; margin: 10% auto; padding: 30px; border-radius: 15px; width: 90%; max-width: 700px; box-shadow: 0 20px 40px rgba(0,0,0,0.2); max-height: 80vh; overflow-y: auto;">
            <div class="modal-header" style="margin-bottom: 25px; padding-bottom: 15px; border-bottom: 2px solid #f1f3f4; display: flex; justify-content: space-between; align-items: center;">
                <h3 style="color: #333; font-weight: 600; margin: 0;">🔍 DNS验证结果</h3>
                <span class="close" onclick="closeDNSVerifyModal()" style="color: #aaa; font-size: 28px; font-weight: bold; cursor: pointer; line-height: 1;">&times;</span>
            </div>
            <div id="dnsVerifyContent">
                <div style="text-align: center; padding: 40px; color: #666;">正在验证DNS配置...</div>
            </div>
        </div>
    </div>

    <script>
        let domains = [];

        async function loadDomains() {
            try {
                const response = await fetch('/api/domains');
                const result = await response.json();
                if (result.success) {
                    domains = result.data || [];
                    renderDomains();
                } else {
                    document.getElementById('domainList').innerHTML = '<div class="empty-state">加载域名失败</div>';
                }
            } catch (error) {
                document.getElementById('domainList').innerHTML = '<div class="empty-state">加载域名失败</div>';
            }
        }

        function renderDomains() {
            const domainList = document.getElementById('domainList');
            if (domains.length === 0) {
                domainList.innerHTML = ` + "`" + `
                    <div class="empty-state">
                        <div class="icon">🌍</div>
                        <h3>暂无域名</h3>
                        <p>系统中还没有配置任何域名</p>
                    </div>
                ` + "`" + `;
                return;
            }

            domainList.innerHTML = domains.map(domain => {
                const domainInitial = domain.name.charAt(0).toUpperCase();
                const statusClass = domain.is_active ? 'status-active' : 'status-inactive';
                const statusText = domain.is_active ? '活跃' : '停用';

                // DNS验证状态
                const dnsStatusClass = domain.dns_verified ? 'status-verified' : 'status-unverified';
                const dnsStatusText = domain.dns_verified ? '✅ 已验证' : '❌ 未验证';
                const dnsStatusColor = domain.dns_verified ? '#28a745' : '#dc3545';

                return ` + "`" + `
                    <div class="domain-item">
                        <div class="domain-info">
                            <div class="domain-icon">${domainInitial}</div>
                            <div class="domain-details">
                                <div class="domain-name">${domain.name}</div>
                                <div style="display: flex; gap: 10px; align-items: center; margin-top: 5px;">
                                    <span class="domain-status ${statusClass}">${statusText}</span>
                                    <span class="dns-status" style="color: ${dnsStatusColor}; font-size: 12px; font-weight: 600;">${dnsStatusText}</span>
                                </div>
                                ${domain.mx_record ? ` + "`" + `<div style="font-size: 11px; color: #666; margin-top: 3px;">MX: ${domain.mx_record}</div>` + "`" + ` : ''}
                            </div>
                        </div>
                        <div class="domain-actions" style="display: flex; gap: 8px;">
                            <button class="btn btn-sm" onclick="showDNSInstructions(${domain.id}, '${domain.name}')" title="DNS配置说明">
                                <span>📋</span>配置
                            </button>
                            <button class="btn btn-sm btn-primary" onclick="verifyDomain(${domain.id}, '${domain.name}')" title="验证DNS">
                                <span>🔍</span>验证
                            </button>
                            <button class="btn btn-sm btn-danger" onclick="deleteDomain(${domain.id})">
                                <span>🗑️</span>删除
                            </button>
                        </div>
                    </div>
                ` + "`" + `;
            }).join('');
        }

        function showCreateDomainModal() {
            document.getElementById('createDomainModal').style.display = 'block';
            document.getElementById('newDomainName').focus();
        }

        function closeCreateDomainModal() {
            document.getElementById('createDomainModal').style.display = 'none';
            document.getElementById('createDomainForm').reset();
        }

        // 处理创建域名表单提交
        document.getElementById('createDomainForm').addEventListener('submit', async (e) => {
            e.preventDefault();

            const name = document.getElementById('newDomainName').value;

            const submitBtn = e.target.querySelector('button[type="submit"]');
            const originalText = submitBtn.innerHTML;

            // 显示加载状态
            submitBtn.innerHTML = '<span>⏳</span>添加中...';
            submitBtn.disabled = true;

            try {
                const response = await fetch('/api/domains', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ name })
                });

                const result = await response.json();
                if (result.success) {
                    // 显示成功状态
                    submitBtn.innerHTML = '<span>✅</span>添加成功';
                    submitBtn.style.background = '#28a745';

                    showNotification('域名添加成功！', 'success');

                    // 2秒后关闭模态框并刷新列表
                    setTimeout(() => {
                        closeCreateDomainModal();
                        loadDomains();
                        submitBtn.innerHTML = originalText;
                        submitBtn.style.background = '';
                        submitBtn.disabled = false;
                    }, 2000);
                } else {
                    throw new Error(result.message || '添加失败');
                }
            } catch (error) {
                // 显示错误状态
                submitBtn.innerHTML = '<span>❌</span>添加失败';
                submitBtn.style.background = '#dc3545';

                showNotification('添加失败: ' + error.message, 'error');

                // 3秒后恢复按钮
                setTimeout(() => {
                    submitBtn.innerHTML = originalText;
                    submitBtn.style.background = '';
                    submitBtn.disabled = false;
                }, 3000);
            }
        });

        async function deleteDomain(domainId) {
            if (!confirm('确定要删除这个域名吗？此操作不可撤销！')) return;

            try {
                const response = await fetch(` + "`" + `/api/domains/${domainId}` + "`" + `, { method: 'DELETE' });
                const result = await response.json();
                if (result.success) {
                    showNotification('域名删除成功', 'success');
                    loadDomains();
                } else {
                    showNotification('删除失败: ' + result.message, 'error');
                }
            } catch (error) {
                showNotification('删除失败，请重试', 'error');
            }
        }

        function refreshDomains() {
            showNotification('正在刷新域名列表...', 'info');
            loadDomains();
        }

        // 显示通知
        function showNotification(message, type = 'info') {
            // 创建通知元素
            const notification = document.createElement('div');
            notification.style.cssText = ` + "`" + `
                position: fixed;
                top: 20px;
                right: 20px;
                padding: 15px 25px;
                border-radius: 10px;
                color: white;
                font-weight: 500;
                z-index: 10000;
                transform: translateX(400px);
                transition: all 0.3s ease;
                box-shadow: 0 5px 20px rgba(0,0,0,0.2);
            ` + "`" + `;

            // 设置不同类型的样式
            switch(type) {
                case 'success':
                    notification.style.background = 'linear-gradient(135deg, #28a745, #20c997)';
                    break;
                case 'error':
                    notification.style.background = 'linear-gradient(135deg, #dc3545, #e74c3c)';
                    break;
                case 'info':
                default:
                    notification.style.background = 'linear-gradient(135deg, #667eea, #764ba2)';
            }

            notification.textContent = message;
            document.body.appendChild(notification);

            // 显示动画
            setTimeout(() => {
                notification.style.transform = 'translateX(0)';
            }, 100);

            // 自动隐藏
            setTimeout(() => {
                notification.style.transform = 'translateX(400px)';
                setTimeout(() => {
                    if (document.body.contains(notification)) {
                        document.body.removeChild(notification);
                    }
                }, 300);
            }, 3000);
        }

        async function logout() {
            if (confirm('确定要退出登录吗？')) {
                try {
                    await fetch('/api/logout', { method: 'POST' });
                    window.location.href = '/login';
                } catch (error) {
                    window.location.href = '/login';
                }
            }
        }

        // 点击模态框外部关闭
        window.onclick = function(event) {
            const modal = document.getElementById('createDomainModal');
            if (event.target == modal) {
                closeCreateDomainModal();
            }
        }

        // DNS验证相关函数
        async function showDNSInstructions(domainId, domainName) {
            document.getElementById('dnsInstructionsModal').style.display = 'block';
            document.getElementById('dnsInstructionsContent').innerHTML = '<div style="text-align: center; padding: 40px; color: #666;">正在加载DNS配置说明...</div>';

            try {
                const response = await fetch(` + "`" + `/api/domains/${domainId}/dns-instructions` + "`" + `);
                const result = await response.json();
                if (result.success) {
                    renderDNSInstructions(result.data, domainId);
                } else {
                    document.getElementById('dnsInstructionsContent').innerHTML = '<div style="text-align: center; padding: 40px; color: #dc3545;">加载失败: ' + result.message + '</div>';
                }
            } catch (error) {
                document.getElementById('dnsInstructionsContent').innerHTML = '<div style="text-align: center; padding: 40px; color: #dc3545;">加载失败，请重试</div>';
            }
        }

        function renderDNSInstructions(data, domainId) {
            const content = ` + "`" + `
                <div style="margin-bottom: 25px;">
                    <h4 style="color: #333; margin-bottom: 15px;">域名: ${data.domain}</h4>
                    <p style="color: #666; margin-bottom: 20px;">服务器IP: <strong>${data.server_ip}</strong></p>
                    <div style="background: #fff3cd; border: 1px solid #ffeaa7; border-radius: 8px; padding: 15px; margin-bottom: 20px;">
                        <strong>⚠️ 重要提示：</strong> 请在您的域名DNS管理面板中添加以下记录，然后等待DNS传播完成（通常需要几分钟到几小时）。
                    </div>
                </div>

                <div style="display: grid; gap: 20px; margin-bottom: 25px;">
                    <div style="border: 2px solid #e9ecef; border-radius: 10px; padding: 20px;">
                        <h5 style="color: #007bff; margin-bottom: 10px;">📧 MX记录（邮件交换记录）</h5>
                        <div style="background: #f8f9fa; padding: 15px; border-radius: 8px; font-family: monospace; margin-bottom: 10px;">
                            <div><strong>类型:</strong> MX</div>
                            <div><strong>名称:</strong> ${data.instructions.mx_record.name}</div>
                            <div><strong>值:</strong> ${data.instructions.mx_record.value}</div>
                            <div><strong>优先级:</strong> ${data.instructions.mx_record.priority}</div>
                            <div><strong>TTL:</strong> ${data.instructions.mx_record.ttl}</div>
                        </div>
                        <p style="color: #666; font-size: 14px;">${data.instructions.mx_record.description}</p>
                    </div>

                    <div style="border: 2px solid #e9ecef; border-radius: 10px; padding: 20px;">
                        <h5 style="color: #28a745; margin-bottom: 10px;">🌐 A记录（域名解析记录）</h5>
                        <div style="background: #f8f9fa; padding: 15px; border-radius: 8px; font-family: monospace; margin-bottom: 10px;">
                            <div><strong>类型:</strong> A</div>
                            <div><strong>名称:</strong> ${data.instructions.a_record.name}</div>
                            <div><strong>值:</strong> ${data.instructions.a_record.value}</div>
                            <div><strong>TTL:</strong> ${data.instructions.a_record.ttl}</div>
                        </div>
                        <p style="color: #666; font-size: 14px;">${data.instructions.a_record.description}</p>
                    </div>

                    <div style="border: 2px solid #e9ecef; border-radius: 10px; padding: 20px;">
                        <h5 style="color: #ffc107; margin-bottom: 10px;">📝 TXT记录（SPF记录，可选）</h5>
                        <div style="background: #f8f9fa; padding: 15px; border-radius: 8px; font-family: monospace; margin-bottom: 10px;">
                            <div><strong>类型:</strong> TXT</div>
                            <div><strong>名称:</strong> ${data.instructions.txt_record.name}</div>
                            <div><strong>值:</strong> ${data.instructions.txt_record.value}</div>
                            <div><strong>TTL:</strong> ${data.instructions.txt_record.ttl}</div>
                        </div>
                        <p style="color: #666; font-size: 14px;">${data.instructions.txt_record.description}</p>
                    </div>
                </div>

                <div style="margin-bottom: 25px;">
                    <h5 style="color: #333; margin-bottom: 15px;">📋 配置步骤</h5>
                    <ol style="color: #666; line-height: 1.6;">
                        ${data.steps.map(step => ` + "`" + `<li>${step}</li>` + "`" + `).join('')}
                    </ol>
                </div>

                <div style="margin-bottom: 25px;">
                    <h5 style="color: #333; margin-bottom: 15px;">🔗 常用DNS服务商</h5>
                    <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 10px;">
                        ${Object.entries(data.common_providers).map(([name, url]) => ` + "`" + `
                            <a href="${url}" target="_blank" style="display: block; padding: 10px 15px; background: #f8f9fa; border-radius: 8px; text-decoration: none; color: #007bff; border: 1px solid #e9ecef; transition: all 0.3s;" onmouseover="this.style.background='#e9ecef'" onmouseout="this.style.background='#f8f9fa'">
                                ${name} →
                            </a>
                        ` + "`" + `).join('')}
                    </div>
                </div>

                <div style="text-align: center; margin-top: 30px;">
                    <button class="btn btn-primary" onclick="verifyDomain(${domainId}, '${data.domain}')">
                        <span>🔍</span>验证DNS配置
                    </button>
                </div>
            ` + "`" + `;

            document.getElementById('dnsInstructionsContent').innerHTML = content;
        }

        function closeDNSInstructionsModal() {
            document.getElementById('dnsInstructionsModal').style.display = 'none';
        }

        async function verifyDomain(domainId, domainName) {
            // 如果从DNS说明模态框调用，先关闭它
            closeDNSInstructionsModal();

            document.getElementById('dnsVerifyModal').style.display = 'block';
            document.getElementById('dnsVerifyContent').innerHTML = ` + "`" + `
                <div style="text-align: center; padding: 40px;">
                    <div style="font-size: 3rem; margin-bottom: 15px;">🔍</div>
                    <h4 style="color: #333; margin-bottom: 10px;">正在验证 ${domainName}</h4>
                    <p style="color: #666;">请稍候，正在检查DNS配置...</p>
                    <div style="margin-top: 20px;">
                        <div style="display: inline-block; width: 40px; height: 40px; border: 4px solid #f3f3f3; border-top: 4px solid #007bff; border-radius: 50%; animation: spin 1s linear infinite;"></div>
                    </div>
                </div>
                <style>
                    @keyframes spin {
                        0% { transform: rotate(0deg); }
                        100% { transform: rotate(360deg); }
                    }
                </style>
            ` + "`" + `;

            try {
                const response = await fetch(` + "`" + `/api/domains/${domainId}/verify` + "`" + `, { method: 'POST' });
                const result = await response.json();
                if (result.success) {
                    renderVerifyResult(result.data);
                    // 刷新域名列表以更新验证状态
                    loadDomains();
                } else {
                    document.getElementById('dnsVerifyContent').innerHTML = ` + "`" + `
                        <div style="text-align: center; padding: 40px;">
                            <div style="font-size: 3rem; margin-bottom: 15px;">❌</div>
                            <h4 style="color: #dc3545;">验证失败</h4>
                            <p style="color: #666;">${result.message}</p>
                        </div>
                    ` + "`" + `;
                }
            } catch (error) {
                document.getElementById('dnsVerifyContent').innerHTML = ` + "`" + `
                    <div style="text-align: center; padding: 40px;">
                        <div style="font-size: 3rem; margin-bottom: 15px;">❌</div>
                        <h4 style="color: #dc3545;">验证失败</h4>
                        <p style="color: #666;">网络错误，请重试</p>
                    </div>
                ` + "`" + `;
            }
        }

        function renderVerifyResult(data) {
            const result = data.verification_result;
            const statusIcon = result.success ? '✅' : '❌';
            const statusColor = result.success ? '#28a745' : '#dc3545';

            const content = ` + "`" + `
                <div style="text-align: center; margin-bottom: 30px;">
                    <div style="font-size: 4rem; margin-bottom: 15px;">${statusIcon}</div>
                    <h4 style="color: ${statusColor}; margin-bottom: 10px;">${result.success ? 'DNS验证成功' : 'DNS验证失败'}</h4>
                    <p style="color: #666; margin-bottom: 20px;">${result.message}</p>
                </div>

                <div style="background: #f8f9fa; border-radius: 10px; padding: 20px; margin-bottom: 20px;">
                    <h5 style="color: #333; margin-bottom: 15px;">验证详情</h5>
                    <div style="display: grid; gap: 15px;">
                        <div>
                            <strong>域名:</strong> ${data.domain_name}
                        </div>
                        <div>
                            <strong>服务器IP:</strong> ${data.server_ip}
                        </div>
                        <div>
                            <strong>DNS验证状态:</strong>
                            <span style="color: ${statusColor}; font-weight: 600;">${result.success ? '通过' : '失败'}</span>
                        </div>
                        <div>
                            <strong>MX记录:</strong> ${result.has_mx ? '已配置' : '未配置'}
                        </div>
                        <div>
                            <strong>指向本服务器:</strong> ${result.points_to_us ? '是' : '否'}
                        </div>
                    </div>
                </div>

                ${result.mx_records.length > 0 ? ` + "`" + `
                    <div style="background: #f8f9fa; border-radius: 10px; padding: 20px; margin-bottom: 20px;">
                        <h5 style="color: #333; margin-bottom: 15px;">检测到的MX记录</h5>
                        <ul style="margin: 0; padding-left: 20px;">
                            ${result.mx_records.map(mx => ` + "`" + `<li style="margin-bottom: 5px;">${mx}</li>` + "`" + `).join('')}
                        </ul>
                    </div>
                ` + "`" + ` : ''}

                ${result.a_records.length > 0 ? ` + "`" + `
                    <div style="background: #f8f9fa; border-radius: 10px; padding: 20px; margin-bottom: 20px;">
                        <h5 style="color: #333; margin-bottom: 15px;">检测到的A记录</h5>
                        <ul style="margin: 0; padding-left: 20px;">
                            ${result.a_records.map(a => ` + "`" + `<li style="margin-bottom: 5px;">${a}</li>` + "`" + `).join('')}
                        </ul>
                    </div>
                ` + "`" + ` : ''}

                <div style="text-align: center; margin-top: 30px;">
                    ${!result.success ? ` + "`" + `
                        <button class="btn btn-primary" onclick="showDNSInstructions(${data.domain_id}, '${data.domain_name}')" style="margin-right: 10px;">
                            <span>📋</span>查看配置说明
                        </button>
                    ` + "`" + ` : ''}
                    <button class="btn" onclick="closeDNSVerifyModal()">关闭</button>
                </div>
            ` + "`" + `;

            document.getElementById('dnsVerifyContent').innerHTML = content;
        }

        function closeDNSVerifyModal() {
            document.getElementById('dnsVerifyModal').style.display = 'none';
        }

        async function batchVerifyDomains() {
            if (!confirm('确定要批量验证所有域名的DNS配置吗？')) return;

            try {
                const response = await fetch('/api/domains/batch-verify', { method: 'POST' });
                const result = await response.json();
                if (result.success) {
                    showNotification(` + "`" + `批量验证完成：成功 ${result.data.success_count} 个，失败 ${result.data.fail_count} 个` + "`" + `, 'success');
                    loadDomains(); // 刷新列表
                } else {
                    showNotification('批量验证失败: ' + result.message, 'error');
                }
            } catch (error) {
                showNotification('批量验证失败，请重试', 'error');
            }
        }

        // 页面加载时获取域名列表
        loadDomains();
    </script>
</body>
</html>`
}

// generateSMTPConfigsPageTemplate 生成SMTP配置管理页面模板
func (s *Server) generateSMTPConfigsPageTemplate(userEmail string, isAdmin bool) string {
	// 使用统一的基础模板
	tmpl := s.getBasePageTemplate(userEmail, isAdmin, "smtp-configs")

	// 定义SMTP配置页面的主内容
	smtpContent := `            <div class="smtp-configs-container" id="configsContainer">
                <div class="page-description" style="padding: 20px; background: #f8f9fa; border-radius: 8px; margin-bottom: 20px; border-left: 4px solid #667eea;">
                    <p style="margin-bottom: 10px;">
                        系统会自动检测您的域名邮箱并生成SMTP配置，包括智能推荐的用户名和强密码。点击"自动配置"扫描现有邮箱，系统将自动生成完整配置。
                    </p>
                    <p style="margin-bottom: 10px;">
                        <strong>🎯 新功能：</strong>支持一键生成用户名和强密码，无需手动填写！
                    </p>
                    <p>
                        <strong>工作原理：</strong>用什么域名的邮箱发件，就自动使用对应域名的SMTP服务器发送。
                    </p>
                </div>

                <div class="empty-state">
                    <div class="empty-state-icon">🔧</div>
                    <div class="empty-state-title">开始配置多域名SMTP</div>
                    <div class="empty-state-description">
                        系统可以自动检测您的域名邮箱并生成SMTP配置<br>
                        点击"自动配置"开始，然后编辑配置添加认证信息
                    </div>
                    <div style="display: flex; gap: 15px; justify-content: center; margin-top: 20px;">
                        <button class="btn btn-primary" onclick="autoConfigSMTP()">
                            <span>🔧</span>自动配置
                        </button>
                        <button class="btn btn-primary" onclick="openAddConfigModal()">
                            <span>➕</span>手动添加
                        </button>
                    </div>
                </div>
            </div>`

	// 替换页面标题
	tmpl = strings.Replace(tmpl, "📥 收件箱", "📮 SMTP配置管理", 1)

	// 替换工具栏按钮
	oldToolbarRight := `                <div class="toolbar-right">
                    <button class="btn" onclick="showMailboxManager()">
                        <span>📮</span>邮箱管理
                    </button>
                    <button class="btn" onclick="refreshEmails()">
                        <span>🔄</span>刷新
                    </button>
                    <button class="btn btn-danger" onclick="deleteSelected()">
                        <span>🗑️</span>删除
                    </button>
                    <button class="btn" onclick="markAsRead()">
                        <span>✅</span>标记已读
                    </button>
                </div>`

	newToolbarRight := `                <div class="toolbar-right">
                    <button class="btn btn-primary" onclick="autoConfigSMTP()">
                        <span>🔧</span>自动配置
                    </button>
                    <button class="btn btn-primary" onclick="openAddConfigModal()">
                        <span>➕</span>手动添加
                    </button>
                </div>`

	tmpl = strings.Replace(tmpl, oldToolbarRight, newToolbarRight, 1)

	// 替换主内容区域
	oldEmailList := `            <div class="email-list" id="emailList">
                <div class="empty-state">
                    <div class="icon">📬</div>
                    <h3>正在加载邮件...</h3>
                    <p>请稍候，正在获取您的邮件</p>
                </div>
            </div>`

	tmpl = strings.Replace(tmpl, oldEmailList, smtpContent, 1)

	// 添加SMTP配置相关的样式
	// 添加SMTP配置相关的样式
	additionalStyles := `
        .configs-grid {
            display: grid;
            gap: 20px;
            margin-bottom: 30px;
        }

        .config-card {
            background: white;
            border: 1px solid #e0e0e0;
            border-radius: 12px;
            padding: 20px;
            box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
            transition: all 0.3s ease;
        }

        .config-card:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
        }

        .config-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 15px;
        }

        .config-domain {
            font-size: 18px;
            font-weight: 600;
            color: #333;
        }

        .config-badge {
            padding: 4px 12px;
            border-radius: 20px;
            font-size: 12px;
            font-weight: 500;
        }

        .badge-default {
            background: #e3f2fd;
            color: #1976d2;
        }

        .badge-custom {
            background: #f3e5f5;
            color: #7b1fa2;
        }

        .config-details {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 15px;
            margin-bottom: 15px;
        }

        .config-item {
            display: flex;
            flex-direction: column;
            gap: 4px;
        }

        .config-label {
            font-size: 12px;
            color: #666;
            font-weight: 500;
            text-transform: uppercase;
        }

        .config-value {
            font-size: 14px;
            color: #333;
            font-family: 'Monaco', 'Menlo', monospace;
        }

        .config-actions {
            display: flex;
            gap: 10px;
            justify-content: flex-end;
        }

        .btn-small {
            padding: 6px 12px;
            font-size: 12px;
        }

        .modal {
            display: none;
            position: fixed;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            background: rgba(0, 0, 0, 0.5);
            z-index: 1000;
        }

        .modal-content {
            position: absolute;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -50%);
            background: white;
            border-radius: 15px;
            padding: 30px;
            width: 90%;
            max-width: 500px;
            box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
        }

        .modal-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 20px;
        }

        .modal-title {
            font-size: 20px;
            font-weight: 600;
            color: #333;
        }

        .close-btn {
            background: none;
            border: none;
            font-size: 24px;
            cursor: pointer;
            color: #666;
            padding: 0;
            width: 30px;
            height: 30px;
            display: flex;
            align-items: center;
            justify-content: center;
            border-radius: 50%;
            transition: all 0.3s ease;
        }

        .close-btn:hover {
            background: #f0f0f0;
            color: #333;
        }

        .form-group {
            margin-bottom: 20px;
        }

        .form-label {
            display: block;
            margin-bottom: 8px;
            font-weight: 500;
            color: #333;
        }

        .form-input {
            width: 100%;
            padding: 12px;
            border: 1px solid #ddd;
            border-radius: 8px;
            font-size: 14px;
            transition: all 0.3s ease;
        }

        .form-input:focus {
            outline: none;
            border-color: #667eea;
            box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
        }

        .form-checkbox {
            display: flex;
            align-items: center;
            gap: 8px;
        }

        .form-checkbox input {
            width: auto;
        }

        .notification {
            position: fixed;
            top: 20px;
            right: 20px;
            padding: 15px 20px;
            border-radius: 8px;
            color: white;
            font-weight: 500;
            z-index: 1001;
            transform: translateX(400px);
            transition: all 0.3s ease;
        }

        .notification.show {
            transform: translateX(0);
        }

        .notification.success {
            background: #4caf50;
        }

        .notification.error {
            background: #f44336;
        }

        .empty-state {
            text-align: center;
            padding: 60px 20px;
            color: #666;
        }

        .empty-state-icon {
            font-size: 48px;
            margin-bottom: 20px;
            opacity: 0.5;
        }

        .empty-state-title {
            font-size: 20px;
            font-weight: 600;
            margin-bottom: 10px;
        }

        .empty-state-description {
            font-size: 14px;
            margin-bottom: 20px;
        }

        .config-card {
            background: white;
            border-radius: 12px;
            padding: 20px;
            margin-bottom: 20px;
            box-shadow: 0 2px 8px rgba(0,0,0,0.1);
            border: 1px solid #e0e0e0;
            transition: all 0.3s ease;
        }

        .config-card:hover {
            box-shadow: 0 4px 16px rgba(0,0,0,0.15);
            transform: translateY(-2px);
        }

        .config-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 15px;
            padding-bottom: 15px;
            border-bottom: 1px solid #f0f0f0;
        }

        .config-domain {
            font-size: 18px;
            font-weight: 600;
            color: #333;
        }

        .config-badge {
            padding: 4px 12px;
            border-radius: 20px;
            font-size: 12px;
            font-weight: 500;
        }

        .badge-default {
            background: #e3f2fd;
            color: #1976d2;
        }

        .badge-custom {
            background: #f3e5f5;
            color: #7b1fa2;
        }

        .config-details {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 15px;
            margin-bottom: 20px;
        }

        .config-item {
            display: flex;
            flex-direction: column;
            gap: 5px;
        }

        .config-label {
            font-size: 12px;
            color: #666;
            font-weight: 500;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }

        .config-value {
            font-size: 14px;
            color: #333;
            font-weight: 500;
        }

        .config-actions {
            display: flex;
            gap: 10px;
            justify-content: flex-end;
            padding-top: 15px;
            border-top: 1px solid #f0f0f0;
        }

        .btn-small {
            padding: 6px 12px;
            font-size: 12px;
        }

        .btn-secondary {
            background: #6c757d;
            color: white;
            border: none;
        }

        .btn-secondary:hover {
            background: #5a6268;
        }

        .btn-danger {
            background: #dc3545;
            color: white;
            border: none;
        }

        .btn-danger:hover {
            background: #c82333;
        }

        /* 模态框样式 */
        .modal-overlay {
            position: fixed;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background: rgba(0, 0, 0, 0.5);
            display: flex;
            align-items: center;
            justify-content: center;
            z-index: 1000;
        }

        .modal-content {
            background: white;
            border-radius: 12px;
            max-width: 600px;
            width: 90%;
            max-height: 90vh;
            overflow-y: auto;
            box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
        }

        .modal-header {
            padding: 20px;
            border-bottom: 1px solid #e0e0e0;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }

        .modal-header h3 {
            margin: 0;
            color: #333;
        }

        .modal-close {
            background: none;
            border: none;
            font-size: 24px;
            cursor: pointer;
            color: #666;
            padding: 0;
            width: 30px;
            height: 30px;
            display: flex;
            align-items: center;
            justify-content: center;
        }

        .modal-close:hover {
            color: #333;
        }

        .modal-body {
            padding: 20px;
        }

        .modal-footer {
            padding: 20px;
            border-top: 1px solid #e0e0e0;
            display: flex;
            gap: 10px;
            justify-content: flex-end;
        }

        /* 表单样式 */
        .form-group {
            margin-bottom: 20px;
        }

        .form-group label {
            display: block;
            margin-bottom: 5px;
            font-weight: 500;
            color: #333;
        }

        .form-control {
            width: 100%;
            padding: 10px;
            border: 1px solid #ddd;
            border-radius: 6px;
            font-size: 14px;
            transition: border-color 0.3s ease;
        }

        .form-control:focus {
            outline: none;
            border-color: #007bff;
            box-shadow: 0 0 0 2px rgba(0, 123, 255, 0.25);
        }

        .form-control-readonly {
            background: #f8f9fa;
            color: #6c757d;
        }

        .form-help {
            display: block;
            margin-top: 5px;
            font-size: 12px;
            color: #666;
        }

        .checkbox-label {
            display: flex;
            align-items: center;
            gap: 8px;
            cursor: pointer;
        }

        .checkbox-label input[type="checkbox"] {
            margin: 0;
        }

        /* DNS验证结果样式 */
        .dns-verification-modal {
            max-width: 800px;
        }

        .verification-summary {
            background: #f8f9fa;
            padding: 15px;
            border-radius: 8px;
            margin-bottom: 20px;
        }

        .summary-item {
            margin-bottom: 8px;
        }

        .status-badge {
            padding: 4px 8px;
            border-radius: 4px;
            font-size: 12px;
            font-weight: 500;
        }

        .status-success {
            background: #d4edda;
            color: #155724;
        }

        .status-warning {
            background: #fff3cd;
            color: #856404;
        }

        .verification-checks h4 {
            margin-bottom: 15px;
            color: #333;
        }

        .check-item {
            border: 1px solid #e0e0e0;
            border-radius: 8px;
            padding: 15px;
            margin-bottom: 10px;
        }

        .check-success {
            border-color: #28a745;
            background: #f8fff9;
        }

        .check-failed {
            border-color: #dc3545;
            background: #fff8f8;
        }

        .check-header {
            display: flex;
            align-items: center;
            gap: 10px;
            margin-bottom: 8px;
        }

        .check-icon {
            font-size: 16px;
        }

        .check-name {
            font-weight: 600;
            color: #333;
        }

        .check-description {
            font-size: 13px;
            color: #666;
            margin-bottom: 8px;
        }

        .check-message {
            font-size: 14px;
            color: #333;
        }

        .check-details {
            font-size: 12px;
            color: #666;
            margin-top: 5px;
            font-family: monospace;
            background: #f8f9fa;
            padding: 5px 8px;
            border-radius: 4px;
        }

        .verification-message {
            background: #e3f2fd;
            border: 1px solid #2196f3;
            border-radius: 8px;
            padding: 15px;
            margin-top: 20px;
            color: #1976d2;
        }

        /* 移动端样式优化 */
        @media (max-width: 768px) {
            .empty-state {
                padding: 30px 15px;
            }

            .empty-state-title {
                font-size: 1.5rem;
            }

            .empty-state-description {
                font-size: 14px;
                margin: 15px 0 25px 0;
            }

            /* 移动端按钮样式优化 */
            .empty-state > div[style*="display: flex"] {
                flex-direction: column !important;
                gap: 12px !important;
                position: sticky;
                bottom: 0;
                background: white;
                padding: 20px 15px;
                margin: 20px -15px -30px -15px;
                border-top: 2px solid #f1f3f4;
                box-shadow: 0 -2px 10px rgba(0,0,0,0.1);
            }

            .empty-state .btn {
                width: 100%;
                padding: 18px 20px;
                font-size: 16px;
                font-weight: 600;
                justify-content: center;
            }

            .config-card {
                margin: 0 -15px 20px -15px;
                border-radius: 0;
                border-left: none;
                border-right: none;
            }

            .config-actions {
                flex-direction: column;
                gap: 8px;
            }

            .config-actions .btn {
                width: 100%;
                justify-content: center;
            }
        }`

	// 将样式插入到模板中
	tmpl = strings.Replace(tmpl, "</style>", additionalStyles+"\n    </style>", 1)

	// 添加SMTP配置相关的JavaScript脚本到模板中
	smtpScript := `
    <script>
        let configs = [];

        // 页面加载时获取配置列表
        document.addEventListener('DOMContentLoaded', function() {
            loadConfigs();
        });

        // 加载配置列表
        async function loadConfigs() {
            try {
                const response = await fetch('/api/smtp-configs');
                const result = await response.json();

                if (result.success) {
                    configs = result.data || [];
                    renderConfigs();
                } else {
                    showNotification('加载配置失败: ' + result.message, 'error');
                }
            } catch (error) {
                showNotification('加载配置失败，请重试', 'error');
            }
        }

        // 渲染配置列表
        function renderConfigs() {
            const container = document.getElementById('configsContainer');

            if (configs.length === 0) {
                container.innerHTML = ` + "`" + `
                    <div class="empty-state">
                        <div class="empty-state-icon">🔧</div>
                        <div class="empty-state-title">开始配置多域名SMTP</div>
                        <div class="empty-state-description">
                            系统可以自动检测您的域名邮箱并生成SMTP配置<br>
                            点击"自动配置"开始，然后编辑配置添加认证信息
                        </div>
                        <div style="display: flex; gap: 15px; justify-content: center; margin-top: 20px;">
                            <button class="btn btn-primary" onclick="autoConfigSMTP()">
                                <span>🔧</span>自动配置
                            </button>
                            <button class="btn btn-primary" onclick="openAddConfigModal()">
                                <span>➕</span>手动添加
                            </button>
                        </div>
                    </div>
                ` + "`" + `;
                return;
            }

            container.innerHTML = configs.map(config => ` + "`" + `
                <div class="config-card">
                    <div class="config-header">
                        <div class="config-domain">${config.domain}</div>
                        <div class="config-badge ${config.is_default ? 'badge-default' : 'badge-custom'}">
                            ${config.is_default ? '默认配置' : '自定义配置'}
                        </div>
                    </div>
                    <div class="config-details">
                        <div class="config-item">
                            <div class="config-label">服务器</div>
                            <div class="config-value">${config.host}</div>
                        </div>
                        <div class="config-item">
                            <div class="config-label">端口</div>
                            <div class="config-value">${config.port}</div>
                        </div>
                        <div class="config-item">
                            <div class="config-label">用户名</div>
                            <div class="config-value">${config.user || '<span style="color: #ff6b6b;">需要配置</span>'}</div>
                        </div>
                        <div class="config-item">
                            <div class="config-label">TLS</div>
                            <div class="config-value">${config.tls ? '启用' : '禁用'}</div>
                        </div>
                        <div class="config-item">
                            <div class="config-label">状态</div>
                            <div class="config-value">${config.user && config.password !== '***' ? '<span style="color: #4caf50;">✅ 已配置</span>' : '<span style="color: #ff9800;">⚠️ 需要认证信息</span>'}</div>
                        </div>
                    </div>
                    <div class="config-actions">
                        ${!config.is_default ? ` + "`" + `<button class="btn btn-primary btn-small" onclick="editConfig('${config.domain}')">编辑</button>` + "`" + ` : ''}
                        <button class="btn btn-secondary btn-small" onclick="verifySMTPDNS('${config.domain}')" title="验证SMTP服务器DNS配置">🔍 DNS验证</button>
                        ${!config.is_default ? ` + "`" + `<button class="btn btn-danger btn-small" onclick="deleteConfig('${config.domain}')">删除</button>` + "`" + ` : ''}
                    </div>
                </div>
            ` + "`" + `).join('');
        }

        // 自动配置SMTP
        async function autoConfigSMTP() {
            if (!confirm('系统将自动检测您的域名邮箱并生成SMTP配置，确定继续吗？')) {
                return;
            }

            try {
                showNotification('正在自动配置SMTP...', 'success');

                const response = await fetch('/api/smtp-configs/auto-config', {
                    method: 'POST'
                });

                const result = await response.json();
                if (result.success) {
                    showNotification(` + "`" + `自动配置完成！检测到 ${result.data.count} 个域名` + "`" + `, 'success');
                    loadConfigs();
                } else {
                    showNotification('自动配置失败: ' + result.message, 'error');
                }
            } catch (error) {
                showNotification('自动配置失败，请重试', 'error');
            }
        }

        // 打开添加配置模态框
        function openAddConfigModal() {
            const modal = document.createElement('div');
            modal.className = 'modal-overlay';
            modal.innerHTML = ` + "`" + `
                <div class="modal-content add-config-modal">
                    <div class="modal-header">
                        <h3>➕ 手动添加SMTP配置</h3>
                        <button class="modal-close" onclick="this.closest('.modal-overlay').remove()">×</button>
                    </div>
                    <div class="modal-body">
                        <form id="addConfigForm">
                            <div class="form-group">
                                <label>域名 *</label>
                                <input type="text" name="domain" class="form-control" required placeholder="例如: example.com">
                                <small class="form-help">输入要配置SMTP的域名</small>
                            </div>
                            <div class="form-group">
                                <label>SMTP服务器 *</label>
                                <input type="text" name="host" class="form-control" required placeholder="例如: mail.example.com">
                                <small class="form-help">SMTP服务器的主机名或IP地址</small>
                            </div>
                            <div class="form-group">
                                <label>端口 *</label>
                                <select name="port" class="form-control" required>
                                    <option value="587">587 (推荐 - STARTTLS)</option>
                                    <option value="465">465 (SSL/TLS)</option>
                                    <option value="25">25 (标准SMTP)</option>
                                    <option value="2525">2525 (备用端口)</option>
                                </select>
                                <small class="form-help">选择SMTP服务器端口</small>
                            </div>
                            <div class="form-group">
                                <label>用户名</label>
                                <input type="text" name="user" class="form-control" placeholder="例如: smtp@example.com">
                                <small class="form-help">SMTP认证用户名，留空表示无需认证</small>
                            </div>
                            <div class="form-group">
                                <label>密码</label>
                                <input type="password" name="password" class="form-control" placeholder="SMTP认证密码">
                                <small class="form-help">SMTP认证密码，留空表示无需认证</small>
                            </div>
                            <div class="form-group">
                                <label class="checkbox-label">
                                    <input type="checkbox" name="tls" checked>
                                    启用TLS加密
                                </label>
                                <small class="form-help">推荐启用以提高安全性</small>
                            </div>
                        </form>
                    </div>
                    <div class="modal-footer">
                        <button class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">取消</button>
                        <button class="btn btn-primary" onclick="saveNewConfig()">添加配置</button>
                    </div>
                </div>
            ` + "`" + `;

            document.body.appendChild(modal);

            // 点击背景关闭模态框
            modal.addEventListener('click', function(e) {
                if (e.target === modal) {
                    modal.remove();
                }
            });

            // 域名输入时自动填充建议的用户名
            const domainInput = modal.querySelector('input[name="domain"]');
            const userInput = modal.querySelector('input[name="user"]');
            const hostInput = modal.querySelector('input[name="host"]');

            domainInput.addEventListener('input', function() {
                const domain = this.value.trim();
                if (domain) {
                    if (!userInput.value) {
                        userInput.value = ` + "`" + `smtp@${domain}` + "`" + `;
                    }
                    if (!hostInput.value) {
                        hostInput.value = ` + "`" + `mail.${domain}` + "`" + `;
                    }
                }
            });
        }

        // 保存新配置
        async function saveNewConfig() {
            const form = document.getElementById('addConfigForm');
            const formData = new FormData(form);

            const configData = {
                domain: formData.get('domain').trim(),
                host: formData.get('host').trim(),
                port: parseInt(formData.get('port')),
                user: formData.get('user').trim(),
                password: formData.get('password'),
                tls: formData.has('tls')
            };

            // 验证必填字段
            if (!configData.domain || !configData.host) {
                showNotification('请填写域名和SMTP服务器', 'error');
                return;
            }

            try {
                const response = await fetch('/api/smtp-configs', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(configData)
                });

                const result = await response.json();
                if (result.success) {
                    showNotification('SMTP配置添加成功！', 'success');
                    document.querySelector('.modal-overlay').remove();
                    loadConfigs(); // 重新加载配置列表
                } else {
                    showNotification('添加失败: ' + result.message, 'error');
                }
            } catch (error) {
                showNotification('添加失败，请重试', 'error');
            }
        }

        // 编辑配置
        function editConfig(domain) {
            const config = configs.find(c => c.domain === domain);
            if (!config) {
                showNotification('未找到配置信息', 'error');
                return;
            }

            const modal = document.createElement('div');
            modal.className = 'modal-overlay';
            modal.innerHTML = ` + "`" + `
                <div class="modal-content edit-config-modal">
                    <div class="modal-header">
                        <h3>✏️ 编辑SMTP配置</h3>
                        <button class="modal-close" onclick="this.closest('.modal-overlay').remove()">×</button>
                    </div>
                    <div class="modal-body">
                        <form id="editConfigForm">
                            <div class="form-group">
                                <label>域名</label>
                                <input type="text" name="domain" value="${config.domain}" readonly class="form-control-readonly">
                            </div>
                            <div class="form-group">
                                <label>SMTP服务器</label>
                                <input type="text" name="host" value="${config.host}" class="form-control" required>
                            </div>
                            <div class="form-group">
                                <label>端口</label>
                                <input type="number" name="port" value="${config.port}" class="form-control" required min="1" max="65535">
                            </div>
                            <div class="form-group">
                                <label>用户名</label>
                                <input type="text" name="user" value="${config.user || ''}" class="form-control" placeholder="例如: smtp@${config.domain}">
                            </div>
                            <div class="form-group">
                                <label>密码</label>
                                <input type="password" name="password" value="" class="form-control" placeholder="留空表示不修改密码">
                            </div>
                            <div class="form-group">
                                <label class="checkbox-label">
                                    <input type="checkbox" name="tls" ${config.tls ? 'checked' : ''}>
                                    启用TLS加密
                                </label>
                            </div>
                        </form>
                    </div>
                    <div class="modal-footer">
                        <button class="btn btn-secondary" onclick="this.closest('.modal-overlay').remove()">取消</button>
                        <button class="btn btn-primary" onclick="saveConfigEdit('${domain}')">保存</button>
                    </div>
                </div>
            ` + "`" + `;

            document.body.appendChild(modal);

            // 点击背景关闭模态框
            modal.addEventListener('click', function(e) {
                if (e.target === modal) {
                    modal.remove();
                }
            });
        }

        // 保存配置编辑
        async function saveConfigEdit(domain) {
            const form = document.getElementById('editConfigForm');
            const formData = new FormData(form);

            const configData = {
                domain: formData.get('domain'),
                host: formData.get('host'),
                port: parseInt(formData.get('port')),
                user: formData.get('user'),
                tls: formData.has('tls')
            };

            // 只有在输入了密码时才包含密码字段
            const password = formData.get('password');
            if (password) {
                configData.password = password;
            }

            try {
                const response = await fetch('/api/smtp-configs', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(configData)
                });

                const result = await response.json();
                if (result.success) {
                    showNotification('SMTP配置更新成功！', 'success');
                    document.querySelector('.modal-overlay').remove();
                    loadConfigs(); // 重新加载配置列表
                } else {
                    showNotification('更新失败: ' + result.message, 'error');
                }
            } catch (error) {
                showNotification('更新失败，请重试', 'error');
            }
        }

        // 验证SMTP DNS配置
        async function verifySMTPDNS(domain) {
            try {
                showNotification('正在验证SMTP DNS配置...', 'info');

                const response = await fetch(` + "`" + `/api/smtp-configs/${domain}/verify` + "`" + `);
                const result = await response.json();

                if (result.success) {
                    showDNSVerificationResult(result.data);
                } else {
                    showNotification('DNS验证失败: ' + result.message, 'error');
                }
            } catch (error) {
                showNotification('DNS验证失败，请重试', 'error');
            }
        }

        // 显示DNS验证结果
        function showDNSVerificationResult(data) {
            const modal = document.createElement('div');
            modal.className = 'modal-overlay';
            modal.innerHTML = ` + "`" + `
                <div class="modal-content dns-verification-modal" style="max-width: 900px; max-height: 90vh; overflow-y: auto;">
                    <div class="modal-header">
                        <h3>🔍 SMTP DNS验证结果</h3>
                        <button class="modal-close" onclick="this.closest('.modal-overlay').remove()">×</button>
                    </div>
                    <div class="modal-body">
                        <div class="verification-summary">
                            <div class="summary-item">
                                <strong>域名:</strong> ${data.domain}
                            </div>
                            <div class="summary-item">
                                <strong>SMTP服务器:</strong> ${data.smtp_host}:${data.smtp_port}
                            </div>
                            <div class="summary-item">
                                <strong>总体状态:</strong>
                                <span class="status-badge ${data.overall ? 'status-success' : 'status-warning'}">
                                    ${data.overall ? '✅ 通过' : '⚠️ 需要注意'}
                                </span>
                            </div>
                        </div>

                        <div class="verification-checks">
                            <h4>详细检查结果:</h4>
                            ${data.checks.map(check => ` + "`" + `
                                <div class="check-item ${check.success ? 'check-success' : 'check-failed'}">
                                    <div class="check-header">
                                        <span class="check-icon">${check.success ? '✅' : '❌'}</span>
                                        <span class="check-name">${check.name}</span>
                                    </div>
                                    <div class="check-description">${check.description}</div>
                                    <div class="check-message">${check.message}</div>
                                    ${check.ips ? ` + "`" + `<div class="check-details">IP地址: ${check.ips.join(', ')}</div>` + "`" + ` : ''}
                                    ${check.mx_records ? ` + "`" + `<div class="check-details">MX记录: ${check.mx_records.join(', ')}</div>` + "`" + ` : ''}
                                </div>
                            ` + "`" + `).join('')}
                        </div>

                        ${data.message ? ` + "`" + `<div class="verification-message">${data.message}</div>` + "`" + ` : ''}

                        <!-- DNS配置说明 -->
                        <div class="dns-config-section" style="margin-top: 30px; padding: 20px; background: #f8f9fa; border-radius: 8px; border-left: 4px solid #007bff;">
                            <h4 style="color: #333; margin-bottom: 15px;">📋 DNS配置说明</h4>
                            <div style="background: #fff3cd; border: 1px solid #ffeaa7; border-radius: 8px; padding: 15px; margin-bottom: 20px;">
                                <strong>⚠️ 重要提示：</strong> 请在您的域名DNS管理面板中添加以下记录，然后等待DNS传播完成（通常需要几分钟到几小时）。
                            </div>

                            <div class="dns-records">
                                <div class="dns-record-item" style="margin-bottom: 20px; padding: 15px; background: white; border-radius: 8px; border: 1px solid #e9ecef;">
                                    <div style="display: flex; align-items: center; margin-bottom: 10px;">
                                        <span style="background: #28a745; color: white; padding: 4px 8px; border-radius: 4px; font-size: 12px; margin-right: 10px;">MX记录</span>
                                        <strong>邮件交换记录</strong>
                                    </div>
                                    <div style="font-family: monospace; background: #f8f9fa; padding: 10px; border-radius: 4px; margin-bottom: 10px;">
                                        类型: MX<br>
                                        名称: @<br>
                                        值: ${data.domain}<br>
                                        优先级: 10<br>
                                        TTL: 3600
                                    </div>
                                    <div style="color: #666; font-size: 14px;">邮件交换记录，用于接收邮件</div>
                                </div>

                                <div class="dns-record-item" style="margin-bottom: 20px; padding: 15px; background: white; border-radius: 8px; border: 1px solid #e9ecef;">
                                    <div style="display: flex; align-items: center; margin-bottom: 10px;">
                                        <span style="background: #17a2b8; color: white; padding: 4px 8px; border-radius: 4px; font-size: 12px; margin-right: 10px;">A记录</span>
                                        <strong>域名解析记录</strong>
                                    </div>
                                    <div style="font-family: monospace; background: #f8f9fa; padding: 10px; border-radius: 4px; margin-bottom: 10px;">
                                        类型: A<br>
                                        名称: @<br>
                                        值: 111.119.198.162<br>
                                        TTL: 3600
                                    </div>
                                    <div style="color: #666; font-size: 14px;">A记录，将域名指向服务器IP</div>
                                </div>

                                <div class="dns-record-item" style="margin-bottom: 20px; padding: 15px; background: white; border-radius: 8px; border: 1px solid #e9ecef;">
                                    <div style="display: flex; align-items: center; margin-bottom: 10px;">
                                        <span style="background: #ffc107; color: black; padding: 4px 8px; border-radius: 4px; font-size: 12px; margin-right: 10px;">TXT记录</span>
                                        <strong>SPF记录（可选）</strong>
                                    </div>
                                    <div style="font-family: monospace; background: #f8f9fa; padding: 10px; border-radius: 4px; margin-bottom: 10px;">
                                        类型: TXT<br>
                                        名称: @<br>
                                        值: v=spf1 ip4:111.119.198.162 ~all<br>
                                        TTL: 3600
                                    </div>
                                    <div style="color: #666; font-size: 14px;">SPF记录，防止邮件被标记为垃圾邮件</div>
                                </div>
                            </div>

                            <div class="config-steps" style="margin-top: 20px;">
                                <h5 style="color: #333; margin-bottom: 10px;">📝 配置步骤</h5>
                                <ol style="color: #666; line-height: 1.6;">
                                    <li>登录您的域名注册商或DNS服务商管理面板</li>
                                    <li>找到DNS记录管理页面（通常叫"域名解析"、"DNS管理"或"解析设置"）</li>
                                    <li>按照上述配置添加MX记录、A记录和TXT记录</li>
                                    <li>保存所有记录并等待DNS传播生效（通常需要几分钟到几小时）</li>
                                    <li>返回本页面重新点击"DNS验证"按钮检查配置是否正确</li>
                                </ol>
                            </div>

                            <div class="common-providers" style="margin-top: 20px;">
                                <h5 style="color: #333; margin-bottom: 10px;">🔗 常用DNS服务商</h5>
                                <div style="display: flex; flex-wrap: wrap; gap: 10px;">
                                    <a href="https://dns.console.aliyun.com/" target="_blank" style="padding: 8px 12px; background: #007bff; color: white; text-decoration: none; border-radius: 4px; font-size: 12px;">阿里云</a>
                                    <a href="https://console.cloud.tencent.com/cns" target="_blank" style="padding: 8px 12px; background: #007bff; color: white; text-decoration: none; border-radius: 4px; font-size: 12px;">腾讯云</a>
                                    <a href="https://console.bce.baidu.com/bcd/" target="_blank" style="padding: 8px 12px; background: #007bff; color: white; text-decoration: none; border-radius: 4px; font-size: 12px;">百度云</a>
                                    <a href="https://console.huaweicloud.com/dns/" target="_blank" style="padding: 8px 12px; background: #007bff; color: white; text-decoration: none; border-radius: 4px; font-size: 12px;">华为云</a>
                                    <a href="https://dash.cloudflare.com/" target="_blank" style="padding: 8px 12px; background: #007bff; color: white; text-decoration: none; border-radius: 4px; font-size: 12px;">Cloudflare</a>
                                </div>
                            </div>
                        </div>
                    </div>
                    <div class="modal-footer">
                        <button class="btn btn-primary" onclick="this.closest('.modal-overlay').remove()">关闭</button>
                    </div>
                </div>
            ` + "`" + `;

            document.body.appendChild(modal);

            // 点击背景关闭模态框
            modal.addEventListener('click', function(e) {
                if (e.target === modal) {
                    modal.remove();
                }
            });
        }

        // 删除配置
        async function deleteConfig(domain) {
            if (!confirm(` + "`" + `确定要删除域名 "${domain}" 的SMTP配置吗？` + "`" + `)) {
                return;
            }

            try {
                const response = await fetch(` + "`" + `/api/smtp-configs/${encodeURIComponent(domain)}` + "`" + `, {
                    method: 'DELETE'
                });

                const result = await response.json();
                if (result.success) {
                    showNotification('SMTP配置删除成功！', 'success');
                    loadConfigs();
                } else {
                    showNotification('删除失败: ' + result.message, 'error');
                }
            } catch (error) {
                showNotification('删除失败，请重试', 'error');
            }
        }

        // 显示通知
        function showNotification(message, type) {
            const notification = document.createElement('div');
            notification.className = ` + "`" + `notification ${type}` + "`" + `;
            notification.textContent = message;
            notification.style.cssText = ` + "`" + `
                position: fixed;
                top: 20px;
                right: 20px;
                padding: 15px 20px;
                border-radius: 8px;
                color: white;
                font-weight: 500;
                z-index: 10000;
                transform: translateX(100%);
                transition: transform 0.3s ease;
            ` + "`" + `;

            if (type === 'success') {
                notification.style.background = '#4caf50';
            } else if (type === 'error') {
                notification.style.background = '#f44336';
            } else {
                notification.style.background = '#2196f3';
            }

            document.body.appendChild(notification);

            setTimeout(() => {
                notification.style.transform = 'translateX(0)';
            }, 100);

            setTimeout(() => {
                notification.style.transform = 'translateX(100%)';
                setTimeout(() => document.body.removeChild(notification), 300);
            }, 3000);
        }

        // 退出登录
        async function logout() {
            if (confirm('确定要退出登录吗？')) {
                try {
                    await fetch('/api/logout', { method: 'POST' });
                    window.location.href = '/login';
                } catch (error) {
                    window.location.href = '/login';
                }
            }
        }
    </script>`

	// 将脚本插入到模板中
	tmpl = strings.Replace(tmpl, "</body>", smtpScript+"\n    </body>", 1)

	return tmpl

}
