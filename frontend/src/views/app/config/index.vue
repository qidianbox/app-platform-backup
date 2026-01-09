<template>
  <div class="app-config">
    <div class="config-header">
      <el-button text @click="$router.push('/apps')">
        <el-icon><ArrowLeft /></el-icon>
        返回APP列表
      </el-button>
    </div>
    
    <div class="config-nav">
      <router-link 
        v-for="tab in tabs" 
        :key="tab.path" 
        :to="`/apps/${appId}/config/${tab.path}`"
        class="nav-tab"
        :class="{ active: $route.path.includes(tab.path) }"
      >
        <el-icon><component :is="tab.icon" /></el-icon>
        <span>{{ tab.name }}</span>
      </router-link>
    </div>
    
    <div class="config-content">
      <router-view />
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { ArrowLeft, Grid, User, Message, Wallet, DataAnalysis, Lock, Connection, Menu } from '@element-plus/icons-vue'

const route = useRoute()
const appId = computed(() => route.params.id)

const tabs = [
  { path: 'dashboard', name: '概览', icon: Grid },
  { path: 'user', name: '用户中心', icon: User },
  { path: 'message', name: '消息推送', icon: Message },
  { path: 'payment', name: '支付', icon: Wallet },
  { path: 'analytics', name: '数据统计', icon: DataAnalysis },
  { path: 'security', name: '安全', icon: Lock },
  { path: 'version', name: '版本', icon: Connection },
  { path: 'modules', name: '模块', icon: Menu }
]
</script>

<style lang="scss" scoped>
.app-config {
  max-width: 1200px;
  margin: 0 auto;
}

.config-header {
  margin-bottom: 15px;
}

.config-nav {
  display: flex;
  background: white;
  border-radius: 8px;
  padding: 5px;
  margin-bottom: 20px;
  overflow-x: auto;
  
  .nav-tab {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 10px 16px;
    border-radius: 6px;
    color: #666;
    text-decoration: none;
    white-space: nowrap;
    transition: all 0.3s;
    
    &:hover {
      background: #f5f7fa;
    }
    
    &.active {
      background: #409eff;
      color: white;
    }
  }
}

.config-content {
  background: white;
  border-radius: 8px;
  padding: 20px;
  min-height: 400px;
}

@media (max-width: 768px) {
  .config-nav {
    .nav-tab {
      padding: 8px 12px;
      font-size: 13px;
      
      span {
        display: none;
      }
    }
  }
  
  .config-content {
    padding: 15px;
  }
}
</style>
