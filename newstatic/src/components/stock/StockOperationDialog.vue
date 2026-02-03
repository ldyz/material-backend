<template>
  <Dialog
    v-model="visible"
    :title="dialogTitle"
    width="580px"
    :loading="loading"
    @confirm="handleSubmit"
    @cancel="handleCancel"
  >
    <div class="stock-dialog-content">
      

      <!-- 物资选择和信息展示 -->
      <div class="material-section">
       

        <!-- 物资详细信息卡片 -->
        <div v-if="selectedMaterial" class="material-card">
          <div class="material-card-header">
            <div class="material-title">
              <span class="material-name">{{ selectedMaterial.name }}</span>
              <el-tag size="small" type="info">{{ selectedMaterial.category || '-' }}</el-tag>
            </div>
          </div>
          <div class="material-card-body">
            <div class="info-grid">
              <div class="info-item">
                <span class="info-label">编码</span>
                <span class="info-value">{{ selectedMaterial.code || '-' }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">规格</span>
                <span class="info-value">{{ selectedMaterial.specification || '-' }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">单位</span>
                <span class="info-value">{{ selectedMaterial.unit || '-' }}</span>
              </div>
              <div class="info-item" v-if="selectedMaterial.price">
                <span class="info-label">单价</span>
                <span class="info-value price">¥{{ formatPrice(selectedMaterial.price) }}</span>
              </div>
            </div>
            <div class="stock-highlight">
              <span class="stock-label">当前库存</span>
              <span class="stock-value">{{ selectedMaterial.quantity || 0 }} {{ selectedMaterial.unit }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 操作表单 -->
      <div class="operation-form">
        <div class="form-row">
          <div class="form-item-wrapper">
            <label class="form-label">
              <span>数量</span>
              <span v-if="formData.type === 'out' && selectedMaterial" class="form-hint">
                最大可出库: {{ selectedMaterial.quantity }}
              </span>
            </label>
            <el-input-number
              v-model="formData.quantity"
              :min="1"
              :max="formData.type === 'out' ? selectedMaterial?.quantity : undefined"
              :step="1"
              :precision="0"
              :disabled="!formData.material_id"
              class="full-width-input"
              placeholder="请输入数量"
            />
          </div>
          
        </div>

        <!-- 金额预览 -->
        <div v-if="showAmountPreview" class="amount-bar">
          <div class="amount-info">
            <span class="amount-label">预计金额</span>
            <span class="amount-value">¥{{ formatPrice(calculateAmount()) }}</span>
          </div>
          <div class="amount-detail">
            <span class="detail-text">{{ formData.quantity }} × ¥{{ formData.price || selectedMaterial?.price || 0 }}</span>
          </div>
        </div>

      
      </div>
    </div>
  </Dialog>
</template>

<script setup>
import { ref, reactive, computed, watch } from 'vue'
import { CirclePlus, Minus, Edit } from '@element-plus/icons-vue'
import Dialog from '@/components/common/Dialog.vue'
import { stockApi, materialApi } from '@/api'
import { ElMessage } from 'element-plus'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  materialId: {
    type: Number,
    default: null
  },
  operationType: {
    type: String,
    default: 'in'
  },
  stockData: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['update:modelValue', 'success'])

const visible = ref(false)
const loading = ref(false)

// 操作类型配置
const operationTypes = [
  { value: 'in', label: '入库', icon: CirclePlus, color: '#67c23a' },
  { value: 'out', label: '出库', icon: Minus, color: '#e6a23c' },
  { value: 'adjust', label: '调整', icon: Edit, color: '#409eff' }
]

// 对话框标题
const dialogTitle = computed(() => {
  const type = operationTypes.find(t => t.value === formData.type)
  return type ? `${type.label}操作` : '库存操作'
})

// 表单数据
const formData = reactive({
  material_id: null,
  quantity: null,
  price: null,
  type: 'in',
  remark: ''
})

// 物资选项
const materialOptions = ref([])

// 选中的物资
const selectedMaterial = computed(() => {
  const material = materialOptions.value.find(m => m.id === formData.material_id)

  if (props.stockData) {
    return {
      id: props.stockData.material_id,
      code: props.stockData.material_code,
      name: props.stockData.material_name,
      category: props.stockData.category,
      specification: props.stockData.specification,
      unit: props.stockData.unit,
      quantity: props.stockData.quantity,
      safety_stock: props.stockData.safety_stock,
      price: material?.price,
      ...material
    }
  }

  return material
})

// 显示金额预览
const showAmountPreview = computed(() => {
  return formData.quantity && (formData.price || selectedMaterial.value?.price)
})

// 格式化价格
const formatPrice = (value) => {
  return Number(value || 0).toFixed(2)
}

// 计算金额
const calculateAmount = () => {
  const price = formData.price || selectedMaterial.value?.price || 0
  const quantity = formData.quantity || 0
  return price * quantity
}

// 获取物资选项
const fetchMaterialOptions = async () => {
  try {
    const { data } = await materialApi.getList({ pageSize: 1000 })
    materialOptions.value = data || []
  } catch (error) {
    console.error('获取物资列表失败:', error)
  }
}

// 处理类型选择
const handleTypeSelect = (type) => {
  formData.type = type
  formData.quantity = null
}

// 物资选择变化
const handleMaterialChange = () => {
  const material = selectedMaterial.value
  if (material) {
    formData.price = material.price
  }
}

// 重置表单
const resetForm = () => {
  Object.assign(formData, {
    material_id: null,
    quantity: null,
    price: null,
    type: props.operationType,
    remark: ''
  })
}

// 验证表单
const validateForm = () => {
  if (!formData.material_id) {
    ElMessage.warning('请选择物资')
    return false
  }
  if (!formData.quantity || formData.quantity <= 0) {
    ElMessage.warning('请输入有效的数量')
    return false
  }
  if (formData.type === 'out' && selectedMaterial.value) {
    if (formData.quantity > selectedMaterial.value.quantity) {
      ElMessage.warning(`出库数量不能超过当前库存 ${selectedMaterial.value.quantity}`)
      return false
    }
  }
  return true
}

// 提交表单
const handleSubmit = async () => {
  if (!validateForm()) return

  // 获取 stock ID
  const stockId = props.stockData?.id
  if (!stockId) {
    ElMessage.error('无法获取库存记录ID')
    return
  }

  try {
    loading.value = true

    // API 只需要 quantity 和 remark
    const data = {
      quantity: formData.quantity,
      remark: formData.remark
    }

    if (formData.type === 'in' || formData.type === 'adjust') {
      await stockApi.in(stockId, data)
      ElMessage.success('入库成功')
    } else {
      await stockApi.out(stockId, data)
      ElMessage.success('出库成功')
    }

    emit('success')
    handleCancel()
  } catch (error) {
    console.error('提交失败:', error)
    ElMessage.error(error.message || '操作失败')
  } finally {
    loading.value = false
  }
}

// 取消操作
const handleCancel = () => {
  visible.value = false
  emit('update:modelValue', false)
}

// 初始化对话框
const initDialog = async () => {
  resetForm()
  await fetchMaterialOptions()

  if (props.materialId) {
    formData.material_id = props.materialId
    formData.type = props.operationType
    handleMaterialChange()
  }
}

// 监听 modelValue 变化
watch(() => props.modelValue, (val) => {
  visible.value = val
  if (val) {
    initDialog()
  }
})

// 监听 visible 变化
watch(visible, (val) => {
  if (!val) {
    emit('update:modelValue', false)
  }
})
</script>

<style scoped>
.stock-dialog-content {
  padding: 0;
}

/* 操作类型选择器 */
.operation-type-selector {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
  padding: 4px;
  background: #f5f7fa;
  border-radius: 8px;
}

.type-option {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 12px 16px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s;
  background: white;
  color: #606266;
  font-weight: 500;
}

.type-option:hover {
  transform: translateY(-2px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.type-option.active {
  background: white;
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.2);
}

.type-option.active:nth-child(1) {
  color: #67c23a;
  border: 2px solid #67c23a;
}

.type-option.active:nth-child(2) {
  color: #e6a23c;
  border: 2px solid #e6a23c;
}

.type-option.active:nth-child(3) {
  color: #409eff;
  border: 2px solid #409eff;
}

.type-icon {
  font-size: 16px;
}

.type-label {
  font-size: 14px;
}

/* 物资区域 */
.material-section {
  margin-bottom: 24px;
}

.section-label {
  font-size: 13px;
  color: #909399;
  margin-bottom: 8px;
  font-weight: 500;
}

.material-select {
  width: 100%;
  margin-bottom: 16px;
}

.option-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.option-name {
  flex: 1;
}

.option-category {
  color: #909399;
  font-size: 12px;
}

/* 物资卡片 */
.material-card {
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  overflow: hidden;
}

.material-card-header {
  background: linear-gradient(135deg, #f5f7fa 0%, #e8eef5 100%);
  padding: 12px 16px;
  border-bottom: 1px solid #e4e7ed;
}

.material-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.material-name {
  font-size: 15px;
  font-weight: 600;
  color: #303133;
}

.material-card-body {
  padding: 16px;
  background: white;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  margin-bottom: 12px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-label {
  font-size: 12px;
  color: #909399;
}

.info-value {
  font-size: 14px;
  color: #303133;
  font-weight: 500;
}

.info-value.price {
  color: #67c23a;
  font-weight: 600;
}

.stock-highlight {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 12px;
  background: linear-gradient(135deg, rgba(64, 158, 255, 0.08) 0%, rgba(64, 158, 255, 0.04) 100%);
  border-radius: 6px;
  border: 1px solid rgba(64, 158, 255, 0.2);
}

.stock-label {
  font-size: 13px;
  color: #606266;
}

.stock-value {
  font-size: 16px;
  font-weight: 600;
  color: #409eff;
}

/* 操作表单 */
.operation-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.form-item-wrapper {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-item-wrapper.full-width {
  grid-column: 1 / -1;
}

.form-label {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 13px;
  color: #606266;
  font-weight: 500;
}

.form-hint {
  font-size: 12px;
  color: #909399;
  font-weight: 400;
}

.optional-tag {
  font-size: 11px;
  color: #909399;
  background: #f5f7fa;
  padding: 2px 6px;
  border-radius: 4px;
}

.full-width-input {
  width: 100%;
}

:deep(.full-width-input .el-input__wrapper) {
  width: 100%;
}

/* 金额预览条 */
.amount-bar {
  background: linear-gradient(135deg, #fff9f0 0%, #fff4e6 100%);
  border-radius: 8px;
  padding: 12px 16px;
  border: 1px solid rgba(230, 162, 60, 0.3);
}

.amount-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.amount-label {
  font-size: 13px;
  color: #606266;
}

.amount-value {
  font-size: 20px;
  font-weight: 600;
  color: #e6a23c;
}

.amount-detail {
  text-align: right;
}

.detail-text {
  font-size: 12px;
  color: #909399;
}

/* 文本域样式 */
:deep(.el-textarea__inner) {
  border-radius: 6px;
}

/* 下拉选择样式 */
:deep(.el-select__wrapper) {
  border-radius: 6px;
}

/* 响应式调整 */
@media (max-width: 480px) {
  .form-row {
    grid-template-columns: 1fr;
  }

  .type-label {
    font-size: 13px;
  }
}
</style>
