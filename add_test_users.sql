-- 添加测试用户数据
INSERT INTO users (id, username, password, email, is_active, contribution, invite_code, created_at, updated_at) 
VALUES 
(1, 'user1', '$2a$10$ehfjJCykP5Sb/DQryLi8FuEIqy83qxLgQyMImKpiwE7UF5zoP4fsi', 'user1@example.com', 1, 0, 'INVITE_USER1', datetime('now'), datetime('now')),
(2, 'user2', '$2a$10$ehfjJCykP5Sb/DQryLi8FuEIqy83qxLgQyMImKpiwE7UF5zoP4fsi', 'user2@example.com', 1, 0, 'INVITE_USER2', datetime('now'), datetime('now')),
(3, 'user3', '$2a$10$ehfjJCykP5Sb/DQryLi8FuEIqy83qxLgQyMImKpiwE7UF5zoP4fsi', 'user3@example.com', 1, 0, 'INVITE_USER3', datetime('now'), datetime('now')),
(4, 'user4', '$2a$10$ehfjJCykP5Sb/DQryLi8FuEIqy83qxLgQyMImKpiwE7UF5zoP4fsi', 'user4@example.com', 1, 0, 'INVITE_USER4', datetime('now'), datetime('now')),
(5, 'kimi', '$2a$10$ehfjJCykP5Sb/DQryLi8FuEIqy83qxLgQyMImKpiwE7UF5zoP4fsi', 'kimi@example.com', 1, 0, 'INVITE_KIMI', datetime('now'), datetime('now'));

-- 查看结果
SELECT COUNT(*) as user_count FROM users;
SELECT id, username, email FROM users;
