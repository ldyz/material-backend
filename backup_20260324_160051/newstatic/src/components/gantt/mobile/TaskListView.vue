<template>
  <div class="task-list-view" ref="containerRef">
    <!-- 搜索栏 -->
    <div v-if="searchKeyword" class="task-list-view__search">
      <el-input
        v-model="localSearchKeyword"
        placeholder="搜索任务..."
        prefix-icon="el-icon-search"
        clearable
        :style="{ width: '100%' }"
      />
    </div>

    <!-- 分组视图 -->
    <template v-if="groupMode && groupedTasks">
      <div
        v-for="(group, key) in groupedTasks"
        :key="key"
        class="task-list-view__group"
      >
        <!-- 分组头部 -->
        <div
          class="task-list-view__group-header"
          @click="toggleGroup(key)"
        >
          <i
            class="collapse-icon"
            :class="collapsedGroups.has(key) ? 'el-icon-caret-right' : 'el-icon-caret-bottom'"
          />
          <span class="group-label">{{ group.label }}</span>
          <span class="group-count">{{ group.tasks.length }}</span>
        </div>

        <!-- 分组内容 -->
        <Transition name="collapse">
          <div v-show="!collapsedGroups.has(key)" class="task-list-view__group-content">
            <TaskRow
              v-for="task in group.tasks"
              :key="task.id"
              :task="task"
              :is-selected="selectedTaskId === task.id"
              @click="handleClick(task)"
            />
          </div>
        </Transition>
      </div>
    </template>

    <!-- 平铺视图 -->
    <template v-else>
      <div class="task-list-view__list">
        <TaskRow
          v-for="task in tasks"
          :key="task.id"
          :task="task"
          :is-selected="selectedTaskId === task.id"
          @click="handleClick(task)"
        />
      </div>
    </template>

    <!-- 空状态 -->
    <div v-if="tasks.length === 0" class="task-list-view__empty">
      <el-empty description="暂无任务" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import TaskRow from './TaskRow.vue'

/**
 * 移动端任务列表视图
 * 显示任务列表，支持分组和搜索
 */
interface Props {
  tasks: any[]
  groupedTasks?: any
  selectedTaskId?: string | number | null
  groupMode?: string
  collapsedGroups?: Set<string>
  searchKeyword?: string
}

const props = withDefaults(defineProps<Props>(), {
  groupedTasks: null,
  selectedTaskId: null,
  groupMode: '',
  collapsedGroups: () => new Set(),
  searchKeyword: ''
})

const emit = defineEmits<{
  rowClick: [task: any]
  toggleGroup: [groupKey: string]
}>()

const containerRef = ref<HTMLElement>()
const localSearchKeyword = ref(props.searchKeyword)

function handleClick(task: any) {
  emit('rowClick', task)
}

function toggleGroup(groupKey: string) {
  emit('toggleGroup', groupKey)
}
</script>

<style scoped>
.task-list-view {
  height: 100%;
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
  background: var(--color-bg-secondary);
}

.task-list-view__search {
  padding: 12px 16px;
  background: var(--color-bg-primary);
  border-bottom: 1px solid var(--color-border-light);
  position: sticky;
  top: 0;
  z-index: 10;
}

.task-list-view__group {
  margin-bottom: 8px;
  background: var(--color-bg-primary);
}

.task-list-view__group-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: var(--color-bg-tertiary);
  border-bottom: 1px solid var(--color-border-light);
  cursor: pointer;
  user-select: none;
}

.collapse-icon {
  font-size: 16px;
  color: var(--color-text-secondary);
  transition: transform var(--transition-fast);
}

.group-label {
  flex: 1;
  font-size: var(--font-size-md);
  font-weight: var(--font-weight-medium);
  color: var(--color-text-primary);
}

.group-count {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  background: var(--color-bg-secondary);
  padding: 2px 8px;
  border-radius: var(--radius-full);
}

.task-list-view__group-content {
  overflow: hidden;
}

.task-list-view__list {
  padding: 8px 0;
}

.task-list-view__empty {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 300px;
}

/* 折叠动画 */
.collapse-enter-active,
.collapse-leave-active {
  transition: all var(--transition-base);
  overflow: hidden;
}

.collapse-enter-from,
.collapse-leave-to {
  max-height: 0;
  opacity: 0;
}

.collapse-enter-to,
.collapse-leave-from {
  max-height: 1000px;
  opacity: 1;
}
</style>
