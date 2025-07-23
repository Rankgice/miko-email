package main

import (
	"fmt"
	"net/smtp"
	"time"
)

func main() {
	fmt.Println("=== Miko邮箱外部发送测试 ===")

	// SMTP服务器配置
	smtpHost := "118.120.221.169"
	smtpPort := "25"

	// 邮箱配置 - 使用你的域名邮箱发送到外部邮箱
	from := "kimi@jbjj.site"            // 你的域名邮箱
	to := []string{"2014131458@qq.com"} // 外部邮箱（QQ邮箱）

	// 邮件内容
	subject := "来自Miko邮箱的外部发送测试"
	body := fmt.Sprintf(`这是一封来自Miko邮箱系统的外部发送测试邮件

发送时间: %s
发件人: %s (Miko邮箱系统)
收件人: %s
SMTP服务器: %s:%s

测试内容:
✅ 使用自己域名邮箱作为发件人
✅ 发送到外部邮箱（QQ邮箱）
✅ 通过MX记录直接发送

如果您收到这封邮件，说明Miko邮箱系统的外部发送功能工作正常！

您可以直接回复此邮件到 %s

---
Miko邮箱系统 - 外部发送测试
服务器: %s`,
		time.Now().Format("2006-01-02 15:04:05"),
		from, to[0], smtpHost, smtpPort, from, smtpHost)

	// 构建邮件消息
	message := fmt.Sprintf(`From: %s
To: %s
Subject: %s
MIME-Version: 1.0
Content-Type: text/plain; charset=UTF-8
Date: %s

%s`, from, to[0], subject, time.Now().Format(time.RFC1123Z), body)

	// SMTP服务器地址
	addr := smtpHost + ":" + smtpPort

	fmt.Printf("发件人: %s\n", from)
	fmt.Printf("收件人: %s\n", to[0])
	fmt.Printf("SMTP服务器: %s\n", addr)
	fmt.Printf("主题: %s\n", subject)

	fmt.Println("\n开始发送外部邮件...")

	// 发送邮件（无认证，让服务器处理MX发送）
	err := smtp.SendMail(addr, nil, from, to, []byte(message))

	if err != nil {
		fmt.Printf("❌ 外部发送失败: %v\n", err)
		fmt.Println("\n可能的原因:")
		fmt.Println("1. 收件人邮箱服务器拒绝接收")
		fmt.Println("2. 你的域名没有正确的MX记录")
		fmt.Println("3. 收件人邮箱服务器的反垃圾邮件策略")
		fmt.Println("4. 网络连接问题")
	} else {
		fmt.Println("✅ 外部邮件发送成功！")
		fmt.Println("\n请检查收件人邮箱:")
		fmt.Printf("- 收件箱: %s\n", to[0])
		fmt.Println("- 如果没有收到，请检查垃圾邮件文件夹")
		fmt.Println("- 外部邮箱可能需要几分钟才能收到邮件")
	}

	fmt.Println("\n=== 测试完成 ===")
	fmt.Println("注意:")
	fmt.Println("1. 发件人必须是你域名下的邮箱")
	fmt.Println("2. 系统会通过MX记录直接发送到外部邮箱")
	fmt.Println("3. 不需要配置第三方SMTP服务器")
	fmt.Println("4. 收件人看到的发件人就是你的域名邮箱")
}
