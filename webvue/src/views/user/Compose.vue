<template>
  <div class="user-compose">
    <div class="page-header">
      <h1>写邮件</h1>
      <div class="header-actions">
        <button class="btn btn-secondary" @click="saveDraft" :disabled="saving">
          <i class="fas fa-save"></i>
          {{ saving ? '保存中...' : '保存草稿' }}
        </button>
        <button class="btn btn-primary" @click="sendEmail" :disabled="sending">
          <i class="fas fa-paper-plane"></i>
          {{ sending ? '发送中...' : '发送邮件' }}
        </button>
      </div>
    </div>

    <div class="compose-form">
      <div class="form-group">
        <label>收件人 *</label>
        <input
          type="email"
          v-model="emailForm.to"
          placeholder="请输入收件人邮箱地址"
          class="form-input"
        />
      </div>

      <div class="form-group">
        <label>主题 *</label>
        <input
          type="text"
          v-model="emailForm.subject"
          placeholder="请输入邮件主题"
          class="form-input"
        />
      </div>

      <div class="form-group">
        <label>邮件内容 *</label>
        <div class="editor-container">
          <textarea
            v-model="emailForm.content"
            class="simple-editor"
            placeholder="请输入邮件内容..."
            rows="15"
          ></textarea>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
// 暂时移除Quill编辑器导入以排查问题
// import { QuillEditor } from '@vueup/vue-quill'
// import '@vueup/vue-quill/dist/vue-quill.snow.css'

// 响应式数据
const emailForm = ref({
  to: '',
  subject: '',
  content: ''
})

const sending = ref(false)
const saving = ref(false)

// 邮件发送
const sendEmail = async () => {
  if (!emailForm.value.to || !emailForm.value.subject) {
    alert('请填写收件人和主题')
    return
  }

  sending.value = true
  try {
    console.log('发送邮件:', emailForm.value)
    alert('邮件发送成功！')
    // 重置表单
    emailForm.value = {
      to: '',
      subject: '',
      content: ''
    }
  } catch (error) {
    console.error('发送失败:', error)
    alert('发送失败，请重试')
  } finally {
    sending.value = false
  }
}

// 保存草稿
const saveDraft = async () => {
  saving.value = true
  try {
    console.log('保存草稿:', emailForm.value)
    alert('草稿保存成功！')
  } catch (error) {
    console.error('保存失败:', error)
    alert('保存失败，请重试')
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.user-compose {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
  background: linear-gradient(135deg, rgba(15, 23, 42, 0.95), rgba(30, 41, 59, 0.95));
  min-height: 100vh;
  color: #e2e8f0;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.page-header h1 {
  color: #f1f5f9;
  margin: 0;
  font-size: 28px;
  font-weight: 600;
}

.header-actions {
  display: flex;
  gap: 15px;
}

.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 8px;
}

.btn-primary {
  background: #3b82f6;
  color: white;
}

.btn-primary:hover {
  background: #2563eb;
  transform: translateY(-1px);
}

.btn-secondary {
  background: rgba(255, 255, 255, 0.1);
  color: #e2e8f0;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.btn-secondary:hover {
  background: rgba(255, 255, 255, 0.15);
}

.compose-form {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  padding: 30px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.form-group {
  margin-bottom: 25px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  color: #f1f5f9;
  font-weight: 500;
  font-size: 14px;
}

.form-input {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.05);
  color: #f1f5f9;
  font-size: 14px;
  transition: all 0.3s ease;
}

.form-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.editor-container {
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  overflow: hidden;
  background: rgba(255, 255, 255, 0.05);
}

.simple-editor {
  width: 100%;
  min-height: 300px;
  padding: 15px;
  border: none;
  background: rgba(255, 255, 255, 0.95);
  color: #333;
  font-size: 14px;
  line-height: 1.6;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  resize: vertical;
  outline: none;
}

.simple-editor:focus {
  background: white;
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.2);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .user-compose {
    padding: 15px;
  }

  .page-header {
    flex-direction: column;
    gap: 15px;
    align-items: stretch;
  }

  .header-actions {
    justify-content: center;
  }

  .compose-form {
    padding: 20px;
  }
}
</style>
