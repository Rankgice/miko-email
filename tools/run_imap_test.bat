@echo off
chcp 65001 >nul
cd /d "%~dp0"

echo ====================================
echo 🧪 IMAP测试工具 - Go版本
echo ====================================
echo.

echo 📝 使用默认测试账号:
echo    用户名: kimi11
echo    邮箱: kimi11@jbjj.site  
echo    密码: 93921438
echo    服务器: localhost:143
echo.

echo 🚀 编译并运行测试...
go run imap_tester.go

echo.
echo 按任意键退出...
pause >nul
