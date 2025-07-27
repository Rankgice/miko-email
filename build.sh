#!/bin/bash

echo "构建Miko邮箱系统..."

# 检查Node.js是否安装
if ! command -v node &> /dev/null; then
    echo "错误: Node.js未安装，请先安装Node.js"
    exit 1
fi

# 检查npm是否安装
if ! command -v npm &> /dev/null; then
    echo "错误: npm未安装，请先安装npm"
    exit 1
fi

# 构建前端Vue应用
echo "构建前端Vue应用..."
cd frontend

# 安装依赖
echo "安装前端依赖..."
npm install

# 构建生产版本
echo "构建生产版本..."
npm run build

# 返回根目录
cd ..

# 构建Go后端
echo "构建Go后端..."
go mod tidy
go build -o miko-email main.go

echo "构建完成！"
echo ""
echo "启动说明："
echo "1. 确保配置文件 config.yaml 存在"
echo "2. 运行: ./miko-email"
echo "3. 访问: http://localhost:8080"
echo ""
echo "注意: 系统已强制使用SQLite数据库，禁止使用MySQL"
