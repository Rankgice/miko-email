package server

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"nbemail/internal/dns"
	"nbemail/internal/models"
)

// handleVerifyDomain 验证域名DNS配置
func (s *Server) handleVerifyDomain(c *gin.Context) {
	if !c.GetBool("is_admin") {
		c.JSON(http.StatusForbidden, models.Response{
			Success: false,
			Message: "权限不足",
		})
		return
	}

	domainID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "域名ID无效",
		})
		return
	}

	// 获取域名信息
	var domain models.Domain
	err = s.db.QueryRow("SELECT id, name FROM domains WHERE id = ?", domainID).Scan(&domain.ID, &domain.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.Response{
				Success: false,
				Message: "域名不存在",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.Response{
				Success: false,
				Message: "查询域名失败",
			})
		}
		return
	}

	// 获取服务器IP
	serverIP, err := dns.GetServerPublicIP()
	if err != nil {
		serverIP = "127.0.0.1" // 使用本地IP作为备选
	}

	// 创建DNS验证器
	verifier := dns.NewDNSVerifier(serverIP)
	result := verifier.VerifyDomain(domain.Name)

	// 更新数据库中的验证状态
	now := time.Now()
	mxRecord := ""
	if len(result.MXRecords) > 0 {
		mxRecord = strings.Join(result.MXRecords, "; ")
	}

	_, err = s.db.Exec(`
		UPDATE domains 
		SET dns_verified = ?, mx_record = ?, last_verified = ?, updated_at = ?
		WHERE id = ?
	`, result.Success, mxRecord, now, now, domainID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "更新验证状态失败",
		})
		return
	}

	// 返回验证结果
	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "DNS验证完成",
		Data: map[string]interface{}{
			"domain_id":    domainID,
			"domain_name":  domain.Name,
			"server_ip":    serverIP,
			"dns_verified": result.Success,
			"verification_result": result,
		},
	})
}

// handleGetDNSInstructions 获取DNS配置说明
func (s *Server) handleGetDNSInstructions(c *gin.Context) {
	if !c.GetBool("is_admin") {
		c.JSON(http.StatusForbidden, models.Response{
			Success: false,
			Message: "权限不足",
		})
		return
	}

	domainID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "域名ID无效",
		})
		return
	}

	// 获取域名信息
	var domainName string
	err = s.db.QueryRow("SELECT name FROM domains WHERE id = ?", domainID).Scan(&domainName)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.Response{
				Success: false,
				Message: "域名不存在",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.Response{
				Success: false,
				Message: "查询域名失败",
			})
		}
		return
	}

	// 获取服务器IP
	serverIP, err := dns.GetServerPublicIP()
	if err != nil {
		serverIP = "127.0.0.1" // 使用本地IP作为备选
	}

	// 生成DNS配置说明
	instructions := dns.GenerateDNSInstructions(domainName, serverIP)

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "获取DNS配置说明成功",
		Data:    instructions,
	})
}

// handleCheckDNSPropagation 检查DNS传播状态
func (s *Server) handleCheckDNSPropagation(c *gin.Context) {
	if !c.GetBool("is_admin") {
		c.JSON(http.StatusForbidden, models.Response{
			Success: false,
			Message: "权限不足",
		})
		return
	}

	domainID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Response{
			Success: false,
			Message: "域名ID无效",
		})
		return
	}

	// 获取域名信息
	var domainName string
	err = s.db.QueryRow("SELECT name FROM domains WHERE id = ?", domainID).Scan(&domainName)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, models.Response{
				Success: false,
				Message: "域名不存在",
			})
		} else {
			c.JSON(http.StatusInternalServerError, models.Response{
				Success: false,
				Message: "查询域名失败",
			})
		}
		return
	}

	// 检查DNS传播状态
	propagation := dns.CheckDNSPropagation(domainName)

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "DNS传播检查完成",
		Data:    propagation,
	})
}

// validateDomainName 验证域名格式
func validateDomainName(domain string) error {
	return dns.ValidateDomainName(domain)
}

// handleBatchVerifyDomains 批量验证所有域名
func (s *Server) handleBatchVerifyDomains(c *gin.Context) {
	if !c.GetBool("is_admin") {
		c.JSON(http.StatusForbidden, models.Response{
			Success: false,
			Message: "权限不足",
		})
		return
	}

	// 获取所有活跃域名
	rows, err := s.db.Query("SELECT id, name FROM domains WHERE is_active = 1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Response{
			Success: false,
			Message: "查询域名失败",
		})
		return
	}
	defer rows.Close()

	// 获取服务器IP
	serverIP, err := dns.GetServerPublicIP()
	if err != nil {
		serverIP = "127.0.0.1" // 使用本地IP作为备选
	}

	// 创建DNS验证器
	verifier := dns.NewDNSVerifier(serverIP)
	
	var results []map[string]interface{}
	var successCount, failCount int

	for rows.Next() {
		var domainID int
		var domainName string
		if err := rows.Scan(&domainID, &domainName); err != nil {
			continue
		}

		// 验证域名
		result := verifier.VerifyDomain(domainName)
		
		// 更新数据库
		now := time.Now()
		mxRecord := ""
		if len(result.MXRecords) > 0 {
			mxRecord = strings.Join(result.MXRecords, "; ")
		}

		s.db.Exec(`
			UPDATE domains 
			SET dns_verified = ?, mx_record = ?, last_verified = ?, updated_at = ?
			WHERE id = ?
		`, result.Success, mxRecord, now, now, domainID)

		// 记录结果
		results = append(results, map[string]interface{}{
			"domain_id":   domainID,
			"domain_name": domainName,
			"success":     result.Success,
			"message":     result.Message,
		})

		if result.Success {
			successCount++
		} else {
			failCount++
		}
	}

	c.JSON(http.StatusOK, models.Response{
		Success: true,
		Message: "批量验证完成",
		Data: map[string]interface{}{
			"total":         len(results),
			"success_count": successCount,
			"fail_count":    failCount,
			"results":       results,
		},
	})
}
