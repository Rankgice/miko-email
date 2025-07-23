package main

import (
	"fmt"
	"strings"

	"miko-email/internal/config"
)

func main() {
	fmt.Println("🎉 Miko邮箱系统域名限制测试总结")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("")

	// 加载配置
	config.Load()

	// 1. 检查配置状态
	fmt.Println("📋 1. 配置检查结果")
	fmt.Println(strings.Repeat("-", 30))

	if config.GlobalYAMLConfig == nil {
		fmt.Println("❌ 未找到config.yaml文件")
		fmt.Println("💡 系统使用默认配置（不限制域名）")
	} else {
		yamlCfg := config.GlobalYAMLConfig
		fmt.Printf("✅ 配置文件: config.yaml\n")
		fmt.Printf("📧 默认域名: %s\n", yamlCfg.Domain.Default)
		fmt.Printf("🔒 域名限制: %v\n", yamlCfg.Domain.EnableDomainRestriction)

		if yamlCfg.Domain.EnableDomainRestriction {
			if len(yamlCfg.Domain.Allowed) == 0 {
				fmt.Printf("📝 允许域名: 无限制 (空列表)\n")
				fmt.Printf("✅ 实际效果: 接受所有域名\n")
			} else {
				fmt.Printf("📝 允许域名: %v\n", yamlCfg.Domain.Allowed)
				fmt.Printf("⚠️  实际效果: 仅接受指定域名\n")
			}
		} else {
			fmt.Printf("📝 允许域名: 不限制\n")
			fmt.Printf("✅ 实际效果: 接受所有域名\n")
		}
	}

	fmt.Println("")

	// 2. 域名验证测试
	fmt.Println("🧪 2. 域名验证测试结果")
	fmt.Println(strings.Repeat("-", 30))

	testDomains := []string{
		"jbjj.site", "gmail.com", "yahoo.com", "example.org", "test.local",
	}

	validCount := 0
	for _, domain := range testDomains {
		isValid := config.IsValidDomain(domain)
		status := "❌"
		if isValid {
			status = "✅"
			validCount++
		}
		fmt.Printf("%s %s\n", status, domain)
	}

	fmt.Printf("\n📊 验证结果: %d/%d 域名被接受\n", validCount, len(testDomains))

	if validCount == len(testDomains) {
		fmt.Println("🎉 所有测试域名都被接受！")
	} else {
		fmt.Println("⚠️  部分域名被拒绝，存在限制")
	}

	fmt.Println("")

	// 3. 服务器状态
	fmt.Println("🚀 3. 服务器配置状态")
	fmt.Println(strings.Repeat("-", 30))

	if config.GlobalYAMLConfig != nil {
		yamlCfg := config.GlobalYAMLConfig
		fmt.Printf("🌐 Web端口: %d\n", yamlCfg.Server.WebPort)
		fmt.Printf("📧 SMTP多端口: %v\n", yamlCfg.Server.SMTP.EnableMultiPort)
		if yamlCfg.Server.SMTP.EnableMultiPort {
			fmt.Printf("📮 SMTP端口: %d, %d, %d\n",
				yamlCfg.Server.SMTP.Port25,
				yamlCfg.Server.SMTP.Port587,
				yamlCfg.Server.SMTP.Port465)
		} else {
			fmt.Printf("📮 SMTP端口: %d\n", yamlCfg.Server.SMTP.Port25)
		}
		fmt.Printf("📬 IMAP端口: %d\n", yamlCfg.Server.IMAP.Port)
		fmt.Printf("📪 POP3端口: %d\n", yamlCfg.Server.POP3.Port)
	}

	fmt.Println("")

	// 4. 管理员信息
	fmt.Println("👤 4. 管理员配置状态")
	fmt.Println(strings.Repeat("-", 30))

	username, password, email, enabled := config.GetAdminCredentials()
	fmt.Printf("👤 用户名: %s\n", username)
	fmt.Printf("📧 邮箱: %s\n", email)
	fmt.Printf("✅ 启用: %v\n", enabled)
	fmt.Printf("🔑 密码: %s\n", maskPassword(password))

	fmt.Println("")

	// 5. 测试建议
	fmt.Println("🔍 5. 测试建议")
	fmt.Println(strings.Repeat("-", 30))

	fmt.Println("📧 邮件发送测试:")
	fmt.Println("   python tools/test_domain_restriction.py multi")
	fmt.Println("")
	fmt.Println("🌐 Web界面测试:")
	fmt.Println("   http://localhost:8080/inbox")
	fmt.Println("")
	fmt.Println("🔧 管理后台测试:")
	fmt.Println("   http://localhost:8080/admin/login")
	fmt.Printf("   用户名: %s, 密码: %s\n", username, password)

	fmt.Println("")

	// 6. 总结
	fmt.Println("📝 6. 测试总结")
	fmt.Println(strings.Repeat("-", 30))

	if validCount == len(testDomains) {
		fmt.Println("🎉 域名限制测试: 通过")
		fmt.Println("✅ 系统可以接受任意域名的邮件")
		fmt.Println("🚀 配置状态: 正常")

		if config.GlobalYAMLConfig != nil && config.GlobalYAMLConfig.Server.SMTP.EnableMultiPort {
			fmt.Println("📡 多端口SMTP: 已启用")
		}

		fmt.Println("")
		fmt.Println("🎯 下一步:")
		fmt.Println("   1. 发送实际测试邮件验证接收功能")
		fmt.Println("   2. 测试不同域名的邮件接收")
		fmt.Println("   3. 验证邮件解码功能是否正常")

	} else {
		fmt.Println("⚠️  域名限制测试: 部分通过")
		fmt.Println("💡 建议检查config.yaml中的域名配置")
		fmt.Println("🔧 确保 enable_domain_restriction: false")
	}

	fmt.Println("")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("🎉 测试总结完成！")
}

func maskPassword(password string) string {
	if len(password) <= 4 {
		return "****"
	}
	return password[:2] + "****" + password[len(password)-2:]
}
