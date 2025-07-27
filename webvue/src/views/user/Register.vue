<template>
  <div class="register-page">
    <!-- 装饰元素 -->
    <div class="decoration circle1"></div>
    <div class="decoration circle2"></div>

    <!-- 浮动图标 -->
    <i class="fas fa-user-plus floating-icon floating-1"></i>
    <i class="fas fa-envelope floating-icon floating-2"></i>
    <i class="fas fa-shield-alt floating-icon floating-3"></i>

    <!-- 顶部导航 -->
    <nav class="top-nav">
      <div class="logo">
        <i class="fas fa-envelope"></i>
        <span>Miko邮箱</span>
      </div>

      <div class="nav-controls">
        <router-link to="/user/login" class="login-link">
          <i class="fas fa-sign-in-alt"></i>
          <span>用户登录</span>
        </router-link>
      </div>
    </nav>

    <!-- 注册容器 -->
    <div class="register-container">
      <div class="register-header">
        <h1 class="register-title">用户注册</h1>
        <p class="register-subtitle">创建您的Miko邮箱账户</p>
      </div>

      <form class="register-form" @submit.prevent="handleRegister">
        <div class="input-group">
          <label class="input-label">
            <i class="fas fa-user input-icon"></i>
            用户名
          </label>
          <div class="input-icon-container">
            <i class="fas fa-user"></i>
          </div>
          <input
            type="text"
            class="form-input"
            placeholder="请输入用户名（3-20个字符）"
            v-model="registerForm.username"
            required
          >
        </div>

        <div class="input-group">
          <label class="input-label">
            <i class="fas fa-envelope input-icon"></i>
            邮箱地址
          </label>
          <div class="input-icon-container">
            <i class="fas fa-envelope"></i>
          </div>
          <input
            type="email"
            class="form-input"
            placeholder="请输入您的邮箱地址"
            v-model="registerForm.email"
            required
          >
        </div>

        <div class="input-group">
          <label class="input-label">
            <i class="fas fa-lock input-icon"></i>
            密码
          </label>
          <div class="input-icon-container">
            <i class="fas fa-lock"></i>
          </div>
          <div class="password-container">
            <input
              :type="showPassword ? 'text' : 'password'"
              class="form-input"
              placeholder="请输入密码（至少6位）"
              v-model="registerForm.password"
              required
            >
            <span class="toggle-password" @click="togglePassword">
              <i :class="showPassword ? 'fas fa-eye-slash' : 'fas fa-eye'"></i>
            </span>
          </div>
        </div>

        <div class="input-group">
          <label class="input-label">
            <i class="fas fa-lock input-icon"></i>
            确认密码
          </label>
          <div class="input-icon-container">
            <i class="fas fa-lock"></i>
          </div>
          <div class="password-container">
            <input
              :type="showConfirmPassword ? 'text' : 'password'"
              class="form-input"
              placeholder="请再次输入密码"
              v-model="registerForm.confirmPassword"
              required
            >
            <span class="toggle-password" @click="toggleConfirmPassword">
              <i :class="showConfirmPassword ? 'fas fa-eye-slash' : 'fas fa-eye'"></i>
            </span>
          </div>
        </div>

        <div class="input-group">
          <label class="input-label">
            <i class="fas fa-shield-alt input-icon"></i>
            验证码
          </label>
          <div class="captcha-container">
            <div class="input-icon-container">
              <i class="fas fa-shield-alt"></i>
            </div>
            <input
              type="text"
              class="form-input captcha-input"
              placeholder="请输入验证码"
              v-model="registerForm.captcha"
              required
            >
            <button type="button" class="captcha-btn" @click="sendCaptcha" :disabled="captchaLoading">
              {{ captchaLoading ? '发送中...' : '获取验证码' }}
            </button>
          </div>
        </div>

        <div class="register-options">
          <div class="agreement">
            <div
              class="agreement-checkbox"
              :class="{ checked: agreeTerms }"
              @click="agreeTerms = !agreeTerms"
            ></div>
            <span>我已阅读并同意 <a href="#" class="terms-link">用户协议</a> 和 <a href="#" class="privacy-link">隐私政策</a></span>
          </div>
        </div>

        <button type="submit" class="register-btn" :disabled="loading || !agreeTerms">
          <i class="fas fa-user-plus"></i>
          {{ loading ? '注册中...' : '立即注册' }}
        </button>
      </form>

      <div class="register-footer">
        <p>已有账户？<router-link to="/user/login" class="login-link">立即登录</router-link></p>
      </div>

      <div class="return-home">
        <router-link to="/" class="home-btn">
          <i class="fas fa-arrow-left"></i>
          返回首页
        </router-link>
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

// 响应式数据
const registerForm = ref({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  captcha: ''
})

const showPassword = ref(false)
const showConfirmPassword = ref(false)
const agreeTerms = ref(false)
const loading = ref(false)
const captchaLoading = ref(false)

// 方法
const togglePassword = () => {
  showPassword.value = !showPassword.value
}

const toggleConfirmPassword = () => {
  showConfirmPassword.value = !showConfirmPassword.value
}

const sendCaptcha = async () => {
  if (!registerForm.value.email) {
    alert('请先输入邮箱地址')
    return
  }

  captchaLoading.value = true

  try {
    // 这里调用发送验证码的API
    // await api.post('/api/captcha/send', { email: registerForm.value.email })
    alert('验证码已发送到您的邮箱')
  } catch (error) {
    alert('验证码发送失败，请重试')
  } finally {
    captchaLoading.value = false
  }
}

const handleRegister = async () => {
  // 表单验证
  if (!registerForm.value.username || !registerForm.value.email ||
      !registerForm.value.password || !registerForm.value.confirmPassword) {
    alert('请填写所有必填项')
    return
  }

  if (registerForm.value.password !== registerForm.value.confirmPassword) {
    alert('两次输入的密码不一致')
    return
  }

  if (registerForm.value.password.length < 6) {
    alert('密码长度至少6位')
    return
  }

  if (!agreeTerms.value) {
    alert('请先同意用户协议和隐私政策')
    return
  }

  loading.value = true

  try {
    const result = await authStore.userRegister(registerForm.value)

    if (result.success) {
      alert('注册成功！请登录您的账户')
      router.push('/user/login')
    } else {
      alert(result.message || '注册失败')
    }
  } catch (error) {
    console.error('注册失败:', error)
    alert('注册失败，请重试')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.register-page {
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
  color: var(--text-primary);
  min-height: 100vh;
  overflow-x: hidden;
  position: relative;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 20px;
}

/* 装饰元素 */
.decoration {
  position: fixed;
  z-index: -1;
  opacity: 0.05;
}

.circle1 {
  width: 400px;
  height: 400px;
  border-radius: 50%;
  background: radial-gradient(circle, var(--primary), transparent);
  top: -200px;
  right: -100px;
}

.circle2 {
  width: 300px;
  height: 300px;
  border-radius: 50%;
  background: radial-gradient(circle, var(--accent), transparent);
  bottom: -150px;
  left: -100px;
}

/* 顶部导航 */
.top-nav {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 30px;
  background-color: rgba(30, 41, 59, 0.7);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid var(--border);
  z-index: 100;
}

.logo {
  font-size: 24px;
  font-weight: 700;
  background: linear-gradient(to right, var(--primary), var(--primary-light));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo i {
  font-size: 28px;
}

.nav-controls {
  display: flex;
  align-items: center;
  gap: 25px;
}

.login-link {
  padding: 8px 15px;
  background-color: rgba(15, 23, 42, 0.5);
  border-radius: 8px;
  border: 1px solid var(--border);
  color: var(--text-primary);
  text-decoration: none;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 8px;
  transition: all 0.3s ease;
}

.login-link:hover {
  background-color: rgba(15, 23, 42, 0.8);
}

/* 注册容器 */
.register-container {
  width: 100%;
  max-width: 500px;
  background: linear-gradient(135deg, rgba(30, 41, 59, 0.9), rgba(15, 23, 42, 0.95));
  border-radius: 16px;
  padding: 40px;
  border: 1px solid var(--border);
  box-shadow: 0 15px 35px rgba(0, 0, 0, 0.25);
  position: relative;
  overflow: hidden;
  z-index: 10;
  margin-top: 80px;
}

.register-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 4px;
  background: linear-gradient(90deg, var(--primary), var(--accent));
}

.register-header {
  text-align: center;
  margin-bottom: 40px;
}

.register-title {
  font-size: 32px;
  font-weight: 700;
  margin-bottom: 10px;
  background: linear-gradient(to right, var(--primary), var(--primary-light));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.register-subtitle {
  font-size: 18px;
  color: var(--text-secondary);
}

.register-form {
  display: flex;
  flex-direction: column;
  gap: 25px;
}

.input-group {
  position: relative;
}

.input-label {
  display: block;
  margin-bottom: 10px;
  font-weight: 500;
  font-size: 16px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.input-icon {
  color: var(--primary);
}

.form-input {
  width: 100%;
  padding: 16px 20px 16px 50px;
  border-radius: 12px;
  border: 1px solid var(--border);
  background-color: rgba(15, 23, 42, 0.5);
  color: var(--text-primary);
  font-size: 16px;
  transition: all 0.3s ease;
}

.form-input:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 3px rgba(0, 180, 216, 0.2);
}

.input-icon-container {
  position: absolute;
  left: 20px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--text-secondary);
  font-size: 18px;
}

.password-container {
  position: relative;
}

.toggle-password {
  position: absolute;
  right: 20px;
  top: 50%;
  transform: translateY(-50%);
  color: var(--text-secondary);
  cursor: pointer;
  font-size: 18px;
}

.captcha-container {
  position: relative;
  display: flex;
  gap: 10px;
}

.captcha-input {
  flex: 1;
}

.captcha-btn {
  padding: 16px 20px;
  border-radius: 12px;
  border: 1px solid var(--primary);
  background: linear-gradient(135deg, var(--primary), #0077b6);
  color: white;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  white-space: nowrap;
}

.captcha-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(0, 180, 216, 0.3);
}

.captcha-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.register-options {
  margin-top: 10px;
}

.agreement {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  font-size: 14px;
  line-height: 1.5;
}

.agreement-checkbox {
  width: 18px;
  height: 18px;
  border-radius: 4px;
  border: 1px solid var(--border);
  background: rgba(15, 23, 42, 0.5);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  margin-top: 2px;
}

.agreement-checkbox.checked::after {
  content: '\f00c';
  font-family: 'Font Awesome 6 Free';
  font-weight: 900;
  font-size: 12px;
  color: var(--primary);
}

.terms-link,
.privacy-link {
  color: var(--primary);
  text-decoration: none;
  transition: all 0.3s ease;
}

.terms-link:hover,
.privacy-link:hover {
  text-decoration: underline;
}

.register-btn {
  padding: 16px;
  border-radius: 12px;
  border: none;
  font-weight: 600;
  font-size: 18px;
  cursor: pointer;
  transition: all 0.3s ease;
  background: linear-gradient(135deg, var(--primary), #0077b6);
  color: white;
  box-shadow: 0 5px 15px rgba(0, 180, 216, 0.2);
  margin-top: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.register-btn:hover:not(:disabled) {
  transform: translateY(-3px);
  box-shadow: 0 8px 20px rgba(0, 180, 216, 0.3);
}

.register-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.register-footer {
  text-align: center;
  margin-top: 30px;
  color: var(--text-secondary);
  font-size: 15px;
}

.register-footer .login-link {
  color: var(--primary);
  text-decoration: none;
  font-weight: 500;
  margin-left: 5px;
  transition: all 0.3s ease;
  padding: 0;
  background: none;
  border: none;
  border-radius: 0;
}

.register-footer .login-link:hover {
  text-decoration: underline;
  background: none;
}

.return-home {
  text-align: center;
  margin-top: 25px;
}

.home-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: var(--text-secondary);
  text-decoration: none;
  transition: all 0.3s ease;
}

.home-btn:hover {
  color: var(--primary);
}

/* 动画效果 */
@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.floating-icon {
  position: absolute;
  font-size: 120px;
  opacity: 0.1;
  z-index: -1;
  animation: float 6s ease-in-out infinite;
}

.floating-1 {
  top: 20%;
  left: 10%;
  color: var(--primary);
}

.floating-2 {
  bottom: 15%;
  right: 12%;
  color: var(--accent);
  animation-delay: 1s;
}

.floating-3 {
  top: 40%;
  right: 20%;
  color: var(--info);
  animation-delay: 2s;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .register-container {
    padding: 30px 20px;
    margin-top: 100px;
  }

  .top-nav {
    flex-direction: column;
    gap: 15px;
    padding: 15px;
  }

  .register-title {
    font-size: 28px;
  }

  .register-subtitle {
    font-size: 16px;
  }

  .captcha-container {
    flex-direction: column;
  }

  .captcha-input {
    margin-bottom: 10px;
  }
}

@media (max-width: 480px) {
  .register-container {
    padding: 25px 15px;
  }

  .register-title {
    font-size: 24px;
  }

  .form-input {
    padding: 14px 15px 14px 45px;
  }

  .input-icon-container {
    left: 15px;
  }
}
</style>
