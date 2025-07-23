@echo off
chcp 65001 >nul
echo === jbjj.site SMTP配置添加脚本 ===
echo.

REM 检查服务是否运行
curl -s http://127.0.0.1:8080 >nul 2>&1
if errorlevel 1 (
    echo ❌ 邮件系统未运行，请先启动服务
    echo 运行: nbemail.exe
    pause
    exit /b 1
)

echo ✅ 邮件系统正在运行
echo.

echo 请确保您已经以管理员身份登录系统
echo 如果未登录，请访问: http://127.0.0.1:8080/login
echo.

echo 请输入 jbjj.site 域名的SMTP配置信息：
echo.

set /p smtp_host="SMTP服务器地址 (默认: mail.jbjj.site): "
if "%smtp_host%"=="" set smtp_host=mail.jbjj.site

set /p smtp_port="SMTP端口 (默认: 587): "
if "%smtp_port%"=="" set smtp_port=587

set /p smtp_user="用户名 (默认: 1keqb385916@jbjj.site): "
if "%smtp_user%"=="" set smtp_user=1keqb385916@jbjj.site

set /p smtp_password="密码: "

set /p enable_tls="启用TLS? (y/n, 默认: y): "
if "%enable_tls%"=="" set enable_tls=y

if /i "%enable_tls%"=="y" (
    set tls_enabled=true
) else (
    set tls_enabled=false
)

echo.
echo 配置信息：
echo 域名: jbjj.site
echo SMTP服务器: %smtp_host%
echo 端口: %smtp_port%
echo 用户名: %smtp_user%
echo 启用TLS: %tls_enabled%
echo.

set /p confirm="确认添加配置? (y/n): "
if /i not "%confirm%"=="y" (
    echo 已取消
    pause
    exit /b 0
)

echo.
echo 正在添加SMTP配置...

REM 创建临时JSON文件
echo { > temp_config.json
echo     "domain": "jbjj.site", >> temp_config.json
echo     "host": "%smtp_host%", >> temp_config.json
echo     "port": %smtp_port%, >> temp_config.json
echo     "user": "%smtp_user%", >> temp_config.json
echo     "password": "%smtp_password%", >> temp_config.json
echo     "tls": %tls_enabled% >> temp_config.json
echo } >> temp_config.json

REM 发送请求
curl -s -X POST -H "Content-Type: application/json" -d @temp_config.json -b "cookies.txt" http://127.0.0.1:8080/api/smtp-configs > response.txt

REM 检查响应
findstr "success.*true" response.txt >nul
if not errorlevel 1 (
    echo ✅ SMTP配置添加成功！
    echo.
    echo 现在您可以使用 @jbjj.site 域名的邮箱发送邮件了
    echo 系统会自动使用刚才配置的SMTP服务器
) else (
    echo ❌ 添加失败
    echo 响应:
    type response.txt
    echo.
    echo 可能的原因：
    echo 1. 未以管理员身份登录
    echo 2. SMTP服务器信息不正确
    echo 3. 网络连接问题
)

REM 清理临时文件
del temp_config.json >nul 2>&1
del response.txt >nul 2>&1

echo.
echo 您也可以通过Web界面管理SMTP配置：
echo 访问: http://127.0.0.1:8080/smtp-configs

pause
