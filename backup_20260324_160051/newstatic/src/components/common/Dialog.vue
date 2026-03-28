<template>
  <el-dialog
    v-model="visible"
    :width="width"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <template #header>
      <div class="dialog-header-custom">
        <span class="el-dialog__title">{{ title }}</span>
        <div class="dialog-header-extra" v-if="$slots.extra">
          <slot name="extra"></slot>
        </div>
      </div>
    </template>

    <slot></slot>
    <template #footer v-if="showFooter">
      <el-button @click="handleClose" v-if="showCancel">{{ cancelText }}</el-button>
      <el-button type="primary" @click="handleConfirm" :loading="loading" v-if="showConfirm">
        {{ confirmText }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: '对话框'
  },
  width: {
    type: String,
    default: '600px'
  },
  showFooter: {
    type: Boolean,
    default: true
  },
  loading: {
    type: Boolean,
    default: false
  },
  cancelText: {
    type: String,
    default: '取消'
  },
  confirmText: {
    type: String,
    default: '确认'
  },
  showCancel: {
    type: Boolean,
    default: true
  },
  showConfirm: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['update:modelValue', 'confirm'])

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const handleClose = () => {
  visible.value = false
}

const handleConfirm = () => {
  emit('confirm')
}
</script>

<style scoped>
.dialog-header-custom {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.dialog-header-custom .el-dialog__title {
  flex: 1;
}

.dialog-header-extra {
  flex-shrink: 0;
  margin-right: 16px;
}
</style>
