<template>
  <div class="inbound-create">
    <van-nav-bar title="新建入库单" left-arrow @click-left="router.back()" />

    <van-form @submit="handleSubmit">
      <van-cell-group inset title="基本信息">
        <!-- 项目选择 -->
        <van-cell
          title="选择项目"
          :value="selectedProject ? selectedProject.name : '请选择项目'"
          is-link
          @click="showProjectPicker = true"
          :rules="[{ required: true, message: '请选择项目' }]"
        />
        <!-- 供应商 -->
        <van-field
          v-model="formData.supplier"
          name="supplier"
          label="供应商"
          placeholder="请输入供应商名称"
          :rules="[{ required: true, message: '请输入供应商名称' }]"
        />
        <!-- 联系人 -->
        <van-field
          v-model="formData.contact"
          name="contact"
          label="联系人"
          placeholder="请输入联系人"
        />
        <!-- 入库日期 -->
        <van-cell
          title="入库日期"
          :value="formData.inbound_date || '请选择日期'"
          is-link
          @click="showDatePicker = true"
        />
        <!-- 备注 -->
        <van-field
          v-model="formData.notes"
          name="notes"
          type="textarea"
          label="备注"
          placeholder="请输入备注（可选）"
          rows="2"
        />
      </van-cell-group>

      <!-- 关联计划（必填） -->
      <van-cell-group inset title="关联物资计划（必填）">
        <van-cell
          title="选择计划"
          :value="selectedPlan ? selectedPlan.plan_name : '请选择计划'"
          is-link
          required
          @click="showPlanPicker = true"
        />
        <van-cell v-if="selectedPlan" title="计划编号" :value="selectedPlan.plan_no" />
        <van-cell
          v-if="selectedPlan"
          title="添加计划物料"
          is-link
          @click="loadPlanItemsForSelection(selectedPlan.id)"
        >
          <template #value>
            <span class="add-plan-material">点击选择物料</span>
          </template>
        </van-cell>
      </van-cell-group>

      <!-- 物资明细 -->
      <van-cell-group inset title="物资明细">
        <template #extra>
          <van-button size="small" type="primary" @click="openMaterialPicker">
            添加物资
          </van-button>
        </template>

        <van-empty v-if="!formData.items.length" description="请添加物资" />

        <div v-for="(item, index) in formData.items" :key="index" class="material-item">
          <van-swipe-cell>
            <div class="item-content">
              <div class="item-header">
                <span class="material-name">{{ item.material_name }}</span>
                <van-icon name="delete" color="#ee0a24" @click="removeItem(index)" />
              </div>
              <div class="item-info">
                <span v-if="item.specification">规格: {{ item.specification }}</span>
                <span v-if="item.material">材质: {{ item.material }}</span>
                <span v-if="item.unit">单位: {{ item.unit }}</span>
              </div>
              <div class="item-inputs">
                <van-field
                  v-model.number="item.quantity"
                  type="number"
                  label="入库数量"
                  placeholder="请输入"
                  :min="0.01"
                  :rules="[{ required: true, message: '请输入数量' }]"
                />
                <van-field
                  v-model.number="item.unit_price"
                  type="number"
                  label="单价(元)"
                  placeholder="请输入"
                  :min="0"
                />
              </div>
            </div>
          </van-swipe-cell>
        </div>
      </van-cell-group>

      <!-- 汇总信息 -->
      <van-cell-group v-if="formData.items.length" inset title="汇总">
        <van-cell title="物资种类" :value="formData.items.length + ' 种'" />
        <van-cell title="入库总额" :value="'¥' + totalAmount.toFixed(2)" />
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
    <van-popup v-model:show="showProjectPicker" position="bottom" :style="{ height: '50%' }">
      <van-picker
        title="选择项目"
        :columns="projectColumns"
        :loading="loadingProjects"
        @confirm="onProjectConfirm"
        @cancel="showProjectPicker = false"
      />
    </van-popup>

    <!-- 计划选择器 -->
    <van-popup v-model:show="showPlanPicker" position="bottom" :style="{ height: '60%' }">
      <div class="picker-header">
        <span>选择物资计划（必填）</span>
      </div>
      <van-list
        v-model:loading="loadingPlans"
        :finished="plansFinished"
        finished-text="没有更多了"
      >
        <van-cell
          v-for="plan in approvedPlans"
          :key="plan.id"
          :title="plan.plan_name"
          :label="`${plan.plan_no} | ${plan.project_name || ''}`"
          is-link
          @click="selectPlan(plan)"
        >
          <template #right-icon>
            <van-icon v-if="selectedPlan?.id === plan.id" name="success" color="#1989fa" />
          </template>
        </van-cell>
      </van-list>
    </van-popup>

    <!-- 物资选择器 -->
    <van-popup v-model:show="showMaterialPicker" position="bottom" :style="{ height: '70%' }">
      <div class="picker-header">
        <van-search
          v-model="materialSearch"
          placeholder="搜索物资名称或编码"
          @update:model-value="debouncedSearchMaterials"
          @clear="searchMaterials"
        />
      </div>
      <van-list
        v-model:loading="loadingMaterials"
        :finished="materialsFinished"
        finished-text="没有更多了"
      >
        <van-cell
          v-for="material in filteredMaterialList"
          :key="material.id"
          :title="material.name"
          :label="`${material.code || ''} | ${material.specification || ''}`"
          is-link
          @click="addMaterial(material)"
        >
          <template #value>
            <span class="material-unit">{{ material.unit }}</span>
          </template>
        </van-cell>
      </van-list>
    </van-popup>

    <!-- 计划物料选择弹窗 -->
    <van-popup v-model:show="showPlanItemPicker" position="bottom" :style="{ height: '70%' }">
      <div class="picker-header">
        <span>选择入库物料</span>
        <van-button size="small" type="primary" @click="selectAllPlanItems">
          {{ isAllPlanItemsSelected ? '取消全选' : '全选' }}
        </van-button>
      </div>

      <van-search
        v-model="planItemSearch"
        placeholder="搜索物资名称"
      />

      <van-list class="plan-item-list">
        <van-checkbox-group v-model="selectedPlanItemIds">
          <van-cell
            v-for="item in filteredPlanItems"
            :key="item.material_id"
            clickable
            @click="togglePlanItem(item.material_id)"
          >
            <template #title>
              <div class="plan-item-title">
                <van-checkbox :name="item.material_id" @click.stop />
                <span>{{ item.material_name }}</span>
              </div>
            </template>
            <template #label>
              <div class="plan-item-info">
                <span v-if="item.specification">规格: {{ item.specification }}</span>
                <span>计划: {{ item.planned_quantity }} {{ item.unit }}</span>
                <span class="pending-qty">待入库: {{ item.pending_quantity }} {{ item.unit }}</span>
              </div>
            </template>
          </van-cell>
        </van-checkbox-group>
      </van-list>

      <div class="picker-footer">
        <van-button block type="primary" @click="confirmPlanItems">
          确认添加 ({{ selectedPlanItemIds.length }}项)
        </van-button>
      </div>
    </van-popup>

    <!-- 日期选择器 -->
    <van-popup v-model:show="showDatePicker" position="bottom">
      <van-date-picker
        v-model="currentDate"
        title="选择入库日期"
        @confirm="onDateConfirm"
        @cancel="showDatePicker = false"
      />
    </van-popup>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { showToast, showLoadingToast, closeToast } from 'vant'
import { createInbound, getProjects, getApprovedPlans, getMaterialMasters, getPlanDetail } from '@/api/inbound'
import { logger } from '@/utils/logger'

const router = useRouter()

const DRAFT_KEY = 'inbound_draft'

const showProjectPicker = ref(false)
const showPlanPicker = ref(false)
const showMaterialPicker = ref(false)
const showDatePicker = ref(false)
const showPlanItemPicker = ref(false)
const submitting = ref(false)
const loadingProjects = ref(false)
const loadingPlans = ref(false)
const loadingMaterials = ref(false)
const plansFinished = ref(false)
const materialsFinished = ref(false)
const hasDraft = ref(false)

const projects = ref([])
const selectedProject = ref(null)
const approvedPlans = ref([])
const selectedPlan = ref(null)
const materialList = ref([])
const materialSearch = ref('')
const planItems = ref([])
const selectedPlanItemIds = ref([])
const planItemSearch = ref('')

const currentDate = ref([new Date().getFullYear(), new Date().getMonth() + 1, new Date().getDate()])
const formData = ref({
  project_id: null,
  supplier: '',
  contact: '',
  inbound_date: '',
  notes: '',
  plan_id: null,
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

// 项目选择器列
const projectColumns = computed(() =>
  projects.value.map(p => ({ text: p.name, value: p.id }))
)

// 计算总金额
const totalAmount = computed(() =>
  formData.value.items.reduce((sum, item) => sum + (item.quantity || 0) * (item.unit_price || 0), 0)
)

// 过滤后的计划物料列表
const filteredPlanItems = computed(() => {
  // 获取已选物料的 material_id 集合（转为字符串确保类型一致）
  const selectedMaterialIds = new Set(
    formData.value.items.map(item => String(item.material_id))
  )

  // 过滤掉已选中的物料
  let items = planItems.value.filter(item =>
    !selectedMaterialIds.has(String(item.material_id))
  )

  // 处理搜索
  if (planItemSearch.value) {
    const search = planItemSearch.value.toLowerCase()
    items = items.filter(item =>
      item.material_name?.toLowerCase().includes(search) ||
      item.specification?.toLowerCase().includes(search)
    )
  }
  return items
})

// 过滤后的普通物资列表（排除已选中的）
const filteredMaterialList = computed(() => {
  const selectedMaterialIds = new Set(
    formData.value.items.map(item => String(item.material_id))
  )
  return materialList.value.filter(item =>
    !selectedMaterialIds.has(String(item.id))
  )
})

// 是否全选
const isAllPlanItemsSelected = computed(() => {
  if (filteredPlanItems.value.length === 0) return false
  return filteredPlanItems.value.every(item =>
    selectedPlanItemIds.value.includes(item.material_id)
  )
})

// 加载项目列表
async function loadProjects() {
  loadingProjects.value = true
  try {
    const response = await getProjects({ page: 1, page_size: 1000, show_all: 'true' })
    // 后端返回分页格式：{ success: true, data: [...], pagination: {...} }
    projects.value = response.data?.data || response.data || []
  } catch (error) {
    showToast({ type: 'fail', message: '加载项目失败' })
  } finally {
    loadingProjects.value = false
  }
}

// 加载已批准的计划列表
async function loadApprovedPlans() {
  if (!selectedProject.value) {
    showToast('请先选择项目')
    return
  }

  loadingPlans.value = true
  plansFinished.value = false
  try {
    const response = await getApprovedPlans({
      project_id: selectedProject.value.id,
      page: 1,
      page_size: 100
    })
    approvedPlans.value = response.data || []
    plansFinished.value = true
  } catch (error) {
    showToast({ type: 'fail', message: '加载计划失败' })
  } finally {
    loadingPlans.value = false
  }
}

// 搜索物资
async function searchMaterials() {
  if (!selectedProject.value) {
    showToast('请先选择项目')
    return
  }

  loadingMaterials.value = true
  materialsFinished.value = false
  try {
    const params = {
      project_id: selectedProject.value.id,
      page: 1,
      page_size: 50
    }
    if (materialSearch.value) {
      params.search = materialSearch.value
    }

    const response = await getMaterialMasters(params)
    materialList.value = response.data || []
    materialsFinished.value = true
  } catch (error) {
    showToast({ type: 'fail', message: '加载物资失败' })
  } finally {
    loadingMaterials.value = false
  }
}

// 防抖搜索
const debouncedSearchMaterials = debounce(searchMaterials, 500)

// 打开物资选择器
function openMaterialPicker() {
  if (!selectedProject.value) {
    showToast('请先选择项目')
    return
  }
  showMaterialPicker.value = true
  searchMaterials()
}

function onProjectConfirm({ selectedOptions }) {
  selectedProject.value = projects.value.find(p => p.id === selectedOptions[0].value)
  formData.value.project_id = selectedOptions[0].value
  showProjectPicker.value = false

  // 清空已选计划和物资
  selectedPlan.value = null
  formData.value.plan_id = null
  formData.value.items = []

  // 加载该项目的已批准计划
  loadApprovedPlans()
}

function selectPlan(plan) {
  // 如果选择了不同的计划，清空已添加的物料
  if (selectedPlan.value && selectedPlan.value.id !== plan.id) {
    formData.value.items = []
  }

  selectedPlan.value = plan
  formData.value.plan_id = plan.id
  showPlanPicker.value = false

  // 加载计划物料并打开选择弹窗
  loadPlanItemsForSelection(plan.id)
}

// 加载计划物料并打开选择弹窗
async function loadPlanItemsForSelection(planId) {
  showLoadingToast({ message: '加载物料...', forbidClick: true })
  try {
    const response = await getPlanDetail(planId)
    planItems.value = (response.data.items || []).map(item => ({
      ...item,
      // 计算待入库数量
      pending_quantity: item.planned_quantity - (item.received_quantity || 0)
    }))
    closeToast()
    showPlanItemPicker.value = true
  } catch (error) {
    closeToast()
    showToast({ type: 'fail', message: '加载物料失败' })
  }
}

// 切换物料选中状态
function togglePlanItem(materialId) {
  const index = selectedPlanItemIds.value.indexOf(materialId)
  if (index === -1) {
    selectedPlanItemIds.value.push(materialId)
  } else {
    selectedPlanItemIds.value.splice(index, 1)
  }
}

// 全选/取消全选
function selectAllPlanItems() {
  if (isAllPlanItemsSelected.value) {
    // 取消全选当前筛选结果
    selectedPlanItemIds.value = selectedPlanItemIds.value.filter(id =>
      !filteredPlanItems.value.some(item => item.material_id === id)
    )
  } else {
    // 全选当前筛选结果
    const newIds = filteredPlanItems.value
      .map(item => item.material_id)
      .filter(id => !selectedPlanItemIds.value.includes(id))
    selectedPlanItemIds.value = [...selectedPlanItemIds.value, ...newIds]
  }
}

// 确认添加选中的物料
function confirmPlanItems() {
  if (selectedPlanItemIds.value.length === 0) {
    showToast('请选择至少一项物料')
    return
  }

  // 获取已选中的物料ID
  const existingIds = formData.value.items.map(item => item.material_id)

  // 将选中的物料填充到表单（排除已存在的）
  const newItems = planItems.value
    .filter(item => selectedPlanItemIds.value.includes(item.material_id) && !existingIds.includes(item.material_id))
    .map(item => ({
      material_id: item.material_id,
      material_name: item.material_name,
      specification: item.specification,
      material: item.material,
      unit: item.unit,
      quantity: item.pending_quantity || item.planned_quantity, // 默认填待入库数量
      unit_price: item.unit_price || 0
    }))

  formData.value.items = [...formData.value.items, ...newItems]

  showPlanItemPicker.value = false
  selectedPlanItemIds.value = []
  planItemSearch.value = ''
}

function clearPlan() {
  selectedPlan.value = null
  formData.value.plan_id = null
  formData.value.items = []
  showPlanPicker.value = false
}

function addMaterial(material) {
  // 检查是否已存在
  if (formData.value.items.some(item => item.material_id === material.id)) {
    showToast('该物资已添加')
    return
  }

  formData.value.items.push({
    material_id: material.id,
    material_name: material.name,
    specification: material.specification,
    material: material.material,
    unit: material.unit,
    quantity: 1,
    unit_price: 0
  })

  showMaterialPicker.value = false
}

function removeItem(index) {
  formData.value.items.splice(index, 1)
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
  formData.value.inbound_date = `${year}-${monthStr}-${dayStr}`
  showDatePicker.value = false
}

// 保存草稿
function saveDraft() {
  const draft = {
    formData: formData.value,
    selectedProjectId: selectedProject.value?.id,
    selectedPlanId: selectedPlan.value?.id,
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

    // 恢复计划选择
    if (draft.selectedPlanId && selectedProject.value) {
      loadApprovedPlans().then(() => {
        selectedPlan.value = approvedPlans.value.find(p => p.id === draft.selectedPlanId)
      })
    }

    hasDraft.value = true
    return true
  } catch (e) {
    logger.error('加载草稿失败:', e)
    return false
  }
}

// 清除草稿
function clearDraft() {
  localStorage.removeItem(DRAFT_KEY)
  hasDraft.value = false
  formData.value = {
    project_id: null,
    supplier: '',
    contact: '',
    inbound_date: new Date().toISOString().split('T')[0],
    notes: '',
    plan_id: null,
    items: []
  }
  selectedProject.value = null
  selectedPlan.value = null
  showToast('草稿已清除')
}

// 监听表单变化自动保存草稿
watch(
  () => formData.value,
  () => {
    if (formData.value.supplier || formData.value.items.length > 0) {
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
    showToast({ type: 'fail', message: '请添加物资明细' })
    return
  }

  // 验证物资数量
  const invalidItems = formData.value.items.filter(item => !item.quantity || item.quantity <= 0)
  if (invalidItems.length) {
    showToast({ type: 'fail', message: '请填写所有物资的数量' })
    return
  }

  submitting.value = true
  showLoadingToast({ message: '提交中...', forbidClick: true })

  try {
    await createInbound({
      project_id: formData.value.project_id,
      plan_id: formData.value.plan_id,
      supplier: formData.value.supplier,
      contact: formData.value.contact,
      notes: formData.value.notes,
      items: formData.value.items.map(item => ({
        material_id: item.material_id,
        quantity: item.quantity,
        unit_price: item.unit_price
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
  formData.value.inbound_date = new Date().toISOString().split('T')[0]

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
.inbound-create {
  min-height: 100vh;
  background-color: #f7f8fa;
  padding-bottom: 80px;
}

.picker-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #fff;
  border-bottom: 1px solid #ebedf0;
}

.material-item {
  background: #fff;
  margin-bottom: 8px;
}

.item-content {
  padding: 12px 16px;
}

.item-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.material-name {
  font-size: 15px;
  font-weight: 500;
  color: #323233;
}

.item-info {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  font-size: 12px;
  color: #969799;
  margin-bottom: 12px;
}

.item-info span {
  background: #f7f8fa;
  padding: 2px 8px;
  border-radius: 4px;
}

.item-inputs {
  display: flex;
  gap: 12px;
}

.item-inputs .van-field {
  flex: 1;
  padding: 8px 0;
}

.material-unit {
  font-size: 12px;
  color: #969799;
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

.plan-item-list {
  height: calc(100% - 140px);
  overflow-y: auto;
}

.plan-item-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.plan-item-info {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 4px;
  font-size: 12px;
  color: #969799;
}

.plan-item-info span {
  background: #f7f8fa;
  padding: 2px 6px;
  border-radius: 4px;
}

.pending-qty {
  color: #1989fa;
  background: #e8f4ff !important;
}

.picker-footer {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 12px 16px;
  background: #fff;
  border-top: 1px solid #ebedf0;
}

.add-plan-material {
  color: #1989fa;
  font-size: 13px;
}
</style>
