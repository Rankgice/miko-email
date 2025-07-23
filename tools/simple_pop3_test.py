#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
简化的POP3测试
"""

import socket
import logging

# 配置日志
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)

def simple_pop3_test():
    """简化的POP3测试"""
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
        
        # 发送STAT命令
        sock.send("STAT\r\n".encode('utf-8'))
        response = sock.recv(1024).decode('utf-8').strip()
        logging.info(f"STAT响应: {response}")
        
        # 发送LIST命令
        logging.info("发送LIST命令...")
        sock.send("LIST\r\n".encode('utf-8'))
        
        # 读取LIST响应
        logging.info("读取LIST响应...")
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
                logging.error("读取LIST响应超时")
                break
        
        response = all_data.decode('utf-8')
        logging.info(f"LIST完整响应: {repr(response)}")
        
        # 解析响应
        lines = response.split('\r\n')
        for i, line in enumerate(lines):
            logging.info(f"第{i+1}行: {repr(line)}")
        
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
    simple_pop3_test()
