<template>
  <div class="home-page">
    <!-- 装饰元素 -->
    <div class="decoration circle1"></div>
    <div class="decoration circle2"></div>

    <!-- 浮动图标 -->
    <i class="fas fa-envelope floating-icon floating-1"></i>
    <i class="fas fa-shield-alt floating-icon floating-2"></i>
    <i class="fas fa-cog floating-icon floating-3"></i>

    <!-- 顶部导航 -->
    <nav class="top-nav">
      <div class="logo">
        <i class="fas fa-envelope"></i>
        <span>Miko邮箱系统</span>
      </div>

      <div class="nav-links">
        <router-link to="/user/login" class="nav-link">
          <i class="fas fa-user"></i>
          <span>用户登录</span>
        </router-link>
        <router-link to="/admin/login" class="nav-link admin">
          <i class="fas fa-shield-alt"></i>
          <span>管理员</span>
        </router-link>
      </div>
    </nav>

    <!-- 主要内容 -->
    <main class="main-content">
      <!-- 英雄区域 -->
      <section class="hero-section">
        <div class="hero-content">
          <h1 class="hero-title">
            欢迎使用
            <span class="highlight">Miko邮箱系统</span>
          </h1>
          <p class="hero-subtitle">
            现代化的邮件管理解决方案，为您提供安全、高效、易用的邮件服务
          </p>
          <div class="hero-actions">
            <router-link to="/user/register" class="btn btn-primary">
              <i class="fas fa-user-plus"></i>
              立即注册
            </router-link>
            <router-link to="/user/login" class="btn btn-secondary">
              <i class="fas fa-sign-in-alt"></i>
              用户登录
            </router-link>
          </div>
        </div>

        <div class="hero-image">
          <div class="email-preview">
            <div class="email-header">
              <div class="email-controls">
                <span class="control red"></span>
                <span class="control yellow"></span>
                <span class="control green"></span>
              </div>
              <span class="email-title">Miko邮箱</span>
            </div>
            <div class="email-content">
              <div class="email-item">
                <div class="email-avatar">M</div>
                <div class="email-info">
                  <h4>欢迎使用Miko邮箱</h4>
                  <p>感谢您选择我们的邮件服务...</p>
                </div>
              </div>
              <div class="email-item">
                <div class="email-avatar">S</div>
                <div class="email-info">
                  <h4>系统通知</h4>
                  <p>您的邮箱已成功激活...</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- 特性介绍 -->
      <section class="features-section">
        <div class="container">
          <h2 class="section-title">核心特性</h2>
          <div class="features-grid">
            <div class="feature-card">
              <div class="feature-icon">
                <i class="fas fa-shield-alt"></i>
              </div>
              <h3>安全可靠</h3>
              <p>采用先进的加密技术，确保您的邮件数据安全</p>
            </div>

            <div class="feature-card">
              <div class="feature-icon">
                <i class="fas fa-bolt"></i>
              </div>
              <h3>高效快速</h3>
              <p>优化的服务器架构，提供快速的邮件收发体验</p>
            </div>

            <div class="feature-card">
              <div class="feature-icon">
                <i class="fas fa-mobile-alt"></i>
              </div>
              <h3>多端同步</h3>
              <p>支持多设备访问，随时随地管理您的邮件</p>
            </div>

            <div class="feature-card">
              <div class="feature-icon">
                <i class="fas fa-cog"></i>
              </div>
              <h3>易于管理</h3>
              <p>直观的管理界面，轻松配置和管理邮箱设置</p>
            </div>
          </div>
        </div>
      </section>

      <!-- 快速入口 -->
      <section class="quick-access-section">
        <div class="container">
          <h2 class="section-title">快速入口</h2>
          <div class="access-grid">
            <router-link to="/user/login" class="access-card user">
              <div class="access-icon">
                <i class="fas fa-user"></i>
              </div>
              <h3>用户入口</h3>
              <p>登录您的邮箱账户，管理邮件和设置</p>
              <span class="access-arrow">
                <i class="fas fa-arrow-right"></i>
              </span>
            </router-link>

            <router-link to="/admin/login" class="access-card admin">
              <div class="access-icon">
                <i class="fas fa-shield-alt"></i>
              </div>
              <h3>管理员入口</h3>
              <p>系统管理员登录，管理用户和系统配置</p>
              <span class="access-arrow">
                <i class="fas fa-arrow-right"></i>
              </span>
            </router-link>
          </div>
        </div>
      </section>
    </main>

    <!-- 底部信息 -->
    <footer class="footer">
      <div class="container">
        <div class="footer-content">
          <div class="footer-info">
            <div class="footer-logo">
              <i class="fas fa-envelope"></i>
              <span>Miko邮箱系统</span>
            </div>
            <p>现代化的邮件管理解决方案</p>
          </div>

          <div class="footer-links">
            <div class="link-group">
              <h4>产品</h4>
              <router-link to="/email-service">邮箱服务</router-link>
              <router-link to="/api-docs">API文档</router-link>
            </div>

            <div class="link-group">
              <h4>支持</h4>
              <router-link to="/help-center">帮助中心</router-link>
              <router-link to="/contact">联系我们</router-link>
              <router-link to="/tech-support">技术支持</router-link>
            </div>

            <div class="link-group">
              <h4>法律</h4>
              <router-link to="/privacy-policy">隐私政策</router-link>
              <router-link to="/terms-of-service">服务条款</router-link>
              <router-link to="/cookie-policy">Cookie政策</router-link>
            </div>
          </div>
        </div>

        <div class="footer-bottom">
          <p>&copy; 2024 Miko邮箱系统. 保留所有权利.</p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
// 首页组件
</script>

<style scoped>
.home-page {
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
  color: var(--text-primary);
  min-height: 100vh;
  overflow-x: hidden;
  position: relative;
}

/* 装饰元素 */
.decoration {
  position: fixed;
  z-index: -1;
  opacity: 0.05;
}

.circle1 {
  width: 600px;
  height: 600px;
  border-radius: 50%;
  background: radial-gradient(circle, var(--primary), transparent);
  top: -300px;
  right: -200px;
}

.circle2 {
  width: 400px;
  height: 400px;
  border-radius: 50%;
  background: radial-gradient(circle, var(--accent), transparent);
  bottom: -200px;
  left: -150px;
}

/* 浮动图标 */
@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-15px); }
}

.floating-icon {
  position: absolute;
  font-size: 150px;
  opacity: 0.08;
  z-index: -1;
  animation: float 8s ease-in-out infinite;
}

.floating-1 {
  top: 15%;
  left: 8%;
  color: var(--primary);
}

.floating-2 {
  bottom: 20%;
  right: 10%;
  color: var(--accent);
  animation-delay: 2s;
}

.floating-3 {
  top: 45%;
  right: 15%;
  color: var(--info);
  animation-delay: 4s;
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
  background: rgba(30, 41, 59, 0.8);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid var(--border);
  z-index: 100;
}

.logo {
  font-size: 28px;
  font-weight: 700;
  background: linear-gradient(to right, var(--primary), var(--primary-light));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  display: flex;
  align-items: center;
  gap: 15px;
}

.logo i {
  font-size: 32px;
}

.nav-links {
  display: flex;
  gap: 25px;
}

.nav-link {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  border-radius: 8px;
  color: var(--text-primary);
  text-decoration: none;
  font-weight: 500;
  transition: all 0.3s ease;
  border: 1px solid transparent;
}

.nav-link:hover {
  background: rgba(0, 180, 216, 0.1);
  border-color: var(--primary);
  color: var(--primary);
}

.nav-link.admin:hover {
  background: rgba(255, 107, 107, 0.1);
  border-color: var(--accent);
  color: var(--accent);
}

/* 主要内容 */
.main-content {
  padding-top: 100px;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 40px;
}

/* 英雄区域 */
.hero-section {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 60px;
  align-items: center;
  min-height: 80vh;
  padding: 60px 40px;
}

.hero-content {
  z-index: 10;
}

.hero-title {
  font-size: 48px;
  font-weight: 700;
  line-height: 1.2;
  margin-bottom: 20px;
}

.highlight {
  background: linear-gradient(to right, var(--primary), var(--primary-light));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.hero-subtitle {
  font-size: 20px;
  color: var(--text-secondary);
  line-height: 1.6;
  margin-bottom: 40px;
}

.hero-actions {
  display: flex;
  gap: 20px;
}

.btn {
  padding: 16px 32px;
  border-radius: 12px;
  border: none;
  font-weight: 600;
  font-size: 16px;
  cursor: pointer;
  transition: all 0.3s ease;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  gap: 10px;
}

.btn-primary {
  background: linear-gradient(135deg, var(--primary), #0077b6);
  color: white;
  box-shadow: 0 8px 25px rgba(0, 180, 216, 0.3);
}

.btn-primary:hover {
  transform: translateY(-3px);
  box-shadow: 0 12px 35px rgba(0, 180, 216, 0.4);
}

.btn-secondary {
  background: rgba(15, 23, 42, 0.8);
  color: var(--text-primary);
  border: 2px solid var(--border);
}

.btn-secondary:hover {
  background: rgba(15, 23, 42, 0.9);
  border-color: var(--primary);
  color: var(--primary);
}

/* 邮件预览 */
.hero-image {
  display: flex;
  justify-content: center;
  align-items: center;
}

.email-preview {
  width: 400px;
  background: linear-gradient(135deg, rgba(30, 41, 59, 0.9), rgba(15, 23, 42, 0.95));
  border-radius: 16px;
  border: 1px solid var(--border);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
  overflow: hidden;
  transform: perspective(1000px) rotateY(-5deg) rotateX(5deg);
}

.email-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 15px 20px;
  background: rgba(15, 23, 42, 0.8);
  border-bottom: 1px solid var(--border);
}

.email-controls {
  display: flex;
  gap: 8px;
}

.control {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.control.red { background: #ff5f57; }
.control.yellow { background: #ffbd2e; }
.control.green { background: #28ca42; }

.email-title {
  font-weight: 600;
  color: var(--text-primary);
}

.email-content {
  padding: 20px;
}

.email-item {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 15px 0;
  border-bottom: 1px solid var(--border);
}

.email-item:last-child {
  border-bottom: none;
}

.email-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--primary), #0077b6);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: 600;
}

.email-info h4 {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 4px;
}

.email-info p {
  font-size: 12px;
  color: var(--text-secondary);
  margin: 0;
}

/* 特性介绍 */
.features-section {
  padding: 100px 0;
  background: rgba(15, 23, 42, 0.5);
}

.section-title {
  text-align: center;
  font-size: 36px;
  font-weight: 700;
  margin-bottom: 60px;
  background: linear-gradient(to right, var(--primary), var(--primary-light));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 40px;
}

.feature-card {
  text-align: center;
  padding: 40px 20px;
  background: linear-gradient(135deg, rgba(30, 41, 59, 0.9), rgba(15, 23, 42, 0.95));
  border-radius: 16px;
  border: 1px solid var(--border);
  transition: all 0.3s ease;
}

.feature-card:hover {
  transform: translateY(-10px);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.2);
}

.feature-icon {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--primary), #0077b6);
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 20px;
  font-size: 32px;
  color: white;
}

.feature-card h3 {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 15px;
  color: var(--text-primary);
}

.feature-card p {
  color: var(--text-secondary);
  line-height: 1.6;
}

/* 快速入口 */
.quick-access-section {
  padding: 100px 0;
}

.access-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 40px;
}

.access-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  padding: 50px 30px;
  background: linear-gradient(135deg, rgba(30, 41, 59, 0.9), rgba(15, 23, 42, 0.95));
  border-radius: 20px;
  border: 2px solid var(--border);
  text-decoration: none;
  color: var(--text-primary);
  transition: all 0.3s ease;
  position: relative;
  overflow: hidden;
}

.access-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 4px;
  background: linear-gradient(90deg, var(--primary), var(--primary-light));
}

.access-card.admin::before {
  background: linear-gradient(90deg, var(--accent), #ff8a80);
}

.access-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 25px 50px rgba(0, 0, 0, 0.25);
  border-color: var(--primary);
}

.access-card.admin:hover {
  border-color: var(--accent);
}

.access-icon {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--primary), #0077b6);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 25px;
  font-size: 40px;
  color: white;
}

.access-card.admin .access-icon {
  background: linear-gradient(135deg, var(--accent), #ff8a80);
}

.access-card h3 {
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 15px;
}

.access-card p {
  color: var(--text-secondary);
  line-height: 1.6;
  margin-bottom: 20px;
}

.access-arrow {
  font-size: 20px;
  color: var(--primary);
  transition: transform 0.3s ease;
}

.access-card:hover .access-arrow {
  transform: translateX(5px);
}

/* 底部信息 */
.footer {
  background: rgba(15, 23, 42, 0.9);
  border-top: 1px solid var(--border);
  padding: 60px 0 20px;
}

.footer-content {
  display: grid;
  grid-template-columns: 1fr 2fr;
  gap: 60px;
  margin-bottom: 40px;
}

.footer-logo {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 15px;
  background: linear-gradient(to right, var(--primary), var(--primary-light));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.footer-info p {
  color: var(--text-secondary);
  line-height: 1.6;
}

.footer-links {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 40px;
}

.link-group h4 {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 20px;
  color: var(--text-primary);
}

.link-group a {
  display: block;
  color: var(--text-secondary);
  text-decoration: none;
  margin-bottom: 10px;
  transition: color 0.3s ease;
}

.link-group a:hover {
  color: var(--primary);
}

.footer-bottom {
  text-align: center;
  padding-top: 20px;
  border-top: 1px solid var(--border);
  color: var(--text-secondary);
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .hero-section {
    grid-template-columns: 1fr;
    text-align: center;
    gap: 40px;
  }

  .hero-title {
    font-size: 40px;
  }

  .features-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .container {
    padding: 0 20px;
  }

  .top-nav {
    padding: 15px 20px;
    flex-direction: column;
    gap: 20px;
  }

  .hero-section {
    padding: 40px 20px;
  }

  .hero-title {
    font-size: 32px;
  }

  .hero-actions {
    flex-direction: column;
    align-items: center;
  }

  .features-grid {
    grid-template-columns: 1fr;
  }

  .access-grid {
    grid-template-columns: 1fr;
  }

  .footer-content {
    grid-template-columns: 1fr;
    gap: 40px;
  }

  .footer-links {
    grid-template-columns: 1fr;
    gap: 30px;
  }
}

@media (max-width: 480px) {
  .hero-title {
    font-size: 28px;
  }

  .hero-subtitle {
    font-size: 18px;
  }

  .btn {
    padding: 14px 24px;
    font-size: 14px;
  }

  .email-preview {
    width: 300px;
  }

  .floating-icon {
    display: none;
  }
}
</style>
