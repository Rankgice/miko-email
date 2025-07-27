<template>
  <div class="cookie-policy-page">
    <!-- 顶部导航 -->
    <nav class="top-nav">
      <div class="logo">
        <i class="fas fa-cookie-bite"></i>
        <span>Cookie政策</span>
      </div>
      <router-link to="/" class="back-home">
        <i class="fas fa-home"></i>
        返回首页
      </router-link>
    </nav>

    <!-- 主要内容 -->
    <main class="main-content">
      <div class="container">
        <div class="policy-content">
          <h1>Cookie政策</h1>
          <p class="last-updated">最后更新时间：2024年1月1日</p>

          <section class="policy-section">
            <h2>1. 什么是Cookie</h2>
            <p>
              Cookie是当您访问网站时存储在您设备上的小型文本文件。它们被广泛用于使网站工作，
              或更高效地工作，以及向网站所有者提供信息。
            </p>
          </section>

          <section class="policy-section">
            <h2>2. 我们如何使用Cookie</h2>
            <p>Miko邮箱系统使用Cookie来：</p>
            <ul>
              <li>保持您的登录状态</li>
              <li>记住您的偏好设置</li>
              <li>提供个性化体验</li>
              <li>分析网站使用情况</li>
              <li>改进我们的服务</li>
              <li>确保网站安全</li>
            </ul>
          </section>

          <section class="policy-section">
            <h2>3. Cookie类型</h2>
            
            <h3>3.1 必要Cookie</h3>
            <p>这些Cookie对于网站的基本功能是必需的，无法禁用：</p>
            <ul>
              <li><strong>会话Cookie：</strong>维护您的登录状态</li>
              <li><strong>安全Cookie：</strong>防止跨站请求伪造攻击</li>
              <li><strong>负载均衡Cookie：</strong>确保请求路由到正确的服务器</li>
            </ul>

            <h3>3.2 功能Cookie</h3>
            <p>这些Cookie增强网站功能和个性化：</p>
            <ul>
              <li><strong>偏好Cookie：</strong>记住您的语言和主题设置</li>
              <li><strong>界面Cookie：</strong>保存您的界面布局偏好</li>
              <li><strong>表单Cookie：</strong>记住表单输入以改善用户体验</li>
            </ul>

            <h3>3.3 分析Cookie</h3>
            <p>这些Cookie帮助我们了解网站使用情况：</p>
            <ul>
              <li><strong>访问统计：</strong>页面访问量和用户行为</li>
              <li><strong>性能监控：</strong>页面加载时间和错误率</li>
              <li><strong>使用模式：</strong>功能使用频率和用户路径</li>
            </ul>
          </section>

          <section class="policy-section">
            <h2>4. 第三方Cookie</h2>
            <p>我们可能使用以下第三方服务，它们可能设置自己的Cookie：</p>
            
            <div class="third-party-service">
              <h4>Google Analytics</h4>
              <p>用于网站分析和性能监控</p>
              <ul>
                <li>收集匿名使用数据</li>
                <li>帮助改进网站性能</li>
                <li>可通过浏览器设置禁用</li>
              </ul>
            </div>

            <div class="third-party-service">
              <h4>CDN服务</h4>
              <p>用于加速内容分发</p>
              <ul>
                <li>提高页面加载速度</li>
                <li>优化用户体验</li>
                <li>不收集个人信息</li>
              </ul>
            </div>
          </section>

          <section class="policy-section">
            <h2>5. Cookie管理</h2>
            
            <h3>5.1 浏览器设置</h3>
            <p>您可以通过浏览器设置管理Cookie：</p>
            <ul>
              <li><strong>Chrome：</strong>设置 > 隐私和安全 > Cookie和其他网站数据</li>
              <li><strong>Firefox：</strong>选项 > 隐私与安全 > Cookie和网站数据</li>
              <li><strong>Safari：</strong>偏好设置 > 隐私 > Cookie和网站数据</li>
              <li><strong>Edge：</strong>设置 > Cookie和网站权限</li>
            </ul>

            <h3>5.2 Cookie偏好中心</h3>
            <div class="cookie-preferences">
              <h4>管理您的Cookie偏好</h4>
              <div class="preference-item">
                <label class="preference-label">
                  <input type="checkbox" checked disabled>
                  <span class="checkmark"></span>
                  <div class="preference-info">
                    <strong>必要Cookie</strong>
                    <p>网站基本功能所需，无法禁用</p>
                  </div>
                </label>
              </div>

              <div class="preference-item">
                <label class="preference-label">
                  <input type="checkbox" v-model="preferences.functional">
                  <span class="checkmark"></span>
                  <div class="preference-info">
                    <strong>功能Cookie</strong>
                    <p>增强网站功能和个性化体验</p>
                  </div>
                </label>
              </div>

              <div class="preference-item">
                <label class="preference-label">
                  <input type="checkbox" v-model="preferences.analytics">
                  <span class="checkmark"></span>
                  <div class="preference-info">
                    <strong>分析Cookie</strong>
                    <p>帮助我们了解网站使用情况</p>
                  </div>
                </label>
              </div>

              <div class="preference-actions">
                <button @click="savePreferences" class="save-btn">
                  <i class="fas fa-save"></i>
                  保存偏好
                </button>
                <button @click="acceptAll" class="accept-btn">
                  <i class="fas fa-check"></i>
                  接受全部
                </button>
              </div>
            </div>
          </section>

          <section class="policy-section">
            <h2>6. Cookie详细信息</h2>
            <div class="cookie-table">
              <table>
                <thead>
                  <tr>
                    <th>Cookie名称</th>
                    <th>类型</th>
                    <th>用途</th>
                    <th>有效期</th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td>session_id</td>
                    <td>必要</td>
                    <td>维护用户会话</td>
                    <td>会话结束</td>
                  </tr>
                  <tr>
                    <td>csrf_token</td>
                    <td>必要</td>
                    <td>防止CSRF攻击</td>
                    <td>会话结束</td>
                  </tr>
                  <tr>
                    <td>user_preferences</td>
                    <td>功能</td>
                    <td>保存用户偏好</td>
                    <td>1年</td>
                  </tr>
                  <tr>
                    <td>theme_setting</td>
                    <td>功能</td>
                    <td>记住主题选择</td>
                    <td>1年</td>
                  </tr>
                  <tr>
                    <td>_ga</td>
                    <td>分析</td>
                    <td>Google Analytics</td>
                    <td>2年</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </section>

          <section class="policy-section">
            <h2>7. 禁用Cookie的影响</h2>
            <p>如果您选择禁用某些Cookie，可能会影响网站功能：</p>
            <ul>
              <li><strong>必要Cookie：</strong>禁用会导致网站无法正常工作</li>
              <li><strong>功能Cookie：</strong>禁用会失去个性化设置</li>
              <li><strong>分析Cookie：</strong>禁用不会影响网站功能</li>
            </ul>
          </section>

          <section class="policy-section">
            <h2>8. 移动设备</h2>
            <p>
              在移动设备上，Cookie的工作方式可能有所不同。您可以通过设备设置或应用内设置
              管理Cookie偏好。
            </p>
          </section>

          <section class="policy-section">
            <h2>9. 政策更新</h2>
            <p>
              我们可能会不时更新此Cookie政策。重大变更会通过网站通知您。
              继续使用网站即表示您接受更新后的政策。
            </p>
          </section>

          <section class="policy-section">
            <h2>10. 联系我们</h2>
            <p>如果您对我们的Cookie使用有任何问题，请联系我们：</p>
            <div class="contact-info">
              <p><strong>邮箱：</strong> privacy@miko-email.com</p>
              <p><strong>电话：</strong> 400-123-4567</p>
            </div>
          </section>
        </div>
      </div>
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

const preferences = ref({
  functional: true,
  analytics: true
})

const savePreferences = () => {
  // 保存Cookie偏好到localStorage或发送到服务器
  localStorage.setItem('cookiePreferences', JSON.stringify(preferences.value))
  alert('Cookie偏好已保存')
}

const acceptAll = () => {
  preferences.value = {
    functional: true,
    analytics: true
  }
  savePreferences()
}
</script>

<style scoped>
/* 基础样式与隐私政策相同，添加特定样式 */
.cookie-policy-page {
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
  color: var(--text-primary);
  min-height: 100vh;
}

/* 复用其他页面的基础样式 */
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

.main-content {
  padding-top: 80px;
}

.container {
  max-width: 900px;
  margin: 0 auto;
  padding: 40px;
}

.policy-content {
  background: linear-gradient(135deg, rgba(30, 41, 59, 0.9), rgba(15, 23, 42, 0.95));
  border: 1px solid var(--border);
  border-radius: 16px;
  padding: 50px 40px;
}

.policy-content h1 {
  font-size: 36px;
  font-weight: 700;
  margin-bottom: 10px;
  background: linear-gradient(to right, var(--primary), var(--primary-light));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  text-align: center;
}

.last-updated {
  text-align: center;
  color: var(--text-secondary);
  margin-bottom: 40px;
  font-style: italic;
}

.policy-section {
  margin-bottom: 40px;
}

.policy-section h2 {
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 20px;
  color: var(--text-primary);
  border-bottom: 2px solid var(--border);
  padding-bottom: 10px;
}

.policy-section h3 {
  font-size: 18px;
  font-weight: 600;
  margin: 20px 0 10px 0;
  color: var(--primary);
}

.policy-section h4 {
  font-size: 16px;
  font-weight: 600;
  margin: 15px 0 8px 0;
  color: var(--text-primary);
}

.policy-section p {
  color: var(--text-secondary);
  line-height: 1.8;
  margin-bottom: 15px;
}

.policy-section ul {
  color: var(--text-secondary);
  line-height: 1.8;
  padding-left: 20px;
  margin-bottom: 15px;
}

.policy-section li {
  margin-bottom: 8px;
}

.policy-section strong {
  color: var(--text-primary);
}

/* 第三方服务样式 */
.third-party-service {
  background: rgba(0, 180, 216, 0.05);
  border: 1px solid rgba(0, 180, 216, 0.2);
  border-radius: 8px;
  padding: 20px;
  margin: 15px 0;
}

/* Cookie偏好中心 */
.cookie-preferences {
  background: rgba(0, 180, 216, 0.1);
  border: 1px solid var(--primary);
  border-radius: 12px;
  padding: 25px;
  margin: 20px 0;
}

.preference-item {
  margin-bottom: 20px;
}

.preference-label {
  display: flex;
  align-items: flex-start;
  gap: 15px;
  cursor: pointer;
}

.preference-label input[type="checkbox"] {
  width: 18px;
  height: 18px;
  margin: 0;
}

.preference-info strong {
  display: block;
  color: var(--text-primary);
  margin-bottom: 5px;
}

.preference-info p {
  margin: 0;
  font-size: 14px;
  color: var(--text-secondary);
}

.preference-actions {
  display: flex;
  gap: 15px;
  margin-top: 25px;
}

.save-btn, .accept-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 24px;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.save-btn {
  background: var(--primary);
  color: white;
}

.accept-btn {
  background: transparent;
  color: var(--primary);
  border: 1px solid var(--primary);
}

.save-btn:hover, .accept-btn:hover {
  transform: translateY(-2px);
}

.accept-btn:hover {
  background: var(--primary);
  color: white;
}

/* Cookie表格 */
.cookie-table {
  margin: 20px 0;
  overflow-x: auto;
}

.cookie-table table {
  width: 100%;
  border-collapse: collapse;
  background: rgba(15, 23, 42, 0.5);
  border-radius: 8px;
  overflow: hidden;
}

.cookie-table th,
.cookie-table td {
  padding: 12px 15px;
  text-align: left;
  border-bottom: 1px solid var(--border);
}

.cookie-table th {
  background: rgba(0, 180, 216, 0.1);
  color: var(--primary);
  font-weight: 600;
}

.cookie-table td {
  color: var(--text-secondary);
}

.contact-info {
  background: rgba(0, 180, 216, 0.1);
  border: 1px solid var(--primary);
  border-radius: 8px;
  padding: 20px;
  margin-top: 15px;
}

.contact-info p {
  margin-bottom: 8px;
  color: var(--text-primary);
}

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

@media (max-width: 768px) {
  .container {
    padding: 20px;
  }

  .policy-content {
    padding: 30px 20px;
  }

  .policy-content h1 {
    font-size: 28px;
  }

  .top-nav {
    padding: 15px 20px;
  }

  .preference-actions {
    flex-direction: column;
  }
}
</style>
