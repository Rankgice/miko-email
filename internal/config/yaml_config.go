package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
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

// ToConfig 将YAML配置转换为原有的Config结构
func (yc *YAMLConfig) ToConfig() *Config {
	return &Config{
		WebPort:         strconv.Itoa(yc.Server.WebPort),
		SMTPPort:        strconv.Itoa(yc.Server.SMTP.Port25),
		SMTPPort587:     strconv.Itoa(yc.Server.SMTP.Port587),
		SMTPPort465:     strconv.Itoa(yc.Server.SMTP.Port465),
		IMAPPort:        strconv.Itoa(yc.Server.IMAP.Port),
		POP3Port:        strconv.Itoa(yc.Server.POP3.Port),
		DatabasePath:    yc.Database.Path,
		SessionKey:      yc.Security.SessionKey,
		Domain:          yc.Domain.Default,
		EnableMultiSMTP: yc.Server.SMTP.EnableMultiPort,
	}
}

// GetAdminCredentials 获取管理员凭据
func (yc *YAMLConfig) GetAdminCredentials() (username, password, email string, enabled bool) {
	return yc.Admin.Username, yc.Admin.Password, yc.Admin.Email, yc.Admin.Enabled
}

// GetSMTPPorts 获取SMTP端口列表
func (yc *YAMLConfig) GetSMTPPorts() []string {
	if yc.Server.SMTP.EnableMultiPort {
		return []string{
			strconv.Itoa(yc.Server.SMTP.Port25),
			strconv.Itoa(yc.Server.SMTP.Port587),
			strconv.Itoa(yc.Server.SMTP.Port465),
		}
	}
	return []string{strconv.Itoa(yc.Server.SMTP.Port25)}
}
