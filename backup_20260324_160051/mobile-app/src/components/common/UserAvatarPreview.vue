<template>
  <van-image-preview
    v-model:show="showPreview"
    :images="previewImages"
    :start-position="0"
    close-on-click-image
    teleport="body"
  />
</template>

<script setup>
import { ref, computed } from 'vue'
import { getAssetUrl } from '@/utils/request'

const props = defineProps({
  avatar: {
    type: String,
    default: ''
  }
})

const showPreview = ref(false)

// 预览图片列表（只有一张）
const previewImages = computed(() => {
  const url = getAssetUrl(props.avatar)
  return [url]
})

// 暴露方法供父组件调用
function show() {
  showPreview.value = true
}

defineExpose({
  show
})
</script>
