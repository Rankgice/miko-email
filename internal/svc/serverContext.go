package svc

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"miko-email/internal/config"
	"miko-email/internal/model"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config            config.Config
	DB                *gorm.DB
	UserModel         *model.UserModel
	AdminModel        *model.AdminModel
	DomainModel       *model.DomainModel
	MailboxModel      *model.MailboxModel
	EmailModel        *model.EmailModel
	EmailForwardModel *model.EmailForwardModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(sqlite.Open("miko_email.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("连接数据库失败", err)
	}

	// 检测数据库文件是否可以正常读写
	const permissionTestTable = "_permission_test"
	if err := db.Exec("CREATE TABLE " + permissionTestTable + " (`id` int);").Error; err != nil {
		log.Fatal("数据库没有写入权限", "error", err.Error())
	}
	if err := db.Exec("DROP TABLE " + permissionTestTable + ";").Error; err != nil {
		log.Fatal("删除测试表失败", err)
	}

	// 自动迁移数据表
	if err = autoMigrate(db); err != nil {
		log.Fatal("数据表迁移失败", err)
	}

	// 创建模型实例用于创建默认管理员
	adminModel := model.NewAdminModel(db)

	// 创建默认管理员
	if err := createDefaultAdmin(adminModel); err != nil {
		log.Fatal("创建默认管理员失败", err)
	}

	return &ServiceContext{
		Config:            c,
		DB:                db,
		UserModel:         model.NewUserModel(db),
		AdminModel:        model.NewAdminModel(db),
		DomainModel:       model.NewDomainModel(db),
		MailboxModel:      model.NewMailboxModel(db),
		EmailModel:        model.NewEmailModel(db),
		EmailForwardModel: model.NewEmailForwardModel(db),
	}
}

// 自动迁移数据表结构
func autoMigrate(db *gorm.DB) error {
	// 自动迁移所有模型
	return db.AutoMigrate(
		&model.User{},
		&model.Admin{},
		&model.Domain{},
		&model.Mailbox{},
		&model.Email{},
		&model.EmailForward{},
	)
}

func createDefaultAdmin(adminModel *model.AdminModel) error {
	// 获取配置中的管理员信息
	username, password, email, enabled := config.GetAdminCredentials()

	// 如果管理员被禁用，跳过创建
	if !enabled {
		return nil
	}

	// 检查是否已存在该用户名的管理员
	existingAdmin, err := adminModel.GetByUsername(username)
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("检查管理员是否存在失败: %w", err)
	}

	if existingAdmin != nil {
		// 管理员已存在，检查是否需要更新信息
		return updateExistingAdmin(adminModel, existingAdmin, password, email, enabled)
	}

	// 创建新管理员
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	inviteCode := generateInviteCode()

	newAdmin := &model.Admin{
		Username:     username,
		Password:     string(hashedPassword),
		Email:        email,
		IsActive:     enabled,
		Contribution: 0,
		InviteCode:   inviteCode,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := adminModel.Create(newAdmin); err != nil {
		return fmt.Errorf("创建管理员失败: %w", err)
	}

	log.Printf("默认管理员创建成功: %s", username)
	return nil
}

func updateExistingAdmin(adminModel *model.AdminModel, existingAdmin *model.Admin, password, email string, enabled bool) error {
	// 更新现有管理员的信息
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	// 更新管理员信息
	updateData := map[string]interface{}{
		"password":   string(hashedPassword),
		"email":      email,
		"is_active":  enabled,
		"updated_at": time.Now(),
	}

	if err := adminModel.MapUpdate(nil, existingAdmin.Id, updateData); err != nil {
		return fmt.Errorf("更新管理员信息失败: %w", err)
	}

	log.Printf("管理员信息更新成功: %s", existingAdmin.Username)
	return nil
}

func generateInviteCode() string {
	// 简单的邀请码生成，实际项目中应该使用更安全的方法
	return fmt.Sprintf("INVITE_%d", time.Now().Unix())
}
