<template>
  <div class="user-dashboard">
    <!-- 欢迎区域 -->
    <div class="welcome-section">
      <div class="welcome-content">
        <h1 class="welcome-title">欢迎回来，{{ userInfo.username || '用户' }}！</h1>
        <p class="welcome-subtitle">您的邮箱系统运行正常</p>
      </div>
      <div class="welcome-actions">
        <button class="action-btn primary" @click="composeEmail">
          <i class="fas fa-edit"></i>
          写邮件
        </button>
        <button class="action-btn secondary" @click="refreshData">
          <i class="fas fa-sync-alt"></i>
          刷新
        </button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon inbox">
          <i class="fas fa-inbox"></i>
        </div>
        <div class="stat-content">
          <h3 class="stat-number">{{ stats.unreadCount || 0 }}</h3>
          <p class="stat-label">未读邮件</p>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon sent">
          <i class="fas fa-paper-plane"></i>
        </div>
        <div class="stat-content">
          <h3 class="stat-number">{{ stats.sentCount || 0 }}</h3>
          <p class="stat-label">已发送</p>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon storage">
          <i class="fas fa-hdd"></i>
        </div>
        <div class="stat-content">
          <h3 class="stat-number">{{ stats.storageUsed || '0MB' }}</h3>
          <p class="stat-label">存储使用</p>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon mailboxes">
          <i class="fas fa-folder"></i>
        </div>
        <div class="stat-content">
          <h3 class="stat-number">{{ stats.mailboxCount || 0 }}</h3>
          <p class="stat-label">邮箱数量</p>
        </div>
      </div>
    </div>

    <!-- 快速操作 -->
    <div class="quick-actions">
      <h2 class="section-title">快速操作</h2>
      <div class="actions-grid">
        <router-link to="/user/inbox" class="quick-action-card">
          <div class="action-icon">
            <i class="fas fa-inbox"></i>
          </div>
          <h3>收件箱</h3>
          <p>查看收到的邮件</p>
        </router-link>

        <router-link to="/user/outbox" class="quick-action-card">
          <div class="action-icon">
            <i class="fas fa-paper-plane"></i>
          </div>
          <h3>发件箱</h3>
          <p>查看已发送的邮件</p>
        </router-link>

        <router-link to="/user/mailboxes" class="quick-action-card">
          <div class="action-icon">
            <i class="fas fa-folder"></i>
          </div>
          <h3>邮箱管理</h3>
          <p>管理您的邮箱账户</p>
        </router-link>

        <router-link to="/user/settings" class="quick-action-card">
          <div class="action-icon">
            <i class="fas fa-cog"></i>
          </div>
          <h3>账户设置</h3>
          <p>修改个人设置</p>
        </router-link>
      </div>
    </div>

    <!-- 最近邮件 -->
    <div class="recent-emails">
      <h2 class="section-title">最近邮件</h2>
      <div class="email-list" v-if="!loading">
        <!-- 有邮件时显示列表 -->
        <div
          class="email-item"
          v-for="email in recentEmails"
          :key="email.id"
          @click="viewEmail(email)"
        >
          <div class="email-avatar">
            <i class="fas fa-user"></i>
          </div>
          <div class="email-content">
            <h4 class="email-subject">{{ email.subject }}</h4>
            <p class="email-sender">来自: {{ email.sender }}</p>
            <p class="email-preview">{{ email.preview }}</p>
          </div>
          <div class="email-meta">
            <span class="email-time">{{ email.time }}</span>
            <span class="email-status" :class="email.status">{{ email.statusText }}</span>
          </div>
        </div>

        <!-- 无邮件时显示空状态 -->
        <div class="empty-state" v-if="recentEmails.length === 0">
          <div class="empty-icon">
            <i class="fas fa-inbox"></i>
          </div>
          <h3>暂无邮件</h3>
          <p>您还没有收到任何邮件</p>
          <router-link to="/user/compose" class="btn btn-primary">
            <i class="fas fa-edit"></i>
            写第一封邮件
          </router-link>
        </div>
      </div>

      <!-- 加载状态 -->
      <div class="loading-state" v-if="loading">
        <div class="loading-spinner">
          <i class="fas fa-spinner fa-spin"></i>
        </div>
        <p>正在加载邮件...</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import userApi from '@/services/userApi'

const router = useRouter()
const authStore = useAuthStore()

// 响应式数据
const stats = ref({
  unreadCount: 0,
  sentCount: 0,
  storageUsed: '0MB',
  mailboxCount: 0
})

const loading = ref(false)

const recentEmails = ref([])

// 计算属性
const userInfo = computed(() => authStore.userInfo)

// 加载仪表盘数据
const loadDashboardData = async () => {
  loading.value = true
  try {
    // 并行加载统计数据和最近邮件
    const [statsResponse, emailsResponse] = await Promise.allSettled([
      userApi.getDashboardStats(),
      userApi.getRecentEmails(5)
    ])

    // 处理统计数据
    if (statsResponse.status === 'fulfilled' && statsResponse.value.data.code === 0) {
      stats.value = statsResponse.value.data.data
    } else {
      // 如果API失败，使用默认值
      stats.value = {
        unreadCount: 0,
        sentCount: 0,
        storageUsed: '0MB',
        mailboxCount: 0
      }
    }

    // 处理最近邮件
    if (emailsResponse.status === 'fulfilled' && emailsResponse.value.data.code === 0) {
      const emails = emailsResponse.value.data.data || []
      recentEmails.value = emails.map(email => ({
        id: email.id,
        subject: email.subject || '无主题',
        sender: email.sender || email.from_email || '未知发件人',
        preview: email.preview || (email.content ? email.content.substring(0, 50) + '...' : '无内容预览'),
        time: formatTime(email.created_at || email.sent_at),
        status: email.is_read ? 'read' : 'unread',
        statusText: email.is_read ? '已读' : '未读'
      }))
    } else {
      recentEmails.value = []
    }
  } catch (error) {
    console.error('加载仪表盘数据失败:', error)
    // 设置默认值
    stats.value = {
      unreadCount: 0,
      sentCount: 0,
      storageUsed: '0MB',
      mailboxCount: 0
    }
    recentEmails.value = []
  } finally {
    loading.value = false
  }
}

// 格式化时间
const formatTime = (timestamp) => {
  if (!timestamp) return '未知时间'

  const now = new Date()
  const time = new Date(timestamp)
  const diff = now - time

  const minutes = Math.floor(diff / (1000 * 60))
  const hours = Math.floor(diff / (1000 * 60 * 60))
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  if (minutes < 60) {
    return `${minutes}分钟前`
  } else if (hours < 24) {
    return `${hours}小时前`
  } else {
    return `${days}天前`
  }
}

// 查看邮件详情
const viewEmail = (email) => {
  router.push(`/user/inbox?email=${email.id}`)
}

// 方法
const composeEmail = () => {
  router.push('/user/compose')
}

const refreshData = async () => {
  await loadDashboardData()
}

// 生命周期
onMounted(() => {
  loadDashboardData()
})
</script>

<style scoped>
.user-dashboard {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

/* 欢迎区域 */
.welcome-section {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  padding: 30px;
  margin-bottom: 30px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  transition: all 0.3s ease;
}

.welcome-section:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(255, 255, 255, 0.2);
}

.welcome-content h1 {
  font-size: 28px;
  font-weight: 700;
  margin-bottom: 8px;
  color: var(--admin-light);
  background: linear-gradient(to right, var(--admin-primary), #60a5fa);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.welcome-content p {
  color: var(--admin-gray);
  font-size: 16px;
}

.welcome-actions {
  display: flex;
  gap: 15px;
}

.action-btn {
  padding: 12px 24px;
  border-radius: 8px;
  border: none;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
}

.action-btn.primary {
  background: var(--admin-primary);
  color: white;
  box-shadow: 0 4px 12px rgba(37, 99, 235, 0.3);
}

.action-btn.primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(37, 99, 235, 0.4);
}

.action-btn.secondary {
  background: rgba(255, 255, 255, 0.1);
  color: var(--admin-light);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.action-btn.secondary:hover {
  background: rgba(255, 255, 255, 0.15);
  transform: translateY(-2px);
  border-color: rgba(255, 255, 255, 0.3);
}

/* 统计卡片 */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 40px;
}

.stat-card {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  padding: 25px;
  display: flex;
  align-items: center;
  gap: 20px;
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

.stat-icon.inbox {
  background: linear-gradient(135deg, var(--admin-primary), #60a5fa);
}

.stat-icon.sent {
  background: linear-gradient(135deg, var(--admin-success), #4ade80);
}

.stat-icon.storage {
  background: linear-gradient(135deg, var(--admin-warning), #fbbf24);
}

.stat-icon.mailboxes {
  background: linear-gradient(135deg, var(--admin-info), #60a5fa);
}

.stat-content h3 {
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 5px;
  color: var(--admin-light);
}

.stat-content p {
  color: var(--admin-gray);
  font-size: 14px;
}

/* 快速操作 */
.quick-actions {
  margin-bottom: 40px;
}

.section-title {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 20px;
  color: var(--admin-light);
}

.actions-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
}

.quick-action-card {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  padding: 25px;
  text-decoration: none;
  color: var(--admin-light);
  transition: all 0.3s ease;
  text-align: center;
}

.quick-action-card:hover {
  transform: translateY(-5px);
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(255, 255, 255, 0.2);
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
}

.action-icon {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--admin-primary), #60a5fa);
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 15px;
  font-size: 20px;
  color: white;
}

.quick-action-card h3 {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 8px;
  color: var(--admin-light);
}

.quick-action-card p {
  font-size: 14px;
  color: var(--admin-gray);
  line-height: 1.4;
}

/* 最近邮件 */
.recent-emails {
  margin-bottom: 40px;
}

.email-list {
  background: linear-gradient(135deg, rgba(30, 41, 59, 0.9), rgba(15, 23, 42, 0.95));
  border-radius: 12px;
  border: 1px solid var(--border);
  overflow: hidden;
}

.email-item {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 20px;
  border-bottom: 1px solid var(--border);
  transition: all 0.3s ease;
  cursor: pointer;
}

.email-item:last-child {
  border-bottom: none;
}

.email-item:hover {
  background: rgba(0, 180, 216, 0.05);
  transform: translateX(5px);
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
  font-size: 16px;
}

.email-content {
  flex: 1;
}

.email-subject {
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 5px;
  color: var(--text-primary);
}

.email-sender {
  font-size: 14px;
  color: var(--text-secondary);
  margin-bottom: 5px;
}

.email-preview {
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.4;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 300px;
}

.email-meta {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 5px;
}

.email-time {
  font-size: 12px;
  color: var(--text-secondary);
}

.email-status {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.email-status.unread {
  background: rgba(0, 180, 216, 0.2);
  color: var(--primary);
}

.email-status.read {
  background: rgba(148, 163, 184, 0.2);
  color: var(--text-secondary);
}

/* 空状态样式 */
.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: var(--text-secondary);
}

.empty-icon {
  font-size: 48px;
  color: var(--text-muted);
  margin-bottom: 20px;
}

.empty-state h3 {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 10px;
  color: var(--text-primary);
}

.empty-state p {
  font-size: 14px;
  margin-bottom: 30px;
  color: var(--text-secondary);
}

.empty-state .btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 12px 24px;
  background: var(--primary);
  color: white;
  text-decoration: none;
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.empty-state .btn:hover {
  background: var(--primary-dark);
  transform: translateY(-2px);
}

/* 加载状态样式 */
.loading-state {
  text-align: center;
  padding: 60px 20px;
  color: var(--text-secondary);
}

.loading-spinner {
  font-size: 32px;
  color: var(--primary);
  margin-bottom: 20px;
}

.loading-state p {
  font-size: 14px;
  color: var(--text-secondary);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .welcome-section {
    flex-direction: column;
    gap: 20px;
    text-align: center;
  }

  .stats-grid {
    grid-template-columns: 1fr;
  }

  .actions-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .email-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 15px;
  }

  .email-meta {
    align-items: flex-start;
    flex-direction: row;
    gap: 15px;
  }
}

@media (max-width: 480px) {
  .actions-grid {
    grid-template-columns: 1fr;
  }

  .welcome-actions {
    flex-direction: column;
    width: 100%;
  }

  .action-btn {
    justify-content: center;
  }
}
</style>
