/**
 * 富文本编辑器组件
 *
 * 基于 Quill.js 的 Vue 3 封装
 * 支持常用的富文本编辑功能：文本格式、列表、图片上传等
 *
 * @module RichTextEditor
 * @author Material Management System
 * @date 2025-01-27
 */
<template>
  <div class="rich-text-editor" :class="{ 'is-disabled': disabled, 'is-readonly': readOnly }">
    <QuillEditor
      v-model:content="content"
      contentType="html"
      :theme="theme"
      :toolbar="toolbar"
      :placeholder="placeholder"
      :disabled="disabled"
      :read-only="readOnly"
      @update:content="handleUpdate"
      @ready="handleReady"
      @focus="handleFocus"
      @blur="handleBlur"
    />

    <!-- 字数统计 -->
    <div v-if="showWordCount && !readOnly" class="word-count">
      <span :class="{ 'text-danger': isOverLimit }">
        {{ currentLength }} / {{ maxLength }}
      </span>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { QuillEditor } from '@vueup/vue-quill'
import '@vueup/vue-quill/dist/vue-quill.snow.css'
import { useAuthStore } from '@/stores/auth'
import axios from 'axios'

const authStore = useAuthStore()

/**
 * 组件 Props
 */
const props = defineProps({
  /**
   * 编辑器内容（HTML格式）
   */
  modelValue: {
    type: String,
    default: ''
  },
  /**
   * 编辑器主题
   */
  theme: {
    type: String,
    default: 'snow' // snow 或 bubble
  },
  /**
   * 占位符文本
   */
  placeholder: {
    type: String,
    default: '请输入内容...'
  },
  /**
   * 是否禁用
   */
  disabled: {
    type: Boolean,
    default: false
  },
  /**
   * 是否只读
   */
  readOnly: {
    type: Boolean,
    default: false
  },
  /**
   * 最小高度（px）
   */
  minHeight: {
    type: Number,
    default: 200
  },
  /**
   * 最大长度（字符数）
   */
  maxLength: {
    type: Number,
    default: 10000
  },
  /**
   * 是否显示字数统计
   */
  showWordCount: {
    type: Boolean,
    default: true
  },
  /**
   * 自定义工具栏
   */
  customToolbar: {
    type: Array,
    default: null
  }
})

/**
 * 组件 Emits
 */
const emit = defineEmits(['update:modelValue', 'focus', 'blur', 'ready'])

// ========== 状态管理 ==========

/**
 * 编辑器内容
 */
const content = ref(props.modelValue || '')

/**
 * Quill 编辑器实例
 */
const quillInstance = ref(null)

// ========== 计算属性 ==========

/**
 * 当前内容长度
 */
const currentLength = computed(() => {
  if (!content.value) return 0
  // 移除HTML标签计算纯文本长度
  const text = content.value.replace(/<[^>]*>/g, '')
  return text.length
})

/**
 * 是否超过长度限制
 */
const isOverLimit = computed(() => {
  return currentLength.value > props.maxLength
})

/**
 * 工具栏配置
 */
const toolbar = computed(() => {
  if (props.customToolbar) {
    return props.customToolbar
  }

  // 默认工具栏配置
  return [
    ['bold', 'italic', 'underline', 'strike'],        // 加粗、斜体、下划线、删除线
    ['blockquote', 'code-block'],                    // 引用、代码块

    [{ 'header': 1 }, { 'header': 2 }],             // 标题1、标题2
    [{ 'list': 'ordered'}, { 'list': 'bullet' }],    // 有序列表、无序列表
    [{ 'script': 'sub'}, { 'script': 'super' }],     // 上标、下标
    [{ 'indent': '-1'}, { 'indent': '+1' }],         // 缩进
    [{ 'direction': 'rtl' }],                        // 文字方向

    [{ 'size': ['small', false, 'large', 'huge'] }], // 字体大小
    [{ 'header': [1, 2, 3, 4, 5, 6, false] }],      // 标题级别

    [{ 'color': [] }, { 'background': [] }],         // 字体颜色、背景颜色
    [{ 'font': [] }],                                // 字体
    [{ 'align': [] }],                               // 对齐方式

    ['clean'],                                       // 清除格式
    ['link', 'image']                                // 链接、图片
  ]
})

// ========== 方法定义 ==========

/**
 * 内容更新事件处理
 */
const handleUpdate = (value) => {
  content.value = value
  emit('update:modelValue', value)
}

/**
 * 编辑器准备就绪
 */
const handleReady = (quill) => {
  quillInstance.value = quill

  // 设置最小高度
  const editor = quill.root
  if (editor) {
    editor.style.minHeight = `${props.minHeight}px`
  }

  // 配置图片上传
  const toolbar = quill.getModule('toolbar')
  if (toolbar) {
    toolbar.addHandler('image', () => {
      // 创建文件选择器
      const input = document.createElement('input')
      input.setAttribute('type', 'file')
      input.setAttribute('accept', 'image/*')

      input.onchange = async () => {
        const file = input.files[0]
        if (!file) return

        // 显示上传中提示
        const range = quill.getSelection(true)
        quill.disable()

        try {
          // 创建FormData
          const formData = new FormData()
          formData.append('file', file)

          // 上传图片
          const response = await axios.post('/api/upload/image', formData, {
            headers: {
              'Authorization': `Bearer ${authStore.token}`,
              'Content-Type': 'multipart/form-data'
            }
          })

          if (response.data && response.data.code === 200) {
            const imageUrl = response.data.data.url

            // 在编辑器中插入图片
            quill.enable()
            quill.setSelection(range.index, range.index + range.length)
            quill.focus()
            quill.insertEmbed(range.index, 'image', imageUrl)
          } else {
            throw new Error(response.data?.message || '上传失败')
          }
        } catch (error) {
          console.error('图片上传失败:', error)
          // 显示错误提示
          alert('图片上传失败: ' + (error.message || '未知错误'))
          quill.enable()
        } finally {
          // 清理input
          input.value = ''
        }
      }

      input.click()
    })
  }

  emit('ready', quill)
}

/**
 * 获得焦点事件
 */
const handleFocus = () => {
  emit('focus')
}

/**
 * 失去焦点事件
 */
const handleBlur = () => {
  emit('blur')
}

/**
 * 获取纯文本内容
 */
const getText = () => {
  if (!quillInstance.value) return ''
  return quillInstance.value.getText().trim()
}

/**
 * 设置内容
 */
const setContent = (html) => {
  content.value = html
}

/**
 * 清空内容
 */
const clear = () => {
  content.value = ''
}

/**
 * 插入文本
 */
const insertText = (text) => {
  if (!quillInstance.value) return
  const selection = quillInstance.value.getSelection()
  const index = selection ? selection.index : 0
  quillInstance.value.insertText(index, text)
  quillInstance.value.setSelection(index + text.length)
}

/**
 * 插入图片
 */
const insertImage = (url) => {
  if (!quillInstance.value) return
  const selection = quillInstance.value.getSelection()
  const index = selection ? selection.index : 0
  quillInstance.value.insertEmbed(index, 'image', url)
  quillInstance.value.setSelection(index + 1)
}

/**
 * 监听外部值变化
 */
watch(() => props.modelValue, (newValue) => {
  if (newValue !== content.value) {
    content.value = newValue
  }
})

// ========== 暴露给父组件的方法 ==========
defineExpose({
  getText,
  setContent,
  clear,
  insertText,
  insertImage
})
</script>

<style scoped>
.rich-text-editor {
  position: relative;
}

.rich-text-editor.is-disabled {
  opacity: 0.6;
  pointer-events: none;
}

/* 只读模式样式 */
.rich-text-editor.is-readonly :deep(.ql-toolbar) {
  pointer-events: none;
  opacity: 0.5;
}

.rich-text-editor.is-readonly :deep(.ql-toolbar button) {
  cursor: not-allowed;
  opacity: 0.5;
}

.rich-text-editor.is-readonly :deep(.ql-editor) {
  background-color: #fafafa;
}

/* 编辑器容器样式 */
.rich-text-editor :deep(.ql-editor) {
  min-height: 200px;
  font-size: 14px;
  line-height: 1.6;
}

.rich-text-editor :deep(.ql-toolbar) {
  border-radius: 4px 4px 0 0;
  background-color: #fafafa;
}

.rich-text-editor :deep(.ql-container) {
  border-radius: 0 0 4px 4px;
  background-color: #fff;
}

/* 字数统计 */
.word-count {
  text-align: right;
  padding: 8px 0;
  font-size: 12px;
  color: #909399;
}

.word-count .text-danger {
  color: #f56c6c;
  font-weight: 600;
}

/* 编辑器禁用状态 */
.rich-text-editor.is-disabled :deep(.ql-toolbar) {
  background-color: #f5f5f5;
  pointer-events: none;
}

/* 图片样式 */
.rich-text-editor :deep(.ql-editor img) {
  max-width: 100%;
  height: auto;
  display: block;
  margin: 10px auto;
  border-radius: 4px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

/* 只读模式下的图片限制 */
.rich-text-editor.is-readonly :deep(.ql-editor img) {
  max-width: 100%;
  max-height: 600px;
  object-fit: contain;
}

/* 代码块样式 */
.rich-text-editor :deep(.ql-editor pre.ql-syntax) {
  background-color: #f5f5f5;
  border-radius: 4px;
  padding: 12px;
  font-family: 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.5;
}

/* 引用样式 */
.rich-text-editor :deep(.ql-editor blockquote) {
  border-left: 4px solid #ccc;
  margin-left: 0;
  padding-left: 16px;
  color: #666;
}

/* 列表样式 */
.rich-text-editor :deep(.ql-editor ul),
.rich-text-editor :deep(.ql-editor ol) {
  padding-left: 20px;
}

.rich-text-editor :deep(.ql-editor li) {
  margin: 4px 0;
}

/* 链接样式 */
.rich-text-editor :deep(.ql-editor a) {
  color: #409eff;
  text-decoration: none;
}

.rich-text-editor :deep(.ql-editor a:hover) {
  text-decoration: underline;
}
</style>
