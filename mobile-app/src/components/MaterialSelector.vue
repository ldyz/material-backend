<template>
  <div class="material-selector safe-area-inset-top">
    <!-- 导航栏 -->
    <van-nav-bar
      title="选择物资"
      left-arrow
      @click-left="onCancel"
    />

    <!-- 筛选条件 -->
    <van-dropdown-menu>
      <van-dropdown-item
        v-model="selectedProject"
        :options="projectOptions"
        @change="onProjectChange"
        title="选择项目"
      />
    </van-dropdown-menu>

    <!-- 搜索 -->
    <van-search
      v-model="searchKeyword"
      placeholder="搜索物资名称或编码"
      @search="onSearch"
    />

    <!-- 物资列表 -->
    <van-checkbox-group v-if="!instantInput" v-model="selectedMaterials">
      <van-cell-group>
        <van-cell
          v-for="item in filteredMaterials"
          :key="item.id"
          clickable
          @click="toggleMaterial(item)"
        >
          <template #title>
            <van-checkbox :name="item.id" ref="checkboxes" @click.stop>
              <div class="material-item">
                <div class="material-name">{{ item.name }}</div>
                <div class="material-code">编码: {{ item.material || '-' }}</div>
                <div class="material-spec">规格: {{ item.specification || '-' }}</div>
                <div class="material-unit">单位: {{ item.unit }}</div>
              </div>
            </van-checkbox>
          </template>
        </van-cell>
      </van-cell-group>
    </van-checkbox-group>

    <!-- 立即输入模式 -->
    <van-cell-group v-else>
      <van-cell
        v-for="item in filteredMaterials"
        :key="item.id"
        clickable
        @click="onMaterialClick(item)"
      >
        <template #title>
          <div class="material-item">
            <div class="material-name">{{ item.name }}</div>
            <div class="material-code">编码: {{ item.material || '-' }}</div>
            <div class="material-spec">规格: {{ item.specification || '-' }}</div>
            <div class="material-unit">单位: {{ item.unit }}</div>
          </div>
        </template>
      </van-cell>
    </van-cell-group>

    <!-- 已选择物资栏 -->
    <van-sticky v-if="selectedMaterials.length > 0" :offset-bottom="0">
      <div class="selected-bar">
        <div class="selected-info">
          已选择 <span class="count">{{ selectedMaterials.length }}</span> 项
        </div>
        <van-button size="small" type="primary" @click="onConfirm">
          确认选择
        </van-button>
      </div>
    </van-sticky>

    <!-- 空状态 -->
    <van-empty
      v-if="filteredMaterials.length === 0"
      description="暂无物资数据"
    />

    <!-- 输入数量弹窗 -->
    <van-dialog
      v-model:show="showQuantityDialog"
      title="输入数量"
      show-cancel-button
      @confirm="onQuantityConfirm"
    >
      <van-field
        v-model="quantityInput"
        type="number"
        label="数量"
        placeholder="请输入数量"
        :rules="[{ required: true, message: '请输入数量' }]"
      >
        <template #button>
          <span class="unit">{{ currentMaterial?.unit }}</span>
        </template>
      </van-field>
      <van-field
        v-model="remarkInput"
        label="备注"
        type="textarea"
        placeholder="备注（可选）"
        rows="2"
      />
    </van-dialog>

    <!-- 购物车浮动按钮 - 使用 Teleport 固定到 body -->
    <teleport to="body">
      <div
        v-if="instantInput && addedItems.length > 0"
        class="cart-float-button"
        @click="showCart = true"
      >
        <van-badge
          :content="addedItems.length"
          :offset="[-5, 5]"
        >
          <van-icon name="cart-o" size="24" color="#fff" />
        </van-badge>
      </div>
    </teleport>

    <!-- 购物车弹窗 -->
    <van-popup
      v-model:show="showCart"
      position="bottom"
      :style="{ height: '70%' }"
      round
    >
      <div class="cart-popup">
        <van-nav-bar
          title="已选物资"
          @click-right="onClearCart"
        >
          <template #right>
            <van-button size="small" type="danger">
              清空
            </van-button>
          </template>
        </van-nav-bar>

        <div class="cart-list">
          <van-cell-group>
            <van-swipe-cell
              v-for="(item, index) in addedItems"
              :key="item.id"
            >
              <van-cell>
                <template #title>
                  <div class="cart-item">
                    <div class="item-name">{{ item.name }}</div>
                    <div class="item-spec">{{ item.specification || item.spec || '-' }}</div>
                    <div class="item-unit">{{ item.unit }}</div>
                    <div class="item-quantity">{{ item.quantity }}</div>
                  </div>
                </template>
              </van-cell>
              <template #right>
                <van-button
                  square
                  type="danger"
                  text="删除"
                  @click="removeFromCart(index)"
                />
              </template>
            </van-swipe-cell>
          </van-cell-group>
        </div>

        <div class="cart-actions">
          <van-button block @click="showCart = false">
            取消
          </van-button>
          <van-button
            block
            type="primary"
            @click="onFinish"
          >
            完成({{ addedItems.length }})
          </van-button>
        </div>
      </div>
    </van-popup>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { showToast } from 'vant'
import { getMaterials } from '@/api/material'

const props = defineProps({
  // 只允许选择特定项目的物资
  projectId: [Number, String],
  // 多选模式
  multiple: {
    type: Boolean,
    default: true
  },
  // 立即输入数量模式（点击物资直接弹出输入框）
  instantInput: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['confirm', 'cancel'])

const materials = ref([])
const selectedMaterials = ref([])
const addedItems = ref([]) // 立即输入模式下已添加的物资
const searchKeyword = ref('')
const selectedProject = ref(props.projectId || '')
const showQuantityDialog = ref(false)
const showCart = ref(false) // 购物车弹窗
const currentMaterial = ref(null)
const quantityInput = ref('')
const remarkInput = ref('')

// 项目选项
const projectOptions = ref([
  { text: '全部项目', value: '' }
])

// 过滤后的物资列表
const filteredMaterials = computed(() => {
  let result = materials.value

  // 按项目筛选
  if (selectedProject.value) {
    result = result.filter(m => m.project_id == selectedProject.value)
  }

  // 按关键词搜索
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    result = result.filter(m =>
      (m.name && m.name.toLowerCase().includes(keyword)) ||
      (m.material && m.material.toLowerCase().includes(keyword)) ||
      (m.specification && m.specification.toLowerCase().includes(keyword))
    )
  }

  return result
})

// 加载物资列表
async function loadMaterials() {
  try {
    const params = {
      filter: 'unstored' // 只显示未完全入库的物资
    }
    if (selectedProject.value) {
      params.project_id = selectedProject.value
    }

    // 适配统一响应格式
    const { data } = await getMaterials(params)

    materials.value = data || []

    // 填充项目选项
    const projects = new Set()
    materials.value.forEach(m => {
      if (m.project_id && m.project_name) {
        projects.add(JSON.stringify({
          value: m.project_id,
          text: m.project_name
        }))
      }
    })

    projectOptions.value = [
      { text: '全部项目', value: '' },
      ...Array.from(projects).map(p => JSON.parse(p))
    ]
  } catch (error) {
    showToast({
      type: 'fail',
      message: '加载物资失败'
    })
  }
}

// 切换物资选择
function toggleMaterial(item) {
  const index = selectedMaterials.value.indexOf(item.id)
  if (index > -1) {
    selectedMaterials.value.splice(index, 1)
  } else {
    if (props.multiple) {
      selectedMaterials.value.push(item.id)
    } else {
      selectedMaterials.value = [item.id]
    }
  }
}

// 立即输入模式：点击物资
function onMaterialClick(item) {
  currentMaterial.value = item
  quantityInput.value = ''
  remarkInput.value = ''
  showQuantityDialog.value = true
}

// 确认数量输入
function onQuantityConfirm() {
  if (!quantityInput.value || quantityInput.value <= 0) {
    showToast('请输入有效的数量')
    return
  }

  const itemWithQuantity = {
    ...currentMaterial.value,
    quantity: parseInt(quantityInput.value),
    remark: remarkInput.value
  }

  // 检查是否已经添加过该物资
  const existsIndex = addedItems.value.findIndex(item => item.id === currentMaterial.value.id)
  if (existsIndex > -1) {
    // 如果已存在，更新数量和备注
    addedItems.value[existsIndex] = itemWithQuantity
    showToast(`已更新 ${currentMaterial.value.name}`)
  } else {
    // 添加新物资
    addedItems.value.push(itemWithQuantity)
    showToast({
      type: 'success',
      message: `已添加 ${currentMaterial.value.name}`
    })
  }

  showQuantityDialog.value = false
  // 不关闭选择器，继续选择
}

// 完成所有选择
function onFinish() {
  showCart.value = false
  emit('confirm', [...addedItems.value])
}

// 从购物车删除
function removeFromCart(index) {
  addedItems.value.splice(index, 1)
  showToast({
    type: 'success',
    message: '已删除'
  })
}

// 清空购物车
function onClearCart() {
  addedItems.value = []
  showCart.value = false
  showToast({
    type: 'success',
    message: '已清空购物车'
  })
}

// 项目变化
function onProjectChange() {
  // 清空已选择
  selectedMaterials.value = []
}

// 搜索
function onSearch() {
  // 搜索会自动触发 computed 更新
}

// 确认选择
function onConfirm() {
  const selectedList = materials.value.filter(m =>
    selectedMaterials.value.includes(m.id)
  )
  emit('confirm', selectedList)
}

// 取消
function onCancel() {
  emit('cancel')
}

onMounted(() => {
  loadMaterials()
})
</script>

<style scoped>
.material-selector {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.material-item {
  padding: 8px 0;
}

.material-name {
  font-weight: bold;
  font-size: 15px;
  margin-bottom: 4px;
}

.material-code,
.material-spec,
.material-unit {
  font-size: 13px;
  color: #969799;
  margin-bottom: 2px;
}

.selected-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: white;
  border-top: 1px solid #ebedf0;
  box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.08);
}

.selected-info {
  font-size: 14px;
  color: #646566;
}

.selected-info .count {
  color: #1989fa;
  font-weight: bold;
  font-size: 16px;
}

/* 购物车浮动按钮 */
.cart-float-button {
  position: fixed;
  bottom: calc(80px + env(safe-area-inset-bottom));
  right: 20px;
  width: 56px;
  height: 56px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  cursor: pointer;
  z-index: 9999;
  transition: all 0.3s;
}

.cart-float-button:active {
  transform: scale(0.95);
}

/* 购物车样式 */
.cart-popup {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.cart-list {
  flex: 1;
  overflow-y: auto;
}

.cart-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 0;
}

.cart-item .item-name {
  flex: 1;
  font-weight: 500;
  font-size: 14px;
  color: #323233;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.cart-item .item-spec {
  width: 80px;
  text-align: center;
  font-size: 13px;
  color: #969799;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.cart-item .item-unit {
  width: 60px;
  text-align: center;
  font-size: 13px;
  color: #969799;
}

.cart-item .item-quantity {
  width: 80px;
  text-align: right;
  font-size: 15px;
  font-weight: bold;
  color: #1989fa;
}

.cart-actions {
  display: flex;
  gap: 8px;
  padding: 12px 16px;
  background: white;
  border-top: 1px solid #ebedf0;
  box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.08);
}
</style>
