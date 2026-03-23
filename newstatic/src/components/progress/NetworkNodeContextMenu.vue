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
        <!-- 节点信息 -->
        <el-dropdown-item disabled class="node-info">
          <div class="node-info-content">
            <span class="node-label">节点 {{ node?.id }}</span>
            <span v-if="node" class="node-date">{{ nodeDateDisplay }}</span>
          </div>
        </el-dropdown-item>

        <el-dropdown-item divided disabled v-if="node">
          <div class="task-count-info">
            <span>关联任务: {{ relatedTaskCount }} 个</span>
          </div>
        </el-dropdown-item>

        <!-- R11：合并节点 -->
        <el-dropdown-item
          command="merge-nodes"
          :icon="Connection"
          :disabled="!canMerge"
          divided
        >
          合并选中节点 (R11)
        </el-dropdown-item>

        <!-- 拆分节点 -->
        <el-dropdown-item
          command="split-node"
          :icon="Remove"
          :disabled="!canSplit"
        >
          拆分节点
        </el-dropdown-item>

        <!-- 编辑关联任务 -->
        <el-dropdown-item
          command="edit-tasks"
          :icon="Edit"
          divided
        >
          编辑关联任务
        </el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
</template>

<script setup>
import { ref, computed, nextTick, watch } from 'vue'
import {
  Connection,
  Remove,
  Edit
} from '@element-plus/icons-vue'
import { formatDate } from '@/utils/dateFormat'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  node: {
    type: Object,
    default: null
  },
  selectedNodes: {
    type: Array,
    default: () => []
  },
  canMerge: {
    type: Boolean,
    default: false
  },
  position: {
    type: Object,
    default: () => ({ x: 0, y: 0 })
  }
})

const emit = defineEmits([
  'update:visible',
  'merge-nodes',
  'split-node',
  'edit-tasks'
])

const dropdownRef = ref(null)

// 创建虚拟触发元素
const virtualTrigger = ref({
  getBoundingClientRect: () => {
    return new DOMRect(
      props.position.x,
      props.position.y,
      0,
      0
    )
  }
})

// 计算关联任务数量
const relatedTaskCount = computed(() => {
  if (!props.node) return 0
  const startCount = props.node.tasks?.start?.length || 0
  const endCount = props.node.tasks?.end?.length || 0
  return startCount + endCount
})

// 是否可以拆分节点（需要关联多个任务）
const canSplit = computed(() => {
  return relatedTaskCount.value > 1
})

// 格式化节点日期显示 - 使用节点的 date 属性
const nodeDateDisplay = computed(() => {
  if (!props.node?.date) return '-'
  // node.date 是格式化的日期字符串，直接返回
  return props.node.date
})

// 监听visible变化，控制菜单显示
watch(() => props.visible, (visible) => {
  if (visible) {
    nextTick(() => {
      dropdownRef.value?.handleOpen()
    })
  } else {
    dropdownRef.value?.handleClose()
  }
})

const handleVisibleChange = (visible) => {
  if (!visible) {
    emit('update:visible', false)
  }
}

const handleCommand = (command) => {
  switch (command) {
    case 'merge-nodes':
      emit('merge-nodes', {
        node: props.node,
        selectedNodes: props.selectedNodes
      })
      break
    case 'split-node':
      emit('split-node', {
        node: props.node,
        nodeId: props.node?.id
      })
      break
    case 'edit-tasks':
      emit('edit-tasks', {
        node: props.node
      })
      break
  }
  emit('update:visible', false)
}

// 更新虚拟触发器位置
watch(() => props.position, () => {
  virtualTrigger.value = {
    getBoundingClientRect: () => {
      return new DOMRect(
        props.position.x,
        props.position.y,
        0,
        0
      )
    }
  }
})
</script>

<style scoped>
.node-info {
  cursor: default !important;
}

.node-info-content {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.node-label {
  font-weight: bold;
  color: #303133;
}

.node-date {
  font-size: 12px;
  color: #909399;
}

.task-count-info {
  display: flex;
  justify-content: space-between;
  gap: 20px;
  font-size: 12px;
  color: #606266;
}
</style>
