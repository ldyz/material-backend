<template>
  <div class="create-outbound-page">
    <van-nav-bar
      title="新建出库单"
      left-arrow
      @click-left="onClickLeft"
    />

    <van-form @submit="onSubmit">
      <!-- 基本信息 -->
      <van-cell-group inset title="基本信息">
        <van-field
          v-model="formData.project_name"
          name="project_id"
          label="项目"
          placeholder="选择项目"
          readonly
          is-link
          @click="showProjectPicker = true"
          :rules="[{ required: true, message: '请选择项目' }]"
        />

        <van-field
          v-model="formData.applicant"
          name="applicant"
          label="申请人"
          placeholder="请输入申请人姓名"
          :rules="[{ required: true, message: '请输入申请人姓名' }]"
        />

    

        <van-field
          v-model="formData.purpose"
          name="purpose"
          label="用途"
          placeholder="请输入用途说明"
          :rules="[{ required: true, message: '请输入用途' }]"
        />

        <van-field
          v-model="formData.remark"
          name="remark"
          label="备注"
          type="textarea"
          placeholder="请输入备注信息"
          rows="3"
          maxlength="500"
          show-word-limit
        />

        <van-field
          name="urgent"
          label="紧急程度"
        >
          <template #input>
            <van-checkbox v-model="formData.urgent">紧急</van-checkbox>
          </template>
        </van-field>
      </van-cell-group>

      <!-- 物资明细 -->
      <van-cell-group inset title="物资明细">
        <van-cell
          title="选择物资"
          is-link
          @click="goToSelectStock"
        >
          <template #value>
            <span v-if="selectedItems.length === 0" style="color: #c8c9cc;">
              点击选择物资
            </span>
            <span v-else style="color: #323233;">
              已选择 {{ selectedItems.length }} 项
            </span>
          </template>
        </van-cell>

        <!-- 已选择的物资列表 -->
        <div
          v-for="(item, index) in selectedItems"
          :key="item.id"
          class="selected-item"
        >
          <div class="item-header">
            <span class="item-name">{{ item.name }}</span>
            <van-button
              type="danger"
              size="mini"
              @click="removeItem(index)"
            >
              删除
            </van-button>
          </div>
          <div class="item-body">
            <div class="item-info">
              <span class="info-label">规格:</span>
              <span class="info-value">{{ item.spec || '-' }}</span>
            </div>
            <div class="item-info">
              <span class="info-label">申请数量:</span>
              <span class="info-value">{{ item.quantity }} {{ item.unit }}</span>
            </div>
            <div v-if="item.remark" class="item-info">
              <span class="info-label">备注:</span>
              <span class="info-value">{{ item.remark }}</span>
            </div>
          </div>
        </div>

        <van-cell
          v-if="selectedItems.length === 0"
          center
          style="color: #c8c9cc; padding: 20px;"
        >
          请点击上方"选择物资"按钮添加物资明细
        </van-cell>
      </van-cell-group>

      <!-- 提交按钮 -->
      <div class="submit-section">
        <van-button
          round
          block
          type="primary"
          native-type="submit"
          :loading="submitting"
        >
          提交出库单
        </van-button>
      </div>
    </van-form>

    <!-- 项目选择器 -->
    <van-popup
      v-model:show="showProjectPicker"
      position="bottom"
      round
    >
      <van-picker
        :columns="projectColumns"
        @confirm="onProjectConfirm"
        @cancel="showProjectPicker = false"
      />
    </van-popup>

    <!-- 库存选择页面 -->
    <van-popup
      v-model:show="showStockSelector"
      position="right"
      :style="{ width: '100%', height: '100%' }"
      round
    >
      <StockSelector
        :project-id="formData.project_id"
        :multiple="false"
        :instant-input="true"
        @confirm="onStockConfirm"
        @cancel="showStockSelector = false"
      />
    </van-popup>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { showToast, showDialog } from 'vant'
import { usePermission } from '@/composables/usePermission'
import { createRequisition } from '@/api/requisition'
import { getProjects } from '@/api/project'
import StockSelector from '@/components/StockSelector.vue'

const router = useRouter()
const { canCreateRequisition } = usePermission()

const showProjectPicker = ref(false)
const showStockSelector = ref(false)
const submitting = ref(false)
const projectColumns = ref([])

const formData = ref({
  project_id: '',
  project_name: '',
  applicant: '',
  department: '',
  purpose: '',
  remark: '',
  urgent: false,
})

const selectedItems = ref([])

// 加载项目列表
async function loadProjects() {
  try {
    // 适配统一响应格式
    const { data } = await getProjects()
    const projects = data || []

    projectColumns.value = [
      { text: '请选择项目', value: '' },
      ...projects.map(p => ({
        text: p.name,
        value: p.id
      }))
    ]
  } catch (error) {
    showToast({
      type: 'fail',
      message: '加载项目列表失败'
    })
  }
}

// 项目确认
function onProjectConfirm({ selectedOptions }) {
  formData.value.project_id = selectedOptions[0].value
  formData.value.project_name = selectedOptions[0].text
  showProjectPicker.value = false

  // 清空已选择的物资（因为项目变了）
  selectedItems.value = []
}

// 去选择库存
function goToSelectStock() {
  if (!formData.value.project_id) {
    showToast('请先选择项目')
    return
  }
  showStockSelector.value = true
}

// 库存确认
function onStockConfirm(stocks) {
  showStockSelector.value = false

  // 添加新选择的库存（已经包含数量）
  let addedCount = 0
  stocks.forEach(stock => {
    const exists = selectedItems.value.find(item => item.id === stock.id)
    if (!exists) {
      selectedItems.value.push({
        id: stock.id,
        name: stock.material_name || stock.name,
        spec: stock.specification || stock.spec,
        unit: stock.unit,
        stock_quantity: stock.quantity,
        quantity: stock.quantity || '',
        remark: stock.remark || '',
        material: stock.material_code || stock.material || stock.code
      })
      addedCount++
    } else {
      // 如果已存在，更新数量和备注
      exists.quantity = stock.quantity
      exists.remark = stock.remark
    }
  })

  if (addedCount > 0) {
    showToast({
      type: 'success',
      message: `成功添加 ${addedCount} 项物资`
    })
  }
}

// 移除物资
function removeItem(index) {
  selectedItems.value.splice(index, 1)
}

// 提交表单
async function onSubmit() {
  if (selectedItems.value.length === 0) {
    showToast('请至少添加一项物资')
    return
  }

  // 验证物资明细
  for (let i = 0; i < selectedItems.value.length; i++) {
    const item = selectedItems.value[i]
    if (!item.quantity || item.quantity <= 0) {
      showToast(`第 ${i + 1} 项物资的数量必须大于0`)
      return
    }
    if (item.quantity > item.stock_quantity) {
      showToast(`第 ${i + 1} 项物资的申请数量不能超过库存 ${item.stock_quantity}`)
      return
    }
  }

  await showDialog({
    title: '确认提交',
    message: `确认提交出库单？\n共 ${selectedItems.value.length} 项物资`,
    showCancelButton: true
  })

  submitting.value = true

  try {
    const data = {
      project_id: parseInt(formData.value.project_id),
      applicant: formData.value.applicant,
      department: formData.value.department,
      purpose: formData.value.purpose,
      remark: formData.value.remark,
      urgent: formData.value.urgent,
      items: selectedItems.value.map(item => ({
        stock_id: item.id,
        name: item.name,
        spec: item.spec,
        unit: item.unit,
        project_id: parseInt(formData.value.project_id),
        material: item.material,
        quantity: parseFloat(item.quantity),
        remark: item.remark
      }))
    }

    await createRequisition(data)

    showToast({
      type: 'success',
      message: '出库单创建成功'
    })

    // 延迟跳转并刷新列表
    setTimeout(() => {
      router.replace('/outbound?refresh=' + Date.now())
    }, 1000)
  } catch (error) {
    showToast({
      type: 'fail',
      message: error.message || '创建失败'
    })
  } finally {
    submitting.value = false
  }
}

function onClickLeft() {
  router.back()
}

// 初始化
loadProjects()
</script>

<style scoped>
.create-outbound-page {
  min-height: 100vh;
  background-color: #f7f8fa;
  padding-bottom: 140px;
}

.submit-section {
  padding: 16px;
  background: white;
  box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.08);
  position: fixed;
  bottom: 50px;
  left: 0;
  right: 0;
  z-index: 999;
}

.selected-item {
  padding: 12px 16px;
  background: #f7f8fa;
  margin-bottom: 8px;
  border-radius: 8px;
}

.item-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.item-name {
  font-weight: bold;
  font-size: 15px;
}

.item-body {
  padding: 0 8px;
}

.item-info {
  display: flex;
  margin-bottom: 4px;
  font-size: 14px;
}

.item-info:last-child {
  margin-bottom: 0;
}

.info-label {
  color: #969799;
  margin-right: 8px;
  min-width: 60px;
}

.info-value {
  color: #323233;
  flex: 1;
}
</style>
