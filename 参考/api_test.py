#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
NBEmail API 测试脚本
用于测试所有可用的 API 端点
"""

import requests
import json
import sys

class NBEmailAPITester:
    def __init__(self, base_url="http://localhost:8080"):
        self.base_url = base_url
        self.session = requests.Session()
        self.session.headers.update({
            'Content-Type': 'application/json; charset=utf-8',
            'User-Agent': 'NBEmail-API-Tester/1.0'
        })
    
    def test_login(self, email, password):
        """测试登录 API"""
        print("=" * 60)
        print("测试登录 API")
        print("=" * 60)
        
        login_data = {
            "email": email,
            "password": password
        }
        
        try:
            response = self.session.post(
                f"{self.base_url}/api/login",
                json=login_data,
                timeout=10
            )
            
            print(f"请求URL: POST {self.base_url}/api/login")
            print(f"请求数据: {json.dumps(login_data, ensure_ascii=False, indent=2)}")
            print(f"响应状态码: {response.status_code}")
            print(f"响应内容: {json.dumps(response.json(), ensure_ascii=False, indent=2)}")
            
            if response.status_code == 200 and response.json().get('success'):
                print("✓ 登录测试成功")
                return True
            else:
                print("✗ 登录测试失败")
                return False
                
        except Exception as e:
            print(f"✗ 登录测试异常: {e}")
            return False
    
    def test_get_emails(self):
        """测试获取邮件列表 API"""
        print("\n" + "=" * 60)
        print("测试获取邮件列表 API")
        print("=" * 60)
        
        try:
            response = self.session.get(
                f"{self.base_url}/api/emails?page=1&limit=5",
                timeout=10
            )
            
            print(f"请求URL: GET {self.base_url}/api/emails?page=1&limit=5")
            print(f"响应状态码: {response.status_code}")
            print(f"响应内容: {json.dumps(response.json(), ensure_ascii=False, indent=2)}")
            
            if response.status_code == 200:
                print("✓ 获取邮件列表测试成功")
                return response.json().get('data', {})
            else:
                print("✗ 获取邮件列表测试失败")
                return None
                
        except Exception as e:
            print(f"✗ 获取邮件列表测试异常: {e}")
            return None
    
    def test_get_mailboxes(self):
        """测试获取邮箱列表 API"""
        print("\n" + "=" * 60)
        print("测试获取邮箱列表 API")
        print("=" * 60)
        
        try:
            response = self.session.get(
                f"{self.base_url}/api/mailboxes",
                timeout=10
            )
            
            print(f"请求URL: GET {self.base_url}/api/mailboxes")
            print(f"响应状态码: {response.status_code}")
            print(f"响应内容: {json.dumps(response.json(), ensure_ascii=False, indent=2)}")
            
            if response.status_code == 200:
                print("✓ 获取邮箱列表测试成功")
                return True
            else:
                print("✗ 获取邮箱列表测试失败")
                return False
                
        except Exception as e:
            print(f"✗ 获取邮箱列表测试异常: {e}")
            return False
    
    def test_get_users(self):
        """测试获取用户列表 API (管理员权限)"""
        print("\n" + "=" * 60)
        print("测试获取用户列表 API")
        print("=" * 60)
        
        try:
            response = self.session.get(
                f"{self.base_url}/api/users",
                timeout=10
            )
            
            print(f"请求URL: GET {self.base_url}/api/users")
            print(f"响应状态码: {response.status_code}")
            print(f"响应内容: {json.dumps(response.json(), ensure_ascii=False, indent=2)}")
            
            if response.status_code == 200:
                print("✓ 获取用户列表测试成功")
                return True
            else:
                print("✗ 获取用户列表测试失败")
                return False
                
        except Exception as e:
            print(f"✗ 获取用户列表测试异常: {e}")
            return False
    
    def test_get_domains(self):
        """测试获取域名列表 API"""
        print("\n" + "=" * 60)
        print("测试获取域名列表 API")
        print("=" * 60)
        
        try:
            response = self.session.get(
                f"{self.base_url}/api/domains",
                timeout=10
            )
            
            print(f"请求URL: GET {self.base_url}/api/domains")
            print(f"响应状态码: {response.status_code}")
            print(f"响应内容: {json.dumps(response.json(), ensure_ascii=False, indent=2)}")
            
            if response.status_code == 200:
                print("✓ 获取域名列表测试成功")
                return True
            else:
                print("✗ 获取域名列表测试失败")
                return False
                
        except Exception as e:
            print(f"✗ 获取域名列表测试异常: {e}")
            return False
    
    def test_logout(self):
        """测试登出 API"""
        print("\n" + "=" * 60)
        print("测试登出 API")
        print("=" * 60)
        
        try:
            response = self.session.post(
                f"{self.base_url}/api/logout",
                timeout=10
            )
            
            print(f"请求URL: POST {self.base_url}/api/logout")
            print(f"响应状态码: {response.status_code}")
            print(f"响应内容: {json.dumps(response.json(), ensure_ascii=False, indent=2)}")
            
            if response.status_code == 200:
                print("✓ 登出测试成功")
                return True
            else:
                print("✗ 登出测试失败")
                return False
                
        except Exception as e:
            print(f"✗ 登出测试异常: {e}")
            return False
    
    def run_all_tests(self, email, password):
        """运行所有测试"""
        print("NBEmail API 测试开始")
        print("服务器地址:", self.base_url)
        print("测试账号:", email)
        
        # 测试登录
        if not self.test_login(email, password):
            print("\n登录失败，无法继续测试")
            return
        
        # 测试各种 API
        self.test_get_emails()
        self.test_get_mailboxes()
        self.test_get_users()
        self.test_get_domains()
        
        # 测试登出
        self.test_logout()
        
        print("\n" + "=" * 60)
        print("所有 API 测试完成")
        print("=" * 60)

def main():
    if len(sys.argv) < 2:
        print("用法: python api_test.py <邮箱地址>.<密码>")
        print("示例: python api_test.py 2014131458@qq.com.tgx123456")
        return
    
    # 解析邮箱和密码
    if '.' in sys.argv[1] and '@' in sys.argv[1]:
        parts = sys.argv[1].rsplit('.', 1)
        if len(parts) == 2:
            email, password = parts
        else:
            print("格式错误，请使用: 邮箱地址.密码")
            return
    else:
        print("格式错误，请使用: 邮箱地址.密码")
        return
    
    # 创建测试器实例
    tester = NBEmailAPITester()
    
    # 运行所有测试
    tester.run_all_tests(email, password)

if __name__ == "__main__":
    main()
