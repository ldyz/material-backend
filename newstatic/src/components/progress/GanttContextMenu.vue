<template>
  <el-dropdown
    ref="dropdownRef"
    :virtual-ref="virtualTrigger"
    virtual-triggering
    trigger="contextmenu"
    @command="handleCommand"
    @visible-change="handleVisibleChange"
  >
    <template #dropdown>
      <el-dropdown-menu>
        <el-dropdown-item command="add-subtask" :icon="Plus" :disabled="!canAddSubtask">
          添加子任务
        </el-dropdown-item>
        <el-dropdown-item command="convert-milestone" :icon="Star" :disabled="task?.is_milestone">
          转换为里程碑
        </el-dropdown-item>
        <el-dropdown-item command="add-dependency" :icon="Connection" divided>
          添加前置任务
        </el-dropdown-item>
        <el-dropdown-item command="view-dependencies" :icon="View">
          查看依赖关系
        </el-dropdown-item>
        <el-dropdown-item command="allocate-resources" :icon="UserFilled" divided>
          分配资源
        </el-dropdown-item>
        <el-dropdown-item command="edit" :icon="Edit">编辑任务</el-dropdown-item>
        <el-dropdown-item command="duplicate" :icon="CopyDocument">复制任务</el-dropdown-item>
        <el-dropdown-item command="delete" :icon="Delete" divided>删除任务</el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
</template>

<script setup>
import { ref, computed, nextTick } from 'vue'
import {
  Plus,
  Star,
  Connection,
  View,
  UserFilled,
  Edit,
  CopyDocument,
  Delete
} from '@element-plus/icons-vue'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  task: {
    type: Object,
    default: null
  },
  position: {
    type: Object,
    default: () => ({ x: 0, y: 0 })
  }
})

const emit = defineEmits([
  'add-subtask',
  'convert-milestone',
  'add-dependency',
  'view-dependencies',
  'allocate-resources',
  'edit',
  'duplicate',
  'delete'
])

const dropdownRef = ref(null)

const virtualTrigger = computed(() => {
  if (props.visible && props.position) {
    return {
      getBoundingClientRect: () => new DOMRect(props.position.x, props.position.y, 0, 0)
    }
  }
  return undefined
})

const canAddSubtask = computed(() => {
  // 里程碑不能有子任务
  if (props.task?.is_milestone) return false
  return true
})

const handleCommand = (command) => {
  emit(command, props.task)
}

const handleVisibleChange = (visible) => {
  if (!visible) {
    emit('update:visible', false)
  }
}

const open = () => {
  nextTick(() => {
    dropdownRef.value?.handleOpen()
  })
}

defineExpose({
  open
})
</script>

<style scoped>
/* ContextMenu styles are handled by Element Plus */
</style>
