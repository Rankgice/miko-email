package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// Proxy 代理IP模型
type Proxy struct {
	Id        int64     `gorm:"column:id;primaryKey;autoIncrement;comment:数据库主键ID" json:"id"`           // 数据库主键ID
	IpAddress string    `gorm:"column:ip_address;comment:IP地址" json:"ip_address"`                       // IP地址
	Port      string    `gorm:"column:port;comment:端口" json:"port"`                                     // 端口
	Username  string    `gorm:"column:username;comment:用户名" json:"username,omitempty"`                  // 用户名
	Password  string    `gorm:"column:password;comment:密码" json:"password,omitempty"`                   // 密码
	Type      string    `gorm:"column:type;comment:代理类型" json:"type"`                                   // 代理类型（http/socks5）
	AddTime   time.Time `gorm:"column:add_time;default:CURRENT_TIMESTAMP;comment:添加时间" json:"add_time"` // 添加时间
	Status    int       `gorm:"column:status;comment:状态" json:"status,omitempty"`                       // 状态 1:禁用 2:启用
	Remark    string    `gorm:"column:remark;comment:备注" json:"remark,omitempty"`                       // 备注
}

// TableName 指定表名
func (Proxy) TableName() string {
	return "proxy"
}

// FormattedProxy 获取格式化后的代理字符串
func (p Proxy) FormattedProxy() string {
	if p.IpAddress == "" {
		return ""
	}
	if p.Username == "" || p.Password == "" {
		return p.Type + "://" + p.IpAddress + ":" + p.Port
	}
	return p.Type + "://" + p.Username + ":" + p.Password + "@" + p.IpAddress + ":" + p.Port
}

// ProxyModel 代理IP模型
type ProxyModel struct {
	db *gorm.DB
}

// NewProxyModel 创建代理IP模型
func NewProxyModel(db *gorm.DB) *ProxyModel {
	return &ProxyModel{
		db: db,
	}
}

// Create 创建代理
func (m *ProxyModel) Create(proxy *Proxy) error {
	return m.db.Create(proxy).Error
}

// Update 更新代理
func (m *ProxyModel) Update(tx *gorm.DB, proxy *Proxy) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Updates(proxy).Error
}

// MapUpdate 更新代理
func (m *ProxyModel) MapUpdate(tx *gorm.DB, id int, data map[string]interface{}) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Model(&Proxy{}).Where("id = ?", id).Updates(data).Error
}

// Save 保存代理
func (m *ProxyModel) Save(tx *gorm.DB, proxy *Proxy) error {
	db := m.db
	if tx != nil {
		db = tx
	}
	return db.Save(proxy).Error
}

// Delete 删除代理
func (m *ProxyModel) Delete(proxy *Proxy) error {
	return m.db.Delete(proxy).Error
}

// GetById 根据ID获取代理
func (m *ProxyModel) GetById(id int64) (*Proxy, error) {
	var proxy Proxy
	if err := m.db.First(&proxy, id).Error; err != nil {
		return nil, err
	}
	return &proxy, nil
}

// List 获取代理列表（统一查询方法）
func (m *ProxyModel) List(params ProxyReq) ([]*Proxy, int64, error) {
	var proxies []*Proxy
	var total int64

	db := m.db.Model(&Proxy{})

	// 添加查询条件
	if params.Id != 0 {
		db = db.Where("id = ?", params.Id)
	}
	if params.IpAddress != "" {
		db = db.Where("ip_address LIKE ?", "%"+params.IpAddress+"%")
	}
	if params.Port != "" {
		db = db.Where("port = ?", params.Port)
	}
	if params.Username != "" {
		db = db.Where("username = ?", params.Username)
	}
	if params.Password != "" {
		db = db.Where("password = ?", params.Password)
	}
	if params.Type != "" {
		db = db.Where("type = ?", params.Type)
	}
	if !params.AddTime.IsZero() {
		db = db.Where("add_time = ?", params.AddTime)
	}
	if params.Status != 0 {
		db = db.Where("status = ?", params.Status)
	}
	if params.Remark != "" {
		db = db.Where("remark LIKE ?", "%"+params.Remark+"%")
	}

	// 分页查询
	if params.Page > 0 && params.PageSize > 0 {
		// 获取总数
		if err := db.Count(&total).Error; err != nil {
			return nil, 0, err
		}
		db = db.Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize)
	}

	if err := db.Find(&proxies).Error; err != nil {
		return nil, 0, err
	}
	if params.Page <= 0 || params.PageSize <= 0 {
		total = int64(len(proxies))
	}
	return proxies, total, nil
}

// BatchDelete 批量删除代理
func (m *ProxyModel) BatchDelete(ids []int) error {
	return m.db.Where("id IN ?", ids).Delete(&Proxy{}).Error
}

// CheckProxyExist 检查代理是否存在
func (m *ProxyModel) CheckProxyExist(ipAddress, port string) (bool, error) {
	var count int64
	err := m.db.Model(&Proxy{}).
		Where("ip_address = ? AND port = ?", ipAddress, port).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// UpdateStatus 更新代理状态
func (m *ProxyModel) UpdateStatus(id int64, status int) error {
	return m.db.Model(&Proxy{}).Where("id = ?", id).Update("status", status).Error
}
