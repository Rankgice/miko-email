#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
调试 API 数据脚本
"""

import requests
import json

def debug_api():
    base_url = "http://me.youddns.site:8080"
    
    # 创建会话
    session = requests.Session()
    session.headers.update({
        'Content-Type': 'application/json; charset=utf-8',
        'User-Agent': 'NBEmail-Debug/1.0'
    })
    
    # 登录
    print("=== 登录测试 ===")
    login_data = {
        "email": "2014131458@qq.com",
        "password": "tgx123456"
    }
    
    response = session.post(f"{base_url}/api/login", json=login_data)
    print(f"登录状态码: {response.status_code}")
    login_result = response.json()
    print(f"登录结果: {json.dumps(login_result, ensure_ascii=False, indent=2)}")
    
    if not login_result.get('success'):
        print("登录失败，退出")
        return
    
    # 获取邮件列表
    print("\n=== 邮件列表测试 ===")
    emails_response = session.get(f"{base_url}/api/emails?page=1&limit=5")
    print(f"邮件列表状态码: {emails_response.status_code}")
    emails_result = emails_response.json()
    print(f"邮件列表结果: {json.dumps(emails_result, ensure_ascii=False, indent=2)}")
    
    # 如果有邮件，获取第一封邮件的详情
    if emails_result.get('success') and emails_result.get('data', {}).get('emails'):
        first_email = emails_result['data']['emails'][0]
        email_id = first_email.get('id')
        
        print(f"\n=== 邮件详情测试 (ID: {email_id}) ===")
        detail_response = session.get(f"{base_url}/api/emails/{email_id}")
        print(f"邮件详情状态码: {detail_response.status_code}")
        detail_result = detail_response.json()
        print(f"邮件详情结果: {json.dumps(detail_result, ensure_ascii=False, indent=2)}")
        
        # 分析字段
        if detail_result.get('success'):
            email_data = detail_result['data']
            print(f"\n=== 字段分析 ===")
            print(f"ID: {email_data.get('id')}")
            print(f"发件人字段 (from_addr): {repr(email_data.get('from_addr'))}")
            print(f"收件人字段 (to_addr): {repr(email_data.get('to_addr'))}")
            print(f"主题字段 (subject): {repr(email_data.get('subject'))}")
            print(f"正文长度: {len(email_data.get('body', ''))}")
            
            # 检查所有字段
            print(f"\n=== 所有字段 ===")
            for key, value in email_data.items():
                if isinstance(value, str) and len(value) > 100:
                    print(f"{key}: {repr(value[:100])}... (长度: {len(value)})")
                else:
                    print(f"{key}: {repr(value)}")

if __name__ == "__main__":
    debug_api()
