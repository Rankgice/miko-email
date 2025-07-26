package handlers

import (
	"miko-email/internal/model"
	"miko-email/internal/result"
	"miko-email/internal/services/dkim"
	"miko-email/internal/services/domain"
	"miko-email/internal/svc"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

type DomainHandler struct {
	domainService *domain.Service
	dkimService   *dkim.Service
	sessionStore  *sessions.CookieStore
	svcCtx        *svc.ServiceContext
}

func NewDomainHandler(domainService *domain.Service, dkimService *dkim.Service, sessionStore *sessions.CookieStore, svcCtx *svc.ServiceContext) *DomainHandler {
	return &DomainHandler{
		domainService: domainService,
		dkimService:   dkimService,
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
