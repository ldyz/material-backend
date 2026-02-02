<template>
  <div class="simple-gantt">
    <!-- 日期标尺 -->
    <div class="gantt-timeline-header" ref="timelineHeaderRef">
      <div class="timeline-tasks">任务列表</div>
      <div class="timeline-dates" ref="timelineDatesRef">
        <div
          v-for="date in dateRange"
          :key="date"
          class="timeline-date"
          :style="{ left: date.position + 'px' }"
        >
          {{ date.label }}
        </div>
      </div>
    </div>

    <!-- 任务条 -->
    <div class="gantt-body">
      <div
        v-for="task in tasks"
        :key="task.id"
        class="gantt-row"
      >
        <div class="gantt-task-name">{{ task.name }}</div>
        <div class="gantt-timeline">
          <div
            class="task-bar"
            :style="getTaskBarStyle(task)"
            :title="`${task.name} (${task.start} ~ ${task.end}) - 进度: ${task.progress}%`"
            @click="handleTaskClick(task)"
          >
            <span class="task-bar-label">{{ task.progress }}%</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="loading-overlay">
      <el-icon class="is-loading" :size="40" />
      <p>加载中...</p>
    </div>

    <!-- 空状态 -->
    <div v-if="!loading && tasks.length === 0" class="empty-overlay">
      <el-empty description="暂无任务数据">
        <el-button type="primary" @click="handleRefresh">刷新</el-button>
      </el-empty>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { ElMessage } from 'element-plus'

const props = defineProps({
  projectId: {
    type: [Number, String],
    required: true
  },
  scheduleData: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['task-updated', 'task-selected'])

const timelineHeaderRef = ref(null)
const timelineDatesRef = ref(null)
const loading = ref(false)

// 日期范围
const dateRange = ref([])

// 格式化任务数据
const tasks = computed(() => {
  const taskList = []
  const activities = props.scheduleData.activities || {}

  console.log('SimpleGantt - 原始 scheduleData:', props.scheduleData)
  console.log('SimpleGantt - activities:', activities)

  for (const [key, activity] of Object.entries(activities)) {
    if (!activity.is_dummy) {
      const startDate = new Date(activity.earliest_start * 1000)
      const endDate = new Date(activity.latest_finish * 1000)

      taskList.push({
        id: activity.id,
        name: activity.name || '未命名任务',
        start: formatDate(startDate),
        end: formatDate(endDate),
        progress: Math.round(activity.progress || 0),
        startDate: startDate,
        endDate: endDate
      })
    }
  }

  // 按开始日期排序
  taskList.sort((a, b) => a.startDate - b.startDate)

  console.log('SimpleGantt - 任务列表:', taskList)

  // 计算日期范围
  if (taskList.length > 0) {
    calculateDateRange(taskList)
  }

  return taskList
})

// 格式化日期
const formatDate = (date) => {
  const d = new Date(date)
  const year = d.getFullYear()
  const month = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

// 计算日期范围
const calculateDateRange = (taskList) => {
  if (taskList.length === 0) return

  const startDates = taskList.map(t => t.startDate.getTime())
  const endDates = taskList.map(t => t.endDate.getTime())

  const minDate = new Date(Math.min(...startDates))
  const maxDate = new Date(Math.max(...endDates))

  // 扩展范围，前后各加几天
  minDate.setDate(minDate.getDate() - 2)
  maxDate.setDate(maxDate.getDate() + 7)

  const dates = []
  const currentDate = new Date(minDate)

  while (currentDate <= maxDate) {
    dates.push({
      date: new Date(currentDate),
      label: formatDate(currentDate),
      position: 0 // 将在下面计算
    })
    currentDate.setDate(currentDate.getDate() + 1)
  }

  // 计算每个日期的位置
  const totalDays = dates.length
  const dayWidth = 40 // 每天40px
  const headerWidth = 250 // 任务名称列宽度

  dateRange.value = dates.map((d, index) => ({
    ...d,
    position: headerWidth + (index * dayWidth)
  }))

  console.log('SimpleGantt - 日期范围:', dateRange.value)
}

// 获取任务条样式
const getTaskBarStyle = (task) => {
  const timelineStart = dateRange.value[0]?.date
  if (!timelineStart) return {}

  const taskStart = new Date(task.start)
  const taskEnd = new Date(task.end)

  const dayWidth = 40
  const headerWidth = 250

  const daysDiff = Math.ceil((taskStart - timelineStart) / (1000 * 60 * 60 * 24))
  const duration = Math.ceil((taskEnd - taskStart) / (1000 * 60 * 60 * 24))

  const left = headerWidth + (daysDiff * dayWidth)
  const width = duration * dayWidth

  // 根据进度设置颜色
  let backgroundColor = '#409eff'
  if (task.progress >= 100) {
    backgroundColor = '#67c23a'
  } else if (task.progress >= 80) {
    backgroundColor = '#e6a23c'
  } else if (task.progress >= 50) {
    backgroundColor = '#409eff'
  } else {
    backgroundColor = '#909399'
  }

  return {
    left: left + 'px',
    width: width + 'px',
    backgroundColor
  }
}

// 任务点击
const handleTaskClick = (task) => {
  emit('task-selected', task)
  ElMessage.success(`选中任务: ${task.name}`)
}

// 刷新
const handleRefresh = () => {
  emit('task-updated', null)
}

// 监听数据变化
watch(
  () => props.scheduleData,
  () => {
    console.log('SimpleGantt - scheduleData changed')
  },
  { deep: true }
)
</script>

<style scoped>
.simple-gantt {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #fff;
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  overflow: hidden;
}

/* 时间轴头部 */
.gantt-timeline-header {
  display: flex;
  background: #f5f7fa;
  border-bottom: 1px solid #dcdfe6;
  font-weight: bold;
  flex-shrink: 0;
}

.timeline-tasks {
  width: 250px;
  padding: 12px;
  border-right: 1px solid #dcdfe6;
  flex-shrink: 0;
}

.timeline-dates {
  flex: 1;
  position: relative;
  height: 50px;
}

.timeline-date {
  position: absolute;
  top: 12px;
  font-size: 12px;
  color: #606266;
  text-align: center;
  width: 40px;
  transform: translateX(-50%);
}

/* 任务列表 */
.gantt-body {
  flex: 1;
  overflow-y: auto;
  min-height: 400px;
}

.gantt-row {
  display: flex;
  border-bottom: 1px solid #ebeef5;
  height: 50px;
  transition: background 0.3s;
}

.gantt-row:hover {
  background: #f5f7fa;
}

.gantt-task-name {
  width: 250px;
  padding: 12px;
  border-right: 1px solid #dcdfe6;
  flex-shrink: 0;
  font-size: 13px;
  color: #303133;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.gantt-timeline {
  flex: 1;
  position: relative;
  min-height: 400px;
  background: linear-gradient(
    to right,
    rgba(0, 0, 0, 0.02) 1px,
    rgba(0, 0, 0, 0.02)
  );
}

.task-bar {
  position: absolute;
  height: 30px;
  top: 10px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 12px;
  font-weight: bold;
  cursor: pointer;
  transition: all 0.3s;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.task-bar:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
  filter: brightness(1.1);
}

.task-bar-label {
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.3);
}

/* 加载/空状态覆盖层 */
.loading-overlay,
.empty-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.95);
  z-index: 10;
}
</style>
