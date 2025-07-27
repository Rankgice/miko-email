# ServiceContext 使用说明

## 概述

`ServiceContext` 是项目的服务上下文，负责初始化和管理所有的数据库模型实例。所有的 model 层表模型都在此处统一初始化，方便在整个项目中使用。

## 结构说明

```go
type ServiceContext struct {
    Config            config.Config              // 配置信息
    DB                *gorm.DB                   // GORM 数据库连接
    ProxyModel        *model.ProxyModel          // 代理IP模型
    UserModel         *model.UserModel           // 用户模型
    AdminModel        *model.AdminModel          // 管理员模型
    DomainModel       *model.DomainModel         // 域名模型
    MailboxModel      *model.MailboxModel        // 邮箱模型
    EmailModel        *model.EmailModel          // 邮件模型
    EmailForwardModel *model.EmailForwardModel   // 邮件转发模型
}
```

## 初始化过程

1. **数据库连接**: 使用 SQLite 数据库，连接到 `data.db` 文件
2. **权限检测**: 检测数据库文件的读写权限
3. **自动迁移**: 自动创建/更新所有表结构
4. **模型初始化**: 初始化所有 model 实例

## 使用方式

### 方式1: 使用全局实例

```go
import "miko-email/internal/svc"

// 直接使用全局实例
func someHandler() {
    // 获取用户列表
    users, total, err := svc.SvcCtx.UserModel.List(model.UserReq{
        Page:     1,
        PageSize: 10,
    })
    
    // 创建新用户
    user := &model.User{
        Username: "test",
        Email:    "test@example.com",
    }
    err = svc.SvcCtx.UserModel.Create(user)
}
```

### 方式2: 依赖注入

```go
type Handler struct {
    svcCtx *svc.ServiceContext
}

func NewHandler(svcCtx *svc.ServiceContext) *Handler {
    return &Handler{
        svcCtx: svcCtx,
    }
}

func (h *Handler) GetUsers() {
    users, total, err := h.svcCtx.UserModel.List(model.UserReq{
        IsActive: &[]bool{true}[0],
        Page:     1,
        PageSize: 10,
    })
    // 处理结果...
}
```

## 自动迁移

系统会自动迁移以下表结构：
- `proxy` - 代理IP表
- `users` - 用户表
- `admins` - 管理员表
- `domains` - 域名表
- `mailboxes` - 邮箱表
- `emails` - 邮件表
- `email_forwards` - 邮件转发表

## 模型使用示例

### 用户管理

```go
// 创建用户
user := &model.User{
    Username:   "newuser",
    Email:      "user@example.com",
    IsActive:   true,
    InviteCode: "INVITE123",
}
err := svc.SvcCtx.UserModel.Create(user)

// 根据用户名查询
user, err := svc.SvcCtx.UserModel.GetByUsername("newuser")

// 分页查询用户
params := model.UserReq{
    IsActive: &[]bool{true}[0],
    Page:     1,
    PageSize: 20,
}
users, total, err := svc.SvcCtx.UserModel.List(params)
```

### 邮箱管理

```go
// 创建邮箱
mailbox := &model.Mailbox{
    UserId:   &userId,
    Email:    "user@domain.com",
    Password: "encrypted_password",
    DomainId: domainId,
    IsActive: true,
}
err := svc.SvcCtx.MailboxModel.Create(mailbox)

// 获取用户的邮箱列表
mailboxes, err := svc.SvcCtx.MailboxModel.GetMailboxesByUserId(userId)
```

### 邮件管理

```go
// 保存邮件
email := &model.Email{
    MailboxId: mailboxId,
    FromAddr:  "sender@example.com",
    ToAddr:    "receiver@example.com",
    Subject:   "邮件主题",
    Body:      "邮件内容",
    Folder:    "inbox",
}
err := svc.SvcCtx.EmailModel.Create(email)

// 获取邮箱的邮件列表
emails, total, err := svc.SvcCtx.EmailModel.GetEmailsByMailboxId(
    mailboxId, "inbox", 1, 20,
)

// 标记邮件为已读
err := svc.SvcCtx.EmailModel.MarkAsRead(emailId)
```

### 域名管理

```go
// 创建域名
domain := &model.Domain{
    Name:       "example.com",
    IsVerified: false,
    IsActive:   true,
    MxRecord:   "mail.example.com",
    ARecord:    "192.168.1.1",
}
err := svc.SvcCtx.DomainModel.Create(domain)

// 获取可用域名列表
domains, err := svc.SvcCtx.DomainModel.GetAvailableDomains()
```

### 邮件转发管理

```go
// 创建转发规则
forward := &model.EmailForward{
    MailboxId:          mailboxId,
    SourceEmail:        "source@example.com",
    TargetEmail:        "target@example.com",
    Enabled:            true,
    KeepOriginal:       true,
    ForwardAttachments: true,
    SubjectPrefix:      "[转发]",
}
err := svc.SvcCtx.EmailForwardModel.Create(forward)

// 获取邮箱的转发规则
forwards, err := svc.SvcCtx.EmailForwardModel.GetForwardsByMailboxId(mailboxId)
```

## 事务支持

所有模型都支持事务操作：

```go
// 开始事务
tx := svc.SvcCtx.DB.Begin()

// 在事务中操作
user := &model.User{Username: "test", Email: "test@example.com"}
err := svc.SvcCtx.UserModel.Create(user)
if err != nil {
    tx.Rollback()
    return err
}

mailbox := &model.Mailbox{UserId: &user.Id, Email: "test@domain.com"}
err = svc.SvcCtx.MailboxModel.Create(mailbox)
if err != nil {
    tx.Rollback()
    return err
}

// 提交事务
tx.Commit()
```

## 注意事项

1. **全局实例**: `svc.SvcCtx` 是全局实例，在 `init()` 函数中自动初始化
2. **数据库文件**: 默认使用 `data.db` 作为数据库文件
3. **自动迁移**: 每次启动时会自动检查并更新表结构
4. **权限检测**: 启动时会检测数据库文件的读写权限
5. **错误处理**: 初始化失败会导致程序退出，请确保数据库配置正确

## 扩展说明

如果需要添加新的模型：

1. 在 `internal/model/` 中创建新的模型文件
2. 在 `ServiceContext` 结构体中添加新的模型字段
3. 在 `NewServiceContext` 函数中初始化新模型
4. 在 `autoMigrate` 函数中添加新模型的迁移
