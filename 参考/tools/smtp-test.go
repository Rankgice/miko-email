package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"
)

// SMTPTestConfig SMTP测试配置
type SMTPTestConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	TLS      bool
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("SMTP连接测试工具")
		fmt.Println("用法: go run smtp-test.go <host:port> [user] [password] [tls]")
		fmt.Println("示例: go run smtp-test.go smtp.gmail.com:587 user@gmail.com password true")
		fmt.Println("示例: go run smtp-test.go smtp.gmail.com:465 user@gmail.com password true")
		os.Exit(1)
	}

	// 解析参数
	hostPort := os.Args[1]
	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 {
		log.Fatal("主机端口格式错误，应为 host:port")
	}

	host := parts[0]
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Fatal("端口号格式错误:", err)
	}

	config := &SMTPTestConfig{
		Host: host,
		Port: port,
		TLS:  false,
	}

	if len(os.Args) > 2 {
		config.User = os.Args[2]
	}
	if len(os.Args) > 3 {
		config.Password = os.Args[3]
	}
	if len(os.Args) > 4 {
		config.TLS, _ = strconv.ParseBool(os.Args[4])
	}

	fmt.Printf("测试SMTP连接: %s:%d\n", config.Host, config.Port)
	fmt.Printf("用户名: %s\n", config.User)
	fmt.Printf("TLS: %v\n", config.TLS)
	fmt.Println(strings.Repeat("-", 50))

	// 执行测试
	testSMTPConnection(config)
}

func testSMTPConnection(config *SMTPTestConfig) {
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	// 1. 测试基本TCP连接
	fmt.Println("1. 测试TCP连接...")
	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		fmt.Printf("❌ TCP连接失败: %v\n", err)
		return
	}
	fmt.Println("✅ TCP连接成功")
	conn.Close()

	// 2. 根据端口和TLS设置选择连接方式
	if config.TLS && config.Port == 465 {
		// SSL直接连接
		fmt.Println("2. 测试SSL直接连接...")
		testSSLConnection(config)
	} else if config.TLS && config.Port == 587 {
		// STARTTLS连接
		fmt.Println("2. 测试STARTTLS连接...")
		testSTARTTLSConnection(config)
	} else {
		// 普通连接
		fmt.Println("2. 测试普通SMTP连接...")
		testPlainConnection(config)
	}
}

func testSSLConnection(config *SMTPTestConfig) {
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	// 创建TLS配置
	tlsConfig := &tls.Config{
		ServerName:         config.Host,
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS12,
		MaxVersion:         tls.VersionTLS13,
	}

	fmt.Printf("   尝试SSL连接到 %s...\n", addr)
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		fmt.Printf("❌ SSL连接失败: %v\n", err)
		return
	}
	defer conn.Close()
	fmt.Println("✅ SSL连接成功")

	// 创建SMTP客户端
	fmt.Println("3. 创建SMTP客户端...")
	client, err := smtp.NewClient(conn, config.Host)
	if err != nil {
		fmt.Printf("❌ 创建SMTP客户端失败: %v\n", err)
		return
	}
	defer client.Close()
	fmt.Println("✅ SMTP客户端创建成功")

	// 测试认证
	if config.User != "" && config.Password != "" {
		testAuthentication(client, config)
	}

	fmt.Println("✅ SSL连接测试完成")
}

func testSTARTTLSConnection(config *SMTPTestConfig) {
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	fmt.Printf("   尝试连接到 %s...\n", addr)
	client, err := smtp.Dial(addr)
	if err != nil {
		fmt.Printf("❌ SMTP连接失败: %v\n", err)
		return
	}
	defer client.Close()
	fmt.Println("✅ SMTP连接成功")

	// 启动TLS
	fmt.Println("3. 启动STARTTLS...")
	tlsConfig := &tls.Config{
		ServerName:         config.Host,
		InsecureSkipVerify: false,
		MinVersion:         tls.VersionTLS12,
		MaxVersion:         tls.VersionTLS13,
	}

	if err = client.StartTLS(tlsConfig); err != nil {
		fmt.Printf("❌ STARTTLS失败: %v\n", err)
		return
	}
	fmt.Println("✅ STARTTLS成功")

	// 测试认证
	if config.User != "" && config.Password != "" {
		testAuthentication(client, config)
	}

	fmt.Println("✅ STARTTLS连接测试完成")
}

func testPlainConnection(config *SMTPTestConfig) {
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	fmt.Printf("   尝试连接到 %s...\n", addr)
	client, err := smtp.Dial(addr)
	if err != nil {
		fmt.Printf("❌ SMTP连接失败: %v\n", err)
		return
	}
	defer client.Close()
	fmt.Println("✅ SMTP连接成功")

	// 测试认证
	if config.User != "" && config.Password != "" {
		testAuthentication(client, config)
	}

	fmt.Println("✅ 普通连接测试完成")
}

func testAuthentication(client *smtp.Client, config *SMTPTestConfig) {
	fmt.Println("4. 测试SMTP认证...")
	auth := smtp.PlainAuth("", config.User, config.Password, config.Host)
	
	if err := client.Auth(auth); err != nil {
		fmt.Printf("❌ SMTP认证失败: %v\n", err)
		return
	}
	fmt.Println("✅ SMTP认证成功")
}
