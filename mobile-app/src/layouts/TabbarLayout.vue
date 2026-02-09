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
      <van-tabbar-item name="dashboard" icon="wap-home-o">首页</van-tabbar-item>
      <van-tabbar-item name="plans" icon="orders-o">计划</van-tabbar-item>
      <van-tabbar-item name="appointments" icon="calendar-o">预约</van-tabbar-item>
      <van-tabbar-item name="inbound" icon="logistics">入库</van-tabbar-item>
      <van-tabbar-item name="profile" icon="user-o">我的</van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

const activeTab = ref('dashboard')

watch(
  () => route.path,
  (path) => {
    if (path === '/' || path.startsWith('/dashboard')) {
      activeTab.value = 'dashboard'
    } else if (path.startsWith('/plans')) {
      activeTab.value = 'plans'
    } else if (path.startsWith('/appointment')) {
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
