package model

import (
	"time"

	"gorm.io/gorm"
)

// Admin 管理员模型
type Admin struct {
	Id           int64     `gorm:"column:id;primaryKey;autoIncrement;comment:数据库主键ID" json:"id"`               // 数据库主键ID
	Username     string    `gorm:"column:username;uniqueIndex;not null;comment:用户名" json:"username"`           // 用户名
	Password     string    `gorm:"column:password;not null;comment:密码" json:"password,omitempty"`              // 密码
	Email        string    `gorm:"column:email;uniqueIndex;not null;comment:邮箱地址" json:"email"`                // 邮箱地址
	IsActive     bool      `gorm:"column:is_active;default:1;comment:是否激活" json:"is_active"`                   // 是否激活
	Contribution int       `gorm:"column:contribution;default:0;comment:贡献度" json:"contribution"`              // 贡献度
	InviteCode   string    `gorm:"column:invite_code;uniqueIndex;not null;comment:邀请码" json:"invite_code"`     // 邀请码
	CreatedAt    time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
	UpdatedAt    time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"` // 更新时间
}

// TableName 指定表名
func (Admin) TableName() string {
	return "admin"
}

// AdminModel 管理员模型
type AdminModel struct {
	db *gorm.DB
}

// NewAdminModel 创建管理员模型
func NewAdminModel(db *gorm.DB) *AdminModel {
	return &AdminModel{
		db: db,
	}
}

// Create 创建管理员
func (m *AdminModel) Create(tx *gorm.DB, admin *Admin) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Create(admin).Error
}

// Update 更新管理员
func (m *AdminModel) Update(tx *gorm.DB, admin *Admin) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Updates(admin).Error
}

// MapUpdate 更新管理员
func (m *AdminModel) MapUpdate(tx *gorm.DB, id int64, data map[string]interface{}) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Model(&Admin{}).Where("id = ?", id).Updates(data).Error
}

// Save 保存管理员
func (m *AdminModel) Save(tx *gorm.DB, admin *Admin) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Save(admin).Error
}

// Delete 删除管理员
func (m *AdminModel) Delete(tx *gorm.DB, admin *Admin) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Delete(admin).Error
}

// GetById 根据ID获取管理员
func (m *AdminModel) GetById(id int64) (*Admin, error) {
	var admin Admin
	if err := m.db.First(&admin, id).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

// GetByUsername 根据用户名获取管理员
func (m *AdminModel) GetByUsername(username string) (*Admin, error) {
	var admin Admin
	if err := m.db.Where("username = ?", username).First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

// GetByEmail 根据邮箱获取管理员
func (m *AdminModel) GetByEmail(email string) (*Admin, error) {
	var admin Admin
	if err := m.db.Where("email = ?", email).First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

// GetByInviteCode 根据邀请码获取管理员
func (m *AdminModel) GetByInviteCode(inviteCode string) (*Admin, error) {
	var admin Admin
	if err := m.db.Where("invite_code = ?", inviteCode).First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

// List 获取管理员列表（统一查询方法）
func (m *AdminModel) List(params AdminReq) ([]*Admin, int64, error) {
	var admins []*Admin
	var total int64

	db := m.db.Model(&Admin{})

	// 添加查询条件
	if params.Id != 0 {
		db = db.Where("id = ?", params.Id)
	}
	if params.Username != "" {
		db = db.Where("username LIKE ?", "%"+params.Username+"%")
	}
	if params.Email != "" {
		db = db.Where("email LIKE ?", "%"+params.Email+"%")
	}
	if params.IsActive != nil {
		db = db.Where("is_active = ?", *params.IsActive)
	}
	if params.Contribution != 0 {
		db = db.Where("contribution = ?", params.Contribution)
	}
	if params.InviteCode != "" {
		db = db.Where("invite_code = ?", params.InviteCode)
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

	if err := db.Find(&admins).Error; err != nil {
		return nil, 0, err
	}
	if params.Page <= 0 || params.PageSize <= 0 {
		total = int64(len(admins))
	}
	return admins, total, nil
}

// BatchDelete 批量删除管理员
func (m *AdminModel) BatchDelete(tx *gorm.DB, ids []int64) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Where("id IN ?", ids).Delete(&Admin{}).Error
}

// CheckUsernameExist 检查用户名是否存在
func (m *AdminModel) CheckUsernameExist(username string) (bool, error) {
	var count int64
	err := m.db.Model(&Admin{}).
		Where("username = ?", username).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// CheckEmailExist 检查邮箱是否存在
func (m *AdminModel) CheckEmailExist(email string) (bool, error) {
	var count int64
	err := m.db.Model(&Admin{}).
		Where("email = ?", email).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// UpdateStatus 更新管理员状态
func (m *AdminModel) UpdateStatus(tx *gorm.DB, id int64, isActive bool) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Model(&Admin{}).Where("id = ?", id).Update("is_active", isActive).Error
}

// UpdateContribution 更新管理员贡献度
func (m *AdminModel) UpdateContribution(tx *gorm.DB, id int64, contribution int) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Model(&Admin{}).Where("id = ?", id).Update("contribution", contribution).Error
}

// GetActiveAdmins 获取激活的管理员列表
func (m *AdminModel) GetActiveAdmins() ([]*Admin, error) {
	var admins []*Admin
	err := m.db.Where("is_active = ?", true).Find(&admins).Error
	return admins, err
}

// AuthenticateAdmin 验证管理员登录
func (m *AdminModel) AuthenticateAdmin(username string) (*Admin, error) {
	var admin Admin
	err := m.db.Where("username = ? AND is_active = ?", username, true).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// UpdatePassword 更新管理员密码
func (m *AdminModel) UpdatePassword(tx *gorm.DB, id int64, password string) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Model(&Admin{}).Where("id = ?", id).Updates(map[string]interface{}{
		"password":   password,
		"updated_at": time.Now(),
	}).Error
}

// UpdateEmail 更新管理员邮箱
func (m *AdminModel) UpdateEmail(tx *gorm.DB, id int64, email string) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Model(&Admin{}).Where("id = ?", id).Updates(map[string]interface{}{
		"email":      email,
		"updated_at": time.Now(),
	}).Error
}
