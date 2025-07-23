@echo off
echo 启动NBEmail系统 - 配置外部SMTP发送功能
echo.

REM 设置外部SMTP配置（请根据您的实际情况修改）
echo 配置外部SMTP服务器...
set OUTBOUND_SMTP_HOST=mail.jbjj.site
set OUTBOUND_SMTP_PORT=587
set OUTBOUND_SMTP_USER=smtp@jbjj.site
set OUTBOUND_SMTP_PASSWORD=your-smtp-password
set OUTBOUND_SMTP_TLS=true

REM 设置域名（请修改为您的实际域名）
set DOMAIN=jbjj.site

echo.
echo 外部SMTP配置:
echo   服务器: %OUTBOUND_SMTP_HOST%
echo   端口: %OUTBOUND_SMTP_PORT%
echo   用户: %OUTBOUND_SMTP_USER%
echo   TLS: %OUTBOUND_SMTP_TLS%
echo   域名: %DOMAIN%
echo.

echo 启动NBEmail服务器...
nbemail.exe --port 8080 --smtp-port 25

pause
