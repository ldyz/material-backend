# 迁移指南

本文档帮助您从旧版甘特图组件迁移到现代化重构版本。

## 目录

- [迁移概述](#迁移概述)
- [破坏性变更](#破坏性变更)
- [分步迁移](#分步迁移)
- [常见问题](#常见问题)

---

## 迁移概述

### 主要变化

| 方面 | 旧版本 | 新版本 |
|------|--------|--------|
| 文件扩展名 | `.js` | `.ts` / `.vue` |
| 类型系统 | JavaScript | TypeScript |
| 状态管理 | Options API | Composition API + Store |
| 组件结构 | 单文件 1300+ 行 | 模块化 50+ 文件 |
| 性能 | 无虚拟滚动 | 虚拟滚动 + 节流 |
| 响应式 | 固定宽度 | 完全响应式 |
| 主题 | 硬编码 | CSS 变量 + 动态切换 |
| 无障碍 | 部分支持 | WCAG 2.1 AA 级别 |

### 迁移收益

✅ **更好的开发体验** - TypeScript 类型提示
✅ **更高的性能** - 虚拟滚动支持 10000+ 任务
✅ **更广的兼容性** - 移动端和桌面端
✅ **更易维护** - 模块化组件架构
✅ **更好的用户体验** - 主题、快捷键、手势

---

## 破坏性变更

### 1. 组件导入路径

**旧版本：**

```javascript
import GanttChart from '@/components/progress/GanttChart.vue'
```

**新版本：**

```typescript
// 选项1: 导入重构版本
import GanttChartRefactored from '@/components/progress/GanttChartRefactored.vue'

// 选项2: 从索引导入
import { GanttChartRefactored as GanttChart } from '@/components/progress'

// 选项3: 使用类型导入
import type { GanttTask, TaskStatus } from '@/types/gantt'
```

### 2. Props 类型

**旧版本（无类型）：**

```javascript
props: {
  projectId: Number,
  scheduleData: Object
}
```

**新版本（TypeScript）：**

```typescript
interface Props {
  projectId: number | string
  scheduleData: ScheduleData
  defaultViewMode?: ViewMode
  defaultTheme?: ThemeMode
}

const props = defineProps<Props>()
```

### 3. Store 访问方式

**旧版本：**

```javascript
import { ganttStore } from '@/stores/ganttStore'

const state = ganttStore.state
const getters = ganttStore.getters
const actions = ganttStore.actions
```

**新版本（兼容旧方式，推荐新方式）：**

```typescript
// 旧方式仍然可用
const store = ganttStore
const { state, getters, actions } = store

// 新方式：使用响应式解构
const loading = computed(() => state.loading)
const filteredTasks = computed(() => state.filteredTasks)
```

### 4. 事件处理

**旧版本：**

```javascript
this.$emit('task-updated', task)
```

**新版本：**

```typescript
// 在 setup 中
const emit = defineEmits<{
  taskUpdated: [task: GanttTask]
  taskSelected: [task: GanttTask]
}>()

emit('taskUpdated', task)
```

### 5. Composable 使用

**旧版本（直接在组件中）：**

```javascript
// 拖拽逻辑写在组件内
handleMouseDown(event) {
  // 拖拽实现...
}
```

**新版本（使用 composable）：**

```typescript
import { useGanttDrag } from '@/composables/useGanttDrag'

const {
  isDragging,
  startDrag,
  cancelDrag
} = useGanttDrag({
  dayWidth: computed(() => state.dayWidth),
  onDragEnd: handleDragEnd
})
```

---

## 分步迁移

### 第 1 步：安装类型定义

```bash
# 安装 TypeScript 类型依赖
npm install --save-dev @types/node @types/lodash-es

# 如果需要更新 tsconfig.json
npm install --save-dev typescript
```

### 第 2 步：更新导入语句

```vue
<!-- 旧版本 -->
<script>
import GanttChart from '@/components/progress/GanttChart.vue'
export default {
  components: { GanttChart }
}
</script>

<!-- 新版本 -->
<script setup lang="ts">
import { GanttChartRefactored } from '@/components/progress'
</script>

<template>
  <GanttChartRefactored
    :project-id="projectId"
    @task-updated="handleUpdate"
  />
</template>
```

### 第 3 步：更新组件 Props

```typescript
// 旧版本
props: {
  projectId: {
    type: Number,
    required: true
  }
}

// 新版本
interface Props {
  projectId: number
  projectName?: string
  scheduleData?: ScheduleData
}

const props = withDefaults(defineProps<Props>(), {
  projectName: '未命名项目',
  scheduleData: () => ({})
})
```

### 第 4 步：迁移事件处理

```typescript
// 旧版本
this.$emit('task-updated', task)

// 新版本
const emit = defineEmits<{
  taskUpdated: [task: GanttTask]
}>()

emit('taskUpdated', task)
```

### 第 5 步：使用新的 Composables

```typescript
// 旧版本：手写拖拽逻辑
// handleMouseDown(event) { ... }

// 新版本：使用 composable
import { useGanttDrag } from '@/composables/useGanttDrag'

const {
  isDragging,
  startDrag,
  cancelDrag
} = useGanttDrag({
  dayWidth: ref(40),
  timelineDays: ref(timelineDays),
  onDragEnd: handleDragEnd,
  enableRAF: true  // 性能优化
})
```

### 第 6 步：添加主题支持

```typescript
// 添加主题切换
import { ThemeToggle } from '@/components/common'

<template>
  <ThemeToggle
    :default-theme="'light'"
    @theme-change="handleThemeChange"
  />
</template>
```

### 第 7 步：添加响应式检测

```typescript
// 添加响应式检测
import { useBreakpoint } from '@/composables'

const { isMobile, isTablet } = useBreakpoint()

// 根据设备类型调整显示
const showMobileView = computed(() => isMobile.value)
```

---

## 常见问题

### Q1: 如何处理自定义样式？

**A:** 使用 CSS 变量覆盖默认样式：

```vue
<style>
.gantt-chart {
  --gantt-row-height: 80px;
  --gantt-day-width: 50px;
  --task-bar-completed: #85ce61;
}
</style>
```

### Q2: 如何禁用虚拟滚动？

**A:** 设置 `enableVirtualScroll` 为 `false` 或调整阈值：

```vue
<GanttChartRefactored
  :enable-virtual-scroll="false"
  :virtual-scroll-threshold="1000"  <!-- 提高阈值 -->
/>
```

### Q3: 如何自定义键盘快捷键？

**A:** 修改 `useGanttKeyboard` 的 `handlers` 参数：

```typescript
const keyboardHandlers = {
  onTaskEdit: (task) => myCustomEdit(task),
  onSave: () => myCustomSave()
}

useGanttKeyboard({
  tasks,
  selectedTask,
  ...keyboardHandlers
})
```

### Q4: 如何集成现有状态管理？

**A:** 新组件完全兼容现有 store，直接使用即可：

```typescript
import { ganttStore } from '@/stores/ganttStore'

const { state, actions } = ganttStore
```

### Q5: 移动端性能如何优化？

**A:** 新组件已内置优化：
- 虚拟滚动自动启用
- 触摸手势使用 RAF 节流
- 图片懒加载
- 代码分割

### Q6: 如何处理大量任务数据？

**A:** 新版本支持 10000+ 任务：

```typescript
// 虚拟滚动自动启用
const useVirtualScroll = computed(() => {
  return tasks.length > 50 || isMobile.value
})
```

### Q7: 如何添加自定义主题？

**A:** 扩展 CSS 变量：

```css
/* src/styles/themes/custom.css */
:root {
  --color-primary: #your-color;
  --gantt-row-height: 80px;
}

/* 然后在组件中导入 */
import '@/styles/themes/custom.css'
```

---

## 快速参考

### 文件重命名对照

| 旧文件 | 新文件 |
|--------|--------|
| `GanttChart.vue` | `GanttChartRefactored.vue` |
| `useGanttDrag.js` | `useGanttDrag.ts` |
| `ganttStore.js` | `ganttStore.js` (保持兼容) |

### 新增文件

```
src/
├── types/gantt.ts              # TypeScript 类型定义
├── config/breakpoints.ts       # 断点配置
├── styles/themes/
│   ├── variables.css           # CSS 变量
│   └── dark.css                # 暗色主题
├── composables/
│   ├── useTheme.ts             # 主题管理
│   ├── useBreakpoint.ts        # 响应式检测
│   ├── useVirtualScroll.ts     # 虚拟滚动
│   ├── useGanttKeyboard.ts     # 键盘快捷键
│   ├── useGanttTooltip.ts      # 工具提示
│   ├── useGanttSelection.ts    # 选择管理
│   ├── useTouchGestures.ts     # 触摸手势
│   └── useA11y.ts              # 无障碍功能
└── components/gantt/
    ├── core/                    # 核心组件
    ├── table/                   # 表格组件
    ├── timeline/                # 时间轴组件
    ├── virtual/                 # 虚拟滚动组件
    └── mobile/                  # 移动端组件
```

---

## 迁移检查清单

### 基础迁移
- [ ] 更新组件导入路径
- [ ] 更新 props 类型定义
- [ ] 更新事件处理方式
- [ ] 更新 store 访问方式

### 功能迁移
- [ ] 添加主题支持
- [ ] 添加响应式支持
- [ ] 添加键盘快捷键
- [ ] 添加无障碍支持

### 测试
- [ ] 单元测试通过
- [ ] E2E 测试通过
- [ ] 性能测试达标
- [ ] 兼容性测试通过

---

## 需要帮助？

如果在迁移过程中遇到问题：

1. 查看 [API 文档](./API.md)
2. 查看 [架构文档](./ARCHITECTURE.md)
3. 在 Issues 中搜索类似问题
4. 提交新的 Issue

---

**最后更新**: 2026-01-31
**版本**: 2.0.0
