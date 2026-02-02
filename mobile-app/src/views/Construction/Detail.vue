<template>
  <div class="construction-detail-page">
    <van-nav-bar
      title="日志详情"
      left-arrow
      @click-left="onClickLeft"
    >
      <template #right>
        <van-icon
          v-if="canEdit"
          name="edit"
          size="18"
          @click="goToEdit"
        />
      </template>
    </van-nav-bar>

    <div v-if="loading" class="loading-container">
      <van-loading type="spinner" size="24" />
    </div>

    <div v-else-if="log" class="detail-content">
      <!-- 日期和天气 -->
      <div class="date-weather-card">
        <div class="date-section">
          <van-icon name="calendar-o" size="20" />
          <span class="date-text">{{ formatDate(log.log_date, 'YYYY年MM月DD日 dddd') }}</span>
        </div>
        <div class="weather-section">
          <van-icon :name="getWeatherIcon(log.weather)" size="20" />
          <span class="weather-text">{{ log.weather_text || '晴' }}</span>
          <span class="temperature" v-if="log.temperature">
            {{ log.temperature }}°C
          </span>
        </div>
      </div>

      <!-- 项目信息 -->
      <van-cell-group inset title="项目信息">
        <van-cell title="项目名称" :value="log.project_name" />
        <van-cell title="记录人" :value="log.created_by_name" />
        <van-cell title="记录时间" :value="formatDate(log.created_at)" />
      </van-cell-group>

      <!-- 日志内容 -->
      <van-cell-group inset title="日志内容">
        <div class="log-content">
          <div v-html="log.content" class="rich-text"></div>
        </div>
      </van-cell-group>

      <!-- 图片展示 -->
      <div v-if="log.images && log.images.length > 0" class="images-section">
        <div class="section-title">现场图片</div>
        <van-image-preview
          v-model:show="showImagePreview"
          :images="log.images"
          :start-position="currentImageIndex"
        >
          <template #default>
            <div class="image-grid">
              <div
                v-for="(img, index) in log.images"
                :key="index"
                class="image-item"
                @click="previewImage(index)"
              >
                <van-image :src="img" fit="cover" />
              </div>
            </div>
          </template>
        </van-image-preview>
      </div>

      <!-- 操作按钮 -->
      <div v-if="canDelete" class="action-buttons">
        <van-button
          type="danger"
          block
          @click="handleDelete"
        >
          删除日志
        </van-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { showDialog, showToast } from 'vant'
import { usePermission } from '@/composables/usePermission'
import { getConstructionLogDetail, deleteConstructionLog } from '@/api/construction'
import { formatDate } from '@/utils/date'

const router = useRouter()
const route = useRoute()
const { canCreateConstructionLog: canEdit, canDeleteConstructionLog: canDelete } = usePermission()

const loading = ref(true)
const log = ref(null)
const showImagePreview = ref(false)
const currentImageIndex = ref(0)

// 获取天气图标
function getWeatherIcon(weather) {
  const iconMap = {
    sunny: 'sun-o',
    cloudy: 'cloud-o',
    rainy: 'gem-o',
    snowy: 'snowflake-o',
  }
  return iconMap[weather] || 'sun-o'
}

// 加载详情
async function loadDetail() {
  try {
    const response = await getConstructionLogDetail(route.params.id)
    // 注意：construction_log接口直接返回对象，不是标准格式
    log.value = response
  } catch (error) {
    showToast({
      type: 'fail',
      message: '加载失败',
    })
    router.back()
  } finally {
    loading.value = false
  }
}

// 返回
function onClickLeft() {
  router.back()
}

// 预览图片
function previewImage(index) {
  currentImageIndex.value = index
  showImagePreview.value = true
}

// 跳转编辑
function goToEdit() {
  router.push(`/construction/${route.params.id}/edit`)
}

// 删除日志
function handleDelete() {
  showDialog({
    title: '确认删除',
    message: '确定要删除这条日志吗？',
    showCancelButton: true,
  })
    .then(async () => {
      try {
        await deleteConstructionLog(route.params.id)
        showToast({
          type: 'success',
          message: '删除成功',
        })
        router.back()
      } catch (error) {
        showToast({
          type: 'fail',
          message: '删除失败',
        })
      }
    })
    .catch(() => {
      // 取消
    })
}

onMounted(() => {
  loadDetail()
})
</script>

<style scoped>
.construction-detail-page {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.loading-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 200px;
}

.detail-content {
  padding: 16px 0 80px 0;
}

.date-weather-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin: 16px;
  padding: 16px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  color: white;
}

.date-section,
.weather-section {
  display: flex;
  align-items: center;
  gap: 8px;
}

.date-text,
.weather-text {
  font-size: 14px;
}

.temperature {
  margin-left: 4px;
  padding: 2px 8px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  font-size: 12px;
}

.log-content {
  padding: 16px;
  background: #f7f8fa;
  min-height: 100px;
}

.rich-text {
  font-size: 15px;
  line-height: 1.8;
  color: #323233;
}

.rich-text :deep(img) {
  max-width: 100%;
  border-radius: 8px;
  margin: 8px 0;
}

.rich-text :deep(p) {
  margin: 8px 0;
}

.images-section {
  margin: 16px;
  padding: 16px;
  background: white;
  border-radius: 8px;
}

.section-title {
  font-size: 16px;
  font-weight: bold;
  color: #323233;
  margin-bottom: 12px;
}

.image-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
}

.image-item {
  aspect-ratio: 1;
  overflow: hidden;
  border-radius: 8px;
  cursor: pointer;
}

.action-buttons {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 12px 16px;
  background: white;
  box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.08);
  z-index: 100;
}
</style>
