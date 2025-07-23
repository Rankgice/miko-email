#!/usr/bin/env python3
"""
测试Steam邮件解码功能
"""
import quopri
import re

def decode_quoted_printable(text):
    """解码quoted-printable编码的文本"""
    try:
        # 使用Python的quopri模块解码
        decoded_bytes = quopri.decodestring(text.encode())
        return decoded_bytes.decode('utf-8', errors='ignore')
    except Exception as e:
        print(f"quoted-printable解码失败: {e}")
        return text

def extract_boundary(body):
    """从邮件中提取boundary"""
    # 查找boundary=
    boundary_match = re.search(r'boundary="([^"]+)"', body, re.IGNORECASE)
    if boundary_match:
        boundary = boundary_match.group(1)
        print(f"找到boundary: {boundary}")
        return boundary
    
    # 自动检测boundary
    lines = body.split('\n')
    for line in lines:
        line = line.strip()
        if line.startswith('--') and len(line) > 10:
            # Steam邮件的boundary特征
            if '_' in line or '=' in line or 'Boundary' in line or len(line) > 15:
                boundary = line[2:]  # 去掉前面的--
                print(f"自动检测到boundary: {boundary}")
                return boundary
    
    return None

def parse_multipart_email(body, boundary):
    """解析多部分邮件"""
    parts = body.split(f'--{boundary}')
    text_parts = []
    html_parts = []
    
    for i, part in enumerate(parts):
        part = part.strip()
        if not part or part == '--':
            continue
            
        print(f"\n处理第 {i} 部分")
        print(f"部分长度: {len(part)}")
        
        # 分离头部和内容
        if '\n\n' in part:
            headers, content = part.split('\n\n', 1)
        elif '\r\n\r\n' in part:
            headers, content = part.split('\r\n\r\n', 1)
        else:
            print("  无法分离头部和内容")
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
        
        print(f"  Content-Type: {content_type}")
        print(f"  Transfer-Encoding: {transfer_encoding}")
        print(f"  Charset: {charset}")
        print(f"  原始内容长度: {len(content)}")
        
        # 解码内容
        if transfer_encoding == 'quoted-printable':
            print(f"  解码quoted-printable...")
            decoded_content = decode_quoted_printable(content)
            print(f"  解码后长度: {len(decoded_content)}")
            content = decoded_content
        
        # 分类存储
        if 'text/plain' in content_type:
            text_parts.append(content)
            print(f"  添加到纯文本部分")
            print(f"  纯文本预览: {content[:200]}...")
        elif 'text/html' in content_type:
            html_parts.append(content)
            print(f"  添加到HTML部分")
            print(f"  HTML预览: {content[:200]}...")
    
    # 优先返回纯文本内容
    if text_parts:
        result = '\n\n'.join(text_parts)
        print(f"\n✅ 返回纯文本内容，总长度: {len(result)}")
        return result, "text"
    elif html_parts:
        result = '\n\n'.join(html_parts)
        print(f"\n✅ 返回HTML内容，总长度: {len(result)}")
        return result, "html"
    
    return body, "raw"

def test_steam_email():
    """测试Steam邮件解码"""
    print("🧪 测试Steam邮件解码功能")
    print("=" * 60)
    
    # 读取Steam邮件文件
    try:
        with open('steam邮件.txt', 'r', encoding='utf-8') as f:
            email_content = f.read()
        print(f"✅ 成功读取邮件文件，长度: {len(email_content)}")
    except Exception as e:
        print(f"❌ 读取邮件文件失败: {e}")
        return
    
    # 提取boundary
    print(f"\n🔍 提取boundary...")
    boundary = extract_boundary(email_content)
    
    if not boundary:
        print("❌ 未找到boundary")
        return
    
    # 解析multipart邮件
    print(f"\n📧 解析multipart邮件...")
    decoded_content, content_type = parse_multipart_email(email_content, boundary)
    
    # 显示解码结果
    print(f"\n📋 解码结果:")
    print(f"内容类型: {content_type}")
    print(f"内容长度: {len(decoded_content)}")
    print(f"\n内容预览:")
    print("-" * 40)
    print(decoded_content[:500])
    print("-" * 40)
    
    # 保存解码结果
    try:
        with open('steam邮件_测试解码结果.txt', 'w', encoding='utf-8') as f:
            f.write(f"解码类型: {content_type}\n")
            f.write(f"内容长度: {len(decoded_content)}\n")
            f.write("=" * 60 + "\n")
            f.write(decoded_content)
        print(f"\n✅ 解码结果已保存到 steam邮件_测试解码结果.txt")
    except Exception as e:
        print(f"❌ 保存解码结果失败: {e}")

def test_individual_parts():
    """测试单独的quoted-printable解码"""
    print("\n🧪 测试单独的quoted-printable解码")
    print("=" * 60)
    
    # 测试纯文本部分的一小段
    test_text = "=0A=0A=E6=82=A8=E5=A5=BD=EF=BC=9A"
    print(f"原始文本: {test_text}")
    
    decoded = decode_quoted_printable(test_text)
    print(f"解码结果: {decoded}")
    print(f"解码结果(repr): {repr(decoded)}")

if __name__ == "__main__":
    test_steam_email()
    test_individual_parts()
