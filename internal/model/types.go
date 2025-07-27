package model

import "time"

// UserReq 用户查询参数
type UserReq struct {
	Id           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	IsActive     *bool     `json:"is_active"`
	Contribution int       `json:"contribution"`
	InviteCode   string    `json:"invite_code"`
	InvitedBy    *int64    `json:"invited_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Page         int       `json:"page"`
	PageSize     int       `json:"page_size"`
}

// AdminReq 管理员查询参数
type AdminReq struct {
	Id           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	IsActive     *bool     `json:"is_active"`
	Contribution int       `json:"contribution"`
	InviteCode   string    `json:"invite_code"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Page         int       `json:"page"`
	PageSize     int       `json:"page_size"`
}

// DomainReq 域名查询参数
type DomainReq struct {
	Id         int64     `json:"id"`
	Name       string    `json:"name"`
	IsVerified *bool     `json:"is_verified"`
	IsActive   *bool     `json:"is_active"`
	MxRecord   string    `json:"mx_record"`
	ARecord    string    `json:"a_record"`
	TxtRecord  string    `json:"txt_record"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Page       int       `json:"page"`
	PageSize   int       `json:"page_size"`
}

// MailboxReq 邮箱查询参数
type MailboxReq struct {
	Id        int64     `json:"id"`
	UserId    *int64    `json:"user_id"`
	AdminId   *int64    `json:"admin_id"`
	Email     string    `json:"email"`
	DomainId  int64     `json:"domain_id"`
	IsActive  *bool     `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Page      int       `json:"page"`
	PageSize  int       `json:"page_size"`
}

// EmailReq 邮件查询参数
type EmailReq struct {
	Id        int64     `json:"id"`
	MailboxId int64     `json:"mailbox_id"`
	FromAddr  string    `json:"from_addr"`
	ToAddr    string    `json:"to_addr"`
	Subject   string    `json:"subject"`
	Body      string    `json:"body"`
	IsRead    *bool     `json:"is_read"`
	Folder    string    `json:"folder"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Page      int       `json:"page"`
	PageSize  int       `json:"page_size"`
}

// EmailForwardReq 邮件转发查询参数
type EmailForwardReq struct {
	Id                 int64      `json:"id"`
	MailboxId          int64      `json:"mailbox_id"`
	SourceEmail        string     `json:"source_email"`
	TargetEmail        string     `json:"target_email"`
	Enabled            *bool      `json:"enabled"`
	KeepOriginal       *bool      `json:"keep_original"`
	ForwardAttachments *bool      `json:"forward_attachments"`
	SubjectPrefix      string     `json:"subject_prefix"`
	Description        string     `json:"description"`
	ForwardCount       int64      `json:"forward_count"`
	LastForwardAt      *time.Time `json:"last_forward_at"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	Page               int        `json:"page"`
	PageSize           int        `json:"page_size"`
}
