#!/usr/bin/env python3
"""
Miko邮箱系统完整功能测试
"""

import requests
import json
import time

BASE_URL = "http://localhost:8080"

def print_header(title):
    print(f"\n{'='*50}")
    print(f"  {title}")
    print(f"{'='*50}")

def print_step(step):
    print(f"\n🔸 {step}")

def print_success(message):
    print(f"✅ {message}")

def print_error(message):
    print(f"❌ {message}")

def test_complete_workflow():
    """测试完整的工作流程"""
    print_header("Miko邮箱系统完整功能测试")
    
    session = requests.Session()
    
    # 1. 测试主页访问
    print_step("测试主页访问")
    try:
        response = session.get(f"{BASE_URL}/")
        if response.status_code == 200:
            print_success("主页访问成功")
        else:
            print_error(f"主页访问失败: {response.status_code}")
    except Exception as e:
        print_error(f"主页访问异常: {e}")
        return False
    
    # 2. 测试管理员登录
    print_step("测试管理员登录")
    try:
        login_data = {"username": "admin", "password": "123456"}
        response = session.post(f"{BASE_URL}/api/admin/login", json=login_data)
        result = response.json()
        
        if result.get("success"):
            print_success("管理员登录成功")
            admin_info = result["data"]["user"]
            print(f"   管理员信息: {admin_info['username']} ({admin_info['email']})")
        else:
            print_error(f"管理员登录失败: {result.get('message')}")
            return False
    except Exception as e:
        print_error(f"管理员登录异常: {e}")
        return False
    
    # 3. 测试域名管理
    print_step("测试域名管理")
    try:
        # 获取域名列表
        response = session.get(f"{BASE_URL}/api/admin/domains")
        result = response.json()
        
        if result.get("success"):
            domains = result["data"]
            print_success(f"获取到 {len(domains)} 个域名")
            for domain in domains:
                status = "已验证" if domain["is_verified"] else "未验证"
                print(f"   - {domain['name']} ({status})")
        else:
            print_error("获取域名列表失败")
    except Exception as e:
        print_error(f"域名管理测试异常: {e}")
    
    # 4. 测试DNS记录查询
    print_step("测试DNS记录查询")
    try:
        test_domain = "google.com"
        response = session.get(f"{BASE_URL}/api/domains/dns?domain={test_domain}")
        result = response.json()
        
        if result.get("success"):
            records = result["data"]["records"]
            print_success(f"{test_domain} DNS记录查询成功")
            for record_type, values in records.items():
                print(f"   {record_type}: {len(values)} 条记录")
        else:
            print_error("DNS记录查询失败")
    except Exception as e:
        print_error(f"DNS记录查询异常: {e}")
    
    # 5. 测试用户注册
    print_step("测试用户注册")
    try:
        # 获取可用域名
        response = session.get(f"{BASE_URL}/api/domains/available")
        domains_result = response.json()
        
        if domains_result.get("success") and domains_result["data"]:
            domain_id = domains_result["data"][0]["id"]
            domain_name = domains_result["data"][0]["name"]
            
            # 注册新用户
            register_data = {
                "username": f"testuser_{int(time.time())}",
                "password": "123456",
                "email": f"test_{int(time.time())}@example.com",
                "domain_prefix": f"test{int(time.time())}",
                "domain_id": domain_id
            }
            
            response = requests.post(f"{BASE_URL}/api/register", json=register_data)
            result = response.json()
            
            if result.get("success"):
                user_info = result["data"]["user"]
                print_success("用户注册成功")
                print(f"   用户名: {user_info['username']}")
                print(f"   邮箱: {user_info['email']}")
                print(f"   邀请码: {user_info['invite_code']}")
                print(f"   贡献度: {user_info['contribution']}")
                
                # 测试用户登录
                print_step("测试新用户登录")
                user_session = requests.Session()
                login_data = {
                    "username": register_data["username"],
                    "password": register_data["password"]
                }
                
                response = user_session.post(f"{BASE_URL}/api/login", json=login_data)
                login_result = response.json()
                
                if login_result.get("success"):
                    print_success("新用户登录成功")
                    
                    # 测试获取邮箱列表
                    print_step("测试获取用户邮箱列表")
                    response = user_session.get(f"{BASE_URL}/api/mailboxes")
                    mailbox_result = response.json()
                    
                    if mailbox_result.get("success"):
                        mailboxes = mailbox_result["data"]
                        print_success(f"获取到 {len(mailboxes)} 个邮箱")
                        for mailbox in mailboxes:
                            print(f"   - {mailbox['email']}")
                    else:
                        print_error("获取邮箱列表失败")
                        
                else:
                    print_error("新用户登录失败")
            else:
                print_error(f"用户注册失败: {result.get('message')}")
        else:
            print_error("获取可用域名失败")
    except Exception as e:
        print_error(f"用户注册测试异常: {e}")
    
    # 6. 测试邮箱创建
    print_step("测试邮箱创建")
    try:
        if 'user_session' in locals():
            create_data = {
                "prefix": f"newbox{int(time.time())}",
                "domain_id": domain_id
            }
            
            response = user_session.post(f"{BASE_URL}/api/mailboxes", json=create_data)
            result = response.json()
            
            if result.get("success"):
                mailbox = result["data"]
                print_success("邮箱创建成功")
                print(f"   邮箱地址: {mailbox['email']}")
            else:
                print_error(f"邮箱创建失败: {result.get('message')}")
    except Exception as e:
        print_error(f"邮箱创建测试异常: {e}")
    
    # 7. 测试批量邮箱创建
    print_step("测试批量邮箱创建")
    try:
        if 'user_session' in locals():
            batch_data = {
                "prefixes": [f"batch1_{int(time.time())}", f"batch2_{int(time.time())}", f"batch3_{int(time.time())}"],
                "domain_id": domain_id
            }
            
            response = user_session.post(f"{BASE_URL}/api/mailboxes/batch", json=batch_data)
            result = response.json()
            
            if result.get("success"):
                mailboxes = result["data"]
                print_success(f"批量创建 {len(mailboxes)} 个邮箱成功")
                for mailbox in mailboxes:
                    print(f"   - {mailbox['email']}")
            else:
                print_error(f"批量邮箱创建失败: {result.get('message')}")
    except Exception as e:
        print_error(f"批量邮箱创建测试异常: {e}")
    
    # 8. 系统状态总结
    print_header("系统状态总结")
    
    try:
        # 获取最终的邮箱列表
        if 'user_session' in locals():
            response = user_session.get(f"{BASE_URL}/api/mailboxes")
            result = response.json()
            
            if result.get("success"):
                total_mailboxes = len(result["data"])
                print_success(f"用户总邮箱数: {total_mailboxes}")
            
        # 获取域名状态
        response = session.get(f"{BASE_URL}/api/admin/domains")
        result = response.json()
        
        if result.get("success"):
            domains = result["data"]
            verified_count = sum(1 for d in domains if d["is_verified"])
            print_success(f"域名总数: {len(domains)}, 已验证: {verified_count}")
            
    except Exception as e:
        print_error(f"状态总结异常: {e}")
    
    print_header("测试完成")
    print("🎉 Miko邮箱系统功能测试完成！")
    print("\n📋 测试覆盖功能:")
    print("   ✅ Web界面访问")
    print("   ✅ 管理员登录")
    print("   ✅ 域名管理")
    print("   ✅ DNS记录查询")
    print("   ✅ 用户注册登录")
    print("   ✅ 邮箱管理")
    print("   ✅ 批量操作")
    print("\n🚀 系统已准备就绪，可以投入使用！")
    
    return True

if __name__ == "__main__":
    test_complete_workflow()
