<template>
  <div class="materials-page">
    <!-- 搜索栏 -->
    <van-sticky>
      <div class="search-bar">
        <van-search
          v-model="searchText"
          placeholder="搜索物资名称、规格"
          @search="onSearch"
          @clear="onClear"
        />
        <div class="filter-tags">
          <van-tag
            v-for="category in categories"
            :key="category"
            :class="{ active: selectedCategory === category }"
            @click="selectCategory(category)"
          >
            {{ category }}
          </van-tag>
        </div>
      </div>
    </van-sticky>

    <!-- 物资列表 -->
    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <van-list
        v-model:loading="loading"
        :finished="finished"
        finished-text="没有更多了"
        @load="onLoad"
      >
        <div
          v-for="material in materials"
          :key="material.id"
          class="material-card"
          @click="goToDetail(material.id)"
        >
          <div class="material-header">
            <div class="material-name">
              <span v-if="material.code" class="material-code">{{ material.code }}</span>
              <span>{{ material.name }}</span>
            </div>
            <van-icon name="arrow" color="#969799" />
          </div>

          <div class="material-spec" v-if="material.specification">
            {{ material.specification }}
          </div>

          <div class="material-info">
            <div class="info-item">
              <van-icon name="bag-o" size="16" color="#969799" />
              <span>单位: {{ material.unit }}</span>
            </div>
            <div class="info-item" v-if="material.price">
              <van-icon name="gold-coin-o" size="16" color="#ff976a" />
              <span>单价: ¥{{ formatPrice(material.price) }}</span>
            </div>
            <div class="info-item" v-if="material.category">
              <van-icon name="label-o" size="16" color="#1989fa" />
              <span>{{ material.category }}</span>
            </div>
          </div>
        </div>
      </van-list>
    </van-pull-refresh>

    <!-- 空状态 -->
    <van-empty
      v-if="!loading && materials.length === 0 && !refreshing"
      description="暂无物资数据"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { showToast } from 'vant'
import { getMasterList } from '@/api/material'

const router = useRouter()

const materials = ref([])
const categories = ref(['全部', '钢材', '水泥', '砂石', '电气材料', '管道材料', '木材', '涂料', '防水材料', '五金配件', '劳保用品', '工具', '其他'])
const selectedCategory = ref('全部')
const searchText = ref('')
const loading = ref(false)
const finished = ref(false)
const refreshing = ref(false)
const currentPage = ref(1)
const pageSize = 20

async function onLoad() {
  if (refreshing.value) {
    materials.value = []
    currentPage.value = 1
    finished.value = false
  }

  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      per_page: pageSize,
    }

    if (searchText.value) {
      params.search = searchText.value
    }

    if (selectedCategory.value !== '全部') {
      params.category = selectedCategory.value
    }

    const response = await getMasterList(params)
    if (response.success) {
      const { data, total } = response.data

      if (refreshing.value) {
        materials.value = data || []
      } else {
        materials.value.push(...(data || []))
      }

      if (materials.value.length >= total) {
        finished.value = true
      } else {
        currentPage.value++
      }
    }
  } catch (error) {
    console.error('加载物资列表失败:', error)
    showToast('加载失败')
  } finally {
    loading.value = false
    refreshing.value = false
  }
}

function onRefresh() {
  finished.value = false
  onLoad()
}

function onSearch() {
  materials.value = []
  currentPage.value = 1
  finished.value = false
  onLoad()
}

function onClear() {
  if (searchText.value === '') return
  searchText.value = ''
  onSearch()
}

function selectCategory(category) {
  if (selectedCategory.value === category) {
    // 取消选择
    selectedCategory.value = '全部'
  } else {
    selectedCategory.value = category
  }
  materials.value = []
  currentPage.value = 1
  finished.value = false
  onLoad()
}

function goToDetail(id) {
  router.push(`/materials/${id}`)
}

function formatPrice(price) {
  if (!price) return '0.00'
  return Number(price).toFixed(2)
}

onMounted(() => {
  onLoad()
})
</script>

<style scoped>
.materials-page {
  min-height: 100vh;
  background-color: #f7f8fa;
  padding-bottom: 60px;
}

.search-bar {
  background: white;
  padding-bottom: 12px;
}

.filter-tags {
  display: flex;
  gap: 8px;
  padding: 0 16px;
  overflow-x: auto;
  white-space: nowrap;
}

.filter-tags .van-tag {
  cursor: pointer;
}

.filter-tags .van-tag.active {
  background-color: #1989fa;
  color: white;
}

.material-card {
  background: white;
  margin: 12px 16px;
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  cursor: pointer;
}

.material-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.material-name {
  font-size: 16px;
  font-weight: 600;
  color: #323233;
}

.material-code {
  font-size: 12px;
  color: #969799;
  margin-right: 8px;
}

.material-spec {
  font-size: 13px;
  color: #646566;
  margin-bottom: 12px;
  padding: 8px 12px;
  background: #f7f8fa;
  border-radius: 6px;
}

.material-info {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: #646566;
}
</style>
