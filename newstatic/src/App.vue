<template>
  <router-view />
</template>

<script setup>
import { onMounted } from 'vue'
import { useAppStore } from '@/stores/app'

// App 根组件
const appStore = useAppStore()

// 应用初始化时加载系统配置
onMounted(async () => {
  try {
    await appStore.loadSystemConfig()
  } catch (error) {
    console.error('加载系统配置失败:', error)
  }
})
</script>

<style>
#app {
  width: 100%;
  min-height: 100vh;
}
</style>
