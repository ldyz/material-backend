<template>
  <div class="plan-create">
    <van-nav-bar title="新建物资计划" left-arrow @click-left="router.back()" />

    <van-form @submit="handleSubmit">
      <van-cell-group inset title="基本信息">
        <van-field
          v-model="formData.plan_name"
          name="plan_name"
          label="计划名称"
          placeholder="请输入计划名称"
          :rules="[{ required: true, message: '请输入计划名称' }]"
        />
        <van-cell
          title="关联项目"
          :value="selectedProject ? selectedProject.name : '请选择项目'"
          is-link
          @click="showProjectPicker = true"
          :rules="[{ required: true, message: '请选择项目' }]"
        />
        <van-field
          name="plan_type"
          label="计划类型"
          readonly
          is-link
          :value="getPlanTypeLabel(formData.plan_type)"
          @click="showPlanTypePicker = true"
        />
        <van-field
          name="priority"
          label="优先级"
          readonly
          is-link
          :value="getPriorityLabel(formData.priority)"
          @click="showPriorityPicker = true"
        />
        <van-field
          v-model="formData.description"
          name="description"
          type="textarea"
          label="描述"
          placeholder="请输入计划描述（可选）"
          rows="2"
        />
      </van-cell-group>

      <van-cell-group inset title="计划时间（可选）">
        <van-field
          v-model="formData.planned_start_date"
          name="start_date"
          label="开始日期"
          readonly
          is-link
          placeholder="请选择开始日期"
          @click="showStartDatePicker = true"
        />
        <van-field
          v-model="formData.planned_end_date"
          name="end_date"
          label="结束日期"
          readonly
          is-link
          placeholder="请选择结束日期"
          @click="showEndDatePicker = true"
        />
      </van-cell-group>

      <!-- 物资明细 -->
      <van-cell-group inset title="物资明细">
        <div v-if="formData.items.length" class="items-section">
          <div
            v-for="(item, index) in formData.items"
            :key="index"
            class="material-item"
          >
            <van-cell
              :title="item.material_name || item.material"
              :label="`规格：${item.specification || '-'} | 数量：${item.planned_quantity} ${item.unit || ''}`"
            >
              <template #right-icon>
                <van-icon name="delete" class="delete-icon" @click="removeItem(index)" />
              </template>
            </van-cell>
            <div class="item-detail">
              <van-field
                v-model.number="item.planned_quantity"
                type="number"
                label="计划数量"
                placeholder="数量"
                :min="1"
              />
              <van-field
                v-model.number="item.unit_price"
                type="number"
                label="单价(元)"
                placeholder="单价"
              />
              <van-field
                v-model="item.remark"
                label="备注"
                placeholder="备注（可选）"
              />
            </div>
          </div>
        </div>

        <van-cell
          title="添加物资"
          :value="formData.items.length ? '继续添加' : '请添加物资'"
          is-link
          @click="showMaterialPicker = true"
        >
          <template #icon>
            <van-icon name="add-o" class="add-icon" />
          </template>
        </van-cell>
      </van-cell-group>

      <div class="submit-section">
        <van-button
          round
          block
          type="primary"
          native-type="submit"
          :loading="submitting"
        >
          提交审批
        </van-button>
      </div>
    </van-form>

    <!-- 项目选择器 -->
    <van-popup v-model:show="showProjectPicker" position="bottom" :style="{ height: '60%' }">
      <van-picker
        :columns="projectColumns"
        @confirm="onProjectConfirm"
        @cancel="showProjectPicker = false"
      />
    </van-popup>

    <!-- 计划类型选择器 -->
    <van-popup v-model:show="showPlanTypePicker" position="bottom">
      <van-picker
        title="选择计划类型"
        :columns="planTypeColumns"
        @confirm="onPlanTypeConfirm"
        @cancel="showPlanTypePicker = false"
      />
    </van-popup>

    <!-- 优先级选择器 -->
    <van-popup v-model:show="showPriorityPicker" position="bottom">
      <van-picker
        title="选择优先级"
        :columns="priorityColumns"
        @confirm="onPriorityConfirm"
        @cancel="showPriorityPicker = false"
      />
    </van-popup>

    <!-- 开始日期选择器 -->
    <van-popup v-model:show="showStartDatePicker" position="bottom">
      <van-date-picker
        v-model="startDate"
        title="选择开始日期"
        :min-date="minDate"
        @confirm="onStartDateConfirm"
        @cancel="showStartDatePicker = false"
      />
    </van-popup>

    <!-- 结束日期选择器 -->
    <van-popup v-model:show="showEndDatePicker" position="bottom">
      <van-date-picker
        v-model="endDate"
        title="选择结束日期"
        :min-date="minDate"
        @confirm="onEndDateConfirm"
        @cancel="showEndDatePicker = false"
      />
    </van-popup>

    <!-- 物资选择弹窗 -->
    <van-popup v-model:show="showMaterialPicker" position="bottom" :style="{ height: '70%' }">
      <div class="material-picker">
        <van-nav-bar title="选择物资" left-text="取消" @click-left="showMaterialPicker = false">
          <template #right>
            <van-button type="primary" size="small" @click="confirmMaterialSelection">
              确定
            </van-button>
          </template>
        </van-nav-bar>
        <van-search
          v-model="materialSearch"
          placeholder="搜索物资名称"
          show-action
          @search="loadMaterials"
          @clear="loadMaterials"
        />
        <van-loading v-if="loadingMaterials" size="24px" />
        <van-empty v-else-if="!filteredMaterials.length" description="暂无可选物资" />
        <van-checkbox-group v-else v-model="tempSelectedMaterialIds">
          <van-cell-group inset>
            <van-cell
              v-for="material in filteredMaterials"
              :key="material.id"
              clickable
              @click="toggleTempMaterial(material.id)"
            >
              <template #title>
                <div class="material-info">
                  <div class="material-name">{{ material.material }}</div>
                  <div class="material-desc">
                    规格：{{ material.specification || '-' }} | 单位：{{ material.unit || '-' }}
                  </div>
                </div>
              </template>
              <template #right-icon>
                <van-checkbox :name="material.id" />
              </template>
            </van-cell>
          </van-cell-group>
        </van-checkbox-group>
      </div>
    </van-popup>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { showToast, showLoadingToast, closeToast } from 'vant'
import { getProjects, createPlan } from '@/api/material_plan'
import { getMaterials } from '@/api/material'

const router = useRouter()

// 表单数据
const formData = ref({
  plan_name: '',
  project_id: null,
  plan_type: 'procurement',
  priority: 'normal',
  description: '',
  planned_start_date: '',
  planned_end_date: '',
  items: []
})

const submitting = ref(false)
const selectedProject = ref(null)

// 选择器显示状态
const showProjectPicker = ref(false)
const showPlanTypePicker = ref(false)
const showPriorityPicker = ref(false)
const showStartDatePicker = ref(false)
const showEndDatePicker = ref(false)
const showMaterialPicker = ref(false)

// 日期选择
const minDate = new Date()
const startDate = ref(['2026', '01', '01'])
const endDate = ref(['2026', '12', '31'])

// 物资选择相关
const materials = ref([])
const loadingMaterials = ref(false)
const materialSearch = ref('')
const tempSelectedMaterialIds = ref([])

// 项目列表
const projects = ref([])
const projectColumns = ref([])

// 选择器选项
const planTypeColumns = [
  { text: '采购计划', value: 'procurement' },
  { text: '使用计划', value: 'usage' },
  { text: '混合计划', value: 'mixed' }
]

const priorityColumns = [
  { text: '低', value: 'low' },
  { text: '普通', value: 'normal' },
  { text: '高', value: 'high' },
  { text: '紧急', value: 'urgent' }
]

// 过滤后的物资列表（排除已选）
const filteredMaterials = computed(() => {
  const selectedIds = formData.value.items.map(item => item.material_id)
  let result = materials.value.filter(m => !selectedIds.includes(m.id))

  if (materialSearch.value) {
    const keyword = materialSearch.value.toLowerCase()
    result = result.filter(m =>
      (m.material && m.material.toLowerCase().includes(keyword)) ||
      (m.specification && m.specification.toLowerCase().includes(keyword))
    )
  }

  return result
})

// 获取项目列表
async function loadProjects() {
  try {
    const response = await getProjects({ page: 1, page_size: 100 })
    projects.value = response.data || []

    // 格式化为选择器列
    projectColumns.value = projects.value.map(project => ({
      text: project.name,
      value: project.id
    }))
  } catch (error) {
    showToast({ type: 'fail', message: '加载项目失败' })
  }
}

// 获取物资列表
async function loadMaterials() {
  loadingMaterials.value = true
  try {
    const response = await getMaterials({
      search: materialSearch.value,
      page: 1,
      page_size: 100
    })
    materials.value = response.data || []
  } catch (error) {
    showToast({ type: 'fail', message: '加载物资失败' })
  } finally {
    loadingMaterials.value = false
  }
}

// 项目选择确认
function onProjectConfirm({ selectedOptions }) {
  const projectId = selectedOptions[0].value
  selectedProject.value = projects.value.find(p => p.id === projectId)
  formData.value.project_id = projectId
  showProjectPicker.value = false
}

// 计划类型选择确认
function onPlanTypeConfirm({ selectedValue }) {
  formData.value.plan_type = selectedValue[0]
  showPlanTypePicker.value = false
}

// 优先级选择确认
function onPriorityConfirm({ selectedValue }) {
  formData.value.priority = selectedValue[0]
  showPriorityPicker.value = false
}

// 开始日期确认
function onStartDateConfirm() {
  formData.value.planned_start_date = startDate.value.join('-')
  showStartDatePicker.value = false
}

// 结束日期确认
function onEndDateConfirm() {
  formData.value.planned_end_date = endDate.value.join('-')
  showEndDatePicker.value = false
}

// 物资选择
function toggleTempMaterial(materialId) {
  const index = tempSelectedMaterialIds.value.indexOf(materialId)
  if (index > -1) {
    tempSelectedMaterialIds.value.splice(index, 1)
  } else {
    tempSelectedMaterialIds.value.push(materialId)
  }
}

// 确认物资选择
function confirmMaterialSelection() {
  tempSelectedMaterialIds.value.forEach(materialId => {
    const material = materials.value.find(m => m.id === materialId)
    if (material && !formData.value.items.some(item => item.material_id === materialId)) {
      formData.value.items.push({
        material_id: material.id,
        material: material.material,
        material_name: material.material,
        material_code: material.material_code,
        specification: material.specification,
        unit: material.unit,
        planned_quantity: 1,
        unit_price: 0,
        remark: ''
      })
    }
  })

  // 清空临时选择
  tempSelectedMaterialIds.value = []
  materialSearch.value = ''
  showMaterialPicker.value = false
}

function removeItem(index) {
  formData.value.items.splice(index, 1)
}

// 工具函数
function getPlanTypeLabel(type) {
  const map = {
    procurement: '采购计划',
    usage: '使用计划',
    mixed: '混合计划'
  }
  return map[type] || type
}

function getPriorityLabel(priority) {
  const map = {
    low: '低',
    normal: '普通',
    high: '高',
    urgent: '紧急'
  }
  return map[priority] || priority
}

// 提交表单
async function handleSubmit() {
  if (!formData.value.project_id) {
    showToast({ type: 'fail', message: '请选择项目' })
    return
  }

  if (!formData.value.items.length) {
    showToast({ type: 'fail', message: '请添加物资明细' })
    return
  }

  // 过滤空数量
  const validItems = formData.value.items.filter(item => item.planned_quantity > 0)
  if (!validItems.length) {
    showToast({ type: 'fail', message: '请填写物资数量' })
    return
  }

  submitting.value = true
  try {
    await createPlan({
      plan_name: formData.value.plan_name,
      project_id: formData.value.project_id,
      plan_type: formData.value.plan_type,
      priority: formData.value.priority,
      description: formData.value.description,
      planned_start_date: formData.value.planned_start_date || undefined,
      planned_end_date: formData.value.planned_end_date || undefined,
      items: validItems
    })

    showToast({ type: 'success', message: '创建成功' })
    setTimeout(() => {
      router.back()
    }, 1000)
  } catch (error) {
    showToast({ type: 'fail', message: error.message || '创建失败' })
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadProjects()
})
</script>

<style scoped>
.plan-create {
  min-height: 100vh;
  background-color: #f7f8fa;
  padding-bottom: 130px;
}

.material-item {
  background: #fff;
  margin-bottom: 8px;
}

.item-detail {
  padding: 8px 16px;
  background: #f7f8fa;
}

.delete-icon {
  color: #ee0a24;
  font-size: 18px;
}

.add-icon {
  color: #1989fa;
  font-size: 18px;
  margin-right: 8px;
}

.material-picker {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.material-picker .van-checkbox-group {
  flex: 1;
  overflow-y: auto;
}

.material-info {
  display: flex;
  flex-direction: column;
}

.material-name {
  font-size: 14px;
  font-weight: 500;
}

.material-desc {
  font-size: 12px;
  color: #969799;
  margin-top: 4px;
}

.submit-section {
  position: fixed;
  bottom: 50px;
  left: 0;
  right: 0;
  padding: 16px;
  background: #fff;
  box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.1);
  padding-bottom: calc(16px + env(safe-area-inset-bottom));
}
</style>
