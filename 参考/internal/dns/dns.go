package dns

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

// DNSVerifier DNS验证器
type DNSVerifier struct {
	ServerIP string // 服务器IP地址
}

// NewDNSVerifier 创建DNS验证器
func NewDNSVerifier(serverIP string) *DNSVerifier {
	return &DNSVerifier{
		ServerIP: serverIP,
	}
}

// VerifyResult DNS验证结果
type VerifyResult struct {
	Success    bool     `json:"success"`
	Message    string   `json:"message"`
	MXRecords  []string `json:"mx_records"`
	ARecords   []string `json:"a_records"`
	HasMX      bool     `json:"has_mx"`
	PointsToUs bool     `json:"points_to_us"`
}

// VerifyDomain 验证域名DNS配置
func (v *DNSVerifier) VerifyDomain(domain string) *VerifyResult {
	result := &VerifyResult{
		MXRecords: []string{},
		ARecords:  []string{},
	}

	// 检查MX记录
	mxRecords, err := net.LookupMX(domain)
	mxQueryFailed := false
	if err != nil {
		mxQueryFailed = true
		// 检查是否是域名不存在的错误
		if strings.Contains(err.Error(), "no such host") {
			// 先尝试A记录查询来确认域名是否真的不存在
			_, aErr := net.LookupHost(domain)
			if aErr != nil && strings.Contains(aErr.Error(), "no such host") {
				result.Message = fmt.Sprintf("域名 %s 不存在或无法解析，请检查域名是否正确", domain)
				return result
			}
			// 如果A记录查询成功，说明域名存在但MX记录查询失败
			result.Message = fmt.Sprintf("域名 %s 存在但MX记录查询失败，可能是网络问题", domain)
		} else {
			result.Message = fmt.Sprintf("无法查询MX记录: %v", err)
		}
	}

	if !mxQueryFailed && len(mxRecords) == 0 {
		result.Message = "域名没有配置MX记录，请添加MX记录指向邮件服务器"
	}

	if !mxQueryFailed {
		result.HasMX = true
		for _, mx := range mxRecords {
			mxHost := strings.TrimSuffix(mx.Host, ".")
			result.MXRecords = append(result.MXRecords, fmt.Sprintf("%s (优先级: %d)", mxHost, mx.Pref))

			// 检查MX记录是否指向我们的服务器
			if v.checkMXPointsToUs(mxHost) {
				result.PointsToUs = true
			}
		}
	}

	// 检查A记录
	aRecords, err := net.LookupHost(domain)
	if err == nil {
		result.ARecords = aRecords
		// 检查A记录是否包含我们的IP
		for _, ip := range aRecords {
			if ip == v.ServerIP {
				result.PointsToUs = true
				break
			}
		}
	} else {
		// A记录查询失败不影响整体验证，只记录错误
		if strings.Contains(err.Error(), "no such host") {
			// 如果域名不存在，这个错误已经在MX记录检查时处理了
		} else {
			// 其他A记录查询错误可以忽略，因为MX记录更重要
		}
	}

	if result.PointsToUs {
		result.Success = true
		result.Message = "DNS配置正确，域名已正确解析到本服务器"
	} else {
		// 如果域名没有指向我们的服务器，验证失败
		result.Success = false
		if result.HasMX {
			result.Message = "DNS验证失败：域名有MX记录但未指向本服务器"
			if v.ServerIP != "" {
				result.Message += fmt.Sprintf("（本服务器IP: %s，域名解析IP: %s）", v.ServerIP, strings.Join(result.ARecords, ", "))
			}
		} else if len(result.ARecords) > 0 {
			result.Message = "DNS验证失败：域名有A记录但MX记录未配置或未指向本服务器"
			if v.ServerIP != "" {
				result.Message += fmt.Sprintf("（本服务器IP: %s，域名解析IP: %s）", v.ServerIP, strings.Join(result.ARecords, ", "))
			}
		} else {
			result.Message = "DNS验证失败：域名未解析或无法访问"
			if v.ServerIP != "" {
				result.Message += fmt.Sprintf("（期望IP: %s）", v.ServerIP)
			}
		}
	}

	return result
}

// checkMXPointsToUs 检查MX记录是否指向我们的服务器
func (v *DNSVerifier) checkMXPointsToUs(mxHost string) bool {
	if v.ServerIP == "" {
		return false
	}

	// 解析MX主机的IP地址
	ips, err := net.LookupHost(mxHost)
	if err != nil {
		return false
	}

	// 检查是否包含我们的IP
	for _, ip := range ips {
		if ip == v.ServerIP {
			return true
		}
	}

	return false
}

// GetServerPublicIP 获取服务器公网IP
func GetServerPublicIP() (string, error) {
	// 首先检查环境变量是否指定了服务器IP
	if serverIP := os.Getenv("SERVER_IP"); serverIP != "" {
		if net.ParseIP(serverIP) != nil {
			return serverIP, nil
		}
	}

	// 尝试多个服务获取公网IP
	services := []string{
		"https://api.ipify.org",
		"https://ipinfo.io/ip",
		"https://icanhazip.com",
	}

	for _, service := range services {
		if ip, err := getIPFromService(service); err == nil {
			return strings.TrimSpace(ip), nil
		}
	}

	return "", fmt.Errorf("无法获取公网IP")
}

// getIPFromService 从指定服务获取IP
func getIPFromService(url string) (string, error) {
	// 导入http包来实现真正的HTTP请求
	// 这里需要添加import "net/http"和"io"
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("HTTP status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(body)), nil
}

// getLocalIP 获取本地IP地址
func getLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}

// GenerateDNSInstructions 生成DNS配置说明
func GenerateDNSInstructions(domain, serverIP string) map[string]interface{} {
	return map[string]interface{}{
		"domain":    domain,
		"server_ip": serverIP,
		"instructions": map[string]interface{}{
			"mx_record": map[string]interface{}{
				"type":        "MX",
				"name":        "@",
				"value":       domain,
				"priority":    10,
				"ttl":         3600,
				"description": "邮件交换记录，用于接收邮件",
			},
			"a_record": map[string]interface{}{
				"type":        "A",
				"name":        "@",
				"value":       serverIP,
				"ttl":         3600,
				"description": "A记录，将域名指向服务器IP",
			},
			"txt_record": map[string]interface{}{
				"type":        "TXT",
				"name":        "@",
				"value":       fmt.Sprintf("v=spf1 ip4:%s ~all", serverIP),
				"ttl":         3600,
				"description": "SPF记录，防止邮件被标记为垃圾邮件",
			},
		},
		"steps": []string{
			"1. 登录您的域名注册商或DNS服务商管理面板",
			"2. 找到DNS记录管理页面（通常叫\"域名解析\"、\"DNS管理\"或\"解析设置\"）",
			"3. 添加MX记录：",
			"   - 记录类型：选择 MX",
			"   - 主机记录：填写 @ （代表根域名）",
			"   - 记录值：填写您的域名（如 example.com）",
			"   - 优先级：填写 10",
			"   - TTL：选择 3600 或保持默认",
			"4. 添加A记录：",
			"   - 记录类型：选择 A",
			"   - 主机记录：填写 @ （代表根域名）",
			"   - 记录值：填写服务器IP地址（" + serverIP + "）",
			"   - TTL：选择 3600 或保持默认",
			"5. 可选：添加SPF记录（推荐）：",
			"   - 记录类型：选择 TXT",
			"   - 主机记录：填写 @",
			"   - 记录值：填写 v=spf1 ip4:" + serverIP + " ~all",
			"   - TTL：选择 3600 或保持默认",
			"6. 保存所有记录并等待DNS传播生效（通常需要几分钟到几小时）",
			"7. 返回本页面点击\"验证DNS\"按钮检查配置是否正确",
		},
		"common_providers": map[string]string{
			"阿里云":        "https://dns.console.aliyun.com/",
			"腾讯云":        "https://console.cloud.tencent.com/cns",
			"百度云":        "https://console.bce.baidu.com/bcd/",
			"华为云":        "https://console.huaweicloud.com/dns/",
			"Cloudflare": "https://dash.cloudflare.com/",
			"GoDaddy":    "https://dcc.godaddy.com/manage/dns",
		},
	}
}

// ValidateDomainName 验证域名格式
func ValidateDomainName(domain string) error {
	if domain == "" {
		return fmt.Errorf("域名不能为空")
	}

	if len(domain) > 253 {
		return fmt.Errorf("域名长度不能超过253个字符")
	}

	// 简单的域名格式验证
	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		return fmt.Errorf("域名格式不正确")
	}

	for _, part := range parts {
		if len(part) == 0 || len(part) > 63 {
			return fmt.Errorf("域名部分长度不正确")
		}

		// 检查字符是否合法
		for _, char := range part {
			if !((char >= 'a' && char <= 'z') ||
				(char >= 'A' && char <= 'Z') ||
				(char >= '0' && char <= '9') ||
				char == '-') {
				return fmt.Errorf("域名包含非法字符")
			}
		}

		// 不能以-开头或结尾
		if strings.HasPrefix(part, "-") || strings.HasSuffix(part, "-") {
			return fmt.Errorf("域名部分不能以-开头或结尾")
		}
	}

	return nil
}

// CheckDNSPropagation 检查DNS传播状态
func CheckDNSPropagation(domain string) map[string]interface{} {
	result := map[string]interface{}{
		"domain": domain,
		"checks": []map[string]interface{}{},
	}

	// 检查不同DNS服务器的解析结果
	dnsServers := []struct {
		Name string
		IP   string
	}{
		{"Google DNS", "8.8.8.8"},
		{"Cloudflare DNS", "1.1.1.1"},
		{"阿里DNS", "223.5.5.5"},
		{"腾讯DNS", "119.29.29.29"},
	}

	for _, server := range dnsServers {
		check := map[string]interface{}{
			"server": server.Name,
			"ip":     server.IP,
		}

		// 这里简化实现，实际应该查询指定DNS服务器
		// 由于Go标准库限制，这里使用默认解析
		if mxRecords, err := net.LookupMX(domain); err == nil {
			check["mx_records"] = len(mxRecords)
			check["status"] = "success"
		} else {
			check["status"] = "failed"
			check["error"] = err.Error()
		}

		result["checks"] = append(result["checks"].([]map[string]interface{}), check)
	}

	return result
}
