package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func main() {
	fmt.Println("=== 单个IMAP认证测试 ===")

	// 测试方式2: 组合认证
	username := "kimi@kimi@jbjj.site"
	password := "06c3c4d1"

	fmt.Printf("测试用户名: %s\n", username)
	fmt.Printf("测试密码: %s\n", password)

	// 连接到IMAP服务器
	conn, err := net.Dial("tcp", "118.120.221.169:143")
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

	// 发送LOGIN命令
	loginCmd := fmt.Sprintf("A001 LOGIN \"%s\" \"%s\"\r\n", username, password)
	fmt.Printf("发送命令: %s", loginCmd)

	_, err = writer.WriteString(loginCmd)
	if err != nil {
		fmt.Printf("❌ 发送LOGIN命令失败: %v\n", err)
		return
	}
	writer.Flush()

	// 读取LOGIN响应
	response, err = reader.ReadString('\n')
	if err != nil {
		fmt.Printf("❌ 读取LOGIN响应失败: %v\n", err)
		return
	}
	fmt.Printf("LOGIN响应: %s", response)

	// 检查是否登录成功
	if strings.Contains(response, "A001 OK") {
		fmt.Println("✅ 认证成功！")
	} else {
		fmt.Println("❌ 认证失败")
	}

	// 发送LOGOUT命令
	_, err = writer.WriteString("A002 LOGOUT\r\n")
	if err == nil {
		writer.Flush()
		// 读取LOGOUT响应
		reader.ReadString('\n')
	}
}
