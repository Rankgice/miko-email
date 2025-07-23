#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
调试IMAP服务器响应
"""

import socket
import logging

# 配置日志
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)

def debug_imap_server():
    """调试IMAP服务器响应"""
    try:
        # 连接到IMAP服务器
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.settimeout(10)
        sock.connect(("localhost", 143))
        
        # 读取欢迎消息
        welcome = sock.recv(1024).decode('utf-8').strip()
        logging.info(f"欢迎消息: {repr(welcome)}")
        
        # 登录 - 使用邮箱地址作为用户名
        login_cmd = "A001 LOGIN kimi11@jbjj.site 93921438\r\n"
        logging.info(f"发送登录命令: {repr(login_cmd)}")
        sock.send(login_cmd.encode('utf-8'))
        
        login_response = sock.recv(1024).decode('utf-8').strip()
        logging.info(f"登录响应: {repr(login_response)}")
        
        # 选择收件箱
        select_cmd = "A002 SELECT INBOX\r\n"
        logging.info(f"发送SELECT命令: {repr(select_cmd)}")
        sock.send(select_cmd.encode('utf-8'))
        
        select_response = sock.recv(1024).decode('utf-8').strip()
        logging.info(f"SELECT响应: {repr(select_response)}")
        
        # 搜索邮件
        search_cmd = "A003 SEARCH ALL\r\n"
        logging.info(f"发送SEARCH命令: {repr(search_cmd)}")
        sock.send(search_cmd.encode('utf-8'))
        
        search_response = sock.recv(1024).decode('utf-8').strip()
        logging.info(f"SEARCH响应: {repr(search_response)}")
        
        # 分析响应
        logging.info("=" * 50)
        logging.info("响应分析:")
        logging.info(f"响应长度: {len(search_response)}")
        logging.info(f"响应行数: {len(search_response.split('\\n'))}")
        
        lines = search_response.split('\n')
        for i, line in enumerate(lines):
            logging.info(f"第{i+1}行: {repr(line)}")
            if line.startswith('*') and 'SEARCH' in line:
                logging.info(f"  -> 找到SEARCH行: {line}")
                parts = line.split()
                logging.info(f"  -> 分割后: {parts}")
                if len(parts) > 2:
                    email_ids = parts[2:]
                    logging.info(f"  -> 邮件ID: {email_ids}")
        
        # 关闭连接
        sock.close()
        
    except Exception as e:
        logging.error(f"调试失败: {str(e)}")
        import traceback
        logging.error(traceback.format_exc())

if __name__ == "__main__":
    debug_imap_server()
