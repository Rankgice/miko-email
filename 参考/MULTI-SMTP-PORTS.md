# NBEmail 多SMTP端口支持

NBEmail现在支持同时启动多个SMTP端口，包括标准的25、587和465端口。

## 支持的SMTP端口

- **端口25**: 标准SMTP端口，用于服务器间邮件传输
- **端口587**: SMTP提交端口，支持STARTTLS，推荐用于邮件客户端
- **端口465**: SMTPS端口，支持SSL/TLS加密连接

## 启用方式

### 1. 命令行参数方式

```bash
# 启用多SMTP端口
./nbemail --multi-smtp

# 或者使用go run
go run main.go --multi-smtp
```

### 2. 环境变量方式

```bash
# 设置环境变量
export ENABLE_MULTI_SMTP=true

# 然后启动程序
./nbemail
```

### 3. 使用提供的脚本

**Windows:**
```cmd
start-multi-smtp.bat
```

**Linux/macOS:**
```bash
./start-multi-smtp.sh
```

## 配置选项

### 环境变量配置

可以通过以下环境变量自定义端口：

```bash
export SMTP_PORT=25          # 主SMTP端口
export SMTP_PORT_587=587     # SMTP提交端口
export SMTP_PORT_465=465     # SMTPS端口
export ENABLE_MULTI_SMTP=true # 启用多端口模式
```

### 配置文件说明

在 `internal/config/config.go` 中新增了以下配置项：

```go
type Config struct {
    SMTPPort        int  `json:"smtp_port"`        // 主SMTP端口 (默认25)
    SMTPPort587     int  `json:"smtp_port_587"`    // SMTP提交端口 (587)
    SMTPPort465     int  `json:"smtp_port_465"`    // SMTPS端口 (465)
    EnableMultiSMTP bool `json:"enable_multi_smtp"` // 是否启用多SMTP端口
    // ... 其他配置项
}
```

## 使用建议

1. **对于现代邮件客户端**: 推荐使用587端口，支持STARTTLS加密
2. **对于需要SSL/TLS的客户端**: 使用465端口
3. **对于服务器间通信**: 使用25端口

## 验证端口状态

启动后可以使用以下命令验证端口是否正常监听：

**Windows:**
```cmd
netstat -an | findstr ":25\|:587\|:465"
```

**Linux/macOS:**
```bash
netstat -an | grep -E ":(25|587|465)"
```

## 注意事项

1. 在某些系统上，使用标准端口（特别是25端口）可能需要管理员权限
2. 确保防火墙允许这些端口的连接
3. 如果端口被其他服务占用，程序会显示启动失败的错误信息

## 故障排除

如果某个端口启动失败，请检查：

1. 端口是否被其他程序占用
2. 是否有足够的权限绑定端口
3. 防火墙设置是否正确

可以通过日志输出查看具体的错误信息。
