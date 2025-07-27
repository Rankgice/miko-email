import axios from 'axios'

// 创建axios实例
const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
  withCredentials: true // 重要：确保发送cookies用于session认证
})

// 请求拦截器
api.interceptors.request.use(
  config => {
    // 基于session的认证不需要添加token，cookies会自动发送
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  response => {
    return response
  },
  error => {
    if (error.response) {
      switch (error.response.status) {
        case 401:
          // 未授权，清除本地存储并跳转到登录页
          localStorage.removeItem('user_token')
          localStorage.removeItem('admin_token')
          localStorage.removeItem('user_info')
          localStorage.removeItem('admin_info')
          
          // 根据当前路径判断跳转到哪个登录页
          const currentPath = window.location.pathname
          if (currentPath.startsWith('/admin')) {
            window.location.href = '/admin/login'
          } else {
            window.location.href = '/user/login'
          }
          break
        case 403:
          console.error('权限不足')
          break
        case 404:
          console.error('请求的资源不存在')
          break
        case 500:
          console.error('服务器内部错误')
          break
        default:
          console.error('请求失败:', error.response.data?.msg || error.message)
      }
    } else if (error.request) {
      console.error('网络错误，请检查网络连接')
    } else {
      console.error('请求配置错误:', error.message)
    }
    
    return Promise.reject(error)
  }
)

export default api
