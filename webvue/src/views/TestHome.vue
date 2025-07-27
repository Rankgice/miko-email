<template>
  <div class="test-home">
    <h1>测试页面</h1>
    <p>如果您能看到这个页面，说明Vue应用正常工作。</p>
    
    <div class="nav-links">
      <router-link to="/user/login" class="nav-link">用户登录</router-link>
      <router-link to="/admin/login" class="nav-link">管理员登录</router-link>
      <router-link to="/user/compose" class="nav-link">写邮件</router-link>
    </div>

    <div class="debug-info">
      <h3>调试信息：</h3>
      <p>当前路由: {{ $route.path }}</p>
      <p>Vue版本: {{ vueVersion }}</p>
      <p>时间: {{ currentTime }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const currentTime = ref('')
const vueVersion = ref('3.x')

onMounted(() => {
  currentTime.value = new Date().toLocaleString()
})
</script>

<style scoped>
.test-home {
  padding: 40px;
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
  color: #f8fafc;
  min-height: 100vh;
}

h1 {
  color: #3b82f6;
  margin-bottom: 20px;
}

.nav-links {
  display: flex;
  gap: 20px;
  margin: 30px 0;
}

.nav-link {
  padding: 10px 20px;
  background: #2563eb;
  color: white;
  text-decoration: none;
  border-radius: 6px;
  transition: all 0.3s ease;
}

.nav-link:hover {
  background: #1d4ed8;
  transform: translateY(-2px);
}

.debug-info {
  margin-top: 40px;
  padding: 20px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.debug-info h3 {
  color: #22c55e;
  margin-bottom: 15px;
}

.debug-info p {
  margin: 8px 0;
  color: #94a3b8;
}
</style>
