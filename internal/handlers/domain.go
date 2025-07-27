package handlers

import (
	"miko-email/internal/model"
	"miko-email/internal/result"
	"miko-email/internal/services/dkim"
	"miko-email/internal/services/domain"
	"miko-email/internal/svc"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

type DomainHandler struct {
	domainService *domain.Service
	dkimService   *dkim.Service
	sessionStore  *sessions.CookieStore
	svcCtx        *svc.ServiceContext
}

func NewDomainHandler(sessionStore *sessions.CookieStore, svcCtx *svc.ServiceContext) *DomainHandler {
	return &DomainHandler{
		domainService: domain.NewService(svcCtx),
		dkimService:   dkim.NewService("./dkim_keys"),
		sessionStore:  sessionStore,
		svcCtx:        svcCtx,
	}
}

type CreateDomainRequest struct {
	Name        string `json:"name" binding:"required"`
	MXRecord    string `json:"mx_record"`
	ARecord     string `json:"a_record"`
	TXTRecord   string `json:"txt_record"`
	SPFRecord   string `json:"spf_record"`
	DMARCRecord string `json:"dmarc_record"`
	DKIMRecord  string `json:"dkim_record"`
	PTRRecord   string `json:"ptr_record"`
}

type UpdateDomainRequest struct {
	MXRecord    string `json:"mx_record"`
	ARecord     string `json:"a_record"`
	TXTRecord   string `json:"txt_record"`
	SPFRecord   string `json:"spf_record"`
	DMARCRecord string `json:"dmarc_record"`
	DKIMRecord  string `json:"dkim_record"`
	PTRRecord   string `json:"ptr_record"`
}

// GetDomains 获取域名列表
func (h *DomainHandler) GetDomains(c *gin.Context) {
	domains, err := h.domainService.GetDomains()
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("获取域名列表失败"))
		return
	}

	c.JSON(http.StatusOK, result.DataResult("", domains))
}

// GetAvailableDomains 获取可用域名列表
func (h *DomainHandler) GetAvailableDomains(c *gin.Context) {
	domains, err := h.domainService.GetAvailableDomains()
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("获取可用域名列表失败"))
		return
	}

	c.JSON(http.StatusOK, result.DataResult("", domains))
}

// CreateDomain 创建域名
func (h *DomainHandler) CreateDomain(c *gin.Context) {
	var req CreateDomainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorReqParam)
		return
	}

	var domain *model.Domain
	var err error

	// 如果有额外的DNS记录，使用完整版本的创建函数
	if req.SPFRecord != "" || req.DMARCRecord != "" || req.DKIMRecord != "" || req.PTRRecord != "" {
		domain, err = h.domainService.CreateDomainWithAllRecords(
			req.Name, req.MXRecord, req.ARecord, req.TXTRecord,
			req.SPFRecord, req.DMARCRecord, req.DKIMRecord, req.PTRRecord)
	} else {
		domain, err = h.domainService.CreateDomain(req.Name, req.MXRecord, req.ARecord, req.TXTRecord)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result.DataResult("域名创建成功", domain))
}

// UpdateDomain 更新域名
func (h *DomainHandler) UpdateDomain(c *gin.Context) {
	domainIDStr := c.Param("id")
	domainID, err := strconv.ParseInt(domainIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("域名ID格式错误"))
		return
	}

	var req UpdateDomainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorReqParam)
		return
	}

	// 如果有额外的DNS记录，使用完整版本的更新函数
	if req.SPFRecord != "" || req.DMARCRecord != "" || req.DKIMRecord != "" || req.PTRRecord != "" {
		err = h.domainService.UpdateDomainWithAllRecords(
			domainID, req.MXRecord, req.ARecord, req.TXTRecord,
			req.SPFRecord, req.DMARCRecord, req.DKIMRecord, req.PTRRecord)
	} else {
		err = h.domainService.UpdateDomain(domainID, req.MXRecord, req.ARecord, req.TXTRecord)
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result.SimpleResult("域名更新成功"))
}

// VerifySenderConfiguration 验证发件配置
func (h *DomainHandler) VerifySenderConfiguration(c *gin.Context) {
	domainIDStr := c.Param("id")
	domainID, err := strconv.ParseInt(domainIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "域名ID格式错误"})
		return
	}

	domain, err := h.domainService.VerifySenderConfiguration(domainID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "发件配置验证完成",
		"data":    domain,
	})
}

// VerifyReceiverConfiguration 验证收件配置
func (h *DomainHandler) VerifyReceiverConfiguration(c *gin.Context) {
	domainIDStr := c.Param("id")
	domainID, err := strconv.ParseInt(domainIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "域名ID格式错误"})
		return
	}

	domain, err := h.domainService.VerifyReceiverConfiguration(domainID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "收件配置验证完成",
		"data":    domain,
	})
}

// DeleteDomain 删除域名
func (h *DomainHandler) DeleteDomain(c *gin.Context) {
	domainIDStr := c.Param("id")
	domainID, err := strconv.ParseInt(domainIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("域名ID格式错误"))
		return
	}

	err = h.domainService.DeleteDomain(domainID)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result.SimpleResult("域名删除成功"))
}

// VerifyDomain 验证域名
func (h *DomainHandler) VerifyDomain(c *gin.Context) {
	domainIDStr := c.Param("id")
	domainID, err := strconv.ParseInt(domainIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("域名ID格式错误"))
		return
	}

	domain, err := h.domainService.VerifyDomain(domainID)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult(err.Error()))
		return
	}

	message := "域名验证完成"
	if domain.IsVerified {
		message = "域名验证成功"
	} else {
		message = "域名验证失败，请检查DNS设置"
	}

	c.JSON(http.StatusOK, result.DataResult(message, domain))
}

// GetDomainDNSRecords 获取域名DNS记录
func (h *DomainHandler) GetDomainDNSRecords(c *gin.Context) {
	domainName := c.Query("domain")
	if domainName == "" {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("域名参数不能为空"))
		return
	}

	records := h.domainService.GetDNSRecords(domainName)

	data := gin.H{
		"domain":  domainName,
		"records": records,
	}

	c.JSON(http.StatusOK, result.DataResult("", data))
}

// GetServerInfo 获取服务器信息（IP地址等）
func (h *DomainHandler) GetServerInfo(c *gin.Context) {
	// 获取服务器公网IP
	publicIP, err := h.domainService.GetPublicIP()
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("获取服务器IP失败: "+err.Error()))
		return
	}

	data := gin.H{
		"public_ip": publicIP,
	}

	c.JSON(http.StatusOK, result.DataResult("", data))
}

// GetDomainDKIMRecord 获取域名的DKIM记录
func (h *DomainHandler) GetDomainDKIMRecord(c *gin.Context) {
	domainName := c.Query("domain")
	if domainName == "" {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("域名参数不能为空"))
		return
	}

	// 生成或获取DKIM记录
	dkimRecord, err := h.dkimService.GenerateDKIMRecord(domainName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("生成DKIM记录失败: "+err.Error()))
		return
	}

	// 获取公钥
	publicKey, err := h.dkimService.GetPublicKey(domainName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("获取公钥失败: "+err.Error()))
		return
	}

	data := gin.H{
		"domain":      domainName,
		"dkim_record": dkimRecord,
		"public_key":  publicKey,
		"selector":    h.dkimService.GetDKIMSelector(),
		"dkim_domain": h.dkimService.GetDKIMDomain(domainName),
	}

	c.JSON(http.StatusOK, result.DataResult("", data))
}

// GetDKIMRecord 获取域名的DKIM记录
func (h *DomainHandler) GetDKIMRecord(c *gin.Context) {
	domainName := c.Query("domain")
	if domainName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "域名参数不能为空"})
		return
	}

	// 生成或获取DKIM记录
	dkimRecord, err := h.dkimService.GenerateDKIMRecord(domainName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "生成DKIM记录失败: " + err.Error()})
		return
	}

	// 获取公钥
	publicKey, err := h.dkimService.GetPublicKey(domainName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取公钥失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"domain":      domainName,
			"selector":    h.dkimService.GetDKIMSelector(),
			"dkim_domain": h.dkimService.GetDKIMDomain(domainName),
			"record":      dkimRecord,
			"public_key":  publicKey,
		},
	})
}

// List 获取域名列表（API）
func (h *DomainHandler) List(c *gin.Context) {
	domains, err := h.domainService.GetDomains()
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("获取域名列表失败"))
		return
	}

	// 转换为API响应格式
	var domainResponses []gin.H
	for _, domain := range domains {
		domainResponses = append(domainResponses, gin.H{
			"id":          domain.Id,
			"name":        domain.Name,
			"description": domain.Name, // 使用域名作为描述
			"active":      domain.IsActive,
			"created_at":  domain.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, result.DataResult("", domainResponses))
}

// Create 创建域名（API）
func (h *DomainHandler) Create(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorReqParam)
		return
	}

	// 创建域名，使用默认的DNS记录
	domain, err := h.domainService.CreateDomain(req.Name, "", "", "")
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult(err.Error()))
		return
	}

	domainResponse := gin.H{
		"id":          domain.Id,
		"name":        domain.Name,
		"description": req.Description,
		"active":      domain.IsActive,
		"created_at":  domain.CreatedAt,
	}

	c.JSON(http.StatusOK, result.DataResult("域名创建成功", domainResponse))
}

// Delete 删除域名（API）
func (h *DomainHandler) Delete(c *gin.Context) {
	domainIDStr := c.Param("id")
	domainID, err := strconv.ParseInt(domainIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("域名ID格式错误"))
		return
	}

	err = h.domainService.DeleteDomain(domainID)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result.SimpleResult("域名删除成功"))
}

// GetAvailable 获取可用域名列表（公共API）
func (h *DomainHandler) GetAvailable(c *gin.Context) {
	domains, err := h.domainService.GetDomains()
	if err != nil {
		c.JSON(http.StatusInternalServerError, result.ErrorSimpleResult("获取域名列表失败"))
		return
	}

	// 转换为用户友好的格式，包含状态信息
	type AvailableDomain struct {
		Name        string `json:"name"`
		IsActive    bool   `json:"is_active"`
		IsVerified  bool   `json:"is_verified"`
		Description string `json:"description"`
		CreatedAt   string `json:"created_at"`
	}

	var availableDomains []AvailableDomain
	for _, domain := range domains {
		// 只返回已验证且活跃的域名
		if domain.IsVerified && domain.IsActive {
			description := "系统域名，可用于创建邮箱"
			if domain.Name == "suxinwl.com" {
				description = "主域名，已验证可用"
			}

			availableDomains = append(availableDomains, AvailableDomain{
				Name:        domain.Name,
				IsActive:    domain.IsActive,
				IsVerified:  domain.IsVerified,
				Description: description,
				CreatedAt:   domain.CreatedAt.Format("2006-01-02 15:04:05"),
			})
		}
	}

	c.JSON(http.StatusOK, result.DataResult("可用域名列表", availableDomains))
}

// GetSystemLogs 获取系统日志
func (h *DomainHandler) GetSystemLogs(c *gin.Context) {
	level := c.Query("level")
	module := c.Query("module")
	search := c.Query("search")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("pageSize", "50")

	// 模拟日志数据，实际应该从日志文件或数据库获取
	logs := []gin.H{
		{
			"id":        1,
			"level":     "info",
			"message":   "系统启动成功",
			"timestamp": time.Now().Add(-5 * time.Minute),
			"module":    "system",
			"details":   "所有服务已正常启动",
			"ip":        "127.0.0.1",
			"user":      "system",
		},
		{
			"id":        2,
			"level":     "warning",
			"message":   "存储空间使用率较高",
			"timestamp": time.Now().Add(-30 * time.Minute),
			"module":    "storage",
			"details":   "当前使用率: 85%，建议清理临时文件",
			"ip":        "127.0.0.1",
			"user":      "system",
		},
		{
			"id":        3,
			"level":     "error",
			"message":   "SMTP连接失败",
			"timestamp": time.Now().Add(-1 * time.Hour),
			"module":    "smtp",
			"details":   "Connection timeout after 30 seconds, 请检查SMTP服务器配置",
			"ip":        "127.0.0.1",
			"user":      "system",
		},
		{
			"id":        4,
			"level":     "info",
			"message":   "用户登录",
			"timestamp": time.Now().Add(-2 * time.Hour),
			"module":    "auth",
			"details":   "用户: admin@example.com 从 192.168.1.100 登录",
			"ip":        "192.168.1.100",
			"user":      "admin@example.com",
		},
		{
			"id":        5,
			"level":     "success",
			"message":   "数据库备份完成",
			"timestamp": time.Now().Add(-24 * time.Hour),
			"module":    "database",
			"details":   "备份文件: backup_20250127.sql，大小: 15.2MB",
			"ip":        "127.0.0.1",
			"user":      "system",
		},
		{
			"id":        6,
			"level":     "warning",
			"message":   "邮件队列积压",
			"timestamp": time.Now().Add(-3 * time.Hour),
			"module":    "mail",
			"details":   "待发送邮件数量: 156，建议检查邮件服务状态",
			"ip":        "127.0.0.1",
			"user":      "system",
		},
		{
			"id":        7,
			"level":     "error",
			"message":   "域名DNS验证失败",
			"timestamp": time.Now().Add(-4 * time.Hour),
			"module":    "dns",
			"details":   "域名 example.com 的SPF记录验证失败",
			"ip":        "127.0.0.1",
			"user":      "admin@example.com",
		},
		{
			"id":        8,
			"level":     "info",
			"message":   "新用户注册",
			"timestamp": time.Now().Add(-6 * time.Hour),
			"module":    "user",
			"details":   "新用户 user@example.com 注册成功",
			"ip":        "192.168.1.101",
			"user":      "user@example.com",
		},
	}

	// 应用过滤器
	filteredLogs := make([]gin.H, 0)
	for _, log := range logs {
		// 级别过滤
		if level != "" && level != "all" && log["level"] != level {
			continue
		}

		// 模块过滤
		if module != "" && module != "all" && log["module"] != module {
			continue
		}

		// 搜索过滤
		if search != "" {
			searchLower := strings.ToLower(search)
			message := strings.ToLower(log["message"].(string))
			details := strings.ToLower(log["details"].(string))
			moduleStr := strings.ToLower(log["module"].(string))

			if !strings.Contains(message, searchLower) &&
				!strings.Contains(details, searchLower) &&
				!strings.Contains(moduleStr, searchLower) {
				continue
			}
		}

		filteredLogs = append(filteredLogs, log)
	}

	data := gin.H{
		"logs":     filteredLogs,
		"total":    len(filteredLogs),
		"page":     page,
		"pageSize": pageSize,
		"modules":  []string{"system", "auth", "smtp", "database", "storage", "mail", "dns", "user"},
		"levels":   []string{"info", "warning", "error", "success"},
	}

	c.JSON(http.StatusOK, result.DataResult("", data))
}

// GetSystemSettings 获取系统设置
func (h *DomainHandler) GetSystemSettings(c *gin.Context) {
	// 模拟系统设置数据，实际应该从配置文件或数据库获取
	settings := gin.H{
		"general": gin.H{
			"systemName":         "Miko邮箱系统",
			"systemDescription":  "企业级邮箱管理系统",
			"defaultTheme":       "dark",
			"timezone":           "Asia/Shanghai",
			"language":           "zh-CN",
			"maxUsers":           1000,
			"enableRegistration": true,
		},
		"security": gin.H{
			"sessionTimeout":    30,
			"minPasswordLength": 8,
			"enableTwoFactor":   false,
			"enableLoginLock":   true,
			"maxLoginAttempts":  5,
			"lockoutDuration":   15,
			"enableSSL":         true,
			"enableFirewall":    true,
		},
		"email": gin.H{
			"smtpHost":          "localhost",
			"smtpPort":          587,
			"imapPort":          143,
			"pop3Port":          110,
			"senderEmail":       "noreply@example.com",
			"emailSignature":    "Miko邮箱系统自动发送，请勿回复。",
			"maxAttachmentSize": 25,
			"enableEncryption":  true,
		},
		"storage": gin.H{
			"maxStorageSpace":   100,
			"enableAutoCleanup": false,
			"cleanupDays":       30,
			"backupEnabled":     true,
			"backupInterval":    24,
			"compressionLevel":  6,
		},
	}

	c.JSON(http.StatusOK, result.DataResult("", settings))
}

// UpdateSystemSettings 更新系统设置
func (h *DomainHandler) UpdateSystemSettings(c *gin.Context) {
	var settings map[string]interface{}
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("参数格式错误: "+err.Error()))
		return
	}

	// 这里应该验证设置并保存到配置文件或数据库
	// 模拟保存过程

	c.JSON(http.StatusOK, result.DataResult("系统设置更新成功", nil))
}

// GetSystemStatus 获取系统状态
func (h *DomainHandler) GetSystemStatus(c *gin.Context) {
	status := gin.H{
		"services": []gin.H{
			{
				"name":        "SMTP服务",
				"status":      "running",
				"port":        587,
				"uptime":      "2天3小时15分钟",
				"connections": 45,
			},
			{
				"name":        "IMAP服务",
				"status":      "running",
				"port":        143,
				"uptime":      "2天3小时15分钟",
				"connections": 23,
			},
			{
				"name":        "POP3服务",
				"status":      "stopped",
				"port":        110,
				"uptime":      "0分钟",
				"connections": 0,
			},
			{
				"name":        "Web服务",
				"status":      "running",
				"port":        8080,
				"uptime":      "2天3小时15分钟",
				"connections": 12,
			},
		},
		"system": gin.H{
			"cpu":    gin.H{"usage": 25.6, "cores": 4},
			"memory": gin.H{"used": 2.1, "total": 8.0, "usage": 26.25},
			"disk":   gin.H{"used": 45.2, "total": 100.0, "usage": 45.2},
			"network": gin.H{
				"bytesIn":    "1.2GB",
				"bytesOut":   "856MB",
				"packetsIn":  125000,
				"packetsOut": 98000,
			},
		},
		"database": gin.H{
			"status":      "running",
			"connections": 8,
			"size":        "156MB",
			"lastBackup":  "2025-01-27 02:00:00",
		},
	}

	c.JSON(http.StatusOK, result.DataResult("", status))
}

// UpdateDomainStatus 更新域名状态
func (h *DomainHandler) UpdateDomainStatus(c *gin.Context) {
	domainIDStr := c.Param("id")
	domainID, err := strconv.ParseInt(domainIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("域名ID格式错误"))
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("请求参数错误"))
		return
	}

	var isActive bool
	var message string

	switch req.Status {
	case "enabled":
		isActive = true
		message = "域名已启用"
	case "disabled":
		isActive = false
		message = "域名已禁用"
	default:
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的状态值，请使用 enabled 或 disabled"))
		return
	}

	err = h.domainService.UpdateDomainStatus(domainID, isActive)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result.DataResult(message, nil))
}
