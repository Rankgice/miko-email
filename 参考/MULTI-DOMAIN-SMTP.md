# NBEmail 多域名SMTP配置指南

## 🎯 功能概述

NBEmail 支持多域名SMTP配置，允许您为不同的域名配置不同的SMTP服务器。这意味着：

- **用什么域名的邮箱发件，就用对应域名的SMTP服务器**
- **支持无限数量的域名配置**
- **自动根据发件人邮箱选择对应的SMTP配置**
- **向后兼容原有的单一SMTP配置**

## 🚀 使用场景

### 场景1：多品牌公司
```
sales@company-a.com    -> 使用 company-a.com 的SMTP服务器
support@company-b.com  -> 使用 company-b.com 的SMTP服务器
info@company-c.com     -> 使用 company-c.com 的SMTP服务器
```

### 场景2：不同部门使用不同邮件服务商
```
hr@company.com         -> 使用腾讯企业邮箱SMTP
tech@company.com       -> 使用阿里云邮箱SMTP  
marketing@company.com  -> 使用Gmail SMTP
```

### 场景3：个人多域名管理
```
me@personal.com        -> 使用个人域名SMTP
work@company.com       -> 使用公司域名SMTP
blog@myblog.org        -> 使用博客域名SMTP
```

## ⚙️ 配置方法

### 方法1：环境变量配置

在 `.env` 文件中添加：

```bash
# 域名SMTP配置格式：
# DOMAIN_SMTP_<域名>=host:port:user:password:tls
# 注意：域名中的点(.)要替换为下划线(_)

# example.com 的SMTP配置
DOMAIN_SMTP_EXAMPLE_COM=mail.example.com:587:smtp@example.com:password123:true

# company.com 的SMTP配置  
DOMAIN_SMTP_COMPANY_COM=smtp.company.com:465:noreply@company.com:company-pass:true

# myblog.org 的SMTP配置
DOMAIN_SMTP_MYBLOG_ORG=mail.myblog.org:587:sender@myblog.org:blog-pass:true

# 默认SMTP配置（可选，向后兼容）
OUTBOUND_SMTP_HOST=mail.default.com
OUTBOUND_SMTP_PORT=587
OUTBOUND_SMTP_USER=default@default.com
OUTBOUND_SMTP_PASSWORD=default-password
OUTBOUND_SMTP_TLS=true
```

### 方法2：Web界面配置

1. 登录NBEmail管理界面
2. 访问 "SMTP配置" 页面
3. 点击 "添加SMTP配置"
4. 填写域名和SMTP服务器信息
5. 保存配置

## 🔧 域名转换规则

环境变量中的域名需要按以下规则转换：

| 原域名 | 环境变量名 |
|--------|------------|
| example.com | DOMAIN_SMTP_EXAMPLE_COM |
| my-site.org | DOMAIN_SMTP_MY_SITE_ORG |
| test.co.uk | DOMAIN_SMTP_TEST_CO_UK |
| 中文域名.com | DOMAIN_SMTP_中文域名_COM |

**转换规则：**
1. 将点(.)替换为下划线(_)
2. 将连字符(-)替换为下划线(_)
3. 转换为大写字母
4. 添加 `DOMAIN_SMTP_` 前缀

## 📋 常见SMTP服务器配置

### Gmail
```bash
DOMAIN_SMTP_GMAIL_COM=smtp.gmail.com:587:your-email@gmail.com:app-password:true
```

### QQ邮箱
```bash
DOMAIN_SMTP_QQ_COM=smtp.qq.com:587:your-email@qq.com:authorization-code:true
```

### 163邮箱
```bash
DOMAIN_SMTP_163_COM=smtp.163.com:587:your-email@163.com:authorization-code:true
```

### 腾讯企业邮箱
```bash
DOMAIN_SMTP_COMPANY_COM=smtp.exmail.qq.com:587:your-email@company.com:password:true
```

### 阿里云邮箱
```bash
DOMAIN_SMTP_COMPANY_COM=smtp.mxhichina.com:587:your-email@company.com:password:true
```

## 🔄 工作原理

1. **邮件发送请求**：用户选择发件人邮箱并发送邮件
2. **域名提取**：系统从发件人邮箱地址提取域名
3. **配置查找**：查找该域名对应的SMTP配置
4. **SMTP选择**：
   - 如果找到域名特定配置 → 使用该配置
   - 如果没有找到 → 使用默认SMTP配置
5. **邮件发送**：使用选定的SMTP服务器发送邮件

## 📊 配置优先级

```
域名特定配置 > 默认SMTP配置 > 系统默认值
```

**示例：**
- 发件人：`user@example.com`
- 查找：`DOMAIN_SMTP_EXAMPLE_COM` 配置
- 如果存在：使用该配置
- 如果不存在：使用 `OUTBOUND_SMTP_*` 默认配置

## 🛠️ 管理界面功能

### 查看配置
- 显示所有已配置的域名SMTP设置
- 区分默认配置和自定义配置
- 隐藏敏感信息（密码）

### 添加配置
- 域名：要配置的邮件域名
- SMTP服务器：邮件服务器地址
- 端口：SMTP端口（通常587或465）
- 用户名：SMTP认证用户名
- 密码：SMTP认证密码
- TLS：是否启用加密传输

### 删除配置
- 支持删除自定义域名配置
- 默认配置不能删除（需要通过环境变量修改）

## 🔍 故障排除

### 1. 邮件发送失败
**检查项目：**
- SMTP服务器地址和端口是否正确
- 用户名和密码是否正确
- 是否需要应用专用密码（Gmail等）
- 防火墙是否阻止SMTP端口

### 2. 配置不生效
**检查项目：**
- 环境变量名称是否正确
- 域名转换是否符合规则
- 配置格式是否正确（host:port:user:password:tls）
- 重启服务器使配置生效

### 3. 认证失败
**解决方案：**
- Gmail：使用应用专用密码
- QQ/163：使用授权码而不是登录密码
- 企业邮箱：确认SMTP功能已启用

## 📝 配置示例

### 完整配置示例
```bash
# 基础配置
WEB_PORT=8080
SMTP_PORT=25
DOMAIN=localhost

# 多域名SMTP配置
DOMAIN_SMTP_COMPANY_COM=smtp.exmail.qq.com:587:hr@company.com:hr-password:true
DOMAIN_SMTP_PERSONAL_NET=mail.personal.net:587:me@personal.net:my-password:true
DOMAIN_SMTP_BLOG_ORG=smtp.gmail.com:587:blog@blog.org:app-password:true

# 默认配置
OUTBOUND_SMTP_HOST=mail.default.com
OUTBOUND_SMTP_PORT=587
OUTBOUND_SMTP_USER=default@default.com
OUTBOUND_SMTP_PASSWORD=default-password
OUTBOUND_SMTP_TLS=true
```

### 启动命令
```bash
# 使用配置文件启动
source .env
./nbemail --port 8080 --smtp-port 2525

# 或使用启动脚本
chmod +x start-multi-domain-smtp.sh
./start-multi-domain-smtp.sh
```

## 🔐 安全建议

1. **使用应用专用密码**：避免使用主账户密码
2. **启用TLS加密**：保护邮件传输安全
3. **定期更换密码**：提高账户安全性
4. **限制SMTP权限**：只授予必要的发送权限
5. **监控发送日志**：及时发现异常活动

## 📈 性能优化

1. **连接复用**：系统自动复用SMTP连接
2. **配置缓存**：SMTP配置在内存中缓存
3. **异步发送**：支持异步邮件发送
4. **错误重试**：自动重试失败的发送请求

## 🆕 版本更新

- **v1.0**：基础SMTP功能
- **v2.0**：多域名SMTP支持
- **v2.1**：Web界面配置管理
- **v2.2**：配置热重载支持

---

**需要帮助？** 查看系统日志或联系技术支持。
