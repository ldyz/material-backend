<template>
  <div class="stock-detail-page">
    <van-nav-bar
      title="库存详情"
      left-arrow
      @click-left="onClickLeft"
    />

    <div v-if="loading" class="loading-container">
      <van-loading type="spinner" size="24" />
    </div>

    <div v-else-if="stock" class="detail-content">
      <!-- 库存信息 -->
      <div class="stock-info-card">
        <div class="info-header">
          <h3>{{ stock.material_name || stock.code || stock.material_code || '库存详情' }}</h3>
          <van-tag :type="getStockStatusType(stock)">
            {{ getStockStatusText(stock) }}
          </van-tag>
        </div>
        <div class="quantity-display">
          <span class="quantity-value">{{ stock.quantity }}</span>
          <span class="quantity-unit">{{ stock.unit }}</span>
        </div>
        <div class="stock-range">
          <p>最小: {{ stock.min_stock || 0 }} {{ stock.unit }}</p>
          <p v-if="stock.max_stock">最大: {{ stock.max_stock }} {{ stock.unit }}</p>
        </div>
      </div>

      <!-- 基本信息 -->
      <van-cell-group inset title="基本信息">
        <van-cell title="规格型号" :value="stock.specification || stock.spec || '-'" />
        <van-cell title="材质" :value="stock.code || stock.material_code || '-'" />
        <van-cell title="最后更新" :value="formatDate(stock.updated_at)" />
      </van-cell-group>

      <!-- 图表 -->
      <van-cell-group inset title="库存趋势">
        <div class="chart-container">
          <div class="chart-bars">
            <div
              v-for="(item, index) in chartData"
              :key="index"
              class="chart-bar"
            >
              <div
                class="bar"
                :style="{ height: `${item.percentage}%` }"
              ></div>
              <span class="bar-label">{{ item.label }}</span>
            </div>
          </div>
        </div>
      </van-cell-group>

      <!-- 操作记录 -->
      <van-cell-group inset title="操作记录">
        <van-tabs v-model:active="activeTab" sticky>
          <van-tab title="最近记录" title-style="font-size: 14px">
            <div class="logs-list">
              <div
                v-for="log in recentLogs"
                :key="log.id"
                class="log-item"
              >
                <div class="log-icon" :class="log.type">
                  <van-icon :name="getLogIcon(log.type)" />
                </div>
                <div class="log-info">
                  <p class="log-title">{{ log.title }}</p>
                  <p class="log-detail">
                    <span
                      v-if="log.parsedDetail?.hasLink"
                      class="order-link"
                      @click="goToOrder(log.parsedDetail.orderType, log.parsedDetail.orderNo)"
                    >
                      {{ log.parsedDetail.displayText }}
                    </span>
                    <span v-else>{{ log.detail }}</span>
                  </p>
                  <p class="log-time">{{ formatRelativeTime(log.created_at) }}</p>
                </div>
                <div class="log-quantity" :class="log.type">
                  {{ log.type === 'in' ? '+' : '-' }}{{ log.quantity }}
                </div>
              </div>
              <van-empty
                v-if="recentLogs.length === 0"
                description="暂无记录"
              />
            </div>
          </van-tab>
          <van-tab title="全部记录" title-style="font-size: 14px">
            <div class="all-logs">
              <div
                v-for="log in allLogs"
                :key="log.id"
                class="timeline-item"
              >
                <div class="timeline-date">{{ formatDate(log.created_at) }}</div>
                <p class="timeline-title">{{ log.title }}</p>
                <p class="timeline-detail">
                  <span
                    v-if="log.parsedDetail?.hasLink"
                    class="order-link"
                    @click="goToOrder(log.parsedDetail.orderType, log.parsedDetail.orderNo)"
                  >
                    {{ log.parsedDetail.displayText }}
                  </span>
                  <span v-else>{{ log.detail }}</span>
                </p>
                <p class="timeline-quantity">数量: {{ log.quantity }} {{ stock.unit }}</p>
              </div>
              <van-empty
                v-if="allLogs.length === 0"
                description="暂无记录"
              />
            </div>
          </van-tab>
        </van-tabs>
      </van-cell-group>

      <!-- 操作按钮 -->
      <div
        v-if="canAdjustStock"
        class="action-buttons"
      >
        <van-button
          type="primary"
          block
          @click="showAdjustDialog = true"
        >
          调整库存
        </van-button>
      </div>
    </div>

    <!-- 调整库存弹窗 -->
    <van-dialog
      v-model:show="showAdjustDialog"
      title="调整库存"
      show-cancel-button
      @confirm="onAdjustConfirm"
    >
      <van-field
        v-model.number="adjustQuantity"
        type="number"
        label="调整数量"
        placeholder="正数为增加，负数为减少"
        :rules="[{ required: true, message: '请输入调整数量' }]"
      />
      <van-field
        v-model="adjustReason"
        type="textarea"
        label="调整原因"
        placeholder="请输入调整原因"
        rows="2"
        maxlength="100"
        show-word-limit
      />
    </van-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { showToast, showConfirmDialog } from 'vant'
import { usePermission } from '@/composables/usePermission'
import { getStockDetail, getStockOpLogs, adjustStock } from '@/api/stock'
import { getInboundOrders } from '@/api/inbound'
import { getRequisitions } from '@/api/requisition'
import { formatDate, formatRelativeTime } from '@/utils/date'

const router = useRouter()
const route = useRoute()
const { canAdjustStock } = usePermission()

const loading = ref(true)
const stock = ref(null)
const recentLogs = ref([])
const allLogs = ref([])
const activeTab = ref(0)
const showAdjustDialog = ref(false)
const adjustQuantity = ref(null)
const adjustReason = ref('')

// 模拟图表数据（最近7天）
const chartData = computed(() => {
  if (!stock.value) return []

  // 模拟最近7天的库存数据
  const days = ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
  const baseQuantity = stock.value.quantity

  return days.map((day, index) => {
    const variation = Math.floor(Math.random() * 20) - 10
    const quantity = Math.max(0, baseQuantity + variation)
    const percentage = stock.value.max_stock
      ? (quantity / stock.value.max_stock) * 100
      : 50

    return {
      label: day,
      quantity,
      percentage: Math.min(100, Math.max(0, percentage)),
    }
  })
})

// 获取库存状态类型
function getStockStatusType(item) {
  if (item.quantity === 0) return 'danger'
  if (item.quantity <= item.min_stock) return 'warning'
  return 'success'
}

// 获取库存状态文本
function getStockStatusText(item) {
  if (item.quantity === 0) return '缺货'
  if (item.quantity <= item.min_stock) return '预警'
  return '正常'
}

// 获取日志图标
function getLogIcon(type) {
  const iconMap = {
    in: 'plus',
    out: 'minus',
    adjust: 'exchange',
  }
  return iconMap[type] || 'notes-o'
}

// 加载详情
async function loadDetail() {
  try {
    // 适配统一响应格式
    const { data } = await getStockDetail(route.params.id)
    stock.value = data
  } catch (error) {
    console.error('加载详情失败:', error)
    showToast({
      type: 'fail',
      message: '加载失败',
    })
    router.back()
  } finally {
    loading.value = false
  }
}

// 解析备注，提取单号信息
function parseRemark(remark) {
  if (!remark) return { text: '-', hasLink: false, orderNo: null, orderType: null }

  // 入库单：支持两种格式
  // 格式1：入库单审核入库：RK20250628080103
  // 格式2：入库单审核入库: RK20250905012310
  const inboundMatch = remark.match(/入库单审核入库[:：]\s*(\w+)/)
  if (inboundMatch) {
    return {
      text: remark,
      hasLink: true,
      orderNo: inboundMatch[1],
      orderType: 'inbound',
      displayText: `入库单 ${inboundMatch[1]}`
    }
  }

  // 出库单：支持三种格式
  // 格式1：出库单发放：CK20250804001，出库 8.00 吨
  // 格式2：出库 8.00 个，出库单号：CK20260126005
  // 格式3：出库单发放: CK20250101001
  let outboundMatch = remark.match(/出库单发放[:：]\s*(\w+)/)
  if (!outboundMatch) {
    // 尝试匹配旧格式：出库单号：CK...
    outboundMatch = remark.match(/出库单号[:：]\s*(\w+)/)
  }
  if (outboundMatch) {
    return {
      text: remark,
      hasLink: true,
      orderNo: outboundMatch[1],
      orderType: 'outbound',
      displayText: `出库单 ${outboundMatch[1]}`
    }
  }

  return { text: remark, hasLink: false, orderNo: null, orderType: null }
}

// 跳转到单据详情
async function goToOrder(orderType, orderNo) {
  try {
    if (orderType === 'inbound') {
      // 适配统一响应格式
      const { data } = await getInboundOrders({ inbound_no: orderNo, per_page: 1 })
      const orders = data || []
      if (orders.length > 0) {
        router.push({ name: 'InboundDetail', params: { id: orders[0].id } })
      } else {
        showToast({ type: 'fail', message: '未找到该入库单' })
      }
    } else if (orderType === 'outbound') {
      // 适配统一响应格式
      const { data } = await getRequisitions({ requisition_no: orderNo, per_page: 1 })
      const orders = data || []
      if (orders.length > 0) {
        router.push({ name: 'OutboundDetail', params: { id: orders[0].id } })
      } else {
        showToast({ type: 'fail', message: '未找到该出库单' })
      }
    }
  } catch (error) {
    console.error('查询单据失败:', error)
    showToast({ type: 'fail', message: '查询单据失败' })
  }
}

// 加载操作日志
async function loadLogs() {
  try {
    // 适配统一响应格式
    const { data } = await getStockOpLogs({
      stock_id: route.params.id,
      per_page: 10,
    })

    const logs = data || []

    // 分类日志
    recentLogs.value = logs.map(log => {
      const parsed = parseRemark(log.remark)
      return {
        id: log.id,
        type: log.type, // in, out (使用 type 字段)
        title: getLogTitle(log.type),
        detail: parsed.text,
        parsedDetail: parsed,
        quantity: log.quantity,
        created_at: log.time || log.created_at,
      }
    })

    // 全部日志（用于时间线）
    allLogs.value = logs.map(log => {
      const parsed = parseRemark(log.remark)
      return {
        id: log.id,
        title: getLogTitle(log.type),
        detail: parsed.text,
        parsedDetail: parsed,
        quantity: log.quantity,
        created_at: log.time || log.created_at,
      }
    })
  } catch (error) {
    console.error('加载日志失败:', error)
  }
}

// 获取日志标题
function getLogTitle(type) {
  const titleMap = {
    in: '入库',
    out: '出库',
    adjust: '调整',
  }
  return titleMap[type] || '操作'
}

// 确认调整库存
async function onAdjustConfirm() {
  if (!adjustQuantity.value && adjustQuantity.value !== 0) {
    showToast('请输入调整数量')
    return
  }

  if (!adjustReason.value.trim()) {
    showToast('请输入调整原因')
    return
  }

  try {
    await adjustStock(route.params.id, {
      quantity: adjustQuantity.value,
      reason: adjustReason.value,
    })

    showToast({
      type: 'success',
      message: '调整成功',
    })

    // 重新加载数据
    loadDetail()
    loadLogs()

    adjustQuantity.value = null
    adjustReason.value = ''
  } catch (error) {
    showToast({
      type: 'fail',
      message: '调整失败',
    })
  }
}

// 返回
function onClickLeft() {
  router.back()
}

onMounted(() => {
  loadDetail()
  loadLogs()
})
</script>

<style scoped>
.stock-detail-page {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.loading-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 200px;
}

.detail-content {
  padding: 16px 0 80px 0;
}

.stock-info-card {
  margin: 16px;
  padding: 16px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  color: white;
}

.info-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.info-header h3 {
  font-size: 18px;
  font-weight: bold;
  margin: 0;
}

.quantity-display {
  text-align: center;
  margin: 20px 0;
}

.quantity-value {
  font-size: 48px;
  font-weight: bold;
}

.quantity-unit {
  font-size: 16px;
  margin-left: 8px;
}

.stock-range {
  display: flex;
  justify-content: space-around;
  font-size: 14px;
  opacity: 0.9;
}

.stock-range p {
  margin: 0;
}

.chart-container {
  padding: 16px;
  background: #f7f8fa;
  border-radius: 8px;
}

.chart-bars {
  display: flex;
  align-items: flex-end;
  justify-content: space-around;
  height: 150px;
  padding-top: 20px;
}

.chart-bar {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex: 1;
}

.bar {
  width: 24px;
  background: linear-gradient(180deg, #667eea 0%, #764ba2 100%);
  border-radius: 4px 4px 0 0;
  transition: height 0.3s;
}

.bar-label {
  font-size: 11px;
  color: #969799;
  margin-top: 8px;
}

.logs-list {
  padding: 16px;
}

.log-item {
  display: flex;
  align-items: flex-start;
  padding: 12px;
  background: white;
  border-radius: 8px;
  margin-bottom: 12px;
}

.log-icon {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 12px;
  flex-shrink: 0;
}

.log-icon.in {
  background: #e6f7ff;
  color: #1890ff;
}

.log-icon.out {
  background: #fff1f0;
  color: #ff4d4f;
}

.log-icon.adjust {
  background: #f6ffed;
  color: #52c41a;
}

.log-info {
  flex: 1;
}

.log-title {
  font-size: 14px;
  font-weight: bold;
  color: #323233;
  margin: 0 0 4px 0;
}

.log-detail {
  font-size: 13px;
  color: #646566;
  margin: 0 0 4px 0;
}

.log-time {
  font-size: 12px;
  color: #969799;
  margin: 0;
}

.log-quantity {
  font-size: 16px;
  font-weight: bold;
  flex-shrink: 0;
}

.log-quantity.in {
  color: #07c160;
}

.log-quantity.out,
.log-quantity.adjust {
  color: #ee0a24;
}

.all-logs {
  padding: 16px;
}

.timeline-item {
  padding: 12px;
  background: white;
  border-radius: 8px;
  margin-bottom: 12px;
}

.timeline-date {
  font-size: 12px;
  color: #969799;
  margin-bottom: 8px;
}

.timeline-title {
  font-size: 14px;
  font-weight: bold;
  color: #323233;
  margin: 0 0 4px 0;
}

.timeline-detail {
  font-size: 13px;
  color: #646566;
  margin: 0 0 4px 0;
}

.timeline-quantity {
  font-size: 12px;
  color: #969799;
  margin: 0;
}

.order-link {
  color: #1989fa;
  text-decoration: underline;
  cursor: pointer;
}

.order-link:active {
  opacity: 0.7;
}

.action-buttons {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 12px 16px;
  background: white;
  box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.08);
  z-index: 100;
}
</style>
