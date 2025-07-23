package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"nbemail/internal/models"
)

// authMiddleware 认证中间件
func (s *Server) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 对于API请求，检查Cookie中的token
		token, err := c.Cookie("token")
		if err != nil || token == "" {
			if c.Request.URL.Path[:4] == "/api" {
				c.JSON(http.StatusUnauthorized, models.Response{
					Success: false,
					Message: "未登录",
				})
				c.Abort()
				return
			} else {
				c.Redirect(http.StatusFound, "/login")
				c.Abort()
				return
			}
		}

		// 验证用户是否存在
		var userID int
		var isAdmin bool
		err = s.db.QueryRow("SELECT id, is_admin FROM users WHERE email = ?", token).Scan(&userID, &isAdmin)
		if err != nil {
			if c.Request.URL.Path[:4] == "/api" {
				c.JSON(http.StatusUnauthorized, models.Response{
					Success: false,
					Message: "用户不存在",
				})
				c.Abort()
				return
			} else {
				c.Redirect(http.StatusFound, "/login")
				c.Abort()
				return
			}
		}

		// 设置用户信息到上下文
		c.Set("user_id", userID)
		c.Set("user_email", token)
		c.Set("is_admin", isAdmin)
		c.Next()
	}
}

// adminMiddleware 管理员中间件
func (s *Server) adminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get("is_admin")
		if !exists || !isAdmin.(bool) {
			if c.Request.URL.Path[:4] == "/api" {
				c.JSON(http.StatusForbidden, models.Response{
					Success: false,
					Message: "权限不足",
				})
				c.Abort()
				return
			} else {
				c.Redirect(http.StatusFound, "/inbox")
				c.Abort()
				return
			}
		}
		c.Next()
	}
}