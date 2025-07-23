#!/usr/bin/env python3
import sqlite3
import quopri
import re
import base64

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
    
    for i, part in enumerate(parts):
        part = part.strip()
        if not part or part == '--':
            continue
            
        print(f"处理第 {i} 部分")
        
        # 分离头部和内容
        if '\n\n' in part:
            headers, content = part.split('\n\n', 1)
        elif '\r\n\r\n' in part:
            headers, content = part.split('\r\n\r\n', 1)
        else:
            continue
            
        # 解析头部
        content_type = ""
        transfer_encoding = ""
        charset = ""
        
        for line in headers.split('\n'):
            line = line.strip()
            if line.lower().startswith('content-type:'):
                content_type = line[13:].strip().lower()
                # 提取charset
                if 'charset=' in content_type:
                    charset_match = re.search(r'charset=([^;\s]+)', content_type)
                    if charset_match:
                        charset = charset_match.group(1).strip('"\'')
            elif line.lower().startswith('content-transfer-encoding:'):
                transfer_encoding = line[26:].strip().lower()
        
        print(f"Content-Type: {content_type}")
        print(f"Transfer-Encoding: {transfer_encoding}")
        print(f"Charset: {charset}")
        
        # 解码内容
        if transfer_encoding == 'quoted-printable':
            content = decode_quoted_printable(content)
        elif transfer_encoding == 'base64':
            try:
                content = base64.b64decode(content.replace('\n', '').replace('\r', '')).decode('utf-8', errors='ignore')
            except:
                pass
        
        # 分类存储
        if 'text/plain' in content_type:
            text_parts.append(content)
        elif 'text/html' in content_type:
            html_parts.append(content)
    
    # 优先返回纯文本内容
    if text_parts:
        result = '\n\n'.join(text_parts)
        print(f"返回纯文本内容，长度: {len(result)}")
        return result
    elif html_parts:
        result = '\n\n'.join(html_parts)
        print(f"返回HTML内容，长度: {len(result)}")
        return result
    
    return body

def extract_boundary(body):
    """从邮件中提取boundary"""
    # 查找boundary=
    boundary_match = re.search(r'boundary=([^;\s\n\r]+)', body, re.IGNORECASE)
    if boundary_match:
        boundary = boundary_match.group(1).strip('"\'')
        print(f"找到boundary: {boundary}")
        return boundary
    
    # 自动检测boundary
    lines = body.split('\n')
    for line in lines:
        line = line.strip()
        if line.startswith('--') and len(line) > 10:
            if '_' in line or '=' in line:
                boundary = line[2:]  # 去掉前面的--
                print(f"自动检测到boundary: {boundary}")
                return boundary
    
    return None

def fix_steam_email(body):
    """修复Steam邮件编码"""
    print(f"原始内容长度: {len(body)}")
    
    # 检查是否是multipart邮件
    if 'multipart/alternative' in body or 'boundary=' in body:
        print("检测到multipart邮件")
        boundary = extract_boundary(body)
        if boundary:
            return parse_multipart_email(body, boundary)
    
    # 检查是否是quoted-printable编码
    if '=' in body and ('=0A' in body or '=0D' in body or '=20' in body):
        print("检测到quoted-printable编码")
        return decode_quoted_printable(body)
    
    # 检查是否是base64编码
    clean_body = body.replace('\n', '').replace('\r', '').strip()
    if len(clean_body) > 0 and len(clean_body) % 4 == 0:
        try:
            decoded = base64.b64decode(clean_body).decode('utf-8', errors='ignore')
            if decoded and len(decoded) > 10:  # 确保解码结果有意义
                print("检测到base64编码")
                return decoded
        except:
            pass
    
    return body

def main():
    # 连接数据库
    conn = sqlite3.connect('miko_email.db')
    cursor = conn.cursor()
    
    # 查找Steam邮件
    cursor.execute("""
        SELECT id, subject, body, sender 
        FROM emails 
        WHERE sender LIKE '%steampowered%' 
        ORDER BY created_at DESC
    """)
    
    emails = cursor.fetchall()
    fixed_count = 0
    
    for email_id, subject, body, sender in emails:
        print(f"\n处理邮件 ID: {email_id}")
        print(f"发件人: {sender}")
        print(f"原始主题: {subject}")
        
        # 修复邮件内容
        fixed_body = fix_steam_email(body)
        
        # 如果内容有变化，更新数据库
        if fixed_body != body:
            cursor.execute("""
                UPDATE emails 
                SET body = ? 
                WHERE id = ?
            """, (fixed_body, email_id))
            
            print(f"✅ 邮件 {email_id} 修复成功")
            print(f"修复后内容预览: {fixed_body[:200]}...")
            fixed_count += 1
        else:
            print(f"邮件 {email_id} 无需修复")
        
        print("-" * 50)
    
    # 提交更改
    conn.commit()
    conn.close()
    
    print(f"\n总共修复了 {fixed_count} 封邮件")

if __name__ == "__main__":
    main()
