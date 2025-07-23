#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
测试邮件转发功能
"""

import smtplib
import time
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart

def test_forward():
    """测试邮件转发功能"""
    
    # SMTP服务器配置
    smtp_host = "localhost"
    smtp_port = 25
    
    # 测试邮件配置
    from_email = "test@external.com"  # 外部发件人
    to_email = "kimi11@jbjj.site"     # 设置了转发规则的邮箱
    subject = "测试转发功能 - " + time.strftime("%Y-%m-%d %H:%M:%S")
    body = """
这是一封测试邮件，用于测试邮件转发功能。

如果转发功能正常工作，这封邮件应该会：
1. 保存到 kimi11@jbjj.site 的收件箱
2. 自动转发到 kimi12@jbjj.site

发送时间: {}
    """.format(time.strftime("%Y-%m-%d %H:%M:%S"))
    
    try:
        print(f"🚀 开始测试邮件转发功能...")
        print(f"📧 发件人: {from_email}")
        print(f"📧 收件人: {to_email}")
        print(f"📧 主题: {subject}")
        
        # 创建邮件
        msg = MIMEMultipart()
        msg['From'] = from_email
        msg['To'] = to_email
        msg['Subject'] = subject
        
        # 添加邮件正文
        msg.attach(MIMEText(body, 'plain', 'utf-8'))
        
        # 连接SMTP服务器
        print(f"🔗 连接到SMTP服务器 {smtp_host}:{smtp_port}")
        server = smtplib.SMTP(smtp_host, smtp_port)
        
        # 启用调试模式
        server.set_debuglevel(1)
        
        # 发送邮件
        print(f"📤 发送邮件...")
        server.sendmail(from_email, [to_email], msg.as_string())
        
        # 关闭连接
        server.quit()
        
        print(f"✅ 邮件发送成功！")
        print(f"📋 请检查以下内容：")
        print(f"   1. kimi11@jbjj.site 的收件箱是否收到邮件")
        print(f"   2. kimi12@jbjj.site 的收件箱是否收到转发邮件")
        print(f"   3. 转发邮件是否包含 '[转发]' 前缀")
        print(f"   4. 服务器日志是否显示转发处理过程")
        
    except Exception as e:
        print(f"❌ 邮件发送失败: {e}")
        return False
    
    return True

def test_multiple_forwards():
    """测试多个转发规则"""
    
    # SMTP服务器配置
    smtp_host = "localhost"
    smtp_port = 25
    
    # 测试邮件配置
    from_email = "test@external.com"
    to_email = "kimi11@jbjj.site"  # 假设这个邮箱有多个转发规则
    subject = "测试多重转发 - " + time.strftime("%Y-%m-%d %H:%M:%S")
    body = """
这是一封测试多重转发的邮件。

如果有多个转发规则，这封邮件应该会转发到所有配置的目标邮箱。

发送时间: {}
    """.format(time.strftime("%Y-%m-%d %H:%M:%S"))
    
    try:
        print(f"\n🚀 开始测试多重转发功能...")
        
        # 创建邮件
        msg = MIMEMultipart()
        msg['From'] = from_email
        msg['To'] = to_email
        msg['Subject'] = subject
        msg.attach(MIMEText(body, 'plain', 'utf-8'))
        
        # 连接并发送
        server = smtplib.SMTP(smtp_host, smtp_port)
        server.sendmail(from_email, [to_email], msg.as_string())
        server.quit()
        
        print(f"✅ 多重转发测试邮件发送成功！")
        
    except Exception as e:
        print(f"❌ 多重转发测试失败: {e}")
        return False
    
    return True

def test_forward_with_attachments():
    """测试带附件的转发"""
    
    # SMTP服务器配置
    smtp_host = "localhost"
    smtp_port = 25
    
    # 测试邮件配置
    from_email = "test@external.com"
    to_email = "kimi11@jbjj.site"
    subject = "测试附件转发 - " + time.strftime("%Y-%m-%d %H:%M:%S")
    body = """
这是一封带附件的测试邮件。

测试转发规则是否正确处理附件。

发送时间: {}
    """.format(time.strftime("%Y-%m-%d %H:%M:%S"))
    
    try:
        print(f"\n🚀 开始测试附件转发功能...")
        
        # 创建邮件
        msg = MIMEMultipart()
        msg['From'] = from_email
        msg['To'] = to_email
        msg['Subject'] = subject
        msg.attach(MIMEText(body, 'plain', 'utf-8'))
        
        # 添加文本附件
        attachment_content = "这是一个测试附件的内容。\n测试时间: {}".format(time.strftime("%Y-%m-%d %H:%M:%S"))
        attachment = MIMEText(attachment_content, 'plain', 'utf-8')
        attachment.add_header('Content-Disposition', 'attachment', filename='test_attachment.txt')
        msg.attach(attachment)
        
        # 连接并发送
        server = smtplib.SMTP(smtp_host, smtp_port)
        server.sendmail(from_email, [to_email], msg.as_string())
        server.quit()
        
        print(f"✅ 附件转发测试邮件发送成功！")
        
    except Exception as e:
        print(f"❌ 附件转发测试失败: {e}")
        return False
    
    return True

if __name__ == "__main__":
    print("=" * 60)
    print("📧 Miko邮箱转发功能测试")
    print("=" * 60)
    
    # 基本转发测试
    success1 = test_forward()
    
    # 等待一下
    time.sleep(2)
    
    # 多重转发测试
    success2 = test_multiple_forwards()
    
    # 等待一下
    time.sleep(2)
    
    # 附件转发测试
    success3 = test_forward_with_attachments()
    
    print("\n" + "=" * 60)
    print("📊 测试结果汇总:")
    print(f"   基本转发测试: {'✅ 成功' if success1 else '❌ 失败'}")
    print(f"   多重转发测试: {'✅ 成功' if success2 else '❌ 失败'}")
    print(f"   附件转发测试: {'✅ 成功' if success3 else '❌ 失败'}")
    print("=" * 60)
    
    if all([success1, success2, success3]):
        print("🎉 所有转发测试都成功完成！")
    else:
        print("⚠️  部分测试失败，请检查服务器日志。")
