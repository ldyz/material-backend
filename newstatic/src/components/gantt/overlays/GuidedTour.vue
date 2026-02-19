<template>
  <div class="guided-tour">
    <!-- Tour Trigger Button (in toolbar) -->
    <el-button
      v-if="showTrigger"
      :icon="QuestionFilled"
      circle
      size="small"
      class="tour-trigger"
      @click="startTour"
    >
    </el-button>

    <!-- Vue Tour Component -->
    <v-tour
      ref="tourRef"
      name="ganttTour"
      :steps="currentSteps"
      :options="tourOptions"
      :callbacks="tourCallbacks"
    />

    <!-- Tour Complete Dialog -->
    <el-dialog
      v-model="showCompleteDialog"
      :title="$t('gantt.tour.complete.title')"
      width="400px"
      :close-on-click-modal="false"
    >
      <div class="tour-complete-content">
        <el-result
          icon="success"
          :title="$t('gantt.tour.complete.title')"
          :sub-title="$t('gantt.tour.complete.message')"
        />
        <div class="tour-actions">
          <el-button @click="restartTour">
            {{ $t('gantt.tour.actions.restart') }}
          </el-button>
          <el-button type="primary" @click="tryAdvancedTour">
            {{ $t('gantt.tour.actions.advanced') }}
          </el-button>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { useTour } from '@vueuse/integrations/useTour'
import { QuestionFilled, CircleCheck, Close } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  getTourSteps,
  tourConfig,
  getCompletedTours,
  markTourCompleted,
  isTourCompleted,
  resetTourProgress
} from '@/utils/tourSteps'

const props = defineProps({
  /**
   * Type of tour to show
   */
  tourType: {
    type: String,
    default: 'basic',
    validator: (value) => ['basic', 'advanced', 'features'].includes(value)
  },
  /**
   * Auto-start tour on mount
   */
  autoStart: {
    type: Boolean,
    default: false
  },
  /**
   * Show tour trigger button
   */
  showTrigger: {
    type: Boolean,
    default: true
  },
  /**
   * Delay before auto-start (ms)
   */
  startDelay: {
    type: Number,
    default: 1000
  }
})

const emit = defineEmits(['tour-start', 'tour-complete', 'tour-skip', 'step-change'])

const { t } = useI18n()
const tourRef = ref(null)
const showCompleteDialog = ref(false)
const currentTourType = ref(props.tourType)

// Current tour steps
const currentSteps = computed(() => getTourSteps(currentTourType.value))

// Tour options
const tourOptions = computed(() => ({
  highlight: true,
  useKeyboardNavigation: true,
  labels: {
    buttonSkip: t('gantt.tour.labels.skip'),
    buttonPrevious: t('gantt.tour.labels.previous'),
    buttonNext: t('gantt.tour.labels.next'),
    buttonStop: t('gantt.tour.labels.finish')
  },
  enableScrolling: true,
  scrollPadding: 100,
  highlightClass: 'gantt-tour-highlight',
  modalClass: 'gantt-tour-modal'
}))

// Tour callbacks
const tourCallbacks = computed(() => ({
  onStart: () => {
    console.log('Tour started')
    emit('tour-start', currentTourType.value)
    ElMessage.info(t('gantt.tour.messages.started'))
  },
  onPreviousStep: (currentStep) => {
    console.log('Previous step:', currentStep)
  },
  onNextStep: (currentStep) => {
    console.log('Next step:', currentStep)
    emit('step-change', currentStep)
  },
  onStop: () => {
    console.log('Tour stopped')
    markTourCompleted(currentTourType.value)
    showCompleteDialog.value = true
    emit('tour-complete', currentTourType.value)
  },
  onFinish: () => {
    console.log('Tour finished')
    markTourCompleted(currentTourType.value)
    showCompleteDialog.value = true
    emit('tour-complete', currentTourType.value)
  },
  onSkip: () => {
    console.log('Tour skipped')
    emit('tour-skip', currentTourType.value)
    ElMessage.info(t('gantt.tour.messages.skipped'))
  }
}))

/**
 * Start the tour
 */
function startTour(tourType = props.tourType) {
  currentTourType.value = tourType

  // Check if tour was already completed
  if (isTourCompleted(tourType) && tourType !== 'basic') {
    ElMessageBox.confirm(
      t('gantt.tour.confirm.alreadyCompleted'),
      t('gantt.tour.confirm.title'),
      {
        confirmButtonText: t('gantt.tour.confirm.resume'),
        cancelButtonText: t('gantt.tour.confirm.cancel'),
        type: 'info'
      }
    )
      .then(() => {
        // Reset and start anyway
        resetTourProgress(tourType)
        launchTour()
      })
      .catch(() => {
        // User cancelled
      })
    return
  }

  launchTour()
}

/**
 * Launch the tour
 */
function launchTour() {
  nextTick(() => {
    if (tourRef.value) {
      tourRef.value.start()
    }
  })
}

/**
 * Restart current tour
 */
function restartTour() {
  showCompleteDialog.value = false
  resetTourProgress(currentTourType.value)
  launchTour()
}

/**
 * Try advanced tour
 */
function tryAdvancedTour() {
  showCompleteDialog.value = false
  if (currentTourType.value !== 'advanced') {
    startTour('advanced')
  }
}

/**
 * Skip tour
 */
function skipTour() {
  if (tourRef.value) {
    tourRef.value.skip()
  }
}

// Watch for tour type changes
watch(() => props.tourType, (newType) => {
  currentTourType.value = newType
})

// Auto-start on mount if configured
onMounted(() => {
  if (props.autoStart && !isTourCompleted(props.tourType)) {
    setTimeout(() => {
      startTour()
    }, props.startDelay)
  }
})

// Expose methods for external use
defineExpose({
  startTour,
  skipTour,
  restartTour,
  isTourCompleted
})
</script>

<script>
export default {
  name: 'GuidedTour'
}
</script>

<style scoped>
.guided-tour {
  display: inline-block;
}

.tour-trigger {
  margin-left: 8px;
}

.tour-complete-content {
  text-align: center;
}

.tour-actions {
  display: flex;
  justify-content: center;
  gap: 12px;
  margin-top: 20px;
}

/* Tour-specific styles */
:deep(.gantt-tour-highlight) {
  position: relative;
  z-index: 9998 !important;
  box-shadow: 0 0 0 9999px rgba(0, 0, 0, 0.5);
  border-radius: 4px;
  transition: all 0.3s ease;
}

:deep(.gantt-tour-highlight::before) {
  content: '';
  position: absolute;
  inset: -4px;
  border: 2px solid #4CAF50;
  border-radius: 6px;
  pointer-events: none;
  animation: pulse 2s infinite;
}

:deep(.gantt-tour-modal) {
  background: rgba(0, 0, 0, 0.7);
  backdrop-filter: blur(2px);
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

/* Tour step styles */
:deep(.v-step) {
  background: white;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
  max-width: 350px;
  z-index: 9999 !important;
}

:deep(.v-step__header) {
  font-size: 18px;
  font-weight: 600;
  margin-bottom: 12px;
  color: #303133;
}

:deep(.v-step__content) {
  font-size: 14px;
  color: #606266;
  line-height: 1.6;
  margin-bottom: 16px;
}

:deep(.v-step__buttons) {
  display: flex;
  justify-content: space-between;
  gap: 8px;
}

:deep(.v-step__button) {
  padding: 8px 16px;
  border-radius: 4px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
}

:deep(.v-step__button--primary) {
  background: #4CAF50;
  color: white;
  border: none;
}

:deep(.v-step__button--primary:hover) {
  background: #45a049;
}

:deep(.v-step__button--secondary) {
  background: #f0f0f0;
  color: #606266;
  border: 1px solid #dcdfe6;
}

:deep(.v-step__button--secondary:hover) {
  background: #e0e0e0;
}

/* Tour progress indicator */
:deep(.v-step__progress) {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
  font-size: 12px;
  color: #909399;
}

:deep(.v-step__progress-bar) {
  flex: 1;
  height: 4px;
  background: #f0f0f0;
  border-radius: 2px;
  overflow: hidden;
}

:deep(.v-step__progress-fill) {
  height: 100%;
  background: linear-gradient(90deg, #4CAF50, #8BC34A);
  transition: width 0.3s ease;
}

/* Tour arrow */
:deep(.v-step__arrow) {
  width: 0;
  height: 0;
  border-style: solid;
  position: absolute;
  margin: 6px;
}

:deep(.v-step__arrow--top) {
  border-width: 0 6px 6px 6px;
  border-color: transparent transparent white transparent;
  bottom: 100%;
  left: 50%;
  transform: translateX(-50%);
}

:deep(.v-step__arrow--bottom) {
  border-width: 6px 6px 0 6px;
  border-color: white transparent transparent transparent;
  top: 100%;
  left: 50%;
  transform: translateX(-50%);
}

:deep(.v-step__arrow--left) {
  border-width: 6px 6px 6px 0;
  border-color: transparent white transparent transparent;
  right: 100%;
  top: 50%;
  transform: translateY(-50%);
}

:deep(.v-step__arrow--right) {
  border-width: 6px 0 6px 6px;
  border-color: transparent transparent transparent white;
  left: 100%;
  top: 50%;
  transform: translateY(-50%);
}

/* Responsive adjustments */
@media (max-width: 768px) {
  :deep(.v-step) {
    max-width: 280px;
    font-size: 13px;
  }

  :deep(.v-step__header) {
    font-size: 16px;
  }

  :deep(.v-step__button) {
    padding: 6px 12px;
    font-size: 13px;
  }
}
</style>
