# API 文档

本文档详细说明甘特图组件的 API 接口。

## 目录

- [组件 Props](#组件-props)
- [组件 Events](#组件-events)
- [组件 Slots](#组件-slots)
- [Composables API](#composables-api)
- [类型定义](#类型定义)

---

## 组件 Props

### GanttChartRefactored

主甘特图组件，包含所有功能。

```typescript
interface GanttChartProps {
  /** 项目ID（必填） */
  projectId: number | string

  /** 项目名称 */
  projectName?: string

  /** 进度数据对象 */
  scheduleData?: ScheduleData

  /** 默认视图模式 */
  defaultViewMode?: 'day' | 'week' | 'month' | 'quarter'

  /** 默认主题 */
  defaultTheme?: 'light' | 'dark' | 'auto'

  /** 是否启用虚拟滚动 */
  enableVirtualScroll?: boolean

  /** 虚拟滚动阈值（任务数） */
  virtualScrollThreshold?: number
}
```

**示例：**

```vue
<GanttChartRefactored
  :project-id="123"
  project-name="我的项目"
  :schedule-data="scheduleData"
  default-view-mode="week"
  default-theme="dark"
  :enable-virtual-scroll="true"
  :virtual-scroll-threshold="100"
  @task-updated="handleUpdate"
/>
```

---

## 组件 Events

### 任务事件

| 事件名 | 参数 | 说明 |
|--------|------|------|
| `task-updated` | `(task: GanttTask)` | 任务数据更新时触发 |
| `task-selected` | `(task: GanttTask)` | 任务被选中时触发 |

### 视图事件

| 事件名 | 参数 | 说明 |
|--------|------|------|
| `view-mode-change` | `(mode: ViewMode)` | 视图模式改变时触发 |
| `theme-change` | `(theme: ThemeMode)` | 主题改变时触发 |

### 编辑事件

| 事件名 | 参数 | 说明 |
|--------|------|------|
| `task-created` | `(task: GanttTask)` | 任务创建时触发 |
| `task-deleted` | `(taskId: string \| number)` | 任务删除时触发 |
| `dependency-created` | `(dep: GanttTaskDependency)` | 依赖关系创建时触发 |

**示例：**

```vue
<GanttChartRefactored
  @task-updated="(task) => console.log('更新:', task)"
  @task-selected="(task) => console.log('选中:', task)"
  @theme-change="(theme) => console.log('主题:', theme)"
/>
```

---

## 组件 Slots

### GanttToolbar Slots

工具栏组件支持插槽扩展。

```vue
<GanttToolbar>
  <!-- 左侧操作区 -->
  <template #actions>
    <el-button>自定义按钮</el-button>
  </template>

  <!-- 右侧工具区 -->
  <template #tools>
    <el-icon><Search /></el-icon>
  </template>
</GanttToolbar>
```

### TaskBar Slots

任务条组件支持插槽自定义。

```vue
<TaskBar>
  <!-- 任务条左侧内容 -->
  <template #left>
    <el-icon><Document /></el-icon>
  </template>

  <!-- 任务条右侧内容 -->
  <template #right>
    <el-tag>重要</el-tag>
  </template>
</TaskBar>
```

---

## Composables API

### useTheme

主题管理 Composable。

```typescript
import { useTheme } from '@/composables'

const {
  mode,              // Ref<ThemeMode> - 当前主题模式
  config,            // Ref<ThemeConfig> - 主题配置
  actualTheme,       // ComputedRef<'light' | 'dark'> - 实际应用的主题
  isDark,            // ComputedRef<boolean> - 是否为暗色主题
  isLight,           // ComputedRef<boolean> - 是否为亮色主题
  systemPreference,  // Ref<'light' | 'dark'> - 系统主题偏好
  setTheme,          // (mode: ThemeMode) => void - 设置主题
  toggleTheme,       // () => void - 切换主题
  setPrimaryColor,   // (color: string) => void - 设置主色调
  setFontSize,       // (size: 'small' | 'medium' | 'large') => void - 设置字体大小
  resetTheme         // () => void - 重置主题
} = useTheme({
  defaultTheme: 'light',
  enableSystem: true,
  storageKey: 'gantt-theme-preference'
})
```

**示例：**

```typescript
// 切换主题
const { toggleTheme } = useTheme()
toggleTheme()

// 设置主色调
const { setPrimaryColor } = useTheme()
setPrimaryColor('#ff6b6b')
```

### useBreakpoint

响应式断点检测 Composable。

```typescript
import { useBreakpoint } from '@/composables'

const {
  current,      // Ref<Breakpoint> - 当前断点
  width,        // Ref<number> - 屏幕宽度
  height,       // Ref<number> - 屏幕高度
  isMobile,     // ComputedRef<boolean> - 是否为移动端
  isTablet,     // ComputedRef<boolean> - 是否为平板
  isDesktop,    // ComputedRef<boolean> - 是否为桌面端
  isTouch,      // ComputedRef<boolean> - 是否支持触摸
  isHighDPI,    // ComputedRef<boolean> - 是否为高DPI
  isGreater,    // (bp: Breakpoint) => boolean - 检测是否大于断点
  isLess,       // (bp: Breakpoint) => boolean - 检测是否小于断点
  isBetween,    // (min: Breakpoint, max: Breakpoint) => boolean - 检测是否在范围内
  is            // (bp: Breakpoint) => boolean - 检测是否匹配断点
} = useBreakpoint()
```

**示例：**

```typescript
const { isMobile, isTablet } = useBreakpoint()

if (isMobile.value) {
  // 移动端逻辑
}
```

### useVirtualScroll

虚拟滚动 Composable。

```typescript
import { useVirtualScroll } from '@/composables'

const {
  visibleRange,    // ComputedRef<VisibleRange> - 可见范围
  visibleItems,    // ComputedRef<T[]> - 可见项
  totalHeight,     // ComputedRef<number> - 总高度
  offsetY,         // ComputedRef<number> - 偏移量
  scrollTop,       // Ref<number> - 滚动位置
  scrollToItem,    // (index: number) => void - 滚动到指定项
  scrollToPosition,// (position: number) => void - 滚动到指定位置
  getItemSize,     // (index: number) => number - 获取项高度
  getItemOffset,   // (index: number) => number - 获取项偏移
  refresh          // () => void - 刷新布局
} = useVirtualScroll({
  items: ref(taskList),
  itemHeight: 60,
  containerHeight: 600,
  overscan: 3
})
```

### useGanttDrag

甘特图拖拽 Composable。

```typescript
import { useGanttDrag } from '@/composables/useGanttDrag'

const {
  isDragging,        // Ref<boolean> - 是否正在拖拽
  dragMode,          // Ref<DragMode> - 拖拽模式
  draggedTask,       // Ref<GanttTask | null> - 被拖拽的任务
  previewTask,       // Ref<GanttTask | null> - 预览任务
  tooltipPosition,   // Ref<{x, y}> - 提示框位置
  tooltipVisible,    // Ref<boolean> - 提示框可见性
  tooltipText,       // ComputedRef<string> - 提示文本
  canResizeLeft,     // ComputedRef<boolean> - 是否可调整左边缘
  canResizeRight,    // ComputedRef<boolean> - 是否可调整右边缘
  startDrag,         // (event, task, taskBar) => void - 开始拖拽
  cancelDrag,        // () => void - 取消拖拽
  getTaskBarClass,   // (task) => string - 获取任务条样式类
  getTaskBarStyle,   // (task) => object - 获取任务条样式
  detectDragMode     // (event, taskBar) => DragMode - 检测拖拽模式
} = useGanttDrag({
  dayWidth: ref(40),
  timelineDays: ref(timelineDays),
  onDragEnd: handleDragEnd,
  onDragChange: handleDragChange,
  enableRAF: true,      // 启用 requestAnimationFrame 优化
  throttleMs: 16        // 节流延迟（~60fps）
})
```

### useGanttKeyboard

键盘快捷键 Composable。

```typescript
import { useGanttKeyboard } from '@/composables'

const {
  showShortcutHelp,  // Ref<boolean> - 帮助面板可见性
  shortcuts,         // KeyboardShortcut[] - 快捷键列表
  getShortcutText,   // (shortcut) => string - 获取快捷键文本
  toggleShortcutHelp,// () => void - 切换帮助面板
  navigate           // (direction: NavigationDirection) => void - 键盘导航
} = useGanttKeyboard({
  tasks: ref(tasks),
  selectedTask: ref(selectedTask),
  onTaskSelect: (task) => {},
  onTaskEdit: (task) => {},
  onTaskDelete: (task) => {},
  onCopy: () => {},
  onPaste: () => {},
  onUndo: () => {},
  onRedo: () => {},
  onSave: () => {},
  onZoomIn: () => {},
  onZoomOut: () => {}
})
```

### useGanttSelection

选择管理 Composable。

```typescript
import { useGanttSelection } from '@/composables'

const {
  selectedTaskIds,   // Ref<Set> - 已选中任务ID
  selectedTasks,     // ComputedRef<GanttTask[]> - 已选中任务
  firstSelectedTask,  // ComputedRef<GanttTask | null> - 第一个选中任务
  lastSelectedTaskId, // Ref - 最后选中任务ID
  anchorTaskId,      // Ref - 锚点任务ID（范围选择）
  isSelected,         // (taskId) => boolean - 是否选中
  selectTask,         // (taskId, addToSelection) => void - 选择任务
  selectRange,        // (startId, endId) => void - 选择范围
  selectAll,          // () => void - 全选
  clearSelection,     // () => void - 清除选择
  invertSelection,    // () => void - 反选
  handleClick,        // (taskId, event) => void - 处理点击
  getSelectionStats,  // () => object - 获取统计
  getSelectionTimeRange // () => object - 获取时间范围
} = useGanttSelection({
  tasks: ref(tasks),
  mode: ref('single'),  // 'single' | 'multiple' | 'range'
  onSelectionChange: (selectedTasks) => {}
})
```

### useTouchGestures

触摸手势 Composable。

```typescript
import { useTouchGestures } from '@/composables'

const {
  bind,    // () => void - 绑定事件
  unbind   // () => void - 解绑事件
} = useTouchGestures({
  elementRef: ref(element),
  handlers: {
    onSwipeLeft: () => console.log('左滑'),
    onSwipeRight: () => console.log('右滑'),
    onPinch: (scale) => console.log('缩放:', scale)
  },
  swipeThreshold: 50,
  longPressDelay: 500
})
```

---

## 类型定义

### GanttTask

```typescript
interface GanttTask {
  id: number | string
  name: string
  start: string                // ISO 日期字符串 YYYY-MM-DD
  end: string                  // ISO 日期字符串 YYYY-MM-DD
  startDate: Date              // JavaScript Date 对象
  endDate: Date                // JavaScript Date 对象
  duration: number             // 工期（天数）
  progress: number             // 进度百分比 0-100
  status: TaskStatus
  priority: TaskPriority
  is_critical: boolean         // 是否关键路径
  is_milestone: boolean        // 是否里程碑
  parent_id?: number | string | null
  sort_order: number
  predecessors: GanttTaskDependency[]
  successors: GanttTaskDependency[]
  resources: GanttResource[]
  description?: string
  assignee?: string
  tags?: string[]
  color?: string
}
```

### TaskStatus

```typescript
enum TaskStatus {
  NOT_STARTED = 'not_started',
  IN_PROGRESS = 'in_progress',
  COMPLETED = 'completed',
  DELAYED = 'delayed'
}
```

### TaskPriority

```typescript
enum TaskPriority {
  LOW = 'low',
  MEDIUM = 'medium',
  HIGH = 'high'
}
```

### ViewMode

```typescript
enum ViewMode {
  DAY = 'day',
  WEEK = 'week',
  MONTH = 'month',
  QUARTER = 'quarter'
}
```

### GanttTaskDependency

```typescript
interface GanttTaskDependency {
  id: number | string
  predecessor_id: number | string
  successor_id: number | string
  type: DependencyType
  lag?: number                 // 延滞天数（负数表示提前）
}
```

### DependencyType

```typescript
enum DependencyType {
  FINISH_TO_START = 'finish_to_start',  // FS
  START_TO_START = 'start_to_start',    // SS
  FINISH_TO_FINISH = 'finish_to_finish',// FF
  START_TO_FINISH = 'start_to_finish'   // SF
}
```

---

## 工具函数

### validateDependencies

验证依赖关系是否有效。

```typescript
import { validateDependencies } from '@/utils/dependencyGraph'

const result = validateDependencies(tasks)

if (!result.valid) {
  console.error('循环依赖:', result.circularPath)
}
```

### getDownstreamTasks

获取任务的所有下游任务。

```typescript
import { getDownstreamTasks } from '@/utils/dependencyGraph'

const downstreamTasks = getDownstreamTasks(tasks, taskId)
console.log('下游任务:', Array.from(downstreamTasks))
```

### formatDate

格式化日期。

```typescript
import { formatDate } from '@/utils/dateFormat'

const dateStr = formatDate(new Date())  // "2026-01-31"
```

---

## 样式定制

### 使用 CSS 变量

```vue
<style>
/* 覆盖默认变量 */
.gantt-chart {
  --gantt-row-height: 80px;
  --gantt-day-width: 50px;
  --task-bar-completed: #85ce61;
}
</style>
```

### 深度定制

```vue
<style>
/* 自定义任务条样式 */
.gantt-chart .task-bar {
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

/* 自定义里程碑样式 */
.gantt-chart .task-bar.is-milestone {
  background: linear-gradient(135deg, #f39c12 0%, #e67e22 100%);
}
</style>
```

---

## 更多信息

- [完整示例](./examples/)
- [最佳实践](./BEST_PRACTICES.md)
- [性能优化](./PERFORMANCE.md)
