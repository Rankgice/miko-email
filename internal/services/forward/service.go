package forward

import (
	"errors"
	"fmt"
	"time"

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

// ForwardRule 转发规则结构（与model.EmailForward保持一致）
type ForwardRule struct {
	ID                 int64      `json:"id"`
	MailboxID          int64      `json:"mailbox_id"`
	SourceEmail        string     `json:"source_email"`
	TargetEmail        string     `json:"target_email"`
	Enabled            bool       `json:"enabled"`
	KeepOriginal       bool       `json:"keep_original"`
	ForwardAttachments bool       `json:"forward_attachments"`
	SubjectPrefix      string     `json:"subject_prefix"`
	Description        string     `json:"description"`
	ForwardCount       int64      `json:"forward_count"`
	LastForwardAt      *time.Time `json:"last_forward_at"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// CreateForwardRuleRequest 创建转发规则请求
type CreateForwardRuleRequest struct {
	SourceEmail        string `json:"source_email" binding:"required"`
	TargetEmail        string `json:"target_email" binding:"required"`
	Enabled            bool   `json:"enabled"`
	KeepOriginal       bool   `json:"keep_original"`
	ForwardAttachments bool   `json:"forward_attachments"`
	SubjectPrefix      string `json:"subject_prefix"`
	Description        string `json:"description"`
}

// convertToForwardRule 将model.EmailForward转换为ForwardRule
func convertToForwardRule(ef *model.EmailForward) ForwardRule {
	return ForwardRule{
		ID:                 ef.Id,
		MailboxID:          ef.MailboxId,
		SourceEmail:        ef.SourceEmail,
		TargetEmail:        ef.TargetEmail,
		Enabled:            ef.Enabled,
		KeepOriginal:       ef.KeepOriginal,
		ForwardAttachments: ef.ForwardAttachments,
		SubjectPrefix:      ef.SubjectPrefix,
		Description:        ef.Description,
		ForwardCount:       ef.ForwardCount,
		LastForwardAt:      ef.LastForwardAt,
		CreatedAt:          ef.CreatedAt,
		UpdatedAt:          ef.UpdatedAt,
	}
}

// GetForwardRulesByUser 获取用户的转发规则
func (s *Service) GetForwardRulesByUser(userID int64) ([]ForwardRule, error) {
	emailForwards, err := s.svcCtx.EmailForwardModel.GetForwardsByUserId(userID)
	if err != nil {
		return nil, fmt.Errorf("查询转发规则失败: %w", err)
	}

	var rules []ForwardRule
	for _, ef := range emailForwards {
		rules = append(rules, convertToForwardRule(ef))
	}

	return rules, nil
}

// GetForwardRuleByID 根据ID获取转发规则
func (s *Service) GetForwardRuleByID(ruleID int64, userID int64) (*ForwardRule, error) {
	emailForward, err := s.svcCtx.EmailForwardModel.GetForwardByIdAndUserId(ruleID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("转发规则不存在")
		}
		return nil, fmt.Errorf("查询转发规则失败: %w", err)
	}

	rule := convertToForwardRule(emailForward)
	return &rule, nil
}

// CreateForwardRule 创建转发规则
func (s *Service) CreateForwardRule(userID int64, req CreateForwardRuleRequest) (*ForwardRule, error) {
	// 首先获取邮箱
	mailbox, err := s.svcCtx.MailboxModel.GetByEmailAndUserId(req.SourceEmail, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("源邮箱不存在或不属于当前用户")
		}
		return nil, fmt.Errorf("查询邮箱失败: %w", err)
	}

	// 检查是否已存在相同的转发规则
	exists, err := s.svcCtx.EmailForwardModel.CheckForwardRuleExistByTarget(mailbox.Id, req.TargetEmail)
	if err != nil {
		return nil, fmt.Errorf("检查转发规则失败: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("转发规则已存在")
	}

	// 创建转发规则
	now := time.Now()
	emailForward := &model.EmailForward{
		MailboxId:          mailbox.Id,
		SourceEmail:        req.SourceEmail,
		TargetEmail:        req.TargetEmail,
		Enabled:            req.Enabled,
		KeepOriginal:       req.KeepOriginal,
		ForwardAttachments: req.ForwardAttachments,
		SubjectPrefix:      req.SubjectPrefix,
		Description:        req.Description,
		ForwardCount:       0,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := s.svcCtx.EmailForwardModel.Create(nil, emailForward); err != nil {
		return nil, fmt.Errorf("创建转发规则失败: %w", err)
	}

	// 返回创建的规则
	return s.GetForwardRuleByID(emailForward.Id, userID)
}

// UpdateForwardRule 更新转发规则
func (s *Service) UpdateForwardRule(ruleID int64, userID int64, req CreateForwardRuleRequest) error {
	// 首先检查规则是否存在且属于当前用户
	_, err := s.GetForwardRuleByID(ruleID, userID)
	if err != nil {
		return err
	}

	// 获取新的邮箱
	mailbox, err := s.svcCtx.MailboxModel.GetByEmailAndUserId(req.SourceEmail, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("源邮箱不存在或不属于当前用户")
		}
		return fmt.Errorf("查询邮箱失败: %w", err)
	}

	// 更新转发规则
	updateData := map[string]interface{}{
		"mailbox_id":          mailbox.Id,
		"source_email":        req.SourceEmail,
		"target_email":        req.TargetEmail,
		"enabled":             req.Enabled,
		"keep_original":       req.KeepOriginal,
		"forward_attachments": req.ForwardAttachments,
		"subject_prefix":      req.SubjectPrefix,
		"description":         req.Description,
		"updated_at":          time.Now(),
	}

	if err := s.svcCtx.EmailForwardModel.MapUpdate(nil, ruleID, updateData); err != nil {
		return fmt.Errorf("更新转发规则失败: %w", err)
	}

	return nil
}

// DeleteForwardRule 删除转发规则
func (s *Service) DeleteForwardRule(ruleID int64, userID int64) error {
	// 首先检查规则是否存在且属于当前用户
	emailForward, err := s.svcCtx.EmailForwardModel.GetForwardByIdAndUserId(ruleID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("转发规则不存在")
		}
		return err
	}

	// 删除转发规则
	if err := s.svcCtx.EmailForwardModel.Delete(nil, emailForward); err != nil {
		return fmt.Errorf("删除转发规则失败: %w", err)
	}

	return nil
}

// ToggleForwardRule 切换转发规则状态
func (s *Service) ToggleForwardRule(ruleID int64, userID int64, enabled bool) error {
	// 首先检查规则是否存在且属于当前用户
	_, err := s.GetForwardRuleByID(ruleID, userID)
	if err != nil {
		return err
	}

	// 更新状态
	if err := s.svcCtx.EmailForwardModel.UpdateStatus(nil, ruleID, enabled); err != nil {
		return fmt.Errorf("更新转发规则状态失败: %w", err)
	}

	return nil
}

// GetForwardStatistics 获取转发统计信息
func (s *Service) GetForwardStatistics(userID int64) (map[string]interface{}, error) {
	return s.svcCtx.EmailForwardModel.GetUserForwardStatistics(userID)
}

// IncrementForwardCount 增加转发次数
func (s *Service) IncrementForwardCount(ruleID int64) error {
	if err := s.svcCtx.EmailForwardModel.IncrementForwardCount(nil, ruleID); err != nil {
		return fmt.Errorf("更新转发次数失败: %w", err)
	}

	return nil
}

// GetActiveForwardRules 获取指定邮箱的活跃转发规则
func (s *Service) GetActiveForwardRules(sourceEmail string) ([]ForwardRule, error) {
	emailForwards, err := s.svcCtx.EmailForwardModel.GetForwardsBySourceEmail(sourceEmail)
	if err != nil {
		return nil, fmt.Errorf("查询活跃转发规则失败: %w", err)
	}

	var rules []ForwardRule
	for _, ef := range emailForwards {
		rules = append(rules, convertToForwardRule(ef))
	}

	return rules, nil
}
