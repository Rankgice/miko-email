#!/usr/bin/env python3
import sqlite3
import quopri
import re

def decode_quoted_printable(text):
    """解码quoted-printable编码的文本"""
    try:
        # 使用Python的quopri模块解码
        decoded_bytes = quopri.decodestring(text.encode())
        return decoded_bytes.decode('utf-8', errors='ignore')
    except Exception as e:
        print(f"解码失败: {e}")
        return text

def parse_multipart_email(body, boundary):
    """解析多部分邮件"""
    parts = body.split(f'--{boundary}')
    text_parts = []
    html_parts = []
    
    for part in parts:
        part = part.strip()
        if not part or part == '--':
            continue
            
        # 分离头部和内容
        if '\n\n' in part:
            headers, content = part.split('\n\n', 1)
        elif '\r\n\r\n' in part:
            headers, content = part.split('\r\n\r\n', 1)
        else:
            continue
            
        # 解析头部
        content_type = ''
        content_transfer_encoding = ''
        
        for line in headers.split('\n'):
            line = line.strip()
            if line.lower().startswith('content-type:'):
                content_type = line[13:].strip().lower()
            elif line.lower().startswith('content-transfer-encoding:'):
                content_transfer_encoding = line[26:].strip().lower()
        
        # 解码内容
        if content_transfer_encoding == 'quoted-printable':
            content = decode_quoted_printable(content)
        
        # 分类内容
        if 'text/plain' in content_type:
            text_parts.append(content)
        elif 'text/html' in content_type:
            html_parts.append(content)
    
    # 优先返回纯文本内容
    if text_parts:
        return '\n\n'.join(text_parts)
    elif html_parts:
        return '\n\n'.join(html_parts)
    else:
        return body

def main():
    # 连接数据库
    conn = sqlite3.connect('miko_email.db')
    cursor = conn.cursor()
    
    # 查询包含quoted-printable编码的邮件
    cursor.execute("""
        SELECT id, body 
        FROM emails 
        WHERE body LIKE '%=0A%' OR body LIKE '%quoted-printable%' OR body LIKE '%multipart%'
        ORDER BY id DESC
    """)
    
    emails = cursor.fetchall()
    processed_count = 0
    
    for email_id, body in emails:
        print(f"处理邮件 ID {email_id}...")
        
        new_body = body
        
        # 检测是否是多部分邮件
        if 'multipart' in body and 'boundary=' in body:
            # 提取boundary
            boundary_match = re.search(r'boundary="?([^"\s;]+)"?', body)
            if boundary_match:
                boundary = boundary_match.group(1)
                print(f"  检测到boundary: {boundary}")
                new_body = parse_multipart_email(body, boundary)
            else:
                new_body = body
        elif '=0A' in body or 'quoted-printable' in body:
            # 直接解码quoted-printable
            new_body = decode_quoted_printable(body)
        
        # 如果内容有变化，更新数据库
        if new_body != body:
            cursor.execute("UPDATE emails SET body = ? WHERE id = ?", (new_body, email_id))
            print(f"  邮件 {email_id} 已更新")
            processed_count += 1
        else:
            print(f"  邮件 {email_id} 无需更新")
    
    # 提交更改
    conn.commit()
    conn.close()
    
    print(f"\n总共处理了 {processed_count} 封邮件")

if __name__ == '__main__':
    main()
