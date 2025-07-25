package handlers

import (
	"miko-email/internal/result"
	"miko-email/internal/svc"
	"net/http"
	"strconv"

	"miko-email/internal/services/domain"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

type DomainHandler struct {
	domainService *domain.Service
	sessionStore  *sessions.CookieStore
	svcCtx        *svc.ServiceContext
}

func NewDomainHandler(domainService *domain.Service, sessionStore *sessions.CookieStore, svcCtx *svc.ServiceContext) *DomainHandler {
	return &DomainHandler{
		domainService: domainService,
		sessionStore:  sessionStore,
		svcCtx:        svcCtx,
	}
}

type CreateDomainRequest struct {
	Name      string `json:"name" binding:"required"`
	MXRecord  string `json:"mx_record"`
	ARecord   string `json:"a_record"`
	TXTRecord string `json:"txt_record"`
}

type UpdateDomainRequest struct {
	MXRecord  string `json:"mx_record"`
	ARecord   string `json:"a_record"`
	TXTRecord string `json:"txt_record"`
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

	domain, err := h.domainService.CreateDomain(req.Name, req.MXRecord, req.ARecord, req.TXTRecord)
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

	err = h.domainService.UpdateDomain(domainID, req.MXRecord, req.ARecord, req.TXTRecord)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result.SimpleResult("域名更新成功"))
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
