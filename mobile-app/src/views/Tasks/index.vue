<template>
  <div class="tasks-page">
    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <!-- 筛选栏 -->
      <van-sticky>
        <div class="filter-bar">
          <van-tabs v-model:active="activeTab" @change="onTabChange" sticky>
            <van-tab title="全部" name="all" />
            <van-tab title="待审批" name="pending" />
            <van-tab title="已处理" name="processed" />
          </van-tabs>
        </div>
      </van-sticky>

      <!-- 任务列表 -->
      <van-list
        v-model:loading="loading"
        :finished="finished"
        finished-text="没有更多了"
        @load="onLoad"
      >
        <div
          v-for="task in tasks"
          :key="task.id"
          class="task-card"
          @click="goToDetail(task)"
        >
          <div class="task-header">
            <div class="task-type">
              <van-tag :type="getTaskTypeColor(task.entity_type)">
                {{ getTaskTypeName(task.entity_type) }}
              </van-tag>
              <span class="task-title">{{ task.title }}</span>
            </div>
            <van-icon name="arrow" color="#969799" />
          </div>

          <div class="task-content">
            <div class="content-row">
              <span class="label">单号:</span>
              <span class="value">{{ task.entity_no || '-' }}</span>
            </div>
            <div class="content-row">
              <span class="label">项目:</span>
              <span class="value">{{ task.project_name || '-' }}</span>
            </div>
            <div class="content-row">
              <span class="label">节点:</span>
              <span class="value">{{ task.node_name }}</span>
            </div>
            <div class="content-row">
              <span class="label">申请人:</span>
              <span class="value">{{ task.creator_name }}</span>
            </div>
            <div class="content-row">
              <span class="label">提交时间:</span>
              <span class="value">{{ formatDateTime(task.created_at) }}</span>
            </div>
          </div>

          <!-- 快捷操作 -->
          <div class="task-actions" v-if="task.status === 'pending'">
            <van-button
              type="success"
              size="small"
              @click.stop="handleApprove(task)"
            >
              通过
            </van-button>
            <van-button
              type="danger"
              size="small"
              @click.stop="handleReject(task)"
            >
              拒绝
            </van-button>
          </div>
        </div>
      </van-list>
    </van-pull-refresh>

    <!-- 空状态 -->
    <van-empty
      v-if="!loading && tasks.length === 0 && !refreshing"
      description="暂无待办任务"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { showToast, showConfirmDialog } from 'vant'
import { getPendingTasks, approveTask, rejectTask } from '@/api/workflow'

const router = useRouter()

const tasks = ref([])
const loading = ref(false)
const finished = ref(false)
const refreshing = ref(false)
const currentPage = ref(1)
const pageSize = 20
const activeTab = ref('all')

async function onLoad() {
  if (refreshing.value) {
    tasks.value = []
    currentPage.value = 1
    finished.value = false
  }

  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      per_page: pageSize,
    }

    // 根据当前tab筛选状态
    if (activeTab.value === 'pending') {
      params.status = 'pending'
    } else if (activeTab.value === 'processed') {
      params.status = 'completed'
    }

    const response = await getPendingTasks(params)
    if (response.success) {
      const { data, total } = response.data

      if (refreshing.value) {
        tasks.value = data || []
      } else {
        tasks.value.push(...(data || []))
      }

      if (tasks.value.length >= total) {
        finished.value = true
      } else {
        currentPage.value++
      }
    }
  } catch (error) {
    console.error('加载任务列表失败:', error)
    showToast('加载失败')
  } finally {
    loading.value = false
    refreshing.value = false
  }
}

function onRefresh() {
  finished.value = false
  onLoad()
}

function onTabChange() {
  tasks.value = []
  currentPage.value = 1
  finished.value = false
  onLoad()
}

function goToDetail(task) {
  // 根据任务类型跳转到不同详情页
  if (task.entity_type === 'material_plan') {
    router.push(`/plans/${task.entity_id}`)
  } else if (task.entity_type === 'inbound_order') {
    router.push(`/inbound/${task.entity_id}`)
  } else if (task.entity_type === 'requisition') {
    router.push(`/outbound/${task.entity_id}`)
  }
}

async function handleApprove(task) {
  const remark = prompt('请输入审批意见（可选）')
  if (remark === null) return // 用户取消

  loading.value = true
  try {
    const response = await approveTask(task.id, { comment: remark || '' })
    if (response.success) {
      showToast('已通过')
      // 刷新列表
      refreshing.value = true
      onRefresh()
    } else {
      showToast(response.message || '操作失败')
    }
  } catch (error) {
    console.error('审批失败:', error)
    showToast('操作失败')
  } finally {
    loading.value = false
  }
}

async function handleReject(task) {
  const { value } = await showConfirmDialog({
    title: '确认拒绝',
    message: '拒绝后流程将被终止，是否继续？',
  })

  if (!value) return

  const remark = prompt('请输入拒绝原因（可选）')
  if (remark === null) return

  loading.value = true
  try {
    const response = await rejectTask(task.id, { reason: remark || '' })
    if (response.success) {
      showToast('已拒绝')
      refreshing.value = true
      onRefresh()
    } else {
      showToast(response.message || '操作失败')
    }
  } catch (error) {
    console.error('拒绝失败:', error)
    showToast('操作失败')
  } finally {
    loading.value = false
  }
}

function getTaskTypeColor(entityType) {
  const colors = {
    material_plan: 'primary',
    inbound_order: 'success',
    requisition: 'warning',
  }
  return colors[entityType] || 'default'
}

function getTaskTypeName(entityType) {
  const names = {
    material_plan: '物资计划',
    inbound_order: '入库单',
    requisition: '出库单',
  }
  return names[entityType] || entityType
}

function formatDateTime(dateStr) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return `${date.getMonth() + 1}月${date.getDate()}日 ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
}

onMounted(() => {
  onLoad()
})
</script>

<style scoped>
.tasks-page {
  min-height: 100vh;
  background-color: #f7f8fa;
  padding-bottom: 60px;
}

.filter-bar {
  background: white;
}

.task-card {
  background: white;
  margin: 12px 16px;
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.task-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.task-type {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
}

.task-title {
  font-size: 15px;
  font-weight: 600;
  color: #323233;
}

.task-content {
  margin-bottom: 12px;
}

.content-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
  font-size: 13px;
}

.content-row:last-child {
  margin-bottom: 0;
}

.label {
  color: #969799;
  margin-right: 8px;
}

.value {
  color: #323233;
  font-weight: 500;
}

.task-actions {
  display: flex;
  gap: 8px;
  padding-top: 12px;
  border-top: 1px solid #ebedf0;
}
</style>
