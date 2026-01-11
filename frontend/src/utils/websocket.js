/**
 * WebSocket客户端
 * 用于实时接收监控数据和告警通知
 */

class WebSocketClient {
  constructor() {
    this.ws = null
    this.url = ''
    this.reconnectAttempts = 0
    this.maxReconnectAttempts = 5
    this.reconnectInterval = 3000
    this.heartbeatInterval = null
    this.listeners = new Map()
    this.isConnected = false
  }

  /**
   * 连接WebSocket
   * @param {string} appId - APP ID
   * @param {string} userId - 用户ID
   */
  connect(appId, userId) {
    // 检查appId是否有效
    if (!appId || appId === '') {
      console.warn('[WebSocket] appId is empty, skipping connection')
      return
    }
    
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    // 使用当前主机，Vite会代理WebSocket连接到后端
    const host = window.location.host
    this.url = `${protocol}//${host}/api/v1/ws?app_id=${appId}&user_id=${userId || ''}`
    console.log('[WebSocket] Connecting to:', this.url)

    this.createConnection()
  }

  /**
   * 创建WebSocket连接
   */
  createConnection() {
    try {
      this.ws = new WebSocket(this.url)

      this.ws.onopen = () => {
        console.log('[WebSocket] Connected')
        this.isConnected = true
        this.reconnectAttempts = 0
        this.startHeartbeat()
        this.emit('connected')
      }

      this.ws.onmessage = (event) => {
        try {
          // 处理多条消息（以换行分隔）
          const messages = event.data.split('\n').filter(m => m.trim())
          messages.forEach(msgStr => {
            const message = JSON.parse(msgStr)
            this.handleMessage(message)
          })
        } catch (e) {
          console.error('[WebSocket] Parse error:', e)
        }
      }

      this.ws.onclose = (event) => {
        console.log('[WebSocket] Disconnected:', event.code, event.reason)
        this.isConnected = false
        this.stopHeartbeat()
        this.emit('disconnected')
        this.attemptReconnect()
      }

      this.ws.onerror = (error) => {
        console.error('[WebSocket] Error:', error)
        this.emit('error', error)
      }
    } catch (e) {
      console.error('[WebSocket] Connection error:', e)
      this.attemptReconnect()
    }
  }

  /**
   * 处理接收到的消息
   */
  handleMessage(message) {
    const { type, data, timestamp } = message

    switch (type) {
      case 'pong':
        // 心跳响应
        break
      case 'monitor':
        this.emit('monitor', data)
        break
      case 'alert':
        this.emit('alert', data)
        this.showAlertNotification(data)
        break
      case 'notification':
        this.emit('notification', data)
        this.showNotification(data)
        break
      case 'log':
        this.emit('log', data)
        break
      default:
        this.emit('message', message)
    }
  }

  /**
   * 显示告警通知
   */
  showAlertNotification(alert) {
    if (Notification.permission === 'granted') {
      const levelEmoji = {
        critical: '🔴',
        warning: '🟡',
        info: '🔵'
      }
      new Notification(`${levelEmoji[alert.level] || '⚪'} ${alert.title}`, {
        body: alert.message,
        icon: '/favicon.ico',
        tag: `alert-${alert.id}`
      })
    }
  }

  /**
   * 显示普通通知
   */
  showNotification(notification) {
    if (Notification.permission === 'granted') {
      new Notification(notification.title, {
        body: notification.message,
        icon: '/favicon.ico'
      })
    }
  }

  /**
   * 发送消息
   */
  send(type, data) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify({ type, data }))
    }
  }

  /**
   * 开始心跳
   */
  startHeartbeat() {
    this.heartbeatInterval = setInterval(() => {
      this.send('ping')
    }, 30000)
  }

  /**
   * 停止心跳
   */
  stopHeartbeat() {
    if (this.heartbeatInterval) {
      clearInterval(this.heartbeatInterval)
      this.heartbeatInterval = null
    }
  }

  /**
   * 尝试重连
   */
  attemptReconnect() {
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++
      console.log(`[WebSocket] Reconnecting... (${this.reconnectAttempts}/${this.maxReconnectAttempts})`)
      setTimeout(() => {
        this.createConnection()
      }, this.reconnectInterval)
    } else {
      console.log('[WebSocket] Max reconnect attempts reached')
      this.emit('maxReconnectReached')
    }
  }

  /**
   * 断开连接
   */
  disconnect() {
    this.stopHeartbeat()
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
    this.isConnected = false
  }

  /**
   * 添加事件监听
   */
  on(event, callback) {
    if (!this.listeners.has(event)) {
      this.listeners.set(event, [])
    }
    this.listeners.get(event).push(callback)
  }

  /**
   * 移除事件监听
   */
  off(event, callback) {
    if (this.listeners.has(event)) {
      const callbacks = this.listeners.get(event)
      const index = callbacks.indexOf(callback)
      if (index > -1) {
        callbacks.splice(index, 1)
      }
    }
  }

  /**
   * 触发事件
   */
  emit(event, data) {
    if (this.listeners.has(event)) {
      this.listeners.get(event).forEach(callback => {
        try {
          callback(data)
        } catch (e) {
          console.error('[WebSocket] Listener error:', e)
        }
      })
    }
  }

  /**
   * 请求通知权限
   */
  static requestNotificationPermission() {
    if ('Notification' in window && Notification.permission === 'default') {
      Notification.requestPermission()
    }
  }
}

// 创建单例实例
const wsClient = new WebSocketClient()

export default wsClient
export { WebSocketClient }
