<template>
  <div class="mobile-menu-container">
    <!-- 汉堡菜单按钮 -->
    <button 
      class="hamburger-btn"
      :class="{ 'is-open': isOpen }"
      @click="toggleMenu"
      aria-label="打开菜单"
      aria-expanded="isOpen"
      aria-controls="mobile-drawer"
    >
      <span class="hamburger-line"></span>
      <span class="hamburger-line"></span>
      <span class="hamburger-line"></span>
    </button>

    <!-- 遮罩层 -->
    <Transition name="fade">
      <div 
        v-if="isOpen" 
        class="drawer-overlay"
        @click="closeMenu"
        aria-hidden="true"
      ></div>
    </Transition>

    <!-- 抽屉式侧边栏 -->
    <Transition name="slide">
      <nav 
        v-if="isOpen"
        id="mobile-drawer"
        class="drawer-sidebar"
        role="navigation"
        aria-label="移动端导航菜单"
      >
        <!-- 抽屉头部 -->
        <div class="drawer-header">
          <div class="drawer-logo">
            <div class="app-icon">
              <span>{{ logoText }}</span>
            </div>
            <span class="app-name">{{ appName }}</span>
          </div>
          <button 
            class="close-btn"
            @click="closeMenu"
            aria-label="关闭菜单"
          >
            <el-icon><Close /></el-icon>
          </button>
        </div>

        <!-- 菜单内容 -->
        <div class="drawer-content">
          <slot></slot>
        </div>

        <!-- 抽屉底部 -->
        <div class="drawer-footer">
          <slot name="footer"></slot>
        </div>
      </nav>
    </Transition>
  </div>
</template>

<script setup>
import { ref, watch, onMounted, onUnmounted } from 'vue'
import { Close } from '@element-plus/icons-vue'

const props = defineProps({
  logoText: {
    type: String,
    default: '拓'
  },
  appName: {
    type: String,
    default: '拓客APP中台'
  },
  modelValue: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'open', 'close'])

const isOpen = ref(props.modelValue)

// 同步外部v-model
watch(() => props.modelValue, (val) => {
  isOpen.value = val
})

watch(isOpen, (val) => {
  emit('update:modelValue', val)
  if (val) {
    emit('open')
    document.body.style.overflow = 'hidden'
  } else {
    emit('close')
    document.body.style.overflow = ''
  }
})

const toggleMenu = () => {
  isOpen.value = !isOpen.value
}

const closeMenu = () => {
  isOpen.value = false
}

// ESC键关闭菜单
const handleKeydown = (e) => {
  if (e.key === 'Escape' && isOpen.value) {
    closeMenu()
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
  document.body.style.overflow = ''
})

// 暴露方法给父组件
defineExpose({
  open: () => { isOpen.value = true },
  close: closeMenu,
  toggle: toggleMenu
})
</script>

<style scoped lang="scss">
.mobile-menu-container {
  display: none;
}

/* 移动端显示 */
@media (max-width: 768px) {
  .mobile-menu-container {
    display: block;
  }
}

/* 汉堡菜单按钮 */
.hamburger-btn {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  width: 44px;
  height: 44px;
  padding: 10px;
  background: transparent;
  border: none;
  cursor: pointer;
  gap: 5px;
  border-radius: 8px;
  transition: background 0.2s;

  &:hover {
    background: rgba(255, 255, 255, 0.1);
  }

  &:focus {
    outline: 2px solid #409eff;
    outline-offset: 2px;
  }
}

.hamburger-line {
  display: block;
  width: 22px;
  height: 2px;
  background: white;
  border-radius: 2px;
  transition: all 0.3s ease;
}

/* 汉堡菜单打开状态动画 */
.hamburger-btn.is-open {
  .hamburger-line:nth-child(1) {
    transform: translateY(7px) rotate(45deg);
  }
  .hamburger-line:nth-child(2) {
    opacity: 0;
  }
  .hamburger-line:nth-child(3) {
    transform: translateY(-7px) rotate(-45deg);
  }
}

/* 遮罩层 */
.drawer-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 999;
}

/* 抽屉侧边栏 */
.drawer-sidebar {
  position: fixed;
  top: 0;
  left: 0;
  width: 280px;
  max-width: 85vw;
  height: 100vh;
  background: white;
  z-index: 1000;
  display: flex;
  flex-direction: column;
  box-shadow: 4px 0 20px rgba(0, 0, 0, 0.15);
}

/* 抽屉头部 */
.drawer-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid #e4e7ed;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
}

.drawer-logo {
  display: flex;
  align-items: center;
  gap: 12px;
}

.app-icon {
  width: 36px;
  height: 36px;
  background: linear-gradient(135deg, #409eff 0%, #66b1ff 100%);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-weight: bold;
  font-size: 16px;
}

.app-name {
  font-size: 16px;
  font-weight: 600;
  color: white;
}

.close-btn {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.1);
  border: none;
  border-radius: 8px;
  color: white;
  cursor: pointer;
  transition: background 0.2s;

  &:hover {
    background: rgba(255, 255, 255, 0.2);
  }

  &:focus {
    outline: 2px solid #409eff;
    outline-offset: 2px;
  }

  .el-icon {
    font-size: 20px;
  }
}

/* 抽屉内容 */
.drawer-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px 0;
}

/* 抽屉底部 */
.drawer-footer {
  border-top: 1px solid #e4e7ed;
  padding: 12px 0;
}

/* 过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.slide-enter-active,
.slide-leave-active {
  transition: transform 0.3s ease;
}

.slide-enter-from,
.slide-leave-to {
  transform: translateX(-100%);
}
</style>
