<template>
  <el-dialog
    v-model="dialogVisible"
    title="裁剪头像"
    width="700px"
    :close-on-click-modal="false"
    @close="handleClose"
    @open="onDialogOpen"
    @opened="onDialogOpened"
  >
    <div v-if="imageSrc" class="cropper-container">
      <!-- 左侧：裁剪区域 -->
      <div class="cropper-main">
        <vue-cropper
          v-if="cropperReady"
          ref="cropperRef"
          :img="imageSrc"
          :output-size="1"
          :output-type="outputType"
          :auto-crop="true"
          :auto-crop-width="200"
          :auto-crop-height="200"
          :fixed="true"
          :fixed-number="[1, 1]"
          :can-move="true"
          :can-move-box="true"
          :can-scale="true"
          :high="true"
          :max-img-size="3000"
          mode="cover"
          @real-time="onRealTime"
          @img-load="onImgLoad"
          @img-error="onImgError"
        />
      </div>

      <!-- 右侧：预览区域 -->
      <div class="preview-sidebar">
        <div class="preview-title">头像预览</div>

        <div class="preview-item">
          <div class="preview-label">200×200</div>
          <div class="preview-box-200">
            <div :style="previewStyle">
              <img v-if="previewUrl" :src="previewUrl" />
            </div>
          </div>
        </div>

        <div class="preview-item">
          <div class="preview-label">100×100</div>
          <div class="preview-box-100">
            <div :style="previewStyleSmall">
              <img v-if="previewUrl" :src="previewUrl" />
            </div>
          </div>
        </div>


        <div class="tips">
          <p>💡 提示：</p>
          <p>• 拖动图片调整位置</p>
          <p>• 滚轮缩放图片大小</p>
          <p>• 拖动方框调整裁剪区域</p>
        </div>
      </div>
    </div>

    <div v-else class="loading-placeholder">
      <el-icon class="is-loading" :size="40"><Loading /></el-icon>
      <p>正在加载图片...</p>
    </div>

    <div class="cropper-controls">
      <el-button-group>
        <el-button @click="rotateLeft" :icon="RefreshLeft">向左旋转</el-button>
        <el-button @click="rotateRight" :icon="RefreshRight">向右旋转</el-button>
      </el-button-group>
      <el-button-group style="margin-left: 12px">
        <el-button @click="zoomIn" :icon="ZoomIn">放大</el-button>
        <el-button @click="zoomOut" :icon="ZoomOut">缩小</el-button>
      </el-button-group>
      <el-button style="margin-left: 12px" @click="reset" :icon="Refresh">重置</el-button>
    </div>

    <template #footer>
      <el-button @click="handleClose">取消</el-button>
      <el-button type="primary" :loading="uploading" @click="handleConfirm">
        确认上传
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue'
import { VueCropper } from 'vue-cropper/next'
import 'vue-cropper/next/dist/index.css'
import { ElMessage } from 'element-plus'
import { RefreshLeft, RefreshRight, ZoomIn, ZoomOut, Refresh, Loading } from '@element-plus/icons-vue'
import { authApi } from '@/api'

const props = defineProps({
  modelValue: Boolean,
  imageFile: File
})

const emit = defineEmits(['update:modelValue', 'success'])

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const cropperRef = ref(null)
const imageSrc = ref('')
const previewUrl = ref('')
const previewStyle = ref({})
const previewStyleSmall = ref({})
const previewStyleTiny = ref({})
const uploading = ref(false)
const outputType = ref('jpeg')
const cropperReady = ref(false)

watch(() => props.imageFile, (file) => {
  if (file) {
    loadImage(file)
  }
})

function onDialogOpen() {
  // 重置状态
  cropperReady.value = false
  if (props.imageFile) {
    loadImage(props.imageFile)
  }
}

function onDialogOpened() {
  // 对话框完全打开后再显示裁剪器
  console.log('Dialog opened, imageSrc:', !!imageSrc.value)
  nextTick(() => {
    setTimeout(() => {
      cropperReady.value = true
      console.log('Cropper ready set to true')
    }, 100)
  })
}

function loadImage(file) {
  console.log('Loading image:', file.name)
  const reader = new FileReader()
  reader.onload = (e) => {
    console.log('Image loaded to base64, length:', e.target.result?.length)
    imageSrc.value = e.target.result
    previewUrl.value = e.target.result
    updatePreviewStyles()
  }
  reader.onerror = () => {
    console.error('Failed to read image')
    ElMessage.error('读取图片失败')
  }
  reader.readAsDataURL(file)
}

function onImgLoad(msg) {
  console.log('图片加载完成', msg)
  // 图片加载完成后更新预览
  nextTick(() => {
    updatePreview()
  })
}

function onImgError(err) {
  console.error('图片加载失败', err)
  ElMessage.error('图片加载失败')
}

function onRealTime(data) {
  // 实时预览 - 当拖动或缩放时调用
  console.log('Real-time event:', data)
  if (!cropperRef.value) return

  // 使用 getCropData 获取裁剪后的预览图
  try {
    cropperRef.value.getCropData((url) => {
      if (url) {
        previewUrl.value = url
      }
    })
  } catch (e) {
    console.log('getCropData in real-time error:', e)
  }
}

function updatePreviewStyles() {
  previewStyle.value = {
    width: '200px',
    height: '200px',
    overflow: 'hidden',
    borderRadius: '50%',
    border: '3px solid #fff',
    boxShadow: '0 2px 8px rgba(0,0,0,0.15)'
  }

  previewStyleSmall.value = {
    width: '100px',
    height: '100px',
    overflow: 'hidden',
    borderRadius: '50%',
    border: '2px solid #fff',
    boxShadow: '0 2px 8px rgba(0,0,0,0.15)'
  }

  previewStyleTiny.value = {
    width: '50px',
    height: '50px',
    overflow: 'hidden',
    borderRadius: '50%',
    border: '2px solid #fff',
    boxShadow: '0 2px 8px rgba(0,0,0,0.15)'
  }
}

function updatePreview() {
  if (!cropperRef.value) return

  try {
    cropperRef.value.getCropData((data) => {
      if (data) {
        previewUrl.value = data
      }
    })
  } catch (e) {
    console.log('getCropData error:', e)
  }
}

function rotateLeft() {
  if (!cropperRef.value) return
  cropperRef.value.rotateLeft()
  nextTick(() => updatePreview())
}

function rotateRight() {
  if (!cropperRef.value) return
  cropperRef.value.rotateRight()
  nextTick(() => updatePreview())
}

function zoomIn() {
  if (!cropperRef.value) return
  cropperRef.value.changeScale(1)
  nextTick(() => updatePreview())
}

function zoomOut() {
  if (!cropperRef.value) return
  cropperRef.value.changeScale(-1)
  nextTick(() => updatePreview())
}

function reset() {
  if (!cropperRef.value) return
  cropperRef.value.refresh()
  nextTick(() => updatePreview())
}

async function handleConfirm() {
  if (!cropperRef.value) {
    ElMessage.warning('请先选择图片')
    return
  }

  uploading.value = true
  try {
    cropperRef.value.getCropBlob((blob) => {
      if (!blob) {
        ElMessage.error('裁剪失败，请重试')
        uploading.value = false
        return
      }

      const formData = new FormData()
      const fileName = `avatar_${Date.now()}.${outputType.value}`
      formData.append('avatar', blob, fileName)

      authApi.uploadAvatar(formData).then(() => {
        ElMessage.success('头像上传成功')
        emit('success')
        handleClose()
      }).catch((error) => {
        console.error('头像上传失败:', error)
        ElMessage.error(error?.message || '头像上传失败')
      }).finally(() => {
        uploading.value = false
      })
    })
  } catch (error) {
    console.error('裁剪失败:', error)
    ElMessage.error('裁剪失败，请重试')
    uploading.value = false
  }
}

function handleClose() {
  imageSrc.value = ''
  previewUrl.value = ''
  cropperReady.value = false
  previewStyle.value = {}
  previewStyleSmall.value = {}
  previewStyleTiny.value = {}
  emit('update:modelValue', false)
}
</script>

<style scoped>
.cropper-container {
  display: flex;
  gap: 20px;
  height: 450px;
}

.cropper-main {
  flex: 1;
  height: 450px;
  background: #f0f0f0;
  border-radius: 8px;
  overflow: hidden;
  position: relative;
}

.preview-sidebar {
  width: 240px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  padding: 16px;
  background: #f5f7fa;
  border-radius: 8px;
}

.preview-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  text-align: center;
}

.preview-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  width: 100%;
}

.preview-label {
  font-size: 13px;
  color: #606266;
  font-weight: 500;
}

.preview-box-200,
.preview-box-100 {
  display: flex;
  align-items: center;
  justify-content: center;
  background: #fff;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  padding: 8px;
  border-radius: 8px;
}

.preview-box-200 {
  width: 216px;
  height: 216px;
}

.preview-box-100 {
  width: 116px;
  height: 116px;
}

.preview-box-200 > div,
.preview-box-100 > div {
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

.preview-box-200 img,
.preview-box-100 img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.tips {
  margin-top: auto;
  padding: 12px;
  background: #fff;
  border-radius: 8px;
  font-size: 12px;
  color: #909399;
  line-height: 1.8;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.tips p {
  margin: 0;
}

.cropper-controls {
  margin-top: 20px;
  display: flex;
  justify-content: center;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}

.loading-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 450px;
  color: #909399;
}

.loading-placeholder p {
  margin-top: 16px;
  font-size: 14px;
}

:deep(.vue-cropper) {
  width: 100%;
  height: 100%;
}
</style>
