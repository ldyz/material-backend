import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'

// 导入 Vant 样式
import 'vant/lib/index.css'

// 导入 Capacitor App 插件
import { App as CapacitorApp } from '@capacitor/app'

const app = createApp(App)

app.use(createPinia())
app.use(router)

// 设置全局 router，供通知跳转使用
window.router = router

app.mount('#app')

// 处理 Android 返回键
if (window.Capacitor) {
  CapacitorApp.addListener('backButton', ({ canGoBack }) => {
    if (canGoBack) {
      // 如果可以后退，则导航到上一页
      router.back()
    } else {
      // 如果不能后退，确认退出
      const currentRoute = router.currentRoute.value
      // 如果在首页，提示退出
      if (currentRoute.path === '/' || currentRoute.name === 'Home') {
        // 显示确认对话框（可以使用 Vant 的 Dialog）
        // 这里直接退出应用，可以根据需要添加确认对话框
        CapacitorApp.exitApp()
      }
    }
  })
}
