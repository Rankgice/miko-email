#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
IMAP连接测试脚本
用于测试NBEmail系统的IMAP功能
"""

import imaplib
import logging
import sys

# 配置日志
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s',
    handlers=[
        logging.StreamHandler(sys.stdout)
    ]
)

def test_imap_connection():
    """测试IMAP连接"""
    
    # 服务器配置
    imap_server = "me.youddns.site"
    imap_port = 143
    
    # 用户凭据
    email = "2014131458@qq.com"
    password = "tgx123456"
    
    logging.info(f"开始测试IMAP连接...")
    logging.info(f"服务器: {imap_server}:{imap_port}")
    logging.info(f"用户: {email}")
    
    try:
        # 连接到IMAP服务器
        logging.info("正在连接到IMAP服务器...")
        mail = imaplib.IMAP4(imap_server, imap_port)
        logging.info("✅ 成功连接到IMAP服务器")
        
        # 尝试登录
        logging.info("正在尝试登录...")
        mail.login(email, password)
        logging.info("✅ 登录成功")
        
        # 列出邮箱
        logging.info("正在获取邮箱列表...")
        status, mailboxes = mail.list()
        if status == 'OK':
            logging.info("✅ 邮箱列表获取成功:")
            for mailbox in mailboxes:
                logging.info(f"  - {mailbox.decode('utf-8')}")
        else:
            logging.warning("⚠️ 无法获取邮箱列表")
        
        # 选择收件箱
        logging.info("正在选择INBOX...")
        status, messages = mail.select('INBOX')
        if status == 'OK':
            message_count = int(messages[0])
            logging.info(f"✅ 成功选择INBOX，邮件数量: {message_count}")
        else:
            logging.warning("⚠️ 无法选择INBOX")
        
        # 搜索邮件
        if status == 'OK':
            logging.info("正在搜索邮件...")
            status, email_ids = mail.search(None, 'ALL')
            if status == 'OK':
                ids = email_ids[0].split()
                logging.info(f"✅ 找到 {len(ids)} 封邮件")
                
                # 获取最新的一封邮件信息
                if ids:
                    latest_id = ids[-1]
                    logging.info(f"正在获取最新邮件 (ID: {latest_id.decode()})...")
                    status, data = mail.fetch(latest_id, '(ENVELOPE)')
                    if status == 'OK':
                        logging.info("✅ 成功获取邮件信息")
                        logging.info(f"邮件数据: {data[0]}")
                    else:
                        logging.warning("⚠️ 无法获取邮件信息")
            else:
                logging.warning("⚠️ 无法搜索邮件")
        
        # 登出
        logging.info("正在登出...")
        mail.logout()
        logging.info("✅ 成功登出")
        
        logging.info("🎉 IMAP测试完成！所有功能正常")
        return True
        
    except imaplib.IMAP4.error as e:
        logging.error(f"❌ IMAP协议错误: {e}")
        return False
    except ConnectionRefusedError:
        logging.error(f"❌ 连接被拒绝，请检查服务器是否运行在 {imap_server}:{imap_port}")
        return False
    except Exception as e:
        logging.error(f"❌ 连接失败: {e}")
        return False

if __name__ == "__main__":
    print("=" * 60)
    print("NBEmail IMAP连接测试")
    print("=" * 60)
    
    success = test_imap_connection()
    
    print("=" * 60)
    if success:
        print("✅ 测试成功！IMAP服务器工作正常")
    else:
        print("❌ 测试失败！请检查服务器配置")
    print("=" * 60)
