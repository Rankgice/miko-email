# NBEmail 外部SMTP发送功能指南

## 问题解决

✅ **已修复**：现在可以使用您域名的邮箱发送邮件到外部邮箱（如Gmail、QQ邮箱等）

## 功能特点

### ✅ 支持的发送方式
- **本地用户间发送**：系统内用户之间的邮件发送
- **外部邮箱发送**：从您的域名邮箱发送到任意外部邮箱

### ✅ 发件人地址
- 使用您域名下的真实邮箱作为发件人
- 收件人看到的发件人是您的域名邮箱（如 `user@yourdomain.com`）
- 支持回复到原始发件人

### ✅ SMTP服务器支持
- **推荐**：使用您自己域名的SMTP服务器
- **备选**：使用第三方SMTP服务器（Gmail、QQ邮箱等）

## 配置方法

### 方法1：环境变量配置

```bash
# Windows
set OUTBOUND_SMTP_HOST=mail.yourdomain.com
set OUTBOUND_SMTP_PORT=587
set OUTBOUND_SMTP_USER=smtp@yourdomain.com
set OUTBOUND_SMTP_PASSWORD=your-password
set OUTBOUND_SMTP_TLS=true

# Linux/Mac
export OUTBOUND_SMTP_HOST=mail.yourdomain.com
export OUTBOUND_SMTP_PORT=587
export OUTBOUND_SMTP_USER=smtp@yourdomain.com
export OUTBOUND_SMTP_PASSWORD=your-password
export OUTBOUND_SMTP_TLS=true
```

### 方法2：使用配置文件

复制并编辑配置模板：
```bash
cp smtp-config-example.env .env
# 编辑 .env 文件，取消注释并填写您的SMTP配置
```

### 方法3：使用启动脚本

```bash
# Windows
start-with-smtp.bat

# Linux/Mac
chmod +x start-with-smtp.sh
./start-with-smtp.sh
```

## 配置示例

### 1. 自己域名的SMTP服务器（推荐）

```bash
OUTBOUND_SMTP_HOST=mail.yourdomain.com
OUTBOUND_SMTP_PORT=587
OUTBOUND_SMTP_USER=smtp@yourdomain.com
OUTBOUND_SMTP_PASSWORD=your-smtp-password
OUTBOUND_SMTP_TLS=true
```

**优势**：
- ✅ 可以使用任意域名邮箱作为发件人
- ✅ 不受第三方邮件服务商限制
- ✅ 完全控制发送过程

### 2. 内网SMTP服务器（无认证）

```bash
OUTBOUND_SMTP_HOST=192.168.1.100
OUTBOUND_SMTP_PORT=25
OUTBOUND_SMTP_USER=
OUTBOUND_SMTP_PASSWORD=
OUTBOUND_SMTP_TLS=false
```

### 3. 第三方SMTP服务器

```bash
# Gmail
OUTBOUND_SMTP_HOST=smtp.gmail.com
OUTBOUND_SMTP_PORT=587
OUTBOUND_SMTP_USER=your-email@gmail.com
OUTBOUND_SMTP_PASSWORD=your-app-password
OUTBOUND_SMTP_TLS=true
```

## 使用流程

1. **配置外部SMTP**：按上述方法配置SMTP服务器
2. **启动系统**：运行 `./nbemail.exe` 或使用启动脚本
3. **登录Web界面**：访问 `http://localhost:8080`
4. **发送邮件**：
   - 选择发件人邮箱（必须是系统中的邮箱）
   - 输入外部收件人邮箱
   - 编写邮件内容
   - 点击发送

## 日志监控

启动时会显示SMTP配置状态：
```
2025/07/19 07:27:24 外部SMTP配置已加载 - Host: mail.yourdomain.com, Port: 587, User: smtp@yourdomain.com, TLS: true
```

发送时会记录发送结果：
```
2025/07/19 07:27:30 外部邮件发送成功 - From: user@yourdomain.com, To: someone@gmail.com, Subject: 测试邮件
```

## 故障排除

### 1. 连接失败
- 检查SMTP服务器地址和端口
- 确认网络连接正常
- 检查防火墙设置

### 2. 认证失败
- 验证用户名和密码
- 对于Gmail，使用应用专用密码
- 对于QQ/163邮箱，使用授权码

### 3. TLS失败
- 检查TLS配置
- 验证服务器证书
- 尝试关闭TLS（仅限内网环境）

## 技术实现

系统会自动判断收件人类型：
- **本地用户**：直接保存到收件箱
- **外部邮箱**：通过配置的SMTP服务器发送
- **发送失败**：邮件仍保存在发件箱，记录错误日志
