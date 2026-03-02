<template>
  <el-dialog
    v-model="dialogVisible"
    title="分配作业人员"
    width="700px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
    <div v-if="appointment">
      <el-alert
        title="预约信息"
        type="info"
        :closable="false"
        style="margin-bottom: 20px"
      >
        <p><strong>预约单号：</strong>{{ appointment.appointment_no }}</p>
        <p><strong>作业时间：</strong>{{ formatDateTime(appointment.work_date, appointment.time_slot) }}</p>
        <p><strong>作业地点：</strong>{{ appointment.work_location }}</p>
      </el-alert>

      <div class="worker-selection">
        <!-- 作业人员选择 -->
        <div class="selection-section">
          <div class="selection-header">
            <span class="label">选择作业人员：</span>
            <span v-if="selectedWorkers.length > 0" class="selected-info">
              已选择 {{ selectedWorkers.length }} 人：{{ getSelectedWorkersNames() }}
            </span>
            <span v-else class="placeholder">请点击下方头像选择作业人员（可多选）</span>
          </div>

          <div v-loading="loading" class="worker-grid">
            <div
              v-for="worker in allWorkers"
              :key="'worker-' + worker.id"
              :class="['worker-card', {
                'selected': isWorkerSelected(worker.id),
                'unavailable': !worker.is_available,
                'is-supervisor': isSupervisorSelected(worker.id)
              }]"
              @click="worker.is_available ? toggleWorker(worker) : null"
            >
              <div class="worker-avatar">
                <el-avatar
                  :size="60"
                  :src="worker.avatar || undefined"
                  :style="!worker.avatar ? { backgroundColor: getAvatarColor(worker.id) } : {}"
                >
                  {{ worker.name?.charAt(0) || '?' }}
                </el-avatar>
                <div v-if="isWorkerSelected(worker.id)" class="selected-badge">
                  <el-icon><Check /></el-icon>
                </div>
                <div v-if="isSupervisorSelected(worker.id)" class="supervisor-badge">
                  <el-icon><User /></el-icon>
                </div>
                <div v-if="!worker.is_available" class="unavailable-mask">
                  <el-icon><Lock /></el-icon>
                </div>
              </div>
              <div class="worker-info">
                <div class="worker-name">{{ worker.name }}</div>
                <el-tag
                  :type="worker.is_available ? 'success' : 'info'"
                  size="small"
                  effect="plain"
                >
                  {{ worker.is_available ? '可用' : '已占用' }}
                </el-tag>
              </div>
            </div>

            <div v-if="allWorkers.length === 0 && !loading" class="empty-state">
              <el-empty description="暂无作业人员" />
            </div>
          </div>
        </div>

        <!-- 监护人选择 -->
        <div class="selection-section supervisor-section">
          <div class="selection-header">
            <span class="label">选择监护人（可选）：</span>
            <span v-if="selectedSupervisor" class="selected-info supervisor-info">
              {{ allWorkers.find(w => w.id === selectedSupervisor)?.name || '已选择' }}
            </span>
            <span v-else class="placeholder">从上述作业人员中选择一名作为监护人</span>
          </div>

          <div class="supervisor-tips">
            <el-icon><InfoFilled /></el-icon>
            <span>监护人负责监督和协调本次作业的安全与质量</span>
          </div>

          <div v-if="selectedWorkers.length > 0" class="available-for-supervisor">
            <div class="supervisor-options">
              <div
                v-for="workerId in selectedWorkers"
                :key="'sup-' + workerId"
                :class="['supervisor-option', { 'active': selectedSupervisor === workerId }]"
                @click="selectSupervisor(workerId)"
              >
                <el-avatar
                  :size="40"
                  :src="allWorkers.find(w => w.id === workerId)?.avatar || undefined"
                  :style="!allWorkers.find(w => w.id === workerId)?.avatar ? { backgroundColor: getAvatarColor(workerId), fontSize: '16px' } : {}"
                >
                  {{ allWorkers.find(w => w.id === workerId)?.name?.charAt(0) || '?' }}
                </el-avatar>
                <span class="supervisor-name">{{ allWorkers.find(w => w.id === workerId)?.name }}</span>
                <el-icon v-if="selectedSupervisor === workerId" class="check-icon"><Check /></el-icon>
              </div>
            </div>
          </div>

          <div v-else class="no-supervisor-tip">
            <el-empty description="请先选择作业人员，然后从中指定监护人" :image-size="60" />
          </div>
        </div>

        <div class="legend">
          <span style="color: #999; font-size: 12px">
            <span style="color: #67c23a">●</span> 绿色表示该时间段可用
            <span style="margin: 0 12px; color: #999">●</span> 灰色表示该时间段已被占用
            <span style="margin: 0 12px; color: #409eff">●</span> 蓝色标记表示已选为监护人
          </span>
        </div>
      </div>
    </div>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button
        type="primary"
        :loading="submitting"
        :disabled="selectedWorkers.length === 0"
        @click="handleSubmit"
      >
        确认分配（{{ selectedWorkers.length }}人{{ selectedSupervisor ? '，含监护人' : '' }}）
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Lock, Check, User, InfoFilled } from '@element-plus/icons-vue'
import { appointmentApi } from '@/api'

const props = defineProps({
  modelValue: Boolean,
  appointment: Object
})

const emit = defineEmits(['update:modelValue', 'success'])

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const loading = ref(false)
const submitting = ref(false)
const allWorkers = ref([])

// 改为支持多选
const selectedWorkers = ref([])
// 监护人选择
const selectedSupervisor = ref(null)

watch(() => props.modelValue, async (val) => {
  if (val && props.appointment) {
    await loadWorkers()
  }
})

async function loadWorkers() {
  if (!props.appointment) return
  loading.value = true
  try {
    const res = await appointmentApi.getAvailableWorkers({
      work_date: props.appointment.work_date,
      time_slot: props.appointment.time_slot
    })

    allWorkers.value = res.data || []

    // 如果已有分配的作业人员，自动选中
    if (props.appointment.assigned_worker_ids && props.appointment.assigned_worker_ids.length > 0) {
      selectedWorkers.value = [...props.appointment.assigned_worker_ids]
    } else if (props.appointment.assigned_worker_id) {
      selectedWorkers.value = [props.appointment.assigned_worker_id]
    } else {
      selectedWorkers.value = []
    }

    // 加载监护人信息
    if (props.appointment.supervisor_id) {
      selectedSupervisor.value = props.appointment.supervisor_id
    }
  } catch (error) {
    console.error('获取作业人员失败:', error)
    ElMessage.error('获取作业人员失败')
  } finally {
    loading.value = false
  }
}

function toggleWorker(worker) {
  // 不可用的作业人员无法被选中（通过@click条件判断阻止）
  if (!worker.is_available) {
    return
  }

  const index = selectedWorkers.value.indexOf(worker.id)
  if (index > -1) {
    // 取消选中
    selectedWorkers.value.splice(index, 1)
    // 如果取消的是监护人，同时取消监护人
    if (selectedSupervisor.value === worker.id) {
      selectedSupervisor.value = null
    }
  } else {
    // 选中
    selectedWorkers.value.push(worker.id)
  }
}

// 选择监护人
function selectSupervisor(workerId) {
  if (selectedSupervisor.value === workerId) {
    // 取消选择监护人
    selectedSupervisor.value = null
  } else {
    // 选择新的监护人
    selectedSupervisor.value = workerId
  }
}

// 检查是否是监护人
function isSupervisorSelected(workerId) {
  return selectedSupervisor.value === workerId
}

function isWorkerSelected(workerId) {
  return selectedWorkers.value.includes(workerId)
}

function getSelectedWorkersNames() {
  const selected = allWorkers.value.filter(w => selectedWorkers.value.includes(w.id))
  return selected.map(w => w.name).join('、')
}

function getAvatarColor(id) {
  const colors = [
    '#409EFF', '#67C23A', '#E6A23C', '#F56C6C',
    '#909399', '#5470c6', '#91cc75', '#fac858',
    '#ee6666', '#73c0de', '#3ba272', '#fc8452'
  ]
  return colors[id % colors.length]
}

async function handleSubmit() {
  if (selectedWorkers.value.length === 0) {
    ElMessage.warning('请选择作业人员')
    return
  }

  submitting.value = true
  try {
    await appointmentApi.assignWorker(props.appointment.id, {
      worker_ids: selectedWorkers.value,
      supervisor_id: selectedSupervisor.value || undefined
    })

    const supervisorName = selectedSupervisor.value
      ? allWorkers.value.find(w => w.id === selectedSupervisor.value)?.name
      : null

    ElMessage.success(
      `成功分配给 ${selectedWorkers.value.length} 位作业人员` +
      (supervisorName ? `，${supervisorName} 为监护人` : '')
    )
    emit('success')
    handleClose()
  } catch (error) {
    ElMessage.error(error.message || '分配失败')
  } finally {
    submitting.value = false
  }
}

function handleClose() {
  selectedWorkers.value = []
  selectedSupervisor.value = null
  emit('update:modelValue', false)
}

function formatDateTime(dateStr, timeSlot) {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  const dateStr2 = date.toLocaleDateString('zh-CN')
  const slots = { morning: '上午', afternoon: '下午', evening: '晚上', full_day: '全天' }
  return `${dateStr2} ${slots[timeSlot] || timeSlot}`
}
</script>

<style scoped>
.worker-selection {
  padding: 10px 0;
}

.selection-header {
  margin-bottom: 20px;
  padding: 12px 16px;
  background: #f5f7fa;
  border-radius: 4px;
  display: flex;
  align-items: center;
}

.selection-header .label {
  font-weight: 500;
  color: #303133;
  margin-right: 12px;
}

.selection-header .selected-info {
  color: #409EFF;
  font-weight: 500;
}

.selection-header .placeholder {
  color: #909399;
  font-size: 14px;
}

.worker-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: 16px;
  min-height: 150px;
}

.worker-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 16px 12px;
  border: 2px solid #e4e7ed;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
  background: #fff;
}

.worker-card:hover {
  border-color: #409EFF;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.2);
}

.worker-card.selected {
  border-color: #409EFF;
  background: #ecf5ff;
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.2);
}

.worker-card.unavailable {
  opacity: 0.6;
  cursor: not-allowed;
}

.worker-card.unavailable:hover {
  border-color: #e4e7ed;
  transform: none;
  box-shadow: none;
}

.worker-avatar {
  position: relative;
  margin-bottom: 12px;
}

.selected-badge {
  position: absolute;
  top: -5px;
  right: -5px;
  background: #67c23a;
  border-radius: 50%;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  border: 2px solid #fff;
}

.unavailable-mask {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.4);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 20px;
}

.worker-info {
  text-align: center;
  width: 100%;
}

.worker-name {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
  margin-bottom: 8px;
  word-break: break-all;
}

.empty-state {
  grid-column: 1 / -1;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 40px 0;
}

/* 监护人选择区域样式 */
.selection-section {
  margin-bottom: 24px;
}

.selection-section:last-child {
  margin-bottom: 0;
}

.supervisor-section {
  padding: 16px;
  background: #f0f9ff;
  border: 1px solid #b3d8ff;
  border-radius: 8px;
}

.selection-header .supervisor-info {
  color: #409eff;
}

.supervisor-tips {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 16px;
  padding: 8px 12px;
  background: #ecf5ff;
  border-radius: 4px;
  color: #409eff;
  font-size: 13px;
}

.available-for-supervisor {
  margin-top: 12px;
}

.supervisor-options {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.supervisor-option {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  background: white;
  border: 2px solid #dcdfe6;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.supervisor-option:hover {
  border-color: #409eff;
  background: #ecf5ff;
}

.supervisor-option.active {
  border-color: #409eff;
  background: #e6f7ff;
  box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.2);
}

.supervisor-name {
  font-size: 13px;
  font-weight: 500;
  color: #303133;
}

.check-icon {
  margin-left: auto;
  color: #409eff;
  font-size: 16px;
}

.no-supervisor-tip {
  padding: 20px 0;
  text-align: center;
}

/* 监护人标记 */
.worker-card.is-supervisor {
  position: relative;
}

.worker-card.is-supervisor::after {
  content: '监护人';
  position: absolute;
  top: 8px;
  right: 8px;
  background: #409eff;
  color: white;
  font-size: 10px;
  padding: 2px 6px;
  border-radius: 10px;
  font-weight: 500;
}

.supervisor-badge {
  position: absolute;
  top: -8px;
  right: -8px;
  width: 22px;
  height: 22px;
  background: #409eff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  border: 2px solid white;
  font-size: 12px;
  z-index: 2;
}

.legend {
  margin-top: 16px;
  text-align: center;
}
</style>
