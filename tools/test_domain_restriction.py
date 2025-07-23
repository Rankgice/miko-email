#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import smtplib
import ssl
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart
from email.header import Header
from datetime import datetime
import sys
import time

def send_test_email(recipient_domain, recipient_user="test"):
    """发送测试邮件到指定域名"""
    
    # 163邮箱SMTP配置
    smtp_server = "smtp.163.com"
    smtp_port = 465  # SSL端口
    sender_email = "18090776855@163.com"
    sender_password = "JTH39ZMMBTennqeQ"
    
    # 构造收件人邮箱
    recipient_email = f"{recipient_user}@{recipient_domain}"
    
    print(f"📧 测试邮件发送到: {recipient_email}")
    
    try:
        # 创建邮件对象
        message = MIMEMultipart()
        message["From"] = Header(f"域名测试工具 <{sender_email}>", 'utf-8')
        message["To"] = Header(recipient_email, 'utf-8')
        message["Subject"] = Header(f"域名限制测试 - {recipient_domain}", 'utf-8')
        
        # 邮件正文
        current_time = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        body = f"""您好！

这是一封域名限制测试邮件，用于验证Miko邮箱系统是否能接收任意域名的邮件。

📧 测试目标域名：{recipient_domain}
📮 收件人邮箱：{recipient_email}
⏰ 发送时间：{current_time}
🔧 测试工具：Python域名限制测试工具

如果您在Miko邮箱系统中看到这封邮件，说明：
✅ 域名限制已成功取消
✅ 系统可以接收任意域名的邮件
✅ SMTP服务器正常工作

测试域名类型：
- 常见域名：gmail.com, yahoo.com, outlook.com
- 自定义域名：jbjj.site, example.org
- 特殊域名：test.local, demo.internal

祝好！
Miko邮箱域名测试工具
"""
        
        # 添加邮件正文
        message.attach(MIMEText(body, 'plain', 'utf-8'))
        
        print("📝 正在构建测试邮件...")
        
        # 创建SSL上下文
        context = ssl.create_default_context()
        
        print("🔐 正在连接SMTP服务器...")
        
        # 连接到服务器并发送邮件
        with smtplib.SMTP_SSL(smtp_server, smtp_port, context=context) as server:
            print("🔑 正在进行身份验证...")
            server.login(sender_email, sender_password)
            
            print("📤 正在发送邮件...")
            text = message.as_string()
            server.sendmail(sender_email, recipient_email, text)
            
        print(f"✅ 测试邮件发送成功！目标域名：{recipient_domain}")
        return True
        
    except Exception as e:
        print(f"❌ 发送失败: {e}")
        return False

def test_multiple_domains():
    """测试多个不同类型的域名"""
    
    print("🚀 Miko邮箱系统域名限制测试工具")
    print("=" * 60)
    print("📝 此工具将测试系统是否能接收各种域名的邮件")
    print("🎯 如果域名限制已取消，所有测试邮件都应该能成功接收")
    print("=" * 60)
    print("")
    
    # 测试域名列表
    test_domains = [
        ("jbjj.site", "kimi", "项目域名"),
        ("gmail.com", "testuser", "常见邮箱域名"),
        ("yahoo.com", "testuser", "常见邮箱域名"),
        ("outlook.com", "testuser", "常见邮箱域名"),
        ("example.org", "admin", "示例域名"),
        ("test.local", "user", "本地测试域名"),
        ("demo.internal", "test", "内部测试域名"),
        ("custom-domain.xyz", "hello", "自定义域名"),
        ("very-long-domain-name-for-testing.com", "user", "长域名测试"),
        ("sub.domain.example.net", "test", "子域名测试")
    ]
    
    successful_tests = 0
    total_tests = len(test_domains)
    
    for i, (domain, user, description) in enumerate(test_domains, 1):
        print(f"🧪 测试 {i}/{total_tests}: {description}")
        print(f"   域名: {domain}")
        print(f"   邮箱: {user}@{domain}")
        
        if send_test_email(domain, user):
            successful_tests += 1
            print(f"   ✅ 成功")
        else:
            print(f"   ❌ 失败")
        
        print("")
        
        # 避免发送过快
        if i < total_tests:
            print("⏳ 等待2秒后继续...")
            time.sleep(2)
    
    print("=" * 60)
    print("📊 测试结果汇总")
    print("=" * 60)
    print(f"✅ 成功发送: {successful_tests}/{total_tests}")
    print(f"❌ 发送失败: {total_tests - successful_tests}/{total_tests}")
    
    if successful_tests == total_tests:
        print("🎉 所有测试邮件发送成功！")
        print("💡 请检查Miko邮箱系统是否收到了这些邮件")
    else:
        print("⚠️  部分测试邮件发送失败")
        print("💡 这可能是网络问题，而非域名限制问题")
    
    print("")
    print("🔍 验证步骤：")
    print("1. 访问：http://localhost:8080/inbox")
    print("2. 登录任意邮箱账号")
    print("3. 查看是否收到来自不同域名的测试邮件")
    print("4. 如果收到邮件，说明域名限制已成功取消")
    print("=" * 60)

def test_single_domain():
    """测试单个域名"""
    
    if len(sys.argv) < 3:
        print("用法: python test_domain_restriction.py single <域名> [用户名]")
        print("示例: python test_domain_restriction.py single example.com testuser")
        return
    
    domain = sys.argv[2]
    user = sys.argv[3] if len(sys.argv) > 3 else "test"
    
    print("🚀 Miko邮箱系统单域名测试")
    print("=" * 40)
    print(f"📧 目标域名: {domain}")
    print(f"👤 用户名: {user}")
    print("=" * 40)
    print("")
    
    if send_test_email(domain, user):
        print("")
        print("🎉 测试邮件发送成功！")
        print("💡 请检查Miko邮箱系统是否收到邮件")
        print(f"📬 收件人: {user}@{domain}")
    else:
        print("")
        print("❌ 测试邮件发送失败")

def show_usage():
    """显示使用说明"""
    print("Miko邮箱系统域名限制测试工具")
    print("")
    print("用法:")
    print("  python tools/test_domain_restriction.py <命令> [参数]")
    print("")
    print("命令:")
    print("  multi   - 测试多个不同类型的域名")
    print("  single  - 测试单个指定域名")
    print("  help    - 显示此帮助信息")
    print("")
    print("示例:")
    print("  python tools/test_domain_restriction.py multi")
    print("  python tools/test_domain_restriction.py single example.com testuser")
    print("")
    print("说明:")
    print("  此工具通过发送测试邮件来验证Miko邮箱系统是否")
    print("  已取消域名限制，能够接收任意域名的邮件。")

if __name__ == "__main__":
    if len(sys.argv) < 2:
        show_usage()
        sys.exit(1)
    
    command = sys.argv[1].lower()
    
    if command == "multi":
        test_multiple_domains()
    elif command == "single":
        test_single_domain()
    elif command == "help":
        show_usage()
    else:
        print(f"未知命令: {command}")
        show_usage()
        sys.exit(1)
