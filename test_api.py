#!/usr/bin/env python3
"""
Miko邮箱系统API测试脚本
"""

import requests
import json

BASE_URL = "http://localhost:8080"

def test_admin_login():
    """测试管理员登录"""
    print("=== 测试管理员登录 ===")
    
    session = requests.Session()
    login_data = {
        "username": "admin",
        "password": "123456"
    }
    
    response = session.post(f"{BASE_URL}/api/admin/login", json=login_data)
    print(f"状态码: {response.status_code}")
    print(f"响应: {response.json()}")
    
    if response.json().get("success"):
        print("✅ 管理员登录成功")
        return session
    else:
        print("❌ 管理员登录失败")
        return None

def test_domain_verification(session):
    """测试域名验证"""
    print("\n=== 测试域名验证 ===")
    
    # 获取域名列表
    response = session.get(f"{BASE_URL}/api/admin/domains")
    print(f"域名列表: {response.json()}")
    
    if response.json().get("success") and response.json().get("data"):
        domain_id = response.json()["data"][0]["id"]
        
        # 验证第一个域名
        verify_response = session.post(f"{BASE_URL}/api/admin/domains/{domain_id}/verify")
        print(f"验证结果: {verify_response.json()}")
        
        if verify_response.json().get("success"):
            print("✅ 域名验证完成")
        else:
            print("❌ 域名验证失败")
    else:
        print("❌ 获取域名列表失败")

def test_dns_records():
    """测试DNS记录查询"""
    print("\n=== 测试DNS记录查询 ===")
    
    test_domains = ["google.com", "github.com", "localhost"]
    
    for domain in test_domains:
        print(f"\n查询域名: {domain}")
        response = requests.get(f"{BASE_URL}/api/domains/dns?domain={domain}")
        
        if response.status_code == 200:
            data = response.json()
            if data.get("success"):
                records = data["data"]["records"]
                print(f"✅ DNS记录查询成功:")
                for record_type, values in records.items():
                    print(f"  {record_type}: {values}")
            else:
                print(f"❌ DNS记录查询失败: {data.get('message')}")
        else:
            print(f"❌ HTTP错误: {response.status_code}")

def test_user_registration():
    """测试用户注册"""
    print("\n=== 测试用户注册 ===")
    
    # 先获取可用域名
    response = requests.get(f"{BASE_URL}/api/domains/available")
    if not response.json().get("success"):
        print("❌ 获取可用域名失败")
        return
    
    domains = response.json()["data"]
    if not domains:
        print("❌ 没有可用域名")
        return
    
    domain_id = domains[0]["id"]
    
    # 注册用户
    register_data = {
        "username": "testuser2",
        "password": "123456",
        "email": "test2@example.com",
        "domain_prefix": "testuser2",
        "domain_id": domain_id
    }
    
    response = requests.post(f"{BASE_URL}/api/register", json=register_data)
    print(f"注册结果: {response.json()}")
    
    if response.json().get("success"):
        print("✅ 用户注册成功")
    else:
        print("❌ 用户注册失败")

def main():
    """主函数"""
    print("🚀 开始测试Miko邮箱系统API")
    
    # 测试管理员登录
    session = test_admin_login()
    
    if session:
        # 测试域名验证
        test_domain_verification(session)
    
    # 测试DNS记录查询
    test_dns_records()
    
    # 测试用户注册
    test_user_registration()
    
    print("\n🎉 API测试完成")

if __name__ == "__main__":
    main()
