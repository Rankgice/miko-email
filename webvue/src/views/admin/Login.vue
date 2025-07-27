<template>
  <div class="admin-login-page">
    <div class="login-container">
      <div class="login-header">
        <h1>管理员登录</h1>
        <p>Miko邮箱系统管理后台</p>
      </div>

      <form @submit.prevent="handleLogin">
        <div class="input-group">
          <label>管理员账号</label>
          <input
            type="text"
            placeholder="请输入管理员账号"
            v-model="loginForm.username"
            required
          />
        </div>

        <div class="input-group">
          <label>管理员密码</label>
          <input
            type="password"
            placeholder="请输入管理员密码"
            v-model="loginForm.password"
            required
          />
        </div>

        <button type="submit" :disabled="loading">
          {{ loading ? '登录中...' : '管理员登录' }}
        </button>
      </form>

      <div class="login-footer">
        <router-link to="/user/login">返回用户登录</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const loginForm = ref({
  username: 'kimi11',
  password: 'tgx1234561'
})

const loading = ref(false)

const handleLogin = async () => {
  loading.value = true
  try {
    console.log('开始管理员登录:', loginForm.value)
    const result = await authStore.adminLogin(loginForm.value)
    console.log('登录结果:', result)

    if (result.success) {
      console.log('登录成功，检查session状态')
      console.log('Cookies:', document.cookie)
      console.log('LocalStorage admin_token:', localStorage.getItem('admin_token'))

      alert('登录成功！')
      router.push('/admin/dashboard')
    } else {
      console.error('登录失败:', result.message)
      alert(result.message || '登录失败')
    }
  } catch (error) {
    console.error('登录异常:', error)
    alert('登录失败，请重试')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.admin-login-page {
  background: linear-gradient(135deg, #1e293b 0%, #0f172a 100%);
  color: white;
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 20px;
}

.login-container {
  background: rgba(15, 23, 42, 0.8);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  padding: 40px;
  width: 100%;
  max-width: 400px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.login-header h1 {
  font-size: 24px;
  margin-bottom: 8px;
  color: #3b82f6;
}

.login-header p {
  color: #94a3b8;
}

.input-group {
  margin-bottom: 20px;
}

.input-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
}

.input-group input {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.05);
  color: white;
  font-size: 16px;
}

.input-group input:focus {
  outline: none;
  border-color: #3b82f6;
}

.input-group input::placeholder {
  color: #94a3b8;
}

button {
  width: 100%;
  padding: 12px;
  background: linear-gradient(135deg, #3b82f6, #2563eb);
  border: none;
  border-radius: 8px;
  color: white;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

button:hover:not(:disabled) {
  transform: translateY(-2px);
}

button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.login-footer {
  text-align: center;
  margin-top: 20px;
}

.login-footer a {
  color: #94a3b8;
  text-decoration: none;
}

.login-footer a:hover {
  color: #3b82f6;
}
</style>
