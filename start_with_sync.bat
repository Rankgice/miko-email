@echo off
echo ========================================
echo    Miko邮箱系统启动脚本
echo ========================================
echo.

echo 🔄 正在同步管理员信息...
go run tools/sync_admin.go sync
if %errorlevel% neq 0 (
    echo ❌ 管理员信息同步失败
    pause
    exit /b 1
)

echo.
echo 📋 显示当前配置信息...
go run tools/config_manager.go show

echo.
echo 🚀 启动Miko邮箱服务器...
echo.
echo 💡 提示:
echo    - Web管理界面: http://localhost:8080
echo    - 管理员登录: http://localhost:8080/admin/login
echo    - 按 Ctrl+C 停止服务器
echo.

go run main.go

pause
