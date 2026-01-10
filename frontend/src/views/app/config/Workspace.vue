<template>
  <div class="workspace">
    <!-- 工作台侧边栏 -->
    <div class="workspace-sidebar">
      <div class="sidebar-menu">
        <div 
          v-for="item in menuItems" 
          :key="item.key"
          class="menu-item"
          :class="{ active: currentMenu === item.key }"
          @click="currentMenu = item.key"
        >
          <el-icon><component :is="item.icon" /></el-icon>
          <span>{{ item.label }}</span>
        </div>
      </div>
      <div class="sidebar-footer">
        <div class="menu-item back-item" @click="$router.push('/apps')">
          <el-icon><ArrowLeft /></el-icon>
          <span>返回APP列表</span>
        </div>
      </div>
    </div>

    <!-- 工作台内容区 -->
    <div class="workspace-content">
      <!-- 数据概览 -->
      <div v-if="currentMenu === 'overview'" class="content-section">
        <div class="section-header">
          <h2>数据概览</h2>
          <p>实时监控应用运行状态</p>
        </div>
        
        <!-- 统计卡片 -->
        <div class="stats-grid">
          <div class="stat-card blue">
            <div class="stat-content">
              <div class="stat-value">{{ stats.userCount.toLocaleString() }}</div>
              <div class="stat-label">总用户数</div>
              <div class="stat-trend up">
                <el-icon><Top /></el-icon>
                <span>+12.5%</span>
              </div>
            </div>
            <div class="stat-icon">
              <el-icon><User /></el-icon>
            </div>
          </div>
          
          <div class="stat-card green">
            <div class="stat-content">
              <div class="stat-value">{{ stats.activeUsers.toLocaleString() }}</div>
              <div class="stat-label">活跃用户</div>
              <div class="stat-trend up">
                <el-icon><Top /></el-icon>
                <span>+8.3%</span>
              </div>
            </div>
            <div class="stat-icon">
              <el-icon><UserFilled /></el-icon>
            </div>
          </div>
          
          <div class="stat-card orange">
            <div class="stat-content">
              <div class="stat-value">{{ stats.todayRequests.toLocaleString() }}</div>
              <div class="stat-label">今日请求</div>
              <div class="stat-trend up">
                <el-icon><Top /></el-icon>
                <span>+15.2%</span>
              </div>
            </div>
            <div class="stat-icon">
              <el-icon><DataLine /></el-icon>
            </div>
          </div>
          
          <div class="stat-card red">
            <div class="stat-content">
              <div class="stat-value">{{ stats.todayErrors }}</div>
              <div class="stat-label">今日异常</div>
              <div class="stat-trend down">
                <el-icon><Bottom /></el-icon>
                <span>-5.1%</span>
              </div>
            </div>
            <div class="stat-icon">
              <el-icon><Warning /></el-icon>
            </div>
          </div>
        </div>

        <!-- 图表区域 -->
        <div class="charts-row">
          <div class="chart-card">
            <div class="chart-header">
              <h3>请求趋势</h3>
              <el-radio-group v-model="chartPeriod" size="small">
                <el-radio-button label="7d">7天</el-radio-button>
                <el-radio-button label="30d">30天</el-radio-button>
              </el-radio-group>
            </div>
            <div class="chart-body" ref="requestChartRef"></div>
          </div>
          
          <div class="chart-card">
            <div class="chart-header">
              <h3>模块调用分布</h3>
            </div>
            <div class="chart-body" ref="moduleChartRef"></div>
          </div>
        </div>
      </div>

      <!-- 用户管理 -->
      <div v-else-if="currentMenu === 'users'" class="content-section">
        <div class="section-header">
          <h2>用户管理</h2>
          <p>管理应用用户数据</p>
        </div>
        
        <!-- 搜索和操作栏 -->
        <div class="toolbar">
          <div class="search-area">
            <el-input 
              v-model="userSearch" 
              placeholder="搜索用户名/手机号/邮箱" 
              prefix-icon="Search"
              clearable
              style="width: 300px"
              @keyup.enter="fetchUserList"
            />
            <el-select v-model="userStatus" placeholder="用户状态" clearable style="width: 120px" @change="fetchUserList">
              <el-option label="全部" value="" />
              <el-option label="正常" value="1" />
              <el-option label="禁用" value="0" />
            </el-select>
            <el-button type="primary" @click="fetchUserList">搜索</el-button>
          </div>
          <div class="action-area">
            <el-button type="primary" :icon="Plus">添加用户</el-button>
            <el-button :icon="Download">导出</el-button>
          </div>
        </div>

        <!-- 用户表格 -->
        <el-table :data="userList" stripe style="width: 100%" v-loading="userLoading">
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="nickname" label="昵称" width="150" />
          <el-table-column prop="phone" label="手机号" width="140" />
          <el-table-column prop="email" label="邮箱" min-width="180" />
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
                {{ row.status === 1 ? '正常' : '禁用' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="注册时间" width="180">
            <template #default="{ row }">
              {{ formatDate(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="150" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link size="small">编辑</el-button>
              <el-button 
                :type="row.status === 1 ? 'danger' : 'success'" 
                link 
                size="small"
                @click="toggleUserStatus(row)"
              >
                {{ row.status === 1 ? '禁用' : '启用' }}
              </el-button>
            </template>
          </el-table-column>
        </el-table>

        <!-- 分页 -->
        <div class="pagination">
          <el-pagination
            v-model:current-page="userPage"
            v-model:page-size="userPageSize"
            :total="userTotal"
            :page-sizes="[10, 20, 50, 100]"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="fetchUserList"
            @current-change="fetchUserList"
          />
        </div>
      </div>

      <!-- 消息推送 -->
      <div v-else-if="currentMenu === 'messages'" class="content-section">
        <div class="section-header">
          <h2>消息推送</h2>
          <p>管理应用消息和推送通知</p>
        </div>

        <el-tabs v-model="messageTab">
          <el-tab-pane label="发送消息" name="send">
            <div class="message-form">
              <el-form :model="messageForm" label-width="100px" :rules="messageRules" ref="messageFormRef">
                <el-form-item label="推送类型" prop="type">
                  <el-radio-group v-model="messageForm.type">
                    <el-radio label="all">全部用户</el-radio>
                    <el-radio label="group">用户分组</el-radio>
                    <el-radio label="user">指定用户</el-radio>
                  </el-radio-group>
                </el-form-item>
                <el-form-item v-if="messageForm.type === 'user'" label="用户ID" prop="userIds">
                  <el-input v-model="messageForm.userIds" placeholder="多个用户ID用逗号分隔" />
                </el-form-item>
                <el-form-item label="消息标题" prop="title">
                  <el-input v-model="messageForm.title" placeholder="请输入消息标题" />
                </el-form-item>
                <el-form-item label="消息内容" prop="content">
                  <el-input v-model="messageForm.content" type="textarea" :rows="4" placeholder="请输入消息内容" />
                </el-form-item>
                <el-form-item>
                  <el-button type="primary" @click="sendMessageNow" :loading="messageSending">立即推送</el-button>
                  <el-button>定时推送</el-button>
                </el-form-item>
              </el-form>
            </div>
          </el-tab-pane>
          
          <el-tab-pane label="推送记录" name="history">
            <el-table :data="messageHistory" stripe v-loading="messageLoading">
              <el-table-column prop="title" label="标题" min-width="200" />
              <el-table-column prop="type" label="推送类型" width="120" />
              <el-table-column prop="content" label="内容" min-width="200" show-overflow-tooltip />
              <el-table-column prop="created_at" label="推送时间" width="180">
                <template #default="{ row }">
                  {{ formatDate(row.created_at) }}
                </template>
              </el-table-column>
              <el-table-column prop="status" label="状态" width="100">
                <template #default="{ row }">
                  <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
                    {{ row.status === 1 ? '已读' : '未读' }}
                  </el-tag>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>
        </el-tabs>
      </div>

      <!-- 日志查询 -->
      <div v-else-if="currentMenu === 'logs'" class="content-section">
        <div class="section-header">
          <h2>日志查询</h2>
          <p>查看应用运行日志和操作记录</p>
        </div>

        <!-- 日志统计卡片 -->
        <div class="log-stats">
          <div class="log-stat-item">
            <span class="label">总日志数</span>
            <span class="value">{{ logStats.total || 0 }}</span>
          </div>
          <div class="log-stat-item error">
            <span class="label">错误日志</span>
            <span class="value">{{ logStats.error_count || 0 }}</span>
          </div>
          <div class="log-stat-item warn">
            <span class="label">警告日志</span>
            <span class="value">{{ logStats.warn_count || 0 }}</span>
          </div>
          <div class="log-stat-item info">
            <span class="label">今日日志</span>
            <span class="value">{{ logStats.today_count || 0 }}</span>
          </div>
        </div>

        <!-- 日志筛选 -->
        <div class="toolbar">
          <div class="search-area">
            <el-date-picker
              v-model="logDateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              style="width: 260px"
              value-format="YYYY-MM-DD"
              @change="fetchLogList"
            />
            <el-select v-model="logLevel" placeholder="日志级别" clearable style="width: 120px" @change="fetchLogList">
              <el-option label="全部" value="" />
              <el-option label="DEBUG" value="debug">
                <el-tag type="info" size="small">DEBUG</el-tag>
              </el-option>
              <el-option label="INFO" value="info">
                <el-tag type="success" size="small">INFO</el-tag>
              </el-option>
              <el-option label="WARN" value="warn">
                <el-tag type="warning" size="small">WARN</el-tag>
              </el-option>
              <el-option label="ERROR" value="error">
                <el-tag type="danger" size="small">ERROR</el-tag>
              </el-option>
            </el-select>
            <el-input 
              v-model="logSearch" 
              placeholder="搜索日志内容" 
              prefix-icon="Search"
              clearable
              style="width: 200px"
              @keyup.enter="fetchLogList"
            />
            <el-button type="primary" @click="fetchLogList">搜索</el-button>
          </div>
          <div class="action-area">
            <el-button :icon="Refresh" @click="fetchLogList">刷新</el-button>
            <el-button :icon="Download" @click="exportLogs">导出</el-button>
          </div>
        </div>

        <!-- 日志列表 -->
        <div class="log-list" v-loading="logLoading">
          <div v-if="logList.length === 0" class="empty-logs">
            <el-empty description="暂无日志数据" />
          </div>
          <div v-else>
            <div v-for="log in logList" :key="log.id" class="log-item" :class="log.level">
              <div class="log-time">{{ formatDate(log.created_at) }}</div>
              <el-tag :type="getLogTagType(log.level)" size="small">{{ (log.level || 'info').toUpperCase() }}</el-tag>
              <div class="log-module" v-if="log.module">{{ log.module }}</div>
              <div class="log-content">{{ log.message }}</div>
            </div>
          </div>
        </div>

        <!-- 日志分页 -->
        <div class="pagination">
          <el-pagination
            v-model:current-page="logPage"
            v-model:page-size="logPageSize"
            :total="logTotal"
            :page-sizes="[20, 50, 100, 200]"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="fetchLogList"
            @current-change="fetchLogList"
          />
        </div>
      </div>

      <!-- 版本管理 -->
      <div v-else-if="currentMenu === 'versions'" class="content-section">
        <div class="section-header">
          <h2>版本管理</h2>
          <p>管理应用版本和更新</p>
        </div>

        <div class="toolbar">
          <div class="search-area">
            <el-select v-model="versionPlatform" placeholder="选择平台" clearable style="width: 120px" @change="fetchVersionList">
              <el-option label="全部" value="" />
              <el-option label="Android" value="android" />
              <el-option label="iOS" value="ios" />
            </el-select>
            <el-select v-model="versionStatus" placeholder="版本状态" clearable style="width: 120px" @change="fetchVersionList">
              <el-option label="全部" value="" />
              <el-option label="已发布" value="published" />
              <el-option label="草稿" value="draft" />
              <el-option label="已下线" value="offline" />
            </el-select>
          </div>
          <div class="action-area">
            <el-button type="primary" :icon="Plus" @click="showVersionDialog = true">发布新版本</el-button>
          </div>
        </div>

        <el-table :data="versionList" stripe v-loading="versionLoading">
          <el-table-column prop="version" label="版本号" width="120" />
          <el-table-column prop="platform" label="平台" width="100">
            <template #default="{ row }">
              <el-tag :type="row.platform === 'android' ? 'success' : 'primary'" size="small">
                {{ row.platform === 'android' ? 'Android' : 'iOS' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="description" label="更新说明" min-width="250" show-overflow-tooltip />
          <el-table-column prop="download_url" label="下载地址" min-width="200" show-overflow-tooltip />
          <el-table-column prop="force_update" label="强制更新" width="100">
            <template #default="{ row }">
              <el-tag :type="row.force_update ? 'danger' : 'info'" size="small">
                {{ row.force_update ? '是' : '否' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="getVersionStatusType(row.status)" size="small">
                {{ getVersionStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="发布时间" width="180">
            <template #default="{ row }">
              {{ formatDate(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="180" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link size="small" @click="editVersion(row)">编辑</el-button>
              <el-button 
                v-if="row.status !== 'published'" 
                type="success" 
                link 
                size="small"
                @click="publishVersionAction(row)"
              >
                发布
              </el-button>
              <el-button 
                v-if="row.status === 'published'" 
                type="danger" 
                link 
                size="small"
                @click="offlineVersionAction(row)"
              >
                下线
              </el-button>
            </template>
          </el-table-column>
        </el-table>

        <!-- 版本分页 -->
        <div class="pagination">
          <el-pagination
            v-model:current-page="versionPage"
            v-model:page-size="versionPageSize"
            :total="versionTotal"
            :page-sizes="[10, 20, 50]"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="fetchVersionList"
            @current-change="fetchVersionList"
          />
        </div>
      </div>
    </div>

    <!-- 发布新版本对话框 -->
    <el-dialog v-model="showVersionDialog" :title="editingVersion ? '编辑版本' : '发布新版本'" width="600px">
      <el-form :model="versionForm" label-width="100px" :rules="versionRules" ref="versionFormRef">
        <el-form-item label="版本号" prop="version">
          <el-input v-model="versionForm.version" placeholder="如：1.0.0" />
        </el-form-item>
        <el-form-item label="平台" prop="platform">
          <el-radio-group v-model="versionForm.platform">
            <el-radio label="android">Android</el-radio>
            <el-radio label="ios">iOS</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="下载地址" prop="download_url">
          <el-input v-model="versionForm.download_url" placeholder="请输入下载地址" />
        </el-form-item>
        <el-form-item label="更新说明" prop="description">
          <el-input v-model="versionForm.description" type="textarea" :rows="4" placeholder="请输入更新说明" />
        </el-form-item>
        <el-form-item label="强制更新">
          <el-switch v-model="versionForm.force_update" />
          <span class="form-tip">开启后，用户必须更新到此版本才能使用</span>
        </el-form-item>
        <el-form-item label="灰度发布">
          <el-switch v-model="versionForm.gray_release" />
        </el-form-item>
        <el-form-item v-if="versionForm.gray_release" label="灰度比例">
          <el-slider v-model="versionForm.gray_percent" :min="1" :max="100" :format-tooltip="val => `${val}%`" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showVersionDialog = false">取消</el-button>
        <el-button type="primary" @click="submitVersion" :loading="versionSubmitting">
          {{ editingVersion ? '保存' : '发布' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import * as echarts from 'echarts'
import { 
  DataLine, User, UserFilled, Warning, Top, Bottom,
  Plus, Download, Refresh, Search,
  House, Management, Bell, Document, Promotion, ArrowLeft
} from '@element-plus/icons-vue'
import {
  getUserList, getUserStats, updateUserStatus,
  getLogList, getLogStats, exportLogs as exportLogsApi,
  getMessageList, sendMessage, batchSendMessage,
  getVersionList, createVersion, publishVersion, offlineVersion
} from '@/api/app'

const props = defineProps({
  appId: String,
  appInfo: Object
})

// 菜单配置
const menuItems = [
  { key: 'overview', label: '数据概览', icon: House },
  { key: 'users', label: '用户管理', icon: User },
  { key: 'messages', label: '消息推送', icon: Bell },
  { key: 'logs', label: '日志查询', icon: Document },
  { key: 'versions', label: '版本管理', icon: Promotion }
]

const currentMenu = ref('overview')
const chartPeriod = ref('7d')

// 统计数据
const stats = ref({
  userCount: 0,
  activeUsers: 0,
  todayRequests: 0,
  todayErrors: 0
})

// 图表引用
const requestChartRef = ref(null)
const moduleChartRef = ref(null)
let requestChart = null
let moduleChart = null

// 用户管理
const userSearch = ref('')
const userStatus = ref('')
const userPage = ref(1)
const userPageSize = ref(10)
const userTotal = ref(0)
const userList = ref([])
const userLoading = ref(false)

// 消息推送
const messageTab = ref('send')
const messageForm = ref({
  type: 'all',
  userIds: '',
  title: '',
  content: ''
})
const messageRules = {
  title: [{ required: true, message: '请输入消息标题', trigger: 'blur' }],
  content: [{ required: true, message: '请输入消息内容', trigger: 'blur' }]
}
const messageFormRef = ref(null)
const messageSending = ref(false)
const messageHistory = ref([])
const messageLoading = ref(false)

// 日志查询
const logDateRange = ref([])
const logLevel = ref('')
const logSearch = ref('')
const logList = ref([])
const logLoading = ref(false)
const logPage = ref(1)
const logPageSize = ref(20)
const logTotal = ref(0)
const logStats = ref({})

// 版本管理
const versionList = ref([])
const versionLoading = ref(false)
const versionPage = ref(1)
const versionPageSize = ref(10)
const versionTotal = ref(0)
const versionPlatform = ref('')
const versionStatus = ref('')
const showVersionDialog = ref(false)
const editingVersion = ref(null)
const versionForm = ref({
  version: '',
  platform: 'android',
  download_url: '',
  description: '',
  force_update: false,
  gray_release: false,
  gray_percent: 10
})
const versionRules = {
  version: [{ required: true, message: '请输入版本号', trigger: 'blur' }],
  platform: [{ required: true, message: '请选择平台', trigger: 'change' }],
  download_url: [{ required: true, message: '请输入下载地址', trigger: 'blur' }],
  description: [{ required: true, message: '请输入更新说明', trigger: 'blur' }]
}
const versionFormRef = ref(null)
const versionSubmitting = ref(false)

// 格式化日期
const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', { 
    year: 'numeric', 
    month: '2-digit', 
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  })
}

const getLogTagType = (level) => {
  const types = { debug: 'info', info: 'success', warn: 'warning', error: 'danger' }
  return types[level] || 'info'
}

const getVersionStatusType = (status) => {
  const types = { published: 'success', draft: 'info', offline: 'danger' }
  return types[status] || 'info'
}

const getVersionStatusText = (status) => {
  const texts = { published: '已发布', draft: '草稿', offline: '已下线' }
  return texts[status] || status
}

// 获取用户列表
const fetchUserList = async () => {
  if (!props.appId) return
  userLoading.value = true
  try {
    const res = await getUserList({
      app_id: props.appId,
      page: userPage.value,
      size: userPageSize.value,
      status: userStatus.value,
      search: userSearch.value
    })
    if (res.code === 0) {
      userList.value = res.data.list || []
      userTotal.value = res.data.total || 0
    }
  } catch (error) {
    console.error('获取用户列表失败:', error)
  } finally {
    userLoading.value = false
  }
}

// 切换用户状态
const toggleUserStatus = async (row) => {
  const newStatus = row.status === 1 ? 0 : 1
  const action = newStatus === 0 ? '禁用' : '启用'
  try {
    await ElMessageBox.confirm(`确定要${action}该用户吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    const res = await updateUserStatus(row.id, newStatus)
    if (res.code === 0) {
      ElMessage.success(`${action}成功`)
      fetchUserList()
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('更新用户状态失败:', error)
    }
  }
}

// 获取用户统计
const fetchUserStats = async () => {
  if (!props.appId) return
  try {
    const res = await getUserStats(props.appId)
    if (res.code === 0) {
      stats.value.userCount = res.data.total || 0
      stats.value.activeUsers = res.data.active || 0
    }
  } catch (error) {
    console.error('获取用户统计失败:', error)
  }
}

// 获取日志列表
const fetchLogList = async () => {
  if (!props.appId) return
  logLoading.value = true
  try {
    const params = {
      app_id: props.appId,
      page: logPage.value,
      size: logPageSize.value
    }
    if (logLevel.value) params.level = logLevel.value
    if (logSearch.value) params.keyword = logSearch.value
    if (logDateRange.value && logDateRange.value.length === 2) {
      params.start_time = logDateRange.value[0]
      params.end_time = logDateRange.value[1]
    }
    
    const res = await getLogList(params)
    if (res.code === 0) {
      logList.value = res.data.list || []
      logTotal.value = res.data.total || 0
    }
  } catch (error) {
    console.error('获取日志列表失败:', error)
  } finally {
    logLoading.value = false
  }
}

// 获取日志统计
const fetchLogStats = async () => {
  if (!props.appId) return
  try {
    const res = await getLogStats(props.appId)
    if (res.code === 0) {
      logStats.value = res.data || {}
      stats.value.todayErrors = res.data.error_count || 0
    }
  } catch (error) {
    console.error('获取日志统计失败:', error)
  }
}

// 导出日志
const exportLogs = async () => {
  if (!props.appId) return
  try {
    const params = {
      app_id: props.appId,
      level: logLevel.value
    }
    if (logDateRange.value && logDateRange.value.length === 2) {
      params.start_time = logDateRange.value[0]
      params.end_time = logDateRange.value[1]
    }
    const res = await exportLogsApi(params)
    if (res.code === 0) {
      ElMessage.success(`导出成功，共${res.data.count}条日志`)
      // 实际项目中这里应该下载文件
    }
  } catch (error) {
    console.error('导出日志失败:', error)
  }
}

// 获取消息列表
const fetchMessageList = async () => {
  if (!props.appId) return
  messageLoading.value = true
  try {
    const res = await getMessageList({
      app_id: props.appId,
      page: 1,
      size: 20
    })
    if (res.code === 0) {
      messageHistory.value = res.data.list || []
    }
  } catch (error) {
    console.error('获取消息列表失败:', error)
  } finally {
    messageLoading.value = false
  }
}

// 发送消息
const sendMessageNow = async () => {
  if (!messageFormRef.value) return
  await messageFormRef.value.validate()
  
  messageSending.value = true
  try {
    let res
    if (messageForm.value.type === 'all') {
      res = await batchSendMessage({
        app_id: parseInt(props.appId),
        title: messageForm.value.title,
        content: messageForm.value.content,
        type: 'system'
      })
    } else if (messageForm.value.type === 'user' && messageForm.value.userIds) {
      const userIds = messageForm.value.userIds.split(',').map(id => parseInt(id.trim()))
      res = await batchSendMessage({
        app_id: parseInt(props.appId),
        user_ids: userIds,
        title: messageForm.value.title,
        content: messageForm.value.content,
        type: 'system'
      })
    } else {
      res = await sendMessage({
        app_id: parseInt(props.appId),
        title: messageForm.value.title,
        content: messageForm.value.content,
        type: 'system'
      })
    }
    
    if (res.code === 0) {
      ElMessage.success('消息推送成功')
      messageForm.value = { type: 'all', userIds: '', title: '', content: '' }
      fetchMessageList()
    }
  } catch (error) {
    console.error('发送消息失败:', error)
  } finally {
    messageSending.value = false
  }
}

// 获取版本列表
const fetchVersionList = async () => {
  if (!props.appId) return
  versionLoading.value = true
  try {
    const params = {
      app_id: props.appId,
      page: versionPage.value,
      size: versionPageSize.value
    }
    if (versionPlatform.value) params.platform = versionPlatform.value
    if (versionStatus.value) params.status = versionStatus.value
    
    const res = await getVersionList(params)
    if (res.code === 0) {
      versionList.value = res.data.list || res.data || []
      versionTotal.value = res.data.total || versionList.value.length
    }
  } catch (error) {
    console.error('获取版本列表失败:', error)
  } finally {
    versionLoading.value = false
  }
}

// 编辑版本
const editVersion = (row) => {
  editingVersion.value = row
  versionForm.value = {
    version: row.version,
    platform: row.platform,
    download_url: row.download_url || '',
    description: row.description,
    force_update: row.force_update,
    gray_release: row.gray_release || false,
    gray_percent: row.gray_percent || 10
  }
  showVersionDialog.value = true
}

// 提交版本
const submitVersion = async () => {
  if (!versionFormRef.value) return
  await versionFormRef.value.validate()
  
  versionSubmitting.value = true
  try {
    const data = {
      app_id: parseInt(props.appId),
      ...versionForm.value
    }
    
    const res = await createVersion(data)
    if (res.code === 0) {
      ElMessage.success(editingVersion.value ? '版本更新成功' : '版本发布成功')
      showVersionDialog.value = false
      editingVersion.value = null
      versionForm.value = {
        version: '',
        platform: 'android',
        download_url: '',
        description: '',
        force_update: false,
        gray_release: false,
        gray_percent: 10
      }
      fetchVersionList()
    }
  } catch (error) {
    console.error('提交版本失败:', error)
  } finally {
    versionSubmitting.value = false
  }
}

// 发布版本
const publishVersionAction = async (row) => {
  try {
    await ElMessageBox.confirm('确定要发布该版本吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    const res = await publishVersion(row.id)
    if (res.code === 0) {
      ElMessage.success('版本发布成功')
      fetchVersionList()
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('发布版本失败:', error)
    }
  }
}

// 下线版本
const offlineVersionAction = async (row) => {
  try {
    await ElMessageBox.confirm('确定要下线该版本吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    const res = await offlineVersion(row.id)
    if (res.code === 0) {
      ElMessage.success('版本已下线')
      fetchVersionList()
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('下线版本失败:', error)
    }
  }
}

// 初始化图表
const initCharts = () => {
  if (requestChartRef.value) {
    requestChart = echarts.init(requestChartRef.value)
    requestChart.setOption({
      tooltip: { trigger: 'axis' },
      grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
      xAxis: {
        type: 'category',
        boundaryGap: false,
        data: ['1/4', '1/5', '1/6', '1/7', '1/8', '1/9', '1/10']
      },
      yAxis: { type: 'value' },
      series: [{
        name: '请求数',
        type: 'line',
        smooth: true,
        areaStyle: { opacity: 0.3 },
        data: [32000, 35000, 38000, 42000, 45000, 43000, 45678],
        itemStyle: { color: '#409eff' }
      }]
    })
  }

  if (moduleChartRef.value) {
    moduleChart = echarts.init(moduleChartRef.value)
    moduleChart.setOption({
      tooltip: { trigger: 'item' },
      legend: { orient: 'vertical', left: 'left' },
      series: [{
        name: '模块调用',
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: { borderRadius: 10, borderColor: '#fff', borderWidth: 2 },
        label: { show: false, position: 'center' },
        emphasis: {
          label: { show: true, fontSize: 16, fontWeight: 'bold' }
        },
        data: [
          { value: 1048, name: '用户管理' },
          { value: 735, name: '消息推送' },
          { value: 580, name: '数据埋点' },
          { value: 484, name: '日志服务' },
          { value: 300, name: '版本管理' }
        ]
      }]
    })
  }
}

// 加载数据
const loadData = () => {
  fetchUserStats()
  fetchLogStats()
}

onMounted(() => {
  setTimeout(initCharts, 100)
  loadData()
})

watch(currentMenu, (val) => {
  if (val === 'overview') {
    setTimeout(initCharts, 100)
    loadData()
  } else if (val === 'users') {
    fetchUserList()
  } else if (val === 'logs') {
    fetchLogList()
    fetchLogStats()
  } else if (val === 'messages') {
    fetchMessageList()
  } else if (val === 'versions') {
    fetchVersionList()
  }
})

watch(() => props.appId, () => {
  loadData()
  if (currentMenu.value === 'users') fetchUserList()
  if (currentMenu.value === 'logs') fetchLogList()
  if (currentMenu.value === 'messages') fetchMessageList()
  if (currentMenu.value === 'versions') fetchVersionList()
})
</script>

<style lang="scss" scoped>
.workspace {
  display: flex;
  height: 100%;
  background: #f5f7fa;
}

.workspace-sidebar {
  width: 200px;
  background: white;
  border-right: 1px solid #e4e7ed;
  padding: 16px 0;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.sidebar-menu {
  flex: 1;
}

.sidebar-footer {
  border-top: 1px solid #e4e7ed;
  padding-top: 8px;
  margin-top: 8px;
  
  .back-item {
    color: #909399;
    
    &:hover {
      color: #409eff;
      background: #f5f7fa;
    }
  }
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 24px;
  cursor: pointer;
  color: #606266;
  font-size: 14px;
  transition: all 0.2s;
  
  &:hover {
    background: #f5f7fa;
    color: #409eff;
  }
  
  &.active {
    background: linear-gradient(90deg, #ecf5ff 0%, transparent 100%);
    color: #409eff;
    border-left: 3px solid #409eff;
    font-weight: 500;
  }
}

.workspace-content {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}

.content-section {
  max-width: 1400px;
}

.section-header {
  margin-bottom: 24px;
  
  h2 {
    font-size: 22px;
    font-weight: 600;
    color: #1a1a2e;
    margin: 0 0 8px;
  }
  
  p {
    font-size: 14px;
    color: #909399;
    margin: 0;
  }
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  margin-bottom: 24px;
}

.stat-card {
  background: white;
  border-radius: 12px;
  padding: 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
  transition: all 0.3s;
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  }
  
  &.blue {
    border-left: 4px solid #409eff;
    .stat-icon { background: #ecf5ff; color: #409eff; }
  }
  
  &.green {
    border-left: 4px solid #67c23a;
    .stat-icon { background: #f0f9eb; color: #67c23a; }
  }
  
  &.orange {
    border-left: 4px solid #e6a23c;
    .stat-icon { background: #fdf6ec; color: #e6a23c; }
  }
  
  &.red {
    border-left: 4px solid #f56c6c;
    .stat-icon { background: #fef0f0; color: #f56c6c; }
  }
  
  .stat-content {
    .stat-value {
      font-size: 32px;
      font-weight: 700;
      color: #1a1a2e;
      line-height: 1;
    }
    
    .stat-label {
      font-size: 14px;
      color: #909399;
      margin-top: 8px;
    }
    
    .stat-trend {
      display: flex;
      align-items: center;
      gap: 4px;
      font-size: 13px;
      margin-top: 8px;
      
      &.up { color: #67c23a; }
      &.down { color: #f56c6c; }
    }
  }
  
  .stat-icon {
    width: 56px;
    height: 56px;
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 28px;
  }
}

.charts-row {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 20px;
}

.chart-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
  
  .chart-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    
    h3 {
      font-size: 16px;
      font-weight: 600;
      color: #1a1a2e;
      margin: 0;
    }
  }
  
  .chart-body {
    height: 300px;
  }
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 16px 20px;
  background: white;
  border-radius: 8px;
  
  .search-area {
    display: flex;
    gap: 12px;
  }
  
  .action-area {
    display: flex;
    gap: 12px;
  }
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.message-form {
  background: white;
  padding: 24px;
  border-radius: 8px;
}

.log-stats {
  display: flex;
  gap: 20px;
  margin-bottom: 20px;
  
  .log-stat-item {
    background: white;
    padding: 16px 24px;
    border-radius: 8px;
    display: flex;
    flex-direction: column;
    gap: 8px;
    min-width: 120px;
    
    .label {
      font-size: 13px;
      color: #909399;
    }
    
    .value {
      font-size: 24px;
      font-weight: 600;
      color: #1a1a2e;
    }
    
    &.error {
      border-left: 3px solid #f56c6c;
      .value { color: #f56c6c; }
    }
    
    &.warn {
      border-left: 3px solid #e6a23c;
      .value { color: #e6a23c; }
    }
    
    &.info {
      border-left: 3px solid #409eff;
      .value { color: #409eff; }
    }
  }
}

.log-list {
  background: white;
  border-radius: 8px;
  overflow: hidden;
  min-height: 200px;
}

.empty-logs {
  padding: 40px;
}

.log-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-bottom: 1px solid #f0f0f0;
  font-size: 13px;
  
  &:last-child {
    border-bottom: none;
  }
  
  &.error {
    background: #fff1f0;
  }
  
  &.warn {
    background: #fffbe6;
  }
  
  .log-time {
    color: #909399;
    font-family: monospace;
    white-space: nowrap;
    min-width: 160px;
  }
  
  .log-module {
    color: #606266;
    background: #f0f0f0;
    padding: 2px 8px;
    border-radius: 4px;
    font-size: 12px;
  }
  
  .log-content {
    flex: 1;
    color: #303133;
  }
}

.form-tip {
  margin-left: 12px;
  font-size: 12px;
  color: #909399;
}

@media (max-width: 1200px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .charts-row {
    grid-template-columns: 1fr;
  }
  
  .log-stats {
    flex-wrap: wrap;
  }
}

@media (max-width: 768px) {
  .workspace {
    flex-direction: column;
  }
  
  .workspace-sidebar {
    width: 100%;
    display: flex;
    overflow-x: auto;
    padding: 0;
    border-right: none;
    border-bottom: 1px solid #e4e7ed;
  }
  
  .menu-item {
    padding: 12px 16px;
    white-space: nowrap;
    border-left: none !important;
    
    &.active {
      border-bottom: 2px solid #409eff;
      background: transparent;
    }
  }
  
  .stats-grid {
    grid-template-columns: 1fr;
  }
  
  .toolbar {
    flex-direction: column;
    gap: 12px;
    
    .search-area, .action-area {
      width: 100%;
      flex-wrap: wrap;
    }
  }
}
</style>
