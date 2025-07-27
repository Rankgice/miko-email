<template>
  <div class="api-test">
    <h1>API连接测试</h1>
    
    <div class="test-section">
      <h2>认证状态</h2>
      <div class="status-info">
        <p><strong>管理员Token:</strong> {{ adminToken || '未设置' }}</p>
        <p><strong>管理员信息:</strong> {{ adminInfo || '未设置' }}</p>
        <p><strong>Cookies:</strong> {{ cookies || '无' }}</p>
      </div>
    </div>

    <div class="test-section">
      <h2>API测试</h2>
      <div class="test-buttons">
        <button @click="testDomainApi" :disabled="loading">
          {{ loading ? '测试中...' : '测试域名API' }}
        </button>
        <button @click="testUserApi" :disabled="loading">
          {{ loading ? '测试中...' : '测试用户API' }}
        </button>
        <button @click="testMailboxApi" :disabled="loading">
          {{ loading ? '测试中...' : '测试邮箱API' }}
        </button>
      </div>
    </div>

    <div class="test-section">
      <h2>测试结果</h2>
      <div class="test-results">
        <pre>{{ testResults }}</pre>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import adminApi from '@/services/adminApi'

const adminToken = ref('')
const adminInfo = ref('')
const cookies = ref('')
const loading = ref(false)
const testResults = ref('')

const addResult = (message) => {
  const timestamp = new Date().toLocaleTimeString()
  testResults.value += `[${timestamp}] ${message}\n`
}

const testDomainApi = async () => {
  loading.value = true
  addResult('开始测试域名API...')
  
  try {
    const response = await adminApi.getAllDomains()
    addResult(`域名API响应: ${JSON.stringify(response.data, null, 2)}`)
  } catch (error) {
    addResult(`域名API错误: ${error.message}`)
    if (error.response) {
      addResult(`错误详情: ${JSON.stringify(error.response.data, null, 2)}`)
    }
  } finally {
    loading.value = false
  }
}

const testUserApi = async () => {
  loading.value = true
  addResult('开始测试用户API...')
  
  try {
    const response = await adminApi.getAllUsers()
    addResult(`用户API响应: ${JSON.stringify(response.data, null, 2)}`)
  } catch (error) {
    addResult(`用户API错误: ${error.message}`)
    if (error.response) {
      addResult(`错误详情: ${JSON.stringify(error.response.data, null, 2)}`)
    }
  } finally {
    loading.value = false
  }
}

const testMailboxApi = async () => {
  loading.value = true
  addResult('开始测试邮箱API...')
  
  try {
    const response = await adminApi.getAllMailboxes()
    addResult(`邮箱API响应: ${JSON.stringify(response.data, null, 2)}`)
  } catch (error) {
    addResult(`邮箱API错误: ${error.message}`)
    if (error.response) {
      addResult(`错误详情: ${JSON.stringify(error.response.data, null, 2)}`)
    }
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  adminToken.value = localStorage.getItem('admin_token')
  adminInfo.value = localStorage.getItem('admin_info')
  cookies.value = document.cookie
  
  addResult('页面加载完成')
  addResult(`管理员Token: ${adminToken.value}`)
  addResult(`Cookies: ${cookies.value}`)
})
</script>

<style scoped>
.api-test {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.test-section {
  margin-bottom: 30px;
  padding: 20px;
  border: 1px solid #ddd;
  border-radius: 8px;
  background: #f9f9f9;
}

.status-info p {
  margin: 10px 0;
  padding: 8px;
  background: #fff;
  border-radius: 4px;
}

.test-buttons {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.test-buttons button {
  padding: 10px 20px;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.test-buttons button:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.test-results {
  background: #000;
  color: #0f0;
  padding: 15px;
  border-radius: 4px;
  font-family: monospace;
  max-height: 400px;
  overflow-y: auto;
}

.test-results pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
}
</style>
