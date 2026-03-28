/**
 * 富文本内容显示组件
 *
 * 用于只读显示富文本内容，支持：
 * - HTML 内容渲染
 * - 图片预览
 * - 样式美化
 *
 * @module RichTextContent
 * @author Material Management System
 * @date 2025-01-27
 */
<template>
  <div class="rich-text-content" v-html="content" @click="handleImageClick"></div>
</template>

<script setup>
import { ref } from 'vue'
import { ElImageViewer } from 'element-plus'

/**
 * 组件 Props
 */
const props = defineProps({
  /**
   * HTML 内容
   */
  content: {
    type: String,
    default: ''
  },
  /**
   * 最大高度（px），超过则显示滚动条
   */
  maxHeight: {
    type: Number,
    default: null
  }
})

// ========== 状态管理 ==========

/**
 * 图片预览列表
 */
const previewImages = ref([])

/**
 * 图片预览器显示状态
 */
const showViewer = ref(false)

/**
 * 当前预览图片索引
 */
const initialIndex = ref(0)

// ========== 方法定义 ==========

/**
 * 处理图片点击事件
 */
const handleImageClick = (event) => {
  // 检查点击的是否为图片
  if (event.target.tagName === 'IMG') {
    // 收集所有图片
    const images = event.currentTarget.querySelectorAll('img')
    previewImages.value = Array.from(images).map(img => img.src)

    // 找到当前点击图片的索引
    initialIndex.value = previewImages.value.indexOf(event.target.src)

    // 显示预览器
    showViewer.value = true
  }
}
</script>

<style scoped>
.rich-text-content {
  font-size: 14px;
  line-height: 1.8;
  color: #303133;
  word-wrap: break-word;
  overflow-wrap: break-word;
}

.rich-text-content:empty:before {
  content: attr(data-empty-text) '暂无内容';
  color: #909399;
  font-style: italic;
}

/* 标题样式 */
.rich-text-content h1,
.rich-text-content h2,
.rich-text-content h3,
.rich-text-content h4,
.rich-text-content h5,
.rich-text-content h6 {
  margin: 16px 0 8px 0;
  font-weight: 600;
  line-height: 1.4;
  color: #303133;
}

.rich-text-content h1 {
  font-size: 24px;
  border-bottom: 2px solid #e4e7ed;
  padding-bottom: 8px;
}

.rich-text-content h2 {
  font-size: 20px;
  border-bottom: 1px solid #e4e7ed;
  padding-bottom: 6px;
}

.rich-text-content h3 {
  font-size: 18px;
}

.rich-text-content h4 {
  font-size: 16px;
}

.rich-text-content h5,
.rich-text-content h6 {
  font-size: 14px;
}

/* 段落样式 */
.rich-text-content p {
  margin: 8px 0;
  text-indent: 0;
}

/* 列表样式 */
.rich-text-content ul,
.rich-text-content ol {
  margin: 8px 0;
  padding-left: 24px;
}

.rich-text-content li {
  margin: 4px 0;
}

.rich-text-content ul li {
  list-style-type: disc;
}

.rich-text-content ol li {
  list-style-type: decimal;
}

/* 引用样式 */
.rich-text-content blockquote {
  margin: 16px 0;
  padding: 12px 16px;
  border-left: 4px solid #409eff;
  background-color: #ecf5ff;
  color: #606266;
  font-style: italic;
}

/* 代码块样式 */
.rich-text-content pre {
  margin: 12px 0;
  padding: 12px;
  background-color: #f5f7fa;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  overflow-x: auto;
  font-family: 'Courier New', Courier, monospace;
  font-size: 13px;
  line-height: 1.5;
  color: #606266;
}

.rich-text-content code {
  padding: 2px 6px;
  background-color: #f5f7fa;
  border: 1px solid #e4e7ed;
  border-radius: 3px;
  font-family: 'Courier New', Courier, monospace;
  font-size: 13px;
  color: #e6a23c;
}

/* 表格样式 */
.rich-text-content table {
  width: 100%;
  margin: 12px 0;
  border-collapse: collapse;
  font-size: 13px;
}

.rich-text-content table th,
.rich-text-content table td {
  border: 1px solid #dcdfe6;
  padding: 8px 12px;
  text-align: left;
}

.rich-text-content table th {
  background-color: #f5f7fa;
  font-weight: 600;
  color: #303133;
}

.rich-text-content table tr:nth-child(even) {
  background-color: #fafafa;
}

/* 图片样式 */
.rich-text-content img {
  max-width: 100%;
  height: auto;
  display: block;
  margin: 12px auto;
  border-radius: 4px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: transform 0.3s, box-shadow 0.3s;
}

.rich-text-content img:hover {
  transform: scale(1.02);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

/* 链接样式 */
.rich-text-content a {
  color: #409eff;
  text-decoration: none;
  transition: color 0.3s;
}

.rich-text-content a:hover {
  color: #66b1ff;
  text-decoration: underline;
}

/* 分隔线样式 */
.rich-text-content hr {
  margin: 20px 0;
  border: none;
  border-top: 1px solid #dcdfe6;
}

/* 粗体、斜体、下划线 */
.rich-text-content strong {
  font-weight: 600;
}

.rich-text-content em {
  font-style: italic;
}

.rich-text-content u {
  text-decoration: underline;
}

/* 文本对齐 */
.rich-text-content .ql-align-left {
  text-align: left;
}

.rich-text-content .ql-align-center {
  text-align: center;
}

.rich-text-content .ql-align-right {
  text-align: right;
}

.rich-text-content .ql-align-justify {
  text-align: justify;
}

/* 文本颜色 */
.rich-text-content .ql-color-red {
  color: #f56c6c;
}

.rich-text-content .ql-color-blue {
  color: #409eff;
}

.rich-text-content .ql-color-green {
  color: #67c23a;
}

/* 背景颜色 */
.rich-text-content .ql-bg-gray {
  background-color: #f5f7fa;
  padding: 2px 4px;
}

.rich-text-content .ql-bg-yellow {
  background-color: #fdf6ec;
  padding: 2px 4px;
}

/* 自定义滚动条样式（当内容超过最大高度时） */
.rich-text-content[style*="max-height"] {
  overflow-y: auto;
  padding-right: 8px;
}

.rich-text-content[style*="max-height"]::-webkit-scrollbar {
  width: 6px;
}

.rich-text-content[style*="max-height"]::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.rich-text-content[style*="max-height"]::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

.rich-text-content[style*="max-height"]::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}
</style>
