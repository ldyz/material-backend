/**
 * WebSocket 客户端工具
 *
 * 提供自动重连、心跳检测、消息队列等功能
 *
 * @module WebSocket
 * @date 2025-02-07
 */

/**
 * 默认配置
 */
const DEFAULT_CONFIG = {
  // 初始重连间隔（毫秒）
  reconnectInterval: 3000,

  // 最大重连次数（0 表示无限重连，但会使用指数退避）
  maxReconnectAttempts: 10,

  // 最大重连间隔（毫秒）- 指数退避上限
  maxReconnectInterval: 60000,

  // 心跳间隔（毫秒）
  heartbeatInterval: 30000,

  // 心跳超时（毫秒）- 应该大于心跳间隔
  heartbeatTimeout: 35000,

  // 消息队列大小
  messageQueueSize: 100,

  // 是否启用调试日志
  debug: false
}

/**
 * WebSocket 管理类
 */
class WebSocketManager {
  constructor(url, options = {}) {
    this.url = url
    this.options = { ...DEFAULT_CONFIG, ...options }

    // WebSocket 实例
    this.ws = null

    // 连接状态
    this.connected = false
    this.reconnecting = false
    this.manualClose = false

    // 重连计数
    this.reconnectAttempts = 0
    this.reconnectTimer = null

    // 连接失败计数（用于检测代理环境）
    this.consecutiveFailures = 0

    // 心跳定时器
    this.heartbeatTimer = null
    this.heartbeatTimeoutTimer = null

    // 消息队列（连接未建立时缓存消息）
    this.messageQueue = []

    // 事件处理器
    this.handlers = {
      onOpen: [],
      onMessage: [],
      onError: [],
      onClose: []
    }
  }

  /**
   * 输出调试日志
   */
  log(...args) {
    if (this.options.debug) {
      console.log('[WebSocket]', ...args)
    }
  }

  /**
   * 连接 WebSocket
   */
  connect() {
    if (this.ws && (this.ws.readyState === WebSocket.CONNECTING || this.ws.readyState === WebSocket.OPEN)) {
      this.log('已连接或正在连接')
      return
    }

    this.manualClose = false
    this.reconnecting = false

    try {
      // 构建 URL（添加 token 参数）
      const token = this.options.token
      const urlWithToken = `${this.url}?token=${encodeURIComponent(token)}`

      this.ws = new WebSocket(urlWithToken)

      // 设置事件处理
      this.ws.onopen = (event) => this.handleOpen(event)
      this.ws.onmessage = (event) => this.handleMessage(event)
      this.ws.onerror = (error) => this.handleError(error)
      this.ws.onclose = (event) => this.handleClose(event)

      this.log('正在连接...', this.url)
    } catch (error) {
      console.error('WebSocket 连接失败:', error)
      this.scheduleReconnect()
    }
  }

  /**
   * 断开 WebSocket 连接
   */
  disconnect() {
    this.manualClose = true
    this.clearTimers()

    if (this.ws) {
      this.ws.close()
      this.ws = null
    }

    this.connected = false
    this.log('已手动断开')
  }

  /**
   * 发送消息
   *
   * @param {Object|string} data - 要发送的数据
   */
  send(data) {
    const message = typeof data === 'string' ? data : JSON.stringify(data)

    if (this.connected) {
      try {
        this.ws.send(message)
      } catch (error) {
        console.error('发送消息失败:', error)
        this.enqueueMessage(message)
      }
    } else {
      // 连接未建立，加入队列
      this.enqueueMessage(message)
    }
  }

  /**
   * 添加消息到队列
   *
   * @param {string} message - 消息内容
   */
  enqueueMessage(message) {
    if (this.messageQueue.length >= this.options.messageQueueSize) {
      // 队列已满，移除最旧的消息
      this.messageQueue.shift()
    }
    this.messageQueue.push(message)
  }

  /**
   * 发送队列中的消息
   */
  flushMessageQueue() {
    while (this.messageQueue.length > 0 && this.connected) {
      const message = this.messageQueue.shift()
      try {
        this.ws.send(message)
      } catch (error) {
        console.error('发送队列消息失败:', error)
        // 失败的消息放回队列开头
        this.messageQueue.unshift(message)
        break
      }
    }
  }

  /**
   * 处理连接打开事件
   */
  handleOpen(event) {
    this.log('连接成功')
    this.connected = true
    this.reconnecting = false
    this.reconnectAttempts = 0
    this.consecutiveFailures = 0  // 重置连续失败计数

    // 发送队列中的消息
    this.flushMessageQueue()

    // 启动心跳
    this.startHeartbeat()

    // 触发回调
    this.emit('open', event)
  }

  /**
   * 处理消息接收事件
   */
  handleMessage(event) {
    // 重置心跳超时
    this.resetHeartbeatTimeout()

    try {
      const data = JSON.parse(event.data)
      this.emit('message', data)
    } catch (error) {
      // 如果不是 JSON，直接传递原始数据
      this.emit('message', event.data)
    }
  }

  /**
   * 处理错误事件
   */
  handleError(error) {
    this.consecutiveFailures++
    this.log('连接错误，连续失败次数:', this.consecutiveFailures)
    this.emit('error', error)
  }

  /**
   * 处理连接关闭事件
   */
  handleClose(event) {
    this.log('连接关闭:', event.code, event.reason)
    this.connected = false
    this.clearTimers()

    // 触发回调
    this.emit('close', event)

    // 如果不是手动关闭，尝试重连
    if (!this.manualClose) {
      this.scheduleReconnect()
    }
  }

  /**
   * 计算重连间隔（指数退避）
   */
  getReconnectInterval() {
    // 基础间隔
    const baseInterval = this.options.reconnectInterval
    // 指数退避因子
    const backoffFactor = Math.min(this.reconnectAttempts, 10)
    // 计算间隔（每次增加1.5倍，最多60秒）
    const interval = Math.min(
      baseInterval * Math.pow(1.5, backoffFactor),
      this.options.maxReconnectInterval
    )
    return Math.floor(interval)
  }

  /**
   * 安排重连
   */
  scheduleReconnect() {
    if (this.manualClose) {
      return
    }

    // 检查最大重连次数
    if (this.options.maxReconnectAttempts > 0 &&
        this.reconnectAttempts >= this.options.maxReconnectAttempts) {
      console.warn('[WebSocket] 重连次数已达上限，停止重连')
      return
    }

    if (this.reconnecting) {
      return
    }

    this.reconnecting = true
    this.reconnectAttempts++

    const interval = this.getReconnectInterval()
    this.log(`将在 ${interval}ms 后尝试重连 (第 ${this.reconnectAttempts} 次)`)

    this.reconnectTimer = setTimeout(() => {
      this.reconnecting = false
      this.connect()
    }, interval)
  }

  /**
   * 启动心跳
   */
  startHeartbeat() {
    this.heartbeatTimer = setInterval(() => {
      if (this.connected) {
        try {
          this.ws.send(JSON.stringify({ type: 'ping' }))
          // 发送ping后，设置超时等待pong
          this.resetHeartbeatTimeout()
        } catch (error) {
          console.error('发送心跳失败:', error)
        }
      }
    }, this.options.heartbeatInterval)
  }

  /**
   * 重置心跳超时
   */
  resetHeartbeatTimeout() {
    this.clearHeartbeatTimeout()

    this.heartbeatTimeoutTimer = setTimeout(() => {
      this.log('心跳超时，关闭连接')
      if (this.ws) {
        this.ws.close()
      }
    }, this.options.heartbeatTimeout)
  }

  /**
   * 清除心跳定时器
   */
  clearHeartbeatTimeout() {
    if (this.heartbeatTimeoutTimer) {
      clearTimeout(this.heartbeatTimeoutTimer)
      this.heartbeatTimeoutTimer = null
    }
  }

  /**
   * 清除所有定时器
   */
  clearTimers() {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer)
      this.reconnectTimer = null
    }

    if (this.heartbeatTimer) {
      clearInterval(this.heartbeatTimer)
      this.heartbeatTimer = null
    }

    this.clearHeartbeatTimeout()
  }

  /**
   * 添加事件监听器
   *
   * @param {string} event - 事件名称 (open/message/error/close)
   * @param {Function} handler - 处理函数
   */
  on(event, handler) {
    if (this.handlers[`on${event.charAt(0).toUpperCase() + event.slice(1)}`]) {
      this.handlers[`on${event.charAt(0).toUpperCase() + event.slice(1)}`].push(handler)
    }
  }

  /**
   * 移除事件监听器
   *
   * @param {string} event - 事件名称
   * @param {Function} handler - 处理函数
   */
  off(event, handler) {
    const eventName = `on${event.charAt(0).toUpperCase() + event.slice(1)}`
    const handlers = this.handlers[eventName]
    if (handlers) {
      const index = handlers.indexOf(handler)
      if (index !== -1) {
        handlers.splice(index, 1)
      }
    }
  }

  /**
   * 触发事件
   *
   * @param {string} event - 事件名称
   * @param {any} data - 事件数据
   */
  emit(event, data) {
    const handlers = this.handlers[`on${event.charAt(0).toUpperCase() + event.slice(1)}`]
    if (handlers) {
      handlers.forEach(handler => {
        try {
          handler(data)
        } catch (error) {
          console.error(`事件处理器错误 (${event}):`, error)
        }
      })
    }
  }

  /**
   * 检查连接状态
   */
  isConnected() {
    return this.connected
  }
}

/**
 * 创建 WebSocket 实例的组合式函数
 *
 * @param {string} url - WebSocket URL
 * @param {Object} options - 配置选项
 * @returns {Object} WebSocket 实例和方法
 *
 * @example
 * const { connect, disconnect, send, isConnected } = useWebSocket('ws://localhost:8088/ws', {
 *   token: 'your-token',
 *   onMessage: (data) => console.log('收到消息:', data)
 * })
 *
 * // 连接
 * connect()
 *
 * // 发送消息
 * send({ type: 'hello', data: 'world' })
 *
 * // 断开连接
 * disconnect()
 */
export function useWebSocket(url, options = {}) {
  const manager = new WebSocketManager(url, options)

  // 设置事件处理
  if (options.onOpen) {
    manager.on('open', options.onOpen)
  }
  if (options.onMessage) {
    manager.on('message', options.onMessage)
  }
  if (options.onError) {
    manager.on('error', options.onError)
  }
  if (options.onClose) {
    manager.on('close', options.onClose)
  }

  return {
    connect: () => manager.connect(),
    disconnect: () => manager.disconnect(),
    send: (data) => manager.send(data),
    on: (event, handler) => manager.on(event, handler),
    off: (event, handler) => manager.off(event, handler),
    isConnected: () => manager.isConnected()
  }
}

export default WebSocketManager
