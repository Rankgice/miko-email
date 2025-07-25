package handlers

import (
	"miko-email/internal/svc"
	"net/http"

	"miko-email/internal/services/auth"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

type AuthHandler struct {
	authService  *auth.Service
	sessionStore *sessions.CookieStore
	svcCtx       *svc.ServiceContext
}

func NewAuthHandler(authService *auth.Service, sessionStore *sessions.CookieStore, svcCtx *svc.ServiceContext) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		sessionStore: sessionStore,
		svcCtx:       svcCtx,
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username     string `json:"username" binding:"required"`
	Password     string `json:"password" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	DomainPrefix string `json:"domain_prefix" binding:"required"`
	DomainID     int64  `json:"domain_id" binding:"required"`
	InviteCode   string `json:"invite_code"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数错误"})
		return
	}

	user, err := h.authService.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": err.Error()})
		return
	}

	// 创建会话
	session, err := h.sessionStore.Get(c.Request, "miko-session")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "会话创建失败"})
		return
	}

	session.Values["user_id"] = user.Id
	session.Values["username"] = user.Username
	session.Values["is_admin"] = false

	if err := session.Save(c.Request, c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "会话保存失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "登录成功",
		"data": gin.H{
			"user": gin.H{
				"id":           user.Id,
				"username":     user.Username,
				"email":        user.Email,
				"contribution": user.Contribution,
				"is_admin":     false,
			},
		},
	})
}

// AdminLogin 管理员登录
func (h *AuthHandler) AdminLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数错误"})
		return
	}

	admin, err := h.authService.AuthenticateAdmin(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": err.Error()})
		return
	}

	// 创建会话
	session, err := h.sessionStore.Get(c.Request, "miko-session")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "会话创建失败"})
		return
	}

	session.Values["user_id"] = admin.Id
	session.Values["username"] = admin.Username
	session.Values["is_admin"] = true

	if err := session.Save(c.Request, c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "会话保存失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "登录成功",
		"data": gin.H{
			"user": gin.H{
				"id":           admin.Id,
				"username":     admin.Username,
				"email":        admin.Email,
				"contribution": admin.Contribution,
				"is_admin":     true,
			},
		},
	})
}

// Register 用户注册
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数错误"})
		return
	}

	user, err := h.authService.RegisterUser(req.Username, req.Password, req.Email, req.DomainPrefix, req.DomainID, req.InviteCode)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "注册成功",
		"data": gin.H{
			"user": gin.H{
				"id":           user.Id,
				"username":     user.Username,
				"email":        user.Email,
				"contribution": user.Contribution,
				"invite_code":  user.InviteCode,
			},
		},
	})
}

// Logout 用户登出
func (h *AuthHandler) Logout(c *gin.Context) {
	session, err := h.sessionStore.Get(c.Request, "miko-session")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "会话错误"})
		return
	}

	// 清除会话
	session.Values = make(map[interface{}]interface{})
	session.Options.MaxAge = -1

	if err := session.Save(c.Request, c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "登出失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "登出成功"})
}

// GetProfile 获取用户信息
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := c.GetInt("user_id")
	isAdmin := c.GetBool("is_admin")

	if isAdmin {
		// 管理员用户
		admin, err := h.svcCtx.AdminModel.GetById(int64(userID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取管理员信息失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"id":           admin.Id,
				"username":     admin.Username,
				"email":        admin.Email,
				"contribution": admin.Contribution,
				"invite_code":  admin.InviteCode,
				"is_admin":     true,
				"created_at":   admin.CreatedAt,
			},
		})
	} else {
		// 普通用户
		user, err := h.svcCtx.UserModel.GetById(int64(userID))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取用户信息失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"id":           user.Id,
				"username":     user.Username,
				"email":        user.Email,
				"contribution": user.Contribution,
				"invite_code":  user.InviteCode,
				"invited_by":   user.InvitedBy,
				"is_admin":     false,
				"created_at":   user.CreatedAt,
			},
		})
	}
}

// ChangePassword 修改密码
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "请求参数错误"})
		return
	}

	userID := int64(c.GetInt("user_id"))
	isAdmin := c.GetBool("is_admin")

	err := h.authService.ChangePassword(userID, req.OldPassword, req.NewPassword, isAdmin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "密码修改成功"})
}
