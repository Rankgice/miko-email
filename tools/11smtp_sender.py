#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import smtplib
import ssl
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart
from email.header import Header
from datetime import datetime
import sys

def send_test_email():
    """使用Miko邮箱发送测试邮件"""

    # Miko邮箱SMTP配置
    smtp_server = "118.120.221.169"  # 你的Miko邮箱服务器
    smtp_port = 25  # 非SSL端口
    sender_email = "kimi@jbjj.site"
    sender_password = "06c3c4d1"
    
    # 收件人
    recipient_email = "2014131458@qq.com"
    
    print("=" * 50)
    print("📧 SMTP发件测试工具")
    print("=" * 50)
    print(f"发件人: {sender_email}")
    print(f"收件人: {recipient_email}")
    print(f"SMTP服务器: {smtp_server}:{smtp_port}")
    print("-" * 50)
    
    try:
        # 创建邮件对象
        message = MIMEMultipart()
        message["From"] = Header(f"Miko邮箱测试 <{sender_email}>", 'utf-8')
        message["To"] = Header(recipient_email, 'utf-8')
        message["Subject"] = Header("SMTP发件测试 - 来自Miko邮箱", 'utf-8')
        
        # 邮件正文
        current_time = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        body = f"""这是一封来自Miko邮箱的测试邮件。

📧 发送目的：测试SMTP发送功能
⏰ 发送时间：{current_time}
🔧 发送工具：Python SMTP测试工具
📍 发件服务器：{smtp_server}

如果您收到这封邮件，说明Miko邮箱的SMTP发送功能正常。

这封邮件是通过Miko邮箱系统的SMTP服务器发送的，
证明了外部邮件发送功能工作正常。

祝好！
Miko邮箱系统测试
"""
        
        # 添加邮件正文
        message.attach(MIMEText(body, 'plain', 'utf-8'))
        
        print("📝 正在构建邮件内容...")
        
        print("🔐 正在连接SMTP服务器...")

        # 连接到服务器并发送邮件（使用普通SMTP，不是SSL）
        with smtplib.SMTP(smtp_server, smtp_port) as server:
            print("🔑 正在进行身份验证...")
            server.login(sender_email, sender_password)

            print("📤 正在发送邮件...")
            text = message.as_string()
            server.sendmail(sender_email, recipient_email, text)
            
        print("✅ 邮件发送成功！")
        print("-" * 50)
        print("📬 请检查QQ邮箱是否收到测试邮件")
        print("📧 如果收到邮件，说明Miko邮箱SMTP发送功能正常")
        print("")
        print("🔄 测试结果：")
        print("   ✅ Miko邮箱SMTP服务器工作正常")
        print("   ✅ 外部邮件发送功能正常")
        print("=" * 50)
        
        return True
        
    except smtplib.SMTPAuthenticationError as e:
        print(f"❌ SMTP认证失败: {e}")
        print("💡 请检查邮箱地址和密码是否正确")
        return False
        
    except smtplib.SMTPConnectError as e:
        print(f"❌ 连接SMTP服务器失败: {e}")
        print("💡 请检查网络连接和服务器地址")
        return False
        
    except smtplib.SMTPException as e:
        print(f"❌ SMTP错误: {e}")
        return False
        
    except Exception as e:
        print(f"❌ 发送失败: {e}")
        return False

def send_to_miko_system():
    """发送邮件到Miko邮箱系统进行接收测试"""
    
    # 163邮箱SMTP配置
    smtp_server = "smtp.163.com"
    smtp_port = 465  # SSL端口
    sender_email = "18090776855@163.com"
    sender_password = "JTH39ZMMBTennqeQ"
    
    # 收件人 - Miko邮箱系统
    recipient_email = "kimi@jbjj.site"
    
    print("=" * 50)
    print("📧 测试Miko邮箱系统接收功能")
    print("=" * 50)
    print(f"发件人: {sender_email}")
    print(f"收件人: {recipient_email}")
    print(f"SMTP服务器: {smtp_server}:{smtp_port}")
    print("-" * 50)
    
    try:
        # 创建邮件对象
        message = MIMEMultipart()
        message["From"] = Header(f"163邮箱测试 <{sender_email}>", 'utf-8')
        message["To"] = Header(recipient_email, 'utf-8')
        message["Subject"] = Header("测试邮件 - 来自163邮箱", 'utf-8')
        
        # 邮件正文
        current_time = datetime.now().strftime("%Y-%m-%d %H:%M:%S")
        body = f"""您好！

这是一封来自163邮箱的测试邮件，用于测试Miko邮箱系统的接收功能。

📧 测试目的：验证邮件系统SMTP接收功能
⏰ 发送时间：{current_time}
🔧 发送工具：Python SMTP测试工具
📍 发件服务器：{smtp_server}
📮 目标系统：Miko邮箱系统 (jbjj.site)

如果您在Miko邮箱系统中看到这封邮件，说明：
✅ SMTP服务器正常监听25端口
✅ 邮件接收功能正常
✅ 邮件解析和保存功能正常

测试内容包括：
- 中文字符编码
- MIME格式解析
- 数据库保存

祝好！
Miko邮箱系统测试
"""
        
        # 添加邮件正文
        message.attach(MIMEText(body, 'plain', 'utf-8'))
        
        print("📝 正在构建邮件内容...")
        
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
            
        print("✅ 邮件发送成功！")
        print("-" * 50)
        print("📬 请检查Miko邮箱系统是否收到测试邮件")
        print("🌐 访问：http://localhost:8080/inbox")
        print("📧 登录邮箱：kimi11@jbjj.site")
        print("=" * 50)
        
        return True
        
    except Exception as e:
        print(f"❌ 发送失败: {e}")
        return False

if __name__ == "__main__":
    print("🚀 Miko邮箱系统 - SMTP测试工具")
    print("")
    
    if len(sys.argv) > 1 and sys.argv[1] == "miko":
        # 直接发送到Miko系统
        send_to_miko_system()
    else:
        # 先发送到QQ邮箱测试
        print("第一步：测试163邮箱SMTP发送功能")
        success = send_test_email()
        
        if success:
            print("")
            choice = input("是否继续测试Miko邮箱系统接收功能？(y/n): ").lower().strip()
            if choice in ['y', 'yes', '是']:
                print("")
                print("第二步：测试Miko邮箱系统接收功能")
                send_to_miko_system()
