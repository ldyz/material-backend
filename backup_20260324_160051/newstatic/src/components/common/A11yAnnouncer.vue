<template>
  <teleport to="body">
    <div
      ref="announcerRef"
      role="status"
      :aria-live="live"
      :aria-atomic="true"
      class="a11y-announcer"
    >
      {{ message }}
    </div>
  </teleport>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

interface Props {
  message: string
  live?: 'polite' | 'assertive'
}

const props = withDefaults(defineProps<Props>(), {
  live: 'polite'
})

const announcerRef = ref<HTMLDivElement>()

/**
 * 监听消息变化，通知屏幕阅读器
 */
watch(() => props.message, (newMessage, oldMessage) => {
  if (newMessage !== oldMessage && announcerRef.value) {
    // 清空后再设置（确保屏幕阅读器能检测到变化）
    announcerRef.value.textContent = ''

    requestAnimationFrame(() => {
      requestAnimationFrame(() => {
        if (announcerRef.value) {
          announcerRef.value.textContent = newMessage
        }
      })
    })
  }
})
</script>

<style scoped>
.a11y-announcer {
  position: absolute;
  left: -10000px;
  width: 1px;
  height: 1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
}
</style>
