package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"net"
	"strings"
	"time"
)

func main() {
	fmt.Println("=== Miko邮箱SMTP认证测试 ===")

	// SMTP服务器配置
	smtpHost := "118.120.221.169"
	smtpPort := "25"

	// 邮箱认证信息
	username := "kimi@jbjj.site" // 邮箱地址
	password := "06c3c4d1"       // 邮箱密码
	from := "kimi@jbjj.site"     // 发件人邮箱
	to := "2014131458@qq.com"    // 外部收件人

	fmt.Printf("服务器: %s:%s\n", smtpHost, smtpPort)
	fmt.Printf("认证用户: %s\n", username)
	fmt.Printf("发件人: %s\n", from)
	fmt.Printf("收件人: %s\n", to)

	// 连接到SMTP服务器
	addr := smtpHost + ":" + smtpPort
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("❌ 连接失败: %v\n", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// 读取欢迎消息
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("❌ 读取欢迎消息失败: %v\n", err)
		return
	}
	fmt.Printf("服务器欢迎: %s", response)

	// 发送EHLO命令
	fmt.Println("\n1. 发送EHLO命令...")
	_, err = writer.WriteString("EHLO test-client\r\n")
	if err != nil {
		fmt.Printf("❌ 发送EHLO失败: %v\n", err)
		return
	}
	writer.Flush()

	// 读取EHLO响应（多行）
	fmt.Println("EHLO响应:")
	for {
		response, err = reader.ReadString('\n')
		if err != nil {
			fmt.Printf("❌ 读取EHLO响应失败: %v\n", err)
			return
		}
		fmt.Printf("  %s", response)
		// 如果响应以"250 "开头（而不是"250-"），表示这是最后一行
		if strings.HasPrefix(response, "250 ") {
			break
		}
	}

	// 发送AUTH PLAIN命令
	fmt.Println("\n2. 发送AUTH PLAIN认证...")

	// 构建PLAIN认证字符串: \0username\0password
	authString := fmt.Sprintf("\x00%s\x00%s", username, password)
	authEncoded := base64.StdEncoding.EncodeToString([]byte(authString))

	_, err = writer.WriteString(fmt.Sprintf("AUTH PLAIN %s\r\n", authEncoded))
	if err != nil {
		fmt.Printf("❌ 发送AUTH命令失败: %v\n", err)
		return
	}
	writer.Flush()

	// 读取认证响应
	response, err = reader.ReadString('\n')
	if err != nil {
		fmt.Printf("❌ 读取认证响应失败: %v\n", err)
		return
	}
	fmt.Printf("认证响应: %s", response)

	if !strings.HasPrefix(response, "235") {
		fmt.Println("❌ 认证失败")
		return
	}
	fmt.Println("✅ 认证成功")

	// 发送MAIL FROM命令
	fmt.Println("\n3. 发送MAIL FROM命令...")
	_, err = writer.WriteString(fmt.Sprintf("MAIL FROM:<%s>\r\n", from))
	if err != nil {
		fmt.Printf("❌ 发送MAIL FROM失败: %v\n", err)
		return
	}
	writer.Flush()

	response, err = reader.ReadString('\n')
	if err != nil {
		fmt.Printf("❌ 读取MAIL FROM响应失败: %v\n", err)
		return
	}
	fmt.Printf("MAIL FROM响应: %s", response)

	if !strings.HasPrefix(response, "250") {
		fmt.Println("❌ MAIL FROM失败")
		return
	}

	// 发送RCPT TO命令（外部邮箱）
	fmt.Println("\n4. 发送RCPT TO命令（外部邮箱）...")
	_, err = writer.WriteString(fmt.Sprintf("RCPT TO:<%s>\r\n", to))
	if err != nil {
		fmt.Printf("❌ 发送RCPT TO失败: %v\n", err)
		return
	}
	writer.Flush()

	response, err = reader.ReadString('\n')
	if err != nil {
		fmt.Printf("❌ 读取RCPT TO响应失败: %v\n", err)
		return
	}
	fmt.Printf("RCPT TO响应: %s", response)

	if !strings.HasPrefix(response, "250") {
		fmt.Printf("❌ RCPT TO失败: %s", response)
		return
	}
	fmt.Println("✅ 外部邮箱收件人验证成功")

	// 发送DATA命令
	fmt.Println("\n5. 发送DATA命令...")
	_, err = writer.WriteString("DATA\r\n")
	if err != nil {
		fmt.Printf("❌ 发送DATA失败: %v\n", err)
		return
	}
	writer.Flush()

	response, err = reader.ReadString('\n')
	if err != nil {
		fmt.Printf("❌ 读取DATA响应失败: %v\n", err)
		return
	}
	fmt.Printf("DATA响应: %s", response)

	if !strings.HasPrefix(response, "354") {
		fmt.Println("❌ DATA命令失败")
		return
	}

	// 发送邮件内容
	fmt.Println("\n6. 发送邮件内容...")
	emailContent := fmt.Sprintf(`From: %s
To: %s
Subject: Miko邮箱SMTP认证测试邮件
MIME-Version: 1.0
Content-Type: text/plain; charset=UTF-8
Date: %s

这是一封通过SMTP认证发送的测试邮件。

发送时间: %s
发件人: %s (已认证)
收件人: %s
认证用户: %s

如果您收到这封邮件，说明Miko邮箱的SMTP认证功能工作正常！

---
Miko邮箱系统 - SMTP认证测试
`, from, to, time.Now().Format(time.RFC1123Z), time.Now().Format("2006-01-02 15:04:05"), from, to, username)

	_, err = writer.WriteString(emailContent)
	if err != nil {
		fmt.Printf("❌ 发送邮件内容失败: %v\n", err)
		return
	}

	// 发送结束标记
	_, err = writer.WriteString("\r\n.\r\n")
	if err != nil {
		fmt.Printf("❌ 发送结束标记失败: %v\n", err)
		return
	}
	writer.Flush()

	// 读取最终响应
	response, err = reader.ReadString('\n')
	if err != nil {
		fmt.Printf("❌ 读取最终响应失败: %v\n", err)
		return
	}
	fmt.Printf("最终响应: %s", response)

	if strings.HasPrefix(response, "250") {
		fmt.Println("✅ 邮件发送成功！")
		fmt.Println("\n测试结果:")
		fmt.Println("✅ SMTP连接成功")
		fmt.Println("✅ 用户认证成功")
		fmt.Println("✅ 发件人权限验证成功")
		fmt.Println("✅ 外部邮箱发送成功")
		fmt.Printf("\n请检查 %s 的收件箱\n", to)
	} else {
		fmt.Printf("❌ 邮件发送失败: %s", response)
	}

	// 发送QUIT命令
	writer.WriteString("QUIT\r\n")
	writer.Flush()

	fmt.Println("\n=== 测试完成 ===")
}
