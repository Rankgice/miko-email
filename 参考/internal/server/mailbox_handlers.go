package server

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"nbemail/internal/models"
)

// handleGetMailboxes 获取用户的邮箱列表
func (s *Server) handleGetMailboxes(c *gin.Context) {
	userID := c.GetInt("user_id")

	query := `
		SELECT m.id, m.user_id, m.email, m.domain_id, m.is_active, m.is_current, 
		       m.created_at, m.updated_at, d.name as domain_name
		FROM mailboxes m
		LEFT JOIN domains d ON m.domain_id = d.id
		WHERE m.user_id = ? AND m.is_active = 1
		ORDER BY m.is_current DESC, m.created_at DESC
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "获取邮箱列表失败",
		})
		return
	}
	defer rows.Close()

	var mailboxes []map[string]interface{}
	for rows.Next() {
		var mailbox models.Mailbox
		var domainName string
		err := rows.Scan(
			&mailbox.ID, &mailbox.UserID, &mailbox.Email, &mailbox.DomainID,
			&mailbox.IsActive, &mailbox.IsCurrent, &mailbox.CreatedAt, &mailbox.UpdatedAt,
			&domainName,
		)
		if err != nil {
			continue
		}

		mailboxes = append(mailboxes, map[string]interface{}{
			"id":          mailbox.ID,
			"user_id":     mailbox.UserID,
			"email":       mailbox.Email,
			"domain_id":   mailbox.DomainID,
			"domain_name": domainName,
			"is_active":   mailbox.IsActive,
			"is_current":  mailbox.IsCurrent,
			"created_at":  mailbox.CreatedAt,
			"updated_at":  mailbox.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "获取邮箱列表成功",
		Data:    mailboxes,
	})
}

// handleGenerateMailboxes 批量生成邮箱
func (s *Server) handleGenerateMailboxes(c *gin.Context) {
	var req models.GenerateMailboxesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	// 获取当前用户ID
	userID := c.GetInt("user_id")

	// 验证域名是否存在且用户有权限使用
	var domainName string
	var domainUserID *int
	err := s.db.QueryRow("SELECT name, user_id FROM domains WHERE id = ? AND is_active = 1", req.DomainID).Scan(&domainName, &domainUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, models.Response{
				Success: false,
				Message: "域名不存在或已停用",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.Response{
				Success: false,
				Message: "验证域名失败",
			})
		}
		return
	}

	// 检查用户是否有权限使用该域名
	// 如果域名有指定用户（domainUserID不为空），则只有该用户可以使用
	// 如果域名没有指定用户（domainUserID为空），则所有用户都可以使用
	if domainUserID != nil && *domainUserID != userID {
		c.JSON(http.StatusForbidden, models.Response{
			Success: false,
			Message: "您没有权限使用该域名",
		})
		return
	}

	// 生成邮箱（属于当前用户）
	var generatedAccounts []map[string]interface{}
	for i := 0; i < req.Count; i++ {
		var email string

		// 如果前缀是完整的邮箱名（不包含@），直接使用
		if req.Prefix != "" && !strings.Contains(req.Prefix, "@") && len(req.Prefix) > 5 {
			// 检查前缀是否看起来像完整的邮箱名（包含数字和字母）
			hasDigit := false
			hasLetter := false
			for _, r := range req.Prefix {
				if r >= '0' && r <= '9' {
					hasDigit = true
				}
				if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
					hasLetter = true
				}
			}

			if hasDigit && hasLetter {
				email = req.Prefix + "@" + domainName
			} else {
				email = s.generateRandomEmail(req.Prefix, domainName)
			}
		} else {
			email = s.generateRandomEmail(req.Prefix, domainName)
		}

		// 检查邮箱是否已存在
		var count int
		err := s.db.QueryRow("SELECT COUNT(*) FROM mailboxes WHERE email = ?", email).Scan(&count)
		if err != nil || count > 0 {
			// 如果已存在，重新生成
			i--
			continue
		}

		// 生成随机密码（用于显示，实际邮箱使用当前用户的登录凭据）
		password := s.generateRandomPassword()

		// 插入邮箱（属于当前用户）
		_, err = s.db.Exec(`
			INSERT INTO mailboxes (user_id, email, domain_id, password, is_active, is_current, created_at, updated_at)
			VALUES (?, ?, ?, ?, 1, 0, ?, ?)
		`, userID, email, req.DomainID, password, time.Now(), time.Now())

		if err != nil {
			log.Printf("创建邮箱失败: %v", err)
			i--
			continue
		}

		generatedAccounts = append(generatedAccounts, map[string]interface{}{
			"email":    email,
			"password": password,
		})
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: fmt.Sprintf("成功生成 %d 个邮箱账号", len(generatedAccounts)),
		Data: map[string]interface{}{
			"accounts": generatedAccounts,
			"count":    len(generatedAccounts),
		},
	})
}

// handleSwitchMailbox 切换当前使用的邮箱
func (s *Server) handleSwitchMailbox(c *gin.Context) {
	userID := c.GetInt("user_id")

	var req models.SwitchMailboxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	// 验证邮箱是否属于当前用户
	var email string
	err := s.db.QueryRow("SELECT email FROM mailboxes WHERE id = ? AND user_id = ? AND is_active = 1", 
		req.MailboxID, userID).Scan(&email)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, models.Response{
				Success: false,
				Message: "邮箱不存在或无权限",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.Response{
				Success: false,
				Message: "验证邮箱失败",
			})
		}
		return
	}

	// 开始事务
	tx, err := s.db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "开始事务失败",
		})
		return
	}
	defer tx.Rollback()

	// 将用户的所有邮箱设为非当前
	_, err = tx.Exec("UPDATE mailboxes SET is_current = 0 WHERE user_id = ?", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "更新邮箱状态失败",
		})
		return
	}

	// 设置指定邮箱为当前
	_, err = tx.Exec("UPDATE mailboxes SET is_current = 1 WHERE id = ? AND user_id = ?", 
		req.MailboxID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "设置当前邮箱失败",
		})
		return
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "提交事务失败",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "切换邮箱成功",
		Data: map[string]interface{}{
			"current_email": email,
		},
	})
}

// handleDeleteMailbox 删除邮箱
func (s *Server) handleDeleteMailbox(c *gin.Context) {
	userID := c.GetInt("user_id")
	mailboxID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "无效的邮箱ID",
		})
		return
	}

	// 验证邮箱是否属于当前用户
	var isCurrent bool
	err = s.db.QueryRow("SELECT is_current FROM mailboxes WHERE id = ? AND user_id = ? AND is_active = 1", 
		mailboxID, userID).Scan(&isCurrent)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, models.Response{
				Success: false,
				Message: "邮箱不存在或无权限",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.Response{
				Success: false,
				Message: "验证邮箱失败",
			})
		}
		return
	}

	// 如果是当前使用的邮箱，不允许删除
	if isCurrent {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "不能删除当前使用的邮箱，请先切换到其他邮箱",
		})
		return
	}

	// 软删除邮箱（设为不活跃）
	_, err = s.db.Exec("UPDATE mailboxes SET is_active = 0 WHERE id = ? AND user_id = ?", 
		mailboxID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "删除邮箱失败",
		})
		return
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "删除邮箱成功",
	})
}

// generateRandomEmail 生成随机邮箱地址
func (s *Server) generateRandomEmail(prefix, domain string) string {
	if prefix == "" {
		prefix = "user"
	}

	// 生成随机数字后缀
	randomNum, _ := rand.Int(rand.Reader, big.NewInt(999999))
	suffix := fmt.Sprintf("%06d", randomNum.Int64())

	// 生成随机字符串
	chars := "abcdefghijklmnopqrstuvwxyz0123456789"
	randomStr := ""
	for i := 0; i < 4; i++ {
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		randomStr += string(chars[idx.Int64()])
	}

	return fmt.Sprintf("%s%s%s@%s", prefix, randomStr, suffix, domain)
}

// generateRandomPassword 生成随机密码
func (s *Server) generateRandomPassword() string {
	const (
		uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		lowercase = "abcdefghijklmnopqrstuvwxyz"
		numbers   = "0123456789"
		length    = 12
	)

	allChars := uppercase + lowercase + numbers
	password := make([]byte, length)

	// 确保至少包含每种类型的字符
	charSets := []string{uppercase, lowercase, numbers}
	for i, charset := range charSets {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		password[i] = charset[n.Int64()]
	}

	// 填充剩余长度
	for i := 3; i < length; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(allChars))))
		password[i] = allChars[n.Int64()]
	}

	// 打乱密码字符顺序
	for i := length - 1; i > 0; i-- {
		j, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		password[i], password[j.Int64()] = password[j.Int64()], password[i]
	}

	return string(password)
}

// handleGetMailboxCredentials 获取邮箱的登录凭据（用于邮件客户端配置）
func (s *Server) handleGetMailboxCredentials(c *gin.Context) {
	userID := c.GetInt("user_id")

	// 获取当前用户的登录凭据
	var email string
	var plainPassword sql.NullString
	err := s.db.QueryRow(`
		SELECT email, plain_password
		FROM users
		WHERE id = ?
	`, userID).Scan(&email, &plainPassword)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.Response{
				Success: false,
				Message: "未找到用户信息",
			})
		} else {
			log.Printf("获取用户信息失败: %v", err)
			c.JSON(http.StatusInternalServerError, models.Response{
				Success: false,
				Message: "获取用户信息失败",
			})
		}
		return
	}

	// 处理明文密码
	password := "未设置明文密码"
	if plainPassword.Valid && plainPassword.String != "" {
		password = plainPassword.String
	}

	// 获取当前用户的邮箱列表（包含密码）
	var mailboxes []map[string]string
	rows, err := s.db.Query(`
		SELECT email, password FROM mailboxes
		WHERE user_id = ?
		ORDER BY is_current DESC, created_at DESC
	`, userID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var mailboxEmail, mailboxPassword string
			if rows.Scan(&mailboxEmail, &mailboxPassword) == nil {
				mailboxes = append(mailboxes, map[string]string{
					"email":    mailboxEmail,
					"password": mailboxPassword,
				})
			}
		}
		log.Printf("用户 %d 的邮箱数量: %d", userID, len(mailboxes))
	} else {
		log.Printf("查询邮箱失败: %v", err)
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "获取邮箱凭据成功",
		Data: map[string]interface{}{
			"user_email": email,
			"password":   password,
			"mailboxes":  mailboxes,
			"smtp_config": map[string]interface{}{
				"host": "localhost",
				"port": s.config.SMTPPort,
				"tls":  false,
				"auth": "plain",
			},
			"imap_config": map[string]interface{}{
				"host": "localhost",
				"port": s.config.IMAPPort,
				"tls":  false,
				"auth": "plain",
			},
			"pop3_config": map[string]interface{}{
				"host": "localhost",
				"port": s.config.POP3Port,
				"tls":  false,
				"auth": "plain",
			},
			"usage_note": "使用您的用户账号凭据（" + email + "）登录邮件客户端，可以收发所有关联邮箱的邮件",
		},
	})
}
