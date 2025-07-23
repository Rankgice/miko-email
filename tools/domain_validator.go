package main

import (
	"fmt"
	"os"

	"miko-email/internal/config"
)

func main() {
	if len(os.Args) < 2 {
		showUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "check":
		checkDomainRestriction()
	case "test":
		testDomainValidation()
	case "status":
		showDomainStatus()
	default:
		fmt.Printf("未知命令: %s\n", command)
		showUsage()
	}
}

func showUsage() {
	fmt.Println("Miko邮箱域名验证工具")
	fmt.Println("")
	fmt.Println("用法:")
	fmt.Println("  go run tools/domain_validator.go <命令>")
	fmt.Println("")
	fmt.Println("命令:")
	fmt.Println("  check  - 检查域名限制配置状态")
	fmt.Println("  test   - 测试域名验证逻辑")
	fmt.Println("  status - 显示域名配置详情")
	fmt.Println("")
	fmt.Println("示例:")
	fmt.Println("  go run tools/domain_validator.go check")
	fmt.Println("  go run tools/domain_validator.go test")
}

func checkDomainRestriction() {
	fmt.Println("=== 域名限制检查 ===")
	fmt.Println("")

	// 加载配置
	config.Load()

	if config.GlobalYAMLConfig == nil {
		fmt.Println("⚠️  未找到config.yaml文件，使用默认配置")
		fmt.Println("📝 默认情况下不限制域名")
		return
	}

	yamlCfg := config.GlobalYAMLConfig

	fmt.Println("📋 域名配置状态:")
	fmt.Printf("  默认域名: %s\n", yamlCfg.Domain.Default)
	fmt.Printf("  启用域名限制: %v\n", yamlCfg.Domain.EnableDomainRestriction)

	if yamlCfg.Domain.EnableDomainRestriction {
		if len(yamlCfg.Domain.Allowed) == 0 {
			fmt.Printf("  允许的域名: 无限制 (空列表)\n")
			fmt.Println("")
			fmt.Println("✅ 结果: 域名不受限制")
			fmt.Println("💡 虽然启用了域名限制，但允许列表为空，等同于不限制")
		} else {
			fmt.Printf("  允许的域名: %v\n", yamlCfg.Domain.Allowed)
			fmt.Println("")
			fmt.Println("⚠️  结果: 域名受到限制")
			fmt.Printf("📝 只允许以下域名的邮件: %v\n", yamlCfg.Domain.Allowed)
		}
	} else {
		fmt.Printf("  允许的域名: 不限制\n")
		fmt.Println("")
		fmt.Println("✅ 结果: 域名不受限制")
		fmt.Println("🎉 系统将接受任何域名的邮件")
	}

	fmt.Println("")
	fmt.Println("🔍 验证方法:")
	fmt.Println("  1. 运行: python tools/test_domain_restriction.py multi")
	fmt.Println("  2. 检查邮箱系统是否收到各种域名的测试邮件")
}

func testDomainValidation() {
	fmt.Println("=== 域名验证逻辑测试 ===")
	fmt.Println("")

	// 测试域名列表
	testDomains := []string{
		"jbjj.site",
		"gmail.com",
		"yahoo.com",
		"outlook.com",
		"example.org",
		"test.local",
		"demo.internal",
		"custom-domain.xyz",
		"very-long-domain-name.com",
		"sub.domain.example.net",
	}

	fmt.Println("🧪 测试域名验证逻辑:")
	fmt.Printf("%-30s %s\n", "域名", "验证结果")
	fmt.Printf("%-30s %s\n", "------------------------------", "----------")

	validCount := 0
	for _, domain := range testDomains {
		isValid := config.IsValidDomain(domain)
		status := "❌ 拒绝"
		if isValid {
			status = "✅ 接受"
			validCount++
		}
		fmt.Printf("%-30s %s\n", domain, status)
	}

	fmt.Println("")
	fmt.Printf("📊 测试结果: %d/%d 域名被接受\n", validCount, len(testDomains))

	if validCount == len(testDomains) {
		fmt.Println("🎉 所有测试域名都被接受，域名限制已取消！")
	} else if validCount == 0 {
		fmt.Println("❌ 所有测试域名都被拒绝，域名限制过于严格")
	} else {
		fmt.Println("⚠️  部分域名被拒绝，存在域名限制")
	}

	fmt.Println("")
	fmt.Println("💡 如果想要接受所有域名，请确保:")
	fmt.Println("   - enable_domain_restriction: false")
	fmt.Println("   - 或者 allowed: [] (空数组)")
}

func showDomainStatus() {
	fmt.Println("=== 域名配置详情 ===")
	fmt.Println("")

	// 加载配置
	config.Load()

	if config.GlobalYAMLConfig == nil {
		fmt.Println("❌ 未找到config.yaml文件")
		fmt.Println("💡 系统将使用默认配置（不限制域名）")
		return
	}

	yamlCfg := config.GlobalYAMLConfig

	fmt.Println("📋 完整域名配置:")
	fmt.Printf("  默认域名: %s\n", yamlCfg.Domain.Default)
	fmt.Printf("  启用域名限制: %v\n", yamlCfg.Domain.EnableDomainRestriction)
	fmt.Printf("  允许的域名数量: %d\n", len(yamlCfg.Domain.Allowed))

	if len(yamlCfg.Domain.Allowed) > 0 {
		fmt.Println("  允许的域名列表:")
		for i, domain := range yamlCfg.Domain.Allowed {
			fmt.Printf("    %d. %s\n", i+1, domain)
		}
	} else {
		fmt.Println("  允许的域名列表: 空 (不限制)")
	}

	fmt.Println("")
	fmt.Println("🔧 配置解释:")

	if !yamlCfg.Domain.EnableDomainRestriction {
		fmt.Println("  ✅ 域名限制已禁用")
		fmt.Println("  📧 系统将接受任何域名的邮件")
		fmt.Println("  🎯 这是推荐的配置")
	} else {
		if len(yamlCfg.Domain.Allowed) == 0 {
			fmt.Println("  ⚠️  域名限制已启用，但允许列表为空")
			fmt.Println("  📧 系统将接受任何域名的邮件")
			fmt.Println("  💡 建议设置 enable_domain_restriction: false")
		} else {
			fmt.Println("  ❌ 域名限制已启用且有具体限制")
			fmt.Println("  📧 系统只接受允许列表中的域名")
			fmt.Println("  💡 如需接受所有域名，请:")
			fmt.Println("     - 设置 enable_domain_restriction: false")
			fmt.Println("     - 或清空 allowed 数组")
		}
	}

	fmt.Println("")
	fmt.Println("🧪 测试建议:")
	fmt.Println("  1. 运行域名验证测试:")
	fmt.Println("     go run tools/domain_validator.go test")
	fmt.Println("  2. 发送实际测试邮件:")
	fmt.Println("     python tools/test_domain_restriction.py multi")
	fmt.Println("  3. 检查邮箱系统收件情况")

	fmt.Println("")
	fmt.Println("📝 配置文件位置: config.yaml")
	fmt.Println("🔄 修改配置后需要重启服务器")
}
