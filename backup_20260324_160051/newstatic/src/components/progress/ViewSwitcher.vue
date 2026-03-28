<template>
  <div class="view-switcher">
    <el-radio-group v-model="currentView" size="large" @change="handleViewChange">
      <el-radio-button value="list">
        <el-icon><List /></el-icon>
        进入列表视图
      </el-radio-button>
      <el-radio-button value="gantt">
        <el-icon><Histogram /></el-icon>
        进入甘特图
      </el-radio-button>
      
    </el-radio-group>

    <div class="view-actions">
      <slot name="actions" />
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { List, Histogram, Share } from '@element-plus/icons-vue'

const props = defineProps({
  modelValue: {
    type: String,
    default: 'list'
  }
})

const emit = defineEmits(['update:modelValue', 'change'])

const currentView = ref(props.modelValue)

watch(() => props.modelValue, (newVal) => {
  currentView.value = newVal
})

const handleViewChange = (value) => {
  emit('update:modelValue', value)
  emit('change', value)
}
</script>

<style scoped>
.view-switcher {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  background: #fff;
  border-radius: 4px;
  margin-bottom: 16px;
}

.view-actions {
  display: flex;
  gap: 12px;
}

:deep(.el-radio-button__inner) {
  display: flex;
  align-items: center;
  gap: 4px;
}

:deep(.el-icon) {
  font-size: 16px;
}
</style>
