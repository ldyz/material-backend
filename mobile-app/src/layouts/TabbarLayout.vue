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
          <NotificationCenter />
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
        v-if="canViewInbound"
        name="inbound"
        icon="logistics"
      >
        入库
      </van-tabbar-item>
      <van-tabbar-item
        v-if="canViewRequisition"
        name="outbound"
        icon="send-gift-o"
      >
        出库
      </van-tabbar-item>
      <van-tabbar-item
        v-if="canViewConstructionLog"
        name="construction"
        icon="notes-o"
      >
        日志
      </van-tabbar-item>
      <van-tabbar-item name="profile" icon="user-o">
        我的
      </van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { usePermission } from '@/composables/usePermission'
import NotificationCenter from '@/components/NotificationCenter.vue'

const route = useRoute()
const router = useRouter()
const { canViewInbound, canViewRequisition, canViewConstructionLog } = usePermission()

const activeTab = ref('home')

// 是否显示顶部导航栏
const showHeader = computed(() => {
  const path = route.path

  // 详情页、创建页、审批页不显示顶部导航（这些页面有自己的导航栏）
  const shouldHide =
    // 详情页：路径中包含单个数字ID的
    /^\/inbound\/\d+$/.test(path) ||
    /^\/outbound\/\d+$/.test(path) ||
    /^\/construction\/\d+$/.test(path) ||
    /^\/stock\/\d+$/.test(path) ||
    // 创建页
    path.endsWith('/create') ||
    // 审批页
    path.includes('/approve') ||
    // 编辑页
    path.includes('/edit')

  // 首页、列表页、我的页面显示顶部导航
  const showHeaderRoutes = ['/', '/home', '/inbound', '/outbound', '/profile', '/construction', '/stock']

  return !shouldHide && showHeaderRoutes.some(p => path === p || path.startsWith(p + '/'))
})

// 页面标题
const pageTitle = computed(() => {
  return route.meta.title || import.meta.env.VITE_APP_TITLE || '材料管理'
})

// 是否显示返回按钮
const showBack = computed(() => {
  // 非首页显示返回按钮
  return route.path !== '/'
})

// 返回上一页
function onClickLeft() {
  if (window.history.state.back) {
    router.back()
  } else {
    router.push('/')
  }
}

// 根据路由更新当前 tab
watch(
  () => route.path,
  (path) => {
    if (path === '/' || path.startsWith('/home')) {
      activeTab.value = 'home'
    } else if (path.startsWith('/inbound')) {
      activeTab.value = 'inbound'
    } else if (path.startsWith('/outbound')) {
      activeTab.value = 'outbound'
    } else if (path.startsWith('/construction')) {
      activeTab.value = 'construction'
    } else if (path.startsWith('/profile')) {
      activeTab.value = 'profile'
    }
  },
  { immediate: true }
)

// 切换 tab
function onTabChange(name) {
  router.push(`/${name === 'home' ? '' : name}`)
}
</script>

<style scoped>
.tabbar-layout {
  min-height: 100vh;
  background-color: #f7f8fa;
}
</style>
