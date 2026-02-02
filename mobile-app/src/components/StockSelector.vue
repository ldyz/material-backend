<template>
  <div class="stock-selector safe-area-inset-top">
    <!-- 导航栏 -->
    <van-nav-bar
      title="选择库存"
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

    <!-- 库存列表 -->
    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <van-list
        v-model:loading="loading"
        :finished="finished"
        finished-text="没有更多了"
        @load="loadStocks"
      >
        <!-- 复选框模式 -->
        <van-checkbox-group v-if="!instantInput" v-model="selectedStocks">
          <van-cell-group>
            <van-cell
              v-for="item in stocks"
              :key="item.id"
              clickable
              @click="toggleStock(item)"
            >
              <template #title>
                <van-checkbox :name="item.id" ref="checkboxes" @click.stop>
                  <div class="stock-item">
                    <div class="stock-name">{{ item.material_name || item.code || item.material_code || '未知材料' }}</div>
                    <div class="stock-code" v-if="item.code || item.material_code">编码: {{ item.code || item.material_code }}</div>
                    <div class="stock-spec">规格: {{ item.specification || '-' }}</div>
                    <div class="stock-unit">单位: {{ item.unit }}</div>
                    <div class="stock-quantity">
                      库存: <span class="quantity">{{ item.quantity }}</span> {{ item.unit }}
                    </div>
                  </div>
                </van-checkbox>
              </template>
            </van-cell>
          </van-cell-group>
        </van-checkbox-group>

        <!-- 立即输入模式 -->
        <van-cell-group v-else>
          <van-cell
            v-for="item in stocks"
            :key="item.id"
            clickable
            @click="onStockClick(item)"
          >
            <template #title>
              <div class="stock-item">
                <div class="stock-name">{{ item.material_name || item.code || item.material_code || '未知材料' }}</div>
                <div class="stock-code" v-if="item.code || item.material_code">编码: {{ item.code || item.material_code }}</div>
                <div class="stock-spec">规格: {{ item.specification || '-' }}</div>
                <div class="stock-unit">单位: {{ item.unit }}</div>
                <div class="stock-quantity">
                  库存: <span class="quantity">{{ item.quantity }}</span> {{ item.unit }}
                </div>
              </div>
            </template>
          </van-cell>
        </van-cell-group>
      </van-list>
    </van-pull-refresh>

    <!-- 已选择物资栏 -->
    <van-sticky v-if="selectedStocks.length > 0" :offset-bottom="0">
      <div class="selected-bar">
        <div class="selected-info">
          已选择 <span class="count">{{ selectedStocks.length }}</span> 项
        </div>
        <van-button
          size="small"
          @click="onCancel"
        >
          取消
        </van-button>
        <van-button
          size="small"
          type="primary"
          @click="onConfirm"
        >
          确认选择
        </van-button>
      </div>
    </van-sticky>

    <!-- 空状态 -->
    <van-empty
      v-if="stocks.length === 0 && !loading"
      description="暂无库存数据"
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
          <span class="unit">{{ currentStock?.unit }}</span>
        </template>
      </van-field>
      <div class="stock-info">
        可用库存: {{ currentStock?.quantity }} {{ currentStock?.unit }}
      </div>
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
        <van-nav-bar title="已选库存">
          <template #right>
            <van-button size="small" type="danger" @click="onClearCart">
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
                    <div class="item-name">{{ item.material_name || item.code || item.material_code }}</div>
                    <div class="item-spec">{{ item.specification || '-' }}</div>
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
import { ref, onMounted } from 'vue'
import { showToast } from 'vant'
import { getStocks } from '@/api/stock'

const props = defineProps({
  // 只允许选择特定项目的库存
  projectId: [Number, String],
  // 多选模式
  multiple: {
    type: Boolean,
    default: true
  },
  // 立即输入数量模式（点击库存直接弹出输入框）
  instantInput: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['confirm', 'cancel'])

const stocks = ref([])
const selectedStocks = ref([])
const addedItems = ref([]) // 立即输入模式下已添加的库存
const searchKeyword = ref('')
const selectedProject = ref(props.projectId || '')
const showQuantityDialog = ref(false)
const showCart = ref(false) // 购物车弹窗
const currentStock = ref(null)
const quantityInput = ref('')
const remarkInput = ref('')

// 分页相关
const loading = ref(false)
const finished = ref(false)
const refreshing = ref(false)
const currentPage = ref(1)

// 项目选项
const projectOptions = ref([
  { text: '全部项目', value: '' }
])

// 加载库存列表
async function loadStocks() {
  if (loading.value || finished.value) return
  
  try {
    loading.value = true
    
    const params = {
      page: currentPage.value,
      per_page: 20,
      quantity_min: 0.01  // 只获取库存大于0的物资
    }

    // 如果指定了项目，只查询该项目的库存
    if (selectedProject.value) {
      params.project_id = selectedProject.value
    }

    // 如果有搜索关键词
    if (searchKeyword.value) {
      params.search = searchKeyword.value
    }

    // 适配统一响应格式
    const { data, pagination } = await getStocks(params)
    const newStocks = data || []

    // 首次加载或刷新
    if (currentPage.value === 1) {
      stocks.value = newStocks
    } else {
      // 追加数据
      stocks.value.push(...newStocks)
    }

    // 提取项目选项（只在首次加载时）
    if (currentPage.value === 1) {
      const projects = new Set()
      newStocks.forEach(s => {
        if (s.project_id && s.project_name) {
          projects.add(JSON.stringify({
            value: s.project_id,
            text: s.project_name
          }))
        }
      })

      projectOptions.value = [
        { text: '全部项目', value: '' },
        ...Array.from(projects).map(p => JSON.parse(p))
      ]
    }

    // 判断是否还有更多数据
    if (pagination) {
      finished.value = currentPage.value >= pagination.pages
    } else {
      finished.value = newStocks.length < 20
    }

    // 如果还有数据，页码+1
    if (!finished.value) {
      currentPage.value++
    }

  } catch (error) {
    console.error('加载库存失败:', error)
    showToast({
      type: 'fail',
      message: '加载库存失败'
    })
  } finally {
    loading.value = false
    refreshing.value = false
  }
}

// 下拉刷新
async function onRefresh() {
  finished.value = false
  currentPage.value = 1
  await loadStocks()
}

// 切换库存选择
function toggleStock(item) {
  const index = selectedStocks.value.indexOf(item.id)
  if (index > -1) {
    selectedStocks.value.splice(index, 1)
  } else {
    if (props.multiple) {
      selectedStocks.value.push(item.id)
    } else {
      selectedStocks.value = [item.id]
    }
  }
}

// 立即输入模式：点击库存
function onStockClick(item) {
  currentStock.value = item
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

  const qty = parseFloat(quantityInput.value)
  if (qty > currentStock.value.quantity) {
    showToast(`数量不能超过库存 ${currentStock.value.quantity}`)
    return
  }

  const stockWithQuantity = {
    ...currentStock.value,
    quantity: qty,
    remark: remarkInput.value
  }

  // 检查是否已经添加过该库存
  const existsIndex = addedItems.value.findIndex(item => item.id === currentStock.value.id)
  if (existsIndex > -1) {
    // 如果已存在，更新数量和备注
    addedItems.value[existsIndex] = stockWithQuantity
    showToast(`已更新 ${currentStock.value.material_name || currentStock.value.code}`)
  } else {
    // 添加新库存
    addedItems.value.push(stockWithQuantity)
    showToast({
      type: 'success',
      message: `已添加 ${currentStock.value.material_name || currentStock.value.code}`
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
  // 清空已选择和数据
  selectedStocks.value = []
  stocks.value = []
  finished.value = false
  currentPage.value = 1
  loadStocks()
}

// 搜索
function onSearch() {
  // 重置并重新加载
  stocks.value = []
  selectedStocks.value = []
  finished.value = false
  currentPage.value = 1
  loadStocks()
}

// 确认选择
function onConfirm() {
  const selectedList = stocks.value.filter(s =>
    selectedStocks.value.includes(s.id)
  )
  emit('confirm', selectedList)
}

// 取消
function onCancel() {
  emit('cancel')
}

onMounted(() => {
  // 如果指定了项目，自动选中
  if (props.projectId) {
    selectedProject.value = props.projectId
  }
  loadStocks()
})
</script>

<style scoped>
.stock-selector {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.stock-item {
  padding: 8px 0;
}

.stock-name {
  font-weight: bold;
  font-size: 15px;
  margin-bottom: 4px;
}

.stock-code,
.stock-spec,
.stock-unit,
.stock-quantity {
  font-size: 13px;
  color: #969799;
  margin-bottom: 2px;
}

.stock-quantity .quantity {
  color: #07c160;
  font-weight: bold;
  font-size: 15px;
}

.selected-bar {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 8px;
  padding: 12px 16px;
  background: white;
  border-top: 1px solid #ebedf0;
  box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.08);
}

.selected-info {
  flex: 1;
  font-size: 14px;
  color: #646566;
}

.selected-info .count {
  color: #1989fa;
  font-weight: bold;
  font-size: 16px;
}

.stock-info {
  padding: 8px 16px;
  font-size: 14px;
  color: #646566;
  background: #f7f8fa;
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
