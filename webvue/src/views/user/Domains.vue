<template>
  <div class="user-domains">
    <div class="page-header">
      <h1>可用域名</h1>
      <div class="header-actions">
        <button class="btn btn-secondary" @click="loadDomains" :disabled="loading">
          <i class="fas fa-sync-alt" :class="{ 'fa-spin': loading }"></i>
          刷新
        </button>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-state">
      <i class="fas fa-spinner fa-spin"></i>
      <span>加载中...</span>
    </div>

    <div v-else class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon primary">
          <i class="fas fa-globe"></i>
        </div>
        <div class="stat-content">
          <h3>{{ domainStats.total }}</h3>
          <p>可用域名数</p>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon success">
          <i class="fas fa-check-circle"></i>
        </div>
        <div class="stat-content">
          <h3>{{ domainStats.active }}</h3>
          <p>活跃域名</p>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon info">
          <i class="fas fa-envelope"></i>
        </div>
        <div class="stat-content">
          <h3>{{ userStats.mailboxes }}</h3>
          <p>我的邮箱数</p>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon warning">
          <i class="fas fa-share"></i>
        </div>
        <div class="stat-content">
          <h3>{{ userStats.forwardRules }}</h3>
          <p>转发规则数</p>
        </div>
      </div>
    </div>

    <!-- 域名列表 -->
    <div class="domains-grid" v-if="!loading">
      <div class="domain-card" v-for="domain in domains" :key="domain.name">
        <div class="domain-header">
          <div class="domain-info">
            <h3>{{ domain.name }}</h3>
            <p>{{ domain.description }}</p>
          </div>
          <div class="domain-status" :class="{ 'active': domain.is_active && domain.is_verified, 'inactive': !domain.is_active || !domain.is_verified }">
            <i :class="domain.is_active && domain.is_verified ? 'fas fa-check-circle' : 'fas fa-exclamation-triangle'"></i>
            <span>{{ getDomainStatusText(domain) }}</span>
          </div>
        </div>

        <div class="domain-stats">
          <div class="stat-item">
            <span class="stat-label">验证状态</span>
            <span class="stat-value" :class="{ 'active': domain.is_verified, 'inactive': !domain.is_verified }">
              {{ domain.is_verified ? '已验证' : '未验证' }}
            </span>
          </div>
          <div class="stat-item">
            <span class="stat-label">我的邮箱</span>
            <span class="stat-value">{{ getMailboxCountForDomain(domain.name) }}</span>
          </div>
          <div class="stat-item">
            <span class="stat-label">创建时间</span>
            <span class="stat-value">{{ formatDate(domain.created_at) }}</span>
          </div>
        </div>

        <div class="domain-actions">
          <button class="action-btn" @click="showDnsInfo(domain)" title="域名信息">
            <i class="fas fa-info-circle"></i>
          </button>
          <button class="action-btn primary" @click="createMailbox(domain.name)" title="创建邮箱" :disabled="!domain.is_active || !domain.is_verified">
            <i class="fas fa-plus"></i>
          </button>
        </div>
      </div>

      <!-- 空状态 -->
      <div class="empty-state" v-if="domains.length === 0">
        <i class="fas fa-globe"></i>
        <h3>暂无可用域名</h3>
        <p>请联系管理员添加域名</p>
      </div>
    </div>

    <!-- DNS信息弹窗 -->
    <div class="modal-overlay" v-if="showDnsModal" @click="closeDnsModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>{{ selectedDomain }} - DNS信息</h3>
          <button class="close-btn" @click="closeDnsModal">
            <i class="fas fa-times"></i>
          </button>
        </div>

        <div class="dns-info">
          <div class="info-section">
            <h4>域名信息</h4>
            <div class="info-grid">
              <div class="info-item">
                <span class="info-label">域名:</span>
                <span class="info-value">{{ selectedDomain }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">状态:</span>
                <span class="info-value active">已验证</span>
              </div>
              <div class="info-item">
                <span class="info-label">类型:</span>
                <span class="info-value">系统域名</span>
              </div>
            </div>
          </div>

          <div class="info-section">
            <h4>使用说明</h4>
            <div class="usage-info">
              <p>您可以使用此域名创建邮箱：</p>
              <ul>
                <li>前往 <router-link to="/user/mailboxes">邮箱管理</router-link> 页面</li>
                <li>点击"添加邮箱"按钮</li>
                <li>选择此域名创建新邮箱</li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import userApi from '@/services/userApi'

const router = useRouter()

// 响应式数据
const domains = ref([])
const loading = ref(false)
const showDnsModal = ref(false)
const selectedDomain = ref('')

const domainStats = ref({
  total: 0,
  active: 0
})

const userStats = ref({
  mailboxes: 0,
  forwardRules: 0
})

// 加载域名列表
const loadDomains = async () => {
  try {
    loading.value = true
    const response = await userApi.getDomains()

    // 统一使用 code: 0 格式
    if (response.data.code === 0) {
      domains.value = response.data.data || []
      updateStats()
    } else {
      const errorMsg = response.data.msg || response.data.message || '获取域名列表失败'
      console.error('获取域名列表失败:', errorMsg)
      alert('获取域名列表失败: ' + errorMsg)
    }
  } catch (error) {
    console.error('加载域名失败:', error)
    const errorMsg = error.response?.data?.msg || error.response?.data?.message || error.message || '网络错误'
    alert('加载域名失败: ' + errorMsg)
  } finally {
    loading.value = false
  }
}

// 更新统计数据
const updateStats = () => {
  domainStats.value = {
    total: domains.value.length,
    active: domains.value.length // 所有可用域名都是活跃的
  }
}

// 获取用户统计数据
const loadUserStats = async () => {
  try {
    // 获取用户邮箱数量
    const mailboxResponse = await userApi.getMailboxes()
    if (mailboxResponse.data.success) {
      userStats.value.mailboxes = mailboxResponse.data.data?.length || 0
    }

    // 获取转发规则数量
    const forwardResponse = await userApi.getForwardRules()
    if (forwardResponse.data.success) {
      userStats.value.forwardRules = forwardResponse.data.data?.length || 0
    }
  } catch (error) {
    console.error('获取用户统计失败:', error)
  }
}

// 获取指定域名下的邮箱数量
const getMailboxCountForDomain = (domainName) => {
  return mailboxes.value.filter(mailbox =>
    mailbox.email && mailbox.email.endsWith('@' + domainName)
  ).length
}

// 获取域名状态文本
const getDomainStatusText = (domain) => {
  if (!domain.is_active) {
    return '已禁用'
  } else if (!domain.is_verified) {
    return '未验证'
  } else {
    return '可用'
  }
}

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN')
}

// 显示DNS信息
const showDnsInfo = (domain) => {
  selectedDomain.value = domain
  showDnsModal.value = true
}

// 关闭DNS信息弹窗
const closeDnsModal = () => {
  showDnsModal.value = false
  selectedDomain.value = ''
}

// 创建邮箱
const createMailbox = (domain) => {
  router.push('/user/mailboxes')
}

// 生命周期
onMounted(() => {
  loadDomains()
  loadUserStats()
})
</script>

<style scoped>
.user-domains {
  padding: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
}

.page-header h1 {
  font-size: 32px;
  font-weight: 700;
  color: white;
  margin: 0;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.header-actions {
  display: flex;
  gap: 12px;
}

.btn {
  padding: 12px 24px;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  text-decoration: none;
}

.btn-secondary {
  background: rgba(255, 255, 255, 0.2);
  color: white;
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.btn-secondary:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.3);
  transform: translateY(-2px);
}

.btn-secondary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: white;
  font-size: 18px;
}

.loading-state i {
  font-size: 48px;
  margin-bottom: 16px;
  color: rgba(255, 255, 255, 0.8);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 24px;
  margin-bottom: 32px;
}

.stat-card {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  padding: 24px;
  display: flex;
  align-items: center;
  gap: 20px;
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-5px);
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(255, 255, 255, 0.2);
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  color: white;
}

.stat-icon.primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.stat-icon.success {
  background: linear-gradient(135deg, #48bb78 0%, #38a169 100%);
}

.stat-icon.info {
  background: linear-gradient(135deg, #4299e1 0%, #3182ce 100%);
}

.stat-icon.warning {
  background: linear-gradient(135deg, #ed8936 0%, #dd6b20 100%);
}

.stat-content h3 {
  font-size: 32px;
  font-weight: 700;
  color: #2d3748;
  margin: 0 0 8px 0;
  line-height: 1;
}

.stat-content p {
  font-size: 14px;
  color: #718096;
  margin: 0;
  font-weight: 500;
}

.domains-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 24px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 16px;
  padding: 32px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
}

.domain-card {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.domain-card:hover {
  transform: translateY(-5px);
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(255, 255, 255, 0.2);
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
}

.domain-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
}

.domain-info h3 {
  font-size: 18px;
  font-weight: 600;
  color: #2d3748;
  margin: 0 0 4px 0;
}

.domain-info p {
  font-size: 14px;
  color: #718096;
  margin: 0;
}

.domain-status {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  font-weight: 600;
  padding: 6px 12px;
  border-radius: 20px;
}

.domain-status.active {
  background: rgba(72, 187, 120, 0.1);
  color: #38a169;
}

.domain-status.inactive {
  background: rgba(245, 101, 101, 0.1);
  color: #e53e3e;
}

.domain-stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  margin-bottom: 20px;
  padding: 16px;
  background: #f7fafc;
  border-radius: 8px;
}

.stat-item {
  text-align: center;
}

.stat-label {
  display: block;
  font-size: 12px;
  color: #718096;
  margin-bottom: 4px;
}

.stat-value {
  display: block;
  font-size: 16px;
  font-weight: 600;
  color: #2d3748;
}

.stat-value.active {
  color: #38a169;
}

.stat-value.inactive {
  color: #e53e3e;
}

.domain-actions {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
}

.action-btn {
  width: 36px;
  height: 36px;
  border: none;
  border-radius: 8px;
  background: #f7fafc;
  color: #4a5568;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.action-btn:hover {
  background: #edf2f7;
  transform: translateY(-1px);
}

.action-btn.primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.action-btn.primary:hover {
  background: linear-gradient(135deg, #5a67d8 0%, #6b46c1 100%);
}

.empty-state {
  grid-column: 1 / -1;
  text-align: center;
  padding: 80px 20px;
  color: #718096;
}

.empty-state i {
  font-size: 64px;
  margin-bottom: 24px;
  color: #cbd5e0;
}

.empty-state h3 {
  font-size: 24px;
  margin: 0 0 12px 0;
  color: #4a5568;
}

.empty-state p {
  font-size: 16px;
  margin: 0;
}

/* 弹窗样式 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  backdrop-filter: blur(4px);
}

.modal-content {
  background: white;
  border-radius: 16px;
  width: 90%;
  max-width: 600px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.2);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px 24px 0 24px;
  border-bottom: 1px solid #e2e8f0;
  margin-bottom: 24px;
}

.modal-header h3 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: #2d3748;
}

.close-btn {
  width: 32px;
  height: 32px;
  border: none;
  background: #f7fafc;
  border-radius: 8px;
  color: #4a5568;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.close-btn:hover {
  background: #edf2f7;
}

.dns-info {
  padding: 0 24px 24px 24px;
}

.info-section {
  margin-bottom: 32px;
}

.info-section:last-child {
  margin-bottom: 0;
}

.info-section h4 {
  margin: 0 0 16px 0;
  font-size: 16px;
  font-weight: 600;
  color: #2d3748;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-label {
  font-size: 12px;
  color: #718096;
  font-weight: 600;
}

.info-value {
  font-size: 14px;
  color: #2d3748;
}

.info-value.active {
  color: #38a169;
  font-weight: 600;
}

.usage-info {
  background: #f7fafc;
  border-radius: 8px;
  padding: 20px;
}

.usage-info p {
  margin: 0 0 12px 0;
  color: #4a5568;
}

.usage-info ul {
  margin: 0;
  padding-left: 20px;
  color: #4a5568;
}

.usage-info li {
  margin-bottom: 8px;
}

.usage-info a {
  color: #667eea;
  text-decoration: none;
  font-weight: 600;
}

.usage-info a:hover {
  text-decoration: underline;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .user-domains {
    padding: 16px;
  }

  .page-header {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }

  .stats-grid {
    grid-template-columns: 1fr;
  }

  .domains-grid {
    grid-template-columns: 1fr;
    padding: 20px;
  }

  .modal-content {
    width: 95%;
    margin: 20px;
  }

  .info-grid {
    grid-template-columns: 1fr;
  }
}
</style>
