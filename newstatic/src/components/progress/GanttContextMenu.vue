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
        <!-- 空行菜单 -->
        <template v-if="!task">
          <el-dropdown-item command="create-task" :icon="Plus">
            新建任务
          </el-dropdown-item>
        </template>

        <!-- 已有任务的菜单 -->
        <template v-else>
          <!-- 层级操作 -->
          <el-dropdown-item command="move-up" :icon="ArrowUp" :disabled="!canMoveUp">
            上移任务
          </el-dropdown-item>
          <el-dropdown-item command="move-down" :icon="ArrowDown" :disabled="!canMoveDown">
            下移任务
          </el-dropdown-item>
          <el-dropdown-item command="convert-to-independent" :icon="Remove" :disabled="!canConvertToIndependent" divided>
            转为独立任务
          </el-dropdown-item>
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
        </template>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
</template>

<script setup>
import { ref, computed, nextTick, watch } from 'vue'
import {
  Plus,
  Star,
  Connection,
  View,
  UserFilled,
  Edit,
  CopyDocument,
  Delete,
  ArrowUp,
  ArrowDown,
  Remove
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
  },
  // 新增：所有任务列表，用于判断是否可以上移/下移
  allTasks: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits([
  'update:visible',
  'create-task',
  'add-subtask',
  'convert-milestone',
  'add-dependency',
  'view-dependencies',
  'allocate-resources',
  'edit',
  'duplicate',
  'delete',
  'move-up',
  'move-down',
  'convert-to-independent'
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

// 判断是否可以上移
const canMoveUp = computed(() => {
  if (!props.task || !props.allTasks.length) return false
  const taskIndex = props.allTasks.findIndex(t => t.id === props.task.id)
  return taskIndex > 0
})

// 判断是否可以下移
const canMoveDown = computed(() => {
  if (!props.task || !props.allTasks.length) return false
  const taskIndex = props.allTasks.findIndex(t => t.id === props.task.id)
  return taskIndex < props.allTasks.length - 1
})

// 判断是否可以转为独立任务（有父任务的才可以）
const canConvertToIndependent = computed(() => {
  return props.task?.parent_id != null
})

const handleCommand = (command) => {
  emit(command, props.task)
}

const handleVisibleChange = (visible) => {
  if (!visible) {
    emit('update:visible', false)
  }
}

// Watch visible prop to open/close dropdown
watch(() => props.visible, (visible) => {
  if (visible) {
    nextTick(() => {
      dropdownRef.value?.handleOpen()
    })
  } else {
    dropdownRef.value?.handleClose()
  }
})

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
