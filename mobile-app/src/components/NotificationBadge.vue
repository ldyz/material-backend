<template>
  <div class="notification-badge" @click="handleClick">
    <van-badge :content="unreadCount" :max="99" :show="unreadCount > 0">
      <van-icon name="bell" size="20" color="#646566" />
    </van-badge>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useNotificationStore } from '@/stores/notification'

const router = useRouter()
const notificationStore = useNotificationStore()
const { unreadCount } = storeToRefs(notificationStore)

const handleClick = () => {
  router.push('/notifications')
}
</script>

<style scoped>
.notification-badge {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 8px;
  cursor: pointer;
}

:deep(.van-badge) {
  transform: translate(4px, -4px);
}
</style>
