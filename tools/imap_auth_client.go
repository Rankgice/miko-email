package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func main() {
	fmt.Println("=== Miko邮箱IMAP认证测试 ===")

	// IMAP服务器配置
	imapHost := "118.120.221.169"
	imapPort := "143"

	fmt.Printf("服务器: %s:%s\n", imapHost, imapPort)

	// 测试不同的认证方式
	testCases := []struct {
		name     string
		username string
		password string
		desc     string
	}{
		{
			name:     "方式1-邮箱直接认证",
			username: "kimi@jbjj.site",
			password: "06c3c4d1",
			desc:     "使用邮箱地址和邮箱密码",
		},
		{
			name:     "方式2-组合认证",
			username: "kimi@kimi@jbjj.site", // 网站用户@邮箱地址
			password: "06c3c4d1",            // 邮箱密码
			desc:     "使用网站用户名@邮箱地址和邮箱密码",
		},
		{
			name:     "方式3-网站用户认证",
			username: "kimi",      // 网站用户名
			password: "tgx123456", // 网站用户密码
			desc:     "使用网站用户名和网站密码",
		},
	}

	for i, testCase := range testCases {
		fmt.Printf("\n=== 测试 %d: %s ===\n", i+1, testCase.name)
		fmt.Printf("描述: %s\n", testCase.desc)
		fmt.Printf("用户名: %s\n", testCase.username)
		fmt.Printf("密码: %s\n", testCase.password)

		success := testIMAPLogin(imapHost, imapPort, testCase.username, testCase.password)
		if success {
			fmt.Printf("✅ %s 认证成功！\n", testCase.name)
		} else {
			fmt.Printf("❌ %s 认证失败\n", testCase.name)
		}
	}

	fmt.Println("\n=== 测试完成 ===")
}

func testIMAPLogin(host, port, username, password string) bool {
	// 连接到IMAP服务器
	addr := host + ":" + port
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("❌ 连接失败: %v\n", err)
		return false
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// 读取欢迎消息
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("❌ 读取欢迎消息失败: %v\n", err)
		return false
	}
	fmt.Printf("服务器欢迎: %s", response)

	// 发送LOGIN命令
	loginCmd := fmt.Sprintf("A001 LOGIN \"%s\" \"%s\"\r\n", username, password)
	_, err = writer.WriteString(loginCmd)
	if err != nil {
		fmt.Printf("❌ 发送LOGIN命令失败: %v\n", err)
		return false
	}
	writer.Flush()

	// 读取LOGIN响应
	response, err = reader.ReadString('\n')
	if err != nil {
		fmt.Printf("❌ 读取LOGIN响应失败: %v\n", err)
		return false
	}
	fmt.Printf("LOGIN响应: %s", response)

	// 检查是否登录成功
	if strings.Contains(response, "A001 OK") {
		// 发送LOGOUT命令
		_, err = writer.WriteString("A002 LOGOUT\r\n")
		if err == nil {
			writer.Flush()
			// 读取LOGOUT响应
			reader.ReadString('\n')
		}
		return true
	}

	return false
}
