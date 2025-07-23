# NBEmail Makefile

.PHONY: build clean run dev test install help

# 默认目标
all: build

# 构建Linux可执行文件
build:
	@echo "🚀 构建 NBEmail Linux 可执行文件..."
	@chmod +x build.sh
	@./build.sh

# 构建Windows可执行文件
build-windows:
	@echo "🚀 构建 NBEmail Windows 可执行文件..."
	@CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o nbemail.exe main.go
	@echo "✅ Windows版本构建完成: nbemail.exe"

# 构建macOS可执行文件
build-macos:
	@echo "🚀 构建 NBEmail macOS 可执行文件..."
	@CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o nbemail-macos main.go
	@echo "✅ macOS版本构建完成: nbemail-macos"

# 构建所有平台
build-all: build build-windows build-macos
	@echo "✅ 所有平台构建完成"

# 清理构建文件
clean:
	@echo "🧹 清理构建文件..."
	@rm -f nbemail nbemail.exe nbemail-macos
	@rm -f nbemail.db
	@echo "✅ 清理完成"

# 运行程序
run: build
	@echo "🚀 启动 NBEmail..."
	@./nbemail

# 开发模式运行
dev:
	@echo "🔧 开发模式启动..."
	@go run main.go

# 运行测试
test:
	@echo "🧪 运行测试..."
	@go test -v ./...

# 安装依赖
install:
	@echo "📦 安装依赖..."
	@go mod tidy
	@go mod download

# 格式化代码
fmt:
	@echo "🎨 格式化代码..."
	@go fmt ./...

# 代码检查
lint:
	@echo "🔍 代码检查..."
	@go vet ./...

# 创建发布包
release: clean build-all
	@echo "📦 创建发布包..."
	@mkdir -p release
	@cp nbemail release/nbemail-linux-amd64
	@cp nbemail.exe release/nbemail-windows-amd64.exe
	@cp nbemail-macos release/nbemail-darwin-amd64
	@cp README.md release/ 2>/dev/null || echo "README.md not found"
	@echo "✅ 发布包已创建在 release/ 目录"

# 显示帮助信息
help:
	@echo "NBEmail 构建工具"
	@echo ""
	@echo "可用命令:"
	@echo "  build         构建Linux可执行文件"
	@echo "  build-windows 构建Windows可执行文件"
	@echo "  build-macos   构建macOS可执行文件"
	@echo "  build-all     构建所有平台"
	@echo "  clean         清理构建文件"
	@echo "  run           构建并运行程序"
	@echo "  dev           开发模式运行"
	@echo "  test          运行测试"
	@echo "  install       安装依赖"
	@echo "  fmt           格式化代码"
	@echo "  lint          代码检查"
	@echo "  release       创建发布包"
	@echo "  help          显示此帮助信息"