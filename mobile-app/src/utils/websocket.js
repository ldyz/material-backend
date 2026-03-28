/**
 * WebSocket Service for Real-time Push Notifications and AI Chat
 *
 * @module WebSocket
 * @date 2025-02-09
 */

import { storage } from './storage'
import { useNotificationStore } from '@/stores/notification'
import { Capacitor } from '@capacitor/core'

class WebSocketService {
  constructor() {
    this.ws = null
    this.reconnectAttempts = 0
    this.maxReconnectAttempts = 5
    this.reconnectDelay = 3000
    this.isManualClose = false
    this.heartbeatInterval = null
    // AI response callbacks
    this.aiCallbacks = new Set()
  }

  /**
   * Get WebSocket URL based on environment
   */
  getWebSocketUrl() {
    const isCapacitor = Capacitor.isNativePlatform()

    if (isCapacitor) {
      // Production: use wss:// with the server domain
      return 'wss://home.mbed.org.cn:9090/api/notification/ws'
    } else {
      // Development: use ws:// with localhost
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
      const host = window.location.host
      return `${protocol}//${host}/api/notification/ws`
    }
  }

  /**
   * Connect to WebSocket server
   */
  connect() {
    const token = storage.getToken()

    if (!token) {
      console.log('No token available, skipping WebSocket connection')
      return
    }

    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      console.log('WebSocket already connected')
      return
    }

    try {
      const wsUrl = this.getWebSocketUrl()
      console.log('Connecting to WebSocket:', wsUrl)

      // Add token to URL as query parameter
      const urlWithToken = `${wsUrl}?token=${encodeURIComponent(token)}`

      this.ws = new WebSocket(urlWithToken)

      this.ws.onopen = () => {
        console.log('WebSocket connected')
        this.reconnectAttempts = 0
        this.isManualClose = false
        this.startHeartbeat()

        // Update notification store connection status
        const notificationStore = useNotificationStore()
        notificationStore.wsConnected = true
      }

      this.ws.onmessage = (event) => {
        this.handleMessage(event.data)
      }

      this.ws.onerror = (error) => {
        console.error('WebSocket error:', error)
      }

      this.ws.onclose = (event) => {
        console.log('WebSocket closed:', event.code, event.reason)

        // Stop heartbeat
        this.stopHeartbeat()

        // Update notification store connection status
        const notificationStore = useNotificationStore()
        notificationStore.wsConnected = false

        // Attempt to reconnect if not manually closed
        if (!this.isManualClose && this.reconnectAttempts < this.maxReconnectAttempts) {
          this.reconnectAttempts++
          const delay = this.reconnectDelay * this.reconnectAttempts
          console.log(`Reconnecting in ${delay}ms... (attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts})`)
          setTimeout(() => this.connect(), delay)
        }
      }
    } catch (error) {
      console.error('Failed to create WebSocket connection:', error)
    }
  }

  /**
   * Handle incoming WebSocket messages
   * 支持多消息合并发送（换行符分隔）
   */
  handleMessage(data) {
    // 按换行符分割，处理多条消息
    const lines = data.split('\n').filter(line => line.trim())

    for (const line of lines) {
      try {
        const message = JSON.parse(line)
        console.log('WebSocket message received:', message)

        switch (message.type) {
          case 'notification':
            this.handleNotification(message.data)
            break
          case 'unread_count':
            this.handleUnreadCount(message.data)
            break
          case 'pong':
            // 收到心跳响应，连接正常
            console.log('WebSocket heartbeat: pong received')
            break
          // AI 相关消息类型
          case 'voice_processing':
          case 'voice_transcript':
          case 'voice_transcript_partial':
          case 'ai_response_start':
          case 'ai_response_chunk':
          case 'ai_response_done':
          case 'error':
            this.handleAIMessage(message)
            break
          default:
            console.log('Unknown message type:', message.type)
        }
      } catch (error) {
        console.error('Failed to parse WebSocket message:', error, 'line:', line)
      }
    }
  }

  /**
   * Handle new notification
   */
  handleNotification(notification) {
    const notificationStore = useNotificationStore()

    // Add notification to store
    notificationStore.addNotification(notification)

    // Vibrate device (if supported and permission granted)
    if ('vibrate' in navigator) {
      navigator.vibrate([200, 100, 200])
    }
  }

  /**
   * Handle unread count update
   */
  handleUnreadCount(data) {
    const notificationStore = useNotificationStore()
    notificationStore.updateUnreadCount(data.count)
  }

  /**
   * Handle AI-related messages
   */
  handleAIMessage(message) {
    console.log('[WebSocket] handleAIMessage received:', JSON.stringify(message))
    // Notify all registered callbacks
    const callbackCount = this.aiCallbacks.size
    console.log('[WebSocket] AI callbacks count:', callbackCount)
    this.aiCallbacks.forEach(callback => {
      try {
        callback(message)
      } catch (error) {
        console.error('AI callback error:', error)
      }
    })
  }

  /**
   * Send message to WebSocket server
   */
  send(data) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(data))
    } else {
      console.warn('WebSocket not connected, cannot send message')
    }
  }

  /**
   * Send voice message for AI processing
   * @param {string} base64Data - Base64 encoded audio data
   * @param {string} mimeType - Audio MIME type (e.g., 'audio/webm')
   * @param {Array} history - Conversation history (optional)
   */
  sendVoice(base64Data, mimeType = 'audio/webm', history = []) {
    this.send({
      type: 'voice',
      data: base64Data,
      mimeType: mimeType,
      history: history
    })
  }

  /**
   * Send text chat message for AI processing
   * @param {string} message - Text message
   * @param {Array} history - Conversation history
   */
  sendChat(message, history = []) {
    this.send({
      type: 'chat',
      message: message,
      history: history
    })
  }

  /**
   * Register callback for AI responses
   * @param {Function} callback - Callback function that receives AI messages
   * @returns {Function} Unregister function
   */
  onAiResponse(callback) {
    this.aiCallbacks.add(callback)
    // Return unregister function
    return () => {
      this.aiCallbacks.delete(callback)
    }
  }

  /**
   * Remove AI response listener
   * @param {Function} callback - Callback to remove
   */
  removeAiListener(callback) {
    this.aiCallbacks.delete(callback)
  }

  /**
   * Clear all AI response listeners
   */
  clearAiListeners() {
    this.aiCallbacks.clear()
  }

  /**
   * Start heartbeat to keep connection alive
   */
  startHeartbeat() {
    this.stopHeartbeat()
    this.heartbeatInterval = setInterval(() => {
      this.send({ type: 'ping' })
    }, 30000) // Send heartbeat every 30 seconds
  }

  /**
   * Stop heartbeat
   */
  stopHeartbeat() {
    if (this.heartbeatInterval) {
      clearInterval(this.heartbeatInterval)
      this.heartbeatInterval = null
    }
  }

  /**
   * Disconnect from WebSocket server
   */
  disconnect() {
    this.isManualClose = true
    this.stopHeartbeat()
    this.clearAiListeners()

    if (this.ws) {
      this.ws.close()
      this.ws = null
    }

    // Update notification store connection status
    const notificationStore = useNotificationStore()
    if (notificationStore) {
      notificationStore.wsConnected = false
    }
  }

  /**
   * Reconnect to WebSocket server
   */
  reconnect() {
    this.disconnect()
    this.reconnectAttempts = 0
    setTimeout(() => this.connect(), 1000)
  }

  /**
   * Check if WebSocket is connected
   */
  isConnected() {
    return this.ws && this.ws.readyState === WebSocket.OPEN
  }
}

// Create singleton instance
const webSocketService = new WebSocketService()

export default webSocketService

/**
 * Initialize WebSocket connection
 * Call this after user logs in
 */
export function initWebSocket() {
  webSocketService.connect()
}

/**
 * Disconnect WebSocket
 * Call this when user logs out
 */
export function disconnectWebSocket() {
  webSocketService.disconnect()
}

/**
 * Get WebSocket service instance
 */
export function getWebSocketService() {
  return webSocketService
}
