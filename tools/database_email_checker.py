#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
直接从数据库检查邮件的工具
"""

import sqlite3
import logging
import time
from datetime import datetime

# 配置日志
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)

class DatabaseEmailChecker:
    def __init__(self, email_address):
        """初始化数据库邮件检查器"""
        self.email_address = email_address
        self.db_path = "../miko_email.db"
        self.processed_email_ids = set()
        self.running = True
        self.poll_interval = 10
        
        logging.info("=" * 60)
        logging.info("数据库邮件检查器")
        logging.info("=" * 60)
        logging.info(f"邮箱地址: {self.email_address}")
        logging.info(f"数据库路径: {self.db_path}")
        logging.info("=" * 60)
    
    def get_mailbox_id(self):
        """获取邮箱ID"""
        try:
            conn = sqlite3.connect(self.db_path)
            cursor = conn.cursor()
            
            cursor.execute("SELECT id FROM mailboxes WHERE email = ?", (self.email_address,))
            result = cursor.fetchone()
            
            conn.close()
            
            if result:
                return result[0]
            else:
                logging.error(f"未找到邮箱: {self.email_address}")
                return None
                
        except Exception as e:
            logging.error(f"获取邮箱ID失败: {str(e)}")
            return None
    
    def check_new_emails(self, mailbox_id):
        """检查新邮件"""
        try:
            conn = sqlite3.connect(self.db_path)
            cursor = conn.cursor()
            
            # 查询收件箱中的所有邮件
            cursor.execute("""
                SELECT id, from_addr, to_addr, subject, body, created_at, is_read
                FROM emails 
                WHERE mailbox_id = ? AND folder = 'inbox'
                ORDER BY created_at DESC
            """, (mailbox_id,))
            
            emails = cursor.fetchall()
            conn.close()
            
            # 找出新邮件
            new_emails = []
            for email in emails:
                email_id = email[0]
                if email_id not in self.processed_email_ids:
                    new_emails.append(email)
                    self.processed_email_ids.add(email_id)
            
            return new_emails
            
        except Exception as e:
            logging.error(f"检查新邮件失败: {str(e)}")
            return []
    
    def display_email(self, email):
        """显示邮件信息"""
        email_id, from_addr, to_addr, subject, body, created_at, is_read = email
        
        logging.info("=" * 50)
        logging.info(f"📧 新邮件 (ID: {email_id})")
        logging.info("=" * 50)
        logging.info(f"发件人: {from_addr}")
        logging.info(f"收件人: {to_addr}")
        logging.info(f"主题: {subject}")
        logging.info(f"时间: {created_at}")
        logging.info(f"已读: {'是' if is_read else '否'}")
        logging.info("-" * 50)
        logging.info("邮件内容:")
        logging.info(body[:200] + "..." if len(body) > 200 else body)
        logging.info("=" * 50)
    
    def mark_email_as_read(self, email_id):
        """标记邮件为已读"""
        try:
            conn = sqlite3.connect(self.db_path)
            cursor = conn.cursor()
            
            cursor.execute("""
                UPDATE emails 
                SET is_read = 1, updated_at = CURRENT_TIMESTAMP 
                WHERE id = ?
            """, (email_id,))
            
            conn.commit()
            conn.close()
            
            logging.info(f"邮件 {email_id} 已标记为已读")
            
        except Exception as e:
            logging.error(f"标记邮件为已读失败: {str(e)}")
    
    def run_polling(self):
        """运行邮件轮询"""
        # 获取邮箱ID
        mailbox_id = self.get_mailbox_id()
        if not mailbox_id:
            logging.error("无法获取邮箱ID，退出")
            return
        
        logging.info(f"邮箱ID: {mailbox_id}")
        logging.info(f"开始邮件轮询，每 {self.poll_interval} 秒检查一次...")
        
        try:
            # 首次获取所有邮件，标记为已处理（避免重复显示历史邮件）
            initial_emails = self.check_new_emails(mailbox_id)
            logging.info(f"初始化：找到 {len(initial_emails)} 封历史邮件，已标记为已处理")
            
            while self.running:
                try:
                    # 检查新邮件
                    new_emails = self.check_new_emails(mailbox_id)
                    
                    if new_emails:
                        logging.info(f"🎉 发现 {len(new_emails)} 封新邮件！")
                        for email in new_emails:
                            self.display_email(email)
                            # 自动标记为已读
                            self.mark_email_as_read(email[0])
                    else:
                        logging.info("📭 没有新邮件")
                    
                    # 等待下一次轮询
                    logging.info(f"等待 {self.poll_interval} 秒后再次检查...")
                    time.sleep(self.poll_interval)
                    
                except Exception as e:
                    logging.error(f"轮询过程中出错: {str(e)}")
                    time.sleep(self.poll_interval)
        
        except KeyboardInterrupt:
            logging.info("用户中断，停止轮询")
        
        finally:
            logging.info("邮件轮询已停止")

def main():
    """主函数"""
    # 从配置文件读取邮箱地址
    email_address = "kimi11@jbjj.site"
    
    checker = DatabaseEmailChecker(email_address)
    checker.run_polling()

if __name__ == "__main__":
    main()
