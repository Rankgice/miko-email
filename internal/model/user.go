package model

import (
	"time"

	"gorm.io/gorm"
)

// User 普通用户模型
type User struct {
	Id           int64     `gorm:"column:id;primaryKey;autoIncrement;comment:数据库主键ID" json:"id"`               // 数据库主键ID
	Username     string    `gorm:"column:username;uniqueIndex;not null;comment:用户名" json:"username"`           // 用户名
	Password     string    `gorm:"column:password;not null;comment:密码" json:"password,omitempty"`              // 密码
	Email        string    `gorm:"column:email;uniqueIndex;not null;comment:邮箱地址" json:"email"`                // 邮箱地址
	IsActive     bool      `gorm:"column:is_active;default:1;comment:是否激活" json:"is_active"`                   // 是否激活
	Contribution int       `gorm:"column:contribution;default:0;comment:贡献度" json:"contribution"`              // 贡献度
	InviteCode   string    `gorm:"column:invite_code;uniqueIndex;not null;comment:邀请码" json:"invite_code"`     // 邀请码
	InvitedBy    *int64    `gorm:"column:invited_by;comment:被谁邀请" json:"invited_by,omitempty"`                 // 被谁邀请
	CreatedAt    time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP;comment:创建时间" json:"created_at"` // 创建时间
	UpdatedAt    time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;comment:更新时间" json:"updated_at"` // 更新时间
}

// TableName 指定表名
func (User) TableName() string {
	return "user"
}

// UserModel 用户模型
type UserModel struct {
	db *gorm.DB
}

// NewUserModel 创建用户模型
func NewUserModel(db *gorm.DB) *UserModel {
	return &UserModel{
		db: db,
	}
}

// Create 创建用户
func (m *UserModel) Create(user *User) error {
	return m.db.Create(user).Error
}

// Update 更新用户
func (m *UserModel) Update(tx *gorm.DB, user *User) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Updates(user).Error
}

// MapUpdate 更新用户
func (m *UserModel) MapUpdate(tx *gorm.DB, id int64, data map[string]interface{}) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Model(&User{}).Where("id = ?", id).Updates(data).Error
}

// Save 保存用户
func (m *UserModel) Save(tx *gorm.DB, user *User) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Save(user).Error
}

// Delete 删除用户
func (m *UserModel) Delete(user *User) error {
	return m.db.Delete(user).Error
}

// GetById 根据ID获取用户
func (m *UserModel) GetById(id int64) (*User, error) {
	var user User
	if err := m.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (m *UserModel) GetByUsername(username string) (*User, error) {
	var user User
	if err := m.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (m *UserModel) GetByEmail(email string) (*User, error) {
	var user User
	if err := m.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByInviteCode 根据邀请码获取用户
func (m *UserModel) GetByInviteCode(inviteCode string) (*User, error) {
	var user User
	if err := m.db.Where("invite_code = ?", inviteCode).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// List 获取用户列表（统一查询方法）
func (m *UserModel) List(params UserReq) ([]*User, int64, error) {
	var users []*User
	var total int64

	db := m.db.Model(&User{})

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
	if params.InvitedBy != nil {
		db = db.Where("invited_by = ?", *params.InvitedBy)
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

	if err := db.Find(&users).Error; err != nil {
		return nil, 0, err
	}
	if params.Page <= 0 || params.PageSize <= 0 {
		total = int64(len(users))
	}
	return users, total, nil
}

// BatchDelete 批量删除用户
func (m *UserModel) BatchDelete(ids []int64) error {
	return m.db.Where("id IN ?", ids).Delete(&User{}).Error
}

// CheckUsernameExist 检查用户名是否存在
func (m *UserModel) CheckUsernameExist(username string) (bool, error) {
	var count int64
	err := m.db.Model(&User{}).
		Where("username = ?", username).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// CheckEmailExist 检查邮箱是否存在
func (m *UserModel) CheckEmailExist(email string) (bool, error) {
	var count int64
	err := m.db.Model(&User{}).
		Where("email = ?", email).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// UpdateStatus 更新用户状态
func (m *UserModel) UpdateStatus(id int64, isActive bool) error {
	return m.db.Model(&User{}).Where("id = ?", id).Update("is_active", isActive).Error
}

// UpdateContribution 更新用户贡献度
func (m *UserModel) UpdateContribution(id int64, contribution int) error {
	return m.db.Model(&User{}).Where("id = ?", id).Update("contribution", contribution).Error
}

// GetActiveUsers 获取激活的用户列表
func (m *UserModel) GetActiveUsers() ([]*User, error) {
	var users []*User
	err := m.db.Where("is_active = ?", true).Find(&users).Error
	return users, err
}

// GetUsersByInviter 根据邀请人获取用户列表
func (m *UserModel) GetUsersByInviter(inviterId int64) ([]*User, error) {
	var users []*User
	err := m.db.Where("invited_by = ?", inviterId).Find(&users).Error
	return users, err
}

// AuthenticateUser 验证用户登录
func (m *UserModel) AuthenticateUser(username string) (*User, error) {
	var user User
	err := m.db.Where("username = ? AND is_active = ?", username, true).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
