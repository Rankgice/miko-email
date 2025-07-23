#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
测试修复后的IMAP服务器
"""

import socket
import logging
import time

# 配置日志
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)

def send_command_and_read_response(sock, command):
    """发送命令并读取完整响应"""
    logging.info(f"发送命令: {repr(command)}")
    sock.send(command.encode('utf-8'))
    
    # 读取响应，直到找到带标签的响应
    response_lines = []
    while True:
        data = sock.recv(1024).decode('utf-8')
        if not data:
            break
        
        lines = data.split('\r\n')
        for line in lines:
            if line.strip():
                response_lines.append(line)
                logging.info(f"响应行: {repr(line)}")
                
                # 检查是否是带标签的最终响应
                if line.startswith(command.split()[0]):
                    return '\r\n'.join(response_lines)
        
        # 短暂等待，避免过快读取
        time.sleep(0.1)
    
    return '\r\n'.join(response_lines)

def test_imap_server():
    """测试IMAP服务器"""
    try:
        logging.info("=" * 60)
        logging.info("测试修复后的IMAP服务器")
        logging.info("=" * 60)
        
        # 连接到IMAP服务器
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.settimeout(10)
        sock.connect(("localhost", 143))
        
        # 读取欢迎消息
        welcome = sock.recv(1024).decode('utf-8').strip()
        logging.info(f"欢迎消息: {repr(welcome)}")
        
        # 1. 登录
        login_response = send_command_and_read_response(sock, "A001 LOGIN kimi11@jbjj.site 93921438\r\n")
        if "OK LOGIN completed" in login_response:
            logging.info("✅ 登录成功")
        else:
            logging.error("❌ 登录失败")
            return
        
        # 2. 选择收件箱
        select_response = send_command_and_read_response(sock, "A002 SELECT INBOX\r\n")
        if "OK" in select_response and "SELECT completed" in select_response:
            logging.info("✅ 选择收件箱成功")
            # 提取邮件数量
            for line in select_response.split('\r\n'):
                if "EXISTS" in line:
                    email_count = line.split()[1]
                    logging.info(f"📧 收件箱中有 {email_count} 封邮件")
        else:
            logging.error("❌ 选择收件箱失败")
            return
        
        # 3. 搜索所有邮件
        search_response = send_command_and_read_response(sock, "A003 SEARCH ALL\r\n")
        if "OK SEARCH completed" in search_response:
            logging.info("✅ 搜索命令成功")
            # 提取邮件ID
            for line in search_response.split('\r\n'):
                if line.startswith("* SEARCH"):
                    if len(line.split()) > 2:
                        email_ids = line.split()[2:]
                        logging.info(f"📬 找到邮件ID: {email_ids}")
                    else:
                        logging.info("📭 没有找到邮件")
        else:
            logging.error("❌ 搜索命令失败")
            return
        
        # 4. 获取第一封邮件的信息
        fetch_response = send_command_and_read_response(sock, "A004 FETCH 1 (ENVELOPE)\r\n")
        if "OK FETCH completed" in fetch_response:
            logging.info("✅ 获取邮件信息成功")
            logging.info(f"📧 邮件信息: {fetch_response[:200]}...")
        else:
            logging.info("ℹ️ 获取邮件信息失败（可能没有邮件）")
        
        # 5. 登出
        logout_response = send_command_and_read_response(sock, "A005 LOGOUT\r\n")
        if "OK LOGOUT completed" in logout_response:
            logging.info("✅ 登出成功")
        
        # 关闭连接
        sock.close()
        
        logging.info("=" * 60)
        logging.info("🎉 IMAP服务器测试完成！")
        logging.info("=" * 60)
        
    except Exception as e:
        logging.error(f"测试失败: {str(e)}")
        import traceback
        logging.error(traceback.format_exc())

if __name__ == "__main__":
    test_imap_server()
