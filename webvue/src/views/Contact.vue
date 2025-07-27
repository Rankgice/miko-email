<template>
  <div class="contact-page">
    <!-- 顶部导航 -->
    <nav class="top-nav">
      <div class="logo">
        <i class="fas fa-phone"></i>
        <span>联系我们</span>
      </div>
      <router-link to="/" class="back-home">
        <i class="fas fa-home"></i>
        返回首页
      </router-link>
    </nav>

    <!-- 主要内容 -->
    <main class="main-content">
      <!-- 英雄区域 -->
      <section class="hero-section">
        <div class="container">
          <h1 class="hero-title">联系我们</h1>
          <p class="hero-subtitle">
            我们随时为您提供专业的服务和支持，欢迎通过以下方式与我们取得联系
          </p>
        </div>
      </section>

      <!-- 联系方式 -->
      <section class="contact-info-section">
        <div class="container">
          <div class="contact-grid">
            <div class="contact-card">
              <div class="contact-icon">
                <i class="fas fa-envelope"></i>
              </div>
              <h3>邮箱联系</h3>
              <p>发送邮件给我们，我们会在24小时内回复</p>
              <div class="contact-details">
                <a href="mailto:support@miko-email.com">support@miko-email.com</a>
                <a href="mailto:sales@miko-email.com">sales@miko-email.com</a>
              </div>
            </div>

            <div class="contact-card">
              <div class="contact-icon">
                <i class="fas fa-phone"></i>
              </div>
              <h3>电话联系</h3>
              <p>工作时间：周一至周五 9:00-18:00</p>
              <div class="contact-details">
                <a href="tel:400-123-4567">400-123-4567</a>
                <span>技术支持热线</span>
              </div>
            </div>

            <div class="contact-card">
              <div class="contact-icon">
                <i class="fas fa-map-marker-alt"></i>
              </div>
              <h3>办公地址</h3>
              <p>欢迎预约到访，我们为您提供面对面的服务</p>
              <div class="contact-details">
                <span>北京市朝阳区科技园区</span>
                <span>创新大厦A座15层</span>
              </div>
            </div>

            <div class="contact-card">
              <div class="contact-icon">
                <i class="fas fa-comments"></i>
              </div>
              <h3>在线客服</h3>
              <p>实时在线支持，快速解决您的问题</p>
              <div class="contact-details">
                <button class="chat-btn" @click="startChat">
                  <i class="fas fa-comment-dots"></i>
                  开始对话
                </button>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- 联系表单 -->
      <section class="contact-form-section">
        <div class="container">
          <div class="form-container">
            <div class="form-header">
              <h2>发送消息</h2>
              <p>填写下面的表单，我们会尽快与您联系</p>
            </div>

            <form @submit.prevent="submitForm" class="contact-form">
              <div class="form-row">
                <div class="form-group">
                  <label for="name">姓名 *</label>
                  <input
                    type="text"
                    id="name"
                    v-model="form.name"
                    required
                    placeholder="请输入您的姓名"
                  >
                </div>
                <div class="form-group">
                  <label for="email">邮箱 *</label>
                  <input
                    type="email"
                    id="email"
                    v-model="form.email"
                    required
                    placeholder="请输入您的邮箱"
                  >
                </div>
              </div>

              <div class="form-row">
                <div class="form-group">
                  <label for="phone">电话</label>
                  <input
                    type="tel"
                    id="phone"
                    v-model="form.phone"
                    placeholder="请输入您的电话号码"
                  >
                </div>
                <div class="form-group">
                  <label for="company">公司</label>
                  <input
                    type="text"
                    id="company"
                    v-model="form.company"
                    placeholder="请输入您的公司名称"
                  >
                </div>
              </div>

              <div class="form-group">
                <label for="subject">主题 *</label>
                <select id="subject" v-model="form.subject" required>
                  <option value="">请选择咨询主题</option>
                  <option value="sales">销售咨询</option>
                  <option value="support">技术支持</option>
                  <option value="billing">账单问题</option>
                  <option value="feature">功能建议</option>
                  <option value="other">其他问题</option>
                </select>
              </div>

              <div class="form-group">
                <label for="message">详细描述 *</label>
                <textarea
                  id="message"
                  v-model="form.message"
                  required
                  rows="6"
                  placeholder="请详细描述您的问题或需求..."
                ></textarea>
              </div>

              <div class="form-group checkbox-group">
                <label class="checkbox-label">
                  <input type="checkbox" v-model="form.agree" required>
                  <span class="checkmark"></span>
                  我同意接收相关的产品信息和服务更新
                </label>
              </div>

              <button type="submit" class="submit-btn" :disabled="submitting">
                <i class="fas fa-paper-plane"></i>
                {{ submitting ? '发送中...' : '发送消息' }}
              </button>
            </form>
          </div>
        </div>
      </section>

      <!-- 常见问题 -->
      <section class="faq-section">
        <div class="container">
          <h2 class="section-title">常见问题</h2>
          <div class="faq-grid">
            <div class="faq-item">
              <h3>
                <i class="fas fa-question-circle"></i>
                如何快速获得技术支持？
              </h3>
              <p>您可以通过在线客服、技术支持热线或发送邮件的方式联系我们。对于紧急问题，建议直接拨打技术支持热线。</p>
            </div>

            <div class="faq-item">
              <h3>
                <i class="fas fa-clock"></i>
                客服工作时间是什么时候？
              </h3>
              <p>我们的客服团队工作时间为周一至周五 9:00-18:00。在非工作时间，您可以通过邮件联系我们，我们会在下一个工作日回复。</p>
            </div>

            <div class="faq-item">
              <h3>
                <i class="fas fa-handshake"></i>
                如何申请产品演示？
              </h3>
              <p>您可以通过联系表单选择"销售咨询"主题，或直接拨打销售热线预约产品演示。我们的销售团队会为您安排专业的产品介绍。</p>
            </div>

            <div class="faq-item">
              <h3>
                <i class="fas fa-shield-alt"></i>
                如何报告安全问题？
              </h3>
              <p>如果您发现安全相关问题，请直接发送邮件至 security@miko-email.com，我们会优先处理并及时回复。</p>
            </div>
          </div>
        </div>
      </section>

      <!-- 其他联系方式 -->
      <section class="social-section">
        <div class="container">
          <h2 class="section-title">关注我们</h2>
          <div class="social-links">
            <a href="#" class="social-link">
              <i class="fab fa-weixin"></i>
              <span>微信公众号</span>
            </a>
            <a href="#" class="social-link">
              <i class="fab fa-weibo"></i>
              <span>新浪微博</span>
            </a>
            <a href="#" class="social-link">
              <i class="fab fa-qq"></i>
              <span>QQ群</span>
            </a>
            <a href="#" class="social-link">
              <i class="fab fa-github"></i>
              <span>GitHub</span>
            </a>
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

const submitting = ref(false)

const form = ref({
  name: '',
  email: '',
  phone: '',
  company: '',
  subject: '',
  message: '',
  agree: false
})

const submitForm = async () => {
  submitting.value = true

  try {
    // 这里应该调用API发送联系表单
    await new Promise(resolve => setTimeout(resolve, 2000)) // 模拟API调用

    alert('消息发送成功！我们会尽快与您联系。')

    // 重置表单
    form.value = {
      name: '',
      email: '',
      phone: '',
      company: '',
      subject: '',
      message: '',
      agree: false
    }
  } catch (error) {
    alert('发送失败，请稍后重试。')
  } finally {
    submitting.value = false
  }
}

const startChat = () => {
  // 这里可以集成在线客服系统
  alert('在线客服功能即将开放，请通过其他方式联系我们。')
}
</script>

<style scoped>
.contact-page {
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

/* 英雄区域 */
.hero-section {
  padding: 80px 0;
  text-align: center;
}

.hero-title {
  font-size: 48px;
  font-weight: 700;
  margin-bottom: 20px;
  background: linear-gradient(to right, var(--primary), var(--primary-light));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.hero-subtitle {
  font-size: 20px;
  color: var(--text-secondary);
  max-width: 600px;
  margin: 0 auto;
  line-height: 1.6;
}

/* 联系信息 */
.contact-info-section {
  padding: 60px 0;
  background: rgba(15, 23, 42, 0.5);
}

.contact-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 30px;
}

.contact-card {
  background: linear-gradient(135deg, rgba(30, 41, 59, 0.9), rgba(15, 23, 42, 0.95));
  border: 1px solid var(--border);
  border-radius: 16px;
  padding: 40px 30px;
  text-align: center;
  transition: all 0.3s ease;
}

.contact-card:hover {
  transform: translateY(-5px);
  border-color: var(--primary);
  box-shadow: 0 20px 40px rgba(0, 180, 216, 0.1);
}

.contact-icon {
  width: 80px;
  height: 80px;
  background: linear-gradient(135deg, var(--primary), var(--primary-light));
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 25px;
  font-size: 32px;
  color: white;
}

.contact-card h3 {
  font-size: 22px;
  font-weight: 600;
  margin-bottom: 15px;
  color: var(--text-primary);
}

.contact-card p {
  color: var(--text-secondary);
  margin-bottom: 20px;
  line-height: 1.6;
}

.contact-details {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.contact-details a {
  color: var(--primary);
  text-decoration: none;
  font-weight: 500;
  transition: all 0.3s ease;
}

.contact-details a:hover {
  color: var(--primary-light);
}

.contact-details span {
  color: var(--text-secondary);
  font-size: 14px;
}

.chat-btn {
  background: var(--primary);
  color: white;
  border: none;
  border-radius: 8px;
  padding: 12px 24px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.chat-btn:hover {
  background: var(--primary-dark);
  transform: translateY(-2px);
}

/* 联系表单 */
.contact-form-section {
  padding: 80px 0;
}

.form-container {
  max-width: 800px;
  margin: 0 auto;
  background: linear-gradient(135deg, rgba(30, 41, 59, 0.9), rgba(15, 23, 42, 0.95));
  border: 1px solid var(--border);
  border-radius: 20px;
  padding: 50px 40px;
}

.form-header {
  text-align: center;
  margin-bottom: 40px;
}

.form-header h2 {
  font-size: 32px;
  font-weight: 700;
  margin-bottom: 15px;
  background: linear-gradient(to right, var(--primary), var(--primary-light));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.form-header p {
  color: var(--text-secondary);
  font-size: 16px;
}

.contact-form {
  display: flex;
  flex-direction: column;
  gap: 25px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group label {
  color: var(--text-primary);
  font-weight: 500;
  font-size: 14px;
}

.form-group input,
.form-group select,
.form-group textarea {
  background: rgba(15, 23, 42, 0.8);
  border: 1px solid var(--border);
  border-radius: 8px;
  padding: 12px 16px;
  color: var(--text-primary);
  font-size: 14px;
  transition: all 0.3s ease;
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(0, 180, 216, 0.1);
}

.form-group input::placeholder,
.form-group textarea::placeholder {
  color: var(--text-secondary);
}

.checkbox-group {
  flex-direction: row;
  align-items: center;
  gap: 12px;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 12px;
  color: var(--text-secondary);
  font-size: 14px;
  cursor: pointer;
}

.checkbox-label input[type="checkbox"] {
  width: 18px;
  height: 18px;
  margin: 0;
}

.submit-btn {
  background: var(--primary);
  color: white;
  border: none;
  border-radius: 8px;
  padding: 16px 32px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  margin-top: 20px;
}

.submit-btn:hover:not(:disabled) {
  background: var(--primary-dark);
  transform: translateY(-2px);
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* 常见问题 */
.faq-section {
  padding: 80px 0;
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

.faq-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 30px;
}

.faq-item {
  background: linear-gradient(135deg, rgba(30, 41, 59, 0.9), rgba(15, 23, 42, 0.95));
  border: 1px solid var(--border);
  border-radius: 12px;
  padding: 30px 25px;
  transition: all 0.3s ease;
}

.faq-item:hover {
  border-color: var(--primary);
}

.faq-item h3 {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 15px;
  color: var(--text-primary);
  display: flex;
  align-items: center;
  gap: 12px;
}

.faq-item h3 i {
  color: var(--primary);
}

.faq-item p {
  color: var(--text-secondary);
  line-height: 1.6;
}

/* 社交媒体 */
.social-section {
  padding: 80px 0;
}

.social-links {
  display: flex;
  justify-content: center;
  gap: 30px;
  flex-wrap: wrap;
}

.social-link {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 25px 20px;
  background: linear-gradient(135deg, rgba(30, 41, 59, 0.9), rgba(15, 23, 42, 0.95));
  border: 1px solid var(--border);
  border-radius: 12px;
  color: var(--text-secondary);
  text-decoration: none;
  transition: all 0.3s ease;
  min-width: 120px;
}

.social-link:hover {
  border-color: var(--primary);
  color: var(--primary);
  transform: translateY(-5px);
}

.social-link i {
  font-size: 32px;
}

.social-link span {
  font-size: 14px;
  font-weight: 500;
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

  .hero-title {
    font-size: 36px;
  }

  .form-row {
    grid-template-columns: 1fr;
  }

  .form-container {
    padding: 30px 20px;
  }

  .contact-grid {
    grid-template-columns: 1fr;
  }

  .faq-grid {
    grid-template-columns: 1fr;
  }

  .social-links {
    gap: 15px;
  }
}
</style>
