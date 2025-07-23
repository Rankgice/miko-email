#!/bin/bash

# NBEmail 自动生成用户名密码功能演示脚本

echo "🎯 NBEmail 自动生成用户名密码功能演示"
echo "========================================"
echo

echo "🚀 新功能亮点："
echo "  ✨ 智能用户名生成 - 根据域名自动推荐格式"
echo "  🔑 强密码自动生成 - 16位超强密码"
echo "  🔧 完全自动化配置 - 一键生成完整认证信息"
echo

echo "🎮 使用方法："
echo "  1. 自动配置时生成 - 点击'自动配置'按钮"
echo "  2. 手动添加时生成 - 点击'🎯 智能生成'按钮"
echo "  3. 编辑配置时生成 - 在编辑框中点击生成按钮"
echo

echo "🎯 生成规则："
echo "  用户名格式："
echo "    • smtp@yourdomain.com"
echo "    • noreply@yourdomain.com"
echo "    • mail@yourdomain.com"
echo "    • sender@yourdomain.com"
echo "    • system@yourdomain.com"
echo
echo "  密码规则："
echo "    • 长度：16位"
echo "    • 字符：大小写字母+数字+符号"
echo "    • 安全：加密级随机生成"
echo

# 设置演示环境
export WEB_PORT=8080
export SMTP_PORT=2525
export DB_PATH="demo-auto-generate.db"
export DOMAIN="localhost"
export ADMIN_EMAIL="admin@localhost"
export ADMIN_PASS="admin123"

echo "🔧 启动演示环境："
echo "  Web界面: http://localhost:8080"
echo "  SMTP配置: http://localhost:8080/smtp-configs"
echo "  管理员: admin@localhost / admin123"
echo

echo "📋 演示步骤："
echo "  1. 访问 http://localhost:8080/smtp-configs"
echo "  2. 点击'自动配置'体验自动生成"
echo "  3. 点击'手动添加'体验手动生成"
echo "  4. 编辑现有配置体验编辑生成"
echo

echo "🎉 开始演示..."
echo "按 Ctrl+C 停止服务器"
echo

# 启动NBEmail
go run .
