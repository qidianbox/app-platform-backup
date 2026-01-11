<template>
  <div class="layout">
    <header class="header">
      <div class="logo">
        <el-icon><Grid /></el-icon>
        <span>APP中台</span>
      </div>
      <div class="user-info">
        <el-dropdown>
          <span class="user-name">
            {{ user?.nickname || 'Admin' }}
            <el-icon><ArrowDown /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="handleLogout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </header>
    
    <nav class="nav">
      <router-link to="/dashboard" class="nav-item" :class="{ active: $route.path === '/dashboard' }">
        <el-icon><DataAnalysis /></el-icon>
        <span>仪表盘</span>
      </router-link>
      <router-link to="/apps" class="nav-item" :class="{ active: $route.path.startsWith('/apps') }">
        <el-icon><Grid /></el-icon>
        <span>APP管理</span>
      </router-link>
      <router-link to="/modules" class="nav-item" :class="{ active: $route.path === '/modules' }">
        <el-icon><Menu /></el-icon>
        <span>模块管理</span>
      </router-link>
      <router-link to="/system/audit" class="nav-item" :class="{ active: $route.path === '/system/audit' }">
        <el-icon><Document /></el-icon>
        <span>审计日志</span>
      </router-link>
    </nav>
    
    <main class="main">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Grid, ArrowDown, DataAnalysis, Menu, Document } from '@element-plus/icons-vue'

const router = useRouter()
const user = ref(null)

onMounted(() => {
  const userData = localStorage.getItem('user')
  if (userData) {
    user.value = JSON.parse(userData)
  }
})

const handleLogout = () => {
  localStorage.removeItem('token')
  localStorage.removeItem('user')
  router.push('/login')
}
</script>

<style lang="scss" scoped>
.layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  height: 60px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  color: white;
  
  .logo {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 18px;
    font-weight: bold;
  }
  
  .user-name {
    display: flex;
    align-items: center;
    gap: 4px;
    cursor: pointer;
    color: white;
  }
}

.nav {
  display: flex;
  background: white;
  border-bottom: 1px solid #eee;
  overflow-x: auto;
  
  .nav-item {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 12px 20px;
    color: #666;
    text-decoration: none;
    white-space: nowrap;
    border-bottom: 2px solid transparent;
    
    &.active {
      color: #667eea;
      border-bottom-color: #667eea;
    }
  }
}

.main {
  flex: 1;
  background: #f5f7fa;
  padding: 20px;
  overflow-y: auto;
}

@media (max-width: 768px) {
  .header {
    padding: 0 15px;
    
    .logo span {
      display: none;
    }
  }
  
  .nav .nav-item {
    padding: 10px 15px;
    font-size: 14px;
  }
  
  .main {
    padding: 15px;
  }
}
</style>
