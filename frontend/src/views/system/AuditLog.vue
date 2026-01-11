<template>
  <div class="audit-log-page">
    <div class="page-header">
      <h2>操作审计日志</h2>
      <p class="description">记录系统所有敏感操作，便于安全审计和问题追溯</p>
    </div>

    <!-- 筛选条件 -->
    <el-card class="filter-card" shadow="never">
      <el-form :model="filters" inline>
        <el-form-item label="操作类型">
          <el-select v-model="filters.action" placeholder="全部" clearable style="width: 120px">
            <el-option label="查看" value="view" />
            <el-option label="创建" value="create" />
            <el-option label="更新" value="update" />
            <el-option label="删除" value="delete" />
            <el-option label="登录" value="login" />
            <el-option label="登出" value="logout" />
            <el-option label="上传" value="upload" />
            <el-option label="下载" value="download" />
          </el-select>
        </el-form-item>
        <el-form-item label="资源类型">
          <el-select v-model="filters.resource" placeholder="全部" clearable style="width: 120px">
            <el-option label="应用" value="app" />
            <el-option label="用户" value="user" />
            <el-option label="文件" value="file" />
            <el-option label="版本" value="version" />
            <el-option label="消息" value="message" />
            <el-option label="配置" value="config" />
            <el-option label="监控" value="monitor" />
          </el-select>
        </el-form-item>
        <el-form-item label="操作人">
          <el-input v-model="filters.user_name" placeholder="用户名" clearable style="width: 150px" />
        </el-form-item>
        <el-form-item label="IP地址">
          <el-input v-model="filters.ip_address" placeholder="IP地址" clearable style="width: 150px" />
        </el-form-item>
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="filters.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
            style="width: 260px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchLogs">查询</el-button>
          <el-button @click="resetFilters">重置</el-button>
          <el-button type="success" @click="exportLogs">导出CSV</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 统计卡片 -->
    <div class="stats-row">
      <el-card class="stat-card" shadow="never">
        <div class="stat-value">{{ stats.total || 0 }}</div>
        <div class="stat-label">总操作数</div>
      </el-card>
      <el-card class="stat-card" shadow="never">
        <div class="stat-value">{{ stats.today || 0 }}</div>
        <div class="stat-label">今日操作</div>
      </el-card>
      <el-card class="stat-card" shadow="never">
        <div class="stat-value">{{ stats.uniqueUsers || 0 }}</div>
        <div class="stat-label">活跃用户</div>
      </el-card>
      <el-card class="stat-card" shadow="never">
        <div class="stat-value warning">{{ stats.sensitiveOps || 0 }}</div>
        <div class="stat-label">敏感操作</div>
      </el-card>
    </div>

    <!-- 日志列表 -->
    <el-card shadow="never">
      <el-table :data="logs" v-loading="loading" stripe style="width: 100%">
        <el-table-column prop="created_at" label="时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column prop="user_name" label="操作人" width="120">
          <template #default="{ row }">
            {{ row.user_name || '系统' }}
          </template>
        </el-table-column>
        <el-table-column prop="action" label="操作类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getActionTagType(row.action)" size="small">
              {{ getActionLabel(row.action) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="resource" label="资源类型" width="100">
          <template #default="{ row }">
            {{ getResourceLabel(row.resource) }}
          </template>
        </el-table-column>
        <el-table-column prop="description" label="操作描述" min-width="200" />
        <el-table-column prop="ip_address" label="IP地址" width="140" />
        <el-table-column prop="status_code" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status_code < 400 ? 'success' : 'danger'" size="small">
              {{ row.status_code }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="duration" label="耗时" width="80">
          <template #default="{ row }">
            {{ row.duration }}ms
          </template>
        </el-table-column>
        <el-table-column label="操作" width="80" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="showDetail(row)">
              详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.size"
          :total="pagination.total"
          :page-sizes="[20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="fetchLogs"
          @current-change="fetchLogs"
        />
      </div>
    </el-card>

    <!-- 详情弹窗 -->
    <el-dialog v-model="detailVisible" title="操作详情" width="600px">
      <el-descriptions :column="2" border v-if="currentLog">
        <el-descriptions-item label="操作时间" :span="2">
          {{ formatTime(currentLog.created_at) }}
        </el-descriptions-item>
        <el-descriptions-item label="操作人">
          {{ currentLog.user_name || '系统' }}
        </el-descriptions-item>
        <el-descriptions-item label="用户ID">
          {{ currentLog.user_id || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="操作类型">
          <el-tag :type="getActionTagType(currentLog.action)" size="small">
            {{ getActionLabel(currentLog.action) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="资源类型">
          {{ getResourceLabel(currentLog.resource) }}
        </el-descriptions-item>
        <el-descriptions-item label="资源ID">
          {{ currentLog.resource_id || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="应用ID">
          {{ currentLog.app_id || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="操作描述" :span="2">
          {{ currentLog.description }}
        </el-descriptions-item>
        <el-descriptions-item label="请求路径" :span="2">
          <code>{{ currentLog.request_method }} {{ currentLog.request_path }}</code>
        </el-descriptions-item>
        <el-descriptions-item label="IP地址">
          {{ currentLog.ip_address }}
        </el-descriptions-item>
        <el-descriptions-item label="状态码">
          <el-tag :type="currentLog.status_code < 400 ? 'success' : 'danger'" size="small">
            {{ currentLog.status_code }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="耗时">
          {{ currentLog.duration }}ms
        </el-descriptions-item>
        <el-descriptions-item label="User-Agent" :span="2">
          <div class="user-agent">{{ currentLog.user_agent }}</div>
        </el-descriptions-item>
        <el-descriptions-item label="请求体" :span="2" v-if="currentLog.request_body">
          <pre class="request-body">{{ formatJSON(currentLog.request_body) }}</pre>
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const loading = ref(false)
const logs = ref([])
const stats = ref({})
const detailVisible = ref(false)
const currentLog = ref(null)

const filters = reactive({
  action: '',
  resource: '',
  user_name: '',
  ip_address: '',
  dateRange: []
})

const pagination = reactive({
  page: 1,
  size: 20,
  total: 0
})

const actionLabels = {
  view: '查看',
  create: '创建',
  update: '更新',
  delete: '删除',
  login: '登录',
  logout: '登出',
  upload: '上传',
  download: '下载',
  export: '导出',
  import: '导入',
  publish: '发布',
  send: '发送'
}

const resourceLabels = {
  app: '应用',
  user: '用户',
  file: '文件',
  version: '版本',
  message: '消息',
  config: '配置',
  monitor: '监控',
  admin: '管理员',
  module: '模块',
  audit: '审计日志',
  alert: '告警'
}

const getActionLabel = (action) => actionLabels[action] || action
const getResourceLabel = (resource) => resourceLabels[resource] || resource

const getActionTagType = (action) => {
  const types = {
    view: 'info',
    create: 'success',
    update: 'warning',
    delete: 'danger',
    login: 'success',
    logout: 'info',
    upload: 'success',
    download: 'info'
  }
  return types[action] || 'info'
}

const formatTime = (time) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

const formatJSON = (str) => {
  try {
    return JSON.stringify(JSON.parse(str), null, 2)
  } catch {
    return str
  }
}

const fetchLogs = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      size: pagination.size,
      ...filters
    }
    if (filters.dateRange && filters.dateRange.length === 2) {
      params.start_date = filters.dateRange[0]
      params.end_date = filters.dateRange[1]
    }
    delete params.dateRange

    const res = await request.get('/api/v1/audit/logs', { params })
    if (res.code === 0) {
      logs.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } catch (error) {
    console.error('获取审计日志失败:', error)
  } finally {
    loading.value = false
  }
}

const fetchStats = async () => {
  try {
    const res = await request.get('/api/v1/audit/stats')
    if (res.code === 0) {
      stats.value = res.data || {}
    }
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

const resetFilters = () => {
  filters.action = ''
  filters.resource = ''
  filters.user_name = ''
  filters.ip_address = ''
  filters.dateRange = []
  pagination.page = 1
  fetchLogs()
}

const showDetail = (row) => {
  currentLog.value = row
  detailVisible.value = true
}

const exportLogs = async () => {
  try {
    const params = { ...filters }
    if (filters.dateRange && filters.dateRange.length === 2) {
      params.start_date = filters.dateRange[0]
      params.end_date = filters.dateRange[1]
    }
    delete params.dateRange

    const res = await request.get('/api/v1/audit/export', { 
      params,
      responseType: 'blob'
    })
    
    const blob = new Blob([res], { type: 'text/csv' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `audit_logs_${new Date().toISOString().slice(0, 10)}.csv`
    link.click()
    window.URL.revokeObjectURL(url)
    
    ElMessage.success('导出成功')
  } catch (error) {
    ElMessage.error('导出失败')
  }
}

onMounted(() => {
  fetchLogs()
  fetchStats()
})
</script>

<style scoped>
.audit-log-page {
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0 0 8px 0;
  font-size: 20px;
  font-weight: 600;
}

.page-header .description {
  margin: 0;
  color: #909399;
  font-size: 14px;
}

.filter-card {
  margin-bottom: 20px;
}

.stats-row {
  display: flex;
  gap: 16px;
  margin-bottom: 20px;
}

.stat-card {
  flex: 1;
  text-align: center;
}

.stat-value {
  font-size: 28px;
  font-weight: 600;
  color: #409EFF;
}

.stat-value.warning {
  color: #E6A23C;
}

.stat-label {
  margin-top: 8px;
  color: #909399;
  font-size: 14px;
}

.pagination-wrapper {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.user-agent {
  word-break: break-all;
  font-size: 12px;
  color: #909399;
}

.request-body {
  margin: 0;
  padding: 10px;
  background: #f5f7fa;
  border-radius: 4px;
  font-size: 12px;
  max-height: 200px;
  overflow: auto;
}
</style>
