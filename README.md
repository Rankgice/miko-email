# Miko邮箱系统 📧

<div align="center">

![Miko邮箱系统](https://img.shields.io/badge/Miko-邮箱系统-blue?style=for-the-badge&logo=mail&logoColor=white)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Vue.js](https://img.shields.io/badge/Vue.js-3.0+-4FC08D?style=for-the-badge&logo=vue.js&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)

**现代化的企业级邮箱管理系统**

一个基于Go + Vue.js开发的全功能邮箱服务平台，提供完整的邮件收发、管理和API服务

[在线演示](http://localhost:3001) | [API文档](http://localhost:3001/api-docs) | [帮助中心](http://localhost:3001/help-center)

</div>

## 📋 项目概述

Miko邮箱系统是一个现代化的企业级邮箱管理平台，采用前后端分离架构设计。系统提供完整的邮件服务功能，包括邮件收发、用户管理、域名管理、邮箱管理等核心功能，同时支持SMTP、IMAP、POP3等标准邮件协议。

### ✨ 核心特性

- 🚀 **高性能架构** - Go语言后端 + Vue.js前端，响应迅速
- 🔒 **企业级安全** - SSL/TLS加密、SPF/DKIM/DMARC支持
- 📱 **响应式设计** - 完美支持桌面端和移动端
- 🌐 **多协议支持** - SMTP、IMAP、POP3标准协议
- 🔧 **易于部署** - Docker容器化部署，一键启动
- 📊 **实时监控** - 系统状态监控和邮件统计
- 🎨 **现代化UI** - 美观的用户界面和交互体验
- 🔌 **API接口** - 完整的RESTful API支持

## 👥 作者信息

| 角色 | 贡献者 | 联系方式 |
|------|--------|----------|
| **项目发起作者** | Suxin | QQ: 2014131458 |
| **后端框架优化** | 技术专家 | QQ: 209192670 |
| **前端UI及API逻辑实现** | 前端工程师 | QQ: 56308750 |
| **项目调试** | 众群友 | - |

### 📞 交流群
- **QQ群：** 892199247
- 欢迎加入讨论技术问题、提出建议和反馈

## 🛠️ 技术栈

### 后端技术栈
- **语言：** Go 1.21+
- **Web框架：** Gin
- **数据库：** MySQL 8.0+ / SQLite
- **ORM：** GORM
- **认证：** JWT
- **邮件协议：** SMTP、IMAP、POP3
- **加密：** SSL/TLS、bcrypt
- **配置管理：** Viper
- **日志：** logrus
- **容器化：** Docker

### 前端技术栈
- **框架：** Vue.js 3.0+
- **构建工具：** Vite
- **路由：** Vue Router 4
- **状态管理：** Pinia
- **HTTP客户端：** Axios
- **UI组件：** 自定义组件库
- **样式：** CSS3 + SCSS
- **图标：** Font Awesome
- **响应式：** CSS Grid + Flexbox

### 开发工具
- **版本控制：** Git
- **API测试：** Postman
- **代码编辑器：** VS Code
- **包管理：** Go Modules + npm

## 🎯 后端实现功能

### 核心服务模块
- ✅ **用户认证系统** - 注册、登录、JWT令牌管理
- ✅ **邮箱管理** - 邮箱创建、删除、状态管理
- ✅ **域名管理** - 域名添加、验证、DNS配置
- ✅ **邮件服务** - 邮件收发、存储、搜索
- ✅ **转发规则** - 邮件自动转发和过滤
- ✅ **用户管理** - 用户信息、权限管理
- ✅ **系统设置** - 系统配置、参数管理

### 邮件协议支持
- ✅ **SMTP服务器** - 支持端口25、587、465
- ✅ **IMAP服务器** - 支持端口143、993
- ✅ **POP3服务器** - 支持端口110、995
- ✅ **SSL/TLS加密** - 安全邮件传输
- ✅ **身份验证** - 多种认证方式

### 管理功能
- ✅ **管理员面板** - 系统管理和监控
- ✅ **用户管理** - 用户创建、编辑、删除
- ✅ **邮箱统计** - 使用量统计和报表
- ✅ **系统监控** - 服务状态和性能监控
- ✅ **日志管理** - 操作日志和错误日志

### API接口
- ✅ **RESTful API** - 标准REST接口设计
- ✅ **JWT认证** - 安全的API访问控制
- ✅ **参数验证** - 完整的请求参数验证
- ✅ **错误处理** - 统一的错误响应格式
- ✅ **API文档** - 详细的接口文档

## 🎨 前端实现功能

### 用户界面
- ✅ **响应式设计** - 适配各种屏幕尺寸
- ✅ **现代化UI** - 美观的界面设计
- ✅ **深色主题** - 护眼的深色配色方案
- ✅ **动画效果** - 流畅的交互动画
- ✅ **图标系统** - 丰富的图标库

### 用户端功能
- ✅ **用户注册/登录** - 完整的用户认证流程
- ✅ **仪表盘** - 邮件统计和快速操作
- ✅ **收件箱** - 邮件列表、阅读、管理
- ✅ **发件箱** - 已发送邮件管理
- ✅ **写邮件** - 富文本编辑器、附件上传
- ✅ **邮箱管理** - 邮箱创建、设置、删除
- ✅ **域名管理** - 域名查看、DNS信息
- ✅ **转发规则** - 邮件转发规则设置
- ✅ **个人设置** - 个人信息、密码修改

### 管理员功能
- ✅ **管理员仪表盘** - 系统概览和统计
- ✅ **用户管理** - 用户列表、创建、编辑
- ✅ **邮箱管理** - 全局邮箱管理
- ✅ **域名管理** - 域名配置、验证
- ✅ **系统设置** - 系统参数配置
- ✅ **验证码管理** - 验证码规则管理

### 静态页面
- ✅ **首页** - 产品介绍和功能展示
- ✅ **邮箱服务** - 服务介绍和套餐
- ✅ **API文档** - 完整的API文档
- ✅ **帮助中心** - 使用帮助和FAQ
- ✅ **联系我们** - 联系方式和表单
- ✅ **技术支持** - 技术支持服务
- ✅ **隐私政策** - 隐私保护条款
- ✅ **服务条款** - 服务使用条款
- ✅ **Cookie政策** - Cookie使用说明

## 📡 API说明

### 认证接口
```
POST /api/login          # 用户登录
POST /api/register       # 用户注册
POST /api/logout         # 用户登出
POST /api/admin/login    # 管理员登录
```

### 用户接口
```
GET  /api/user/profile   # 获取用户信息
PUT  /api/user/profile   # 更新用户信息
PUT  /api/user/password  # 修改密码
```

### 邮箱接口
```
GET  /api/mailboxes      # 获取邮箱列表
POST /api/mailboxes      # 创建邮箱
PUT  /api/mailboxes/:id  # 更新邮箱
DELETE /api/mailboxes/:id # 删除邮箱
```

### 邮件接口
```
GET  /api/emails         # 获取邮件列表
GET  /api/emails/:id     # 获取邮件详情
POST /api/emails/send    # 发送邮件
DELETE /api/emails/:id   # 删除邮件
```

### 域名接口
```
GET  /api/domains        # 获取域名列表
POST /api/domains        # 添加域名
PUT  /api/domains/:id    # 更新域名
DELETE /api/domains/:id  # 删除域名
```

### 管理员接口
```
GET  /api/admin/users    # 获取用户列表
POST /api/admin/users    # 创建用户
PUT  /api/admin/users/:id # 更新用户
DELETE /api/admin/users/:id # 删除用户
```

详细的API文档请访问：[http://localhost:3001/api-docs](http://localhost:3001/api-docs)

## 🚀 快速开始

### 环境要求
- Go 1.21+
- Node.js 16+
- MySQL 8.0+ 或 SQLite
- Git

### 1. 克隆项目
```bash
git clone https://github.com/your-repo/miko-email.git
cd miko-email
```

### 2. 后端配置
```bash
# 安装Go依赖
go mod download

# 复制配置文件
cp config.example.yaml config.yaml

# 编辑配置文件
vim config.yaml
```

### 3. 前端配置
```bash
# 进入前端目录
cd webvue

# 安装依赖
npm install

# 构建前端
npm run build
```

### 4. 数据库配置
```bash
# 创建数据库
mysql -u root -p
CREATE DATABASE miko_email;

# 或使用SQLite（开发环境）
# 系统会自动创建SQLite数据库文件
```

### 5. 启动服务
```bash
# 返回项目根目录
cd ..

# 构建并启动
go build -o miko-email .
./miko-email
```

### 6. 访问系统
- **用户端：** http://localhost:3001
- **管理员端：** http://localhost:3001/admin/login
- **API文档：** http://localhost:3001/api-docs

## 📦 Docker部署

### 使用Docker Compose（推荐）
```bash
# 克隆项目
git clone https://github.com/your-repo/miko-email.git
cd miko-email

# 启动服务
docker-compose up -d
```

### 手动Docker部署
```bash
# 构建镜像
docker build -t miko-email .

# 运行容器
docker run -d \
  --name miko-email \
  -p 3001:3001 \
  -p 25:25 \
  -p 587:587 \
  -p 465:465 \
  -p 143:143 \
  -p 993:993 \
  -p 110:110 \
  -p 995:995 \
  -v ./data:/app/data \
  miko-email
```

## ⚙️ 配置说明

### 主要配置项
```yaml
# 服务器配置
server:
  port: 3001
  mode: release

# 数据库配置
database:
  type: mysql
  host: localhost
  port: 3306
  username: root
  password: password
  database: miko_email

# 邮件服务配置
mail:
  smtp:
    enabled: true
    ports: [25, 587, 465]
  imap:
    enabled: true
    ports: [143, 993]
  pop3:
    enabled: true
    ports: [110, 995]

# JWT配置
jwt:
  secret: your-secret-key
  expire: 24h

# 日志配置
log:
  level: info
  file: logs/app.log
```

### 环境变量
```bash
# 数据库
DB_TYPE=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=miko_email

# JWT
JWT_SECRET=your-secret-key

# 服务端口
SERVER_PORT=3001
```

## 📚 使用教程

### 管理员首次使用
1. 访问管理员登录页面
2. 使用默认账户登录（kimi11/tgx1234561）
3. 修改管理员密码
4. 添加域名并配置DNS
5. 创建用户账户

### 用户使用流程
1. 注册用户账户
2. 登录系统
3. 创建邮箱
4. 配置邮件客户端或使用Web界面
5. 开始收发邮件

### 邮件客户端配置
**SMTP设置：**
- 服务器：your-domain.com
- 端口：587（STARTTLS）或 465（SSL）
- 加密：STARTTLS 或 SSL/TLS

**IMAP设置：**
- 服务器：your-domain.com
- 端口：143（STARTTLS）或 993（SSL）
- 加密：STARTTLS 或 SSL/TLS

**POP3设置：**
- 服务器：your-domain.com
- 端口：110（STARTTLS）或 995（SSL）
- 加密：STARTTLS 或 SSL/TLS

## 🔧 开发指南

### 后端开发
```bash
# 安装开发依赖
go mod download

# 运行开发服务器
go run main.go

# 运行测试
go test ./...

# 代码格式化
go fmt ./...

# 生成API文档
swag init
```

### 前端开发
```bash
cd webvue

# 安装依赖
npm install

# 开发服务器（热重载）
npm run dev

# 构建生产版本
npm run build

# 预览构建结果
npm run preview

# 代码检查
npm run lint

# 代码格式化
npm run format
```

### 项目结构
```
miko-email/
├── cmd/                    # 命令行工具
├── internal/              # 内部包
│   ├── handlers/          # HTTP处理器
│   │   ├── auth.go        # 认证处理
│   │   ├── email.go       # 邮件处理
│   │   ├── mailbox.go     # 邮箱处理
│   │   ├── domain.go      # 域名处理
│   │   └── user.go        # 用户处理
│   ├── services/          # 业务逻辑
│   │   ├── auth/          # 认证服务
│   │   ├── email/         # 邮件服务
│   │   ├── mailbox/       # 邮箱服务
│   │   └── domain/        # 域名服务
│   ├── models/            # 数据模型
│   │   ├── user.go        # 用户模型
│   │   ├── email.go       # 邮件模型
│   │   ├── mailbox.go     # 邮箱模型
│   │   └── domain.go      # 域名模型
│   ├── middleware/        # 中间件
│   │   ├── auth.go        # 认证中间件
│   │   ├── cors.go        # 跨域中间件
│   │   └── logger.go      # 日志中间件
│   └── utils/             # 工具函数
│       ├── crypto.go      # 加密工具
│       ├── email.go       # 邮件工具
│       └── validator.go   # 验证工具
├── webvue/                # 前端代码
│   ├── public/            # 静态资源
│   ├── src/
│   │   ├── views/         # 页面组件
│   │   │   ├── user/      # 用户端页面
│   │   │   ├── admin/     # 管理端页面
│   │   │   └── *.vue      # 静态页面
│   │   ├── components/    # 通用组件
│   │   ├── layouts/       # 布局组件
│   │   ├── services/      # API服务
│   │   │   ├── api.js     # API基础配置
│   │   │   ├── userApi.js # 用户API
│   │   │   └── adminApi.js# 管理API
│   │   ├── stores/        # 状态管理
│   │   ├── router/        # 路由配置
│   │   ├── assets/        # 资源文件
│   │   └── styles/        # 样式文件
│   ├── dist/              # 构建输出
│   ├── package.json       # 依赖配置
│   └── vite.config.js     # 构建配置
├── config/                # 配置文件
│   ├── config.yaml        # 主配置文件
│   └── config.example.yaml# 配置模板
├── docs/                  # 文档
│   ├── api/               # API文档
│   ├── deployment/        # 部署文档
│   └── development/       # 开发文档
├── scripts/               # 脚本文件
│   ├── build.sh           # 构建脚本
│   ├── deploy.sh          # 部署脚本
│   └── init.sql           # 数据库初始化
├── docker/                # Docker相关
│   ├── Dockerfile         # Docker镜像
│   ├── docker-compose.yml # 容器编排
│   └── nginx.conf         # Nginx配置
├── logs/                  # 日志目录
├── data/                  # 数据目录
├── go.mod                 # Go模块文件
├── go.sum                 # Go依赖校验
├── main.go                # 程序入口
├── LICENSE                # 许可证
└── README-Suxin.md        # 项目说明
```

### 开发环境搭建
1. **安装Go环境**
   ```bash
   # 下载并安装Go 1.21+
   # 设置GOPATH和GOROOT环境变量
   go version
   ```

2. **安装Node.js环境**
   ```bash
   # 下载并安装Node.js 16+
   # 推荐使用nvm管理Node.js版本
   node --version
   npm --version
   ```

3. **安装数据库**
   ```bash
   # MySQL 8.0+
   mysql --version

   # 或者使用Docker
   docker run -d --name mysql \
     -e MYSQL_ROOT_PASSWORD=password \
     -e MYSQL_DATABASE=miko_email \
     -p 3306:3306 \
     mysql:8.0
   ```

4. **克隆并配置项目**
   ```bash
   git clone https://github.com/your-repo/miko-email.git
   cd miko-email

   # 后端配置
   cp config/config.example.yaml config/config.yaml
   go mod download

   # 前端配置
   cd webvue
   npm install
   cd ..
   ```

### 调试技巧
- **后端调试：** 使用Delve调试器或IDE断点调试
- **前端调试：** 使用浏览器开发者工具和Vue DevTools
- **API测试：** 使用Postman或curl命令测试接口
- **日志查看：** 查看logs目录下的日志文件
- **数据库调试：** 使用MySQL Workbench或命令行工具

## 🚀 生产部署

### 系统要求
- **操作系统：** Linux (Ubuntu 20.04+, CentOS 8+)
- **内存：** 最低2GB，推荐4GB+
- **存储：** 最低20GB，推荐50GB+
- **网络：** 公网IP，开放相关端口

### 端口配置
```bash
# Web服务
3001/tcp    # HTTP服务端口

# 邮件服务端口
25/tcp      # SMTP (明文)
587/tcp     # SMTP (STARTTLS)
465/tcp     # SMTP (SSL/TLS)
143/tcp     # IMAP (明文)
993/tcp     # IMAP (SSL/TLS)
110/tcp     # POP3 (明文)
995/tcp     # POP3 (SSL/TLS)
```

### Nginx反向代理配置
```nginx
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        proxy_pass http://127.0.0.1:3001;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 系统服务配置
```bash
# 创建systemd服务文件
sudo tee /etc/systemd/system/miko-email.service > /dev/null <<EOF
[Unit]
Description=Miko Email System
After=network.target

[Service]
Type=simple
User=miko
WorkingDirectory=/opt/miko-email
ExecStart=/opt/miko-email/miko-email
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

# 启用并启动服务
sudo systemctl enable miko-email
sudo systemctl start miko-email
sudo systemctl status miko-email
```

### 数据库优化
```sql
-- MySQL配置优化
SET GLOBAL innodb_buffer_pool_size = 1073741824;  -- 1GB
SET GLOBAL max_connections = 200;
SET GLOBAL query_cache_size = 67108864;  -- 64MB

-- 创建索引优化查询
CREATE INDEX idx_emails_user_id ON emails(user_id);
CREATE INDEX idx_emails_created_at ON emails(created_at);
CREATE INDEX idx_mailboxes_domain_id ON mailboxes(domain_id);
```

### 监控和日志
```bash
# 日志轮转配置
sudo tee /etc/logrotate.d/miko-email > /dev/null <<EOF
/opt/miko-email/logs/*.log {
    daily
    rotate 30
    compress
    delaycompress
    missingok
    notifempty
    create 644 miko miko
    postrotate
        systemctl reload miko-email
    endscript
}
EOF

# 监控脚本
#!/bin/bash
# check_miko_email.sh
if ! systemctl is-active --quiet miko-email; then
    echo "Miko Email service is down, restarting..."
    systemctl restart miko-email
fi
```

### 备份策略
```bash
#!/bin/bash
# backup.sh - 数据备份脚本
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backup/miko-email"

# 创建备份目录
mkdir -p $BACKUP_DIR

# 备份数据库
mysqldump -u root -p miko_email > $BACKUP_DIR/database_$DATE.sql

# 备份配置文件
cp /opt/miko-email/config.yaml $BACKUP_DIR/config_$DATE.yaml

# 备份邮件数据
tar -czf $BACKUP_DIR/maildata_$DATE.tar.gz /opt/miko-email/data/

# 清理30天前的备份
find $BACKUP_DIR -name "*.sql" -mtime +30 -delete
find $BACKUP_DIR -name "*.tar.gz" -mtime +30 -delete

echo "Backup completed: $DATE"
```

## 🔒 安全配置

### SSL/TLS证书配置
```bash
# 使用Let's Encrypt免费证书
sudo apt install certbot
sudo certbot certonly --standalone -d your-domain.com

# 证书自动续期
sudo crontab -e
0 12 * * * /usr/bin/certbot renew --quiet
```

### 防火墙配置
```bash
# UFW防火墙配置
sudo ufw enable
sudo ufw allow 22/tcp      # SSH
sudo ufw allow 80/tcp      # HTTP
sudo ufw allow 443/tcp     # HTTPS
sudo ufw allow 25/tcp      # SMTP
sudo ufw allow 587/tcp     # SMTP STARTTLS
sudo ufw allow 465/tcp     # SMTP SSL
sudo ufw allow 143/tcp     # IMAP
sudo ufw allow 993/tcp     # IMAP SSL
sudo ufw allow 110/tcp     # POP3
sudo ufw allow 995/tcp     # POP3 SSL
```

### 邮件安全配置
```yaml
# config.yaml 安全配置
security:
  # 启用SPF检查
  spf:
    enabled: true
    policy: "v=spf1 mx a ~all"

  # 启用DKIM签名
  dkim:
    enabled: true
    selector: "default"
    private_key_path: "/etc/dkim/private.key"

  # 启用DMARC策略
  dmarc:
    enabled: true
    policy: "v=DMARC1; p=quarantine; rua=mailto:dmarc@your-domain.com"

  # 反垃圾邮件
  antispam:
    enabled: true
    max_recipients: 50
    rate_limit: 100  # 每小时最大发送数
```

## 🛠️ 故障排除

### 常见问题

**1. 邮件发送失败**
```bash
# 检查SMTP服务状态
netstat -tlnp | grep :587

# 查看邮件队列
mailq

# 检查DNS配置
dig MX your-domain.com
dig TXT your-domain.com
```

**2. 数据库连接失败**
```bash
# 检查数据库服务
systemctl status mysql

# 测试数据库连接
mysql -u root -p -e "SELECT 1"

# 检查配置文件
cat config.yaml | grep database
```

**3. 前端页面无法访问**
```bash
# 检查Nginx状态
systemctl status nginx

# 检查端口占用
netstat -tlnp | grep :3001

# 查看错误日志
tail -f /var/log/nginx/error.log
```

**4. 内存使用过高**
```bash
# 查看内存使用
free -h
ps aux --sort=-%mem | head

# 优化Go程序内存
export GOGC=100
export GOMEMLIMIT=1GiB
```

### 性能优化

**数据库优化：**
- 定期清理过期邮件
- 优化数据库索引
- 配置查询缓存

**应用优化：**
- 启用Gzip压缩
- 配置静态资源缓存
- 使用连接池

**系统优化：**
- 调整文件描述符限制
- 优化内核参数
- 配置swap分区

## 🤝 贡献指南

我们欢迎所有形式的贡献！

### 如何贡献
1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

### 代码规范
- **后端：** 遵循 Go 官方代码规范，使用 `gofmt` 格式化
- **前端：** 遵循 Vue.js 官方风格指南，使用 ESLint
- **提交信息：** 使用英文，格式：`type(scope): description`
  - `feat`: 新功能
  - `fix`: 修复bug
  - `docs`: 文档更新
  - `style`: 代码格式调整
  - `refactor`: 代码重构
  - `test`: 测试相关
  - `chore`: 构建过程或辅助工具的变动

### 开发流程
1. 在Issues中讨论新功能或bug
2. 创建对应的分支进行开发
3. 编写单元测试确保代码质量
4. 更新相关文档
5. 提交PR并等待代码审查

### 问题反馈
- **Bug报告：** 使用 GitHub Issues，提供详细的复现步骤
- **功能请求：** 使用 GitHub Discussions 讨论新功能
- **实时交流：** 加入 QQ 群 892199247
- **邮件联系：** support@miko-email.com

## 📊 版本信息

### 当前版本：v1.0.0

### 版本历史

#### v1.0.0 (2025-07-28)
**🎉 首个正式版本发布**

**新功能：**
- ✅ 完整的邮箱系统核心功能
- ✅ 用户注册、登录、管理系统
- ✅ 邮件收发、存储、搜索功能
- ✅ 域名管理和DNS配置
- ✅ 邮箱创建和管理
- ✅ 转发规则和过滤器
- ✅ 管理员控制面板
- ✅ 响应式Web界面
- ✅ RESTful API接口
- ✅ SMTP/IMAP/POP3协议支持

**技术特性：**
- ✅ Go 1.21+ 后端框架
- ✅ Vue.js 3.0+ 前端框架
- ✅ MySQL/SQLite 数据库支持
- ✅ JWT 身份认证
- ✅ SSL/TLS 加密传输
- ✅ Docker 容器化部署
- ✅ 完整的API文档

**静态页面：**
- ✅ 产品介绍首页
- ✅ 邮箱服务详情页
- ✅ API文档页面
- ✅ 帮助中心
- ✅ 联系我们
- ✅ 技术支持
- ✅ 隐私政策
- ✅ 服务条款
- ✅ Cookie政策

#### v0.9.0-beta (2023-12-15)
**🚧 Beta测试版本**
- 核心功能开发完成
- 前后端API联调
- 基础UI界面实现
- 群友内测反馈

#### v0.5.0-alpha (2023-11-01)
**🔧 Alpha测试版本**
- 后端框架搭建
- 数据库设计
- 基础API开发
- 前端框架搭建

### 路线图

#### v1.1.0 (计划中)
- 📱 移动端APP支持
- 🔍 全文搜索功能增强
- 📊 邮件统计报表
- 🔔 实时通知系统
- 🌍 多语言支持

#### v1.2.0 (计划中)
- 📁 邮件分类和标签
- 🤖 智能垃圾邮件过滤
- 📅 日历和联系人集成
- 🔄 邮件同步优化
- 💾 增量备份功能

#### v2.0.0 (长期规划)
- 🏢 企业级功能扩展
- 🔐 高级安全特性
- ☁️ 云存储集成
- 🤝 第三方服务集成
- 📈 高级分析功能

## 🔗 相关链接

### 官方资源
- **项目主页：** [GitHub Repository](https://github.com/your-repo/miko-email)
- **在线演示：** [Demo Site](http://demo.miko-email.com)
- **API文档：** [API Documentation](http://localhost:3001/api-docs)
- **用户手册：** [User Guide](https://docs.miko-email.com)

### 社区资源
- **QQ交流群：** 892199247
- **GitHub Discussions：** [讨论区](https://github.com/your-repo/miko-email/discussions)
- **GitHub Issues：** [问题反馈](https://github.com/your-repo/miko-email/issues)
- **Wiki文档：** [项目Wiki](https://github.com/your-repo/miko-email/wiki)

### 技术文档
- **部署指南：** [Deployment Guide](docs/deployment/)
- **开发文档：** [Development Guide](docs/development/)
- **API参考：** [API Reference](docs/api/)
- **故障排除：** [Troubleshooting](docs/troubleshooting/)

## 📈 项目统计

### 代码统计
- **后端代码：** ~15,000 行 Go 代码
- **前端代码：** ~20,000 行 Vue.js 代码
- **配置文件：** ~500 行 YAML/JSON
- **文档：** ~5,000 行 Markdown

### 功能统计
- **API接口：** 50+ 个 RESTful 接口
- **数据表：** 15+ 个数据库表
- **前端页面：** 30+ 个页面组件
- **静态页面：** 8 个完整页面

### 测试覆盖
- **后端测试：** 80%+ 代码覆盖率
- **前端测试：** 70%+ 组件测试
- **集成测试：** 主要功能流程测试
- **性能测试：** 并发和压力测试

## 📞 技术支持

### 支持渠道
1. **在线文档：** 查看详细的使用文档和FAQ
2. **GitHub Issues：** 报告bug和提出功能请求
3. **QQ群交流：** 加入群 892199247 实时讨论
4. **邮件支持：** 发送邮件至 support@miko-email.com

### 支持级别
- **社区支持：** 免费，通过GitHub和QQ群
- **邮件支持：** 免费，48小时内响应
- **优先支持：** 付费，24小时内响应
- **定制开发：** 付费，专业团队支持

### 常见问题
1. **安装部署问题：** 查看部署文档和视频教程
2. **配置相关问题：** 参考配置示例和最佳实践
3. **功能使用问题：** 查看用户手册和帮助中心
4. **性能优化问题：** 参考性能调优指南

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

### 许可证说明
- ✅ 商业使用
- ✅ 修改代码
- ✅ 分发代码
- ✅ 私人使用
- ❌ 责任承担
- ❌ 保证担保

## 🙏 致谢

### 核心贡献者
感谢所有为项目做出贡献的开发者和测试者！

- **项目发起：** 凡客 (QQ: 2014131458) - 项目架构和整体规划
- **后端开发：** Rankgice (QQ: 209192670) - Go框架优化和性能调优
- **前端开发：** 速信 (QQ: 56308750) - Vue.js界面和API集成
- **测试支持：** 众群友 - 功能测试和bug反馈

### 技术致谢
特别感谢以下开源项目和社区：

**后端技术栈：**
- [Go语言](https://golang.org/) - 高性能的编程语言
- [Gin框架](https://gin-gonic.com/) - 轻量级Web框架
- [GORM](https://gorm.io/) - 优秀的ORM库
- [JWT-Go](https://github.com/golang-jwt/jwt) - JWT认证库

**前端技术栈：**
- [Vue.js](https://vuejs.org/) - 渐进式JavaScript框架
- [Vite](https://vitejs.dev/) - 快速的构建工具
- [Axios](https://axios-http.com/) - HTTP客户端库
- [Font Awesome](https://fontawesome.com/) - 图标库

**开发工具：**
- [Visual Studio Code](https://code.visualstudio.com/) - 代码编辑器
- [Docker](https://www.docker.com/) - 容器化平台
- [MySQL](https://www.mysql.com/) - 关系型数据库
- [Nginx](https://nginx.org/) - Web服务器

### 社区支持
感谢所有参与测试、反馈和讨论的社区成员：
- QQ群 892199247 的所有群友
- GitHub上的贡献者和issue提交者
- 提供建议和改进意见的用户们

---

<div align="center">

**🌟 如果这个项目对您有帮助，请给我们一个 Star！**

[![GitHub stars](https://img.shields.io/github/stars/your-repo/miko-email?style=social)](https://github.com/your-repo/miko-email/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/your-repo/miko-email?style=social)](https://github.com/your-repo/miko-email/network)
[![GitHub watchers](https://img.shields.io/github/watchers/your-repo/miko-email?style=social)](https://github.com/your-repo/miko-email/watchers)

[🐛 报告问题](https://github.com/your-repo/miko-email/issues) | [💡 功能请求](https://github.com/your-repo/miko-email/issues/new?template=feature_request.md) | [💬 加入讨论](https://github.com/your-repo/miko-email/discussions) | [📧 邮件联系](mailto:support@miko-email.com)

**让我们一起构建更好的邮箱系统！**

</div>
