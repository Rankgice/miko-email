# Miko邮箱系统配置管理

## 📋 概述

Miko邮箱系统现在支持通过 `config.yaml` 文件进行配置管理，包括服务器端口、管理员账号、域名设置等。

## 🔧 配置文件

### config.yaml

主配置文件，包含所有系统设置：

```yaml
# 服务器端口配置
server:
  web_port: 8080
  smtp:
    enable_multi_port: true
    port_25: 25
    port_587: 587
    port_465: 465
  imap:
    port: 143
  pop3:
    port: 110

# 管理员账号配置
admin:
  username: "admin"
  password: "your_password"
  email: "admin@yourdomain.com"
  enabled: true

# 域名配置
domain:
  default: "localhost"
  allowed: []  # 空数组表示不限制域名
  enable_domain_restriction: false
```

## 🛠️ 管理工具

### 1. 配置管理工具

```bash
# 显示当前配置
go run tools/config_manager.go show

# 显示管理员信息
go run tools/config_manager.go admin

# 显示端口配置
go run tools/config_manager.go ports

# 显示功能开关
go run tools/config_manager.go features

# 测试配置文件
go run tools/config_manager.go test
```

### 2. 管理员同步工具

```bash
# 从config.yaml同步管理员信息到数据库
go run tools/sync_admin.go sync

# 显示数据库中的管理员信息
go run tools/sync_admin.go show

# 重置管理员密码为配置文件中的密码
go run tools/sync_admin.go reset
```

## 🚀 启动方式

### 方式1：使用启动脚本（推荐）

**Windows:**
```cmd
start_with_sync.bat
```

**Linux/Mac:**
```bash
./start_with_sync.sh
```

启动脚本会自动：
1. 同步管理员信息到数据库
2. 显示当前配置
3. 启动邮件服务器

### 方式2：手动启动

```bash
# 1. 同步管理员信息（可选）
go run tools/sync_admin.go sync

# 2. 启动服务器
go run main.go
```

## 📝 配置修改流程

1. **修改配置文件**
   ```bash
   # 编辑 config.yaml 文件
   nano config.yaml
   ```

2. **同步管理员信息**（如果修改了管理员配置）
   ```bash
   go run tools/sync_admin.go sync
   ```

3. **重启服务器**
   ```bash
   # 停止当前服务器 (Ctrl+C)
   # 重新启动
   go run main.go
   ```

## 🔑 管理员账号管理

### 修改管理员密码

1. 编辑 `config.yaml` 中的 `admin.password`
2. 运行同步命令：
   ```bash
   go run tools/sync_admin.go sync
   ```

### 修改管理员邮箱

1. 编辑 `config.yaml` 中的 `admin.email`
2. 运行同步命令：
   ```bash
   go run tools/sync_admin.go sync
   ```

### 禁用管理员账号

1. 设置 `config.yaml` 中的 `admin.enabled: false`
2. 运行同步命令：
   ```bash
   go run tools/sync_admin.go sync
   ```

## 🌐 域名配置

### 不限制域名（默认）

```yaml
domain:
  enable_domain_restriction: false
  allowed: []
```

### 限制特定域名

```yaml
domain:
  enable_domain_restriction: true
  allowed:
    - "yourdomain.com"
    - "example.com"
```

## 📡 端口配置

### 多SMTP端口（推荐）

```yaml
server:
  smtp:
    enable_multi_port: true
    port_25: 25    # 标准SMTP端口
    port_587: 587  # SMTP提交端口
    port_465: 465  # SMTPS安全端口
```

### 单SMTP端口

```yaml
server:
  smtp:
    enable_multi_port: false
    port_25: 25
```

## 🔍 故障排除

### 配置文件格式错误

```bash
# 测试配置文件格式
go run tools/config_manager.go test
```

### 管理员登录失败

```bash
# 检查管理员信息
go run tools/sync_admin.go show

# 重置管理员密码
go run tools/sync_admin.go reset
```

### 端口被占用

```bash
# 检查端口配置
go run tools/config_manager.go ports

# 修改config.yaml中的端口设置
```

## 📚 配置项说明

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| `server.web_port` | Web管理界面端口 | 8080 |
| `server.smtp.port_25` | 标准SMTP端口 | 25 |
| `server.smtp.port_587` | SMTP提交端口 | 587 |
| `server.smtp.port_465` | SMTPS安全端口 | 465 |
| `admin.username` | 管理员用户名 | admin |
| `admin.password` | 管理员密码 | admin123456 |
| `admin.email` | 管理员邮箱 | admin@localhost |
| `domain.enable_domain_restriction` | 是否限制域名 | false |
| `features.allow_registration` | 是否允许用户注册 | true |

## 🎯 最佳实践

1. **定期备份配置文件**
2. **使用强密码作为管理员密码**
3. **根据需要调整端口配置**
4. **定期检查管理员账号状态**
5. **修改配置后及时同步到数据库**
