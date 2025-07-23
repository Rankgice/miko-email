#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
NBEmail 收件箱查看器 - 图形界面版本
功能：登录/退出、邮件列表、查看邮件、自动刷新
"""

import tkinter as tk
from tkinter import ttk, messagebox, scrolledtext
import requests
import json
import threading
import time
from datetime import datetime
from html.parser import HTMLParser
import html
import random

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

class NBEmailGUI:
    def __init__(self, root):
        self.root = root
        self.root.title("🌟 NBEmail 炫彩收件箱查看器 🌟")
        self.root.geometry("1200x800")

        # 设置炫彩主题
        self.setup_theme()

        # 渐变色彩配置
        self.colors = {
            'primary': '#667eea',      # 主色调 - 紫蓝色
            'secondary': '#764ba2',    # 次色调 - 深紫色
            'accent': '#f093fb',       # 强调色 - 粉紫色
            'success': '#4facfe',      # 成功色 - 蓝色
            'warning': '#f6d365',      # 警告色 - 黄色
            'danger': '#ff6b6b',       # 危险色 - 红色
            'dark': '#2c3e50',         # 深色
            'light': '#ecf0f1',        # 浅色
            'gradient_start': '#667eea',
            'gradient_end': '#764ba2'
        }

        # 设置窗口背景渐变
        self.root.configure(bg='#1a1a2e')
        
        # 网络会话
        self.session = requests.Session()
        self.session.headers.update({
            'Content-Type': 'application/json; charset=utf-8',
            'User-Agent': 'NBEmail-GUI-Viewer/1.0'
        })
        
        # 状态变量
        self.logged_in = False
        self.current_user = ""
        self.base_url = "http://me.youddns.site:8080"
        self.emails = []
        self.selected_email = None
        
        # 自动刷新相关
        self.auto_refresh = True
        self.refresh_interval = 15  # 15秒
        self.countdown = 0
        self.refresh_thread = None
        self.stop_refresh = False
        
        self.setup_ui()
        self.start_refresh_timer()
        self.start_color_animation()

    def setup_theme(self):
        """设置炫彩主题"""
        style = ttk.Style()

        # 配置主题
        style.theme_use('clam')

        # 自定义样式
        style.configure('Title.TLabel',
                       font=('Arial', 16, 'bold'),
                       foreground='#667eea',
                       background='#1a1a2e')

        style.configure('Accent.TLabel',
                       font=('Arial', 10, 'bold'),
                       foreground='#f093fb',
                       background='#1a1a2e')

        style.configure('Success.TLabel',
                       font=('Arial', 10, 'bold'),
                       foreground='#4facfe',
                       background='#1a1a2e')

        style.configure('Warning.TLabel',
                       font=('Arial', 10, 'bold'),
                       foreground='#f6d365',
                       background='#1a1a2e')

        style.configure('Danger.TLabel',
                       font=('Arial', 10, 'bold'),
                       foreground='#ff6b6b',
                       background='#1a1a2e')

        # 按钮样式
        style.configure('Accent.TButton',
                       font=('Arial', 10, 'bold'),
                       foreground='white',
                       background='#667eea',
                       borderwidth=0,
                       focuscolor='none')

        style.map('Accent.TButton',
                 background=[('active', '#764ba2'),
                           ('pressed', '#f093fb')])

        # 框架样式
        style.configure('Card.TLabelFrame',
                       background='#16213e',
                       foreground='#667eea',
                       borderwidth=2,
                       relief='raised')

        style.configure('Card.TLabelFrame.Label',
                       background='#16213e',
                       foreground='#667eea',
                       font=('Arial', 12, 'bold'))

    def start_color_animation(self):
        """启动颜色动画效果"""
        def animate_colors():
            colors = ['#667eea', '#764ba2', '#f093fb', '#4facfe', '#f6d365']
            while not self.stop_refresh:
                for color in colors:
                    if self.stop_refresh:
                        break
                    try:
                        self.root.after(0, lambda c=color: self.update_accent_color(c))
                        time.sleep(2)
                    except:
                        break

        threading.Thread(target=animate_colors, daemon=True).start()

    def update_accent_color(self, color):
        """更新强调色"""
        try:
            style = ttk.Style()
            style.configure('Title.TLabel', foreground=color)
        except:
            pass
        
    def setup_ui(self):
        """设置用户界面"""
        # 主框架
        main_frame = tk.Frame(self.root, bg='#1a1a2e')
        main_frame.pack(fill=tk.BOTH, expand=True, padx=15, pady=15)

        # 标题
        title_label = ttk.Label(main_frame, text="🌟 NBEmail 炫彩收件箱查看器 🌟",
                               style='Title.TLabel')
        title_label.pack(pady=(0, 15))

        # 顶部登录框架
        login_frame = ttk.LabelFrame(main_frame, text="🔐 登录信息", padding=15)
        login_frame.pack(fill=tk.X, pady=(0, 15))
        # login_frame.configure(style='Card.TLabelFrame')
        
        # 登录控件
        login_inner = tk.Frame(login_frame, bg='#16213e')
        login_inner.pack(fill=tk.X)

        ttk.Label(login_inner, text="📧 邮箱:", style='Accent.TLabel').grid(row=0, column=0, sticky=tk.W, padx=(0, 8))
        self.email_var = tk.StringVar(value="3123717439@qq.com")
        self.email_entry = tk.Entry(login_inner, textvariable=self.email_var, width=25,
                                   font=('Arial', 10), bg='#2c3e50', fg='white',
                                   insertbackground='white', relief='flat', bd=5)
        self.email_entry.grid(row=0, column=1, padx=(0, 15), ipady=5)

        ttk.Label(login_inner, text="🔑 密码:", style='Accent.TLabel').grid(row=0, column=2, sticky=tk.W, padx=(0, 8))
        self.password_var = tk.StringVar(value="12345678")
        self.password_entry = tk.Entry(login_inner, textvariable=self.password_var, show="*", width=20,
                                      font=('Arial', 10), bg='#2c3e50', fg='white',
                                      insertbackground='white', relief='flat', bd=5)
        self.password_entry.grid(row=0, column=3, padx=(0, 15), ipady=5)

        self.login_btn = ttk.Button(login_inner, text="🚀 登录", command=self.login, style='Accent.TButton')
        self.login_btn.grid(row=0, column=4, padx=(0, 10), ipady=3)

        self.logout_btn = ttk.Button(login_inner, text="🚪 退出", command=self.logout,
                                   state=tk.DISABLED, style='Accent.TButton')
        self.logout_btn.grid(row=0, column=5, ipady=3)
        
        # 状态标签
        status_frame = tk.Frame(login_inner, bg='#16213e')
        status_frame.grid(row=1, column=0, columnspan=6, sticky=tk.EW, pady=(10, 0))

        self.status_var = tk.StringVar(value="❌ 未登录")
        self.status_label = ttk.Label(status_frame, textvariable=self.status_var, style='Danger.TLabel')
        self.status_label.pack(side=tk.LEFT)

        # 刷新倒计时标签
        self.countdown_var = tk.StringVar(value="")
        self.countdown_label = ttk.Label(status_frame, textvariable=self.countdown_var, style='Success.TLabel')
        self.countdown_label.pack(side=tk.RIGHT)
        
        # 中间内容框架
        content_frame = tk.Frame(main_frame, bg='#1a1a2e')
        content_frame.pack(fill=tk.BOTH, expand=True)

        # 左侧邮件列表
        left_frame = ttk.LabelFrame(content_frame, text="📬 邮件列表", padding=10)
        left_frame.pack(side=tk.LEFT, fill=tk.BOTH, expand=True, padx=(0, 10))
        # left_frame.configure(style='Card.TLabelFrame')
        
        # 邮件列表树形控件
        tree_frame = tk.Frame(left_frame, bg='#16213e')
        tree_frame.pack(fill=tk.BOTH, expand=True)

        columns = ("ID", "日期", "发件人", "主题", "状态")
        self.email_tree = ttk.Treeview(tree_frame, columns=columns, show="headings", height=18)

        # 设置列标题和宽度
        self.email_tree.heading("ID", text="🆔 ID")
        self.email_tree.heading("日期", text="📅 日期")
        self.email_tree.heading("发件人", text="👤 发件人")
        self.email_tree.heading("主题", text="📝 主题")
        self.email_tree.heading("状态", text="📊 状态")

        self.email_tree.column("ID", width=60)
        self.email_tree.column("日期", width=120)
        self.email_tree.column("发件人", width=180)
        self.email_tree.column("主题", width=250)
        self.email_tree.column("状态", width=80)

        # 配置树形控件样式
        style = ttk.Style()
        style.configure("Treeview",
                       background="#2c3e50",
                       foreground="white",
                       fieldbackground="#2c3e50",
                       font=('Arial', 9))
        style.configure("Treeview.Heading",
                       background="#667eea",
                       foreground="white",
                       font=('Arial', 10, 'bold'))
        style.map("Treeview",
                 background=[('selected', '#764ba2')])
        style.map("Treeview.Heading",
                 background=[('active', '#f093fb')])
        
        # 滚动条
        tree_scroll = ttk.Scrollbar(tree_frame, orient=tk.VERTICAL, command=self.email_tree.yview)
        self.email_tree.configure(yscrollcommand=tree_scroll.set)

        self.email_tree.pack(side=tk.LEFT, fill=tk.BOTH, expand=True)
        tree_scroll.pack(side=tk.RIGHT, fill=tk.Y)

        # 绑定选择事件
        self.email_tree.bind("<<TreeviewSelect>>", self.on_email_select)

        # 右侧邮件内容
        right_frame = ttk.LabelFrame(content_frame, text="📖 邮件内容", padding=10)
        right_frame.pack(side=tk.RIGHT, fill=tk.BOTH, expand=True)
        # right_frame.configure(style='Card.TLabelFrame')

        # 邮件详情显示
        detail_frame = tk.Frame(right_frame, bg='#16213e')
        detail_frame.pack(fill=tk.BOTH, expand=True)

        self.email_detail = scrolledtext.ScrolledText(detail_frame, wrap=tk.WORD, width=45, height=28,
                                                     bg='#2c3e50', fg='white',
                                                     insertbackground='white',
                                                     font=('Consolas', 10),
                                                     relief='flat', bd=0)
        self.email_detail.pack(fill=tk.BOTH, expand=True, padx=5, pady=5)
        
        # 底部按钮框架
        button_frame = tk.Frame(main_frame, bg='#1a1a2e')
        button_frame.pack(fill=tk.X, pady=(15, 0))

        # 创建炫彩按钮
        self.refresh_btn = ttk.Button(button_frame, text="🔄 立即刷新",
                                    command=self.manual_refresh, state=tk.DISABLED,
                                    style='Accent.TButton')
        self.refresh_btn.pack(side=tk.LEFT, padx=(0, 15))

        # 自动刷新开关
        self.auto_refresh_var = tk.BooleanVar(value=True)
        refresh_frame = tk.Frame(button_frame, bg='#1a1a2e')
        refresh_frame.pack(side=tk.LEFT)

        self.auto_refresh_cb = tk.Checkbutton(refresh_frame, text="⚡ 自动刷新(15秒)",
                                            variable=self.auto_refresh_var,
                                            command=self.toggle_auto_refresh,
                                            bg='#1a1a2e', fg='#667eea',
                                            selectcolor='#764ba2',
                                            activebackground='#1a1a2e',
                                            activeforeground='#f093fb',
                                            font=('Arial', 10, 'bold'))
        self.auto_refresh_cb.pack()

        # 添加装饰性元素
        decoration_frame = tk.Frame(button_frame, bg='#1a1a2e')
        decoration_frame.pack(side=tk.RIGHT)

        ttk.Label(decoration_frame, text="✨ 炫彩邮箱 ✨", style='Accent.TLabel').pack()
        
    def html_to_text(self, html_content):
        """将 HTML 内容转换为纯文本"""
        if not html_content or not html_content.strip():
            return ""
            
        try:
            extractor = HTMLTextExtractor()
            extractor.feed(html_content)
            text = extractor.get_text()
            
            lines = [line.strip() for line in text.split('\n')]
            lines = [line for line in lines if line]
            
            return '\n'.join(lines)
        except Exception as e:
            import re
            text = re.sub(r'<[^>]+>', '', html_content)
            text = html.unescape(text)
            return text.strip()
    
    def login(self):
        """登录"""
        email = self.email_var.get().strip()
        password = self.password_var.get().strip()
        
        if not email or not password:
            messagebox.showerror("错误", "请输入邮箱和密码")
            return
            
        try:
            login_data = {"email": email, "password": password}
            response = self.session.post(f"{self.base_url}/api/login", json=login_data, timeout=10)
            
            if response.status_code == 200:
                result = response.json()
                if result.get('success'):
                    self.logged_in = True
                    self.current_user = email
                    self.status_var.set(f"✅ 已登录: {email}")
                    self.status_label.config(style='Success.TLabel')
                    
                    # 更新按钮状态
                    self.login_btn.config(state=tk.DISABLED)
                    self.logout_btn.config(state=tk.NORMAL)
                    self.refresh_btn.config(state=tk.NORMAL)
                    self.email_entry.config(state=tk.DISABLED)
                    self.password_entry.config(state=tk.DISABLED)
                    
                    # 立即刷新邮件
                    self.load_emails()
                    messagebox.showinfo("🎉 成功", "登录成功！欢迎使用炫彩邮箱！")
                else:
                    messagebox.showerror("登录失败", result.get('message', '未知错误'))
            else:
                messagebox.showerror("登录失败", f"HTTP {response.status_code}")
                
        except requests.exceptions.RequestException as e:
            messagebox.showerror("网络错误", f"登录失败: {e}")
    
    def logout(self):
        """退出登录"""
        try:
            self.session.post(f"{self.base_url}/api/logout", timeout=5)
        except:
            pass
            
        self.logged_in = False
        self.current_user = ""
        self.emails = []
        
        # 更新界面
        self.status_var.set("❌ 未登录")
        self.status_label.config(style='Danger.TLabel')
        
        # 更新按钮状态
        self.login_btn.config(state=tk.NORMAL)
        self.logout_btn.config(state=tk.DISABLED)
        self.refresh_btn.config(state=tk.DISABLED)
        self.email_entry.config(state=tk.NORMAL)
        self.password_entry.config(state=tk.NORMAL)
        
        # 清空显示
        self.email_tree.delete(*self.email_tree.get_children())
        self.email_detail.delete(1.0, tk.END)
        
        messagebox.showinfo("👋 提示", "已退出登录，感谢使用炫彩邮箱！")

    def load_emails(self):
        """加载邮件列表"""
        if not self.logged_in:
            return

        try:
            response = self.session.get(f"{self.base_url}/api/emails?folder=inbox&page=1&limit=50", timeout=10)

            if response.status_code == 200:
                result = response.json()
                if result.get('success'):
                    data = result.get('data', {})
                    self.emails = data.get('emails', [])
                    self.update_email_list()
                else:
                    print(f"获取邮件失败: {result.get('message')}")
            else:
                print(f"获取邮件请求失败: HTTP {response.status_code}")

        except requests.exceptions.RequestException as e:
            print(f"网络错误: {e}")

    def update_email_list(self):
        """更新邮件列表显示"""
        # 清空现有项目
        self.email_tree.delete(*self.email_tree.get_children())

        for email in self.emails:
            email_id = email.get('id', '')
            created_at = email.get('created_at', '')
            from_addr = email.get('from', '')
            subject = email.get('subject', '无主题')
            is_read = email.get('is_read', False)

            # 格式化日期
            try:
                if created_at:
                    dt = datetime.fromisoformat(created_at.replace('Z', '+00:00'))
                    date_str = dt.strftime('%m-%d %H:%M')
                else:
                    date_str = ''
            except:
                date_str = created_at[:16] if len(created_at) > 16 else created_at

            # 状态
            status = "✅ 已读" if is_read else "🔥 未读"

            # 插入到树形控件
            item = self.email_tree.insert("", tk.END, values=(email_id, date_str, from_addr, subject, status))

            # 未读邮件用特殊标记显示
            if not is_read:
                self.email_tree.set(item, "主题", f"🌟 {subject}")
                # 可以在这里添加更多未读邮件的视觉效果

    def on_email_select(self, event):
        """邮件选择事件"""
        selection = self.email_tree.selection()
        if not selection:
            return

        item = selection[0]
        values = self.email_tree.item(item, 'values')
        if not values:
            return

        email_id = values[0]

        # 查找对应的邮件
        self.selected_email = None
        for email in self.emails:
            if str(email.get('id')) == str(email_id):
                self.selected_email = email
                break

        if self.selected_email:
            self.show_email_detail()

    def show_email_detail(self):
        """显示邮件详情"""
        if not self.selected_email:
            return

        # 清空内容
        self.email_detail.delete(1.0, tk.END)

        # 显示邮件头部信息
        email_info = f"""🌟 ═══════════════════════════════════════════════════════════ 🌟
📧 邮件详情
🌟 ═══════════════════════════════════════════════════════════ 🌟

🆔 邮件ID: {self.selected_email.get('id', 'N/A')}
📅 日期: {self.format_date(self.selected_email.get('created_at', ''))}
👤 发件人: {self.selected_email.get('from', 'N/A')}
📮 收件人: {self.selected_email.get('to', 'N/A')}
📝 主题: {self.selected_email.get('subject', '无主题')}
📊 状态: {'✅ 已读' if self.selected_email.get('is_read') else '🔥 未读'}

✨ ═══════════════════════════════════════════════════════════ ✨
📖 邮件内容
✨ ═══════════════════════════════════════════════════════════ ✨

"""
        self.email_detail.insert(tk.END, email_info)

        # 显示邮件正文
        raw_body = self.selected_email.get('body', '')
        if raw_body:
            # 检查是否包含HTML标签，如果是则转换为纯文本
            if '<html>' in raw_body.lower() or '<div>' in raw_body.lower() or '<p>' in raw_body.lower():
                clean_body = self.html_to_text(raw_body)
            else:
                clean_body = raw_body

            if clean_body:
                self.email_detail.insert(tk.END, clean_body)
                self.email_detail.insert(tk.END, f"\n\n🌟 ═══════════════════════════════════════════════════════════ 🌟\n✨ 邮件内容结束 ✨\n🌟 ═══════════════════════════════════════════════════════════ 🌟")
            else:
                self.email_detail.insert(tk.END, "❌ (无法解析邮件内容)")
        else:
            self.email_detail.insert(tk.END, "📭 (无邮件内容)")

    def format_date(self, date_str):
        """格式化日期"""
        if not date_str:
            return "未知时间"

        try:
            dt = datetime.fromisoformat(date_str.replace('Z', '+00:00'))
            return dt.strftime('%Y-%m-%d %H:%M:%S')
        except:
            return date_str

    def manual_refresh(self):
        """手动刷新"""
        if self.logged_in:
            self.load_emails()
            self.countdown = self.refresh_interval  # 重置倒计时

    def toggle_auto_refresh(self):
        """切换自动刷新"""
        self.auto_refresh = self.auto_refresh_var.get()
        if self.auto_refresh:
            self.countdown = self.refresh_interval

    def start_refresh_timer(self):
        """启动刷新计时器"""
        def refresh_worker():
            while True:
                if self.stop_refresh:
                    break

                if self.auto_refresh and self.logged_in:
                    if self.countdown <= 0:
                        # 刷新邮件
                        self.root.after(0, self.load_emails)
                        self.countdown = self.refresh_interval
                    else:
                        self.countdown -= 1

                    # 更新倒计时显示
                    if self.logged_in:
                        if self.countdown <= 5:
                            countdown_text = f"🔥 下次刷新: {self.countdown}秒"
                        elif self.countdown <= 10:
                            countdown_text = f"⚡ 下次刷新: {self.countdown}秒"
                        else:
                            countdown_text = f"⏰ 下次刷新: {self.countdown}秒"
                    else:
                        countdown_text = ""
                    self.root.after(0, lambda: self.countdown_var.set(countdown_text))
                else:
                    self.root.after(0, lambda: self.countdown_var.set(""))

                time.sleep(1)

        self.countdown = self.refresh_interval
        self.refresh_thread = threading.Thread(target=refresh_worker, daemon=True)
        self.refresh_thread.start()

    def on_closing(self):
        """窗口关闭事件"""
        self.stop_refresh = True
        if self.logged_in:
            try:
                self.session.post(f"{self.base_url}/api/logout", timeout=2)
            except:
                pass
        self.root.destroy()

def main():
    root = tk.Tk()
    app = NBEmailGUI(root)

    # 绑定窗口关闭事件
    root.protocol("WM_DELETE_WINDOW", app.on_closing)

    # 启动主循环
    root.mainloop()

if __name__ == "__main__":
    main()
