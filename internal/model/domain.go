package model

import (
	"time"

	"gorm.io/gorm"
)

// Domain 域名模型
type Domain struct {
	Id         int64     `gorm:"column:id;primaryKey;autoIncrement;comment:数据库主键ID" json:"id"`               // 数据库主键ID
	Name       string    `gorm:"column:name;uniqueIndex;not null;comment:域名" json:"name"`                    // 域名
	IsVerified bool      `gorm:"column:is_verified;default:0;comment:是否已验证" json:"is_verified"`              // 是否已验证
	IsActive   bool      `gorm:"column:is_active;default:1;comment:是否激活" json:"is_active"`                   // 是否激活
	MxRecord   string    `gorm:"column:mx_record;comment:MX记录" json:"mx_record,omitempty"`                   // MX记录
	ARecord    string    `gorm:"column:a_record;comment:A记录" json:"a_record,omitempty"`                      // A记录
	TxtRecord  string    `gorm:"column:txt_record;comment:TXT记录" json:"txt_record,omitempty"`                // TXT记录
	CreatedAt  time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
	UpdatedAt  time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"` // 更新时间
}

// TableName 指定表名
func (Domain) TableName() string {
	return "domains"
}

// DomainModel 域名模型
type DomainModel struct {
	db *gorm.DB
}

// NewDomainModel 创建域名模型
func NewDomainModel(db *gorm.DB) *DomainModel {
	return &DomainModel{
		db: db,
	}
}

// Create 创建域名
func (m *DomainModel) Create(domain *Domain) error {
	return m.db.Create(domain).Error
}

// Update 更新域名
func (m *DomainModel) Update(tx *gorm.DB, domain *Domain) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Updates(domain).Error
}

// MapUpdate 更新域名
func (m *DomainModel) MapUpdate(tx *gorm.DB, id int64, data map[string]interface{}) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Model(&Domain{}).Where("id = ?", id).Updates(data).Error
}

// Save 保存域名
func (m *DomainModel) Save(tx *gorm.DB, domain *Domain) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Save(domain).Error
}

// Delete 删除域名
func (m *DomainModel) Delete(domain *Domain) error {
	return m.db.Delete(domain).Error
}

// GetById 根据ID获取域名
func (m *DomainModel) GetById(id int64) (*Domain, error) {
	var domain Domain
	if err := m.db.First(&domain, id).Error; err != nil {
		return nil, err
	}
	return &domain, nil
}

// GetByName 根据域名获取记录
func (m *DomainModel) GetByName(name string) (*Domain, error) {
	var domain Domain
	if err := m.db.Where("name = ?", name).First(&domain).Error; err != nil {
		return nil, err
	}
	return &domain, nil
}

// List 获取域名列表（统一查询方法）
func (m *DomainModel) List(params DomainReq) ([]*Domain, int64, error) {
	var domains []*Domain
	var total int64

	db := m.db.Model(&Domain{})

	// 添加查询条件
	if params.Id != 0 {
		db = db.Where("id = ?", params.Id)
	}
	if params.Name != "" {
		db = db.Where("name LIKE ?", "%"+params.Name+"%")
	}
	if params.IsVerified != nil {
		db = db.Where("is_verified = ?", *params.IsVerified)
	}
	if params.IsActive != nil {
		db = db.Where("is_active = ?", *params.IsActive)
	}
	if params.MxRecord != "" {
		db = db.Where("mx_record LIKE ?", "%"+params.MxRecord+"%")
	}
	if params.ARecord != "" {
		db = db.Where("a_record LIKE ?", "%"+params.ARecord+"%")
	}
	if params.TxtRecord != "" {
		db = db.Where("txt_record LIKE ?", "%"+params.TxtRecord+"%")
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

	if err := db.Find(&domains).Error; err != nil {
		return nil, 0, err
	}
	if params.Page <= 0 || params.PageSize <= 0 {
		total = int64(len(domains))
	}
	return domains, total, nil
}

// BatchDelete 批量删除域名
func (m *DomainModel) BatchDelete(ids []int64) error {
	return m.db.Where("id IN ?", ids).Delete(&Domain{}).Error
}

// CheckDomainExist 检查域名是否存在
func (m *DomainModel) CheckDomainExist(name string) (bool, error) {
	var count int64
	err := m.db.Model(&Domain{}).
		Where("name = ?", name).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// UpdateStatus 更新域名状态
func (m *DomainModel) UpdateStatus(id int64, isActive bool) error {
	return m.db.Model(&Domain{}).Where("id = ?", id).Update("is_active", isActive).Error
}

// UpdateVerificationStatus 更新域名验证状态
func (m *DomainModel) UpdateVerificationStatus(id int64, isVerified bool) error {
	return m.db.Model(&Domain{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_verified": isVerified,
		"updated_at":  time.Now(),
	}).Error
}

// GetActiveDomains 获取激活的域名列表
func (m *DomainModel) GetActiveDomains() ([]*Domain, error) {
	var domains []*Domain
	err := m.db.Where("is_active = ?", true).Find(&domains).Error
	return domains, err
}

// GetVerifiedDomains 获取已验证的域名列表
func (m *DomainModel) GetVerifiedDomains() ([]*Domain, error) {
	var domains []*Domain
	err := m.db.Where("is_verified = ?", true).Find(&domains).Error
	return domains, err
}

// GetAvailableDomains 获取可用的域名列表（已验证且激活的）
func (m *DomainModel) GetAvailableDomains() ([]*Domain, error) {
	var domains []*Domain
	err := m.db.Where("is_active = ? AND is_verified = ?", true, true).Find(&domains).Error
	return domains, err
}

// UpdateDNSRecords 更新DNS记录
func (m *DomainModel) UpdateDNSRecords(id int64, mxRecord, aRecord, txtRecord string) error {
	return m.db.Model(&Domain{}).Where("id = ?", id).Updates(map[string]interface{}{
		"mx_record":  mxRecord,
		"a_record":   aRecord,
		"txt_record": txtRecord,
		"updated_at": time.Now(),
	}).Error
}

// UpdateMXRecord 更新MX记录
func (m *DomainModel) UpdateMXRecord(id int64, mxRecord string) error {
	return m.db.Model(&Domain{}).Where("id = ?", id).Updates(map[string]interface{}{
		"mx_record":  mxRecord,
		"updated_at": time.Now(),
	}).Error
}

// UpdateARecord 更新A记录
func (m *DomainModel) UpdateARecord(id int64, aRecord string) error {
	return m.db.Model(&Domain{}).Where("id = ?", id).Updates(map[string]interface{}{
		"a_record":   aRecord,
		"updated_at": time.Now(),
	}).Error
}

// UpdateTXTRecord 更新TXT记录
func (m *DomainModel) UpdateTXTRecord(id int64, txtRecord string) error {
	return m.db.Model(&Domain{}).Where("id = ?", id).Updates(map[string]interface{}{
		"txt_record": txtRecord,
		"updated_at": time.Now(),
	}).Error
}

// GetDomainsByStatus 根据状态获取域名列表
func (m *DomainModel) GetDomainsByStatus(isActive, isVerified bool) ([]*Domain, error) {
	var domains []*Domain
	err := m.db.Where("is_active = ? AND is_verified = ?", isActive, isVerified).Find(&domains).Error
	return domains, err
}
