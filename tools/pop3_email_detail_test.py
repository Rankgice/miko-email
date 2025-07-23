#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
POP3邮件详情测试 - 查看邮件的完整内容
"""

import socket
import logging
import re

# 配置日志
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)

def send_command(sock, command):
    """发送POP3命令"""
    logging.info(f"发送命令: {command}")
    sock.send((command + "\r\n").encode('utf-8'))

def read_response(sock):
    """读取单行响应"""
    response = sock.recv(1024).decode('utf-8').strip()
    logging.info(f"响应: {response}")
    return response

def read_multiline_response(sock):
    """读取多行响应（以.结束）"""
    all_data = b""
    while True:
        data = sock.recv(1024)
        if not data:
            break
        all_data += data
        
        # 检查是否收到完整响应（以.\r\n结束）
        if b".\r\n" in all_data:
            break
    
    response = all_data.decode('utf-8')
    return response

def parse_email_content(email_content):
    """解析邮件内容"""
    lines = email_content.split('\r\n')
    
    # 跳过第一行的+OK响应
    if lines[0].startswith('+OK'):
        lines = lines[1:]
    
    # 移除最后的.标记
    if lines and lines[-1] == '.':
        lines = lines[:-1]
    
    # 分离邮件头和正文
    headers = {}
    body_start = 0
    
    for i, line in enumerate(lines):
        if line == '':  # 空行分隔头部和正文
            body_start = i + 1
            break
        
        # 解析邮件头
        if ':' in line:
            key, value = line.split(':', 1)
            headers[key.strip()] = value.strip()
    
    # 获取邮件正文
    body = '\r\n'.join(lines[body_start:])
    
    return headers, body

def extract_verification_code(text):
    """提取验证码"""
    # 常见的验证码模式
    patterns = [
        r'验证码[：:]\s*(\d{4,8})',
        r'验证码[是为]\s*(\d{4,8})',
        r'code[：:]\s*(\d{4,8})',
        r'(\d{6})',  # 6位数字
        r'(\d{4})',  # 4位数字
    ]
    
    for pattern in patterns:
        matches = re.findall(pattern, text, re.IGNORECASE)
        if matches:
            return matches[0]
    
    return None

def pop3_email_detail_test():
    """POP3邮件详情测试"""
    logging.info("=" * 80)
    logging.info("📧 POP3邮件详情测试")
    logging.info("=" * 80)
    
    try:
        # 连接到POP3服务器
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.settimeout(30)
        sock.connect(("localhost", 110))
        
        # 读取欢迎消息
        welcome = read_response(sock)
        if not welcome.startswith("+OK"):
            logging.error("❌ 连接失败")
            return
        
        # 登录
        send_command(sock, "USER kimi11@jbjj.site")
        response = read_response(sock)
        if not response.startswith("+OK"):
            logging.error("❌ USER命令失败")
            return
        
        send_command(sock, "PASS 93921438")
        response = read_response(sock)
        if not response.startswith("+OK"):
            logging.error("❌ PASS命令失败")
            return
        
        logging.info("✅ 登录成功")
        
        # 获取邮件列表
        send_command(sock, "LIST")
        response = read_multiline_response(sock)
        lines = response.split('\r\n')
        
        messages = []
        if lines[0].startswith("+OK"):
            for line in lines[1:]:
                if line == "." or line == "":
                    break
                if line:
                    parts = line.split()
                    if len(parts) >= 2:
                        msg_num = int(parts[0])
                        msg_size = int(parts[1])
                        messages.append((msg_num, msg_size))
        
        if not messages:
            logging.error("❌ 没有找到邮件")
            return
        
        logging.info(f"📋 找到 {len(messages)} 封邮件")
        
        # 显示邮件列表供用户选择
        print("\n" + "=" * 60)
        print("📬 邮件列表:")
        for i, (num, size) in enumerate(messages[:10]):  # 只显示前10封
            print(f"  {num}. 邮件 {num} ({size} 字节)")
        print("=" * 60)
        
        # 选择要查看的邮件（默认查看最新的几封）
        test_emails = [messages[0][0], messages[-1][0]]  # 第一封和最后一封
        if len(messages) > 2:
            test_emails.append(messages[len(messages)//2][0])  # 中间一封
        
        for msg_num in test_emails[:3]:  # 最多测试3封邮件
            logging.info(f"\n📖 正在获取邮件 {msg_num} 的详细内容...")
            
            # 获取完整邮件内容
            send_command(sock, f"RETR {msg_num}")
            email_content = read_multiline_response(sock)
            
            if not email_content.startswith("+OK"):
                logging.error(f"❌ 获取邮件 {msg_num} 失败")
                continue
            
            # 解析邮件内容
            headers, body = parse_email_content(email_content)
            
            print("\n" + "=" * 80)
            print(f"📧 邮件 {msg_num} 详细信息:")
            print("=" * 80)
            
            # 显示邮件头信息
            print("📋 邮件头信息:")
            important_headers = ['From', 'To', 'Subject', 'Date']
            for header in important_headers:
                if header in headers:
                    print(f"  {header}: {headers[header]}")
            
            # 显示其他头信息
            other_headers = {k: v for k, v in headers.items() if k not in important_headers}
            if other_headers:
                print("  其他头信息:")
                for key, value in other_headers.items():
                    print(f"    {key}: {value}")
            
            print("\n📄 邮件正文:")
            print("-" * 60)
            
            # 显示邮件正文
            if body:
                # 如果正文很长，只显示前500字符
                if len(body) > 500:
                    print(body[:500])
                    print(f"\n... (正文总长度: {len(body)} 字符，已截断显示)")
                else:
                    print(body)
            else:
                print("(无正文内容)")
            
            print("-" * 60)
            
            # 尝试提取验证码
            full_text = ' '.join(headers.values()) + ' ' + body
            verification_code = extract_verification_code(full_text)
            if verification_code:
                print(f"🔑 检测到验证码: {verification_code}")
            
            # 显示邮件统计信息
            print(f"📊 邮件统计:")
            print(f"  - 头部字段数: {len(headers)}")
            print(f"  - 正文长度: {len(body)} 字符")
            print(f"  - 总大小: {len(email_content)} 字节")
            
            print("=" * 80)
        
        # 退出
        send_command(sock, "QUIT")
        response = read_response(sock)
        
        sock.close()
        
        logging.info("\n🎉 邮件详情测试完成！")
        
    except Exception as e:
        logging.error(f"❌ 测试失败: {str(e)}")
        import traceback
        logging.error(traceback.format_exc())

if __name__ == "__main__":
    pop3_email_detail_test()
