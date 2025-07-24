package auth

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

// AuthenticateUser 验证普通用户
func (s *Service) AuthenticateUser(username, password string) (*model.User, error) {
	// 使用UserModel获取用户
	user, err := s.svcCtx.UserModel.GetByUsername(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("用户不存在")
		}
		return nil, err
	} else if !user.IsActive {
		return nil, fmt.Errorf("用户已被禁用")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("密码错误")
	}

	return user, nil
}

// AuthenticateAdmin 验证管理员
func (s *Service) AuthenticateAdmin(username, password string) (*model.Admin, error) {
	// 使用AdminModel获取管理员
	admin, err := s.svcCtx.AdminModel.GetByUsername(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("管理员不存在或已被禁用")
		}
		return nil, err
	} else if !admin.IsActive { // 检查管理员是否激活
		return nil, fmt.Errorf("管理员不存在或已被禁用")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("密码错误")
	}

	return admin, nil
}

// RegisterUser 注册普通用户
func (s *Service) RegisterUser(username, password, email, domainPrefix string, domainID int64, inviteCode string) (*model.User, error) {
	// 检查用户名是否已存在
	existingUser, err := s.svcCtx.UserModel.GetByUsername(username)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if existingUser != nil {
		return nil, fmt.Errorf("用户名已存在")
	}

	// 检查邮箱是否已存在
	existingUser, err = s.svcCtx.UserModel.GetByEmail(email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if existingUser != nil {
		return nil, fmt.Errorf("邮箱已存在")
	}

	// 验证邀请码并获取邀请人ID
	var invitedBy *int64
	if inviteCode != "" {
		// 先检查普通用户的邀请码
		inviterUser, err := s.svcCtx.UserModel.GetByInviteCode(inviteCode)
		if err != nil && err != gorm.ErrRecordNotFound {
			return nil, err
		}
		if inviterUser != nil {
			invitedBy = &inviterUser.Id
		} else {
			// 再检查管理员的邀请码
			inviterAdmin, err := s.svcCtx.AdminModel.GetByInviteCode(inviteCode)
			if err != nil && err != gorm.ErrRecordNotFound {
				return nil, err
			}
			if inviterAdmin != nil {
				invitedBy = &inviterAdmin.Id
			} else {
				return nil, fmt.Errorf("无效的邀请码")
			}
		}
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 生成用户邀请码
	userInviteCode := uuid.New().String()

	// 开始事务
	tx := s.svcCtx.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建用户
	newUser := &model.User{
		Username:   username,
		Password:   string(hashedPassword),
		Email:      email,
		IsActive:   true,
		InviteCode: userInviteCode,
		InvitedBy:  invitedBy,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	userModel := model.NewUserModel(tx)
	if err := userModel.Create(tx, newUser); err != nil {
		tx.Rollback()
		return nil, err
	}

	// 创建用户邮箱
	domain, err := s.svcCtx.DomainModel.GetById(domainID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("域名不存在")
	}

	fullEmail := fmt.Sprintf("%s@%s", domainPrefix, domain.Name)
	mailboxPassword := uuid.New().String()[:8] // 生成8位随机密码

	newMailbox := &model.Mailbox{
		UserId:    &newUser.Id,
		Email:     fullEmail,
		Password:  mailboxPassword,
		DomainId:  domainID,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mailboxModel := model.NewMailboxModel(tx)
	if err := mailboxModel.Create(tx, newMailbox); err != nil {
		tx.Rollback()
		return nil, err
	}

	// 如果有邀请人，增加邀请人的贡献度
	if invitedBy != nil {
		// 先尝试更新普通用户的贡献度
		userModel := model.NewUserModel(tx)
		inviterUser, err := userModel.GetById(*invitedBy)
		if err == nil {
			// 是普通用户，更新贡献度
			updateData := map[string]interface{}{
				"contribution": inviterUser.Contribution + 1,
				"updated_at":   time.Now(),
			}
			if err := userModel.MapUpdate(tx, *invitedBy, updateData); err != nil {
				tx.Rollback()
				return nil, err
			}
		} else {
			// 尝试更新管理员的贡献度
			adminModel := model.NewAdminModel(tx)
			inviterAdmin, err := adminModel.GetById(*invitedBy)
			if err == nil {
				updateData := map[string]interface{}{
					"contribution": inviterAdmin.Contribution + 1,
					"updated_at":   time.Now(),
				}
				if err := adminModel.MapUpdate(tx, *invitedBy, updateData); err != nil {
					tx.Rollback()
					return nil, err
				}
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	newUser.Password = ""
	return newUser, nil
}

// ChangePassword 修改密码
func (s *Service) ChangePassword(userID int64, oldPassword, newPassword string, isAdmin bool) error {
	if isAdmin {
		// 获取管理员信息
		admin, err := s.svcCtx.AdminModel.GetById(userID)
		if err != nil {
			return fmt.Errorf("管理员不存在")
		}

		// 验证旧密码
		if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(oldPassword)); err != nil {
			return fmt.Errorf("旧密码错误")
		}

		// 加密新密码
		newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		// 更新密码
		updateData := map[string]interface{}{
			"password":   string(newHashedPassword),
			"updated_at": time.Now(),
		}
		return s.svcCtx.AdminModel.MapUpdate(nil, userID, updateData)
	} else {
		// 获取用户信息
		user, err := s.svcCtx.UserModel.GetById(userID)
		if err != nil {
			return fmt.Errorf("用户不存在")
		}

		// 验证旧密码
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
			return fmt.Errorf("旧密码错误")
		}

		// 加密新密码
		newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		// 更新密码
		updateData := map[string]interface{}{
			"password":   string(newHashedPassword),
			"updated_at": time.Now(),
		}
		return s.svcCtx.UserModel.MapUpdate(nil, userID, updateData)
	}
}
