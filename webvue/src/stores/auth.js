import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/utils/api'

export const useAuthStore = defineStore('auth', () => {
  // 用户状态
  const userToken = ref(localStorage.getItem('user_token') || '')
  const userInfo = ref(JSON.parse(localStorage.getItem('user_info') || '{}'))
  const isUserLoggedIn = ref(!!userToken.value)

  // 管理员状态
  const adminToken = ref(localStorage.getItem('admin_token') || '')
  const adminInfo = ref(JSON.parse(localStorage.getItem('admin_info') || '{}'))
  const isAdminLoggedIn = ref(!!adminToken.value)

  // 用户登录
  const userLogin = async (credentials) => {
    try {
      const response = await api.post('/login', credentials)

      // 适配统一的 code: 0 格式
      if (response.data.code === 0) {
        const { user } = response.data.data

        // 基于session的认证，不需要token
        userToken.value = 'session-based'
        userInfo.value = user
        isUserLoggedIn.value = true

        localStorage.setItem('user_token', 'session-based')
        localStorage.setItem('user_info', JSON.stringify(user))

        return { success: true }
      } else {
        return { success: false, message: response.data.msg || response.data.message || '登录失败' }
      }
    } catch (error) {
      console.error('用户登录错误:', error)
      return {
        success: false,
        message: error.response?.data?.msg || error.response?.data?.message || error.message || '登录失败'
      }
    }
  }

  // 管理员登录
  const adminLogin = async (credentials) => {
    try {
      const response = await api.post('/admin/login', credentials)

      // 适配统一的 code: 0 格式
      if (response.data.code === 0) {
        // 后端返回的用户信息可能在不同的位置，先检查
        const user = response.data.data?.user || { username: credentials.username }

        adminToken.value = 'admin-logged-in'
        adminInfo.value = user
        isAdminLoggedIn.value = true

        localStorage.setItem('admin_token', 'admin-logged-in')
        localStorage.setItem('admin_info', JSON.stringify(user))

        return { success: true }
      } else {
        return { success: false, message: response.data.msg || response.data.message || '登录失败' }
      }
    } catch (error) {
      console.error('管理员登录错误:', error)
      return {
        success: false,
        message: error.response?.data?.msg || error.response?.data?.message || error.message || '登录失败'
      }
    }
  }

  // 用户注册
  const userRegister = async (userData) => {
    try {
      const response = await api.post('/register', userData)
      if (response.data.code === 0) {
        return { success: true }
      } else {
        return { success: false, message: response.data.msg }
      }
    } catch (error) {
      return {
        success: false,
        message: error.response?.data?.msg || '注册失败'
      }
    }
  }

  // 用户登出
  const userLogout = async () => {
    try {
      await api.post('/logout')
    } catch (error) {
      console.error('登出请求失败:', error)
    } finally {
      userToken.value = ''
      userInfo.value = {}
      isUserLoggedIn.value = false

      localStorage.removeItem('user_token')
      localStorage.removeItem('user_info')
    }
  }

  // 管理员登出
  const adminLogout = async () => {
    try {
      await api.post('/admin/logout')
    } catch (error) {
      console.error('登出请求失败:', error)
    } finally {
      adminToken.value = ''
      adminInfo.value = {}
      isAdminLoggedIn.value = false

      localStorage.removeItem('admin_token')
      localStorage.removeItem('admin_info')
    }
  }

  return {
    // 状态
    userToken,
    userInfo,
    isUserLoggedIn,
    adminToken,
    adminInfo,
    isAdminLoggedIn,

    // 方法
    userLogin,
    adminLogin,
    userRegister,
    userLogout,
    adminLogout
  }
})
