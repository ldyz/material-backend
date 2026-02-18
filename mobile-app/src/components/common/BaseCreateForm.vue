<template>
  <div class="base-create-form">
    <van-form @submit="handleSubmit">
      <!-- 基础信息区插槽 -->
      <div class="form-section">
        <slot name="basic-fields" :form-data="formData">
          <!-- 默认内容：如果没有提供插槽，使用此默认内容 -->
          <van-cell-group inset title="基本信息">
            <slot name="default-fields"></slot>
          </van-cell-group>
        </slot>
      </div>

      <!-- 物料明细区插槽（可选） -->
      <div v-if="showItems" class="form-section">
        <slot name="items" :form-data="formData" :items="items">
          <!-- 默认内容：如果没有提供插槽，使用此默认内容 -->
          <van-cell-group inset title="物料明细">
            <slot name="default-items"></slot>
          </van-cell-group>
        </slot>
      </div>

      <!-- 底部占位，防止内容被按钮遮挡 -->
      <div class="form-footer-spacer"></div>

      <!-- 固定底部提交按钮 -->
      <van-submit-bar
        :button-text="submitText"
        :loading="submitting"
        @submit="handleSubmit"
      />
    </van-form>
  </div>
</template>

<script setup>
const props = defineProps({
  /**
   * 提交按钮文本
   */
  submitText: {
    type: String,
    default: '提交'
  },
  /**
   * 提交中状态
   */
  submitting: {
    type: Boolean,
    default: false
  },
  /**
   * 是否显示物料明细区
   */
  showItems: {
    type: Boolean,
    default: false
  },
  /**
   * 表单数据（用于插槽访问）
   */
  formData: {
    type: Object,
    default: () => ({})
  },
  /**
   * 物料明细数据（用于插槽访问）
   */
  items: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['submit', 'update:formData', 'update:items'])

function handleSubmit() {
  emit('submit')
}
</script>

<style scoped>
.base-create-form {
  min-height: 100vh;
  background-color: #f7f8fa;
  padding-bottom: constant(safe-area-inset-bottom);
  padding-bottom: env(safe-area-inset-bottom);
}

.form-section {
  margin-bottom: 12px;
}

.form-footer-spacer {
  height: 60px;
}
</style>
