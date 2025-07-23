#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
IMAP测试工具
用于测试Miko邮箱系统的IMAP功能
登录方式：网站登录账号 + 域名邮箱 + 邮箱密码
"""

import imaplib
import sys
import argparse
import ssl
import socket
from datetime import datetime

class IMAPTester:
    def __init__(self, host='localhost', port=143, use_ssl=False):
        self.host = host
        self.port = port
        self.use_ssl = use_ssl
        self.imap = None
        
    def connect(self):
        """连接到IMAP服务器"""
        try:
            print(f"🔗 正在连接到IMAP服务器 {self.host}:{self.port}")
            
            if self.use_ssl:
                # 使用SSL连接
                context = ssl.create_default_context()
                context.check_hostname = False
                context.verify_mode = ssl.CERT_NONE
                self.imap = imaplib.IMAP4_SSL(self.host, self.port, ssl_context=context)
                print("✅ SSL连接成功")
            else:
                # 使用普通连接
                self.imap = imaplib.IMAP4(self.host, self.port)
                print("✅ 连接成功")
                
            return True
            
        except Exception as e:
            print(f"❌ 连接失败: {e}")
            return False
    
    def login(self, username, email, password):
        """
        登录IMAP服务器
        username: 网站登录账号
        email: 域名邮箱
        password: 邮箱密码
        """
        if not self.imap:
            print("❌ 请先连接到服务器")
            return False
            
        try:
            print(f"🔐 正在登录...")
            print(f"   用户名: {username}")
            print(f"   邮箱: {email}")
            print(f"   密码: {'*' * len(password)}")
            
            # 尝试不同的登录方式
            login_attempts = [
                (username, password),  # 使用用户名登录
                (email, password),     # 使用邮箱登录
                (f"{username}@{email.split('@')[1]}", password),  # 用户名@域名
            ]
            
            for attempt_user, attempt_pass in login_attempts:
                try:
                    print(f"🔄 尝试登录: {attempt_user}")
                    result = self.imap.login(attempt_user, attempt_pass)
                    
                    if result[0] == 'OK':
                        print(f"✅ 登录成功! 使用凭据: {attempt_user}")
                        return True
                        
                except Exception as login_error:
                    print(f"⚠️  登录尝试失败 ({attempt_user}): {login_error}")
                    continue
            
            print("❌ 所有登录尝试都失败了")
            return False
            
        except Exception as e:
            print(f"❌ 登录过程出错: {e}")
            return False
    
    def list_folders(self):
        """列出所有文件夹"""
        if not self.imap:
            return False
            
        try:
            print("\n📁 获取文件夹列表...")
            result, folders = self.imap.list()
            
            if result == 'OK':
                print("✅ 文件夹列表:")
                for folder in folders:
                    folder_str = folder.decode('utf-8') if isinstance(folder, bytes) else str(folder)
                    print(f"   📂 {folder_str}")
                return True
            else:
                print(f"❌ 获取文件夹失败: {folders}")
                return False
                
        except Exception as e:
            print(f"❌ 列出文件夹时出错: {e}")
            return False
    
    def select_inbox(self):
        """选择收件箱"""
        if not self.imap:
            return False
            
        try:
            print("\n📥 选择收件箱...")
            result, data = self.imap.select('INBOX')
            
            if result == 'OK':
                message_count = data[0].decode('utf-8') if data[0] else '0'
                print(f"✅ 收件箱选择成功，共有 {message_count} 封邮件")
                return True
            else:
                print(f"❌ 选择收件箱失败: {data}")
                return False
                
        except Exception as e:
            print(f"❌ 选择收件箱时出错: {e}")
            return False
    
    def search_emails(self, criteria='ALL'):
        """搜索邮件"""
        if not self.imap:
            return []
            
        try:
            print(f"\n🔍 搜索邮件 (条件: {criteria})...")
            result, data = self.imap.search(None, criteria)
            
            if result == 'OK':
                email_ids = data[0].split() if data[0] else []
                print(f"✅ 找到 {len(email_ids)} 封邮件")
                return email_ids
            else:
                print(f"❌ 搜索邮件失败: {data}")
                return []
                
        except Exception as e:
            print(f"❌ 搜索邮件时出错: {e}")
            return []
    
    def fetch_email_headers(self, email_id):
        """获取邮件头部信息"""
        if not self.imap:
            return None
            
        try:
            result, data = self.imap.fetch(email_id, '(RFC822.HEADER)')
            
            if result == 'OK' and data[0]:
                header = data[0][1].decode('utf-8', errors='ignore')
                return header
            else:
                return None
                
        except Exception as e:
            print(f"❌ 获取邮件头部时出错: {e}")
            return None
    
    def show_recent_emails(self, count=5):
        """显示最近的邮件"""
        email_ids = self.search_emails('ALL')
        
        if not email_ids:
            print("📭 没有找到邮件")
            return
        
        print(f"\n📧 显示最近 {min(count, len(email_ids))} 封邮件:")
        
        # 获取最后几封邮件
        recent_ids = email_ids[-count:] if len(email_ids) > count else email_ids
        
        for i, email_id in enumerate(reversed(recent_ids), 1):
            try:
                result, data = self.imap.fetch(email_id, '(RFC822.HEADER)')
                if result == 'OK' and data[0]:
                    header = data[0][1].decode('utf-8', errors='ignore')
                    
                    # 提取主题和发件人
                    subject = "无主题"
                    sender = "未知发件人"
                    date = "未知日期"
                    
                    for line in header.split('\n'):
                        if line.lower().startswith('subject:'):
                            subject = line[8:].strip()
                        elif line.lower().startswith('from:'):
                            sender = line[5:].strip()
                        elif line.lower().startswith('date:'):
                            date = line[5:].strip()
                    
                    print(f"   {i}. ID: {email_id.decode()}")
                    print(f"      主题: {subject}")
                    print(f"      发件人: {sender}")
                    print(f"      日期: {date}")
                    print()
                    
            except Exception as e:
                print(f"   ❌ 获取邮件 {email_id} 失败: {e}")
    
    def logout(self):
        """登出并关闭连接"""
        if self.imap:
            try:
                print("\n👋 正在登出...")
                self.imap.logout()
                print("✅ 已安全登出")
            except Exception as e:
                print(f"⚠️  登出时出现警告: {e}")
            finally:
                self.imap = None

def main():
    parser = argparse.ArgumentParser(description='IMAP测试工具')
    parser.add_argument('--host', default='localhost', help='IMAP服务器地址 (默认: localhost)')
    parser.add_argument('--port', type=int, default=143, help='IMAP端口 (默认: 143)')
    parser.add_argument('--ssl', action='store_true', help='使用SSL连接')
    parser.add_argument('--username', help='网站登录账号')
    parser.add_argument('--email', help='域名邮箱')
    parser.add_argument('--password', help='邮箱密码')
    
    args = parser.parse_args()
    
    # 如果没有提供参数，使用默认测试账号
    if not args.username:
        print("📝 使用默认测试账号:")
        args.username = 'kimi11'
        args.email = 'kimi11@jbjj.site'
        args.password = '93921438'
    
    print("=" * 60)
    print("🧪 IMAP测试工具")
    print("=" * 60)
    print(f"服务器: {args.host}:{args.port}")
    print(f"SSL: {'是' if args.ssl else '否'}")
    print(f"测试时间: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
    print("=" * 60)
    
    # 创建测试器
    tester = IMAPTester(args.host, args.port, args.ssl)
    
    try:
        # 连接
        if not tester.connect():
            return 1
        
        # 登录
        if not tester.login(args.username, args.email, args.password):
            return 1
        
        # 列出文件夹
        tester.list_folders()
        
        # 选择收件箱
        if tester.select_inbox():
            # 显示最近邮件
            tester.show_recent_emails(5)
        
        print("\n" + "=" * 60)
        print("✅ IMAP测试完成!")
        print("=" * 60)
        
    except KeyboardInterrupt:
        print("\n⚠️  用户中断测试")
    except Exception as e:
        print(f"\n❌ 测试过程中出现错误: {e}")
        return 1
    finally:
        # 登出
        tester.logout()
    
    return 0

if __name__ == '__main__':
    sys.exit(main())
