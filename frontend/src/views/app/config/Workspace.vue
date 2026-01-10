<template>
  <div class="workspace">
    <!-- 工作台侧边栏 -->
    <div class="workspace-sidebar">
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
            />
            <el-select v-model="userStatus" placeholder="用户状态" clearable style="width: 120px">
              <el-option label="全部" value="" />
              <el-option label="正常" value="1" />
              <el-option label="禁用" value="0" />
            </el-select>
          </div>
          <div class="action-area">
            <el-button type="primary" :icon="Plus">添加用户</el-button>
            <el-button :icon="Download">导出</el-button>
          </div>
        </div>

        <!-- 用户表格 -->
        <el-table :data="userList" stripe style="width: 100%">
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="username" label="用户名" width="150" />
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
          <el-table-column prop="created_at" label="注册时间" width="180" />
          <el-table-column label="操作" width="150" fixed="right">
            <template #default>
              <el-button type="primary" link size="small">编辑</el-button>
              <el-button type="danger" link size="small">禁用</el-button>
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
              <el-form :model="messageForm" label-width="100px">
                <el-form-item label="推送类型">
                  <el-radio-group v-model="messageForm.type">
                    <el-radio label="all">全部用户</el-radio>
                    <el-radio label="group">用户分组</el-radio>
                    <el-radio label="user">指定用户</el-radio>
                  </el-radio-group>
                </el-form-item>
                <el-form-item v-if="messageForm.type === 'user'" label="用户ID">
                  <el-input v-model="messageForm.userIds" placeholder="多个用户ID用逗号分隔" />
                </el-form-item>
                <el-form-item label="消息标题">
                  <el-input v-model="messageForm.title" placeholder="请输入消息标题" />
                </el-form-item>
                <el-form-item label="消息内容">
                  <el-input v-model="messageForm.content" type="textarea" :rows="4" placeholder="请输入消息内容" />
                </el-form-item>
                <el-form-item>
                  <el-button type="primary">立即推送</el-button>
                  <el-button>定时推送</el-button>
                </el-form-item>
              </el-form>
            </div>
          </el-tab-pane>
          
          <el-tab-pane label="推送记录" name="history">
            <el-table :data="messageHistory" stripe>
              <el-table-column prop="title" label="标题" min-width="200" />
              <el-table-column prop="type" label="推送类型" width="120" />
              <el-table-column prop="target_count" label="目标用户" width="100" />
              <el-table-column prop="success_count" label="成功数" width="100" />
              <el-table-column prop="created_at" label="推送时间" width="180" />
              <el-table-column prop="status" label="状态" width="100">
                <template #default="{ row }">
                  <el-tag :type="row.status === 'success' ? 'success' : 'warning'" size="small">
                    {{ row.status === 'success' ? '已完成' : '进行中' }}
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
            />
            <el-select v-model="logLevel" placeholder="日志级别" clearable style="width: 120px">
              <el-option label="全部" value="" />
              <el-option label="INFO" value="info" />
              <el-option label="WARN" value="warn" />
              <el-option label="ERROR" value="error" />
            </el-select>
            <el-input 
              v-model="logSearch" 
              placeholder="搜索日志内容" 
              prefix-icon="Search"
              clearable
              style="width: 200px"
            />
          </div>
          <div class="action-area">
            <el-button :icon="Refresh" @click="refreshLogs">刷新</el-button>
            <el-button :icon="Download">导出</el-button>
          </div>
        </div>

        <!-- 日志列表 -->
        <div class="log-list">
          <div v-for="log in logList" :key="log.id" class="log-item" :class="log.level">
            <div class="log-time">{{ log.time }}</div>
            <el-tag :type="getLogTagType(log.level)" size="small">{{ log.level.toUpperCase() }}</el-tag>
            <div class="log-content">{{ log.content }}</div>
          </div>
        </div>
      </div>

      <!-- 版本管理 -->
      <div v-else-if="currentMenu === 'versions'" class="content-section">
        <div class="section-header">
          <h2>版本管理</h2>
          <p>管理应用版本和更新</p>
        </div>

        <div class="toolbar">
          <div class="search-area"></div>
          <div class="action-area">
            <el-button type="primary" :icon="Plus">发布新版本</el-button>
          </div>
        </div>

        <el-table :data="versionList" stripe>
          <el-table-column prop="version" label="版本号" width="120" />
          <el-table-column prop="platform" label="平台" width="100">
            <template #default="{ row }">
              <el-tag size="small">{{ row.platform }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="description" label="更新说明" min-width="250" />
          <el-table-column prop="download_count" label="下载量" width="100" />
          <el-table-column prop="force_update" label="强制更新" width="100">
            <template #default="{ row }">
              <el-tag :type="row.force_update ? 'danger' : 'info'" size="small">
                {{ row.force_update ? '是' : '否' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="发布时间" width="180" />
          <el-table-column label="操作" width="150" fixed="right">
            <template #default>
              <el-button type="primary" link size="small">编辑</el-button>
              <el-button type="danger" link size="small">下架</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import * as echarts from 'echarts'
import { 
  DataLine, User, UserFilled, Warning, Top, Bottom,
  Plus, Download, Refresh, Search,
  House, Management, Bell, Document, Promotion
} from '@element-plus/icons-vue'

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
  userCount: 12580,
  activeUsers: 3420,
  todayRequests: 45678,
  todayErrors: 23
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
const userTotal = ref(1258)
const userList = ref([
  { id: 1, username: 'user001', nickname: '张三', phone: '138****1234', email: 'zhang***@example.com', status: 1, created_at: '2025-01-01 10:00:00' },
  { id: 2, username: 'user002', nickname: '李四', phone: '139****5678', email: 'li***@example.com', status: 1, created_at: '2025-01-02 11:00:00' },
  { id: 3, username: 'user003', nickname: '王五', phone: '137****9012', email: 'wang***@example.com', status: 0, created_at: '2025-01-03 12:00:00' },
])

// 消息推送
const messageTab = ref('send')
const messageForm = ref({
  type: 'all',
  userIds: '',
  title: '',
  content: ''
})
const messageHistory = ref([
  { title: '系统维护通知', type: '全部用户', target_count: 12580, success_count: 12500, created_at: '2025-01-10 10:00:00', status: 'success' },
  { title: '新功能上线', type: '活跃用户', target_count: 3420, success_count: 3400, created_at: '2025-01-09 15:00:00', status: 'success' },
])

// 日志查询
const logDateRange = ref([])
const logLevel = ref('')
const logSearch = ref('')
const logList = ref([
  { id: 1, time: '2025-01-10 14:30:25', level: 'info', content: '用户 user001 登录成功' },
  { id: 2, time: '2025-01-10 14:30:20', level: 'warn', content: '用户 user002 登录失败，密码错误' },
  { id: 3, time: '2025-01-10 14:30:15', level: 'error', content: '数据库连接超时，重试中...' },
  { id: 4, time: '2025-01-10 14:30:10', level: 'info', content: '推送任务 #1234 执行完成' },
])

// 版本管理
const versionList = ref([
  { version: '2.1.0', platform: 'Android', description: '修复已知问题，优化性能', download_count: 1234, force_update: false, created_at: '2025-01-08 10:00:00' },
  { version: '2.1.0', platform: 'iOS', description: '修复已知问题，优化性能', download_count: 890, force_update: false, created_at: '2025-01-08 10:00:00' },
  { version: '2.0.0', platform: 'Android', description: '全新UI设计，新增消息中心', download_count: 5678, force_update: true, created_at: '2025-01-01 10:00:00' },
])

const getLogTagType = (level) => {
  const types = { info: 'info', warn: 'warning', error: 'danger' }
  return types[level] || 'info'
}

const refreshLogs = () => {
  // 刷新日志
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

onMounted(() => {
  setTimeout(initCharts, 100)
})

watch(currentMenu, (val) => {
  if (val === 'overview') {
    setTimeout(initCharts, 100)
  }
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

.log-list {
  background: white;
  border-radius: 8px;
  overflow: hidden;
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
  }
  
  .log-content {
    flex: 1;
    color: #303133;
  }
}

@media (max-width: 1200px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .charts-row {
    grid-template-columns: 1fr;
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
