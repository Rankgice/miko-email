#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
POP3完整邮件内容测试 - 查看指定邮件的完整内容
"""

import socket
import logging
import re
import html

# 配置日志
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)

def send_command(sock, command):
    """发送POP3命令"""
    sock.send((command + "\r\n").encode('utf-8'))

def read_response(sock):
    """读取单行响应"""
    response = sock.recv(1024).decode('utf-8').strip()
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

def extract_text_from_html(html_content):
    """从HTML中提取纯文本"""
    # 简单的HTML标签移除
    text = re.sub(r'<[^>]+>', '', html_content)
    # 解码HTML实体
    text = html.unescape(text)
    # 清理多余的空白
    text = re.sub(r'\s+', ' ', text).strip()
    return text

def pop3_full_email_test(email_num=15):
    """查看指定邮件的完整内容"""
    print("=" * 100)
    print(f"📧 POP3完整邮件内容测试 - 邮件 {email_num}")
    print("=" * 100)
    
    try:
        # 连接到POP3服务器
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.settimeout(30)
        sock.connect(("localhost", 110))
        
        # 读取欢迎消息
        welcome = read_response(sock)
        
        # 登录
        send_command(sock, "USER kimi11@jbjj.site")
        read_response(sock)
        
        send_command(sock, "PASS 93921438")
        read_response(sock)
        
        print("✅ 登录成功")
        
        # 获取指定邮件的完整内容
        print(f"\n📖 正在获取邮件 {email_num} 的完整内容...")
        
        send_command(sock, f"RETR {email_num}")
        email_content = read_multiline_response(sock)
        
        if not email_content.startswith("+OK"):
            print(f"❌ 获取邮件 {email_num} 失败")
            return
        
        # 解析邮件内容
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
        
        print("\n" + "=" * 100)
        print("📋 邮件头信息:")
        print("=" * 100)
        
        for key, value in headers.items():
            print(f"{key}: {value}")
        
        print("\n" + "=" * 100)
        print("📄 邮件正文 (原始内容):")
        print("=" * 100)
        print(body)
        
        # 如果是HTML邮件，提取纯文本
        if '<html>' in body.lower() or '<div>' in body.lower():
            print("\n" + "=" * 100)
            print("📄 邮件正文 (纯文本提取):")
            print("=" * 100)
            
            text_content = extract_text_from_html(body)
            print(text_content)
            
            # 查找验证码
            verification_patterns = [
                r'验证码[：:]\s*(\d{4,8})',
                r'验证码[是为]\s*(\d{4,8})',
                r'code[：:]\s*(\d{4,8})',
                r'<div class="code">(\d+)</div>',
                r'(\d{6})',  # 6位数字
            ]
            
            print("\n🔍 验证码搜索:")
            found_codes = set()
            for pattern in verification_patterns:
                matches = re.findall(pattern, body, re.IGNORECASE)
                for match in matches:
                    if match not in found_codes and len(match) >= 4:
                        found_codes.add(match)
                        print(f"  🔑 找到验证码: {match}")
            
            if not found_codes:
                print("  ❌ 未找到验证码")
        
        print("\n" + "=" * 100)
        print("📊 邮件统计信息:")
        print("=" * 100)
        print(f"头部字段数: {len(headers)}")
        print(f"正文长度: {len(body)} 字符")
        print(f"总大小: {len(email_content)} 字节")
        print(f"行数: {len(lines)}")
        
        # 退出
        send_command(sock, "QUIT")
        read_response(sock)
        sock.close()
        
        print("\n🎉 邮件内容查看完成！")
        
    except Exception as e:
        print(f"❌ 测试失败: {str(e)}")
        import traceback
        print(traceback.format_exc())

if __name__ == "__main__":
    # 查看邮件15（YouDDNS验证码邮件）
    pop3_full_email_test(15)
