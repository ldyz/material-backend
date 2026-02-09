<template>
  <div
    class="notification-item"
    :class="{ 'is-unread': !notification.is_read }"
    @click="handleClick"
  >
    <div class="notification-icon">
      <el-icon :color="iconColor">
        <component :is="iconComponent" />
      </el-icon>
    </div>
    <div class="notification-content">
      <div class="notification-header">
        <span class="notification-title">{{ notification.title }}</span>
        <span class="notification-time">{{ timeAgo }}</span>
      </div>
      <div class="notification-text">{{ notification.content }}</div>
      <div v-if="notification.is_read" class="notification-read-badge">已读</div>
    </div>
    <div class="notification-actions">
      <el-dropdown trigger="click" @command="handleCommand">
        <el-icon class="more-icon" @click.stop>
          <MoreFilled />
        </el-icon>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item v-if="!notification.is_read" command="markRead">
              标记为已读
            </el-dropdown-item>
            <el-dropdown-item command="delete" divided>
              删除
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import {
  Check,
  Warning,
  InfoFilled,
  MoreFilled
} from '@element-plus/icons-vue'

const props = defineProps({
  notification: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['mark-read', 'delete', 'click'])

const router = useRouter()

// 根据通知类型获取图标
const iconComponent = computed(() => {
  const type = props.notification.type
  if (type.includes('approve')) {
    return Check
  } else if (type.includes('reject')) {
    return Warning
  } else if (type.includes('stock_alert')) {
    return Warning
  } else {
    return InfoFilled
  }
})

// 根据通知类型获取图标颜色
const iconColor = computed(() => {
  const type = props.notification.type
  if (type.includes('approve')) {
    return '#67c23a'
  } else if (type.includes('reject')) {
    return '#f56c6c'
  } else if (type.includes('stock_alert')) {
    return '#e6a23c'
  } else {
    return '#409eff'
  }
})

// 时间格式化
const timeAgo = computed(() => {
  if (!props.notification.created_at) return ''

  const now = new Date()
  const created = new Date(props.notification.created_at)
  const diffMs = now - created
  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMins / 60)
  const diffDays = Math.floor(diffHours / 24)

  if (diffMins < 1) return '刚刚'
  if (diffMins < 60) return `${diffMins}分钟前`
  if (diffHours < 24) return `${diffHours}小时前`
  if (diffDays < 7) return `${diffDays}天前`

  return created.toLocaleDateString()
})

// 点击通知
const handleClick = () => {
  // 标记为已读
  if (!props.notification.is_read) {
    emit('mark-read', props.notification.id)
  }

  // 处理跳转
  handleNavigation()

  // 发送点击事件
  emit('click', props.notification)
}

// 处理导航跳转
const handleNavigation = () => {
  const data = parseNotificationData()
  if (!data) return

  // 根据业务类型跳转到相应页面
  // 使用查询参数而不是路径参数，因为列表页支持查询参数打开详情
  switch (data.business_type) {
    case 'inbound_order':
      // 支持两种方式：通过 order_no 或 id 查询
      if (data.order_no) {
        router.push({ path: '/inbound', query: { order_no: data.order_no } })
      } else if (data.inbound_no) {
        router.push({ path: '/inbound', query: { order_no: data.inbound_no } })
      } else {
        router.push({ path: '/inbound', query: { id: data.business_id } })
      }
      break
    case 'requisition':
      // 支持两种方式：通过 requisition_no 或 id 查询
      if (data.requisition_no) {
        router.push({ path: '/requisitions', query: { requisition_no: data.requisition_no } })
      } else {
        router.push({ path: '/requisitions', query: { id: data.business_id } })
      }
      break
    case 'material_plan':
      // 物资计划暂时跳转到列表页
      router.push('/material-plans')
      break
    default:
      break
  }
}

// 解析通知数据
const parseNotificationData = () => {
  try {
    return JSON.parse(props.notification.data || '{}')
  } catch {
    return null
  }
}

// 处理下拉菜单命令
const handleCommand = (command) => {
  switch (command) {
    case 'markRead':
      emit('mark-read', props.notification.id)
      break
    case 'delete':
      emit('delete', props.notification.id)
      break
  }
}
</script>

<style scoped>
.notification-item {
  display: flex;
  align-items: flex-start;
  padding: 12px 16px;
  cursor: pointer;
  transition: background-color 0.2s;
  border-bottom: 1px solid #f0f0f0;
  position: relative;
}

.notification-item:hover {
  background-color: #f5f7fa;
}

.notification-item.is-unread {
  background-color: #ecf5ff;
}

.notification-item.is-unread::before {
  content: '';
  position: absolute;
  left: 4px;
  top: 50%;
  transform: translateY(-50%);
  width: 4px;
  height: 40px;
  background-color: #409eff;
  border-radius: 2px;
}

.notification-icon {
  flex-shrink: 0;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 12px;
  background-color: #f0f9ff;
  border-radius: 50%;
}

.notification-content {
  flex: 1;
  min-width: 0;
}

.notification-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.notification-title {
  font-weight: 500;
  color: #303133;
  font-size: 14px;
}

.notification-time {
  font-size: 12px;
  color: #909399;
  white-space: nowrap;
  margin-left: 12px;
}

.notification-text {
  font-size: 13px;
  color: #606266;
  line-height: 1.5;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.notification-read-badge {
  display: inline-block;
  margin-top: 4px;
  padding: 2px 6px;
  background-color: #f0f2f5;
  color: #909399;
  font-size: 11px;
  border-radius: 3px;
}

.notification-actions {
  flex-shrink: 0;
  margin-left: 8px;
}

.more-icon {
  font-size: 18px;
  color: #909399;
  cursor: pointer;
  transition: color 0.2s;
}

.more-icon:hover {
  color: #409eff;
}
</style>
