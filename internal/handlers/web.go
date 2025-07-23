package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

type WebHandler struct {
	sessionStore *sessions.CookieStore
}

func NewWebHandler(sessionStore *sessions.CookieStore) *WebHandler {
	return &WebHandler{sessionStore: sessionStore}
}

// Home 首页
func (h *WebHandler) Home(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", gin.H{
		"title": "Miko邮箱系统",
	})
}

// LoginPage 登录页面
func (h *WebHandler) LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "用户登录",
	})
}

// RegisterPage 注册页面
func (h *WebHandler) RegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{
		"title": "用户注册",
	})
}

// AdminLoginPage 管理员登录页面
func (h *WebHandler) AdminLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_login.html", gin.H{
		"title": "管理员登录",
	})
}

// Dashboard 用户仪表板
func (h *WebHandler) Dashboard(c *gin.Context) {
	username := c.GetString("username")
	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"title":    "用户中心",
		"username": username,
	})
}

// ComposePage 写邮件页面
func (h *WebHandler) ComposePage(c *gin.Context) {
	username := c.GetString("username")
	c.HTML(http.StatusOK, "compose.html", gin.H{
		"title":    "写邮件",
		"username": username,
	})
}

// ForwardPage 转邮件页面
func (h *WebHandler) ForwardPage(c *gin.Context) {
	username := c.GetString("username")
	c.HTML(http.StatusOK, "forward.html", gin.H{
		"title":    "转邮件",
		"username": username,
	})
}

// InboxPage 收件箱页面
func (h *WebHandler) InboxPage(c *gin.Context) {
	username := c.GetString("username")
	c.HTML(http.StatusOK, "inbox.html", gin.H{
		"title":    "收件箱",
		"username": username,
	})
}

// SentPage 已发送页面
func (h *WebHandler) SentPage(c *gin.Context) {
	username := c.GetString("username")
	c.HTML(http.StatusOK, "sent.html", gin.H{
		"title":    "已发送",
		"username": username,
	})
}

// SettingsPage 设置页面
func (h *WebHandler) SettingsPage(c *gin.Context) {
	username := c.GetString("username")
	c.HTML(http.StatusOK, "settings.html", gin.H{
		"title":    "设置",
		"username": username,
	})
}

// MailboxesPage 邮箱管理页面
func (h *WebHandler) MailboxesPage(c *gin.Context) {
	username := c.GetString("username")
	c.HTML(http.StatusOK, "mailboxes.html", gin.H{
		"title":    "邮箱管理",
		"username": username,
	})
}

// AdminDashboard 管理员仪表板
func (h *WebHandler) AdminDashboard(c *gin.Context) {
	username := c.GetString("username")
	c.HTML(http.StatusOK, "admin_dashboard.html", gin.H{
		"title":    "管理员中心",
		"username": username,
	})
}

// UsersPage 用户管理页面
func (h *WebHandler) UsersPage(c *gin.Context) {
	username := c.GetString("username")
	c.HTML(http.StatusOK, "admin_users.html", gin.H{
		"title":    "用户管理",
		"username": username,
	})
}

// DomainsPage 域名管理页面
func (h *WebHandler) DomainsPage(c *gin.Context) {
	username := c.GetString("username")
	c.HTML(http.StatusOK, "admin_domains.html", gin.H{
		"title":    "域名管理",
		"username": username,
	})
}

// AdminMailboxesPage 管理员邮箱管理页面
func (h *WebHandler) AdminMailboxesPage(c *gin.Context) {
	username := c.GetString("username")
	c.HTML(http.StatusOK, "admin_mailboxes.html", gin.H{
		"title":    "邮箱管理",
		"username": username,
	})
}
