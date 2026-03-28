<template>
  <van-popup
    v-model:show="dialogVisible"
    position="bottom"
    :style="{ height: '85%' }"
    round
    @close="handleClose"
    @opened="onPopupOpened"
  >
    <div class="cropper-popup">
      <van-nav-bar
        title="裁剪头像"
        left-text="取消"
        right-text="确认"
        @click-left="handleClose"
        @click-right="handleConfirm"
      />

      <!-- 裁剪区域 -->
      <div class="cropper-container">
        <div ref="cropperContainerRef" class="cropper-main">
          <img ref="imageRef" :src="imageSrc" style="display: none" />
        </div>
      </div>

      <!-- 预览区域 -->
      <div class="preview-area">
        <div class="preview-title">头像预览</div>
        <div class="preview-list">
          <div class="preview-item">
            <div class="preview-label">120×120</div>
            <div class="preview-box-120">
              <div class="preview-circular">
                <img v-if="previewUrl" :src="previewUrl" />
              </div>
            </div>
          </div>
          <div class="preview-item">
            <div class="preview-label">60×60</div>
            <div class="preview-box-60">
              <div class="preview-circular">
                <img v-if="previewUrl" :src="previewUrl" />
              </div>
            </div>
          </div>
          <div class="preview-item">
            <div class="preview-label">40×40</div>
            <div class="preview-box-40">
              <div class="preview-circular">
                <img v-if="previewUrl" :src="previewUrl" />
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 控制按钮 -->
      <div class="cropper-controls">
        <div class="button-row">
          <van-button icon="replay" size="small" @click="rotateLeft">左旋转</van-button>
          <van-button icon="replay" style="transform: scaleX(-1)" size="small" @click="rotateRight">右旋转</van-button>
        </div>
        <div class="button-row">
          <van-button icon="plus" size="small" @click="zoomIn">放大</van-button>
          <van-button icon="minus" size="small" @click="zoomOut">缩小</van-button>
        </div>
        <van-button block size="small" icon="replay" @click="reset">重置</van-button>
      </div>

      <!-- 加载提示 -->
      <div v-if="!imageLoaded" class="loading-placeholder">
        <van-loading size="24px">正在加载图片...</van-loading>
      </div>

      <van-loading v-if="uploading" class="uploading-loader" size="24px">上传中...</van-loading>
    </div>
  </van-popup>
</template>

<script setup>
import { ref, computed, watch, nextTick, onBeforeUnmount } from 'vue'
import Cropper from 'cropperjs'
import 'cropperjs/dist/cropper.css'
import { showToast } from 'vant'
import * as authApi from '@/api/auth'

const props = defineProps({
  modelValue: Boolean,
  imageFile: File
})

const emit = defineEmits(['update:modelValue', 'success'])

const dialogVisible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

const imageRef = ref(null)
const cropperContainerRef = ref(null)
const imageSrc = ref('')
const previewUrl = ref('')
const uploading = ref(false)
const imageLoaded = ref(false)
let cropper = null

watch(() => props.imageFile, (file) => {
  if (file) {
    loadImage(file)
  }
})

function loadImage(file) {
  imageLoaded.value = false
  if (cropper) {
    cropper.destroy()
    cropper = null
  }

  const reader = new FileReader()
  reader.onload = (e) => {
    imageSrc.value = e.target.result
    previewUrl.value = e.target.result
  }
  reader.onerror = () => {
    showToast({ type: 'fail', message: '读取图片失败' })
  }
  reader.readAsDataURL(file)
}

function onPopupOpened() {
  // Popup 完全打开后初始化裁剪器
  if (imageSrc.value && !cropper) {
    initCropper()
  }
}

function initCropper() {
  nextTick(() => {
    if (!imageRef.value || !cropperContainerRef.value) return

    const image = imageRef.value
    image.style.display = 'block'

    cropper = new Cropper(image, {
      aspectRatio: 1,
      viewMode: 1,
      dragMode: 'move',
      autoCropArea: 1,
      restore: false,
      guides: true,
      center: true,
      highlight: true,
      cropBoxMovable: false,
      cropBoxResizable: false,
      toggleDragModeOnDblclick: false,
      movable: true,
      scalable: true,
      zoomable: true,
      zoomOnWheel: true,
      wheelZoomRatio: 0.1,
      ready() {
        imageLoaded.value = true
        updatePreview()
      },
      crop(event) {
        updatePreview()
      }
    })
  })
}

function updatePreview() {
  if (!cropper) return

  const canvas = cropper.getCroppedCanvas({
    width: 200,
    height: 200,
    imageSmoothingQuality: 'high'
  })

  if (canvas) {
    previewUrl.value = canvas.toDataURL('image/jpeg', 0.9)
  }
}

function rotateLeft() {
  if (!cropper) return
  cropper.rotate(-90)
}

function rotateRight() {
  if (!cropper) return
  cropper.rotate(90)
}

function zoomIn() {
  if (!cropper) return
  cropper.zoom(0.1)
}

function zoomOut() {
  if (!cropper) return
  cropper.zoom(-0.1)
}

function reset() {
  if (!cropper) return
  cropper.reset()
}

async function handleConfirm() {
  if (!cropper) {
    showToast('请先选择图片')
    return
  }

  uploading.value = true
  try {
    const canvas = cropper.getCroppedCanvas({
      width: 400,
      height: 400,
      imageSmoothingQuality: 'high'
    })

    if (!canvas) {
      showToast({ type: 'fail', message: '裁剪失败，请重试' })
      uploading.value = false
      return
    }

    canvas.toBlob((blob) => {
      if (!blob) {
        showToast({ type: 'fail', message: '裁剪失败，请重试' })
        uploading.value = false
        return
      }

      const formData = new FormData()
      const fileName = `avatar_${Date.now()}.jpeg`
      formData.append('avatar', blob, fileName)

      authApi.uploadAvatar(formData).then(() => {
        showToast({ type: 'success', message: '头像上传成功' })
        emit('success')
        handleClose()
      }).catch((error) => {
        console.error('头像上传失败:', error)
        showToast({ type: 'fail', message: error?.message || '头像上传失败' })
      }).finally(() => {
        uploading.value = false
      })
    }, 'image/jpeg', 0.9)
  } catch (error) {
    console.error('裁剪失败:', error)
    showToast({ type: 'fail', message: '裁剪失败，请重试' })
    uploading.value = false
  }
}

function handleClose() {
  if (cropper) {
    cropper.destroy()
    cropper = null
  }
  imageSrc.value = ''
  previewUrl.value = ''
  imageLoaded.value = false
  emit('update:modelValue', false)
}

onBeforeUnmount(() => {
  if (cropper) {
    cropper.destroy()
    cropper = null
  }
})
</script>

<style scoped>
.cropper-popup {
  height: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.cropper-container {
  flex: 1;
  background: #000;
  overflow: hidden;
  position: relative;
}

.cropper-main {
  width: 100%;
  height: 100%;
}

:deep(.cropper-container) {
  width: 100%;
  height: 100%;
}

/* Cropper.js 默认样式修复 */
:deep(.cropper-container) {
  direction: ltr;
  font-size: 0;
  line-height: 0;
  position: relative;
  touch-action: none;
  user-select: none;
}

:deep(.cropper-wrap-box),
:deep(.cropper-canvas),
:deep(.cropper-drag-box),
:deep(.cropper-box),
:deep(.cropper-modal) {
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
}

:deep(.cropper-wrap-box) {
  z-index: 1;
}

:deep(.cropper-canvas) {
  z-index: 2;
}

:deep(.cropper-drag-box) {
  z-index: 3;
}

:deep(.cropper-box) {
  z-index: 4;
}

:deep(.cropper-modal) {
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 5;
}

/* 裁剪框边框 - 使用多重边框确保可见性 */
:deep(.cropper-view-box) {
  display: block;
  height: 100%;
  overflow: hidden;
  width: 100%;
  z-index: 6;
  position: relative;
  /* 使用 box-shadow 创建遮罩效果和边框 */
  box-shadow:
    0 0 0 1px #fff,
    0 0 0 3px #39f,
    0 0 0 4px #fff,
    0 0 0 9999px rgba(0, 0, 0, 0.5);
}

/* 参考线 - 三分线，使用更亮的颜色 */
:deep(.cropper-line) {
  background-color: #fff;
  z-index: 7;
}

/* 九宫格辅助点 - 增大尺寸并使用高对比度颜色 */
:deep(.cropper-point) {
  background-color: #39f;
  height: 10px;
  opacity: 1;
  width: 10px;
  z-index: 8;
  border: 2px solid #fff;
  box-shadow: 0 0 2px rgba(0, 0, 0, 0.5);
}

/* 确保裁剪框边框始终在最上层 */
:deep(.cropper-face) {
  background-color: transparent;
  left: 0;
  top: 0;
  z-index: 9;
}

:deep(.cropper-point.invisible) {
  opacity: 0;
}

:deep(.cropper-bg) {
  background-image: repeating-linear-gradient(
    45deg,
    rgba(255, 255, 255, 0.05) 0.25rem,
    rgba(255, 255, 255, 0.05) 0.25rem 0.5rem,
    transparent 0.5rem,
    transparent 0.5rem 0.75rem,
    rgba(255, 255, 255, 0.05) 0.75rem,
    rgba(255, 255, 255, 0.05) 0.75rem 1rem,
    transparent 1rem,
    transparent 1rem 1.25rem,
    rgba(255, 255, 255, 0.05) 1.25rem,
    rgba(255, 255, 255, 0.05) 1.25rem 1.5rem,
    transparent 1.5rem,
    transparent 1.5rem 1.75rem,
    rgba(255, 255, 255, 0.05) 1.75rem,
    rgba(255, 255, 255, 0.05) 1.75rem 2rem,
    transparent 2rem,
    transparent 2rem 2.25rem,
    rgba(255, 255, 255, 0.05) 2.25rem,
    rgba(255, 255, 255, 0.05) 2.25rem 2.5rem,
    transparent 2.5rem,
    transparent 2.5rem 2.75rem,
    rgba(255, 255, 255, 0.05) 2.75rem,
    rgba(255, 255, 255, 0.05) 2.75rem 3rem,
    transparent 3rem,
    transparent 3rem 3.25rem,
    rgba(255, 255, 255, 0.05) 3.25rem,
    rgba(255, 255, 255, 0.05) 3.25rem 3.5rem,
    transparent 3.5rem,
    transparent 3.5rem 3.75rem,
    rgba(255, 255, 255, 0.05) 3.75rem,
    rgba(255, 255, 255, 0.05) 3.75rem 4rem,
    transparent 4rem,
    transparent 4rem 4.25rem,
    rgba(255, 255, 255, 0.05) 4.25rem,
    rgba(255, 255, 255, 0.05) 4.25rem 4.5rem,
    transparent 4.5rem,
    transparent 4.5rem 4.75rem,
    rgba(255, 255, 255, 0.05) 4.75rem,
    rgba(255, 255, 255, 0.05) 4.75rem 5rem
  );
}

:deep(.cropper-line),
:deep(.cropper-point) {
  display: block;
  pointer-events: none;
  position: absolute;
}

.preview-area {
  flex-shrink: 0;
  background: #f7f8fa;
  padding: 8px 12px;
}

.preview-title {
  font-size: 12px;
  font-weight: 600;
  color: #323233;
  text-align: center;
  margin-bottom: 8px;
}

.preview-list {
  display: flex;
  justify-content: space-around;
  align-items: center;
  gap: 8px;
}

.preview-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.preview-label {
  font-size: 10px;
  color: #646566;
  font-weight: 500;
}

.preview-box-120,
.preview-box-60,
.preview-box-40 {
  display: flex;
  align-items: center;
  justify-content: center;
  background: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 4px;
  border-radius: 8px;
  overflow: hidden;
}

.preview-box-120 {
  width: 80px;
  height: 80px;
}

.preview-box-60 {
  width: 50px;
  height: 50px;
}

.preview-box-40 {
  width: 36px;
  height: 36px;
}

.preview-circular {
  border-radius: 50%;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid #fff;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
  width: 100%;
  height: 100%;
}

.preview-box-120 .preview-circular {
  width: 70px;
  height: 70px;
}

.preview-box-60 .preview-circular {
  width: 42px;
  height: 42px;
}

.preview-box-40 .preview-circular {
  width: 28px;
  height: 28px;
}

.preview-circular img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.cropper-controls {
  flex-shrink: 0;
  padding: 8px 12px;
  background: #fff;
  border-top: 1px solid #ebedf0;
}

.button-row {
  display: flex;
  gap: 6px;
  margin-bottom: 8px;
}

.button-row .van-button {
  flex: 1;
}

.tips {
  margin-top: 8px;
  padding: 6px;
  background: #f7f8fa;
  border-radius: 6px;
  font-size: 10px;
  color: #969799;
  text-align: center;
  line-height: 1.4;
}

.tips p {
  margin: 0;
}

.loading-placeholder {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  z-index: 10;
}

.uploading-loader {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: rgba(0, 0, 0, 0.7);
  padding: 20px;
  border-radius: 8px;
  color: #fff;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  z-index: 9999;
}
</style>
