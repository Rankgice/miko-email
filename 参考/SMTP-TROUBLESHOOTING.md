# SMTP发送故障排除指南

## 常见错误及解决方案

### 1. SSL连接失败

**错误信息：**
- "SSL连接失败"
- "bind: An operation on a socket could not be performed because the system lacked sufficient buffer space or because a queue was full"

**可能原因：**
- 网络连接不稳定
- 防火墙阻止连接
- SMTP服务器配置错误
- 系统资源不足

**解决方案：**
1. **检查网络连接**
   ```bash
   # Windows
   ping smtp.gmail.com
   telnet smtp.gmail.com 465
   
   # Linux/macOS
   ping smtp.gmail.com
   telnet smtp.gmail.com 465
   ```

2. **使用SMTP测试工具**
   ```bash
   # Windows
   test-smtp.bat smtp.gmail.com:465 your-email@gmail.com your-password true
   
   # Linux/macOS
   ./test-smtp.sh smtp.gmail.com:465 your-email@gmail.com your-password true
   ```

3. **检查防火墙设置**
   - 确保防火墙允许出站连接到SMTP端口（25, 587, 465）
   - 临时关闭防火墙测试连接

4. **尝试不同的端口**
   - 465端口（SSL）
   - 587端口（STARTTLS）
   - 25端口（普通连接）

### 2. SMTP认证失败

**错误信息：**
- "SMTP认证失败"
- "Authentication failed"

**解决方案：**
1. **检查用户名和密码**
   - 确保用户名是完整的邮箱地址
   - 检查密码是否正确

2. **启用"允许不够安全的应用"**
   - Gmail: 在Google账户设置中启用
   - 或使用应用专用密码

3. **检查邮箱服务商要求**
   - 某些邮箱需要特殊设置
   - 查看邮箱服务商的SMTP配置文档

### 3. TLS/STARTTLS失败

**错误信息：**
- "启动TLS失败"
- "TLS handshake failed"

**解决方案：**
1. **检查TLS版本支持**
   - 确保服务器支持TLS 1.2或更高版本
   - 更新系统的TLS库

2. **尝试不同的加密方式**
   - SSL（端口465）
   - STARTTLS（端口587）
   - 无加密（端口25，不推荐）

### 4. 连接超时

**错误信息：**
- "连接超时"
- "Connection timeout"

**解决方案：**
1. **检查网络延迟**
   ```bash
   ping smtp.gmail.com
   ```

2. **增加超时时间**
   - 系统已设置30秒超时
   - 检查网络是否过慢

3. **检查DNS解析**
   ```bash
   nslookup smtp.gmail.com
   ```

## 常见SMTP服务器配置

### Gmail
- **SMTP服务器**: smtp.gmail.com
- **端口**: 587 (STARTTLS) 或 465 (SSL)
- **认证**: 需要
- **注意**: 需要启用"允许不够安全的应用"或使用应用专用密码

### Outlook/Hotmail
- **SMTP服务器**: smtp-mail.outlook.com
- **端口**: 587 (STARTTLS)
- **认证**: 需要

### QQ邮箱
- **SMTP服务器**: smtp.qq.com
- **端口**: 587 (STARTTLS) 或 465 (SSL)
- **认证**: 需要
- **注意**: 需要在QQ邮箱设置中开启SMTP服务

### 163邮箱
- **SMTP服务器**: smtp.163.com
- **端口**: 587 (STARTTLS) 或 465 (SSL)
- **认证**: 需要
- **注意**: 需要使用授权码而不是登录密码

## 测试工具使用

### 基本连接测试
```bash
# Windows
test-smtp.bat smtp.gmail.com:587

# Linux/macOS
./test-smtp.sh smtp.gmail.com:587
```

### 完整认证测试
```bash
# Windows
test-smtp.bat smtp.gmail.com:587 your-email@gmail.com your-password true

# Linux/macOS
./test-smtp.sh smtp.gmail.com:587 your-email@gmail.com your-password true
```

## 系统级解决方案

### Windows
1. **检查Windows Defender防火墙**
2. **更新网络驱动程序**
3. **重置网络设置**
   ```cmd
   netsh winsock reset
   netsh int ip reset
   ```

### Linux
1. **检查iptables规则**
   ```bash
   sudo iptables -L
   ```
2. **检查系统日志**
   ```bash
   journalctl -f
   ```

### 网络诊断命令

```bash
# 检查端口连通性
telnet smtp.gmail.com 587

# 检查DNS解析
nslookup smtp.gmail.com

# 检查路由
traceroute smtp.gmail.com  # Linux/macOS
tracert smtp.gmail.com     # Windows
```

## 联系支持

如果以上方法都无法解决问题，请：

1. 运行SMTP测试工具并保存输出结果
2. 检查系统日志中的错误信息
3. 记录具体的错误信息和重现步骤
4. 联系系统管理员或技术支持
