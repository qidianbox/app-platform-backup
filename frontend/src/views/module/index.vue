<template>
  <div class="module-page">
    <div class="page-header">
      <h2>模块管理</h2>
      <p>管理系统中所有可用的功能模块</p>
    </div>
    
    <div class="module-grid">
      <el-card v-for="mod in modules" :key="mod.code" class="module-card">
        <div class="module-header">
          <div class="module-icon" :style="{ background: getModuleColor(mod.code) }">
            <el-icon><component :is="getModuleIcon(mod.icon)" /></el-icon>
          </div>
          <div class="module-info">
            <h3>{{ mod.name }}</h3>
            <p>{{ mod.description }}</p>
          </div>
        </div>
        
        <div class="module-functions">
          <div class="function-title">功能列表 ({{ mod.functions?.length || 0 }})</div>
          <div class="function-list">
            <el-tag 
              v-for="fn in mod.functions" 
              :key="fn.code" 
              :type="fn.type === 'active' ? 'primary' : 'info'"
              size="small"
            >
              {{ fn.name }}
            </el-tag>
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { User, Message, Bell, DataAnalysis, Document, Monitor, Folder, Setting, Connection } from '@element-plus/icons-vue'
import request from '@/utils/request'

const modules = ref([])

const colors = ['#409eff', '#67c23a', '#e6a23c', '#f56c6c', '#909399', '#9c27b0', '#00bcd4', '#ff5722', '#795548']

const getModuleColor = (code) => {
  const index = code?.charCodeAt(0) % colors.length || 0
  return colors[index]
}

const iconMap = {
  user: User,
  message: Message,
  bell: Bell,
  chart: DataAnalysis,
  'file-text': Document,
  monitor: Monitor,
  folder: Folder,
  settings: Setting,
  'git-branch': Connection
}

const getModuleIcon = (icon) => iconMap[icon] || Setting

const loadModules = async () => {
  try {
    const res = await request.get('/system/modules')
    modules.value = res.modules || []
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  loadModules()
})
</script>

<style lang="scss" scoped>
.module-page {
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 20px;
  
  h2 {
    margin: 0 0 5px;
  }
  
  p {
    margin: 0;
    color: #999;
  }
}

.module-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 20px;
}

.module-card {
  .module-header {
    display: flex;
    gap: 15px;
    margin-bottom: 15px;
    
    .module-icon {
      width: 50px;
      height: 50px;
      border-radius: 10px;
      display: flex;
      align-items: center;
      justify-content: center;
      color: white;
      font-size: 24px;
      flex-shrink: 0;
    }
    
    .module-info {
      flex: 1;
      
      h3 {
        margin: 0 0 5px;
        font-size: 16px;
      }
      
      p {
        margin: 0;
        color: #999;
        font-size: 13px;
      }
    }
  }
  
  .module-functions {
    .function-title {
      font-size: 13px;
      color: #666;
      margin-bottom: 10px;
    }
    
    .function-list {
      display: flex;
      flex-wrap: wrap;
      gap: 8px;
    }
  }
}

@media (max-width: 768px) {
  .module-grid {
    grid-template-columns: 1fr;
  }
}
</style>
