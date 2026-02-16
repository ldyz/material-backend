<template>
  <van-cell class="list-item-card" is-link @click="handleClick">
    <template #title>
      <div class="item-title">
        <StatusTag :status="item.status" :type="type" />
        <span class="item-number">{{ itemNumber }}</span>
      </div>
    </template>
    <template #label>
      <div class="item-info">
        <div v-if="item.project_name" class="info-row">
          <van-icon name="shop-o" size="14" />
          <span>{{ item.project_name || '-' }}</span>
        </div>
        <div v-if="itemDate" class="info-row">
          <van-icon name="calendar-o" size="14" />
          <span>{{ itemDate }}</span>
        </div>
        <div v-if="item.supplier_name" class="info-row">
          <van-icon name="user-o" size="14" />
          <span>{{ item.supplier_name }}</span>
        </div>
        <div v-if="item.applicant_name" class="info-row">
          <van-icon name="contact" size="14" />
          <span>申请人：{{ item.applicant_name }}</span>
        </div>
        <div v-if="item.work_type" class="info-row">
          <van-icon name="bookmark-o" size="14" />
          <span>{{ item.work_type }}</span>
        </div>
      </div>
    </template>
  </van-cell>
</template>

<script setup>
import { computed } from 'vue'
import StatusTag from './StatusTag.vue'
import { formatDate, formatAppointmentDate } from '@/composables/useDateTime'

const props = defineProps({
  /**
   * 列表项数据
   */
  item: {
    type: Object,
    required: true
  },
  /**
   * 业务类型
   * 可选值: 'inbound', 'plan', 'requisition', 'appointment'
   */
  type: {
    type: String,
    default: 'inbound',
    validator: (value) => ['inbound', 'plan', 'requisition', 'appointment'].includes(value)
  },
  /**
   * 点击事件处理函数
   */
  onClick: {
    type: Function,
    default: null
  }
})

const emit = defineEmits(['click'])

// 单据编号
const itemNumber = computed(() => {
  const numberFields = {
    inbound: 'order_number',
    plan: 'plan_number',
    requisition: 'requisition_number',
    appointment: 'appointment_number'
  }
  const field = numberFields[props.type] || 'order_number'
  return props.item[field] || props.item.number || '-'
})

// 日期显示
const itemDate = computed(() => {
  if (props.type === 'appointment') {
    return formatAppointmentDate(props.item.appointment_date, props.item.time_slot)
  }

  const dateFields = {
    inbound: 'inbound_date',
    plan: 'plan_date',
    requisition: 'requisition_date'
  }
  const field = dateFields[props.type]
  if (field && props.item[field]) {
    return formatDate(props.item[field])
  }

  return null
})

// 处理点击事件
function handleClick() {
  if (props.onClick) {
    props.onClick(props.item)
  } else {
    emit('click', props.item)
  }
}
</script>

<style scoped>
.list-item-card {
  margin-bottom: 8px;
  border-radius: 8px;
  overflow: hidden;
}

.item-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.item-number {
  font-size: 15px;
  font-weight: 500;
  color: #323233;
  flex: 1;
}

.item-info {
  margin-top: 8px;
}

.info-row {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: #646566;
  margin-bottom: 4px;
}

.info-row:last-child {
  margin-bottom: 0;
}

.info-row .van-icon {
  color: #969799;
}
</style>
