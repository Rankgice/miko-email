@echo off
echo 启动NBEmail系统 - 多SMTP端口模式
echo.
echo 此模式将启动以下SMTP端口:
echo   - 2525 (对应标准25端口)
echo   - 2587 (对应标准587端口，支持STARTTLS)
echo   - 2465 (对应标准465端口，支持SSL/TLS)
echo.
echo 注意: 使用高端口号避免权限问题
echo 如需使用标准端口(25,587,465)，请以管理员身份运行
echo.

REM 设置多SMTP端口配置
set ENABLE_MULTI_SMTP=true
set SMTP_PORT=25
set SMTP_PORT_587=587
set SMTP_PORT_465=465

REM 设置外部SMTP配置（可选）
REM set OUTBOUND_SMTP_HOST=smtp.yourdomain.com
REM set OUTBOUND_SMTP_PORT=587
REM set OUTBOUND_SMTP_USER=your-email@yourdomain.com
REM set OUTBOUND_SMTP_PASSWORD=your-password
REM set OUTBOUND_SMTP_TLS=true

echo 启动NBEmail服务器...
nbemail.exe

pause
