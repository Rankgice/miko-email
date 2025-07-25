package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type YConfig struct {
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

// NewConfig 创建配置
func NewConfig(path string) YConfig {
	conf, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("读取配置文件失败：", err)
	}
	var c YConfig
	if err := yaml.Unmarshal(conf, &c); err != nil {
		log.Fatal("解析配置文件失败：", err)
	}
	return c
}

//		WebPort:         getEnv("WEB_PORT", "8080"),
//		SMTPPort:        getEnv("SMTP_PORT", "25"),
//		SMTPPort587:     getEnv("SMTP_PORT_587", "587"), // SMTP提交端口
//		SMTPPort465:     getEnv("SMTP_PORT_465", "465"), // SMTPS端口
//		IMAPPort:        getEnv("IMAP_PORT", "143"),
//		POP3Port:        getEnv("POP3_PORT", "110"),
//		DatabasePath:    getEnv("DATABASE_PATH", "./miko_email.db"),
//		SessionKey:      getEnv("SESSION_KEY", "miko-email-secret-key-change-in-production"),
//		Domain:          getEnv("DOMAIN", "localhost"),
//		EnableMultiSMTP: getEnvBool("ENABLE_MULTI_SMTP", true), // 默认启用多SMTP端口
