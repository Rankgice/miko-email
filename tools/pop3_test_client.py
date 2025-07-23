#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
POP3测试客户端
"""

import socket
import logging
import time

# 配置日志
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)

class POP3TestClient:
    def __init__(self, host="localhost", port=110):
        self.host = host
        self.port = port
        self.socket = None
    
    def connect(self):
        """连接到POP3服务器"""
        try:
            self.socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            self.socket.settimeout(30)
            self.socket.connect((self.host, self.port))
            
            # 读取欢迎消息
            response = self.read_response()
            logging.info(f"服务器欢迎消息: {response}")
            
            if response.startswith("+OK"):
                return True
            else:
                logging.error(f"连接失败: {response}")
                return False
                
        except Exception as e:
            logging.error(f"连接失败: {str(e)}")
            return False
    
    def send_command(self, command):
        """发送POP3命令"""
        try:
            logging.info(f"发送命令: {command}")
            self.socket.send((command + "\r\n").encode('utf-8'))
            return True
        except Exception as e:
            logging.error(f"发送命令失败: {str(e)}")
            return False
    
    def read_response(self):
        """读取POP3响应"""
        try:
            response = self.socket.recv(1024).decode('utf-8').strip()
            logging.info(f"服务器响应: {response}")
            return response
        except Exception as e:
            logging.error(f"读取响应失败: {str(e)}")
            return ""
    
    def read_multiline_response(self):
        """读取多行响应（以.结束）"""
        try:
            lines = []
            while True:
                line = self.socket.recv(1024).decode('utf-8')
                if not line:
                    break
                
                lines.append(line)
                if line.strip().endswith('.'):
                    break
            
            response = ''.join(lines)
            logging.info(f"多行响应: {response[:200]}...")
            return response
        except Exception as e:
            logging.error(f"读取多行响应失败: {str(e)}")
            return ""
    
    def login(self, username, password):
        """登录POP3服务器"""
        # 发送USER命令
        if not self.send_command(f"USER {username}"):
            return False
        
        response = self.read_response()
        if not response.startswith("+OK"):
            logging.error(f"USER命令失败: {response}")
            return False
        
        # 发送PASS命令
        if not self.send_command(f"PASS {password}"):
            return False
        
        response = self.read_response()
        if response.startswith("+OK"):
            logging.info("✅ POP3登录成功")
            return True
        else:
            logging.error(f"❌ POP3登录失败: {response}")
            return False
    
    def stat(self):
        """获取邮箱统计信息"""
        if not self.send_command("STAT"):
            return None
        
        response = self.read_response()
        if response.startswith("+OK"):
            parts = response.split()
            if len(parts) >= 3:
                count = int(parts[1])
                size = int(parts[2])
                logging.info(f"📊 邮箱统计: {count} 封邮件, 总大小 {size} 字节")
                return count, size
        
        logging.error(f"STAT命令失败: {response}")
        return None
    
    def list_messages(self):
        """列出所有邮件"""
        if not self.send_command("LIST"):
            return []
        
        response = self.read_response()
        if not response.startswith("+OK"):
            logging.error(f"LIST命令失败: {response}")
            return []
        
        # 读取完整的多行响应
        all_data = b""
        while True:
            data = self.socket.recv(1024)
            if not data:
                break
            all_data += data

            # 检查是否收到完整响应（以.\r\n结束）
            if b".\r\n" in all_data:
                break

        response = all_data.decode('utf-8')
        lines = response.split('\r\n')

        messages = []
        for line in lines[1:]:  # 跳过第一行的+OK
            if line == "." or line == "":
                break
            if line:
                parts = line.split()
                if len(parts) >= 2:
                    msg_num = int(parts[0])
                    msg_size = int(parts[1])
                    messages.append((msg_num, msg_size))
                    logging.info(f"📧 邮件 {msg_num}: {msg_size} 字节")

        return messages
    
    def retrieve_message(self, msg_num):
        """获取指定邮件"""
        if not self.send_command(f"RETR {msg_num}"):
            return None
        
        response = self.read_response()
        if not response.startswith("+OK"):
            logging.error(f"RETR命令失败: {response}")
            return None
        
        # 读取邮件内容
        content = []
        while True:
            line = self.socket.recv(1024).decode('utf-8')
            if not line:
                break
            content.append(line)
            if line.strip() == ".":
                break
        
        email_content = ''.join(content)
        logging.info(f"📬 获取邮件 {msg_num} 成功，大小: {len(email_content)} 字节")
        return email_content
    
    def delete_message(self, msg_num):
        """删除指定邮件"""
        if not self.send_command(f"DELE {msg_num}"):
            return False
        
        response = self.read_response()
        if response.startswith("+OK"):
            logging.info(f"🗑️ 邮件 {msg_num} 已标记为删除")
            return True
        else:
            logging.error(f"删除邮件失败: {response}")
            return False
    
    def quit(self):
        """退出POP3会话"""
        if not self.send_command("QUIT"):
            return False
        
        response = self.read_response()
        if response.startswith("+OK"):
            logging.info("✅ POP3会话结束")
            return True
        else:
            logging.error(f"QUIT命令失败: {response}")
            return False
    
    def close(self):
        """关闭连接"""
        if self.socket:
            self.socket.close()

def test_pop3_server():
    """测试POP3服务器"""
    logging.info("=" * 60)
    logging.info("🧪 POP3服务器测试")
    logging.info("=" * 60)
    
    # 测试用户信息
    username = "kimi11@jbjj.site"
    password = "93921438"
    
    client = POP3TestClient()
    
    try:
        # 1. 连接服务器
        if not client.connect():
            logging.error("❌ 无法连接到POP3服务器")
            return
        
        # 2. 登录
        if not client.login(username, password):
            logging.error("❌ POP3登录失败")
            return
        
        # 3. 获取邮箱统计
        stat_result = client.stat()
        if not stat_result:
            logging.error("❌ 获取邮箱统计失败")
            return
        
        count, total_size = stat_result
        
        # 4. 列出邮件
        messages = client.list_messages()
        logging.info(f"📋 找到 {len(messages)} 封邮件")
        
        # 5. 获取前几封邮件的内容
        for i, (msg_num, msg_size) in enumerate(messages[:3]):  # 只获取前3封
            logging.info(f"📖 正在获取邮件 {msg_num}...")
            content = client.retrieve_message(msg_num)
            if content:
                # 显示邮件头部信息
                lines = content.split('\n')
                for line in lines[:10]:  # 显示前10行
                    if line.strip():
                        logging.info(f"  {line.strip()}")
                logging.info("  ...")
        
        # 6. 测试删除功能（可选）
        # if messages:
        #     logging.info(f"🗑️ 测试删除邮件 {messages[-1][0]}...")
        #     client.delete_message(messages[-1][0])
        
        # 7. 退出
        client.quit()
        
        logging.info("=" * 60)
        logging.info("🎉 POP3测试完成！")
        logging.info("=" * 60)
        
    except Exception as e:
        logging.error(f"❌ 测试过程中出错: {str(e)}")
        import traceback
        logging.error(traceback.format_exc())
    
    finally:
        client.close()

if __name__ == "__main__":
    test_pop3_server()
