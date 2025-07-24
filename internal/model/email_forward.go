package model

import (
	"database/sql"
	"errors"
	"time"

	"gorm.io/gorm"
)

// EmailForward 邮件转发模型
type EmailForward struct {
	Id                 int64      `gorm:"column:id;primaryKey;autoIncrement;comment:数据库主键ID" json:"id"`                    // 数据库主键ID
	MailboxId          int64      `gorm:"column:mailbox_id;not null;comment:邮箱ID" json:"mailbox_id"`                       // 邮箱ID
	SourceEmail        string     `gorm:"column:source_email;not null;comment:源邮箱" json:"source_email"`                    // 源邮箱
	TargetEmail        string     `gorm:"column:target_email;not null;comment:目标邮箱" json:"target_email"`                   // 目标邮箱
	Enabled            bool       `gorm:"column:enabled;default:1;comment:是否启用" json:"enabled"`                            // 是否启用
	KeepOriginal       bool       `gorm:"column:keep_original;default:1;comment:是否保留原邮件" json:"keep_original"`             // 是否保留原邮件
	ForwardAttachments bool       `gorm:"column:forward_attachments;default:1;comment:是否转发附件" json:"forward_attachments"`  // 是否转发附件
	SubjectPrefix      string     `gorm:"column:subject_prefix;default:[转发];comment:主题前缀" json:"subject_prefix,omitempty"` // 主题前缀
	Description        string     `gorm:"column:description;comment:描述" json:"description,omitempty"`                      // 描述
	ForwardCount       int64      `gorm:"column:forward_count;default:0;comment:转发次数" json:"forward_count"`                // 转发次数
	LastForwardAt      *time.Time `gorm:"column:last_forward_at;comment:最后转发时间" json:"last_forward_at,omitempty"`          // 最后转发时间
	CreatedAt          time.Time  `gorm:"column:created_at;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"`      // 创建时间
	UpdatedAt          time.Time  `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"`      // 更新时间
}

// TableName 指定表名
func (EmailForward) TableName() string {
	return "email_forwards"
}

// EmailForwardModel 邮件转发模型
type EmailForwardModel struct {
	db *gorm.DB
}

// NewEmailForwardModel 创建邮件转发模型
func NewEmailForwardModel(db *gorm.DB) *EmailForwardModel {
	return &EmailForwardModel{
		db: db,
	}
}

// Create 创建邮件转发
func (m *EmailForwardModel) Create(emailForward *EmailForward) error {
	return m.db.Create(emailForward).Error
}

// Update 更新邮件转发
func (m *EmailForwardModel) Update(tx *gorm.DB, emailForward *EmailForward) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Updates(emailForward).Error
}

// MapUpdate 更新邮件转发
func (m *EmailForwardModel) MapUpdate(tx *gorm.DB, id int64, data map[string]interface{}) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Model(&EmailForward{}).Where("id = ?", id).Updates(data).Error
}

// Save 保存邮件转发
func (m *EmailForwardModel) Save(tx *gorm.DB, emailForward *EmailForward) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Save(emailForward).Error
}

// Delete 删除邮件转发
func (m *EmailForwardModel) Delete(emailForward *EmailForward) error {
	return m.db.Delete(emailForward).Error
}

// GetById 根据ID获取邮件转发
func (m *EmailForwardModel) GetById(id int64) (*EmailForward, error) {
	var emailForward EmailForward
	if err := m.db.First(&emailForward, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 如果没有找到记录，返回nil
		}
		return nil, err
	}
	return &emailForward, nil
}

// GetByIdAndMailboxId 根据ID和邮箱ID获取邮件转发
func (m *EmailForwardModel) GetByIdAndMailboxId(id, mailboxId int64) (*EmailForward, error) {
	var emailForward EmailForward
	if err := m.db.Where("id = ? AND mailbox_id = ?", id, mailboxId).First(&emailForward).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &emailForward, nil
}

// List 获取邮件转发列表（统一查询方法）
func (m *EmailForwardModel) List(params EmailForwardReq) ([]*EmailForward, int64, error) {
	var emailForwards []*EmailForward
	var total int64

	db := m.db.Model(&EmailForward{})

	// 添加查询条件
	if params.Id != 0 {
		db = db.Where("id = ?", params.Id)
	}
	if params.MailboxId != 0 {
		db = db.Where("mailbox_id = ?", params.MailboxId)
	}
	if params.SourceEmail != "" {
		db = db.Where("source_email LIKE ?", "%"+params.SourceEmail+"%")
	}
	if params.TargetEmail != "" {
		db = db.Where("target_email LIKE ?", "%"+params.TargetEmail+"%")
	}
	if params.Enabled != nil {
		db = db.Where("enabled = ?", *params.Enabled)
	}
	if params.KeepOriginal != nil {
		db = db.Where("keep_original = ?", *params.KeepOriginal)
	}
	if params.ForwardAttachments != nil {
		db = db.Where("forward_attachments = ?", *params.ForwardAttachments)
	}
	if params.SubjectPrefix != "" {
		db = db.Where("subject_prefix LIKE ?", "%"+params.SubjectPrefix+"%")
	}
	if params.Description != "" {
		db = db.Where("description LIKE ?", "%"+params.Description+"%")
	}
	if params.ForwardCount != 0 {
		db = db.Where("forward_count = ?", params.ForwardCount)
	}
	if params.LastForwardAt != nil {
		db = db.Where("last_forward_at = ?", *params.LastForwardAt)
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
	if err := db.Order("created_at DESC").Find(&emailForwards).Error; err != nil {
		return nil, 0, err
	}
	if params.Page <= 0 || params.PageSize <= 0 {
		total = int64(len(emailForwards))
	}
	return emailForwards, total, nil
}

// BatchDelete 批量删除邮件转发
func (m *EmailForwardModel) BatchDelete(ids []int64) error {
	return m.db.Where("id IN ?", ids).Delete(&EmailForward{}).Error
}

// UpdateStatus 更新转发状态
func (m *EmailForwardModel) UpdateStatus(id int64, enabled bool) error {
	return m.db.Model(&EmailForward{}).Where("id = ?", id).Updates(map[string]interface{}{
		"enabled":    enabled,
		"updated_at": time.Now(),
	}).Error
}

// GetForwardsByMailboxId 根据邮箱ID获取转发规则列表
func (m *EmailForwardModel) GetForwardsByMailboxId(mailboxId int64) ([]*EmailForward, error) {
	var emailForwards []*EmailForward
	err := m.db.Where("mailbox_id = ?", mailboxId).Order("created_at DESC").Find(&emailForwards).Error
	return emailForwards, err
}

// GetEnabledForwardsByMailboxId 根据邮箱ID获取启用的转发规则列表
func (m *EmailForwardModel) GetEnabledForwardsByMailboxId(mailboxId int64) ([]*EmailForward, error) {
	var emailForwards []*EmailForward
	err := m.db.Where("mailbox_id = ? AND enabled = ?", mailboxId, true).
		Order("created_at DESC").Find(&emailForwards).Error
	return emailForwards, err
}

// GetForwardsBySourceEmail 根据源邮箱获取转发规则
func (m *EmailForwardModel) GetForwardsBySourceEmail(sourceEmail string) ([]*EmailForward, error) {
	var emailForwards []*EmailForward
	err := m.db.Where("source_email = ? AND enabled = ?", sourceEmail, true).
		Order("created_at DESC").Find(&emailForwards).Error
	return emailForwards, err
}

// IncrementForwardCount 增加转发次数
func (m *EmailForwardModel) IncrementForwardCount(id int64) error {
	now := time.Now()
	return m.db.Model(&EmailForward{}).Where("id = ?", id).Updates(map[string]interface{}{
		"forward_count":   gorm.Expr("forward_count + 1"),
		"last_forward_at": &now,
		"updated_at":      now,
	}).Error
}

// UpdateForwardSettings 更新转发设置
func (m *EmailForwardModel) UpdateForwardSettings(id int64, keepOriginal, forwardAttachments bool, subjectPrefix string) error {
	return m.db.Model(&EmailForward{}).Where("id = ?", id).Updates(map[string]interface{}{
		"keep_original":       keepOriginal,
		"forward_attachments": forwardAttachments,
		"subject_prefix":      subjectPrefix,
		"updated_at":          time.Now(),
	}).Error
}

// CheckForwardRuleExist 检查转发规则是否存在
func (m *EmailForwardModel) CheckForwardRuleExist(mailboxId int64, sourceEmail, targetEmail string) (bool, error) {
	var count int64
	err := m.db.Model(&EmailForward{}).
		Where("mailbox_id = ? AND source_email = ? AND target_email = ?", mailboxId, sourceEmail, targetEmail).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// DeleteForwardsByMailboxId 删除邮箱下的所有转发规则
func (m *EmailForwardModel) DeleteForwardsByMailboxId(mailboxId int64) error {
	return m.db.Where("mailbox_id = ?", mailboxId).Delete(&EmailForward{}).Error
}

// GetForwardStatistics 获取转发统计信息
func (m *EmailForwardModel) GetForwardStatistics(mailboxId int64) (map[string]interface{}, error) {
	var result struct {
		TotalRules    int64 `json:"total_rules"`
		EnabledRules  int64 `json:"enabled_rules"`
		DisabledRules int64 `json:"disabled_rules"`
		TotalForwards int64 `json:"total_forwards"`
	}

	// 获取总规则数
	if err := m.db.Model(&EmailForward{}).Where("mailbox_id = ?", mailboxId).Count(&result.TotalRules).Error; err != nil {
		return nil, err
	}

	// 获取启用规则数
	if err := m.db.Model(&EmailForward{}).Where("mailbox_id = ? AND enabled = ?", mailboxId, true).Count(&result.EnabledRules).Error; err != nil {
		return nil, err
	}

	// 获取禁用规则数
	result.DisabledRules = result.TotalRules - result.EnabledRules

	// 获取总转发次数
	var totalForwards sql.NullInt64
	if err := m.db.Model(&EmailForward{}).Where("mailbox_id = ?", mailboxId).Select("SUM(forward_count)").Scan(&totalForwards).Error; err != nil {
		return nil, err
	}
	if totalForwards.Valid {
		result.TotalForwards = totalForwards.Int64
	}

	return map[string]interface{}{
		"total_rules":    result.TotalRules,
		"enabled_rules":  result.EnabledRules,
		"disabled_rules": result.DisabledRules,
		"total_forwards": result.TotalForwards,
	}, nil
}
