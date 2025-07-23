package domain

import (
	"database/sql"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"
	
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

// GetDomains 获取域名列表
func (s *Service) GetDomains() ([]models.Domain, error) {
	query := `
		SELECT id, name, is_verified, is_active, mx_record, a_record, txt_record, created_at, updated_at
		FROM domains
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var domains []models.Domain
	for rows.Next() {
		var domain models.Domain
		err = rows.Scan(&domain.ID, &domain.Name, &domain.IsVerified, &domain.IsActive,
			&domain.MXRecord, &domain.ARecord, &domain.TXTRecord,
			&domain.CreatedAt, &domain.UpdatedAt)
		if err != nil {
			return nil, err
		}
		domains = append(domains, domain)
	}

	return domains, nil
}

// CreateDomain 创建域名
func (s *Service) CreateDomain(name, mxRecord, aRecord, txtRecord string) (*models.Domain, error) {
	// 检查域名是否已存在
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM domains WHERE name = ?", name).Scan(&count)
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, fmt.Errorf("域名已存在")
	}

	// 插入域名
	result, err := s.db.Exec(`
		INSERT INTO domains (name, mx_record, a_record, txt_record, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, name, mxRecord, aRecord, txtRecord, time.Now(), time.Now())

	if err != nil {
		return nil, err
	}

	domainID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	domain := &models.Domain{
		ID:         int(domainID),
		Name:       name,
		IsVerified: false,
		IsActive:   true,
		MXRecord:   mxRecord,
		ARecord:    aRecord,
		TXTRecord:  txtRecord,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return domain, nil
}

// VerifyDomain 验证域名DNS设置
func (s *Service) VerifyDomain(domainID int) (*models.Domain, error) {
	// 获取域名信息
	var domain models.Domain
	query := `
		SELECT id, name, is_verified, is_active, mx_record, a_record, txt_record, created_at, updated_at
		FROM domains
		WHERE id = ?
	`

	err := s.db.QueryRow(query, domainID).Scan(
		&domain.ID, &domain.Name, &domain.IsVerified, &domain.IsActive,
		&domain.MXRecord, &domain.ARecord, &domain.TXTRecord,
		&domain.CreatedAt, &domain.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("域名不存在")
		}
		return nil, err
	}

	// 验证DNS记录
	verified := true

	// 验证MX记录
	if domain.MXRecord != "" {
		if !s.verifyMXRecord(domain.Name, domain.MXRecord) {
			verified = false
		}
	}

	// 验证A记录
	if domain.ARecord != "" {
		if !s.verifyARecord(domain.Name, domain.ARecord) {
			verified = false
		}
	}

	// 验证TXT记录
	if domain.TXTRecord != "" {
		if !s.verifyTXTRecord(domain.Name, domain.TXTRecord) {
			verified = false
		}
	}

	// 更新验证状态
	_, err = s.db.Exec("UPDATE domains SET is_verified = ?, updated_at = ? WHERE id = ?",
		verified, time.Now(), domainID)
	if err != nil {
		return nil, err
	}

	domain.IsVerified = verified
	domain.UpdatedAt = time.Now()

	return &domain, nil
}

// verifyMXRecord 验证MX记录
func (s *Service) verifyMXRecord(domain, expectedMX string) bool {
	// 使用标准库验证
	mxRecords, err := net.LookupMX(domain)
	if err == nil {
		for _, mx := range mxRecords {
			if strings.TrimSuffix(mx.Host, ".") == strings.TrimSuffix(expectedMX, ".") {
				return true
			}
		}
	}

	// 使用DNS库进行更详细的验证
	return s.verifyDNSRecord(domain, dns.TypeMX, expectedMX)
}

// verifyARecord 验证A记录
func (s *Service) verifyARecord(domain, expectedIP string) bool {
	// 使用标准库验证
	ips, err := net.LookupIP(domain)
	if err == nil {
		for _, ip := range ips {
			if ip.String() == expectedIP {
				return true
			}
		}
	}

	// 使用DNS库进行更详细的验证
	return s.verifyDNSRecord(domain, dns.TypeA, expectedIP)
}

// verifyTXTRecord 验证TXT记录
func (s *Service) verifyTXTRecord(domain, expectedTXT string) bool {
	// 使用标准库验证
	txtRecords, err := net.LookupTXT(domain)
	if err == nil {
		for _, txt := range txtRecords {
			if txt == expectedTXT {
				return true
			}
		}
	}

	// 使用DNS库进行更详细的验证
	return s.verifyDNSRecord(domain, dns.TypeTXT, expectedTXT)
}

// verifyDNSRecord 使用DNS库验证DNS记录
func (s *Service) verifyDNSRecord(domain string, recordType uint16, expectedValue string) bool {
	c := dns.Client{
		Timeout: time.Second * 5,
	}

	// 构造DNS查询
	m := dns.Msg{}
	m.SetQuestion(dns.Fqdn(domain), recordType)

	// 查询DNS服务器
	dnsServers := []string{"8.8.8.8:53", "1.1.1.1:53", "114.114.114.114:53"}

	for _, server := range dnsServers {
		r, _, err := c.Exchange(&m, server)
		if err != nil {
			continue
		}

		// 检查响应
		for _, ans := range r.Answer {
			switch recordType {
			case dns.TypeMX:
				if mx, ok := ans.(*dns.MX); ok {
					if strings.TrimSuffix(mx.Mx, ".") == strings.TrimSuffix(expectedValue, ".") {
						return true
					}
				}
			case dns.TypeA:
				if a, ok := ans.(*dns.A); ok {
					if a.A.String() == expectedValue {
						return true
					}
				}
			case dns.TypeTXT:
				if txt, ok := ans.(*dns.TXT); ok {
					for _, t := range txt.Txt {
						if t == expectedValue {
							return true
						}
					}
				}
			}
		}
	}

	return false
}

// GetDNSRecords 获取域名的所有DNS记录信息
func (s *Service) GetDNSRecords(domain string) map[string][]string {
	records := make(map[string][]string)

	// 获取MX记录
	if mxRecords, err := net.LookupMX(domain); err == nil {
		var mxList []string
		for _, mx := range mxRecords {
			mxList = append(mxList, fmt.Sprintf("%s (优先级: %d)", strings.TrimSuffix(mx.Host, "."), mx.Pref))
		}
		records["MX"] = mxList
	}

	// 获取A记录
	if ips, err := net.LookupIP(domain); err == nil {
		var aList []string
		for _, ip := range ips {
			if ip.To4() != nil { // 只获取IPv4地址
				aList = append(aList, ip.String())
			}
		}
		records["A"] = aList
	}

	// 获取TXT记录
	if txtRecords, err := net.LookupTXT(domain); err == nil {
		records["TXT"] = txtRecords
	}

	// 获取CNAME记录
	if cname, err := net.LookupCNAME(domain); err == nil && cname != domain+"." {
		records["CNAME"] = []string{strings.TrimSuffix(cname, ".")}
	}

	return records
}

// GetDomainByID 根据ID获取域名
func (s *Service) GetDomainByID(domainID int) (*models.Domain, error) {
	var domain models.Domain
	query := `
		SELECT id, name, is_verified, is_active, mx_record, a_record, txt_record, created_at, updated_at
		FROM domains
		WHERE id = ?
	`

	err := s.db.QueryRow(query, domainID).Scan(
		&domain.ID, &domain.Name, &domain.IsVerified, &domain.IsActive,
		&domain.MXRecord, &domain.ARecord, &domain.TXTRecord,
		&domain.CreatedAt, &domain.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("域名不存在")
		}
		return nil, err
	}

	return &domain, nil
}

// UpdateDomain 更新域名信息
func (s *Service) UpdateDomain(domainID int, mxRecord, aRecord, txtRecord string) error {
	_, err := s.db.Exec(`
		UPDATE domains
		SET mx_record = ?, a_record = ?, txt_record = ?, updated_at = ?
		WHERE id = ?
	`, mxRecord, aRecord, txtRecord, time.Now(), domainID)

	return err
}

// DeleteDomain 删除域名
func (s *Service) DeleteDomain(domainID int) error {
	// 检查是否有邮箱使用此域名
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM mailboxes WHERE domain_id = ? AND is_active = 1", domainID).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("该域名下还有邮箱，无法删除")
	}

	// 真正删除域名记录
	_, err = s.db.Exec("DELETE FROM domains WHERE id = ?", domainID)
	return err
}

// GetAvailableDomains 获取可用的域名列表（已验证且激活的）
func (s *Service) GetAvailableDomains() ([]models.Domain, error) {
	query := `
		SELECT id, name, is_verified, is_active, mx_record, a_record, txt_record, created_at, updated_at
		FROM domains
		WHERE is_active = 1 AND is_verified = 1
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var domains []models.Domain
	for rows.Next() {
		var domain models.Domain
		err = rows.Scan(&domain.ID, &domain.Name, &domain.IsVerified, &domain.IsActive,
			&domain.MXRecord, &domain.ARecord, &domain.TXTRecord,
			&domain.CreatedAt, &domain.UpdatedAt)
		if err != nil {
			return nil, err
		}
		domains = append(domains, domain)
	}

	return domains, nil
}
