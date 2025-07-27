<template>
  <div class="email-service-page">
    <!-- 顶部导航 -->
    <nav class="top-nav">
      <div class="logo">
        <i class="fas fa-envelope"></i>
        <span>Miko邮箱系统</span>
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
          <h1 class="hero-title">
            <i class="fas fa-envelope"></i>
            邮箱服务
          </h1>
          <p class="hero-subtitle">
            专业的企业级邮箱解决方案，为您的业务提供安全、稳定、高效的邮件服务
          </p>
        </div>
      </section>

      <!-- 服务特性 -->
      <section class="features-section">
        <div class="container">
          <h2 class="section-title">服务特性</h2>
          <div class="features-grid">
            <div class="feature-card">
              <div class="feature-icon">
                <i class="fas fa-shield-alt"></i>
              </div>
              <h3>企业级安全</h3>
              <p>采用SSL/TLS加密传输，多重身份验证，确保邮件数据安全</p>
              <ul>
                <li>端到端加密</li>
                <li>反垃圾邮件</li>
                <li>病毒扫描</li>
                <li>数据备份</li>
              </ul>
            </div>

            <div class="feature-card">
              <div class="feature-icon">
                <i class="fas fa-server"></i>
              </div>
              <h3>高可用性</h3>
              <p>99.9%服务可用性保证，分布式架构，确保服务稳定运行</p>
              <ul>
                <li>负载均衡</li>
                <li>故障转移</li>
                <li>实时监控</li>
                <li>自动恢复</li>
              </ul>
            </div>

            <div class="feature-card">
              <div class="feature-icon">
                <i class="fas fa-mobile-alt"></i>
              </div>
              <h3>多端同步</h3>
              <p>支持Web、移动端、桌面客户端，随时随地访问您的邮件</p>
              <ul>
                <li>Web邮箱</li>
                <li>移动应用</li>
                <li>IMAP/POP3</li>
                <li>Exchange同步</li>
              </ul>
            </div>

            <div class="feature-card">
              <div class="feature-icon">
                <i class="fas fa-cogs"></i>
              </div>
              <h3>灵活管理</h3>
              <p>强大的管理控制台，轻松管理用户、域名、邮箱配置</p>
              <ul>
                <li>用户管理</li>
                <li>域名配置</li>
                <li>邮箱监控</li>
                <li>统计报表</li>
              </ul>
            </div>
          </div>
        </div>
      </section>

      <!-- 服务套餐 -->
      <section class="pricing-section">
        <div class="container">
          <h2 class="section-title">服务套餐</h2>
          <div class="pricing-grid">
            <div class="pricing-card">
              <div class="pricing-header">
                <h3>基础版</h3>
                <div class="price">
                  <span class="currency">¥</span>
                  <span class="amount">99</span>
                  <span class="period">/月</span>
                </div>
              </div>
              <div class="pricing-features">
                <ul>
                  <li><i class="fas fa-check"></i> 50个邮箱账户</li>
                  <li><i class="fas fa-check"></i> 5GB存储空间/账户</li>
                  <li><i class="fas fa-check"></i> 基础安全防护</li>
                  <li><i class="fas fa-check"></i> Web邮箱访问</li>
                  <li><i class="fas fa-check"></i> 邮件客户端支持</li>
                </ul>
              </div>
              <button class="pricing-btn">选择套餐</button>
            </div>

            <div class="pricing-card featured">
              <div class="pricing-badge">推荐</div>
              <div class="pricing-header">
                <h3>专业版</h3>
                <div class="price">
                  <span class="currency">¥</span>
                  <span class="amount">299</span>
                  <span class="period">/月</span>
                </div>
              </div>
              <div class="pricing-features">
                <ul>
                  <li><i class="fas fa-check"></i> 200个邮箱账户</li>
                  <li><i class="fas fa-check"></i> 20GB存储空间/账户</li>
                  <li><i class="fas fa-check"></i> 高级安全防护</li>
                  <li><i class="fas fa-check"></i> 移动应用支持</li>
                  <li><i class="fas fa-check"></i> 管理控制台</li>
                  <li><i class="fas fa-check"></i> 24/7技术支持</li>
                </ul>
              </div>
              <button class="pricing-btn">选择套餐</button>
            </div>

            <div class="pricing-card">
              <div class="pricing-header">
                <h3>企业版</h3>
                <div class="price">
                  <span class="currency">¥</span>
                  <span class="amount">599</span>
                  <span class="period">/月</span>
                </div>
              </div>
              <div class="pricing-features">
                <ul>
                  <li><i class="fas fa-check"></i> 无限邮箱账户</li>
                  <li><i class="fas fa-check"></i> 50GB存储空间/账户</li>
                  <li><i class="fas fa-check"></i> 企业级安全</li>
                  <li><i class="fas fa-check"></i> API接口支持</li>
                  <li><i class="fas fa-check"></i> 定制化服务</li>
                  <li><i class="fas fa-check"></i> 专属客户经理</li>
                </ul>
              </div>
              <button class="pricing-btn">联系销售</button>
            </div>
          </div>
        </div>
      </section>

      <!-- 技术规格 -->
      <section class="specs-section">
        <div class="container">
          <h2 class="section-title">技术规格</h2>
          <div class="specs-grid">
            <div class="spec-item">
              <h4>协议支持</h4>
              <p>SMTP, IMAP, POP3, Exchange ActiveSync</p>
            </div>
            <div class="spec-item">
              <h4>安全标准</h4>
              <p>SSL/TLS, SPF, DKIM, DMARC</p>
            </div>
            <div class="spec-item">
              <h4>存储容量</h4>
              <p>5GB - 50GB 可扩展</p>
            </div>
            <div class="spec-item">
              <h4>附件大小</h4>
              <p>最大支持 100MB 附件</p>
            </div>
            <div class="spec-item">
              <h4>服务可用性</h4>
              <p>99.9% SLA 保证</p>
            </div>
            <div class="spec-item">
              <h4>数据中心</h4>
              <p>国内多地部署，就近访问</p>
            </div>
          </div>
        </div>
      </section>

      <!-- 联系我们 -->
      <section class="contact-section">
        <div class="container">
          <div class="contact-content">
            <h2>开始使用邮箱服务</h2>
            <p>立即联系我们的销售团队，获取专业的邮箱解决方案</p>
            <div class="contact-actions">
              <router-link to="/contact" class="btn btn-primary">
                <i class="fas fa-phone"></i>
                联系销售
              </router-link>
              <router-link to="/user/register" class="btn btn-secondary">
                <i class="fas fa-user-plus"></i>
                免费试用
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
// 邮箱服务页面
</script>

<style scoped>
.email-service-page {
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
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 20px;
}

.hero-subtitle {
  font-size: 20px;
  color: var(--text-secondary);
  max-width: 600px;
  margin: 0 auto;
  line-height: 1.6;
}

/* 特性部分 */
.features-section {
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

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 40px;
}

.feature-card {
  background: linear-gradient(135deg, rgba(30, 41, 59, 0.9), rgba(15, 23, 42, 0.95));
  border-radius: 16px;
  border: 1px solid var(--border);
  padding: 40px 30px;
  text-align: center;
  transition: all 0.3s ease;
}

.feature-card:hover {
  transform: translateY(-5px);
  border-color: var(--primary);
  box-shadow: 0 20px 40px rgba(0, 180, 216, 0.1);
}

.feature-icon {
  width: 80px;
  height: 80px;
  background: linear-gradient(135deg, var(--primary), var(--primary-light));
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 20px;
  font-size: 32px;
  color: white;
}

.feature-card h3 {
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 15px;
  color: var(--text-primary);
}

.feature-card p {
  color: var(--text-secondary);
  margin-bottom: 20px;
  line-height: 1.6;
}

.feature-card ul {
  list-style: none;
  padding: 0;
  text-align: left;
}

.feature-card li {
  color: var(--text-secondary);
  margin-bottom: 8px;
  padding-left: 20px;
  position: relative;
}

.feature-card li::before {
  content: "✓";
  position: absolute;
  left: 0;
  color: var(--primary);
  font-weight: bold;
}

/* 价格部分 */
.pricing-section {
  padding: 80px 0;
}

.pricing-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 40px;
}

.pricing-card {
  background: linear-gradient(135deg, rgba(30, 41, 59, 0.9), rgba(15, 23, 42, 0.95));
  border-radius: 16px;
  border: 1px solid var(--border);
  padding: 40px 30px;
  text-align: center;
  position: relative;
  transition: all 0.3s ease;
}

.pricing-card.featured {
  border-color: var(--primary);
  transform: scale(1.05);
}

.pricing-badge {
  position: absolute;
  top: -15px;
  left: 50%;
  transform: translateX(-50%);
  background: var(--primary);
  color: white;
  padding: 8px 20px;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 600;
}

.pricing-header h3 {
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 20px;
  color: var(--text-primary);
}

.price {
  margin-bottom: 30px;
}

.currency {
  font-size: 20px;
  color: var(--text-secondary);
}

.amount {
  font-size: 48px;
  font-weight: 700;
  color: var(--primary);
}

.period {
  font-size: 16px;
  color: var(--text-secondary);
}

.pricing-features ul {
  list-style: none;
  padding: 0;
  margin-bottom: 30px;
}

.pricing-features li {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 12px;
  color: var(--text-secondary);
}

.pricing-features i {
  color: var(--primary);
}

.pricing-btn {
  width: 100%;
  padding: 15px;
  background: var(--primary);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.pricing-btn:hover {
  background: var(--primary-dark);
  transform: translateY(-2px);
}

/* 技术规格 */
.specs-section {
  padding: 80px 0;
  background: rgba(15, 23, 42, 0.5);
}

.specs-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 30px;
}

.spec-item {
  background: linear-gradient(135deg, rgba(30, 41, 59, 0.9), rgba(15, 23, 42, 0.95));
  border-radius: 12px;
  border: 1px solid var(--border);
  padding: 30px 20px;
  text-align: center;
}

.spec-item h4 {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 10px;
  color: var(--primary);
}

.spec-item p {
  color: var(--text-secondary);
  line-height: 1.6;
}

/* 联系部分 */
.contact-section {
  padding: 80px 0;
}

.contact-content {
  text-align: center;
  max-width: 600px;
  margin: 0 auto;
}

.contact-content h2 {
  font-size: 36px;
  font-weight: 700;
  margin-bottom: 20px;
  background: linear-gradient(to right, var(--primary), var(--primary-light));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.contact-content p {
  font-size: 18px;
  color: var(--text-secondary);
  margin-bottom: 40px;
  line-height: 1.6;
}

.contact-actions {
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

  .hero-title {
    font-size: 36px;
    flex-direction: column;
    gap: 10px;
  }

  .contact-actions {
    flex-direction: column;
    align-items: center;
  }

  .pricing-card.featured {
    transform: none;
  }
}
</style>
