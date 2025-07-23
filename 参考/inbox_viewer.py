#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
NBEmail 收件箱查看器 - 基于 Web API
使用方法: python inbox_viewer.py <邮箱地址> <密码>
示例: python inbox_viewer.py 2014131458@qq.com tgx123456
"""

import requests
import json
import sys
import time
import threading
import os
import base64
import re
from datetime import datetime
import html
from email.mime.text import MIMEText
from email import message_from_string
from html.parser import HTMLParser

class HTMLTextExtractor(HTMLParser):
    """HTML 到纯文本转换器"""
    def __init__(self):
        super().__init__()
        self.text = []

    def handle_data(self, data):
        self.text.append(data)

    def handle_starttag(self, tag, attrs):
        if tag.lower() in ['br', 'p', 'div']:
            self.text.append('\n')

    def handle_endtag(self, tag):
        if tag.lower() in ['p', 'div']:
            self.text.append('\n')

    def get_text(self):
        return ''.join(self.text).strip()

class NBEmailViewer:
    def __init__(self, base_url="http://me.youddns.site:8080"):
        self.base_url = base_url
        self.session = requests.Session()
        self.session.headers.update({
            'Content-Type': 'application/json; charset=utf-8',
            'User-Agent': 'NBEmail-Python-Viewer/1.0'
        })
        self.auto_refresh = False
        self.refresh_interval = 15  # 15秒刷新间隔
        self.refresh_thread = None
        self.stop_refresh = False

    def html_to_text(self, html_content):
        """将 HTML 内容转换为纯文本"""
        if not html_content or not html_content.strip():
            return ""

        try:
            # 使用 HTMLTextExtractor 提取纯文本
            extractor = HTMLTextExtractor()
            extractor.feed(html_content)
            text = extractor.get_text()

            # 清理多余的空行和空格
            lines = [line.strip() for line in text.split('\n')]
            lines = [line for line in lines if line]  # 移除空行

            return '\n'.join(lines)
        except Exception as e:
            # 如果解析失败，使用简单的正则表达式清理
            import re
            text = re.sub(r'<[^>]+>', '', html_content)  # 移除HTML标签
            text = html.unescape(text)  # 解码HTML实体
            return text.strip()

    def login(self, email, password):
        """登录到 NBEmail 系统"""
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
            
            if response.status_code == 200:
                result = response.json()
                if result.get('success'):
                    print(f"✓ 登录成功: {email}")
                    return True
                else:
                    print(f"✗ 登录失败: {result.get('message', '未知错误')}")
                    return False
            else:
                print(f"✗ 登录请求失败: HTTP {response.status_code}")
                return False
                
        except requests.exceptions.RequestException as e:
            print(f"✗ 网络请求失败: {e}")
            return False
    
    def get_emails(self, page=1, limit=20, folder="inbox", search=""):
        """获取邮件列表"""
        params = {
            "page": page,
            "limit": limit,
            "folder": folder
        }
        if search:
            params["search"] = search
        
        try:
            response = self.session.get(
                f"{self.base_url}/api/emails",
                params=params,
                timeout=10
            )
            
            if response.status_code == 200:
                result = response.json()
                if result.get('success'):
                    return result.get('data', {})
                else:
                    print(f"✗ 获取邮件失败: {result.get('message', '未知错误')}")
                    return None
            else:
                print(f"✗ 获取邮件请求失败: HTTP {response.status_code}")
                return None
                
        except requests.exceptions.RequestException as e:
            print(f"✗ 网络请求失败: {e}")
            return None
    
    def get_email_detail(self, email_id):
        """获取单个邮件详情"""
        try:
            response = self.session.get(
                f"{self.base_url}/api/emails/{email_id}",
                timeout=10
            )
            
            if response.status_code == 200:
                result = response.json()
                if result.get('success'):
                    return result.get('data')
                else:
                    print(f"✗ 获取邮件详情失败: {result.get('message', '未知错误')}")
                    return None
            else:
                print(f"✗ 获取邮件详情请求失败: HTTP {response.status_code}")
                return None
                
        except requests.exceptions.RequestException as e:
            print(f"✗ 网络请求失败: {e}")
            return None
    
    def format_date(self, date_str):
        """格式化日期"""
        try:
            # 尝试解析不同的日期格式
            for fmt in ["%Y-%m-%dT%H:%M:%S.%fZ", "%Y-%m-%dT%H:%M:%SZ", "%Y-%m-%d %H:%M:%S"]:
                try:
                    dt = datetime.strptime(date_str, fmt)
                    return dt.strftime("%Y-%m-%d %H:%M:%S")
                except ValueError:
                    continue
            return date_str
        except:
            return date_str
    
    def clean_text(self, text):
        """清理文本内容"""
        if not text:
            return ""
        # 解码 HTML 实体
        text = html.unescape(text)
        # 移除多余的空白字符
        text = ' '.join(text.split())
        return text

    def decode_email_body(self, body):
        """解码邮件正文"""
        if not body:
            return ""

        # 如果包含 MIME 多部分内容，尝试解码
        if "Content-Type: text/plain" in body and "Content-Transfer-Encoding: base64" in body:
            # 使用正则表达式提取 base64 内容
            pattern = r'Content-Type: text/plain[^;]*(?:; charset="?([^";\r\n]*)"?)?\s*\r?\n(?:Content-Transfer-Encoding: base64\s*\r?\n)?\s*\r?\n([A-Za-z0-9+/=\s]+?)(?:\r?\n--|\r?\n\r?\n|$)'
            match = re.search(pattern, body, re.MULTILINE | re.DOTALL)

            if match:
                charset = match.group(1) or 'utf-8'
                base64_content = match.group(2).strip()

                try:
                    # 解码 base64
                    decoded_bytes = base64.b64decode(base64_content)

                    # 尝试用指定的字符集解码
                    try:
                        decoded_text = decoded_bytes.decode(charset)
                    except (UnicodeDecodeError, LookupError):
                        # 如果指定字符集失败，尝试常见字符集
                        for fallback_charset in ['utf-8', 'gbk', 'gb2312', 'latin1']:
                            try:
                                decoded_text = decoded_bytes.decode(fallback_charset)
                                break
                            except (UnicodeDecodeError, LookupError):
                                continue
                        else:
                            decoded_text = decoded_bytes.decode('utf-8', errors='ignore')

                    return decoded_text.strip()

                except Exception as e:
                    print(f"Base64解码失败: {e}")

        # 如果包含 quoted-printable 编码
        if "Content-Transfer-Encoding: quoted-printable" in body:
            try:
                import quopri
                # 提取 quoted-printable 内容
                pattern = r'Content-Transfer-Encoding: quoted-printable\s*\r?\n\s*\r?\n(.*?)(?:\r?\n--|\r?\n\r?\n|$)'
                match = re.search(pattern, body, re.MULTILINE | re.DOTALL)
                if match:
                    qp_content = match.group(1)
                    decoded_bytes = quopri.decodestring(qp_content.encode())
                    return decoded_bytes.decode('utf-8', errors='ignore').strip()
            except Exception as e:
                print(f"Quoted-printable解码失败: {e}")

        # 如果是普通文本或无法解码，返回清理后的原文本
        # 移除 MIME 头部信息，只保留可读内容
        if "--_" in body and "Content-Type:" in body:
            # 尝试提取纯文本部分
            lines = body.split('\n')
            content_lines = []
            in_content = False

            for line in lines:
                line = line.strip()
                if line.startswith('Content-Type: text/plain'):
                    in_content = True
                    continue
                elif line.startswith('Content-Transfer-Encoding:'):
                    continue
                elif line.startswith('--_') and in_content:
                    break
                elif in_content and line and not line.startswith('Content-'):
                    content_lines.append(line)

            if content_lines:
                return '\n'.join(content_lines)

        # 返回原始内容的前200个字符
        return body[:200] + "..." if len(body) > 200 else body
    
    def display_emails(self, emails_data):
        """显示邮件列表"""
        if not emails_data:
            print("没有找到邮件")
            return
        
        emails = emails_data.get('emails', [])
        total = emails_data.get('total', 0)
        page = emails_data.get('page', 1)
        limit = emails_data.get('limit', 20)
        
        print("=" * 80)
        print(f"NBEmail 收件箱查看器 - 第 {page} 页")
        print(f"总计: {total} 封邮件")
        print("=" * 80)
        
        if not emails:
            print("收件箱为空")
            return
        
        for i, email in enumerate(emails, 1):
            print(f"\n邮件 {i} (ID: {email.get('id', 'N/A')})")
            print(f"日期: {self.format_date(email.get('created_at', ''))}")
            print(f"发件人: {email.get('from', 'N/A')}")
            print(f"收件人: {email.get('to', 'N/A')}")
            print(f"主题: {self.clean_text(email.get('subject', '无主题'))}")
            
            # 显示邮件正文预览
            raw_body = email.get('body', '')
            if raw_body:
                # 先解码邮件内容
                decoded_body = self.decode_email_body(raw_body)
                # 再清理文本
                clean_body = self.clean_text(decoded_body)

                # 如果清理后的内容为空或者包含HTML标签，尝试HTML转换
                if not clean_body or '<html>' in clean_body.lower() or '<div>' in clean_body.lower() or '<p>' in clean_body.lower():
                    # 尝试将HTML转换为纯文本
                    text_from_html = self.html_to_text(decoded_body)
                    if text_from_html and len(text_from_html.strip()) > 0:
                        clean_body = text_from_html

                if clean_body:
                    preview = clean_body[:200] + "..." if len(clean_body) > 200 else clean_body
                    print(f"内容预览: {preview}")
                else:
                    print("内容预览: (解码失败，显示原始内容)")
                    preview = self.clean_text(raw_body)[:200] + "..."
                    print(f"原始内容: {preview}")
            else:
                print("内容预览: (无内容)")
            
            print("-" * 80)
    
    def show_inbox(self, email, password, page=1, limit=10):
        """显示收件箱"""
        if not self.login(email, password):
            return False

        print(f"\n正在获取收件箱邮件...")
        emails_data = self.get_emails(page=page, limit=limit, folder="inbox")

        if emails_data:
            self.display_emails(emails_data)

            # 如果收件箱为空，尝试获取所有邮件
            if emails_data.get('total', 0) == 0:
                print("\n收件箱为空，尝试获取所有邮件...")
                all_emails_data = self.get_all_emails(page=page, limit=limit)
                if all_emails_data:
                    self.display_emails(all_emails_data)

            return True
        else:
            print("获取邮件失败")
            return False

    def get_all_emails(self, page=1, limit=20):
        """获取所有邮件（不限制文件夹）"""
        params = {
            "page": page,
            "limit": limit
        }

        try:
            response = self.session.get(
                f"{self.base_url}/api/emails",
                params=params,
                timeout=10
            )

            if response.status_code == 200:
                result = response.json()
                if result.get('success'):
                    return result.get('data', {})
                else:
                    print(f"✗ 获取所有邮件失败: {result.get('message', '未知错误')}")
                    return None
            else:
                print(f"✗ 获取所有邮件请求失败: HTTP {response.status_code}")
                return None

        except requests.exceptions.RequestException as e:
            print(f"✗ 网络请求失败: {e}")
            return None

    def interactive_mode(self, email, password):
        """交互式模式"""
        if not self.login(email, password):
            return

        while True:
            print("\n" + "="*50)
            print("NBEmail 收件箱查看器 - 交互模式")
            print("="*50)
            print("1. 查看收件箱")
            print("2. 查看所有邮件")
            print("3. 搜索邮件")
            print("4. 查看邮件详情")
            print("5. 手动刷新")
            print("6. 自动刷新模式 (每15秒)")
            print("0. 退出")
            print("-"*50)

            choice = input("请选择操作 (0-6): ").strip()

            if choice == "0":
                print("再见！")
                break
            elif choice == "1":
                page = int(input("页码 (默认1): ") or "1")
                limit = int(input("每页数量 (默认10): ") or "10")
                emails_data = self.get_emails(page=page, limit=limit, folder="inbox")
                if emails_data:
                    self.display_emails(emails_data)
            elif choice == "2":
                page = int(input("页码 (默认1): ") or "1")
                limit = int(input("每页数量 (默认10): ") or "10")
                emails_data = self.get_all_emails(page=page, limit=limit)
                if emails_data:
                    self.display_emails(emails_data)
            elif choice == "3":
                search_term = input("搜索关键词: ").strip()
                if search_term:
                    emails_data = self.get_emails(search=search_term)
                    if emails_data:
                        self.display_emails(emails_data)
            elif choice == "4":
                email_id = input("邮件ID: ").strip()
                if email_id.isdigit():
                    email_detail = self.get_email_detail(int(email_id))
                    if email_detail:
                        self.display_email_detail(email_detail)
            elif choice == "5":
                print("手动刷新中...")
                continue
            elif choice == "6":
                print("启动自动刷新模式...")
                self.start_auto_refresh(email, password)
                print("自动刷新模式已退出，返回主菜单")
            else:
                print("无效选择，请重试")

    def display_email_detail(self, email):
        """显示邮件详情"""
        print("\n" + "="*80)
        print("邮件详情")
        print("="*80)
        print(f"ID: {email.get('id', 'N/A')}")
        print(f"日期: {self.format_date(email.get('created_at', ''))}")
        print(f"发件人: {email.get('from', 'N/A')}")
        print(f"收件人: {email.get('to', 'N/A')}")
        print(f"主题: {self.clean_text(email.get('subject', '无主题'))}")
        print(f"已读: {'是' if email.get('is_read') else '否'}")
        print("-"*80)
        print("邮件内容:")
        raw_body = email.get('body', '')
        if raw_body:
            decoded_body = self.decode_email_body(raw_body)
            clean_body = self.clean_text(decoded_body)

            # 如果清理后的内容为空或者包含HTML标签，尝试HTML转换
            if not clean_body or '<html>' in clean_body.lower() or '<div>' in clean_body.lower() or '<p>' in clean_body.lower():
                # 尝试将HTML转换为纯文本
                text_from_html = self.html_to_text(decoded_body)
                if text_from_html and len(text_from_html.strip()) > 0:
                    clean_body = text_from_html

            if clean_body:
                print(clean_body)
            else:
                print("(解码失败，显示原始内容)")
                print(self.clean_text(raw_body))
        else:
            print("无内容")
        print("="*80)

    def clear_screen(self):
        """清屏"""
        os.system('cls' if os.name == 'nt' else 'clear')

    def auto_refresh_emails(self, email, password, page=1, limit=10):
        """自动刷新邮件"""
        while not self.stop_refresh:
            try:
                self.clear_screen()
                print(f"🔄 自动刷新模式 - 每{self.refresh_interval}秒刷新")
                print(f"⏰ 当前时间: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
                print("按 Ctrl+C 退出自动刷新模式")
                print("=" * 80)

                # 获取邮件
                emails_data = self.get_emails(page=page, limit=limit, folder="inbox")
                if emails_data:
                    self.display_emails(emails_data)

                    # 如果收件箱为空，尝试获取所有邮件
                    if emails_data.get('total', 0) == 0:
                        print("\n收件箱为空，尝试获取所有邮件...")
                        all_emails_data = self.get_all_emails(page=page, limit=limit)
                        if all_emails_data:
                            self.display_emails(all_emails_data)
                else:
                    print("获取邮件失败")

                print(f"\n⏳ 下次刷新时间: {datetime.now().strftime('%H:%M:%S')} + {self.refresh_interval}秒")

                # 等待刷新间隔
                for i in range(self.refresh_interval):
                    if self.stop_refresh:
                        break
                    time.sleep(1)

            except KeyboardInterrupt:
                print("\n\n用户中断，退出自动刷新模式")
                self.stop_refresh = True
                break
            except Exception as e:
                print(f"\n自动刷新出错: {e}")
                time.sleep(5)  # 出错后等待5秒再重试

    def start_auto_refresh(self, email, password, page=1, limit=10):
        """启动自动刷新模式"""
        if not self.login(email, password):
            return False

        print(f"🚀 启动自动刷新模式，每{self.refresh_interval}秒刷新一次")
        print("按 Ctrl+C 可以随时退出")
        time.sleep(2)

        self.stop_refresh = False
        try:
            self.auto_refresh_emails(email, password, page, limit)
        except KeyboardInterrupt:
            print("\n\n👋 自动刷新已停止")
        finally:
            self.stop_refresh = True

        return True

def main():
    if len(sys.argv) < 2:
        print("用法: python inbox_viewer.py <邮箱地址>.<密码> [选项]")
        print("示例: python inbox_viewer.py 2014131458@qq.com.tgx123456")
        print("选项:")
        print("  -i, --interactive  交互模式")
        print("  -a, --auto-refresh 自动刷新模式 (每15秒)")
        print("  -p, --page N       指定页码")
        print("  -l, --limit N      每页邮件数量")
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

    # 解析选项
    interactive = False
    auto_refresh = False
    page = 1
    limit = 20

    for i, arg in enumerate(sys.argv[2:], 2):
        if arg in ['-i', '--interactive']:
            interactive = True
        elif arg in ['-a', '--auto-refresh']:
            auto_refresh = True
        elif arg in ['-p', '--page'] and i + 1 < len(sys.argv):
            try:
                page = int(sys.argv[i + 1])
            except ValueError:
                print("页码必须是数字")
                return
        elif arg in ['-l', '--limit'] and i + 1 < len(sys.argv):
            try:
                limit = int(sys.argv[i + 1])
            except ValueError:
                print("每页数量必须是数字")
                return

    # 创建查看器实例
    viewer = NBEmailViewer()

    if interactive:
        # 交互模式
        viewer.interactive_mode(email, password)
    elif auto_refresh:
        # 自动刷新模式
        viewer.start_auto_refresh(email, password, page=page, limit=limit)
    else:
        # 单次查看模式
        success = viewer.show_inbox(email, password, page=page, limit=limit)

        if success:
            print(f"\n✓ 收件箱查看完成")
        else:
            print(f"\n✗ 收件箱查看失败")

if __name__ == "__main__":
    main()
