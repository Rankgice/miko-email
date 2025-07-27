<template>
  <div class="settings-page">
    <div class="page-header">
      <h1 class="page-title">账户设置</h1>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-container">
      <div class="loading-spinner"></div>
      <p>加载用户设置中...</p>
    </div>

    <div v-else class="settings-content">
      <!-- 个人信息 -->
      <div class="settings-section">
        <h2 class="section-title">个人信息</h2>
        <div class="settings-card">
          <div class="form-group">
            <label class="form-label">用户名</label>
            <input type="text" class="form-input" v-model="userSettings.username" readonly>
          </div>

          <div class="form-group">
            <label class="form-label">邮箱地址</label>
            <input type="email" class="form-input" v-model="userSettings.email">
          </div>

          <div class="form-group">
            <label class="form-label">显示名称</label>
            <input type="text" class="form-input" v-model="userSettings.displayName" placeholder="请输入显示名称">
          </div>

          <div class="form-actions">
            <button class="btn btn-primary" @click="saveProfile" :disabled="saving">
              <i class="fas fa-save" v-if="!saving"></i>
              <i class="fas fa-spinner fa-spin" v-if="saving"></i>
              {{ saving ? '保存中...' : '保存更改' }}
            </button>
          </div>
        </div>
      </div>

      <!-- 密码设置 -->
      <div class="settings-section">
        <h2 class="section-title">密码设置</h2>
        <div class="settings-card">
          <div class="form-group">
            <label class="form-label">当前密码</label>
            <input type="password" class="form-input" v-model="passwordForm.currentPassword" placeholder="请输入当前密码">
          </div>

          <div class="form-group">
            <label class="form-label">新密码</label>
            <input type="password" class="form-input" v-model="passwordForm.newPassword" placeholder="请输入新密码">
          </div>

          <div class="form-group">
            <label class="form-label">确认新密码</label>
            <input type="password" class="form-input" v-model="passwordForm.confirmPassword" placeholder="请再次输入新密码">
          </div>

          <div class="form-actions">
            <button class="btn btn-warning" @click="changePassword" :disabled="saving">
              <i class="fas fa-key" v-if="!saving"></i>
              <i class="fas fa-spinner fa-spin" v-if="saving"></i>
              {{ saving ? '修改中...' : '修改密码' }}
            </button>
          </div>
        </div>
      </div>

      <!-- 邮件设置 -->
      <div class="settings-section">
        <h2 class="section-title">邮件设置</h2>
        <div class="settings-card">
          <div class="form-group">
            <label class="form-label">邮件签名</label>
            <textarea class="form-textarea" v-model="userSettings.signature" placeholder="请输入邮件签名" rows="4"></textarea>
          </div>

          <div class="form-group">
            <div class="checkbox-group">
              <label class="checkbox-label">
                <input type="checkbox" v-model="userSettings.autoReply">
                <span class="checkmark"></span>
                启用自动回复
              </label>
            </div>
          </div>

          <div class="form-group" v-if="userSettings.autoReply">
            <label class="form-label">自动回复内容</label>
            <textarea class="form-textarea" v-model="userSettings.autoReplyMessage" placeholder="请输入自动回复内容" rows="3"></textarea>
          </div>

          <div class="form-actions">
            <button class="btn btn-primary" @click="saveEmailSettings" :disabled="saving">
              <i class="fas fa-save" v-if="!saving"></i>
              <i class="fas fa-spinner fa-spin" v-if="saving"></i>
              {{ saving ? '保存中...' : '保存设置' }}
            </button>
          </div>
        </div>
      </div>

      <!-- 通知设置 -->
      <div class="settings-section">
        <h2 class="section-title">通知设置</h2>
        <div class="settings-card">
          <div class="form-group">
            <div class="checkbox-group">
              <label class="checkbox-label">
                <input type="checkbox" v-model="userSettings.emailNotifications">
                <span class="checkmark"></span>
                邮件通知
              </label>
            </div>
          </div>

          <div class="form-group">
            <div class="checkbox-group">
              <label class="checkbox-label">
                <input type="checkbox" v-model="userSettings.browserNotifications">
                <span class="checkmark"></span>
                浏览器通知
              </label>
            </div>
          </div>

          <div class="form-actions">
            <button class="btn btn-primary" @click="saveNotificationSettings" :disabled="saving">
              <i class="fas fa-save" v-if="!saving"></i>
              <i class="fas fa-spinner fa-spin" v-if="saving"></i>
              {{ saving ? '保存中...' : '保存设置' }}
            </button>
          </div>
        </div>
      </div>

      <!-- 界面设置 -->
      <div class="settings-section">
        <h2 class="section-title">界面设置</h2>
        <div class="settings-card">
          <div class="form-group">
            <label class="form-label">主题模式</label>
            <select class="form-input" v-model="userSettings.theme">
              <option value="dark">深色主题</option>
              <option value="light">浅色主题</option>
              <option value="auto">跟随系统</option>
            </select>
          </div>

          <div class="form-group">
            <div class="checkbox-group">
              <label class="checkbox-label">
                <input type="checkbox" v-model="userSettings.compactMode">
                <span class="checkmark"></span>
                紧凑模式
              </label>
            </div>
          </div>

          <div class="form-group">
            <div class="checkbox-group">
              <label class="checkbox-label">
                <input type="checkbox" v-model="userSettings.animations">
                <span class="checkmark"></span>
                启用动画效果
              </label>
            </div>
          </div>

          <div class="form-actions">
            <button class="btn btn-primary" @click="saveThemeSettings" :disabled="saving">
              <i class="fas fa-save" v-if="!saving"></i>
              <i class="fas fa-spinner fa-spin" v-if="saving"></i>
              {{ saving ? '保存中...' : '保存设置' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import userApi from '@/services/userApi'

const authStore = useAuthStore()

// 响应式数据
const userSettings = ref({
  username: '',
  email: '',
  displayName: '',
  signature: '',
  autoReply: false,
  autoReplyMessage: '',
  emailNotifications: true,
  browserNotifications: false,
  theme: 'dark',
  compactMode: false,
  animations: true
})

const passwordForm = ref({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const loading = ref(false)
const saving = ref(false)

// 加载用户设置
const loadUserSettings = async () => {
  try {
    loading.value = true
    const response = await userApi.getUserSettings()

    if (response.data.code === 0) {
      const settings = response.data.data
      userSettings.value = {
        username: settings.profile.username || '',
        email: settings.profile.email || '',
        displayName: settings.profile.displayName || settings.profile.username || '',
        signature: settings.profile.signature || '',
        autoReply: false, // 暂时硬编码，后续可以从API获取
        autoReplyMessage: '',
        emailNotifications: settings.notifications.newEmail || true,
        browserNotifications: settings.notifications.security || false,
        theme: settings.theme || 'dark',
        compactMode: settings.interface?.compact || false,
        animations: settings.interface?.animations !== false
      }
    } else {
      console.error('获取用户设置失败:', response.data.msg)
      alert('获取用户设置失败: ' + response.data.msg)
    }
  } catch (error) {
    console.error('加载用户设置失败:', error)
    alert('加载用户设置失败: ' + (error.response?.data?.msg || error.message))
  } finally {
    loading.value = false
  }
}

// 保存个人信息
const saveProfile = async () => {
  try {
    saving.value = true
    const profileData = {
      email: userSettings.value.email
    }

    const response = await userApi.updateUserProfile(profileData)

    if (response.data.code === 0) {
      alert('个人信息保存成功')
    } else {
      alert('保存失败: ' + (response.data.msg || '未知错误'))
    }
  } catch (error) {
    console.error('保存个人信息失败:', error)
    alert('保存失败: ' + (error.response?.data?.msg || error.message))
  } finally {
    saving.value = false
  }
}

const changePassword = async () => {
  if (!passwordForm.value.currentPassword || !passwordForm.value.newPassword) {
    alert('请填写所有密码字段')
    return
  }

  if (passwordForm.value.newPassword !== passwordForm.value.confirmPassword) {
    alert('新密码和确认密码不一致')
    return
  }

  if (passwordForm.value.newPassword.length < 6) {
    alert('新密码长度至少6位')
    return
  }

  try {
    saving.value = true
    const passwordData = {
      currentPassword: passwordForm.value.currentPassword,
      newPassword: passwordForm.value.newPassword
    }

    const response = await userApi.changePassword(passwordData)

    if (response.data.code === 0) {
      alert('密码修改成功')
      // 清空表单
      passwordForm.value = {
        currentPassword: '',
        newPassword: '',
        confirmPassword: ''
      }
    } else {
      alert('修改密码失败: ' + (response.data.msg || '未知错误'))
    }
  } catch (error) {
    console.error('修改密码失败:', error)
    alert('修改密码失败: ' + (error.response?.data?.msg || error.message))
  } finally {
    saving.value = false
  }
}

const saveEmailSettings = async () => {
  try {
    saving.value = true
    // 邮件设置暂时只保存签名，自动回复功能需要后端支持
    const profileData = {
      email: userSettings.value.email
    }

    const response = await userApi.updateUserProfile(profileData)

    if (response.data.code === 0) {
      alert('邮件设置保存成功')
    } else {
      alert('保存失败: ' + (response.data.msg || '未知错误'))
    }
  } catch (error) {
    console.error('保存邮件设置失败:', error)
    alert('保存失败: ' + (error.response?.data?.msg || error.message))
  } finally {
    saving.value = false
  }
}

const saveNotificationSettings = async () => {
  try {
    saving.value = true
    const notificationData = {
      newEmail: userSettings.value.emailNotifications,
      emailSent: true,
      forwardRule: true,
      security: userSettings.value.browserNotifications,
      maintenance: false
    }

    const response = await userApi.updateNotifications(notificationData)

    if (response.data.code === 0) {
      alert('通知设置保存成功')
    } else {
      alert('保存失败: ' + (response.data.msg || '未知错误'))
    }
  } catch (error) {
    console.error('保存通知设置失败:', error)
    alert('保存失败: ' + (error.response?.data?.msg || error.message))
  } finally {
    saving.value = false
  }
}

const saveThemeSettings = async () => {
  try {
    saving.value = true
    const themeData = {
      theme: userSettings.value.theme,
      interface: {
        compact: userSettings.value.compactMode,
        animations: userSettings.value.animations,
        showAvatars: true
      }
    }

    const response = await userApi.updateTheme(themeData)

    if (response.data.code === 0) {
      alert('主题设置保存成功')
    } else {
      alert('保存失败: ' + (response.data.msg || '未知错误'))
    }
  } catch (error) {
    console.error('保存主题设置失败:', error)
    alert('保存失败: ' + (error.response?.data?.msg || error.message))
  } finally {
    saving.value = false
  }
}

// 页面加载时获取用户设置
onMounted(() => {
  loadUserSettings()
})
</script>

<style scoped>
.settings-page {
  max-width: 800px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 30px;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0;
}

.settings-content {
  display: flex;
  flex-direction: column;
  gap: 30px;
}

.settings-section {
  background: linear-gradient(135deg, rgba(30, 41, 59, 0.9), rgba(15, 23, 42, 0.95));
  border-radius: 12px;
  border: 1px solid var(--border);
  overflow: hidden;
}

.section-title {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
  padding: 20px 25px;
  margin: 0;
  background: rgba(15, 23, 42, 0.5);
  border-bottom: 1px solid var(--border);
}

.settings-card {
  padding: 25px;
}

.form-group {
  margin-bottom: 20px;
}

.form-label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: var(--text-primary);
  font-size: 14px;
}

.form-input,
.form-textarea {
  width: 100%;
  padding: 12px 16px;
  border-radius: 8px;
  border: 1px solid var(--border);
  background-color: rgba(15, 23, 42, 0.5);
  color: var(--text-primary);
  font-size: 14px;
  transition: all 0.3s ease;
  resize: vertical;
}

.form-input:focus,
.form-textarea:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(0, 180, 216, 0.2);
}

.form-input[readonly] {
  background-color: rgba(15, 23, 42, 0.3);
  color: var(--text-secondary);
  cursor: not-allowed;
}

.form-textarea {
  min-height: 80px;
  font-family: inherit;
}

.checkbox-group {
  display: flex;
  align-items: center;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  color: var(--text-primary);
  font-size: 14px;
}

.checkbox-label input[type="checkbox"] {
  display: none;
}

.checkmark {
  width: 20px;
  height: 20px;
  border: 2px solid var(--border);
  border-radius: 4px;
  position: relative;
  transition: all 0.3s ease;
}

.checkbox-label input[type="checkbox"]:checked + .checkmark {
  background: var(--primary);
  border-color: var(--primary);
}

.checkbox-label input[type="checkbox"]:checked + .checkmark::after {
  content: '\f00c';
  font-family: 'Font Awesome 6 Free';
  font-weight: 900;
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  color: white;
  font-size: 12px;
}

.form-actions {
  margin-top: 25px;
  padding-top: 20px;
  border-top: 1px solid var(--border);
}

.btn {
  padding: 12px 24px;
  border-radius: 8px;
  border: none;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
}

.btn-primary {
  background: linear-gradient(135deg, var(--primary), #0077b6);
  color: white;
  box-shadow: 0 5px 15px rgba(0, 180, 216, 0.2);
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(0, 180, 216, 0.3);
}

.btn-warning {
  background: linear-gradient(135deg, var(--warning), #d97706);
  color: white;
  box-shadow: 0 5px 15px rgba(245, 158, 11, 0.2);
}

.btn-warning:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(245, 158, 11, 0.3);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .settings-page {
    padding: 0 15px;
  }

  .settings-card {
    padding: 20px;
  }

  .section-title {
    padding: 15px 20px;
    font-size: 18px;
  }

  .page-title {
    font-size: 24px;
  }
}

@media (max-width: 480px) {
  .settings-card {
    padding: 15px;
  }

  .form-input,
  .form-textarea {
    padding: 10px 12px;
  }

  .btn {
    width: 100%;
    justify-content: center;
  }
}

/* 加载状态样式 */
.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: var(--text-secondary);
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid rgba(59, 130, 246, 0.3);
  border-top: 3px solid #3b82f6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 16px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* 按钮禁用状态 */
.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  pointer-events: none;
}

.btn .fa-spinner {
  animation: spin 1s linear infinite;
}
</style>
