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

    <!-- 自定义底部导航栏 -->
    <div class="custom-tabbar" :class="{ 'no-voice': !hasAIPermission }">
      <!-- 首页按钮 -->
      <div
        class="tabbar-item"
        :class="{ active: activeTab === 'dashboard' }"
        @click="handleTabClick('dashboard')"
      >
        <van-icon name="wap-home-o" size="22" />
        <span class="tabbar-label">首页</span>
      </div>

      <!-- 中间 AI 助手按钮 - 仅对有AI权限的用户显示 -->
      <div v-if="hasAIPermission" class="ai-button-wrapper">
        <div class="ai-button" @click="openAiDialog">
          <van-icon name="chat-o" size="24" />
        </div>
        <span class="ai-label">AI助手</span>
      </div>

      <!-- 我的按钮 -->
      <div
        class="tabbar-item"
        :class="{ active: activeTab === 'profile' }"
        @click="handleTabClick('profile')"
      >
        <van-icon name="user-o" size="22" />
        <span class="tabbar-label">我的</span>
      </div>
    </div>

    <!-- AI 聊天弹窗 -->
    <AiChatPopup
      v-model:show="showAiDialog"
      context="dashboard"
    />

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
import AiChatPopup from '@/components/AiChatPopup.vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const activeTab = ref('dashboard')
const showUpdateDialog = ref(false)
const showAiDialog = ref(false)

// 检查是否有AI权限
const hasAIPermission = computed(() => {
  // 管理员有所有权限
  if (authStore.isAdmin) return true
  // 检查是否有AI相关权限
  const permissions = authStore.permissions
  return permissions.some(p =>
    p === 'ai_agent_query' ||
    p === 'ai_agent_operate' ||
    p === 'ai_agent_workflow' ||
    p === 'ai_agent_view'
  )
})

// 使用自动更新检测
const { hasUpdate, forceUpdate, performCheck } = useAutoUpdate({
  autoCheck: true,
  checkOnMount: true,
  checkInterval: 24 * 60 * 60 * 1000
})

// 处理标签点击
function handleTabClick(name) {
  if (name === activeTab.value) return
  activeTab.value = name
  const targetPath = name === 'dashboard' ? '/' : `/${name}`
  if (route.path !== targetPath) {
    router.push(targetPath)
  }
}

// 打开 AI 对话框
function openAiDialog() {
  showAiDialog.value = true
}

// 监听是否有更新
watch([hasUpdate, forceUpdate], ([hasUpdate, forceUpdate]) => {
  if (hasUpdate) {
    if (forceUpdate) {
      showUpdateDialog.value = true
    } else {
      setTimeout(() => {
        showUpdateDialog.value = true
      }, 3000)
    }
  }
})

watch(
  () => route.path,
  (path) => {
    if (path === '/' || path.startsWith('/dashboard')) {
      activeTab.value = 'dashboard'
    } else if (path.startsWith('/profile')) {
      activeTab.value = 'profile'
    }
  },
  { immediate: true }
)

onMounted(() => {
  window.checkAppUpdate = performCheck
})
</script>

<style scoped>
.tabbar-layout {
  min-height: 100vh;
  background-color: #f7f8fa;
  padding-bottom: 70px;
  padding-top: var(--capacitor-status-bar-height, 0px);
}

/* 自定义底部导航栏 */
.custom-tabbar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 60px;
  background: #fff;
  display: flex;
  justify-content: space-around;
  align-items: center;
  box-shadow: 0 -2px 12px rgba(0, 0, 0, 0.08);
  padding-bottom: env(safe-area-inset-bottom);
  z-index: 100;
}

.tabbar-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 8px 0;
  cursor: pointer;
  color: #969799;
  transition: color 0.2s;
}

.tabbar-item.active {
  color: #1989fa;
}

.tabbar-label {
  font-size: 12px;
  margin-top: 4px;
}

/* 无语音按钮时的布局 */
.custom-tabbar.no-voice {
  justify-content: center;
}

.custom-tabbar.no-voice .tabbar-item {
  flex: none;
  min-width: 120px;
}

/* AI 助手按钮 */
.ai-button-wrapper {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-top: -25px;
}

.ai-button {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: linear-gradient(135deg, #1989fa 0%, #0d7ce9 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 12px rgba(25, 137, 250, 0.4);
  cursor: pointer;
  transition: all 0.2s;
  border: 3px solid #fff;
  color: #fff;
}

.ai-button:active {
  transform: scale(0.95);
}

.ai-label {
  font-size: 10px;
  color: #969799;
  margin-top: 4px;
}

/* Capacitor 状态栏高度变量 */
:global(html) {
  --capacitor-status-bar-height: 0px;
}

@supports (padding: max(0px)) {
  .tabbar-layout {
    padding-top: max(var(--capacitor-status-bar-height), env(safe-area-inset-top));
  }
}
</style>
