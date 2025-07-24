package model

import (
	"time"

	"gorm.io/gorm"
)

// Mailbox 邮箱模型
type Mailbox struct {
	Id        int64     `gorm:"column:id;primaryKey;autoIncrement;comment:数据库主键ID" json:"id"`               // 数据库主键ID
	UserId    *int64    `gorm:"column:user_id;comment:普通用户ID" json:"user_id,omitempty"`                     // 普通用户ID
	AdminId   *int64    `gorm:"column:admin_id;comment:管理员ID" json:"admin_id,omitempty"`                    // 管理员ID
	Email     string    `gorm:"column:email;uniqueIndex;not null;comment:完整邮箱地址" json:"email"`              // 完整邮箱地址
	Password  string    `gorm:"column:password;not null;comment:邮箱密码" json:"password,omitempty"`            // 邮箱密码
	DomainId  int64     `gorm:"column:domain_id;not null;comment:域名ID" json:"domain_id"`                    // 域名ID
	IsActive  bool      `gorm:"column:is_active;default:1;comment:是否激活" json:"is_active"`                   // 是否激活
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"` // 更新时间
}

// TableName 指定表名
func (Mailbox) TableName() string {
	return "mailboxes"
}

// MailboxModel 邮箱模型
type MailboxModel struct {
	db *gorm.DB
}

// NewMailboxModel 创建邮箱模型
func NewMailboxModel(db *gorm.DB) *MailboxModel {
	return &MailboxModel{
		db: db,
	}
}

// Create 创建邮箱
func (m *MailboxModel) Create(mailbox *Mailbox) error {
	return m.db.Create(mailbox).Error
}

// Update 更新邮箱
func (m *MailboxModel) Update(tx *gorm.DB, mailbox *Mailbox) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Updates(mailbox).Error
}

// MapUpdate 更新邮箱
func (m *MailboxModel) MapUpdate(tx *gorm.DB, id int64, data map[string]interface{}) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Model(&Mailbox{}).Where("id = ?", id).Updates(data).Error
}

// Save 保存邮箱
func (m *MailboxModel) Save(tx *gorm.DB, mailbox *Mailbox) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Save(mailbox).Error
}

// Delete 删除邮箱
func (m *MailboxModel) Delete(mailbox *Mailbox) error {
	return m.db.Delete(mailbox).Error
}

// GetById 根据ID获取邮箱
func (m *MailboxModel) GetById(id int64) (*Mailbox, error) {
	var mailbox Mailbox
	if err := m.db.First(&mailbox, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 如果没有找到记录，返回nil
		}
		return nil, err
	}
	return &mailbox, nil
}

// GetByEmail 根据邮箱地址获取记录
func (m *MailboxModel) GetByEmail(email string) (*Mailbox, error) {
	var mailbox Mailbox
	if err := m.db.Where("email = ?", email).First(&mailbox).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &mailbox, nil
}

// List 获取邮箱列表（统一查询方法）
func (m *MailboxModel) List(params MailboxReq) ([]*Mailbox, int64, error) {
	var mailboxes []*Mailbox
	var total int64

	db := m.db.Model(&Mailbox{})

	// 添加查询条件
	if params.Id != 0 {
		db = db.Where("id = ?", params.Id)
	}
	if params.UserId != nil {
		db = db.Where("user_id = ?", *params.UserId)
	}
	if params.AdminId != nil {
		db = db.Where("admin_id = ?", *params.AdminId)
	}
	if params.Email != "" {
		db = db.Where("email LIKE ?", "%"+params.Email+"%")
	}
	if params.DomainId != 0 {
		db = db.Where("domain_id = ?", params.DomainId)
	}
	if params.IsActive != nil {
		db = db.Where("is_active = ?", *params.IsActive)
	}
	if !params.CreatedAt.IsZero() {
		db = db.Where("created_at = ?", params.CreatedAt)
	}
	if !params.UpdatedAt.IsZero() {
		db = db.Where("updated_at = ?", params.UpdatedAt)
	}

	// 分页查询
	if params.Page > 0 && params.PageSize > 0 {
		// 获取总数
		if err := db.Count(&total).Error; err != nil {
			return nil, 0, err
		}
		db = db.Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize)
	}

	if err := db.Find(&mailboxes).Error; err != nil {
		return nil, 0, err
	}
	if params.Page <= 0 || params.PageSize <= 0 {
		total = int64(len(mailboxes))
	}
	return mailboxes, total, nil
}

// BatchDelete 批量删除邮箱
func (m *MailboxModel) BatchDelete(ids []int64) error {
	return m.db.Where("id IN ?", ids).Delete(&Mailbox{}).Error
}

// CheckEmailExist 检查邮箱是否存在
func (m *MailboxModel) CheckEmailExist(email string) (bool, error) {
	var count int64
	err := m.db.Model(&Mailbox{}).
		Where("email = ?", email).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// UpdateStatus 更新邮箱状态
func (m *MailboxModel) UpdateStatus(id int64, isActive bool) error {
	return m.db.Model(&Mailbox{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_active":  isActive,
		"updated_at": time.Now(),
	}).Error
}

// GetActiveMailboxes 获取激活的邮箱列表
func (m *MailboxModel) GetActiveMailboxes() ([]*Mailbox, error) {
	var mailboxes []*Mailbox
	err := m.db.Where("is_active = ?", true).Find(&mailboxes).Error
	return mailboxes, err
}

// GetMailboxesByUserId 根据用户ID获取邮箱列表
func (m *MailboxModel) GetMailboxesByUserId(userId int64) ([]*Mailbox, error) {
	var mailboxes []*Mailbox
	err := m.db.Where("user_id = ? AND is_active = ?", userId, true).Find(&mailboxes).Error
	return mailboxes, err
}

// GetMailboxesByAdminId 根据管理员ID获取邮箱列表
func (m *MailboxModel) GetMailboxesByAdminId(adminId int64) ([]*Mailbox, error) {
	var mailboxes []*Mailbox
	err := m.db.Where("admin_id = ? AND is_active = ?", adminId, true).Find(&mailboxes).Error
	return mailboxes, err
}

// GetMailboxesByDomainId 根据域名ID获取邮箱列表
func (m *MailboxModel) GetMailboxesByDomainId(domainId int64) ([]*Mailbox, error) {
	var mailboxes []*Mailbox
	err := m.db.Where("domain_id = ? AND is_active = ?", domainId, true).Find(&mailboxes).Error
	return mailboxes, err
}

// UpdatePassword 更新邮箱密码
func (m *MailboxModel) UpdatePassword(id int64, password string) error {
	return m.db.Model(&Mailbox{}).Where("id = ?", id).Updates(map[string]interface{}{
		"password":   password,
		"updated_at": time.Now(),
	}).Error
}

// SoftDelete 软删除邮箱
func (m *MailboxModel) SoftDelete(id int64) error {
	return m.db.Model(&Mailbox{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_active":  false,
		"updated_at": time.Now(),
	}).Error
}

// HardDelete 硬删除邮箱
func (m *MailboxModel) HardDelete(id int64) error {
	return m.db.Where("id = ?", id).Delete(&Mailbox{}).Error
}

// GetMailboxWithOwner 获取邮箱及其所有者信息
func (m *MailboxModel) GetMailboxWithOwner(id int64) (*Mailbox, error) {
	var mailbox Mailbox
	err := m.db.Preload("User").Preload("Admin").First(&mailbox, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &mailbox, nil
}

// CountMailboxesByUserId 统计用户的邮箱数量
func (m *MailboxModel) CountMailboxesByUserId(userId int64) (int64, error) {
	var count int64
	err := m.db.Model(&Mailbox{}).Where("user_id = ? AND is_active = ?", userId, true).Count(&count).Error
	return count, err
}

// CountMailboxesByAdminId 统计管理员的邮箱数量
func (m *MailboxModel) CountMailboxesByAdminId(adminId int64) (int64, error) {
	var count int64
	err := m.db.Model(&Mailbox{}).Where("admin_id = ? AND is_active = ?", adminId, true).Count(&count).Error
	return count, err
}

// CountMailboxesByDomainId 统计域名下的邮箱数量
func (m *MailboxModel) CountMailboxesByDomainId(domainId int64) (int64, error) {
	var count int64
	err := m.db.Model(&Mailbox{}).Where("domain_id = ? AND is_active = ?", domainId, true).Count(&count).Error
	return count, err
}
