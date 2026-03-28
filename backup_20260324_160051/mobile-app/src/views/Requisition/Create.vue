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
          v-model="formData.notes"
          name="notes"
          type="textarea"
          label="备注"
          placeholder="请输入备注（可选）"
          rows="2"
        />
      </van-cell-group>

      <van-cell-group inset title="库存物资">
        <van-empty v-if="!stockList.length" description="请先选择项目查看库存" />
        <div
          v-for="(stock, index) in stockList"
          :key="stock.id"
          class="stock-item"
        >
          <van-cell
            :title="stock.material_name"
            :label="`库存：${stock.quantity} ${stock.unit || ''}`"
            is-link
            @click="toggleItem(stock)"
          >
            <template #right-icon>
              <van-checkbox
                :model-value="isItemSelected(stock)"
                @click.stop="toggleItem(stock)"
              />
            </template>
          </van-cell>
          <div v-if="isItemSelected(stock)" class="item-quantity">
            <van-field
              v-model.number="getFormItem(stock).requested_quantity"
              type="number"
              label="申请数量"
              placeholder="请输入数量"
              :min="1"
              :max="stock.quantity"
            />
          </div>
        </div>
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

    <!-- 日期选择器 -->
    <van-popup v-model:show="showDatePicker" position="bottom">
      <van-date-picker
        v-model="currentDate"
        @confirm="onDateConfirm"
        @cancel="showDatePicker = false"
      />
    </van-popup>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { showToast, showLoadingToast, closeToast } from 'vant'
import { createRequisition, getProjects, getStock } from '@/api/requisition'

const router = useRouter()

const showProjectPicker = ref(false)
const showDatePicker = ref(false)
const submitting = ref(false)

const projects = ref([])
const selectedProject = ref(null)
const stockList = ref([])

const currentDate = ref(new Date())
const formData = ref({
  project_id: null,
  department_name: '',
  requisition_date: '',
  notes: '',
  items: []
})

const projectColumns = ref([])

async function loadProjects() {
  try {
    const response = await getProjects({ page: 1, page_size: 100 })
    projects.value = response.data || []

    projectColumns.value = projects.value.map(project => ({
      text: project.name,
      value: project.id
    }))
  } catch (error) {
    showToast({ type: 'fail', message: '加载项目失败' })
  }
}

async function onProjectConfirm({ selectedOptions }) {
  const projectId = selectedOptions[0].value
  selectedProject.value = projects.value.find(p => p.id === projectId)
  formData.value.project_id = projectId

  showLoadingToast({
    message: '加载库存...',
    forbidClick: true,
  })

  try {
    const response = await getStock({
      project_id: projectId,
      page: 1,
      page_size: 100
    })

    stockList.value = (response.data || []).filter(s => s.quantity > 0)
    formData.value.items = []

    closeToast()
    showProjectPicker.value = false
  } catch (error) {
    closeToast()
    showToast({ type: 'fail', message: '加载库存失败' })
  }
}

function onDateConfirm() {
  formData.value.requisition_date = currentDate.value.toISOString().split('T')[0]
  showDatePicker.value = false
}

function isItemSelected(stock) {
  return formData.value.items.some(item => item.material_id === stock.material_id)
}

function getFormItem(stock) {
  return formData.value.items.find(item => item.material_id === stock.material_id) || {}
}

function toggleItem(stock) {
  if (isItemSelected(stock)) {
    formData.value.items = formData.value.items.filter(
      item => item.material_id !== stock.material_id
    )
  } else {
    formData.value.items.push({
      material_id: stock.material_id,
      material_name: stock.material_name,
      specification: stock.specification,
      unit: stock.unit,
      requested_quantity: 1
    })
  }
}

async function handleSubmit() {
  if (!selectedProject.value) {
    showToast({ type: 'fail', message: '请选择项目' })
    return
  }

  if (!formData.value.items.length) {
    showToast({ type: 'fail', message: '请选择物资' })
    return
  }

  submitting.value = true
  try {
    await createRequisition({
      project_id: formData.value.project_id,
      department_name: formData.value.department_name,
      requisition_date: formData.value.requisition_date,
      notes: formData.value.notes,
      items: formData.value.items.map(item => ({
        material_id: item.material_id,
        requested_quantity: item.requested_quantity
      }))
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
  // 默认今天
  formData.value.requisition_date = new Date().toISOString().split('T')[0]
})
</script>

<style scoped>
.requisition-create {
  min-height: 100vh;
  background-color: #f7f8fa;
  padding-bottom: 80px;
}

.stock-item {
  background: #fff;
  margin-bottom: 8px;
}

.item-quantity {
  padding: 8px 16px;
  background: #f7f8fa;
}

.submit-section {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 16px;
  background: #fff;
  box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.1);
}
</style>
