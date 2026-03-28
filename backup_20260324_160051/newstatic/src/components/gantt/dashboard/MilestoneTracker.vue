<template>
  <div class="milestone-tracker">
    <div class="milestone-timeline">
      <div
        v-for="(milestone, index) in sortedMilestones"
        :key="milestone.id"
        class="milestone-item"
        :class="{
          'is-completed': milestone.status === 'completed',
          'is-upcoming': milestone.status === 'upcoming',
          'is-overdue': milestone.status === 'overdue'
        }"
      >
        <!-- Timeline Line -->
        <div class="milestone-line"></div>

        <!-- Milestone Node -->
        <div class="milestone-node">
          <div class="milestone-marker">
            <el-icon v-if="milestone.status === 'completed'">
              <CircleCheck />
            </el-icon>
            <el-icon v-else-if="milestone.status === 'overdue'">
              <Warning />
            </el-icon>
            <div v-else class="milestone-dot"></div>
          </div>
        </div>

        <!-- Milestone Content -->
        <div class="milestone-content">
          <div class="milestone-header">
            <h4 class="milestone-title">{{ milestone.name }}</h4>
            <el-tag
              :type="getStatusType(milestone.status)"
              size="small"
            >
              {{ getStatusLabel(milestone.status) }}
            </el-tag>
          </div>

          <div class="milestone-date">
            <el-icon><Calendar /></el-icon>
            <span>{{ formatDate(milestone.date) }}</span>
            <span v-if="milestone.status !== 'completed'" class="days-remaining">
              {{ getDaysRemaining(milestone.date) }}
            </span>
          </div>

          <div v-if="milestone.description" class="milestone-description">
            {{ milestone.description }}
          </div>

          <!-- Progress Bar -->
          <div v-if="milestone.progress !== undefined" class="milestone-progress">
            <el-progress
              :percentage="milestone.progress"
              :status="getProgressStatus(milestone.status)"
              :stroke-width="8"
            />
          </div>

          <!-- Associated Tasks -->
          <div v-if="milestone.tasks?.length" class="milestone-tasks">
            <div class="tasks-label">
              <strong>{{ milestone.tasks.filter(t => t.status === 'completed').length }}</strong>
              / {{ milestone.tasks.length }} tasks completed
            </div>
            <div class="task-avatars">
              <el-tooltip
                v-for="task in milestone.tasks.slice(0, 5)"
                :key="task.id"
                :content="task.name"
                placement="top"
              >
                <el-avatar
                  :size="28"
                  :style="{ backgroundColor: getTaskColor(task.status) }"
                >
                  <el-icon v-if="task.status === 'completed'"><Check /></el-icon>
                  <span v-else>{{ task.name.charAt(0) }}</span>
                </el-avatar>
              </el-tooltip>
              <el-avatar
                v-if="milestone.tasks.length > 5"
                :size="28"
                style="background: #909399"
              >
                +{{ milestone.tasks.length - 5 }}
              </el-avatar>
            </div>
          </div>

          <!-- Next Milestone Indicator -->
          <div v-if="isNextMilestone(milestone)" class="next-milestone-badge">
            <el-tag type="warning" size="small" effect="plain">
              <el-icon><Star /></el-icon>
              Next Milestone
            </el-tag>
          </div>
        </div>
      </div>
    </div>

    <!-- Overall Progress -->
    <div class="milestone-overall">
      <div class="overall-header">
        <h4>Overall Progress</h4>
        <span class="overall-percentage">{{ overallPercentage }}%</span>
      </div>
      <el-progress
        :percentage="overallPercentage"
        :status="overallPercentage >= 100 ? 'success' : undefined"
        :stroke-width="12"
      />
      <div class="overall-stats">
        <div class="stat">
          <span class="stat-label">Completed:</span>
          <span class="stat-value stat-success">{{ completedCount }}</span>
        </div>
        <div class="stat">
          <span class="stat-label">Upcoming:</span>
          <span class="stat-value stat-info">{{ upcomingCount }}</span>
        </div>
        <div class="stat">
          <span class="stat-label">Overdue:</span>
          <span class="stat-value stat-danger">{{ overdueCount }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import {
  Calendar,
  CircleCheck,
  Warning,
  Check,
  Star
} from '@element-plus/icons-vue'
import { format, differenceInDays, parseISO } from 'date-fns'

/**
 * MilestoneTracker Component
 * Displays milestone timeline with status indicators
 *
 * @props {Array} milestones - Array of milestone objects
 * @props {Array} dateRange - Current date range
 * @props {Number} height - Component height
 */

const props = defineProps({
  milestones: {
    type: Array,
    default: () => []
  },
  dateRange: {
    type: Array,
    default: () => [new Date(), new Date()]
  },
  height: {
    type: Number,
    default: 300
  }
})

// Computed
const sortedMilestones = computed(() => {
  return [...props.milestones].sort((a, b) => {
    return new Date(a.date) - new Date(b.date)
  })
})

const completedCount = computed(() => {
  return props.milestones.filter(m => m.status === 'completed').length
})

const upcomingCount = computed(() => {
  return props.milestones.filter(m => m.status === 'upcoming').length
})

const overdueCount = computed(() => {
  return props.milestones.filter(m => m.status === 'overdue').length
})

const overallPercentage = computed(() => {
  if (props.milestones.length === 0) return 0
  return Math.round((completedCount.value / props.milestones.length) * 100)
})

// Methods
const getStatusType = (status) => {
  const types = {
    completed: 'success',
    upcoming: 'info',
    overdue: 'danger'
  }
  return types[status] || 'info'
}

const getStatusLabel = (status) => {
  const labels = {
    completed: 'Completed',
    upcoming: 'Upcoming',
    overdue: 'Overdue'
  }
  return labels[status] || 'Unknown'
}

const getProgressStatus = (status) => {
  if (status === 'completed') return 'success'
  if (status === 'overdue') return 'exception'
  return undefined
}

const formatDate = (dateStr) => {
  return format(parseISO(dateStr), 'MMMM d, yyyy')
}

const getDaysRemaining = (dateStr) => {
  const date = parseISO(dateStr)
  const now = new Date()
  const days = differenceInDays(date, now)

  if (days < 0) {
    return `${Math.abs(days)} days overdue`
  } else if (days === 0) {
    return 'Due today'
  } else if (days === 1) {
    return 'Due tomorrow'
  } else {
    return `${days} days remaining`
  }
}

const getTaskColor = (status) => {
  const colors = {
    completed: '#67C23A',
    in_progress: '#409EFF',
    not_started: '#909399',
    overdue: '#F56C6C'
  }
  return colors[status] || '#909399'
}

const isNextMilestone = (milestone) => {
  if (milestone.status === 'completed') return false

  const now = new Date()
  const hasIncompleteBefore = sortedMilestones.value.some(m =>
    m.id !== milestone.id &&
    m.status !== 'completed' &&
    new Date(m.date) < new Date(milestone.date)
  )

  return !hasIncompleteBefore
}
</script>

<style scoped lang="scss">
.milestone-tracker {
  width: 100%;
  height: 100%;
  overflow-y: auto;
  padding: 16px;

  .milestone-timeline {
    position: relative;

    .milestone-item {
      display: flex;
      gap: 16px;
      margin-bottom: 24px;
      position: relative;

      &:last-child {
        margin-bottom: 0;

        .milestone-line {
          display: none;
        }
      }

      .milestone-line {
        position: absolute;
        left: 19px;
        top: 40px;
        bottom: -24px;
        width: 2px;
        background: #ebeef5;
      }

      .milestone-node {
        flex-shrink: 0;
        width: 40px;
        display: flex;
        justify-content: center;
        padding-top: 4px;
      }

      .milestone-marker {
        width: 40px;
        height: 40px;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 20px;
        color: #fff;
        z-index: 1;
        transition: all 0.3s;
      }

      .milestone-dot {
        width: 16px;
        height: 16px;
        border-radius: 50%;
        background: currentColor;
      }

      &.is-completed {
        .milestone-marker {
          background: #67C23A;
          box-shadow: 0 0 0 4px rgba(103, 194, 58, 0.2);
        }
        .milestone-line {
          background: #67C23A;
        }
      }

      &.is-upcoming {
        .milestone-marker {
          background: #409EFF;
          box-shadow: 0 0 0 4px rgba(64, 158, 255, 0.2);
        }
      }

      &.is-overdue {
        .milestone-marker {
          background: #F56C6C;
          box-shadow: 0 0 0 4px rgba(245, 108, 108, 0.2);
        }
      }

      .milestone-content {
        flex: 1;
        padding: 12px;
        background: #fff;
        border-radius: 8px;
        border: 1px solid #ebeef5;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
        transition: all 0.3s;

        &:hover {
          box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
        }

        .milestone-header {
          display: flex;
          justify-content: space-between;
          align-items: center;
          margin-bottom: 8px;

          .milestone-title {
            margin: 0;
            font-size: 15px;
            font-weight: 600;
            color: #303133;
          }
        }

        .milestone-date {
          display: flex;
          align-items: center;
          gap: 6px;
          font-size: 13px;
          color: #606266;
          margin-bottom: 8px;

          .days-remaining {
            margin-left: auto;
            padding: 2px 8px;
            background: #f5f7fa;
            border-radius: 4px;
            font-size: 12px;
            font-weight: 600;
            color: #409EFF;
          }
        }

        .milestone-description {
          font-size: 13px;
          color: #909399;
          line-height: 1.5;
          margin-bottom: 12px;
        }

        .milestone-progress {
          margin-bottom: 12px;
        }

        .milestone-tasks {
          display: flex;
          justify-content: space-between;
          align-items: center;
          padding-top: 12px;
          border-top: 1px solid #ebeef5;

          .tasks-label {
            font-size: 12px;
            color: #606266;

            strong {
              color: #303133;
            }
          }

          .task-avatars {
            display: flex;
            gap: 4px;

            .el-avatar {
              cursor: pointer;
              transition: transform 0.2s;

              &:hover {
                transform: scale(1.1);
              }
            }
          }
        }

        .next-milestone-badge {
          margin-top: 12px;
          padding-top: 12px;
          border-top: 1px solid #ebeef5;

          :deep(.el-tag) {
            display: inline-flex;
            align-items: center;
            gap: 4px;
          }
        }
      }
    }
  }

  .milestone-overall {
    margin-top: 24px;
    padding: 16px;
    background: #fff;
    border-radius: 8px;
    border: 1px solid #ebeef5;

    .overall-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 12px;

      h4 {
        margin: 0;
        font-size: 14px;
        font-weight: 600;
        color: #303133;
      }

      .overall-percentage {
        font-size: 18px;
        font-weight: 700;
        color: #409EFF;
      }
    }

    .overall-stats {
      display: flex;
      justify-content: space-around;
      margin-top: 16px;

      .stat {
        text-align: center;

        .stat-label {
          display: block;
          font-size: 12px;
          color: #909399;
          margin-bottom: 4px;
        }

        .stat-value {
          display: block;
          font-size: 16px;
          font-weight: 700;
          color: #303133;

          &.stat-success {
            color: #67C23A;
          }

          &.stat-info {
            color: #409EFF;
          }

          &.stat-danger {
            color: #F56C6C;
          }
        }
      }
    }
  }
}

/* Custom Scrollbar */
.milestone-tracker::-webkit-scrollbar {
  width: 6px;
}

.milestone-tracker::-webkit-scrollbar-track {
  background: #f5f7fa;
}

.milestone-tracker::-webkit-scrollbar-thumb {
  background: #dcdfe6;
  border-radius: 3px;
}

.milestone-tracker::-webkit-scrollbar-thumb:hover {
  background: #c0c4cc;
}
</style>
