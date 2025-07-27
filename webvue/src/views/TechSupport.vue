<template>
  <div class="tech-support-page">
    <!-- 顶部导航 -->
    <nav class="top-nav">
      <div class="logo">
        <i class="fas fa-headset"></i>
        <span>技术支持</span>
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
          <h1 class="hero-title">技术支持</h1>
          <p class="hero-subtitle">
            我们的技术团队为您提供专业的技术支持服务，确保您的邮箱系统稳定运行
          </p>
        </div>
      </section>

      <!-- 支持服务 -->
      <section class="support-services-section">
        <div class="container">
          <h2 class="section-title">支持服务</h2>
          <div class="services-grid">
            <div class="service-card">
              <div class="service-icon">
                <i class="fas fa-phone-alt"></i>
              </div>
              <h3>电话支持</h3>
              <p>7×24小时技术支持热线</p>
              <div class="service-details">
                <strong>400-123-4567</strong>
                <span>紧急故障专线</span>
              </div>
            </div>

            <div class="service-card">
              <div class="service-icon">
                <i class="fas fa-envelope"></i>
              </div>
              <h3>邮件支持</h3>
              <p>详细问题描述，专业技术回复</p>
              <div class="service-details">
                <strong>tech@miko-email.com</strong>
                <span>2小时内响应</span>
              </div>
            </div>

            <div class="service-card">
              <div class="service-icon">
                <i class="fas fa-desktop"></i>
              </div>
              <h3>远程协助</h3>
              <p>远程桌面技术支持服务</p>
              <div class="service-details">
                <strong>预约远程支持</strong>
                <span>专业工程师一对一</span>
              </div>
            </div>

            <div class="service-card">
              <div class="service-icon">
                <i class="fas fa-users"></i>
              </div>
              <h3>现场服务</h3>
              <p>重要客户现场技术支持</p>
              <div class="service-details">
                <strong>企业级客户</strong>
                <span>24小时内到达</span>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- 技术文档 -->
      <section class="tech-docs-section">
        <div class="container">
          <h2 class="section-title">技术文档</h2>
          <div class="docs-grid">
            <div class="doc-category">
              <h3>
                <i class="fas fa-cog"></i>
                系统配置
              </h3>
              <ul>
                <li><a href="#">邮箱服务器配置指南</a></li>
                <li><a href="#">DNS设置教程</a></li>
                <li><a href="#">SSL证书配置</a></li>
                <li><a href="#">防火墙设置</a></li>
              </ul>
            </div>

            <div class="doc-category">
              <h3>
                <i class="fas fa-tools"></i>
                故障排除
              </h3>
              <ul>
                <li><a href="#">邮件发送失败解决方案</a></li>
                <li><a href="#">登录问题排查</a></li>
                <li><a href="#">性能优化指南</a></li>
                <li><a href="#">常见错误代码说明</a></li>
              </ul>
            </div>

            <div class="doc-category">
              <h3>
                <i class="fas fa-code"></i>
                API集成
              </h3>
              <ul>
                <li><a href="#">API接口文档</a></li>
                <li><a href="#">SDK下载</a></li>
                <li><a href="#">集成示例代码</a></li>
                <li><a href="#">Webhook配置</a></li>
              </ul>
            </div>

            <div class="doc-category">
              <h3>
                <i class="fas fa-shield-alt"></i>
                安全配置
              </h3>
              <ul>
                <li><a href="#">安全策略配置</a></li>
                <li><a href="#">反垃圾邮件设置</a></li>
                <li><a href="#">数据备份方案</a></li>
                <li><a href="#">安全审计日志</a></li>
              </ul>
            </div>
          </div>
        </div>
      </section>

      <!-- 系统状态 -->
      <section class="system-status-section">
        <div class="container">
          <h2 class="section-title">系统状态</h2>
          <div class="status-grid">
            <div class="status-item">
              <div class="status-indicator online"></div>
              <div class="status-info">
                <h4>邮件服务</h4>
                <span class="status-text">正常运行</span>
                <span class="uptime">99.9% 可用性</span>
              </div>
            </div>

            <div class="status-item">
              <div class="status-indicator online"></div>
              <div class="status-info">
                <h4>Web服务</h4>
                <span class="status-text">正常运行</span>
                <span class="uptime">99.8% 可用性</span>
              </div>
            </div>

            <div class="status-item">
              <div class="status-indicator online"></div>
              <div class="status-info">
                <h4>API服务</h4>
                <span class="status-text">正常运行</span>
                <span class="uptime">99.9% 可用性</span>
              </div>
            </div>

            <div class="status-item">
              <div class="status-indicator maintenance"></div>
              <div class="status-info">
                <h4>数据库</h4>
                <span class="status-text">维护中</span>
                <span class="uptime">预计2小时完成</span>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- 提交工单 -->
      <section class="ticket-section">
        <div class="container">
          <div class="ticket-container">
            <h2>提交技术支持工单</h2>
            <p>遇到技术问题？提交工单，我们的技术团队会尽快为您解决</p>
            
            <form @submit.prevent="submitTicket" class="ticket-form">
              <div class="form-row">
                <div class="form-group">
                  <label>问题类型</label>
                  <select v-model="ticket.type" required>
                    <option value="">请选择问题类型</option>
                    <option value="email">邮件收发问题</option>
                    <option value="login">登录认证问题</option>
                    <option value="performance">性能问题</option>
                    <option value="security">安全问题</option>
                    <option value="api">API集成问题</option>
                    <option value="other">其他技术问题</option>
                  </select>
                </div>
                <div class="form-group">
                  <label>优先级</label>
                  <select v-model="ticket.priority" required>
                    <option value="">请选择优先级</option>
                    <option value="low">低</option>
                    <option value="medium">中</option>
                    <option value="high">高</option>
                    <option value="urgent">紧急</option>
                  </select>
                </div>
              </div>

              <div class="form-group">
                <label>问题标题</label>
                <input type="text" v-model="ticket.title" required placeholder="简要描述您遇到的问题">
              </div>

              <div class="form-group">
                <label>详细描述</label>
                <textarea v-model="ticket.description" required rows="6" 
                  placeholder="请详细描述问题的具体情况、出现时间、错误信息等"></textarea>
              </div>

              <div class="form-group">
                <label>系统环境</label>
                <textarea v-model="ticket.environment" rows="3" 
                  placeholder="操作系统、浏览器版本、网络环境等信息"></textarea>
              </div>

              <button type="submit" class="submit-btn" :disabled="submitting">
                <i class="fas fa-ticket-alt"></i>
                {{ submitting ? '提交中...' : '提交工单' }}
              </button>
            </form>
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

const ticket = ref({
  type: '',
  priority: '',
  title: '',
  description: '',
  environment: ''
})

const submitTicket = async () => {
  submitting.value = true
  
  try {
    // 这里应该调用API提交工单
    await new Promise(resolve => setTimeout(resolve, 2000))
    
    alert('工单提交成功！工单号：#' + Math.random().toString(36).substr(2, 9).toUpperCase())
    
    // 重置表单
    ticket.value = {
      type: '',
      priority: '',
      title: '',
      description: '',
      environment: ''
    }
  } catch (error) {
    alert('提交失败，请稍后重试。')
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
/* 样式将在下一个编辑中添加 */
</style>
