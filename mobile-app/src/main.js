import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import { StatusBar, StatusBarStyle } from '@capacitor/status-bar'

// 导入 Vant 样式
import 'vant/lib/index.css'

// 导入 WebSocket 推送通知
import { initWebSocket } from '@/utils/websocket'

// 导入 setRouter 以避免循环依赖
import { setRouter } from '@/utils/request'
import { Capacitor } from '@capacitor/core'

// 全局错误处理
window.addEventListener('error', (event) => {
  console.error('Global error:', event.error)
})

window.addEventListener('unhandledrejection', (event) => {
  console.error('Unhandled promise rejection:', event.reason)
})

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)

// 设置 router 实例到 request 模块
setRouter(router)

// 添加全局错误处理器
app.config.errorHandler = (err, instance, info) => {
  console.error('Vue error:', err, info)
}

app.mount('#app')

// 设置状态栏
async function setupStatusBar() {
  if (Capacitor.isNativePlatform()) {
    try {
      // 获取状态栏信息
      const info = await StatusBar.getInfo()
      console.log('Status bar height:', info.height)

      // 设置 CSS 变量
      document.documentElement.style.setProperty('--capacitor-status-bar-height', `${info.height}px`)

      // 设置状态栏样式为浅色（深色文字）
      await StatusBar.setStyle({ style: StatusBarStyle.Light })

      // 隐藏状态栏覆盖（可选）
      await StatusBar.setOverlaysWebView({ overlay: false })
    } catch (error) {
      console.error('Failed to setup status bar:', error)
    }
  }
}

// Initialize WebSocket notifications after app mount
initWebSocket()

// 设置状态栏
setupStatusBar()
