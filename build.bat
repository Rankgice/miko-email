@echo off
echo 构建Miko邮箱系统...

REM 检查Node.js是否安装
where node >nul 2>nul
if %errorlevel% neq 0 (
    echo 错误: Node.js未安装，请先安装Node.js
    pause
    exit /b 1
)

REM 检查npm是否安装
where npm >nul 2>nul
if %errorlevel% neq 0 (
    echo 错误: npm未安装，请先安装npm
    pause
    exit /b 1
)

REM 构建前端Vue应用
echo 构建前端Vue应用...
cd frontend

REM 安装依赖
echo 安装前端依赖...
npm install

REM 构建生产版本
echo 构建生产版本...
npm run build

REM 返回根目录
cd ..

REM 构建Go后端
echo 构建Go后端...
go mod tidy
go build -o miko-email.exe main.go

echo 构建完成！
echo.
echo 启动说明：
echo 1. 确保配置文件 config.yaml 存在
echo 2. 运行: miko-email.exe
echo 3. 访问: http://localhost:8080
echo.
echo 注意: 系统已强制使用SQLite数据库，禁止使用MySQL
pause
