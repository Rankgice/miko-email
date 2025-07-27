<template>
  <div class="user-layout">
    <!-- 侧边栏 -->
    <aside class="sidebar" :class="{ collapsed: sidebarCollapsed }">
      <div class="sidebar-header">
        <div class="logo">
          <i class="fas fa-envelope"></i>
          <span v-show="!sidebarCollapsed">Miko邮箱</span>
        </div>
        <button class="collapse-btn" @click="toggleSidebar">
          <i class="fas fa-bars"></i>
        </button>
      </div>

      <nav class="sidebar-nav">
        <router-link to="/user/dashboard" class="nav-item">
          <i class="fas fa-tachometer-alt"></i>
          <span v-show="!sidebarCollapsed">仪表盘</span>
        </router-link>
        <router-link to="/user/inbox" class="nav-item">
          <i class="fas fa-inbox"></i>
          <span v-show="!sidebarCollapsed">收件箱</span>
        </router-link>
        <router-link to="/user/outbox" class="nav-item">
          <i class="fas fa-paper-plane"></i>
          <span v-show="!sidebarCollapsed">发件箱</span>
        </router-link>
        <router-link to="/user/compose" class="nav-item">
          <i class="fas fa-edit"></i>
          <span v-show="!sidebarCollapsed">写邮件</span>
        </router-link>
        <router-link to="/user/forward-rules" class="nav-item">
          <i class="fas fa-share"></i>
          <span v-show="!sidebarCollapsed">转发管理</span>
        </router-link>
        <router-link to="/user/mailboxes" class="nav-item">
          <i class="fas fa-folder"></i>
          <span v-show="!sidebarCollapsed">邮箱管理</span>
        </router-link>
        <router-link to="/user/domains" class="nav-item">
          <i class="fas fa-globe"></i>
          <span v-show="!sidebarCollapsed">域名管理</span>
        </router-link>
        <router-link to="/user/settings" class="nav-item">
          <i class="fas fa-cog"></i>
          <span v-show="!sidebarCollapsed">设置</span>
        </router-link>
      </nav>

      <div class="sidebar-footer">
        <div class="user-info">
          <div class="user-avatar">
            <i class="fas fa-user"></i>
          </div>
          <div class="user-details" v-show="!sidebarCollapsed">
            <span class="username">{{ userInfo.username || '用户' }}</span>
            <span class="user-email">{{ userInfo.email || 'user@example.com' }}</span>
          </div>
        </div>
      </div>
    </aside>

    <!-- 主内容区域 -->
    <div class="main-wrapper" :class="{ collapsed: sidebarCollapsed }">
      <!-- 顶部导航栏 -->
      <header class="header">
        <div class="header-left">
          <div class="breadcrumb">
            <span class="current-page">{{ getCurrentPageName() }}</span>
          </div>
        </div>

        <div class="header-right">
          <!-- 主题切换按钮 -->
          <button class="action-btn theme-toggle-btn" @click="toggleTheme" :title="currentTheme === 'dark' ? '切换到亮色主题' : '切换到暗色主题'">
            <i :class="currentTheme === 'dark' ? 'fas fa-sun' : 'fas fa-moon'"></i>
          </button>

          <!-- 通知按钮 -->
          <button class="action-btn" @click="toggleNotifications">
            <i class="fas fa-bell"></i>
            <span class="notification-badge" v-if="unreadNotifications > 0">{{ unreadNotifications }}</span>
          </button>

          <!-- 用户菜单 -->
          <div class="user-menu" @click="toggleUserMenu">
            <div class="user-avatar">
              <i class="fas fa-user"></i>
            </div>
            <i class="fas fa-chevron-down"></i>
          </div>

          <!-- 用户菜单下拉 -->
          <div class="user-dropdown" v-show="showUserMenu">
            <router-link to="/user/settings" class="dropdown-item">
              <i class="fas fa-cog"></i>
              <span>账户设置</span>
            </router-link>
            <div class="dropdown-divider"></div>
            <button class="dropdown-item logout-btn" @click="handleLogout">
              <i class="fas fa-sign-out-alt"></i>
              <span>退出登录</span>
            </button>
          </div>

          <!-- 通知面板 -->
          <div class="notification-panel" v-show="showNotifications">
            <div class="notification-header">
              <h3>通知</h3>
              <button @click="markAllAsRead" class="mark-read-btn">全部已读</button>
            </div>
            <div class="notification-list">
              <div v-if="notifications.length === 0" class="empty-notifications">
                <i class="fas fa-bell-slash"></i>
                <p>暂无通知</p>
              </div>
              <div v-else class="notification-item" v-for="notification in notifications" :key="notification.id" :class="{ unread: !notification.read }">
                <div class="notification-icon">
                  <i :class="notification.icon"></i>
                </div>
                <div class="notification-content">
                  <p class="notification-title">{{ notification.title }}</p>
                  <p class="notification-message">{{ notification.message }}</p>
                  <span class="notification-time">{{ notification.time }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </header>

      <!-- 主要内容区域 -->
      <main class="main-content">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import userApi from '@/services/userApi'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

// 响应式数据
const sidebarCollapsed = ref(false)
const showUserMenu = ref(false)
const showNotifications = ref(false)
const currentTheme = ref('dark')

// 通知数据 - 初始为空，从API加载
const notifications = ref([])

// 计算属性
const userInfo = computed(() => authStore.userInfo)
const unreadNotifications = computed(() =>
  notifications.value.filter(n => !n.read).length
)

// 获取当前页面名称
const getCurrentPageName = () => {
  const routeNames = {
    'UserDashboard': '仪表盘',
    'UserInbox': '收件箱',
    'UserOutbox': '发件箱',
    'UserCompose': '写邮件',
    'UserForwardRules': '转发管理',
    'UserMailboxes': '邮箱管理',
    'UserDomains': '域名管理',
    'UserSettings': '设置'
  }
  return routeNames[route.name] || '用户中心'
}

// 方法
const toggleSidebar = () => {
  sidebarCollapsed.value = !sidebarCollapsed.value
}

const toggleUserMenu = () => {
  showUserMenu.value = !showUserMenu.value
  showNotifications.value = false
}

const toggleNotifications = () => {
  showNotifications.value = !showNotifications.value
  showUserMenu.value = false
}

// 切换主题
const toggleTheme = () => {
  currentTheme.value = currentTheme.value === 'dark' ? 'light' : 'dark'
  document.documentElement.setAttribute('data-theme', currentTheme.value)
  localStorage.setItem('user-theme', currentTheme.value)
}

// 加载通知数据
const loadNotifications = async () => {
  try {
    const response = await userApi.getNotifications()
    if (response.data.code === 0) {
      notifications.value = response.data.data || []
    }
  } catch (error) {
    console.error('加载通知失败:', error)
    // API失败时显示空列表
    notifications.value = []
  }
}

const markAllAsRead = async () => {
  try {
    await userApi.markAllNotificationsAsRead()
    notifications.value.forEach(n => n.read = true)
  } catch (error) {
    console.error('标记通知已读失败:', error)
    // 即使API失败，也在前端标记为已读
    notifications.value.forEach(n => n.read = true)
  }
}

const handleLogout = async () => {
  if (confirm('确定要退出登录吗？')) {
    await authStore.userLogout()
    router.push('/user/login')
  }
}

// 点击外部关闭菜单
const handleClickOutside = (event) => {
  if (!event.target.closest('.user-menu') && !event.target.closest('.user-dropdown')) {
    showUserMenu.value = false
  }
  if (!event.target.closest('.action-btn') && !event.target.closest('.notification-panel')) {
    showNotifications.value = false
  }
}

// 初始化主题
const initTheme = () => {
  const savedTheme = localStorage.getItem('user-theme') || 'dark'
  currentTheme.value = savedTheme
  document.documentElement.setAttribute('data-theme', savedTheme)
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
  loadNotifications()
  initTheme()
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.user-layout {
  min-height: 100vh;
  background: var(--admin-dark);
  color: var(--admin-light);
  display: flex;
}

/* 侧边栏样式 */
.sidebar {
  width: 280px;
  background: linear-gradient(180deg, rgba(15, 23, 42, 0.95) 0%, rgba(30, 41, 59, 0.95) 100%);
  border-right: 1px solid rgba(255, 255, 255, 0.1);
  height: 100vh;
  position: fixed;
  left: 0;
  top: 0;
  z-index: 1000;
  transition: all 0.3s ease;
  display: flex;
  flex-direction: column;
}

.sidebar.collapsed {
  width: 80px;
}

.sidebar-header {
  padding: 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: space-between;
  min-height: 80px;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 20px;
  font-weight: 700;
  color: var(--admin-light);
}

.logo i {
  font-size: 24px;
  color: var(--admin-primary);
}

.collapse-btn {
  background: none;
  border: none;
  color: var(--admin-gray);
  cursor: pointer;
  padding: 8px;
  border-radius: 6px;
  transition: all 0.3s ease;
}

.collapse-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: var(--admin-light);
}

.sidebar-nav {
  flex: 1;
  padding: 20px 0;
  overflow-y: auto;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 20px;
  color: var(--admin-gray);
  text-decoration: none;
  font-weight: 500;
  transition: all 0.3s ease;
  border-left: 3px solid transparent;
  margin: 2px 0;
}

.nav-item:hover {
  background: rgba(255, 255, 255, 0.05);
  color: var(--admin-light);
  border-left-color: var(--admin-primary);
}

.nav-item.router-link-active {
  background: rgba(37, 99, 235, 0.1);
  color: var(--admin-primary);
  border-left-color: var(--admin-primary);
}

.nav-item i {
  width: 20px;
  text-align: center;
  font-size: 16px;
}

.sidebar-footer {
  padding: 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

/* 用户信息 */
.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-radius: 8px;
  transition: all 0.3s ease;
}

.user-info:hover {
  background: rgba(255, 255, 255, 0.05);
}

.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--admin-primary), #60a5fa);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 16px;
}

.user-details {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

.username {
  font-weight: 600;
  color: var(--admin-light);
  font-size: 14px;
}

.user-email {
  color: var(--admin-gray);
  font-size: 12px;
}

/* 主内容区域 */
.main-wrapper {
  flex: 1;
  margin-left: 280px;
  transition: all 0.3s ease;
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.main-wrapper.collapsed {
  margin-left: 80px;
}

.header {
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  padding: 0 30px;
  height: 70px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-left {
  display: flex;
  align-items: center;
}

.breadcrumb {
  display: flex;
  align-items: center;
  gap: 8px;
}

.current-page {
  font-size: 18px;
  font-weight: 600;
  color: var(--admin-light);
}

.header-right {
  display: flex;
  align-items: center;
  gap: 15px;
  position: relative;
}

/* 操作按钮 */
.action-btn {
  position: relative;
  background: none;
  border: none;
  color: var(--admin-gray);
  cursor: pointer;
  padding: 10px;
  border-radius: 8px;
  transition: all 0.3s ease;
  font-size: 16px;
}

.action-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: var(--admin-light);
}

/* 主题切换按钮特殊样式 */
.theme-toggle-btn {
  position: relative;
  overflow: hidden;
}

.theme-toggle-btn::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  width: 0;
  height: 0;
  background: radial-gradient(circle, var(--admin-primary) 0%, transparent 70%);
  transition: all 0.6s ease;
  transform: translate(-50%, -50%);
  border-radius: 50%;
  opacity: 0;
}

.theme-toggle-btn:hover::before {
  width: 100px;
  height: 100px;
  opacity: 0.1;
}

.notification-badge {
  position: absolute;
  top: 6px;
  right: 6px;
  background: var(--admin-danger);
  color: white;
  border-radius: 50%;
  width: 18px;
  height: 18px;
  font-size: 10px;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 用户菜单 */
.user-menu {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.user-menu:hover {
  background: rgba(255, 255, 255, 0.1);
}

.user-menu i {
  color: var(--admin-gray);
  transition: transform 0.3s ease;
}

.user-menu:hover i {
  transform: rotate(180deg);
}

.user-dropdown {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 8px;
  background: linear-gradient(135deg, rgba(15, 23, 42, 0.95), rgba(30, 41, 59, 0.95));
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
  min-width: 200px;
  overflow: hidden;
  z-index: 200;
}

.dropdown-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  color: var(--admin-light);
  text-decoration: none;
  transition: all 0.3s ease;
  border: none;
  background: none;
  width: 100%;
  text-align: left;
  cursor: pointer;
  font-size: 14px;
}

.dropdown-item:hover {
  background: rgba(37, 99, 235, 0.1);
  color: var(--admin-primary);
}

.dropdown-divider {
  height: 1px;
  background: rgba(255, 255, 255, 0.1);
  margin: 8px 0;
}

.logout-btn {
  color: var(--admin-danger) !important;
}

.logout-btn:hover {
  background: rgba(239, 68, 68, 0.1) !important;
}

/* 通知面板 */
.notification-panel {
  position: absolute;
  top: 100%;
  right: 60px;
  margin-top: 8px;
  background: linear-gradient(135deg, rgba(15, 23, 42, 0.95), rgba(30, 41, 59, 0.95));
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
  width: 350px;
  max-height: 400px;
  overflow: hidden;
  z-index: 200;
}

.notification-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.notification-header h3 {
  color: var(--admin-light);
  margin: 0;
  font-size: 16px;
  font-weight: 600;
}

.mark-read-btn {
  background: none;
  border: none;
  color: var(--admin-primary);
  cursor: pointer;
  font-size: 12px;
  padding: 4px 8px;
  border-radius: 4px;
  transition: all 0.3s ease;
}

.mark-read-btn:hover {
  background: rgba(37, 99, 235, 0.1);
}

.notification-list {
  max-height: 300px;
  overflow-y: auto;
}

.notification-item {
  display: flex;
  gap: 12px;
  padding: 15px 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  transition: all 0.3s ease;
}

.notification-item:last-child {
  border-bottom: none;
}

.notification-item:hover {
  background: rgba(255, 255, 255, 0.05);
}

.notification-item.unread {
  background: rgba(37, 99, 235, 0.05);
  border-left: 3px solid var(--admin-primary);
}

.notification-icon {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: rgba(37, 99, 235, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--admin-primary);
  font-size: 16px;
  flex-shrink: 0;
}

.notification-content {
  flex: 1;
}

.notification-title {
  color: var(--admin-light);
  font-weight: 600;
  font-size: 14px;
  margin: 0 0 4px 0;
}

.notification-message {
  color: var(--admin-gray);
  font-size: 13px;
  margin: 0 0 6px 0;
  line-height: 1.4;
}

.notification-time {
  color: var(--admin-gray);
  font-size: 11px;
}

/* 空通知状态 */
.empty-notifications {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  color: var(--admin-gray);
  text-align: center;
}

.empty-notifications i {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-notifications p {
  font-size: 14px;
  margin: 0;
}

/* 主要内容区域 */
.main-content {
  flex: 1;
  padding: 30px;
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 50%, #334155 100%);
  min-height: calc(100vh - 70px);
  overflow-y: auto;
}

/* 主题切换样式 */
[data-theme="light"] .user-layout {
  background: var(--bg-primary);
  color: var(--text-primary);
}

[data-theme="light"] .sidebar {
  background: linear-gradient(180deg, var(--bg-secondary) 0%, var(--bg-tertiary) 100%);
  border-right: 1px solid var(--border-color);
}

[data-theme="light"] .header {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid var(--border-color);
}

[data-theme="light"] .main-content {
  background: linear-gradient(135deg, var(--bg-primary) 0%, var(--bg-secondary) 100%);
}

[data-theme="light"] .nav-item {
  color: var(--text-secondary);
}

[data-theme="light"] .nav-item:hover {
  background: rgba(37, 99, 235, 0.1);
  color: var(--admin-primary);
}

[data-theme="light"] .nav-item.router-link-active {
  background: rgba(37, 99, 235, 0.15);
  color: var(--admin-primary);
}

[data-theme="light"] .action-btn {
  background: rgba(0, 0, 0, 0.05);
  border: 1px solid var(--border-color);
  color: var(--text-secondary);
}

[data-theme="light"] .action-btn:hover {
  background: rgba(37, 99, 235, 0.1);
  color: var(--admin-primary);
}

[data-theme="light"] .user-menu {
  background: rgba(0, 0, 0, 0.05);
  border: 1px solid var(--border-color);
}

[data-theme="light"] .user-dropdown {
  background: rgba(255, 255, 255, 0.95);
  border: 1px solid var(--border-color);
}

[data-theme="light"] .notification-panel {
  background: rgba(255, 255, 255, 0.95);
  border: 1px solid var(--border-color);
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .sidebar {
    width: 260px;
  }

  .main-wrapper {
    margin-left: 260px;
  }

  .main-wrapper.collapsed {
    margin-left: 80px;
  }

  .main-content {
    padding: 20px;
  }
}

@media (max-width: 768px) {
  .sidebar {
    transform: translateX(-100%);
    z-index: 1001;
  }

  .sidebar.collapsed {
    transform: translateX(-100%);
  }

  .main-wrapper,
  .main-wrapper.collapsed {
    margin-left: 0;
  }

  .header {
    padding: 0 15px;
  }

  .main-content {
    padding: 15px;
  }

  .notification-panel {
    width: 300px;
    right: 15px;
  }
}

@media (max-width: 480px) {
  .header {
    height: 60px;
  }

  .current-page {
    font-size: 16px;
  }

  .main-content {
    padding: 10px;
    min-height: calc(100vh - 60px);
  }

  .notification-panel {
    width: calc(100vw - 30px);
    right: 15px;
    left: 15px;
  }
}
</style>
