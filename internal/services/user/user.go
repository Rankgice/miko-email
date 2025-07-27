package user

import (
	"errors"
	"fmt"

	"miko-email/internal/model"
	"miko-email/internal/svc"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	svcCtx *svc.ServiceContext
}

func NewService(svcCtx *svc.ServiceContext) *Service {
	return &Service{svcCtx: svcCtx}
}

// UserWithStats 用户统计信息（使用model中的定义并添加Status字段）
type UserWithStats struct {
	model.UserWithStats
	Status string `json:"status"`
}

// GetUsers 获取用户列表
func (s *Service) GetUsers() ([]UserWithStats, error) {
	modelUsers, err := s.svcCtx.UserModel.GetUsersWithStats()
	if err != nil {
		return nil, err
	}

	// 转换为service层的UserWithStats并设置状态
	var users []UserWithStats
	for _, modelUser := range modelUsers {
		user := UserWithStats{
			UserWithStats: modelUser,
		}
		if modelUser.IsActive {
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
	modelUser, err := s.svcCtx.UserModel.GetUserWithStatsByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, err
	}

	// 转换为service层的UserWithStats并设置状态
	user := &UserWithStats{
		UserWithStats: *modelUser,
	}
	if modelUser.IsActive {
		user.Status = "active"
	} else {
		user.Status = "inactive"
	}

	return user, nil
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

// CreateUser 创建用户（管理员）
func (s *Service) CreateUser(username, email, password string, isActive bool) (*model.User, error) {
	// 检查用户名是否已存在
	existingUser, err := s.svcCtx.UserModel.GetByUsername(username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, fmt.Errorf("用户名已存在")
	}

	// 检查邮箱是否已存在
	existingUser, err = s.svcCtx.UserModel.GetByEmail(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, fmt.Errorf("邮箱已存在")
	}

	// 加密密码
	hashedPassword, err := s.hashPassword(password)
	if err != nil {
		return nil, err
	}

	// 生成用户邀请码
	inviteCode := s.generateInviteCode()

	// 创建用户
	newUser := &model.User{
		Username:   username,
		Password:   hashedPassword,
		Email:      email,
		IsActive:   isActive,
		InviteCode: inviteCode,
	}

	if err := s.svcCtx.UserModel.Create(nil, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

// hashPassword 加密密码
func (s *Service) hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// generateInviteCode 生成邀请码
func (s *Service) generateInviteCode() string {
	return uuid.New().String()
}
