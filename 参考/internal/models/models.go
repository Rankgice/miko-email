package models

import (
	"time"
)

// User 用户模型
type User struct {
	ID        int       `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"`
	Name      string    `json:"name" db:"name"`
	IsAdmin   bool      `json:"is_admin" db:"is_admin"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Email 邮件模型
type Email struct {
	ID          int       `json:"id" db:"id"`
	MessageID   string    `json:"message_id" db:"message_id"`
	From        string    `json:"from" db:"from_addr"`
	To          string    `json:"to" db:"to_addr"`
	Subject     string    `json:"subject" db:"subject"`
	Body        string    `json:"body" db:"body"`
	HTMLBody    string    `json:"html_body" db:"html_body"`
	IsRead      bool      `json:"is_read" db:"is_read"`
	IsDeleted   bool      `json:"is_deleted" db:"is_deleted"`
	IsSent      bool      `json:"is_sent" db:"is_sent"`
	UserID      int       `json:"user_id" db:"user_id"`
	Size        int       `json:"size" db:"size"`
	Attachments string    `json:"attachments" db:"attachments"` // JSON格式存储附件信息
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Domain 域名模型
type Domain struct {
	ID           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Description  string    `json:"description" db:"description"` // 域名描述
	UserID       *int      `json:"user_id" db:"user_id"`         // 所属用户ID，可为空表示公共域名
	IsActive     bool      `json:"is_active" db:"is_active"`
	DNSVerified  bool      `json:"dns_verified" db:"dns_verified"`   // DNS验证状态
	MXRecord     string    `json:"mx_record" db:"mx_record"`         // MX记录
	LastVerified time.Time `json:"last_verified" db:"last_verified"` // 最后验证时间
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Mailbox 邮箱模型
type Mailbox struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Email     string    `json:"email" db:"email"`
	DomainID  int       `json:"domain_id" db:"domain_id"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	IsCurrent bool      `json:"is_current" db:"is_current"` // 是否为当前使用的邮箱
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Attachment 附件模型
type Attachment struct {
	ID       int    `json:"id"`
	Filename string `json:"filename"`
	Size     int    `json:"size"`
	MimeType string `json:"mime_type"`
	Content  []byte `json:"-"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// SendEmailRequest 发送邮件请求
type SendEmailRequest struct {
	To      string `json:"to" binding:"required"`
	Subject string `json:"subject" binding:"required"`
	Body    string `json:"body" binding:"required"`
	From    string `json:"from"` // 可选的发件人邮箱
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password,omitempty"`
}

// GenerateMailboxesRequest 批量生成邮箱请求
type GenerateMailboxesRequest struct {
	DomainID int    `json:"domain_id" binding:"required"`
	Count    int    `json:"count" binding:"required,min=1,max=100"`
	Prefix   string `json:"prefix"` // 可选的前缀
}

// SwitchMailboxRequest 切换邮箱请求
type SwitchMailboxRequest struct {
	MailboxID int `json:"mailbox_id" binding:"required"`
}

// CreateDomainRequest 创建域名请求
type CreateDomainRequest struct {
	Name string `json:"name" binding:"required"`
}

// VerifyDomainRequest 验证域名请求
type VerifyDomainRequest struct {
	DomainID int `json:"domain_id" binding:"required"`
}

// EmailListResponse 邮件列表响应
type EmailListResponse struct {
	Emails []Email `json:"emails"`
	Total  int     `json:"total"`
	Page   int     `json:"page"`
	Limit  int     `json:"limit"`
}

// Response 通用响应
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
