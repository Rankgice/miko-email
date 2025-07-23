@echo off
chcp 65001 >nul
echo SMTP连接测试工具
echo ==================
echo.

if "%1"=="" (
    echo 用法: test-smtp.bat ^<host:port^> [user] [password] [tls]
    echo.
    echo 示例:
    echo   test-smtp.bat smtp.gmail.com:587 user@gmail.com password true
    echo   test-smtp.bat smtp.gmail.com:465 user@gmail.com password true
    echo   test-smtp.bat mail.example.com:25
    echo.
    echo 常见SMTP服务器配置:
    echo   Gmail:     smtp.gmail.com:587 (STARTTLS) 或 smtp.gmail.com:465 (SSL)
    echo   Outlook:   smtp-mail.outlook.com:587 (STARTTLS)
    echo   QQ邮箱:    smtp.qq.com:587 (STARTTLS) 或 smtp.qq.com:465 (SSL)
    echo   163邮箱:   smtp.163.com:587 (STARTTLS) 或 smtp.163.com:465 (SSL)
    echo.
    pause
    exit /b 1
)

echo 正在测试SMTP连接...
echo.

go run tools/smtp-test.go %*

echo.
echo 测试完成！
echo.
echo 如果连接失败，请检查:
echo 1. 网络连接是否正常
echo 2. SMTP服务器地址和端口是否正确
echo 3. 用户名和密码是否正确
echo 4. 是否需要启用"允许不够安全的应用"或使用应用专用密码
echo 5. 防火墙是否阻止了连接
echo.
pause
