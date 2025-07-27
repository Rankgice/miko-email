<template>
  <div class="forward-rules-page">
    <div class="page-header">
      <h1 class="page-title">转发规则</h1>
      <button class="btn btn-primary" @click="showAddRuleModal" :disabled="loading">
        <i class="fas fa-plus"></i>
        添加规则
      </button>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-state">
      <i class="fas fa-spinner fa-spin"></i>
      <span>加载中...</span>
    </div>

    <div class="rules-content" v-else>
      <div class="rules-container">
        <!-- 规则列表 -->
        <div class="rules-list" v-if="forwardRules.length > 0">
          <div class="rule-card" v-for="rule in forwardRules" :key="rule.id">
          <div class="rule-header">
            <div class="rule-status" :class="rule.enabled ? 'enabled' : 'disabled'">
              <i :class="rule.enabled ? 'fas fa-toggle-on' : 'fas fa-toggle-off'"></i>
            </div>
            <h3 class="rule-name">{{ rule.description || `${rule.source_email} → ${rule.target_email}` }}</h3>
            <div class="rule-actions">
              <button class="action-btn" @click="editRule(rule)" title="编辑">
                <i class="fas fa-edit"></i>
              </button>
              <button class="action-btn" @click="testRule(rule)" title="测试">
                <i class="fas fa-vial"></i>
              </button>
              <button class="action-btn" @click="toggleRule(rule)" :title="rule.enabled ? '禁用' : '启用'">
                <i :class="rule.enabled ? 'fas fa-pause' : 'fas fa-play'"></i>
              </button>
              <button class="action-btn danger" @click="deleteRule(rule)" title="删除">
                <i class="fas fa-trash"></i>
              </button>
            </div>
          </div>

          <div class="rule-details">
            <div class="rule-info">
              <div class="info-item">
                <span class="info-label">源邮箱:</span>
                <span class="info-value">{{ rule.source_email }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">目标邮箱:</span>
                <span class="info-value">{{ rule.target_email }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">保留原邮件:</span>
                <span class="info-value">{{ rule.keep_original ? '是' : '否' }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">转发附件:</span>
                <span class="info-value">{{ rule.forward_attachments ? '是' : '否' }}</span>
              </div>
              <div class="info-item" v-if="rule.subject_prefix">
                <span class="info-label">主题前缀:</span>
                <span class="info-value">{{ rule.subject_prefix }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">转发次数:</span>
                <span class="info-value">{{ rule.forward_count || 0 }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">创建时间:</span>
                <span class="info-value">{{ formatDate(rule.created_at) }}</span>
              </div>
            </div>
          </div>
        </div>
        </div>

        <!-- 空状态 -->
        <div class="empty-state" v-else>
          <i class="fas fa-share"></i>
          <h3>暂无转发规则</h3>
          <p>您还没有创建任何转发规则</p>
          <button class="btn btn-primary" @click="showAddRuleModal">
            <i class="fas fa-plus"></i>
            创建第一个规则
          </button>
        </div>
      </div>
    </div>

    <!-- 添加/编辑规则弹窗 -->
    <div class="modal-overlay" v-if="showModal" @click="closeModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>{{ isEditing ? '编辑转发规则' : '添加转发规则' }}</h3>
          <button class="close-btn" @click="closeModal">
            <i class="fas fa-times"></i>
          </button>
        </div>

        <form @submit.prevent="saveRule" class="rule-form">
          <div class="form-group">
            <label for="sourceEmail">源邮箱 *</label>
            <select
              id="sourceEmail"
              v-model="ruleForm.source_email"
              required
            >
              <option value="">请选择源邮箱</option>
              <option
                v-for="mailbox in userMailboxes"
                :key="mailbox.id"
                :value="mailbox.email"
              >
                {{ mailbox.email }}
              </option>
            </select>
            <p class="form-help" v-if="userMailboxes.length === 0">
              您还没有创建任何邮箱，请先到 <router-link to="/user/mailboxes">邮箱管理</router-link> 页面创建邮箱。
            </p>
          </div>

          <div class="form-group">
            <label for="targetEmail">目标邮箱 *</label>
            <input
              type="email"
              id="targetEmail"
              v-model="ruleForm.target_email"
              required
              placeholder="target@domain.com"
            >
          </div>

          <div class="form-group">
            <label for="description">规则描述</label>
            <input
              type="text"
              id="description"
              v-model="ruleForm.description"
              placeholder="规则描述（可选）"
            >
          </div>

          <div class="form-group">
            <label for="subjectPrefix">主题前缀</label>
            <input
              type="text"
              id="subjectPrefix"
              v-model="ruleForm.subject_prefix"
              placeholder="例如：[转发]"
            >
          </div>

          <div class="form-group checkbox-group">
            <label class="checkbox-label">
              <input type="checkbox" v-model="ruleForm.enabled">
              <span class="checkmark"></span>
              启用规则
            </label>
          </div>

          <div class="form-group checkbox-group">
            <label class="checkbox-label">
              <input type="checkbox" v-model="ruleForm.keep_original">
              <span class="checkmark"></span>
              保留原邮件
            </label>
          </div>

          <div class="form-group checkbox-group">
            <label class="checkbox-label">
              <input type="checkbox" v-model="ruleForm.forward_attachments">
              <span class="checkmark"></span>
              转发附件
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

    <!-- 测试规则弹窗 -->
    <div class="modal-overlay" v-if="showTestModal" @click="closeTestModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>测试转发规则</h3>
          <button class="close-btn" @click="closeTestModal">
            <i class="fas fa-times"></i>
          </button>
        </div>

        <form @submit.prevent="runTest" class="test-form">
          <div class="form-group">
            <label for="testSubject">测试邮件主题</label>
            <input
              type="text"
              id="testSubject"
              v-model="testForm.subject"
              placeholder="输入测试邮件主题"
              required
            >
          </div>

          <div class="form-group">
            <label for="testContent">测试邮件内容</label>
            <textarea
              id="testContent"
              v-model="testForm.content"
              placeholder="输入测试邮件内容"
              rows="4"
              required
            ></textarea>
          </div>

          <div class="form-actions">
            <button type="button" class="btn btn-secondary" @click="closeTestModal">取消</button>
            <button type="submit" class="btn btn-primary" :disabled="testing">
              <i v-if="testing" class="fas fa-spinner fa-spin"></i>
              {{ testing ? '测试中...' : '开始测试' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import userApi from '@/services/userApi'

// 响应式数据
const forwardRules = ref([])
const userMailboxes = ref([])
const loading = ref(false)
const showModal = ref(false)
const showTestModal = ref(false)
const isEditing = ref(false)
const submitting = ref(false)
const testing = ref(false)
const currentRule = ref(null)

// 表单数据
const ruleForm = ref({
  source_email: '',
  target_email: '',
  description: '',
  subject_prefix: '',
  enabled: true,
  keep_original: true,
  forward_attachments: true
})

const testForm = ref({
  subject: '',
  content: ''
})

// 重置表单
const resetForm = () => {
  ruleForm.value = {
    source_email: '',
    target_email: '',
    description: '',
    subject_prefix: '',
    enabled: true,
    keep_original: true,
    forward_attachments: true
  }
}

// 获取转发规则列表
const fetchForwardRules = async () => {
  try {
    loading.value = true
    const response = await userApi.getForwardRules()
    if (response.data.success) {
      forwardRules.value = response.data.data || []
    } else {
      console.error('获取转发规则失败:', response.data.message)
    }
  } catch (error) {
    console.error('获取转发规则失败:', error)
  } finally {
    loading.value = false
  }
}

// 获取用户邮箱列表
const fetchUserMailboxes = async () => {
  try {
    const response = await userApi.getMailboxes()
    if (response.data.success) {
      userMailboxes.value = response.data.data || []
    } else {
      console.error('获取邮箱列表失败:', response.data.message)
    }
  } catch (error) {
    console.error('获取邮箱列表失败:', error)
  }
}

// 显示添加规则弹窗
const showAddRuleModal = () => {
  isEditing.value = false
  resetForm()
  showModal.value = true
}

// 编辑规则
const editRule = (rule) => {
  isEditing.value = true
  currentRule.value = rule
  ruleForm.value = {
    source_email: rule.source_email,
    target_email: rule.target_email,
    description: rule.description || '',
    subject_prefix: rule.subject_prefix || '',
    enabled: rule.enabled,
    keep_original: rule.keep_original,
    forward_attachments: rule.forward_attachments
  }
  showModal.value = true
}

// 保存规则
const saveRule = async () => {
  try {
    submitting.value = true

    if (isEditing.value) {
      // 更新规则
      const response = await userApi.updateForwardRule(currentRule.value.id, ruleForm.value)
      if (response.data.code === 0) {
        await fetchForwardRules()
        closeModal()
        alert('规则更新成功')
      } else {
        alert('更新失败: ' + (response.data.msg || response.data.message || '更新失败'))
      }
    } else {
      // 创建规则
      const response = await userApi.createForwardRule(ruleForm.value)
      if (response.data.code === 0) {
        await fetchForwardRules()
        closeModal()
        alert('规则创建成功')
      } else {
        alert('创建失败: ' + (response.data.msg || response.data.message || '创建失败'))
      }
    }
  } catch (error) {
    console.error('保存规则失败:', error)
    alert('保存失败: ' + (error.response?.data?.message || error.message))
  } finally {
    submitting.value = false
  }
}

// 切换规则状态
const toggleRule = async (rule) => {
  try {
    const newStatus = !rule.enabled
    const response = await userApi.toggleForwardRule(rule.id, newStatus)
    if (response.data.code === 0) {
      rule.enabled = newStatus
      alert(`规则已${newStatus ? '启用' : '禁用'}`)
    } else {
      alert('操作失败: ' + (response.data.msg || response.data.message || '操作失败'))
    }
  } catch (error) {
    console.error('切换规则状态失败:', error)
    alert('操作失败: ' + (error.response?.data?.msg || error.response?.data?.message || error.message))
  }
}

// 删除规则
const deleteRule = async (rule) => {
  if (!confirm(`确定要删除规则 "${rule.description || rule.source_email}" 吗？`)) {
    return
  }

  try {
    const response = await userApi.deleteForwardRule(rule.id)
    if (response.data.success) {
      await fetchForwardRules()
      alert('规则删除成功')
    } else {
      alert('删除失败: ' + response.data.message)
    }
  } catch (error) {
    console.error('删除规则失败:', error)
    alert('删除失败: ' + (error.response?.data?.message || error.message))
  }
}

// 测试规则
const testRule = (rule) => {
  currentRule.value = rule
  testForm.value = {
    subject: '测试邮件主题',
    content: '这是一封测试邮件内容'
  }
  showTestModal.value = true
}

// 运行测试
const runTest = async () => {
  try {
    testing.value = true
    const response = await userApi.testForwardRule(currentRule.value.id, testForm.value)
    if (response.data.success) {
      alert('测试成功！转发规则工作正常。')
      closeTestModal()
    } else {
      alert('测试失败: ' + response.data.message)
    }
  } catch (error) {
    console.error('测试规则失败:', error)
    alert('测试失败: ' + (error.response?.data?.message || error.message))
  } finally {
    testing.value = false
  }
}

// 关闭弹窗
const closeModal = () => {
  showModal.value = false
  isEditing.value = false
  currentRule.value = null
  resetForm()
}

const closeTestModal = () => {
  showTestModal.value = false
  currentRule.value = null
  testForm.value = {
    subject: '',
    content: ''
  }
}

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

// 生命周期
onMounted(() => {
  fetchForwardRules()
  fetchUserMailboxes()
})
</script>

<style scoped>
.forward-rules-page {
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

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0;
}

.btn {
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

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-primary {
  background: linear-gradient(135deg, var(--primary), #0077b6);
  color: white;
  box-shadow: 0 5px 15px rgba(0, 180, 216, 0.2);
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(0, 180, 216, 0.3);
}

.btn-secondary {
  background: var(--bg-secondary);
  color: var(--text-secondary);
  border: 1px solid var(--border);
}

.btn-secondary:hover {
  background: var(--bg-hover);
}

.loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 60px;
  color: var(--text-secondary);
  font-size: 16px;
}

.rules-content {
  min-height: 400px;
}

.rules-container {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  padding: 32px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(10px);
  transition: all 0.3s ease;
}

.rules-container:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(255, 255, 255, 0.2);
}

.rules-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.rule-card {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  padding: 25px;
  transition: all 0.3s ease;
}

.rule-card:hover {
  transform: translateY(-5px);
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(255, 255, 255, 0.2);
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
}

.rule-header {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 20px;
}

.rule-status {
  font-size: 24px;
}

.rule-status.enabled {
  color: var(--success);
}

.rule-status.disabled {
  color: var(--text-secondary);
}

.rule-name {
  flex: 1;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.rule-actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  width: 36px;
  height: 36px;
  border-radius: 6px;
  background: rgba(15, 23, 42, 0.5);
  border: 1px solid var(--border);
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.action-btn:hover {
  background: rgba(0, 180, 216, 0.1);
  color: var(--primary);
  border-color: var(--primary);
}

.action-btn.danger:hover {
  background: rgba(255, 107, 107, 0.1);
  color: var(--accent);
  border-color: var(--accent);
}

.rule-details {
  margin-top: 20px;
}

.rule-info {
  background: rgba(15, 23, 42, 0.5);
  border-radius: 8px;
  padding: 20px;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 15px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.info-item:last-child {
  border-bottom: none;
}

.info-label {
  font-size: 14px;
  color: var(--text-secondary);
  font-weight: 500;
}

.info-value {
  font-size: 14px;
  color: var(--text-primary);
  font-weight: 600;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: var(--text-secondary);
}

.empty-state i {
  font-size: 64px;
  margin-bottom: 20px;
  opacity: 0.5;
}

.empty-state h3 {
  font-size: 20px;
  margin-bottom: 10px;
  color: var(--text-primary);
}

.empty-state p {
  font-size: 16px;
  margin-bottom: 20px;
}

/* 弹窗样式 */
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
  padding: 20px;
}

.modal-content {
  background: var(--bg-primary);
  border-radius: 12px;
  border: 1px solid var(--border);
  max-width: 500px;
  width: 100%;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 25px;
  border-bottom: 1px solid var(--border);
}

.modal-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
}

.close-btn {
  background: none;
  border: none;
  color: var(--text-secondary);
  cursor: pointer;
  padding: 5px;
  border-radius: 4px;
  transition: all 0.3s ease;
}

.close-btn:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

/* 表单样式 */
.rule-form,
.test-form {
  padding: 25px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
}

.form-group input,
.form-group textarea,
.form-group select {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid var(--border);
  border-radius: 8px;
  background: var(--bg-secondary);
  color: var(--text-primary);
  font-size: 14px;
  transition: all 0.3s ease;
  box-sizing: border-box;
}

.form-group input:focus,
.form-group textarea:focus,
.form-group select:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(0, 180, 216, 0.1);
}

.form-help {
  margin-top: 8px;
  font-size: 12px;
  color: var(--text-secondary);
  line-height: 1.4;
}

.form-help a {
  color: var(--primary);
  text-decoration: none;
}

.form-help a:hover {
  text-decoration: underline;
}

.form-group textarea {
  resize: vertical;
  min-height: 100px;
}

.checkbox-group {
  display: flex;
  align-items: center;
  gap: 10px;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  font-size: 14px;
  color: var(--text-primary);
  margin-bottom: 0 !important;
}

.checkbox-label input[type="checkbox"] {
  width: auto;
  margin: 0;
}

.checkmark {
  position: relative;
}

.form-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  margin-top: 30px;
  padding-top: 20px;
  border-top: 1px solid var(--border);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 20px;
    align-items: stretch;
  }

  .rule-info {
    grid-template-columns: 1fr;
  }

  .rule-header {
    flex-wrap: wrap;
    gap: 10px;
  }

  .rule-actions {
    order: -1;
    width: 100%;
    justify-content: center;
  }

  .modal-content {
    margin: 10px;
    max-width: none;
  }

  .form-actions {
    flex-direction: column;
  }
}

@media (max-width: 480px) {
  .forward-rules-page {
    padding: 10px;
  }

  .rule-card {
    padding: 20px;
  }

  .rule-info {
    padding: 15px;
  }

  .info-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 5px;
  }

  .modal-header,
  .rule-form,
  .test-form {
    padding: 20px;
  }
}
</style>
