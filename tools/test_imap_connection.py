#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
测试IMAP连接 - 适配Miko邮箱系统
用户凭据: kimi11, kimi11@jbjj.site, 93921438
"""

import socket
import time
import sys

def test_imap_connection():
    """测试IMAP连接"""
    host = 'localhost'
    port = 143
    username = 'kimi11'
    email = 'kimi11@jbjj.site'
    password = '93921438'
    
    print("=" * 60)
    print("Miko邮箱 IMAP连接测试")
    print("=" * 60)
    print(f"服务器: {host}:{port}")
    print(f"用户名: {username}")
    print(f"邮箱: {email}")
    print(f"密码: {password}")
    print("=" * 60)
    
    try:
        # 创建socket连接
        print(f"\n1. 正在连接到 {host}:{port}...")
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.settimeout(10)
        sock.connect((host, port))
        print("✅ 连接成功!")
        
        # 读取欢迎消息
        print("\n2. 读取服务器欢迎消息...")
        welcome = sock.recv(1024).decode('utf-8').strip()
        print(f"服务器响应: {welcome}")
        
        # 发送LOGIN命令
        print(f"\n3. 尝试登录用户: {username}")
        login_cmd = f"A001 LOGIN {username} {password}\r\n"
        print(f"发送命令: {login_cmd.strip()}")
        sock.send(login_cmd.encode('utf-8'))
        
        # 读取登录响应
        response = sock.recv(1024).decode('utf-8').strip()
        print(f"登录响应: {response}")
        
        if "OK" in response:
            print("✅ 登录成功!")
            
            # 测试几个基本命令
            commands = [
                "A002 CAPABILITY",
                "A003 LIST \"\" \"*\"",
                "A004 SELECT INBOX",
                "A005 SEARCH ALL"
            ]
            
            print("\n4. 测试基本IMAP命令...")
            for cmd in commands:
                print(f"\n发送命令: {cmd}")
                sock.send((cmd + "\r\n").encode('utf-8'))
                time.sleep(0.5)
                
                try:
                    response = sock.recv(1024).decode('utf-8').strip()
                    print(f"响应: {response}")
                except:
                    print("无响应或连接中断")
            
            # 登出
            print("\n5. 登出...")
            sock.send("A999 LOGOUT\r\n".encode('utf-8'))
            try:
                response = sock.recv(1024).decode('utf-8').strip()
                print(f"登出响应: {response}")
            except:
                pass
            
        else:
            print("❌ 登录失败!")
            
            # 尝试用邮箱地址登录
            print(f"\n尝试用邮箱地址登录: {email}")
            login_cmd2 = f"A002 LOGIN {email} {password}\r\n"
            print(f"发送命令: {login_cmd2.strip()}")
            sock.send(login_cmd2.encode('utf-8'))
            
            response2 = sock.recv(1024).decode('utf-8').strip()
            print(f"登录响应: {response2}")
            
            if "OK" in response2:
                print("✅ 邮箱地址登录成功!")
            else:
                print("❌ 邮箱地址登录也失败!")
        
        # 关闭连接
        sock.close()
        print("\n6. 连接已关闭")
        
        print("\n" + "=" * 60)
        print("✅ IMAP连接测试完成!")
        print("=" * 60)
        
        print("\n测试结果:")
        print("✅ IMAP服务器在143端口正常运行")
        print("✅ 能够建立TCP连接")
        print("✅ 服务器返回欢迎消息")
        print("✅ 能够发送和接收IMAP命令")
        print("ℹ️  当前为简单IMAP实现，适合基础测试")
        
    except socket.timeout:
        print("❌ 连接超时!")
        print("请检查IMAP服务器是否在运行")
    except ConnectionRefusedError:
        print("❌ 连接被拒绝!")
        print("请检查IMAP服务器是否在143端口监听")
    except Exception as e:
        print(f"❌ 连接出错: {e}")
    
    return True

if __name__ == "__main__":
    test_imap_connection()
