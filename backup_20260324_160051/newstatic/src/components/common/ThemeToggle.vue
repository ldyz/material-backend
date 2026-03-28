<template>
  <div class="theme-toggle">
    <!-- 下拉菜单版本 -->
    <el-dropdown
      trigger="click"
      placement="bottom-end"
      @command="handleCommand"
    >
      <button
        class="theme-toggle__button"
        :aria-label="`当前主题: ${currentThemeLabel}, 点击切换主题`"
        type="button"
      >
        <span class="theme-icon">{{ currentThemeIcon }}</span>
        <span v-if="showLabel" class="theme-label">{{ currentThemeLabel }}</span>
      </button>

      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item
            v-for="(preset, key) in THEME_PRESETS"
            :key="key"
            :command="preset.mode"
            :class="{ 'is-active': mode === preset.mode }"
          >
            <span class="theme-option__icon">{{ preset.icon }}</span>
            <span class="theme-option__label">{{ preset.name }}</span>
            <span v-if="mode === preset.mode" class="theme-option__check">✓</span>
          </el-dropdown-item>

          <el-dropdown-item divided :command="'reset'">
            <span class="theme-option__icon">🔄</span>
            <span class="theme-option__label">重置默认</span>
          </el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </div>
</template>

<script setup lang="ts">
import { computed, PropType } from 'vue'
import { useTheme, THEME_PRESETS, getThemeIcon } from '@/composables/useTheme'
import type { ThemeMode } from '@/types/gantt'

interface Props {
  showLabel?: boolean
  defaultTheme?: ThemeMode
}

const props = withDefaults(defineProps<Props>(), {
  showLabel: false,
  defaultTheme: 'light'
})

const emit = defineEmits<{
  themeChange: [mode: ThemeMode]
}>()

const { mode, setTheme, toggleTheme, resetTheme } = useTheme({
  defaultTheme: props.defaultTheme
})

/**
 * 当前主题图标
 */
const currentThemeIcon = computed(() => {
  return getThemeIcon(mode.value)
})

/**
 * 当前主题标签
 */
const currentThemeLabel = computed(() => {
  return THEME_PRESETS[mode.value]?.name || '未知'
})

/**
 * 带图标的预设
 */
const THEME_PRESETS_WITH_ICONS = computed(() => {
  return {
    light: { ...THEME_PRESETS.light, icon: '☀️' },
    dark: { ...THEME_PRESETS.dark, icon: '🌙' },
    auto: { ...THEME_PRESETS.auto, icon: '🖥️' }
  }
})

const THEME_PRESETS = THEME_PRESETS_WITH_ICONS.value

/**
 * 处理下拉菜单命令
 */
function handleCommand(command: string | ThemeMode) {
  if (command === 'reset') {
    resetTheme()
  } else {
    setTheme(command as ThemeMode)
  }
  emit('themeChange', mode.value)
}

// 暴露方法供父组件调用
defineExpose({
  toggleTheme,
  setTheme,
  resetTheme,
  mode
})
</script>

<style scoped>
.theme-toggle {
  display: inline-block;
}

.theme-toggle__button {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border-base);
  border-radius: var(--radius-md);
  color: var(--color-text-primary);
  font-size: var(--font-size-sm);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.theme-toggle__button:hover {
  background: var(--color-bg-tertiary);
  border-color: var(--color-primary);
}

.theme-toggle__button:active {
  transform: scale(0.98);
}

.theme-icon {
  font-size: 18px;
  line-height: 1;
}

.theme-label {
  font-weight: var(--font-weight-medium);
}

/* 下拉选项样式 */
.theme-option__icon {
  margin-right: 8px;
  font-size: 16px;
}

.theme-option__label {
  flex: 1;
}

.theme-option__check {
  color: var(--color-primary);
  font-weight: bold;
}

.el-dropdown-menu__item.is-active {
  background: var(--color-primary-lighter);
  color: var(--color-primary);
}

/* 响应式 */
@media (max-width: 768px) {
  .theme-toggle__button {
    padding: 10px;
  }

  .theme-label {
    display: none;
  }
}
</style>
