<template>
  <div class="inbound-create">
    <van-nav-bar title="新建入库单" left-arrow @click-left="router.back()" />

    <van-form @submit="handleSubmit">
      <van-cell-group inset title="基本信息">
        <van-field
          v-model="formData.supplier_name"
          name="supplier"
          label="供应商"
          placeholder="请输入供应商名称"
          :rules="[{ required: true, message: '请输入供应商名称' }]"
        />
        <van-field
          v-model="formData.inbound_date"
          name="date"
          label="入库日期"
          placeholder="请选择入库日期"
          readonly
          is-link
          @click="showDatePicker = true"
          :rules="[{ required: true, message: '请选择入库日期' }]"
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

      <van-cell-group inset title="选择物资计划">
        <van-cell
          title="选择计划"
          :value="selectedPlan ? selectedPlan.plan_number : '请选择已批准的计划'"
          is-link
          @click="showPlanPicker = true"
          :rules="[{ required: true, message: '请选择物资计划' }]"
        />
      </van-cell-group>

      <van-cell-group v-if="formData.items.length" inset title="物资明细">
        <van-card
          v-for="(item, index) in formData.items"
          :key="index"
          class="item-card"
        >
          <template #title>
            {{ item.material_name }}
          </template>
          <template #desc>
            规格：{{ item.specification || '-' }}
          </template>
          <template #footer>
            <van-field
              v-model.number="item.quantity"
              type="number"
              label="入库数量"
              placeholder="请输入数量"
              :min="1"
            />
            <van-field
              v-model.number="item.unit_price"
              type="number"
              label="单价"
              placeholder="请输入单价"
            />
            <van-button
              size="small"
              type="danger"
              plain
              @click="removeItem(index)"
            >
              移除
            </van-button>
          </template>
        </van-card>
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

    <!-- 计划选择器 -->
    <van-popup v-model:show="showPlanPicker" position="bottom" :style="{ height: '60%' }">
      <van-picker
        :columns="planColumns"
        @confirm="onPlanConfirm"
        @cancel="showPlanPicker = false"
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
import { createInbound, getApprovedPlans } from '@/api/inbound'
import { getPlanDetail } from '@/api/material_plan'

const router = useRouter()

const showPlanPicker = ref(false)
const showDatePicker = ref(false)
const submitting = ref(false)

const plans = ref([])
const selectedPlan = ref(null)

const currentDate = ref(new Date())
const formData = ref({
  supplier_name: '',
  inbound_date: '',
  notes: '',
  items: []
})

// 计划选择器列
const planColumns = ref([])

async function loadPlans() {
  try {
    const response = await getApprovedPlans({ page: 1, page_size: 100 })
    plans.value = response.data || []

    // 格式化为选择器列
    planColumns.value = plans.value.map(plan => ({
      text: `${plan.plan_number} - ${plan.project_name || '未命名项目'}`,
      value: plan.id
    }))
  } catch (error) {
    showToast({ type: 'fail', message: '加载计划失败' })
  }
}

async function onPlanConfirm({ selectedOptions }) {
  const planId = selectedOptions[0].value
  selectedPlan.value = plans.value.find(p => p.id === planId)

  showLoadingToast({
    message: '加载物料...',
    forbidClick: true,
  })

  try {
    const response = await getPlanDetail(planId)
    const plan = response.data

    // 自动填充物料列表
    formData.value.items = (plan.items || []).map(item => ({
      material_id: item.material_id,
      material_name: item.material_name,
      specification: item.specification,
      unit: item.unit,
      quantity: item.quantity || 0,
      unit_price: 0
    }))

    closeToast()
    showPlanPicker.value = false
  } catch (error) {
    closeToast()
    showToast({ type: 'fail', message: '加载物料失败' })
  }
}

function onDateConfirm() {
  formData.value.inbound_date = currentDate.value.toISOString().split('T')[0]
  showDatePicker.value = false
}

function removeItem(index) {
  formData.value.items.splice(index, 1)
}

async function handleSubmit() {
  if (!selectedPlan.value) {
    showToast({ type: 'fail', message: '请选择物资计划' })
    return
  }

  if (!formData.value.items.length) {
    showToast({ type: 'fail', message: '请添加物资明细' })
    return
  }

  submitting.value = true
  try {
    await createInbound({
      plan_id: selectedPlan.value.id,
      supplier_name: formData.value.supplier_name,
      inbound_date: formData.value.inbound_date,
      notes: formData.value.notes,
      items: formData.value.items.map(item => ({
        material_id: item.material_id,
        quantity: item.quantity,
        unit_price: item.unit_price
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
  loadPlans()
  // 默认今天
  formData.value.inbound_date = new Date().toISOString().split('T')[0]
})
</script>

<style scoped>
.inbound-create {
  min-height: 100vh;
  background-color: #f7f8fa;
  padding-bottom: 80px;
}

.item-card {
  margin: 8px 16px;
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
