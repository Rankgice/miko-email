package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"strings"
	"time"
)

// SMTP配置
type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func main() {
	// 163邮箱SMTP配置
	config := SMTPConfig{
		Host:     "smtp.163.com",
		Port:     "465",
		Username: "18090776855@163.com",
		Password: "JTH39ZMMBTennqeQ",
		From:     "18090776855@163.com",
	}

	// 测试发送到QQ邮箱
	testRecipient := "2014131458@qq.com"

	fmt.Printf("开始测试发送邮件...\n")
	fmt.Printf("发件人: %s\n", config.From)
	fmt.Printf("收件人: %s\n", testRecipient)
	fmt.Printf("SMTP服务器: %s:%s\n", config.Host, config.Port)

	err := sendEmail(config, testRecipient, "SMTP发件测试", "这是一封来自163邮箱的测试邮件，用于测试SMTP发送功能。\n\n发送时间: "+time.Now().Format("2006-01-02 15:04:05"))

	if err != nil {
		log.Printf("❌ 发送失败: %v", err)
	} else {
		fmt.Printf("✅ 邮件发送成功！\n")
		fmt.Printf("\n现在可以测试您的邮件系统能否接收来自163邮箱的邮件了。\n")
		fmt.Printf("请从163邮箱发送邮件到: kimi11@jbjj.site\n")
	}
}

func sendEmail(config SMTPConfig, to, subject, body string) error {
	// 构建邮件内容
	message := buildMessage(config.From, to, subject, body)

	// 连接到SMTP服务器
	addr := config.Host + ":" + config.Port

	// 创建TLS连接
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         config.Host,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("连接SMTP服务器失败: %v", err)
	}
	defer conn.Close()

	// 创建SMTP客户端
	client, err := smtp.NewClient(conn, config.Host)
	if err != nil {
		return fmt.Errorf("创建SMTP客户端失败: %v", err)
	}
	defer client.Quit()

	// 认证
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP认证失败: %v", err)
	}

	// 设置发件人
	if err = client.Mail(config.From); err != nil {
		return fmt.Errorf("设置发件人失败: %v", err)
	}

	// 设置收件人
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("设置收件人失败: %v", err)
	}

	// 发送邮件内容
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("开始发送数据失败: %v", err)
	}

	_, err = writer.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("写入邮件内容失败: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("完成邮件发送失败: %v", err)
	}

	return nil
}

func buildMessage(from, to, subject, body string) string {
	var msg strings.Builder

	msg.WriteString(fmt.Sprintf("From: %s\r\n", from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	msg.WriteString("Content-Transfer-Encoding: 8bit\r\n")
	msg.WriteString(fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z)))
	msg.WriteString("\r\n")
	msg.WriteString(body)

	return msg.String()
}
