package config

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

// SMTPConfig SMTP服务器配置
type SMTPConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	TLS      bool   `json:"tls"`
}

// Config 应用配置
type Config struct {
	WebPort     int    `json:"web_port"`
	SMTPPort    int    `json:"smtp_port"`     // 主SMTP端口 (默认25)
	SMTPPort587 int    `json:"smtp_port_587"` // SMTP提交端口 (587)
	SMTPPort465 int    `json:"smtp_port_465"` // SMTPS端口 (465)
	IMAPPort    int    `json:"imap_port"`
	POP3Port    int    `json:"pop3_port"`
	DBPath      string `json:"db_path"`
	Domain      string `json:"domain"`
	AdminEmail  string `json:"admin_email"`
	AdminPass   string `json:"admin_pass"`
	JWTSecret   string `json:"jwt_secret"`
	// SMTP多端口支持
	EnableMultiSMTP bool `json:"enable_multi_smtp"` // 是否启用多SMTP端口
	// 外部SMTP发送配置（向后兼容）
	OutboundSMTPHost     string `json:"outbound_smtp_host"`
	OutboundSMTPPort     int    `json:"outbound_smtp_port"`
	OutboundSMTPUser     string `json:"outbound_smtp_user"`
	OutboundSMTPPassword string `json:"outbound_smtp_password"`
	OutboundSMTPTLS      bool   `json:"outbound_smtp_tls"`
	// 多域名SMTP配置
	DomainSMTPConfigs map[string]*SMTPConfig `json:"domain_smtp_configs"`
}

// LoadFromEnv 从环境变量加载配置
func (c *Config) LoadFromEnv() {
	if port := os.Getenv("WEB_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			c.WebPort = p
		}
	}

	if port := os.Getenv("SMTP_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			c.SMTPPort = p
		}
	}

	if port := os.Getenv("SMTP_PORT_587"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			c.SMTPPort587 = p
		}
	}

	if port := os.Getenv("SMTP_PORT_465"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			c.SMTPPort465 = p
		}
	}

	if multiSMTP := os.Getenv("ENABLE_MULTI_SMTP"); multiSMTP == "true" {
		c.EnableMultiSMTP = true
	}

	if port := os.Getenv("IMAP_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			c.IMAPPort = p
		}
	}

	if port := os.Getenv("POP3_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			c.POP3Port = p
		}
	}

	if dbPath := os.Getenv("DB_PATH"); dbPath != "" {
		c.DBPath = dbPath
	}

	if domain := os.Getenv("DOMAIN"); domain != "" {
		c.Domain = domain
	} else {
		c.Domain = "localhost"
	}

	if adminEmail := os.Getenv("ADMIN_EMAIL"); adminEmail != "" {
		c.AdminEmail = adminEmail
	} else if c.AdminEmail == "" {
		c.AdminEmail = "admin@" + c.Domain
	}

	if adminPass := os.Getenv("ADMIN_PASS"); adminPass != "" {
		c.AdminPass = adminPass
	} else if c.AdminPass == "" {
		c.AdminPass = "admin123"
	}

	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		c.JWTSecret = jwtSecret
	} else {
		c.JWTSecret = "nbemail-secret-key-2024"
	}

	// 外部SMTP配置
	if host := os.Getenv("OUTBOUND_SMTP_HOST"); host != "" {
		c.OutboundSMTPHost = host
	}

	if port := os.Getenv("OUTBOUND_SMTP_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			c.OutboundSMTPPort = p
		}
	}

	if user := os.Getenv("OUTBOUND_SMTP_USER"); user != "" {
		c.OutboundSMTPUser = user
	}

	if password := os.Getenv("OUTBOUND_SMTP_PASSWORD"); password != "" {
		c.OutboundSMTPPassword = password
	}

	if tls := os.Getenv("OUTBOUND_SMTP_TLS"); tls == "true" {
		c.OutboundSMTPTLS = true
	}
}

// GetDefaults 获取默认配置
func GetDefaults() *Config {
	cfg := &Config{
		WebPort:         8080, // Web界面端口
		SMTPPort:        25,   // 主SMTP服务器端口
		SMTPPort587:     587,  // SMTP提交端口
		SMTPPort465:     465,  // SMTPS端口
		IMAPPort:        143,  // IMAP服务器端口
		POP3Port:        110,  // POP3服务器端口
		DBPath:          "nbemail.db",
		Domain:          "localhost",
		AdminEmail:      "2014131458@qq.com",
		AdminPass:       "tgx123456",
		JWTSecret:       "nbemail-secret-key-2024",
		EnableMultiSMTP: true, // 默认启用多SMTP端口
		// 外部SMTP默认配置（可通过环境变量覆盖）
		OutboundSMTPHost:     "",
		OutboundSMTPPort:     587,
		OutboundSMTPUser:     "",
		OutboundSMTPPassword: "",
		OutboundSMTPTLS:      false,
		// 初始化多域名SMTP配置
		DomainSMTPConfigs: make(map[string]*SMTPConfig),
	}
	cfg.LoadFromEnv()
	cfg.LoadDomainSMTPConfigs()
	return cfg
}

// GetSMTPPorts 获取所有SMTP端口
func (c *Config) GetSMTPPorts() []int {
	if c.EnableMultiSMTP {
		return []int{c.SMTPPort, c.SMTPPort587, c.SMTPPort465}
	}
	return []int{c.SMTPPort}
}

// GetSMTPConfigForEmail 根据邮箱地址获取对应的SMTP配置
func (c *Config) GetSMTPConfigForEmail(email string) *SMTPConfig {
	if email == "" {
		return nil
	}

	// 提取域名
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return nil
	}
	domain := parts[1]

	// 查找域名对应的SMTP配置
	if smtpConfig, exists := c.DomainSMTPConfigs[domain]; exists {
		return smtpConfig
	}

	// 如果没有找到域名特定配置，使用默认配置
	if c.OutboundSMTPHost != "" {
		return &SMTPConfig{
			Host:     c.OutboundSMTPHost,
			Port:     c.OutboundSMTPPort,
			User:     c.OutboundSMTPUser,
			Password: c.OutboundSMTPPassword,
			TLS:      c.OutboundSMTPTLS,
		}
	}

	return nil
}

// GetSMTPConfigForDomain 根据域名获取对应的SMTP配置
func (c *Config) GetSMTPConfigForDomain(domain string) *SMTPConfig {
	if domain == "" {
		return nil
	}

	domain = strings.ToLower(domain)

	// 查找域名对应的SMTP配置
	if smtpConfig, exists := c.DomainSMTPConfigs[domain]; exists {
		return smtpConfig
	}

	// 如果没有找到域名特定配置，使用默认配置
	if c.OutboundSMTPHost != "" {
		return &SMTPConfig{
			Host:     c.OutboundSMTPHost,
			Port:     c.OutboundSMTPPort,
			User:     c.OutboundSMTPUser,
			Password: c.OutboundSMTPPassword,
			TLS:      c.OutboundSMTPTLS,
		}
	}

	return nil
}

// LoadDomainSMTPConfigs 从环境变量或配置文件加载域名SMTP配置
func (c *Config) LoadDomainSMTPConfigs() {
	// 可以从环境变量加载多个域名配置
	// 格式: DOMAIN_SMTP_<DOMAIN>=host:port:user:password:tls
	// 例如: DOMAIN_SMTP_EXAMPLE_COM=smtp.example.com:587:user@example.com:password:true

	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "DOMAIN_SMTP_") {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) != 2 {
				continue
			}

			// 提取域名
			envKey := parts[0]
			domain := strings.ToLower(strings.Replace(envKey[12:], "_", ".", -1)) // 去掉 DOMAIN_SMTP_ 前缀

			// 解析配置
			configParts := strings.Split(parts[1], ":")
			if len(configParts) != 5 {
				continue
			}

			port, err := strconv.Atoi(configParts[1])
			if err != nil {
				continue
			}

			tls, err := strconv.ParseBool(configParts[4])
			if err != nil {
				tls = true // 默认启用TLS
			}

			c.DomainSMTPConfigs[domain] = &SMTPConfig{
				Host:     configParts[0],
				Port:     port,
				User:     configParts[2],
				Password: configParts[3],
				TLS:      tls,
			}
		}
	}
}

// AddDomainSMTPConfig 添加域名SMTP配置
func (c *Config) AddDomainSMTPConfig(domain string, config *SMTPConfig) {
	if c.DomainSMTPConfigs == nil {
		c.DomainSMTPConfigs = make(map[string]*SMTPConfig)
	}
	c.DomainSMTPConfigs[strings.ToLower(domain)] = config
}

// RemoveDomainSMTPConfig 删除域名SMTP配置
func (c *Config) RemoveDomainSMTPConfig(domain string) {
	if c.DomainSMTPConfigs != nil {
		delete(c.DomainSMTPConfigs, strings.ToLower(domain))
	}
}

// AutoConfigureDomainSMTP 自动配置域名SMTP
func (c *Config) AutoConfigureDomainSMTP(domains []string) {
	if c.DomainSMTPConfigs == nil {
		c.DomainSMTPConfigs = make(map[string]*SMTPConfig)
	}

	for _, domain := range domains {
		domain = strings.ToLower(domain)

		// 如果已经有配置，跳过
		if _, exists := c.DomainSMTPConfigs[domain]; exists {
			continue
		}

		// 生成常见的SMTP配置
		smtpConfig := c.generateSMTPConfigForDomain(domain)
		if smtpConfig != nil {
			c.DomainSMTPConfigs[domain] = smtpConfig
		}
	}
}

// generateSMTPConfigForDomain 为域名生成SMTP配置
func (c *Config) generateSMTPConfigForDomain(domain string) *SMTPConfig {
	// 常见的SMTP服务器模式（优先使用25端口）
	smtpPatterns := []struct {
		pattern string
		port    int
		tls     bool
	}{
		{"mail." + domain, 25, false},  // mail.domain.com (标准SMTP端口)
		{"smtp." + domain, 25, false},  // smtp.domain.com (标准SMTP端口)
		{"mx." + domain, 25, false},    // mx.domain.com (标准SMTP端口)
		{"email." + domain, 25, false}, // email.domain.com (标准SMTP端口)
	}

	// 生成推荐的用户名和密码
	suggestedUser := fmt.Sprintf("smtp@%s", domain)
	suggestedPassword := c.generateStrongPassword()

	// 尝试每种模式
	for _, pattern := range smtpPatterns {
		// 检查DNS记录是否存在
		if c.checkSMTPServer(pattern.pattern, pattern.port) {
			return &SMTPConfig{
				Host:     pattern.pattern,
				Port:     pattern.port,
				User:     suggestedUser,
				Password: suggestedPassword,
				TLS:      pattern.tls,
			}
		}
	}

	// 如果没有找到，返回默认配置
	return &SMTPConfig{
		Host:     "mail." + domain,
		Port:     25,
		User:     suggestedUser,
		Password: suggestedPassword,
		TLS:      false,
	}
}

// generateStrongPassword 生成强密码
func (c *Config) generateStrongPassword() string {
	const (
		uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		lowercase = "abcdefghijklmnopqrstuvwxyz"
		numbers   = "0123456789"
		symbols   = "!@#$%^&*()_+-=[]{}|;:,.<>?"
		length    = 16
	)

	allChars := uppercase + lowercase + numbers + symbols
	password := make([]byte, length)

	// 确保至少包含每种类型的字符
	charSets := []string{uppercase, lowercase, numbers, symbols}
	for i, charset := range charSets {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		password[i] = charset[n.Int64()]
	}

	// 填充剩余长度
	for i := 4; i < length; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(allChars))))
		password[i] = allChars[n.Int64()]
	}

	// 打乱字符顺序
	for i := range password {
		j, _ := rand.Int(rand.Reader, big.NewInt(int64(len(password))))
		password[i], password[j.Int64()] = password[j.Int64()], password[i]
	}

	return string(password)
}

// checkSMTPServer 检查SMTP服务器是否可用
func (c *Config) checkSMTPServer(host string, port int) bool {
	// 简单的连接测试
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), 5*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// GetDomainsFromMailboxes 从邮箱地址中提取域名列表
func (c *Config) GetDomainsFromMailboxes(mailboxes []string) []string {
	domainSet := make(map[string]bool)
	var domains []string

	// 定义要过滤的测试域名（只过滤明显的测试域名）
	testDomains := map[string]bool{
		"localhost":     true,
		"test.local":    true,
		"example.test":  true,
		"localhost.com": true,
	}

	for _, mailbox := range mailboxes {
		if strings.Contains(mailbox, "@") {
			parts := strings.Split(mailbox, "@")
			if len(parts) == 2 {
				domain := strings.ToLower(parts[1])
				// 过滤掉测试域名和已存在的域名
				if !domainSet[domain] && !testDomains[domain] {
					domainSet[domain] = true
					domains = append(domains, domain)
				}
			}
		}
	}

	return domains
}
