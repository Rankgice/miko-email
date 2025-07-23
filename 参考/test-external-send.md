# 测试外部SMTP发送功能

## 配置说明

现在系统支持使用自己域名的SMTP服务器发送邮件到外部邮箱。

### 1. 配置自己域名的SMTP服务器

编辑环境变量或创建 `.env` 文件：

```bash
# 使用您自己的SMTP服务器
OUTBOUND_SMTP_HOST=mail.yourdomain.com
OUTBOUND_SMTP_PORT=587
OUTBOUND_SMTP_USER=smtp@yourdomain.com
OUTBOUND_SMTP_PASSWORD=your-smtp-password
OUTBOUND_SMTP_TLS=true
```

### 2. 发送邮件测试

通过Web界面发送邮件时：
- **发件人**：可以使用您域名下的任意邮箱（如 `user@yourdomain.com`）
- **收件人**：可以是任意外部邮箱（如 `someone@gmail.com`）

### 3. 工作原理

1. 用户在Web界面选择发件人邮箱（必须是系统中的邮箱）
2. 系统检查收件人是否为本地用户
3. 如果是外部邮箱，系统使用配置的SMTP服务器发送
4. 发件人地址保持为原始的域名邮箱

### 4. 优势

- ✅ 可以使用您域名的任意邮箱作为发件人
- ✅ 收件人看到的发件人是真实的域名邮箱
- ✅ 支持回复到原始发件人
- ✅ 不依赖第三方邮件服务商的限制

### 5. 配置示例

#### 企业邮箱服务器
```bash
OUTBOUND_SMTP_HOST=smtp.company.com
OUTBOUND_SMTP_PORT=587
OUTBOUND_SMTP_USER=noreply@company.com
OUTBOUND_SMTP_PASSWORD=smtp-password
OUTBOUND_SMTP_TLS=true
```

#### 内网SMTP服务器（无认证）
```bash
OUTBOUND_SMTP_HOST=192.168.1.100
OUTBOUND_SMTP_PORT=25
OUTBOUND_SMTP_USER=
OUTBOUND_SMTP_PASSWORD=
OUTBOUND_SMTP_TLS=false
```

### 6. 故障排除

如果发送失败，检查日志中的错误信息：
- 连接失败：检查SMTP服务器地址和端口
- 认证失败：检查用户名和密码
- TLS失败：检查TLS配置和证书

### 7. 测试命令

启动服务时设置环境变量：
```bash
# Windows
set OUTBOUND_SMTP_HOST=mail.yourdomain.com && set OUTBOUND_SMTP_PORT=587 && ./nbemail.exe

# Linux/Mac
OUTBOUND_SMTP_HOST=mail.yourdomain.com OUTBOUND_SMTP_PORT=587 ./nbemail
```
