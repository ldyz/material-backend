<template>
  <div class="tabbar-layout">
    <!-- 顶部导航栏 -->
    <van-sticky v-if="showHeader">
      <van-nav-bar
        :title="pageTitle"
        :left-arrow="showBack"
        @click-left="onClickLeft"
      >
        <template #right>
          <van-icon
            :badge="unreadCount > 0 ? unreadCount : null"
            name="bell"
            size="20"
            @click="goToNotifications"
          />
        </template>
      </van-nav-bar>
    </van-sticky>

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

    <van-tabbar
      v-model="activeTab"
      :safe-area-inset-bottom="true"
      @change="onTabChange"
    >
      <van-tabbar-item name="home" icon="wap-home-o">
        首页
      </van-tabbar-item>
      <van-tabbar-item
        name="plans"
        icon="orders-o"
      >
        计划
      </van-tabbar-item>
      <van-tabbar-item
        name="materials"
        icon="apps-o"
      >
        物资
      </van-tabbar-item>
      <van-tabbar-item
        name="tasks"
        icon="todo-list-o"
        :badge="pendingTaskCount > 0 ? pendingTaskCount : null"
      >
        待办
      </van-tabbar-item>
      <van-tabbar-item name="profile" icon="user-o">
        我的
      </van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useNotification } from '@/composables/useNotification'

const route = useRoute()
const router = useRouter()
const { unreadCount } = useNotification()

const activeTab = ref('home')
const pendingTaskCount = ref(0)

// 是否显示顶部导航栏
const showHeader = computed(() => {
  const path = route.path

  // 详情页、创建页、审批页不显示顶部导航
  const shouldHide =
    /^\/(plans|inbound|outbound|materials)\/\d+$/.test(path) ||
    path.endsWith('/create') ||
    path.includes('/approve') ||
    path.includes('/edit') ||
    path.includes('/issue')

  // 显示顶部导航的路由
  const showHeaderRoutes = ['/', '/home', '/plans', '/inbound', '/outbound', '/materials', '/stock', '/tasks', '/profile']

  return !shouldHide && showHeaderRoutes.some(p => path === p || path.startsWith(p + '/'))
})

// 页面标题
const pageTitle = computed(() => {
  return route.meta.title || import.meta.env.VITE_APP_TITLE || '材料管理'
})

// 是否显示返回按钮
const showBack = computed(() => {
  return route.path !== '/' && route.path !== '/home'
})

// 返回上一页
function onClickLeft() {
  if (window.history.state.back) {
    router.back()
  } else {
    router.push('/')
  }
}

// 前往通知中心
function goToNotifications() {
  router.push('/notifications')
}

// 加载待办数量
async function loadPendingCount() {
  try {
    // 从 localStorage 或 API 获取待办数量
    const stored = localStorage.getItem('pending_task_count')
    if (stored) {
      pendingTaskCount.value = parseInt(stored, 10)
    }
  } catch (error) {
    console.error('加载待办数量失败:', error)
  }
}

// 根据路由更新当前 tab
watch(
  () => route.path,
  (path) => {
    if (path === '/' || path.startsWith('/home')) {
      activeTab.value = 'home'
    } else if (path.startsWith('/plans')) {
      activeTab.value = 'plans'
    } else if (path.startsWith('/materials') || path.startsWith('/stock')) {
      activeTab.value = 'materials'
    } else if (path.startsWith('/tasks')) {
      activeTab.value = 'tasks'
    } else if (path.startsWith('/profile')) {
      activeTab.value = 'profile'
    }
  },
  { immediate: true }
)

// 切换 tab
function onTabChange(name) {
  if (name === 'home') {
    router.push('/')
  } else {
    router.push(`/${name}`)
  }
}

onMounted(() => {
  loadPendingCount()
  // 定期更新待办数量
  setInterval(loadPendingCount, 30000)
})
</script>

<style scoped>
.tabbar-layout {
  min-height: 100vh;
  background-color: #f7f8fa;
  padding-bottom: env(safe-area-inset-bottom);
}
</style>
