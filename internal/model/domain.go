package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Domain 域名模型
type Domain struct {
	Id                         int64     `gorm:"column:id;primaryKey;autoIncrement;comment:数据库主键ID" json:"id"`               // 数据库主键ID
	Name                       string    `gorm:"column:name;uniqueIndex;not null;comment:域名" json:"name"`                    // 域名
	IsVerified                 bool      `gorm:"column:is_verified;default:0;comment:是否已验证" json:"is_verified"`              // 是否已验证
	IsActive                   bool      `gorm:"column:is_active;default:1;comment:是否激活" json:"is_active"`                   // 是否激活
	MxRecord                   string    `gorm:"column:mx_record;comment:MX记录" json:"mx_record,omitempty"`                   // MX记录
	ARecord                    string    `gorm:"column:a_record;comment:A记录" json:"a_record,omitempty"`                      // A记录
	TxtRecord                  string    `gorm:"column:txt_record;comment:TXT记录" json:"txt_record,omitempty"`                // TXT记录
	SPFRecord                  string    `json:"spf_record" db:"spf_record"`                                                 // SPF记录
	DMARCRecord                string    `json:"dmarc_record" db:"dmarc_record"`                                             // DMARC记录
	DKIMRecord                 string    `json:"dkim_record" db:"dkim_record"`                                               // DKIM记录
	PTRRecord                  string    `json:"ptr_record" db:"ptr_record"`                                                 // PTR记录
	SenderVerificationStatus   string    `json:"sender_verification_status" db:"sender_verification_status"`                 // 发件验证状态
	ReceiverVerificationStatus string    `json:"receiver_verification_status" db:"receiver_verification_status"`             // 收件验证状态
	CreatedAt                  time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
	UpdatedAt                  time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"` // 更新时间
}

// TableName 指定表名
func (Domain) TableName() string {
	return "domain"
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
func (m *DomainModel) Create(tx *gorm.DB, domain *Domain) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Create(domain).Error
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
func (m *DomainModel) Delete(tx *gorm.DB, domain *Domain) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Delete(domain).Error
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
func (m *DomainModel) BatchDelete(tx *gorm.DB, ids []int64) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Where("id IN ?", ids).Delete(&Domain{}).Error
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
func (m *DomainModel) UpdateStatus(tx *gorm.DB, id int64, isActive bool) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Model(&Domain{}).Where("id = ?", id).Update("is_active", isActive).Error
}

// UpdateVerificationStatus 更新域名验证状态
func (m *DomainModel) UpdateVerificationStatus(tx *gorm.DB, id int64, isVerified bool) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Model(&Domain{}).Where("id = ?", id).Updates(map[string]interface{}{
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
func (m *DomainModel) UpdateDNSRecords(tx *gorm.DB, id int64, mxRecord, aRecord, txtRecord string) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Model(&Domain{}).Where("id = ?", id).Updates(map[string]interface{}{
		"mx_record":  mxRecord,
		"a_record":   aRecord,
		"txt_record": txtRecord,
		"updated_at": time.Now(),
	}).Error
}

// UpdateMXRecord 更新MX记录
func (m *DomainModel) UpdateMXRecord(tx *gorm.DB, id int64, mxRecord string) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Model(&Domain{}).Where("id = ?", id).Updates(map[string]interface{}{
		"mx_record":  mxRecord,
		"updated_at": time.Now(),
	}).Error
}

// UpdateARecord 更新A记录
func (m *DomainModel) UpdateARecord(tx *gorm.DB, id int64, aRecord string) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Model(&Domain{}).Where("id = ?", id).Updates(map[string]interface{}{
		"a_record":   aRecord,
		"updated_at": time.Now(),
	}).Error
}

// UpdateTXTRecord 更新TXT记录
func (m *DomainModel) UpdateTXTRecord(tx *gorm.DB, id int64, txtRecord string) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Model(&Domain{}).Where("id = ?", id).Updates(map[string]interface{}{
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

// UpdateSenderVerificationStatus 更新发件验证状态
func (m *DomainModel) UpdateSenderVerificationStatus(tx *gorm.DB, id int64, status string) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Model(&Domain{}).Where("id = ?", id).Updates(map[string]interface{}{
		"sender_verification_status": status,
		"updated_at":                 time.Now(),
	}).Error
}

// UpdateReceiverVerificationStatus 更新收件验证状态
func (m *DomainModel) UpdateReceiverVerificationStatus(tx *gorm.DB, id int64, status string) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Model(&Domain{}).Where("id = ?", id).Updates(map[string]interface{}{
		"receiver_verification_status": status,
		"updated_at":                   time.Now(),
	}).Error
}

// UpdateAllDNSRecords 更新所有DNS记录
func (m *DomainModel) UpdateAllDNSRecords(tx *gorm.DB, id int64, mxRecord, aRecord, txtRecord, spfRecord, dmarcRecord, dkimRecord, ptrRecord string) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Model(&Domain{}).Where("id = ?", id).Updates(map[string]interface{}{
		"mx_record":    mxRecord,
		"a_record":     aRecord,
		"txt_record":   txtRecord,
		"spf_record":   spfRecord,
		"dmarc_record": dmarcRecord,
		"dkim_record":  dkimRecord,
		"ptr_record":   ptrRecord,
		"updated_at":   time.Now(),
	}).Error
}

// CreateWithAllRecords 创建包含所有DNS记录的域名
func (m *DomainModel) CreateWithAllRecords(tx *gorm.DB, domain *Domain) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Create(domain).Error
}

// CheckActiveDomainExists 检查活跃域名是否存在
func (m *DomainModel) CheckActiveDomainExists(name string) (bool, error) {
	var count int64
	err := m.db.Model(&Domain{}).
		Where("name = ? AND is_active = ?", name, true).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetFirstActiveDomain 获取第一个活跃域名（排除localhost）
func (m *DomainModel) GetFirstActiveDomain() (string, error) {
	var domain Domain
	err := m.db.Where("is_active = ? AND name != ?", true, "localhost").
		Order("id").
		First(&domain).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil // 没有找到活跃域名，返回空字符串
		}
		return "", err
	}

	return domain.Name, nil
}
