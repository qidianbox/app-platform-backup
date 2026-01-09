<template>
  <div class="app-list-container">
    <div class="page-header">
      <h2 class="page-title">我的APP项目</h2>
      <el-input
        v-model="searchKeyword"
        placeholder="搜索APP..."
        prefix-icon="Search"
        clearable
        style="width: 300px; margin-right: 16px"
        @input="handleSearch"
      />
    </div>

    <div class="app-cards-container" v-loading="loading">
      <!-- 新建APP卡片 -->
      <el-card class="app-card create-card" shadow="hover" @click="handleCreate">
        <div class="create-card-content">
          <el-icon :size="48" color="#409EFF"><Plus /></el-icon>
          <div class="create-text">创建新APP</div>
        </div>
      </el-card>

      <!-- APP项目卡片 -->
      <el-card
        v-for="app in paginatedApps"
        :key="app.id"
        class="app-card"
        shadow="hover"
      >
        <div class="app-card-header">
          <div class="app-logo">
            <img v-if="app.logo" :src="app.logo" alt="logo" />
            <el-icon v-else :size="40" color="#409EFF"><Platform /></el-icon>
          </div>
          <div class="app-info">
            <h3 class="app-name">{{ app.app_name }}</h3>
            <div class="app-key" :title="app.app_key">{{ shortenAppKey(app.app_key) }}</div>
          </div>
          <el-tag :type="app.status === 1 ? 'success' : 'danger'" size="small">
            {{ app.status === 1 ? '运行中' : '已停用' }}
          </el-tag>
        </div>

        <div class="app-card-body">
          <div class="app-stats">
            <div class="stat-item">
              <el-icon><Grid /></el-icon>
              <span>{{ app.module_count || 0 }} 个模块</span>
            </div>
            <div class="stat-item">
              <el-icon><User /></el-icon>
              <span>{{ app.user_count || 0 }} 用户</span>
            </div>
          </div>
          <div class="app-description">
            {{ app.description || '暂无描述' }}
          </div>
        </div>

        <div class="app-card-footer">
          <el-button type="primary" size="small" @click="handleEnterApp(app)">
            <el-icon><Right /></el-icon>
            进入配置
          </el-button>
          <el-button-group size="small">
            <el-button @click="handleEdit(app)" title="编辑">
              <el-icon><Edit /></el-icon>
            </el-button>
            <el-button @click="handleManageModules(app)" title="管理模块">
              <el-icon><Grid /></el-icon>
            </el-button>
            <el-button @click="handleResetSecret(app)" title="重置密钥">
              <el-icon><Key /></el-icon>
            </el-button>
            <el-button @click="handleDelete(app)" type="danger" title="删除">
              <el-icon><Delete /></el-icon>
            </el-button>
          </el-button-group>
        </div>
      </el-card>
    </div>

    <!-- 分页 -->
    <div class="pagination-container" v-if="filteredApps.length > pageSize">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[12, 24, 48, 96]"
        :total="filteredApps.length"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="handlePageChange"
      />
    </div>

    <!-- 创建/编辑APP对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="900px"
      @close="handleDialogClose"
    >
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item label="APP名称" prop="app_name">
          <el-input v-model="formData.app_name" placeholder="请输入APP名称" />
        </el-form-item>
        <el-form-item label="APP标识" prop="app_key">
          <el-input v-model="formData.app_key" placeholder="请输入APP标识（英文）" />
        </el-form-item>
        <el-form-item label="包名" prop="package_name">
          <el-input v-model="formData.package_name" placeholder="com.example.app" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :rows="3"
            placeholder="请输入APP描述"
          />
        </el-form-item>
        <el-form-item label="选择模块" prop="modules" v-if="!isEdit">
          <div class="module-selection-container">
            <div class="module-selection-header">
              <span class="selected-count">已选择 {{ formData.modules.length }} 个模块</span>
              <el-button text type="primary" @click="handleSelectAll">全选</el-button>
              <el-button text @click="handleClearAll">清空</el-button>
            </div>
            
            <!-- 分组展示模块 -->
            <div v-for="group in groupedModuleTemplates" :key="group.id" class="module-group">
              <div v-if="group.modules.length > 0" class="module-group-header">
                <span class="group-icon">{{ group.icon }}</span>
                <span class="group-name">{{ group.name }}</span>
                <span class="group-desc">{{ group.description }}</span>
                <el-button text size="small" type="primary" @click="handleSelectGroup(group)">
                  全选本组
                </el-button>
              </div>
              <div class="module-grid">
                <div
                  v-for="module in group.modules"
                  :key="module.module_code"
                  class="module-card"
                  :class="{ 'is-selected': formData.modules.includes(module.module_code) }"
                  @click="toggleModule(module.module_code)"
                >
                  <div class="module-card-header">
                    <el-icon :size="32" :color="formData.modules.includes(module.module_code) ? '#409EFF' : '#909399'">
                      <component :is="module.icon" />
                    </el-icon>
                    <el-checkbox
                      :model-value="formData.modules.includes(module.module_code)"
                      @click.stop
                      @change="toggleModule(module.module_code)"
                    />
                  </div>
                  <div class="module-card-body">
                    <div class="module-name">{{ module.module_name }}</div>
                    <div class="module-desc">{{ module.description }}</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">
          {{ isEdit ? '保存' : '创建' }}
        </el-button>
      </template>
    </el-dialog>

    <!-- 模块管理对话框 -->
    <el-dialog v-model="moduleDialogVisible" title="管理模块" width="700px">
      <el-table :data="currentAppModules" border>
        <el-table-column prop="module_name" label="模块名称" width="150" />
        <el-table-column prop="description" label="描述" />
        <el-table-column label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.is_enabled"
              @change="handleToggleModule(row)"
            />
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Plus,
  Platform,
  Grid,
  User,
  Right,
  MoreFilled,
  Edit,
  Delete,
  Key,
  Search,
  FolderOpened,
  ChatDotRound,
  BellFilled,
  DataAnalysis,
  Tickets,
  Setting,
  Document,
  Monitor,
  Wallet
} from '@element-plus/icons-vue'
import { getAppList, createApp, updateApp, deleteApp, resetAppSecret } from '@/api/app'
import { getModuleTemplates, getAppModules, updateAppModule } from '@/api/module'
import { getGroupedModules } from '@/config/moduleCategories'

const router = useRouter()

// 数据
const loading = ref(false)
const appList = ref([])
const searchKeyword = ref('')
const moduleTemplates = ref([])
const groupedModuleTemplates = ref([])
const currentPage = ref(1)
const pageSize = ref(12)
const total = ref(0)

// 对话框
const dialogVisible = ref(false)
const dialogTitle = ref('')
const isEdit = ref(false)
const submitting = ref(false)
const formRef = ref(null)
const formData = reactive({
  app_name: '',
  app_key: '',
  package_name: '',
  description: '',
  modules: []
})
const formRules = {
  app_name: [{ required: true, message: '请输入APP名称', trigger: 'blur' }],
  app_key: [{ required: true, message: '请输入APP标识', trigger: 'blur' }],
  modules: [{ type: 'array', min: 1, message: '请至少选择一个模块', trigger: 'change' }]
}

// 模块管理对话框
const moduleDialogVisible = ref(false)
const currentApp = ref(null)
const currentAppModules = ref([])

// 缩短APP Key显示
const shortenAppKey = (appKey) => {
  if (!appKey) return ''
  // 取前8位 + ... + 后4位
  if (appKey.length <= 16) return appKey
  return `${appKey.slice(0, 8)}...${appKey.slice(-4)}`
}

// 计算属性
const filteredApps = computed(() => {
  if (!searchKeyword.value) return appList.value
  const keyword = searchKeyword.value.toLowerCase()
  return appList.value.filter(app =>
    app.app_name.toLowerCase().includes(keyword) ||
    app.app_key.toLowerCase().includes(keyword)
  )
})

const paginatedApps = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredApps.value.slice(start, end)
})

const handlePageChange = (page) => {
  currentPage.value = page
}

// 方法
const fetchAppList = async () => {
  loading.value = true
  try {
    const res = await getAppList()
    appList.value = res.data.list || []
  } catch (error) {
    ElMessage.error('获取APP列表失败')
  } finally {
    loading.value = false
  }
}

const fetchModuleTemplates = async () => {
  try {
    const res = await getModuleTemplates()
    moduleTemplates.value = res.data || []
    // 按分组整理模块
    groupedModuleTemplates.value = getGroupedModules(moduleTemplates.value)
  } catch (error) {
    ElMessage.error('获取模块模板失败')
  }
}

const handleSearch = () => {
  // 搜索在computed中实现
}

const handleCreate = () => {
  dialogTitle.value = '创建新APP'
  isEdit.value = false
  resetForm()
  dialogVisible.value = true
}

const handleEdit = (app) => {
  dialogTitle.value = '编辑APP'
  isEdit.value = true
  Object.assign(formData, {
    id: app.id,
    app_name: app.app_name,
    app_key: app.app_key,
    package_name: app.package_name,
    description: app.description
  })
  dialogVisible.value = true
}

const handleEnterApp = (app) => {
  router.push(`/app/${app.id}/dashboard`)
}

const handleCommand = (command, app) => {
  switch (command) {
    case 'edit':
      handleEdit(app)
      break
    case 'modules':
      handleManageModules(app)
      break
    case 'reset':
      handleResetSecret(app)
      break
    case 'delete':
      handleDelete(app)
      break
  }
}

const handleManageModules = async (app) => {
  currentApp.value = app
  try {
    const res = await getAppModules(app.id)
    currentAppModules.value = res.data || []
    moduleDialogVisible.value = true
  } catch (error) {
    ElMessage.error('获取模块列表失败')
  }
}

const handleToggleModule = async (module) => {
  try {
    await updateAppModule(currentApp.value.id, module.module_code, {
      is_enabled: module.is_enabled
    })
    ElMessage.success('模块状态更新成功')
    fetchAppList()
  } catch (error) {
    ElMessage.error('更新失败')
    module.is_enabled = !module.is_enabled
  }
}

const handleResetSecret = async (app) => {
  try {
    await ElMessageBox.confirm('确定要重置该APP的密钥吗？', '提示', {
      type: 'warning'
    })
    await resetAppSecret(app.id)
    ElMessage.success('密钥重置成功')
    fetchAppList()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('重置失败')
    }
  }
}

const handleDelete = async (app) => {
  try {
    await ElMessageBox.confirm(`确定要删除APP "${app.app_name}" 吗？`, '提示', {
      type: 'warning'
    })
    await deleteApp(app.id)
    ElMessage.success('删除成功')
    fetchAppList()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitting.value = true
    try {
      if (isEdit.value) {
        await updateApp(formData.id, formData)
        ElMessage.success('更新成功')
      } else {
        await createApp(formData)
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      fetchAppList()
    } catch (error) {
      ElMessage.error(isEdit.value ? '更新失败' : '创建失败')
    } finally {
      submitting.value = false
    }
  })
}

const handleDialogClose = () => {
  resetForm()
}

const resetForm = () => {
  formData.id = null
  formData.app_name = ''
  formData.app_key = ''
  formData.package_name = ''
  formData.description = ''
  formData.modules = []
  formRef.value?.resetFields()
}

// 模块选择方法
const toggleModule = (moduleCode) => {
  const index = formData.modules.indexOf(moduleCode)
  if (index > -1) {
    formData.modules.splice(index, 1)
  } else {
    formData.modules.push(moduleCode)
  }
}

const handleSelectAll = () => {
  formData.modules = moduleTemplates.value.map(m => m.module_code)
}

const handleClearAll = () => {
  formData.modules = []
}

const handleSelectGroup = (group) => {
  // 全选本组模块
  const groupModuleCodes = group.modules.map(m => m.module_code)
  groupModuleCodes.forEach(code => {
    if (!formData.modules.includes(code)) {
      formData.modules.push(code)
    }
  })
}

// 生命周期
onMounted(() => {
  fetchAppList()
  fetchModuleTemplates()
})
</script>

<style scoped lang="scss">
.app-list-container {
  padding: 24px;
  min-height: 100vh;
  background: #f5f7fa;

  @media (max-width: 768px) {
    padding: 12px;
  }
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 12px;

  .page-title {
    font-size: 24px;
    font-weight: 600;
    color: #303133;
    margin: 0;

    @media (max-width: 768px) {
      font-size: 20px;
    }
  }

  @media (max-width: 768px) {
    margin-bottom: 16px;
  }
}

.app-cards-container {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;

  @media (max-width: 768px) {
    grid-template-columns: 1fr;
    gap: 12px;
  }
}

.app-card {
  cursor: pointer;
  transition: all 0.3s;

  &:hover {
    transform: translateY(-4px);
  }

  .app-card-header {
    display: flex;
    align-items: flex-start;
    gap: 12px;
    margin-bottom: 16px;

    .app-logo {
      width: 50px;
      height: 50px;
      border-radius: 8px;
      background: #f0f2f5;
      display: flex;
      align-items: center;
      justify-content: center;

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
        border-radius: 8px;
      }
    }

    .app-info {
      flex: 1;
      min-width: 0;

      .app-name {
        font-size: 18px;
        font-weight: 600;
        color: #303133;
        margin: 0 0 4px 0;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }

      .app-key {
        font-size: 12px;
        color: #909399;
        font-family: monospace;
      }
    }
  }

  .app-card-body {
    margin-bottom: 16px;

    .app-stats {
      display: flex;
      gap: 16px;
      margin-bottom: 12px;

      .stat-item {
        display: flex;
        align-items: center;
        gap: 4px;
        font-size: 14px;
        color: #606266;

        .el-icon {
          color: #909399;
        }
      }
    }

    .app-description {
      font-size: 13px;
      color: #909399;
      line-height: 1.6;
      display: -webkit-box;
      -webkit-line-clamp: 2;
      -webkit-box-orient: vertical;
      overflow: hidden;
    }
  }

  .app-card-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-top: 16px;
    border-top: 1px solid #ebeef5;
  }
}

.create-card {
  border: 2px dashed #dcdfe6;
  background: #fafafa;

  &:hover {
    border-color: #409eff;
    background: #f0f9ff;
  }

  .create-card-content {
    height: 240px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 12px;

    .create-text {
      font-size: 16px;
      color: #606266;
      font-weight: 500;
    }
  }
}

.module-selection-container {
  width: 100%;

  .module-group {
    margin-bottom: 24px;

    &:last-child {
      margin-bottom: 0;
    }

    .module-group-header {
      display: flex;
      align-items: center;
      gap: 8px;
      padding: 12px 16px;
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      border-radius: 8px;
      margin-bottom: 16px;
      color: #fff;

      .group-icon {
        font-size: 20px;
      }

      .group-name {
        font-size: 15px;
        font-weight: 600;
      }

      .group-desc {
        flex: 1;
        font-size: 12px;
        opacity: 0.9;
        margin-left: 8px;
      }

      .el-button {
        color: #fff;
        
        &:hover {
          background: rgba(255, 255, 255, 0.2);
        }
      }
    }
  }

  .module-selection-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 16px;
    padding: 12px 16px;
    background: #f5f7fa;
    border-radius: 4px;

    .selected-count {
      font-size: 14px;
      font-weight: 500;
      color: #409EFF;
    }
  }

  .module-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 16px;
    max-height: 500px;
    overflow-y: auto;
    padding: 4px;

    @media (max-width: 1200px) {
      grid-template-columns: repeat(2, 1fr);
    }

    @media (max-width: 768px) {
      grid-template-columns: 1fr;
    }
  }

  .module-card {
    border: 2px solid #e4e7ed;
    border-radius: 8px;
    padding: 16px;
    cursor: pointer;
    transition: all 0.3s;
    background: #fff;

    &:hover {
      border-color: #409EFF;
      box-shadow: 0 2px 12px rgba(64, 158, 255, 0.15);
      transform: translateY(-2px);
    }

    &.is-selected {
      border-color: #409EFF;
      background: #f0f9ff;
      box-shadow: 0 2px 8px rgba(64, 158, 255, 0.2);
    }

    .module-card-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: 12px;
    }

    .module-card-body {
      .module-name {
        font-size: 15px;
        font-weight: 600;
        color: #303133;
        margin-bottom: 8px;
      }

      .module-desc {
        font-size: 13px;
        color: #606266;
        line-height: 1.5;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
      }
    }
  }
}
</style>

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 32px;
  padding: 20px 0;
}

// 移动端适配
@media (max-width: 768px) {
  .app-list-container {
    .page-header {
      flex-direction: column;
      align-items: flex-start;
      gap: 12px;

      .page-title {
        font-size: 18px;
      }

      .el-input {
        width: 100% !important;
        margin-right: 0 !important;
      }
    }

    .app-cards-container {
      grid-template-columns: 1fr;
    }
  }

  .app-card {
    .app-card-header {
      .app-logo {
        width: 40px;
        height: 40px;
      }

      .app-info {
        .app-name {
          font-size: 16px;
        }

        .app-key {
          font-size: 12px;
        }
      }
    }

    .app-card-body {
      .app-stats {
        flex-direction: column;
        align-items: flex-start;
        gap: 8px;

        .stat-item {
          font-size: 13px;
        }
      }

      .app-description {
        font-size: 13px;
      }
    }

    .app-card-footer {
      flex-direction: column;
      gap: 8px;

      .el-button {
        width: 100%;
      }

      .el-button-group {
        width: 100%;
        display: flex;

        .el-button {
          flex: 1;
        }
      }
    }
  }

  .create-card {
    .create-card-content {
      height: 150px;

      .el-icon {
        font-size: 36px !important;
      }

      .create-text {
        font-size: 14px;
      }
    }
  }

  // 对话框移动端适配
  :deep(.el-dialog) {
    width: 95% !important;
    max-width: 95% !important;
    margin: 10px auto !important;

    .el-dialog__body {
      padding: 16px;
    }

    .el-form-item__label {
      width: 100% !important;
      text-align: left;
      margin-bottom: 8px;
    }

    .el-form-item__content {
      margin-left: 0 !important;
    }
  }

  .module-selection-container {
    .module-group-header {
      flex-wrap: wrap;
      padding: 10px 12px;

      .group-desc {
        width: 100%;
        margin-left: 0;
        margin-top: 4px;
      }
    }

    .module-grid {
      grid-template-columns: 1fr !important;
      max-height: 300px;
    }

    .module-card {
      padding: 12px;

      .module-card-header {
        .el-icon {
          font-size: 24px !important;
        }
      }

      .module-card-body {
        .module-name {
          font-size: 14px;
        }

        .module-desc {
          font-size: 12px;
        }
      }
    }
  }

  .pagination-container {
    margin-top: 20px;
    padding: 12px 0;

    :deep(.el-pagination) {
      justify-content: center;
      flex-wrap: wrap;
      gap: 8px;

      .el-pagination__sizes {
        margin: 0;
      }

      .el-pagination__jump {
        margin-left: 0;
      }
    }
  }
}
