# Model 包说明文档

本包包含了项目中所有数据表的 GORM 模型定义，统一了数据库操作接口。

## 文件列表

### types.go (新增)
- **功能**: 统一管理所有查询参数结构体
- **包含**: 所有 `xxxReq` 查询参数结构体定义

### 1. proxy.go (已存在)
- **模型**: `Proxy` - 代理IP模型
- **表名**: `proxy`
- **功能**: 代理IP管理，包括IP地址、端口、用户名、密码、类型等
- **特殊方法**: `FormattedProxy()` - 获取格式化的代理字符串

### 2. user.go (新生成)
- **模型**: `User` - 普通用户模型
- **表名**: `users`
- **功能**: 普通用户管理，包括用户名、密码、邮箱、激活状态、贡献度、邀请码等
- **特殊方法**: 
  - `GetByUsername()` - 根据用户名获取用户
  - `GetByEmail()` - 根据邮箱获取用户
  - `GetByInviteCode()` - 根据邀请码获取用户
  - `AuthenticateUser()` - 验证用户登录

### 3. admin.go (新生成)
- **模型**: `Admin` - 管理员模型
- **表名**: `admins`
- **功能**: 管理员用户管理，字段与普通用户类似但无邀请人字段
- **特殊方法**:
  - `GetByUsername()` - 根据用户名获取管理员
  - `GetByEmail()` - 根据邮箱获取管理员
  - `AuthenticateAdmin()` - 验证管理员登录
  - `UpdatePassword()` - 更新管理员密码

### 4. domain.go (新生成)
- **模型**: `Domain` - 域名模型
- **表名**: `domains`
- **功能**: 域名管理，包括域名、验证状态、激活状态、DNS记录等
- **特殊方法**:
  - `GetByName()` - 根据域名获取记录
  - `GetActiveDomains()` - 获取激活的域名列表
  - `GetVerifiedDomains()` - 获取已验证的域名列表
  - `GetAvailableDomains()` - 获取可用的域名列表
  - `UpdateDNSRecords()` - 更新DNS记录

### 5. mailbox.go (新生成)
- **模型**: `Mailbox` - 邮箱模型
- **表名**: `mailboxes`
- **功能**: 邮箱管理，包括用户ID、管理员ID、邮箱地址、密码、域名ID等
- **特殊方法**:
  - `GetByEmail()` - 根据邮箱地址获取记录
  - `GetMailboxesByUserId()` - 根据用户ID获取邮箱列表
  - `GetMailboxesByAdminId()` - 根据管理员ID获取邮箱列表
  - `GetMailboxesByDomainId()` - 根据域名ID获取邮箱列表
  - `CountMailboxesByUserId()` - 统计用户的邮箱数量

### 6. email.go (新生成)
- **模型**: `Email` - 邮件模型
- **表名**: `emails`
- **功能**: 邮件管理，包括邮箱ID、发件人、收件人、主题、内容、已读状态、文件夹等
- **特殊方法**:
  - `GetByIdAndMailboxId()` - 根据ID和邮箱ID获取邮件
  - `MarkAsRead()/MarkAsUnread()` - 标记已读/未读
  - `MoveToFolder()` - 移动到指定文件夹
  - `GetEmailsByMailboxId()` - 根据邮箱ID获取邮件列表
  - `GetUnreadCount()` - 获取未读邮件数量
  - `SearchEmails()` - 搜索邮件

### 7. email_forward.go (新生成)
- **模型**: `EmailForward` - 邮件转发模型
- **表名**: `email_forwards`
- **功能**: 邮件转发规则管理，包括源邮箱、目标邮箱、启用状态、转发设置等
- **特殊方法**:
  - `GetForwardsByMailboxId()` - 根据邮箱ID获取转发规则
  - `GetEnabledForwardsByMailboxId()` - 获取启用的转发规则
  - `GetForwardsBySourceEmail()` - 根据源邮箱获取转发规则
  - `IncrementForwardCount()` - 增加转发次数
  - `GetForwardStatistics()` - 获取转发统计信息

## 统一的方法接口

每个模型都包含以下标准方法：

### 基础 CRUD 方法
- `Create(model)` - 创建记录
- `Update(tx, model)` - 更新记录
- `MapUpdate(tx, id, data)` - 使用 map 更新记录
- `Save(tx, model)` - 保存记录
- `Delete(model)` - 删除记录
- `GetById(id)` - 根据ID获取记录

### 查询方法
- `List(params)` - 统一的列表查询方法，支持条件查询和分页
- `BatchDelete(ids)` - 批量删除

### 状态管理方法
- `UpdateStatus(id, status)` - 更新状态（适用于有状态字段的表）
- `CheckXXXExist()` - 检查记录是否存在

## 字段命名规范

根据用户要求，所有形如 `xxxID` 的字段都改为 `xxxId`：
- `UserId` 而不是 `UserID`
- `AdminId` 而不是 `AdminID`  
- `DomainId` 而不是 `DomainID`
- `MailboxId` 而不是 `MailboxID`

## GORM 标签说明

每个字段都使用了完整的 GORM 标签：
- `column` - 指定数据库列名
- `primaryKey` - 主键标识
- `autoIncrement` - 自增标识
- `uniqueIndex` - 唯一索引
- `not null` - 非空约束
- `default` - 默认值
- `comment` - 字段注释

## 使用示例

```go
// 创建模型实例
db := gorm.Open(...)
userModel := NewUserModel(db)

// 创建用户
user := &User{
    Username: "test",
    Email: "test@example.com",
    // ...
}
err := userModel.Create(user)

// 查询用户
user, err := userModel.GetByUsername("test")

// 列表查询
params := UserReq{
    IsActive: &[]bool{true}[0],
    Page: 1,
    PageSize: 10,
}
users, total, err := userModel.List(params)
```

## 注意事项

1. 所有模型都支持事务操作，通过 `tx` 参数传入事务对象
2. 查询参数中的指针类型字段（如 `*bool`）用于区分零值和未设置
3. 分页查询当 `Page` 或 `PageSize` 为 0 时不进行分页
4. 所有时间字段使用 `time.Time` 类型
5. 密码等敏感字段在 JSON 序列化时会被忽略（`omitempty` 标签）

## 数据库表关系

- `users` ↔ `mailboxes` (一对多)
- `admins` ↔ `mailboxes` (一对多)
- `domains` ↔ `mailboxes` (一对多)
- `mailboxes` ↔ `emails` (一对多)
- `mailboxes` ↔ `email_forwards` (一对多)
- `users` ↔ `users` (自关联，邀请关系)

## 重要变更说明

### 表名保持兼容
- 保持与现有数据库的兼容性，表名仍使用复数形式
- 例如：`users`, `admins`, `domains`, `mailboxes`, `emails`, `email_forwards`

### 查询参数结构体变更
- 所有查询参数结构体从各自的模型文件移动到 `types.go`
- 命名从 `xxxQueryParams` 改为 `xxxReq`
- 例如：`UserQueryParams` → `UserReq`

### 文件结构优化
- 新增 `types.go` 统一管理查询参数结构体
- 各模型文件只保留核心的模型定义和方法
- 提高了代码的组织性和可维护性
