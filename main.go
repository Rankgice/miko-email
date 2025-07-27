package main

import (
	"log"
	"miko-email/internal/config"
	"miko-email/internal/server"
	"miko-email/internal/services/email"
	"miko-email/internal/svc"
)

func main() {
	// 加载配置
	cfg := config.Load()

	//初始化上下文
	svcCtx := svc.NewServiceContext(*cfg)

	// 从GORM获取原生SQL数据库连接
	sqlDB, err := svcCtx.DB.DB()
	if err != nil {
		log.Fatal("Failed to get SQL DB from GORM:", err)
	}

	// 创建邮件服务
	emailService := email.NewService(svcCtx)

	// 启动邮件服务器（SMTP, IMAP, POP3）
	// 启动多个SMTP端口
	smtpPorts := cfg.GetSMTPPorts()
	if cfg.EnableMultiSMTP {
		log.Printf("启用多SMTP端口模式，端口: %v", smtpPorts)
	}

	for _, port := range smtpPorts {
		go func(smtpPort string) {
			log.Printf("SMTP server starting on port %s", smtpPort)
			if err := emailService.StartSMTPServer(smtpPort); err != nil {
				log.Printf("SMTP server error (port %s): %v", smtpPort, err)
			}
		}(port)
	}

	go func() {
		if err := emailService.StartIMAPServer(cfg.IMAPPort); err != nil {
			log.Printf("IMAP server error: %v", err)
		}
	}()

	go func() {
		if err := emailService.StartPOP3Server(cfg.POP3Port); err != nil {
			log.Printf("POP3 server error: %v", err)
		}
	}()

	// 启动Web服务器
	webServer := server.New(sqlDB, cfg, svcCtx)
	log.Printf("Starting web server on port %s", cfg.WebPort)
	if err := webServer.Start(); err != nil {
		log.Fatal("Failed to start web server:", err)
	}
}
