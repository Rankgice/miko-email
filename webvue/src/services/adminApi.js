import api from '@/utils/api'

// 管理员API服务
export const adminApi = {
  // 仪表盘相关
  getDashboardStats() {
    return api.get('/admin/dashboard/stats')
  },

  getSystemHealth() {
    return api.get('/admin/system/health')
  },

  getRecentActivities() {
    return api.get('/admin/activities/recent')
  },

  // 用户管理
  getAllUsers() {
    return api.get('/admin/users')
  },

  getUserById(userId) {
    return api.get(`/admin/users/${userId}`)
  },

  createUser(userData) {
    return api.post('/admin/users', userData)
  },

  updateUser(userId, userData) {
    return api.put(`/admin/users/${userId}`, userData)
  },

  updateUserStatus(userId, status) {
    return api.put(`/admin/users/${userId}/status`, { status })
  },

  deleteUser(userId) {
    return api.delete(`/admin/users/${userId}`)
  },

  // 邮箱管理
  getAllMailboxes() {
    return api.get('/admin/mailboxes')
  },

  getMailboxById(mailboxId) {
    return api.get(`/admin/mailboxes/${mailboxId}`)
  },

  createMailbox(mailboxData) {
    return api.post('/admin/mailboxes', mailboxData)
  },

  updateMailbox(mailboxId, mailboxData) {
    return api.put(`/admin/mailboxes/${mailboxId}`, mailboxData)
  },

  deleteMailbox(mailboxId) {
    return api.delete(`/admin/mailboxes/${mailboxId}`)
  },

  // 域名管理
  getAllDomains() {
    return api.get('/admin/domains')
  },

  getDomainById(domainId) {
    return api.get(`/admin/domains/${domainId}`)
  },

  createDomain(domainData) {
    return api.post('/admin/domains', domainData)
  },

  updateDomain(domainId, domainData) {
    return api.put(`/admin/domains/${domainId}`, domainData)
  },

  deleteDomain(domainId) {
    return api.delete(`/admin/domains/${domainId}`)
  },

  verifyDomain(domainId) {
    return api.post(`/admin/domains/${domainId}/verify`)
  },

  getDomainDNSRecords(domainId) {
    return api.get(`/admin/domains/${domainId}/dns-records`)
  },

  verifySenderConfiguration(domainId) {
    return api.post(`/admin/domains/${domainId}/verify-sender`)
  },

  verifyReceiverConfiguration(domainId) {
    return api.post(`/admin/domains/${domainId}/verify-receiver`)
  },

  getDomainDKIMRecord(domain) {
    return api.get(`/admin/domains/dkim-record?domain=${domain}`)
  },

  getServerInfo() {
    return api.get('/admin/domains/server-info')
  },

  // 系统管理相关API
  getSystemLogs(params = {}) {
    return api.get('/admin/system/logs', { params })
  },

  getSystemSettings() {
    return api.get('/admin/system/settings')
  },

  updateSystemSettings(settings) {
    return api.put('/admin/system/settings', settings)
  },

  getSystemStatus() {
    return api.get('/admin/system/status')
  },

  // 域名状态管理
  updateDomainStatus(domainId, data) {
    return api.put(`/admin/domains/${domainId}/status`, data)
  },

  // 验证码管理
  getCaptchaCodes() {
    return api.get('/admin/captcha/codes')
  },

  getCaptchaStats() {
    return api.get('/admin/captcha/stats')
  },

  generateTestCaptcha() {
    return api.post('/admin/captcha/generate-test')
  },

  deleteCaptchaCode(codeId) {
    return api.delete(`/admin/captcha/codes/${codeId}`)
  },

  clearExpiredCaptcha() {
    return api.delete('/admin/captcha/expired')
  },

  resendCaptcha(codeId) {
    return api.post(`/admin/captcha/resend/${codeId}`)
  },

  // 验证码规则管理
  getCaptchaRules() {
    return api.get('/admin/captcha/rules')
  },

  updateCaptchaRuleStatus(ruleId, enabled) {
    return api.put(`/admin/captcha/rules/${ruleId}/status`, { enabled })
  },

  createCaptchaRule(ruleData) {
    return api.post('/admin/captcha/rules', ruleData)
  },

  updateCaptchaRule(ruleId, ruleData) {
    return api.put(`/admin/captcha/rules/${ruleId}`, ruleData)
  },

  deleteCaptchaRule(ruleId) {
    return api.delete(`/admin/captcha/rules/${ruleId}`)
  },

  testCaptchaRule(ruleId, testData) {
    return api.post(`/admin/captcha/rules/${ruleId}/test`, testData)
  },

  // 系统设置
  getAdminSettings() {
    return api.get('/admin/settings')
  },

  updateAdminSettings(settings) {
    return api.put('/admin/settings', settings)
  },

  // 系统设置管理
  getSystemSettings() {
    return api.get('/admin/system/settings')
  },

  updateSystemSettings(settings) {
    return api.put('/admin/system/settings', settings)
  },

  // 系统状态
  getSystemStatus() {
    return api.get('/admin/system/status')
  },

  // 系统日志
  getSystemLogs(params) {
    return api.get('/admin/system/logs', { params })
  },

  // 管理员通知
  getAdminNotifications() {
    return api.get('/admin/notifications')
  }
}

export default adminApi
