# Miko邮箱系统 - Vue3前端

这是Miko邮箱系统的Vue3前端项目，提供现代化的邮件管理界面。

## 🚀 项目特性

- **Vue 3** - 使用最新的Vue 3 Composition API
- **Vue Router 4** - 现代化的路由管理
- **Pinia** - 轻量级状态管理
- **Vite** - 快速的构建工具
- **Sass** - CSS预处理器
- **响应式设计** - 支持多种设备尺寸
- **现代化UI** - 美观的用户界面设计

## 📁 项目结构

```
webvue/
├── src/
│   ├── components/          # 公共组件
│   ├── layouts/            # 布局组件
│   │   ├── UserLayout.vue  # 用户端布局
│   │   └── AdminLayout.vue # 管理员布局
│   ├── views/              # 页面组件
│   │   ├── Home.vue        # 首页
│   │   ├── user/           # 用户端页面
│   │   │   ├── Login.vue   # 用户登录
│   │   │   ├── Register.vue # 用户注册
│   │   │   ├── Dashboard.vue # 用户仪表盘
│   │   │   ├── Inbox.vue   # 收件箱
│   │   │   ├── Outbox.vue  # 发件箱
│   │   │   ├── Settings.vue # 设置
│   │   │   ├── Mailboxes.vue # 邮箱管理
│   │   │   └── ForwardRules.vue # 转发规则
│   │   └── admin/          # 管理员页面
│   │       ├── Login.vue   # 管理员登录
│   │       ├── Dashboard.vue # 管理员仪表盘
│   │       ├── Users.vue   # 用户管理
│   │       ├── Mailboxes.vue # 邮箱管理
│   │       ├── Domains.vue # 域名管理
│   │       └── Captcha.vue # 验证码管理
│   ├── router/             # 路由配置
│   ├── stores/             # 状态管理
│   ├── utils/              # 工具函数
│   ├── styles/             # 全局样式
│   ├── App.vue             # 根组件
│   └── main.js             # 入口文件
├── public/                 # 静态资源
├── index.html              # HTML模板
├── package.json            # 项目配置
├── vite.config.js          # Vite配置
└── README.md               # 项目说明
```

## 🎨 页面功能

### 用户端功能
- **首页** - 项目介绍和快速入口
- **用户登录/注册** - 用户身份验证
- **仪表盘** - 邮件统计和快速操作
- **收件箱** - 查看收到的邮件
- **发件箱** - 查看已发送的邮件
- **邮箱管理** - 管理多个邮箱账户
- **转发规则** - 设置邮件转发规则
- **账户设置** - 个人信息和偏好设置

### 管理员功能
- **管理员登录** - 管理员身份验证
- **管理仪表盘** - 系统统计和状态监控
- **用户管理** - 管理系统用户
- **邮箱管理** - 管理所有邮箱账户
- **域名管理** - 管理邮件域名
- **验证码管理** - 管理验证码配置

## 🛠️ 开发环境

### 环境要求
- Node.js 16+
- npm 或 yarn

### 安装依赖
```bash
npm install
```

### 启动开发服务器
```bash
npm run dev
```

项目将在 http://localhost:3001 启动

### 构建生产版本
```bash
npm run build
```

### 预览生产版本
```bash
npm run preview
```

## 🎯 API集成

项目使用Axios进行HTTP请求，API基础配置在 `src/utils/api.js` 中。

### API配置
- 基础URL: `/api`
- 超时时间: 10秒
- 自动携带Cookies用于session认证
- 自动处理401/403等错误状态

### 状态管理
使用Pinia进行状态管理，主要store包括：
- `auth.js` - 用户认证状态管理

## 🎨 样式系统

### CSS变量
项目使用CSS变量定义主题色彩：

```css
:root {
  /* 用户端主题 */
  --primary: #00B4D8;
  --accent: #FF6B6B;
  --success: #10b981;
  --warning: #f59e0b;
  
  /* 管理员主题 */
  --admin-primary: #2563eb;
  --admin-secondary: #8b5cf6;
  --admin-success: #10b981;
  --admin-danger: #ef4444;
}
```

### 响应式设计
- 移动端优先设计
- 支持平板和桌面端
- 使用CSS Grid和Flexbox布局

## 🔧 配置说明

### Vite配置 (vite.config.js)
- Vue插件配置
- 路径别名 (@指向src目录)
- 开发服务器配置 (端口3001)
- 代理配置 (API请求代理到后端)

### 路由配置
- 用户端路由 (/user/*)
- 管理员路由 (/admin/*)
- 路由守卫 (认证检查)

## 📱 浏览器支持

- Chrome 90+
- Firefox 88+
- Safari 14+
- Edge 90+

## 🤝 开发指南

### 代码规范
- 使用Vue 3 Composition API
- 组件使用PascalCase命名
- 文件使用kebab-case命名
- 使用ESLint进行代码检查

### 组件开发
- 优先使用Composition API
- 合理拆分组件
- 注意组件的可复用性
- 添加适当的注释

### 样式开发
- 使用scoped样式
- 遵循BEM命名规范
- 使用CSS变量
- 注意响应式设计

## 📄 许可证

MIT License

## 👥 贡献

欢迎提交Issue和Pull Request来改进项目。

## 📞 联系方式

如有问题，请联系开发团队。
