/**
 * 前端错误自动收集器
 * 自动收集JS错误、API错误、Promise拒绝等，并发送到Manus通知系统
 */

// Manus通知API配置
const MANUS_NOTIFY_CONFIG = {
  // 使用后端代理发送通知，避免跨域问题
  endpoint: '/api/v1/system/error-report',
  // 错误收集开关
  enabled: true,
  // 错误去重时间窗口（毫秒）
  dedupeWindow: 60000,
  // 批量发送间隔（毫秒）
  batchInterval: 5000,
  // 最大批量大小
  maxBatchSize: 10
}

class ErrorCollector {
  constructor() {
    this.errors = []
    this.sentErrors = new Map() // 用于去重
    this.batchTimer = null
    this.isInitialized = false
  }

  /**
   * 初始化错误收集器
   */
  init() {
    if (this.isInitialized) return
    this.isInitialized = true

    // 捕获全局JS错误
    window.addEventListener('error', (event) => {
      this.collectError({
        type: 'js_error',
        message: event.message,
        filename: event.filename,
        lineno: event.lineno,
        colno: event.colno,
        stack: event.error?.stack
      })
    })

    // 捕获未处理的Promise拒绝
    window.addEventListener('unhandledrejection', (event) => {
      this.collectError({
        type: 'promise_rejection',
        message: String(event.reason),
        stack: event.reason?.stack
      })
    })

    // 捕获资源加载错误
    window.addEventListener('error', (event) => {
      if (event.target && (event.target.tagName === 'IMG' || event.target.tagName === 'SCRIPT' || event.target.tagName === 'LINK')) {
        this.collectError({
          type: 'resource_error',
          message: `Failed to load ${event.target.tagName.toLowerCase()}`,
          url: event.target.src || event.target.href
        })
      }
    }, true)

    // 监控控制台错误
    const originalConsoleError = console.error
    console.error = (...args) => {
      this.collectError({
        type: 'console_error',
        message: args.map(arg => typeof arg === 'object' ? JSON.stringify(arg) : String(arg)).join(' ')
      })
      originalConsoleError.apply(console, args)
    }

    console.log('[ErrorCollector] Initialized')
  }

  /**
   * 收集API错误
   */
  collectApiError(error) {
    this.collectError({
      type: 'api_error',
      method: error.config?.method?.toUpperCase(),
      url: error.config?.url,
      status: error.response?.status,
      statusText: error.response?.statusText,
      message: error.response?.data?.message || error.message,
      errorCode: error.response?.data?.code,
      requestData: this.sanitizeData(error.config?.data),
      responseData: error.response?.data
    })
  }

  /**
   * 收集错误
   */
  collectError(errorInfo) {
    if (!MANUS_NOTIFY_CONFIG.enabled) return

    const error = {
      id: this.generateId(),
      timestamp: new Date().toISOString(),
      url: window.location.href,
      userAgent: navigator.userAgent,
      ...errorInfo
    }

    // 去重检查
    const errorKey = this.getErrorKey(error)
    const lastSent = this.sentErrors.get(errorKey)
    if (lastSent && Date.now() - lastSent < MANUS_NOTIFY_CONFIG.dedupeWindow) {
      return // 跳过重复错误
    }

    this.errors.push(error)
    this.sentErrors.set(errorKey, Date.now())

    // 启动批量发送定时器
    if (!this.batchTimer) {
      this.batchTimer = setTimeout(() => this.sendBatch(), MANUS_NOTIFY_CONFIG.batchInterval)
    }

    // 如果错误数量达到阈值，立即发送
    if (this.errors.length >= MANUS_NOTIFY_CONFIG.maxBatchSize) {
      this.sendBatch()
    }
  }

  /**
   * 批量发送错误到Manus
   */
  async sendBatch() {
    if (this.batchTimer) {
      clearTimeout(this.batchTimer)
      this.batchTimer = null
    }

    if (this.errors.length === 0) return

    const errorsToSend = this.errors.splice(0, MANUS_NOTIFY_CONFIG.maxBatchSize)
    
    try {
      const response = await fetch(MANUS_NOTIFY_CONFIG.endpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('token') || ''}`
        },
        body: JSON.stringify({
          errors: errorsToSend,
          metadata: {
            appName: 'APP中台管理系统',
            environment: window.location.hostname.includes('localhost') ? 'development' : 'production',
            timestamp: new Date().toISOString(),
            totalErrors: errorsToSend.length
          }
        })
      })

      if (!response.ok) {
        console.warn('[ErrorCollector] Failed to send errors:', response.status)
        // 发送失败，将错误放回队列
        this.errors.unshift(...errorsToSend)
      } else {
        console.log(`[ErrorCollector] Sent ${errorsToSend.length} errors to Manus`)
      }
    } catch (e) {
      console.warn('[ErrorCollector] Error sending to Manus:', e)
      // 网络错误，将错误放回队列
      this.errors.unshift(...errorsToSend)
    }

    // 如果还有待发送的错误，继续定时发送
    if (this.errors.length > 0) {
      this.batchTimer = setTimeout(() => this.sendBatch(), MANUS_NOTIFY_CONFIG.batchInterval)
    }
  }

  /**
   * 生成唯一ID
   */
  generateId() {
    return Date.now().toString(36) + Math.random().toString(36).substr(2, 9)
  }

  /**
   * 生成错误去重键
   */
  getErrorKey(error) {
    return `${error.type}:${error.message}:${error.url || ''}:${error.filename || ''}`
  }

  /**
   * 脱敏数据
   */
  sanitizeData(data) {
    if (!data) return data
    try {
      const parsed = typeof data === 'string' ? JSON.parse(data) : data
      const sanitized = { ...parsed }
      // 脱敏敏感字段
      const sensitiveFields = ['password', 'token', 'secret', 'key', 'authorization']
      for (const field of sensitiveFields) {
        if (sanitized[field]) {
          sanitized[field] = '[REDACTED]'
        }
      }
      return sanitized
    } catch {
      return data
    }
  }

  /**
   * 手动报告错误
   */
  report(message, data = {}) {
    this.collectError({
      type: 'manual_report',
      message,
      ...data
    })
  }

  /**
   * 获取收集的错误列表
   */
  getErrors() {
    return [...this.errors]
  }

  /**
   * 清空错误队列
   */
  clear() {
    this.errors = []
    this.sentErrors.clear()
  }

  /**
   * 立即发送所有待发送的错误
   */
  flush() {
    return this.sendBatch()
  }
}

// 创建单例实例
const errorCollector = new ErrorCollector()

// 自动初始化
if (typeof window !== 'undefined') {
  errorCollector.init()
}

// 暴露到全局，方便调试
window.__errorCollector = errorCollector

export default errorCollector
export { ErrorCollector }
