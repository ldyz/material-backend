/**
 * 项目选择器组件
 *
 * 功能说明：
 * - 支持树形层级显示项目
 * - 支持搜索过滤
 * - 支持清除选择
 * - 支持选择主项目或子项目
 * - 统一的样式和交互
 *
 * @module ProjectSelector
 * @author Material Management System
 * @date 2026-02-01
 */
<template>
  <el-tree-select
    :model-value="modelValue"
    :data="projectTreeData"
    :props="treeProps"
    :placeholder="placeholder"
    :clearable="clearable"
    :filterable="filterable"
    :check-strictly="checkStrictly"
    :render-after-expand="renderAfterExpand"
    :multiple="multiple"
    :disabled="disabled || loading"
    :loading="loading"
    :size="size"
    :style="{ width: width }"
    @change="handleChange"
    @focus="$emit('focus', $event)"
    @visible-change="$emit('visible-change', $event)"
  />
</template>

<script setup>
import { computed } from 'vue'

/**
 * 组件 Props
 */
const props = defineProps({
  /**
   * 模型值，v-model 绑定
   */
  modelValue: {
    type: [String, Number, Array],
    default: null
  },
  /**
   * 项目列表数据（扁平或树形）
   */
  projects: {
    type: Array,
    default: () => []
  },
  /**
   * 占位文本
   */
  placeholder: {
    type: String,
    default: '请选择项目'
  },
  /**
   * 是否可清除
   */
  clearable: {
    type: Boolean,
    default: true
  },
  /**
   * 是否可搜索
   */
  filterable: {
    type: Boolean,
    default: true
  },
  /**
   * 是否任意选择一级
   */
  checkStrictly: {
    type: Boolean,
    default: true
  },
  /**
   * 是否展开子节点
   */
  renderAfterExpand: {
    type: Boolean,
    default: false
  },
  /**
   * 是否多选
   */
  multiple: {
    type: Boolean,
    default: false
  },
  /**
   * 是否禁用
   */
  disabled: {
    type: Boolean,
    default: false
  },
  /**
   * 尺寸
   */
  size: {
    type: String,
    default: 'default',
    validator: (value) => ['large', 'default', 'small'].includes(value)
  },
  /**
   * 宽度
   */
  width: {
    type: String,
    default: '100%'
  },
  /**
   * 是否加载中
   */
  loading: {
    type: Boolean,
    default: false
  }
})

/**
 * 组件 Emits
 */
const emit = defineEmits(['update:modelValue', 'change', 'focus', 'visible-change'])

/**
 * 树形配置
 */
const treeProps = {
  label: 'name',
  value: 'id',
  children: 'children'
}

/**
 * 项目树形数据
 * 将扁平列表转换为树形结构
 */
const projectTreeData = computed(() => {
  console.log('[ProjectSelector] projects prop 变化:', props.projects)

  if (!props.projects || props.projects.length === 0) {
    console.log('[ProjectSelector] 项目列表为空')
    return []
  }

  // 如果已经是树形结构（有children属性），直接返回
  if (props.projects[0] && props.projects[0].children) {
    console.log('[ProjectSelector] 已经是树形结构，直接返回')
    return props.projects
  }

  // 构建树形结构
  const tree = buildProjectTree(props.projects)
  console.log('[ProjectSelector] 构建后的树形结构:', tree)
  return tree
})

/**
 * 构建项目树形结构
 */
const buildProjectTree = (projects) => {
  if (!projects || projects.length === 0) return []

  // 创建映射表
  const projectMap = new Map()
  projects.forEach(project => {
    projectMap.set(project.id, { ...project, children: [] })
  })

  // 构建树形结构
  const roots = []
  projects.forEach(project => {
    const node = projectMap.get(project.id)
    if (!project.parent_id) {
      roots.push(node)
    } else {
      const parent = projectMap.get(project.parent_id)
      if (parent) {
        parent.children.push(node)
      } else {
        // 找不到父节点，作为根节点
        roots.push(node)
      }
    }
  })

  return roots
}

/**
 * 变化事件处理
 */
const handleChange = (value) => {
  emit('update:modelValue', value)
  emit('change', value)
}
</script>

<style scoped>
/* 组件样式继承自 el-tree-select */
</style>
