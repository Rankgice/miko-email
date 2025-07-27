<template>
  <div class="admin-users">
    <div class="page-header">
      <h1>用户管理</h1>
      <button class="btn btn-primary" @click="addUser">
        <i class="fas fa-user-plus"></i>
        添加用户
      </button>
    </div>

    <div class="users-table">
      <div class="table-header">
        <div>用户名</div>
        <div>邮箱</div>
        <div>状态</div>
        <div>注册时间</div>
        <div>操作</div>
      </div>

      <div class="table-row" v-for="user in users" :key="user.id">
        <div>{{ user.username }}</div>
        <div>{{ user.email }}</div>
        <div>
          <span class="status-badge" :class="user.status">
            {{ user.status === 'enabled' ? '启用' : '禁用' }}
          </span>
        </div>
        <div>{{ user.registerTime }}</div>
        <div class="actions">
          <button @click="viewUser(user)" class="btn-view">
            <i class="fas fa-eye"></i>
            查看
          </button>
          <button @click="editUserAction(user)" class="btn-edit">
            <i class="fas fa-edit"></i>
            编辑
          </button>
          <button @click="deleteUser(user)" class="btn-delete danger">
            <i class="fas fa-trash"></i>
            删除
          </button>
        </div>
      </div>
    </div>

    <!-- 创建用户对话框 -->
    <div class="modal-overlay" v-if="showCreateDialog" @click="showCreateDialog = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>创建新用户</h3>
          <button class="close-btn" @click="showCreateDialog = false">
            <i class="fas fa-times"></i>
          </button>
        </div>

        <div class="modal-body">
          <div class="form-group">
            <label>用户名</label>
            <input type="text" v-model="newUser.username" placeholder="请输入用户名" />
          </div>

          <div class="form-group">
            <label>邮箱地址</label>
            <input type="email" v-model="newUser.email" placeholder="请输入邮箱地址" />
          </div>

          <div class="form-group">
            <label>密码</label>
            <input type="password" v-model="newUser.password" placeholder="请输入密码" />
          </div>

          <div class="form-group">
            <label>确认密码</label>
            <input type="password" v-model="newUser.confirmPassword" placeholder="请再次输入密码" />
          </div>
        </div>

        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showCreateDialog = false">取消</button>
          <button class="btn btn-primary" @click="createUser">创建用户</button>
        </div>
      </div>
    </div>

    <!-- 查看用户对话框 -->
    <div class="modal-overlay" v-if="showViewDialog" @click="showViewDialog = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>用户详情</h3>
          <button class="close-btn" @click="showViewDialog = false">
            <i class="fas fa-times"></i>
          </button>
        </div>

        <div class="modal-body">
          <div class="detail-group">
            <label>用户ID</label>
            <div class="detail-value">{{ selectedUser.id }}</div>
          </div>

          <div class="detail-group">
            <label>用户名</label>
            <div class="detail-value">{{ selectedUser.username }}</div>
          </div>

          <div class="detail-group">
            <label>邮箱地址</label>
            <div class="detail-value">{{ selectedUser.email }}</div>
          </div>

          <div class="detail-group">
            <label>状态</label>
            <div class="detail-value">
              <span :class="['status-badge', selectedUser.status]">
                {{ getStatusText(selectedUser.status) }}
              </span>
            </div>
          </div>

          <div class="detail-group">
            <label>注册时间</label>
            <div class="detail-value">{{ selectedUser.registerTime }}</div>
          </div>

          <div class="detail-group">
            <label>邮箱数量</label>
            <div class="detail-value">{{ selectedUser.mailboxCount || 0 }} 个</div>
          </div>
        </div>

        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showViewDialog = false">关闭</button>
          <button class="btn btn-primary" @click="showViewDialog = false; editUserAction(selectedUser)">编辑用户</button>
        </div>
      </div>
    </div>

    <!-- 编辑用户对话框 -->
    <div class="modal-overlay" v-if="showEditDialog" @click="showEditDialog = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>编辑用户</h3>
          <button class="close-btn" @click="showEditDialog = false">
            <i class="fas fa-times"></i>
          </button>
        </div>

        <div class="modal-body">
          <div class="form-group">
            <label>用户名</label>
            <input type="text" v-model="editUser.username" placeholder="请输入用户名" readonly />
            <small class="form-hint">用户名不可修改</small>
          </div>

          <div class="form-group">
            <label>邮箱地址</label>
            <input type="email" v-model="editUser.email" placeholder="请输入邮箱地址" readonly />
            <small class="form-hint">邮箱地址不可修改</small>
          </div>

          <div class="form-group">
            <label>状态</label>
            <select v-model="editUser.status">
              <option value="enabled">启用</option>
              <option value="disabled">禁用</option>
            </select>
            <small class="form-hint">禁用后该邮箱将无法发送和接收邮件</small>
          </div>
        </div>

        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showEditDialog = false">取消</button>
          <button class="btn btn-primary" @click="submitEditUser">保存更改</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import adminApi from '@/services/adminApi'

const users = ref([])
const loading = ref(false)
const showCreateDialog = ref(false)
const showViewDialog = ref(false)
const showEditDialog = ref(false)

const newUser = ref({
  username: '',
  email: '',
  password: '',
  confirmPassword: ''
})

const selectedUser = ref({})
const editUser = ref({
  id: null,
  username: '',
  email: '',
  status: 'enabled' // 默认启用状态
})

// 加载用户列表
const loadUsers = async () => {
  loading.value = true
  try {
    const response = await adminApi.getAllUsers()
    if (response.data.code === 0) {
      users.value = response.data.data.map(user => ({
        id: user.id,
        username: user.username,
        email: user.email,
        status: user.status || (user.is_active ? 'enabled' : 'disabled'), // 使用后端返回的status字段
        registerTime: new Date(user.created_at).toLocaleDateString(),
        mailboxCount: user.mailbox_count || 0
      }))
    }
  } catch (error) {
    console.error('加载用户列表失败:', error)
    alert('加载用户列表失败')
  } finally {
    loading.value = false
  }
}

const addUser = () => {
  showCreateDialog.value = true
  newUser.value = {
    username: '',
    email: '',
    password: '',
    confirmPassword: ''
  }
}

const createUser = async () => {
  if (!newUser.value.username || !newUser.value.email || !newUser.value.password) {
    alert('请填写所有必填字段')
    return
  }

  if (newUser.value.password !== newUser.value.confirmPassword) {
    alert('密码和确认密码不一致')
    return
  }

  if (newUser.value.password.length < 6) {
    alert('密码长度至少6位')
    return
  }

  try {
    const response = await adminApi.createUser({
      username: newUser.value.username,
      email: newUser.value.email,
      password: newUser.value.password,
      active: true // 默认创建启用状态的用户
    })

    if (response.data.code === 0) {
      alert('用户创建成功')
      showCreateDialog.value = false
      // 重置表单
      newUser.value = {
        username: '',
        email: '',
        password: '',
        confirmPassword: ''
      }
      loadUsers() // 重新加载用户列表
    } else {
      alert(response.data.msg || '创建用户失败')
    }
  } catch (error) {
    console.error('创建用户失败:', error)
    alert('创建用户失败')
  }
}

const viewUser = async (user) => {
  try {
    selectedUser.value = { ...user }
    showViewDialog.value = true
  } catch (error) {
    console.error('查看用户详情失败:', error)
    alert('查看用户详情失败')
  }
}

const editUserAction = async (user) => {
  try {
    editUser.value = {
      id: user.id,
      username: user.username,
      email: user.email,
      status: user.status // 直接使用状态字符串
    }
    showEditDialog.value = true
  } catch (error) {
    console.error('编辑用户失败:', error)
    alert('编辑用户失败')
  }
}

const submitEditUser = async () => {
  if (!editUser.value.username || !editUser.value.email) {
    alert('请填写用户名和邮箱')
    return
  }

  try {
    // 确保状态值正确
    const statusValue = editUser.value.status === 'enabled' ? 'enabled' : 'disabled'
    console.log('更新用户状态:', editUser.value.id, statusValue)

    const response = await adminApi.updateUserStatus(editUser.value.id, statusValue)

    if (response.data.code === 0) {
      alert('邮箱状态更新成功')
      showEditDialog.value = false
      loadUsers() // 重新加载用户列表
    } else {
      alert(response.data.msg || '更新邮箱状态失败')
    }
  } catch (error) {
    console.error('更新用户状态失败:', error)
    if (error.response && error.response.data) {
      alert(`更新用户状态失败: ${error.response.data.msg || '未知错误'}`)
    } else {
      alert('更新用户状态失败，请检查网络连接')
    }
  }
}

const deleteUser = async (user) => {
  if (!confirm(`确定删除用户 ${user.username}？`)) {
    return
  }

  try {
    const response = await adminApi.deleteUser(user.id)
    if (response.data.code === 0) {
      alert('用户删除成功')
      loadUsers() // 重新加载用户列表
    } else {
      alert(response.data.msg || '删除用户失败')
    }
  } catch (error) {
    console.error('删除用户失败:', error)
    alert('删除用户失败')
  }
}

// 状态文本转换
const getStatusText = (status) => {
  const statusMap = {
    'enabled': '启用',
    'disabled': '禁用',
    'active': '启用', // 兼容旧的状态值
    'suspended': '禁用', // 兼容旧的状态值
    'inactive': '禁用' // 兼容旧的状态值
  }
  return statusMap[status] || '未知'
}

// 生命周期
onMounted(() => {
  loadUsers()
})
</script>

<style scoped>
.admin-users {
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

.users-table {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.table-header {
  display: grid;
  grid-template-columns: 1fr 2fr 1fr 1fr 1fr;
  gap: 20px;
  padding: 20px;
  background: rgba(255, 255, 255, 0.05);
  font-weight: 600;
  color: var(--admin-light);
}

.table-row {
  display: grid;
  grid-template-columns: 1fr 2fr 1fr 1fr 1fr;
  gap: 20px;
  padding: 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
  color: var(--admin-light);
  align-items: center;
}

.table-row:hover {
  background: rgba(59, 130, 246, 0.05);
}

.status-badge {
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.status-badge.enabled,
.status-badge.active {
  background: rgba(16, 185, 129, 0.2);
  color: #10b981;
}

.status-badge.disabled,
.status-badge.inactive,
.status-badge.suspended {
  background: rgba(156, 163, 175, 0.2);
  color: #9ca3af;
}

.actions {
  display: flex;
  gap: 8px;
}

.actions button {
  padding: 6px 12px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 12px;
  background: rgba(59, 130, 246, 0.2);
  color: #3b82f6;
}

.actions button.danger {
  background: rgba(239, 68, 68, 0.2);
  color: #ef4444;
}

.actions button:hover {
  opacity: 0.8;
}

/* 模态框样式 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background: linear-gradient(135deg, rgba(30, 41, 59, 0.95), rgba(15, 23, 42, 0.98));
  border-radius: 12px;
  width: 90%;
  max-width: 500px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.modal-header h3 {
  color: var(--admin-light);
  margin: 0;
}

.close-btn {
  background: none;
  border: none;
  color: var(--admin-gray);
  font-size: 18px;
  cursor: pointer;
  padding: 5px;
  border-radius: 4px;
  transition: all 0.3s ease;
}

.close-btn:hover {
  color: var(--admin-light);
  background: rgba(255, 255, 255, 0.1);
}

.modal-body {
  padding: 20px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  color: var(--admin-light);
  font-weight: 500;
}

.form-group input {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.05);
  color: var(--admin-light);
  font-size: 14px;
  transition: all 0.3s ease;
}

.form-group input:focus {
  outline: none;
  border-color: var(--admin-primary);
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.1);
}

.form-group input::placeholder {
  color: var(--admin-gray);
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.btn-secondary {
  background: rgba(255, 255, 255, 0.1);
  color: var(--admin-gray);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.btn-secondary:hover {
  background: rgba(255, 255, 255, 0.15);
  color: var(--admin-light);
}

/* 操作按钮样式 */
.btn-view {
  background: rgba(34, 197, 94, 0.2);
  color: #22c55e;
  border: 1px solid rgba(34, 197, 94, 0.3);
}

.btn-view:hover {
  background: rgba(34, 197, 94, 0.3);
}

.btn-edit {
  background: rgba(59, 130, 246, 0.2);
  color: #3b82f6;
  border: 1px solid rgba(59, 130, 246, 0.3);
}

.btn-edit:hover {
  background: rgba(59, 130, 246, 0.3);
}

.btn-delete {
  background: rgba(239, 68, 68, 0.2);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.btn-delete:hover {
  background: rgba(239, 68, 68, 0.3);
}

/* 详情显示样式 */
.detail-group {
  margin-bottom: 20px;
}

.detail-group label {
  display: block;
  margin-bottom: 8px;
  color: var(--admin-light);
  font-weight: 500;
  font-size: 14px;
}

.detail-value {
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  color: var(--admin-light);
  font-size: 14px;
}

.status-badge {
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 500;
}

.status-badge.active {
  background: rgba(34, 197, 94, 0.2);
  color: #22c55e;
  border: 1px solid rgba(34, 197, 94, 0.3);
}

.status-badge.inactive,
.status-badge.suspended {
  background: rgba(156, 163, 175, 0.2);
  color: #9ca3af;
  border: 1px solid rgba(156, 163, 175, 0.3);
}
</style>
