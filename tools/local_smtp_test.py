#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
直接向本地SMTP服务器发送测试邮件
"""

import smtplib
import logging
from email.mime.text import MIMEText
from email.mime.multipart import MIMEMultipart
from datetime import datetime

# 配置日志
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)

def send_test_email():
    """直接向本地SMTP服务器发送测试邮件"""
    try:
        # 邮件配置
        smtp_server = "localhost"
        smtp_port = 25
        
        # 发件人和收件人
        sender = "test@example.com"
        recipient = "kimi11@jbjj.site"
        
        # 创建邮件
        msg = MIMEMultipart()
        msg['From'] = sender
        msg['To'] = recipient
        msg['Subject'] = f"本地SMTP测试邮件 - {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}"
        
        # 邮件正文
        body = f"""
这是一封测试邮件，用于验证本地SMTP服务器功能。

发送时间: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}
发件人: {sender}
收件人: {recipient}
SMTP服务器: {smtp_server}:{smtp_port}

如果您收到这封邮件，说明本地SMTP服务器工作正常！
        """
        
        msg.attach(MIMEText(body, 'plain', 'utf-8'))
        
        logging.info("=" * 60)
        logging.info("本地SMTP测试邮件发送器")
        logging.info("=" * 60)
        logging.info(f"SMTP服务器: {smtp_server}:{smtp_port}")
        logging.info(f"发件人: {sender}")
        logging.info(f"收件人: {recipient}")
        logging.info("=" * 60)
        
        # 连接到SMTP服务器
        logging.info("正在连接到本地SMTP服务器...")
        server = smtplib.SMTP(smtp_server, smtp_port)
        
        # 启用调试模式
        server.set_debuglevel(1)
        
        # 发送邮件
        logging.info("正在发送邮件...")
        text = msg.as_string()
        server.sendmail(sender, recipient, text)
        
        # 关闭连接
        server.quit()
        
        logging.info("✅ 邮件发送成功！")
        logging.info("请检查IMAP接收器是否检测到新邮件")
        
    except Exception as e:
        logging.error(f"❌ 邮件发送失败: {str(e)}")
        import traceback
        logging.error(traceback.format_exc())

if __name__ == "__main__":
    send_test_email()
