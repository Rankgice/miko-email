# Miko邮箱系统 API 接口文档

## 基础信息

- **基础URL**: `http://localhost:8080`
- **内容类型**: `application/json; charset=utf-8`
- **认证方式**: Session Cookie (登录后自动设置)
- **响应格式**: 统一JSON格式

## 通用响应格式

### 成功响应
```json
{
    "success": true,
    "message": "操作成功",
    "data": {...}
}
```

### 错误响应
```json
{
    "success": false,
    "message": "错误信息"
}
```

## 1. 认证相关 API

### 1.1 用户登录
- **URL**: `POST /api/login`
- **描述**: 用户登录获取会话
- **请求体**:
```json
{
    "username": "testuser",
    "password": "password123"
}
```
- **响应**:
```json
{
    "success": true,
    "message": "登录成功",
    "data": {
        "user": {
            "id": 1,
            "username": "testuser",
            "email": "test@example.com",
            "status": "active",
            "contribution": 0,
            "invite_code": "ABC123",
            "created_at": "2025-01-01T00:00:00Z"
        }
    }
}
```

### 1.2 用户注册
- **URL**: `POST /api/register`
- **描述**: 用户注册
- **请求体**:
```json
{
    "username": "newuser",
    "password": "password123",
    "email": "newuser@example.com",
    "domain_prefix": "newuser",
    "domain_id": 1,
    "invite_code": "INVITE123"
}
```

### 1.3 管理员登录
- **URL**: `POST /api/admin/login`
- **描述**: 管理员登录
- **请求体**:
```json
{
    "username": "admin",
    "password": "admin123"
}
```

### 1.4 用户登出
- **URL**: `POST /api/logout`
- **描述**: 用户登出清除会话
- **响应**:
```json
{
    "success": true,
    "message": "登出成功"
}
```

### 1.5 获取用户信息
- **URL**: `GET /api/profile`
- **描述**: 获取当前登录用户信息
- **需要认证**: 是

### 1.6 修改密码
- **URL**: `PUT /api/profile/password`
- **描述**: 修改用户密码
- **需要认证**: 是
- **请求体**:
```json
{
    "old_password": "oldpass123",
    "new_password": "newpass123"
}
```

## 2. 邮箱管理 API

### 2.1 获取邮箱列表
- **URL**: `GET /api/mailboxes`
- **描述**: 获取用户的邮箱列表
- **需要认证**: 是
- **响应**:
```json
{
    "success": true,
    "data": [
        {
            "id": 1,
            "email": "user@example.com",
            "domain_id": 1,
            "status": "active",
            "created_at": "2025-01-01T00:00:00Z"
        }
    ]
}
```

### 2.2 创建邮箱
- **URL**: `POST /api/mailboxes`
- **描述**: 创建单个邮箱
- **需要认证**: 是
- **请求体**:
```json
{
    "prefix": "newbox",
    "domain_id": 1
}
```

### 2.3 批量创建邮箱
- **URL**: `POST /api/mailboxes/batch`
- **描述**: 批量创建邮箱
- **需要认证**: 是
- **请求体**:
```json
{
    "prefixes": ["box1", "box2", "box3"],
    "domain_id": 1
}
```

### 2.4 获取邮箱密码
- **URL**: `GET /api/mailboxes/:id/password`
- **描述**: 获取指定邮箱的密码
- **需要认证**: 是

### 2.5 删除邮箱
- **URL**: `DELETE /api/mailboxes/:id`
- **描述**: 删除指定邮箱
- **需要认证**: 是

## 3. 邮件管理 API

### 3.1 获取邮件列表
- **URL**: `GET /api/emails`
- **描述**: 获取邮件列表
- **需要认证**: 是
- **查询参数**:
  - `mailbox`: 邮箱地址 (可选)
  - `type`: 邮件类型 (inbox/sent/trash, 默认: inbox)
  - `page`: 页码 (默认: 1)
  - `limit`: 每页数量 (默认: 20, 最大: 100)
- **响应**:
```json
{
    "success": true,
    "data": [
        {
            "id": 1,
            "from_addr": "sender@example.com",
            "to_addr": "receiver@example.com",
            "subject": "邮件主题",
            "body": "邮件内容",
            "is_read": false,
            "created_at": "2025-01-01T00:00:00Z"
        }
    ],
    "total": 100,
    "page": 1,
    "limit": 20
}
```

### 3.2 获取单个邮件
- **URL**: `GET /api/emails/:id`
- **描述**: 获取指定邮件详情
- **需要认证**: 是

### 3.3 发送邮件
- **URL**: `POST /api/emails/send`
- **描述**: 发送邮件
- **需要认证**: 是
- **请求体**:
```json
{
    "from": "sender@example.com",
    "to": "receiver@example.com,another@example.com",
    "subject": "邮件主题",
    "content": "邮件内容"
}
```

## 4. 转发规则管理 API

### 4.1 获取转发规则列表
- **URL**: `GET /api/forward-rules`
- **描述**: 获取用户的转发规则列表
- **需要认证**: 是

### 4.2 创建转发规则
- **URL**: `POST /api/forward-rules`
- **描述**: 创建新的转发规则
- **需要认证**: 是
- **请求体**:
```json
{
    "source_email": "source@example.com",
    "target_email": "target@example.com",
    "enabled": true,
    "keep_original": true,
    "forward_attachments": true,
    "subject_prefix": "[转发]",
    "description": "转发规则描述"
}
```

### 4.3 获取单个转发规则
- **URL**: `GET /api/forward-rules/:id`
- **描述**: 获取指定转发规则详情
- **需要认证**: 是

### 4.4 更新转发规则
- **URL**: `PUT /api/forward-rules/:id`
- **描述**: 更新转发规则
- **需要认证**: 是

### 4.5 删除转发规则
- **URL**: `DELETE /api/forward-rules/:id`
- **描述**: 删除转发规则
- **需要认证**: 是

### 4.6 切换转发规则状态
- **URL**: `PATCH /api/forward-rules/:id/toggle`
- **描述**: 启用/禁用转发规则
- **需要认证**: 是

### 4.7 测试转发规则
- **URL**: `POST /api/forward-rules/:id/test`
- **描述**: 测试转发规则
- **需要认证**: 是

### 4.8 获取转发统计
- **URL**: `GET /api/forward-statistics`
- **描述**: 获取转发统计信息
- **需要认证**: 是

## 5. 域名管理 API (公共)

### 5.1 获取可用域名
- **URL**: `GET /api/domains/available`
- **描述**: 获取可用于注册的域名列表
- **需要认证**: 否

### 5.2 获取域名DNS记录
- **URL**: `GET /api/domains/dns`
- **描述**: 查询域名DNS记录
- **需要认证**: 否
- **查询参数**:
  - `domain`: 域名

## 6. 管理员 API (需要管理员权限)

### 6.1 域名管理

#### 6.1.1 获取域名列表
- **URL**: `GET /api/admin/domains`
- **描述**: 获取所有域名列表
- **需要认证**: 是 (管理员)

#### 6.1.2 创建域名
- **URL**: `POST /api/admin/domains`
- **描述**: 创建新域名
- **需要认证**: 是 (管理员)
- **请求体**:
```json
{
    "name": "example.com",
    "mx_record": "mail.example.com",
    "a_record": "192.168.1.1",
    "txt_record": "v=spf1 include:example.com ~all"
}
```

#### 6.1.3 更新域名
- **URL**: `PUT /api/admin/domains/:id`
- **描述**: 更新域名信息
- **需要认证**: 是 (管理员)

#### 6.1.4 删除域名
- **URL**: `DELETE /api/admin/domains/:id`
- **描述**: 删除域名
- **需要认证**: 是 (管理员)

#### 6.1.5 验证域名
- **URL**: `POST /api/admin/domains/:id/verify`
- **描述**: 验证域名DNS配置
- **需要认证**: 是 (管理员)

### 6.2 用户管理

#### 6.2.1 获取用户列表
- **URL**: `GET /api/admin/users`
- **描述**: 获取所有用户列表
- **需要认证**: 是 (管理员)
- **响应**:
```json
{
    "success": true,
    "data": [
        {
            "id": 1,
            "username": "testuser",
            "email": "test@example.com",
            "status": "active",
            "contribution": 0,
            "mailbox_count": 3,
            "inviter_name": "admin",
            "created_at": "2025-01-01T00:00:00Z"
        }
    ]
}
```

#### 6.2.2 获取单个用户
- **URL**: `GET /api/admin/users/:id`
- **描述**: 获取指定用户详情
- **需要认证**: 是 (管理员)

#### 6.2.3 获取用户邮箱
- **URL**: `GET /api/admin/users/:id/mailboxes`
- **描述**: 获取指定用户的邮箱列表
- **需要认证**: 是 (管理员)

#### 6.2.4 更新用户状态
- **URL**: `PUT /api/admin/users/:id/status`
- **描述**: 更新用户状态 (活跃/暂停)
- **需要认证**: 是 (管理员)
- **请求体**:
```json
{
    "status": "active"
}
```
- **状态值**:
  - `"active"`: 活跃状态
  - `"suspended"`: 暂停状态
- **响应**:
```json
{
    "success": true,
    "message": "用户状态更新成功"
}
```

#### 6.2.5 删除用户
- **URL**: `DELETE /api/admin/users/:id`
- **描述**: 删除用户 (软删除)
- **需要认证**: 是 (管理员)

### 6.3 用户状态管理说明

#### 6.3.1 用户状态类型
- **active**: 活跃状态
  - 用户可以正常登录
  - 可以使用所有功能
  - 可以创建和管理邮箱
  - 可以收发邮件

- **suspended**: 暂停状态
  - 用户无法登录
  - 现有会话会被终止
  - 邮箱仍然可以接收邮件
  - 无法发送邮件或管理邮箱

- **deleted**: 已删除状态
  - 软删除，数据保留但不可访问
  - 用户无法登录
  - 邮箱停止接收邮件

#### 6.3.2 状态切换规则
- `active` ↔ `suspended`: 可以相互切换
- `active/suspended` → `deleted`: 可以删除
- `deleted` → `active/suspended`: 需要特殊恢复操作

## 7. 错误代码说明

| HTTP状态码 | 说明 | 常见原因 |
|-----------|------|---------|
| 200 | 请求成功 | - |
| 400 | 请求参数错误 | 参数格式不正确、必需参数缺失 |
| 401 | 未认证或认证失败 | 未登录、会话过期、用户被暂停 |
| 403 | 权限不足 | 非管理员访问管理员接口 |
| 404 | 资源不存在 | 用户ID不存在、邮箱ID不存在 |
| 500 | 服务器内部错误 | 数据库连接失败、系统异常 |

### 7.1 常见错误响应示例

#### 7.1.1 用户未找到
```json
{
    "success": false,
    "message": "用户不存在"
}
```

#### 7.1.2 用户状态错误
```json
{
    "success": false,
    "message": "用户已被暂停"
}
```

#### 7.1.3 权限不足
```json
{
    "success": false,
    "message": "权限不足"
}
```

## 8. 认证说明

### Session认证
- 登录成功后，服务器会设置名为 `miko-session` 的Cookie
- 后续请求会自动携带此Cookie进行身份验证
- 登出时会清除此Cookie

### 权限级别
1. **公共接口**: 无需认证
2. **用户接口**: 需要用户登录
3. **管理员接口**: 需要管理员权限

## 9. 使用示例

### JavaScript (Axios)
```javascript
// 登录
const loginResponse = await axios.post('/api/login', {
    username: 'testuser',
    password: 'password123'
});

// 获取邮箱列表
const mailboxes = await axios.get('/api/mailboxes');

// 发送邮件
const sendResponse = await axios.post('/api/emails/send', {
    from: 'sender@example.com',
    to: 'receiver@example.com',
    subject: '测试邮件',
    content: '这是一封测试邮件'
});
```

### cURL
```bash
# 登录
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}' \
  -c cookies.txt

# 获取邮箱列表
curl -X GET http://localhost:8080/api/mailboxes \
  -b cookies.txt

# 创建邮箱
curl -X POST http://localhost:8080/api/mailboxes \
  -H "Content-Type: application/json" \
  -d '{"prefix":"newbox","domain_id":1}' \
  -b cookies.txt
```

## 10. 注意事项

1. **字符编码**: 所有请求和响应均使用UTF-8编码
2. **时间格式**: 使用ISO 8601格式 (RFC3339)
3. **分页**: 默认每页20条记录，最大100条
4. **邮件地址**: 支持多个收件人，用逗号分隔
5. **域名验证**: 创建域名后需要进行DNS验证才能使用
6. **转发规则**: 支持多种转发选项，包括保留原件、转发附件等

## 11. 数据模型

### 11.1 用户模型 (User)
```json
{
    "id": 1,
    "username": "testuser",
    "email": "test@example.com",
    "status": "active",
    "contribution": 100,
    "invite_code": "ABC123",
    "invited_by": 1,
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
}
```
**状态说明**:
- `"active"`: 活跃用户，可以正常使用所有功能
- `"suspended"`: 暂停用户，无法登录和使用功能
- `"deleted"`: 已删除用户（软删除）

### 11.2 邮箱模型 (Mailbox)
```json
{
    "id": 1,
    "user_id": 1,
    "admin_id": null,
    "email": "user@example.com",
    "domain_id": 1,
    "status": "active",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
}
```

### 11.3 域名模型 (Domain)
```json
{
    "id": 1,
    "name": "example.com",
    "mx_record": "mail.example.com",
    "a_record": "192.168.1.1",
    "txt_record": "v=spf1 include:example.com ~all",
    "is_verified": true,
    "is_active": true,
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
}
```

### 11.4 邮件模型 (Email)
```json
{
    "id": 1,
    "mailbox_id": 1,
    "from_addr": "sender@example.com",
    "to_addr": "receiver@example.com",
    "subject": "邮件主题",
    "body": "邮件内容",
    "folder": "inbox",
    "is_read": false,
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
}
```

### 11.5 转发规则模型 (ForwardRule)
```json
{
    "id": 1,
    "mailbox_id": 1,
    "source_email": "source@example.com",
    "target_email": "target@example.com",
    "enabled": true,
    "keep_original": true,
    "forward_attachments": true,
    "subject_prefix": "[转发]",
    "description": "转发规则描述",
    "forward_count": 10,
    "last_forwarded_at": "2025-01-01T00:00:00Z",
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
}
```

## 12. 邮件服务器配置

### 12.1 SMTP配置
- **服务器**: localhost (或您的域名)
- **端口**: 25
- **加密**: 无
- **认证**: PLAIN
- **用户名**: 您的用户名
- **密码**: 您的密码

### 12.2 IMAP配置
- **服务器**: localhost (或您的域名)
- **端口**: 143
- **加密**: 无
- **认证**: PLAIN
- **用户名**: 您的用户名
- **密码**: 您的密码

### 12.3 POP3配置
- **服务器**: localhost (或您的域名)
- **端口**: 110
- **加密**: 无
- **认证**: PLAIN
- **用户名**: 您的用户名
- **密码**: 您的密码

## 13. 开发环境搭建

### 13.1 依赖要求
- Go 1.19+
- SQLite3
- Git

### 13.2 启动服务
```bash
# 克隆项目
git clone <repository-url>
cd miko-email

# 编译项目
go build -o nbemail.exe .

# 运行服务
./nbemail.exe
```

### 13.3 默认端口
- **Web服务**: 8080
- **SMTP服务**: 25
- **IMAP服务**: 143
- **POP3服务**: 110

## 14. API测试工具推荐

### 14.1 Postman
推荐使用Postman进行API测试，可以方便地管理Cookie和会话。

### 14.2 cURL
适合命令行测试和脚本自动化。

### 14.3 JavaScript/Axios
适合前端开发和集成测试。

## 15. 常见问题

### 15.1 认证失败
- 检查用户名和密码是否正确
- 确认用户账号是否已激活
- 检查Cookie是否正确设置

### 15.2 邮箱创建失败
- 检查域名是否已验证
- 确认邮箱前缀是否已存在
- 检查用户权限

### 15.3 邮件发送失败
- 检查发件邮箱是否属于当前用户
- 确认收件人邮箱格式是否正确
- 检查SMTP服务是否正常运行

### 15.4 转发规则不生效
- 检查转发规则是否已启用
- 确认源邮箱和目标邮箱是否正确
- 检查邮件转发服务是否正常运行

### 15.5 用户管理功能异常
- **暂停/启用按钮不工作**: 检查用户ID类型匹配，确保JavaScript中正确处理数字和字符串转换
- **邮箱收件箱按钮无法跳转**: 检查邮箱ID类型匹配，确保前端正确查找邮箱对象
- **用户状态显示不正确**: 检查API返回的状态字段名称，前端应使用 `status` 而不是 `is_active`

### 15.6 前端显示问题
- **状态徽章显示错误**: 确保前端代码使用正确的状态字段名称
- **按钮状态不更新**: 检查页面刷新逻辑，确保操作后重新加载数据
- **类型匹配错误**: JavaScript中字符串和数字比较时使用 `parseInt()` 进行类型转换

## 16. 更新日志

### v1.0.1 (2025-07-22)
- **修复**: 用户管理页面暂停/启用按钮功能
- **修复**: 邮箱管理页面收件箱按钮跳转问题
- **修复**: 仪表板页面邮箱状态显示问题
- **改进**: 统一用户状态字段为 `status`，支持 `active`/`suspended`/`deleted` 状态
- **改进**: 增强前端类型匹配处理，避免字符串和数字比较错误
- **改进**: 优化按钮显示逻辑，根据状态动态显示不同图标和文本
- **文档**: 更新API文档，完善用户状态管理说明

### v1.0.0 (2025-01-01)
- 初始版本发布
- 支持用户注册、登录、邮箱管理
- 支持邮件收发、转发规则
- 支持域名管理和DNS验证
- 支持管理员功能

---

**文档版本**: v1.0.1
**最后更新**: 2025-07-22
**维护者**: Miko邮箱系统开发团队
