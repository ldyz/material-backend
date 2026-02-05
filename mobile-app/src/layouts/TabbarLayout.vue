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
      <van-tabbar-item name="inbound" icon="logistics">入库</van-tabbar-item>
      <van-tabbar-item name="requisition" icon="send-gift-o">出库</van-tabbar-item>
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

watch(activeTab, (name) => {
  if (name === 'dashboard') {
    router.push('/')
  } else {
    router.push(`/${name}`)
  }
})
</script>

<style scoped>
.tabbar-layout {
  min-height: 100vh;
  background-color: #f7f8fa;
  padding-bottom: 50px;
}
</style>
