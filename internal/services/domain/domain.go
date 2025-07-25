package domain

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"
	"gorm.io/gorm"
	"miko-email/internal/model"
	"miko-email/internal/svc"
)

type Service struct {
	svcCtx *svc.ServiceContext
}

func NewService(svcCtx *svc.ServiceContext) *Service {
	return &Service{svcCtx: svcCtx}
}

// GetDomains 获取域名列表
func (s *Service) GetDomains() ([]*model.Domain, error) {
	// 使用DomainModel的List方法获取所有域名，按创建时间倒序
	domains, _, err := s.svcCtx.DomainModel.List(model.DomainReq{})
	if err != nil {
		return nil, fmt.Errorf("查询域名列表失败: %v", err)
	}
	// 补充默认状态
	for i, domain := range domains {
		if domain.SenderVerificationStatus == "" {
			domains[i].SenderVerificationStatus = "pending"
		}
		if domain.ReceiverVerificationStatus == "" {
			domains[i].ReceiverVerificationStatus = "pending"
		}
	}
	return domains, nil
}

// CreateDomain 创建域名
func (s *Service) CreateDomain(name, mxRecord, aRecord, txtRecord string) (*model.Domain, error) {
	// 检查域名是否已存在
	exists, err := s.svcCtx.DomainModel.CheckDomainExist(name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("域名已存在")
	}

	// 创建域名
	domain := &model.Domain{
		Name:                       name,
		IsVerified:                 false,
		IsActive:                   true,
		MxRecord:                   mxRecord,
		ARecord:                    aRecord,
		TxtRecord:                  txtRecord,
		SenderVerificationStatus:   "pending",
		ReceiverVerificationStatus: "pending",
		CreatedAt:                  time.Now(),
		UpdatedAt:                  time.Now(),
	}

	if err := s.svcCtx.DomainModel.Create(nil, domain); err != nil {
		return nil, err
	}

	return domain, nil
}

// CreateDomainWithAllRecords 创建域名（包含所有DNS记录）
func (s *Service) CreateDomainWithAllRecords(name, mxRecord, aRecord, txtRecord, spfRecord, dmarcRecord, dkimRecord, ptrRecord string) (*model.Domain, error) {
	// 检查域名是否已存在
	exists, err := s.svcCtx.DomainModel.CheckDomainExist(name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("域名已存在")
	}

	// 创建域名
	domain := &model.Domain{
		Name:                       name,
		IsVerified:                 false,
		IsActive:                   true,
		MxRecord:                   mxRecord,
		ARecord:                    aRecord,
		TxtRecord:                  txtRecord,
		SPFRecord:                  spfRecord,
		DMARCRecord:                dmarcRecord,
		DKIMRecord:                 dkimRecord,
		PTRRecord:                  ptrRecord,
		SenderVerificationStatus:   "pending",
		ReceiverVerificationStatus: "pending",
		CreatedAt:                  time.Now(),
		UpdatedAt:                  time.Now(),
	}

	if err := s.svcCtx.DomainModel.Create(nil, domain); err != nil {
		return nil, err
	}

	return domain, nil
}

// VerifyDomain 验证域名DNS设置
func (s *Service) VerifyDomain(domainID int64) (*model.Domain, error) {
	// 获取域名信息
	domain, err := s.svcCtx.DomainModel.GetById(domainID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("域名不存在")
		}
		return nil, fmt.Errorf("查询域名失败: %v", err)
	}
	// 补充默认状态
	if domain.SenderVerificationStatus == "" {
		domain.SenderVerificationStatus = "pending"
	}
	if domain.ReceiverVerificationStatus == "" {
		domain.ReceiverVerificationStatus = "pending"
	}

	// 验证DNS记录
	verified := true

	// 验证MX记录
	if domain.MxRecord != "" {
		if !s.verifyMXRecord(domain.Name, domain.MxRecord) {
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
	if domain.TxtRecord != "" {
		if !s.verifyTXTRecord(domain.Name, domain.TxtRecord) {
			verified = false
		}
	}

	// 更新验证状态
	if err := s.svcCtx.DomainModel.UpdateVerificationStatus(nil, domainID, verified); err != nil {
		return nil, err
	}

	domain.IsVerified = verified
	domain.UpdatedAt = time.Now()

	return domain, nil
}

// VerifySenderConfiguration 验证发件配置
func (s *Service) VerifySenderConfiguration(domainID int64) (*model.Domain, error) {
	// 获取域名信息
	domain, err := s.GetDomainByID(domainID)
	if err != nil {
		return nil, err
	}

	senderStatus := "verified"

	// 验证SPF记录
	if domain.SPFRecord != "" {
		if !s.verifySPFRecord(domain.Name, domain.SPFRecord) {
			senderStatus = "failed"
		}
	}

	// 验证DKIM记录
	if domain.DKIMRecord != "" {
		if !s.verifyDKIMRecord(domain.Name, domain.DKIMRecord) {
			senderStatus = "failed"
		}
	}

	// 验证DMARC记录
	if domain.DMARCRecord != "" {
		if !s.verifyDMARCRecord(domain.Name, domain.DMARCRecord) {
			senderStatus = "failed"
		}
	}

	// 更新发件验证状态
	if err := s.svcCtx.DomainModel.UpdateSenderVerificationStatus(nil, domainID, senderStatus); err != nil {
		return nil, err
	}

	domain.SenderVerificationStatus = senderStatus
	domain.UpdatedAt = time.Now()

	return domain, nil
}

// VerifyReceiverConfiguration 验证收件配置
func (s *Service) VerifyReceiverConfiguration(domainID int64) (*model.Domain, error) {
	// 获取域名信息
	domain, err := s.GetDomainByID(domainID)
	if err != nil {
		return nil, err
	}

	receiverStatus := "verified"

	// 验证MX记录
	if domain.MxRecord != "" {
		if !s.verifyMXRecord(domain.Name, domain.MxRecord) {
			receiverStatus = "failed"
		}
	}

	// 验证A记录
	if domain.ARecord != "" {
		if !s.verifyARecord(domain.Name, domain.ARecord) {
			receiverStatus = "failed"
		}
	}

	// 验证PTR记录
	if domain.PTRRecord != "" {
		if !s.verifyPTRRecord(domain.ARecord, domain.PTRRecord) {
			receiverStatus = "failed"
		}
	}

	// 更新收件验证状态
	if err := s.svcCtx.DomainModel.UpdateReceiverVerificationStatus(nil, domainID, receiverStatus); err != nil {
		return nil, err
	}

	domain.ReceiverVerificationStatus = receiverStatus
	domain.UpdatedAt = time.Now()

	return domain, nil
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

// verifySPFRecord 验证SPF记录
func (s *Service) verifySPFRecord(domain, expectedSPF string) bool {
	// SPF记录通常在TXT记录中
	txtRecords, err := net.LookupTXT(domain)
	if err == nil {
		for _, txt := range txtRecords {
			if strings.HasPrefix(txt, "v=spf1") && strings.Contains(txt, expectedSPF) {
				return true
			}
		}
	}
	return s.verifyDNSRecord(domain, dns.TypeTXT, expectedSPF)
}

// verifyDMARCRecord 验证DMARC记录
func (s *Service) verifyDMARCRecord(domain, expectedDMARC string) bool {
	// DMARC记录在_dmarc子域名的TXT记录中
	dmarcDomain := "_dmarc." + domain
	txtRecords, err := net.LookupTXT(dmarcDomain)
	if err == nil {
		for _, txt := range txtRecords {
			if strings.HasPrefix(txt, "v=DMARC1") && strings.Contains(txt, expectedDMARC) {
				return true
			}
		}
	}
	return s.verifyDNSRecord(dmarcDomain, dns.TypeTXT, expectedDMARC)
}

// verifyDKIMRecord 验证DKIM记录
func (s *Service) verifyDKIMRecord(domain, expectedDKIM string) bool {
	// DKIM记录通常在selector._domainkey.domain的TXT记录中
	// 这里假设使用default作为selector
	dkimDomain := "default._domainkey." + domain
	txtRecords, err := net.LookupTXT(dkimDomain)
	if err == nil {
		for _, txt := range txtRecords {
			if strings.Contains(txt, "v=DKIM1") && strings.Contains(txt, expectedDKIM) {
				return true
			}
		}
	}
	return s.verifyDNSRecord(dkimDomain, dns.TypeTXT, expectedDKIM)
}

// verifyPTRRecord 验证PTR记录
func (s *Service) verifyPTRRecord(ip, expectedPTR string) bool {
	// 反向DNS查询
	names, err := net.LookupAddr(ip)
	if err == nil {
		for _, name := range names {
			if strings.TrimSuffix(name, ".") == strings.TrimSuffix(expectedPTR, ".") {
				return true
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
func (s *Service) GetDomainByID(domainID int64) (*model.Domain, error) {
	domain, err := s.svcCtx.DomainModel.GetById(domainID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("域名不存在")
		}
		return nil, err
	}
	if domain.SenderVerificationStatus == "" {
		domain.SenderVerificationStatus = "pending"
	}
	if domain.ReceiverVerificationStatus == "" {
		domain.ReceiverVerificationStatus = "pending"
	}
	return domain, nil
}

// UpdateDomain 更新域名信息
func (s *Service) UpdateDomain(domainID int64, mxRecord, aRecord, txtRecord string) error {
	return s.svcCtx.DomainModel.UpdateDNSRecords(nil, domainID, mxRecord, aRecord, txtRecord)
}

// UpdateDomainWithAllRecords 更新域名信息（包含所有DNS记录）
func (s *Service) UpdateDomainWithAllRecords(domainID int64, mxRecord, aRecord, txtRecord, spfRecord, dmarcRecord, dkimRecord, ptrRecord string) error {
	return s.svcCtx.DomainModel.UpdateAllDNSRecords(nil, int64(domainID), mxRecord, aRecord, txtRecord, spfRecord, dmarcRecord, dkimRecord, ptrRecord)
}

// DeleteDomain 删除域名
func (s *Service) DeleteDomain(domainID int64) error {
	// 检查是否有邮箱使用此域名
	count, err := s.svcCtx.MailboxModel.CountMailboxesByDomainId(domainID)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("该域名下还有邮箱，无法删除")
	}

	// 获取域名对象用于删除
	domain, err := s.svcCtx.DomainModel.GetById(domainID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("域名不存在")
		}
		return err
	}

	// 删除域名记录
	return s.svcCtx.DomainModel.Delete(nil, domain)
}

// GetAvailableDomains 获取可用的域名列表（已验证且激活的）
func (s *Service) GetAvailableDomains() ([]*model.Domain, error) {
	domains, err := s.svcCtx.DomainModel.GetAvailableDomains()
	if err != nil {
		return nil, err
	}
	// 补充默认状态
	for i, domain := range domains {
		if domain.SenderVerificationStatus == "" {
			domains[i].SenderVerificationStatus = "pending"
		}
		if domain.ReceiverVerificationStatus == "" {
			domains[i].ReceiverVerificationStatus = "pending"
		}
	}
	return domains, nil
}
