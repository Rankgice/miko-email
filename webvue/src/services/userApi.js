import api from '@/utils/api'

/**
 * 用户端API服务
 */
class UserApiService {
  // 仪表盘相关
  getDashboardStats() {
    return api.get('/dashboard/stats')
  }

  getRecentEmails(limit = 5) {
    return api.get(`/emails/recent?limit=${limit}`)
  }

  // 邮件相关
  getEmails(params = {}) {
    return api.get('/emails', { params })
  }

  getEmail(id) {
    return api.get(`/emails/${id}`)
  }

  sendEmail(emailData) {
    return api.post('/emails/send', emailData)
  }

  saveDraft(emailData) {
    return api.post('/emails/draft', emailData)
  }

  markAsRead(emailId) {
    return api.put(`/emails/${emailId}/read`)
  }

  markAsUnread(emailId) {
    return api.put(`/emails/${emailId}/unread`)
  }

  deleteEmail(emailId) {
    return api.delete(`/emails/${emailId}`)
  }

  // 邮箱管理
  getMailboxes() {
    return api.get('/mailboxes')
  }

  createMailbox(mailboxData) {
    return api.post('/mailboxes', mailboxData)
  }

  updateMailbox(mailboxId, mailboxData) {
    return api.put(`/mailboxes/${mailboxId}`, mailboxData)
  }

  deleteMailbox(mailboxId) {
    return api.delete(`/mailboxes/${mailboxId}`)
  }

  toggleMailboxStatus(mailboxId, isActive) {
    return api.put(`/mailboxes/${mailboxId}/status`, { is_active: isActive })
  }

  // 域名管理
  getDomains() {
    return api.get('/domains/available')
  }

  // 注意：普通用户没有域名管理权限，这些方法仅用于兼容性
  createDomain(domainData) {
    return Promise.reject(new Error('普通用户无权限创建域名，请联系管理员'))
  }

  updateDomain(domainId, domainData) {
    return Promise.reject(new Error('普通用户无权限修改域名，请联系管理员'))
  }

  deleteDomain(domainId) {
    return Promise.reject(new Error('普通用户无权限删除域名，请联系管理员'))
  }

  verifyDomain(domainId) {
    return Promise.reject(new Error('普通用户无权限验证域名，请联系管理员'))
  }

  getDomainDNSRecords(domainName) {
    return Promise.reject(new Error('普通用户无权限查看DNS记录，请联系管理员'))
  }

  // 转发规则
  getForwardRules() {
    return api.get('/forward-rules')
  }

  createForwardRule(ruleData) {
    return api.post('/forward-rules', ruleData)
  }

  getForwardRule(ruleId) {
    return api.get(`/forward-rules/${ruleId}`)
  }

  updateForwardRule(ruleId, ruleData) {
    return api.put(`/forward-rules/${ruleId}`, ruleData)
  }

  deleteForwardRule(ruleId) {
    return api.delete(`/forward-rules/${ruleId}`)
  }

  toggleForwardRule(ruleId, enabled) {
    return api.put(`/forward-rules/${ruleId}/toggle`, { enabled })
  }

  testForwardRule(ruleId, testData) {
    return api.post(`/forward-rules/${ruleId}/test`, testData)
  }

  // 用户设置
  getUserSettings() {
    return api.get('/user/settings')
  }

  updateUserProfile(profileData) {
    return api.put('/user/profile', profileData)
  }

  changePassword(passwordData) {
    return api.put('/user/password', passwordData)
  }

  updateNotifications(notificationData) {
    return api.put('/user/notifications', notificationData)
  }

  updateTheme(themeData) {
    return api.put('/user/theme', themeData)
  }

  // 通知相关
  getNotifications() {
    return api.get('/user/notifications')
  }

  markNotificationAsRead(notificationId) {
    return api.put(`/user/notifications/${notificationId}/read`)
  }

  markAllNotificationsAsRead() {
    return api.put('/user/notifications/read-all')
  }

  deleteNotification(notificationId) {
    return api.delete(`/user/notifications/${notificationId}`)
  }

  // 文件上传
  uploadAttachment(formData, config = {}) {
    return api.post('/upload/attachment', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      },
      ...config
    })
  }

  // 下载附件
  downloadAttachment(attachmentId) {
    return api.get(`/attachments/${attachmentId}/download`, {
      responseType: 'blob'
    })
  }

  // 搜索
  searchEmails(query, params = {}) {
    return api.post('/emails/search', {
      query: query,
      ...params
    })
  }

  // 统计信息
  getEmailStats() {
    return api.get('/stats/emails')
  }

  getStorageStats() {
    return api.get('/stats/storage')
  }
}

// 创建实例
const userApi = new UserApiService()

export default userApi
