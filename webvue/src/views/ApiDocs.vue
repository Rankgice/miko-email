<template>
  <div class="api-docs-page">
    <!-- 顶部导航 -->
    <nav class="top-nav">
      <div class="logo">
        <i class="fas fa-code"></i>
        <span>API文档</span>
      </div>
      <router-link to="/" class="back-home">
        <i class="fas fa-home"></i>
        返回首页
      </router-link>
    </nav>

    <!-- 侧边栏导航 -->
    <div class="docs-layout">
      <aside class="docs-sidebar">
        <div class="sidebar-content">
          <h3>API文档</h3>
          <nav class="docs-nav">
            <a href="#introduction" @click="scrollTo('introduction')" :class="{ active: activeSection === 'introduction' }">
              <i class="fas fa-info-circle"></i>
              介绍
            </a>
            <a href="#authentication" @click="scrollTo('authentication')" :class="{ active: activeSection === 'authentication' }">
              <i class="fas fa-key"></i>
              身份验证
            </a>
            <a href="#users" @click="scrollTo('users')" :class="{ active: activeSection === 'users' }">
              <i class="fas fa-users"></i>
              用户管理
            </a>
            <a href="#emails" @click="scrollTo('emails')" :class="{ active: activeSection === 'emails' }">
              <i class="fas fa-envelope"></i>
              邮件管理
            </a>
            <a href="#mailboxes" @click="scrollTo('mailboxes')" :class="{ active: activeSection === 'mailboxes' }">
              <i class="fas fa-inbox"></i>
              邮箱管理
            </a>
            <a href="#domains" @click="scrollTo('domains')" :class="{ active: activeSection === 'domains' }">
              <i class="fas fa-globe"></i>
              域名管理
            </a>
            <a href="#errors" @click="scrollTo('errors')" :class="{ active: activeSection === 'errors' }">
              <i class="fas fa-exclamation-triangle"></i>
              错误代码
            </a>
          </nav>
        </div>
      </aside>

      <!-- 主要内容 -->
      <main class="docs-content">
        <!-- 介绍 -->
        <section id="introduction" class="docs-section">
          <h1>Miko邮箱系统 API 文档</h1>
          <p class="intro-text">
            欢迎使用Miko邮箱系统API！本文档将帮助您快速集成和使用我们的邮件服务API。
          </p>

          <div class="info-card">
            <h3>基础信息</h3>
            <ul>
              <li><strong>API版本:</strong> v1.0</li>
              <li><strong>基础URL:</strong> <code>https://api.miko-email.com/api</code></li>
              <li><strong>协议:</strong> HTTPS</li>
              <li><strong>数据格式:</strong> JSON</li>
            </ul>
          </div>

          <div class="code-example">
            <h4>快速开始</h4>
            <pre><code>curl -X GET "https://api.miko-email.com/api/user/profile" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json"</code></pre>
          </div>
        </section>

        <!-- 身份验证 -->
        <section id="authentication" class="docs-section">
          <h2>身份验证</h2>
          <p>所有API请求都需要进行身份验证。我们支持以下认证方式：</p>

          <div class="auth-method">
            <h3>Bearer Token</h3>
            <p>在请求头中包含访问令牌：</p>
            <div class="code-example">
              <pre><code>Authorization: Bearer YOUR_ACCESS_TOKEN</code></pre>
            </div>
          </div>

          <div class="auth-method">
            <h3>获取访问令牌</h3>
            <div class="api-endpoint">
              <div class="method post">POST</div>
              <div class="url">/auth/login</div>
            </div>
            <div class="code-example">
              <h4>请求示例</h4>
              <pre><code>{
  "email": "user@example.com",
  "password": "your_password"
}</code></pre>
            </div>
            <div class="code-example">
              <h4>响应示例</h4>
              <pre><code>{
  "code": 0,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "email": "user@example.com",
      "name": "用户名"
    }
  }
}</code></pre>
            </div>
          </div>
        </section>

        <!-- 用户管理 -->
        <section id="users" class="docs-section">
          <h2>用户管理</h2>

          <div class="api-group">
            <h3>获取用户信息</h3>
            <div class="api-endpoint">
              <div class="method get">GET</div>
              <div class="url">/user/profile</div>
            </div>
            <p>获取当前登录用户的详细信息。</p>

            <div class="code-example">
              <h4>响应示例</h4>
              <pre><code>{
  "code": 0,
  "message": "获取成功",
  "data": {
    "id": 1,
    "email": "user@example.com",
    "name": "用户名",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}</code></pre>
            </div>
          </div>

          <div class="api-group">
            <h3>更新用户信息</h3>
            <div class="api-endpoint">
              <div class="method put">PUT</div>
              <div class="url">/user/profile</div>
            </div>
            <p>更新当前登录用户的信息。</p>

            <div class="code-example">
              <h4>请求示例</h4>
              <pre><code>{
  "name": "新用户名",
  "phone": "13800138000"
}</code></pre>
            </div>
          </div>
        </section>

        <!-- 邮件管理 -->
        <section id="emails" class="docs-section">
          <h2>邮件管理</h2>

          <div class="api-group">
            <h3>获取邮件列表</h3>
            <div class="api-endpoint">
              <div class="method get">GET</div>
              <div class="url">/emails</div>
            </div>
            <p>获取用户的邮件列表，支持分页和筛选。</p>

            <div class="params-table">
              <h4>查询参数</h4>
              <table>
                <thead>
                  <tr>
                    <th>参数</th>
                    <th>类型</th>
                    <th>必填</th>
                    <th>说明</th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td>page</td>
                    <td>integer</td>
                    <td>否</td>
                    <td>页码，默认为1</td>
                  </tr>
                  <tr>
                    <td>limit</td>
                    <td>integer</td>
                    <td>否</td>
                    <td>每页数量，默认为20</td>
                  </tr>
                  <tr>
                    <td>folder</td>
                    <td>string</td>
                    <td>否</td>
                    <td>文件夹类型：inbox, sent, draft, trash</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

          <div class="api-group">
            <h3>发送邮件</h3>
            <div class="api-endpoint">
              <div class="method post">POST</div>
              <div class="url">/emails/send</div>
            </div>
            <p>发送新邮件。</p>

            <div class="code-example">
              <h4>请求示例</h4>
              <pre><code>{
  "to": ["recipient@example.com"],
  "cc": ["cc@example.com"],
  "bcc": ["bcc@example.com"],
  "subject": "邮件主题",
  "content": "邮件内容",
  "attachments": []
}</code></pre>
            </div>
          </div>
        </section>

        <!-- 邮箱管理 -->
        <section id="mailboxes" class="docs-section">
          <h2>邮箱管理</h2>

          <div class="api-group">
            <h3>获取邮箱列表</h3>
            <div class="api-endpoint">
              <div class="method get">GET</div>
              <div class="url">/mailboxes</div>
            </div>
            <p>获取用户的邮箱列表。</p>
          </div>

          <div class="api-group">
            <h3>创建邮箱</h3>
            <div class="api-endpoint">
              <div class="method post">POST</div>
              <div class="url">/mailboxes</div>
            </div>
            <p>创建新的邮箱。</p>

            <div class="code-example">
              <h4>请求示例</h4>
              <pre><code>{
  "email": "newbox@example.com",
  "password": "secure_password",
  "domain_id": 1
}</code></pre>
            </div>
          </div>
        </section>

        <!-- 域名管理 -->
        <section id="domains" class="docs-section">
          <h2>域名管理</h2>

          <div class="api-group">
            <h3>获取域名列表</h3>
            <div class="api-endpoint">
              <div class="method get">GET</div>
              <div class="url">/domains</div>
            </div>
            <p>获取用户的域名列表。</p>
          </div>

          <div class="api-group">
            <h3>添加域名</h3>
            <div class="api-endpoint">
              <div class="method post">POST</div>
              <div class="url">/domains</div>
            </div>
            <p>添加新的域名。</p>

            <div class="code-example">
              <h4>请求示例</h4>
              <pre><code>{
  "domain": "example.com",
  "description": "公司域名"
}</code></pre>
            </div>
          </div>
        </section>

        <!-- 错误代码 -->
        <section id="errors" class="docs-section">
          <h2>错误代码</h2>
          <p>API使用标准的HTTP状态码和自定义错误代码。</p>

          <div class="error-table">
            <table>
              <thead>
                <tr>
                  <th>HTTP状态码</th>
                  <th>错误代码</th>
                  <th>说明</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td>200</td>
                  <td>0</td>
                  <td>请求成功</td>
                </tr>
                <tr>
                  <td>400</td>
                  <td>1001</td>
                  <td>参数错误</td>
                </tr>
                <tr>
                  <td>401</td>
                  <td>1002</td>
                  <td>未授权访问</td>
                </tr>
                <tr>
                  <td>403</td>
                  <td>1003</td>
                  <td>权限不足</td>
                </tr>
                <tr>
                  <td>404</td>
                  <td>1004</td>
                  <td>资源不存在</td>
                </tr>
                <tr>
                  <td>500</td>
                  <td>1005</td>
                  <td>服务器内部错误</td>
                </tr>
              </tbody>
            </table>
          </div>
        </section>
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const activeSection = ref('introduction')

const scrollTo = (sectionId) => {
  const element = document.getElementById(sectionId)
  if (element) {
    element.scrollIntoView({ behavior: 'smooth' })
  }
}

const handleScroll = () => {
  const sections = ['introduction', 'authentication', 'users', 'emails', 'mailboxes', 'domains', 'errors']

  for (const section of sections) {
    const element = document.getElementById(section)
    if (element) {
      const rect = element.getBoundingClientRect()
      if (rect.top <= 100 && rect.bottom >= 100) {
        activeSection.value = section
        break
      }
    }
  }
}

onMounted(() => {
  window.addEventListener('scroll', handleScroll)
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>

<style scoped>
.api-docs-page {
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
  color: var(--text-primary);
  min-height: 100vh;
}

/* 顶部导航 */
.top-nav {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 40px;
  background: rgba(30, 41, 59, 0.9);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid var(--border);
  z-index: 100;
}

.logo {
  font-size: 24px;
  font-weight: 700;
  background: linear-gradient(to right, var(--primary), var(--primary-light));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  display: flex;
  align-items: center;
  gap: 12px;
}

.back-home {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  background: rgba(0, 180, 216, 0.1);
  border: 1px solid var(--primary);
  border-radius: 8px;
  color: var(--primary);
  text-decoration: none;
  transition: all 0.3s ease;
}

.back-home:hover {
  background: var(--primary);
  color: white;
}

/* 文档布局 */
.docs-layout {
  display: flex;
  padding-top: 80px;
  min-height: calc(100vh - 80px);
}

/* 侧边栏 */
.docs-sidebar {
  width: 280px;
  background: rgba(15, 23, 42, 0.9);
  border-right: 1px solid var(--border);
  position: fixed;
  height: calc(100vh - 80px);
  overflow-y: auto;
}

.sidebar-content {
  padding: 30px 20px;
}

.sidebar-content h3 {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 20px;
  color: var(--text-primary);
}

.docs-nav {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.docs-nav a {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  color: var(--text-secondary);
  text-decoration: none;
  border-radius: 8px;
  transition: all 0.3s ease;
  font-weight: 500;
}

.docs-nav a:hover {
  background: rgba(0, 180, 216, 0.1);
  color: var(--primary);
}

.docs-nav a.active {
  background: rgba(0, 180, 216, 0.15);
  color: var(--primary);
  border-left: 3px solid var(--primary);
}

/* 主要内容 */
.docs-content {
  flex: 1;
  margin-left: 280px;
  padding: 40px;
  max-width: calc(100% - 280px);
}

.docs-section {
  margin-bottom: 60px;
}

.docs-section h1 {
  font-size: 36px;
  font-weight: 700;
  margin-bottom: 20px;
  background: linear-gradient(to right, var(--primary), var(--primary-light));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.docs-section h2 {
  font-size: 28px;
  font-weight: 600;
  margin-bottom: 20px;
  color: var(--text-primary);
  border-bottom: 2px solid var(--border);
  padding-bottom: 10px;
}

.docs-section h3 {
  font-size: 22px;
  font-weight: 600;
  margin-bottom: 15px;
  color: var(--text-primary);
}

.docs-section h4 {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 10px;
  color: var(--text-primary);
}

.intro-text {
  font-size: 18px;
  color: var(--text-secondary);
  line-height: 1.6;
  margin-bottom: 30px;
}

/* 信息卡片 */
.info-card {
  background: linear-gradient(135deg, rgba(30, 41, 59, 0.9), rgba(15, 23, 42, 0.95));
  border: 1px solid var(--border);
  border-radius: 12px;
  padding: 25px;
  margin-bottom: 30px;
}

.info-card h3 {
  margin-bottom: 15px;
  color: var(--primary);
}

.info-card ul {
  list-style: none;
  padding: 0;
}

.info-card li {
  margin-bottom: 8px;
  color: var(--text-secondary);
}

.info-card code {
  background: rgba(0, 180, 216, 0.1);
  color: var(--primary);
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
}

/* 代码示例 */
.code-example {
  margin: 20px 0;
}

.code-example h4 {
  margin-bottom: 10px;
  color: var(--text-primary);
}

.code-example pre {
  background: rgba(15, 23, 42, 0.9);
  border: 1px solid var(--border);
  border-radius: 8px;
  padding: 20px;
  overflow-x: auto;
  margin: 0;
}

.code-example code {
  color: #e2e8f0;
  font-family: 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.5;
}

/* API端点 */
.api-endpoint {
  display: flex;
  align-items: center;
  gap: 15px;
  margin: 15px 0;
  padding: 15px;
  background: rgba(15, 23, 42, 0.5);
  border-radius: 8px;
  border-left: 4px solid var(--primary);
}

.method {
  padding: 6px 12px;
  border-radius: 6px;
  font-weight: 600;
  font-size: 12px;
  text-transform: uppercase;
  min-width: 60px;
  text-align: center;
}

.method.get {
  background: #10b981;
  color: white;
}

.method.post {
  background: #3b82f6;
  color: white;
}

.method.put {
  background: #f59e0b;
  color: white;
}

.method.delete {
  background: #ef4444;
  color: white;
}

.url {
  font-family: 'Courier New', monospace;
  font-size: 16px;
  color: var(--text-primary);
  font-weight: 500;
}

/* API组 */
.api-group {
  margin-bottom: 40px;
  padding: 25px;
  background: linear-gradient(135deg, rgba(30, 41, 59, 0.5), rgba(15, 23, 42, 0.7));
  border: 1px solid var(--border);
  border-radius: 12px;
}

/* 参数表格 */
.params-table, .error-table {
  margin: 20px 0;
}

.params-table table, .error-table table {
  width: 100%;
  border-collapse: collapse;
  background: rgba(15, 23, 42, 0.5);
  border-radius: 8px;
  overflow: hidden;
}

.params-table th, .error-table th,
.params-table td, .error-table td {
  padding: 12px 15px;
  text-align: left;
  border-bottom: 1px solid var(--border);
}

.params-table th, .error-table th {
  background: rgba(0, 180, 216, 0.1);
  color: var(--primary);
  font-weight: 600;
}

.params-table td, .error-table td {
  color: var(--text-secondary);
}

/* 认证方法 */
.auth-method {
  margin-bottom: 30px;
  padding: 25px;
  background: linear-gradient(135deg, rgba(30, 41, 59, 0.5), rgba(15, 23, 42, 0.7));
  border: 1px solid var(--border);
  border-radius: 12px;
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .docs-sidebar {
    transform: translateX(-100%);
    transition: transform 0.3s ease;
  }

  .docs-content {
    margin-left: 0;
    max-width: 100%;
    padding: 20px;
  }
}

@media (max-width: 768px) {
  .top-nav {
    padding: 15px 20px;
  }

  .docs-content {
    padding: 15px;
  }

  .api-endpoint {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
}
</style>
