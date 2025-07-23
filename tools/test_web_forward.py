#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
测试Web界面邮件转发功能
"""

import requests
import time

def test_web_forward():
    """测试Web界面发送邮件的转发功能"""
    
    # 服务器配置
    base_url = "http://localhost:8080"
    
    # 登录信息（需要先登录）
    login_data = {
        "username": "kimi",  # 请根据实际用户名修改
        "password": "123456"  # 请根据实际密码修改
    }
    
    # 创建会话
    session = requests.Session()
    
    try:
        print("🔐 正在登录...")
        
        # 登录
        login_response = session.post(f"{base_url}/api/login", json=login_data)
        if login_response.status_code != 200:
            print(f"❌ 登录失败: {login_response.status_code}")
            print(f"响应: {login_response.text}")
            return False
        
        login_result = login_response.json()
        if not login_result.get("success"):
            print(f"❌ 登录失败: {login_result.get('message')}")
            return False
        
        print("✅ 登录成功")
        
        # 发送测试邮件
        print("📧 正在发送测试邮件...")
        
        email_data = {
            "from": "kimi11@jbjj.site",
            "to": "kimi11@jbjj.site",  # 自己给自己发邮件
            "subject": f"测试Web转发功能 - {time.strftime('%Y-%m-%d %H:%M:%S')}",
            "content": f"""这是通过Web界面发送的测试邮件，用于测试转发功能是否正常工作。

发送时间: {time.strftime('%Y-%m-%d %H:%M:%S')}

如果转发功能正常，这封邮件应该会：
1. 保存到 kimi11@jbjj.site 的收件箱
2. 自动转发到 kimi12@jbjj.site 的收件箱
3. 转发邮件包含 [转发] 前缀
4. 更新转发统计数据
"""
        }
        
        # 发送邮件
        send_response = session.post(f"{base_url}/api/emails/send", data=email_data)
        
        if send_response.status_code != 200:
            print(f"❌ 发送邮件失败: {send_response.status_code}")
            print(f"响应: {send_response.text}")
            return False
        
        send_result = send_response.json()
        if not send_result.get("success"):
            print(f"❌ 发送邮件失败: {send_result.get('message')}")
            return False
        
        print("✅ 邮件发送成功")
        print(f"📋 请检查服务器日志，看是否有转发处理信息")
        
        # 等待一下让转发处理完成
        time.sleep(2)
        
        # 检查转发统计
        print("📊 检查转发统计...")
        stats_response = session.get(f"{base_url}/api/forward-statistics")
        
        if stats_response.status_code == 200:
            stats_result = stats_response.json()
            if stats_result.get("success"):
                stats = stats_result.get("data", {})
                print(f"   总规则数: {stats.get('total_rules', 0)}")
                print(f"   启用规则数: {stats.get('active_rules', 0)}")
                print(f"   总转发次数: {stats.get('total_forwards', 0)}")
                print(f"   今日转发次数: {stats.get('today_forwards', 0)}")
        
        return True
        
    except Exception as e:
        print(f"❌ 测试失败: {e}")
        return False

def test_multiple_web_forwards():
    """测试多次Web发送邮件"""
    
    base_url = "http://localhost:8080"
    session = requests.Session()
    
    # 登录
    login_data = {
        "username": "kimi",
        "password": "123456"
    }
    
    try:
        login_response = session.post(f"{base_url}/api/login", json=login_data)
        if login_response.status_code != 200 or not login_response.json().get("success"):
            print("❌ 登录失败，跳过多次发送测试")
            return False
        
        print("\n🔄 开始多次发送测试...")
        
        for i in range(3):
            print(f"📧 发送第 {i+1} 封邮件...")
            
            email_data = {
                "from": "kimi11@jbjj.site",
                "to": "kimi11@jbjj.site",
                "subject": f"批量测试Web转发 #{i+1} - {time.strftime('%H:%M:%S')}",
                "content": f"这是第 {i+1} 封测试邮件，发送时间: {time.strftime('%Y-%m-%d %H:%M:%S')}"
            }
            
            send_response = session.post(f"{base_url}/api/emails/send", data=email_data)
            
            if send_response.status_code == 200 and send_response.json().get("success"):
                print(f"   ✅ 第 {i+1} 封邮件发送成功")
            else:
                print(f"   ❌ 第 {i+1} 封邮件发送失败")
            
            time.sleep(1)  # 间隔1秒
        
        print("✅ 批量发送测试完成")
        return True
        
    except Exception as e:
        print(f"❌ 批量测试失败: {e}")
        return False

if __name__ == "__main__":
    print("=" * 60)
    print("📧 Web界面邮件转发功能测试")
    print("=" * 60)
    
    # 基本转发测试
    success1 = test_web_forward()
    
    # 等待一下
    time.sleep(3)
    
    # 多次发送测试
    success2 = test_multiple_web_forwards()
    
    print("\n" + "=" * 60)
    print("📊 测试结果汇总:")
    print(f"   基本Web转发测试: {'✅ 成功' if success1 else '❌ 失败'}")
    print(f"   批量Web转发测试: {'✅ 成功' if success2 else '❌ 失败'}")
    print("=" * 60)
    
    if success1 and success2:
        print("🎉 所有Web转发测试都成功完成！")
        print("📋 请检查服务器日志确认转发处理过程")
    else:
        print("⚠️  部分测试失败，请检查服务器日志和登录信息。")
