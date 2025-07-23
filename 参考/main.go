package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"nbemail/internal/config"
	"nbemail/internal/database"
	"nbemail/internal/imap"
	"nbemail/internal/pop3"
	"nbemail/internal/server"
	"nbemail/internal/smtp"
)

//go:embed web/static
var staticFiles embed.FS

//go:embed web/templates
var templateFiles embed.FS

func main() {
	var port = flag.Int("port", 8080, "Web服务端口")
	var smtpPort = flag.Int("smtp-port", 25, "主SMTP服务端口")
	var enableMultiPorts = flag.Bool("multi-smtp", false, "启用多SMTP端口(25,587,465)")
	var dbPath = flag.String("db", "nbemail.db", "数据库文件路径")
	flag.Parse()

	// 加载配置
	cfg := config.GetDefaults()

	// 显示外部SMTP配置状态
	if cfg.OutboundSMTPHost != "" {
		log.Printf("外部SMTP配置已加载 - Host: %s, Port: %d, User: %s, TLS: %v",
			cfg.OutboundSMTPHost, cfg.OutboundSMTPPort, cfg.OutboundSMTPUser, cfg.OutboundSMTPTLS)
	} else {
		log.Printf("外部SMTP未配置，只能在本地用户间发送邮件")
	}

	// 覆盖命令行参数
	if *port != 8080 {
		cfg.WebPort = *port
	}
	if *smtpPort != 25 {
		cfg.SMTPPort = *smtpPort
	}
	if *dbPath != "nbemail.db" {
		cfg.DBPath = *dbPath
	}

	// 初始化数据库
	db, err := database.Init(cfg.DBPath, cfg.AdminEmail, cfg.AdminPass)
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	defer db.Close()

	// 启动SMTP服务器
	var smtpServers []*smtp.Server

	// 如果命令行指定了多端口模式，则启用配置中的多端口
	if *enableMultiPorts {
		cfg.EnableMultiSMTP = true
	}

	// 获取要启动的SMTP端口列表
	smtpPorts := cfg.GetSMTPPorts()
	if cfg.EnableMultiSMTP {
		log.Printf("启用多SMTP端口模式，端口: %v", smtpPorts)
	}
	for _, smtpPortNum := range smtpPorts {
		// 为每个端口创建独立的配置
		portCfg := *cfg // 复制配置
		portCfg.SMTPPort = smtpPortNum

		smtpServer := smtp.NewServer(&portCfg, db)
		smtpServers = append(smtpServers, smtpServer)
		go func(port int, server *smtp.Server) {
			log.Printf("启动SMTP服务器，端口: %d", port)
			if err := server.Start(); err != nil {
				log.Printf("SMTP服务器启动失败 (端口 %d): %v", port, err)
			}
		}(smtpPortNum, smtpServer)
	}

	// 启动IMAP服务器
	imapServer := imap.NewServer(cfg, db)
	go func() {
		log.Printf("启动IMAP服务器，端口: %d", cfg.IMAPPort)
		if err := imapServer.Start(); err != nil {
			log.Fatalf("IMAP服务器启动失败: %v", err)
		}
	}()

	// 启动POP3服务器
	pop3Server := pop3.NewServer(cfg, db)
	go func() {
		log.Printf("启动POP3服务器，端口: %d", cfg.POP3Port)
		if err := pop3Server.Start(); err != nil {
			log.Fatalf("POP3服务器启动失败: %v", err)
		}
	}()

	// 启动Web服务器
	webServer := server.NewServer(cfg, db, staticFiles, templateFiles)
	go func() {
		log.Printf("启动Web服务器，端口: %d", cfg.WebPort)
		if err := webServer.Start(); err != nil {
			log.Fatalf("Web服务器启动失败: %v", err)
		}
	}()

	// 等待信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("\n正在关闭NBEmail服务器...")
	// 停止所有SMTP服务器
	for _, smtpServer := range smtpServers {
		smtpServer.Stop()
	}
	imapServer.Stop()
	pop3Server.Stop()
	webServer.Stop()
	fmt.Println("NBEmail服务器已关闭")
}
