import { ref, onMounted } from 'vue'
import { App } from '@capacitor/app'
import { Device } from '@capacitor/device'
import axios from 'axios'

const isChecking = ref(false)
const hasUpdate = ref(false)
const latestVersion = ref('')
const downloadUrl = ref('')
const currentVersion = ref('')
const forceUpdate = ref(false)
const updateMessage = ref('')

// 检测是否在浏览器环境中
const isWeb = ref(
  typeof window !== 'undefined' &&
  (window.location.protocol === 'http:' || window.location.protocol === 'https:' ||
   window.location.hostname === 'localhost' || window.location.hostname.startsWith('192.168'))
)

// 获取当前应用版本
// 获取当前应用版本
async function getCurrentVersion() {
  // 在 Web 环境中，从 package.json 读取版本
  if (isWeb.value) {
    try {
      // 尝试从 package.json 读取版本
      const response = await fetch('/package.json', { cache: 'no-store' })
      const packageJson = await response.json()
      const version = packageJson.version || '1.0.3'
      currentVersion.value = version
      return version
    } catch (error) {
      console.error('无法读取 package.json，使用默认版本:', error)
      currentVersion.value = '1.0.3'
      return '1.0.3'
    }
  }

  try {
    const info = await App.getInfo()
    currentVersion.value = info.version
    return info.version
  } catch (error) {
    console.error('获取应用版本失败:', error)
    return '1.0.3'
  }
}

// 检查更新
async function checkUpdate() {
  if (isChecking.value) return

  isChecking.value = true

  try {
    // 获取当前版本
    await getCurrentVersion()

    // 确定平台类型
    let platform = 'android' // 默认为 android
    if (!isWeb.value) {
      try {
        const deviceInfo = await Device.getInfo()
        platform = deviceInfo.platform === 'android' ? 'android' : 'ios'
      } catch (error) {
        console.warn('无法获取设备信息，使用默认平台:', error)
      }
    }

    // 从后端获取最新版本信息
    const response = await axios.get('/api/app/version', {
      params: {
        platform: platform,
        current_version: currentVersion.value
      }
    })

    const data = response.data.data

    if (data.has_update) {
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

    hasUpdate.value = false
    return { hasUpdate: false }
  } catch (error) {
    console.error('检查更新失败:', error)
    hasUpdate.value = false
    return { hasUpdate: false }
  } finally {
    isChecking.value = false
  }
}

// 下载并安装更新（仅 Android）
async function downloadAndInstall() {
  try {
    // 在 Web 环境中，直接打开下载链接
    if (isWeb.value) {
      if (downloadUrl.value) {
        window.open(downloadUrl.value, '_blank')
      }
      return
    }

    const deviceInfo = await Device.getInfo()

    if (deviceInfo.platform !== 'android') {
      // iOS 需要跳转到 App Store
      if (downloadUrl.value) {
        window.open(downloadUrl.value, '_blank')
      }
      return
    }

    // Android 直接下载 APK
    if (downloadUrl.value) {
      // 创建一个隐藏的 a 标签来下载
      const link = document.createElement('a')
      link.href = downloadUrl.value
      link.download = `material-management-${latestVersion.value}.apk`
      document.body.appendChild(link)

      // 添加调试日志
      console.log('开始下载 APK:', downloadUrl.value)
      console.log('下载文件名:', `material-management-${latestVersion.value}.apk`)

      // 尝试点击下载
      link.click()

      // 等待一段时间后检查
      setTimeout(() => {
        document.body.removeChild(link)
        console.log('下载链接已触发')
      }, 500)
    }
  } catch (error) {
    console.error('下载更新失败:', error)
  }
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
    // 如果已跳过此版本且不是强制更新，则不检查
    if (!forceUpdate.value && hasSkippedVersion()) {
      return
    }

    await checkUpdate()
  }

  // 启动定时检查
  function startPeriodicCheck() {
    if (checkTimer) {
      clearInterval(checkTimer)
    }

    if (!autoCheck) return

    checkTimer = setInterval(() => {
      performCheck()
    }, checkInterval)
  }

  // 停止定时检查
  function stopPeriodicCheck() {
    if (checkTimer) {
      clearInterval(checkTimer)
      checkTimer = null
    }
  }

  onMounted(() => {
    if (checkOnMount) {
      performCheck()
    }
    startPeriodicCheck()
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
