import imaplib
import email
import os
import time
import logging
import re
import email.header
import ssl
from email.parser import BytesParser
from email.policy import default
from datetime import datetime

# 配置日志
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler("mail_receiver.log", encoding='utf-8'),
        logging.StreamHandler()
    ]
)

class EmailReceiver:
    def __init__(self, username, email_address, password):
        """初始化IMAP服务器连接参数"""
        # 使用本地Miko邮箱IMAP服务器
        self.imap_server = "localhost"
        self.imap_port = 143  # 非SSL端口
        self.imap_ssl_port = 993  # SSL端口（备用）
        self.download_folder = "downloaded_emails"

        # 用户凭据
        self.username = username
        self.email_address = email_address
        self.password = password

        # 打印用户信息用于调试
        logging.info(f"初始化用户: {self.username}")
        logging.info(f"邮箱地址: {self.email_address}")
        logging.info(f"密码: {self.password}")

        # 轮询间隔（秒）
        self.poll_interval = 10

        # 已处理的邮件ID集合，用于避免重复处理
        self.processed_email_ids = set()

        # 运行状态
        self.running = True

        # 确保下载文件夹存在
        if not os.path.exists(self.download_folder):
            os.makedirs(self.download_folder)
    
    def connect_to_server(self):
        """连接到IMAP服务器，使用标准IMAP协议"""
        # 尝试不同的连接方式
        methods = [
            {"type": "非SSL", "ssl": False, "port": self.imap_port},
            {"type": "SSL", "ssl": True, "port": self.imap_ssl_port}
        ]

        for method in methods:
            try:
                logging.info(f"[{self.username}] 尝试使用{method['type']}连接到 {self.imap_server}:{method['port']}...")
                logging.info(f"[{self.username}] 使用密码: {self.password}")

                if method['ssl']:
                    # 创建SSL上下文
                    context = ssl.create_default_context()
                    # 如果服务器证书有问题，可以禁用证书验证
                    context.check_hostname = False
                    context.verify_mode = ssl.CERT_NONE
                    mail = imaplib.IMAP4_SSL(self.imap_server, method['port'], ssl_context=context)
                else:
                    mail = imaplib.IMAP4(self.imap_server, method['port'])

                # 尝试不同的认证方式（适配Miko邮箱系统）
                auth_methods = [
                    {"name": "邮箱登录", "func": lambda: mail.login(self.email_address, self.password)},
                    {"name": "用户名登录", "func": lambda: mail.login(self.username, self.password)},
                ]

                for auth in auth_methods:
                    try:
                        logging.info(f"[{self.username}] 尝试 {auth['name']}...")
                        auth['func']()
                        logging.info(f"[{self.username}] 成功登录 (使用{method['type']}连接和{auth['name']})")
                        return mail
                    except Exception as e:
                        logging.warning(f"[{self.username}] {auth['name']}失败: {str(e)}")

            except Exception as e:
                logging.warning(f"[{self.username}] {method['type']}连接失败: {str(e)}")

        logging.error(f"[{self.username}] 所有连接和认证方式都失败")
        return None

    def get_mailbox_id(self):
        """从数据库获取邮箱ID"""
        try:
            conn = sqlite3.connect(self.db_path)
            cursor = conn.cursor()

            cursor.execute("SELECT id FROM mailboxes WHERE email = ?", (self.email_address,))
            result = cursor.fetchone()

            conn.close()

            if result:
                self.mailbox_id = result[0]
                logging.info(f"[{self.username}] 邮箱ID: {self.mailbox_id}")
                return True
            else:
                logging.error(f"[{self.username}] 未找到邮箱: {self.email_address}")
                return False

        except Exception as e:
            logging.error(f"[{self.username}] 获取邮箱ID失败: {str(e)}")
            return False


    
    def check_new_emails(self, mail):
        """检查新邮件"""
        try:
            # 选择收件箱
            mail.select('INBOX')

            # 搜索所有邮件
            status, messages = mail.search(None, 'ALL')

            if status != 'OK':
                logging.error(f"[{self.username}] 无法搜索邮件")
                return []

            # 获取邮件ID列表
            email_ids = messages[0].split()

            # 找出新邮件（未处理过的邮件）
            new_email_ids = []
            for email_id in email_ids:
                email_id_str = email_id.decode('utf-8')
                if email_id_str not in self.processed_email_ids:
                    new_email_ids.append(email_id)
                    self.processed_email_ids.add(email_id_str)

            if new_email_ids:
                logging.info(f"[{self.username}] 发现 {len(new_email_ids)} 封新邮件")
            else:
                logging.info(f"[{self.username}] 没有新邮件")

            return new_email_ids

        except Exception as e:
            logging.error(f"[{self.username}] 检查新邮件失败: {str(e)}")
            return []
    
    def fetch_email(self, mail, email_id):
        """获取单个邮件内容"""
        try:
            status, data = mail.fetch(email_id, '(RFC822)')

            if status != 'OK':
                logging.error(f"[{self.email_address}] 无法获取邮件 ID: {email_id}")
                return None

            # 解析邮件内容
            raw_email = data[0][1]
            email_message = email.message_from_bytes(raw_email)

            # 获取邮件信息
            subject = self.decode_header(email_message['Subject'])
            from_address = self.decode_header(email_message['From'])
            date_str = email_message['Date']

            # 获取邮件正文
            body = self.get_email_body(email_message)

            # 保存附件
            attachments = self.save_attachments(email_message, email_id)

            email_data = {
                'id': email_id,
                'subject': subject,
                'from': from_address,
                'date': date_str,
                'body': body,
                'attachments': attachments
            }

            logging.info(f"[{self.email_address}] 已获取邮件: {subject}")
            return email_data

        except Exception as e:
            logging.error(f"[{self.email_address}] 获取邮件失败: {str(e)}")
            return None

    def decode_header(self, header):
        """解码邮件头信息"""
        if header is None:
            return ""

        try:
            decoded_header = email.header.decode_header(header)
            header_parts = []

            for part, encoding in decoded_header:
                if isinstance(part, bytes):
                    # 处理编码
                    if encoding and encoding.lower() not in ['unknown-8bit', 'unknown']:
                        try:
                            header_parts.append(part.decode(encoding))
                            continue
                        except (UnicodeDecodeError, LookupError):
                            pass

                    # 尝试多种编码方式
                    for enc in ['utf-8', 'gbk', 'gb2312', 'big5', 'latin1']:
                        try:
                            decoded_part = part.decode(enc)
                            header_parts.append(decoded_part)
                            break
                        except UnicodeDecodeError:
                            continue
                    else:
                        # 如果所有编码都失败，使用错误处理
                        header_parts.append(part.decode('utf-8', errors='replace'))
                else:
                    header_parts.append(str(part))

            return " ".join(header_parts)
        except Exception as e:
            logging.warning(f"[{self.email_address}] 解码邮件头失败: {str(e)}")
            # 如果解码完全失败，尝试直接返回原始字符串
            try:
                return str(header)
            except:
                return "邮件头解码失败"
    
    def get_email_body(self, email_message):
        """获取邮件正文内容"""
        body = ""
        
        if email_message.is_multipart():
            # 如果邮件包含多个部分，遍历所有部分
            for part in email_message.walk():
                content_type = part.get_content_type()
                content_disposition = str(part.get("Content-Disposition"))
                
                # 跳过附件
                if "attachment" in content_disposition:
                    continue
                
                # 获取文本内容
                if content_type == "text/plain" and "attachment" not in content_disposition:
                    try:
                        payload = part.get_payload(decode=True)
                        if payload:
                            body += self.decode_payload(payload, part.get_content_charset())
                    except Exception as e:
                        logging.error(f"[{self.email_address}] 解析邮件正文失败: {str(e)}")

                # 如果没有纯文本，尝试获取HTML内容
                elif content_type == "text/html" and not body and "attachment" not in content_disposition:
                    try:
                        payload = part.get_payload(decode=True)
                        if payload:
                            body += self.decode_payload(payload, part.get_content_charset())
                    except Exception as e:
                        logging.error(f"[{self.email_address}] 解析HTML邮件正文失败: {str(e)}")
        else:
            # 如果邮件不是多部分的
            content_type = email_message.get_content_type()
            
            if content_type == "text/plain" or content_type == "text/html":
                try:
                    payload = email_message.get_payload(decode=True)
                    if payload:
                        body = self.decode_payload(payload, email_message.get_content_charset())
                except Exception as e:
                    logging.error(f"[{self.email_address}] 解析单部分邮件正文失败: {str(e)}")
        
        return body

    def decode_payload(self, payload, charset):
        """解码邮件载荷"""
        if not payload:
            return ""

        if isinstance(payload, str):
            return payload

        # 尝试使用指定的字符集
        if charset and charset.lower() not in ['unknown-8bit', 'unknown']:
            try:
                return payload.decode(charset)
            except (UnicodeDecodeError, LookupError):
                pass

        # 尝试多种编码方式
        for encoding in ['utf-8', 'gbk', 'gb2312', 'big5', 'latin1']:
            try:
                return payload.decode(encoding)
            except UnicodeDecodeError:
                continue

        # 如果所有编码都失败，使用错误处理
        return payload.decode('utf-8', errors='replace')

    def save_attachments(self, email_message, email_id):
        """保存邮件附件"""
        saved_attachments = []
        
        for part in email_message.walk():
            if part.get_content_maintype() == 'multipart':
                continue
            
            if part.get('Content-Disposition') is None:
                continue
            
            filename = part.get_filename()
            
            if filename:
                # 解码文件名
                filename = self.decode_header(filename)
                
                # 创建一个唯一的文件名，避免覆盖
                timestamp = datetime.now().strftime("%Y%m%d%H%M%S")
                unique_filename = f"{timestamp}_{self.email_address}_{filename}"
                filepath = os.path.join(self.download_folder, unique_filename)
                
                # 保存附件
                with open(filepath, 'wb') as f:
                    f.write(part.get_payload(decode=True))
                
                logging.info(f"[{self.email_address}] 已保存附件: {filepath}")
                saved_attachments.append(filepath)
        
        return saved_attachments
    
    def extract_verification_code(self, email_body):
        """从邮件正文中提取验证码"""
        # 尝试匹配6位数字验证码
        match = re.search(r'\b\d{6}\b', email_body)
        if match:
            return match.group(0)
        
        # 尝试匹配4位数字验证码
        match = re.search(r'\b\d{4}\b', email_body)
        if match:
            return match.group(0)
        
        return None
    
    def process_email(self, email_data):
        """处理单个邮件"""
        if not email_data:
            return
        
        logging.info("-" * 50)
        logging.info(f"[{self.email_address}] 处理邮件: {email_data['subject']}")
        logging.info(f"[{self.email_address}] 发件人: {email_data['from']}")
        logging.info(f"[{self.email_address}] 日期: {email_data['date']}")
        
        # 提取验证码
        verification_code = self.extract_verification_code(email_data['body'])
        if verification_code:
            logging.info(f"[{self.email_address}] 找到验证码: {verification_code}")
        
        # 如果有附件
        if email_data['attachments']:
            logging.info(f"[{self.email_address}] 附件数量: {len(email_data['attachments'])}")
            for attachment in email_data['attachments']:
                logging.info(f"[{self.email_address}]   - {attachment}")
        
        logging.info("-" * 50)
    
    def stop(self):
        """停止轮询"""
        self.running = False
    
    def run_polling(self):
        """运行邮件轮询"""
        logging.info(f"[{self.email_address}] 开始处理邮箱")

        # 连接到服务器
        mail = self.connect_to_server()
        if not mail:
            logging.error(f"[{self.email_address}] 无法连接到邮件服务器，退出轮询")
            return

        logging.info(f"[{self.email_address}] 开始邮件轮询，每 {self.poll_interval} 秒检查一次新邮件...")

        try:
            # 首次获取所有邮件ID，标记为已处理
            self.check_new_emails(mail)

            while self.running:
                try:
                    # 检查新邮件
                    new_email_ids = self.check_new_emails(mail)

                    # 处理每封新邮件
                    for email_id in new_email_ids:
                        email_data = self.fetch_email(mail, email_id)
                        if email_data:
                            self.process_email(email_data)

                    # 等待下一次轮询
                    logging.info(f"[{self.email_address}] 等待 {self.poll_interval} 秒后再次检查...")
                    time.sleep(self.poll_interval)

                except imaplib.IMAP4.abort:
                    logging.warning(f"[{self.email_address}] 连接中断，尝试重新连接...")
                    # 重新连接
                    try:
                        mail.logout()
                    except:
                        pass

                    mail = self.connect_to_server()
                    if not mail:
                        logging.error(f"[{self.email_address}] 重新连接失败，退出轮询")
                        break

                except Exception as e:
                    logging.error(f"[{self.email_address}] 轮询过程中出错: {str(e)}")
                    time.sleep(self.poll_interval)

        except KeyboardInterrupt:
            logging.info(f"[{self.email_address}] 用户中断，停止轮询")

        finally:
            # 关闭连接
            try:
                mail.close()
                mail.logout()
            except:
                pass

            logging.info(f"[{self.email_address}] 邮件轮询已停止")


def read_email_list(file_path):
    """读取文件中的邮箱列表，支持CSV和TXT格式"""
    accounts = []
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            # 读取文件内容
            content = f.read().strip()
            lines = content.split('\n')
            
            for line in lines:
                if ',' in line:
                    # 按逗号分割，格式: 用户名,邮箱地址,密码
                    parts = line.split(',')
                    if len(parts) >= 3:
                        username = parts[0].strip()
                        email = parts[1].strip()
                        password = parts[2].strip()
                        
                        # 确保邮箱格式正确
                        if '@' in email:
                            accounts.append({
                                'username': username,
                                'email': email,
                                'password': password
                            })
                            logging.info(f"解析账号: 用户名={username}, 邮箱={email}, 密码={password}")
            
    except Exception as e:
        logging.error(f"读取文件失败: {str(e)}")
    
    return accounts


def main():
    try:
        # 使用指定的Miko邮箱凭据
        username = "kimi11"
        email_address = "kimi11@jbjj.site"
        password = "93921438"

        logging.info("=" * 60)
        logging.info("Miko邮箱 IMAP 邮件接收器")
        logging.info("=" * 60)
        logging.info(f"用户名: {username}")
        logging.info(f"邮箱: {email_address}")
        logging.info(f"IMAP服务器: localhost:143")
        logging.info("=" * 60)

        # 创建接收器
        receiver = EmailReceiver(username, email_address, password)

        # 运行轮询
        receiver.run_polling()

    except KeyboardInterrupt:
        logging.info("用户中断，停止轮询")
    except Exception as e:
        logging.error(f"程序运行出错: {str(e)}")


if __name__ == "__main__":
    main() 