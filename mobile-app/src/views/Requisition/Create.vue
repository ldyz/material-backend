<template>
  <div class="requisition-create">
    <van-nav-bar title="新建出库单" left-arrow @click-left="router.back()" />

    <van-form @submit="handleSubmit">
      <van-cell-group inset title="基本信息">
        <van-cell
          title="选择项目"
          :value="selectedProject ? selectedProject.name : '请选择项目'"
          is-link
          @click="showProjectPicker = true"
        />
        <van-field
          v-model="formData.requisition_date"
          name="date"
          label="申请日期"
          placeholder="请选择申请日期"
          readonly
          is-link
          @click="showDatePicker = true"
        />
        <van-field
          v-model="formData.department_name"
          name="department"
          label="申请部门"
          placeholder="请输入申请部门"
        />
        <van-field
          v-model="formData.purpose"
          name="purpose"
          label="用途"
          placeholder="请输入领料用途"
          rows="1"
          autosize
        />
        <!-- 紧急程度开关 -->
        <van-cell title="紧急程度" center>
          <template #right-icon>
            <van-switch v-model="formData.urgent" size="20" />
          </template>
        </van-cell>
        <van-field
          v-model="formData.notes"
          name="notes"
          type="textarea"
          label="备注"
          placeholder="请输入备注（可选）"
          rows="2"
        />
      </van-cell-group>

      <!-- 关联物资计划（必填） -->
      <van-cell-group inset title="关联物资计划（必填）">
        <van-cell
          title="选择计划"
          :value="selectedPlan ? selectedPlan.plan_name : '请选择计划'"
          is-link
          required
          :disabled="!selectedProject"
          @click="openPlanPicker"
        />
        <van-cell v-if="selectedPlan" title="计划编号" :value="selectedPlan.plan_no" />
        <van-cell
          v-if="selectedPlan"
          title="添加计划物料"
          is-link
          @click="openPlanItemPicker"
        >
          <template #value>
            <span class="plan-item-count">{{ formData.items.filter(i => i.from_plan).length }} 项已选</span>
          </template>
        </van-cell>
        <van-cell v-if="!selectedProject" title="请先选择项目" />
      </van-cell-group>

      <!-- 已选物资列表 -->
      <van-cell-group v-if="formData.items.length" inset title="已选物资">
        <van-swipe-cell v-for="(item, index) in formData.items" :key="item.material_id">
          <van-cell :title="item.material_name" :label="getItemLabel(item)">
            <template #value>
              <van-field
                v-model.number="item.requested_quantity"
                type="number"
                :min="0.01"
                :max="item.max_quantity"
                class="quantity-input"
                @blur="validateItemQuantity(item)"
              />
              <span class="unit-label">{{ item.unit }}</span>
            </template>
          </van-cell>
          <template #right>
            <van-button square type="danger" text="删除" @click="removeItem(index)" />
          </template>
        </van-swipe-cell>
        <van-cell title="合计" :value="totalRequestedQuantity" />
      </van-cell-group>

      <!-- 草稿状态提示 -->
      <div v-if="hasDraft" class="draft-tip">
        <van-icon name="info-o" />
        <span>已自动保存草稿</span>
        <van-button size="small" plain type="danger" @click="clearDraft">清除草稿</van-button>
      </div>
    </van-form>

    <!-- 浮动提交按钮 -->
    <div class="submit-section">
      <van-button
        round
        block
        type="primary"
        :loading="submitting"
        @click="handleSubmit"
      >
        提交审批
      </van-button>
    </div>

    <!-- 项目选择器 -->
    <van-popup v-model:show="showProjectPicker" position="bottom" :style="{ height: '60%' }">
      <van-picker
        :columns="projectColumns"
        :loading="loadingProjects"
        @confirm="onProjectConfirm"
        @cancel="showProjectPicker = false"
      />
    </van-popup>

    <!-- 日期选择器 -->
    <van-popup v-model:show="showDatePicker" position="bottom">
      <van-date-picker
        v-model="currentDate"
        @confirm="onDateConfirm"
        @cancel="showDatePicker = false"
      />
    </van-popup>

    <!-- 计划选择器弹窗 -->
    <van-popup v-model:show="showPlanPicker" position="bottom" :style="{ height: '60%' }">
      <div class="plan-picker">
        <van-nav-bar title="选择物资计划" left-text="取消" @click-left="showPlanPicker = false" />
        <van-loading v-if="loadingPlans" class="loading-center" size="24px">加载计划中...</van-loading>
        <van-empty v-else-if="approvedPlans.length === 0" description="暂无可用的已批准计划" />
        <div v-else class="plan-list">
          <van-cell
            v-for="plan in approvedPlans"
            :key="plan.id"
            :title="plan.plan_name"
            :label="plan.project_name || ''"
            is-link
            @click="selectPlan(plan)"
          >
            <template #value>
              <span class="plan-status">{{ formatPlanStatus(plan.status) }}</span>
            </template>
          </van-cell>
        </div>
      </div>
    </van-popup>

    <!-- 计划物料选择弹窗 -->
    <van-popup v-model:show="showPlanItemPicker" position="bottom" :style="{ height: '70%' }">
      <div class="plan-item-picker">
        <van-nav-bar title="选择计划物料" left-text="取消" right-text="确认" @click-left="showPlanItemPicker = false" @click-right="confirmPlanItems" />
        <van-search
          v-model="planItemSearch"
          placeholder="搜索物资名称"
          @update:model-value="filterPlanItems"
        />
        <div class="plan-item-actions">
          <van-button size="small" plain type="primary" @click="toggleSelectAllPlanItems">
            {{ isAllPlanItemsSelected ? '取消全选' : '全选' }}
          </van-button>
        </div>
        <van-empty v-if="filteredPlanItems.length === 0" description="暂无待领料物料" />
        <div v-else class="plan-item-list">
          <van-checkbox-group v-model="selectedPlanItemIds">
            <van-cell
              v-for="item in filteredPlanItems"
              :key="item.material_id"
              clickable
              @click="togglePlanItem(item)"
            >
              <template #title>
                <div class="plan-item-title">
                  <span>{{ item.material_name }}</span>
                </div>
              </template>
              <template #label>
                <div class="plan-item-detail">
                  <span v-if="item.specification">规格: {{ item.specification }}</span>
                  <span>计划: {{ item.planned_quantity }} {{ item.unit }}</span>
                  <span class="pending-quantity">待领: {{ item.pending_quantity }} {{ item.unit }}</span>
                </div>
              </template>
              <template #right-icon>
                <van-checkbox :name="item.material_id" @click.stop />
              </template>
            </van-cell>
          </van-checkbox-group>
        </div>
      </div>
    </van-popup>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { showToast, showLoadingToast, closeToast } from 'vant'
import { createRequisition, getProjects } from '@/api/requisition'
import { getApprovedPlans, getPlanDetail } from '@/api/material_plan'

const router = useRouter()

const DRAFT_KEY = 'requisition_draft'

const showProjectPicker = ref(false)
const showDatePicker = ref(false)
const submitting = ref(false)
const loadingProjects = ref(false)
const hasDraft = ref(false)

const projects = ref([])
const selectedProject = ref(null)

// 计划相关状态
const approvedPlans = ref([])
const selectedPlan = ref(null)
const showPlanPicker = ref(false)
const loadingPlans = ref(false)

// 计划物料选择
const planItems = ref([])
const selectedPlanItemIds = ref([])
const showPlanItemPicker = ref(false)
const planItemSearch = ref('')

const currentDate = ref([new Date().getFullYear(), new Date().getMonth() + 1, new Date().getDate()])
const formData = ref({
  project_id: null,
  plan_id: null,
  department_name: '',
  requisition_date: '',
  purpose: '',
  urgent: false,
  notes: '',
  items: []
})

// 防抖定时器
let debounceTimer = null

// 防抖函数
function debounce(fn, delay = 500) {
  return function(...args) {
    if (debounceTimer) {
      clearTimeout(debounceTimer)
    }
    debounceTimer = setTimeout(() => {
      fn.apply(this, args)
    }, delay)
  }
}

const projectColumns = computed(() =>
  projects.value.map(p => ({ text: p.name, value: p.id }))
)

// 计算申请总量
const totalRequestedQuantity = computed(() => {
  return formData.value.items.reduce((sum, item) => {
    return sum + (item.requested_quantity || 0)
  }, 0).toFixed(2)
})

// 加载项目列表
async function loadProjects() {
  loadingProjects.value = true
  try {
    const response = await getProjects({ page: 1, page_size: 100, show_all: 'true' })
    // 后端返回分页格式：{ success: true, data: [...], pagination: {...} }
    projects.value = response.data?.data || response.data || []
  } catch (error) {
    showToast({ type: 'fail', message: '加载项目失败' })
  } finally {
    loadingProjects.value = false
  }
}

async function onProjectConfirm({ selectedOptions }) {
  const projectId = selectedOptions[0].value
  selectedProject.value = projects.value.find(p => p.id === projectId)
  formData.value.project_id = projectId
  showProjectPicker.value = false
}

function onDateConfirm(result) {
  let dateArray
  if (Array.isArray(result)) {
    dateArray = result
  } else if (result && result.selectedValues) {
    dateArray = result.selectedValues
  } else {
    dateArray = currentDate.value
  }
  const [year, month, day] = dateArray
  const monthStr = month.toString().padStart(2, '0')
  const dayStr = day.toString().padStart(2, '0')
  formData.value.requisition_date = `${year}-${monthStr}-${dayStr}`
  showDatePicker.value = false
}

// ========== 计划相关函数 ==========

// 格式化计划状态
function formatPlanStatus(status) {
  const statusMap = {
    'draft': '草稿',
    'pending': '待审批',
    'approved': '已批准',
    'rejected': '已拒绝'
  }
  return statusMap[status] || status
}

// 过滤后的计划物料列表（只显示待领数量 > 0 且未被选中的物料）
const filteredPlanItems = computed(() => {
  // 获取已选物料的 material_id 集合
  const selectedMaterialIds = new Set(
    formData.value.items.map(item => item.material_id)
  )

  // 先过滤掉已领完的和已选中的物料
  let items = planItems.value.filter(item =>
    item.pending_quantity > 0 && !selectedMaterialIds.has(item.material_id)
  )

  // 再处理搜索
  if (planItemSearch.value) {
    const searchLower = planItemSearch.value.toLowerCase()
    items = items.filter(item =>
      item.material_name?.toLowerCase().includes(searchLower)
    )
  }
  return items
})

// 是否全选计划物料
const isAllPlanItemsSelected = computed(() => {
  return filteredPlanItems.value.length > 0 &&
    filteredPlanItems.value.every(item => selectedPlanItemIds.value.includes(item.material_id))
})

// 打开计划选择器
async function openPlanPicker() {
  if (!selectedProject.value) {
    showToast('请先选择项目')
    return
  }

  showPlanPicker.value = true
  await loadApprovedPlans()
}

// 加载已批准的计划列表
async function loadApprovedPlans() {
  if (!selectedProject.value) return

  loadingPlans.value = true
  try {
    const response = await getApprovedPlans({
      project_id: selectedProject.value.id,
      page: 1,
      page_size: 100
    })
    approvedPlans.value = response.data || []
  } catch (error) {
    showToast({ type: 'fail', message: '加载计划失败' })
  } finally {
    loadingPlans.value = false
  }
}

// 选择计划
function selectPlan(plan) {
  // 如果切换了计划，清空之前选择的计划物料
  if (selectedPlan.value && selectedPlan.value.id !== plan.id) {
    formData.value.items = formData.value.items.filter(item => !item.from_plan)
  }

  selectedPlan.value = plan
  formData.value.plan_id = plan.id
  showPlanPicker.value = false

  // 选择计划后自动打开物料选择
  openPlanItemPicker()
}

// 打开计划物料选择弹窗
async function openPlanItemPicker() {
  if (!selectedPlan.value) {
    showToast('请先选择计划')
    return
  }

  showLoadingToast({ message: '加载物料...', forbidClick: true })
  try {
    const response = await getPlanDetail(selectedPlan.value.id)
    planItems.value = (response.data.items || []).map(item => ({
      ...item,
      // 计算待领料数量 = 已入库 - 已领料
      pending_quantity: (item.received_quantity || 0) - (item.issued_quantity || 0)
    }))
    closeToast()
    showPlanItemPicker.value = true
  } catch (error) {
    closeToast()
    showToast({ type: 'fail', message: '加载物料失败' })
  }
}

// 切换计划物料选择
function togglePlanItem(item) {
  const index = selectedPlanItemIds.value.indexOf(item.material_id)
  if (index > -1) {
    selectedPlanItemIds.value.splice(index, 1)
  } else {
    selectedPlanItemIds.value.push(item.material_id)
  }
}

// 全选/取消全选计划物料
function toggleSelectAllPlanItems() {
  if (isAllPlanItemsSelected.value) {
    // 取消全选
    selectedPlanItemIds.value = []
  } else {
    // 全选
    selectedPlanItemIds.value = filteredPlanItems.value.map(item => item.material_id)
  }
}

// 搜索过滤计划物料
function filterPlanItems() {
  // 计算属性自动处理，这里不需要额外逻辑
}

// 确认添加计划物料
function confirmPlanItems() {
  if (selectedPlanItemIds.value.length === 0) {
    showToast('请选择至少一项物料')
    return
  }

  // 先移除之前选择的计划物料
  formData.value.items = formData.value.items.filter(item => !item.from_plan)

  // 直接添加选中的物料
  const newItems = planItems.value
    .filter(item => selectedPlanItemIds.value.includes(item.material_id))
    .map(item => ({
      material_id: item.material_id,
      material_name: item.material_name,
      specification: item.specification,
      material: item.material,
      unit: item.unit,
      requested_quantity: item.pending_quantity || 1,
      max_quantity: item.pending_quantity,  // 添加最大可领数量
      from_plan: true,
      plan_item_id: item.id
    }))

  formData.value.items = [...formData.value.items, ...newItems]
  showPlanItemPicker.value = false
  selectedPlanItemIds.value = []

  if (newItems.length > 0) {
    showToast({ type: 'success', message: `已添加 ${newItems.length} 项物料` })
  }
}

// 获取物料标签信息
function getItemLabel(item) {
  const parts = []
  if (item.specification) parts.push(`规格: ${item.specification}`)
  parts.push(`待领: ${item.max_quantity} ${item.unit || ''}`)
  return parts.join(' | ')
}

// 验证物料数量
function validateItemQuantity(item) {
  if (!item.requested_quantity || item.requested_quantity <= 0) {
    item.requested_quantity = 1
  }
  if (item.max_quantity && item.requested_quantity > item.max_quantity) {
    item.requested_quantity = item.max_quantity
    showToast({ type: 'fail', message: `最大可领 ${item.max_quantity}` })
  }
}

// 删除物料
function removeItem(index) {
  formData.value.items.splice(index, 1)
  showToast('已移除')
}

// 保存草稿
function saveDraft() {
  const draft = {
    formData: formData.value,
    selectedProjectId: selectedProject.value?.id,
    savedAt: Date.now()
  }
  localStorage.setItem(DRAFT_KEY, JSON.stringify(draft))
  hasDraft.value = true
}

// 加载草稿
function loadDraft() {
  try {
    const draftStr = localStorage.getItem(DRAFT_KEY)
    if (!draftStr) return false

    const draft = JSON.parse(draftStr)
    // 检查草稿是否过期（24小时）
    if (Date.now() - draft.savedAt > 24 * 60 * 60 * 1000) {
      localStorage.removeItem(DRAFT_KEY)
      return false
    }

    formData.value = draft.formData

    // 恢复项目选择
    if (draft.selectedProjectId) {
      selectedProject.value = projects.value.find(p => p.id === draft.selectedProjectId)
    }

    hasDraft.value = true
    return true
  } catch (e) {
    console.error('加载草稿失败:', e)
    return false
  }
}

// 清除草稿
function clearDraft() {
  localStorage.removeItem(DRAFT_KEY)
  hasDraft.value = false
  formData.value = {
    project_id: null,
    plan_id: null,
    department_name: '',
    requisition_date: new Date().toISOString().split('T')[0],
    purpose: '',
    urgent: false,
    notes: '',
    items: []
  }
  selectedProject.value = null
  selectedPlan.value = null
  approvedPlans.value = []
  planItems.value = []
  selectedPlanItemIds.value = []
  showToast('草稿已清除')
}

// 监听表单变化自动保存草稿
watch(
  () => formData.value,
  () => {
    if (formData.value.department_name || formData.value.purpose || formData.value.items.length > 0) {
      saveDraft()
    }
  },
  { deep: true }
)

async function handleSubmit() {
  if (!selectedProject.value) {
    showToast({ type: 'fail', message: '请选择项目' })
    return
  }

  if (!selectedPlan.value) {
    showToast({ type: 'fail', message: '请选择物资计划' })
    return
  }

  if (!formData.value.items.length) {
    showToast({ type: 'fail', message: '请选择物资' })
    return
  }

  // 验证数量
  const invalidItems = formData.value.items.filter(item => {
    return !item.requested_quantity || item.requested_quantity <= 0
  })

  if (invalidItems.length) {
    showToast({ type: 'fail', message: '请检查物资数量是否正确' })
    return
  }

  submitting.value = true
  showLoadingToast({ message: '提交中...', forbidClick: true })

  try {
    await createRequisition({
      project_id: formData.value.project_id,
      plan_id: formData.value.plan_id,
      department: formData.value.department_name,  // 字段名改为 department
      purpose: formData.value.purpose,
      urgent: formData.value.urgent,
      remark: formData.value.notes,  // 字段名改为 remark
      items: formData.value.items.map(item => ({
        material_id: item.material_id,
        requested_quantity: item.requested_quantity,
        plan_item_id: item.plan_item_id
      }))
    })

    closeToast()
    // 提交成功后清除草稿
    localStorage.removeItem(DRAFT_KEY)
    showToast({ type: 'success', message: '创建成功' })
    setTimeout(() => {
      router.back()
    }, 1000)
  } catch (error) {
    closeToast()
    showToast({ type: 'fail', message: error.message || error.error || '创建失败' })
  } finally {
    submitting.value = false
  }
}

onMounted(async () => {
  await loadProjects()
  // 默认今天
  formData.value.requisition_date = new Date().toISOString().split('T')[0]

  // 尝试加载草稿
  loadDraft()
})

onUnmounted(() => {
  // 清理防抖定时器
  if (debounceTimer) {
    clearTimeout(debounceTimer)
  }
})
</script>

<style scoped>
.requisition-create {
  min-height: 100vh;
  background-color: #f7f8fa;
  padding-bottom: 80px;
}

.draft-tip {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px;
  background: #fff7e6;
  color: #fa8c16;
  font-size: 13px;
}

.submit-section {
  position: fixed;
  bottom: calc(60px + env(safe-area-inset-bottom));
  left: 0;
  right: 0;
  padding: 16px;
  background: #fff;
  box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.1);
  z-index: 99;
}

/* 计划相关样式 */
.plan-picker {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.plan-list {
  flex: 1;
  overflow-y: auto;
}

.plan-status {
  font-size: 12px;
  color: #52c41a;
}

.plan-item-picker {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.plan-item-actions {
  padding: 8px 16px;
  background: #fff;
  border-bottom: 1px solid #ebedf0;
}

.plan-item-list {
  flex: 1;
  overflow-y: auto;
}

.plan-item-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.plan-item-detail {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 4px;
  font-size: 12px;
  color: #969799;
}

.plan-item-detail span {
  background: #f7f8fa;
  padding: 2px 6px;
  border-radius: 4px;
}

.pending-quantity {
  color: #ff9761 !important;
  background: #fff7e6 !important;
}

.plan-item-count {
  color: #1989fa;
  font-size: 13px;
}

/* 已选物资列表 */
.quantity-input {
  width: 60px;
  padding: 0 8px;
  text-align: center;
}

.unit-label {
  font-size: 12px;
  color: #969799;
  margin-left: 4px;
}
</style>
