package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"miko-email/internal/result"
	"miko-email/internal/svc"
)

type SystemHandler struct {
	svcCtx *svc.ServiceContext
}

func NewSystemHandler(svcCtx *svc.ServiceContext) *SystemHandler {
	return &SystemHandler{
		svcCtx: svcCtx,
	}
}

// SystemSettings 系统设置结构
type SystemSettings struct {
	General struct {
		SystemName         string `json:"systemName"`
		SystemDescription  string `json:"systemDescription"`
		DefaultTheme       string `json:"defaultTheme"`
		Timezone           string `json:"timezone"`
		Language           string `json:"language"`
		MaxUsers           int    `json:"maxUsers"`
		EnableRegistration bool   `json:"enableRegistration"`
	} `json:"general"`
	Security struct {
		SessionTimeout    int  `json:"sessionTimeout"`
		MinPasswordLength int  `json:"minPasswordLength"`
		EnableTwoFactor   bool `json:"enableTwoFactor"`
		EnableLoginLock   bool `json:"enableLoginLock"`
		MaxLoginAttempts  int  `json:"maxLoginAttempts"`
		LockoutDuration   int  `json:"lockoutDuration"`
		EnableSSL         bool `json:"enableSSL"`
		EnableFirewall    bool `json:"enableFirewall"`
	} `json:"security"`
	Email struct {
		SMTPHost      string `json:"smtpHost"`
		SMTPPort      int    `json:"smtpPort"`
		SMTPUser      string `json:"smtpUser"`
		SMTPPassword  string `json:"smtpPassword"`
		SenderEmail   string `json:"senderEmail"`
		SenderName    string `json:"senderName"`
		EnableTLS     bool   `json:"enableTLS"`
		EmailSignature string `json:"emailSignature"`
	} `json:"email"`
	Storage struct {
		MaxFileSize        int64 `json:"maxFileSize"`
		AllowedExtensions  []string `json:"allowedExtensions"`
		StoragePath        string `json:"storagePath"`
		EnableAutoCleanup  bool   `json:"enableAutoCleanup"`
		BackupEnabled      bool   `json:"backupEnabled"`
		BackupInterval     int    `json:"backupInterval"`
		RetentionDays      int    `json:"retentionDays"`
	} `json:"storage"`
}

// SystemStatus 系统状态结构
type SystemStatus struct {
	Services []ServiceStatus `json:"services"`
	System   SystemInfo      `json:"system"`
	Database DatabaseInfo    `json:"database"`
}

type ServiceStatus struct {
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	Uptime      string    `json:"uptime"`
	LastCheck   time.Time `json:"lastCheck"`
	Description string    `json:"description"`
}

type SystemInfo struct {
	CPUUsage    float64 `json:"cpuUsage"`
	MemoryUsage float64 `json:"memoryUsage"`
	DiskUsage   float64 `json:"diskUsage"`
	Uptime      string  `json:"uptime"`
	Version     string  `json:"version"`
}

type DatabaseInfo struct {
	Status      string `json:"status"`
	Connections int    `json:"connections"`
	Size        string `json:"size"`
	LastBackup  string `json:"lastBackup"`
}

// SystemLog 系统日志结构
type SystemLog struct {
	ID        int64     `json:"id"`
	Level     string    `json:"level"`
	Module    string    `json:"module"`
	Message   string    `json:"message"`
	Details   string    `json:"details"`
	Timestamp time.Time `json:"timestamp"`
	UserID    int64     `json:"userId,omitempty"`
	IP        string    `json:"ip,omitempty"`
}

// 全局系统设置存储
var globalSystemSettings = SystemSettings{
	General: struct {
		SystemName         string `json:"systemName"`
		SystemDescription  string `json:"systemDescription"`
		DefaultTheme       string `json:"defaultTheme"`
		Timezone           string `json:"timezone"`
		Language           string `json:"language"`
		MaxUsers           int    `json:"maxUsers"`
		EnableRegistration bool   `json:"enableRegistration"`
	}{
		SystemName:         "Miko邮箱系统",
		SystemDescription:  "企业级邮箱管理系统",
		DefaultTheme:       "dark",
		Timezone:           "Asia/Shanghai",
		Language:           "zh-CN",
		MaxUsers:           1000,
		EnableRegistration: true,
	},
	Security: struct {
		SessionTimeout    int  `json:"sessionTimeout"`
		MinPasswordLength int  `json:"minPasswordLength"`
		EnableTwoFactor   bool `json:"enableTwoFactor"`
		EnableLoginLock   bool `json:"enableLoginLock"`
		MaxLoginAttempts  int  `json:"maxLoginAttempts"`
		LockoutDuration   int  `json:"lockoutDuration"`
		EnableSSL         bool `json:"enableSSL"`
		EnableFirewall    bool `json:"enableFirewall"`
	}{
		SessionTimeout:    30,
		MinPasswordLength: 8,
		EnableTwoFactor:   false,
		EnableLoginLock:   true,
		MaxLoginAttempts:  5,
		LockoutDuration:   15,
		EnableSSL:         true,
		EnableFirewall:    true,
	},
	Email: struct {
		SMTPHost      string `json:"smtpHost"`
		SMTPPort      int    `json:"smtpPort"`
		SMTPUser      string `json:"smtpUser"`
		SMTPPassword  string `json:"smtpPassword"`
		SenderEmail   string `json:"senderEmail"`
		SenderName    string `json:"senderName"`
		EnableTLS     bool   `json:"enableTLS"`
		EmailSignature string `json:"emailSignature"`
	}{
		SMTPHost:      "smtp.example.com",
		SMTPPort:      587,
		SMTPUser:      "noreply@example.com",
		SMTPPassword:  "password",
		SenderEmail:   "noreply@example.com",
		SenderName:    "Miko邮箱系统",
		EnableTLS:     true,
		EmailSignature: "此邮件由Miko邮箱系统自动发送，请勿回复。",
	},
	Storage: struct {
		MaxFileSize        int64 `json:"maxFileSize"`
		AllowedExtensions  []string `json:"allowedExtensions"`
		StoragePath        string `json:"storagePath"`
		EnableAutoCleanup  bool   `json:"enableAutoCleanup"`
		BackupEnabled      bool   `json:"backupEnabled"`
		BackupInterval     int    `json:"backupInterval"`
		RetentionDays      int    `json:"retentionDays"`
	}{
		MaxFileSize:       10 * 1024 * 1024, // 10MB
		AllowedExtensions: []string{".jpg", ".png", ".pdf", ".doc", ".docx"},
		StoragePath:       "./storage",
		EnableAutoCleanup: true,
		BackupEnabled:     true,
		BackupInterval:    24, // 24小时
		RetentionDays:     30,
	},
}

// GetSystemSettings 获取系统设置
func (h *SystemHandler) GetSystemSettings(c *gin.Context) {
	// 检查管理员权限
	isAdmin := c.GetBool("is_admin")
	if !isAdmin {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("需要管理员权限"))
		return
	}

	c.JSON(http.StatusOK, result.DataResult("获取系统设置成功", globalSystemSettings))
}

// UpdateSystemSettings 更新系统设置
func (h *SystemHandler) UpdateSystemSettings(c *gin.Context) {
	// 检查管理员权限
	isAdmin := c.GetBool("is_admin")
	if !isAdmin {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("需要管理员权限"))
		return
	}

	var settings SystemSettings
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, result.ErrorSimpleResult("参数错误"))
		return
	}

	// 更新全局设置
	globalSystemSettings = settings

	c.JSON(http.StatusOK, result.SimpleResult("系统设置更新成功"))
}

// GetSystemStatus 获取系统状态
func (h *SystemHandler) GetSystemStatus(c *gin.Context) {
	// 检查管理员权限
	isAdmin := c.GetBool("is_admin")
	if !isAdmin {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("需要管理员权限"))
		return
	}

	status := SystemStatus{
		Services: []ServiceStatus{
			{
				Name:        "邮件服务",
				Status:      "running",
				Uptime:      "15天 8小时",
				LastCheck:   time.Now(),
				Description: "SMTP/IMAP邮件服务",
			},
			{
				Name:        "数据库服务",
				Status:      "running",
				Uptime:      "30天 12小时",
				LastCheck:   time.Now(),
				Description: "MySQL数据库服务",
			},
			{
				Name:        "缓存服务",
				Status:      "running",
				Uptime:      "7天 3小时",
				LastCheck:   time.Now(),
				Description: "Redis缓存服务",
			},
		},
		System: SystemInfo{
			CPUUsage:    25.6,
			MemoryUsage: 68.2,
			DiskUsage:   45.8,
			Uptime:      "30天 12小时",
			Version:     "v1.0.0",
		},
		Database: DatabaseInfo{
			Status:      "healthy",
			Connections: 15,
			Size:        "256MB",
			LastBackup:  "2025-01-27 02:00:00",
		},
	}

	c.JSON(http.StatusOK, result.DataResult("获取系统状态成功", status))
}

// GetSystemLogs 获取系统日志
func (h *SystemHandler) GetSystemLogs(c *gin.Context) {
	// 检查管理员权限
	isAdmin := c.GetBool("is_admin")
	if !isAdmin {
		c.JSON(http.StatusForbidden, result.ErrorSimpleResult("需要管理员权限"))
		return
	}

	// 模拟日志数据
	logs := []SystemLog{
		{
			ID:        1,
			Level:     "info",
			Module:    "auth",
			Message:   "用户登录成功",
			Details:   "用户 admin 从 192.168.1.100 登录",
			Timestamp: time.Now().Add(-1 * time.Hour),
			UserID:    1,
			IP:        "192.168.1.100",
		},
		{
			ID:        2,
			Level:     "warning",
			Module:    "email",
			Message:   "邮件发送延迟",
			Details:   "SMTP服务器响应时间超过5秒",
			Timestamp: time.Now().Add(-2 * time.Hour),
			IP:        "127.0.0.1",
		},
		{
			ID:        3,
			Level:     "error",
			Module:    "database",
			Message:   "数据库连接失败",
			Details:   "连接超时，已自动重连",
			Timestamp: time.Now().Add(-3 * time.Hour),
			IP:        "127.0.0.1",
		},
	}

	c.JSON(http.StatusOK, result.DataResult("获取系统日志成功", logs))
}
