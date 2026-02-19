<template>
  <div class="history-panel">
    <!-- Header -->
    <div class="history-header">
      <h3>Change History</h3>
      <div class="header-actions">
        <el-dropdown trigger="click" @command="handleFilterCommand">
          <el-button :icon="Filter" size="small">
            Filter
            <el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="all">All Changes</el-dropdown-item>
              <el-dropdown-item command="task">Tasks Only</el-dropdown-item>
              <el-dropdown-item command="dependency">Dependencies Only</el-dropdown-item>
              <el-dropdown-item command="resource">Resources Only</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>

        <el-dropdown trigger="click" @command="handleSortCommand">
          <el-button :icon="Sort" size="small">
            Sort
            <el-icon class="el-icon--right"><ArrowDown /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="newest">Newest First</el-dropdown-item>
              <el-dropdown-item command="oldest">Oldest First</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>

        <el-button
          :icon="Download"
          size="small"
          @click="exportHistory"
        >
          Export
        </el-button>
      </div>
    </div>

    <!-- Filters -->
    <div v-if="activeFilter !== 'all'" class="history-filters">
      <el-tag closable @close="clearFilter">
        Filter: {{ filterLabel }}
      </el-tag>
    </div>

    <!-- History List -->
    <el-scrollbar ref="scrollbarRef" class="history-list">
      <div v-if="loading" class="history-loading">
        <el-skeleton :rows="5" animated />
      </div>

      <div v-else-if="filteredHistory.length === 0" class="history-empty">
        <el-empty description="No change history" :image-size="80">
          <p class="empty-hint">Changes to tasks will appear here</p>
        </el-empty>
      </div>

      <el-timeline v-else class="history-timeline">
        <el-timeline-item
          v-for="item in filteredHistory"
          :key="item.id"
          :timestamp="formatDate(item.created_at)"
          placement="top"
          :type="getActionType(item.action_type)"
          :icon="getActionIcon(item.action_type)"
          size="large"
        >
          <HistoryItem
            :change="item"
            :user="getUser(item.user_id)"
            :show-diff="expandedItems.has(item.id)"
            @toggle-diff="toggleDiff(item.id)"
            @rollback="handleRollback"
          />
        </el-timeline-item>
      </el-timeline>
    </el-scrollbar>

    <!-- Load More -->
    <div v-if="hasMore" class="history-footer">
      <el-button
        :loading="loadingMore"
        @click="loadMore"
      >
        Load More
      </el-button>
    </div>
  </div>
</template>

<script setup>
/**
 * HistoryPanel.vue
 *
 * Change history panel with timeline visualization,
 * diff display, and rollback functionality.
 */

import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Filter,
  Sort,
  Download,
  ArrowDown,
  Plus,
  Edit,
  Delete,
  RefreshRight
} from '@element-plus/icons-vue'
import { format } from 'date-fns'
import { progressApi } from '@/api'
import HistoryItem from './HistoryItem.vue'

// ==================== Props ====================

const props = defineProps({
  taskId: {
    type: Number,
    default: null
  },
  projectId: {
    type: Number,
    required: true
  },
  projectUsers: {
    type: Array,
    default: () => []
  }
})

// ==================== State ====================

const scrollbarRef = ref(null)
const loading = ref(false)
const loadingMore = ref(false)
const history = ref([])
const expandedItems = ref(new Set())
const activeFilter = ref('all')
const sortOrder = ref('newest')
const page = ref(1)
const hasMore = ref(false)

// ==================== Computed ====================

/**
 * Get filter label
 */
const filterLabel = computed(() => {
  const labels = {
    task: 'Tasks',
    dependency: 'Dependencies',
    resource: 'Resources'
  }
  return labels[activeFilter.value] || activeFilter.value
})

/**
 * Filter and sort history
 */
const filteredHistory = computed(() => {
  let items = [...history.value]

  // Apply filter
  if (activeFilter.value !== 'all') {
    items = items.filter(item => item.entity_type === activeFilter.value)
  }

  // Apply sort
  if (sortOrder.value === 'newest') {
    items.sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
  } else {
    items.sort((a, b) => new Date(a.created_at) - new Date(b.created_at))
  }

  return items
})

// ==================== Methods ====================

/**
 * Get user by ID
 */
function getUser(userId) {
  return props.projectUsers.find(u => u.id === userId) || {
    name: 'Unknown User',
    avatar: null
  }
}

/**
 * Format date
 */
function formatDate(dateString) {
  try {
    return format(new Date(dateString), 'MMM d, yyyy HH:mm')
  } catch {
    return 'Unknown date'
  }
}

/**
 * Get action type for timeline
 */
function getActionType(actionType) {
  const types = {
    create: 'success',
    update: 'primary',
    delete: 'danger'
  }
  return types[actionType] || 'info'
}

/**
 * Get action icon
 */
function getActionIcon(actionType) {
  const icons = {
    create: Plus,
    update: Edit,
    delete: Delete
  }
  return icons[actionType] || RefreshRight
}

/**
 * Load change history
 */
async function loadHistory(reset = true) {
  if (reset) {
    page.value = 1
    history.value = []
  }

  loading.value = reset
  loadingMore.value = !reset

  try {
    const response = props.taskId
      ? await progressApi.getTaskHistory(props.taskId, { page: page.value })
      : await progressApi.getProjectHistory(props.projectId, { page: page.value })

    if (reset) {
      history.value = response.data || []
    } else {
      history.value.push(...(response.data || []))
    }

    hasMore.value = response.data?.length === 20 // Assuming page size of 20
  } catch (error) {
    console.error('[HistoryPanel] Error loading history:', error)
    ElMessage.error('Failed to load change history')
  } finally {
    loading.value = false
    loadingMore.value = false
  }
}

/**
 * Load more history
 */
function loadMore() {
  page.value++
  loadHistory(false)
}

/**
 * Handle filter command
 */
function handleFilterCommand(command) {
  activeFilter.value = command
}

/**
 * Clear filter
 */
function clearFilter() {
  activeFilter.value = 'all'
}

/**
 * Handle sort command
 */
function handleSortCommand(command) {
  sortOrder.value = command
}

/**
 * Toggle diff visibility
 */
function toggleDiff(itemId) {
  if (expandedItems.value.has(itemId)) {
    expandedItems.value.delete(itemId)
  } else {
    expandedItems.value.add(itemId)
  }
}

/**
 * Handle rollback
 */
async function handleRollback(changeId) {
  try {
    await ElMessageBox.confirm(
      'This will revert the changes made in this entry. Continue?',
      'Confirm Rollback',
      {
        type: 'warning',
        confirmButtonText: 'Rollback',
        cancelButtonText: 'Cancel'
      }
    )

    await progressApi.rollbackChange(changeId)

    ElMessage.success('Change rolled back successfully')
    loadHistory()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('[HistoryPanel] Error rolling back:', error)
      ElMessage.error('Failed to rollback change')
    }
  }
}

/**
 * Export history
 */
function exportHistory() {
  const items = filteredHistory.value

  // Convert to CSV
  const headers = ['Date', 'User', 'Action', 'Entity', 'Changes']
  const rows = items.map(item => [
    formatDate(item.created_at),
    getUser(item.user_id).name,
    item.action_type,
    `${item.entity_type} #${item.entity_id}`,
    JSON.stringify(item.changes)
  ])

  const csv = [
    headers.join(','),
    ...rows.map(row => row.map(cell => `"${cell}"`).join(','))
  ].join('\n')

  // Download
  const blob = new Blob([csv], { type: 'text/csv' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = `history-${props.projectId}-${Date.now()}.csv`
  link.click()
  URL.revokeObjectURL(url)

  ElMessage.success('History exported successfully')
}

// ==================== Lifecycle ====================

onMounted(() => {
  loadHistory()
})
</script>

<style scoped lang="scss">
.history-panel {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #fff;
  border-radius: 4px;
}

.history-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  border-bottom: 1px solid #e4e7ed;

  h3 {
    margin: 0;
    font-size: 16px;
    font-weight: 600;
    color: #303133;
  }
}

.header-actions {
  display: flex;
  gap: 8px;
}

.history-filters {
  padding: 12px 16px;
  border-bottom: 1px solid #e4e7ed;
}

.history-list {
  flex: 1;
  padding: 16px;
}

.history-loading {
  padding: 16px;
}

.history-empty {
  padding: 40px 16px;
  text-align: center;

  .empty-hint {
    margin-top: 8px;
    font-size: 12px;
    color: #909399;
  }
}

.history-timeline {
  padding-left: 8px;
}

.history-footer {
  padding: 16px;
  border-top: 1px solid #e4e7ed;
  text-align: center;
}
</style>
