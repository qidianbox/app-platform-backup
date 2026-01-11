import axios from 'axios'
import { ElMessage, ElNotification } from 'element-plus'
import { ErrorCodes, getErrorMessage, isAuthError, isRetryableError, isRateLimitError } from './errorCodes'
import errorCollector from './errorCollector'

// 根据环境自动选择API地址
const getBaseURL = () => {
  const hostname = window.location.hostname
  const origin = window.location.origin
  
  // 如果是通过代理访问（开发模式），使用相对路径
  if (hostname === 'localhost' || hostname === '127.0.0.1') {
    return '/api/v1'
  }
  
  // Manus沙箱环境：将5173端口替换为8080端口
  // 格式: https://5173-xxx.manus.computer -> https://8080-xxx.manus.computer
  if (hostname.includes('.manus.computer')) {
    const apiHost = origin.replace(/5173-/, '8080-').replace(/5174-/, '8080-')
    return `${apiHost}/api/v1`
  }
  
  // 其他生产环境，假设API和前端同域或使用环境变量
  const apiBaseUrl = window.__API_BASE_URL__ || origin.replace(/:\d+$/, ':8080')
  return `${apiBaseUrl}/api/v1`
}

// 重试配置
const RETRY_CONFIG = {
  maxRetries: 3,           // 最大重试次数
  retryDelay: 1000,        // 初始重试延迟（毫秒）
  maxRetryDelay: 10000,    // 最大重试延迟
  backoffMultiplier: 2,    // 退避倍数
  retryStatusCodes: [408, 429, 500, 502, 503, 504]  // 可重试的HTTP状态码
}

// 计算重试延迟（指数退避）
const getRetryDelay = (retryCount) => {
  const delay = RETRY_CONFIG.retryDelay * Math.pow(RETRY_CONFIG.backoffMultiplier, retryCount)
  return Math.min(delay, RETRY_CONFIG.maxRetryDelay)
}

// 延迟函数
const delay = (ms) => new Promise(resolve => setTimeout(resolve, ms))

// 系统日志收集器
const systemLogger = {
  logs: [],
  errors: [],
  requests: [],
  maxLogs: 200,
  
  // 记录通用日志
  log(level, module, message, data = null) {
    const log = {
      id: Date.now() + Math.random().toString(36).substr(2, 9),
      timestamp: new Date().toISOString(),
      level,
      module,
      message,
      data,
      url: window.location.href
    }
    this.logs.push(log)
    if (this.logs.length > this.maxLogs) {
      this.logs.shift()
    }
    
    // 控制台输出
    const consoleMethod = level === 'error' ? 'error' : level === 'warn' ? 'warn' : 'log'
    console[consoleMethod](`[${level.toUpperCase()}] [${module}] ${message}`, data || '')
    
    return log
  },
  
  info(module, message, data) {
    return this.log('info', module, message, data)
  },
  
  warn(module, message, data) {
    return this.log('warn', module, message, data)
  },
  
  error(module, message, data) {
    const log = this.log('error', module, message, data)
    this.errors.push(log)
    if (this.errors.length > this.maxLogs) {
      this.errors.shift()
    }
    return log
  },
  
  // 记录API请求
  logRequest(config) {
    const log = {
      id: Date.now() + Math.random().toString(36).substr(2, 9),
      timestamp: new Date().toISOString(),
      type: 'request',
      method: config.method?.toUpperCase(),
      url: config.url,
      fullUrl: config.baseURL + config.url,
      params: config.params,
      data: config.data,
      retryCount: config._retryCount || 0,
      headers: {
        Authorization: config.headers?.Authorization ? '[REDACTED]' : undefined,
        'Content-Type': config.headers?.['Content-Type']
      }
    }
    this.requests.push(log)
    if (this.requests.length > this.maxLogs) {
      this.requests.shift()
    }
    const retryInfo = log.retryCount > 0 ? ` (重试 #${log.retryCount})` : ''
    this.info('API', `Request: ${log.method} ${log.url}${retryInfo}`, { params: log.params })
    return log
  },
  
  // 记录API响应
  logResponse(response, requestLog) {
    const log = {
      id: Date.now() + Math.random().toString(36).substr(2, 9),
      timestamp: new Date().toISOString(),
      type: 'response',
      method: response.config?.method?.toUpperCase(),
      url: response.config?.url,
      status: response.status,
      statusText: response.statusText,
      code: response.data?.code,
      duration: requestLog ? Date.now() - new Date(requestLog.timestamp).getTime() : null,
      dataSize: JSON.stringify(response.data || {}).length
    }
    this.info('API', `Response: ${log.method} ${log.url} - ${log.status} [code:${log.code}] (${log.duration}ms)`, {
      status: log.status,
      code: log.code,
      dataSize: log.dataSize
    })
    return log
  },
  
  // 记录API错误
  logApiError(error) {
    const log = {
      id: Date.now() + Math.random().toString(36).substr(2, 9),
      timestamp: new Date().toISOString(),
      type: 'api_error',
      method: error.config?.method?.toUpperCase(),
      url: error.config?.url,
      status: error.response?.status,
      statusText: error.response?.statusText,
      message: error.response?.data?.message || error.message,
      errorCode: error.response?.data?.code,
      responseData: error.response?.data,
      retryCount: error.config?._retryCount || 0,
      stack: error.stack
    }
    this.error('API', `Error: ${log.method} ${log.url} - ${log.status}: ${log.message}`, log)
    return log
  },
  
  // 获取最近的错误日志
  getRecentErrors(count = 20) {
    return this.errors.slice(-count)
  },
  
  // 获取最近的请求日志
  getRecentRequests(count = 20) {
    return this.requests.slice(-count)
  },
  
  // 获取所有日志
  getAllLogs() {
    return {
      logs: this.logs,
      errors: this.errors,
      requests: this.requests
    }
  },
  
  // 导出日志为JSON
  exportLogs() {
    const data = {
      exportTime: new Date().toISOString(),
      userAgent: navigator.userAgent,
      url: window.location.href,
      ...this.getAllLogs()
    }
    return JSON.stringify(data, null, 2)
  },
  
  // 清空日志
  clear() {
    this.logs = []
    this.errors = []
    this.requests = []
    this.info('System', 'Logs cleared')
  }
}

// 将日志收集器暴露到全局，方便调试
window.__systemLogger = systemLogger
window.__exportLogs = () => {
  const data = systemLogger.exportLogs()
  const blob = new Blob([data], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `system-logs-${new Date().toISOString().replace(/[:.]/g, '-')}.json`
  a.click()
  URL.revokeObjectURL(url)
  console.log('Logs exported successfully')
}

// 全局错误捕获
window.addEventListener('error', (event) => {
  systemLogger.error('Global', `Uncaught error: ${event.message}`, {
    filename: event.filename,
    lineno: event.lineno,
    colno: event.colno,
    error: event.error?.stack
  })
})

window.addEventListener('unhandledrejection', (event) => {
  systemLogger.error('Global', `Unhandled promise rejection: ${event.reason}`, {
    reason: event.reason?.stack || event.reason
  })
})

// 记录页面加载
systemLogger.info('System', 'Application initialized', {
  baseURL: getBaseURL(),
  userAgent: navigator.userAgent,
  timestamp: new Date().toISOString()
})

const request = axios.create({
  baseURL: getBaseURL(),
  timeout: 30000
})

systemLogger.info('API', `API client created with base URL: ${getBaseURL()}`)

// 请求拦截器
request.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  config._requestLog = systemLogger.logRequest(config)
  config._startTime = Date.now()
  return config
}, error => {
  systemLogger.error('API', 'Request interceptor error', error)
  return Promise.reject(error)
})

// 响应拦截器
request.interceptors.response.use(
  response => {
    systemLogger.logResponse(response, response.config._requestLog)
    
    // 处理统一响应格式
    const data = response.data
    
    // 如果响应包含code字段，检查是否成功
    if (data && typeof data.code !== 'undefined') {
      if (data.code === ErrorCodes.SUCCESS) {
        // 成功，返回data字段
        return data.data !== undefined ? data.data : data
      } else {
        // 业务错误，根据错误码处理
        const errorCode = data.code
        const errorMessage = data.message || getErrorMessage(errorCode)
        
        // 认证错误处理
        if (isAuthError(errorCode)) {
          handleAuthError(errorCode, errorMessage)
          return Promise.reject(new Error(errorMessage))
        }
        
        // 限流错误处理
        if (isRateLimitError(errorCode)) {
          ElMessage.warning({
            message: '请求过于频繁，请稍后重试',
            duration: 5000
          })
          return Promise.reject(new Error(errorMessage))
        }
        
        // 其他业务错误
        ElMessage.error(errorMessage)
        return Promise.reject(new Error(errorMessage))
      }
    }
    
    // 兼容旧格式响应
    return data
  },
  async error => {
    const config = error.config
    
    // 检查是否可以重试
    if (shouldRetry(error, config)) {
      config._retryCount = (config._retryCount || 0) + 1
      const retryDelay = getRetryDelay(config._retryCount - 1)
      
      systemLogger.warn('API', `Retrying request (${config._retryCount}/${RETRY_CONFIG.maxRetries}) after ${retryDelay}ms`, {
        url: config.url,
        status: error.response?.status,
        retryCount: config._retryCount
      })
      
      // 显示重试提示
      ElMessage.info({
        message: `网络异常，正在重试 (${config._retryCount}/${RETRY_CONFIG.maxRetries})...`,
        duration: retryDelay
      })
      
      await delay(retryDelay)
      return request(config)
    }
    
    const errorLog = systemLogger.logApiError(error)
    
    // 发送到错误收集器
    errorCollector.collectApiError(error)
    
    // 详细的错误处理
    let errorMessage = '请求失败'
    let showNotification = false
    
    if (error.response) {
      const status = error.response.status
      const data = error.response.data
      const errorCode = data?.code
      
      // 优先使用后端返回的错误消息
      if (data?.message) {
        errorMessage = data.message
      } else if (errorCode) {
        errorMessage = getErrorMessage(errorCode)
      }
      
      // 根据状态码处理
      if (status === 401 || isAuthError(errorCode)) {
        handleAuthError(errorCode || ErrorCodes.UNAUTHORIZED, errorMessage)
        return Promise.reject(error)
      } else if (status === 403) {
        errorMessage = errorMessage || '没有权限执行此操作'
      } else if (status === 404) {
        errorMessage = errorMessage || 'API接口不存在'
      } else if (status === 429) {
        errorMessage = '请求过于频繁，请稍后重试'
        ElMessage.warning({
          message: errorMessage,
          duration: 5000
        })
        return Promise.reject(error)
      } else if (status >= 500) {
        errorMessage = errorMessage || '服务器内部错误'
        showNotification = true
      }
    } else if (error.code === 'ECONNABORTED') {
      errorMessage = '请求超时，请检查网络连接'
      showNotification = true
    } else if (!navigator.onLine) {
      errorMessage = '网络连接已断开'
      showNotification = true
    }
    
    // 对于严重错误，显示通知而不是简单的消息
    if (showNotification) {
      ElNotification({
        title: '系统错误',
        message: `${errorMessage}\n\n错误ID: ${errorLog.id}\n\n您可以在控制台输入 __exportLogs() 导出日志`,
        type: 'error',
        duration: 10000
      })
    } else {
      ElMessage.error(errorMessage)
    }
    
    return Promise.reject(error)
  }
)

// 判断是否应该重试
function shouldRetry(error, config) {
  // 已达到最大重试次数
  if ((config._retryCount || 0) >= RETRY_CONFIG.maxRetries) {
    return false
  }
  
  // 请求被取消
  if (axios.isCancel(error)) {
    return false
  }
  
  // 网络错误或超时
  if (!error.response) {
    return true
  }
  
  // 检查状态码是否可重试
  const status = error.response.status
  if (RETRY_CONFIG.retryStatusCodes.includes(status)) {
    return true
  }
  
  // 检查业务错误码是否可重试
  const errorCode = error.response.data?.code
  if (errorCode && isRetryableError(errorCode)) {
    return true
  }
  
  return false
}

// 处理认证错误
function handleAuthError(errorCode, errorMessage) {
  let message = errorMessage
  
  switch (errorCode) {
    case ErrorCodes.TOKEN_EXPIRED:
      message = '登录已过期，请重新登录'
      break
    case ErrorCodes.TOKEN_INVALID:
      message = '登录凭证无效，请重新登录'
      break
    case ErrorCodes.UNAUTHORIZED:
      message = '请先登录'
      break
    case ErrorCodes.PERMISSION_DENIED:
      message = '没有操作权限'
      ElMessage.error(message)
      return // 权限不足不需要跳转登录
  }
  
  ElMessage.error(message)
  localStorage.removeItem('token')
  
  // 延迟跳转，让用户看到提示
  setTimeout(() => {
    window.location.href = '/login'
  }, 1500)
}

// 带重试的请求方法
export const requestWithRetry = async (config, customRetryConfig = {}) => {
  const retryConfig = { ...RETRY_CONFIG, ...customRetryConfig }
  config._retryCount = 0
  
  const makeRequest = async () => {
    try {
      return await request(config)
    } catch (error) {
      if (config._retryCount < retryConfig.maxRetries && shouldRetry(error, config)) {
        config._retryCount++
        const retryDelay = getRetryDelay(config._retryCount - 1)
        await delay(retryDelay)
        return makeRequest()
      }
      throw error
    }
  }
  
  return makeRequest()
}

// 批量请求（带并发控制）
export const batchRequest = async (requests, concurrency = 3) => {
  const results = []
  const executing = []
  
  for (const req of requests) {
    const p = Promise.resolve().then(() => request(req))
    results.push(p)
    
    if (concurrency <= requests.length) {
      const e = p.then(() => executing.splice(executing.indexOf(e), 1))
      executing.push(e)
      if (executing.length >= concurrency) {
        await Promise.race(executing)
      }
    }
  }
  
  return Promise.allSettled(results)
}

// 导出日志收集器供其他模块使用
export const logger = systemLogger

// 导出错误码相关工具
export { ErrorCodes, getErrorMessage, isAuthError, isRetryableError, isRateLimitError }

export default request
