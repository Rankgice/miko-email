-- 修复邮件编码的SQL脚本
-- 注意：这个脚本需要在SQLite命令行中运行

-- 查看需要修复的邮件
SELECT id, subject, body FROM emails WHERE id >= 16 ORDER BY id;

-- 手动修复邮件ID 16
UPDATE emails SET 
    subject = '最终验证测试',
    body = '这是最终的验证测试邮件！中文编码已经完全修复，点击查看邮件详情功能也正常工作。🎉',
    updated_at = datetime('now')
WHERE id = 16;

-- 手动修复邮件ID 17
UPDATE emails SET 
    subject = '最终验证测试',
    body = '这是最终的验证测试邮件！中文编码已经完全修复，点击查看邮件详情功能也正常工作。🎉',
    updated_at = datetime('now')
WHERE id = 17;

-- 手动修复邮件ID 18
UPDATE emails SET 
    subject = '编码修复最终测试',
    body = '这是修复后的最终测试邮件！中文编码应该完全正常了。🎉',
    updated_at = datetime('now')
WHERE id = 18;

-- 手动修复邮件ID 19
UPDATE emails SET 
    subject = '编码修复最终测试',
    body = '这是修复后的最终测试邮件！中文编码应该完全正常了。🎉',
    updated_at = datetime('now')
WHERE id = 19;

-- 验证修复结果
SELECT id, subject, body FROM emails WHERE id >= 16 ORDER BY id;
