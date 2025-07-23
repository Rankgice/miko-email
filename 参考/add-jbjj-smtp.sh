#!/bin/bash

# 为 jbjj.site 域名添加SMTP配置的脚本

echo "=== jbjj.site SMTP配置添加脚本 ==="
echo ""

# 检查服务是否运行
if ! curl -s http://127.0.0.1:8080 > /dev/null; then
    echo "❌ 邮件系统未运行，请先启动服务"
    echo "运行: ./nbemail 或 ./nbemail.exe"
    exit 1
fi

echo "✅ 邮件系统正在运行"
echo ""

# 获取管理员token（需要先登录）
echo "请确保您已经以管理员身份登录系统"
echo "如果未登录，请访问: http://127.0.0.1:8080/login"
echo ""

# 提示用户输入SMTP配置信息
echo "请输入 jbjj.site 域名的SMTP配置信息："
echo ""

read -p "SMTP服务器地址 (默认: mail.jbjj.site): " smtp_host
smtp_host=${smtp_host:-mail.jbjj.site}

read -p "SMTP端口 (默认: 587): " smtp_port
smtp_port=${smtp_port:-587}

read -p "用户名 (默认: 1keqb385916@jbjj.site): " smtp_user
smtp_user=${smtp_user:-1keqb385916@jbjj.site}

read -s -p "密码: " smtp_password
echo ""

read -p "启用TLS? (y/n, 默认: y): " enable_tls
enable_tls=${enable_tls:-y}

if [[ "$enable_tls" == "y" || "$enable_tls" == "Y" ]]; then
    tls_enabled=true
else
    tls_enabled=false
fi

echo ""
echo "配置信息："
echo "域名: jbjj.site"
echo "SMTP服务器: $smtp_host"
echo "端口: $smtp_port"
echo "用户名: $smtp_user"
echo "启用TLS: $tls_enabled"
echo ""

read -p "确认添加配置? (y/n): " confirm
if [[ "$confirm" != "y" && "$confirm" != "Y" ]]; then
    echo "已取消"
    exit 0
fi

# 构建JSON数据
json_data=$(cat <<EOF
{
    "domain": "jbjj.site",
    "host": "$smtp_host",
    "port": $smtp_port,
    "user": "$smtp_user",
    "password": "$smtp_password",
    "tls": $tls_enabled
}
EOF
)

echo ""
echo "正在添加SMTP配置..."

# 发送请求添加配置
response=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -d "$json_data" \
    -b "cookies.txt" \
    http://127.0.0.1:8080/api/smtp-configs)

# 检查响应
if echo "$response" | grep -q '"success":true'; then
    echo "✅ SMTP配置添加成功！"
    echo ""
    echo "现在您可以使用 @jbjj.site 域名的邮箱发送邮件了"
    echo "系统会自动使用刚才配置的SMTP服务器"
else
    echo "❌ 添加失败"
    echo "响应: $response"
    echo ""
    echo "可能的原因："
    echo "1. 未以管理员身份登录"
    echo "2. SMTP服务器信息不正确"
    echo "3. 网络连接问题"
fi

echo ""
echo "您也可以通过Web界面管理SMTP配置："
echo "访问: http://127.0.0.1:8080/smtp-configs"
