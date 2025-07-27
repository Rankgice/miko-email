<template>
  <div class="outbox-page">
    <div class="page-header">
      <h1 class="page-title">发件箱</h1>
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

    <div class="outbox-content">
      <!-- 加载状态 -->
      <div class="loading-state" v-if="loading && sentEmails.length === 0">
        <div class="loading-spinner">
          <i class="fas fa-spinner fa-spin"></i>
        </div>
        <p>正在加载已发送邮件...</p>
      </div>

      <!-- 邮件列表 -->
      <div class="email-list" v-else-if="sentEmails.length > 0">
        <div class="email-item"
             v-for="email in sentEmails"
             :key="email.id"
             @click="viewEmail(email)">
          <div class="email-checkbox" @click.stop>
            <input type="checkbox" v-model="email.selected">
          </div>
          <div class="email-recipient">
            <div class="recipient-avatar">
              <i class="fas fa-user"></i>
            </div>
            <span class="recipient-name">{{ email.to_name || email.to_email || '未知收件人' }}</span>
          </div>
          <div class="email-subject">{{ email.subject || '无主题' }}</div>
          <div class="email-preview">{{ email.preview || getEmailPreview(email.content) }}</div>
          <div class="email-status" :class="email.status">
            <i :class="getStatusIcon(email.status)"></i>
            {{ getStatusText(email.status) }}
          </div>
          <div class="email-time">{{ formatTime(email.sent_at || email.created_at) }}</div>
        </div>
      </div>

      <!-- 空状态 -->
      <div class="empty-state" v-else>
        <div class="empty-icon">
          <i class="fas fa-paper-plane"></i>
        </div>
        <h3>发件箱为空</h3>
        <p>您还没有发送任何邮件</p>
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
              <label>收件人:</label>
              <span>{{ selectedEmail?.to_name || selectedEmail?.to_email }}</span>
            </div>
            <div class="meta-item" v-if="selectedEmail?.cc_email">
              <label>抄送:</label>
              <span>{{ selectedEmail.cc_email }}</span>
            </div>
            <div class="meta-item" v-if="selectedEmail?.bcc_email">
              <label>密送:</label>
              <span>{{ selectedEmail.bcc_email }}</span>
            </div>
            <div class="meta-item">
              <label>发送时间:</label>
              <span>{{ formatTime(selectedEmail?.sent_at || selectedEmail?.created_at) }}</span>
            </div>
            <div class="meta-item">
              <label>状态:</label>
              <span class="status-badge" :class="selectedEmail?.status">
                <i :class="getStatusIcon(selectedEmail?.status)"></i>
                {{ getStatusText(selectedEmail?.status) }}
              </span>
            </div>
          </div>

          <div class="email-content-detail">
            <h4>邮件内容</h4>
            <div class="content-body" v-html="selectedEmail?.content || selectedEmail?.html_content"></div>
          </div>

          <!-- 附件列表 -->
          <div class="attachments-section" v-if="selectedEmail?.attachments && selectedEmail.attachments.length > 0">
            <h4>附件</h4>
            <div class="attachment-list">
              <div class="attachment-item" v-for="attachment in selectedEmail.attachments" :key="attachment.id">
                <div class="attachment-icon">
                  <i class="fas fa-file"></i>
                </div>
                <div class="attachment-info">
                  <span class="attachment-name">{{ attachment.filename }}</span>
                  <span class="attachment-size">({{ formatFileSize(attachment.size) }})</span>
                </div>
                <button class="download-btn" @click="downloadAttachment(attachment)">
                  <i class="fas fa-download"></i>
                </button>
              </div>
            </div>
          </div>
        </div>

        <div class="modal-footer">
          <button class="btn btn-secondary" @click="closeEmailDetail">关闭</button>
          <button class="btn btn-primary" @click="forwardEmail" v-if="selectedEmail">
            <i class="fas fa-share"></i>
            转发
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
const sentEmails = ref([])
const loading = ref(false)
const showEmailDetail = ref(false)
const selectedEmail = ref(null)

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

// 格式化文件大小
const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 加载已发送邮件列表
const loadSentEmails = async () => {
  loading.value = true
  try {
    const response = await userApi.getEmails({ type: 'sent' })
    if (response.data.code === 0) {
      const emailList = response.data.data || []

      // 处理邮件数据
      sentEmails.value = emailList.map(email => ({
        ...email,
        selected: false,
        status: email.delivery_status || 'sent' // 使用delivery_status或默认为sent
      }))
    } else {
      console.error('获取已发送邮件失败:', response.data.msg)
      sentEmails.value = []
    }
  } catch (error) {
    console.error('加载已发送邮件失败:', error)
    sentEmails.value = []
  } finally {
    loading.value = false
  }
}

// 查看邮件详情
const viewEmail = (email) => {
  selectedEmail.value = email
  showEmailDetail.value = true
}

// 关闭邮件详情
const closeEmailDetail = () => {
  showEmailDetail.value = false
  selectedEmail.value = null
}

// 转发邮件
const forwardEmail = () => {
  if (selectedEmail.value) {
    router.push({
      path: '/user/compose',
      query: {
        forward: selectedEmail.value.id,
        subject: `Fwd: ${selectedEmail.value.subject}`,
        content: `\n\n---------- 转发邮件 ----------\n发件人: ${selectedEmail.value.from_email}\n收件人: ${selectedEmail.value.to_email}\n主题: ${selectedEmail.value.subject}\n时间: ${formatTime(selectedEmail.value.sent_at)}\n\n${selectedEmail.value.content}`
      }
    })
  }
}

// 下载附件
const downloadAttachment = async (attachment) => {
  try {
    const response = await userApi.downloadAttachment(attachment.id)

    // 创建下载链接
    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', attachment.filename)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
  } catch (error) {
    console.error('下载附件失败:', error)
    alert('下载附件失败')
  }
}

// 写邮件
const composeEmail = () => {
  router.push('/user/compose')
}

// 刷新邮件
const refreshEmails = async () => {
  await loadSentEmails()
}

// 获取状态图标
const getStatusIcon = (status) => {
  const icons = {
    sent: 'fas fa-paper-plane',
    delivered: 'fas fa-check',
    read: 'fas fa-check-double',
    failed: 'fas fa-exclamation-triangle',
    pending: 'fas fa-clock',
    bounced: 'fas fa-exclamation-circle'
  }
  return icons[status] || 'fas fa-question'
}

// 获取状态文本
const getStatusText = (status) => {
  const texts = {
    sent: '已发送',
    delivered: '已送达',
    read: '已读',
    failed: '发送失败',
    pending: '发送中',
    bounced: '退回'
  }
  return texts[status] || '未知'
}

// 生命周期
onMounted(() => {
  loadSentEmails()
})
</script>

<style scoped>
.outbox-page {
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
}

.btn-primary {
  background: var(--admin-primary);
  color: white;
  box-shadow: 0 4px 12px rgba(37, 99, 235, 0.3);
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(37, 99, 235, 0.4);
}

.btn-secondary {
  background: rgba(255, 255, 255, 0.1);
  color: var(--admin-light);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.btn-secondary:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.15);
  transform: translateY(-2px);
  border-color: rgba(255, 255, 255, 0.3);
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.outbox-content {
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
  grid-template-columns: 40px 200px 1fr 200px 120px 120px;
  gap: 20px;
  padding: 20px;
  border-bottom: 1px solid var(--border);
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

.email-checkbox input {
  width: 16px;
  height: 16px;
  cursor: pointer;
}

.email-recipient {
  display: flex;
  align-items: center;
  gap: 12px;
}

.recipient-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--success), #059669);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 14px;
}

.recipient-name {
  font-weight: 500;
  color: var(--text-primary);
  font-size: 14px;
}

.email-subject {
  font-weight: 600;
  color: var(--text-primary);
  font-size: 16px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.email-preview {
  color: var(--text-secondary);
  font-size: 14px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.email-status {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  font-weight: 500;
  padding: 4px 8px;
  border-radius: 6px;
}

.email-status.sent {
  background: rgba(59, 130, 246, 0.2);
  color: var(--info);
}

.email-status.delivered {
  background: rgba(16, 185, 129, 0.2);
  color: var(--success);
}

.email-status.read {
  background: rgba(16, 185, 129, 0.3);
  color: var(--success);
}

.email-status.failed {
  background: rgba(239, 68, 68, 0.2);
  color: var(--accent);
}

.email-time {
  color: var(--text-secondary);
  font-size: 12px;
  text-align: right;
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

/* 空状态 */
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
}

.modal-content {
  background: linear-gradient(135deg, rgba(15, 23, 42, 0.95), rgba(30, 41, 59, 0.95));
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  width: 90%;
  max-width: 800px;
  max-height: 80vh;
  overflow-y: auto;
  backdrop-filter: blur(10px);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 25px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.modal-header h3 {
  color: var(--admin-light);
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.close-btn {
  background: none;
  border: none;
  color: var(--admin-gray);
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
  padding: 25px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 15px;
  padding: 20px 25px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

/* 邮件元信息 */
.email-meta-info {
  margin-bottom: 25px;
}

.meta-item {
  display: flex;
  margin-bottom: 12px;
  font-size: 14px;
}

.meta-item label {
  min-width: 80px;
  color: var(--admin-gray);
  font-weight: 500;
}

.meta-item span {
  color: var(--admin-light);
}

.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.status-badge.sent {
  background: rgba(59, 130, 246, 0.2);
  color: #60a5fa;
}

.status-badge.delivered {
  background: rgba(34, 197, 94, 0.2);
  color: #4ade80;
}

.status-badge.read {
  background: rgba(16, 185, 129, 0.2);
  color: #10b981;
}

.status-badge.failed {
  background: rgba(239, 68, 68, 0.2);
  color: #f87171;
}

.status-badge.pending {
  background: rgba(245, 158, 11, 0.2);
  color: #fbbf24;
}

/* 邮件内容 */
.email-content-detail h4 {
  color: var(--admin-light);
  margin-bottom: 15px;
  font-size: 16px;
  font-weight: 600;
}

.content-body {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  padding: 20px;
  color: var(--admin-light);
  line-height: 1.6;
  max-height: 300px;
  overflow-y: auto;
}

/* 附件样式 */
.attachments-section {
  margin-top: 25px;
}

.attachments-section h4 {
  color: var(--admin-light);
  margin-bottom: 15px;
  font-size: 16px;
  font-weight: 600;
}

.attachment-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.attachment-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  transition: all 0.3s ease;
}

.attachment-item:hover {
  background: rgba(255, 255, 255, 0.08);
}

.attachment-icon {
  width: 32px;
  height: 32px;
  background: rgba(59, 130, 246, 0.2);
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #60a5fa;
}

.attachment-info {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.attachment-name {
  color: var(--admin-light);
  font-weight: 500;
  font-size: 14px;
}

.attachment-size {
  color: var(--admin-gray);
  font-size: 12px;
}

.download-btn {
  background: rgba(34, 197, 94, 0.1);
  border: 1px solid rgba(34, 197, 94, 0.3);
  color: #4ade80;
  padding: 6px 8px;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.download-btn:hover {
  background: rgba(34, 197, 94, 0.2);
  transform: translateY(-1px);
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

  .email-recipient,
  .email-subject,
  .email-preview,
  .email-status,
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
