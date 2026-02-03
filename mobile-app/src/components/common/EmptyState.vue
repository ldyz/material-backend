<template>
  <div class="empty-state">
    <svg class="empty-image" viewBox="0 0 200 200" xmlns="http://www.w3.org/2000/svg">
      <!-- 背景圆 -->
      <circle cx="100" cy="100" r="80" fill="#f7f8fa" />

      <!-- 文档图标 -->
      <rect x="70" y="50" width="60" height="80" rx="4" fill="#dcdee0" />
      <rect x="80" y="70" width="40" height="4" rx="2" fill="#ffffff" />
      <rect x="80" y="80" width="40" height="4" rx="2" fill="#ffffff" />
      <rect x="80" y="90" width="30" height="4" rx="2" fill="#ffffff" />

      <!-- 放大镜 -->
      <circle cx="120" cy="120" r="20" fill="none" stroke="#969799" stroke-width="3" />
      <line x1="135" y1="135" x2="150" y2="150" stroke="#969799" stroke-width="3" stroke-linecap="round" />
    </svg>

    <p class="empty-text">{{ text }}</p>

    <van-button
      v-if="showAction"
      type="primary"
      size="small"
      @click="handleAction"
    >
      {{ actionText }}
    </van-button>
  </div>
</template>

<script setup>
defineProps({
  text: {
    type: String,
    default: '暂无数据',
  },
  showAction: {
    type: Boolean,
    default: false,
  },
  actionText: {
    type: String,
    default: '去添加',
  },
})

const emit = defineEmits(['action'])

function handleAction() {
  emit('action')
}
</script>

<style scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  min-height: 300px;
}

.empty-image {
  width: 160px;
  height: 160px;
  margin-bottom: 20px;
  opacity: 0.6;
}

.empty-text {
  color: #969799;
  font-size: 14px;
  margin: 0 0 20px 0;
}
</style>
