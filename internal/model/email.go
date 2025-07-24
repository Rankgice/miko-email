package model

import (
	"time"

	"gorm.io/gorm"
)

// Email 邮件模型
type Email struct {
	Id        int64     `gorm:"column:id;primaryKey;autoIncrement;comment:数据库主键ID" json:"id"`               // 数据库主键ID
	MailboxId int64     `gorm:"column:mailbox_id;not null;comment:邮箱ID" json:"mailbox_id"`                  // 邮箱ID
	FromAddr  string    `gorm:"column:from_addr;not null;comment:发件人" json:"from_addr"`                     // 发件人
	ToAddr    string    `gorm:"column:to_addr;not null;comment:收件人" json:"to_addr"`                         // 收件人
	Subject   string    `gorm:"column:subject;comment:主题" json:"subject,omitempty"`                         // 主题
	Body      string    `gorm:"column:body;comment:邮件内容" json:"body,omitempty"`                             // 邮件内容
	IsRead    bool      `gorm:"column:is_read;default:0;comment:是否已读" json:"is_read"`                       // 是否已读
	Folder    string    `gorm:"column:folder;default:inbox;comment:文件夹" json:"folder"`                      // 文件夹 (inbox, sent, trash)
	CreatedAt time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"` // 更新时间
}

// TableName 指定表名
func (Email) TableName() string {
	return "emails"
}

// EmailModel 邮件模型
type EmailModel struct {
	db *gorm.DB
}

// NewEmailModel 创建邮件模型
func NewEmailModel(db *gorm.DB) *EmailModel {
	return &EmailModel{
		db: db,
	}
}

// Create 创建邮件
func (m *EmailModel) Create(email *Email) error {
	return m.db.Create(email).Error
}

// Update 更新邮件
func (m *EmailModel) Update(tx *gorm.DB, email *Email) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Updates(email).Error
}

// MapUpdate 更新邮件
func (m *EmailModel) MapUpdate(tx *gorm.DB, id int64, data map[string]interface{}) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Model(&Email{}).Where("id = ?", id).Updates(data).Error
}

// Save 保存邮件
func (m *EmailModel) Save(tx *gorm.DB, email *Email) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Save(email).Error
}

// Delete 删除邮件
func (m *EmailModel) Delete(email *Email) error {
	return m.db.Delete(email).Error
}

// GetById 根据ID获取邮件
func (m *EmailModel) GetById(id int64) (*Email, error) {
	var email Email
	if err := m.db.First(&email, id).Error; err != nil {
		return nil, err
	}
	return &email, nil
}

// GetByIdAndMailboxId 根据ID和邮箱ID获取邮件
func (m *EmailModel) GetByIdAndMailboxId(id, mailboxId int64) (*Email, error) {
	var email Email
	if err := m.db.Where("id = ? AND mailbox_id = ?", id, mailboxId).First(&email).Error; err != nil {
		return nil, err
	}
	return &email, nil
}

// List 获取邮件列表（统一查询方法）
func (m *EmailModel) List(params EmailReq) ([]*Email, int64, error) {
	var emails []*Email
	var total int64

	db := m.db.Model(&Email{})

	// 添加查询条件
	if params.Id != 0 {
		db = db.Where("id = ?", params.Id)
	}
	if params.MailboxId != 0 {
		db = db.Where("mailbox_id = ?", params.MailboxId)
	}
	if params.FromAddr != "" {
		db = db.Where("from_addr LIKE ?", "%"+params.FromAddr+"%")
	}
	if params.ToAddr != "" {
		db = db.Where("to_addr LIKE ?", "%"+params.ToAddr+"%")
	}
	if params.Subject != "" {
		db = db.Where("subject LIKE ?", "%"+params.Subject+"%")
	}
	if params.Body != "" {
		db = db.Where("body LIKE ?", "%"+params.Body+"%")
	}
	if params.IsRead != nil {
		db = db.Where("is_read = ?", *params.IsRead)
	}
	if params.Folder != "" {
		db = db.Where("folder = ?", params.Folder)
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

	// 按创建时间倒序排列
	if err := db.Order("created_at DESC").Find(&emails).Error; err != nil {
		return nil, 0, err
	}
	if params.Page <= 0 || params.PageSize <= 0 {
		total = int64(len(emails))
	}
	return emails, total, nil
}

// BatchDelete 批量删除邮件
func (m *EmailModel) BatchDelete(ids []int64) error {
	return m.db.Where("id IN ?", ids).Delete(&Email{}).Error
}

// MarkAsRead 标记邮件为已读
func (m *EmailModel) MarkAsRead(id int64) error {
	return m.db.Model(&Email{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_read":    true,
		"updated_at": time.Now(),
	}).Error
}

// MarkAsUnread 标记邮件为未读
func (m *EmailModel) MarkAsUnread(id int64) error {
	return m.db.Model(&Email{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_read":    false,
		"updated_at": time.Now(),
	}).Error
}

// MoveToFolder 移动邮件到指定文件夹
func (m *EmailModel) MoveToFolder(id int64, folder string) error {
	return m.db.Model(&Email{}).Where("id = ?", id).Updates(map[string]interface{}{
		"folder":     folder,
		"updated_at": time.Now(),
	}).Error
}

// GetEmailsByMailboxId 根据邮箱ID获取邮件列表
func (m *EmailModel) GetEmailsByMailboxId(mailboxId int64, folder string, page, pageSize int) ([]*Email, int64, error) {
	var emails []*Email
	var total int64

	db := m.db.Model(&Email{}).Where("mailbox_id = ?", mailboxId)

	if folder != "" {
		db = db.Where("folder = ?", folder)
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if page > 0 && pageSize > 0 {
		db = db.Offset((page - 1) * pageSize).Limit(pageSize)
	}

	if err := db.Order("created_at DESC").Find(&emails).Error; err != nil {
		return nil, 0, err
	}

	return emails, total, nil
}

// GetUnreadCount 获取未读邮件数量
func (m *EmailModel) GetUnreadCount(mailboxId int64, folder string) (int64, error) {
	var count int64
	db := m.db.Model(&Email{}).Where("mailbox_id = ? AND is_read = ?", mailboxId, false)

	if folder != "" {
		db = db.Where("folder = ?", folder)
	}

	err := db.Count(&count).Error
	return count, err
}

// GetEmailsByFolder 根据文件夹获取邮件列表
func (m *EmailModel) GetEmailsByFolder(mailboxId int64, folder string) ([]*Email, error) {
	var emails []*Email
	err := m.db.Where("mailbox_id = ? AND folder = ?", mailboxId, folder).
		Order("created_at DESC").Find(&emails).Error
	return emails, err
}

// SearchEmails 搜索邮件
func (m *EmailModel) SearchEmails(mailboxId int64, keyword string, page, pageSize int) ([]*Email, int64, error) {
	var emails []*Email
	var total int64

	db := m.db.Model(&Email{}).Where("mailbox_id = ?", mailboxId)

	if keyword != "" {
		db = db.Where("subject LIKE ? OR body LIKE ? OR from_addr LIKE ? OR to_addr LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if page > 0 && pageSize > 0 {
		db = db.Offset((page - 1) * pageSize).Limit(pageSize)
	}

	if err := db.Order("created_at DESC").Find(&emails).Error; err != nil {
		return nil, 0, err
	}

	return emails, total, nil
}

// DeleteEmailsByMailboxId 删除邮箱下的所有邮件
func (m *EmailModel) DeleteEmailsByMailboxId(mailboxId int64) error {
	return m.db.Where("mailbox_id = ?", mailboxId).Delete(&Email{}).Error
}

// BatchMarkAsRead 批量标记为已读
func (m *EmailModel) BatchMarkAsRead(ids []int64) error {
	return m.db.Model(&Email{}).Where("id IN ?", ids).Updates(map[string]interface{}{
		"is_read":    true,
		"updated_at": time.Now(),
	}).Error
}

// BatchMoveToFolder 批量移动到文件夹
func (m *EmailModel) BatchMoveToFolder(ids []int64, folder string) error {
	return m.db.Model(&Email{}).Where("id IN ?", ids).Updates(map[string]interface{}{
		"folder":     folder,
		"updated_at": time.Now(),
	}).Error
}
