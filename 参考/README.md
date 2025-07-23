# NBEmail - 现代化邮箱系统

NBEmail 是一个使用 Go 语言开发的现代化邮箱系统，支持 SMTP 邮件服务器和 Web 管理界面，可以打包成单一的 Linux 可执行文件。

## ✨ 特性

- 🚀 **单一可执行文件** - 无需复杂部署，一个文件包含所有功能
- 📧 **完整SMTP服务器** - 支持25端口邮件接收和发送
- 🌐 **现代化Web界面** - 响应式设计，支持桌面和移动端
- 👥 **用户管理** - 支持多用户，管理员权限控制
- 🏷️ **域名管理** - 支持多域名邮件服务
- 💾 **SQLite数据库** - 轻量级，无需额外数据库服务
- 🔒 **安全认证** - 密码加密存储，Cookie会话管理
- 📱 **响应式设计** - 适配各种屏幕尺寸

## 🚀 快速开始

### 方式一：使用预编译版本

1. 下载最新版本的 `nbemail` 可执行文件
2. 给予执行权限：`chmod +x nbemail`
3. 运行：`./nbemail`
4. 访问：http://localhost:8080

### 方式二：从源码构建

```bash
# 克隆项目
git clone <repository-url>
cd nbemail

# 构建
make build

# 运行
./nbemail
```

## 📖 使用说明

### 启动参数

```bash
./nbemail --help                    # 查看帮助
./nbemail                           # 使用默认配置启动
./nbemail --port 8080               # 指定Web端口
./nbemail --smtp-port 25            # 指定SMTP端口
./nbemail --db /path/to/nbemail.db  # 指定数据库文件
```

### 默认配置

- **Web端口**: 8080
- **SMTP端口**: 25
- **数据库**: nbemail.db（当前目录）
- **默认管理员**: admin@localhost / admin123

### 环境变量配置

```bash
export WEB_PORT=8080
export SMTP_PORT=25
export DB_PATH=nbemail.db
export DOMAIN=localhost
export ADMIN_EMAIL=admin@localhost
export ADMIN_PASS=admin123
export JWT_SECRET=nbemail-secret-key-2024

# 外部SMTP发送配置（可选）
export OUTBOUND_SMTP_HOST=smtp.gmail.com
export OUTBOUND_SMTP_PORT=587
export OUTBOUND_SMTP_USER=your-email@gmail.com
export OUTBOUND_SMTP_PASSWORD=your-app-password
export OUTBOUND_SMTP_TLS=true
```

### 外部SMTP配置

为了能够向外部邮箱（如Gmail、QQ邮箱等）发送邮件，需要配置外部SMTP服务器：

1. **复制配置模板**：
   ```bash
   cp smtp-config-example.env .env
   ```

2. **编辑配置文件**，取消注释并填写相应的SMTP配置：
   - Gmail: 使用应用专用密码
   - QQ邮箱/163邮箱: 使用授权码
   - 企业邮箱: 使用相应的SMTP设置

3. **重启服务**使配置生效

**注意**：如果不配置外部SMTP，系统只能在本地用户之间发送邮件。

### 快速启动（带外部SMTP）

1. **Windows用户**：
   ```bash
   # 编辑 start-with-smtp.bat 文件，修改SMTP配置
   start-with-smtp.bat
   ```

2. **Linux/Mac用户**：
   ```bash
   # 编辑 start-with-smtp.sh 文件，修改SMTP配置
   chmod +x start-with-smtp.sh
   ./start-with-smtp.sh
   ```

### 外部SMTP发送原理

- ✅ **使用您的域名邮箱作为发件人**：系统支持使用您域名下的任意邮箱作为发件人
- ✅ **真实的发件人地址**：收件人看到的发件人是您的真实域名邮箱
- ✅ **支持回复**：收件人可以直接回复到原始发件人邮箱
- ✅ **不受第三方限制**：不依赖Gmail等第三方邮件服务商的发件人限制

## 🔧 功能说明

### SMTP邮件服务器

- 监听25端口（可配置）
- 支持标准SMTP协议
- 自动为本地用户接收邮件
- 支持邮件发送功能

### Web管理界面

#### 用户功能
- 📥 **收件箱** - 查看接收的邮件
- 📤 **发件箱** - 查看已发送的邮件
- ✏️ **写邮件** - 撰写和发送邮件
- 🔍 **邮件搜索** - 快速查找邮件

#### 管理员功能
- 👥 **用户管理** - 创建、编辑、删除用户
- 🌐 **域名管理** - 管理邮件域名
- 📊 **系统监控** - 查看系统状态

## 🏗️ 项目结构

```
nbemail/
├── main.go                 # 主程序入口
├── go.mod                  # Go模块文件
├── build.sh               # 构建脚本
├── Makefile               # 构建工具
├── README.md              # 项目文档
├── internal/              # 内部包
│   ├── config/           # 配置管理
│   ├── database/         # 数据库操作
│   ├── models/           # 数据模型
│   ├── server/           # Web服务器
│   ├── smtp/             # SMTP服务器
│   └── auth/             # 认证模块
└── web/                   # Web资源
    ├── static/           # 静态文件
    └── templates/        # 模板文件
```

## 🔨 开发

### 环境要求

- Go 1.21+
- GCC（用于SQLite）

### 开发命令

```bash
make dev          # 开发模式运行
make test         # 运行测试
make fmt          # 格式化代码
make lint         # 代码检查
make clean        # 清理构建文件
```

### 构建命令

```bash
make build        # 构建Linux版本
make build-windows # 构建Windows版本
make build-macos  # 构建macOS版本
make build-all    # 构建所有平台
make release      # 创建发布包
```

## 📧 邮件客户端配置

### 接收邮件（IMAP）
NBEmail 目前不支持 IMAP，请使用 Web 界面查看邮件。

### 发送邮件（SMTP）
- **服务器**: 你的服务器地址
- **端口**: 25
- **加密**: 无
- **认证**: 无需认证（本地用户）

## 🛡️ 安全说明

- 默认管理员密码请及时修改
- 建议在防火墙后运行
- 生产环境请配置HTTPS
- 定期备份数据库文件

## 📝 许可证

MIT License

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📞 支持

如有问题，请提交 Issue 或联系开发者。