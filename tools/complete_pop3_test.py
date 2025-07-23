#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
完整的POP3功能测试
"""

import socket
import logging

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

def complete_pop3_test():
    """完整的POP3功能测试"""
    logging.info("=" * 60)
    logging.info("🧪 完整POP3功能测试")
    logging.info("=" * 60)
    
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
        
        logging.info("✅ 连接成功")
        
        # 1. 测试USER命令
        send_command(sock, "USER kimi11@jbjj.site")
        response = read_response(sock)
        if not response.startswith("+OK"):
            logging.error("❌ USER命令失败")
            return
        logging.info("✅ USER命令成功")
        
        # 2. 测试PASS命令
        send_command(sock, "PASS 93921438")
        response = read_response(sock)
        if not response.startswith("+OK"):
            logging.error("❌ PASS命令失败")
            return
        logging.info("✅ PASS命令成功，登录完成")
        
        # 3. 测试STAT命令
        send_command(sock, "STAT")
        response = read_response(sock)
        if response.startswith("+OK"):
            parts = response.split()
            if len(parts) >= 3:
                count = int(parts[1])
                size = int(parts[2])
                logging.info(f"✅ STAT命令成功: {count} 封邮件, 总大小 {size} 字节")
            else:
                logging.error("❌ STAT响应格式错误")
                return
        else:
            logging.error("❌ STAT命令失败")
            return
        
        # 4. 测试LIST命令
        send_command(sock, "LIST")
        response = read_multiline_response(sock)
        lines = response.split('\r\n')
        
        if lines[0].startswith("+OK"):
            messages = []
            for line in lines[1:]:
                if line == "." or line == "":
                    break
                if line:
                    parts = line.split()
                    if len(parts) >= 2:
                        msg_num = int(parts[0])
                        msg_size = int(parts[1])
                        messages.append((msg_num, msg_size))
            
            logging.info(f"✅ LIST命令成功: 找到 {len(messages)} 封邮件")
            for i, (num, size) in enumerate(messages[:5]):  # 只显示前5封
                logging.info(f"  邮件 {num}: {size} 字节")
        else:
            logging.error("❌ LIST命令失败")
            return
        
        # 5. 测试RETR命令（获取第一封邮件）
        if messages:
            first_msg = messages[0][0]
            send_command(sock, f"RETR {first_msg}")
            response = read_multiline_response(sock)
            
            if response.startswith("+OK"):
                lines = response.split('\r\n')
                logging.info(f"✅ RETR命令成功: 获取邮件 {first_msg}")
                logging.info("  邮件头部:")
                for line in lines[1:6]:  # 显示前5行
                    if line and line != ".":
                        logging.info(f"    {line}")
                logging.info("  ...")
            else:
                logging.error("❌ RETR命令失败")
        
        # 6. 测试UIDL命令
        send_command(sock, "UIDL")
        response = read_multiline_response(sock)
        
        if response.startswith("+OK"):
            logging.info("✅ UIDL命令成功")
        else:
            logging.error("❌ UIDL命令失败")
        
        # 7. 测试NOOP命令
        send_command(sock, "NOOP")
        response = read_response(sock)
        if response.startswith("+OK"):
            logging.info("✅ NOOP命令成功")
        else:
            logging.error("❌ NOOP命令失败")
        
        # 8. 测试TOP命令（获取邮件头部）
        if messages:
            first_msg = messages[0][0]
            send_command(sock, f"TOP {first_msg} 5")
            response = read_multiline_response(sock)
            
            if response.startswith("+OK"):
                logging.info(f"✅ TOP命令成功: 获取邮件 {first_msg} 的头部和前5行")
            else:
                logging.error("❌ TOP命令失败")
        
        # 9. 测试DELE命令（删除最后一封邮件）
        # 注意：这会真的删除邮件，所以要小心
        # if messages:
        #     last_msg = messages[-1][0]
        #     send_command(sock, f"DELE {last_msg}")
        #     response = read_response(sock)
        #     if response.startswith("+OK"):
        #         logging.info(f"✅ DELE命令成功: 标记删除邮件 {last_msg}")
        #     else:
        #         logging.error("❌ DELE命令失败")
        
        # 10. 测试RSET命令
        send_command(sock, "RSET")
        response = read_response(sock)
        if response.startswith("+OK"):
            logging.info("✅ RSET命令成功")
        else:
            logging.error("❌ RSET命令失败")
        
        # 11. 测试QUIT命令
        send_command(sock, "QUIT")
        response = read_response(sock)
        if response.startswith("+OK"):
            logging.info("✅ QUIT命令成功")
        else:
            logging.error("❌ QUIT命令失败")
        
        sock.close()
        
        logging.info("=" * 60)
        logging.info("🎉 POP3功能测试完成！")
        logging.info("=" * 60)
        
    except Exception as e:
        logging.error(f"❌ 测试失败: {str(e)}")
        import traceback
        logging.error(traceback.format_exc())

if __name__ == "__main__":
    complete_pop3_test()
