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
import { useAuthStore } from '@/stores/auth'
import AppUpdateDialog from '@/components/AppUpdateDialog.vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const activeTab = ref('dashboard')
const showUpdateDialog = ref(false)

// 使用自动更新检测
const { hasUpdate, forceUpdate, performCheck } = useAutoUpdate({
  autoCheck: true,
  checkOnMount: true,
  checkInterval: 24 * 60 * 60 * 1000 // 每24小时检查一次
})

// 根据权限动态生成菜单项
const menuItems = computed(() => {
  const items = [
    // 首页 - 所有人可见
    { name: 'dashboard', label: '首页', icon: 'wap-home-o', permissions: [] }
  ]

  // 预约管理 - 需要预约查看权限
  if (authStore.hasPermission('appointment_view')) {
    items.push({
      name: 'appointments',
      label: '预约',
      icon: 'calendar-o',
      permissions: ['appointment_view']
    })
  }

  // 物资计划 - 需要物资计划查看权限
  if (authStore.hasPermission('material_plan_view')) {
    items.push({
      name: 'plans',
      label: '计划',
      icon: 'todo-list-o',
      permissions: ['material_plan_view']
    })
  }

  // 入库管理 - 需要入库查看权限
  if (authStore.hasPermission('inbound_view')) {
    items.push({
      name: 'inbound',
      label: '入库',
      icon: 'logistics',
      permissions: ['inbound_view']
    })
  }

  // 出库管理 - 需要出库查看权限
  if (authStore.hasPermission('requisition_view')) {
    items.push({
      name: 'requisition',
      label: '出库',
      icon: 'send-gift-o',
      permissions: ['requisition_view']
    })
  }

  // 施工日志 - 需要施工日志查看权限
  if (authStore.hasAnyPermission(['constructionlog_view', 'constructionlog_create'])) {
    items.push({
      name: 'construction-log',
      label: '日志',
      icon: 'notes-o',
      permissions: ['constructionlog_view', 'constructionlog_create']
    })
  }

  // 进度管理 - 需要进度查看权限
  if (authStore.hasPermission('progress_view')) {
    items.push({
      name: 'progress',
      label: '进度',
      icon: 'bar-chart-o',
      permissions: ['progress_view']
    })
  }

  // 我的 - 所有人可见
  items.push({ name: 'profile', label: '我的', icon: 'user-o', permissions: [] })

  return items
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
    } else if (path.startsWith('/construction-log')) {
      activeTab.value = 'construction-log'
    } else if (path.startsWith('/progress')) {
      activeTab.value = 'progress'
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
