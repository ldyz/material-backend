import { ref, onMounted } from 'vue'
import { App } from '@capacitor/app'
import { Device } from '@capacitor/device'
import { Filesystem } from '@capacitor/filesystem'
import { showToast } from 'vant'
import request from '@/utils/request'
import { logger } from '@/utils/logger'

const isChecking = ref(false)
const hasUpdate = ref(false)
const latestVersion = ref('')
const downloadUrl = ref('')
const currentVersion = ref('')
const forceUpdate = ref(false)
const updateMessage = ref('')
const isDownloading = ref(false)
const downloadProgress = ref(0)

// 检测是否在浏览器环境中
// 注意：Capacitor 环境中 protocol 是 'file:'，需要通过 Capacitor API 来判断
const isWeb = ref(
  typeof window !== 'undefined' &&
  !window.Capacitor && // 没有 Capacitor 对象则是纯 Web 环境
  (window.location.protocol === 'http:' || window.location.protocol === 'https:')
)

// 获取当前应用版本
async function getCurrentVersion() {
  logger.log('[更新检测] 开始获取当前版本...')
  logger.log('[更新检测] window.Capacitor:', window.Capacitor)
  logger.log('[更新检测] isWeb.value:', isWeb.value)
  logger.log('[更新检测] protocol:', window.location?.protocol)

  // 在 Web 环境中，从专门的版本文件读取
  if (isWeb.value) {
    logger.log('[更新检测] 使用 Web 环境，从 version.json 读取')
    try {
      // 尝试从版本文件读取
      const response = await fetch('/version.json', { cache: 'no-store' })
      const versionData = await response.json()
      const version = versionData.version || '1.0.3'
      currentVersion.value = version
      logger.log('[更新检测] 从 version.json 读取到版本:', version)
      return version
    } catch (error) {
      logger.error('[更新检测] 无法读取版本文件，使用默认版本 1.0.3:', error)
      currentVersion.value = '1.0.3'
      return '1.0.3'
    }
  }

  // 原生环境，使用 Capacitor App API
  logger.log('[更新检测] 使用原生环境，调用 App.getInfo()')
  try {
    const info = await App.getInfo()
    logger.log('[更新检测] App.getInfo() 返回:', info)
    currentVersion.value = info.version
    logger.log('[更新检测] 获取到版本号:', info.version)
    return info.version
  } catch (error) {
    logger.error('[更新检测] 获取应用版本失败:', error)
    return '1.0.3'
  }
}

// 检查更新
async function checkUpdate() {
  if (isChecking.value) return

  isChecking.value = true
  logger.log('[更新检测] ========== 开始检查更新 ==========')

  try {
    // 获取当前版本
    await getCurrentVersion()
    logger.log('[更新检测] 当前版本:', currentVersion.value)

    // 确定平台类型
    let platform = 'android' // 默认为 android
    if (!isWeb.value) {
      try {
        const deviceInfo = await Device.getInfo()
        platform = deviceInfo.platform === 'android' ? 'android' : 'ios'
        logger.log('[更新检测] 设备平台:', deviceInfo.platform)
      } catch (error) {
        logger.warn('[更新检测] 无法获取设备信息，使用默认平台:', error)
      }
    }

    // 从后端获取最新版本信息
    logger.log('[更新检测] 请求参数:', { platform, current_version: currentVersion.value })
    const response = await request.get('/app/version', {
      params: {
        platform: platform,
        current_version: currentVersion.value
      }
    })

    logger.log('[更新检测] API 响应:', response)
    const data = response.data

    if (data.has_update) {
      logger.log('[更新检测] ✓ 发现新版本:', data.latest_version)
      latestVersion.value = data.latest_version
      downloadUrl.value = data.download_url
      forceUpdate.value = data.force_update || false
      updateMessage.value = data.update_message || '发现新版本，请更新以获得更好的体验'
      hasUpdate.value = true
      return {
        hasUpdate: true,
        version: data.latest_version,
        forceUpdate: data.force_update,
        message: data.update_message,
        downloadUrl: data.download_url
      }
    }

    logger.log('[更新检测] × 已是最新版本')
    hasUpdate.value = false
    return { hasUpdate: false }
  } catch (error) {
    logger.error('检查更新失败:', error)
    hasUpdate.value = false
    return { hasUpdate: false }
  } finally {
    isChecking.value = false
  }
}

// 下载并安装更新（仅 Android）
async function downloadAndInstall() {
  logger.log('[更新检测] ========== 开始下载 APK ==========')

  try {
    // 在 Web 环境中，直接打开下载链接
    if (isWeb.value) {
      logger.log('[更新检测] Web 环境，打开下载链接')
      if (downloadUrl.value) {
        window.open(downloadUrl.value, '_blank')
      }
      return
    }

    const deviceInfo = await Device.getInfo()
    logger.log('[更新检测] 设备平台:', deviceInfo.platform)

    if (deviceInfo.platform !== 'android') {
      logger.log('[更新检测] 非 Android 设备，打开浏览器')
      if (downloadUrl.value) {
        window.open(downloadUrl.value, '_blank')
      }
      return
    }

    // Android 使用系统下载管理器（通过浏览器）
    if (!downloadUrl.value) {
      logger.error('[更新检测] 下载链接为空')
      return
    }

    logger.log('[更新检测] 使用系统下载器:', downloadUrl.value)

    // 显示提示
    showToast({
      type: 'success',
      message: '正在打开浏览器下载...',
      duration: 2000
    })

    // 延迟打开，让用户看到提示
    setTimeout(() => {
      // 使用系统浏览器下载
      const link = document.createElement('a')
      link.href = downloadUrl.value
      link.target = '_blank'
      link.download = `material-management-${latestVersion.value}.apk`

      // 添加到 DOM 并点击
      document.body.appendChild(link)
      link.click()

      // 清理
      setTimeout(() => {
        document.body.removeChild(link)
      }, 100)

      logger.log('[更新检测] 已触发浏览器下载')

      // 提示用户查看下载进度
      setTimeout(() => {
        showToast({
          type: 'success',
          message: '请在通知栏查看下载进度',
          duration: 3000
        })
      }, 500)
    }, 500)

  } catch (error) {
    logger.error('[更新检测] 下载失败:', error)
    logger.error('[更新检测] 错误详情:', {
      message: error.message,
      stack: error.stack,
      name: error.name
    })
    showToast({ type: 'fail', message: `下载失败: ${error.message}` })

    // 降级方案
    setTimeout(() => {
      if (downloadUrl.value) {
        window.open(downloadUrl.value, '_blank')
      }
    }, 1500)
  }
}

// 将 ArrayBuffer 转换为 base64
function arrayBufferToBase64(buffer) {
  let binary = ''
  const bytes = new Uint8Array(buffer)
  const len = bytes.byteLength
  for (let i = 0; i < len; i++) {
    binary += String.fromCharCode(bytes[i])
  }
  return window.btoa(binary)
}

// 跳转到应用商店
function openStore() {
  if (downloadUrl.value) {
    window.open(downloadUrl.value, '_blank')
  }
}

// 稍后提醒
function remindLater() {
  hasUpdate.value = false
}

// 跳过此版本
function skipVersion() {
  const skippedVersion = latestVersion.value
  localStorage.setItem(`skipped_version_${currentVersion.value}`, skippedVersion)
  hasUpdate.value = false
}

// 检查是否已跳过此版本
function hasSkippedVersion() {
  const skippedVersion = localStorage.getItem(`skipped_version_${currentVersion.value}`)
  return skippedVersion === latestVersion.value
}

export function useAppUpdate() {
  return {
    isChecking,
    hasUpdate,
    latestVersion,
    currentVersion,
    downloadUrl,
    forceUpdate,
    updateMessage,
    isDownloading,
    downloadProgress,
    checkUpdate,
    downloadAndInstall,
    openStore,
    remindLater,
    skipVersion,
    hasSkippedVersion
  }
}

// 自动检查更新的 hook
export function useAutoUpdate(options = {}) {
  const {
    autoCheck = true,
    checkOnMount = true,
    checkInterval = 24 * 60 * 60 * 1000 // 24小时
  } = options

  const {
    isChecking,
    hasUpdate,
    forceUpdate,
    checkUpdate
  } = useAppUpdate()

  let checkTimer = null

  // 执行更新检查
  async function performCheck() {
    logger.log('[自动更新] ========== 执行更新检查 ==========')
    logger.log('[自动更新] isChecking:', isChecking.value)

    // 执行检查
    const result = await checkUpdate()

    // 如果已跳过此版本且不是强制更新，则隐藏更新提示
    if (result && result.hasUpdate && !result.forceUpdate) {
      if (hasSkippedVersion()) {
        logger.log('[自动更新] 用户已跳过此版本，不显示提示')
        hasUpdate.value = false
      }
    }

    return result
  }

  // 启动定时检查
  function startPeriodicCheck() {
    if (checkTimer) {
      clearInterval(checkTimer)
    }

    if (!autoCheck) {
      logger.log('[自动更新] 自动检查已禁用')
      return
    }

    logger.log('[自动更新] 启动定时检查，间隔:', checkInterval / 1000 / 60 / 60, '小时')
    checkTimer = setInterval(() => {
      logger.log('[自动更新] 定时检查触发')
      performCheck()
    }, checkInterval)
  }

  // 停止定时检查
  function stopPeriodicCheck() {
    if (checkTimer) {
      clearInterval(checkTimer)
      checkTimer = null
      logger.log('[自动更新] 定时检查已停止')
    }
  }

  // 应用从后台恢复时检查
  function setupAppStateListener() {
    // 监听应用状态变化（需要 Capacitor App 插件）
    if (window.Capacitor && window.Capacitor.App) {
      window.Capacitor.App.addListener('appStateChange', (state) => {
        if (state.isActive) {
          logger.log('[自动更新] 应用恢复活动状态，检查更新')
          // 延迟 3 秒检查，避免频繁检查
          setTimeout(() => {
            performCheck()
          }, 3000)
        }
      })
    }
  }

  onMounted(() => {
    logger.log('[自动更新] ========== 组件已挂载 ==========')
    logger.log('[自动更新] autoCheck:', autoCheck)
    logger.log('[自动更新] checkOnMount:', checkOnMount)

    if (checkOnMount) {
      // 延迟 2 秒执行首次检查，等待应用完全加载
      setTimeout(() => {
        logger.log('[自动更新] 执行首次检查')
        performCheck()
      }, 2000)
    }

    startPeriodicCheck()
    setupAppStateListener()
  })

  return {
    isChecking,
    hasUpdate,
    forceUpdate,
    performCheck,
    startPeriodicCheck,
    stopPeriodicCheck
  }
}
