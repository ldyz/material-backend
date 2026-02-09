/**
 * Push Notifications Utility
 *
 * Handles push notification registration and handling for Capacitor apps
 *
 * @module PushNotifications
 * @date 2025-02-07
 */

import { PushNotifications } from '@capacitor/push-notifications'
import { LocalNotifications } from '@capacitor/local-notifications'
import { Capacitor } from '@capacitor/core'
import { registerPushToken, unregisterPushToken } from '@/api/notification'
import { storage } from './storage'

/**
 * Push notification service
 */
class PushNotificationService {
  constructor() {
    this.isRegistered = false
    this.token = null
  }

  /**
   * Initialize push notifications
   */
  async init() {
    try {
      console.log('Initializing push notifications...')

      // Check if running in Capacitor
      const isCapacitor = typeof window !== 'undefined' && window.Capacitor
      if (!isCapacitor) {
        console.log('Not running in Capacitor, skipping push notifications')
        return
      }

      // Check if PushNotifications plugin is available
      if (!PushNotifications) {
        console.log('PushNotifications plugin not available, skipping')
        return
      }

      // Request permissions
      const result = await PushNotifications.requestPermissions()
      if (result.receive === 'granted') {
        console.log('Push notification permissions granted')

        // Register for push notifications
        await this.register()

        // Set up listeners
        this.setupListeners()
      } else {
        console.warn('Push notification permissions denied')
      }

      // Configure local notifications
      await LocalNotifications.configure({
        notifications: [
          {
            id: 1,
            title: 'Notification',
            body: 'Local notification',
            largeBody: 'Large body text',
            summaryText: 'Summary',
          }
        ]
      })
    } catch (error) {
      // Don't throw - just log the error and continue
      console.error('Failed to initialize push notifications:', error)
      // Silently fail - app should work without push notifications
    }
  }

  /**
   * Register for push notifications
   */
  async register() {
    try {
      // Register with the push service
      await PushNotifications.register()

      // Listen for registration
      PushNotifications.addListener('registration', async (token) => {
        console.log('Push notification registration successful:', token)
        this.token = token.value
        this.isRegistered = true

        // Send token to server
        await this.sendTokenToServer(token.value)
      })

      // Listen for registration error
      PushNotifications.addListener('registrationError', (error) => {
        console.error('Push notification registration error:', error)
      })

    } catch (error) {
      console.error('Failed to register for push notifications:', error)
      // Don't throw - allow app to continue without push notifications
    }
  }

  /**
   * Set up push notification listeners
   */
  setupListeners() {
    // Handle incoming push notifications
    PushNotifications.addListener('pushNotificationReceived', (notification) => {
      console.log('Push notification received:', notification)
    })

    // Handle push notification action performed
    PushNotifications.addListener('pushNotificationActionPerformed', (notification) => {
      console.log('Push notification action performed:', notification)

      // Navigate based on notification data
      this.handleNotificationAction(notification)
    })

    // Handle local notification action performed
    LocalNotifications.addListener('localNotificationActionPerformed', (notification) => {
      console.log('Local notification action performed:', notification)

      // Navigate based on notification data
      this.handleNotificationAction(notification)
    })

    // Handle app state change (app opened from notification)
    document.addEventListener('visibilitychange', () => {
      if (!document.hidden) {
        // App came to foreground, check for pending notifications
        this.checkPendingNotifications()
      }
    })
  }

  /**
   * Send push token to server
   */
  async sendTokenToServer(token) {
    try {
      // Get platform info
      const platform = await this.getPlatform()
      const deviceId = await this.getDeviceId()

      await registerPushToken({
        token,
        platform,
        device_id: deviceId
      })

      // Store token locally
      storage.setPushToken(token)
      console.log('Push token registered with server')
    } catch (error) {
      console.error('Failed to register push token with server:', error)
    }
  }

  /**
   * Unregister push token from server
   */
  async unregister() {
    try {
      const token = storage.getPushToken()
      if (token) {
        await unregisterPushToken({ token })
        storage.removePushToken()
        console.log('Push token unregistered from server')
      }
    } catch (error) {
      console.error('Failed to unregister push token:', error)
    }
  }

  /**
   * Get platform type
   */
  async getPlatform() {
    const isCapacitor = typeof window !== 'undefined' && window.Capacitor
    if (!isCapacitor) return 'web'

    const platform = await Capacitor.getPlatform()
    return platform || 'web'
  }

  /**
   * Get device ID
   */
  async getDeviceId() {
    try {
      // Generate a random device ID and store it locally
      let deviceId = storage.getPushToken()
      if (!deviceId) {
        deviceId = 'device_' + Date.now() + '_' + Math.random().toString(36).substring(2, 15)
        // Use a separate storage key for device ID
        localStorage.setItem('device_id', deviceId)
      } else {
        deviceId = localStorage.getItem('device_id')
      }
      return deviceId
    } catch (error) {
      console.error('Failed to get device ID:', error)
      return null
    }
  }

  /**
   * Handle notification action (tap)
   */
  handleNotificationAction(notification) {
    try {
      const data = notification.notification?.data || notification.data

      if (data) {
        // Parse data if it's a string
        const notificationData = typeof data === 'string' ? JSON.parse(data) : data

        // Navigate to appropriate page
        this.navigateToNotification(notificationData)
      }
    } catch (error) {
      console.error('Failed to handle notification action:', error)
    }
  }

  /**
   * Navigate based on notification data
   */
  navigateToNotification(data) {
    // Import router dynamically to avoid circular dependencies
    import('@/router').then(({ default: router }) => {
      if (data.business_type && data.business_id) {
        switch (data.business_type) {
          case 'inbound_order':
            router.push(`/inbound/${data.business_id}`)
            break
          case 'requisition':
            router.push(`/requisition/${data.business_id}`)
            break
          case 'material_plan':
            router.push(`/plans/${data.business_id}`)
            break
        }
      } else {
        // Default to notifications page
        router.push('/notifications')
      }
    })
  }

  /**
   * Schedule a local notification
   */
  async scheduleLocalNotification(options) {
    try {
      await LocalNotifications.schedule({
        notifications: [
          {
            id: Date.now(),
            title: options.title || 'Notification',
            body: options.body || '',
            schedule: { at: new Date(options.when || Date.now() + 1000) },
            sound: options.sound || 'default',
            attachments: options.attachments,
            actionTypeId: '',
            extra: options.extra || null
          }
        ]
      })
    } catch (error) {
      console.error('Failed to schedule local notification:', error)
    }
  }

  /**
   * Show a local notification immediately
   */
  async showLocalNotification(options) {
    try {
      await LocalNotifications.schedule({
        notifications: [
          {
            id: Date.now(),
            title: options.title || 'Notification',
            body: options.body || '',
            sound: options.sound || 'default',
            attachments: options.attachments,
            extra: options.extra || null
          }
        ]
      })
    } catch (error) {
      console.error('Failed to show local notification:', error)
    }
  }

  /**
   * Check for pending notifications
   */
  async checkPendingNotifications() {
    try {
      const pending = await LocalNotifications.getPending()
      console.log('Pending notifications:', pending)
    } catch (error) {
      console.error('Failed to check pending notifications:', error)
    }
  }

  /**
   * Clear all notifications
   */
  async clearAllNotifications() {
    try {
      await LocalNotifications.cancel()
      await LocalNotifications.removeAllDeliveredNotifications()
    } catch (error) {
      console.error('Failed to clear notifications:', error)
    }
  }

  /**
   * Clean up when logging out
   */
  async cleanup() {
    await this.unregister()
    await this.clearAllNotifications()
    this.isRegistered = false
    this.token = null
  }
}

// Create singleton instance
const pushNotificationService = new PushNotificationService()

export default pushNotificationService

/**
 * Initialize push notifications on app start
 * Call this in main.js after app initialization
 */
export async function initPushNotifications() {
  await pushNotificationService.init()
}

/**
 * Clean up push notifications on logout
 */
export async function cleanupPushNotifications() {
  await pushNotificationService.cleanup()
}
