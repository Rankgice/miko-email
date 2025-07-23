#!/bin/bash

echo "🚀 NBEmail 自动SMTP配置演示"
echo "================================"
echo

# 检查可执行文件
if [ ! -f "./nbemail" ]; then
    echo "❌ 未找到 nbemail 可执行文件"
    echo "   请先编译项目: go build -o nbemail"
    exit 1
fi

echo "📋 功能演示："
echo "  1. 系统自动扫描现有域名邮箱"
echo "  2. 自动生成对应的SMTP配置"
echo "  3. 用户只需填写认证信息"
echo "  4. 发件时自动选择正确的SMTP服务器"
echo

echo "🎯 使用场景："
echo "  • 个人多域名邮箱管理"
echo "  • 企业多部门邮件系统"
echo "  • 多品牌邮件发送"
echo

echo "⚡ 自动配置原理："
echo "  user@example.com  -> 自动配置 mail.example.com:587"
echo "  admin@company.org -> 自动配置 smtp.company.org:587"
echo "  info@mysite.net   -> 自动配置 mail.mysite.net:587"
echo

# 设置基础配置
export WEB_PORT=8080
export SMTP_PORT=2525
export DB_PATH="demo-nbemail.db"
export DOMAIN="localhost"
export ADMIN_EMAIL="admin@localhost"
export ADMIN_PASS="admin123"

echo "🔧 启动配置："
echo "  Web界面: http://localhost:8080"
echo "  SMTP服务: localhost:2525"
echo "  管理员: admin@localhost / admin123"
echo

echo "📖 使用步骤："
echo "  1. 访问 http://localhost:8080"
echo "  2. 登录管理界面"
echo "  3. 创建一些不同域名的邮箱账户"
echo "  4. 访问 'SMTP配置' 页面"
echo "  5. 点击 '自动配置' 按钮"
echo "  6. 编辑配置添加认证信息"
echo "  7. 开始发送邮件测试"
echo

echo "🎉 特色功能："
echo "  ✅ 一键自动配置多域名SMTP"
echo "  ✅ 智能检测SMTP服务器地址"
echo "  ✅ 自动匹配发件人与SMTP配置"
echo "  ✅ 可视化配置管理界面"
echo "  ✅ 支持TLS加密和认证"
echo

echo "启动NBEmail演示系统..."
echo "按 Ctrl+C 停止服务"
echo

# 启动服务器
./nbemail --port 8080 --smtp-port 2525
