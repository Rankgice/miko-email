<template>
  <div class="terms-page">
    <!-- 顶部导航 -->
    <nav class="top-nav">
      <div class="logo">
        <i class="fas fa-file-contract"></i>
        <span>服务条款</span>
      </div>
      <router-link to="/" class="back-home">
        <i class="fas fa-home"></i>
        返回首页
      </router-link>
    </nav>

    <!-- 主要内容 -->
    <main class="main-content">
      <div class="container">
        <div class="terms-content">
          <h1>服务条款</h1>
          <p class="last-updated">最后更新时间：2024年1月1日</p>

          <section class="terms-section">
            <h2>1. 服务条款的接受</h2>
            <p>
              欢迎使用Miko邮箱系统！通过访问或使用我们的服务，您同意受本服务条款的约束。
              如果您不同意这些条款，请不要使用我们的服务。
            </p>
          </section>

          <section class="terms-section">
            <h2>2. 服务描述</h2>
            <p>Miko邮箱系统提供以下服务：</p>
            <ul>
              <li>电子邮件收发服务</li>
              <li>邮箱管理功能</li>
              <li>域名管理服务</li>
              <li>邮件转发和过滤</li>
              <li>API接口服务</li>
              <li>技术支持服务</li>
            </ul>
          </section>

          <section class="terms-section">
            <h2>3. 用户账户</h2>
            <h3>3.1 账户注册</h3>
            <ul>
              <li>您必须提供准确、完整的注册信息</li>
              <li>您有责任维护账户信息的准确性</li>
              <li>您必须保护账户密码的安全</li>
              <li>一个人只能注册一个账户</li>
            </ul>

            <h3>3.2 账户责任</h3>
            <ul>
              <li>您对账户下的所有活动负责</li>
              <li>如发现账户被盗用，应立即通知我们</li>
              <li>不得与他人共享账户凭据</li>
            </ul>
          </section>

          <section class="terms-section">
            <h2>4. 使用规则</h2>
            <h3>4.1 禁止行为</h3>
            <p>使用我们的服务时，您不得：</p>
            <ul>
              <li>发送垃圾邮件或恶意邮件</li>
              <li>传播病毒、恶意软件或有害代码</li>
              <li>侵犯他人的知识产权</li>
              <li>发布非法、有害或不当内容</li>
              <li>尝试未经授权访问系统</li>
              <li>干扰或破坏服务的正常运行</li>
              <li>使用自动化工具进行大量操作</li>
            </ul>

            <h3>4.2 内容责任</h3>
            <ul>
              <li>您对通过服务发送的所有内容负责</li>
              <li>我们不对用户内容进行预先审查</li>
              <li>我们有权删除违规内容</li>
            </ul>
          </section>

          <section class="terms-section">
            <h2>5. 服务可用性</h2>
            <ul>
              <li>我们努力提供99.9%的服务可用性</li>
              <li>可能因维护、升级等原因暂停服务</li>
              <li>不可抗力因素可能影响服务可用性</li>
              <li>我们会提前通知计划内的服务中断</li>
            </ul>
          </section>

          <section class="terms-section">
            <h2>6. 费用和付款</h2>
            <h3>6.1 服务费用</h3>
            <ul>
              <li>基础服务免费提供</li>
              <li>高级功能可能需要付费</li>
              <li>费用标准在网站上公布</li>
              <li>价格可能会调整，会提前通知</li>
            </ul>

            <h3>6.2 付款条款</h3>
            <ul>
              <li>付费服务需要预付费</li>
              <li>支持多种付款方式</li>
              <li>逾期付款可能导致服务暂停</li>
              <li>退款政策详见相关条款</li>
            </ul>
          </section>

          <section class="terms-section">
            <h2>7. 知识产权</h2>
            <ul>
              <li>我们拥有服务相关的所有知识产权</li>
              <li>您保留对自己内容的知识产权</li>
              <li>使用服务不转让任何知识产权</li>
              <li>未经许可不得复制或分发我们的内容</li>
            </ul>
          </section>

          <section class="terms-section">
            <h2>8. 隐私保护</h2>
            <p>
              我们重视您的隐私，详细的隐私保护措施请参阅我们的
              <router-link to="/privacy-policy" class="link">隐私政策</router-link>。
            </p>
          </section>

          <section class="terms-section">
            <h2>9. 服务终止</h2>
            <h3>9.1 用户终止</h3>
            <ul>
              <li>您可以随时停止使用我们的服务</li>
              <li>删除账户前请备份重要数据</li>
              <li>账户删除后数据无法恢复</li>
            </ul>

            <h3>9.2 我们的终止权</h3>
            <p>在以下情况下，我们可能暂停或终止您的服务：</p>
            <ul>
              <li>违反服务条款</li>
              <li>长期不使用账户</li>
              <li>涉嫌非法活动</li>
              <li>拖欠服务费用</li>
            </ul>
          </section>

          <section class="terms-section">
            <h2>10. 免责声明</h2>
            <ul>
              <li>服务按"现状"提供，不提供任何明示或暗示的保证</li>
              <li>我们不保证服务完全无错误或不中断</li>
              <li>您使用服务的风险由您自行承担</li>
              <li>我们不对间接损失承担责任</li>
            </ul>
          </section>

          <section class="terms-section">
            <h2>11. 责任限制</h2>
            <p>
              在法律允许的最大范围内，我们对因使用或无法使用服务而产生的任何损害的责任
              不超过您在过去12个月内支付给我们的费用总额。
            </p>
          </section>

          <section class="terms-section">
            <h2>12. 争议解决</h2>
            <ul>
              <li>本条款受中华人民共和国法律管辖</li>
              <li>争议应首先通过友好协商解决</li>
              <li>协商不成的，提交北京仲裁委员会仲裁</li>
              <li>仲裁裁决是终局的，对双方都有约束力</li>
            </ul>
          </section>

          <section class="terms-section">
            <h2>13. 条款修改</h2>
            <p>
              我们可能会不时修改这些服务条款。重大修改会通过邮件或网站通知您。
              继续使用服务即表示您接受修改后的条款。
            </p>
          </section>

          <section class="terms-section">
            <h2>14. 其他条款</h2>
            <ul>
              <li>如果条款的任何部分被认定无效，其余部分仍然有效</li>
              <li>我们未行使任何权利不构成对该权利的放弃</li>
              <li>本条款构成您与我们之间的完整协议</li>
            </ul>
          </section>

          <section class="terms-section">
            <h2>15. 联系信息</h2>
            <p>如果您对这些服务条款有任何问题，请联系我们：</p>
            <div class="contact-info">
              <p><strong>邮箱：</strong> legal@miko-email.com</p>
              <p><strong>电话：</strong> 400-123-4567</p>
              <p><strong>地址：</strong> 北京市朝阳区科技园区创新大厦A座15层</p>
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
// 服务条款页面
</script>

<style scoped>
.terms-page {
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
  color: var(--text-primary);
  min-height: 100vh;
}

/* 复用隐私政策的样式 */
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

.terms-content {
  background: linear-gradient(135deg, rgba(30, 41, 59, 0.9), rgba(15, 23, 42, 0.95));
  border: 1px solid var(--border);
  border-radius: 16px;
  padding: 50px 40px;
}

.terms-content h1 {
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

.terms-section {
  margin-bottom: 40px;
}

.terms-section h2 {
  font-size: 24px;
  font-weight: 600;
  margin-bottom: 20px;
  color: var(--text-primary);
  border-bottom: 2px solid var(--border);
  padding-bottom: 10px;
}

.terms-section h3 {
  font-size: 18px;
  font-weight: 600;
  margin: 20px 0 10px 0;
  color: var(--primary);
}

.terms-section p {
  color: var(--text-secondary);
  line-height: 1.8;
  margin-bottom: 15px;
}

.terms-section ul {
  color: var(--text-secondary);
  line-height: 1.8;
  padding-left: 20px;
  margin-bottom: 15px;
}

.terms-section li {
  margin-bottom: 8px;
}

.terms-section strong {
  color: var(--text-primary);
}

.link {
  color: var(--primary);
  text-decoration: none;
  transition: all 0.3s ease;
}

.link:hover {
  color: var(--primary-light);
  text-decoration: underline;
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

  .terms-content {
    padding: 30px 20px;
  }

  .terms-content h1 {
    font-size: 28px;
  }

  .top-nav {
    padding: 15px 20px;
  }
}
</style>
