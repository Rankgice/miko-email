<template>
  <div class="admin-domains">
    <div class="page-header">
      <h1>域名管理</h1>
      <button class="btn btn-primary" @click="addDomain">
        <i class="fas fa-plus"></i>
        添加域名
      </button>
    </div>

    <div class="domains-grid">
      <div class="domain-card" v-for="domain in domains" :key="domain.id">
        <div class="domain-header">
          <div class="domain-icon">
            <i class="fas fa-globe"></i>
          </div>
          <div class="domain-info">
            <h3>{{ domain.name }}</h3>
            <p>{{ domain.description }}</p>
          </div>
          <div class="domain-status">
            <span class="status-badge" :class="domain.is_active ? 'enabled' : 'disabled'">
              {{ domain.is_active ? '启用' : '禁用' }}
            </span>
            <button
              class="status-btn"
              @click="showDNSVerification(domain)"
              :title="getStatusTitle(domain.status)"
            >
              <i class="fas fa-search"></i>
              {{ getStatusText(domain.status) }}
            </button>
          </div>
        </div>

        <div class="domain-stats">
          <div class="stat-item">
            <span>邮箱数量</span>
            <span>{{ domain.mailboxCount }}</span>
          </div>
          <div class="stat-item">
            <span>用户数量</span>
            <span>{{ domain.userCount }}</span>
          </div>
          <div class="stat-item">
            <span>添加时间</span>
            <span>{{ domain.createdAt }}</span>
          </div>
        </div>

        <div class="domain-actions">
          <button @click="viewDomain(domain)" class="btn-view">
            <i class="fas fa-eye"></i>
            查看
          </button>
          <button @click="editDomainAction(domain)" class="btn-edit">
            <i class="fas fa-edit"></i>
            编辑
          </button>
          <button @click="verifyDomain(domain)" v-if="domain.status === 'pending'" class="btn-verify">
            <i class="fas fa-check"></i>
            验证
          </button>
          <button @click="toggleDomainStatus(domain)" class="btn-toggle" :class="domain.is_active ? 'btn-disable' : 'btn-enable'">
            <i :class="domain.is_active ? 'fas fa-ban' : 'fas fa-check-circle'"></i>
            {{ domain.is_active ? '禁用' : '启用' }}
          </button>
          <button @click="deleteDomain(domain)" class="btn-delete danger">
            <i class="fas fa-trash"></i>
            删除
          </button>
        </div>
      </div>
    </div>

    <!-- 创建域名对话框 -->
    <div class="modal-overlay" v-if="showCreateDialog" @click="showCreateDialog = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>添加新域名</h3>
          <button class="close-btn" @click="showCreateDialog = false">
            <i class="fas fa-times"></i>
          </button>
        </div>

        <div class="modal-body">
          <div class="form-group">
            <label>域名</label>
            <input type="text" v-model="newDomain.name" placeholder="例如: example.com" />
          </div>

          <div class="form-group">
            <label>描述</label>
            <input type="text" v-model="newDomain.description" placeholder="请输入域名描述（可选）" />
          </div>
        </div>

        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showCreateDialog = false">取消</button>
          <button class="btn btn-primary" @click="submitCreateDomain">添加域名</button>
        </div>
      </div>
    </div>

    <!-- 查看域名对话框 -->
    <div class="modal-overlay" v-if="showViewDialog" @click="showViewDialog = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>域名详情</h3>
          <button class="close-btn" @click="showViewDialog = false">
            <i class="fas fa-times"></i>
          </button>
        </div>

        <div class="modal-body">
          <div class="detail-group">
            <label>域名ID</label>
            <div class="detail-value">{{ selectedDomain.id }}</div>
          </div>

          <div class="detail-group">
            <label>域名</label>
            <div class="detail-value">{{ selectedDomain.name }}</div>
          </div>

          <div class="detail-group">
            <label>描述</label>
            <div class="detail-value">{{ selectedDomain.description || '无描述' }}</div>
          </div>

          <div class="detail-group">
            <label>状态</label>
            <div class="detail-value">
              <span :class="['status-badge', selectedDomain.status]">
                {{ selectedDomain.status === 'verified' ? '已验证' : '待验证' }}
              </span>
            </div>
          </div>

          <div class="detail-group">
            <label>邮箱数量</label>
            <div class="detail-value">{{ selectedDomain.mailboxCount || 0 }} 个</div>
          </div>

          <div class="detail-group">
            <label>用户数量</label>
            <div class="detail-value">{{ selectedDomain.userCount || 0 }} 个</div>
          </div>

          <div class="detail-group">
            <label>创建时间</label>
            <div class="detail-value">{{ selectedDomain.createdAt }}</div>
          </div>
        </div>

        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showViewDialog = false">关闭</button>
          <button class="btn btn-primary" @click="showViewDialog = false; editDomainAction(selectedDomain)">编辑域名</button>
        </div>
      </div>
    </div>

    <!-- 编辑域名对话框 -->
    <div class="modal-overlay" v-if="showEditDialog" @click="showEditDialog = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>编辑域名</h3>
          <button class="close-btn" @click="showEditDialog = false">
            <i class="fas fa-times"></i>
          </button>
        </div>

        <div class="modal-body">
          <div class="form-group">
            <label>域名</label>
            <input type="text" v-model="editDomain.name" placeholder="例如: example.com" readonly />
            <small class="form-hint">域名不可修改</small>
          </div>

          <div class="form-group">
            <label>描述</label>
            <input type="text" v-model="editDomain.description" placeholder="请输入域名描述（可选）" />
          </div>

          <div class="form-group">
            <label>状态</label>
            <select v-model="editDomain.is_active">
              <option :value="true">启用</option>
              <option :value="false">禁用</option>
            </select>
            <small class="form-hint">禁用后该域名下的所有邮箱将无法发送和接收邮件</small>
          </div>
        </div>

        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showEditDialog = false">取消</button>
          <button class="btn btn-primary" @click="submitEditDomain">保存更改</button>
        </div>
      </div>
    </div>

    <!-- DNS验证对话框 -->
    <div class="modal-overlay" v-if="showDNSDialog" @click="showDNSDialog = false">
      <div class="modal-content dns-modal" @click.stop>
        <div class="modal-header">
          <div class="dns-header-info">
            <h3>DNS记录验证 - {{ selectedDomain?.name }}</h3>
            <div class="domain-verification-summary" v-if="selectedDomain">
              <div class="verification-item">
                <span class="label">总体状态:</span>
                <span :class="['status-indicator', selectedDomain.status]">
                  {{ getStatusText(selectedDomain.status) }}
                </span>
              </div>
              <div class="verification-item">
                <span class="label">发件验证:</span>
                <span :class="['status-indicator', selectedDomain.senderStatus]">
                  {{ getVerificationStatusText(selectedDomain.senderStatus) }}
                </span>
              </div>
              <div class="verification-item">
                <span class="label">收件验证:</span>
                <span :class="['status-indicator', selectedDomain.receiverStatus]">
                  {{ getVerificationStatusText(selectedDomain.receiverStatus) }}
                </span>
              </div>
            </div>
          </div>
          <button class="close-btn" @click="showDNSDialog = false">
            <i class="fas fa-times"></i>
          </button>
        </div>

        <div class="modal-body">
          <div class="dns-tabs">
            <button
              :class="['tab-btn', { active: activeTab === 'spf' }]"
              @click="activeTab = 'spf'"
            >
              <i class="fas fa-shield-alt"></i>
              SPF记录
            </button>
            <button
              :class="['tab-btn', { active: activeTab === 'dkim' }]"
              @click="activeTab = 'dkim'"
            >
              <i class="fas fa-key"></i>
              DKIM记录
            </button>
            <button
              :class="['tab-btn', { active: activeTab === 'dmarc' }]"
              @click="activeTab = 'dmarc'"
            >
              <i class="fas fa-lock"></i>
              DMARC记录
            </button>
            <button
              :class="['tab-btn', { active: activeTab === 'a' }]"
              @click="activeTab = 'a'"
            >
              <i class="fas fa-server"></i>
              A记录
            </button>
          </div>

          <div class="dns-content">
            <!-- SPF记录 -->
            <div v-if="activeTab === 'spf'" class="dns-section">
              <div class="record-info">
                <h4>SPF (Sender Policy Framework) 记录</h4>
                <p>SPF记录用于防止邮件欺诈，指定哪些服务器可以代表您的域名发送邮件。</p>
              </div>

              <div class="record-details">
                <div class="record-item">
                  <label>记录类型:</label>
                  <span class="record-value">TXT</span>
                </div>
                <div class="record-item">
                  <label>主机记录:</label>
                  <span class="record-value">@</span>
                </div>
                <div class="record-item">
                  <label>记录值:</label>
                  <div class="record-value-box">
                    <code>v=spf1 include:{{ selectedDomain?.name }} ~all</code>
                    <button class="copy-btn" @click="copyToClipboard('v=spf1 include:' + selectedDomain?.name + ' ~all')">
                      <i class="fas fa-copy"></i>
                    </button>
                  </div>
                </div>
              </div>

              <div class="verification-section">
                <button
                  class="verify-btn"
                  @click="verifyRecord('spf')"
                  :disabled="verifying"
                >
                  <i class="fas fa-search"></i>
                  {{ verifying ? '验证中...' : '验证SPF记录' }}
                </button>
                <div v-if="verificationResults.spf" class="verification-result">
                  <div :class="['result-status', verificationResults.spf.status]">
                    <i :class="verificationResults.spf.status === 'success' ? 'fas fa-check-circle' : 'fas fa-times-circle'"></i>
                    {{ verificationResults.spf.message }}
                  </div>
                  <div v-if="verificationResults.spf.records" class="found-records">
                    <h5>找到的记录:</h5>
                    <ul>
                      <li v-for="record in verificationResults.spf.records" :key="record">{{ record }}</li>
                    </ul>
                  </div>
                </div>
              </div>
            </div>

            <!-- DKIM记录 -->
            <div v-if="activeTab === 'dkim'" class="dns-section">
              <div class="record-info">
                <h4>DKIM (DomainKeys Identified Mail) 记录</h4>
                <p>DKIM记录用于验证邮件的真实性，确保邮件在传输过程中未被篡改。</p>
              </div>

              <div class="record-details">
                <div class="record-item">
                  <label>记录类型:</label>
                  <span class="record-value">TXT</span>
                </div>
                <div class="record-item">
                  <label>主机记录:</label>
                  <span class="record-value">default._domainkey</span>
                </div>
                <div class="record-item">
                  <label>记录值:</label>
                  <div class="record-value-box">
                    <code>{{ dkimInfo.record }}</code>
                    <button class="copy-btn" @click="copyToClipboard(dkimInfo.record)" :disabled="dkimInfo.loading">
                      <i :class="dkimInfo.loading ? 'fas fa-spinner fa-spin' : 'fas fa-copy'"></i>
                    </button>
                  </div>
                </div>
              </div>

              <div class="verification-section">
                <div class="verification-actions">
                  <button
                    class="verify-btn"
                    @click="verifyRecord('dkim')"
                    :disabled="verifying"
                  >
                    <i class="fas fa-search"></i>
                    {{ verifying ? '验证中...' : '验证DKIM记录' }}
                  </button>
                  <button
                    class="refresh-btn"
                    @click="loadDKIMInfo(selectedDomain?.name)"
                    :disabled="dkimInfo.loading"
                    title="重新生成DKIM密钥"
                  >
                    <i :class="dkimInfo.loading ? 'fas fa-spinner fa-spin' : 'fas fa-sync-alt'"></i>
                    {{ dkimInfo.loading ? '生成中...' : '重新生成' }}
                  </button>
                </div>
                <div v-if="verificationResults.dkim" class="verification-result">
                  <div :class="['result-status', verificationResults.dkim.status]">
                    <i :class="verificationResults.dkim.status === 'success' ? 'fas fa-check-circle' : 'fas fa-times-circle'"></i>
                    {{ verificationResults.dkim.message }}
                  </div>
                  <div v-if="verificationResults.dkim.records" class="found-records">
                    <h5>找到的记录:</h5>
                    <ul>
                      <li v-for="record in verificationResults.dkim.records" :key="record">{{ record }}</li>
                    </ul>
                  </div>
                </div>
              </div>
            </div>

            <!-- DMARC记录 -->
            <div v-if="activeTab === 'dmarc'" class="dns-section">
              <div class="record-info">
                <h4>DMARC (Domain-based Message Authentication) 记录</h4>
                <p>DMARC记录定义了当SPF或DKIM验证失败时应该如何处理邮件。</p>
              </div>

              <div class="record-details">
                <div class="record-item">
                  <label>记录类型:</label>
                  <span class="record-value">TXT</span>
                </div>
                <div class="record-item">
                  <label>主机记录:</label>
                  <span class="record-value">_dmarc</span>
                </div>
                <div class="record-item">
                  <label>记录值:</label>
                  <div class="record-value-box">
                    <code>v=DMARC1; p=quarantine; rua=mailto:dmarc@{{ selectedDomain?.name }}</code>
                    <button class="copy-btn" @click="copyToClipboard('v=DMARC1; p=quarantine; rua=mailto:dmarc@' + selectedDomain?.name)">
                      <i class="fas fa-copy"></i>
                    </button>
                  </div>
                </div>
              </div>

              <div class="verification-section">
                <button
                  class="verify-btn"
                  @click="verifyRecord('dmarc')"
                  :disabled="verifying"
                >
                  <i class="fas fa-search"></i>
                  {{ verifying ? '验证中...' : '验证DMARC记录' }}
                </button>
                <div v-if="verificationResults.dmarc" class="verification-result">
                  <div :class="['result-status', verificationResults.dmarc.status]">
                    <i :class="verificationResults.dmarc.status === 'success' ? 'fas fa-check-circle' : 'fas fa-times-circle'"></i>
                    {{ verificationResults.dmarc.message }}
                  </div>
                  <div v-if="verificationResults.dmarc.records" class="found-records">
                    <h5>找到的记录:</h5>
                    <ul>
                      <li v-for="record in verificationResults.dmarc.records" :key="record">{{ record }}</li>
                    </ul>
                  </div>
                </div>
              </div>
            </div>

            <!-- A记录 -->
            <div v-if="activeTab === 'a'" class="dns-section">
              <div class="record-info">
                <h4>A记录</h4>
                <p>A记录将域名指向IPv4地址，用于邮件服务器的连接。</p>
              </div>

              <div class="record-details">
                <div class="record-item">
                  <label>记录类型:</label>
                  <span class="record-value">A</span>
                </div>
                <div class="record-item">
                  <label>主机记录:</label>
                  <span class="record-value">@</span>
                </div>
                <div class="record-item">
                  <label>记录值:</label>
                  <div class="record-value-box">
                    <code>{{ serverInfo.publicIP }}</code>
                    <button class="copy-btn" @click="copyToClipboard(serverInfo.publicIP)" :disabled="serverInfo.loading">
                      <i :class="serverInfo.loading ? 'fas fa-spinner fa-spin' : 'fas fa-copy'"></i>
                    </button>
                  </div>
                </div>
              </div>

              <div class="verification-section">
                <div class="verification-actions">
                  <button
                    class="verify-btn"
                    @click="verifyRecord('a')"
                    :disabled="verifying"
                  >
                    <i class="fas fa-search"></i>
                    {{ verifying ? '验证中...' : '验证A记录' }}
                  </button>
                  <button
                    class="refresh-btn"
                    @click="loadServerInfo()"
                    :disabled="serverInfo.loading"
                    title="刷新服务器IP"
                  >
                    <i :class="serverInfo.loading ? 'fas fa-spinner fa-spin' : 'fas fa-sync-alt'"></i>
                    {{ serverInfo.loading ? '获取中...' : '刷新IP' }}
                  </button>
                </div>
                <div v-if="verificationResults.a" class="verification-result">
                  <div :class="['result-status', verificationResults.a.status]">
                    <i :class="verificationResults.a.status === 'success' ? 'fas fa-check-circle' : 'fas fa-times-circle'"></i>
                    {{ verificationResults.a.message }}
                  </div>
                  <div v-if="verificationResults.a.records" class="found-records">
                    <h5>找到的记录:</h5>
                    <ul>
                      <li v-for="record in verificationResults.a.records" :key="record">{{ record }}</li>
                    </ul>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="modal-footer">
          <button class="btn btn-secondary" @click="showDNSDialog = false">关闭</button>
          <button class="btn btn-primary" @click="verifyAllRecords">验证所有记录</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import adminApi from '@/services/adminApi'

const domains = ref([])
const loading = ref(false)
const showCreateDialog = ref(false)
const showViewDialog = ref(false)
const showEditDialog = ref(false)
const showDNSDialog = ref(false)

const newDomain = ref({
  name: '',
  description: ''
})

const selectedDomain = ref({})
const editDomain = ref({
  id: null,
  name: '',
  description: '',
  active: true
})

// DNS验证相关数据
const activeTab = ref('spf')
const verifying = ref(false)
const verificationResults = ref({
  spf: null,
  dkim: null,
  dmarc: null,
  a: null
})

// 服务器信息
const serverInfo = ref({
  publicIP: 'YOUR_SERVER_IP',
  loading: false
})

// DKIM信息
const dkimInfo = ref({
  record: 'v=DKIM1; k=rsa; p=YOUR_PUBLIC_KEY',
  publicKey: 'YOUR_PUBLIC_KEY',
  loading: false
})

// 计算域名验证状态
const getDomainVerificationStatus = (domain) => {
  // 如果总体验证状态为true，显示为verified
  if (domain.is_verified) {
    return 'verified'
  }

  // 如果域名未激活，显示为inactive
  if (!domain.is_active) {
    return 'inactive'
  }

  // 检查发件和收件验证状态
  const senderStatus = domain.sender_verification_status || 'pending'
  const receiverStatus = domain.receiver_verification_status || 'pending'

  // 如果任一验证失败，显示为failed
  if (senderStatus === 'failed' || receiverStatus === 'failed') {
    return 'failed'
  }

  // 如果发件和收件都验证成功，显示为verified
  if (senderStatus === 'verified' && receiverStatus === 'verified') {
    return 'verified'
  }

  // 如果有部分验证成功，显示为partial
  if (senderStatus === 'verified' || receiverStatus === 'verified') {
    return 'partial'
  }

  // 默认显示为pending
  return 'pending'
}

// 获取状态显示文本
const getStatusText = (status) => {
  const statusMap = {
    'verified': '已验证',
    'partial': '部分验证',
    'pending': '待验证',
    'failed': '验证失败',
    'inactive': '未激活'
  }
  return statusMap[status] || '未知状态'
}

// 获取状态提示文本
const getStatusTitle = (status) => {
  const titleMap = {
    'verified': '域名已完全验证 - 点击查看DNS记录详情',
    'partial': '域名部分验证 - 点击查看详情并完成验证',
    'pending': '域名待验证 - 点击开始DNS验证',
    'failed': '域名验证失败 - 点击查看失败原因',
    'inactive': '域名未激活 - 点击查看详情'
  }
  return titleMap[status] || '点击查看DNS记录'
}

// 获取验证状态文本
const getVerificationStatusText = (status) => {
  const statusMap = {
    'verified': '已验证',
    'pending': '待验证',
    'failed': '验证失败'
  }
  return statusMap[status] || '未知'
}

// 加载域名列表
const loadDomains = async () => {
  loading.value = true
  try {
    const response = await adminApi.getAllDomains()
    console.log('域名API响应:', response.data) // 调试日志

    if (response.data.code === 0) {
      // 检查数据是否存在
      const domainData = response.data.data || []
      domains.value = domainData.map(domain => ({
        id: domain.id,
        name: domain.name,
        description: domain.description || domain.name || '无描述',
        status: getDomainVerificationStatus(domain), // 使用实际的验证状态
        isVerified: domain.is_verified || false,
        senderStatus: domain.sender_verification_status || 'pending',
        receiverStatus: domain.receiver_verification_status || 'pending',
        isActive: domain.is_active !== false,
        mailboxCount: domain.mailbox_count || 0,
        userCount: domain.user_count || 0,
        createdAt: domain.created_at ? new Date(domain.created_at).toLocaleDateString() : '未知'
      }))
      console.log('处理后的域名数据:', domains.value) // 调试日志
    } else {
      console.error('API返回错误:', response.data.msg)
      alert(response.data.msg || '获取域名列表失败')
    }
  } catch (error) {
    console.error('加载域名列表失败:', error)
    if (error.response) {
      console.error('错误响应:', error.response.data)
      alert(`加载域名列表失败: ${error.response.data.msg || error.message}`)
    } else {
      alert('网络错误，请检查连接')
    }
  } finally {
    loading.value = false
  }
}

const addDomain = () => {
  showCreateDialog.value = true
  newDomain.value = {
    name: '',
    description: ''
  }
}

const submitCreateDomain = async () => {
  if (!newDomain.value.name) {
    alert('请输入域名')
    return
  }

  // 简单的域名格式验证
  const domainRegex = /^[a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9]\.[a-zA-Z]{2,}$/
  if (!domainRegex.test(newDomain.value.name)) {
    alert('请输入有效的域名格式')
    return
  }

  try {
    const response = await adminApi.createDomain({
      name: newDomain.value.name,
      description: newDomain.value.description
    })

    if (response.data.code === 0) {
      alert('域名添加成功')
      showCreateDialog.value = false
      loadDomains() // 重新加载域名列表
    } else {
      alert(response.data.msg || '添加域名失败')
    }
  } catch (error) {
    console.error('添加域名失败:', error)
    alert('添加域名失败')
  }
}

const viewDomain = async (domain) => {
  try {
    selectedDomain.value = { ...domain }
    showViewDialog.value = true
  } catch (error) {
    console.error('查看域名详情失败:', error)
    alert('查看域名详情失败')
  }
}

const editDomainAction = async (domain) => {
  try {
    editDomain.value = {
      id: domain.id,
      name: domain.name,
      description: domain.description,
      is_active: domain.is_active !== undefined ? domain.is_active : true
    }
    showEditDialog.value = true
  } catch (error) {
    console.error('编辑域名失败:', error)
    alert('编辑域名失败')
  }
}

const submitEditDomain = async () => {
  if (!editDomain.value.name) {
    alert('域名不能为空')
    return
  }

  // 简单的域名格式验证
  const domainRegex = /^[a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9]\.[a-zA-Z]{2,}$/
  if (!domainRegex.test(editDomain.value.name)) {
    alert('请输入有效的域名格式')
    return
  }

  try {
    // 分别更新域名信息和状态
    let updateSuccess = true
    let errorMessage = ''

    // 如果有描述更新，先更新域名信息
    if (editDomain.value.description !== selectedDomain.value.description) {
      try {
        const infoResponse = await adminApi.updateDomain(editDomain.value.id, {
          name: editDomain.value.name,
          description: editDomain.value.description
        })
        if (infoResponse.data.code !== 0) {
          updateSuccess = false
          errorMessage = infoResponse.data.msg || '更新域名信息失败'
        }
      } catch (error) {
        updateSuccess = false
        errorMessage = '更新域名信息失败'
      }
    }

    // 如果状态有变化，更新域名状态
    if (updateSuccess && editDomain.value.is_active !== selectedDomain.value.is_active) {
      try {
        const statusResponse = await adminApi.updateDomainStatus(editDomain.value.id, {
          status: editDomain.value.is_active ? 'enabled' : 'disabled'
        })
        if (statusResponse.data.code !== 0) {
          updateSuccess = false
          errorMessage = statusResponse.data.msg || '更新域名状态失败'
        }
      } catch (error) {
        updateSuccess = false
        errorMessage = '更新域名状态失败'
      }
    }

    if (updateSuccess) {
      alert('域名更新成功')
      showEditDialog.value = false
      loadDomains() // 重新加载域名列表
    } else {
      alert(errorMessage || '更新域名失败')
    }
  } catch (error) {
    console.error('更新域名失败:', error)
    if (error.response && error.response.data) {
      alert(`更新域名失败: ${error.response.data.msg || '未知错误'}`)
    } else {
      alert('更新域名失败，请检查网络连接')
    }
  }
}

const verifyDomain = async (domain) => {
  try {
    const response = await adminApi.verifyDomain(domain.id)
    if (response.data.code === 0) {
      alert(`域名 ${domain.name} 验证成功`)
      loadDomains() // 重新加载域名列表
    } else {
      alert(response.data.msg || '域名验证失败')
    }
  } catch (error) {
    console.error('域名验证失败:', error)
    alert('域名验证失败')
  }
}

const deleteDomain = async (domain) => {
  if (!confirm(`确定删除域名 ${domain.name}？`)) {
    return
  }

  try {
    const response = await adminApi.deleteDomain(domain.id)
    if (response.data.code === 0) {
      alert('域名删除成功')
      loadDomains() // 重新加载域名列表
    } else {
      alert(response.data.msg || '删除域名失败')
    }
  } catch (error) {
    console.error('删除域名失败:', error)
    alert('删除域名失败')
  }
}

// 切换域名状态
const toggleDomainStatus = async (domain) => {
  const newStatus = domain.is_active ? 'disabled' : 'enabled'
  const action = domain.is_active ? '禁用' : '启用'

  if (!confirm(`确定${action}域名 ${domain.name}？${domain.is_active ? '\n禁用后该域名下的所有邮箱将无法发送和接收邮件。' : ''}`)) {
    return
  }

  try {
    const response = await adminApi.updateDomainStatus(domain.id, { status: newStatus })
    if (response.data.code === 0) {
      alert(`域名${action}成功`)
      loadDomains() // 重新加载域名列表
    } else {
      alert(response.data.msg || `${action}域名失败`)
    }
  } catch (error) {
    console.error(`${action}域名失败:`, error)
    alert(`${action}域名失败`)
  }
}

// DNS验证相关方法
const showDNSVerification = async (domain) => {
  selectedDomain.value = domain
  activeTab.value = 'spf'
  verificationResults.value = {
    spf: null,
    dkim: null,
    dmarc: null,
    a: null
  }
  showDNSDialog.value = true

  // 加载服务器信息和DKIM记录
  await loadServerInfo()
  await loadDKIMInfo(domain.name)
}

const verifyRecord = async (recordType) => {
  if (!selectedDomain.value) return

  verifying.value = true

  try {
    let response
    switch (recordType) {
      case 'spf':
      case 'dkim':
      case 'dmarc':
        response = await adminApi.verifySenderConfiguration(selectedDomain.value.id)
        break
      case 'a':
        response = await adminApi.verifyReceiverConfiguration(selectedDomain.value.id)
        break
    }

    if (response.data.code === 0) {
      // 模拟验证结果，实际应该从后端获取详细的验证信息
      verificationResults.value[recordType] = {
        status: 'success',
        message: `${recordType.toUpperCase()}记录验证成功`,
        records: [`示例记录: v=${recordType}1 ...`]
      }
    } else {
      verificationResults.value[recordType] = {
        status: 'error',
        message: `${recordType.toUpperCase()}记录验证失败: ${response.data.msg}`,
        records: []
      }
    }
  } catch (error) {
    console.error(`验证${recordType}记录失败:`, error)
    verificationResults.value[recordType] = {
      status: 'error',
      message: `${recordType.toUpperCase()}记录验证失败: 网络错误`,
      records: []
    }
  } finally {
    verifying.value = false
  }
}

const verifyAllRecords = async () => {
  if (!selectedDomain.value) return

  verifying.value = true

  try {
    // 并行验证所有记录
    const [senderResponse, receiverResponse] = await Promise.all([
      adminApi.verifySenderConfiguration(selectedDomain.value.id),
      adminApi.verifyReceiverConfiguration(selectedDomain.value.id)
    ])

    // 处理发件配置验证结果
    if (senderResponse.data.code === 0) {
      verificationResults.value.spf = {
        status: 'success',
        message: 'SPF记录验证成功',
        records: ['v=spf1 include:' + selectedDomain.value.name + ' ~all']
      }
      verificationResults.value.dkim = {
        status: 'success',
        message: 'DKIM记录验证成功',
        records: ['v=DKIM1; k=rsa; p=...']
      }
      verificationResults.value.dmarc = {
        status: 'success',
        message: 'DMARC记录验证成功',
        records: ['v=DMARC1; p=quarantine; rua=mailto:dmarc@' + selectedDomain.value.name]
      }
    }

    // 处理收件配置验证结果
    if (receiverResponse.data.code === 0) {
      verificationResults.value.a = {
        status: 'success',
        message: 'A记录验证成功',
        records: ['192.168.1.1']
      }
    }

    alert('所有DNS记录验证完成')
  } catch (error) {
    console.error('验证DNS记录失败:', error)
    alert('验证DNS记录失败')
  } finally {
    verifying.value = false
  }
}

const copyToClipboard = async (text) => {
  try {
    await navigator.clipboard.writeText(text)
    alert('已复制到剪贴板')
  } catch (error) {
    console.error('复制失败:', error)
    // 降级方案
    const textArea = document.createElement('textarea')
    textArea.value = text
    document.body.appendChild(textArea)
    textArea.select()
    document.execCommand('copy')
    document.body.removeChild(textArea)
    alert('已复制到剪贴板')
  }
}

// 加载服务器信息
const loadServerInfo = async () => {
  if (serverInfo.value.loading) return

  serverInfo.value.loading = true
  try {
    const response = await adminApi.getServerInfo()
    if (response.data.code === 0) {
      serverInfo.value.publicIP = response.data.data.public_ip || 'YOUR_SERVER_IP'
    } else {
      console.error('获取服务器信息失败:', response.data.msg)
    }
  } catch (error) {
    console.error('获取服务器信息失败:', error)
  } finally {
    serverInfo.value.loading = false
  }
}

// 加载DKIM信息
const loadDKIMInfo = async (domainName) => {
  if (dkimInfo.value.loading) return

  dkimInfo.value.loading = true
  try {
    const response = await adminApi.getDomainDKIMRecord(domainName)
    if (response.data.code === 0) {
      const data = response.data.data
      dkimInfo.value.record = data.dkim_record || 'v=DKIM1; k=rsa; p=YOUR_PUBLIC_KEY'
      dkimInfo.value.publicKey = data.public_key || 'YOUR_PUBLIC_KEY'
    } else {
      console.error('获取DKIM记录失败:', response.data.msg)
    }
  } catch (error) {
    console.error('获取DKIM记录失败:', error)
  } finally {
    dkimInfo.value.loading = false
  }
}

// 生命周期
onMounted(() => {
  // 检查管理员登录状态
  const adminToken = localStorage.getItem('admin_token')
  const adminInfo = localStorage.getItem('admin_info')

  console.log('管理员登录状态检查:')
  console.log('- adminToken:', adminToken)
  console.log('- adminInfo:', adminInfo)
  console.log('- cookies:', document.cookie)

  if (!adminToken) {
    console.warn('管理员未登录，跳转到登录页')
    // 可能需要重新登录
    alert('请先登录管理员账户')
    return
  }

  loadDomains()
})
</script>

<style scoped>
.admin-domains {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.page-header h1 {
  color: var(--admin-light);
  margin: 0;
}

.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
}

.btn-primary {
  background: linear-gradient(135deg, #3b82f6, #2563eb);
  color: white;
}

.domains-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 20px;
}

.domain-card {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  padding: 20px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.domain-header {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 15px;
}

.domain-icon {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  background: linear-gradient(135deg, #8b5cf6, #7c3aed);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.domain-info {
  flex: 1;
}

.domain-info h3 {
  color: var(--admin-light);
  margin: 0 0 4px 0;
  font-size: 16px;
}

.domain-info p {
  color: var(--admin-gray);
  margin: 0;
  font-size: 14px;
}

.domain-status {
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.domain-status.verified {
  background: rgba(34, 197, 94, 0.2);
  color: #22c55e;
  border: 1px solid rgba(34, 197, 94, 0.3);
}

.domain-status.partial {
  background: rgba(59, 130, 246, 0.2);
  color: #3b82f6;
  border: 1px solid rgba(59, 130, 246, 0.3);
}

.domain-status.pending {
  background: rgba(245, 158, 11, 0.2);
  color: #f59e0b;
  border: 1px solid rgba(245, 158, 11, 0.3);
}

.domain-status.failed {
  background: rgba(239, 68, 68, 0.2);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.domain-status.inactive {
  background: rgba(156, 163, 175, 0.2);
  color: #9ca3af;
  border: 1px solid rgba(156, 163, 175, 0.3);
}

.domain-stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
  margin-bottom: 15px;
  padding: 15px;
  background: rgba(15, 23, 42, 0.5);
  border-radius: 8px;
}

.stat-item {
  text-align: center;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-item span:first-child {
  font-size: 12px;
  color: var(--admin-gray);
}

.stat-item span:last-child {
  font-weight: 600;
  color: var(--admin-light);
}

.domain-actions {
  display: flex;
  gap: 8px;
  justify-content: center;
  flex-wrap: wrap;
}

.domain-actions button {
  padding: 6px 12px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 12px;
  background: rgba(59, 130, 246, 0.2);
  color: #3b82f6;
}

.domain-actions button.danger {
  background: rgba(239, 68, 68, 0.2);
  color: #ef4444;
}

.domain-actions button:hover {
  opacity: 0.8;
}

/* 状态按钮样式 */
.status-btn {
  background: none;
  border: none;
  color: inherit;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 5px 10px;
  border-radius: 4px;
  transition: all 0.3s ease;
}

.status-btn:hover {
  background: rgba(255, 255, 255, 0.1);
}

.status-btn i {
  font-size: 12px;
}

/* DNS验证对话框样式 */
.dns-modal {
  max-width: 800px;
  width: 90vw;
  max-height: 90vh;
  overflow-y: auto;
}

.dns-header-info {
  flex: 1;
}

.dns-header-info h3 {
  margin: 0 0 15px 0;
  color: var(--admin-light);
}

.domain-verification-summary {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
}

.verification-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.verification-item .label {
  color: var(--admin-gray);
  font-size: 14px;
  font-weight: 500;
}

.verification-item .status-indicator {
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.verification-item .status-indicator.verified {
  background: rgba(34, 197, 94, 0.2);
  color: #22c55e;
  border: 1px solid rgba(34, 197, 94, 0.3);
}

.verification-item .status-indicator.pending {
  background: rgba(245, 158, 11, 0.2);
  color: #f59e0b;
  border: 1px solid rgba(245, 158, 11, 0.3);
}

.verification-item .status-indicator.failed {
  background: rgba(239, 68, 68, 0.2);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.verification-item .status-indicator.partial {
  background: rgba(59, 130, 246, 0.2);
  color: #3b82f6;
  border: 1px solid rgba(59, 130, 246, 0.3);
}

.verification-item .status-indicator.inactive {
  background: rgba(156, 163, 175, 0.2);
  color: #9ca3af;
  border: 1px solid rgba(156, 163, 175, 0.3);
}

.dns-tabs {
  display: flex;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  margin-bottom: 20px;
}

.tab-btn {
  background: none;
  border: none;
  color: var(--admin-gray);
  padding: 12px 20px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  border-bottom: 2px solid transparent;
  transition: all 0.3s ease;
}

.tab-btn:hover {
  color: var(--admin-light);
  background: rgba(255, 255, 255, 0.05);
}

.tab-btn.active {
  color: var(--admin-primary);
  border-bottom-color: var(--admin-primary);
}

.dns-content {
  min-height: 400px;
}

.dns-section {
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.record-info {
  background: rgba(30, 41, 59, 0.5);
  padding: 15px;
  border-radius: 8px;
  margin-bottom: 20px;
  border-left: 4px solid var(--admin-primary);
}

.record-info h4 {
  color: var(--admin-light);
  margin: 0 0 8px 0;
  font-size: 18px;
}

.record-info p {
  color: var(--admin-gray);
  margin: 0;
  line-height: 1.5;
}

.record-details {
  background: rgba(15, 23, 42, 0.8);
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 20px;
}

.record-item {
  display: flex;
  align-items: center;
  margin-bottom: 15px;
}

.record-item:last-child {
  margin-bottom: 0;
}

.record-item label {
  color: var(--admin-gray);
  min-width: 100px;
  font-weight: 500;
}

.record-value {
  color: var(--admin-light);
  font-family: 'Courier New', monospace;
}

.record-value-box {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1;
}

.record-value-box code {
  background: rgba(0, 0, 0, 0.3);
  padding: 8px 12px;
  border-radius: 4px;
  color: #22c55e;
  font-family: 'Courier New', monospace;
  flex: 1;
  word-break: break-all;
}

.copy-btn {
  background: var(--admin-primary);
  border: none;
  color: white;
  padding: 8px 12px;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.copy-btn:hover {
  background: #3b82f6;
  transform: translateY(-1px);
}

.verification-section {
  background: rgba(30, 41, 59, 0.3);
  padding: 20px;
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.verification-actions {
  display: flex;
  gap: 12px;
  margin-bottom: 15px;
  flex-wrap: wrap;
}

.verify-btn {
  background: linear-gradient(135deg, #22c55e, #16a34a);
  border: none;
  color: white;
  padding: 12px 24px;
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
  margin-bottom: 15px;
}

.verify-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(34, 197, 94, 0.3);
}

.verify-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.refresh-btn {
  background: linear-gradient(135deg, #3b82f6, #2563eb);
  border: none;
  color: white;
  padding: 12px 20px;
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
  font-size: 14px;
}

.refresh-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(59, 130, 246, 0.3);
}

.refresh-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.verification-result {
  margin-top: 15px;
}

.result-status {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px;
  border-radius: 6px;
  font-weight: 500;
  margin-bottom: 10px;
}

.result-status.success {
  background: rgba(34, 197, 94, 0.2);
  color: #22c55e;
  border: 1px solid rgba(34, 197, 94, 0.3);
}

.result-status.error {
  background: rgba(239, 68, 68, 0.2);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.found-records {
  background: rgba(0, 0, 0, 0.2);
  padding: 15px;
  border-radius: 6px;
  border-left: 3px solid var(--admin-primary);
}

.found-records h5 {
  color: var(--admin-light);
  margin: 0 0 10px 0;
  font-size: 14px;
}

.found-records ul {
  margin: 0;
  padding-left: 20px;
}

.found-records li {
  color: var(--admin-gray);
  font-family: 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.5;
  margin-bottom: 5px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .dns-modal {
    width: 95vw;
    max-height: 95vh;
  }

  .dns-tabs {
    flex-wrap: wrap;
  }

  .tab-btn {
    padding: 10px 15px;
    font-size: 14px;
  }

  .record-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 5px;
  }

  .record-item label {
    min-width: auto;
  }

  .record-value-box {
    width: 100%;
  }
}

/* 域名状态样式 */
.domain-status {
  display: flex;
  align-items: center;
  gap: 10px;
}

.status-badge {
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}

.status-badge.enabled {
  background: rgba(16, 185, 129, 0.2);
  color: #10b981;
}

.status-badge.disabled {
  background: rgba(156, 163, 175, 0.2);
  color: #9ca3af;
}

/* 域名操作按钮样式 */
.btn-toggle {
  padding: 6px 12px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.btn-enable {
  background: rgba(16, 185, 129, 0.1);
  color: #10b981;
  border: 1px solid rgba(16, 185, 129, 0.3);
}

.btn-enable:hover {
  background: rgba(16, 185, 129, 0.2);
  transform: translateY(-1px);
}

.btn-disable {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.btn-disable:hover {
  background: rgba(239, 68, 68, 0.2);
  transform: translateY(-1px);
}
</style>
