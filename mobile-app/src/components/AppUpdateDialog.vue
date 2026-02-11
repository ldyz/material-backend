<template>
  <van-dialog
    v-model:show="visible"
    :title="forceUpdate ? '需要更新' : '发现新版本'"
    :show-confirm-button="true"
    :show-cancel-button="!forceUpdate"
    :confirm-button-text="confirmButtonText"
    :cancel-button-text="cancelButtonText"
    :close-on-click-overlay="!forceUpdate"
    @confirm="handleConfirm"
    @cancel="handleCancel"
    class="update-dialog"
  >
    <div class="update-content">
      <div class="update-icon">
        <van-icon name="upgrade" size="48" color="#1989fa" />
      </div>

      <div class="version-info">
        <div class="current-version">当前版本: {{ currentVersion }}</div>
        <div class="latest-version">最新版本: {{ latestVersion }}</div>
      </div>

      <div class="update-message">{{ updateMessage }}</div>

      <div v-if="updateNotes" class="update-notes">
        <div class="notes-title">更新内容:</div>
        <ul class="notes-list">
          <li v-for="(note, index) in updateNotes" :key="index">{{ note }}</li>
        </ul>
      </div>

      <div v-if="isDownloading" class="downloading">
        <van-loading size="24" />
        <span>正在下载...</span>
      </div>
    </div>
  </van-dialog>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useAppUpdate } from '@/composables/useAppUpdate'

const props = defineProps({
  show: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:show'])

const {
  hasUpdate,
  latestVersion,
  currentVersion,
  forceUpdate,
  updateMessage,
  downloadUrl,
  downloadAndInstall,
  openStore,
  remindLater,
  skipVersion
} = useAppUpdate()

const isDownloading = ref(false)
const updateNotes = ref([])

const visible = computed({
  get: () => props.show && hasUpdate.value,
  set: (val) => emit('update:show', val)
})

const confirmButtonText = computed(() => {
  if (isDownloading.value) return '下载中...'
  return forceUpdate.value ? '立即更新' : '立即更新'
})

const cancelButtonText = computed(() => {
  return '稍后提醒'
})

// 监听显示状态，获取更新日志
watch(visible, async (newVal) => {
  if (newVal && latestVersion.value) {
    // 这里可以从 API 获取更新日志
    // updateNotes.value = await getUpdateNotes(latestVersion.value)
    updateNotes.value = [
      '修复已知问题',
      '优化用户体验',
      '提升应用性能'
    ]
  }
})

async function handleConfirm() {
  if (isDownloading.value) return

  isDownloading.value = true

  try {
    // 对于 Android，直接下载 APK
    // 对于 iOS，跳转到 App Store
    await downloadAndInstall()
  } catch (error) {
    console.error('下载更新失败:', error)
  } finally {
    isDownloading.value = false
  }

  // 强制更新不关闭对话框
  if (!forceUpdate.value) {
    visible.value = false
  }
}

function handleCancel() {
  // 强制更新不能取消
  if (forceUpdate.value) return

  remindLater()
  visible.value = false
}

// 跳过此版本
function handleSkip() {
  skipVersion()
  visible.value = false
}
</script>

<style scoped>
.update-dialog :deep(.van-dialog__content) {
  max-height: 60vh;
  overflow-y: auto;
}

.update-content {
  padding: 20px;
  text-align: center;
}

.update-icon {
  margin-bottom: 16px;
}

.version-info {
  margin-bottom: 16px;
  padding: 12px;
  background-color: #f7f8fa;
  border-radius: 8px;
}

.current-version {
  font-size: 14px;
  color: #646566;
  margin-bottom: 4px;
}

.latest-version {
  font-size: 14px;
  color: #323233;
  font-weight: 500;
}

.update-message {
  font-size: 15px;
  color: #323233;
  line-height: 1.5;
  margin-bottom: 16px;
}

.update-notes {
  margin-bottom: 16px;
  text-align: left;
}

.notes-title {
  font-size: 14px;
  font-weight: 500;
  color: #323233;
  margin-bottom: 8px;
}

.notes-list {
  margin: 0;
  padding-left: 20px;
}

.notes-list li {
  font-size: 13px;
  color: #646566;
  line-height: 1.6;
  margin-bottom: 4px;
}

.downloading {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  color: #1989fa;
  font-size: 14px;
}
</style>
