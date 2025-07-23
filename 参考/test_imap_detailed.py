#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
详细的IMAP测试脚本
用于调试SEARCH命令的具体响应
"""

import imaplib
import logging
import sys

# 配置日志
logging.basicConfig(
    level=logging.DEBUG,
    format='%(asctime)s - %(levelname)s - %(message)s',
    handlers=[
        logging.StreamHandler(sys.stdout)
    ]
)

def test_imap_search():
    """详细测试IMAP SEARCH命令"""
    
    # 服务器配置
    imap_server = "me.youddns.site"
    imap_port = 143
    
    # 用户凭据
    email = "2014131458@qq.com"
    password = "tgx123456"
    
    logging.info(f"开始详细测试IMAP SEARCH...")
    logging.info(f"服务器: {imap_server}:{imap_port}")
    logging.info(f"用户: {email}")
    
    try:
        # 连接到IMAP服务器
        logging.info("正在连接到IMAP服务器...")
        mail = imaplib.IMAP4(imap_server, imap_port)
        logging.info("✅ 成功连接到IMAP服务器")
        
        # 启用调试模式
        mail.debug = 4
        
        # 尝试登录
        logging.info("正在尝试登录...")
        mail.login(email, password)
        logging.info("✅ 登录成功")
        
        # 选择收件箱
        logging.info("正在选择INBOX...")
        status, messages = mail.select('INBOX')
        if status == 'OK':
            message_count = int(messages[0])
            logging.info(f"✅ 成功选择INBOX，邮件数量: {message_count}")
        else:
            logging.warning("⚠️ 无法选择INBOX")
            return False
        
        # 测试SEARCH命令
        logging.info("正在测试SEARCH ALL命令...")
        try:
            status, email_ids = mail.search(None, 'ALL')
            logging.info(f"SEARCH状态: {status}")
            logging.info(f"SEARCH结果: {email_ids}")
            
            if status == 'OK':
                ids = email_ids[0].split() if email_ids[0] else []
                logging.info(f"✅ SEARCH成功，找到 {len(ids)} 封邮件")
                if ids:
                    logging.info(f"邮件ID列表: {[id.decode() for id in ids]}")
                else:
                    logging.info("📭 收件箱为空")
            else:
                logging.error(f"❌ SEARCH失败: {status}")
                return False
                
        except Exception as e:
            logging.error(f"❌ SEARCH命令异常: {e}")
            return False
        
        # 登出
        logging.info("正在登出...")
        mail.logout()
        logging.info("✅ 成功登出")
        
        logging.info("🎉 详细IMAP测试完成！")
        return True
        
    except Exception as e:
        logging.error(f"❌ 测试失败: {e}")
        return False

if __name__ == "__main__":
    print("=" * 60)
    print("NBEmail IMAP详细测试 - SEARCH命令调试")
    print("=" * 60)
    
    success = test_imap_search()
    
    print("=" * 60)
    if success:
        print("✅ 详细测试成功！")
    else:
        print("❌ 详细测试失败！")
    print("=" * 60)
