<template>
  <!-- 作业人员选择器 -->
  <van-popup v-model:show="show" position="bottom" :style="{ height: height }" round>
    <div class="worker-picker">
      <van-nav-bar
        :title="title"
        left-text="取消"
        :right-text="multiple ? '确认' : ''"
        @click-left="handleCancel"
        @click-right="handleConfirm"
      />
      <div class="worker-selection-info" v-if="multiple && selectedCount > 0">
        <span class="selected-count">已选择 {{ selectedCount }} 人</span>
      </div>
      <div v-if="!multiple && selectedWorker" class="worker-selection-info">
        <span class="selected-count">已选择: {{ selectedWorkerName }}</span>
      </div>
      <van-checkbox-group v-if="multiple" v-model="selectedWorkersList">
        <div class="worker-grid">
          <div
            v-for="worker in workerList"
            :key="worker.id"
            :class="['worker-card', { 'selected': isWorkerSelected(worker.id) }]"
            @click="toggleWorkerSelection(worker.id)"
          >
            <div class="worker-avatar">
              <van-image
                v-if="worker.avatar"
                round
                width="50"
                height="50"
                :src="getAssetUrl(worker.avatar)"
              />
              <div
                v-else
                class="worker-avatar-placeholder"
                :style="{ backgroundColor: getAvatarColor(worker.id) }"
              >
                {{ (worker.full_name || worker.username || '?').charAt(0) }}
              </div>
              <div v-if="isWorkerSelected(worker.id)" class="selected-badge">
                <van-icon name="success" />
              </div>
            </div>
            <div class="worker-name">{{ worker.full_name || worker.username }}</div>
            <van-checkbox :name="worker.id" ref="workerCheckboxes" @click.stop />
          </div>
          <div v-if="workerList.length === 0" class="empty-workers">
            <van-empty description="暂无作业人员" />
          </div>
        </div>
      </van-checkbox-group>

      <!-- 单选模式 -->
      <van-radio-group v-else v-model="selectedWorkerId">
        <div class="worker-grid">
          <div
            v-for="worker in workerList"
            :key="worker.id"
            :class="['worker-card', { 'selected': selectedWorkerId === worker.id }]"
            @click="selectWorker(worker)"
          >
            <div class="worker-avatar">
              <van-image
                v-if="worker.avatar"
                round
                width="50"
                height="50"
                :src="getAssetUrl(worker.avatar)"
              />
              <div
                v-else
                class="worker-avatar-placeholder"
                :style="{ backgroundColor: getAvatarColor(worker.id) }"
              >
                {{ (worker.full_name || worker.username || '?').charAt(0) }}
              </div>
              <div v-if="selectedWorkerId === worker.id" class="selected-badge">
                <van-icon name="success" />
              </div>
            </div>
            <div class="worker-name">{{ worker.full_name || worker.username }}</div>
            <van-radio :name="worker.id" @click.stop />
          </div>
          <div v-if="workerList.length === 0" class="empty-workers">
            <van-empty description="暂无作业人员" />
          </div>
        </div>
      </van-radio-group>
    </div>
  </van-popup>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { getWorkersList, getAvailableWorkers } from '@/api/appointment'
import { getAssetUrl } from '@/utils/request'

const props = defineProps({
  // 是否显示
  modelValue: Boolean,
  // 标题
  title: {
    type: String,
    default: '选择作业人员'
  },
  // 是否多选
  multiple: {
    type: Boolean,
    default: false
  },
  // 弹窗高度
  height: {
    type: String,
    default: '70%'
  },
  // 工作日期（用于获取可用作业人员）
  workDate: String,
  // 时间段（用于获取可用作业人员）
  timeSlot: String
})

const emit = defineEmits(['update:modelValue', 'confirm', 'cancel'])

const show = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

// 作业人员列表
const workerList = ref([])

// 多选模式
const selectedWorkersList = ref([])
const selectedCount = computed(() => selectedWorkersList.value.length)

// 单选模式
const selectedWorkerId = ref(null)
const selectedWorker = computed(() => {
  if (!selectedWorkerId.value) return null
  return workerList.value.find(w => w.id === selectedWorkerId.value)
})
const selectedWorkerName = computed(() => {
  return selectedWorker.value?.full_name || selectedWorker.value?.username || ''
})

// 监听显示状态，加载作业人员列表
watch(() => props.modelValue, async (val) => {
  if (val) {
    await loadWorkers()
  }
})

// 加载作业人员列表
async function loadWorkers() {
  try {
    let response
    // 如果提供了日期和时间段，获取可用作业人员
    if (props.workDate && props.timeSlot) {
      response = await getAvailableWorkers({
        work_date: props.workDate,
        time_slot: props.timeSlot
      })
    } else {
      // 否则获取所有作业人员
      response = await getWorkersList()
    }
    workerList.value = response.data || []
  } catch (error) {
    console.error('获取作业人员列表失败:', error)
    workerList.value = []
  }
}

// 多选：判断是否选中
function isWorkerSelected(workerId) {
  return selectedWorkersList.value.includes(workerId)
}

// 多选：切换选择状态
function toggleWorkerSelection(workerId) {
  const index = selectedWorkersList.value.indexOf(workerId)
  if (index > -1) {
    selectedWorkersList.value.splice(index, 1)
  } else {
    selectedWorkersList.value.push(workerId)
  }
}

// 单选：选择作业人员
function selectWorker(worker) {
  selectedWorkerId.value = worker.id
}

// 获取头像颜色
function getAvatarColor(id) {
  const colors = [
    '#f56c6c', '#5ac8fa', '#5cd65e', '#fab94b',
    '#ff976a', '#67c23a', '#e6a23c', '#909399'
  ]
  return colors[id % colors.length]
}

// 确认选择
function handleConfirm() {
  if (props.multiple) {
    emit('confirm', selectedWorkersList.value, workerList.value)
  } else {
    emit('confirm', selectedWorkerId.value, selectedWorker.value)
  }
  show.value = false
}

// 取消
function handleCancel() {
  emit('cancel')
  show.value = false
}

// 暴露方法给父组件
defineExpose({
  loadWorkers,
  clearSelection: () => {
    selectedWorkersList.value = []
    selectedWorkerId.value = null
  }
})
</script>

<style scoped>
.worker-picker {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.worker-selection-info {
  padding: 8px 16px;
  background: #f5f5f5;
  text-align: center;
}

.selected-count {
  color: #1989fa;
  font-size: 14px;
}

.worker-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
  padding: 16px;
  overflow-y: auto;
  flex: 1;
}

.worker-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 12px 8px;
  background: white;
  border: 2px solid #e5e5e5;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.worker-card:active {
  transform: scale(0.95);
}

.worker-card.selected {
  border-color: #1989fa;
  background: #ecf5ff;
}

.worker-avatar {
  position: relative;
  margin-bottom: 8px;
}

.worker-avatar-placeholder {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 20px;
  font-weight: bold;
}

.selected-badge {
  position: absolute;
  top: -4px;
  right: -4px;
  width: 18px;
  height: 18px;
  background: #1989fa;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 12px;
}

.worker-name {
  font-size: 12px;
  color: #323233;
  text-align: center;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  width: 100%;
}

.worker-card .van-checkbox,
.worker-card .van-radio {
  position: absolute;
  top: 8px;
  right: 8px;
  opacity: 0;
}

.empty-workers {
  grid-column: 1 / -1;
  padding: 40px 0;
}
</style>
