<template>
  <div class="dashboard-container">
    <el-row :gutter="20" class="stat-row">
      <el-col :xs="12" :sm="12" :md="6" :lg="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background-color: #409eff;">
              <el-icon :size="32"><Grid /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.totalApps }}</div>
              <div class="stat-label">APP总数</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :xs="12" :sm="12" :md="6" :lg="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background-color: #67c23a;">
              <el-icon :size="32"><User /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.totalUsers }}</div>
              <div class="stat-label">用户总数</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :xs="12" :sm="12" :md="6" :lg="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background-color: #e6a23c;">
              <el-icon :size="32"><TrendCharts /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.todayNewUsers }}</div>
              <div class="stat-label">今日新增</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :xs="12" :sm="12" :md="6" :lg="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" style="background-color: #f56c6c;">
              <el-icon :size="32"><Money /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stats.activeApps }}</div>
              <div class="stat-label">活跃APP</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :xs="24" :sm="24" :md="12" :lg="12">
        <el-card>
          <template #header>
            <span>用户增长趋势</span>
          </template>
          <div ref="userTrendChart" style="height: 300px;"></div>
        </el-card>
      </el-col>

      <el-col :xs="24" :sm="24" :md="12" :lg="12">
        <el-card>
          <template #header>
            <span>APP使用统计</span>
          </template>
          <div ref="appUsageChart" style="height: 300px;"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :span="24">
        <el-card>
          <template #header>
            <span>快捷操作</span>
          </template>
          <div class="quick-actions">
            <el-button type="primary" icon="Plus" @click="$router.push('/app/list')">
              创建APP
            </el-button>
            <el-button type="success" icon="Grid" @click="$router.push('/app/list')">
              APP列表
            </el-button>
            <el-button type="info" icon="Setting" @click="$router.push('/system/admin')">
              系统管理
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import * as echarts from 'echarts'
import { getAppList } from '@/api/app'
import request from '@/utils/request'

const stats = ref({
  totalApps: 4,
  totalUsers: 0,
  todayNewUsers: 0,
  activeApps: 3
})

const userTrendChart = ref(null)
const appUsageChart = ref(null)

const loadStats = async () => {
  try {
    const res = await request({
      url: '/stats',
      method: 'get'
    })
    // request.js已解包，res直接是数据对象
    stats.value.totalApps = res.app_count || 0
    stats.value.totalUsers = res.user_count || 0
    stats.value.todayNew = res.today_new_apps || 0
    stats.value.activeApps = res.active_apps || 0
  } catch (error) {
    console.error('Failed to load stats:', error)
  }
}

const initUserTrendChart = () => {
  if (!userTrendChart.value) return
  const chart = echarts.init(userTrendChart.value)
  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日'],
      axisTick: {
        alignWithLabel: true
      }
    },
    yAxis: {
      type: 'value'
    },
    series: [{
      name: '用户增长',
      type: 'line',
      smooth: true,
      data: [120, 132, 101, 134, 90, 230, 210],
      itemStyle: {
        color: '#409EFF'
      },
      areaStyle: {
        color: {
          type: 'linear',
          x: 0,
          y: 0,
          x2: 0,
          y2: 1,
          colorStops: [{
            offset: 0, color: 'rgba(64, 158, 255, 0.3)'
          }, {
            offset: 1, color: 'rgba(64, 158, 255, 0)'
          }]
        }
      }
    }]
  }
  chart.setOption(option)
  window.addEventListener('resize', () => chart.resize())
}

const initAppUsageChart = () => {
  if (!appUsageChart.value) return
  const chart = echarts.init(appUsageChart.value)
  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      left: 10,
      data: ['测试APP', '演示APP', '测试模块关联', '测试新APP']
    },
    series: [{
      name: 'APP使用统计',
      type: 'pie',
      radius: ['40%', '70%'],
      avoidLabelOverlap: false,
      itemStyle: {
        borderRadius: 10,
        borderColor: '#fff',
        borderWidth: 2
      },
      label: {
        show: false,
        position: 'center'
      },
      emphasis: {
        label: {
          show: true,
          fontSize: 20,
          fontWeight: 'bold'
        }
      },
      labelLine: {
        show: false
      },
      data: [
        { value: 1048, name: '测试APP', itemStyle: { color: '#409EFF' } },
        { value: 735, name: '演示APP', itemStyle: { color: '#67C23A' } },
        { value: 580, name: '测试模块关联', itemStyle: { color: '#E6A23C' } },
        { value: 484, name: '测试新APP', itemStyle: { color: '#F56C6C' } }
      ]
    }]
  }
  chart.setOption(option)
  window.addEventListener('resize', () => chart.resize())
}

onMounted(() => {
  loadStats()
  setTimeout(() => {
    initUserTrendChart()
    initAppUsageChart()
  }, 100)
})
</script>

<style scoped lang="scss">
.dashboard-container {
  .stat-card {
    .stat-content {
      display: flex;
      align-items: center;
      gap: 20px;

      .stat-icon {
        width: 60px;
        height: 60px;
        border-radius: 10px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: white;
      }

      .stat-info {
        flex: 1;

        .stat-value {
          font-size: 28px;
          font-weight: 600;
          color: #303133;
          margin-bottom: 8px;
        }

        .stat-label {
          font-size: 14px;
          color: #909399;
        }
      }
    }
  }

  .quick-actions {
    display: flex;
    gap: 12px;
    flex-wrap: wrap;
  }
}

// 移动端适配
@media (max-width: 768px) {
  .dashboard-container {
    .stat-card {
      .stat-content {
        gap: 12px;

        .stat-icon {
          width: 48px;
          height: 48px;

          .el-icon {
            font-size: 24px !important;
          }
        }

        .stat-info {
          .stat-value {
            font-size: 22px;
          }

          .stat-label {
            font-size: 13px;
          }
        }
      }
    }

    .quick-actions {
      .el-button {
        flex: 1;
        min-width: calc(50% - 6px);
      }
    }
  }

  // 图表高度调整
  .chart-container {
    height: 280px !important;
  }
}
</style>
