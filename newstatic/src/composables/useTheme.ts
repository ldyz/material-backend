/**
 * 主题管理 Composable
 * 支持亮色/暗色主题切换和持久化
 */

import { ref, computed, watch, onMounted } from 'vue'
import type { ThemeMode, ThemeConfig } from '@/types/gantt'

/**
 * 主题持久化存储键
 */
const THEME_STORAGE_KEY = 'gantt-theme-preference'
const THEME_SYSTEM_KEY = 'gantt-theme-system'

/**
 * 使用主题管理
 * @param options 配置选项
 * @returns 主题状态和方法
 */
export function useTheme(options: {
  defaultTheme?: ThemeMode
  enableSystem?: boolean
  storageKey?: string
} = {}) {
  const {
    defaultTheme = 'light',
    enableSystem = true,
    storageKey = THEME_STORAGE_KEY
  } = options

  // 当前模式
  const mode = ref<ThemeMode>(defaultTheme)

  // 系统主题偏好
  const systemPreference = ref<'light' | 'dark'>('light')

  // 主题配置
  const config = ref<ThemeConfig>({
    mode: defaultTheme,
    primaryColor: '#409eff',
    fontSize: 'medium'
  })

  /**
   * 实际应用的主题（处理 auto 模式）
   */
  const actualTheme = computed<'light' | 'dark'>(() => {
    if (mode.value === 'auto') {
      return systemPreference.value
    }
    return mode.value
  })

  /**
   * 是否为暗色主题
   */
  const isDark = computed(() => {
    return actualTheme.value === 'dark'
  })

  /**
   * 是否为亮色主题
   */
  const isLight = computed(() => {
    return actualTheme.value === 'light'
  })

  /**
   * 设置主题模式
   */
  function setTheme(newMode: ThemeMode) {
    mode.value = newMode

    // 持久化到 localStorage
    try {
      localStorage.setItem(storageKey, newMode)
    } catch (e) {
      console.warn('Failed to save theme preference:', e)
    }

    // 应用到 DOM
    applyTheme()
  }

  /**
   * 切换主题（在 light/dark 之间切换）
   */
  function toggleTheme() {
    if (mode.value === 'light') {
      setTheme('dark')
    } else if (mode.value === 'dark') {
      setTheme(enableSystem ? 'auto' : 'light')
    } else {
      setTheme('light')
    }
  }

  /**
   * 应用主题到 DOM
   */
  function applyTheme() {
    const root = document.documentElement

    // 移除所有主题类
    root.classList.remove('theme-light', 'theme-dark')

    // 添加当前主题类
    root.classList.add(`theme-${actualTheme.value}`)

    // 设置 data-theme 属性（供 CSS 选择器使用）
    root.setAttribute('data-theme', actualTheme.value)

    // 更新 meta theme-color（移动端浏览器栏颜色）
    updateMetaThemeColor()
  }

  /**
   * 更新 meta theme-color
   */
  function updateMetaThemeColor() {
    let metaTag = document.querySelector('meta[name="theme-color"]') as HTMLMetaElement

    if (!metaTag) {
      metaTag = document.createElement('meta')
      metaTag.name = 'theme-color'
      document.head.appendChild(metaTag)
    }

    const colors = {
      light: '#ffffff',
      dark: '#1a1a1a'
    }

    metaTag.content = colors[actualTheme.value]
  }

  /**
   * 检测系统主题偏好
   */
  function detectSystemPreference() {
    if (!enableSystem) return

    // 检测媒体查询
    const darkModeQuery = window.matchMedia('(prefers-color-scheme: dark)')
    systemPreference.value = darkModeQuery.matches ? 'dark' : 'light'

    // 监听系统主题变化
    darkModeQuery.addEventListener('change', (e) => {
      systemPreference.value = e.matches ? 'dark' : 'light'
      if (mode.value === 'auto') {
        applyTheme()
      }
    })
  }

  /**
   * 从 localStorage 加载主题偏好
   */
  function loadSavedTheme() {
    try {
      const saved = localStorage.getItem(storageKey) as ThemeMode
      if (saved && ['light', 'dark', 'auto'].includes(saved)) {
        mode.value = saved
        return true
      }
    } catch (e) {
      console.warn('Failed to load theme preference:', e)
    }
    return false
  }

  /**
   * 设置主色调
   */
  function setPrimaryColor(color: string) {
    config.value.primaryColor = color
    document.documentElement.style.setProperty('--color-primary', color)
  }

  /**
   * 设置字体大小
   */
  function setFontSize(size: 'small' | 'medium' | 'large') {
    config.value.fontSize = size

    const sizes = {
      small: '13px',
      medium: '14px',
      large: '16px'
    }

    document.documentElement.style.setProperty('--font-size-base', sizes[size])
  }

  /**
   * 重置主题
   */
  function resetTheme() {
    setTheme(defaultTheme)
    setPrimaryColor('#409eff')
    setFontSize('medium')
  }

  // 监听模式变化，自动应用
  watch(mode, () => {
    applyTheme()
  })

  // 初始化
  onMounted(() => {
    // 加载保存的主题
    if (!loadSavedTheme()) {
      mode.value = defaultTheme
    }

    // 检测系统偏好
    detectSystemPreference()

    // 应用主题
    applyTheme()
  })

  return {
    // 状态
    mode,
    config,
    actualTheme,
    isDark,
    isLight,
    systemPreference,

    // 方法
    setTheme,
    toggleTheme,
    setPrimaryColor,
    setFontSize,
    resetTheme,
    applyTheme
  }
}

/**
 * 主题预定义配置
 */
export const THEME_PRESETS = {
  light: {
    name: '亮色',
    mode: 'light' as ThemeMode,
    primaryColor: '#409eff',
    background: '#ffffff',
    text: '#303133'
  },
  dark: {
    name: '暗色',
    mode: 'dark' as ThemeMode,
    primaryColor: '#409eff',
    background: '#1a1a1a',
    text: '#e5eaf3'
  },
  auto: {
    name: '跟随系统',
    mode: 'auto' as ThemeMode,
    primaryColor: '#409eff',
    background: 'auto',
    text: 'auto'
  }
}

/**
 * 获取主题图标
 */
export function getThemeIcon(mode: ThemeMode): string {
  switch (mode) {
    case 'light':
      return '☀️'
    case 'dark':
      return '🌙'
    case 'auto':
      return '🖥️'
    default:
      return '☀️'
  }
}

/**
 * 获取主题 Element Plus 图标
 */
export function getThemeElementPlusIcon(mode: ThemeMode): string {
  switch (mode) {
    case 'light':
      return 'Sunny'
    case 'dark':
      return 'Moon'
    case 'auto':
      return 'Monitor'
    default:
      return 'Sunny'
  }
}
