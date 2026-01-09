<template>
  <div class="dashboard-container">
    <h2 class="page-title">APP概览</h2>

    <el-row :gutter="20" class="stat-row">
      <el-col :xs="12" :sm="12" :md="6" :lg="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <el-icon :size="40" color="#409EFF" class="stat-icon"><User /></el-icon>
            <div class="stat-info">
              <div class="stat-value">{{ stats.userCount }}</div>
              <div class="stat-label">用户数</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :xs="12" :sm="12" :md="6" :lg="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <el-icon :size="40" color="#67C23A" class="stat-icon"><Grid /></el-icon>
            <div class="stat-info">
              <div class="stat-value">{{ stats.moduleCount }}</div>
              <div class="stat-label">启用模块</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :xs="12" :sm="12" :md="6" :lg="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <el-icon :size="40" color="#E6A23C" class="stat-icon"><DataLine /></el-icon>
            <div class="stat-info">
              <div class="stat-value">{{ stats.requestCount }}</div>
              <div class="stat-label">今日请求</div>
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :xs="12" :sm="12" :md="6" :lg="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <el-icon :size="40" color="#F56C6C" class="stat-icon"><Warning /></el-icon>
            <div class="stat-info">
              <div class="stat-value">{{ stats.errorCount }}</div>
              <div class="stat-label">今日异常</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :xs="24" :sm="24" :md="12" :lg="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>APP信息</span>
            </div>
          </template>
          <el-descriptions :column="1" border>
            <el-descriptions-item label="APP名称">
              {{ app?.app_name }}
            </el-descriptions-item>
            <el-descriptions-item label="APP标识">
              {{ app?.app_key }}
            </el-descriptions-item>
            <el-descriptions-item label="AppSecret">
              <span class="secret-text">{{ app?.app_secret }}</span>
              <el-button size="small" text @click="handleCopySecret">
                <el-icon><CopyDocument /></el-icon>
              </el-button>
            </el-descriptions-item>
            <el-descriptions-item label="包名">
              {{ app?.package_name || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="app?.status === 1 ? 'success' : 'danger'">
                {{ app?.status === 1 ? '运行中' : '已停用' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="创建时间">
              {{ app?.created_at }}
            </el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>

      <el-col :xs="24" :sm="24" :md="12" :lg="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>已启用模块</span>
            </div>
          </template>
          <div class="module-list">
            <el-tag
              v-for="module in enabledModules"
              :key="module.module_code"
              size="large"
              style="margin: 4px"
            >
              {{ module.module_name }}
            </el-tag>
            <el-empty v-if="enabledModules.length === 0" description="暂无启用的模块" />
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { User, Grid, DataLine, Warning, CopyDocument } from '@element-plus/icons-vue'
import { getAppModules } from '@/api/module'

const props = defineProps({
  appId: {
    type: Number,
    required: true
  },
  app: {
    type: Object,
    default: () => ({})
  }
})

const stats = ref({
  userCount: 0,
  moduleCount: 0,
  requestCount: 0,
  errorCount: 0
})

const enabledModules = ref([])

const fetchEnabledModules = async () => {
  try {
    const res = await getAppModules(props.appId, true)
    enabledModules.value = res.data || []
    stats.value.moduleCount = enabledModules.value.length
  } catch (error) {
    ElMessage.error('获取模块列表失败')
  }
}

const handleCopySecret = () => {
  if (props.app?.app_secret) {
    navigator.clipboard.writeText(props.app.app_secret)
    ElMessage.success('已复制到剪贴板')
  }
}

onMounted(() => {
  fetchEnabledModules()
})
</script>

<style scoped lang="scss">
.dashboard-container {
  .page-title {
    font-size: 24px;
    font-weight: 600;
    color: #303133;
    margin: 0 0 24px 0;
  }

  .stat-card {
    .stat-content {
      display: flex;
      align-items: center;
      gap: 16px;

      .stat-info {
        flex: 1;

        .stat-value {
          font-size: 28px;
          font-weight: 600;
          color: #303133;
          margin-bottom: 4px;
        }

        .stat-label {
          font-size: 14px;
          color: #909399;
        }
      }
    }
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .secret-text {
    font-family: monospace;
    font-size: 12px;
    color: #606266;
  }

  .module-list {
    min-height: 200px;
  }
}

// 移动端适配
@media (max-width: 768px) {
  .dashboard-container {
    .page-title {
      font-size: 18px;
      margin-bottom: 16px;
    }

    .stat-row {
      .el-col {
        margin-bottom: 12px;
      }
    }

    .stat-card {
      .stat-content {
        gap: 10px;
        flex-direction: column;
        text-align: center;
        padding: 8px 0;

        .stat-icon {
          font-size: 28px !important;
        }

        .stat-info {
          .stat-value {
            font-size: 20px;
          }

          .stat-label {
            font-size: 12px;
          }
        }
      }
    }

    .secret-text {
      word-break: break-all;
      font-size: 11px;
    }

    .module-list {
      min-height: 100px;
    }
  }
}
</style>
