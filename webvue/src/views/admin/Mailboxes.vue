<template>
  <div class="admin-mailboxes">
    <div class="page-header">
      <h1>邮箱管理</h1>
      <button class="btn btn-primary" @click="createMailbox">
        <i class="fas fa-plus"></i>
        创建邮箱
      </button>
    </div>

    <div class="mailboxes-grid">
      <div class="mailbox-card" v-for="mailbox in mailboxes" :key="mailbox.id">
        <div class="mailbox-header">
          <div class="mailbox-icon">
            <i class="fas fa-envelope"></i>
          </div>
          <div class="mailbox-info">
            <h3>{{ mailbox.email }}</h3>
            <p>{{ mailbox.domain }}</p>
          </div>
          <div class="mailbox-status" :class="mailbox.status">
            <span class="status-badge" :class="mailbox.status === 'enabled' || mailbox.status === 'active' ? 'enabled' : 'disabled'">
              {{ (mailbox.status === 'enabled' || mailbox.status === 'active') ? '启用' : '禁用' }}
            </span>
          </div>
        </div>

        <div class="mailbox-stats">
          <div class="stat-item">
            <span>收件数</span>
            <span>{{ mailbox.inboxCount }}</span>
          </div>
          <div class="stat-item">
            <span>发件数</span>
            <span>{{ mailbox.sentCount }}</span>
          </div>
          <div class="stat-item">
            <span>存储</span>
            <span>{{ mailbox.storage }}</span>
          </div>
        </div>

        <div class="mailbox-actions">
          <button @click="viewMailbox(mailbox)" class="btn-view">
            <i class="fas fa-eye"></i>
            查看
          </button>
          <button @click="toggleMailboxStatus(mailbox)" class="btn-toggle" :class="(mailbox.status === 'enabled' || mailbox.status === 'active') ? 'btn-disable' : 'btn-enable'">
            <i :class="(mailbox.status === 'enabled' || mailbox.status === 'active') ? 'fas fa-ban' : 'fas fa-check-circle'"></i>
            {{ (mailbox.status === 'enabled' || mailbox.status === 'active') ? '禁用' : '启用' }}
          </button>
          <button @click="editMailboxAction(mailbox)" class="btn-edit">
            <i class="fas fa-edit"></i>
            编辑
          </button>
          <button @click="deleteMailbox(mailbox)" class="btn-delete danger">
            <i class="fas fa-trash"></i>
            删除
          </button>
        </div>
      </div>
    </div>

    <!-- 创建邮箱对话框 -->
    <div class="modal-overlay" v-if="showCreateDialog" @click="showCreateDialog = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>创建新邮箱</h3>
          <button class="close-btn" @click="showCreateDialog = false">
            <i class="fas fa-times"></i>
          </button>
        </div>

        <div class="modal-body">
          <div class="form-group">
            <label>选择用户</label>
            <select v-model="newMailbox.userId">
              <option v-for="user in availableUsers" :key="user.id" :value="user.id">
                {{ user.username }} ({{ user.email }})
              </option>
            </select>
          </div>

          <div class="form-group">
            <label>邮箱前缀</label>
            <input type="text" v-model="newMailbox.prefix" placeholder="例如: john, admin, support" />
            <small class="form-hint">只能包含字母、数字、点号、下划线和连字符</small>
          </div>

          <div class="form-group">
            <label>选择域名</label>
            <select v-model="newMailbox.domain">
              <option value="">请选择域名</option>
              <option v-for="domain in availableDomains" :key="domain.id" :value="domain.name">
                {{ domain.name }}
              </option>
            </select>
          </div>

          <div class="form-group">
            <label>邮箱预览</label>
            <div class="email-preview">
              {{ newMailbox.prefix || '[前缀]' }}@{{ newMailbox.domain || '[域名]' }}
            </div>
          </div>

          <div class="form-group">
            <label>密码</label>
            <input type="password" v-model="newMailbox.password" placeholder="请输入密码（至少6位）" />
          </div>
        </div>

        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showCreateDialog = false">取消</button>
          <button class="btn btn-primary" @click="submitCreateMailbox">创建邮箱</button>
        </div>
      </div>
    </div>

    <!-- 查看邮箱对话框 -->
    <div class="modal-overlay" v-if="showViewDialog" @click="showViewDialog = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>邮箱详情</h3>
          <button class="close-btn" @click="showViewDialog = false">
            <i class="fas fa-times"></i>
          </button>
        </div>

        <div class="modal-body">
          <div class="detail-group">
            <label>邮箱ID</label>
            <div class="detail-value">{{ selectedMailbox.id }}</div>
          </div>

          <div class="detail-group">
            <label>邮箱地址</label>
            <div class="detail-value">{{ selectedMailbox.email }}</div>
          </div>

          <div class="detail-group">
            <label>域名</label>
            <div class="detail-value">{{ selectedMailbox.domain }}</div>
          </div>

          <div class="detail-group">
            <label>状态</label>
            <div class="detail-value">
              <span :class="['status-badge', selectedMailbox.status]">
                {{ selectedMailbox.status === 'active' ? '活跃' : '非活跃' }}
              </span>
            </div>
          </div>

          <div class="detail-group">
            <label>收件箱邮件</label>
            <div class="detail-value">{{ selectedMailbox.inboxCount || 0 }} 封</div>
          </div>

          <div class="detail-group">
            <label>发件箱邮件</label>
            <div class="detail-value">{{ selectedMailbox.sentCount || 0 }} 封</div>
          </div>

          <div class="detail-group">
            <label>存储使用</label>
            <div class="detail-value">{{ selectedMailbox.storage || '0MB' }}</div>
          </div>

          <div class="detail-group">
            <label>所属用户</label>
            <div class="detail-value">{{ selectedMailbox.username || '未知' }}</div>
          </div>
        </div>

        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showViewDialog = false">关闭</button>
          <button class="btn btn-primary" @click="showViewDialog = false; editMailboxAction(selectedMailbox)">编辑邮箱</button>
        </div>
      </div>
    </div>

    <!-- 编辑邮箱对话框 -->
    <div class="modal-overlay" v-if="showEditDialog" @click="showEditDialog = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>编辑邮箱</h3>
          <button class="close-btn" @click="showEditDialog = false">
            <i class="fas fa-times"></i>
          </button>
        </div>

        <div class="modal-body">
          <div class="form-group">
            <label>邮箱地址</label>
            <input type="email" v-model="editMailbox.email" placeholder="请输入邮箱地址" readonly />
            <small class="form-hint">邮箱地址不可修改</small>
          </div>

          <div class="form-group">
            <label>新密码</label>
            <input type="password" v-model="editMailbox.password" placeholder="留空则不修改密码" />
            <small class="form-hint">留空则不修改密码，设置新密码需至少6位</small>
          </div>

          <div class="form-group">
            <label>状态</label>
            <select v-model="editMailbox.active">
              <option :value="true">活跃</option>
              <option :value="false">非活跃</option>
            </select>
          </div>
        </div>

        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showEditDialog = false">取消</button>
          <button class="btn btn-primary" @click="submitEditMailbox">保存更改</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import adminApi from '@/services/adminApi'

const mailboxes = ref([])
const loading = ref(false)
const showCreateDialog = ref(false)
const showViewDialog = ref(false)
const showEditDialog = ref(false)

const newMailbox = ref({
  prefix: '',
  domain: '',
  password: '',
  userId: 1 // 默认用户ID，后续可以改为选择用户
})

const selectedMailbox = ref({})
const editMailbox = ref({
  id: null,
  email: '',
  password: '',
  active: true
})

const availableDomains = ref([])
const availableUsers = ref([])

// 加载邮箱列表
const loadMailboxes = async () => {
  loading.value = true
  try {
    const response = await adminApi.getAllMailboxes()
    if (response.data.code === 0) {
      mailboxes.value = response.data.data.map(mailbox => ({
        id: mailbox.id,
        email: mailbox.email,
        domain: mailbox.domain,
        status: mailbox.is_active ? 'active' : 'inactive',
        inboxCount: mailbox.inbox_count || 0,
        sentCount: mailbox.sent_count || 0,
        storage: formatStorage(mailbox.storage_used || 0),
        username: mailbox.username
      }))
    }
  } catch (error) {
    console.error('加载邮箱列表失败:', error)
    alert('加载邮箱列表失败')
  } finally {
    loading.value = false
  }
}

const formatStorage = (bytes) => {
  if (bytes === 0) return '0MB'
  const mb = bytes / (1024 * 1024)
  if (mb < 1024) {
    return `${mb.toFixed(1)}MB`
  }
  return `${(mb / 1024).toFixed(1)}GB`
}

// 加载可用的域名和用户数据
const loadAvailableData = async () => {
  try {
    // 加载域名列表
    const domainsResponse = await adminApi.getAllDomains()
    if (domainsResponse.data.code === 0) {
      availableDomains.value = domainsResponse.data.data.map(domain => ({
        id: domain.id,
        name: domain.name,
        active: domain.active || domain.is_active
      })).filter(domain => domain.active) // 只显示活跃的域名
    }

    // 加载用户列表
    const usersResponse = await adminApi.getAllUsers()
    if (usersResponse.data.code === 0) {
      availableUsers.value = usersResponse.data.data.map(user => ({
        id: user.id,
        username: user.username,
        email: user.email
      }))
    }
  } catch (error) {
    console.error('加载可用数据失败:', error)
  }
}

const createMailbox = async () => {
  // 加载可用域名和用户
  await loadAvailableData()

  showCreateDialog.value = true
  newMailbox.value = {
    prefix: '',
    domain: '',
    password: '',
    userId: availableUsers.value.length > 0 ? availableUsers.value[0].id : 1
  }
}

const submitCreateMailbox = async () => {
  if (!newMailbox.value.prefix || !newMailbox.value.domain || !newMailbox.value.password) {
    alert('请填写邮箱前缀、选择域名和设置密码')
    return
  }

  if (newMailbox.value.password.length < 6) {
    alert('密码长度至少6位')
    return
  }

  // 验证邮箱前缀格式
  const prefixRegex = /^[a-zA-Z0-9][a-zA-Z0-9._-]*[a-zA-Z0-9]$|^[a-zA-Z0-9]$/
  if (!prefixRegex.test(newMailbox.value.prefix)) {
    alert('邮箱前缀格式不正确，只能包含字母、数字、点号、下划线和连字符')
    return
  }

  try {
    // 找到选中域名的ID
    const selectedDomain = availableDomains.value.find(d => d.name === newMailbox.value.domain)
    if (!selectedDomain) {
      alert('请选择有效的域名')
      return
    }

    const response = await adminApi.createMailbox({
      user_id: parseInt(newMailbox.value.userId),
      prefix: newMailbox.value.prefix,
      domain_id: selectedDomain.id,
      password: newMailbox.value.password
    })

    if (response.data.code === 0) {
      alert('邮箱创建成功')
      showCreateDialog.value = false
      loadMailboxes() // 重新加载邮箱列表
    } else {
      alert(response.data.msg || '创建邮箱失败')
    }
  } catch (error) {
    console.error('创建邮箱失败:', error)
    if (error.response && error.response.data) {
      alert(`创建邮箱失败: ${error.response.data.msg || error.response.data.message || '未知错误'}`)
    } else {
      alert('创建邮箱失败，请检查网络连接')
    }
  }
}

const viewMailbox = async (mailbox) => {
  try {
    selectedMailbox.value = { ...mailbox }
    showViewDialog.value = true
  } catch (error) {
    console.error('查看邮箱详情失败:', error)
    alert('查看邮箱详情失败')
  }
}

const editMailboxAction = async (mailbox) => {
  try {
    editMailbox.value = {
      id: mailbox.id,
      email: mailbox.email,
      password: '', // 密码不显示，需要重新设置
      active: mailbox.status === 'active'
    }
    showEditDialog.value = true
  } catch (error) {
    console.error('编辑邮箱失败:', error)
    alert('编辑邮箱失败')
  }
}

const submitEditMailbox = async () => {
  if (!editMailbox.value.email) {
    alert('邮箱地址不能为空')
    return
  }

  if (editMailbox.value.password && editMailbox.value.password.length < 6) {
    alert('密码长度至少6位')
    return
  }

  try {
    const updateData = {
      email: editMailbox.value.email,
      active: editMailbox.value.active
    }

    // 只有在设置了新密码时才发送密码
    if (editMailbox.value.password) {
      updateData.password = editMailbox.value.password
    }

    console.log('更新邮箱数据:', updateData)
    const response = await adminApi.updateMailbox(editMailbox.value.id, updateData)

    if (response.data.code === 0) {
      alert('邮箱更新成功')
      showEditDialog.value = false
      loadMailboxes() // 重新加载邮箱列表
    } else {
      alert(response.data.msg || '更新邮箱失败')
    }
  } catch (error) {
    console.error('更新邮箱失败:', error)
    if (error.response && error.response.data) {
      alert(`更新邮箱失败: ${error.response.data.msg || '未知错误'}`)
    } else {
      alert('更新邮箱失败，请检查网络连接')
    }
  }
}

const deleteMailbox = async (mailbox) => {
  if (!confirm(`确定删除邮箱 ${mailbox.email}？`)) {
    return
  }

  try {
    const response = await adminApi.deleteMailbox(mailbox.id)
    if (response.data.code === 0) {
      alert('邮箱删除成功')
      loadMailboxes() // 重新加载邮箱列表
    } else {
      alert(response.data.msg || '删除邮箱失败')
    }
  } catch (error) {
    console.error('删除邮箱失败:', error)
    alert('删除邮箱失败')
  }
}

// 切换邮箱状态
const toggleMailboxStatus = async (mailbox) => {
  const isCurrentlyEnabled = mailbox.status === 'enabled' || mailbox.status === 'active'
  const newStatus = isCurrentlyEnabled ? 'disabled' : 'enabled'
  const action = isCurrentlyEnabled ? '禁用' : '启用'

  if (!confirm(`确定${action}邮箱 ${mailbox.email}？${isCurrentlyEnabled ? '\n禁用后该邮箱将无法发送和接收邮件。' : ''}`)) {
    return
  }

  try {
    // 使用用户状态更新API，因为邮箱实际上就是用户
    const response = await adminApi.updateUserStatus(mailbox.id, newStatus)
    if (response.data.code === 0) {
      alert(`邮箱${action}成功`)
      loadMailboxes() // 重新加载邮箱列表
    } else {
      alert(response.data.msg || `${action}邮箱失败`)
    }
  } catch (error) {
    console.error(`${action}邮箱失败:`, error)
    alert(`${action}邮箱失败`)
  }
}

// 生命周期
onMounted(() => {
  loadMailboxes()
})
</script>

<style scoped>
.admin-mailboxes {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.page-header h1 {
  color: var(--admin-light);
  margin: 0;
}

.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
}

.btn-primary {
  background: linear-gradient(135deg, #3b82f6, #2563eb);
  color: white;
}

.mailboxes-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.mailbox-card {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  padding: 20px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.mailbox-header {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 15px;
}

.mailbox-icon {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  background: linear-gradient(135deg, #3b82f6, #2563eb);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.mailbox-info {
  flex: 1;
}

.mailbox-info h3 {
  color: var(--admin-light);
  margin: 0 0 4px 0;
  font-size: 16px;
}

.mailbox-info p {
  color: var(--admin-gray);
  margin: 0;
  font-size: 14px;
}

.mailbox-status {
  font-size: 18px;
}

.mailbox-status.active {
  color: #10b981;
}

.mailbox-status.inactive {
  color: #9ca3af;
}

.mailbox-stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
  margin-bottom: 15px;
  padding: 15px;
  background: rgba(15, 23, 42, 0.5);
  border-radius: 8px;
}

.stat-item {
  text-align: center;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-item span:first-child {
  font-size: 12px;
  color: var(--admin-gray);
}

.stat-item span:last-child {
  font-weight: 600;
  color: var(--admin-light);
}

.mailbox-actions {
  display: flex;
  gap: 8px;
  justify-content: center;
}

.mailbox-actions button {
  padding: 6px 12px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 12px;
  background: rgba(59, 130, 246, 0.2);
  color: #3b82f6;
}

.mailbox-actions button.danger {
  background: rgba(239, 68, 68, 0.2);
  color: #ef4444;
}

.mailbox-actions button:hover {
  opacity: 0.8;
}

/* 新增表单样式 */
.form-hint {
  display: block;
  margin-top: 4px;
  font-size: 12px;
  color: var(--admin-gray);
  font-style: italic;
}

.email-preview {
  padding: 8px 12px;
  background: rgba(37, 99, 235, 0.1);
  border: 1px solid rgba(37, 99, 235, 0.3);
  border-radius: 6px;
  color: var(--admin-primary);
  font-family: monospace;
  font-size: 14px;
  min-height: 20px;
}

.form-group select {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.05);
  color: var(--admin-light);
  font-size: 14px;
  transition: all 0.3s ease;
}

.form-group select:focus {
  outline: none;
  border-color: var(--admin-primary);
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1);
}

.form-group select option {
  background: var(--admin-dark);
  color: var(--admin-light);
}

/* 邮箱状态样式 */
.status-badge {
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.status-badge.enabled {
  background: rgba(16, 185, 129, 0.2);
  color: #10b981;
}

.status-badge.disabled {
  background: rgba(156, 163, 175, 0.2);
  color: #9ca3af;
}

/* 邮箱操作按钮样式 */
.btn-toggle {
  padding: 6px 12px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.btn-enable {
  background: rgba(16, 185, 129, 0.1);
  color: #10b981;
  border: 1px solid rgba(16, 185, 129, 0.3);
}

.btn-enable:hover {
  background: rgba(16, 185, 129, 0.2);
  transform: translateY(-1px);
}

.btn-disable {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.btn-disable:hover {
  background: rgba(239, 68, 68, 0.2);
  transform: translateY(-1px);
}
</style>
