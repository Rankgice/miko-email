#!/bin/bash

# Miko邮箱系统部署脚本
# 使用方法: ./deploy.sh [production|development]

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查Go环境
check_go() {
    log_info "检查Go环境..."
    if ! command -v go &> /dev/null; then
        log_error "Go未安装，请先安装Go 1.21或更高版本"
        exit 1
    fi
    
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    log_success "Go版本: $GO_VERSION"
}

# 检查依赖
check_dependencies() {
    log_info "检查项目依赖..."
    go mod tidy
    log_success "依赖检查完成"
}

# 初始化数据库
init_database() {
    log_info "初始化数据库..."
    if [ ! -f "miko_email.db" ]; then
        go run cmd/init/main.go
        log_success "数据库初始化完成"
    else
        log_warning "数据库文件已存在，跳过初始化"
    fi
}

# 构建应用
build_app() {
    log_info "构建应用程序..."
    
    # 设置构建信息
    BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
    GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
    
    # 构建参数
    LDFLAGS="-X main.BuildTime=$BUILD_TIME -X main.GitCommit=$GIT_COMMIT"
    
    if [ "$1" = "production" ]; then
        log_info "生产环境构建..."
        go build -ldflags "$LDFLAGS -s -w" -o miko-email main.go
    else
        log_info "开发环境构建..."
        go build -ldflags "$LDFLAGS" -o miko-email main.go
    fi
    
    log_success "应用构建完成"
}

# 创建配置文件
create_config() {
    log_info "创建配置文件..."
    
    if [ ! -f ".env" ]; then
        cat > .env << EOF
# Miko邮箱系统配置文件

# Web服务配置
WEB_PORT=8080

# 邮件服务端口配置
SMTP_PORT=25
IMAP_PORT=143
POP3_PORT=110

# 数据库配置
DATABASE_PATH=./miko_email.db

# 会话密钥（生产环境请修改）
SESSION_KEY=miko-email-secret-key-change-in-production

# 域名配置
DOMAIN=localhost

# 运行模式 (development/production)
GIN_MODE=release
EOF
        log_success "配置文件创建完成: .env"
    else
        log_warning "配置文件已存在，跳过创建"
    fi
}

# 创建systemd服务文件
create_systemd_service() {
    log_info "创建systemd服务文件..."
    
    CURRENT_DIR=$(pwd)
    USER=$(whoami)
    
    cat > miko-email.service << EOF
[Unit]
Description=Miko Email System
After=network.target

[Service]
Type=simple
User=$USER
WorkingDirectory=$CURRENT_DIR
ExecStart=$CURRENT_DIR/miko-email
Restart=always
RestartSec=5
Environment=PATH=/usr/local/bin:/usr/bin:/bin
EnvironmentFile=$CURRENT_DIR/.env

# 安全设置
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=$CURRENT_DIR

[Install]
WantedBy=multi-user.target
EOF
    
    log_success "systemd服务文件创建完成: miko-email.service"
    log_info "要安装服务，请运行: sudo cp miko-email.service /etc/systemd/system/"
    log_info "然后运行: sudo systemctl enable miko-email && sudo systemctl start miko-email"
}

# 创建Docker配置
create_docker_config() {
    log_info "创建Docker配置文件..."
    
    # Dockerfile
    cat > Dockerfile << EOF
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o miko-email main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/miko-email .
COPY --from=builder /app/web ./web
COPY --from=builder /app/scripts ./scripts

EXPOSE 8080 25 143 110

CMD ["./miko-email"]
EOF

    # docker-compose.yml
    cat > docker-compose.yml << EOF
version: '3.8'

services:
  miko-email:
    build: .
    ports:
      - "8080:8080"
      - "25:25"
      - "143:143"
      - "110:110"
    volumes:
      - ./data:/root/data
    environment:
      - WEB_PORT=8080
      - SMTP_PORT=25
      - IMAP_PORT=143
      - POP3_PORT=110
      - DATABASE_PATH=/root/data/miko_email.db
      - DOMAIN=localhost
      - GIN_MODE=release
    restart: unless-stopped
    
volumes:
  data:
EOF

    log_success "Docker配置文件创建完成"
}

# 运行测试
run_tests() {
    log_info "运行测试..."
    
    if command -v python3 &> /dev/null; then
        python3 test_api.py
        log_success "API测试完成"
    else
        log_warning "Python3未安装，跳过API测试"
    fi
}

# 主函数
main() {
    local mode=${1:-development}
    
    log_info "开始部署Miko邮箱系统 (模式: $mode)"
    
    # 检查环境
    check_go
    check_dependencies
    
    # 初始化
    init_database
    create_config
    
    # 构建
    build_app $mode
    
    # 创建部署文件
    if [ "$mode" = "production" ]; then
        create_systemd_service
        create_docker_config
    fi
    
    # 测试
    if [ "$mode" = "development" ]; then
        log_info "启动开发服务器进行测试..."
        ./miko-email &
        SERVER_PID=$!
        
        sleep 3
        run_tests
        
        kill $SERVER_PID 2>/dev/null || true
    fi
    
    log_success "部署完成！"
    
    if [ "$mode" = "production" ]; then
        log_info "生产环境部署说明:"
        log_info "1. 修改 .env 文件中的配置"
        log_info "2. 安装systemd服务: sudo cp miko-email.service /etc/systemd/system/"
        log_info "3. 启动服务: sudo systemctl enable miko-email && sudo systemctl start miko-email"
        log_info "4. 或使用Docker: docker-compose up -d"
    else
        log_info "开发环境启动命令: ./miko-email"
    fi
}

# 脚本入口
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
