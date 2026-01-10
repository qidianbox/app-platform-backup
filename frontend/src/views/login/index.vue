<template>
  <div class="login-container">
    <div class="login-left">
      <div class="brand">
        <div class="logo">
          <svg viewBox="0 0 48 48" fill="none" xmlns="http://www.w3.org/2000/svg">
            <rect width="48" height="48" rx="8" fill="white" fill-opacity="0.1"/>
            <path d="M24 8L36 14V26L24 32L12 26V14L24 8Z" stroke="white" stroke-width="2" fill="none"/>
            <path d="M24 20L30 23V29L24 32L18 29V23L24 20Z" fill="white"/>
          </svg>
        </div>
        <h1>APP中台管理系统</h1>
        <p>统一管理多个APP的后台系统</p>
      </div>
      <div class="features">
        <div class="feature-item">
          <div class="feature-icon">📊</div>
          <div class="feature-text">
            <h3>数据统计</h3>
            <p>实时监控应用数据</p>
          </div>
        </div>
        <div class="feature-item">
          <div class="feature-icon">🔐</div>
          <div class="feature-text">
            <h3>安全管理</h3>
            <p>多维度安全防护</p>
          </div>
        </div>
        <div class="feature-item">
          <div class="feature-icon">⚡</div>
          <div class="feature-text">
            <h3>高效运维</h3>
            <p>一站式运维管理</p>
          </div>
        </div>
      </div>
    </div>
    
    <div class="login-right">
      <div class="login-card">
        <div class="card-header">
          <h2>欢迎登录</h2>
          <p>请输入您的账号信息</p>
        </div>
        
        <el-form :model="form" @submit.prevent="handleLogin" class="login-form">
          <el-form-item>
            <label class="form-label">用户名</label>
            <el-input 
              v-model="form.username" 
              placeholder="请输入用户名" 
              size="large"
              :prefix-icon="User"
            />
          </el-form-item>
          <el-form-item>
            <label class="form-label">密码</label>
            <el-input 
              v-model="form.password" 
              type="password" 
              placeholder="请输入密码" 
              size="large"
              :prefix-icon="Lock"
              show-password
            />
          </el-form-item>
          <el-form-item>
            <el-button 
              type="primary" 
              size="large" 
              :loading="loading" 
              @click="handleLogin"
              class="login-btn"
            >
              {{ loading ? '登录中...' : '登 录' }}
            </el-button>
          </el-form-item>
        </el-form>
        
        <div class="login-footer">
          <p class="hint">默认账号：admin / admin123</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'
import request from '@/utils/request'

const router = useRouter()
const loading = ref(false)
const form = ref({
  username: 'admin',
  password: 'admin123'
})

const handleLogin = async () => {
  if (!form.value.username || !form.value.password) {
    ElMessage.warning('请输入用户名和密码')
    return
  }
  
  loading.value = true
  try {
    const res = await request.post('/admin/login', form.value)
    localStorage.setItem('token', res.token)
    localStorage.setItem('user', JSON.stringify(res.user))
    ElMessage.success('登录成功')
    router.push('/dashboard')
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}
</script>

<style lang="scss" scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  background: #f0f2f5;
}

.login-left {
  flex: 1;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 60px;
  color: white;
  position: relative;
  overflow: hidden;
  
  &::before {
    content: '';
    position: absolute;
    top: -50%;
    left: -50%;
    width: 200%;
    height: 200%;
    background: radial-gradient(circle, rgba(255,255,255,0.03) 0%, transparent 70%);
    animation: rotate 30s linear infinite;
  }
  
  @keyframes rotate {
    from { transform: rotate(0deg); }
    to { transform: rotate(360deg); }
  }
}

.brand {
  position: relative;
  z-index: 1;
  margin-bottom: 60px;
  
  .logo {
    width: 64px;
    height: 64px;
    margin-bottom: 24px;
    
    svg {
      width: 100%;
      height: 100%;
    }
  }
  
  h1 {
    font-size: 32px;
    font-weight: 600;
    margin: 0 0 12px;
    letter-spacing: 1px;
  }
  
  p {
    font-size: 16px;
    color: rgba(255, 255, 255, 0.7);
    margin: 0;
  }
}

.features {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  transition: all 0.3s ease;
  
  &:hover {
    background: rgba(255, 255, 255, 0.1);
    transform: translateX(8px);
  }
  
  .feature-icon {
    font-size: 28px;
    width: 48px;
    height: 48px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(255, 255, 255, 0.1);
    border-radius: 10px;
  }
  
  .feature-text {
    h3 {
      font-size: 16px;
      font-weight: 500;
      margin: 0 0 4px;
    }
    
    p {
      font-size: 13px;
      color: rgba(255, 255, 255, 0.6);
      margin: 0;
    }
  }
}

.login-right {
  width: 480px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px;
  background: white;
}

.login-card {
  width: 100%;
  max-width: 360px;
}

.card-header {
  margin-bottom: 32px;
  
  h2 {
    font-size: 28px;
    font-weight: 600;
    color: #1a1a2e;
    margin: 0 0 8px;
  }
  
  p {
    font-size: 14px;
    color: #8c8c8c;
    margin: 0;
  }
}

.login-form {
  .form-label {
    display: block;
    font-size: 14px;
    font-weight: 500;
    color: #333;
    margin-bottom: 8px;
  }
  
  :deep(.el-input) {
    .el-input__wrapper {
      padding: 0 15px;
      height: 48px;
      border-radius: 8px;
      box-shadow: 0 0 0 1px #e0e0e0;
      
      &:hover {
        box-shadow: 0 0 0 1px #1a1a2e;
      }
      
      &.is-focus {
        box-shadow: 0 0 0 2px rgba(26, 26, 46, 0.2);
      }
    }
    
    .el-input__inner {
      height: 48px;
      font-size: 15px;
    }
  }
  
  .el-form-item {
    margin-bottom: 24px;
  }
}

.login-btn {
  width: 100%;
  height: 48px;
  font-size: 16px;
  font-weight: 500;
  border-radius: 8px;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
  border: none;
  
  &:hover {
    background: linear-gradient(135deg, #16213e 0%, #0f3460 100%);
  }
}

.login-footer {
  margin-top: 24px;
  text-align: center;
  
  .hint {
    font-size: 13px;
    color: #8c8c8c;
    margin: 0;
  }
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .login-left {
    display: none;
  }
  
  .login-right {
    width: 100%;
    min-height: 100vh;
    background: linear-gradient(135deg, #f5f7fa 0%, #e4e8eb 100%);
  }
  
  .login-card {
    background: white;
    padding: 40px;
    border-radius: 16px;
    box-shadow: 0 4px 24px rgba(0, 0, 0, 0.08);
  }
}

@media (max-width: 480px) {
  .login-right {
    padding: 20px;
  }
  
  .login-card {
    padding: 30px 24px;
  }
  
  .card-header h2 {
    font-size: 24px;
  }
}
</style>
