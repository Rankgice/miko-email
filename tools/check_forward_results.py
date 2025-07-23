#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
检查邮件转发结果
"""

import sqlite3
import os

def check_forward_results():
    """检查转发结果"""
    
    # 数据库路径
    db_path = "../miko_email.db"
    
    if not os.path.exists(db_path):
        print("❌ 数据库文件不存在")
        return
    
    try:
        conn = sqlite3.connect(db_path)
        cursor = conn.cursor()
        
        print("=" * 60)
        print("📧 邮件转发结果检查")
        print("=" * 60)
        
        # 检查 kimi11@jbjj.site 的收件箱
        print("\n📥 检查 kimi11@jbjj.site 的收件箱:")
        cursor.execute("""
            SELECT e.id, e.from_addr, e.subject, e.created_at
            FROM emails e
            JOIN mailboxes m ON e.mailbox_id = m.id
            WHERE m.email = 'kimi11@jbjj.site' AND e.folder = 'inbox'
            ORDER BY e.created_at DESC
            LIMIT 10
        """)
        
        emails = cursor.fetchall()
        if emails:
            for email in emails:
                print(f"   📧 ID: {email[0]}, 发件人: {email[1]}, 主题: {email[2]}, 时间: {email[3]}")
        else:
            print("   ❌ 没有找到邮件")
        
        # 检查 kimi12@jbjj.site 的收件箱
        print("\n📥 检查 kimi12@jbjj.site 的收件箱 (转发目标):")
        cursor.execute("""
            SELECT e.id, e.from_addr, e.subject, e.created_at
            FROM emails e
            JOIN mailboxes m ON e.mailbox_id = m.id
            WHERE m.email = 'kimi12@jbjj.site' AND e.folder = 'inbox'
            ORDER BY e.created_at DESC
            LIMIT 10
        """)
        
        forward_emails = cursor.fetchall()
        if forward_emails:
            for email in forward_emails:
                print(f"   📧 ID: {email[0]}, 发件人: {email[1]}, 主题: {email[2]}, 时间: {email[3]}")
        else:
            print("   ❌ 没有找到转发邮件")
        
        # 检查转发规则统计
        print("\n📊 检查转发规则统计:")
        cursor.execute("""
            SELECT ef.id, ef.source_email, ef.target_email, ef.forward_count, ef.last_forward_at
            FROM email_forwards ef
            WHERE ef.enabled = 1
        """)
        
        rules = cursor.fetchall()
        if rules:
            for rule in rules:
                print(f"   📋 规则ID: {rule[0]}, {rule[1]} -> {rule[2]}, 转发次数: {rule[3]}, 最后转发: {rule[4]}")
        else:
            print("   ❌ 没有找到转发规则")
        
        # 统计总数
        print("\n📈 统计信息:")
        
        # 原始邮件数量
        cursor.execute("SELECT COUNT(*) FROM emails WHERE mailbox_id = (SELECT id FROM mailboxes WHERE email = 'kimi11@jbjj.site')")
        original_count = cursor.fetchone()[0]
        print(f"   📧 kimi11@jbjj.site 总邮件数: {original_count}")
        
        # 转发邮件数量
        cursor.execute("SELECT COUNT(*) FROM emails WHERE mailbox_id = (SELECT id FROM mailboxes WHERE email = 'kimi12@jbjj.site')")
        forward_count = cursor.fetchone()[0]
        print(f"   📧 kimi12@jbjj.site 总邮件数: {forward_count}")
        
        # 转发规则总转发次数
        cursor.execute("SELECT SUM(forward_count) FROM email_forwards WHERE enabled = 1")
        total_forwards = cursor.fetchone()[0] or 0
        print(f"   📊 转发规则总转发次数: {total_forwards}")
        
        print("\n" + "=" * 60)
        
        # 验证转发是否成功
        if original_count > 0 and forward_count > 0 and total_forwards > 0:
            print("🎉 转发功能测试成功！")
            print("✅ 原始邮件已保存")
            print("✅ 转发邮件已送达")
            print("✅ 转发统计已更新")
        else:
            print("⚠️  转发功能可能存在问题")
            if original_count == 0:
                print("❌ 没有找到原始邮件")
            if forward_count == 0:
                print("❌ 没有找到转发邮件")
            if total_forwards == 0:
                print("❌ 转发统计未更新")
        
        conn.close()
        
    except Exception as e:
        print(f"❌ 检查失败: {e}")

def show_email_content():
    """显示最新邮件的内容"""
    
    db_path = "../miko_email.db"
    
    if not os.path.exists(db_path):
        print("❌ 数据库文件不存在")
        return
    
    try:
        conn = sqlite3.connect(db_path)
        cursor = conn.cursor()
        
        print("\n" + "=" * 60)
        print("📧 最新邮件内容预览")
        print("=" * 60)
        
        # 获取最新的原始邮件
        cursor.execute("""
            SELECT e.subject, e.body, e.from_addr, e.created_at
            FROM emails e
            JOIN mailboxes m ON e.mailbox_id = m.id
            WHERE m.email = 'kimi11@jbjj.site'
            ORDER BY e.created_at DESC
            LIMIT 1
        """)
        
        original = cursor.fetchone()
        if original:
            print(f"\n📧 原始邮件 (kimi11@jbjj.site):")
            print(f"   主题: {original[0]}")
            print(f"   发件人: {original[2]}")
            print(f"   时间: {original[3]}")
            print(f"   内容预览: {original[1][:100]}...")
        
        # 获取最新的转发邮件
        cursor.execute("""
            SELECT e.subject, e.body, e.from_addr, e.created_at
            FROM emails e
            JOIN mailboxes m ON e.mailbox_id = m.id
            WHERE m.email = 'kimi12@jbjj.site'
            ORDER BY e.created_at DESC
            LIMIT 1
        """)
        
        forwarded = cursor.fetchone()
        if forwarded:
            print(f"\n📧 转发邮件 (kimi12@jbjj.site):")
            print(f"   主题: {forwarded[0]}")
            print(f"   发件人: {forwarded[2]}")
            print(f"   时间: {forwarded[3]}")
            print(f"   内容预览: {forwarded[1][:200]}...")
        
        conn.close()
        
    except Exception as e:
        print(f"❌ 显示邮件内容失败: {e}")

if __name__ == "__main__":
    check_forward_results()
    show_email_content()
