<template>
  <div class="help-center-page">
    <!-- 顶部导航 -->
    <nav class="top-nav">
      <div class="logo">
        <i class="fas fa-question-circle"></i>
        <span>帮助中心</span>
      </div>
      <router-link to="/" class="back-home">
        <i class="fas fa-home"></i>
        返回首页
      </router-link>
    </nav>

    <!-- 主要内容 -->
    <main class="main-content">
      <!-- 搜索区域 -->
      <section class="search-section">
        <div class="container">
          <h1 class="search-title">我们能为您提供什么帮助？</h1>
          <div class="search-box">
            <i class="fas fa-search"></i>
            <input
              type="text"
              v-model="searchQuery"
              placeholder="搜索帮助文档..."
              @keyup.enter="searchHelp"
            >
            <button @click="searchHelp" class="search-btn">搜索</button>
          </div>
        </div>
      </section>

      <!-- 快速导航 -->
      <section class="quick-nav-section">
        <div class="container">
          <div class="quick-nav-grid">
            <div class="nav-card" @click="scrollToSection('getting-started')">
              <div class="nav-icon">
                <i class="fas fa-rocket"></i>
              </div>
              <h3>快速开始</h3>
              <p>了解如何快速设置和使用邮箱系统</p>
            </div>

            <div class="nav-card" @click="scrollToSection('account-management')">
              <div class="nav-icon">
                <i class="fas fa-user-cog"></i>
              </div>
              <h3>账户管理</h3>
              <p>管理您的账户设置和个人信息</p>
            </div>

            <div class="nav-card" @click="scrollToSection('email-features')">
              <div class="nav-icon">
                <i class="fas fa-envelope-open"></i>
              </div>
              <h3>邮件功能</h3>
              <p>了解邮件发送、接收和管理功能</p>
            </div>

            <div class="nav-card" @click="scrollToSection('troubleshooting')">
              <div class="nav-icon">
                <i class="fas fa-tools"></i>
              </div>
              <h3>故障排除</h3>
              <p>解决常见问题和技术故障</p>
            </div>
          </div>
        </div>
      </section>

      <!-- 帮助内容 -->
      <section class="help-content-section">
        <div class="container">
          <!-- 快速开始 -->
          <div id="getting-started" class="help-category">
            <h2>快速开始</h2>
            <div class="help-items">
              <div class="help-item">
                <h3>
                  <i class="fas fa-user-plus"></i>
                  如何注册账户？
                </h3>
                <div class="help-content">
                  <p>注册Miko邮箱账户非常简单：</p>
                  <ol>
                    <li>访问首页，点击"立即注册"按钮</li>
                    <li>填写您的邮箱地址和密码</li>
                    <li>验证邮箱地址</li>
                    <li>完成注册，开始使用</li>
                  </ol>
                </div>
              </div>

              <div class="help-item">
                <h3>
                  <i class="fas fa-inbox"></i>
                  如何创建第一个邮箱？
                </h3>
                <div class="help-content">
                  <p>登录后创建邮箱的步骤：</p>
                  <ol>
                    <li>进入"邮箱管理"页面</li>
                    <li>点击"创建邮箱"按钮</li>
                    <li>选择域名并设置邮箱名称</li>
                    <li>设置邮箱密码</li>
                    <li>确认创建</li>
                  </ol>
                </div>
              </div>

              <div class="help-item">
                <h3>
                  <i class="fas fa-paper-plane"></i>
                  如何发送第一封邮件？
                </h3>
                <div class="help-content">
                  <p>发送邮件的简单步骤：</p>
                  <ol>
                    <li>进入"收件箱"页面</li>
                    <li>点击"写邮件"按钮</li>
                    <li>填写收件人、主题和内容</li>
                    <li>可选择添加附件</li>
                    <li>点击"发送"</li>
                  </ol>
                </div>
              </div>
            </div>
          </div>

          <!-- 账户管理 -->
          <div id="account-management" class="help-category">
            <h2>账户管理</h2>
            <div class="help-items">
              <div class="help-item">
                <h3>
                  <i class="fas fa-key"></i>
                  如何修改密码？
                </h3>
                <div class="help-content">
                  <p>修改账户密码的步骤：</p>
                  <ol>
                    <li>进入"设置"页面</li>
                    <li>选择"安全设置"</li>
                    <li>点击"修改密码"</li>
                    <li>输入当前密码和新密码</li>
                    <li>确认修改</li>
                  </ol>
                </div>
              </div>

              <div class="help-item">
                <h3>
                  <i class="fas fa-user-edit"></i>
                  如何更新个人信息？
                </h3>
                <div class="help-content">
                  <p>更新个人资料的方法：</p>
                  <ol>
                    <li>点击右上角头像</li>
                    <li>选择"个人资料"</li>
                    <li>编辑需要修改的信息</li>
                    <li>保存更改</li>
                  </ol>
                </div>
              </div>

              <div class="help-item">
                <h3>
                  <i class="fas fa-bell"></i>
                  如何设置通知偏好？
                </h3>
                <div class="help-content">
                  <p>自定义通知设置：</p>
                  <ol>
                    <li>进入"设置"页面</li>
                    <li>选择"通知设置"</li>
                    <li>选择您希望接收的通知类型</li>
                    <li>设置通知方式（邮件、短信等）</li>
                    <li>保存设置</li>
                  </ol>
                </div>
              </div>
            </div>
          </div>

          <!-- 邮件功能 -->
          <div id="email-features" class="help-category">
            <h2>邮件功能</h2>
            <div class="help-items">
              <div class="help-item">
                <h3>
                  <i class="fas fa-share"></i>
                  如何设置邮件转发？
                </h3>
                <div class="help-content">
                  <p>设置邮件自动转发：</p>
                  <ol>
                    <li>进入"转发规则"页面</li>
                    <li>点击"创建规则"</li>
                    <li>设置转发条件（发件人、主题等）</li>
                    <li>设置转发目标邮箱</li>
                    <li>启用规则</li>
                  </ol>
                </div>
              </div>

              <div class="help-item">
                <h3>
                  <i class="fas fa-filter"></i>
                  如何使用邮件过滤器？
                </h3>
                <div class="help-content">
                  <p>创建邮件过滤规则：</p>
                  <ol>
                    <li>进入"邮件设置"</li>
                    <li>选择"过滤器"</li>
                    <li>创建新的过滤规则</li>
                    <li>设置过滤条件和动作</li>
                    <li>测试并启用规则</li>
                  </ol>
                </div>
              </div>

              <div class="help-item">
                <h3>
                  <i class="fas fa-paperclip"></i>
                  如何发送大附件？
                </h3>
                <div class="help-content">
                  <p>发送大文件的方法：</p>
                  <ol>
                    <li>写邮件时点击"添加附件"</li>
                    <li>选择要发送的文件（最大100MB）</li>
                    <li>等待文件上传完成</li>
                    <li>发送邮件</li>
                  </ol>
                  <p><strong>注意：</strong>单个附件最大支持100MB，建议压缩大文件后发送。</p>
                </div>
              </div>
            </div>
          </div>

          <!-- 故障排除 -->
          <div id="troubleshooting" class="help-category">
            <h2>故障排除</h2>
            <div class="help-items">
              <div class="help-item">
                <h3>
                  <i class="fas fa-exclamation-triangle"></i>
                  无法登录怎么办？
                </h3>
                <div class="help-content">
                  <p>登录问题的解决方法：</p>
                  <ol>
                    <li>检查用户名和密码是否正确</li>
                    <li>确认账户是否已激活</li>
                    <li>尝试重置密码</li>
                    <li>清除浏览器缓存和Cookie</li>
                    <li>如仍有问题，请联系技术支持</li>
                  </ol>
                </div>
              </div>

              <div class="help-item">
                <h3>
                  <i class="fas fa-envelope-open"></i>
                  邮件发送失败怎么办？
                </h3>
                <div class="help-content">
                  <p>邮件发送问题的排查：</p>
                  <ol>
                    <li>检查收件人邮箱地址是否正确</li>
                    <li>确认附件大小是否超限</li>
                    <li>检查网络连接是否正常</li>
                    <li>查看发件箱是否有错误提示</li>
                    <li>尝试重新发送</li>
                  </ol>
                </div>
              </div>

              <div class="help-item">
                <h3>
                  <i class="fas fa-download"></i>
                  无法接收邮件怎么办？
                </h3>
                <div class="help-content">
                  <p>邮件接收问题的解决：</p>
                  <ol>
                    <li>检查垃圾邮件文件夹</li>
                    <li>确认邮箱存储空间是否充足</li>
                    <li>检查邮件过滤规则设置</li>
                    <li>刷新收件箱</li>
                    <li>联系发件人确认是否已发送</li>
                  </ol>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- 联系支持 -->
      <section class="support-section">
        <div class="container">
          <div class="support-content">
            <h2>还需要帮助？</h2>
            <p>如果您在帮助文档中没有找到答案，我们的支持团队随时为您提供帮助。</p>
            <div class="support-actions">
              <router-link to="/contact" class="btn btn-primary">
                <i class="fas fa-envelope"></i>
                联系我们
              </router-link>
              <router-link to="/tech-support" class="btn btn-secondary">
                <i class="fas fa-headset"></i>
                技术支持
              </router-link>
            </div>
          </div>
        </div>
      </section>
    </main>

    <!-- 底部 -->
    <footer class="footer">
      <div class="container">
        <p>&copy; 2024 Miko邮箱系统. 保留所有权利.</p>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const searchQuery = ref('')

const searchHelp = () => {
  if (searchQuery.value.trim()) {
    // 这里可以实现搜索功能
    console.log('搜索:', searchQuery.value)
    // 可以跳转到搜索结果页面或过滤当前内容
  }
}

const scrollToSection = (sectionId) => {
  const element = document.getElementById(sectionId)
  if (element) {
    element.scrollIntoView({ behavior: 'smooth' })
  }
}
</script>

<style scoped>
.help-center-page {
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

/* 主要内容 */
.main-content {
  padding-top: 80px;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 40px;
}

/* 搜索区域 */
.search-section {
  padding: 80px 0;
  text-align: center;
}

.search-title {
  font-size: 42px;
  font-weight: 700;
  margin-bottom: 40px;
  background: linear-gradient(to right, var(--primary), var(--primary-light));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.search-box {
  display: flex;
  align-items: center;
  max-width: 600px;
  margin: 0 auto;
  background: rgba(30, 41, 59, 0.9);
  border: 1px solid var(--border);
  border-radius: 12px;
  padding: 8px;
  gap: 15px;
}

.search-box i {
  color: var(--text-secondary);
  font-size: 18px;
  margin-left: 15px;
}

.search-box input {
  flex: 1;
  background: transparent;
  border: none;
  outline: none;
  color: var(--text-primary);
  font-size: 16px;
  padding: 12px 0;
}

.search-box input::placeholder {
  color: var(--text-secondary);
}

.search-btn {
  background: var(--primary);
  color: white;
  border: none;
  border-radius: 8px;
  padding: 12px 24px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.search-btn:hover {
  background: var(--primary-dark);
}

/* 快速导航 */
.quick-nav-section {
  padding: 60px 0;
  background: rgba(15, 23, 42, 0.5);
}

.quick-nav-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 30px;
}

.nav-card {
  background: linear-gradient(135deg, rgba(30, 41, 59, 0.9), rgba(15, 23, 42, 0.95));
  border: 1px solid var(--border);
  border-radius: 16px;
  padding: 30px 25px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
}

.nav-card:hover {
  transform: translateY(-5px);
  border-color: var(--primary);
  box-shadow: 0 20px 40px rgba(0, 180, 216, 0.1);
}

.nav-icon {
  width: 70px;
  height: 70px;
  background: linear-gradient(135deg, var(--primary), var(--primary-light));
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 20px;
  font-size: 28px;
  color: white;
}

.nav-card h3 {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 10px;
  color: var(--text-primary);
}

.nav-card p {
  color: var(--text-secondary);
  line-height: 1.6;
}

/* 帮助内容 */
.help-content-section {
  padding: 80px 0;
}

.help-category {
  margin-bottom: 60px;
}

.help-category h2 {
  font-size: 32px;
  font-weight: 700;
  margin-bottom: 40px;
  color: var(--text-primary);
  border-bottom: 2px solid var(--border);
  padding-bottom: 15px;
}

.help-items {
  display: grid;
  gap: 25px;
}

.help-item {
  background: linear-gradient(135deg, rgba(30, 41, 59, 0.9), rgba(15, 23, 42, 0.95));
  border: 1px solid var(--border);
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.3s ease;
}

.help-item:hover {
  border-color: var(--primary);
}

.help-item h3 {
  background: rgba(0, 180, 216, 0.1);
  padding: 20px 25px;
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.help-item h3:hover {
  background: rgba(0, 180, 216, 0.15);
}

.help-item h3 i {
  color: var(--primary);
  font-size: 16px;
}

.help-content {
  padding: 25px;
}

.help-content p {
  color: var(--text-secondary);
  line-height: 1.6;
  margin-bottom: 15px;
}

.help-content ol {
  color: var(--text-secondary);
  padding-left: 20px;
}

.help-content li {
  margin-bottom: 8px;
  line-height: 1.6;
}

.help-content strong {
  color: var(--text-primary);
}

/* 支持部分 */
.support-section {
  padding: 80px 0;
  background: rgba(15, 23, 42, 0.5);
}

.support-content {
  text-align: center;
  max-width: 600px;
  margin: 0 auto;
}

.support-content h2 {
  font-size: 32px;
  font-weight: 700;
  margin-bottom: 20px;
  background: linear-gradient(to right, var(--primary), var(--primary-light));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.support-content p {
  font-size: 18px;
  color: var(--text-secondary);
  margin-bottom: 40px;
  line-height: 1.6;
}

.support-actions {
  display: flex;
  gap: 20px;
  justify-content: center;
}

.btn {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  padding: 15px 30px;
  border-radius: 8px;
  text-decoration: none;
  font-weight: 600;
  transition: all 0.3s ease;
}

.btn-primary {
  background: var(--primary);
  color: white;
  border: 1px solid var(--primary);
}

.btn-primary:hover {
  background: var(--primary-dark);
  transform: translateY(-2px);
}

.btn-secondary {
  background: transparent;
  color: var(--primary);
  border: 1px solid var(--primary);
}

.btn-secondary:hover {
  background: var(--primary);
  color: white;
}

/* 底部 */
.footer {
  background: rgba(15, 23, 42, 0.9);
  border-top: 1px solid var(--border);
  padding: 40px 0;
  text-align: center;
}

.footer p {
  color: var(--text-secondary);
  margin: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .container {
    padding: 0 20px;
  }

  .top-nav {
    padding: 15px 20px;
  }

  .search-title {
    font-size: 32px;
  }

  .search-box {
    flex-direction: column;
    gap: 10px;
    padding: 15px;
  }

  .search-box input {
    text-align: center;
  }

  .support-actions {
    flex-direction: column;
    align-items: center;
  }

  .quick-nav-grid {
    grid-template-columns: 1fr;
  }
}
</style>
