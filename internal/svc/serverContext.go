package svc

import (
	"log"
	"miko-email/internal/config"
	"miko-email/internal/model"

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
