<template>
  <el-dialog
    v-model="dialogVisible"
    title="预约单详情"
    width="900px"
    @close="handleClose"
  >
    <el-descriptions :column="2" border v-if="appointment">
      <el-descriptions-item label="预约单号">
        {{ appointment.appointment_no }}
      </el-descriptions-item>
      <el-descriptions-item label="状态">
        <el-tag :type="getStatusType(appointment.status)">
          {{ getStatusLabel(appointment.status) }}
        </el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="申请人">
        {{ appointment.applicant_name }}
      </el-descriptions-item>
      <el-descriptions-item label="联系电话">
        {{ appointment.contact_phone || '-' }}
      </el-descriptions-item>
      <el-descriptions-item label="作业时间" :span="2">
        {{ formatDateTime(appointment.work_date, appointment.time_slot) }}
      </el-descriptions-item>
      <el-descriptions-item label="作业地点" :span="2">
        {{ appointment.work_location }}
      </el-descriptions-item>
      <el-descriptions-item label="作业内容" :span="2">
        {{ appointment.work_content }}
      </el-descriptions-item>
      <el-descriptions-item label="作业类型" v-if="appointment.work_type">
        {{ formatWorkType(appointment.work_type) }}
      </el-descriptions-item>
      <el-descriptions-item label="作业人员">
        <span v-if="appointment.assigned_worker_names">{{ appointment.assigned_worker_names }}</span>
        <span v-else-if="appointment.assigned_worker_name">{{ appointment.assigned_worker_name }}</span>
        <span v-else>未分配</span>
      </el-descriptions-item>
      <el-descriptions-item label="是否加急">
        <el-tag v-if="appointment.is_urgent" type="danger" size="small">是</el-tag>
        <el-tag v-else type="info" size="small">否</el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="优先级" v-if="appointment.is_urgent">
        {{ appointment.priority }}
      </el-descriptions-item>
      <el-descriptions-item label="加急原因" v-if="appointment.urgent_reason" :span="2">
        {{ appointment.urgent_reason }}
      </el-descriptions-item>
      <el-descriptions-item label="创建时间">
        {{ formatFullDateTime(appointment.created_at) }}
      </el-descriptions-item>
      <el-descriptions-item label="提交时间" v-if="appointment.submitted_at">
        {{ formatFullDateTime(appointment.submitted_at) }}
      </el-descriptions-item>
      <el-descriptions-item label="审批时间" v-if="appointment.approved_at">
        {{ formatFullDateTime(appointment.approved_at) }}
      </el-descriptions-item>
      <el-descriptions-item label="完成时间" v-if="appointment.completed_at">
        {{ formatFullDateTime(appointment.completed_at) }}
      </el-descriptions-item>
    </el-descriptions>

    <!-- 审批历史 -->
    <div class="approval-history" v-if="approvalLogs.length > 0">
      <h4>审批历史</h4>
      <el-timeline>
        <el-timeline-item
          v-for="(log, index) in approvalLogs"
          :key="index"
          :timestamp="formatFullDateTime(log.created_at)"
        >
          <div class="timeline-content">
            <div class="timeline-title">{{ log.node_name }}</div>
            <div class="timeline-action">
              <el-tag :type="getActionType(log.action)" size="small">
                {{ getActionText(log.action) }}
              </el-tag>
              <span class="approver">{{ log.approver_name }}</span>
            </div>
            <div v-if="log.remark" class="timeline-remark">{{ log.remark }}</div>
          </div>
        </el-timeline-item>
      </el-timeline>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button
          v-if="canApprove()"
          type="success"
          :icon="Check"
          @click="handleApprove"
        >
          审批
        </el-button>
        <el-button @click="handleClose">关闭</el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { Check } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'
import { appointmentApi } from '@/api'
import eventBus from '@/utils/eventBus'

const props = defineProps({
  modelValue: Boolean,
  appointmentId: [Number, String]
})

const emit = defineEmits(['update:modelValue', 'approve'])

const authStore = useAuthStore()

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const appointment = ref(null)
const approvalLogs = ref([])

watch(() => props.appointmentId, (id) => {
  if (id && props.modelValue) {
    loadDetail()
  }
}, { immediate: true })

watch(() => props.modelValue, (val) => {
  if (val && props.appointmentId) {
    loadDetail()
    setupEventBusListener()
    setupVisibilityListener()
  } else {
    cleanupEventBusListener()
    cleanupVisibilityListener()
  }
})

// EventBus 事件处理函数
function handleApprovalUpdate(updateData) {
  console.log('[EventBus] 收到审批更新事件:', updateData)

  // 检查是否是当前预约单的更新
  if (updateData.appointment_id === parseInt(props.appointmentId)) {
    console.log('[EventBus] 当前预约单有更新，正在刷新审批历史...')

    // 播放提示音
    playNotificationSound()

    // 显示 Toast 提示
    ElMessage.success('审批状态已更新')

    // 刷新预约单详情
    loadDetail()
  }
}

// 播放提示音
function playNotificationSound() {
  try {
    // 使用 Web Audio API 生成简单的"叮"声
    if ('AudioContext' in window || 'webkitAudioContext' in window) {
      const AudioContext = window.AudioContext || window.webkitAudioContext
      const audioCtx = new AudioContext()

      // 创建振荡器（发声源）
      const oscillator = audioCtx.createOscillator()
      const gainNode = audioCtx.createGain()

      // 设置音调（800Hz）和类型（正弦波）
      oscillator.type = 'sine'
      oscillator.frequency.value = 800

      // 设置音量（开始时 0.3）
      gainNode.gain.value = 0.3

      // 连接节点
      oscillator.connect(gainNode)
      gainNode.connect(audioCtx.destination)

      // 播放 200ms 的提示音
      oscillator.start()
      setTimeout(() => {
        oscillator.stop()
        audioCtx.close()
      }, 200)

      console.log('[提示音] 播放成功')
    } else {
      console.log('[提示音] 浏览器不支持 Web Audio API')
    }
  } catch (error) {
    console.log('[提示音] 播放失败:', error)
  }
}

// 设置 EventBus 监听
let unsubscribeFromEventBus = null

function setupEventBusListener() {
  // 清除之前的监听
  cleanupEventBusListener()

  // 监听审批更新事件
  unsubscribeFromEventBus = eventBus.on('appointment:approval-updated', handleApprovalUpdate)
  console.log('[EventBus] 已注册审批更新监听')
}

// 清理 EventBus 监听
function cleanupEventBusListener() {
  if (unsubscribeFromEventBus) {
    unsubscribeFromEventBus()
    unsubscribeFromEventBus = null
    console.log('[EventBus] 已清理审批更新监听')
  }
}

// 页面可见性检测
function setupVisibilityListener() {
  // 监听页面可见性变化
  document.addEventListener('visibilitychange', handleVisibilityChange)

  // 监听对话框获得焦点
  window.addEventListener('focus', handlePageFocus)
}

// 清理页面可见性监听
function cleanupVisibilityListener() {
  document.removeEventListener('visibilitychange', handleVisibilityChange)
  window.removeEventListener('focus', handlePageFocus)
}

// 处理页面可见性变化
function handleVisibilityChange() {
  // 对话框打开且页面重新可见时刷新
  if (!document.hidden && props.modelValue && props.appointmentId) {
    console.log('[页面可见] 页面重新可见，刷新数据...')
    refreshDataIfNeeded()
  }
}

// 处理页面获得焦点
function handlePageFocus() {
  // 对话框打开时刷新
  if (props.modelValue && props.appointmentId) {
    console.log('[页面焦点] 页面获得焦点，刷新数据...')
    refreshDataIfNeeded()
  }
}

// 刷新数据（如果需要）
function refreshDataIfNeeded() {
  if (props.appointmentId) {
    loadDetail()
  }
}

// 组件卸载时清理
onUnmounted(() => {
  cleanupEventBusListener()
  cleanupVisibilityListener()
  console.log('[清理] 已移除所有监听器')
})

async function loadDetail() {
  try {
    const response = await appointmentApi.getDetail(props.appointmentId)
    appointment.value = response.data

    // 加载审批历史
    const historyRes = await appointmentApi.getApprovalHistory(props.appointmentId)
    approvalLogs.value = historyRes.data || []
  } catch (error) {
    console.error('加载详情失败:', error)
  }
}

function handleClose() {
  emit('update:modelValue', false)
}

function formatDateTime(dateStr, timeSlot) {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  const dateStr2 = date.toLocaleDateString('zh-CN')
  const slots = { morning: '上午', afternoon: '下午', evening: '晚上', full_day: '全天' }
  return `${dateStr2} ${slots[timeSlot] || timeSlot}`
}

function formatWorkType(workType) {
  if (!workType) return '-'
  const workTypeMap = {
    general: '一般作业',
    hot_work: '动火作业',
    high_work: '高处作业',
    excavation: '动土作业',
    confined_space: '受限空间',
    electrical: '临时用电',
    lifting: '吊装作业',
    blind_plate: '盲板抽堵'
  }
  // work_type 是逗号分隔的字符串，如 "general,hot_work"
  return workType.split(',').map(type => workTypeMap[type] || type).join('、')
}

function formatFullDateTime(dateStr) {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}

function getStatusType(status) {
  return appointmentApi.getStatusType(status)
}

function getStatusLabel(status) {
  return appointmentApi.getStatusLabel(status)
}

function getActionType(action) {
  const types = {
    approve: 'success',
    reject: 'danger',
    comment: 'info'
  }
  return types[action] || ''
}

function getActionText(action) {
  const texts = {
    approve: '通过',
    reject: '拒绝',
    comment: '评论'
  }
  return texts[action] || action
}

function canApprove() {
  if (!appointment.value) return false
  return appointment.value.status === 'pending' && authStore.hasPermission('appointment_approve')
}

function handleApprove() {
  emit('approve', appointment.value)
}

</script>

<style scoped>
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

.approval-history {
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #eee;
}

.approval-history h4 {
  margin-bottom: 16px;
  font-size: 16px;
  font-weight: 500;
}

.timeline-content {
  padding-bottom: 8px;
}

.timeline-title {
  font-weight: 500;
  margin-bottom: 4px;
}

.timeline-action {
  display: flex;
  align-items: center;
  gap: 8px;
}

.timeline-remark {
  margin-top: 4px;
  color: #666;
  font-size: 13px;
}
</style>
