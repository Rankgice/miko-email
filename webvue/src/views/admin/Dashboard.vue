<template>
  <div class="admin-dashboard">
    <h1>管理员仪表盘</h1>

    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon">
          <i class="fas fa-users"></i>
        </div>
        <div class="stat-content">
          <h3>{{ stats.totalUsers }}</h3>
          <p>总用户数</p>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">
          <i class="fas fa-envelope"></i>
        </div>
        <div class="stat-content">
          <h3>{{ stats.totalMailboxes }}</h3>
          <p>邮箱数量</p>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">
          <i class="fas fa-globe"></i>
        </div>
        <div class="stat-content">
          <h3>{{ stats.totalDomains }}</h3>
          <p>域名数量</p>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">
          <i class="fas fa-hdd"></i>
        </div>
        <div class="stat-content">
          <h3>{{ stats.storageUsed }}</h3>
          <p>存储使用</p>
        </div>
      </div>
    </div>

    <!-- 系统信息 -->
    <div class="system-info">
      <h2>系统信息</h2>
      <div class="info-grid">
        <div class="info-card">
          <div class="info-header">
            <i class="fas fa-info-circle"></i>
            <h3>系统版本</h3>
          </div>
          <div class="info-content">
            <span class="info-value">{{ systemInfo.version }}</span>
          </div>
        </div>

        <div class="info-card">
          <div class="info-header">
            <i class="fas fa-power-off"></i>
            <h3>运行状态</h3>
          </div>
          <div class="info-content">
            <span :class="['status-indicator', systemInfo.status]">
              {{ systemInfo.status === 'running' ? '运行中' : '已停止' }}
            </span>
          </div>
        </div>

        <div class="info-card">
          <div class="info-header">
            <i class="fas fa-envelope"></i>
            <h3>SMTP服务</h3>
          </div>
          <div class="info-content">
            <span :class="['status-indicator', serviceStatus.smtp]">
              {{ getServiceStatusText(serviceStatus.smtp) }}
            </span>
          </div>
        </div>

        <div class="info-card">
          <div class="info-header">
            <i class="fas fa-inbox"></i>
            <h3>IMAP服务</h3>
          </div>
          <div class="info-content">
            <span :class="['status-indicator', serviceStatus.imap]">
              {{ getServiceStatusText(serviceStatus.imap) }}
            </span>
          </div>
        </div>

        <div class="info-card">
          <div class="info-header">
            <i class="fas fa-download"></i>
            <h3>POP服务</h3>
          </div>
          <div class="info-content">
            <span :class="['status-indicator', serviceStatus.pop]">
              {{ getServiceStatusText(serviceStatus.pop) }}
            </span>
          </div>
        </div>

        <div class="info-card">
          <div class="info-header">
            <i class="fas fa-globe"></i>
            <h3>WEB服务</h3>
          </div>
          <div class="info-content">
            <span :class="['status-indicator', serviceStatus.web]">
              {{ getServiceStatusText(serviceStatus.web) }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <div class="quick-actions">
      <h2>快速操作</h2>
      <div class="actions-grid">
        <div class="action-card" @click="addUser">
          <i class="fas fa-user-plus"></i>
          <h3>添加用户</h3>
          <p>创建新的用户账户</p>
        </div>

        <div class="action-card" @click="addDomain">
          <i class="fas fa-plus-circle"></i>
          <h3>添加域名</h3>
          <p>添加新的邮件域名</p>
        </div>

        <div class="action-card" @click="viewLogs">
          <i class="fas fa-list-alt"></i>
          <h3>查看日志</h3>
          <p>系统操作日志</p>
        </div>

        <div class="action-card" @click="openSystemSettings">
          <i class="fas fa-cog"></i>
          <h3>系统设置</h3>
          <p>配置系统参数</p>
        </div>
      </div>
    </div>

    <!-- 系统日志对话框 -->
    <div class="modal-overlay" v-if="showLogsDialog" @click="showLogsDialog = false">
      <div class="modal-content logs-modal" @click.stop>
        <div class="modal-header">
          <h3>系统日志</h3>
          <button class="close-btn" @click="showLogsDialog = false">
            <i class="fas fa-times"></i>
          </button>
        </div>

        <div class="modal-body">
          <!-- 日志过滤器 -->
          <div class="log-filters">
            <div class="filter-group">
              <label>日志级别:</label>
              <select v-model="logFilters.level" @change="loadSystemLogs">
                <option value="all">全部</option>
                <option value="info">信息</option>
                <option value="warning">警告</option>
                <option value="error">错误</option>
                <option value="success">成功</option>
              </select>
            </div>
            <div class="filter-group">
              <label>模块:</label>
              <select v-model="logFilters.module" @change="loadSystemLogs">
                <option value="all">全部</option>
                <option value="system">系统</option>
                <option value="auth">认证</option>
                <option value="smtp">SMTP</option>
                <option value="database">数据库</option>
                <option value="storage">存储</option>
                <option value="mail">邮件</option>
                <option value="dns">DNS</option>
                <option value="user">用户</option>
              </select>
            </div>
            <div class="filter-group">
              <label>搜索:</label>
              <input type="text" v-model="logFilters.search" placeholder="搜索日志内容..." @input="debounceSearch" />
            </div>
            <button class="refresh-logs-btn" @click="loadSystemLogs" :disabled="logsLoading">
              <i :class="logsLoading ? 'fas fa-spinner fa-spin' : 'fas fa-sync-alt'"></i>
              刷新
            </button>
          </div>

          <!-- 日志列表 -->
          <div class="logs-container">
            <div class="log-item" v-for="log in filteredLogs" :key="log.id" :class="getLogLevelClass(log.level)">
              <div class="log-header">
                <span class="log-level">{{ log.level.toUpperCase() }}</span>
                <span class="log-module">[{{ log.module }}]</span>
                <span class="log-user" v-if="log.user">{{ log.user }}</span>
                <span class="log-ip" v-if="log.ip">{{ log.ip }}</span>
                <span class="log-time">{{ formatLogTime(log.timestamp) }}</span>
              </div>
              <div class="log-message">{{ log.message }}</div>
              <div class="log-details" v-if="log.details">{{ log.details }}</div>
            </div>

            <div class="no-logs" v-if="filteredLogs.length === 0 && !logsLoading">
              <i class="fas fa-file-alt"></i>
              <p>没有找到匹配的日志</p>
            </div>

            <div class="logs-loading" v-if="logsLoading">
              <i class="fas fa-spinner fa-spin"></i>
              <p>加载中...</p>
            </div>
          </div>
        </div>

        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showLogsDialog = false">关闭</button>
          <button class="btn btn-primary" @click="exportLogs">导出日志</button>
        </div>
      </div>
    </div>

    <!-- 系统设置对话框 -->
    <div class="modal-overlay" v-if="showSystemSettingsDialog" @click="showSystemSettingsDialog = false">
      <div class="modal-content settings-modal" @click.stop>
        <div class="modal-header">
          <h3>系统设置</h3>
          <button class="close-btn" @click="showSystemSettingsDialog = false">
            <i class="fas fa-times"></i>
          </button>
        </div>

        <div class="modal-body">
          <div class="settings-tabs">
            <button :class="['tab-btn', { active: activeSettingsTab === 'general' }]" @click="activeSettingsTab = 'general'">
              <i class="fas fa-cog"></i>
              常规设置
            </button>
            <button :class="['tab-btn', { active: activeSettingsTab === 'security' }]" @click="activeSettingsTab = 'security'">
              <i class="fas fa-shield-alt"></i>
              安全设置
            </button>
            <button :class="['tab-btn', { active: activeSettingsTab === 'email' }]" @click="activeSettingsTab = 'email'">
              <i class="fas fa-envelope"></i>
              邮件设置
            </button>
            <button :class="['tab-btn', { active: activeSettingsTab === 'storage' }]" @click="activeSettingsTab = 'storage'">
              <i class="fas fa-database"></i>
              存储设置
            </button>
          </div>

          <div class="settings-content">
            <!-- 常规设置 -->
            <div v-if="activeSettingsTab === 'general'" class="settings-section">
              <div class="form-group">
                <label>系统名称</label>
                <input type="text" v-model="systemSettings.general.systemName" placeholder="请输入系统名称" />
              </div>
              <div class="form-group">
                <label>系统描述</label>
                <textarea v-model="systemSettings.general.systemDescription" placeholder="请输入系统描述"></textarea>
              </div>
              <div class="form-group">
                <label>默认主题</label>
                <select v-model="systemSettings.general.defaultTheme">
                  <option value="dark">暗色主题</option>
                  <option value="light">亮色主题</option>
                </select>
              </div>
              <div class="form-group">
                <label>时区设置</label>
                <select v-model="systemSettings.general.timezone">
                  <option value="Asia/Shanghai">Asia/Shanghai (UTC+8)</option>
                  <option value="UTC">UTC (UTC+0)</option>
                </select>
              </div>
              <div class="form-group">
                <label>最大用户数</label>
                <input type="number" v-model="systemSettings.general.maxUsers" min="1" />
              </div>
              <div class="form-group">
                <label>允许用户注册</label>
                <div class="switch-group">
                  <label class="switch">
                    <input type="checkbox" v-model="systemSettings.general.enableRegistration" />
                    <span class="slider"></span>
                  </label>
                  <span>允许新用户自主注册</span>
                </div>
              </div>
            </div>

            <!-- 安全设置 -->
            <div v-if="activeSettingsTab === 'security'" class="settings-section">
              <div class="form-group">
                <label>会话超时时间 (分钟)</label>
                <input type="number" v-model="systemSettings.security.sessionTimeout" min="5" max="1440" />
              </div>
              <div class="form-group">
                <label>密码最小长度</label>
                <input type="number" v-model="systemSettings.security.minPasswordLength" min="6" max="32" />
              </div>
              <div class="form-group">
                <label>最大登录尝试次数</label>
                <input type="number" v-model="systemSettings.security.maxLoginAttempts" min="3" max="10" />
              </div>
              <div class="form-group">
                <label>锁定持续时间 (分钟)</label>
                <input type="number" v-model="systemSettings.security.lockoutDuration" min="5" max="60" />
              </div>
              <div class="form-group">
                <label>启用双因素认证</label>
                <div class="switch-group">
                  <label class="switch">
                    <input type="checkbox" v-model="systemSettings.security.enableTwoFactor" />
                    <span class="slider"></span>
                  </label>
                  <span>启用后管理员登录需要双因素验证</span>
                </div>
              </div>
              <div class="form-group">
                <label>登录失败锁定</label>
                <div class="switch-group">
                  <label class="switch">
                    <input type="checkbox" v-model="systemSettings.security.enableLoginLock" />
                    <span class="slider"></span>
                  </label>
                  <span>连续登录失败后锁定账户</span>
                </div>
              </div>
              <div class="form-group">
                <label>启用SSL</label>
                <div class="switch-group">
                  <label class="switch">
                    <input type="checkbox" v-model="systemSettings.security.enableSSL" />
                    <span class="slider"></span>
                  </label>
                  <span>启用SSL加密连接</span>
                </div>
              </div>
              <div class="form-group">
                <label>启用防火墙</label>
                <div class="switch-group">
                  <label class="switch">
                    <input type="checkbox" v-model="systemSettings.security.enableFirewall" />
                    <span class="slider"></span>
                  </label>
                  <span>启用内置防火墙保护</span>
                </div>
              </div>
            </div>

            <!-- 邮件设置 -->
            <div v-if="activeSettingsTab === 'email'" class="settings-section">
              <div class="form-group">
                <label>SMTP服务器</label>
                <input type="text" v-model="systemSettings.email.smtpHost" placeholder="请输入SMTP服务器地址" />
              </div>
              <div class="form-group">
                <label>SMTP端口</label>
                <input type="number" v-model="systemSettings.email.smtpPort" placeholder="请输入SMTP端口" />
              </div>
              <div class="form-group">
                <label>IMAP端口</label>
                <input type="number" v-model="systemSettings.email.imapPort" placeholder="请输入IMAP端口" />
              </div>
              <div class="form-group">
                <label>POP3端口</label>
                <input type="number" v-model="systemSettings.email.pop3Port" placeholder="请输入POP3端口" />
              </div>
              <div class="form-group">
                <label>发件人邮箱</label>
                <input type="email" v-model="systemSettings.email.senderEmail" placeholder="请输入发件人邮箱" />
              </div>
              <div class="form-group">
                <label>最大附件大小 (MB)</label>
                <input type="number" v-model="systemSettings.email.maxAttachmentSize" min="1" max="100" />
              </div>
              <div class="form-group">
                <label>邮件签名</label>
                <textarea v-model="systemSettings.email.emailSignature" placeholder="请输入邮件签名"></textarea>
              </div>
              <div class="form-group">
                <label>启用邮件加密</label>
                <div class="switch-group">
                  <label class="switch">
                    <input type="checkbox" v-model="systemSettings.email.enableEncryption" />
                    <span class="slider"></span>
                  </label>
                  <span>启用邮件传输加密</span>
                </div>
              </div>
            </div>

            <!-- 存储设置 -->
            <div v-if="activeSettingsTab === 'storage'" class="settings-section">
              <div class="form-group">
                <label>最大存储空间 (GB)</label>
                <input type="number" v-model="systemSettings.storage.maxStorageSpace" min="1" />
              </div>
              <div class="form-group">
                <label>清理天数</label>
                <input type="number" v-model="systemSettings.storage.cleanupDays" min="1" max="365" />
              </div>
              <div class="form-group">
                <label>备份间隔 (小时)</label>
                <input type="number" v-model="systemSettings.storage.backupInterval" min="1" max="168" />
              </div>
              <div class="form-group">
                <label>压缩级别 (1-9)</label>
                <input type="number" v-model="systemSettings.storage.compressionLevel" min="1" max="9" />
              </div>
              <div class="form-group">
                <label>自动清理</label>
                <div class="switch-group">
                  <label class="switch">
                    <input type="checkbox" v-model="systemSettings.storage.enableAutoCleanup" />
                    <span class="slider"></span>
                  </label>
                  <span>自动清理过期文件</span>
                </div>
              </div>
              <div class="form-group">
                <label>启用备份</label>
                <div class="switch-group">
                  <label class="switch">
                    <input type="checkbox" v-model="systemSettings.storage.backupEnabled" />
                    <span class="slider"></span>
                  </label>
                  <span>启用自动备份功能</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showSystemSettingsDialog = false">取消</button>
          <button class="btn btn-primary" @click="saveSettings">保存设置</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import adminApi from '@/services/adminApi'

const router = useRouter()

// 响应式数据
const stats = ref({
  totalUsers: 0,
  totalDomains: 0,
  totalMailboxes: 0,
  storageUsed: '0GB'
})

const systemInfo = ref({
  version: 'v1.0.0',
  uptime: '0天',
  status: 'running'
})

const serviceStatus = ref({
  smtp: 'running',
  imap: 'running',
  pop: 'running',
  web: 'running'
})

const systemHealth = ref({
  emailService: 'running',
  database: 'normal',
  systemResources: 'warning'
})

const recentActivities = ref([])
const loading = ref(false)

// 日志和设置对话框
const showLogsDialog = ref(false)
const showSystemSettingsDialog = ref(false)
const logsLoading = ref(false)
const activeLogTab = ref('system')

// 系统日志数据
const systemLogs = ref([])
const logFilters = ref({
  level: 'all',
  module: 'all',
  search: ''
})

// 系统设置数据
const systemSettings = ref({
  general: {},
  security: {},
  email: {},
  storage: {}
})
const settingsLoading = ref(false)
const activeSettingsTab = ref('general')

// 系统状态数据
const systemStatus = ref({
  services: [],
  system: {},
  database: {}
})

// 方法
const addUser = () => {
  router.push('/admin/users')
}

const addDomain = () => {
  router.push('/admin/domains')
}

const viewLogs = () => {
  showLogsDialog.value = true
  loadSystemLogs()
}

const openSystemSettings = () => {
  showSystemSettingsDialog.value = true
  loadSystemSettings()
}

// 搜索防抖
let searchTimeout = null
const debounceSearch = () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    loadSystemLogs()
  }, 500)
}

const getServiceStatusText = (status) => {
  const statusMap = {
    'running': '运行中',
    'stopped': '已停止',
    'error': '错误',
    'warning': '警告'
  }
  return statusMap[status] || '未知'
}

// 加载系统日志
const loadSystemLogs = async () => {
  logsLoading.value = true
  try {
    const params = {
      level: logFilters.value.level !== 'all' ? logFilters.value.level : '',
      module: logFilters.value.module !== 'all' ? logFilters.value.module : '',
      search: logFilters.value.search
    }

    const response = await adminApi.getSystemLogs(params)
    if (response.data.code === 0) {
      const data = response.data.data
      // 后端直接返回日志数组
      systemLogs.value = (Array.isArray(data) ? data : []).map(log => ({
        ...log,
        timestamp: new Date(log.timestamp)
      }))
    } else {
      console.error('获取系统日志失败:', response.data.msg)
      alert('获取系统日志失败: ' + response.data.msg)
    }
  } catch (error) {
    console.error('加载系统日志失败:', error)
    alert('加载系统日志失败，请检查网络连接')
  } finally {
    logsLoading.value = false
  }
}

// 格式化时间
const formatLogTime = (timestamp) => {
  return timestamp.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

// 获取日志级别样式
const getLogLevelClass = (level) => {
  const classMap = {
    'info': 'log-info',
    'warning': 'log-warning',
    'error': 'log-error',
    'success': 'log-success'
  }
  return classMap[level] || 'log-info'
}

// 过滤日志
const filteredLogs = computed(() => {
  let logs = systemLogs.value

  // 按级别过滤
  if (logFilters.value.level !== 'all') {
    logs = logs.filter(log => log.level === logFilters.value.level)
  }

  // 按搜索关键词过滤
  if (logFilters.value.search) {
    const search = logFilters.value.search.toLowerCase()
    logs = logs.filter(log =>
      log.message.toLowerCase().includes(search) ||
      log.module.toLowerCase().includes(search) ||
      log.details.toLowerCase().includes(search)
    )
  }

  return logs
})

// 导出日志
const exportLogs = () => {
  const logs = filteredLogs.value
  const csvContent = [
    ['时间', '级别', '模块', '消息', '详情'].join(','),
    ...logs.map(log => [
      formatLogTime(log.timestamp),
      log.level,
      log.module,
      `"${log.message}"`,
      `"${log.details}"`
    ].join(','))
  ].join('\n')

  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
  const link = document.createElement('a')
  const url = URL.createObjectURL(blob)
  link.setAttribute('href', url)
  link.setAttribute('download', `system_logs_${new Date().toISOString().split('T')[0]}.csv`)
  link.style.visibility = 'hidden'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

// 加载系统设置
const loadSystemSettings = async () => {
  settingsLoading.value = true
  try {
    const response = await adminApi.getSystemSettings()
    if (response.data.code === 0) {
      systemSettings.value = response.data.data
    } else {
      console.error('获取系统设置失败:', response.data.msg)
      alert('获取系统设置失败: ' + response.data.msg)
    }
  } catch (error) {
    console.error('加载系统设置失败:', error)
    alert('加载系统设置失败，请检查网络连接')
  } finally {
    settingsLoading.value = false
  }
}

// 保存设置
const saveSettings = async () => {
  settingsLoading.value = true
  try {
    const response = await adminApi.updateSystemSettings(systemSettings.value)
    if (response.data.code === 0) {
      alert('设置保存成功')
      showSystemSettingsDialog.value = false
    } else {
      alert('保存设置失败: ' + response.data.msg)
    }
  } catch (error) {
    console.error('保存设置失败:', error)
    alert('保存设置失败，请检查网络连接')
  } finally {
    settingsLoading.value = false
  }
}

// 加载系统状态
const loadSystemStatus = async () => {
  try {
    const response = await adminApi.getSystemStatus()
    if (response.data.code === 0) {
      systemStatus.value = response.data.data
    } else {
      console.error('获取系统状态失败:', response.data.msg)
    }
  } catch (error) {
    console.error('加载系统状态失败:', error)
  }
}

const loadDashboardData = async () => {
  loading.value = true
  try {
    // 并行加载所有数据
    const [statsRes, healthRes, activitiesRes, usersRes, domainsRes, mailboxesRes] = await Promise.all([
      adminApi.getDashboardStats().catch(() => ({ data: { data: {} } })),
      adminApi.getSystemHealth().catch(() => ({ data: { data: {} } })),
      adminApi.getRecentActivities().catch(() => ({ data: { data: [] } })),
      adminApi.getAllUsers().catch(() => ({ data: { data: [] } })),
      adminApi.getAllDomains().catch(() => ({ data: { data: [] } })),
      adminApi.getAllMailboxes().catch(() => ({ data: { data: [] } }))
    ])

    // 计算真实的统计数据
    const realStats = {
      totalUsers: 0,
      totalDomains: 0,
      totalMailboxes: 0,
      storageUsed: '0GB'
    }

    // 统计用户数量
    if (usersRes.data.code === 0) {
      realStats.totalUsers = usersRes.data.data.length
    }

    // 统计域名数量
    if (domainsRes.data.code === 0) {
      realStats.totalDomains = domainsRes.data.data.length
    }

    // 统计邮箱数量
    if (mailboxesRes.data.code === 0) {
      realStats.totalMailboxes = mailboxesRes.data.data.length
    }

    // 从后端统计数据获取存储使用量
    if (statsRes.data.code === 0) {
      const data = statsRes.data.data
      realStats.storageUsed = data.storageUsed ? `${data.storageUsed}GB` : '0GB'
    }

    stats.value = realStats

    // 更新系统信息
    if (healthRes.data.code === 0) {
      const healthData = healthRes.data.data

      // 更新系统信息
      if (healthData.systemInfo) {
        systemInfo.value = {
          version: healthData.systemInfo.version || 'v1.0.0',
          uptime: healthData.systemInfo.uptime || '0天',
          status: healthData.status || 'running'
        }
      }

      // 更新服务状态 (从后端获取真实数据)
      if (healthData.services) {
        serviceStatus.value = {
          smtp: healthData.services.smtp || 'unknown',
          imap: healthData.services.imap || 'unknown',
          pop: healthData.services.pop || 'unknown',
          web: healthData.services.web || 'unknown'
        }
      } else {
        // 如果后端没有返回服务状态，使用默认值
        serviceStatus.value = {
          smtp: 'unknown',
          imap: 'unknown',
          pop: 'unknown',
          web: 'unknown'
        }
      }

      systemHealth.value = healthData
    }

    // 更新最近活动
    if (activitiesRes.data.code === 0) {
      recentActivities.value = activitiesRes.data.data
    }

    console.log('仪表盘数据加载完成:', {
      stats: stats.value,
      systemInfo: systemInfo.value,
      serviceStatus: serviceStatus.value
    })

  } catch (error) {
    console.error('加载仪表盘数据失败:', error)
  } finally {
    loading.value = false
  }
}

// 生命周期
onMounted(() => {
  loadDashboardData()
})
</script>

<style scoped>
.admin-dashboard {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

h1 {
  color: var(--admin-light);
  margin-bottom: 30px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 40px;
}

.stat-card {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 15px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.stat-icon {
  width: 50px;
  height: 50px;
  border-radius: 10px;
  background: linear-gradient(135deg, #3b82f6, #2563eb);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 20px;
}

.stat-content h3 {
  font-size: 24px;
  font-weight: 700;
  color: var(--admin-light);
  margin-bottom: 4px;
}

.stat-content p {
  color: var(--admin-gray);
  font-size: 14px;
  margin: 0;
}

.quick-actions h2 {
  color: var(--admin-light);
  margin-bottom: 20px;
}

.actions-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
}

.action-card {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  padding: 20px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.action-card:hover {
  transform: translateY(-3px);
  background: rgba(59, 130, 246, 0.1);
}

.action-card i {
  font-size: 30px;
  color: #3b82f6;
  margin-bottom: 15px;
}

.action-card h3 {
  color: var(--admin-light);
  margin-bottom: 8px;
}

.action-card p {
  color: var(--admin-gray);
  font-size: 14px;
  margin: 0;
}

/* 系统信息样式 */
.system-info {
  margin-bottom: 30px;
}

.system-info h2 {
  color: var(--admin-light);
  margin-bottom: 20px;
  font-size: 24px;
  font-weight: 600;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 20px;
}

.info-card {
  background: linear-gradient(135deg, rgba(30, 41, 59, 0.8), rgba(15, 23, 42, 0.9));
  border-radius: 12px;
  padding: 20px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  transition: all 0.3s ease;
}

.info-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.3);
  border-color: rgba(255, 255, 255, 0.2);
}

.info-header {
  display: flex;
  align-items: center;
  margin-bottom: 15px;
}

.info-header i {
  font-size: 20px;
  color: var(--admin-primary);
  margin-right: 12px;
  width: 24px;
  text-align: center;
}

.info-header h3 {
  color: var(--admin-light);
  margin: 0;
  font-size: 16px;
  font-weight: 500;
}

.info-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.info-value {
  color: var(--admin-light);
  font-size: 18px;
  font-weight: 600;
  font-family: 'Courier New', monospace;
}

.status-indicator {
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.status-indicator.running {
  background: rgba(34, 197, 94, 0.2);
  color: #22c55e;
  border: 1px solid rgba(34, 197, 94, 0.3);
}

.status-indicator.stopped {
  background: rgba(239, 68, 68, 0.2);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.status-indicator.error {
  background: rgba(239, 68, 68, 0.2);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.status-indicator.warning {
  background: rgba(245, 158, 11, 0.2);
  color: #f59e0b;
  border: 1px solid rgba(245, 158, 11, 0.3);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .info-grid {
    grid-template-columns: 1fr;
  }

  .info-card {
    padding: 15px;
  }

  .info-header i {
    font-size: 18px;
  }

  .info-header h3 {
    font-size: 14px;
  }

  .info-value {
    font-size: 16px;
  }
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

.logs-modal {
  width: 900px;
}

.settings-modal {
  width: 1000px;
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

/* 日志样式 */
.log-filters {
  display: flex;
  gap: 20px;
  align-items: center;
  margin-bottom: 20px;
  padding: 15px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 8px;
  flex-wrap: wrap;
}

.filter-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.filter-group label {
  color: var(--admin-light);
  font-size: 14px;
  font-weight: 500;
}

.filter-group select,
.filter-group input {
  padding: 6px 12px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 6px;
  color: var(--admin-light);
  font-size: 13px;
}

.refresh-logs-btn {
  background: var(--admin-primary);
  border: none;
  color: white;
  padding: 6px 12px;
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  transition: all 0.3s ease;
}

.refresh-logs-btn:hover:not(:disabled) {
  background: #3b82f6;
}

.refresh-logs-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.logs-container {
  max-height: 400px;
  overflow-y: auto;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
}

.log-item {
  padding: 15px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  transition: all 0.3s ease;
}

.log-item:hover {
  background: rgba(255, 255, 255, 0.03);
}

.log-item:last-child {
  border-bottom: none;
}

.log-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.log-level {
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
}

.log-info .log-level {
  background: rgba(59, 130, 246, 0.2);
  color: #3b82f6;
}

.log-warning .log-level {
  background: rgba(245, 158, 11, 0.2);
  color: #f59e0b;
}

.log-error .log-level {
  background: rgba(239, 68, 68, 0.2);
  color: #ef4444;
}

.log-success .log-level {
  background: rgba(34, 197, 94, 0.2);
  color: #22c55e;
}

.log-module {
  color: var(--admin-primary);
  font-weight: 500;
  font-size: 13px;
}

.log-user {
  color: var(--admin-info);
  font-size: 12px;
  background: rgba(59, 130, 246, 0.1);
  padding: 2px 6px;
  border-radius: 4px;
}

.log-ip {
  color: var(--admin-gray);
  font-size: 11px;
  font-family: 'Courier New', monospace;
}

.log-time {
  color: var(--admin-gray);
  font-size: 12px;
  margin-left: auto;
}

.log-message {
  color: var(--admin-light);
  font-size: 14px;
  margin-bottom: 4px;
}

.log-details {
  color: var(--admin-gray);
  font-size: 12px;
  font-family: 'Courier New', monospace;
  background: rgba(0, 0, 0, 0.2);
  padding: 8px;
  border-radius: 4px;
  margin-top: 8px;
}

.no-logs,
.logs-loading {
  text-align: center;
  padding: 40px;
  color: var(--admin-gray);
}

.no-logs i,
.logs-loading i {
  font-size: 32px;
  margin-bottom: 10px;
  opacity: 0.5;
}

/* 设置样式 */
.settings-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 20px;
}

.setting-card {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  padding: 20px;
  transition: all 0.3s ease;
}

.setting-card:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(255, 255, 255, 0.2);
}

.setting-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.setting-header i {
  font-size: 20px;
  color: var(--admin-primary);
}

.setting-header h4 {
  margin: 0;
  color: var(--admin-light);
  font-size: 16px;
  font-weight: 600;
}

.setting-item {
  margin-bottom: 15px;
}

.setting-item:last-child {
  margin-bottom: 0;
}

.setting-item label {
  display: block;
  color: var(--admin-light);
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 8px;
}

.setting-item input,
.setting-item select {
  width: 100%;
  padding: 8px 12px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 6px;
  color: var(--admin-light);
  font-size: 13px;
}

.switch-group {
  display: flex;
  align-items: center;
}

.switch {
  position: relative;
  display: inline-block;
  width: 44px;
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
  transform: translateX(20px);
}

.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s ease;
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
</style>
