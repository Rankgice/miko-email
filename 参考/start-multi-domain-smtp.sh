#!/bin/bash

echo "启动NBEmail系统 - 多域名SMTP配置"
echo "========================================"
echo

# 检查是否存在配置文件
if [ -f ".env" ]; then
    echo "✅ 发现配置文件 .env"
    source .env
else
    echo "⚠️  未发现 .env 配置文件，使用示例配置"
    echo "   请复制 multi-domain-smtp-example.env 为 .env 并修改配置"
    echo
fi

# 设置多域名SMTP配置示例
echo "配置多域名SMTP服务器..."

# 域名1：example.com
export DOMAIN_SMTP_EXAMPLE_COM="mail.example.com:587:smtp@example.com:your-password:true"

# 域名2：company.com  
export DOMAIN_SMTP_COMPANY_COM="smtp.company.com:465:noreply@company.com:company-password:true"

# 域名3：mysite.org
export DOMAIN_SMTP_MYSITE_ORG="mail.mysite.org:587:sender@mysite.org:mysite-password:true"

# 默认SMTP配置（向后兼容）
export OUTBOUND_SMTP_HOST="mail.yourdomain.com"
export OUTBOUND_SMTP_PORT=587
export OUTBOUND_SMTP_USER="default@yourdomain.com"
export OUTBOUND_SMTP_PASSWORD="default-password"
export OUTBOUND_SMTP_TLS=true

# 设置域名
export DOMAIN="yourdomain.com"

echo
echo "多域名SMTP配置："
echo "  example.com -> mail.example.com:587"
echo "  company.com -> smtp.company.com:465"  
echo "  mysite.org  -> mail.mysite.org:587"
echo "  默认配置    -> mail.yourdomain.com:587"
echo
echo "工作原理："
echo "  📧 user@example.com 发邮件 -> 使用 mail.example.com"
echo "  📧 admin@company.com 发邮件 -> 使用 smtp.company.com"
echo "  📧 info@mysite.org 发邮件 -> 使用 mail.mysite.org"
echo "  📧 其他邮箱发邮件 -> 使用默认SMTP服务器"
echo

echo "启动NBEmail服务器..."
echo "Web界面: http://localhost:8080"
echo "SMTP服务: localhost:2525"
echo

# 检查可执行文件
if [ ! -f "./nbemail" ]; then
    echo "❌ 未找到 nbemail 可执行文件"
    echo "   请先编译项目: go build -o nbemail"
    exit 1
fi

# 启动服务器
./nbemail --port 8080 --smtp-port 2525
