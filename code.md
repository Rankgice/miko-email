# AI编码规范文档

## 📋 目录
- [Go语言编码规范](#go语言编码规范)
- [数据库事务使用规范](#数据库事务使用规范)
- [Model层使用规范](#model层使用规范)
- [错误处理规范](#错误处理规范)
- [项目特定规范](#项目特定规范)

## Go语言编码规范

### 1. 必要的包导入
在使用GORM错误处理时，需要导入errors包：

```go
import (
    "errors"
    "gorm.io/gorm"
)
```

### 3. 变量声明优先级
当没有非err变量，或者有但变量作用域只在此if-else中时，优先使用以下写法：

```go
// ✅ 推荐写法
if o, err := newObject(); err != nil {
    log.Fatalln("连接数据库失败", "error", err.Error())
}

// ❌ 不推荐写法
o, err := newObject()
if err != nil {
    log.Fatalln("连接数据库失败", "error", err.Error())
}
```

### 4. 包管理
- 使用适当的包管理器进行依赖管理，而非手动编辑配置文件
- JavaScript/Node.js: 使用 `npm install`, `yarn add`, `pnpm add`
- Python: 使用 `pip install`, `poetry add`
- Go: 使用 `go get`, `go mod tidy`

## 数据库事务使用规范

### 标准事务模式
所有涉及多个数据库操作的业务逻辑都应使用以下标准事务模式：

```go
// 开始事务
tx := s.svcCtx.DB.Begin()
defer func() {
    if tx != nil {
        tx.Rollback()
    }
}()

// 业务代码
if err := s.svcCtx.UserModel.Create(tx, user); err != nil {
    return nil, err
}

if err := s.svcCtx.MailboxModel.Create(tx, mailbox); err != nil {
    return nil, err
}

// 提交事务
if err := tx.Commit().Error; err != nil {
    return nil, err
}
tx = nil

return result, nil
```

### 事务使用要点

1. **defer回滚**：使用defer确保异常情况下事务能够回滚
2. **提交后置nil**：提交成功后将tx设为nil，避免defer中重复回滚
3. **错误即返回**：任何操作失败立即返回，依赖defer进行回滚
4. **简洁调用**：直接使用svcCtx中的Model，不要创建新的Model实例

### 事务vs非事务

```go
// 单个操作 - 不使用事务
err := s.svcCtx.UserModel.Create(nil, user)

// 多个操作 - 使用事务
tx := s.svcCtx.DB.Begin()
defer func() {
    if tx != nil {
        tx.Rollback()
    }
}()

err := s.svcCtx.UserModel.Create(tx, user)
// ... 其他操作
tx.Commit()
tx = nil
```

## Model层使用规范

### 1. 方法参数规范
所有增删改方法的第一个参数必须是`tx *gorm.DB`：

```go
// ✅ 正确的方法签名
func (m *UserModel) Create(tx *gorm.DB, user *User) error
func (m *UserModel) Update(tx *gorm.DB, user *User) error
func (m *UserModel) Delete(tx *gorm.DB, user *User) error
```

### 2. 事务逻辑实现
所有增删改方法内部都应实现统一的事务逻辑：

```go
func (m *UserModel) Create(tx *gorm.DB, user *User) error {
    db := m.db
    if tx != nil {
        db = tx
    }
    return db.Create(user).Error
}
```

### 3. 表名规范
- 所有模型的`TableName()`方法返回单数形式表名
- 例如：`user`, `admin`, `domain`, `mailbox`, `email`, `email_forward`

### 4. 调用规范

```go
// ✅ 推荐写法 - 直接使用svcCtx中的Model
err := s.svcCtx.UserModel.Create(tx, user)

// ❌ 不推荐写法 - 创建新的Model实例
userModel := model.NewUserModel(tx)
err := userModel.Create(tx, user)
```

## 错误处理规范

### 1. 错误包装
使用`fmt.Errorf`进行错误包装，提供上下文信息：

```go
if err != nil {
    return fmt.Errorf("创建用户失败: %w", err)
}
```

### 2. GORM错误判断
使用`errors.Is`进行GORM错误判断，而非直接比较：

```go
// ✅ 推荐写法
if !errors.Is(err, gorm.ErrRecordNotFound) {
    return err
}

if errors.Is(err, gorm.ErrRecordNotFound) {
    return fmt.Errorf("记录不存在")
}

// ❌ 不推荐写法
if err != gorm.ErrRecordNotFound {
    return err
}

if err == gorm.ErrRecordNotFound {
    return fmt.Errorf("记录不存在")
}
```

**原因**: `errors.Is`能够正确处理错误包装链，更加健壮和准确。

### 3. 业务错误vs系统错误
- 业务错误：返回用户友好的错误信息
- 系统错误：记录详细日志，返回通用错误信息

```go
// 业务错误
if existingUser != nil {
    return fmt.Errorf("用户名已存在")
}

// 系统错误
if err != nil {
    log.Printf("数据库查询失败: %v", err)
    return fmt.Errorf("系统错误，请稍后重试")
}
```

## 项目特定规范

### 1. 密码处理
- 使用bcrypt进行密码加密
- 返回用户信息时清空密码字段

```go
// 加密密码
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

// 清空敏感信息
user.Password = ""
return user, nil
```

### 2. 时间处理
- 统一使用`time.Now()`设置创建和更新时间
- 数据库字段使用`CURRENT_TIMESTAMP`作为默认值

### 3. ID类型
- 统一使用`int64`作为ID类型
- GORM模型中使用`Id`字段名（首字母大写）

### 4. JSON标签
- 敏感字段使用`omitempty`标签
- 密码字段使用`json:"-"`或`json:"password,omitempty"`

### 5. 文档更新
- 对项目进行较大修改后，更新根目录的README.md文件
- 小改动不需要构建并测试

## 代码审查清单

### 提交前检查
- [ ] 是否遵循事务使用规范
- [ ] 是否正确处理错误
- [ ] 是否清空敏感信息
- [ ] 是否使用正确的Model调用方式
- [ ] 是否更新相关文档

### 性能考虑
- [ ] 避免在循环中进行数据库操作
- [ ] 合理使用事务，避免长时间锁定
- [ ] 查询时使用适当的索引

### 安全考虑
- [ ] 密码正确加密
- [ ] 输入参数验证
- [ ] SQL注入防护（GORM自动处理）
- [ ] 敏感信息不记录日志

---

**注意**: 此文档会随着项目发展持续更新，请定期查看最新版本。
