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
        />
        <van-cell
          title="计划类型"
          :value="planTypeText"
          is-link
          @click="showPlanTypePicker = true"
        />
        <van-cell
          title="优先级"
          :value="priorityText"
          is-link
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

        <van-cell
          title="导入物资"
          value="从Excel导入"
          is-link
          @click="showImportPopup = true"
        >
          <template #icon>
            <van-icon name="description" class="import-icon" />
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
              v-for="item in filteredMaterials"
              :key="item.id"
              clickable
            >
              <template #title>
                <div class="material-info" @click="toggleTempMaterial(item.id)">
                  <div class="material-name">{{ item.name }}</div>
                  <div class="material-desc">
                    规格：{{ item.specification || '-' }} | 单位：{{ item.unit || '-' }}
                  </div>
                </div>
              </template>
              <template #right-icon>
                <van-checkbox :name="item.id" @click.stop />
              </template>
            </van-cell>
          </van-cell-group>
        </van-checkbox-group>
      </div>
    </van-popup>

    <!-- Excel 导入弹窗 -->
    <van-popup
      v-model:show="showImportPopup"
      position="bottom"
      :style="{ height: '90%' }"
      :close-on-click-overlay="false"
    >
      <div class="import-popup">
        <van-nav-bar title="导入物资" left-text="取消" @click-left="closeImportPopup" />

        <!-- 步骤指示器 -->
        <van-steps :active="importStep" class="import-steps">
          <van-step>上传文件</van-step>
          <van-step>选择表头</van-step>
          <van-step>字段映射</van-step>
          <van-step>预览导入</van-step>
        </van-steps>

        <!-- 步骤1: 上传文件 -->
        <div v-show="importStep === 0" class="import-step-content">
          <van-uploader
            v-model="importFileList"
            :max-count="1"
            accept=".xlsx,.xls"
            :after-read="handleFileRead"
          >
            <template #default>
              <div class="upload-area">
                <van-icon name="description" size="40" color="#1989fa" />
                <div class="upload-text">点击上传 Excel 文件</div>
                <div class="upload-hint">支持 .xlsx / .xls 格式</div>
              </div>
            </template>
          </van-uploader>
          <div class="import-tips">
            <div class="tips-title">支持的字段：</div>
            <div class="tips-list">
              <span class="required">*物资名称</span>、
              <span class="required">*单位</span>、
              <span class="required">*计划数量</span>、
              物资编码、规格型号、材质、分类、单价、优先级、需求日期、备注
            </div>
          </div>
        </div>

        <!-- 步骤2: 选择标题行 -->
        <div v-show="importStep === 1" class="import-step-content">
          <div class="step-header">
            <div class="step-title">选择标题行（当前：第 {{ headerRowIndex }} 行）</div>
            <div class="step-desc">点击选择包含列标题的行</div>
          </div>
          <div class="header-preview">
            <van-radio-group v-model="headerRowIndex">
              <van-cell-group inset>
                <van-cell
                  v-for="(row, index) in headerPreviewRows"
                  :key="index"
                  clickable
                  @click="headerRowIndex = index + 1"
                >
                  <template #title>
                    <div class="header-row-content">
                      <span class="row-label">{{ row.join(' | ') }}</span>
                    </div>
                  </template>
                  <template #right-icon>
                    <van-radio :name="index + 1" />
                  </template>
                </van-cell>
              </van-cell-group>
            </van-radio-group>
          </div>
        </div>

        <!-- 步骤3: 字段映射 -->
        <div v-show="importStep === 2" class="import-step-content">
          <div class="step-header">
            <div class="step-title">字段映射</div>
            <div class="step-desc">将 Excel 列映射到对应字段</div>
          </div>
          <div class="mapping-list">
            <van-cell-group inset>
              <van-cell
                v-for="field in importFields"
                :key="field.field"
                :title="field.label"
                is-link
                @click="openFieldMapping(field)"
              >
                <template #label>
                  <span :class="{ 'required-hint': field.required }">
                    {{ field.hint }}
                  </span>
                </template>
                <template #value>
                  <span :class="{ 'mapped': field.mappedColumn, 'unmapped': !field.mappedColumn && field.required }">
                    {{ field.mappedColumn || '未映射' }}
                  </span>
                </template>
              </van-cell>
            </van-cell-group>
          </div>
        </div>

        <!-- 步骤4: 预览导入 -->
        <div v-show="importStep === 3" class="import-step-content">
          <div class="step-header">
            <div class="step-title">数据预览</div>
            <div class="step-desc">共 {{ previewData.length }} 条数据待导入</div>
          </div>
          <div class="preview-list">
            <van-cell-group inset>
              <van-cell
                v-for="(item, index) in previewData"
                :key="index"
                :title="item.material_name || '-'"
                :label="`数量: ${item.planned_quantity || 0} ${item.unit || ''}`"
              >
                <template #value>
                  <span class="preview-price">{{ item.unit_price ? '¥' + item.unit_price : '' }}</span>
                </template>
              </van-cell>
            </van-cell-group>
          </div>
        </div>

        <!-- 底部操作按钮 -->
        <div class="import-actions">
          <van-button v-if="importStep > 0" block @click="importStep--">上一步</van-button>
          <van-button
            v-if="importStep < 3"
            block
            type="primary"
            :disabled="!canGoNext"
            @click="goNextStep"
          >
            下一步
          </van-button>
          <van-button
            v-if="importStep === 3"
            block
            type="primary"
            :loading="importing"
            @click="confirmImport"
          >
            确认导入 ({{ previewData.length }} 条)
          </van-button>
        </div>
      </div>
    </van-popup>

    <!-- 字段映射选择器 -->
    <van-popup v-model:show="showFieldMappingPicker" position="bottom">
      <van-picker
        :title="'映射 ' + currentMappingField?.label"
        :columns="excelColumnOptions"
        @confirm="onFieldMappingConfirm"
        @cancel="showFieldMappingPicker = false"
      />
    </van-popup>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { showToast, showLoadingToast, closeToast } from 'vant'
import { getProjects, createPlan } from '@/api/material_plan'
import { getMaterials } from '@/api/material'
import * as XLSX from 'xlsx'
import { logger } from '@/utils/logger'

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

// 选择器显示文本
const planTypeText = ref('请选择')
const priorityText = ref('请选择')

// 选择器显示状态
const showProjectPicker = ref(false)
const showPlanTypePicker = ref(false)
const showPriorityPicker = ref(false)
const showStartDatePicker = ref(false)
const showEndDatePicker = ref(false)
const showMaterialPicker = ref(false)

// 日期选择
const minDate = new Date()
const today = new Date()
const startDate = ref([
  today.getFullYear().toString(),
  (today.getMonth() + 1).toString().padStart(2, '0'),
  today.getDate().toString().padStart(2, '0')
])
const endDate = ref([
  (today.getFullYear() + 1).toString(),
  '12',
  '31'
])

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

// ========== 导入相关 ==========
const showImportPopup = ref(false)
const importStep = ref(0)
const importFileList = ref([])
const importing = ref(false)

// Excel 数据
const excelRawData = ref([])
const excelColumns = ref([])
const headerRowIndex = ref(1)

// 字段定义（参考 web 端）
const importFields = ref([
  { field: 'material_name', label: '物资名称', required: true, hint: '必填', mappedColumn: null },
  { field: 'material_code', label: '物资编码', required: false, hint: '物资编码', mappedColumn: null },
  { field: 'specification', label: '规格型号', required: false, hint: '规格型号', mappedColumn: null },
  { field: 'material', label: '材质', required: false, hint: '材质说明', mappedColumn: null },
  { field: 'category', label: '分类', required: false, hint: '物资分类', mappedColumn: null },
  { field: 'unit', label: '单位', required: true, hint: '必填', mappedColumn: null },
  { field: 'planned_quantity', label: '计划数量', required: true, hint: '必填', mappedColumn: null },
  { field: 'unit_price', label: '单价', required: false, hint: '单价（元）', mappedColumn: null },
  { field: 'priority', label: '优先级', required: false, hint: '紧急/高/普通/低', mappedColumn: null },
  { field: 'required_date', label: '需求日期', required: false, hint: 'YYYY-MM-DD', mappedColumn: null },
  { field: 'remark', label: '备注', required: false, hint: '备注说明', mappedColumn: null }
])

// 字段映射选择器
const showFieldMappingPicker = ref(false)
const currentMappingField = ref(null)
const excelColumnOptions = computed(() => {
  const options = [{ text: '不映射', value: '' }]
  excelColumns.value.forEach(col => {
    if (col) {
      options.push({ text: col, value: col })
    }
  })
  return options
})

// 标题行预览（前10行）
const headerPreviewRows = computed(() => {
  return excelRawData.value.slice(0, 10).map(row => {
    return row.map(cell => String(cell || '').trim()).slice(0, 5)
  })
})

// 预览数据
const previewData = computed(() => {
  const dataRows = excelRawData.value.slice(headerRowIndex.value)
  return dataRows.map(row => {
    const item = {}
    importFields.value.forEach(field => {
      if (field.mappedColumn) {
        const colIndex = excelColumns.value.indexOf(field.mappedColumn)
        let value = row[colIndex]

        // 数据类型转换
        if (field.field === 'planned_quantity' || field.field === 'unit_price') {
          value = parseFloat(value) || 0
        } else if (value) {
          value = String(value).trim()
        }

        item[field.field] = value
      }
    })
    return item
  }).filter(item => item.material_name) // 过滤掉没有物资名称的行
})

// 是否可以进入下一步
const canGoNext = computed(() => {
  if (importStep.value === 0) {
    return excelRawData.value.length > 0
  }
  if (importStep.value === 1) {
    return headerRowIndex.value > 0
  }
  if (importStep.value === 2) {
    // 检查必填字段是否已映射
    const requiredFields = importFields.value.filter(f => f.required)
    return requiredFields.every(f => f.mappedColumn)
  }
  return true
})

// 过滤后的物资列表（排除已选）
const filteredMaterials = computed(() => {
  const selectedIds = formData.value.items.map(item => item.material_id)
  let result = materials.value.filter(m => !selectedIds.includes(m.id))

  if (materialSearch.value) {
    const keyword = materialSearch.value.toLowerCase()
    result = result.filter(m =>
      (m.name && m.name.toLowerCase().includes(keyword)) ||
      (m.specification && m.specification.toLowerCase().includes(keyword))
    )
  }

  return result
})

// 获取项目列表
async function loadProjects() {
  try {
    const response = await getProjects({ page: 1, page_size: 100, show_all: 'true' })
    // 后端返回分页格式：{ success: true, data: [...], pagination: {...} }
    projects.value = response.data?.data || response.data || []

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
function onPlanTypeConfirm(event) {
  logger.debug('onPlanTypeConfirm event:', event)
  const { selectedOptions, selectedValues } = event
  // 优先从 selectedOptions 获取 text
  if (selectedOptions && selectedOptions.length > 0) {
    formData.value.plan_type = selectedOptions[0].value
    planTypeText.value = selectedOptions[0].text
    logger.debug('set planTypeText:', planTypeText.value)
  } else if (selectedValues && selectedValues.length > 0) {
    // 备用方案：从 selectedValues 获取值，然后查找对应的 text
    const value = selectedValues[0]
    formData.value.plan_type = value
    const column = planTypeColumns.find(c => c.value === value)
    planTypeText.value = column ? column.text : value
    logger.debug('set planTypeText from column:', planTypeText.value)
  }
  showPlanTypePicker.value = false
}

// 优先级选择确认
function onPriorityConfirm(event) {
  logger.debug('onPriorityConfirm event:', event)
  const { selectedOptions, selectedValues } = event
  // 优先从 selectedOptions 获取 text
  if (selectedOptions && selectedOptions.length > 0) {
    formData.value.priority = selectedOptions[0].value
    priorityText.value = selectedOptions[0].text
    logger.debug('set priorityText:', priorityText.value)
  } else if (selectedValues && selectedValues.length > 0) {
    // 备用方案：从 selectedValues 获取值，然后查找对应的 text
    const value = selectedValues[0]
    formData.value.priority = value
    const column = priorityColumns.find(c => c.value === value)
    priorityText.value = column ? column.text : value
    logger.debug('set priorityText from column:', priorityText.value)
  }
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
        material: material.name,
        material_name: material.name,
        material_code: material.code || '',
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

// ========== 导入相关函数 ==========

// 关闭导入弹窗并重置状态
function closeImportPopup() {
  showImportPopup.value = false
  importStep.value = 0
  importFileList.value = []
  excelRawData.value = []
  excelColumns.value = []
  headerRowIndex.value = 1
  // 重置字段映射
  importFields.value.forEach(f => {
    f.mappedColumn = null
  })
}

// 处理文件读取
function handleFileRead(file) {
  showLoadingToast({ message: '正在读取文件...', forbidClick: true })

  const reader = new FileReader()
  reader.onload = (e) => {
    try {
      const data = new Uint8Array(e.target.result)
      const workbook = XLSX.read(data, { type: 'array' })

      // 读取第一个工作表
      const firstSheetName = workbook.SheetNames[0]
      const worksheet = workbook.Sheets[firstSheetName]

      // 转换为数组格式（包含表头）
      const jsonData = XLSX.utils.sheet_to_json(worksheet, {
        header: 1,
        defval: ''
      })

      if (jsonData.length < 2) {
        closeToast()
        showToast({ type: 'fail', message: 'Excel 文件没有数据行' })
        return
      }

      // 保存原始数据
      excelRawData.value = jsonData

      // 默认第一行为标题行，初始化列名
      updateExcelColumns()

      // 自动匹配字段
      autoMatchFields()

      closeToast()
      showToast({ type: 'success', message: `成功读取 ${jsonData.length - 1} 行数据` })
    } catch (error) {
      closeToast()
      showToast({ type: 'fail', message: '读取文件失败：' + error.message })
    }
  }

  reader.readAsArrayBuffer(file.file)
}

// 根据标题行更新列名
function updateExcelColumns() {
  const rowIndex = headerRowIndex.value - 1
  if (excelRawData.value[rowIndex]) {
    excelColumns.value = excelRawData.value[rowIndex].map(col => String(col || '').trim())
  }
}

// 自动匹配字段
function autoMatchFields() {
  importFields.value.forEach(field => {
    let bestMatch = ''
    let bestScore = 0

    excelColumns.value.forEach(col => {
      const score = calculateSimilarity(field.field, field.label, col)
      if (score > bestScore && score >= 60) {
        bestScore = score
        bestMatch = col
      }
    })

    field.mappedColumn = bestMatch || null
  })
}

// 计算相似度
function calculateSimilarity(fieldKey, fieldLabel, excelColumn) {
  const key = fieldKey.toLowerCase()
  const label = fieldLabel.toLowerCase()
  const col = excelColumn.toLowerCase().trim()

  // 完全匹配
  if (col === key || col === label) {
    return 100
  }

  // 包含匹配
  if (col.includes(key) || key.includes(col) || col.includes(label) || label.includes(col)) {
    return 90
  }

  // 别名匹配
  const aliases = {
    'material_name': ['物资名称', '名称', '品名', 'name', 'title'],
    'material_code': ['物资编码', '编码', '编号', 'code', 'no'],
    'specification': ['规格型号', '规格', '型号', 'specification', 'spec', 'model'],
    'material': ['材质', '材料', 'material'],
    'category': ['分类', '类别', 'category', 'type'],
    'unit': ['单位', '计量单位', 'unit'],
    'planned_quantity': ['计划数量', '数量', 'quantity', 'qty'],
    'unit_price': ['单价', '价格', 'price'],
    'priority': ['优先级', 'priority'],
    'required_date': ['需求日期', '日期', 'date'],
    'remark': ['备注', '说明', 'remark', 'note']
  }

  if (aliases[key]) {
    for (const alias of aliases[key]) {
      if (col.includes(alias.toLowerCase()) || alias.toLowerCase().includes(col)) {
        return 80
      }
    }
  }

  return 0
}

// 下一步
function goNextStep() {
  if (importStep.value === 1) {
    // 选择标题行后，更新列名并重新匹配
    updateExcelColumns()
    autoMatchFields()
  }
  importStep.value++
}

// 打开字段映射选择器
function openFieldMapping(field) {
  currentMappingField.value = field
  showFieldMappingPicker.value = true
}

// 确认字段映射
function onFieldMappingConfirm({ selectedOptions }) {
  if (currentMappingField.value) {
    currentMappingField.value.mappedColumn = selectedOptions[0].value
  }
  showFieldMappingPicker.value = false
}

// 确认导入
function confirmImport() {
  if (previewData.value.length === 0) {
    showToast({ type: 'fail', message: '没有可导入的数据' })
    return
  }

  importing.value = true

  try {
    // 将预览数据添加到表单
    let imported = 0
    previewData.value.forEach(item => {
      formData.value.items.push({
        material_name: item.material_name || '',
        material_code: item.material_code || '',
        specification: item.specification || '',
        material: item.material || '',
        category: item.category || '',
        unit: item.unit || '',
        planned_quantity: item.planned_quantity || 1,
        unit_price: item.unit_price || 0,
        priority: item.priority || 'normal',
        required_date: item.required_date || '',
        remark: item.remark || ''
      })
      imported++
    })

    showToast({ type: 'success', message: `成功导入 ${imported} 条物资` })
    closeImportPopup()
  } catch (error) {
    showToast({ type: 'fail', message: '导入失败：' + error.message })
  } finally {
    importing.value = false
  }
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

// 监听物资选择弹窗打开，自动加载物资列表
watch(showMaterialPicker, (val) => {
  if (val && materials.value.length === 0) {
    loadMaterials()
  }
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

.import-icon {
  color: #07c160;
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
  bottom: calc(60px + env(safe-area-inset-bottom));
  left: 0;
  right: 0;
  padding: 16px;
  background: #fff;
  box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.1);
  z-index: 99;
}

/* 导入弹窗样式 */
.import-popup {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.import-steps {
  padding: 16px;
  background: #fff;
}

.import-step-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.upload-area {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  background: #f7f8fa;
  border-radius: 8px;
  border: 1px dashed #dcdee0;
}

.upload-text {
  margin-top: 12px;
  font-size: 14px;
  color: #323233;
}

.upload-hint {
  margin-top: 8px;
  font-size: 12px;
  color: #969799;
}

.import-tips {
  margin-top: 20px;
  padding: 12px;
  background: #fff;
  border-radius: 8px;
}

.tips-title {
  font-size: 14px;
  font-weight: 500;
  color: #323233;
  margin-bottom: 8px;
}

.tips-list {
  font-size: 12px;
  color: #646566;
  line-height: 1.6;
}

.required {
  color: #ee0a24;
}

.step-header {
  margin-bottom: 16px;
}

.step-title {
  font-size: 16px;
  font-weight: 500;
  color: #323233;
}

.step-desc {
  font-size: 12px;
  color: #969799;
  margin-top: 4px;
}

.header-preview {
  max-height: 400px;
  overflow-y: auto;
}

.header-row-content {
  font-size: 13px;
}

.row-label {
  color: #646566;
  word-break: break-all;
}

.mapping-list {
  max-height: 400px;
  overflow-y: auto;
}

.required-hint {
  color: #ee0a24;
}

.mapped {
  color: #07c160;
}

.unmapped {
  color: #ee0a24;
}

.preview-list {
  max-height: 400px;
  overflow-y: auto;
}

.preview-price {
  color: #ee0a24;
  font-size: 13px;
}

.import-actions {
  padding: 16px;
  background: #fff;
  display: flex;
  gap: 12px;
}

.import-actions .van-button {
  flex: 1;
}
</style>
