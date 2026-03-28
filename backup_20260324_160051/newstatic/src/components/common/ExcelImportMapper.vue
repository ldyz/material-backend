/**
 * Excel 导入映射组件
 *
 * 功能说明：
 * - 支持上传任意 Excel 文件
 * - 自动读取表头列
 * - 智能匹配 Excel 列和数据库字段
 * - 支持手动调整映射关系
 * - 提供数据预览功能
 *
 * 注意：这是一个纯前端组件，不涉及后端 API 调用
 * - importApi 只是数据处理函数，由父组件提供
 * - 导入的数据通过 success 事件返回给父组件
 * - 父组件决定如何处理这些数据（添加到表单、调用 API 等）
 *
 * @module ExcelImportMapper
 * @author Material Management System
 * @date 2025-01-27
 */
<template>
  <div class="excel-import-mapper">
    <!-- 步骤指示器 -->
    <el-steps :active="currentStep" finish-status="success" class="steps">
      <el-step v-if="showProjectStep" title="选择项目" />
      <el-step :title="showProjectStep ? '上传文件' : '选择文件'" />
      <el-step :title="showProjectStep ? '选择标题行' : '选择表头'" />
      <el-step title="映射字段" />
      <el-step title="预览数据" />
    </el-steps>

    <!-- 步骤1: 选择项目 (仅在有showProjectStep时显示) -->
    <div v-if="showProjectStep && currentStep === 0" class="step-content">
      <el-alert
        title="选择项目"
        type="info"
        :closable="false"
        class="mb-20"
      >
        请选择要导入物资的所属项目，导入的物资将关联到该项目。支持选择主项目或子项目。
      </el-alert>

      <div class="project-selector-container">
        <el-form label-width="120px">
          <el-form-item label="所属项目" required>
            <ProjectSelector
              v-model="selectedProjectId"
              :projects="projects"
              placeholder="选择项目（支持层级显示）"
            />
          </el-form-item>
        </el-form>
      </div>

      <div class="step-actions">
        <el-button type="primary" :disabled="!selectedProjectId" @click="currentStep = 1">
          下一步
        </el-button>
      </div>
    </div>

    <!-- 步骤2: 上传文件 -->
    <div v-show="(!showProjectStep && currentStep === 0) || (showProjectStep && currentStep === 1)" class="step-content">
      <el-upload
        ref="uploadRef"
        :auto-upload="false"
        :limit="1"
        accept=".xlsx,.xls"
        drag
        :before-upload="handleBeforeUpload"
        :on-change="handleFileChange"
        :on-remove="handleFileRemove"
        class="upload-area"
      >
        <el-icon :size="80" class="upload-icon"><UploadFilled /></el-icon>
        <div class="el-upload__text">
          将 Excel 文件拖到此处，或<em>点击上传</em>
        </div>
        <template #tip>
          <div class="el-upload__tip">
            支持 .xlsx/.xls 格式，文件大小不超过 10MB
          </div>
        </template>
      </el-upload>

      <div class="step-actions">
        <el-button v-if="showProjectStep" @click="currentStep--">上一步</el-button>
        <el-button type="primary" :disabled="!uploadedFile" @click="handleUploadNext">
          下一步
        </el-button>
      </div>
    </div>

    <!-- 步骤3: 选择标题行 -->
    <div v-show="(!showProjectStep && currentStep === 1) || (showProjectStep && currentStep === 2)" class="step-content">
      <el-alert
        title="选择标题行"
        type="info"
        :closable="false"
        class="mb-20"
      >
        请点击下表中的某一行作为标题行（列标题），选中的行将高亮显示
      </el-alert>

      <!-- 当前选择的标题行提示 -->
      <div class="current-header-info">
        <el-tag type="success" size="large">
          当前选择的标题行：第 {{ headerRowIndex }} 行
        </el-tag>
        <span class="header-content-preview">
          内容：{{ excelData[headerRowIndex - 1]?.join(' | ') || '无' }}
        </span>
      </div>

      <!-- 预览Excel内容 -->
      <div class="excel-preview-container">
        <el-table
          :data="previewExcelData"
          border
          stripe
          max-height="400"
          class="preview-table"
          :row-class-name="getRowClassName"
          @row-click="handleRowClick"
        >
          <el-table-column label="行号" width="80" align="center">
            <template #default="scope">
              <el-tag v-if="scope.row._isSelected" type="success">标题行</el-tag>
              <span v-else class="row-number">{{ scope.$index + 1 }}</span>
            </template>
          </el-table-column>
          <el-table-column
            v-for="col in maxColumns"
            :key="col"
            :prop="String(col - 1)"
            :label="`列${col}`"
            min-width="120"
            show-overflow-tooltip
          >
            <template #default="scope">
              <span
                :class="{
                  'header-row': scope.row._isSelected
                }"
              >
                {{ scope.row[col - 1] || '' }}
              </span>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <div class="step-actions">
        <el-button @click="currentStep--">上一步</el-button>
        <el-button type="primary" @click="handleHeaderRowNext">
          下一步
        </el-button>
      </div>
    </div>

    <!-- 步骤4: 映射字段 -->
    <div v-show="(!showProjectStep && currentStep === 2) || (showProjectStep && currentStep === 3)" class="step-content">
      <el-alert
        title="字段映射"
        type="info"
        :closable="false"
        class="mb-20"
      >
        请将 Excel 列映射到对应的数据库字段，标有 * 的为必填字段
      </el-alert>

      <div class="mapping-table-container">
        <el-table
          :data="mappingFields"
          border
          stripe
          max-height="400"
        >
          <el-table-column label="数据库字段" width="180">
            <template #default="scope">
              <span v-if="scope.row.required" class="required">*</span>
              {{ scope.row.label }}
              <div class="field-hint">{{ scope.row.hint }}</div>
            </template>
          </el-table-column>

          <el-table-column label="Excel 列" min-width="200">
            <template #default="scope">
              <el-select
                v-model="scope.row.mappedColumn"
                placeholder="请选择 Excel 列"
                filterable
                clearable
                @change="handleMappingChange(scope.row)"
              >
                <el-option
                  v-for="col in excelColumns"
                  :key="col"
                  :label="col"
                  :value="col"
                >
                  <span>{{ col }}</span>
                  <span v-if="isAutoMatched(col, scope.row)" class="auto-match-tag">
                    自动匹配
                  </span>
                </el-option>
              </el-select>
            </template>
          </el-table-column>

          <el-table-column label="数据预览" min-width="150">
            <template #default="scope">
              <div v-if="scope.row.mappedColumn" class="preview-cell">
                {{ getPreviewData(scope.row.mappedColumn) }}
              </div>
              <span v-else class="text-placeholder">-</span>
            </template>
          </el-table-column>

          <el-table-column label="匹配度" width="100">
            <template #default="scope">
              <el-tag
                v-if="scope.row.matchScore"
                :type="getMatchScoreType(scope.row.matchScore)"
                size="small"
              >
                {{ scope.row.matchScore }}%
              </el-tag>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <div class="step-actions">
        <el-button @click="currentStep--">上一步</el-button>
        <el-button type="primary" @click="handleMappingNext">下一步</el-button>
      </div>
    </div>

    <!-- 步骤5: 预览数据 -->
    <div v-show="(!showProjectStep && currentStep === 3) || (showProjectStep && currentStep === 4)" class="step-content">
      <el-alert
        title="数据预览"
        type="info"
        :closable="false"
        class="mb-20"
      >
        <template #default>
          <div>将导入 <strong>{{ selectedRowIndexes.length }}</strong> 条数据（共 {{ previewData.length }} 条），请仔细核对。</div>
          <div style="margin-top: 8px; color: #909399;">
            <el-icon style="vertical-align: middle"><InfoFilled /></el-icon>
            默认已全选所有数据，您可以取消勾选不需要导入的行。切换分页时选中状态会自动保持。
          </div>
        </template>
      </el-alert>

      <el-table
        ref="previewTableRef"
        :data="pagedPreviewData"
        row-key="_rowKey"
        border
        stripe
        max-height="400"
        class="preview-table"
        @selection-change="handlePreviewSelectionChange"
      >
        <el-table-column
          type="selection"
          width="55"
        />
        <el-table-column
          type="index"
          label="序号"
          width="60"
          align="center"
        >
          <template #default="scope">
            {{ (previewPagination.currentPage - 1) * previewPagination.pageSize + scope.$index + 1 }}
          </template>
        </el-table-column>
        <el-table-column
          v-for="field in mappedFields"
          :key="field.key"
          :prop="field.key"
          :label="field.label"
          :min-width="field.width || 120"
          show-overflow-tooltip
        />
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="previewPagination.currentPage"
          v-model:page-size="previewPagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="previewData.length"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handlePreviewPageSizeChange"
          @current-change="handlePreviewPageChange"
        />
      </div>

      <div class="step-actions">
        <el-button @click="currentStep--">上一步</el-button>
        <el-button type="primary" :loading="importing" @click="handleImport">
          开始导入 ({{ selectedRowIndexes.length }} 条)
        </el-button>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { UploadFilled, InfoFilled } from '@element-plus/icons-vue'
import * as XLSX from 'xlsx'
import ProjectSelector from './ProjectSelector.vue'

/**
 * 组件 Props
 */
const props = defineProps({
  /**
   * 字段定义
   */
  fields: {
    type: Array,
    required: true
  },
  /**
   * 数据处理函数（由父组件提供）
   * 注意：这不是后端 API，只是一个数据处理函数
   * 接收格式: { items: [...] }
   * 返回格式: { success: true, data: { items: [...] } }
   */
  importApi: {
    type: Function,
    required: true
  },
  /**
   * 项目列表（仅在选择项目步骤时需要）
   */
  projects: {
    type: Array,
    default: () => []
  },
  /**
   * 是否显示选择项目步骤
   */
  showProjectStep: {
    type: Boolean,
    default: true
  }
})

/**
 * 组件 Emits
 */
const emit = defineEmits(['success', 'close'])

// ========== 状态管理 ==========

/**
 * 当前步骤
 * 0: 选择项目 (仅在 showProjectStep=true 时)
 * 1: 上传文件
 * 2: 选择标题行
 * 3: 映射字段
 * 4: 预览数据
 */
const currentStep = ref(0)

/**
 * 选择的项目ID
 */
const selectedProjectId = ref(null)

/**
 * 上传的文件
 */
const uploadedFile = ref(null)

/**
 * 标题行索引（从1开始）
 */
const headerRowIndex = ref(1)

/**
 * Excel 工作簿数据
 */
const workbook = ref(null)

/**
 * Excel 工作表名称
 */
const sheetName = ref('')

/**
 * Excel 原始数据（包含表头）
 */
const excelData = ref([])

/**
 * Excel 表头列
 */
const excelColumns = ref([])

/**
 * Excel 数据行（不含表头）
 */
const excelRows = ref([])

/**
 * 字段映射配置
 */
const mappingFields = ref([])

/**
 * 导入中状态
 */
const importing = ref(false)

/**
 * 导入结果
 */
const importResult = ref({
  success: false,
  title: '',
  successCount: 0,
  failCount: 0,
  errors: []
})

/**
 * 预览表格引用
 */
const previewTableRef = ref(null)

/**
 * 选中的预览行
 */
const selectedPreviewRows = ref([])

/**
 * 选中行的原始索引
 */
const selectedRowIndexes = ref([])

/**
 * 预览分页配置
 */
const previewPagination = ref({
  currentPage: 1,
  pageSize: 10
})

// ========== 计算属性 ==========

/**
 * Excel最大列数
 */
const maxColumns = computed(() => {
  if (excelData.value.length === 0) return 0
  const maxPreviewRows = Math.min(excelData.value.length, 10)
  return Math.max(...excelData.value.slice(0, maxPreviewRows).map(row => row.length), 0)
})

/**
 * 预览Excel数据（用于选择标题行）
 */
const previewExcelData = computed(() => {
  const maxPreviewRows = Math.min(excelData.value.length, 10)
  const preview = []

  for (let i = 0; i < maxPreviewRows; i++) {
    const row = excelData.value[i] || []
    const rowObj = {
      _isHeader: i === headerRowIndex.value - 1,
      _isSelected: i === headerRowIndex.value - 1,
      _rowIndex: i  // 添加原始行索引
    }

    // 将数组转换为对象，键为列索引
    for (let j = 0; j < maxColumns.value; j++) {
      rowObj[j] = row[j] || ''
    }

    preview.push(rowObj)
  }

  return preview
})

/**
 * 已映射的字段列表
 */
const mappedFields = computed(() => {
  return mappingFields.value
    .filter(field => field.mappedColumn)
    .map(field => ({
      key: field.field,
      label: field.label,
      width: field.field === 'remark' ? 200 : 120
    }))
})

/**
 * 预览数据
 */
const previewData = computed(() => {
  return excelRows.value.map(row => {
    const obj = {}
    mappingFields.value.forEach(field => {
      if (field.mappedColumn) {
        const colIndex = excelColumns.value.indexOf(field.mappedColumn)
        obj[field.field] = row[colIndex]
      }
    })
    return obj
  })
})

/**
 * 分页后的预览数据
 */
const pagedPreviewData = computed(() => {
  const start = (previewPagination.value.currentPage - 1) * previewPagination.value.pageSize
  const end = start + previewPagination.value.pageSize
  return previewData.value.slice(start, end).map((item, index) => ({
    ...item,
    _rowKey: `row_${start + index}`  // 添加唯一标识符
  }))
})

// ========== 方法定义 ==========

/**
 * 读取文件头（魔数）验证文件类型
 */
const validateFileSignature = (file) => {
  return new Promise((resolve) => {
    const reader = new FileReader()

    reader.onload = (e) => {
      const arr = new Uint8Array(e.target.result)
      let isValid = false

      // 检查文件头签名
      // XLS 文件: D0 CF 11 E0 A1 B1 1A E1 (OLE Compound Document)
      if (arr[0] === 0xD0 && arr[1] === 0xCF && arr[2] === 0x11 && arr[3] === 0xE0) {
        isValid = true
      }
      // XLSX 文件: 50 4B 03 04 (ZIP 文件，xlsx 实际上是 zip 格式)
      else if (arr[0] === 0x50 && arr[1] === 0x4B && arr[2] === 0x03 && arr[3] === 0x04) {
        isValid = true
      }

      resolve(isValid)
    }

    reader.onerror = () => resolve(false)

    // 只读取前 8 字节用于验证
    reader.readAsArrayBuffer(file.slice(0, 8))
  })
}

/**
 * 验证文件类型
 */
const validateFileType = async (file) => {
  // 检查文件扩展名
  const fileName = file.name || ''
  const lastDotIndex = fileName.lastIndexOf('.')

  // 如果没有扩展名，直接拒绝
  if (lastDotIndex === -1) {
    ElMessage.error('文件没有扩展名，请选择 .xlsx 或 .xls 文件')
    return false
  }

  const fileExt = fileName.substring(lastDotIndex).toLowerCase()
  const validExtensions = ['.xlsx', '.xls']

  if (!validExtensions.includes(fileExt)) {
    ElMessage.error(`只支持 Excel 文件格式（.xlsx/.xls），当前文件：${fileName}`)
    return false
  }

  // 检查文件大小（10MB）
  const maxSize = 10 * 1024 * 1024
  if (file.size > maxSize) {
    ElMessage.error('文件大小不能超过 10MB')
    return false
  }

  // 验证文件头签名（防止用户修改扩展名伪装文件）
  const isValidSignature = await validateFileSignature(file)
  if (!isValidSignature) {
    ElMessage.error(`文件内容不是有效的 Excel 格式，请检查文件是否损坏`)
    return false
  }

  return true
}

/**
 * 上传前验证
 * 在选择文件时立即触发，返回 false 阻止无效文件
 */
const handleBeforeUpload = async (file) => {
  return await validateFileType(file)
}

/**
 * 文件选择变化事件
 */
const handleFileChange = async (file) => {
  // 再次验证文件类型（双重保险）
  const isValid = await validateFileType(file.raw)
  if (!isValid) {
    // 使用 el-upload 的方法移除文件
    setTimeout(() => {
      uploadRef.value?.handleRemove(file)
    }, 0)
    uploadedFile.value = null
    return
  }

  uploadedFile.value = file
  parseExcelFile(file.raw)
}

/**
 * 文件移除事件
 */
const handleFileRemove = () => {
  uploadedFile.value = null
  workbook.value = null
  sheetName.value = ''
  excelData.value = []
  excelColumns.value = []
  excelRows.value = []
  headerRowIndex.value = 1  // 重置标题行索引
  resetMappingFields()
}

/**
 * 解析 Excel 文件
 */
const parseExcelFile = (file) => {
  const reader = new FileReader()

  reader.onload = (e) => {
    try {
      const data = new Uint8Array(e.target.result)
      const wb = XLSX.read(data, { type: 'array' })

      // 读取第一个工作表
      const firstSheetName = wb.SheetNames[0]
      const worksheet = wb.Sheets[firstSheetName]

      // 转换为 JSON 数据
      const jsonData = XLSX.utils.sheet_to_json(worksheet, {
        header: 1,  // 使用数组格式，包含表头
        defval: ''  // 空单元格默认值
      })

      if (jsonData.length < 2) {
        ElMessage.error('Excel 文件没有数据行')
        return
      }

      // 保存数据
      workbook.value = wb
      sheetName.value = firstSheetName
      excelData.value = jsonData

      // 默认第一行为标题行，初始化列名
      updateColumnsByHeaderRow()

      // 初始化映射字段（但不自动匹配）
      resetMappingFields()

      ElMessage.success(`成功读取 ${jsonData.length - 1} 行数据`)
    } catch (error) {
      console.error('解析 Excel 文件失败:', error)
      ElMessage.error('解析 Excel 文件失败，请检查文件格式')
    }
  }

  reader.readAsArrayBuffer(file)
}

/**
 * 根据标题行更新列名
 */
const updateColumnsByHeaderRow = () => {
  const rowIndex = headerRowIndex.value - 1
  if (excelData.value[rowIndex]) {
    excelColumns.value = excelData.value[rowIndex].map(col => String(col).trim())
    // 更新数据行（标题行之后的所有行）
    excelRows.value = excelData.value.slice(rowIndex + 1)
  }
}

/**
 * 获取行样式类名
 */
const getRowClassName = ({ row, rowIndex }) => {
  if (row._isSelected) {
    return 'selected-header-row'
  }
  return ''
}

/**
 * 表格行点击事件
 */
const handleRowClick = (row, column, event) => {
  // 直接使用行对象中存储的索引
  if (row._rowIndex !== undefined) {
    const newRowIndex = row._rowIndex + 1

    // 如果选择的是同一行，不需要更新
    if (newRowIndex === headerRowIndex.value) {
      return
    }

    headerRowIndex.value = newRowIndex

    // 立即更新列名和数据行
    updateColumnsByHeaderRow()

    // 重置并重新自动匹配字段
    resetMappingFields()
    autoMatchFields()

    ElMessage.success(`已选择第 ${headerRowIndex.value} 行作为标题行`)
  }
}

/**
 * 自动匹配字段
 *
 * 使用相似度算法自动匹配 Excel 列和数据库字段
 */
const autoMatchFields = () => {
  // 初始化映射字段
  mappingFields.value = props.fields.map(field => ({
    ...field,
    mappedColumn: null,
    matchScore: 0
  }))

  // 自动匹配
  mappingFields.value.forEach(field => {
    let bestMatch = ''
    let bestScore = 0

    excelColumns.value.forEach(col => {
      const score = calculateSimilarity(field.field, field.label, col)
      if (score > bestScore && score >= 60) {
        bestScore = score
        bestMatch = col
      }
    })

    if (bestMatch) {
      field.mappedColumn = bestMatch
      field.matchScore = bestScore
    }
  })
}

/**
 * 计算相似度
 *
 * @param {string} fieldKey - 数据库字段键
 * @param {string} fieldLabel - 数据库字段标签
 * @param {string} excelColumn - Excel 列名
 * @returns {number} 相似度分数 (0-100)
 */
const calculateSimilarity = (fieldKey, fieldLabel, excelColumn) => {
  const key = fieldKey.toLowerCase()
  const label = fieldLabel.toLowerCase()
  const col = excelColumn.toLowerCase().trim()

  // 完全匹配
  if (col === key || col === label) {
    return 100
  }

  // 包含匹配
  if (col.includes(key) || key.includes(col) ||
      col.includes(label) || label.includes(col)) {
    return 90
  }

  // 拼音/简写匹配（可以扩展）
  const aliases = {
    'code': ['编码', '编号', '代码', 'code', 'no'],
    'name': ['名称', '品名', 'name', 'title'],
    'category': ['分类', '类别', 'category', 'type'],
    'specification': ['规格', '型号', 'specification', 'spec', 'model'],
    'material': ['材质', '材料', 'material'],
    'unit': ['单位', '计量单位', 'unit'],
    'price': ['单价', '价格', 'price', '金额'],
    'stock': ['库存', '数量', 'stock', 'quantity', 'qty'],
    'remark': ['备注', '说明', 'remark', 'note', '描述']
  }

  if (aliases[key]) {
    for (const alias of aliases[key]) {
      if (col.includes(alias) || alias.includes(col)) {
        return 80
      }
    }
  }

  // 编辑距离算法（Levenshtein Distance）
  const distance = levenshteinDistance(label, col)
  const maxLen = Math.max(label.length, col.length)
  const similarity = ((maxLen - distance) / maxLen) * 100

  return Math.round(similarity)
}

/**
 * 计算编辑距离
 */
const levenshteinDistance = (str1, str2) => {
  const m = str1.length
  const n = str2.length
  const dp = Array(m + 1).fill(null).map(() => Array(n + 1).fill(0))

  for (let i = 0; i <= m; i++) dp[i][0] = i
  for (let j = 0; j <= n; j++) dp[0][j] = j

  for (let i = 1; i <= m; i++) {
    for (let j = 1; j <= n; j++) {
      if (str1[i - 1] === str2[j - 1]) {
        dp[i][j] = dp[i - 1][j - 1]
      } else {
        dp[i][j] = Math.min(
          dp[i - 1][j] + 1,
          dp[i][j - 1] + 1,
          dp[i - 1][j - 1] + 1
        )
      }
    }
  }

  return dp[m][n]
}

/**
 * 判断是否为自动匹配
 */
const isAutoMatched = (excelColumn, field) => {
  return field.mappedColumn === excelColumn && field.matchScore >= 60
}

/**
 * 获取匹配度标签类型
 */
const getMatchScoreType = (score) => {
  if (score >= 90) return 'success'
  if (score >= 70) return 'warning'
  if (score >= 60) return 'info'
  return ''
}

/**
 * 获取预览数据
 */
const getPreviewData = (column) => {
  const colIndex = excelColumns.value.indexOf(column)
  if (colIndex === -1) return '-'
  return excelRows.value[0]?.[colIndex] || '-'
}

/**
 * 映射变化事件
 */
const handleMappingChange = (field) => {
  // 可以在这里添加验证逻辑
}

/**
 * 上传文件后的下一步处理
 */
const handleUploadNext = () => {
  // 跳转到选择标题行步骤
  // showProjectStep=true: 步骤0(选择项目) -> 1(上传文件) -> 2(选择标题行)
  // showProjectStep=false: 步骤0(上传文件) -> 1(选择标题行)
  const nextStep = props.showProjectStep ? 2 : 1
  currentStep.value = nextStep
}

/**
 * 选择标题行后的下一步处理
 */
const handleHeaderRowNext = () => {
  // 跳转到映射字段步骤
  // showProjectStep=true: 步骤2(选择标题行) -> 3(映射字段)
  // showProjectStep=false: 步骤1(选择标题行) -> 2(映射字段)
  const nextStep = props.showProjectStep ? 3 : 2
  currentStep.value = nextStep
}

/**
 * 映射字段后的下一步处理
 */
const handleMappingNext = () => {
  // 验证必填字段
  const requiredFields = mappingFields.value.filter(f => f.required)
  const unmappedRequired = requiredFields.filter(f => !f.mappedColumn)

  if (unmappedRequired.length > 0) {
    ElMessage.warning(`请先映射必填字段: ${unmappedRequired.map(f => f.label).join(', ')}`)
    return
  }

  // 先清空之前的选中状态
  selectedRowIndexes.value = []
  selectedPreviewRows.value = []

  // 跳转到预览数据步骤
  // showProjectStep=true: 步骤3(映射字段) -> 4(预览数据)
  // showProjectStep=false: 步骤2(映射字段) -> 3(预览数据)
  const nextStep = props.showProjectStep ? 4 : 3
  currentStep.value = nextStep
}

/**
 * 处理预览表格选择变化
 */
const handlePreviewSelectionChange = (selection) => {
  // 获取当前页所有数据的索引
  const currentPageData = pagedPreviewData.value
  const currentPageIndexes = currentPageData.map(row => {
    const match = row._rowKey?.match(/row_(\d+)/)
    return match ? parseInt(match[1]) : -1
  }).filter(idx => idx >= 0)

  // 获取当前页被选中的行的索引
  const selectedIndexesInPage = selection.map(row => {
    const match = row._rowKey?.match(/row_(\d+)/)
    return match ? parseInt(match[1]) : -1
  }).filter(idx => idx >= 0)

  // 从 selectedRowIndexes 中移除当前页的所有索引
  selectedRowIndexes.value = selectedRowIndexes.value.filter(idx => !currentPageIndexes.includes(idx))

  // 添加当前页中被选中的索引
  selectedRowIndexes.value.push(...selectedIndexesInPage)

  // 更新选中的行（用于显示）
  selectedPreviewRows.value = selection
}

/**
 * 处理预览分页大小变化
 */
const handlePreviewPageSizeChange = async (pageSize) => {
  previewPagination.value.pageSize = pageSize
  previewPagination.value.currentPage = 1
  // 等待 DOM 更新后同步选中状态（需要多次 nextTick 确保表格完全渲染）
  await nextTick()
  await nextTick()
  syncSelectionWithTable()
}

/**
 * 处理预览页码变化
 */
const handlePreviewPageChange = async (page) => {
  previewPagination.value.currentPage = page
  // 等待 DOM 更新后同步选中状态（需要多次 nextTick 确保表格完全渲染）
  await nextTick()
  await nextTick()
  syncSelectionWithTable()
}

/**
 * 同步选中状态到表格
 * 根据 selectedRowIndexes 自动选中当前页中应该被选中的行
 */
const syncSelectionWithTable = () => {
  if (!previewTableRef.value) {
    console.warn('previewTableRef is not available')
    return
  }

  try {
    // 获取当前页的数据
    const currentPageData = pagedPreviewData.value

    // 先获取当前表格的选中状态
    const currentSelection = previewTableRef.value.getSelectionRows()

    // 遍历当前页的每一行，精确控制选中状态
    currentPageData.forEach(row => {
      // 从 rowKey 中提取原始索引
      const match = row._rowKey?.match(/row_(\d+)/)
      if (match) {
        const rowIndex = parseInt(match[1])
        const shouldBeSelected = selectedRowIndexes.value.includes(rowIndex)
        const isCurrentlySelected = currentSelection.some(s => s._rowKey === row._rowKey)

        // 只有当状态不一致时才切换
        if (shouldBeSelected && !isCurrentlySelected) {
          previewTableRef.value.toggleRowSelection(row, true)
        } else if (!shouldBeSelected && isCurrentlySelected) {
          previewTableRef.value.toggleRowSelection(row, false)
        }
      }
    })

    console.log(`同步选中状态完成: 总选中 ${selectedRowIndexes.value.length} 行, 当前页显示 ${currentPageData.length} 行`)
  } catch (error) {
    console.error('同步选中状态失败:', error)
  }
}

/**
 * 执行数据处理
 *
 * 注意：这不是后端 API 调用，而是调用父组件提供的数据处理函数
 * 处理后的数据通过 success 事件返回给父组件
 * 父组件决定如何处理这些数据（添加到表单、调用 API 等）
 */
const handleImport = async () => {
  // 检查是否选中了行
  if (selectedRowIndexes.value.length === 0) {
    ElMessage.warning('请至少选择一行数据进行导入')
    return
  }

  // 只有在显示项目选择步骤时才需要检查项目
  if (props.showProjectStep && !selectedProjectId.value) {
    ElMessage.error('请先选择项目')
    currentStep.value = 0
    return
  }

  importing.value = true

  try {
    // 只处理选中的行
    const importData = selectedRowIndexes.value.map(rowIndex => {
      const row = excelRows.value[rowIndex]
      const obj = {}
      // 只有在显示项目选择步骤时才添加 project_id
      if (props.showProjectStep) {
        obj.project_id = selectedProjectId.value
      }
      mappingFields.value.forEach(field => {
        if (field.mappedColumn) {
          const colIndex = excelColumns.value.indexOf(field.mappedColumn)
          let value = row[colIndex]

          // 数据类型转换
          if (field.field === 'price' || field.field === 'quantity' || field.field === 'unit_price' || field.field === 'planned_quantity') {
            value = parseFloat(value) || 0
          } else if (value) {
            value = String(value).trim()
          }

          obj[field.field] = value
        }
      })
      return obj
    })

    // 调用父组件提供的数据处理函数
    // 注意：这不是后端 API，只是一个数据处理函数
    // 根据是否显示项目步骤决定传递的参数格式
    const dataParams = props.showProjectStep
      ? { materials: importData }
      : { items: importData }

    // 获取处理结果
    const processResult = await props.importApi(dataParams)

    // 处理成功
    importResult.value = {
      success: processResult?.success !== false,  // 默认成功
      title: processResult?.title || '处理成功',
      successCount: processResult?.data?.success || importData.length,
      failCount: processResult?.data?.failed || 0,
      errors: processResult?.data?.errors || []
    }

    // 将完整的结果传递给父组件，包含实际的数据项
    emit('success', processResult || {
      success: true,
      data: {
        total: importData.length,
        success: importData.length,
        failed: 0,
        items: importData
      }
    })

    // 立即关闭对话框并清空所有状态数据
    emit('close')
    resetState()
  } catch (error) {
    console.error('数据处理失败:', error)

    // 处理失败，显示错误消息
    ElMessage.error(error.message || '数据处理失败，请检查数据格式后重试')

    importResult.value = {
      success: false,
      title: '处理失败',
      successCount: 0,
      failCount: selectedRowIndexes.value.length,
      errors: [
        {
          row: '全部',
          message: error.message || '处理错误'
        }
      ]
    }
  } finally {
    importing.value = false
  }
}

/**
 * 重置映射字段
 */
const resetMappingFields = () => {
  mappingFields.value = props.fields.map(field => ({
    ...field,
    mappedColumn: null,
    matchScore: 0
  }))
}

/**
 * 重置所有状态数据
 * 用于导入完成后清空组件状态，准备下次导入
 */
const resetState = () => {
  // 重置步骤
  currentStep.value = 0

  // 重置项目选择
  selectedProjectId.value = null

  // 重置文件上传
  uploadedFile.value = null

  // 重置标题行
  headerRowIndex.value = 1

  // 清空 Excel 数据
  workbook.value = null
  sheetName.value = ''
  excelData.value = []
  excelColumns.value = []
  excelRows.value = []

  // 重置字段映射
  resetMappingFields()

  // 重置导入结果
  importResult.value = {
    success: false,
    title: '',
    successCount: 0,
    failCount: 0,
    errors: []
  }

  // 重置导入状态
  importing.value = false

  // 重置选择状态
  selectedPreviewRows.value = []
  selectedRowIndexes.value = []

  // 重置分页配置
  previewPagination.value = {
    currentPage: 1,
    pageSize: 10
  }
}

// 监听步骤变化，当进入预览步骤时自动全选
watch(currentStep, async (newStep) => {
  const previewStep = props.showProjectStep ? 4 : 3
  if (newStep === previewStep) {
    // 先设置所有行为选中状态
    selectedRowIndexes.value = previewData.value.map((_, index) => index)
    // 等待表格渲染
    await nextTick()
    // 直接全选当前表格中的所有行
    if (previewTableRef.value) {
      // 使用 toggleAllSelection 全选
      pagedPreviewData.value.forEach(row => {
        previewTableRef.value.toggleRowSelection(row, true)
      })
    }
  }
})

// 初始化
resetMappingFields()
</script>

<style scoped>
.excel-import-mapper {
  padding: 20px 0;
}

.steps {
  margin-bottom: 40px;
}

.step-content {
  min-height: 400px;
}

.project-selector-container {
  margin: 30px 0;
  padding: 30px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.tip-list {
  margin: 10px 0 0 20px;
  font-size: 14px;
  line-height: 1.8;
}

.upload-area {
  margin: 30px 0;
}

.upload-icon {
  color: #409eff;
}

.mapping-table-container {
  margin-bottom: 20px;
}

.required {
  color: #f56c6c;
  margin-right: 4px;
}

.field-hint {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.preview-cell {
  font-size: 13px;
  color: #606266;
}

.text-placeholder {
  color: #c0c4cc;
}

.auto-match-tag {
  float: right;
  font-size: 12px;
  color: #67c23a;
}

.step-actions {
  margin-top: 30px;
  text-align: right;
}

.preview-table {
  margin-bottom: 20px;
}

.mb-20 {
  margin-bottom: 20px;
}

.current-header-info {
  margin: 20px 0;
  padding: 15px 20px;
  background-color: #f0f9ff;
  border: 1px solid #d1e7dd;
  border-radius: 4px;
  display: flex;
  align-items: center;
  gap: 15px;
}

.header-content-preview {
  color: #606266;
  font-size: 14px;
}

.excel-preview-container {
  margin: 20px 0;
}

.preview-table {
  cursor: pointer;
}

.preview-table :deep(.el-table__body-wrapper .el-table__body tr) {
  cursor: pointer;
  transition: background-color 0.2s;
}

.preview-table :deep(.el-table__body-wrapper .el-table__body tr:hover) {
  background-color: #e6f7ff !important;
}

.preview-table :deep(.selected-header-row) {
  background-color: #f0f9ff !important;
}

.preview-table :deep(.selected-header-row:hover) {
  background-color: #d1e7dd !important;
}

.row-number {
  color: #909399;
  font-size: 13px;
}

.header-row {
  font-weight: bold;
  color: #67c23a;
  background-color: #f0f9ff !important;
}

.pagination-container {
  margin: 20px 0;
  display: flex;
  justify-content: center;
}
</style>
