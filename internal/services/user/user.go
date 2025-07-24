package user

import (
	"database/sql"
	"errors"
	"fmt"
	"miko-email/internal/models"

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

// UserWithStats 用户统计信息
type UserWithStats struct {
	models.User
	MailboxCount int    `json:"mailbox_count"`
	Status       string `json:"status"`
	InviterName  string `json:"inviter_name"`
}

// GetUsers 获取用户列表
func (s *Service) GetUsers() ([]UserWithStats, error) {
	query := `
		SELECT 
			u.id, u.username, u.email, u.is_active, u.contribution, 
			u.invite_code, u.invited_by, u.created_at, u.updated_at,
			COUNT(m.id) as mailbox_count,
			COALESCE(inviter.username, admin_inviter.username, '') as inviter_name
		FROM users u
		LEFT JOIN mailboxes m ON u.id = m.user_id AND m.is_active = 1
		LEFT JOIN users inviter ON u.invited_by = inviter.id
		LEFT JOIN admins admin_inviter ON u.invited_by = admin_inviter.id
		GROUP BY u.id, u.username, u.email, u.is_active, u.contribution, 
				 u.invite_code, u.invited_by, u.created_at, u.updated_at,
				 inviter.username, admin_inviter.username
		ORDER BY u.created_at DESC
	`
	db, err := s.svcCtx.DB.DB()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []UserWithStats
	for rows.Next() {
		var user UserWithStats
		err = rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.IsActive, &user.Contribution,
			&user.InviteCode, &user.InvitedBy, &user.CreatedAt, &user.UpdatedAt,
			&user.MailboxCount, &user.InviterName,
		)
		if err != nil {
			return nil, err
		}

		// 设置状态
		if user.IsActive {
			user.Status = "active"
		} else {
			user.Status = "inactive"
		}

		users = append(users, user)
	}

	return users, nil
}

// GetUserByID 根据ID获取用户
func (s *Service) GetUserByID(userID int64) (*UserWithStats, error) {
	query := `
		SELECT 
			u.id, u.username, u.email, u.is_active, u.contribution, 
			u.invite_code, u.invited_by, u.created_at, u.updated_at,
			COUNT(m.id) as mailbox_count,
			COALESCE(inviter.username, admin_inviter.username, '') as inviter_name
		FROM users u
		LEFT JOIN mailboxes m ON u.id = m.user_id AND m.is_active = 1
		LEFT JOIN users inviter ON u.invited_by = inviter.id
		LEFT JOIN admins admin_inviter ON u.invited_by = admin_inviter.id
		WHERE u.id = ?
		GROUP BY u.id, u.username, u.email, u.is_active, u.contribution, 
				 u.invite_code, u.invited_by, u.created_at, u.updated_at,
				 inviter.username, admin_inviter.username
	`

	var user UserWithStats
	db, err := s.svcCtx.DB.DB()
	if err != nil {
		return nil, err
	}
	err = db.QueryRow(query, userID).Scan(&user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, err
	}

	// 设置状态
	if user.IsActive {
		user.Status = "active"
	} else {
		user.Status = "inactive"
	}

	return &user, nil
}

// GetUserMailboxes 获取用户的邮箱列表
func (s *Service) GetUserMailboxes(userID int64) ([]*model.Mailbox, error) {
	var params model.MailboxReq
	isActive := true
	params.IsActive = &isActive
	params.UserId = &userID

	mailboxes, _, err := s.svcCtx.MailboxModel.List(params)
	if err != nil {
		return nil, err
	}

	return mailboxes, nil
}

// UpdateUserStatus 更新用户状态
func (s *Service) UpdateUserStatus(userID int64, isActive bool) error {
	return s.svcCtx.UserModel.UpdateStatus(nil, userID, isActive)
}

// DeleteUser 删除用户（硬删除）
func (s *Service) DeleteUser(userID int64) error {
	// 检查用户是否存在
	_, err := s.svcCtx.UserModel.GetById(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("用户不存在")
		}
		return err
	}

	// 开始事务
	tx := s.svcCtx.DB.Begin()
	defer func() {
		if tx != nil {
			tx.Rollback()
		}
	}()

	// 获取用户的所有邮箱
	var params model.MailboxReq
	params.UserId = &userID
	mailboxes, _, err := s.svcCtx.MailboxModel.List(params)
	if err != nil {
		return err
	}

	// 删除每个邮箱的相关数据
	for _, mailbox := range mailboxes {
		// 1. 删除邮件转发规则
		if err := s.svcCtx.EmailForwardModel.DeleteForwardsByMailboxId(tx, mailbox.Id); err != nil {
			return err
		}

		// 2. 删除邮件
		if err := s.svcCtx.EmailModel.DeleteEmailsByMailboxId(tx, mailbox.Id); err != nil {
			return err
		}

		// 3. 删除邮箱
		if err := s.svcCtx.MailboxModel.HardDelete(tx, mailbox.Id); err != nil {
			return err
		}
	}

	// 4. 删除用户记录
	if err := s.svcCtx.UserModel.HardDelete(tx, userID); err != nil {
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return err
	}
	tx = nil

	return nil
}
