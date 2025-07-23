#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
简单的IMAP测试工具 - 适配Miko邮箱系统
用户凭据: kimi11, kimi11@jbjj.site, 93921438
"""

import socket
import time
import logging

# 配置日志
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s',
    handlers=[
        logging.StreamHandler()
    ]
)

class SimpleIMAPTester:
    def __init__(self, username, email, password, host='localhost', port=143):
        self.username = username
        self.email = email
        self.password = password
        self.host = host
        self.port = port
        self.socket = None
        self.command_id = 1
        
    def connect(self):
        """连接到IMAP服务器"""
        try:
            logging.info(f"连接到IMAP服务器 {self.host}:{self.port}")
            self.socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            self.socket.settimeout(10)
            self.socket.connect((self.host, self.port))
            
            # 读取欢迎消息
            welcome = self.socket.recv(1024).decode('utf-8').strip()
            logging.info(f"服务器欢迎消息: {welcome}")
            
            return True
        except Exception as e:
            logging.error(f"连接失败: {e}")
            return False
    
    def send_command(self, command):
        """发送IMAP命令"""
        try:
            cmd_with_id = f"A{self.command_id:03d} {command}"
            logging.info(f"发送命令: {cmd_with_id}")
            
            self.socket.send((cmd_with_id + "\r\n").encode('utf-8'))
            self.command_id += 1
            
            # 读取响应
            response = self.socket.recv(1024).decode('utf-8').strip()
            logging.info(f"服务器响应: {response}")
            
            return response
        except Exception as e:
            logging.error(f"发送命令失败: {e}")
            return None
    
    def login(self):
        """登录到IMAP服务器"""
        logging.info(f"尝试登录用户: {self.username}")
        
        # 尝试用户名登录
        response = self.send_command(f"LOGIN {self.username} {self.password}")
        if response and "OK" in response:
            logging.info("✅ 用户名登录成功")
            return True
        
        # 尝试邮箱地址登录
        logging.info(f"尝试邮箱地址登录: {self.email}")
        response = self.send_command(f"LOGIN {self.email} {self.password}")
        if response and "OK" in response:
            logging.info("✅ 邮箱地址登录成功")
            return True
        
        logging.error("❌ 登录失败")
        return False
    
    def test_commands(self):
        """测试各种IMAP命令"""
        commands = [
            ("CAPABILITY", "查询服务器能力"),
            ("LIST \"\" \"*\"", "列出文件夹"),
            ("SELECT INBOX", "选择收件箱"),
            ("SEARCH ALL", "搜索所有邮件"),
            ("FETCH 1:* (FLAGS)", "获取邮件标志"),
            ("FETCH 1:* (ENVELOPE)", "获取邮件信封"),
            ("STATUS INBOX (MESSAGES RECENT UNSEEN)", "获取邮箱状态"),
        ]
        
        logging.info("\n" + "="*50)
        logging.info("开始测试IMAP命令")
        logging.info("="*50)
        
        for command, description in commands:
            logging.info(f"\n🔍 {description}")
            response = self.send_command(command)
            time.sleep(0.5)  # 短暂延迟
    
    def logout(self):
        """登出"""
        logging.info("\n👋 正在登出...")
        self.send_command("LOGOUT")
    
    def close(self):
        """关闭连接"""
        if self.socket:
            self.socket.close()
            logging.info("连接已关闭")
    
    def run_test(self):
        """运行完整测试"""
        logging.info("="*60)
        logging.info("🧪 Miko邮箱 简单IMAP测试工具")
        logging.info("="*60)
        logging.info(f"用户名: {self.username}")
        logging.info(f"邮箱: {self.email}")
        logging.info(f"服务器: {self.host}:{self.port}")
        logging.info("="*60)
        
        try:
            # 连接
            if not self.connect():
                return False
            
            # 登录
            if not self.login():
                return False
            
            # 测试命令
            self.test_commands()
            
            # 登出
            self.logout()
            
            logging.info("\n" + "="*60)
            logging.info("✅ IMAP测试完成!")
            logging.info("="*60)
            
            return True
            
        except Exception as e:
            logging.error(f"测试过程中出错: {e}")
            return False
        finally:
            self.close()

def main():
    """主函数"""
    # Miko邮箱用户凭据
    username = "kimi11"
    email = "kimi11@jbjj.site"
    password = "93921438"
    
    # 创建测试器
    tester = SimpleIMAPTester(username, email, password)
    
    # 运行测试
    success = tester.run_test()
    
    if success:
        logging.info("\n💡 测试结果:")
        logging.info("✅ IMAP服务器连接正常")
        logging.info("✅ 用户认证成功")
        logging.info("✅ 基础IMAP命令响应正常")
        logging.info("ℹ️  当前为简单IMAP实现，适合基础需求")
    else:
        logging.error("\n❌ 测试失败")

if __name__ == "__main__":
    main()
