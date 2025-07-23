#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
简化的POP3 RETR测试
"""

import socket
import logging

# 配置日志
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)

def simple_retr_test():
    """简化的POP3 RETR测试"""
    try:
        # 连接到POP3服务器
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.settimeout(10)
        sock.connect(("localhost", 110))
        
        # 读取欢迎消息
        welcome = sock.recv(1024).decode('utf-8').strip()
        logging.info(f"欢迎消息: {welcome}")
        
        # 登录
        sock.send("USER kimi11@jbjj.site\r\n".encode('utf-8'))
        response = sock.recv(1024).decode('utf-8').strip()
        logging.info(f"USER响应: {response}")
        
        sock.send("PASS 93921438\r\n".encode('utf-8'))
        response = sock.recv(1024).decode('utf-8').strip()
        logging.info(f"PASS响应: {response}")
        
        # 发送RETR命令
        logging.info("发送RETR 1命令...")
        sock.send("RETR 1\r\n".encode('utf-8'))
        
        # 读取RETR响应
        logging.info("读取RETR响应...")
        all_data = b""
        while True:
            try:
                data = sock.recv(1024)
                if not data:
                    break
                all_data += data
                
                # 检查是否收到完整响应（以.\r\n结束）
                if b".\r\n" in all_data:
                    break
                    
            except socket.timeout:
                logging.error("读取RETR响应超时")
                break
        
        response = all_data.decode('utf-8')
        logging.info(f"RETR完整响应长度: {len(response)}")
        logging.info(f"RETR响应前200字符: {repr(response[:200])}")
        
        # 发送QUIT命令
        sock.send("QUIT\r\n".encode('utf-8'))
        response = sock.recv(1024).decode('utf-8').strip()
        logging.info(f"QUIT响应: {response}")
        
        sock.close()
        
    except Exception as e:
        logging.error(f"测试失败: {str(e)}")
        import traceback
        logging.error(traceback.format_exc())

if __name__ == "__main__":
    simple_retr_test()
