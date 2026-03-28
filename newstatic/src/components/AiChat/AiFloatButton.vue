<template>
  <div class="ai-float-button" @click="$emit('click')">
    <el-badge :value="unreadCount" :hidden="!unreadCount" :max="99">
      <el-icon :size="28">
        <ChatDotRound />
      </el-icon>
    </el-badge>
  </div>
</template>

<script setup>
import { ChatDotRound } from '@element-plus/icons-vue'

defineProps({
  unreadCount: {
    type: Number,
    default: 0
  }
})

defineEmits(['click'])
</script>

<style scoped>
.ai-float-button {
  position: fixed;
  right: 24px;
  bottom: 24px;
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: linear-gradient(135deg, #409eff, #66b1ff);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.4);
  transition: all 0.3s;
  z-index: 1000;
}

.ai-float-button:hover {
  transform: scale(1.1);
  box-shadow: 0 6px 16px rgba(64, 158, 255, 0.5);
}

.ai-float-button:active {
  transform: scale(0.95);
}
</style>
