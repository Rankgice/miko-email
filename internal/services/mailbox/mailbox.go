package mailbox

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"miko-email/internal/model"
	"miko-email/internal/svc"
)

type Service struct {
	svcCtx *svc.ServiceContext
}

// MailboxResponse 邮箱响应结构体
type MailboxResponse struct {
	ID        int64     `json:"id"`
	UserID    *int64    `json:"user_id,omitempty"`
	AdminID   *int64    `json:"admin_id,omitempty"`
	Email     string    `json:"email"`
	DomainID  int64     `json:"domain_id"`
	Status    string    `json:"status"` // 转换后的状态字段
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewService(svcCtx *svc.ServiceContext) *Service {
	return &Service{svcCtx: svcCtx}
}

// GetUserMailboxes 获取用户的邮箱列表
func (s *Service) GetUserMailboxes(userID int64, isAdmin bool) ([]MailboxResponse, error) {
	var params model.MailboxReq
	isActive := true
	params.IsActive = &isActive

	if isAdmin {
		params.AdminId = &userID
	} else {
		params.UserId = &userID
	}

	mailboxes, _, err := s.svcCtx.MailboxModel.List(params)
	if err != nil {
		return nil, err
	}

	// 初始化为空数组而不是nil，确保JSON序列化时返回[]而不是null
	responses := make([]MailboxResponse, 0)
	for _, mailbox := range mailboxes {
		response := MailboxResponse{
			ID:        mailbox.Id,
			UserID:    mailbox.UserId,
			AdminID:   mailbox.AdminId,
			Email:     mailbox.Email,
			DomainID:  mailbox.DomainId,
			CreatedAt: mailbox.CreatedAt,
			UpdatedAt: mailbox.UpdatedAt,
		}

		// 转换状态字段
		if mailbox.IsActive {
			response.Status = "active"
		} else {
			response.Status = "deleted"
		}

		responses = append(responses, response)
	}

	return responses, nil
}

// GetUserMailboxesRaw 获取用户的邮箱列表（返回原始model.Mailbox类型）
func (s *Service) GetUserMailboxesRaw(userID int64, isAdmin bool) ([]*model.Mailbox, error) {
	var params model.MailboxReq
	isActive := true
	params.IsActive = &isActive

	if isAdmin {
		params.AdminId = &userID
	} else {
		params.UserId = &userID
	}

	mailboxes, _, err := s.svcCtx.MailboxModel.List(params)
	if err != nil {
		return nil, err
	}

	return mailboxes, nil
}

// CreateMailbox 创建邮箱
func (s *Service) CreateMailbox(userID int64, prefix string, domainID int64, isAdmin bool) (*model.Mailbox, error) {
	// 获取域名
	domain, err := s.svcCtx.DomainModel.GetById(domainID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("域名不存在或已禁用")
		}
		return nil, err
	}

	if !domain.IsActive {
		return nil, fmt.Errorf("域名不存在或已禁用")
	}

	fullEmail := fmt.Sprintf("%s@%s", prefix, domain.Name)

	// 检查邮箱是否已存在
	exists, err := s.svcCtx.MailboxModel.CheckEmailExist(fullEmail)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("邮箱已存在")
	}

	// 生成邮箱密码
	password := uuid.New().String()[:8]

	// 创建邮箱
	mailbox := &model.Mailbox{
		Email:     fullEmail,
		Password:  password,
		DomainId:  domainID,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if isAdmin {
		mailbox.AdminId = &userID
	} else {
		mailbox.UserId = &userID
	}

	if err := s.svcCtx.MailboxModel.Create(nil, mailbox); err != nil {
		return nil, err
	}

	return mailbox, nil
}

// CreateMailboxWithPassword 创建邮箱（使用自定义密码）
func (s *Service) CreateMailboxWithPassword(userID int64, prefix string, password string, domainID int64, isAdmin bool) (*model.Mailbox, error) {
	// 获取域名
	domain, err := s.svcCtx.DomainModel.GetById(domainID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("域名不存在或已禁用")
		}
		return nil, err
	}

	if !domain.IsActive {
		return nil, fmt.Errorf("域名不存在或已禁用")
	}

	fullEmail := fmt.Sprintf("%s@%s", prefix, domain.Name)

	// 检查邮箱是否已存在
	exists, err := s.svcCtx.MailboxModel.CheckEmailExist(fullEmail)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("邮箱已存在")
	}

	// 创建邮箱
	mailbox := &model.Mailbox{
		Email:     fullEmail,
		Password:  password,
		DomainId:  domainID,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if isAdmin {
		mailbox.AdminId = &userID
	} else {
		mailbox.UserId = &userID
	}

	if err := s.svcCtx.MailboxModel.Create(nil, mailbox); err != nil {
		return nil, err
	}

	return mailbox, nil
}

// BatchCreateMailboxes 批量创建邮箱
func (s *Service) BatchCreateMailboxes(userID int64, prefixes []string, domainID int64, isAdmin bool) ([]*model.Mailbox, error) {
	// 获取域名
	domain, err := s.svcCtx.DomainModel.GetById(domainID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("域名不存在或已禁用")
		}
		return nil, err
	}

	if !domain.IsActive {
		return nil, fmt.Errorf("域名不存在或已禁用")
	}

	// 开始事务
	tx := s.svcCtx.DB.Begin()
	defer func() {
		if tx != nil {
			tx.Rollback()
		}
	}()

	var mailboxes []*model.Mailbox
	for _, prefix := range prefixes {
		fullEmail := fmt.Sprintf("%s@%s", prefix, domain.Name)

		// 检查邮箱是否已存在
		exists, err := s.svcCtx.MailboxModel.CheckEmailExist(fullEmail)
		if err != nil {
			return nil, err
		}
		if exists {
			continue // 跳过已存在的邮箱
		}

		// 生成邮箱密码
		password := uuid.New().String()[:8]

		// 创建邮箱
		mailbox := &model.Mailbox{
			Email:     fullEmail,
			Password:  password,
			DomainId:  domainID,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if isAdmin {
			mailbox.AdminId = &userID
		} else {
			mailbox.UserId = &userID
		}

		if err := s.svcCtx.MailboxModel.Create(tx, mailbox); err != nil {
			return nil, err
		}

		mailboxes = append(mailboxes, mailbox)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	tx = nil

	return mailboxes, nil
}

// GetMailboxByEmail 根据邮箱地址获取邮箱信息
func (s *Service) GetMailboxByEmail(email string) (*model.Mailbox, error) {
	mailbox, err := s.svcCtx.MailboxModel.GetByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("邮箱不存在")
		}
		return nil, err
	}

	if !mailbox.IsActive {
		return nil, fmt.Errorf("邮箱不存在")
	}

	return mailbox, nil
}

// GetMailboxByID 根据ID获取邮箱信息
func (s *Service) GetMailboxByID(mailboxID int64) (*model.Mailbox, error) {
	mailbox, err := s.svcCtx.MailboxModel.GetById(mailboxID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("邮箱不存在")
		}
		return nil, err
	}

	if !mailbox.IsActive {
		return nil, fmt.Errorf("邮箱不存在")
	}

	return mailbox, nil
}

// GetMailboxPassword 获取邮箱密码
func (s *Service) GetMailboxPassword(mailboxID int64, userID int64, isAdmin bool) (string, error) {
	// 获取邮箱信息
	mailbox, err := s.svcCtx.MailboxModel.GetById(mailboxID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("邮箱不存在")
		}
		return "", err
	}

	if !mailbox.IsActive {
		return "", fmt.Errorf("邮箱不存在")
	}

	// 检查权限
	if isAdmin {
		// 管理员需要检查是否是自己创建的邮箱
		if mailbox.AdminId == nil || *mailbox.AdminId != userID {
			return "", fmt.Errorf("无权限访问此邮箱")
		}
	} else {
		// 普通用户需要检查是否是自己的邮箱
		if mailbox.UserId == nil || *mailbox.UserId != userID {
			return "", fmt.Errorf("无权限访问此邮箱")
		}
	}

	return mailbox.Password, nil
}

// DeleteMailbox 删除邮箱
func (s *Service) DeleteMailbox(mailboxID int64, userID int64, isAdmin bool) error {
	// 获取邮箱信息
	mailbox, err := s.svcCtx.MailboxModel.GetById(mailboxID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("邮箱不存在")
		}
		return err
	}

	// 检查权限
	if isAdmin {
		// 管理员需要检查是否是自己创建的邮箱
		if mailbox.AdminId == nil || *mailbox.AdminId != userID {
			return fmt.Errorf("无权限删除此邮箱")
		}
	} else {
		// 普通用户需要检查是否是自己的邮箱
		if mailbox.UserId == nil || *mailbox.UserId != userID {
			return fmt.Errorf("无权限删除此邮箱")
		}
	}

	// 开始事务
	tx := s.svcCtx.DB.Begin()
	defer func() {
		if tx != nil {
			tx.Rollback()
		}
	}()

	// 删除相关邮件
	if err := tx.Where("mailbox_id = ?", mailboxID).Delete(&model.Email{}).Error; err != nil {
		return err
	}

	// 删除转发规则
	if err := tx.Where("mailbox_id = ?", mailboxID).Delete(&model.EmailForward{}).Error; err != nil {
		return err
	}

	// 删除邮箱
	if err := s.svcCtx.MailboxModel.Delete(tx, mailbox); err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}
	tx = nil

	return nil
}

// 管理员邮箱管理方法

// AdminMailboxResponse 管理员邮箱响应结构体
type AdminMailboxResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	UserID    *int64    `json:"user_id,omitempty"`
	AdminID   *int64    `json:"admin_id,omitempty"`
	DomainID  int64     `json:"domain_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// MailboxStats 邮箱统计信息
type MailboxStats struct {
	InboxCount   int        `json:"inbox_count"`
	SentCount    int        `json:"sent_count"`
	LastActivity *time.Time `json:"last_activity,omitempty"`
}

// GetAllMailboxes 获取所有邮箱列表（管理员）
func (s *Service) GetAllMailboxes() ([]AdminMailboxResponse, error) {
	query := `
		SELECT m.id, m.email, m.user_id, m.admin_id, m.domain_id, m.is_active, m.created_at, m.updated_at,
		       COALESCE(u.username, a.username, '未知用户') as username
		FROM mailboxes m
		LEFT JOIN users u ON m.user_id = u.id
		LEFT JOIN admins a ON m.admin_id = a.id
		ORDER BY m.created_at DESC
	`
	db, err := s.svcCtx.DB.DB()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 初始化为空数组而不是nil，确保JSON序列化时返回[]而不是null
	mailboxes := make([]AdminMailboxResponse, 0)
	for rows.Next() {
		var mailbox AdminMailboxResponse
		var isActive bool

		err := rows.Scan(
			&mailbox.ID, &mailbox.Email, &mailbox.UserID, &mailbox.AdminID,
			&mailbox.DomainID, &isActive, &mailbox.CreatedAt, &mailbox.UpdatedAt,
			&mailbox.Username,
		)
		if err != nil {
			continue
		}

		// 转换状态
		if isActive {
			mailbox.Status = "active"
		} else {
			mailbox.Status = "suspended"
		}

		mailboxes = append(mailboxes, mailbox)
	}

	return mailboxes, nil
}

// UpdateMailboxStatus 更新邮箱状态（管理员）
func (s *Service) UpdateMailboxStatus(mailboxID int64, status string) error {
	isActive := status == "active"
	return s.svcCtx.MailboxModel.UpdateStatus(nil, mailboxID, isActive)
}

// DeleteMailboxAdmin 删除邮箱（管理员）
func (s *Service) DeleteMailboxAdmin(mailboxID int64) error {
	// 开始事务
	tx := s.svcCtx.DB.Begin()
	defer func() {
		if tx != nil {
			tx.Rollback()
		}
	}()

	// 删除相关邮件
	if err := tx.Where("mailbox_id = ?", mailboxID).Delete(&model.Email{}).Error; err != nil {
		return err
	}

	// 删除转发规则
	if err := tx.Where("mailbox_id = ?", mailboxID).Delete(&model.EmailForward{}).Error; err != nil {
		return err
	}

	// 删除邮箱
	if err := s.svcCtx.MailboxModel.HardDelete(tx, mailboxID); err != nil {
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}
	tx = nil

	return nil
}

// GetMailboxStats 获取邮箱统计信息（管理员）
func (s *Service) GetMailboxStats(mailboxID int64) (*MailboxStats, error) {
	stats := &MailboxStats{}

	// 获取收件数量
	var inboxCount int64
	if err := s.svcCtx.DB.Model(&model.Email{}).
		Where("mailbox_id = ? AND folder = ?", mailboxID, "inbox").
		Count(&inboxCount).Error; err != nil {
		return nil, err
	}
	stats.InboxCount = int(inboxCount)

	// 获取发件数量
	var sentCount int64
	if err := s.svcCtx.DB.Model(&model.Email{}).
		Where("mailbox_id = ? AND folder = ?", mailboxID, "sent").
		Count(&sentCount).Error; err != nil {
		return nil, err
	}
	stats.SentCount = int(sentCount)

	// 获取最后活动时间
	var lastEmail model.Email
	if err := s.svcCtx.DB.Where("mailbox_id = ?", mailboxID).
		Order("created_at DESC").
		First(&lastEmail).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		// 没有邮件记录，lastActivity 保持为 nil
	} else {
		stats.LastActivity = &lastEmail.CreatedAt
	}

	return stats, nil
}

// UserStats 用户统计信息
type UserStats struct {
	TotalMailboxes int `json:"total_mailboxes"`
	UnreadEmails   int `json:"unread_emails"`
	SentEmails     int `json:"sent_emails"`
	TotalEmails    int `json:"total_emails"`
}

// GetUserStats 获取用户统计信息
func (s *Service) GetUserStats(userID int64) (*UserStats, error) {
	stats := &UserStats{}

	// 获取用户的邮箱数量
	var totalMailboxes int64
	if err := s.svcCtx.DB.Model(&model.Mailbox{}).
		Where("user_id = ? AND is_active = ?", userID, true).
		Count(&totalMailboxes).Error; err != nil {
		return nil, err
	}
	stats.TotalMailboxes = int(totalMailboxes)

	// 获取未读邮件数量（用户所有邮箱的未读邮件）
	var unreadEmails int64
	if err := s.svcCtx.DB.Raw(`
		SELECT COUNT(*) FROM email e
		JOIN mailbox m ON e.mailbox_id = m.id
		WHERE m.user_id = ? AND m.is_active = 1 AND e.folder = 'inbox' AND e.is_read = 0
	`, userID).Scan(&unreadEmails).Error; err != nil {
		return nil, err
	}
	stats.UnreadEmails = int(unreadEmails)

	// 获取已发送邮件数量
	var sentEmails int64
	if err := s.svcCtx.DB.Raw(`
		SELECT COUNT(*) FROM email e
		JOIN mailbox m ON e.mailbox_id = m.id
		WHERE m.user_id = ? AND m.is_active = 1 AND e.folder = 'sent'
	`, userID).Scan(&sentEmails).Error; err != nil {
		return nil, err
	}
	stats.SentEmails = int(sentEmails)

	// 获取总邮件数量
	var totalEmails int64
	if err := s.svcCtx.DB.Raw(`
		SELECT COUNT(*) FROM email e
		JOIN mailbox m ON e.mailbox_id = m.id
		WHERE m.user_id = ? AND m.is_active = 1
	`, userID).Scan(&totalEmails).Error; err != nil {
		return nil, err
	}
	stats.TotalEmails = int(totalEmails)

	return stats, nil
}
