-- 初始化默认域名
INSERT OR IGNORE INTO domains (name, is_verified, is_active, mx_record, a_record, txt_record, created_at, updated_at) 
VALUES 
('localhost', 1, 1, 'localhost', '127.0.0.1', 'v=spf1 ip4:127.0.0.1 ~all', datetime('now'), datetime('now')),
('example.com', 0, 1, 'mail.example.com', '192.168.1.100', 'v=spf1 ip4:192.168.1.100 ~all', datetime('now'), datetime('now')),
('test.local', 1, 1, 'mail.test.local', '127.0.0.1', 'v=spf1 ip4:127.0.0.1 ~all', datetime('now'), datetime('now'));
