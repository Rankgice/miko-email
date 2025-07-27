<template>
  <div class="mailboxes-page">
    <div class="page-header">
      <h1 class="page-title">邮箱管理</h1>
      <button class="btn btn-primary" @click="showAddMailboxModal" :disabled="loading">
        <i class="fas fa-plus"></i>
        添加邮箱
      </button>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-state">
      <i class="fas fa-spinner fa-spin"></i>
      <span>加载中...</span>
    </div>

    <div class="mailboxes-content" v-else>
      <!-- 邮箱列表 -->
      <div class="mailbox-grid" v-if="mailboxes.length > 0">
        <div class="mailbox-card" v-for="mailbox in mailboxes" :key="mailbox.id">
          <div class="mailbox-header">
            <div class="mailbox-icon">
              <i class="fas fa-envelope"></i>
            </div>
            <div class="mailbox-info">
              <h3 class="mailbox-email">{{ mailbox.email }}</h3>
              <p class="mailbox-domain">{{ getDomainFromEmail(mailbox.email) }}</p>
            </div>
            <div class="mailbox-status" :class="mailbox.is_active ? 'active' : 'inactive'">
              <i :class="mailbox.is_active ? 'fas fa-check-circle' : 'fas fa-times-circle'"></i>
              <span>{{ mailbox.is_active ? '活跃' : '禁用' }}</span>
            </div>
          </div>

          <div class="mailbox-stats">
            <div class="stat-item">
              <span class="stat-label">收件箱</span>
              <span class="stat-value">{{ mailbox.inbox_count || 0 }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">已发送</span>
              <span class="stat-value">{{ mailbox.sent_count || 0 }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">存储</span>
              <span class="stat-value">{{ formatStorage(mailbox.storage_used || 0) }}</span>
            </div>
          </div>

          <div class="mailbox-actions">
            <button class="action-btn" @click="editMailbox(mailbox)" title="编辑">
              <i class="fas fa-edit"></i>
            </button>
            <button class="action-btn" @click="viewMailbox(mailbox)" title="查看详情">
              <i class="fas fa-eye"></i>
            </button>
            <button class="action-btn" @click="toggleMailboxStatus(mailbox)" :title="mailbox.is_active ? '禁用' : '启用'">
              <i :class="mailbox.is_active ? 'fas fa-pause' : 'fas fa-play'"></i>
            </button>
            <button class="action-btn danger" @click="deleteMailbox(mailbox)" title="删除">
              <i class="fas fa-trash"></i>
            </button>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div class="empty-state" v-else>
        <i class="fas fa-folder-open"></i>
        <h3>暂无邮箱</h3>
        <p>您还没有创建任何邮箱</p>
        <button class="btn btn-primary" @click="showAddMailboxModal">
          <i class="fas fa-plus"></i>
          创建第一个邮箱
        </button>
      </div>
    </div>

    <!-- 添加/编辑邮箱弹窗 -->
    <div class="modal-overlay" v-if="showModal" @click="closeModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>{{ isEditing ? '编辑邮箱' : '添加邮箱' }}</h3>
          <button class="close-btn" @click="closeModal">
            <i class="fas fa-times"></i>
          </button>
        </div>

        <form @submit.prevent="saveMailbox" class="mailbox-form">
          <div class="form-group" v-if="!isEditing">
            <label for="emailPrefix">邮箱前缀 *</label>
            <input
              type="text"
              id="emailPrefix"
              v-model="mailboxForm.prefix"
              required
              placeholder="例如：user"
              :disabled="isEditing"
            >
          </div>

          <div class="form-group" v-if="!isEditing">
            <label for="domain">域名 *</label>
            <select
              id="domain"
              v-model="mailboxForm.domain_id"
              required
              :disabled="isEditing"
            >
              <option value="">请选择域名</option>
              <option
                v-for="domain in availableDomains"
                :key="domain.id"
                :value="domain.id"
              >
                {{ domain.name }}
              </option>
            </select>
            <p class="form-help" v-if="availableDomains.length === 0">
              没有可用域名，请先到 <router-link to="/user/domains">域名管理</router-link> 页面添加域名。
            </p>
          </div>

          <div class="form-group" v-if="isEditing">
            <label for="email">邮箱地址</label>
            <input
              type="email"
              id="email"
              v-model="mailboxForm.email"
              readonly
              class="readonly-input"
            >
          </div>

          <div class="form-group">
            <label for="password">{{ isEditing ? '新密码（留空不修改）' : '密码 *' }}</label>
            <input
              type="password"
              id="password"
              v-model="mailboxForm.password"
              :required="!isEditing"
              placeholder="至少6位字符"
              minlength="6"
            >
          </div>

          <div class="form-group checkbox-group" v-if="isEditing">
            <label class="checkbox-label">
              <input type="checkbox" v-model="mailboxForm.is_active">
              <span class="checkmark"></span>
              启用邮箱
            </label>
          </div>

          <div class="form-actions">
            <button type="button" class="btn btn-secondary" @click="closeModal">取消</button>
            <button type="submit" class="btn btn-primary" :disabled="submitting">
              <i v-if="submitting" class="fas fa-spinner fa-spin"></i>
              {{ submitting ? '保存中...' : '保存' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- 查看邮箱详情弹窗 -->
    <div class="modal-overlay" v-if="showViewModal" @click="closeViewModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>邮箱详情</h3>
          <button class="close-btn" @click="closeViewModal">
            <i class="fas fa-times"></i>
          </button>
        </div>

        <div class="mailbox-details" v-if="selectedMailbox">
          <div class="detail-section">
            <h4>基本信息</h4>
            <div class="detail-grid">
              <div class="detail-item">
                <span class="detail-label">邮箱地址:</span>
                <span class="detail-value">{{ selectedMailbox.email }}</span>
              </div>
              <div class="detail-item">
                <span class="detail-label">状态:</span>
                <span class="detail-value" :class="selectedMailbox.is_active ? 'status-active' : 'status-inactive'">
                  {{ selectedMailbox.is_active ? '活跃' : '禁用' }}
                </span>
              </div>
              <div class="detail-item">
                <span class="detail-label">创建时间:</span>
                <span class="detail-value">{{ formatDate(selectedMailbox.created_at) }}</span>
              </div>
              <div class="detail-item">
                <span class="detail-label">最后更新:</span>
                <span class="detail-value">{{ formatDate(selectedMailbox.updated_at) }}</span>
              </div>
            </div>
          </div>

          <div class="detail-section">
            <h4>邮件统计</h4>
            <div class="stats-grid">
              <div class="stats-card">
                <div class="stats-icon">
                  <i class="fas fa-inbox"></i>
                </div>
                <div class="stats-info">
                  <span class="stats-number">{{ selectedMailbox.inbox_count || 0 }}</span>
                  <span class="stats-label">收件箱</span>
                </div>
              </div>
              <div class="stats-card">
                <div class="stats-icon">
                  <i class="fas fa-paper-plane"></i>
                </div>
                <div class="stats-info">
                  <span class="stats-number">{{ selectedMailbox.sent_count || 0 }}</span>
                  <span class="stats-label">已发送</span>
                </div>
              </div>
              <div class="stats-card">
                <div class="stats-icon">
                  <i class="fas fa-hdd"></i>
                </div>
                <div class="stats-info">
                  <span class="stats-number">{{ formatStorage(selectedMailbox.storage_used || 0) }}</span>
                  <span class="stats-label">存储空间</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import userApi from '@/services/userApi'

// 响应式数据
const mailboxes = ref([])
const availableDomains = ref([])
const loading = ref(false)
const showModal = ref(false)
const showViewModal = ref(false)
const isEditing = ref(false)
const submitting = ref(false)
const selectedMailbox = ref(null)

// 表单数据
const mailboxForm = ref({
  prefix: '',
  domain_id: '',
  email: '',
  password: '',
  is_active: true
})

// 获取邮箱列表
const fetchMailboxes = async () => {
  try {
    loading.value = true
    const response = await userApi.getMailboxes()

    // 统一使用 code: 0 格式
    if (response.data.code === 0) {
      mailboxes.value = response.data.data || []
    } else {
      const errorMsg = response.data.msg || response.data.message || '获取邮箱列表失败'
      console.error('获取邮箱列表失败:', errorMsg)
      alert('获取邮箱列表失败: ' + errorMsg)
    }
  } catch (error) {
    console.error('获取邮箱列表失败:', error)
    const errorMsg = error.response?.data?.msg || error.response?.data?.message || error.message || '网络错误'
    alert('获取邮箱列表失败: ' + errorMsg)
  } finally {
    loading.value = false
  }
}

// 获取可用域名列表
const fetchAvailableDomains = async () => {
  try {
    const response = await userApi.getDomains()

    // 统一使用 code: 0 格式
    if (response.data.code === 0) {
      // API现在返回域名对象数组
      const domains = response.data.data || []
      availableDomains.value = domains.map((domain, index) => ({
        id: index + 1, // 临时ID，实际创建时使用域名名称
        name: domain.name || domain, // 兼容新旧格式
        is_active: domain.is_active !== undefined ? domain.is_active : true,
        is_verified: domain.is_verified !== undefined ? domain.is_verified : true
      }))
    } else {
      const errorMsg = response.data.msg || response.data.message || '获取域名列表失败'
      console.error('获取域名列表失败:', errorMsg)
    }
  } catch (error) {
    console.error('获取域名列表失败:', error)
  }
}

// 显示添加邮箱弹窗
const showAddMailboxModal = () => {
  isEditing.value = false
  resetForm()
  showModal.value = true
  fetchAvailableDomains()
}

// 编辑邮箱
const editMailbox = (mailbox) => {
  isEditing.value = true
  mailboxForm.value = {
    id: mailbox.id,
    email: mailbox.email,
    password: '',
    is_active: mailbox.is_active
  }
  showModal.value = true
}

// 查看邮箱详情
const viewMailbox = (mailbox) => {
  selectedMailbox.value = mailbox
  showViewModal.value = true
}

// 切换邮箱状态
const toggleMailboxStatus = async (mailbox) => {
  const newStatus = !mailbox.is_active
  const action = newStatus ? '启用' : '禁用'

  if (!confirm(`确定${action}邮箱 ${mailbox.email}？${!newStatus ? '\n禁用后该邮箱将无法发送和接收邮件。' : ''}`)) {
    return
  }

  try {
    const response = await userApi.toggleMailboxStatus(mailbox.id, newStatus)

    // 统一使用 code: 0 格式
    if (response.data.code === 0) {
      alert(`邮箱${action}成功`)
      fetchMailboxes() // 重新加载邮箱列表
    } else {
      const errorMsg = response.data.msg || response.data.message || `${action}失败`
      alert(`${action}失败: ` + errorMsg)
    }
  } catch (error) {
    console.error(`${action}邮箱失败:`, error)
    const errorMsg = error.response?.data?.msg || error.response?.data?.message || error.message || '网络错误'
    alert(`${action}邮箱失败: ` + errorMsg)
  }
}

// 删除邮箱
const deleteMailbox = async (mailbox) => {
  if (!confirm(`确定删除邮箱 ${mailbox.email}？\n删除后将无法恢复！`)) {
    return
  }

  try {
    const response = await userApi.deleteMailbox(mailbox.id)

    // 统一使用 code: 0 格式
    if (response.data.code === 0) {
      alert('邮箱删除成功')
      fetchMailboxes()
    } else {
      const errorMsg = response.data.msg || response.data.message || '删除失败'
      alert('删除失败: ' + errorMsg)
    }
  } catch (error) {
    console.error('删除邮箱失败:', error)
    const errorMsg = error.response?.data?.msg || error.response?.data?.message || error.message || '网络错误'
    alert('删除邮箱失败: ' + errorMsg)
  }
}

// 保存邮箱
const saveMailbox = async () => {
  try {
    submitting.value = true

    if (isEditing.value) {
      // 更新邮箱
      const updateData = {
        password: mailboxForm.value.password,
        is_active: mailboxForm.value.is_active
      }

      const response = await userApi.updateMailbox(mailboxForm.value.id, updateData)

      // 统一使用 code: 0 格式
      if (response.data.code === 0) {
        alert('邮箱更新成功')
        closeModal()
        fetchMailboxes()
      } else {
        const errorMsg = response.data.msg || response.data.message || '更新失败'
        alert('更新失败: ' + errorMsg)
      }
    } else {
      // 创建邮箱
      const selectedDomain = availableDomains.value.find(d => d.id === mailboxForm.value.domain_id)
      if (!selectedDomain) {
        alert('请选择有效的域名')
        return
      }

      const createData = {
        email: `${mailboxForm.value.prefix}@${selectedDomain.name}`,
        password: mailboxForm.value.password
      }

      const response = await userApi.createMailbox(createData)

      // 统一使用 code: 0 格式
      if (response.data.code === 0) {
        alert('邮箱创建成功')
        closeModal()
        fetchMailboxes()
      } else {
        const errorMsg = response.data.msg || response.data.message || '创建失败'
        alert('创建失败: ' + errorMsg)
      }
    }
  } catch (error) {
    console.error('保存邮箱失败:', error)
    const errorMsg = error.response?.data?.msg || error.response?.data?.message || error.message || '网络错误'
    alert('保存失败: ' + errorMsg)
  } finally {
    submitting.value = false
  }
}

// 工具函数
const getDomainFromEmail = (email) => {
  return email.split('@')[1] || ''
}

// 工具函数

const formatStorage = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const formatDate = (dateString) => {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleString('zh-CN')
}

// 重置表单
const resetForm = () => {
  mailboxForm.value = {
    prefix: '',
    domain_id: '',
    email: '',
    password: '',
    is_active: true
  }
}

// 关闭弹窗
const closeModal = () => {
  showModal.value = false
  isEditing.value = false
  resetForm()
}

const closeViewModal = () => {
  showViewModal.value = false
  selectedMailbox.value = null
}

// 生命周期
onMounted(() => {
  fetchMailboxes()
})
</script>

<style scoped>
.mailboxes-page {
  padding: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
}

.page-title {
  font-size: 32px;
  font-weight: 700;
  color: white;
  margin: 0;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
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

.btn-primary {
  background: linear-gradient(135deg, #00b4db 0%, #0083b0 100%);
  color: white;
  box-shadow: 0 4px 15px rgba(0, 180, 216, 0.3);
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(0, 180, 216, 0.4);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-secondary {
  background: #6c757d;
  color: white;
}

.btn-secondary:hover {
  background: #5a6268;
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

.mailboxes-content {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  padding: 32px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
  transition: all 0.3s ease;
}

.mailboxes-content:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(255, 255, 255, 0.2);
}

.mailbox-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 24px;
}

.mailbox-card {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.mailbox-card:hover {
  transform: translateY(-5px);
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(255, 255, 255, 0.2);
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
}

.mailbox-header {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
}

.mailbox-icon {
  width: 48px;
  height: 48px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16px;
  flex-shrink: 0;
}

.mailbox-icon i {
  color: white;
  font-size: 20px;
}

.mailbox-info {
  flex: 1;
  min-width: 0;
}

.mailbox-email {
  font-size: 16px;
  font-weight: 600;
  color: #2d3748;
  margin: 0 0 4px 0;
  word-break: break-all;
}

.mailbox-domain {
  font-size: 14px;
  color: #718096;
  margin: 0;
}

.mailbox-status {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  font-weight: 600;
  padding: 6px 12px;
  border-radius: 20px;
  margin-left: 12px;
}

.mailbox-status.active {
  background: rgba(72, 187, 120, 0.1);
  color: #38a169;
}

.mailbox-status.inactive {
  background: rgba(245, 101, 101, 0.1);
  color: #e53e3e;
}

.mailbox-stats {
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
  font-size: 18px;
  font-weight: 700;
  color: #2d3748;
}

.mailbox-actions {
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

.action-btn.danger {
  color: #e53e3e;
}

.action-btn.danger:hover {
  background: rgba(245, 101, 101, 0.1);
}

.empty-state {
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
  margin: 0 0 32px 0;
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
  max-width: 500px;
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

.mailbox-form {
  padding: 0 24px 24px 24px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 600;
  color: #2d3748;
  font-size: 14px;
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  background: white;
  color: #2d3748;
  font-size: 14px;
  transition: all 0.3s ease;
  box-sizing: border-box;
}

.form-group input:focus,
.form-group select:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.readonly-input {
  background: #f7fafc !important;
  color: #718096 !important;
  cursor: not-allowed;
}

.form-help {
  margin-top: 8px;
  font-size: 12px;
  color: #718096;
  line-height: 1.4;
}

.form-help a {
  color: #667eea;
  text-decoration: none;
}

.form-help a:hover {
  text-decoration: underline;
}

.checkbox-group {
  display: flex;
  align-items: center;
}

.checkbox-label {
  display: flex;
  align-items: center;
  cursor: pointer;
  font-weight: normal !important;
  margin-bottom: 0 !important;
}

.checkbox-label input[type="checkbox"] {
  width: auto;
  margin-right: 8px;
}

.form-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  margin-top: 32px;
  padding-top: 20px;
  border-top: 1px solid #e2e8f0;
}

/* 详情页面样式 */
.mailbox-details {
  padding: 0 24px 24px 24px;
}

.detail-section {
  margin-bottom: 32px;
}

.detail-section:last-child {
  margin-bottom: 0;
}

.detail-section h4 {
  margin: 0 0 16px 0;
  font-size: 16px;
  font-weight: 600;
  color: #2d3748;
}

.detail-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.detail-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.detail-label {
  font-size: 12px;
  color: #718096;
  font-weight: 600;
}

.detail-value {
  font-size: 14px;
  color: #2d3748;
  word-break: break-all;
}

.status-active {
  color: #38a169 !important;
  font-weight: 600;
}

.status-inactive {
  color: #e53e3e !important;
  font-weight: 600;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 16px;
}

.stats-card {
  background: #f7fafc;
  border-radius: 12px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
}

.stats-icon {
  width: 48px;
  height: 48px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stats-icon i {
  color: white;
  font-size: 20px;
}

.stats-info {
  display: flex;
  flex-direction: column;
}

.stats-number {
  font-size: 24px;
  font-weight: 700;
  color: #2d3748;
  line-height: 1;
}

.stats-label {
  font-size: 12px;
  color: #718096;
  margin-top: 4px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .mailboxes-page {
    padding: 16px;
  }

  .page-header {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }

  .mailbox-grid {
    grid-template-columns: 1fr;
  }

  .modal-content {
    width: 95%;
    margin: 20px;
  }

  .detail-grid {
    grid-template-columns: 1fr;
  }

  .stats-grid {
    grid-template-columns: 1fr;
  }
}
</style>
