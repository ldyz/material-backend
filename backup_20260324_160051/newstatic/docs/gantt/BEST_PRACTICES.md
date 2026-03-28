# 甘特图组件 - 最佳实践

本文档提供使用甘特图组件的最佳实践和推荐模式。

## 目录

- [性能优化](#性能优化)
- [样式定制](#样式定制)
- [数据处理](#数据处理)
- [错误处理](#错误处理)
- [测试策略](#测试策略)
- [部署建议](#部署建议)

---

## 性能优化

### 1. 大数据集处理

#### 启用虚拟滚动

```vue
<template>
  <GanttChartRefactored
    :enable-virtual-scroll="true"
    :virtual-scroll-threshold="50"
    :schedule-data="scheduleData"
  />
</template>
```

#### 数据分页加载

```typescript
// 对于超大数据集，考虑分页
const PAGE_SIZE = 100

async function loadTasks(page = 1) {
  const response = await progressApi.list({
    project_id: projectId,
    page,
    page_size: PAGE_SIZE
  })

  return response.data
}
```

### 2. 减少重渲染

#### 使用计算属性

```typescript
// ✅ 好的做法
const filteredTasks = computed(() => {
  return tasks.value.filter(t => t.status === 'active')
})

// ❌ 避免
const filteredTasks = ref([])
watch(tasks, () => {
  filteredTasks.value = tasks.value.filter(/*...*/)
}, { deep: true })
```

#### 防抖搜索

```typescript
import { debounce } from '@/utils'

const debouncedSearch = debounce((keyword) => {
  actions.filterTasks(keyword)
}, 300)

watch(searchKeyword, debouncedSearch)
```

### 3. 懒加载组件

```typescript
// 懒加载大型组件
const TaskEditDialog = defineAsyncComponent(() =>
  import('@/components/progress/TaskEditDialog.vue')
)
```

---

## 样式定制

### 1. 使用 CSS 变量覆盖

```vue
<style>
/* 在父组件中覆盖变量 */
.my-gantt-chart {
  --gantt-row-height: 80px;
  --gantt-day-width: 50px;
  --color-primary: #custom-color;

  /* 自定义任务条颜色 */
  --task-bar-completed: #custom-green;
  --task-bar-delayed: #custom-red;
}
</style>

<template>
  <GanttChartRefactored
    class="my-gantt-chart"
    ...
  />
</template>
```

### 2. 深度定制

```vue
<style>
/* 自定义里程碑样式 */
.my-gantt-chart .task-bar.is-milestone {
  background: linear-gradient(135deg, #9b59b6 0%, #8e44ad 100%);
  border-radius: 6px;
}

/* 自定义关键路径样式 */
.my-gantt-chart .task-bar.is-critical {
  box-shadow: 0 0 0 2px #e74c3c;
}
</style>
```

### 3. 主题扩展

```typescript
// 扩展主题配置
const customTheme = {
  mode: 'dark',
  primaryColor: '#ff6b6b',
  fontSize: 'large'
}

// 应用自定义主题
const { setTheme, setPrimaryColor, setFontSize } = useTheme()

setPrimaryColor(customTheme.primaryColor)
setFontSize(customTheme.fontSize)
```

---

## 数据处理

### 1. 数据预处理

```typescript
// 预处理任务数据
function preprocessTasks(rawData) {
  return rawData.map(task => ({
    ...task,
    startDate: new Date(task.start_date),
    endDate: new Date(task.end_date),
    duration: calculateDuration(task.start_date, task.end_date)
  }))
}
```

### 2. 数据验证

```typescript
import type { GanttTask } from '@/types/gantt'

function validateTask(task: any): task is GanttTask {
  return (
    typeof task.id !== 'undefined' &&
    typeof task.name === 'string' &&
    typeof task.start === 'string' &&
    typeof task.end === 'string' &&
    typeof task.progress === 'number'
  )
}
```

### 3. 错误数据处理

```typescript
async function loadTasks() {
  try {
    const response = await progressApi.list()
    const validTasks = response.data.filter(validateTask)
    state.tasks = validTasks
  } catch (error) {
    console.error('加载任务失败:', error)
    ElMessage.error('加载任务失败，请重试')
    state.tasks = []  // 使用空数组作为降级方案
  }
}
```

---

## 错误处理

### 1. API 错误处理

```typescript
import { progressApi } from '@/api'

async function updateTask(taskId, data) {
  try {
    await progressApi.update(taskId, data)
    ElMessage.success('保存成功')
  } catch (error) {
    // 用户友好的错误信息
    const userMessage = error.response?.data?.message || error.message
    ElMessage.error(`保存失败: ${userMessage}`)

    // 记录详细错误
    console.error('任务更新失败:', {
      taskId,
      data,
      error: error.toString()
    })
  }
}
```

### 2. 边界情况处理

```typescript
// 处理空数据
function renderTasks(tasks) {
  if (!tasks || tasks.length === 0) {
    return <EmptyState message="暂无任务" />
  }
  return <TaskList :tasks="tasks" />
}

// 处理极端值
function clampDateRange(start, end) {
  const maxDuration = 365 * 5  // 最多5年
  const duration = diffDays(start, end)

  if (duration > maxDuration) {
    throw new Error('任务工期不能超过5年')
  }
}
```

### 3. 用户提示

```typescript
// 使用状态栏提示
statusBarRef.value?.showStatus('正在保存...', 'loading')

// 成功提示
statusBarRef.value?.showStatus('保存成功', 'success', 2000)

// 错误提示
statusBarRef.value?.showStatus('保存失败', 'error', 2000)
```

---

## 测试策略

### 1. 单元测试示例

```typescript
// composables/useTheme.test.ts
import { describe, it, expect } from 'vitest'
import { useTheme } from '@/composables/useTheme'

describe('useTheme', () => {
  it('should set theme to dark', () => {
    const { mode, setTheme } = useTheme()
    setTheme('dark')
    expect(mode.value).toBe('dark')
  })

  it('should toggle theme', () => {
    const { mode, toggleTheme } = useTheme()
    const original = mode.value
    toggleTheme()
    expect(mode.value).not.toBe(original)
  })
})
```

### 2. 组件测试示例

```vue
<!-- GanttChart.spec.ts -->
<script setup>
import { mount } from '@vue/test-utils'
import { describe, it, expect } from 'vitest'
import GanttChartRefactored from './GanttChartRefactored.vue'

describe('GanttChart', () => {
  it('renders tasks correctly', () => {
    const wrapper = mount(GanttChartRefactored, {
      props: {
        projectId: 1,
        scheduleData: mockScheduleData
      }
    })

    expect(wrapper.find('.task-bar').exists()).toBe(true)
  })

  it('emits task-updated event', async () => {
    const wrapper = mount(GanttChartRefactored, {
      props: {
        projectId: 1,
        scheduleData: mockScheduleData
      }
    })

    await wrapper.vm.handleSaveTask(mockTaskData)
    expect(wrapper.emitted('task-updated')).toBeTruthy()
  })
})
</script>
```

### 3. E2E 测试示例

```typescript
// tests/e2e/gantt.spec.ts
import { test, expect } from '@playwright/test'

test('should create and delete task', async ({ page }) => {
  await page.goto('/gantt')
  await page.click('[aria-label="添加任务"]')

  // 填写任务表单
  await page.fill('[name="task-name"]', '新任务')
  await page.click('button:has-text("确定")')

  // 验证任务创建成功
  await expect(page.locator('.task-bar').first()).toBeVisible()

  // 删除任务
  await page.click('.task-bar.first')
  await page.click('[aria-label="删除任务"]')
  await page.click('button:has-text("确定")')

  // 验证任务已删除
  await expect(page.locator('.task-bar').first()).not.toBeVisible()
})
```

---

## 部署建议

### 1. 构建优化

```bash
# 分析构建产物
npm run build -- --analyze

# 代码分割优化
# vite.config.js
export default {
  build: {
    rollupOptions: {
      output: {
      manualChunks: {
        'vendor': ['vue', 'element-plus'],
        'gantt': ['./src/components/gantt/**'],
        'utils': ['./src/utils/**']
      }
    }
  }
}
```

### 2. 按需加载

```typescript
// 路由级代码分割
const routes = [
  {
    path: '/gantt',
    component: () => import('@/views/GanttView.vue')
  }
]
```

### 3. 资源优化

```typescript
// 图片懒加载
<img :src="imageSrc" loading="lazy" />

// 字体子集化
@font-face {
  font-family: 'Custom Font';
  font-display: swap;
  unicode-range: U+0020-007E;
}
```

### 4. CDN 加速

```html
<!-- 使用 CDN 加速 Vue -->
<script src="https://unpkg.com/vue@3/dist/vue.global.prod.js"></script>
```

---

## 常见问题解决方案

### 问题1: 大数据量卡顿

**解决方案：**

```vue
<GanttChartRefactored
  :enable-virtual-scroll="true"
  :virtual-scroll-threshold="50"
/>
```

### 问题2: 移动端显示异常

**解决方案：**

```css
/* 确保移动端样式优先 */
@media (max-width: 768px) {
  .gantt-chart {
    min-width: 100%;
  }

  .task-table {
    display: none;
  }
}
```

### 问题3: 主题不生效

**解决方案：**

```typescript
// 确保 CSS 变量已加载
import '@/styles/themes/variables.css'

// 设置 data-theme 属性
document.documentElement.setAttribute('data-theme', 'dark')
```

---

## 性能指标

### 目标指标

| 指标 | 目标值 | 说明 |
|------|--------|------|
| 首屏加载 (FCP) | < 1.5s | 首次内容绘制 |
| 最大内容绘制 (LCP) | < 2.5s | 最大内容绘制 |
| 首次输入延迟 (FID) | < 100ms | 首次交互延迟 |
| 累积布局偏移 (CLS) | < 0.1 | 布局稳定性 |
| 1000 任务渲染 | < 500ms | 大数据量性能 |

### 监控方法

```typescript
import { onMounted } from 'vue'

onMounted(() => {
  // 性能监控
  if (window.performance) {
    const timing = window.performance.timing
    console.log('FCP:', timing.domContentLoadedEventEnd - timing.fetchStart)
    console.log('LCP:', timing.loadEventEnd - timing.fetchStart)
  }
})
```

---

## 安全建议

### 1. 输入验证

```typescript
function sanitizeInput(input: string): string {
  // 移除 HTML 标签
  return input.replace(/<[^>]*>/g, '')
}

function validateDateRange(start: string, end: string): boolean {
  const startDate = new Date(start)
  const endDate = new Date(end)
  return startDate < endDate
}
```

### 2. XSS 防护

```vue
<!-- 使用 v-text 而不是 v-html -->
<div v-text="userContent"></div>

<!-- 必须使用 v-html 时进行清理 -->
<div v-html="sanitizedHtml"></div>
```

### 3. CSRF 防护

```typescript
import { getCsrfToken } from '@/api'

async function makeRequest(data: any) {
  const csrfToken = getCsrfToken()

  await fetch('/api/tasks', {
    method: 'POST',
    headers: {
      'X-CSRF-TOKEN': csrfToken
    },
    body: JSON.stringify(data)
  })
}
```

---

## 更多资源

- [API 文档](./API.md)
- [迁移指南](./MIGRATION.md)
- [架构文档](./ARCHITECTURE.md)

---

**更新日期**: 2026-01-31
**版本**: 2.0.0
