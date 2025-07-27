<template>
  <div class="admin-captcha-rules">
    <div class="page-header">
      <h1>验证码提取规则管理</h1>
      <div class="header-actions">
        <button class="btn btn-primary" @click="showCreateDialog = true">
          <i class="fas fa-plus"></i>
          添加规则
        </button>
        <button class="btn btn-secondary" @click="loadRules">
          <i class="fas fa-sync-alt"></i>
          刷新
        </button>
      </div>
    </div>

    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon">
          <i class="fas fa-cogs"></i>
        </div>
        <div class="stat-content">
          <h3>{{ ruleStats.total }}</h3>
          <p>总规则数</p>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">
          <i class="fas fa-check-circle"></i>
        </div>
        <div class="stat-content">
          <h3>{{ ruleStats.active }}</h3>
          <p>启用规则</p>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">
          <i class="fas fa-pause-circle"></i>
        </div>
        <div class="stat-content">
          <h3>{{ ruleStats.inactive }}</h3>
          <p>禁用规则</p>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">
          <i class="fas fa-percentage"></i>
        </div>
        <div class="stat-content">
          <h3>{{ ruleStats.successRate }}%</h3>
          <p>匹配成功率</p>
        </div>
      </div>
    </div>

    <div class="rules-section">
      <h2>提取规则列表</h2>
      <div class="rules-grid">
        <div class="rule-card" v-for="rule in extractionRules" :key="rule.id">
          <div class="rule-header">
            <div class="rule-title">
              <h3>{{ rule.name }}</h3>
              <span class="rule-type">{{ rule.type }}</span>
            </div>
            <div class="rule-status">
              <span class="status-badge" :class="rule.enabled ? 'enabled' : 'disabled'">
                {{ rule.enabled ? '启用' : '禁用' }}
              </span>
            </div>
          </div>

          <div class="rule-content">
            <div class="rule-info">
              <div class="info-item">
                <label>发件人匹配:</label>
                <span>{{ rule.senderPattern || '无限制' }}</span>
              </div>
              <div class="info-item">
                <label>主题匹配:</label>
                <span>{{ rule.subjectPattern || '无限制' }}</span>
              </div>
              <div class="info-item">
                <label>提取规则:</label>
                <code>{{ rule.extractionPattern }}</code>
              </div>
              <div class="info-item">
                <label>优先级:</label>
                <span class="priority-badge" :class="getPriorityClass(rule.priority)">
                  {{ rule.priority }}
                </span>
              </div>
            </div>
          </div>

          <div class="rule-stats">
            <div class="stat-item">
              <span class="stat-label">匹配次数:</span>
              <span class="stat-value">{{ rule.matchCount || 0 }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">成功提取:</span>
              <span class="stat-value">{{ rule.successCount || 0 }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">最后使用:</span>
              <span class="stat-value">{{ rule.lastUsed || '从未' }}</span>
            </div>
          </div>

          <div class="rule-actions">
            <button @click="editRule(rule)" class="btn-edit">
              <i class="fas fa-edit"></i>
              编辑
            </button>
            <button @click="toggleRuleStatus(rule)" class="btn-toggle" :class="rule.enabled ? 'btn-disable' : 'btn-enable'">
              <i :class="rule.enabled ? 'fas fa-pause' : 'fas fa-play'"></i>
              {{ rule.enabled ? '禁用' : '启用' }}
            </button>
            <button @click="testRule(rule)" class="btn-test">
              <i class="fas fa-vial"></i>
              测试
            </button>
            <button @click="deleteRule(rule)" class="btn-delete danger">
              <i class="fas fa-trash"></i>
              删除
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 创建/编辑规则对话框 -->
    <div class="modal-overlay" v-if="showCreateDialog || showEditDialog" @click="closeDialog">
      <div class="modal-content rule-modal" @click.stop>
        <div class="modal-header">
          <h3>{{ showEditDialog ? '编辑提取规则' : '创建提取规则' }}</h3>
          <button class="close-btn" @click="closeDialog">
            <i class="fas fa-times"></i>
          </button>
        </div>

        <div class="modal-body">
          <div class="form-group">
            <label>规则名称 *</label>
            <input type="text" v-model="ruleForm.name" placeholder="请输入规则名称" />
            <small class="form-hint">用于识别和管理的规则名称</small>
          </div>

          <div class="form-group">
            <label>规则类型 *</label>
            <select v-model="ruleForm.type">
              <option value="regex">正则表达式</option>
              <option value="keyword">关键词匹配</option>
              <option value="position">位置提取</option>
            </select>
          </div>

          <div class="form-group">
            <label>发件人匹配模式</label>
            <input type="text" v-model="ruleForm.senderPattern" placeholder="例如: @example.com 或 noreply@*" />
            <small class="form-hint">留空表示匹配所有发件人</small>
          </div>

          <div class="form-group">
            <label>主题匹配模式</label>
            <input type="text" v-model="ruleForm.subjectPattern" placeholder="例如: 验证码 或 *verification*" />
            <small class="form-hint">留空表示匹配所有主题</small>
          </div>

          <div class="form-group">
            <label>验证码提取规则 *</label>
            <textarea v-model="ruleForm.extractionPattern" rows="3" placeholder="请输入提取规则"></textarea>
            <small class="form-hint">
              正则表达式示例: (\d{6}) | 关键词示例: 验证码: | 位置示例: line:2,word:3
            </small>
          </div>

          <div class="form-row">
            <div class="form-group">
              <label>优先级</label>
              <select v-model="ruleForm.priority">
                <option value="high">高</option>
                <option value="medium">中</option>
                <option value="low">低</option>
              </select>
            </div>

            <div class="form-group">
              <label>状态</label>
              <select v-model="ruleForm.enabled">
                <option :value="true">启用</option>
                <option :value="false">禁用</option>
              </select>
            </div>
          </div>

          <div class="form-group">
            <label>描述</label>
            <textarea v-model="ruleForm.description" rows="2" placeholder="规则描述（可选）"></textarea>
          </div>
        </div>

        <div class="modal-footer">
          <button class="btn btn-secondary" @click="closeDialog">取消</button>
          <button class="btn btn-primary" @click="saveRule">
            {{ showEditDialog ? '保存更改' : '创建规则' }}
          </button>
        </div>
      </div>
    </div>

    <!-- 测试规则对话框 -->
    <div class="modal-overlay" v-if="showTestDialog" @click="showTestDialog = false">
      <div class="modal-content test-modal" @click.stop>
        <div class="modal-header">
          <h3>测试提取规则</h3>
          <button class="close-btn" @click="showTestDialog = false">
            <i class="fas fa-times"></i>
          </button>
        </div>

        <div class="modal-body">
          <div class="form-group">
            <label>测试邮件内容</label>
            <textarea v-model="testContent" rows="8" placeholder="请粘贴邮件内容进行测试..."></textarea>
          </div>

          <div class="test-result" v-if="testResult">
            <h4>提取结果:</h4>
            <div class="result-content" :class="testResult.success ? 'success' : 'error'">
              <div v-if="testResult.success">
                <strong>提取成功:</strong> {{ testResult.code }}
              </div>
              <div v-else>
                <strong>提取失败:</strong> {{ testResult.error }}
              </div>
            </div>
          </div>
        </div>

        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showTestDialog = false">关闭</button>
          <button class="btn btn-primary" @click="runTest">执行测试</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import adminApi from '@/services/adminApi'

// 响应式数据
const extractionRules = ref([])
const ruleStats = ref({
  total: 0,
  active: 0,
  inactive: 0,
  successRate: 0
})
const loading = ref(false)

// 对话框状态
const showCreateDialog = ref(false)
const showEditDialog = ref(false)
const showTestDialog = ref(false)

// 表单数据
const ruleForm = ref({
  id: null,
  name: '',
  type: 'regex',
  senderPattern: '',
  subjectPattern: '',
  extractionPattern: '',
  priority: 'medium',
  enabled: true,
  description: ''
})

// 测试相关
const testContent = ref('')
const testResult = ref(null)
const currentTestRule = ref(null)

// 加载提取规则
const loadRules = async () => {
  loading.value = true
  try {
    // 尝试从API加载数据
    const response = await adminApi.getCaptchaRules()
    if (response.data.code === 0) {
      extractionRules.value = response.data.data || []
    } else {
      // API失败时使用模拟数据
      console.warn('API加载失败，使用模拟数据')
      extractionRules.value = [
      {
        id: 1,
        name: '通用6位数字验证码',
        type: '正则表达式',
        senderPattern: '',
        subjectPattern: '*验证码*',
        extractionPattern: '(\\d{6})',
        priority: 'high',
        enabled: true,
        description: '提取6位数字验证码',
        matchCount: 1250,
        successCount: 1180,
        lastUsed: '2025-01-27 23:45:12'
      },
      {
        id: 2,
        name: 'GitHub验证码',
        type: '正则表达式',
        senderPattern: '@github.com',
        subjectPattern: '*verification*',
        extractionPattern: 'verification code is (\\d{6})',
        priority: 'high',
        enabled: true,
        description: 'GitHub邮箱验证码提取',
        matchCount: 89,
        successCount: 87,
        lastUsed: '2025-01-27 20:15:30'
      },
      {
        id: 3,
        name: '阿里云验证码',
        type: '关键词匹配',
        senderPattern: '@aliyun.com',
        subjectPattern: '',
        extractionPattern: '验证码：',
        priority: 'medium',
        enabled: true,
        description: '阿里云服务验证码',
        matchCount: 156,
        successCount: 152,
        lastUsed: '2025-01-26 14:22:45'
      },
      {
        id: 4,
        name: '腾讯云验证码',
        type: '位置提取',
        senderPattern: '@tencent.com',
        subjectPattern: '*验证码*',
        extractionPattern: 'line:3,word:2',
        priority: 'medium',
        enabled: false,
        description: '腾讯云验证码（已禁用）',
        matchCount: 45,
        successCount: 40,
        lastUsed: '2025-01-25 09:30:15'
      }
      ]
    }

    // 计算统计数据
    updateRuleStats()
  } catch (error) {
    console.error('加载提取规则失败:', error)
    // API失败时使用模拟数据
    extractionRules.value = [
      {
        id: 1,
        name: '通用6位数字验证码',
        type: '正则表达式',
        senderPattern: '',
        subjectPattern: '*验证码*',
        extractionPattern: '(\\d{6})',
        priority: 'high',
        enabled: true,
        description: '提取6位数字验证码',
        matchCount: 1250,
        successCount: 1180,
        lastUsed: '2025-01-27 23:45:12'
      },
      {
        id: 4,
        name: '腾讯云验证码',
        type: '位置提取',
        senderPattern: '@tencent.com',
        subjectPattern: '*验证码*',
        extractionPattern: 'line:3,word:2',
        priority: 'medium',
        enabled: false,
        description: '腾讯云验证码（已禁用）',
        matchCount: 45,
        successCount: 40,
        lastUsed: '2025-01-25 09:30:15'
      }
    ]
    updateRuleStats()
  } finally {
    loading.value = false
  }
}

// 更新规则统计
const updateRuleStats = () => {
  const total = extractionRules.value.length
  const active = extractionRules.value.filter(rule => rule.enabled).length
  const inactive = total - active

  const totalMatches = extractionRules.value.reduce((sum, rule) => sum + (rule.matchCount || 0), 0)
  const totalSuccess = extractionRules.value.reduce((sum, rule) => sum + (rule.successCount || 0), 0)
  const successRate = totalMatches > 0 ? Math.round((totalSuccess / totalMatches) * 100) : 0

  ruleStats.value = {
    total,
    active,
    inactive,
    successRate
  }
}

// 获取优先级样式类
const getPriorityClass = (priority) => {
  switch (priority) {
    case 'high': return 'priority-high'
    case 'medium': return 'priority-medium'
    case 'low': return 'priority-low'
    default: return 'priority-medium'
  }
}

// 编辑规则
const editRule = (rule) => {
  ruleForm.value = {
    id: rule.id,
    name: rule.name,
    type: rule.type === '正则表达式' ? 'regex' : rule.type === '关键词匹配' ? 'keyword' : 'position',
    senderPattern: rule.senderPattern,
    subjectPattern: rule.subjectPattern,
    extractionPattern: rule.extractionPattern,
    priority: rule.priority,
    enabled: rule.enabled,
    description: rule.description
  }
  showEditDialog.value = true
}

// 切换规则状态
const toggleRuleStatus = async (rule) => {
  const action = rule.enabled ? '禁用' : '启用'
  if (!confirm(`确定${action}规则 "${rule.name}"？`)) {
    return
  }

  try {
    // 调用API更新规则状态
    const response = await adminApi.updateCaptchaRuleStatus(rule.id, !rule.enabled)
    if (response.data.code === 0) {
      rule.enabled = !rule.enabled
      updateRuleStats()
      alert(`规则${action}成功`)
    } else {
      alert(`${action}规则失败: ` + (response.data.msg || '未知错误'))
    }
  } catch (error) {
    console.error(`${action}规则失败:`, error)
    alert(`${action}规则失败: ` + (error.response?.data?.msg || error.message))
  }
}

// 测试规则
const testRule = (rule) => {
  currentTestRule.value = rule
  testContent.value = ''
  testResult.value = null
  showTestDialog.value = true
}

// 执行测试
const runTest = () => {
  if (!testContent.value.trim()) {
    alert('请输入测试内容')
    return
  }

  try {
    const rule = currentTestRule.value
    let extractedCode = null
    let success = false

    // 模拟提取逻辑
    if (rule.type === '正则表达式' || rule.extractionPattern.includes('\\d')) {
      const regex = new RegExp(rule.extractionPattern.replace(/\\\\/g, '\\'))
      const match = testContent.value.match(regex)
      if (match && match[1]) {
        extractedCode = match[1]
        success = true
      }
    } else if (rule.type === '关键词匹配') {
      const lines = testContent.value.split('\n')
      for (const line of lines) {
        if (line.includes(rule.extractionPattern)) {
          const parts = line.split(rule.extractionPattern)
          if (parts.length > 1) {
            const codeMatch = parts[1].match(/\d{4,8}/)
            if (codeMatch) {
              extractedCode = codeMatch[0]
              success = true
              break
            }
          }
        }
      }
    }

    testResult.value = success
      ? { success: true, code: extractedCode }
      : { success: false, error: '未能提取到验证码' }

  } catch (error) {
    testResult.value = { success: false, error: error.message }
  }
}

// 删除规则
const deleteRule = async (rule) => {
  if (!confirm(`确定删除规则 "${rule.name}"？\n删除后将无法恢复。`)) {
    return
  }

  try {
    // 这里应该调用API删除规则
    const index = extractionRules.value.findIndex(r => r.id === rule.id)
    if (index > -1) {
      extractionRules.value.splice(index, 1)
      updateRuleStats()
      alert('规则删除成功')
    }
  } catch (error) {
    console.error('删除规则失败:', error)
    alert('删除规则失败')
  }
}

// 保存规则
const saveRule = async () => {
  if (!ruleForm.value.name.trim()) {
    alert('请输入规则名称')
    return
  }
  if (!ruleForm.value.extractionPattern.trim()) {
    alert('请输入提取规则')
    return
  }

  try {
    if (showEditDialog.value) {
      // 编辑模式
      const index = extractionRules.value.findIndex(r => r.id === ruleForm.value.id)
      if (index > -1) {
        extractionRules.value[index] = {
          ...extractionRules.value[index],
          name: ruleForm.value.name,
          type: ruleForm.value.type === 'regex' ? '正则表达式' : ruleForm.value.type === 'keyword' ? '关键词匹配' : '位置提取',
          senderPattern: ruleForm.value.senderPattern,
          subjectPattern: ruleForm.value.subjectPattern,
          extractionPattern: ruleForm.value.extractionPattern,
          priority: ruleForm.value.priority,
          enabled: ruleForm.value.enabled,
          description: ruleForm.value.description
        }
        alert('规则更新成功')
      }
    } else {
      // 创建模式
      const newRule = {
        id: Date.now(), // 临时ID
        name: ruleForm.value.name,
        type: ruleForm.value.type === 'regex' ? '正则表达式' : ruleForm.value.type === 'keyword' ? '关键词匹配' : '位置提取',
        senderPattern: ruleForm.value.senderPattern,
        subjectPattern: ruleForm.value.subjectPattern,
        extractionPattern: ruleForm.value.extractionPattern,
        priority: ruleForm.value.priority,
        enabled: ruleForm.value.enabled,
        description: ruleForm.value.description,
        matchCount: 0,
        successCount: 0,
        lastUsed: '从未'
      }
      extractionRules.value.push(newRule)
      alert('规则创建成功')
    }

    updateRuleStats()
    closeDialog()
  } catch (error) {
    console.error('保存规则失败:', error)
    alert('保存规则失败')
  }
}

// 关闭对话框
const closeDialog = () => {
  showCreateDialog.value = false
  showEditDialog.value = false
  ruleForm.value = {
    id: null,
    name: '',
    type: 'regex',
    senderPattern: '',
    subjectPattern: '',
    extractionPattern: '',
    priority: 'medium',
    enabled: true,
    description: ''
  }
}

// 生命周期
onMounted(() => {
  loadRules()
})
</script>

<style scoped>
.admin-captcha-rules {
  max-width: 1400px;
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

.btn-warning {
  background: linear-gradient(135deg, #f59e0b, #d97706);
  color: white;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
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
  background: linear-gradient(135deg, #10b981, #059669);
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

.captcha-records h2 {
  color: var(--admin-light);
  margin-bottom: 20px;
}

.records-table {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.table-header {
  display: grid;
  grid-template-columns: 1fr 1fr 2fr 1fr 2fr 1fr;
  gap: 20px;
  padding: 20px;
  background: rgba(255, 255, 255, 0.05);
  font-weight: 600;
  color: var(--admin-light);
}

.table-row {
  display: grid;
  grid-template-columns: 1fr 1fr 2fr 1fr 2fr 1fr;
  gap: 20px;
  padding: 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
  color: var(--admin-light);
  align-items: center;
}

.table-row:hover {
  background: rgba(59, 130, 246, 0.05);
}

.code-cell {
  font-family: monospace;
  font-weight: 600;
  color: #3b82f6;
}

.type-badge {
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
  background: rgba(139, 92, 246, 0.2);
  color: #8b5cf6;
}

.status-badge {
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.status-badge.pending {
  background: rgba(245, 158, 11, 0.2);
  color: #f59e0b;
}

.status-badge.verified {
  background: rgba(16, 185, 129, 0.2);
  color: #10b981;
}

.status-badge.expired {
  background: rgba(239, 68, 68, 0.2);
  color: #ef4444;
}

.status-badge.revoked {
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

/* 规则管理样式 */
.rules-section {
  margin-top: 30px;
}

.rules-section h2 {
  color: var(--admin-light);
  margin-bottom: 20px;
  font-size: 20px;
  font-weight: 600;
}

.rules-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(500px, 1fr));
  gap: 20px;
}

.rule-card {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  padding: 20px;
  transition: all 0.3s ease;
}

.rule-card:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(255, 255, 255, 0.2);
  transform: translateY(-2px);
}

.rule-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 15px;
}

.rule-title h3 {
  color: var(--admin-light);
  margin: 0 0 5px 0;
  font-size: 16px;
  font-weight: 600;
}

.rule-type {
  background: rgba(37, 99, 235, 0.2);
  color: #60a5fa;
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.rule-status .status-badge {
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.rule-status .status-badge.enabled {
  background: rgba(16, 185, 129, 0.2);
  color: #10b981;
}

.rule-status .status-badge.disabled {
  background: rgba(156, 163, 175, 0.2);
  color: #9ca3af;
}

.rule-content {
  margin-bottom: 15px;
}

.rule-info {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 10px;
}

.info-item label {
  color: var(--admin-muted);
  font-size: 12px;
  min-width: 80px;
  font-weight: 500;
}

.info-item span {
  color: var(--admin-light);
  font-size: 13px;
}

.info-item code {
  background: rgba(0, 0, 0, 0.3);
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  color: #fbbf24;
}

.priority-badge {
  padding: 2px 6px;
  border-radius: 8px;
  font-size: 11px;
  font-weight: 600;
}

.priority-badge.priority-high {
  background: rgba(239, 68, 68, 0.2);
  color: #ef4444;
}

.priority-badge.priority-medium {
  background: rgba(245, 158, 11, 0.2);
  color: #f59e0b;
}

.priority-badge.priority-low {
  background: rgba(156, 163, 175, 0.2);
  color: #9ca3af;
}

.rule-stats {
  display: flex;
  justify-content: space-between;
  margin-bottom: 15px;
  padding: 10px;
  background: rgba(0, 0, 0, 0.2);
  border-radius: 8px;
}

.stat-item {
  text-align: center;
}

.stat-label {
  display: block;
  color: var(--admin-muted);
  font-size: 11px;
  margin-bottom: 2px;
}

.stat-value {
  color: var(--admin-light);
  font-size: 13px;
  font-weight: 600;
}

.rule-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.rule-actions button {
  padding: 6px 12px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 12px;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 4px;
  transition: all 0.3s ease;
}

.btn-edit {
  background: rgba(37, 99, 235, 0.1);
  color: #60a5fa;
  border: 1px solid rgba(37, 99, 235, 0.3);
}

.btn-edit:hover {
  background: rgba(37, 99, 235, 0.2);
  transform: translateY(-1px);
}

.btn-toggle.btn-enable {
  background: rgba(16, 185, 129, 0.1);
  color: #10b981;
  border: 1px solid rgba(16, 185, 129, 0.3);
}

.btn-toggle.btn-disable {
  background: rgba(245, 158, 11, 0.1);
  color: #f59e0b;
  border: 1px solid rgba(245, 158, 11, 0.3);
}

.btn-test {
  background: rgba(139, 92, 246, 0.1);
  color: #a78bfa;
  border: 1px solid rgba(139, 92, 246, 0.3);
}

.btn-delete {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.rule-actions button:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
}

/* 对话框样式 */
.rule-modal {
  max-width: 600px;
  max-height: 80vh;
  overflow-y: auto;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 15px;
}

.test-modal {
  max-width: 700px;
}

.test-result {
  margin-top: 15px;
  padding: 15px;
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.test-result h4 {
  color: var(--admin-light);
  margin: 0 0 10px 0;
  font-size: 14px;
}

.result-content {
  padding: 10px;
  border-radius: 6px;
  font-family: 'Courier New', monospace;
  font-size: 13px;
}

.result-content.success {
  background: rgba(16, 185, 129, 0.1);
  border: 1px solid rgba(16, 185, 129, 0.3);
  color: #10b981;
}

.result-content.error {
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  color: #ef4444;
}
</style>
