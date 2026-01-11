import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import App from './App.vue'
import router from './router'
import './styles/index.scss'

// 导入错误收集器
import errorCollector from './utils/errorCollector'

const app = createApp(App)

// Vue全局错误处理
app.config.errorHandler = (err, instance, info) => {
  console.error('[Vue Error]', err, info)
  errorCollector.collectError({
    type: 'vue_error',
    message: err.message,
    stack: err.stack,
    componentInfo: info,
    componentName: instance?.$options?.name || 'Unknown'
  })
}

// Vue警告处理（开发环境）
app.config.warnHandler = (msg, instance, trace) => {
  console.warn('[Vue Warning]', msg)
  // 警告不发送到Manus，只记录到控制台
}

app.use(createPinia())
app.use(router)
app.use(ElementPlus)
app.mount('#app')

// 页面卸载前发送剩余错误
window.addEventListener('beforeunload', () => {
  errorCollector.flush()
})

console.log('[App] Error collector initialized')
