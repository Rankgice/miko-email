package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strconv"
)

// YAMLConfig YAML配置文件结构
type YAMLConfig struct {
	Server struct {
		WebPort int `yaml:"web_port"`
		SMTP    struct {
			EnableMultiPort bool `yaml:"enable_multi_port"`
			Port25          int  `yaml:"port_25"`
			Port587         int  `yaml:"port_587"`
			Port465         int  `yaml:"port_465"`
		} `yaml:"smtp"`
		IMAP struct {
			Port       int `yaml:"port"`
			SecurePort int `yaml:"secure_port"`
		} `yaml:"imap"`
		POP3 struct {
			Port       int `yaml:"port"`
			SecurePort int `yaml:"secure_port"`
		} `yaml:"pop3"`
	} `yaml:"server"`

	Admin struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Email    string `yaml:"email"`
		Enabled  bool   `yaml:"enabled"`
	} `yaml:"admin"`

	Database struct {
		Path  string `yaml:"path"`
		Debug bool   `yaml:"debug"`
	} `yaml:"database"`

	Domain struct {
		Default                 string   `yaml:"default"`
		Allowed                 []string `yaml:"allowed"`
		EnableDomainRestriction bool     `yaml:"enable_domain_restriction"`
	} `yaml:"domain"`

	Security struct {
		SessionKey     string `yaml:"session_key"`
		JWTSecret      string `yaml:"jwt_secret"`
		SessionTimeout int    `yaml:"session_timeout"`
		EnableHTTPS    bool   `yaml:"enable_https"`
		SSLCert        string `yaml:"ssl_cert"`
		SSLKey         string `yaml:"ssl_key"`
	} `yaml:"security"`

	Email struct {
		MaxSize             int  `yaml:"max_size"`
		MaxMailboxesPerUser int  `yaml:"max_mailboxes_per_user"`
		RetentionDays       int  `yaml:"retention_days"`
		EnableForwarding    bool `yaml:"enable_forwarding"`
	} `yaml:"email"`

	Logging struct {
		Level     string `yaml:"level"`
		ToFile    bool   `yaml:"to_file"`
		FilePath  string `yaml:"file_path"`
		AccessLog bool   `yaml:"access_log"`
	} `yaml:"logging"`

	Performance struct {
		MaxConnections int `yaml:"max_connections"`
		ReadTimeout    int `yaml:"read_timeout"`
		WriteTimeout   int `yaml:"write_timeout"`
		IdleTimeout    int `yaml:"idle_timeout"`
	} `yaml:"performance"`

	Features struct {
		AllowRegistration bool `yaml:"allow_registration"`
		EnableSearch      bool `yaml:"enable_search"`
		EnableAttachments bool `yaml:"enable_attachments"`
		EnableSpamFilter  bool `yaml:"enable_spam_filter"`
	} `yaml:"features"`
}

type Config struct {
	WebPort         string
	SMTPPort        string
	SMTPPort587     string // SMTP提交端口
	SMTPPort465     string // SMTPS端口
	IMAPPort        string
	POP3Port        string
	DatabasePath    string
	SessionKey      string
	Domain          string
	EnableMultiSMTP bool // 是否启用多SMTP端口
}

// GlobalYAMLConfig 全局YAML配置
var GlobalYAMLConfig *YAMLConfig

func Load() *Config {
	// 尝试加载YAML配置文件
	if yamlConfig, err := LoadYAMLConfig("config.yaml"); err == nil {
		GlobalYAMLConfig = yamlConfig
		return &Config{
			WebPort:         strconv.Itoa(yamlConfig.Server.WebPort),
			SMTPPort:        strconv.Itoa(yamlConfig.Server.SMTP.Port25),
			SMTPPort587:     strconv.Itoa(yamlConfig.Server.SMTP.Port587),
			SMTPPort465:     strconv.Itoa(yamlConfig.Server.SMTP.Port465),
			IMAPPort:        strconv.Itoa(yamlConfig.Server.IMAP.Port),
			POP3Port:        strconv.Itoa(yamlConfig.Server.POP3.Port),
			DatabasePath:    yamlConfig.Database.Path,
			SessionKey:      yamlConfig.Security.SessionKey,
			Domain:          yamlConfig.Domain.Default,
			EnableMultiSMTP: yamlConfig.Server.SMTP.EnableMultiPort,
		}
	}

	// 如果YAML配置文件不存在或加载失败，使用环境变量配置
	return &Config{
		WebPort:         getEnv("WEB_PORT", "8080"),
		SMTPPort:        getEnv("SMTP_PORT", "25"),
		SMTPPort587:     getEnv("SMTP_PORT_587", "587"), // SMTP提交端口
		SMTPPort465:     getEnv("SMTP_PORT_465", "465"), // SMTPS端口
		IMAPPort:        getEnv("IMAP_PORT", "143"),
		POP3Port:        getEnv("POP3_PORT", "110"),
		DatabasePath:    getEnv("DATABASE_PATH", "./miko_email.db"),
		SessionKey:      getEnv("SESSION_KEY", "miko-email-secret-key-change-in-production"),
		Domain:          getEnv("DOMAIN", "localhost"),
		EnableMultiSMTP: getEnvBool("ENABLE_MULTI_SMTP", true), // 默认启用多SMTP端口
	}
}

// LoadYAMLConfig 加载YAML配置文件
func LoadYAMLConfig(configPath string) (*YAMLConfig, error) {
	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("配置文件不存在: %s", configPath)
	}

	// 读取配置文件
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析YAML
	var config YAMLConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析YAML配置失败: %v", err)
	}

	return &config, nil
}
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1" || value == "yes"
	}
	return defaultValue
}

// GetSMTPPorts 获取所有SMTP端口
func (c *Config) GetSMTPPorts() []string {
	// 如果有YAML配置，优先使用YAML配置
	if GlobalYAMLConfig != nil {
		if GlobalYAMLConfig.Server.SMTP.EnableMultiPort {
			return []string{
				strconv.Itoa(GlobalYAMLConfig.Server.SMTP.Port25),
				strconv.Itoa(GlobalYAMLConfig.Server.SMTP.Port587),
				strconv.Itoa(GlobalYAMLConfig.Server.SMTP.Port465),
			}
		}
		return []string{strconv.Itoa(GlobalYAMLConfig.Server.SMTP.Port25)}
	}

	// 否则使用原有逻辑
	if c.EnableMultiSMTP {
		return []string{c.SMTPPort, c.SMTPPort587, c.SMTPPort465}
	}
	return []string{c.SMTPPort}
}

// GetAdminCredentials 获取管理员凭据
func GetAdminCredentials() (username, password, email string, enabled bool) {
	if GlobalYAMLConfig != nil {
		return GlobalYAMLConfig.Admin.Username, GlobalYAMLConfig.Admin.Password, GlobalYAMLConfig.Admin.Email, GlobalYAMLConfig.Admin.Enabled
	}
	// 默认管理员凭据
	return getEnv("ADMIN_USERNAME", "admin"),
		getEnv("ADMIN_PASSWORD", "admin123456"),
		getEnv("ADMIN_EMAIL", "admin@localhost"),
		getEnvBool("ADMIN_ENABLED", true)
}
