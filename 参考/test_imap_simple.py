#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
简单的IMAP测试脚本
用于调试IMAP协议问题
"""

import imaplib
import sys

def test_imap_connection():
    """测试IMAP连接和基本操作"""
    try:
        # 连接到本地IMAP服务器
        print("连接到IMAP服务器...")
        mail = imaplib.IMAP4("127.0.0.1", 143)
        
        # 启用调试模式
        mail.debug = 4
        
        print("尝试登录...")
        mail.login("2014131458@qq.com", "tgx123456")
        print("登录成功")
        
        print("选择收件箱...")
        mail.select('INBOX')
        
        print("搜索邮件...")
        status, messages = mail.search(None, 'ALL')
        print(f"搜索状态: {status}")
        print(f"邮件列表: {messages}")
        
        if status == 'OK' and messages[0]:
            email_ids = messages[0].split()
            print(f"找到 {len(email_ids)} 封邮件")
            
            if email_ids:
                # 只测试第一封邮件
                first_email_id = email_ids[0]
                print(f"获取邮件 ID: {first_email_id}")
                
                # 尝试获取邮件
                status, data = mail.fetch(first_email_id, '(RFC822)')
                print(f"FETCH状态: {status}")
                print(f"FETCH数据类型: {type(data)}")
                print(f"FETCH数据长度: {len(data) if data else 0}")
                
                if data:
                    for i, item in enumerate(data):
                        print(f"数据项 {i}: {type(item)} - {len(str(item)) if item else 0} 字符")
                        if isinstance(item, tuple) and len(item) > 1:
                            print(f"  元组第一项: {type(item[0])}")
                            print(f"  元组第二项: {type(item[1])}")
        
        print("登出...")
        mail.logout()
        print("测试完成")
        
    except Exception as e:
        print(f"测试失败: {e}")
        import traceback
        traceback.print_exc()

if __name__ == "__main__":
    test_imap_connection()
