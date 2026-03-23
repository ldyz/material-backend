/**
 * Collaboration Pinia Store
 *
 * Manages real-time collaboration state including connected users,
 * cursor positions, typing indicators, and conflict resolution.
 */

import { defineStore } from 'pinia'
import { ref, reactive, computed } from 'vue'
import WebSocketManager from '@/utils/websocketManager'

/**
 * @typedef {Object} ConnectedUser
 * @property {string} id - User ID
 * @property {string} name - User name
 * @property {string} avatar - User avatar URL
 * @property {string} color - User's cursor color
 * @property {Date} lastSeen - Last activity timestamp
 */

/**
 * @typedef {Object} CursorPosition
 * @property {number} x - X coordinate
 * @property {number} y - Y coordinate
 * @property {number|null} taskId - Task ID if hovering over task
 * @property {string} userId - User ID
 */

export const useCollaborationStore = defineStore('collaboration', () => {
  // ==================== State ====================

  /** @type {import('vue').Ref<WebSocketManager|null>} */
  const wsManager = ref(null)

  /** @type {import('vue').Ref<string>} */
  const projectId = ref(null)

  /** @type {import('vue').Ref<ConnectionStatus>} */
  const connectionStatus = ref('disconnected')

  /** @type {import('vue').Ref<Map<string, ConnectedUser>>} */
  const connectedUsers = ref(new Map())

  /** @type {import('vue').Ref<Map<string, CursorPosition>>} */
  const remoteCursors = ref(new Map())

  /** @type {import('vue').Ref<Map<string, {isTyping: boolean, taskId: string|null}>>} */
  const typingUsers = ref(new Map())

  /** @type {import('vue').Ref<Object>} */
  const myCursor = reactive({
    x: 0,
    y: 0,
    taskId: null
  })

  /** @type {import('vue).Ref<Array>} */
  const pendingUpdates = ref([])

  /** @type {import('vue').Ref<Array>} */
  const conflicts = ref([])

  // ==================== Getters ====================

  /**
   * Check if connected to WebSocket server
   */
  const isConnected = computed(() => {
    return connectionStatus.value === 'connected' && wsManager.value?.isConnected()
  })

  /**
   * Get list of connected users (excluding current user)
   */
  const otherUsers = computed(() => {
    const currentUserId = localStorage.getItem('userId')
    return Array.from(connectedUsers.value.values())
      .filter(user => user.id !== currentUserId)
  })

  /**
   * Get user by ID
   */
  const getUserById = computed(() => {
    return (userId) => connectedUsers.value.get(userId)
  })

  /**
   * Get cursor by user ID
   */
  const getCursorByUserId = computed(() => {
    return (userId) => remoteCursors.value.get(userId)
  })

  /**
   * Check if any user is typing on a task
   */
  const getUsersTypingOnTask = computed(() => {
    return (taskId) => {
      return Array.from(typingUsers.value.entries())
        .filter(([_, data]) => data.isTyping && data.taskId === taskId)
        .map(([userId, _]) => userId)
    }
  })

  /**
   * Get connection status text
   */
  const connectionStatusText = computed(() => {
    const statusMap = {
      disconnected: 'Disconnected',
      connecting: 'Connecting...',
      connected: 'Connected',
      error: 'Connection Error'
    }
    return statusMap[connectionStatus.value] || 'Unknown'
  })

  /**
   * Get connection status color
   */
  const connectionStatusColor = computed(() => {
    const colorMap = {
      disconnected: 'info',
      connecting: 'warning',
      connected: 'success',
      error: 'danger'
    }
    return colorMap[connectionStatus.value] || 'info'
  })

  // ==================== Actions ====================

  /**
   * Initialize WebSocket connection
   *
   * @param {string} url - WebSocket server URL
   * @param {number} pid - Project ID
   */
  async function connect(url, pid) {
    if (wsManager.value) {
      await disconnect()
    }

    projectId.value = pid
    connectionStatus.value = 'connecting'

    try {
      wsManager.value = new WebSocketManager(url, pid)

      // Setup event listeners
      _setupEventListeners()

      // Connect
      await wsManager.value.connect()
      connectionStatus.value = 'connected'
    } catch (error) {
      connectionStatus.value = 'error'
      console.error('[Collaboration] Connection failed:', error)
      throw error
    }
  }

  /**
   * Disconnect from WebSocket server
   */
  async function disconnect() {
    if (wsManager.value) {
      wsManager.value.disconnect()
      wsManager.value = null
    }

    // Reset state
    connectionStatus.value = 'disconnected'
    connectedUsers.value.clear()
    remoteCursors.value.clear()
    typingUsers.value.clear()
    pendingUpdates.value = []
    conflicts.value = []
  }

  /**
   * Register event handler
   *
   * @param {string} event - Event name
   * @param {Function} callback - Event handler
   * @returns {Function} Unsubscribe function
   */
  function on(event, callback) {
    if (!wsManager.value) {
      console.warn('[Collaboration] Cannot register event: not connected')
      return () => {}
    }

    return wsManager.value.on(event, callback)
  }

  /**
   * Emit event to server
   *
   * @param {string} event - Event name
   * @param {object} data - Event data
   */
  function emit(event, data) {
    if (!wsManager.value || !isConnected.value) {
      console.warn('[Collaboration] Cannot emit event: not connected')
      return
    }

    wsManager.value.emit(event, data)
  }

  /**
   * Update local cursor position
   *
   * @param {number} x - X coordinate
   * @param {number} y - Y coordinate
   * @param {number|null} taskId - Task ID if hovering over task
   */
  function updateCursor(x, y, taskId = null) {
    myCursor.x = x
    myCursor.y = y
    myCursor.taskId = taskId

    if (isConnected.value) {
      wsManager.value.broadcastCursor({ x, y, taskId })
    }
  }

  /**
   * Send typing indicator
   *
   * @param {boolean} isTyping - Whether user is typing
   * @param {string|null} taskId - Task ID being edited
   */
  function sendTyping(isTyping, taskId = null) {
    if (isConnected.value) {
      wsManager.value.sendTypingIndicator(isTyping, taskId)
    }
  }

  /**
   * Add user to connected users
   *
   * @param {ConnectedUser} user - User data
   */
  function addUser(user) {
    connectedUsers.value.set(user.id, {
      ...user,
      lastSeen: new Date()
    })
  }

  /**
   * Remove user from connected users
   *
   * @param {string} userId - User ID
   */
  function removeUser(userId) {
    connectedUsers.value.delete(userId)
    remoteCursors.value.delete(userId)
    typingUsers.value.delete(userId)
  }

  /**
   * Update remote cursor position
   *
   * @param {string} userId - User ID
   * @param {CursorPosition} cursor - Cursor position
   */
  function updateRemoteCursor(userId, cursor) {
    remoteCursors.value.set(userId, cursor)

    // Remove cursor after 5 seconds of inactivity
    setTimeout(() => {
      if (remoteCursors.value.get(userId)?.timestamp === cursor.timestamp) {
        remoteCursors.value.delete(userId)
      }
    }, 5000)
  }

  /**
   * Update user typing status
   *
   * @param {string} userId - User ID
   * @param {boolean} isTyping - Whether user is typing
   * @param {string|null} taskId - Task ID being edited
   */
  function updateTypingStatus(userId, isTyping, taskId = null) {
    if (isTyping) {
      typingUsers.value.set(userId, { isTyping, taskId })

      // Auto-clear after 3 seconds
      setTimeout(() => {
        const current = typingUsers.value.get(userId)
        if (current?.taskId === taskId) {
          typingUsers.value.delete(userId)
        }
      }, 3000)
    } else {
      typingUsers.value.delete(userId)
    }
  }

  /**
   * Add pending update
   *
   * @param {object} update - Update data
   */
  function addPendingUpdate(update) {
    pendingUpdates.value.push({
      ...update,
      timestamp: Date.now()
    })
  }

  /**
   * Remove pending update
   *
   * @param {string} updateId - Update ID
   */
  function removePendingUpdate(updateId) {
    const index = pendingUpdates.value.findIndex(u => u.id === updateId)
    if (index !== -1) {
      pendingUpdates.value.splice(index, 1)
    }
  }

  /**
   * Add conflict
   *
   * @param {object} conflict - Conflict data
   */
  function addConflict(conflict) {
    conflicts.value.push({
      ...conflict,
      timestamp: Date.now(),
      resolved: false
    })
  }

  /**
   * Resolve conflict
   *
   * @param {string} conflictId - Conflict ID
   */
  function resolveConflict(conflictId) {
    const conflict = conflicts.value.find(c => c.id === conflictId)
    if (conflict) {
      conflict.resolved = true
    }
  }

  /**
   * Get user color for visualization
   *
   * @param {string} userId - User ID
   * @returns {string} Hex color code
   */
  function getUserColor(userId) {
    const user = connectedUsers.value.get(userId)
    if (user?.color) {
      return user.color
    }

    // Generate consistent color from user ID
    const hash = userId.split('').reduce((acc, char) => {
      return acc + char.charCodeAt(0)
    }, 0)

    const colors = [
      '#FF6B6B', '#4ECDC4', '#45B7D1', '#FFA07A',
      '#98D8C8', '#F7DC6F', '#BB8FCE', '#85C1E2',
      '#F8B739', '#52B788'
    ]

    return colors[hash % colors.length]
  }

  // ==================== Private Methods ====================

  /**
   * Setup WebSocket event listeners
   * @private
   */
  function _setupEventListeners() {
    if (!wsManager.value) return

    // User joined
    wsManager.value.on('user:joined', (data) => {
      console.log('[Collaboration] User joined:', data)
      addUser(data.user)
    })

    // User left
    wsManager.value.on('user:left', (data) => {
      console.log('[Collaboration] User left:', data)
      removeUser(data.userId)
    })

    // Cursor updates
    wsManager.value.on('cursor:update', (data) => {
      if (data.userId !== localStorage.getItem('userId')) {
        updateRemoteCursor(data.userId, data.cursor)
      }
    })

    // User typing
    wsManager.value.on('user:typing', (data) => {
      if (data.userId !== localStorage.getItem('userId')) {
        updateTypingStatus(data.userId, data.isTyping, data.taskId)
      }
    })

    // Task updates (for conflict detection)
    wsManager.value.on('task:update', (data) => {
      // Check for conflicts with pending updates
      const conflictingUpdate = pendingUpdates.value.find(
        u => u.taskId === data.taskId && u.userId !== data.userId
      )

      if (conflictingUpdate) {
        addConflict({
          id: `conflict-${Date.now()}`,
          type: 'task_update',
          taskId: data.taskId,
          localChanges: conflictingUpdate.changes,
          remoteChanges: data.changes
        })
      }
    })
  }

  return {
    // State
    wsManager,
    projectId,
    connectionStatus,
    connectedUsers,
    remoteCursors,
    typingUsers,
    myCursor,
    pendingUpdates,
    conflicts,

    // Getters
    isConnected,
    otherUsers,
    getUserById,
    getCursorByUserId,
    getUsersTypingOnTask,
    connectionStatusText,
    connectionStatusColor,

    // Actions
    connect,
    disconnect,
    on,
    emit,
    updateCursor,
    sendTyping,
    addUser,
    removeUser,
    updateRemoteCursor,
    updateTypingStatus,
    addPendingUpdate,
    removePendingUpdate,
    addConflict,
    resolveConflict,
    getUserColor
  }
})
