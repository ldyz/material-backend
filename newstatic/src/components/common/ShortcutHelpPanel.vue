<template>
  <el-dialog
    v-model="visible"
    title="⌨️ 键盘快捷键"
    :width="700"
    :modal="true"
    :close-on-click-modal="true"
    :close-on-press-escape="true"
    class="shortcut-help-dialog"
    @close="handleClose"
  >
    <div class="shortcut-help">
      <!-- 搜索框 -->
      <div class="shortcut-help__search">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索快捷键..."
          prefix-icon="el-icon-search"
          clearable
        />
      </div>

      <!-- 快捷键分组 -->
      <div class="shortcut-help__groups">
        <div
          v-for="group in filteredGroups"
          :key="group.title"
          class="shortcut-group"
        >
          <h3 class="shortcut-group__title">{{ group.title }}</h3>

          <div class="shortcut-group__items">
            <div
              v-for="shortcut in group.shortcuts"
              :key="shortcut.key"
              class="shortcut-item"
            >
              <div class="shortcut-item__keys">
                <kbd
                  v-for="(key, index) in shortcut.keys"
                  :key="index"
                  class="key-badge"
                >
                  {{ key }}
                </kbd>
              </div>
              <div class="shortcut-item__description">
                {{ shortcut.description }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 提示信息 -->
      <div class="shortcut-help__tips">
        <el-alert
          type="info"
          :closable="false"
          show-icon
        >
          <template #title>
            <span>提示：按 <kbd>?</kbd> 或 <kbd>Ctrl+/</kbd> 可随时打开此帮助面板</span>
          </template>
        </el-alert>
      </div>
    </div>

    <template #footer>
      <span class="dialog-footer">
        <el-checkbox v-model="dontShowAgain">不再显示此提示</el-checkbox>
        <el-button @click="handleClose">关闭</el-button>
        <el-button type="primary" @click="handleClose">知道了</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'

interface Shortcut {
  keys: string[]
  description: string
}

interface ShortcutGroup {
  title: string
  shortcuts: Shortcut[]
}

interface Props {
  modelValue: boolean
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
}>()

const searchKeyword = ref('')
const dontShowAgain = ref(false)

/**
 * 对话框可见性
 */
const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

/**
 * 所有快捷键分组
 */
const allGroups: ShortcutGroup[] = [
  {
    title: '🧭 导航',
    shortcuts: [
      { keys: ['↑'], description: '选择上一个任务' },
      { keys: ['↓'], description: '选择下一个任务' },
      { keys: ['←'], description: '选择前驱任务' },
      { keys: ['→'], description: '选择后继任务' },
      { keys: ['Home'], description: '跳到第一个任务' },
      { keys: ['End'], description: '跳到最后一个任务' },
      { keys: ['Page Up'], description: '向上翻页' },
      { keys: ['Page Down'], description: '向下翻页' }
    ]
  },
  {
    title: '✏️ 编辑',
    shortcuts: [
      { keys: ['Enter'], description: '编辑选中的任务' },
      { keys: ['Delete', '⌫'], description: '删除选中的任务' },
      { keys: ['Escape'], description: '取消操作 / 关闭对话框' },
      { keys: ['Ctrl', 'C'], description: '复制任务' },
      { keys: ['Ctrl', 'V'], description: '粘贴任务' },
      { keys: ['Ctrl', 'Z'], description: '撤销' },
      { keys: ['Ctrl', 'Shift', 'Z'], description: '重做' },
      { keys: ['Ctrl', 'S'], description: '保存更改' }
    ]
  },
  {
    title: '👁️ 视图',
    shortcuts: [
      { keys: ['Ctrl', '+'], description: '放大时间轴' },
      { keys: ['Ctrl', '-'], description: '缩小时间轴' },
      { keys: ['Ctrl', '0'], description: '重置缩放' },
      { keys: ['Alt', 'D'], description: '切换依赖关系显示' },
      { keys: ['Alt', 'P'], description: '切换关键路径显示' },
      { keys: ['Alt', 'B'], description: '切换基线显示' },
      { keys: ['G'], description: '切换分组模式' }
    ]
  },
  {
    title: '📅 日期导航',
    shortcuts: [
      { keys: ['Alt', '↑'], description: '上一周' },
      { keys: ['Alt', '↓'], description: '下一周' },
      { keys: ['Alt', '←'], description: '前一天' },
      { keys: ['Alt', '→'], description: '后一天' },
      { keys: ['Alt', 'H'], description: '回到今天' }
    ]
  },
  {
    title: '🔍 搜索与筛选',
    shortcuts: [
      { keys: ['Ctrl', 'F'], description: '聚焦搜索框' },
      { keys: ['F3'], description: '查找下一个' },
      { keys: ['Shift', 'F3'], description: '查找上一个' },
      { keys: ['Ctrl', 'Shift', 'F'], description: '打开高级筛选' }
    ]
  },
  {
    title: '⚙️ 其他',
    shortcuts: [
      { keys: ['?'], description: '显示/隐藏此帮助' },
      { keys: ['Ctrl', '/'], description: '显示/隐藏此帮助' },
      { keys: ['F11'], description: '全屏模式' },
      { keys: ['Ctrl', 'P'], description: '打印 / 导出PDF' }
    ]
  }
]

/**
 * 过滤后的分组
 */
const filteredGroups = computed(() => {
  if (!searchKeyword.value) {
    return allGroups
  }

  const keyword = searchKeyword.value.toLowerCase()

  return allGroups
    .map(group => ({
      title: group.title,
      shortcuts: group.shortcuts.filter(shortcut => {
        const searchText = shortcut.keys.join(' ') + ' ' + shortcut.description
        return searchText.toLowerCase().includes(keyword)
      })
    }))
    .filter(group => group.shortcuts.length > 0)
})

/**
 * 关闭对话框
 */
function handleClose() {
  visible.value = false

  // 保存"不再显示"偏好
  if (dontShowAgain.value) {
    try {
      localStorage.setItem('gantt-hide-shortcut-help', 'true')
    } catch (e) {
      console.warn('Failed to save shortcut help preference:', e)
    }
  }
}

// 监听对话框关闭，重置状态
watch(visible, (newVal) => {
  if (!newVal) {
    searchKeyword.value = ''
    dontShowAgain.value = false
  }
})
</script>

<style scoped>
.shortcut-help {
  max-height: 60vh;
  overflow-y: auto;
}

.shortcut-help__search {
  margin-bottom: 20px;
}

.shortcut-help__groups {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.shortcut-group {
  background: var(--color-bg-secondary);
  border-radius: var(--radius-md);
  padding: 16px;
}

.shortcut-group__title {
  margin: 0 0 12px 0;
  font-size: var(--font-size-md);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-primary);
  display: flex;
  align-items: center;
  gap: 8px;
}

.shortcut-group__items {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.shortcut-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 12px;
  background: var(--color-bg-primary);
  border-radius: var(--radius-sm);
  transition: background var(--transition-fast);
}

.shortcut-item:hover {
  background: var(--color-bg-tertiary);
}

.shortcut-item__keys {
  display: flex;
  gap: 4px;
  min-width: 120px;
}

.key-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 28px;
  height: 28px;
  padding: 0 8px;
  background: linear-gradient(180deg, #fafafa 0%, #e8e8e8 100%);
  border: 1px solid #c4c4c4;
  border-radius: var(--radius-sm);
  box-shadow: 0 2px 0 #888, inset 0 1px 0 #fff;
  font-family: monospace;
  font-size: 12px;
  font-weight: var(--font-weight-medium);
  color: var(--color-text-primary);
  text-align: center;
}

.key-badge:only-of-type {
  min-width: 60px;
}

.shortcut-item__description {
  flex: 1;
  font-size: var(--font-size-sm);
  color: var(--color-text-regular);
}

.shortcut-help__tips {
  margin-top: 20px;
}

.dialog-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

/* 暗色主题下的按键样式 */
[data-theme="dark"] .key-badge {
  background: linear-gradient(180deg, #4a4a4a 0%, #2a2a2a 100%);
  border-color: #5a5a5a;
  box-shadow: 0 2px 0 #1a1a1a, inset 0 1px 0 #6a6a6a;
  color: #e5eaf3;
}

/* 响应式 */
@media (max-width: 768px) {
  .shortcut-help__groups {
    gap: 16px;
  }

  .shortcut-group {
    padding: 12px;
  }

  .shortcut-item {
    flex-wrap: wrap;
  }

  .shortcut-item__keys {
    min-width: auto;
    width: 100%;
  }

  .key-badge {
    flex: 1;
  }

  .dialog-footer {
    flex-direction: column;
    gap: 12px;
    align-items: stretch;
  }
}
</style>

<style>
.shortcut-help-dialog .el-dialog__body {
  padding: 16px 20px;
}
</style>
