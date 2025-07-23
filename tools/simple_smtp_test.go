package main

import (
	"fmt"
	"net/smtp"
	"time"
)

func main() {
	fmt.Println("=== Miko邮箱SMTP发送测试 ===")

	// SMTP服务器配置
	smtpHost := "118.120.221.169"
	smtpPort := "25"

	// 邮箱配置 - 使用你的邮箱账号
	from := "kimi@jbjj.site"
	password := "06c3c4d1"           // 请替换为实际密码
	to := []string{"kimi@jbjj.site"} // 发送给自己测试

	// 邮件内容
	subject := "SMTP测试邮件"
	body := fmt.Sprintf(`这是一封SMTP测试邮件

发送时间: %s
发件人: %s
收件人: %s
SMTP服务器: %s:%s

如果收到这封邮件，说明SMTP发送功能正常！

---
Miko邮箱系统`, time.Now().Format("2006-01-02 15:04:05"), from, to[0], smtpHost, smtpPort)

	// 构建邮件消息
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, to[0], subject, body)

	// SMTP服务器地址
	addr := smtpHost + ":" + smtpPort

	fmt.Printf("发件人: %s\n", from)
	fmt.Printf("收件人: %s\n", to[0])
	fmt.Printf("SMTP服务器: %s\n", addr)
	fmt.Printf("主题: %s\n", subject)

	// 创建认证（如果需要的话）
	// 注意：有些SMTP服务器可能不需要认证，或者认证方式不同
	var auth smtp.Auth
	if password != "你的邮箱密码" && password != "" {
		auth = smtp.PlainAuth("", from, password, smtpHost)
		fmt.Println("使用认证发送...")
	} else {
		fmt.Println("无认证发送...")
	}

	// 发送邮件
	fmt.Println("开始发送邮件...")
	err := smtp.SendMail(addr, auth, from, to, []byte(message))

	if err != nil {
		fmt.Printf("❌ 发送失败: %v\n", err)

		// 尝试不同的方法
		fmt.Println("\n尝试无认证发送...")
		err2 := smtp.SendMail(addr, nil, from, to, []byte(message))
		if err2 != nil {
			fmt.Printf("❌ 无认证发送也失败: %v\n", err2)
		} else {
			fmt.Println("✅ 无认证发送成功！")
		}
	} else {
		fmt.Println("✅ 发送成功！")
	}

	fmt.Println("\n=== 测试完成 ===")
	fmt.Println("请检查邮箱收件箱是否收到测试邮件")
}
