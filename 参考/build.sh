#!/bin/bash

# NBEmail 构建脚本
# 用于构建单一的Linux可执行文件

set -e

echo "🚀 开始构建 NBEmail..."

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "❌ 错误: 未找到Go环境，请先安装Go"
    exit 1
fi

# 创建必要的目录
mkdir -p web/static/css
mkdir -p web/static/js
mkdir -p web/static/images
mkdir -p web/templates

# 创建基础的静态文件（如果不存在）
if [ ! -f "web/static/css/style.css" ]; then
    echo "/* NBEmail 样式文件 */" > web/static/css/style.css
fi

if [ ! -f "web/static/js/app.js" ]; then
    echo "// NBEmail 应用脚本" > web/static/js/app.js
fi

if [ ! -f "web/templates/base.html" ]; then
    echo "<!-- NBEmail 基础模板 -->" > web/templates/base.html
fi

# 设置构建环境变量
export CGO_ENABLED=1
export GOOS=linux
export GOARCH=amd64

# 构建信息
BUILD_TIME=$(date -u '+%Y-%m-%d %H:%M:%S UTC')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
VERSION="1.0.0"

# 构建标志
LDFLAGS="-s -w -X 'main.Version=${VERSION}' -X 'main.BuildTime=${BUILD_TIME}' -X 'main.GitCommit=${GIT_COMMIT}'"

echo "📦 正在构建..."
echo "   版本: ${VERSION}"
echo "   构建时间: ${BUILD_TIME}"
echo "   Git提交: ${GIT_COMMIT}"
echo "   目标平台: ${GOOS}/${GOARCH}"

# 下载依赖
echo "📥 下载依赖..."
go mod tidy
go mod download

# 构建可执行文件
echo "🔨 编译中..."
go build -ldflags="${LDFLAGS}" -o nbemail main.go

# 检查构建结果
if [ -f "nbemail" ]; then
    echo "✅ 构建成功!"
    echo "📁 可执行文件: $(pwd)/nbemail"
    echo "📊 文件大小: $(du -h nbemail | cut -f1)"
    echo ""
    echo "🚀 使用方法:"
    echo "   ./nbemail --help                    # 查看帮助"
    echo "   ./nbemail                           # 使用默认配置启动"
    echo "   ./nbemail --port 8080 --smtp-port 25  # 指定端口启动"
    echo "   ./nbemail --db /path/to/nbemail.db  # 指定数据库文件"
    echo ""
    echo "🌐 默认访问地址: http://localhost:8080"
    echo "👤 默认管理员账户: admin@localhost / admin123"
    echo ""
    echo "📧 SMTP服务器配置:"
    echo "   服务器: localhost"
    echo "   端口: 25"
    echo "   认证: 无需认证（本地用户自动接收）"
else
    echo "❌ 构建失败!"
    exit 1
fi