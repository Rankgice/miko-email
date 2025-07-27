package handlers

import (
	"crypto/rand"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"miko-email/internal/result"
	"miko-email/internal/svc"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

// CaptchaHandler 验证码处理器
type CaptchaHandler struct {
	sessionStore sessions.Store
	svcCtx       *svc.ServiceContext
}

// NewCaptchaHandler 创建验证码处理器
func NewCaptchaHandler(sessionStore sessions.Store, svcCtx *svc.ServiceContext) *CaptchaHandler {
	return &CaptchaHandler{
		sessionStore: sessionStore,
		svcCtx:       svcCtx,
	}
}

// CaptchaCode 验证码结构
type CaptchaCode struct {
	ID            int64      `json:"id"`
	Email         string     `json:"email"`
	Code          string     `json:"code"`
	SentAt        time.Time  `json:"sentAt"`
	ExpiryMinutes int        `json:"expiryMinutes"`
	IsUsed        bool       `json:"isUsed"`
	UsedAt        *time.Time `json:"usedAt,omitempty"`
	Purpose       string     `json:"purpose"` // login, register, reset_password
}

// CaptchaRule 验证码提取规则结构
type CaptchaRule struct {
	ID                int64     `json:"id"`
	Name              string    `json:"name"`
	Type              string    `json:"type"`              // regex, keyword, position
	SenderPattern     string    `json:"senderPattern"`     // 发件人匹配模式
	SubjectPattern    string    `json:"subjectPattern"`    // 主题匹配模式
	ExtractionPattern string    `json:"extractionPattern"` // 提取规则
	Priority          string    `json:"priority"`          // high, medium, low
	Enabled           bool      `json:"enabled"`
	Description       string    `json:"description"`
	MatchCount        int       `json:"matchCount"`
	SuccessCount      int       `json:"successCount"`
	LastUsed          time.Time `json:"lastUsed"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

// CaptchaStats 验证码统计
type CaptchaStats struct {
	TotalCodes       int     `json:"totalCodes"`
	TodaySent        int     `json:"todaySent"`
	ActiveCodes      int     `json:"activeCodes"`
	AvgSendTime      float64 `json:"avgSendTime"`
	TodayIncrease    int     `json:"todayIncrease"`
	SentIncrease     int     `json:"sentIncrease"`
	ActivePercentage float64 `json:"activePercentage"`
	SpeedImprovement int     `json:"speedImprovement"`
}

// 内存存储验证码（生产环境应使用数据库）
var captchaCodes = make(map[int64]*CaptchaCode)
var codeCounter int64 = 0

// GetCodes 获取验证码列表
func (h *CaptchaHandler) GetCodes(c *gin.Context) {
	// 检查管理员权限
	isAdmin := c.GetBool("is_admin")
	if !isAdmin {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("需要管理员权限"))
		return
	}

	// 转换为切片并排序
	codes := make([]*CaptchaCode, 0, len(captchaCodes))
	for _, code := range captchaCodes {
		codes = append(codes, code)
	}

	// 按发送时间倒序排序
	for i := 0; i < len(codes)-1; i++ {
		for j := i + 1; j < len(codes); j++ {
			if codes[i].SentAt.Before(codes[j].SentAt) {
				codes[i], codes[j] = codes[j], codes[i]
			}
		}
	}

	c.JSON(http.StatusOK, result.DataResult("验证码列表", codes))
}

// GetStats 获取验证码统计
func (h *CaptchaHandler) GetStats(c *gin.Context) {
	// 检查管理员权限
	isAdmin := c.GetBool("is_admin")
	if !isAdmin {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("需要管理员权限"))
		return
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	totalCodes := len(captchaCodes)
	todaySent := 0
	activeCodes := 0

	for _, code := range captchaCodes {
		// 统计今日发送
		if code.SentAt.After(today) {
			todaySent++
		}

		// 统计有效验证码
		if !h.isCodeExpired(code) && !code.IsUsed {
			activeCodes++
		}
	}

	activePercentage := 0.0
	if totalCodes > 0 {
		activePercentage = float64(activeCodes) / float64(totalCodes) * 100
	}

	stats := CaptchaStats{
		TotalCodes:       totalCodes,
		TodaySent:        todaySent,
		ActiveCodes:      activeCodes,
		AvgSendTime:      2.4, // 模拟平均发送时间
		TodayIncrease:    12,  // 模拟增长率
		SentIncrease:     5,   // 模拟发送增长率
		ActivePercentage: activePercentage,
		SpeedImprovement: 20, // 模拟速度提升
	}

	c.JSON(http.StatusOK, result.DataResult("验证码统计", stats))
}

// GenerateTestData 生成测试验证码
func (h *CaptchaHandler) GenerateTestData(c *gin.Context) {
	// 检查管理员权限
	isAdmin := c.GetBool("is_admin")
	if !isAdmin {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("需要管理员权限"))
		return
	}

	// 生成测试邮箱列表
	testEmails := []string{
		"test1@example.com",
		"test2@domain.org",
		"user@company.io",
		"admin@business.net",
		"support@service.com",
	}

	purposes := []string{"login", "register", "reset_password"}

	// 为每个邮箱生成验证码
	for _, email := range testEmails {
		codeCounter++
		code := &CaptchaCode{
			ID:            codeCounter,
			Email:         email,
			Code:          h.generateRandomCode(),
			SentAt:        time.Now().Add(-time.Duration(codeCounter) * time.Minute),
			ExpiryMinutes: 10,
			IsUsed:        false,
			Purpose:       purposes[int(codeCounter)%len(purposes)],
		}
		captchaCodes[codeCounter] = code
	}

	c.JSON(http.StatusOK, result.DataResult("测试数据生成成功", gin.H{
		"generated": len(testEmails),
	}))
}

// DeleteCode 删除验证码
func (h *CaptchaHandler) DeleteCode(c *gin.Context) {
	// 检查管理员权限
	isAdmin := c.GetBool("is_admin")
	if !isAdmin {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("需要管理员权限"))
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的验证码ID"))
		return
	}

	if _, exists := captchaCodes[id]; !exists {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("验证码不存在"))
		return
	}

	delete(captchaCodes, id)
	c.JSON(http.StatusOK, result.SimpleResult("验证码删除成功"))
}

// ClearExpired 清理过期验证码
func (h *CaptchaHandler) ClearExpired(c *gin.Context) {
	// 检查管理员权限
	isAdmin := c.GetBool("is_admin")
	if !isAdmin {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("需要管理员权限"))
		return
	}

	deletedCount := 0
	for id, code := range captchaCodes {
		if h.isCodeExpired(code) {
			delete(captchaCodes, id)
			deletedCount++
		}
	}

	c.JSON(http.StatusOK, result.DataResult("过期验证码清理完成", gin.H{
		"deletedCount": deletedCount,
	}))
}

// ResendCode 重新发送验证码
func (h *CaptchaHandler) ResendCode(c *gin.Context) {
	// 检查管理员权限
	isAdmin := c.GetBool("is_admin")
	if !isAdmin {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("需要管理员权限"))
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的验证码ID"))
		return
	}

	code, exists := captchaCodes[id]
	if !exists {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("验证码不存在"))
		return
	}

	// 生成新的验证码
	code.Code = h.generateRandomCode()
	code.SentAt = time.Now()
	code.IsUsed = false
	code.UsedAt = nil

	// 这里应该实际发送邮件
	// 暂时只更新记录

	c.JSON(http.StatusOK, result.SimpleResult("验证码重新发送成功"))
}

// SendCaptcha 发送验证码（供其他模块调用）
func (h *CaptchaHandler) SendCaptcha(c *gin.Context) {
	var req struct {
		Email   string `json:"email" binding:"required,email"`
		Purpose string `json:"purpose" binding:"required"` // login, register, reset_password
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("参数错误"))
		return
	}

	// 检查是否已有未过期的验证码
	for _, code := range captchaCodes {
		if code.Email == req.Email && !h.isCodeExpired(code) && !code.IsUsed {
			c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("验证码尚未过期，请稍后再试"))
			return
		}
	}

	// 生成新验证码
	codeCounter++
	newCode := &CaptchaCode{
		ID:            codeCounter,
		Email:         req.Email,
		Code:          h.generateRandomCode(),
		SentAt:        time.Now(),
		ExpiryMinutes: 10,
		IsUsed:        false,
		Purpose:       req.Purpose,
	}

	captchaCodes[codeCounter] = newCode

	// 这里应该实际发送邮件
	// 暂时只返回成功响应

	c.JSON(http.StatusOK, result.DataResult("验证码发送成功", gin.H{
		"email":     req.Email,
		"expiresIn": "10分钟",
	}))
}

// VerifyCaptcha 验证验证码（供其他模块调用）
func (h *CaptchaHandler) VerifyCaptcha(c *gin.Context) {
	var req struct {
		Email   string `json:"email" binding:"required,email"`
		Code    string `json:"code" binding:"required"`
		Purpose string `json:"purpose" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("参数错误"))
		return
	}

	// 查找匹配的验证码
	for _, code := range captchaCodes {
		if code.Email == req.Email &&
			strings.EqualFold(code.Code, req.Code) &&
			code.Purpose == req.Purpose &&
			!code.IsUsed &&
			!h.isCodeExpired(code) {

			// 标记为已使用
			code.IsUsed = true
			now := time.Now()
			code.UsedAt = &now

			c.JSON(http.StatusOK, result.SimpleResult("验证码验证成功"))
			return
		}
	}

	c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("验证码无效或已过期"))
}

// 生成随机验证码
func (h *CaptchaHandler) generateRandomCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 6

	result := make([]byte, length)
	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[num.Int64()]
	}

	return string(result)
}

// 检查验证码是否过期
func (h *CaptchaHandler) isCodeExpired(code *CaptchaCode) bool {
	expiryTime := code.SentAt.Add(time.Duration(code.ExpiryMinutes) * time.Minute)
	return time.Now().After(expiryTime)
}

// GetCaptchaByEmailAndPurpose 根据邮箱和用途获取验证码（内部使用）
func (h *CaptchaHandler) GetCaptchaByEmailAndPurpose(email, purpose string) *CaptchaCode {
	for _, code := range captchaCodes {
		if code.Email == email && code.Purpose == purpose && !code.IsUsed && !h.isCodeExpired(code) {
			return code
		}
	}
	return nil
}

// 初始化一些测试数据
func init() {
	// 添加一些初始测试数据
	testCodes := []*CaptchaCode{
		{
			ID:            1,
			Email:         "admin@company.com",
			Code:          "A3F8K9",
			SentAt:        time.Now().Add(-5 * time.Minute),
			ExpiryMinutes: 10,
			IsUsed:        false,
			Purpose:       "login",
		},
		{
			ID:            2,
			Email:         "user123@example.com",
			Code:          "B5R2T7",
			SentAt:        time.Now().Add(-15 * time.Minute),
			ExpiryMinutes: 10,
			IsUsed:        false,
			Purpose:       "register",
		},
		{
			ID:            3,
			Email:         "support@domain.org",
			Code:          "C9X1Y4",
			SentAt:        time.Now().Add(-2 * time.Minute),
			ExpiryMinutes: 10,
			IsUsed:        false,
			Purpose:       "reset_password",
		},
	}

	codeCounter = 3
	for _, code := range testCodes {
		captchaCodes[code.ID] = code
	}
}

// 全局规则存储
var (
	captchaRules = make(map[int64]*CaptchaRule)
	ruleCounter  int64
)

// 初始化默认规则
func init() {
	defaultRules := []*CaptchaRule{
		{
			ID:                1,
			Name:              "通用6位数字验证码",
			Type:              "regex",
			SenderPattern:     "",
			SubjectPattern:    "*验证码*",
			ExtractionPattern: "(\\d{6})",
			Priority:          "high",
			Enabled:           true,
			Description:       "提取6位数字验证码",
			MatchCount:        1250,
			SuccessCount:      1180,
			LastUsed:          time.Now().Add(-2 * time.Hour),
			CreatedAt:         time.Now().Add(-30 * 24 * time.Hour),
			UpdatedAt:         time.Now().Add(-1 * time.Hour),
		},
		{
			ID:                2,
			Name:              "GitHub验证码",
			Type:              "regex",
			SenderPattern:     "@github.com",
			SubjectPattern:    "*verification*",
			ExtractionPattern: "verification code is (\\d{6})",
			Priority:          "high",
			Enabled:           true,
			Description:       "GitHub邮箱验证码提取",
			MatchCount:        89,
			SuccessCount:      87,
			LastUsed:          time.Now().Add(-3 * time.Hour),
			CreatedAt:         time.Now().Add(-20 * 24 * time.Hour),
			UpdatedAt:         time.Now().Add(-2 * time.Hour),
		},
		{
			ID:                3,
			Name:              "阿里云验证码",
			Type:              "keyword",
			SenderPattern:     "@aliyun.com",
			SubjectPattern:    "",
			ExtractionPattern: "验证码：",
			Priority:          "medium",
			Enabled:           true,
			Description:       "阿里云服务验证码",
			MatchCount:        156,
			SuccessCount:      152,
			LastUsed:          time.Now().Add(-1 * 24 * time.Hour),
			CreatedAt:         time.Now().Add(-15 * 24 * time.Hour),
			UpdatedAt:         time.Now().Add(-12 * time.Hour),
		},
		{
			ID:                4,
			Name:              "腾讯云验证码",
			Type:              "position",
			SenderPattern:     "@tencent.com",
			SubjectPattern:    "*验证码*",
			ExtractionPattern: "line:3,word:2",
			Priority:          "medium",
			Enabled:           false,
			Description:       "腾讯云验证码（已禁用）",
			MatchCount:        45,
			SuccessCount:      40,
			LastUsed:          time.Now().Add(-2 * 24 * time.Hour),
			CreatedAt:         time.Now().Add(-10 * 24 * time.Hour),
			UpdatedAt:         time.Now().Add(-1 * 24 * time.Hour),
		},
	}

	ruleCounter = 4
	for _, rule := range defaultRules {
		captchaRules[rule.ID] = rule
	}
}

// GetCaptchaRules 获取验证码提取规则列表
func (h *CaptchaHandler) GetCaptchaRules(c *gin.Context) {
	// 检查管理员权限
	isAdmin := c.GetBool("is_admin")
	if !isAdmin {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("需要管理员权限"))
		return
	}

	rules := make([]*CaptchaRule, 0, len(captchaRules))
	for _, rule := range captchaRules {
		rules = append(rules, rule)
	}

	c.JSON(http.StatusOK, result.DataResult("获取规则列表成功", rules))
}

// UpdateCaptchaRuleStatus 更新验证码规则状态
func (h *CaptchaHandler) UpdateCaptchaRuleStatus(c *gin.Context) {
	// 检查管理员权限
	isAdmin := c.GetBool("is_admin")
	if !isAdmin {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("需要管理员权限"))
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("无效的规则ID"))
		return
	}

	var req struct {
		Enabled bool `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("参数错误"))
		return
	}

	rule, exists := captchaRules[id]
	if !exists {
		c.JSON(http.StatusNotFound, result.ErrorSimpleResult("规则不存在"))
		return
	}

	rule.Enabled = req.Enabled
	rule.UpdatedAt = time.Now()

	c.JSON(http.StatusOK, result.SimpleResult("规则状态更新成功"))
}
