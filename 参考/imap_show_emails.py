#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
IMAP邮件显示脚本
显示邮箱中所有邮件的详细信息
"""

import imaplib
import email
import sys
import logging
import email.header
import ssl
from email.parser import BytesParser
from email.policy import default

# 配置日志
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s',
    handlers=[
        logging.StreamHandler(sys.stdout)
    ]
)

class EmailViewer:
    def __init__(self, email_address, password):
        """初始化IMAP连接参数"""
        self.imap_server = "127.0.0.1"
        self.imap_port = 143
        self.email_address = email_address
        self.password = password
        
    def connect_to_server(self):
        """连接到IMAP服务器"""
        try:
            logging.info(f"[{self.email_address}] 连接到IMAP服务器...")
            mail = imaplib.IMAP4(self.imap_server, self.imap_port)
            
            logging.info(f"[{self.email_address}] 尝试登录...")
            mail.login(self.email_address, self.password)
            logging.info(f"[{self.email_address}] 登录成功")
            
            return mail
        except Exception as e:
            logging.error(f"[{self.email_address}] 连接失败: {e}")
            return None
    
    def decode_header(self, header_value):
        """解码邮件头"""
        if not header_value:
            return ""
        
        try:
            decoded_parts = email.header.decode_header(header_value)
            decoded_string = ""
            for part, encoding in decoded_parts:
                if isinstance(part, bytes):
                    if encoding:
                        decoded_string += part.decode(encoding)
                    else:
                        decoded_string += part.decode('utf-8', errors='ignore')
                else:
                    decoded_string += part
            return decoded_string
        except Exception as e:
            logging.warning(f"解码头部失败: {e}")
            return str(header_value)
    
    def get_email_body(self, email_message):
        """获取邮件正文"""
        body = ""
        try:
            if email_message.is_multipart():
                for part in email_message.walk():
                    content_type = part.get_content_type()
                    if content_type == "text/plain":
                        charset = part.get_content_charset() or 'utf-8'
                        payload = part.get_payload(decode=True)
                        if isinstance(payload, bytes):
                            body = payload.decode(charset, errors='ignore')
                        else:
                            body = str(payload)
                        break
                    elif content_type == "text/html" and not body:
                        charset = part.get_content_charset() or 'utf-8'
                        payload = part.get_payload(decode=True)
                        if isinstance(payload, bytes):
                            body = payload.decode(charset, errors='ignore')
                        else:
                            body = str(payload)
            else:
                charset = email_message.get_content_charset() or 'utf-8'
                payload = email_message.get_payload(decode=True)
                if isinstance(payload, bytes):
                    body = payload.decode(charset, errors='ignore')
                else:
                    body = str(payload)
        except Exception as e:
            logging.warning(f"获取邮件正文失败: {e}")
            body = "无法解析邮件正文"

        return body
    
    def fetch_email(self, mail, email_id):
        """获取单个邮件内容"""
        try:
            status, data = mail.fetch(email_id, '(RFC822)')
            
            if status != 'OK':
                logging.error(f"无法获取邮件 ID: {email_id}")
                return None
            
            # 解析邮件内容
            raw_email = data[0][1]
            email_message = email.message_from_bytes(raw_email)
            
            # 获取邮件信息
            subject = self.decode_header(email_message['Subject'])
            from_address = self.decode_header(email_message['From'])
            to_address = self.decode_header(email_message['To'])
            date_str = email_message['Date']
            
            # 获取邮件正文
            body = self.get_email_body(email_message)
            
            return {
                'id': email_id.decode() if isinstance(email_id, bytes) else str(email_id),
                'subject': subject,
                'from': from_address,
                'to': to_address,
                'date': date_str,
                'body': body[:500] + "..." if len(body) > 500 else body  # 限制显示长度
            }
            
        except Exception as e:
            logging.error(f"获取邮件失败: {e}")
            return None
    
    def show_all_emails(self):
        """显示所有邮件"""
        mail = self.connect_to_server()
        if not mail:
            return
        
        try:
            # 选择收件箱
            mail.select('INBOX')
            
            # 搜索所有邮件
            status, messages = mail.search(None, 'ALL')
            
            if status != 'OK':
                logging.error("无法搜索邮件")
                return
            
            # 获取邮件ID列表
            email_ids = messages[0].split()
            
            if not email_ids:
                logging.info("📭 收件箱为空")
                return
            
            logging.info(f"📧 找到 {len(email_ids)} 封邮件")
            print("=" * 80)
            
            # 显示每封邮件的详细信息
            for i, email_id in enumerate(reversed(email_ids), 1):  # 从最新的开始显示
                email_data = self.fetch_email(mail, email_id)
                if email_data:
                    print(f"\n📧 邮件 {i}/{len(email_ids)} (ID: {email_data['id']})")
                    print(f"📅 日期: {email_data['date']}")
                    print(f"👤 发件人: {email_data['from']}")
                    print(f"👥 收件人: {email_data['to']}")
                    print(f"📝 主题: {email_data['subject']}")
                    print(f"📄 内容预览:")
                    print(f"   {email_data['body']}")
                    print("-" * 80)
            
            # 登出
            mail.logout()
            logging.info("✅ 邮件显示完成")
            
        except Exception as e:
            logging.error(f"显示邮件失败: {e}")

def main():
    if len(sys.argv) < 3:
        print("用法: python imap_show_emails.py <邮箱地址> <密码>")
        print("示例: python imap_show_emails.py 2014131458@qq.com tgx123456")
        return
    
    email_address = sys.argv[1]
    password = sys.argv[2]
    
    print("=" * 80)
    print(f"📧 NBEmail IMAP邮件查看器")
    print(f"📮 邮箱: {email_address}")
    print("=" * 80)
    
    viewer = EmailViewer(email_address, password)
    viewer.show_all_emails()

if __name__ == "__main__":
    main()
