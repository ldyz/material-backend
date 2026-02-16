<template>
  <div class="tabbar-layout">
    <router-view v-slot="{ Component, route }">
      <keep-alive>
        <component
          :is="Component"
          v-if="route.meta.keepAlive"
          :key="route.path"
        />
      </keep-alive>
      <component
        :is="Component"
        v-if="!route.meta.keepAlive"
        :key="route.path"
      />
    </router-view>

    <van-tabbar v-model="activeTab" :safe-area-inset-bottom="true">
      <!-- 动态渲染的菜单项 -->
      <template v-for="item in menuItems" :key="item.name">
        <van-tabbar-item :name="item.name" :icon="item.icon">
          {{ item.label }}
        </van-tabbar-item>
      </template>
    </van-tabbar>

    <!-- 自动更新对话框 -->
    <AppUpdateDialog v-model:show="showUpdateDialog" />
  </div>
</template>

<script setup>
import { ref, watch, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAutoUpdate } from '@/composables/useAppUpdate'
import AppUpdateDialog from '@/components/AppUpdateDialog.vue'

const router = useRouter()
const route = useRoute()

const activeTab = ref('dashboard')
const showUpdateDialog = ref(false)

// 使用自动更新检测
const { hasUpdate, forceUpdate, performCheck } = useAutoUpdate({
  autoCheck: true,
  checkOnMount: true,
  checkInterval: 24 * 60 * 60 * 1000 // 每24小时检查一次
})

// 从 localStorage 获取用户信息
const userInfo = ref(JSON.parse(localStorage.getItem('userInfo') || '{}'))
const userRole = computed(() => userInfo.value.role || '')

// 根据角色动态生成菜单项
const menuItems = computed(() => {
  const role = userRole.value

  // 默认菜单（所有人可见）
  const defaultMenus = [
    { name: 'dashboard', label: '首页', icon: 'wap-home-o' },
    { name: 'appointments', label: '预约', icon: 'calendar-o' },
    { name: 'profile', label: '我的', icon: 'user-o' }
  ]

  // 管理员菜单
  if (role === 'admin' || role === 'project_manager') {
    return [
      { name: 'dashboard', label: '首页', icon: 'wap-home-o' },
      { name: 'appointments', label: '预约管理', icon: 'calendar-o' },
      { name: 'inbound', label: '入库管理', icon: 'logistics' },
      { name: 'profile', label: '我的', icon: 'user-o' }
    ]
  }

  // 作业人员菜单
  if (role === 'worker') {
    return [
      { name: 'dashboard', label: '首页', icon: 'wap-home-o' },
      { name: 'profile', label: '我的', icon: 'user-o' }
    ]
  }

  // 默认返回普通菜单
  return defaultMenus
})

// 监听是否有更新
watch([hasUpdate, forceUpdate], ([hasUpdate, forceUpdate]) => {
  if (hasUpdate) {
    // 如果是强制更新，立即显示
    // 如果是可选更新，延迟显示
    if (forceUpdate) {
      showUpdateDialog.value = true
    } else {
      setTimeout(() => {
        showUpdateDialog.value = true
      }, 3000) // 3秒后显示
    }
  }
})

watch(
  () => route.path,
  (path) => {
    if (path === '/' || path.startsWith('/dashboard')) {
      activeTab.value = 'dashboard'
    } else if (path.startsWith('/plans')) {
      activeTab.value = 'plans'
    } else if (path.startsWith('/appointments')) {
      activeTab.value = 'appointments'
    } else if (path.startsWith('/inbound')) {
      activeTab.value = 'inbound'
    } else if (path.startsWith('/requisition')) {
      activeTab.value = 'requisition'
    } else if (path.startsWith('/profile')) {
      activeTab.value = 'profile'
    }
  },
  { immediate: true }
)

watch(activeTab, (name, oldName) => {
  if (name === oldName) return

  const targetPath = name === 'dashboard' ? '/' : `/${name}`
  if (route.path !== targetPath) {
    router.push(targetPath)
  }
})

// 在"我的"页面添加手动检查更新
onMounted(() => {
  // 暴露全局方法供其他页面调用
  window.checkAppUpdate = performCheck
})
</script>

<style scoped>
.tabbar-layout {
  min-height: 100vh;
  background-color: #f7f8fa;
  padding-bottom: 50px;
  padding-top: var(--capacitor-status-bar-height, 0px);
}

/* Capacitor 状态栏高度变量 */
:global(html) {
  --capacitor-status-bar-height: 0px;
}

/* 在 Android 原生环境中添加状态栏高度 */
@supports (padding: max(0px)) {
  .tabbar-layout {
    padding-top: max(var(--capacitor-status-bar-height), env(safe-area-inset-top));
  }
}
</style>
