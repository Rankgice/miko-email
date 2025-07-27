<template>
  <div class="admin-layout">
    <!-- 侧边栏 -->
    <aside class="sidebar" :class="{ collapsed: sidebarCollapsed }">
      <div class="sidebar-header">
        <div class="logo">
          <i class="fas fa-shield-alt"></i>
          <span v-show="!sidebarCollapsed">管理后台</span>
        </div>
        <button class="collapse-btn" @click="toggleSidebar">
          <i class="fas fa-bars"></i>
        </button>
      </div>

      <nav class="sidebar-nav">
        <router-link to="/admin/dashboard" class="nav-item">
          <i class="fas fa-tachometer-alt"></i>
          <span v-show="!sidebarCollapsed">仪表盘</span>
        </router-link>
        <router-link to="/admin/users" class="nav-item">
          <i class="fas fa-users"></i>
          <span v-show="!sidebarCollapsed">用户管理</span>
        </router-link>
        <router-link to="/admin/mailboxes" class="nav-item">
          <i class="fas fa-envelope"></i>
          <span v-show="!sidebarCollapsed">邮箱管理</span>
        </router-link>
        <router-link to="/admin/domains" class="nav-item">
          <i class="fas fa-globe"></i>
          <span v-show="!sidebarCollapsed">域名管理</span>
        </router-link>
        <router-link to="/admin/captcha" class="nav-item">
          <i class="fas fa-shield-alt"></i>
          <span v-show="!sidebarCollapsed">验证码管理</span>
        </router-link>
      </nav>

      <div class="sidebar-footer">
        <div class="admin-info">
          <div class="admin-avatar">
            <i class="fas fa-user-shield"></i>
          </div>
          <div class="admin-details" v-show="!sidebarCollapsed">
            <span class="admin-name">{{ adminInfo.username || '管理员' }}</span>
            <span class="admin-role">系统管理员</span>
          </div>
        </div>
      </div>
    </aside>

    <!-- 主要内容区域 -->
    <div class="main-wrapper" :class="{ expanded: sidebarCollapsed }">
      <!-- 顶部导航栏 -->
      <header class="header">
        <div class="header-left">
          <h1 class="page-title">{{ pageTitle }}</h1>
        </div>

        <div class="header-right">
          <div class="header-actions">
            <button class="action-btn" @click="refreshData" title="刷新数据" :disabled="refreshing">
              <i :class="refreshing ? 'fas fa-spinner fa-spin' : 'fas fa-sync-alt'"></i>
            </button>
            <button class="action-btn" @click="showNotifications" title="消息通知" :class="{ active: showNotificationPanel }">
              <i class="fas fa-bell"></i>
              <span class="notification-badge" v-if="unreadNotifications > 0">{{ unreadNotifications }}</span>
            </button>
            <button class="action-btn" @click="toggleTheme" :title="currentTheme === 'dark' ? '切换到亮色主题' : '切换到暗色主题'">
              <i :class="currentTheme === 'dark' ? 'fas fa-sun' : 'fas fa-moon'"></i>
            </button>
          </div>

          <div class="admin-menu" @click="toggleAdminMenu">
            <div class="admin-avatar-small">
              <i class="fas fa-user-shield"></i>
            </div>
            <span class="admin-name-small">{{ adminInfo.username || '管理员' }}</span>
            <i class="fas fa-chevron-down"></i>
          </div>

          <!-- 管理员菜单下拉 -->
          <div class="admin-dropdown" v-show="showAdminMenu">
            <div class="dropdown-item" @click="showProfile">
              <i class="fas fa-user"></i>
              <span>个人资料</span>
            </div>
            <div class="dropdown-divider"></div>
            <button class="dropdown-item logout-btn" @click="handleLogout">
              <i class="fas fa-sign-out-alt"></i>
              <span>退出登录</span>
            </button>
          </div>
        </div>

        <!-- 通知面板 -->
        <div class="notification-panel" v-show="showNotificationPanel" @click.stop>
          <div class="notification-header">
            <h3>消息通知</h3>
            <button class="close-btn" @click="showNotificationPanel = false">
              <i class="fas fa-times"></i>
            </button>
          </div>
          <div class="notification-content">
            <div class="notification-item" v-for="notification in notifications" :key="notification.id" :class="{ unread: !notification.read }">
              <div class="notification-icon" :class="notification.type">
                <i :class="getNotificationIcon(notification.type)"></i>
              </div>
              <div class="notification-body">
                <div class="notification-title">{{ notification.title }}</div>
                <div class="notification-message">{{ notification.message }}</div>
                <div class="notification-time">{{ formatTime(notification.time) }}</div>
              </div>
              <button class="mark-read-btn" v-if="!notification.read" @click="markAsRead(notification.id)">
                <i class="fas fa-check"></i>
              </button>
            </div>
            <div class="no-notifications" v-if="notifications.length === 0">
              <i class="fas fa-bell-slash"></i>
              <p>暂无通知</p>
            </div>
          </div>
          <div class="notification-footer">
            <button class="btn btn-sm" @click="markAllAsRead">全部标记为已读</button>
            <button class="btn btn-sm btn-primary" @click="viewAllNotifications">查看全部</button>
          </div>
        </div>
      </header>

      <!-- 主要内容 -->
      <main class="main-content">
        <router-view />
      </main>
    </div>

    <!-- 个人资料对话框 -->
    <div class="modal-overlay" v-if="showProfileDialog" @click="showProfileDialog = false">
      <div class="modal-content profile-modal" @click.stop>
        <div class="modal-header">
          <h3>个人资料</h3>
          <button class="close-btn" @click="showProfileDialog = false">
            <i class="fas fa-times"></i>
          </button>
        </div>
        <div class="modal-body">
          <div class="profile-section">
            <div class="profile-avatar-section">
              <div class="profile-avatar-large">
                <i class="fas fa-user-shield"></i>
              </div>
              <button class="change-avatar-btn">更换头像</button>
            </div>
            <div class="profile-form">
              <div class="form-group">
                <label>用户名</label>
                <input type="text" v-model="profileForm.username" placeholder="请输入用户名" />
              </div>
              <div class="form-group">
                <label>邮箱地址</label>
                <input type="email" v-model="profileForm.email" placeholder="请输入邮箱地址" />
              </div>
              <div class="form-group">
                <label>显示名称</label>
                <input type="text" v-model="profileForm.displayName" placeholder="请输入显示名称" />
              </div>
              <div class="form-group">
                <label>角色</label>
                <input type="text" :value="profileForm.role" readonly />
              </div>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showProfileDialog = false">取消</button>
          <button class="btn btn-primary" @click="updateProfile">保存更改</button>
        </div>
      </div>
    </div>


  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import adminApi from '@/services/adminApi'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

// 响应式数据
const sidebarCollapsed = ref(false)
const showAdminMenu = ref(false)
const showNotificationPanel = ref(false)
const showProfileDialog = ref(false)
const refreshing = ref(false)
const currentTheme = ref('dark')

// 通知相关数据
const notifications = ref([])

const unreadNotifications = computed(() => {
  return notifications.value.filter(n => !n.read).length
})

// 个人资料表单
const profileForm = ref({
  username: '',
  email: '',
  displayName: '',
  role: '系统管理员'
})



// 计算属性
const adminInfo = computed(() => authStore.adminInfo)

const pageTitle = computed(() => {
  const titles = {
    '/admin/dashboard': '仪表盘',
    '/admin/users': '用户管理',
    '/admin/mailboxes': '邮箱管理',
    '/admin/domains': '域名管理',
    '/admin/captcha': '验证码管理'
  }
  return titles[route.path] || '管理后台'
})

// 方法
const toggleSidebar = () => {
  sidebarCollapsed.value = !sidebarCollapsed.value
}

const toggleAdminMenu = () => {
  showAdminMenu.value = !showAdminMenu.value
}

const handleLogout = async () => {
  if (confirm('确定要退出登录吗？')) {
    await authStore.adminLogout()
    router.push('/admin/login')
  }
}

// 刷新数据
const refreshData = async () => {
  if (refreshing.value) return

  refreshing.value = true
  try {
    // 触发当前页面的数据刷新
    window.location.reload()
  } catch (error) {
    console.error('刷新数据失败:', error)
  } finally {
    setTimeout(() => {
      refreshing.value = false
    }, 1000)
  }
}

// 显示通知面板
const showNotifications = () => {
  showNotificationPanel.value = !showNotificationPanel.value
  showAdminMenu.value = false
}

// 切换主题
const toggleTheme = () => {
  currentTheme.value = currentTheme.value === 'dark' ? 'light' : 'dark'
  document.documentElement.setAttribute('data-theme', currentTheme.value)
  localStorage.setItem('admin-theme', currentTheme.value)
}

// 显示个人资料
const showProfile = () => {
  profileForm.value.username = adminInfo.value.username || ''
  profileForm.value.email = adminInfo.value.email || ''
  profileForm.value.displayName = adminInfo.value.displayName || ''
  showProfileDialog.value = true
  showAdminMenu.value = false
}



// 更新个人资料
const updateProfile = async () => {
  try {
    // 这里应该调用API更新个人资料
    console.log('更新个人资料:', profileForm.value)
    alert('个人资料更新成功')
    showProfileDialog.value = false
  } catch (error) {
    console.error('更新个人资料失败:', error)
    alert('更新个人资料失败')
  }
}



// 通知相关方法
const getNotificationIcon = (type) => {
  const icons = {
    info: 'fas fa-info-circle',
    warning: 'fas fa-exclamation-triangle',
    success: 'fas fa-check-circle',
    error: 'fas fa-times-circle'
  }
  return icons[type] || 'fas fa-bell'
}

const formatTime = (time) => {
  const now = new Date()
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

// 加载通知数据
const loadNotifications = async () => {
  try {
    const response = await adminApi.getAdminNotifications()
    if (response.data.code === 0) {
      notifications.value = response.data.data.map(notification => ({
        ...notification,
        time: new Date(notification.time)
      }))
    }
  } catch (error) {
    console.error('加载通知失败:', error)
    // 使用默认通知
    notifications.value = [
      {
        id: 1,
        type: 'info',
        title: '系统状态',
        message: '系统运行正常',
        time: new Date(),
        read: false
      }
    ]
  }
}

// 保存通知状态到localStorage
const saveNotifications = () => {
  try {
    localStorage.setItem('admin_notifications', JSON.stringify(notifications.value))
  } catch (error) {
    console.error('保存通知数据失败:', error)
  }
}

const markAsRead = (id) => {
  const notification = notifications.value.find(n => n.id === id)
  if (notification) {
    notification.read = true
    saveNotifications()
  }
}

const markAllAsRead = () => {
  notifications.value.forEach(n => n.read = true)
  saveNotifications()
}

const viewAllNotifications = () => {
  // 跳转到通知页面或显示更多通知
  console.log('查看全部通知')
  showNotificationPanel.value = false
}

// 点击外部关闭菜单和面板
const handleClickOutside = (event) => {
  if (!event.target.closest('.admin-menu') && !event.target.closest('.admin-dropdown')) {
    showAdminMenu.value = false
  }
  if (!event.target.closest('.notification-panel') && !event.target.closest('.action-btn')) {
    showNotificationPanel.value = false
  }
}

// 初始化主题
const initTheme = () => {
  const savedTheme = localStorage.getItem('admin-theme') || 'dark'
  currentTheme.value = savedTheme
  document.documentElement.setAttribute('data-theme', savedTheme)
}

// 初始化个人资料
const initProfile = () => {
  if (adminInfo.value) {
    profileForm.value.username = adminInfo.value.username || ''
    profileForm.value.email = adminInfo.value.email || ''
    profileForm.value.displayName = adminInfo.value.displayName || adminInfo.value.username || ''
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
  initTheme()
  initProfile()
  loadNotifications()

  // 监听通知数据变化，自动保存到localStorage
  watch(notifications, () => {
    saveNotifications()
  }, { deep: true })
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.admin-layout {
  display: flex;
  min-height: 100vh;
  background: var(--admin-dark);
  color: var(--admin-light);
}

/* 侧边栏 */
.sidebar {
  width: 280px;
  background: linear-gradient(180deg, var(--admin-darker) 0%, var(--admin-dark) 100%);
  border-right: 1px solid rgba(255, 255, 255, 0.1);
  transition: width 0.3s ease;
  position: relative;
  z-index: 100;
}

.sidebar.collapsed {
  width: 80px;
}

.sidebar-header {
  padding: 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 20px;
  font-weight: 700;
  background: linear-gradient(135deg, var(--admin-primary), var(--admin-secondary));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.logo i {
  font-size: 24px;
  color: var(--admin-primary);
}

.collapse-btn {
  background: none;
  border: none;
  color: var(--admin-gray);
  font-size: 18px;
  cursor: pointer;
  padding: 8px;
  border-radius: 6px;
  transition: all 0.3s ease;
}

.collapse-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: var(--admin-primary);
}

.sidebar-nav {
  padding: 20px 0;
  flex: 1;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 15px 20px;
  color: var(--admin-gray);
  text-decoration: none;
  transition: all 0.3s ease;
  border-left: 3px solid transparent;
  margin: 2px 0;
}

.nav-item:hover {
  background: rgba(37, 99, 235, 0.1);
  color: var(--admin-primary);
  border-left-color: var(--admin-primary);
}

.nav-item.router-link-active {
  background: rgba(37, 99, 235, 0.15);
  color: var(--admin-primary);
  border-left-color: var(--admin-primary);
}

.nav-item i {
  font-size: 18px;
  width: 20px;
  text-align: center;
}

.sidebar-footer {
  padding: 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.admin-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.admin-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--admin-primary), var(--admin-secondary));
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 16px;
}

.admin-details {
  display: flex;
  flex-direction: column;
}

.admin-name {
  font-weight: 600;
  font-size: 14px;
  color: var(--admin-light);
}

.admin-role {
  font-size: 12px;
  color: var(--admin-gray);
}

/* 主要内容区域 */
.main-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  transition: margin-left 0.3s ease;
}

.main-wrapper.expanded {
  margin-left: 0;
}

/* 顶部导航栏 */
.header {
  background: rgba(15, 23, 42, 0.8);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  padding: 0 30px;
  height: 70px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: sticky;
  top: 0;
  z-index: 90;
}

.header-left {
  display: flex;
  align-items: center;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: var(--admin-light);
  margin: 0;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 20px;
  position: relative;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.action-btn {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: var(--admin-gray);
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.action-btn:hover {
  background: rgba(37, 99, 235, 0.1);
  color: var(--admin-primary);
}

.action-btn.active {
  background: rgba(37, 99, 235, 0.2);
  color: var(--admin-primary);
}

.action-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.notification-badge {
  position: absolute;
  top: -5px;
  right: -5px;
  background: var(--admin-danger);
  color: white;
  font-size: 10px;
  padding: 2px 6px;
  border-radius: 10px;
  min-width: 16px;
  text-align: center;
}

.admin-menu {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 16px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  cursor: pointer;
  transition: all 0.3s ease;
}

.admin-menu:hover {
  background: rgba(255, 255, 255, 0.08);
}

.admin-avatar-small {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--admin-primary), var(--admin-secondary));
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 14px;
}

.admin-name-small {
  font-weight: 500;
  color: var(--admin-light);
}

.admin-dropdown {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 8px;
  background: rgba(15, 23, 42, 0.95);
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

/* 通知面板样式 */
.notification-panel {
  position: absolute;
  top: 100%;
  right: 80px;
  margin-top: 8px;
  width: 350px;
  max-height: 500px;
  background: rgba(15, 23, 42, 0.95);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
  overflow: hidden;
  z-index: 200;
}

.notification-header {
  padding: 15px 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.notification-header h3 {
  margin: 0;
  color: var(--admin-light);
  font-size: 16px;
  font-weight: 600;
}

.notification-content {
  max-height: 350px;
  overflow-y: auto;
}

.notification-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 15px 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  transition: all 0.3s ease;
}

.notification-item:hover {
  background: rgba(255, 255, 255, 0.03);
}

.notification-item.unread {
  background: rgba(37, 99, 235, 0.05);
  border-left: 3px solid var(--admin-primary);
}

.notification-icon {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  flex-shrink: 0;
}

.notification-icon.info {
  background: rgba(59, 130, 246, 0.2);
  color: #3b82f6;
}

.notification-icon.warning {
  background: rgba(245, 158, 11, 0.2);
  color: #f59e0b;
}

.notification-icon.success {
  background: rgba(34, 197, 94, 0.2);
  color: #22c55e;
}

.notification-icon.error {
  background: rgba(239, 68, 68, 0.2);
  color: #ef4444;
}

.notification-body {
  flex: 1;
}

.notification-title {
  font-weight: 600;
  color: var(--admin-light);
  font-size: 14px;
  margin-bottom: 4px;
}

.notification-message {
  color: var(--admin-gray);
  font-size: 13px;
  line-height: 1.4;
  margin-bottom: 6px;
}

.notification-time {
  color: var(--admin-gray);
  font-size: 12px;
}

.mark-read-btn {
  background: none;
  border: none;
  color: var(--admin-gray);
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  transition: all 0.3s ease;
}

.mark-read-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: var(--admin-primary);
}

.no-notifications {
  text-align: center;
  padding: 40px 20px;
  color: var(--admin-gray);
}

.no-notifications i {
  font-size: 32px;
  margin-bottom: 10px;
  opacity: 0.5;
}

.notification-footer {
  padding: 15px 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  display: flex;
  justify-content: space-between;
  gap: 10px;
}

.btn {
  padding: 8px 16px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 12px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.btn-sm {
  padding: 6px 12px;
  font-size: 11px;
}

.btn-secondary {
  background: rgba(255, 255, 255, 0.1);
  color: var(--admin-light);
}

.btn-secondary:hover {
  background: rgba(255, 255, 255, 0.15);
}

.btn-primary {
  background: var(--admin-primary);
  color: white;
}

.btn-primary:hover {
  background: #3b82f6;
}

/* 主要内容 */
.main-content {
  flex: 1;
  padding: 30px;
  background: linear-gradient(135deg, var(--admin-dark) 0%, var(--admin-darker) 100%);
  min-height: calc(100vh - 70px);
}

/* 模态框样式 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  backdrop-filter: blur(5px);
}

.modal-content {
  background: linear-gradient(135deg, rgba(30, 41, 59, 0.95), rgba(15, 23, 42, 0.95));
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
  max-width: 90vw;
  max-height: 90vh;
  overflow: hidden;
  backdrop-filter: blur(20px);
}

.profile-modal {
  width: 600px;
}

.settings-modal {
  width: 800px;
}

.modal-header {
  padding: 20px 30px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h3 {
  margin: 0;
  color: var(--admin-light);
  font-size: 20px;
  font-weight: 600;
}

.close-btn {
  background: none;
  border: none;
  color: var(--admin-gray);
  font-size: 18px;
  cursor: pointer;
  padding: 8px;
  border-radius: 6px;
  transition: all 0.3s ease;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: var(--admin-light);
}

.modal-body {
  padding: 30px;
  max-height: 60vh;
  overflow-y: auto;
}

.modal-footer {
  padding: 20px 30px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

/* 个人资料样式 */
.profile-section {
  display: flex;
  gap: 30px;
}

.profile-avatar-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 15px;
}

.profile-avatar-large {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--admin-primary), var(--admin-secondary));
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 36px;
}

.change-avatar-btn {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  color: var(--admin-light);
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 12px;
  transition: all 0.3s ease;
}

.change-avatar-btn:hover {
  background: rgba(255, 255, 255, 0.15);
}

.profile-form {
  flex: 1;
}

/* 表单样式 */
.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  color: var(--admin-light);
  font-weight: 500;
  margin-bottom: 8px;
  font-size: 14px;
}

.form-group input,
.form-group select,
.form-group textarea {
  width: 100%;
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  color: var(--admin-light);
  font-size: 14px;
  transition: all 0.3s ease;
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  outline: none;
  border-color: var(--admin-primary);
  background: rgba(255, 255, 255, 0.08);
}

.form-group input[readonly] {
  background: rgba(255, 255, 255, 0.02);
  color: var(--admin-gray);
  cursor: not-allowed;
}

.form-group textarea {
  resize: vertical;
  min-height: 80px;
}

/* 设置标签页样式 */
.settings-tabs {
  display: flex;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  margin-bottom: 30px;
}

.tab-btn {
  background: none;
  border: none;
  color: var(--admin-gray);
  padding: 15px 20px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  border-bottom: 2px solid transparent;
  transition: all 0.3s ease;
  font-size: 14px;
}

.tab-btn:hover {
  color: var(--admin-light);
  background: rgba(255, 255, 255, 0.05);
}

.tab-btn.active {
  color: var(--admin-primary);
  border-bottom-color: var(--admin-primary);
}

.settings-section {
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

/* 开关样式 */
.switch-group {
  display: flex;
  align-items: center;
  gap: 12px;
}

.switch {
  position: relative;
  display: inline-block;
  width: 50px;
  height: 24px;
}

.switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

.slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(255, 255, 255, 0.2);
  transition: 0.3s;
  border-radius: 24px;
}

.slider:before {
  position: absolute;
  content: "";
  height: 18px;
  width: 18px;
  left: 3px;
  bottom: 3px;
  background-color: white;
  transition: 0.3s;
  border-radius: 50%;
}

input:checked + .slider {
  background-color: var(--admin-primary);
}

input:checked + .slider:before {
  transform: translateX(26px);
}

.switch-group span {
  color: var(--admin-gray);
  font-size: 13px;
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .sidebar {
    width: 80px;
  }

  .sidebar .logo span,
  .sidebar .nav-item span,
  .sidebar .admin-details {
    display: none;
  }

  .header {
    padding: 0 20px;
  }

  .main-content {
    padding: 20px;
  }
}

@media (max-width: 768px) {
  .sidebar {
    position: fixed;
    left: -280px;
    height: 100vh;
    z-index: 200;
  }

  .sidebar.collapsed {
    left: 0;
    width: 280px;
  }

  .main-wrapper {
    margin-left: 0;
  }

  .header {
    padding: 0 15px;
  }

  .page-title {
    font-size: 20px;
  }

  .header-actions {
    gap: 5px;
  }

  .admin-name-small {
    display: none;
  }
}

@media (max-width: 480px) {
  .main-content {
    padding: 15px;
  }

  .page-title {
    font-size: 18px;
  }

  .action-btn {
    width: 36px;
    height: 36px;
  }
}
</style>
