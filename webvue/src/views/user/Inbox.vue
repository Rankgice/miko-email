<template>
  <div class="inbox-page">
    <div class="page-header">
      <h1 class="page-title">收件箱</h1>
      <div class="header-actions">
        <button class="btn btn-primary" @click="composeEmail">
          <i class="fas fa-edit"></i>
          写邮件
        </button>
        <button class="btn btn-secondary" @click="refreshEmails" :disabled="loading">
          <i class="fas fa-sync-alt" :class="{ 'fa-spin': loading }"></i>
          {{ loading ? '刷新中...' : '刷新' }}
        </button>
      </div>
    </div>

    <div class="inbox-content">
      <!-- 加载状态 -->
      <div class="loading-state" v-if="loading && emails.length === 0">
        <div class="loading-spinner">
          <i class="fas fa-spinner fa-spin"></i>
        </div>
        <p>正在加载邮件...</p>
      </div>

      <!-- 邮件列表 -->
      <div class="email-list" v-else-if="emails.length > 0">
        <div class="email-item"
             v-for="email in emails"
             :key="email.id"
             :class="{ unread: !email.is_read }"
             @click="selectEmail(email)">
          <div class="email-checkbox" @click.stop>
            <input type="checkbox" v-model="email.selected">
          </div>
          <div class="email-sender">
            <div class="sender-avatar">
              <i class="fas fa-user"></i>
            </div>
            <span class="sender-name">{{ email.from_name || email.from_email || '未知发件人' }}</span>
          </div>
          <div class="email-subject">{{ email.subject || '无主题' }}</div>
          <div class="email-preview">{{ email.preview || getEmailPreview(email.content) }}</div>

          <!-- 验证码显示 -->
          <div class="verification-codes" v-if="email.verificationCodes && email.verificationCodes.length > 0">
            <div class="code-label">验证码:</div>
            <div class="code-list">
              <span
                v-for="(code, index) in email.verificationCodes"
                :key="index"
                class="verification-code"
                @click.stop="copyCode(code)"
                :title="'点击复制: ' + code"
              >
                {{ code }}
                <i class="fas fa-copy"></i>
              </span>
            </div>
          </div>

          <div class="email-time">{{ formatTime(email.created_at || email.received_at) }}</div>
        </div>
      </div>

      <!-- 空状态 -->
      <div class="empty-state" v-else>
        <div class="empty-icon">
          <i class="fas fa-inbox"></i>
        </div>
        <h3>收件箱为空</h3>
        <p>您还没有收到任何邮件</p>
        <router-link to="/user/compose" class="btn btn-primary">
          <i class="fas fa-edit"></i>
          写第一封邮件
        </router-link>
      </div>
    </div>

    <!-- 邮件详情弹窗 -->
    <div class="modal-overlay" v-if="showEmailDetail" @click="closeEmailDetail">
      <div class="modal-content email-detail-modal" @click.stop>
        <div class="modal-header">
          <h3>{{ selectedEmail?.subject || '无主题' }}</h3>
          <button class="close-btn" @click="closeEmailDetail">
            <i class="fas fa-times"></i>
          </button>
        </div>

        <div class="modal-body">
          <div class="email-meta-info">
            <div class="meta-item">
              <label>发件人:</label>
              <span>{{ selectedEmail?.from_name || selectedEmail?.from_email }}</span>
            </div>
            <div class="meta-item">
              <label>收件人:</label>
              <span>{{ selectedEmail?.to_email }}</span>
            </div>
            <div class="meta-item">
              <label>时间:</label>
              <span>{{ formatTime(selectedEmail?.created_at) }}</span>
            </div>
          </div>

          <!-- 验证码区域 -->
          <div class="verification-section" v-if="selectedEmail?.verificationCodes && selectedEmail.verificationCodes.length > 0">
            <h4>验证码</h4>
            <div class="verification-codes-detail">
              <div
                v-for="(code, index) in selectedEmail.verificationCodes"
                :key="index"
                class="verification-code-item"
                @click="copyCode(code)"
              >
                <span class="code-value">{{ code }}</span>
                <button class="copy-btn">
                  <i class="fas fa-copy"></i>
                  复制
                </button>
              </div>
            </div>
          </div>

          <div class="email-content-detail">
            <h4>邮件内容</h4>
            <div class="content-body" v-html="selectedEmail?.content || selectedEmail?.html_content"></div>
          </div>
        </div>

        <div class="modal-footer">
          <button class="btn btn-secondary" @click="closeEmailDetail">关闭</button>
          <button class="btn btn-primary" @click="replyEmail" v-if="selectedEmail">
            <i class="fas fa-reply"></i>
            回复
          </button>
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
const emails = ref([])
const loading = ref(false)
const showEmailDetail = ref(false)
const selectedEmail = ref(null)

// 验证码正则表达式
const verificationCodePatterns = [
  /\b\d{4,8}\b/g,                    // 4-8位数字
  /验证码[：:\s]*(\d{4,8})/gi,        // 验证码：123456
  /code[：:\s]*(\d{4,8})/gi,         // code: 123456
  /\b[A-Z0-9]{4,8}\b/g,              // 4-8位大写字母数字组合
  /\b\d{6}\b/g,                      // 常见的6位验证码
]

// 提取验证码
const extractVerificationCodes = (content) => {
  if (!content) return []

  const codes = new Set()

  // 移除HTML标签
  const textContent = content.replace(/<[^>]*>/g, ' ')

  verificationCodePatterns.forEach(pattern => {
    const matches = textContent.match(pattern)
    if (matches) {
      matches.forEach(match => {
        // 清理匹配结果
        const cleanCode = match.replace(/验证码[：:\s]*/gi, '').replace(/code[：:\s]*/gi, '').trim()
        if (cleanCode.length >= 4 && cleanCode.length <= 8) {
          codes.add(cleanCode)
        }
      })
    }
  })

  return Array.from(codes)
}

// 获取邮件预览
const getEmailPreview = (content) => {
  if (!content) return '无内容预览'

  // 移除HTML标签
  const textContent = content.replace(/<[^>]*>/g, ' ')
  // 移除多余空白
  const cleanText = textContent.replace(/\s+/g, ' ').trim()
  // 截取前100个字符
  return cleanText.length > 100 ? cleanText.substring(0, 100) + '...' : cleanText
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
  } else if (days < 7) {
    return `${days}天前`
  } else {
    return time.toLocaleDateString()
  }
}

// 加载邮件列表
const loadEmails = async () => {
  loading.value = true
  try {
    const response = await userApi.getEmails({ type: 'inbox' })
    if (response.data.code === 0) {
      const emailList = response.data.data || []

      // 处理邮件数据并提取验证码
      emails.value = emailList.map(email => ({
        ...email,
        selected: false,
        verificationCodes: extractVerificationCodes(email.content || email.html_content)
      }))
    } else {
      console.error('获取邮件失败:', response.data.msg)
      emails.value = []
    }
  } catch (error) {
    console.error('加载邮件失败:', error)
    emails.value = []
  } finally {
    loading.value = false
  }
}

// 复制验证码
const copyCode = async (code) => {
  try {
    await navigator.clipboard.writeText(code)
    // 显示复制成功提示
    showToast(`验证码 ${code} 已复制到剪贴板`)
  } catch (error) {
    console.error('复制失败:', error)
    // 降级方案
    const textArea = document.createElement('textarea')
    textArea.value = code
    document.body.appendChild(textArea)
    textArea.select()
    document.execCommand('copy')
    document.body.removeChild(textArea)
    showToast(`验证码 ${code} 已复制到剪贴板`)
  }
}

// 显示提示消息
const showToast = (message) => {
  // 创建提示元素
  const toast = document.createElement('div')
  toast.className = 'toast-message'
  toast.textContent = message
  toast.style.cssText = `
    position: fixed;
    top: 20px;
    right: 20px;
    background: #4CAF50;
    color: white;
    padding: 12px 20px;
    border-radius: 6px;
    z-index: 10000;
    font-size: 14px;
    box-shadow: 0 4px 12px rgba(0,0,0,0.3);
  `

  document.body.appendChild(toast)

  // 3秒后移除
  setTimeout(() => {
    if (document.body.contains(toast)) {
      document.body.removeChild(toast)
    }
  }, 3000)
}

// 选择邮件
const selectEmail = async (email) => {
  selectedEmail.value = email
  showEmailDetail.value = true

  // 标记为已读
  if (!email.is_read) {
    try {
      await userApi.markAsRead(email.id)
      email.is_read = true
    } catch (error) {
      console.error('标记已读失败:', error)
    }
  }
}

// 关闭邮件详情
const closeEmailDetail = () => {
  showEmailDetail.value = false
  selectedEmail.value = null
}

// 回复邮件
const replyEmail = () => {
  if (selectedEmail.value) {
    router.push({
      path: '/user/compose',
      query: {
        reply_to: selectedEmail.value.id,
        subject: `Re: ${selectedEmail.value.subject}`,
        to: selectedEmail.value.from_email
      }
    })
  }
}

// 写邮件
const composeEmail = () => {
  router.push('/user/compose')
}

// 刷新邮件
const refreshEmails = async () => {
  await loadEmails()
}

// 生命周期
onMounted(() => {
  loadEmails()
})
</script>

<style scoped>
.inbox-page {
  max-width: 1200px;
  margin: 0 auto;
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
  color: var(--admin-light);
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 15px;
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

.btn-primary {
  background: var(--admin-primary);
  color: white;
  box-shadow: 0 4px 12px rgba(37, 99, 235, 0.3);
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(37, 99, 235, 0.4);
}

.btn-secondary {
  background: rgba(255, 255, 255, 0.1);
  color: var(--admin-light);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.btn-secondary:hover {
  background: rgba(255, 255, 255, 0.15);
  transform: translateY(-2px);
  border-color: rgba(255, 255, 255, 0.3);
}

.inbox-content {
  background: linear-gradient(135deg, rgba(30, 41, 59, 0.9), rgba(15, 23, 42, 0.95));
  border-radius: 12px;
  border: 1px solid var(--border);
  overflow: hidden;
}

.email-list {
  display: flex;
  flex-direction: column;
}

.email-item {
  display: grid;
  grid-template-columns: 40px 200px 1fr 150px 120px;
  gap: 20px;
  padding: 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  cursor: pointer;
  transition: all 0.3s ease;
  align-items: center;
}

.email-item:last-child {
  border-bottom: none;
}

.email-item:hover {
  background: rgba(0, 180, 216, 0.05);
}

.email-item.unread {
  background: rgba(37, 99, 235, 0.05);
  border-left: 3px solid var(--admin-primary);
}

.email-checkbox input {
  width: 16px;
  height: 16px;
  cursor: pointer;
}

.email-sender {
  display: flex;
  align-items: center;
  gap: 12px;
}

.sender-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--admin-primary), #60a5fa);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 14px;
}

.sender-name {
  font-weight: 500;
  color: var(--admin-light);
  font-size: 14px;
}

.email-subject {
  font-weight: 600;
  color: var(--admin-light);
  font-size: 16px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.email-item.unread .email-subject {
  color: var(--admin-primary);
}

.email-preview {
  color: var(--text-secondary);
  font-size: 14px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.email-time {
  color: var(--admin-gray);
  font-size: 12px;
  text-align: right;
}

/* 验证码样式 */
.verification-codes {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.code-label {
  font-size: 12px;
  color: var(--admin-gray);
  font-weight: 500;
}

.code-list {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.verification-code {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  background: linear-gradient(135deg, #10b981, #059669);
  color: white;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 600;
  font-family: 'Courier New', monospace;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.verification-code:hover {
  background: linear-gradient(135deg, #059669, #047857);
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(16, 185, 129, 0.3);
}

.verification-code i {
  font-size: 10px;
  opacity: 0.8;
}

/* 加载状态 */
.loading-state {
  text-align: center;
  padding: 60px 20px;
  color: var(--admin-gray);
}

.loading-spinner {
  font-size: 32px;
  color: var(--admin-primary);
  margin-bottom: 20px;
}

.loading-state p {
  font-size: 14px;
  color: var(--admin-gray);
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: var(--admin-gray);
}

.empty-icon {
  font-size: 64px;
  color: var(--admin-gray);
  margin-bottom: 20px;
  opacity: 0.5;
}

.empty-state h3 {
  font-size: 20px;
  margin-bottom: 10px;
  color: var(--admin-light);
  font-weight: 600;
}

.empty-state p {
  font-size: 14px;
  margin-bottom: 30px;
  color: var(--admin-gray);
}

.empty-state .btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  text-decoration: none;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 20px;
    align-items: stretch;
  }

  .header-actions {
    justify-content: center;
  }

  .email-item {
    grid-template-columns: 1fr;
    gap: 10px;
    padding: 15px;
  }

  .email-sender,
  .email-subject,
  .email-preview,
  .email-time {
    grid-column: 1;
  }

  .email-checkbox {
    position: absolute;
    top: 15px;
    right: 15px;
  }
}
</style>
