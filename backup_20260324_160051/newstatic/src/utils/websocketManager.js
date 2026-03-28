/**
 * WebSocket Connection Manager
 *
 * Manages WebSocket connections using socket.io-client for real-time collaboration.
 * Handles connection lifecycle, event processing, and conflict resolution using OT.
 *
 * @class WebSocketManager
 * @example
 * const wsManager = new WebSocketManager('http://localhost:8080', 123);
 * await wsManager.connect();
 * wsManager.on('task:update', (data) => console.log('Task updated:', data));
 */

import { io } from 'socket.io-client'
import { ElNotification } from 'element-plus'

/**
 * @typedef {Object} ConnectionStatus
 * @property {string} status - 'disconnected' | 'connecting' | 'connected' | 'error'
 * @property {number} attempts - Reconnection attempt count
 * @property {string|null} lastError - Last error message
 */

/**
 * @typedef {Object} OTOperation
 * @property {string} type - 'insert' | 'delete' | 'retain'
 * @property {number} position - Position in document
 * @property {string|number} value - Value to insert/delete or count to retain
 */

/**
 * @typedef {Object} CollaborationEvent
 * @property {string} type - Event type
 * @property {number} taskId - Task ID
 * @property {object} payload - Event data
 * @property {string} userId - User ID who made the change
 * @property {string} timestamp - ISO timestamp
 */

class WebSocketManager {
  /**
   * Creates a new WebSocket manager instance
   *
   * @param {string} url - WebSocket server URL
   * @param {number} projectId - Project ID for room management
   * @param {object} options - Configuration options
   * @param {number} options.timeout - Connection timeout in ms (default: 10000)
   * @param {number} options.reconnectAttempts - Max reconnection attempts (default: 5)
   * @param {number} options.reconnectDelay - Delay between attempts in ms (default: 3000)
   */
  constructor(url, projectId, options = {}) {
    this.url = url
    this.projectId = projectId
    this.socket = null
    this.eventHandlers = new Map()

    // Connection state
    this._connectionStatus = {
      status: 'disconnected',
      attempts: 0,
      lastError: null
    }

    // Configuration
    this.options = {
      timeout: options.timeout || 10000,
      reconnectAttempts: options.reconnectAttempts || 5,
      reconnectDelay: options.reconnectDelay || 3000,
      ...options
    }

    // OT state
    this.pendingOperations = []
    this.revision = 0
  }

  /**
   * Establishes WebSocket connection
   *
   * @returns {Promise<void>}
   * @throws {Error} If connection fails after max attempts
   *
   * @example
   * try {
   *   await wsManager.connect()
   *   console.log('Connected!')
   * } catch (error) {
   *   console.error('Connection failed:', error)
   * }
   */
  async connect() {
    if (this.isConnected()) {
      console.warn('[WebSocket] Already connected')
      return
    }

    this._setConnectionStatus('connecting')

    try {
      this.socket = io(this.url, {
        query: { projectId: this.projectId },
        timeout: this.options.timeout,
        reconnection: true,
        reconnectionAttempts: this.options.reconnectAttempts,
        reconnectionDelay: this.options.reconnectDelay,
        transports: ['websocket', 'polling']
      })

      // Setup event listeners
      this._setupEventListeners()

      // Wait for connection
      await new Promise((resolve, reject) => {
        const timeout = setTimeout(() => {
          reject(new Error('Connection timeout'))
        }, this.options.timeout)

        this.socket.once('connect', () => {
          clearTimeout(timeout)
          resolve()
        })

        this.socket.once('connect_error', (error) => {
          clearTimeout(timeout)
          reject(error)
        })
      })

      // Join project room
      this.socket.emit('join:project', { projectId: this.projectId })

      this._setConnectionStatus('connected')
      this._connectionStatus.attempts = 0

      ElNotification.success({
        title: 'Connected',
        message: 'Real-time collaboration enabled',
        duration: 2000
      })

      console.log('[WebSocket] Connected to project', this.projectId)
    } catch (error) {
      this._setConnectionStatus('error', error.message)
      this._connectionStatus.attempts++

      ElNotification.error({
        title: 'Connection Failed',
        message: error.message,
        duration: 5000
      })

      throw error
    }
  }

  /**
   * Closes WebSocket connection
   *
   * @example
   * wsManager.disconnect()
   */
  disconnect() {
    if (!this.socket) return

    // Leave project room
    this.socket.emit('leave:project', { projectId: this.projectId })

    // Remove all listeners
    this.socket.removeAllListeners()

    // Disconnect
    this.socket.disconnect()
    this.socket = null

    this._setConnectionStatus('disconnected')
    console.log('[WebSocket] Disconnected from project', this.projectId)
  }

  /**
   * Emits an event to the server
   *
   * @param {string} event - Event name
   * @param {object} data - Event data
   *
   * @example
   * wsManager.emit('task:update', { taskId: 123, name: 'New name' })
   */
  emit(event, data) {
    if (!this.isConnected()) {
      console.warn('[WebSocket] Cannot emit event: not connected')
      return
    }

    // Add metadata
    const eventData = {
      ...data,
      projectId: this.projectId,
      timestamp: new Date().toISOString(),
      revision: ++this.revision
    }

    // Store operation for OT
    this.pendingOperations.push({
      event,
      data: eventData,
      revision: this.revision
    })

    this.socket.emit(event, eventData)
    console.log('[WebSocket] Emitted:', event, eventData)
  }

  /**
   * Registers an event handler
   *
   * @param {string} event - Event name
   * @param {Function} callback - Event handler function
   * @returns {Function} Unsubscribe function
   *
   * @example
   * const unsubscribe = wsManager.on('task:update', (data) => {
   *   console.log('Task updated:', data)
   * })
   * // Later: unsubscribe()
   */
  on(event, callback) {
    if (!this.eventHandlers.has(event)) {
      this.eventHandlers.set(event, new Set())
    }

    this.eventHandlers.get(event).add(callback)

    // Also register with socket
    if (this.socket) {
      this.socket.on(event, callback)
    }

    // Return unsubscribe function
    return () => this.off(event, callback)
  }

  /**
   * Removes an event handler
   *
   * @param {string} event - Event name
   * @param {Function} callback - Event handler to remove
   */
  off(event, callback) {
    if (this.eventHandlers.has(event)) {
      this.eventHandlers.get(event).delete(callback)
    }

    if (this.socket) {
      this.socket.off(event, callback)
    }
  }

  /**
   * Checks if connected to WebSocket server
   *
   * @returns {boolean}
   */
  isConnected() {
    return this.socket?.connected === true
  }

  /**
   * Gets current connection status
   *
   * @returns {ConnectionStatus}
   */
  getConnectionStatus() {
    return { ...this._connectionStatus }
  }

  /**
   * Applies Operational Transformation to resolve conflicts
   *
   * @param {OTOperation[]} localOps - Local operations
   * @param {OTOperation[]} remoteOps - Remote operations
   * @returns {OTOperation[]} Transformed operations
   */
  transform(localOps, remoteOps) {
    // Simple OT implementation based on operational transformation
    // In production, use a library like ShareJS or Yjs

    const transformed = [...localOps]

    for (const remote of remoteOps) {
      for (const local of transformed) {
        if (remote.position <= local.position) {
          // Shift local operation position
          if (remote.type === 'insert') {
            local.position += remote.value?.length || 1
          } else if (remote.type === 'delete') {
            local.position -= remote.value?.length || 1
          }
        }
      }
    }

    return transformed
  }

  /**
   * Broadcasts cursor position to other users
   *
   * @param {object} cursor - Cursor position
   * @param {number} cursor.x - X coordinate
   * @param {number} cursor.y - Y coordinate
   * @param {number} cursor.taskId - Task ID if hovering over task
   */
  broadcastCursor(cursor) {
    this.emit('cursor:move', cursor)
  }

  /**
   * Sends typing indicator
   *
   * @param {boolean} isTyping - Whether user is typing
   * @param {string} taskId - Task ID being edited
   */
  sendTypingIndicator(isTyping, taskId) {
    this.emit('user:typing', { isTyping, taskId })
  }

  // ==================== Private Methods ====================

  /**
   * Sets up internal event listeners
   * @private
   */
  _setupEventListeners() {
    // Connection events
    this.socket.on('connect', () => {
      console.log('[WebSocket] Connected')
      this._setConnectionStatus('connected')
    })

    this.socket.on('disconnect', (reason) => {
      console.log('[WebSocket] Disconnected:', reason)
      this._setConnectionStatus('disconnected')
    })

    this.socket.on('reconnect', (attemptNumber) => {
      console.log('[WebSocket] Reconnected after', attemptNumber, 'attempts')
      this._setConnectionStatus('connected')

      ElNotification.success({
        title: 'Reconnected',
        message: 'Real-time collaboration restored',
        duration: 2000
      })

      // Re-join project room
      this.socket.emit('join:project', { projectId: this.projectId })
    })

    this.socket.on('reconnect_attempt', (attemptNumber) => {
      console.log('[WebSocket] Reconnection attempt', attemptNumber)
      this._setConnectionStatus('connecting')
    })

    this.socket.on('reconnect_failed', () => {
      console.error('[WebSocket] Reconnection failed')
      this._setConnectionStatus('error', 'Reconnection failed')

      ElNotification.error({
        title: 'Connection Lost',
        message: 'Unable to reconnect to server',
        duration: 5000
      })
    })

    this.socket.on('connect_error', (error) => {
      console.error('[WebSocket] Connection error:', error)
    })

    this.socket.on('error', (error) => {
      console.error('[WebSocket] Error:', error)
    })

    // Project events
    this.socket.on('user:joined', (data) => {
      console.log('[WebSocket] User joined:', data)
      this._triggerEvent('user:joined', data)
    })

    this.socket.on('user:left', (data) => {
      console.log('[WebSocket] User left:', data)
      this._triggerEvent('user:left', data)
    })

    // Gantt events
    this.socket.on('task:update', (data) => {
      console.log('[WebSocket] Task update received:', data)
      this._triggerEvent('task:update', data)
    })

    this.socket.on('task:create', (data) => {
      console.log('[WebSocket] Task created:', data)
      this._triggerEvent('task:create', data)
    })

    this.socketOn('task:delete', (data) => {
      console.log('[WebSocket] Task deleted:', data)
      this._triggerEvent('task:delete', data)
    })

    // Cursor events
    this.socket.on('cursor:update', (data) => {
      this._triggerEvent('cursor:update', data)
    })
  }

  /**
   * Triggers event handlers for an event
   * @private
   */
  _triggerEvent(event, data) {
    if (this.eventHandlers.has(event)) {
      this.eventHandlers.get(event).forEach(callback => {
        try {
          callback(data)
        } catch (error) {
          console.error('[WebSocket] Error in event handler:', error)
        }
      })
    }
  }

  /**
   * Updates connection status
   * @private
   */
  _setConnectionStatus(status, error = null) {
    this._connectionStatus.status = status
    this._connectionStatus.lastError = error
  }
}

export default WebSocketManager
