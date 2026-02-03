<template>
  <div class="material-detail-page">
    <van-nav-bar
      title="物资详情"
      left-text="返回"
      @click-left="onClickLeft"
    />

    <div v-if="loading" class="loading-container">
      <van-loading size="24" vertical>加载中...</van-loading>
    </div>

    <div v-else-if="material">
      <!-- 基本信息 -->
      <van-cell-group inset class="info-section">
        <van-cell title="物资编码" :value="material.code" />
        <van-cell title="物资名称" :value="material.name" />
        <van-cell title="规格型号" :value="material.specification" />
        <van-cell title="计量单位" :value="material.unit" />
        <van-cell
          title="参考价格"
          :value="material.price ? `¥${formatPrice(material.price)}` : '-'"
        />
        <van-cell title="物资分类" :value="material.category" />
        <van-cell title="安全库存" :value="material.safety_stock || '-'" />
        <van-cell title="描述" :value="material.description || '无'" />
      </van-cell-group>

      <!-- 库存信息 -->
      <van-cell-group inset title="库存信息" class="stock-section">
        <van-cell
          title="当前库存"
          :value="stockInfo.quantity || 0"
        />
        <van-cell
          title="库存状态"
          :value="getStockStatus(stockInfo.quantity, material.safety_stock)"
        >
          <template #extra>
            <van-tag
              :type="getStockStatusType(stockInfo.quantity, material.safety_stock)"
            >
              {{ getStockStatusText(stockInfo.quantity, material.safety_stock) }}
            </van-tag>
          </template>
        </van-cell>
      </van-cell-group>

      <!-- 操作按钮 -->
      <div class="action-buttons">
        <van-button
          type="primary"
          icon="cart-o"
          @click="goToStock"
        >
          查看库存
        </van-button>
        <van-button
          icon="todo-list-o"
          @click="goToPlans"
        >
          查看计划
        </van-button>
      </div>
    </div>

    <van-empty v-else description="物资不存在" />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { showToast } from 'vant'
import { getMasterDetail, getProjectMaterials } from '@/api/material'
import { getStockList } from '@/api/stock'

const router = useRouter()
const route = useRoute()

const material = ref(null)
const stockInfo = ref({})
const loading = ref(true)

async function loadDetail() {
  loading.value = true
  try {
    const id = route.params.id
    const response = await getMasterDetail(id)
    if (response.success) {
      material.value = response.data
      // 加载库存信息
      loadStockInfo(material.value.name)
    } else {
      showToast('加载失败')
    }
  } catch (error) {
    console.error('加载物资详情失败:', error)
    showToast('加载失败')
  } finally {
    loading.value = false
  }
}

async function loadStockInfo(materialName) {
  try {
    const response = await getStockList({
      search: materialName,
      page: 1,
      per_page: 1
    })
    if (response.success && response.data.data && response.data.data.length > 0) {
      stockInfo.value = response.data.data[0]
    }
  } catch (error) {
    console.error('加载库存信息失败:', error)
  }
}

function getStockStatus(quantity, safetyStock) {
  if (!quantity || !safetyStock) return '未知'
  const q = Number(quantity)
  const s = Number(safetyStock)
  if (q === 0) return '无库存'
  if (q < s * 0.5) return '严重不足'
  if (q < s) return '库存偏低'
  return '库存充足'
}

function getStockStatusType(quantity, safetyStock) {
  if (!quantity || !safetyStock) return 'default'
  const q = Number(quantity)
  const s = Number(safetyStock)
  if (q === 0) return 'danger'
  if (q < s * 0.5) return 'danger'
  if (q < s) return 'warning'
  return 'success'
}

function getStockStatusText(quantity, safetyStock) {
  if (!quantity || !safetyStock) return '未知'
  const q = Number(quantity)
  const s = Number(safetyStock)
  if (q === 0) return '无库存'
  if (q < s * 0.5) return '严重不足'
  if (q < s) return '库存偏低'
  return '正常'
}

function formatPrice(price) {
  if (!price) return '0.00'
  return Number(price).toFixed(2)
}

function onClickLeft() {
  router.back()
}

function goToStock() {
  router.push('/stock?search=' + encodeURIComponent(material.value.name))
}

function goToPlans() {
  router.push('/plans?search=' + encodeURIComponent(material.value.name))
}

onMounted(() => {
  loadDetail()
})
</script>

<style scoped>
.material-detail-page {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.loading-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 50vh;
}

.info-section,
.stock-section {
  margin: 16px;
}

.action-buttons {
  display: flex;
  gap: 12px;
  padding: 16px;
}
</style>
