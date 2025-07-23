package server

// getBasePageTemplate 获取基础页面模板
func (s *Server) getBasePageTemplate(userEmail string, isAdmin bool, activeNav string) string {
	// 这里返回收件箱页面的基础模板，可以根据activeNav参数调整
	return s.generateInboxTemplate(userEmail, isAdmin, activeNav)
}

// generateInboxTemplate 生成收件箱模板
func (s *Server) generateInboxTemplate(userEmail string, isAdmin bool, activeNav string) string {
	// 现代化的收件箱模板
	tmpl := `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>NBEmail - 收件箱</title>
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
            border-radius: 0;
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
            min-height: 0; /* 允许flex子元素收缩 */
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
        }
        .toolbar-right {
            display: flex;
            gap: 12px;
            flex-wrap: wrap;
        }
        .btn {
            padding: 8px 16px;
            border: 2px solid #e9ecef;
            background: white;
            border-radius: 8px;
            cursor: pointer;
            font-weight: 500;
            transition: all 0.3s;
            display: flex;
            align-items: center;
            gap: 6px;
            font-size: 13px;
            color: #495057;
            white-space: nowrap;
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

        /* 邮件列表 */
        .email-list {
            background: white;
            border-radius: 15px;
            overflow-y: auto; /* 允许垂直滚动 */
            box-shadow: 0 5px 20px rgba(0,0,0,0.08);
            flex: 1;
            min-height: 0; /* 允许flex子元素收缩 */
        }

        /* 自定义滚动条样式 */
        .email-list::-webkit-scrollbar {
            width: 8px;
        }
        .email-list::-webkit-scrollbar-track {
            background: #f1f1f1;
            border-radius: 10px;
        }
        .email-list::-webkit-scrollbar-thumb {
            background: #c1c1c1;
            border-radius: 10px;
        }
        .email-list::-webkit-scrollbar-thumb:hover {
            background: #a8a8a8;
        }
        .email-item {
            padding: 20px 25px;
            border-bottom: 1px solid #f1f3f4;
            cursor: pointer;
            display: flex;
            align-items: center;
            gap: 20px;
            transition: all 0.3s;
        }
        .email-item:hover {
            background: #f8f9fa;
            transform: translateX(5px);
        }
        .email-item:last-child {
            border-bottom: none;
        }
        .email-item.unread {
            background: linear-gradient(90deg, rgba(102,126,234,0.05) 0%, rgba(255,255,255,1) 100%);
            border-left: 4px solid #667eea;
        }
        .email-item.selected {
            background: linear-gradient(90deg, rgba(33,150,243,0.1) 0%, rgba(227,242,253,0.8) 100%);
            border-left: 4px solid #2196f3;
            box-shadow: 0 2px 8px rgba(33,150,243,0.15);
        }
        .email-item.selected.unread {
            background: linear-gradient(90deg, rgba(25,118,210,0.15) 0%, rgba(227,242,253,0.9) 100%);
            border-left: 4px solid #1976d2;
        }
        .email-item.selected {
            background: linear-gradient(90deg, rgba(33,150,243,0.1) 0%, rgba(227,242,253,0.8) 100%);
            border-left: 4px solid #2196f3;
            box-shadow: 0 2px 8px rgba(33,150,243,0.15);
        }
        .email-item.selected.unread {
            background: linear-gradient(90deg, rgba(25,118,210,0.15) 0%, rgba(227,242,253,0.9) 100%);
            border-left: 4px solid #1976d2;
        }
        .email-checkbox {
            transform: scale(1.2);
            accent-color: #667eea;
        }
        .email-avatar {
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
        .email-content {
            flex: 1;
            display: flex;
            flex-direction: column;
            gap: 5px;
        }
        .email-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .email-from {
            font-weight: 600;
            color: #333;
            font-size: 15px;
        }
        .email-date {
            color: #8e9aaf;
            font-size: 13px;
            font-weight: 500;
        }
        .email-subject {
            font-weight: 500;
            color: #5a6c7d;
            font-size: 14px;
            margin-bottom: 3px;
        }
        .email-preview {
            color: #8e9aaf;
            font-size: 13px;
            line-height: 1.4;
        }
        .email-status {
            display: flex;
            gap: 8px;
            align-items: center;
        }
        .status-badge {
            padding: 3px 8px;
            border-radius: 12px;
            font-size: 11px;
            font-weight: 600;
            text-transform: uppercase;
        }
        .status-unread {
            background: #e3f2fd;
            color: #1976d2;
        }
        .status-important {
            background: #fff3e0;
            color: #f57c00;
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

        /* 搜索框样式 */
        .search-container {
            margin-bottom: 20px;
        }
        .search-box {
            display: flex;
            background: white;
            border-radius: 15px;
            padding: 15px 20px;
            box-shadow: 0 5px 20px rgba(0,0,0,0.08);
            gap: 10px;
            align-items: center;
        }
        .search-box input {
            flex: 1;
            border: none;
            outline: none;
            font-size: 16px;
            color: #333;
        }
        .search-box input::placeholder {
            color: #8e9aaf;
        }
        .search-btn, .clear-btn {
            background: none;
            border: none;
            font-size: 18px;
            cursor: pointer;
            padding: 5px;
            border-radius: 8px;
            transition: background-color 0.2s;
        }
        .search-btn:hover {
            background: #f1f3f4;
        }
        .clear-btn:hover {
            background: #fee;
        }

        /* 分页样式 */
        .pagination-container {
            margin-top: 20px;
            display: flex;
            justify-content: center;
        }
        .pagination {
            display: flex;
            gap: 8px;
            align-items: center;
            background: white;
            padding: 15px 20px;
            border-radius: 15px;
            box-shadow: 0 5px 20px rgba(0,0,0,0.08);
        }
        .page-btn {
            padding: 8px 12px;
            border: none;
            background: #f8f9fa;
            color: #6c757d;
            border-radius: 8px;
            cursor: pointer;
            font-size: 14px;
            transition: all 0.2s;
            min-width: 40px;
        }
        .page-btn:hover {
            background: #e9ecef;
            color: #495057;
        }
        .page-btn.active {
            background: #007bff;
            color: white;
        }
        .page-ellipsis {
            padding: 8px 4px;
            color: #6c757d;
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
            .container {
                flex-direction: column;
                padding: 0 10px;
            }
            .sidebar {
                width: 100%;
                order: 2;
                margin-top: 20px;
                background: white;
                border-radius: 12px;
                box-shadow: 0 2px 8px rgba(0,0,0,0.1);
                padding: 20px;
            }
            .main {
                order: 1;
                min-height: 60vh;
                width: 100%;
            }
            .toolbar-right {
                gap: 8px;
            }
            .btn {
                padding: 6px 12px;
                font-size: 12px;
                gap: 4px;
            }
            .email-item {
                padding: 15px;
                gap: 15px;
            }
            .email-avatar {
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
                <a href="/inbox" class="nav-item` + func() string {
		if activeNav == "inbox" {
			return " active"
		}
		return ""
	}() + `">
                    <span class="icon">📥</span>收件箱
                </a>
                <a href="/sent" class="nav-item` + func() string {
		if activeNav == "sent" {
			return " active"
		}
		return ""
	}() + `">
                    <span class="icon">📤</span>已发送
                </a>
                <a href="/compose" class="nav-item` + func() string {
		if activeNav == "compose" {
			return " active"
		}
		return ""
	}() + `">
                    <span class="icon">✏️</span>写邮件
                </a>
            </div>`

	if isAdmin {
		tmpl += `
            <div class="nav-divider"></div>
            <div class="nav-section">
                <div class="nav-section-title">管理</div>
                <a href="/users" class="nav-item` + func() string {
			if activeNav == "users" {
				return " active"
			}
			return ""
		}() + `">
                    <span class="icon">👥</span>用户管理
                </a>
                <a href="/domains" class="nav-item` + func() string {
			if activeNav == "domains" {
				return " active"
			}
			return ""
		}() + `">
                    <span class="icon">🌐</span>域名管理
                </a>
                <a href="/smtp-configs" class="nav-item` + func() string {
			if activeNav == "smtp-configs" {
				return " active"
			}
			return ""
		}() + `">
                    <span class="icon">📮</span>SMTP配置
                </a>
            </div>`
	}

	tmpl += `
        </div>
        <div class="main">
            <div class="toolbar">
                <div class="toolbar-left">
                    <h2 style="margin: 0; color: #333; font-weight: 600;">` + func() string {
		switch activeNav {
		case "sent":
			return "📤 已发送"
		case "compose":
			return "✏️ 写邮件"
		default:
			return "📥 收件箱"
		}
	}() + `</h2>
                    <div class="current-mailbox" id="currentMailbox" style="margin-left: 20px; padding: 8px 15px; background: #e3f2fd; border-radius: 20px; font-size: 14px; color: #1976d2; display: flex; align-items: center; gap: 8px;">
                        <span>📧</span>
                        <span id="currentEmail">加载中...</span>
                        <button onclick="showMailboxManager()" style="background: none; border: none; color: #1976d2; cursor: pointer; padding: 2px 6px; border-radius: 4px; font-size: 12px;" title="管理邮箱">⚙️</button>
                    </div>
                </div>
                <div class="toolbar-right">
                    <button class="btn" onclick="showMailboxManager()">
                        <span>📮</span>邮箱管理
                    </button>
                    <button class="btn" onclick="refreshEmails()">
                        <span>🔄</span>刷新
                    </button>
                    <button class="btn" onclick="selectAllCurrentPage()">
                        <span>☑️</span>全选当前页
                    </button>
                    <button class="btn" onclick="deselectAll()">
                        <span>⬜</span>取消全选
                    </button>
                    <button class="btn btn-danger" onclick="deleteSelected()">
                        <span>🗑️</span>删除选中
                    </button>
                    <button class="btn btn-danger" onclick="deleteAllEmails()">
                        <span>🗑️💥</span>删除全部
                    </button>
                    <button class="btn" onclick="markAsRead()">
                        <span>✅</span>标记已读
                    </button>
                </div>
            </div>

            <!-- 搜索框 -->
            <div class="search-container">
                <div class="search-box">
                    <input type="text" id="searchInput" placeholder="搜索邮件（主题、发件人、收件人、内容）" onkeypress="if(event.key==='Enter') performSearch()">
                    <button class="search-btn" onclick="performSearch()">🔍</button>
                    <button class="clear-btn" onclick="clearSearch()">✖️</button>
                </div>
            </div>

            <div class="email-list" id="emailList">
                <div class="empty-state">
                    <div class="icon">📬</div>
                    <h3>正在加载邮件...</h3>
                    <p>请稍候，正在获取您的邮件</p>
                </div>
            </div>

            <!-- 分页容器 -->
            <div id="pagination" class="pagination-container"></div>
        </div>
    </div>

    <!-- 邮箱管理模态框 -->
    <div id="mailboxModal" class="modal" style="display: none; position: fixed; z-index: 1000; left: 0; top: 0; width: 100%; height: 100%; background-color: rgba(0,0,0,0.5); backdrop-filter: blur(5px);">
        <div class="modal-content" style="background-color: white; margin: 5% auto; padding: 30px; border-radius: 15px; width: 90%; max-width: 800px; box-shadow: 0 20px 40px rgba(0,0,0,0.2); max-height: 80vh; overflow-y: auto;">
            <div class="modal-header" style="margin-bottom: 25px; padding-bottom: 15px; border-bottom: 2px solid #f1f3f4; display: flex; justify-content: space-between; align-items: center;">
                <h3 style="color: #333; font-weight: 600; margin: 0;">📮 邮箱管理</h3>
                <span class="close" onclick="closeMailboxModal()" style="color: #aaa; font-size: 28px; font-weight: bold; cursor: pointer; line-height: 1;">&times;</span>
            </div>

            <div class="mailbox-actions" style="margin-bottom: 20px; display: flex; gap: 10px; flex-wrap: wrap;">
                <button class="btn btn-primary" onclick="showGenerateModal()">
                    <span>➕</span>批量生成邮箱
                </button>
                <button class="btn" onclick="showCredentials()">
                    <span>🔑</span>查看登录凭据
                </button>
                <button class="btn" onclick="refreshMailboxes()">
                    <span>🔄</span>刷新列表
                </button>
            </div>

            <div class="current-mailbox-info" style="background: #e8f5e8; padding: 15px; border-radius: 10px; margin-bottom: 20px;">
                <h4 style="margin: 0 0 10px 0; color: #2e7d32;">当前使用的邮箱</h4>
                <div id="currentMailboxInfo" style="font-weight: 600; color: #1b5e20;">加载中...</div>
            </div>

            <div class="mailbox-list" id="mailboxList">
                <div style="text-align: center; padding: 40px; color: #666;">正在加载邮箱列表...</div>
            </div>
        </div>
    </div>

    <!-- 批量生成邮箱模态框 -->
    <div id="generateModal" class="modal" style="display: none; position: fixed; z-index: 1001; left: 0; top: 0; width: 100%; height: 100%; background-color: rgba(0,0,0,0.5); backdrop-filter: blur(5px);">
        <div class="modal-content" style="background-color: white; margin: 15% auto; padding: 30px; border-radius: 15px; width: 90%; max-width: 500px; box-shadow: 0 20px 40px rgba(0,0,0,0.2);">
            <div class="modal-header" style="margin-bottom: 25px; padding-bottom: 15px; border-bottom: 2px solid #f1f3f4;">
                <span class="close" onclick="closeGenerateModal()" style="color: #aaa; float: right; font-size: 28px; font-weight: bold; cursor: pointer; line-height: 1;">&times;</span>
                <h3 style="color: #333; font-weight: 600;">➕ 批量生成邮箱</h3>
            </div>
            <form id="generateForm">
                <div class="form-group" style="margin-bottom: 20px;">
                    <label style="display: block; margin-bottom: 8px; font-weight: 600; color: #333;">选择域名 *</label>
                    <select id="domainSelect" required style="width: 100%; padding: 12px 15px; border: 2px solid #e9ecef; border-radius: 8px; font-size: 14px;">
                        <option value="">请选择域名</option>
                    </select>
                </div>
                <div class="form-group" style="margin-bottom: 20px;">
                    <label style="display: block; margin-bottom: 8px; font-weight: 600; color: #333;">生成数量 *</label>
                    <input type="number" id="generateCount" required min="1" max="100" value="10" style="width: 100%; padding: 12px 15px; border: 2px solid #e9ecef; border-radius: 8px; font-size: 14px;">
                </div>
                <div class="form-group" style="margin-bottom: 20px;">
                    <label style="display: block; margin-bottom: 8px; font-weight: 600; color: #333;">邮箱前缀（可选）</label>
                    <input type="text" id="emailPrefix" placeholder="例如：user、test" style="width: 100%; padding: 12px 15px; border: 2px solid #e9ecef; border-radius: 8px; font-size: 14px;">
                </div>
                <div style="display: flex; gap: 10px; justify-content: flex-end; margin-top: 25px;">
                    <button type="button" class="btn" onclick="closeGenerateModal()">取消</button>
                    <button type="submit" class="btn btn-primary">生成邮箱</button>
                </div>
            </form>
        </div>
    </div>

    <script>
        let emails = [];
        let selectedEmails = [];
        let mailboxes = [];
        let domains = [];
        let currentPage = 1;
        let totalPages = 1;
        let currentSearch = '';
        let currentFolder = '` + activeNav + `';

        async function loadEmails(page = 1, search = '') {
            try {
                currentPage = page;
                currentSearch = search;

                let url = '/api/emails?folder=' + currentFolder + '&page=' + page + '&limit=20';
                if (search) {
                    url += '&search=' + encodeURIComponent(search);
                }

                const response = await fetch(url);
                const result = await response.json();
                if (result.success) {
                    emails = result.data.emails || [];
                    totalPages = Math.ceil(result.data.total / result.data.limit);
                    renderEmails();
                    renderPagination();
                } else {
                    document.getElementById('emailList').innerHTML = '<div class="empty-state">加载邮件失败</div>';
                }
            } catch (error) {
                document.getElementById('emailList').innerHTML = '<div class="empty-state">加载邮件失败</div>';
            }
        }

        function renderEmails() {
            const emailList = document.getElementById('emailList');
            if (emails.length === 0) {
                emailList.innerHTML = ` + "`" + `
                    <div class="empty-state">
                        <div class="icon">📭</div>
                        <h3>暂无邮件</h3>
                        <p>您的收件箱是空的，暂时没有新邮件</p>
                    </div>
                ` + "`" + `;
                return;
            }

            emailList.innerHTML = emails.map(email => {
                const date = new Date(email.created_at);
                const dateStr = date.toLocaleDateString('zh-CN', {
                    month: 'short',
                    day: 'numeric',
                    hour: '2-digit',
                    minute: '2-digit'
                });
                const unreadClass = email.is_read ? '' : 'unread';
                const fromInitial = email.from.charAt(0).toUpperCase();
                const preview = email.body ? email.body.substring(0, 100) + '...' : '';

                const selectedClass = selectedEmails.includes(email.id) ? 'selected' : '';
                return ` + "`" + `
                    <div class="email-item ${unreadClass} ${selectedClass}" data-email-id="${email.id}" onclick="openEmail(${email.id})">
                        <input type="checkbox" class="email-checkbox" onclick="event.stopPropagation(); toggleSelect(${email.id})" ${selectedEmails.includes(email.id) ? 'checked' : ''}>
                        <div class="email-avatar">${fromInitial}</div>
                        <div class="email-content">
                            <div class="email-header">
                                <div class="email-from">${email.from}</div>
                                <div class="email-status">
                                    ${!email.is_read ? '<span class="status-badge status-unread">未读</span>' : ''}
                                    <div class="email-date">${dateStr}</div>
                                </div>
                            </div>
                            <div class="email-subject">${email.subject}</div>
                            <div class="email-preview">${preview}</div>
                        </div>
                    </div>
                ` + "`" + `;
            }).join('');
        }

        function toggleSelect(emailId) {
            const index = selectedEmails.indexOf(emailId);
            if (index > -1) {
                selectedEmails.splice(index, 1);
            } else {
                selectedEmails.push(emailId);
            }
        }

        function refreshEmails() {
            loadEmails(currentPage, currentSearch);
        }

        function renderPagination() {
            const paginationContainer = document.getElementById('pagination');
            if (!paginationContainer) return;

            if (totalPages <= 1) {
                paginationContainer.innerHTML = '';
                return;
            }

            let paginationHTML = '<div class="pagination">';

            // 上一页
            if (currentPage > 1) {
                paginationHTML += '<button onclick="loadEmails(' + (currentPage - 1) + ', currentSearch)" class="page-btn">上一页</button>';
            }

            // 页码
            const startPage = Math.max(1, currentPage - 2);
            const endPage = Math.min(totalPages, currentPage + 2);

            if (startPage > 1) {
                paginationHTML += ` + "`" + `<button onclick="loadEmails(1, currentSearch)" class="page-btn">1</button>` + "`" + `;
                if (startPage > 2) {
                    paginationHTML += '<span class="page-ellipsis">...</span>';
                }
            }

            for (let i = startPage; i <= endPage; i++) {
                const activeClass = i === currentPage ? ' active' : '';
                paginationHTML += '<button onclick="loadEmails(' + i + ', currentSearch)" class="page-btn' + activeClass + '">' + i + '</button>';
            }

            if (endPage < totalPages) {
                if (endPage < totalPages - 1) {
                    paginationHTML += '<span class="page-ellipsis">...</span>';
                }
                paginationHTML += '<button onclick="loadEmails(' + totalPages + ', currentSearch)" class="page-btn">' + totalPages + '</button>';
            }

            // 下一页
            if (currentPage < totalPages) {
                paginationHTML += '<button onclick="loadEmails(' + (currentPage + 1) + ', currentSearch)" class="page-btn">下一页</button>';
            }

            paginationHTML += '</div>';
            paginationContainer.innerHTML = paginationHTML;
        }

        function performSearch() {
            const searchInput = document.getElementById('searchInput');
            const searchTerm = searchInput ? searchInput.value.trim() : '';
            loadEmails(1, searchTerm);
        }

        function clearSearch() {
            const searchInput = document.getElementById('searchInput');
            if (searchInput) {
                searchInput.value = '';
            }
            loadEmails(1, '');
        }

        function selectAllCurrentPage() {
            selectedEmails = [];
            emails.forEach(email => {
                selectedEmails.push(email.id);
            });
            updateEmailSelection();
            showNotification('已选中当前页 ' + selectedEmails.length + ' 封邮件', 'success');
        }

        function deselectAll() {
            selectedEmails = [];
            updateEmailSelection();
            showNotification('已取消全部选择', 'info');
        }

        function updateEmailSelection() {
            // 更新邮件项的选中状态显示
            emails.forEach(email => {
                const emailElement = document.querySelector('[data-email-id="' + email.id + '"]');
                if (emailElement) {
                    const checkbox = emailElement.querySelector('input[type="checkbox"]');
                    if (checkbox) {
                        checkbox.checked = selectedEmails.includes(email.id);
                    }
                    // 更新邮件项的视觉状态
                    if (selectedEmails.includes(email.id)) {
                        emailElement.classList.add('selected');
                    } else {
                        emailElement.classList.remove('selected');
                    }
                }
            });
        }

        async function deleteSelected() {
            if (selectedEmails.length === 0) {
                alert('请选择要删除的邮件');
                return;
            }

            if (!confirm('确定要删除选中的邮件吗？')) {
                return;
            }

            for (const emailId of selectedEmails) {
                try {
                    await fetch('/api/emails/' + emailId, { method: 'DELETE' });
                } catch (error) {
                    console.error('删除邮件失败:', error);
                }
            }

            selectedEmails = [];
            loadEmails(currentPage, currentSearch);
        }

        async function markAsRead() {
            if (selectedEmails.length === 0) {
                alert('请选择要标记的邮件');
                return;
            }

            for (const emailId of selectedEmails) {
                try {
                    await fetch('/api/emails/' + emailId + '/read', { method: 'PUT' });
                } catch (error) {
                    console.error('标记邮件失败:', error);
                }
            }

            selectedEmails = [];
            loadEmails(currentPage, currentSearch);
        }

        async function deleteAllEmails() {
            const folderName = currentFolder === 'inbox' ? '收件箱' : '已发送';
            const confirmMessage = '⚠️ 危险操作警告！\\n\\n您即将删除' + folderName + '中的所有邮件！\\n这个操作不可撤销！\\n\\n请输入 "DELETE ALL" 来确认删除：';

            const userInput = prompt(confirmMessage);
            if (userInput !== 'DELETE ALL') {
                if (userInput !== null) {
                    alert('输入不正确，操作已取消');
                }
                return;
            }

            const secondConfirm = confirm('最后确认：您真的要删除' + folderName + '中的所有邮件吗？\\n\\n点击"确定"将永久删除所有邮件\\n点击"取消"将中止操作');
            if (!secondConfirm) {
                return;
            }

            try {
                showNotification('正在删除所有邮件，请稍候...', 'info');

                // 获取所有邮件ID
                const response = await fetch('/api/emails?folder=' + currentFolder + '&limit=1000');
                const result = await response.json();

                if (!result.success) {
                    throw new Error('获取邮件列表失败');
                }

                const allEmails = result.data.emails || [];
                let deletedCount = 0;
                let failedCount = 0;

                // 批量删除邮件
                for (const email of allEmails) {
                    try {
                        const deleteResponse = await fetch('/api/emails/' + email.id, { method: 'DELETE' });
                        if (deleteResponse.ok) {
                            deletedCount++;
                        } else {
                            failedCount++;
                        }
                    } catch (error) {
                        console.error('删除邮件失败:', error);
                        failedCount++;
                    }
                }

                selectedEmails = [];
                loadEmails(1, currentSearch);

                if (failedCount === 0) {
                    showNotification('✅ 成功删除 ' + deletedCount + ' 封邮件', 'success');
                } else {
                    showNotification('⚠️ 删除完成：成功 ' + deletedCount + ' 封，失败 ' + failedCount + ' 封', 'warning');
                }

            } catch (error) {
                console.error('删除全部邮件失败:', error);
                showNotification('❌ 删除失败，请重试', 'error');
            }
        }

        function openEmail(emailId) {
            window.open('/email/' + emailId, '_blank');
        }

        async function logout() {
            try {
                await fetch('/api/logout', { method: 'POST' });
                window.location.href = '/login';
            } catch (error) {
                window.location.href = '/login';
            }
        }

        // 邮箱管理相关函数
        async function loadCurrentMailbox() {
            try {
                const response = await fetch('/api/mailboxes');
                const result = await response.json();
                if (result.success && result.data.length > 0) {
                    const currentMailbox = result.data.find(m => m.is_current);
                    if (currentMailbox) {
                        document.getElementById('currentEmail').textContent = currentMailbox.email;
                    } else {
                        document.getElementById('currentEmail').textContent = '未设置';
                    }
                } else {
                    document.getElementById('currentEmail').textContent = '无邮箱';
                }
            } catch (error) {
                document.getElementById('currentEmail').textContent = '加载失败';
            }
        }

        async function showMailboxManager() {
            document.getElementById('mailboxModal').style.display = 'block';
            await loadMailboxes();
            await loadDomains();
        }

        function closeMailboxModal() {
            document.getElementById('mailboxModal').style.display = 'none';
        }

        async function loadMailboxes() {
            try {
                const response = await fetch('/api/mailboxes');
                const result = await response.json();
                if (result.success) {
                    mailboxes = result.data || [];
                    renderMailboxes();
                } else {
                    document.getElementById('mailboxList').innerHTML = '<div style="text-align: center; padding: 40px; color: #666;">加载邮箱失败</div>';
                }
            } catch (error) {
                document.getElementById('mailboxList').innerHTML = '<div style="text-align: center; padding: 40px; color: #666;">加载邮箱失败</div>';
            }
        }

        function renderMailboxes() {
            const mailboxList = document.getElementById('mailboxList');
            const currentMailboxInfo = document.getElementById('currentMailboxInfo');

            if (mailboxes.length === 0) {
                mailboxList.innerHTML = ` + "`" + `
                    <div style="text-align: center; padding: 40px; color: #666;">
                        <div style="font-size: 3rem; margin-bottom: 15px;">📭</div>
                        <h3>暂无邮箱</h3>
                        <p>请先生成一些邮箱</p>
                    </div>
                ` + "`" + `;
                currentMailboxInfo.textContent = '未设置当前邮箱';
                return;
            }

            const currentMailbox = mailboxes.find(m => m.is_current);
            if (currentMailbox) {
                currentMailboxInfo.innerHTML = ` + "`" + `
                    <div style="display: flex; align-items: center; gap: 10px;">
                        <span style="font-size: 1.2em;">📧</span>
                        <span>${currentMailbox.email}</span>
                        <span style="background: #4caf50; color: white; padding: 2px 8px; border-radius: 12px; font-size: 12px;">当前</span>
                    </div>
                ` + "`" + `;
            } else {
                currentMailboxInfo.textContent = '未设置当前邮箱';
            }

            mailboxList.innerHTML = ` + "`" + `
                <div style="display: grid; gap: 10px;">
                    ${mailboxes.map(mailbox => ` + "`" + `
                        <div style="display: flex; align-items: center; justify-content: space-between; padding: 15px; border: 2px solid ${mailbox.is_current ? '#4caf50' : '#e9ecef'}; border-radius: 10px; background: ${mailbox.is_current ? '#f1f8e9' : 'white'};">
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
                            <div style="display: flex; gap: 8px;">
                                ${!mailbox.is_current ? ` + "`" + `<button onclick="switchMailbox(${mailbox.id})" style="padding: 6px 12px; background: #007bff; color: white; border: none; border-radius: 6px; cursor: pointer; font-size: 12px;">切换</button>` + "`" + ` : ''}
                                ${!mailbox.is_current ? ` + "`" + `<button onclick="deleteMailbox(${mailbox.id})" style="padding: 6px 12px; background: #dc3545; color: white; border: none; border-radius: 6px; cursor: pointer; font-size: 12px;">删除</button>` + "`" + ` : ''}
                            </div>
                        </div>
                    ` + "`" + `).join('')}
                </div>
            ` + "`" + `;
        }

        async function switchMailbox(mailboxId) {
            try {
                const response = await fetch('/api/mailboxes/switch', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ mailbox_id: mailboxId })
                });

                const result = await response.json();
                if (result.success) {
                    showNotification('邮箱切换成功！', 'success');
                    await loadMailboxes();
                    await loadCurrentMailbox();
                } else {
                    showNotification('切换失败: ' + result.message, 'error');
                }
            } catch (error) {
                showNotification('切换失败，请重试', 'error');
            }
        }

        async function deleteMailbox(mailboxId) {
            if (!confirm('确定要删除这个邮箱吗？此操作不可撤销！')) return;

            try {
                const response = await fetch(` + "`" + `/api/mailboxes/${mailboxId}` + "`" + `, { method: 'DELETE' });
                const result = await response.json();
                if (result.success) {
                    showNotification('邮箱删除成功', 'success');
                    await loadMailboxes();
                } else {
                    showNotification('删除失败: ' + result.message, 'error');
                }
            } catch (error) {
                showNotification('删除失败，请重试', 'error');
            }
        }

        function refreshMailboxes() {
            loadMailboxes();
        }

        // 批量生成邮箱相关函数
        async function showGenerateModal() {
            document.getElementById('generateModal').style.display = 'block';
            // 每次打开模态框时重新加载域名列表
            await loadDomains();
        }

        function closeGenerateModal() {
            document.getElementById('generateModal').style.display = 'none';
            document.getElementById('generateForm').reset();
            // 清除生成结果显示
            const resultDiv = document.getElementById('generateResult');
            if (resultDiv) {
                resultDiv.remove();
            }
        }

        function showGeneratedAccounts(accounts) {
            // 移除之前的结果显示
            const existingResult = document.getElementById('generateResult');
            if (existingResult) {
                existingResult.remove();
            }

            // 创建结果显示区域
            const resultDiv = document.createElement('div');
            resultDiv.id = 'generateResult';
            resultDiv.style.cssText = ` + "`" + `
                margin-top: 20px;
                padding: 15px;
                background: #f8f9fa;
                border-radius: 8px;
                border: 1px solid #e9ecef;
                max-height: 300px;
                overflow-y: auto;
            ` + "`" + `;

            let html = ` + "`" + `
                <h4 style="margin: 0 0 15px 0; color: #28a745;">
                    <span>✅</span> 生成成功！以下是账号信息：
                </h4>
                <div style="font-size: 12px; color: #6c757d; margin-bottom: 10px;">
                    请保存这些信息，关闭窗口后将无法再次查看密码
                </div>
                <div style="display: grid; gap: 10px;">
            ` + "`" + `;

            accounts.forEach((account, index) => {
                html += ` + "`" + `
                    <div style="
                        padding: 10px;
                        background: white;
                        border-radius: 6px;
                        border: 1px solid #dee2e6;
                        display: grid;
                        grid-template-columns: 1fr auto;
                        gap: 10px;
                        align-items: center;
                    ">
                        <div>
                            <div style="font-weight: 500; color: #495057;">
                                ${account.email}
                            </div>
                            <div style="font-size: 12px; color: #6c757d; font-family: monospace;">
                                密码: ${account.password}
                            </div>
                        </div>
                        <button onclick="copyAccountInfo('${account.email}', '${account.password}')"
                                style="
                                    padding: 5px 10px;
                                    background: #007bff;
                                    color: white;
                                    border: none;
                                    border-radius: 4px;
                                    cursor: pointer;
                                    font-size: 12px;
                                ">
                            复制
                        </button>
                    </div>
                ` + "`" + `;
            });

            html += ` + "`" + `
                </div>
                <div style="margin-top: 15px; text-align: center;">
                    <button onclick="copyAllAccounts()" style="
                        padding: 8px 16px;
                        background: #28a745;
                        color: white;
                        border: none;
                        border-radius: 4px;
                        cursor: pointer;
                    ">
                        复制所有账号信息
                    </button>
                </div>
            ` + "`" + `;

            resultDiv.innerHTML = html;

            // 插入到表单后面
            const form = document.getElementById('generateForm');
            form.parentNode.insertBefore(resultDiv, form.nextSibling);

            // 存储账号信息供复制使用
            window.generatedAccounts = accounts;
        }

        function copyAccountInfo(email, password) {
            const text = ` + "`" + `邮箱: ${email}
密码: ${password}` + "`" + `;

            navigator.clipboard.writeText(text).then(() => {
                showNotification('账号信息已复制到剪贴板', 'success');
            }).catch(() => {
                // 降级方案
                const textArea = document.createElement('textarea');
                textArea.value = text;
                document.body.appendChild(textArea);
                textArea.select();
                document.execCommand('copy');
                document.body.removeChild(textArea);
                showNotification('账号信息已复制到剪贴板', 'success');
            });
        }

        function copyAllAccounts() {
            if (!window.generatedAccounts) return;

            let text = '批量生成的邮箱账号信息：\\n\\n';
            window.generatedAccounts.forEach((account, index) => {
                text += ` + "`" + `${index + 1}. 邮箱: ${account.email}
   密码: ${account.password}

` + "`" + `;
            });

            navigator.clipboard.writeText(text).then(() => {
                showNotification('所有账号信息已复制到剪贴板', 'success');
            }).catch(() => {
                // 降级方案
                const textArea = document.createElement('textarea');
                textArea.value = text;
                document.body.appendChild(textArea);
                textArea.select();
                document.execCommand('copy');
                document.body.removeChild(textArea);
                showNotification('所有账号信息已复制到剪贴板', 'success');
            });
        }

        async function loadDomains() {
            try {
                // 普通用户调用用户域名API
                const response = await fetch('/api/user/domains');
                const result = await response.json();
                if (result.success) {
                    domains = result.data || [];
                    const domainSelect = document.getElementById('domainSelect');
                    domainSelect.innerHTML = '<option value="">请选择域名</option>';
                    domains.forEach(domain => {
                        if (domain.is_active) {
                            domainSelect.innerHTML += ` + "`" + `<option value="${domain.id}">${domain.name}</option>` + "`" + `;
                        }
                    });

                    // 如果没有可用域名，显示提示
                    if (domains.length === 0) {
                        domainSelect.innerHTML = '<option value="">暂无可用域名，请联系管理员分配</option>';
                    }
                } else {
                    console.error('加载域名失败:', result.message);
                    const domainSelect = document.getElementById('domainSelect');
                    domainSelect.innerHTML = '<option value="">加载域名失败</option>';
                }
            } catch (error) {
                console.error('加载域名失败:', error);
                const domainSelect = document.getElementById('domainSelect');
                domainSelect.innerHTML = '<option value="">网络错误</option>';
            }
        }

        // 处理批量生成表单提交
        document.getElementById('generateForm').addEventListener('submit', async (e) => {
            e.preventDefault();

            const domainId = document.getElementById('domainSelect').value;
            const count = parseInt(document.getElementById('generateCount').value);
            const prefix = document.getElementById('emailPrefix').value;

            if (!domainId) {
                showNotification('请选择域名', 'error');
                return;
            }

            const submitBtn = e.target.querySelector('button[type="submit"]');
            const originalText = submitBtn.innerHTML;

            // 显示加载状态
            submitBtn.innerHTML = '<span>⏳</span>生成中...';
            submitBtn.disabled = true;

            try {
                const response = await fetch('/api/mailboxes/generate', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        domain_id: parseInt(domainId),
                        count: count,
                        prefix: prefix || undefined
                    })
                });

                const result = await response.json();
                if (result.success) {
                    // 显示成功状态
                    submitBtn.innerHTML = '<span>✅</span>生成成功';
                    submitBtn.style.background = '#28a745';

                    // 显示生成的账号信息
                    showGeneratedAccounts(result.data.accounts);

                    showNotification(` + "`" + `成功生成 ${result.data.count} 个邮箱账号！` + "`" + `, 'success');

                    // 5秒后关闭模态框并刷新列表
                    setTimeout(() => {
                        closeGenerateModal();
                        loadMailboxes();
                        submitBtn.innerHTML = originalText;
                        submitBtn.style.background = '';
                        submitBtn.disabled = false;
                    }, 5000);
                } else {
                    throw new Error(result.message || '生成失败');
                }
            } catch (error) {
                // 显示错误状态
                submitBtn.innerHTML = '<span>❌</span>生成失败';
                submitBtn.style.background = '#dc3545';

                showNotification('生成失败: ' + error.message, 'error');

                // 3秒后恢复按钮
                setTimeout(() => {
                    submitBtn.innerHTML = originalText;
                    submitBtn.style.background = '';
                    submitBtn.disabled = false;
                }, 3000);
            }
        });

        // 显示邮箱登录凭据
        async function showCredentials() {
            try {
                const response = await fetch('/api/mailboxes/credentials', {
                    method: 'GET',
                    headers: {
                        'Authorization': 'Bearer ' + localStorage.getItem('token')
                    }
                });

                const result = await response.json();
                if (result.success) {
                    const data = result.data;

                    // 构建邮箱列表显示
                    let mailboxList = '';
                    if (data.mailboxes && data.mailboxes.length > 0) {
                        mailboxList = '\n\n关联邮箱账号:\n' + data.mailboxes.map((mailbox, index) => {
                            if (typeof mailbox === 'string') {
                                return (index + 1) + '. ' + mailbox;
                            } else {
                                return (index + 1) + '. ' + mailbox.email + ' (密码: ' + mailbox.password + ')';
                            }
                        }).join('\n');
                    }

                    // 构建完整的凭据信息
                    const credentialsInfo = '登录凭据信息:\n\n' +
                        '用户邮箱: ' + data.user_email + '\n' +
                        '登录密码: ' + (data.password || '未设置') + mailboxList + '\n\n' +
                        '邮件服务器配置:\n' +
                        'SMTP服务器: ' + data.smtp_config.host + ':' + data.smtp_config.port + '\n' +
                        'IMAP服务器: ' + data.imap_config.host + ':' + data.imap_config.port + '\n' +
                        'POP3服务器: ' + data.pop3_config.host + ':' + data.pop3_config.port + '\n\n' +
                        (data.usage_note || '');

                    alert(credentialsInfo);
                } else {
                    showNotification('获取凭据失败: ' + result.message, 'error');
                }
            } catch (error) {
                showNotification('获取凭据失败: ' + error.message, 'error');
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
                    if (document.body.contains(notification)) {
                        document.body.removeChild(notification);
                    }
                }, 300);
            }, 3000);
        }

        // 点击模态框外部关闭
        window.onclick = function(event) {
            const mailboxModal = document.getElementById('mailboxModal');
            const generateModal = document.getElementById('generateModal');
            if (event.target == mailboxModal) {
                closeMailboxModal();
            }
            if (event.target == generateModal) {
                closeGenerateModal();
            }
        }

        // 页面加载时获取邮件和当前邮箱
        loadEmails(1, '');
        loadCurrentMailbox();

        // 页面加载时也加载域名列表
        loadDomains();

        // 页面获得焦点时重新加载域名列表（处理从其他页面返回的情况）
        window.addEventListener('focus', loadDomains);

        // 页面可见性改变时重新加载域名列表
        document.addEventListener('visibilitychange', function() {
            if (!document.hidden) {
                loadDomains();
            }
        });
    </script>
</body>
</html>`

	return tmpl
}

// generateEmailDetailPageTemplate 生成邮件详情页面模板
func (s *Server) generateEmailDetailPageTemplate(userEmail string, isAdmin bool, emailID string) string {
	tmpl := `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>NBEmail - 邮件详情</title>
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
            max-width: 1000px;
            margin: 20px auto;
            padding: 0 20px;
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
            margin-bottom: 20px;
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
            text-decoration: none;
            color: #333;
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

        /* 邮件内容 */
        .email-detail {
            background: white;
            border-radius: 15px;
            overflow: hidden;
            box-shadow: 0 5px 20px rgba(0,0,0,0.08);
        }
        .email-header {
            padding: 30px;
            border-bottom: 1px solid #f1f3f4;
            background: linear-gradient(135deg, rgba(102,126,234,0.05) 0%, rgba(255,255,255,1) 100%);
        }
        .email-subject {
            font-size: 24px;
            font-weight: 600;
            color: #333;
            margin-bottom: 20px;
            line-height: 1.3;
        }
        .email-meta {
            display: grid;
            grid-template-columns: auto 1fr;
            gap: 15px 20px;
            font-size: 14px;
        }
        .meta-label {
            font-weight: 600;
            color: #5a6c7d;
            text-align: right;
        }
        .meta-value {
            color: #333;
        }
        .email-body {
            padding: 30px;
            line-height: 1.6;
            font-size: 16px;
            color: #333;
            white-space: pre-wrap;
            word-wrap: break-word;
        }

        /* 加载状态 */
        .loading {
            text-align: center;
            padding: 80px 20px;
            color: #8e9aaf;
        }
        .loading .icon {
            font-size: 4rem;
            margin-bottom: 20px;
            opacity: 0.5;
        }

        /* 错误状态 */
        .error {
            text-align: center;
            padding: 80px 20px;
            color: #dc3545;
        }
        .error .icon {
            font-size: 4rem;
            margin-bottom: 20px;
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
            .container {
                flex-direction: column;
                padding: 0 10px;
            }
            .sidebar {
                width: 100%;
                order: 2;
                margin-top: 20px;
                background: white;
                border-radius: 12px;
                box-shadow: 0 2px 8px rgba(0,0,0,0.1);
                padding: 20px;
            }
            .main {
                order: 1;
            }
            .toolbar {
                flex-direction: column;
                gap: 15px;
            }
            .toolbar-left, .toolbar-right {
                width: 100%;
                justify-content: center;
            }
            .email-meta {
                grid-template-columns: 1fr;
                gap: 10px;
            }
            .meta-label {
                text-align: left;
                font-weight: 600;
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
        <div class="toolbar">
            <div class="toolbar-left">
                <a href="/inbox" class="btn">
                    <span>⬅️</span>返回收件箱
                </a>
                <h2 style="margin: 0; color: #333; font-weight: 600;">📧 邮件详情</h2>
            </div>
            <div class="toolbar-right">
                <button class="btn" onclick="replyEmail()">
                    <span>↩️</span>回复
                </button>
                <button class="btn" onclick="forwardEmail()">
                    <span>↗️</span>转发
                </button>
                <button class="btn btn-danger" onclick="deleteEmail()">
                    <span>🗑️</span>删除
                </button>
                <button class="btn" onclick="markAsRead()">
                    <span>✅</span>标记已读
                </button>
            </div>
        </div>
        <div class="email-detail" id="emailDetail">
            <div class="loading">
                <div class="icon">📬</div>
                <h3>正在加载邮件...</h3>
                <p>请稍候，正在获取邮件内容</p>
            </div>
        </div>
    </div>

    <script>
        const emailId = '` + emailID + `';
        let emailData = null;

        async function loadEmail() {
            try {
                const response = await fetch('/api/emails/' + emailId);
                const result = await response.json();
                if (result.success) {
                    emailData = result.data;
                    renderEmail();
                    // 自动标记为已读
                    if (!emailData.is_read) {
                        markAsRead();
                    }
                } else {
                    showError('邮件不存在或已被删除');
                }
            } catch (error) {
                showError('加载邮件失败，请重试');
            }
        }

        function renderEmail() {
            const emailDetail = document.getElementById('emailDetail');
            const date = new Date(emailData.created_at);
            const dateStr = date.toLocaleString('zh-CN', {
                year: 'numeric',
                month: 'long',
                day: 'numeric',
                hour: '2-digit',
                minute: '2-digit'
            });

            emailDetail.innerHTML = ` + "`" + `
                <div class="email-header">
                    <div class="email-subject">${emailData.subject}</div>
                    <div class="email-meta">
                        <div class="meta-label">发件人:</div>
                        <div class="meta-value">${emailData.from}</div>
                        <div class="meta-label">收件人:</div>
                        <div class="meta-value">${emailData.to}</div>
                        <div class="meta-label">时间:</div>
                        <div class="meta-value">${dateStr}</div>
                        <div class="meta-label">大小:</div>
                        <div class="meta-value">${formatSize(emailData.size)}</div>
                        <div class="meta-label">状态:</div>
                        <div class="meta-value">
                            ${emailData.is_read ? '<span style="color: #28a745;">已读</span>' : '<span style="color: #007bff;">未读</span>'}
                            ${emailData.is_sent ? '<span style="color: #6c757d; margin-left: 10px;">已发送</span>' : ''}
                        </div>
                    </div>
                </div>
                <div class="email-body">${emailData.body}</div>
            ` + "`" + `;
        }

        function showError(message) {
            const emailDetail = document.getElementById('emailDetail');
            emailDetail.innerHTML = ` + "`" + `
                <div class="error">
                    <div class="icon">❌</div>
                    <h3>加载失败</h3>
                    <p>${message}</p>
                    <button class="btn" onclick="loadEmail()" style="margin-top: 20px;">
                        <span>🔄</span>重试
                    </button>
                </div>
            ` + "`" + `;
        }

        function formatSize(bytes) {
            if (bytes === 0) return '0 B';
            const k = 1024;
            const sizes = ['B', 'KB', 'MB', 'GB'];
            const i = Math.floor(Math.log(bytes) / Math.log(k));
            return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
        }

        async function markAsRead() {
            if (!emailData || emailData.is_read) return;

            try {
                const response = await fetch('/api/emails/' + emailId + '/read', { method: 'PUT' });
                const result = await response.json();
                if (result.success) {
                    emailData.is_read = true;
                    renderEmail();
                }
            } catch (error) {
                console.error('标记已读失败:', error);
            }
        }

        async function deleteEmail() {
            if (!confirm('确定要删除这封邮件吗？')) return;

            try {
                const response = await fetch('/api/emails/' + emailId, { method: 'DELETE' });
                const result = await response.json();
                if (result.success) {
                    alert('邮件已删除');
                    window.close();
                } else {
                    alert('删除失败: ' + result.message);
                }
            } catch (error) {
                alert('删除失败，请重试');
            }
        }

        function replyEmail() {
            if (!emailData) return;
            const replySubject = emailData.subject.startsWith('Re: ') ? emailData.subject : 'Re: ' + emailData.subject;
            const replyTo = emailData.from;
            const replyBody = '\\n\\n--- 原始邮件 ---\\n' +
                             '发件人: ' + emailData.from + '\\n' +
                             '时间: ' + new Date(emailData.created_at).toLocaleString('zh-CN') + '\\n' +
                             '主题: ' + emailData.subject + '\\n\\n' +
                             emailData.body;

            const composeUrl = '/compose?to=' + encodeURIComponent(replyTo) +
                              '&subject=' + encodeURIComponent(replySubject) +
                              '&body=' + encodeURIComponent(replyBody);
            window.open(composeUrl, '_blank');
        }

        function forwardEmail() {
            if (!emailData) return;
            const forwardSubject = emailData.subject.startsWith('Fwd: ') ? emailData.subject : 'Fwd: ' + emailData.subject;
            const forwardBody = '\\n\\n--- 转发邮件 ---\\n' +
                               '发件人: ' + emailData.from + '\\n' +
                               '收件人: ' + emailData.to + '\\n' +
                               '时间: ' + new Date(emailData.created_at).toLocaleString('zh-CN') + '\\n' +
                               '主题: ' + emailData.subject + '\\n\\n' +
                               emailData.body;

            const composeUrl = '/compose?subject=' + encodeURIComponent(forwardSubject) +
                              '&body=' + encodeURIComponent(forwardBody);
            window.open(composeUrl, '_blank');
        }

        async function logout() {
            try {
                await fetch('/api/logout', { method: 'POST' });
                window.location.href = '/login';
            } catch (error) {
                window.location.href = '/login';
            }
        }

        // 页面加载时获取邮件
        loadEmail();
    </script>
</body>
</html>`

	return tmpl
}
