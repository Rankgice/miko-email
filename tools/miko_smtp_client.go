package main

import (
	"fmt"
	"net/smtp"
	"strings"
	"time"
)

// Miko邮箱SMTP配置
const (
	SMTPHost = "118.120.221.169" // 你的SMTP服务器地址
	SMTPPort = "25"              // SMTP端口
)

// 邮件配置
type EmailConfig struct {
	From     string
	Password string
	To       []string
	Subject  string
	Body     string
}

// 发送邮件函数
func sendEmail(config EmailConfig) error {
	// 连接到SMTP服务器
	addr := SMTPHost + ":" + SMTPPort

	// 创建认证
	auth := smtp.PlainAuth("", config.From, config.Password, SMTPHost)

	// 构建邮件内容
	msg := buildMessage(config)

	// 发送邮件
	err := smtp.SendMail(addr, auth, config.From, config.To, []byte(msg))
	if err != nil {
		return fmt.Errorf("发送邮件失败: %v", err)
	}

	return nil
}

// 构建邮件消息
func buildMessage(config EmailConfig) string {
	msg := fmt.Sprintf("From: %s\r\n", config.From)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(config.To, ","))
	msg += fmt.Sprintf("Subject: %s\r\n", config.Subject)
	msg += "MIME-Version: 1.0\r\n"
	msg += "Content-Type: text/plain; charset=UTF-8\r\n"
	msg += fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z))
	msg += "\r\n"
	msg += config.Body
	return msg
}

// 测试不同端口的连接
func testSMTPConnection(port string) {
	fmt.Printf("\n=== 测试SMTP端口 %s ===\n", port)

	addr := SMTPHost + ":" + port

	// 尝试连接
	conn, err := smtp.Dial(addr)
	if err != nil {
		fmt.Printf("❌ 连接端口 %s 失败: %v\n", port, err)
		return
	}
	defer conn.Quit()

	fmt.Printf("✅ 成功连接到端口 %s\n", port)

	// 测试EHLO命令
	err = conn.Hello("test-client")
	if err != nil {
		fmt.Printf("❌ EHLO命令失败: %v\n", err)
		return
	}

	fmt.Printf("✅ EHLO命令成功\n")
}

func main() {
	fmt.Println("=== Miko邮箱SMTP测试工具 ===")
	fmt.Printf("测试服务器: %s\n", SMTPHost)
	fmt.Printf("测试时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))

	// 测试各个端口的连接
	fmt.Println("\n=== 端口连接测试 ===")
	ports := []string{"25", "587", "465"}
	for _, port := range ports {
		testSMTPConnection(port)
	}

	// 邮件发送测试
	fmt.Println("\n=== 邮件发送测试 ===")

	// 配置测试邮件 - 请根据实际情况修改这些参数
	config := EmailConfig{
		From:     "test@jbjj.site",           // 请替换为实际的发件邮箱
		Password: "test123",                  // 请替换为实际的邮箱密码
		To:       []string{"kimi@jbjj.site"}, // 请替换为实际的收件邮箱
		Subject:  "SMTP测试邮件 - " + time.Now().Format("2006-01-02 15:04:05"),
		Body: `这是一封SMTP测试邮件。

发送时间: ` + time.Now().Format("2006-01-02 15:04:05") + `
服务器: ` + SMTPHost + `
端口: ` + SMTPPort + `

如果您收到这封邮件，说明SMTP服务器工作正常！

测试内容:
- 连接测试: ✅
- 认证测试: ✅  
- 发送测试: ✅

---
Miko邮箱系统 SMTP测试工具`,
	}

	fmt.Printf("发件人: %s\n", config.From)
	fmt.Printf("收件人: %s\n", strings.Join(config.To, ", "))
	fmt.Printf("主题: %s\n", config.Subject)

	// 尝试发送邮件
	fmt.Println("\n开始发送测试邮件...")
	err := sendEmail(config)
	if err != nil {
		fmt.Printf("❌ 邮件发送失败: %v\n", err)
		fmt.Println("\n可能的原因:")
		fmt.Println("1. 邮箱账号或密码不正确")
		fmt.Println("2. SMTP服务器配置问题")
		fmt.Println("3. 网络连接问题")
		fmt.Println("4. 邮箱不存在")
	} else {
		fmt.Println("✅ 邮件发送成功！")
		fmt.Println("\n请检查收件箱是否收到测试邮件")
		fmt.Println("如果没有收到，请检查垃圾邮件文件夹")
	}

	fmt.Println("\n=== 测试完成 ===")
	fmt.Println("注意: 请确保在发送测试前:")
	fmt.Println("1. 修改代码中的邮箱账号和密码")
	fmt.Println("2. 确保邮箱账号在Miko邮箱系统中存在")
	fmt.Println("3. 确保SMTP服务器正在运行")
}
