# 系统架构文档

本文档详细说明甘特图组件的架构设计、技术选型和最佳实践。

## 目录

- [架构概览](#架构概览)
- [技术栈](#技术栈)
- [目录结构](#目录结构)
- [数据流](#数据流)
- [组件架构](#组件架构)
- [状态管理](#状态管理)
- [性能优化](#性能优化)
- [设计模式](#设计模式)

---

## 架构概览

### 整体架构图

```
┌─────────────────────────────────────────────────────────────┐
│                         应用层 (App)                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ 路由 (Router) │  │ 状态 (Store)  │  │ 工具 (Utils)  │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                       业务逻辑层 (Logic)                     │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ Composables  │  │ 类型 (Types)  │  │ 工厂 (Utils)  │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                         视图层 (View)                        │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ 核心组件      │  │ 表格组件      │  │ 时间轴组件    │      │
│  │ (Core)        │  │ (Table)       │  │ (Timeline)    │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ 虚拟滚动      │  │ 移动端组件    │  │ 通用组件      │      │
│  │ (Virtual)     │  │ (Mobile)      │  │ (Common)      │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                         样式层 (Style)                        │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ CSS 变量      │  │ 主题样式      │  │ 组件样式      │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
```

---

## 技术栈

### 前端框架

| 技术 | 版本 | 用途 |
|------|------|------|
| Vue | 3.x | 渐染框架 |
| TypeScript | 5.x | 类型系统 |
| Vite | 5.x | 构建工具 |
| Element Plus | 2.x | UI 组件库 |

### 状态管理

| 技术 | 模式 | 说明 |
|------|------|------|
| Pinia | - | 状态管理（推荐） |
| ganttStore | Custom | 自定义集中式 Store |

### 工具库

| 技术 | 版本 | 用途 |
|------|------|------|
| lodash-es | ^4.17 | 工具函数 |
| date-fns | ^3.0 | 日期处理 |
| html2canvas | ^1.4 | 导出图片 |
| jspdf | ^2.5 | 导出 PDF |

---

## 目录结构

### 完整目录树

```
src/
├── components/
│   ├── gantt/
│   │   ├── core/                    # 核心组件
│   │   │   ├── GanttChart.vue       # 主容器（待重构 <300行）
│   │   │   ├── GanttToolbar.vue     # 工具栏
│   │   │   ├── GanttBody.vue        # 内容区容器
│   │   │   └── GanttHeader.vue      # 时间轴表头
│   │   ├── table/                   # 表格组件
│   │   │   ├── TaskTable.vue        # 任务表格
│   │   │   ├── TaskRow.vue          # 单个任务行
│   │   │   └── TaskCell.vue         # 可编辑单元格
│   │   ├── timeline/                # 时间轴组件
│   │   │   ├── TaskTimeline.vue      # 时间轴
│   │   │   ├── TaskBar.vue          # 任务条
│   │   │   ├── DependencyLine.vue   # 依赖线
│   │   │   └── TimelineGrid.vue     # 时间网格
│   │   ├── virtual/                 # 虚拟滚动
│   │   │   ├── VirtualList.vue      # 虚拟列表
│   │   │   └── VirtualTimeline.vue  # 虚拟时间轴
│   │   ├── mobile/                  # 移动端
│   │   │   ├── MobileGanttView.vue  # 移动端视图
│   │   │   ├── TaskListView.vue     # 任务列表
│   │   │   └── TimelineSwipeView.vue # 时间轴
│   │   └── dialogs/                 # 对话框
│   │       ├── TaskEditDialog.vue
│   │       ├── ResourceDialog.vue
│   │       └── SettingsDialog.vue
│   ├── common/                      # 通用组件
│   │   ├── ThemeToggle.vue         # 主题切换
│   │   ├── ShortcutHelpPanel.vue   # 快捷键帮助
│   │   ├── A11yAnnouncer.vue       # 屏幕阅读通知
│   │   ├── FocusTrap.vue           # 焦点陷阱
│   │   └── SkipLink.vue            # 跳过链接
│   └── progress/                    # 进度组件
│       ├── GanttChart.vue          # 原始组件（保留兼容）
│       ├── GanttChartRefactored.vue # 重构组件
│       ├── GanttToolbar.vue
│       ├── GanttStats.vue
│       └── ... (其他现有组件)
├── composables/                      # 组合式函数
│   ├── useGanttDrag.ts            # 拖拽逻辑
│   ├── useVirtualScroll.ts        # 虚拟滚动
│   ├── useBreakpoint.ts           # 响应式检测
│   ├── useTouchGestures.ts        # 触摸手势
│   ├── useGanttKeyboard.ts        # 键盘快捷键
│   ├── useGanttTooltip.ts         # 工具提示
│   ├── useGanttSelection.ts       # 选择管理
│   ├── useTheme.ts                # 主题管理
│   └── useA11y.ts                 # 无障碍功能
├── stores/                          # 状态管理
│   └── ganttStore.js              # 甘特图 Store
├── types/                           # 类型定义
│   └── gantt.ts                   # 甘特图类型
├── utils/                           # 工具函数
│   ├── dateFormat.ts             # 日期格式化
│   ├── ganttHelpers.ts           # 甘特图辅助函数
│   ├── dependencyGraph.ts        # 依赖图优化
│   └── eventBus.ts               # 事件总线
├── config/                          # 配置文件
│   └── breakpoints.ts             # 断点配置
└── styles/                          # 样式文件
      └── themes/
          ├── variables.css        # CSS 变量
          └── dark.css             # 暗色主题
```

---

## 数据流

### 组件通信流程

```
┌──────────────┐
│  Parent View  │
└──────┬───────┘
       │ props
       ↓
┌─────────────────────┐
│  GanttChart Component │
└──────┬──────────────┘
       │
       ├─────────────────────┬──────────────────────┐
       ↓                     ↓                      ↓
┌──────────────┐   ┌──────────────┐   ┌──────────────┐
│  GanttToolbar │   │  GanttStats  │   │ GanttLegend  │
└──────────────┘   └──────────────┘   └──────────────┘
       │
       ↓
┌─────────────────────┐
│     GanttBody       │
└──────┬──────────────┘
       │
       ├──────────┬──────────────┐
       ↓          ↓              ↓
┌──────────┐  ┌──────────┐  ┌──────────┐
│TaskTable │  │Timeline  │  │ Modals   │
└──────────┘  └──────────┘  └──────────┘
       │          │              │
       └──────────┴──────────────┘
                  ↓
         ┌──────────────┐
         │ ganttStore   │
         └──────────────┘
```

### 状态流向

```
┌─────────────┐
│  User Action │
└──────┬──────┘
       │ emit
       ↓
┌──────────────┐
│  Component   │
└──────┬───────┘
       │ action
       ↓
┌──────────────┐      ┌─────────────┐
│ ganttStore   │ ←── │  API Layer  │
└──────────────┘      └─────────────┘
       │
       ↓ reactive
┌──────────────┐
│   Component  │
└──────────────┘
       │ render
       ↓
┌──────────────┐
│     DOM      │
└──────────────┘
```

---

## 组件架构

### 组件层次结构

```
GanttChartRefactored (根容器)
│
├─ GanttToolbar (工具栏)
│  ├─ 搜索框
│  ├─ 视图切换
│  ├─ 缩放控制
│  ├─ 显示选项
│  └─ actions 插槽
│
├─ GanttStats (统计信息)
│
├─ GanttBody (内容区容器)
│  │
│  ├─ GanttHeader (表头)
│  │  └─ 时间刻度
│  │
│  └─ 表格+时间轴 (可滚动)
│     ├─ TaskTableRefactored (任务表格)
│     │  └─ TaskRow (任务行)
│     │
│     └─ TaskTimelineRefactored (时间轴)
│        ├─ TaskBar (任务条)
│        ├─ DependencyLine (依赖线)
│        └─ TimelineGrid (时间网格)
│
├─ 调整大小手柄
│
├─ GanttLegend (图例)
│
├─ 对话框组
│  ├─ TaskEditDialog
│  ├─ ResourceAllocationDialog
│  ├─ ResourceManagementDialog
│  └─ GanttContextMenu
│
├─ TaskDetailDrawer (任务详情)
│
├─ GanttStatusBar (状态栏)
│
└─ 覆盖层
   ├─ ShortcutHelpPanel
   ├─ ThemeToggle
   ├─ A11yAnnouncer
   └─ SkipLink
```

### 组件职责划分

| 组件 | 职责 | 行数限制 |
|------|------|----------|
| `GanttChart` | 容器、事件协调 | < 300 |
| `GanttToolbar` | 工具栏 UI | < 200 |
| `GanttBody` | 表格+时间轴容器 | < 150 |
| `TaskTable` | 表格渲染 | < 400 |
| `TaskTimeline` | 时间轴渲染 | < 300 |
| `TaskBar` | 单个任务条 | < 100 |
| `VirtualList` | 虚拟滚动 | < 150 |
| `MobileGanttView` | 移动端视图 | < 200 |

---

## 状态管理

### ganttStore 结构

```javascript
{
  // 项目信息
  state: {
    projectId: null,
    projectName: '',
    scheduleData: {},
    tasks: [],
    filteredTasks: [],
    // ... 其他状态
  },

  // 计算属性 (getters)
  getters: {
    selectedTask: computed(() => {}),
    taskStats: computed(() => {}),
    timelineDays: computed(() => {}),
    // ...
  },

  // 操作方法 (actions)
  actions: {
    setProject() {},
    formatTasks() {},
    selectTask() {},
    // ...
  }
}
```

### 状态更新流程

```
用户操作
   ↓
组件事件 (emit)
   ↓
Store Action
   ↓
State 更新
   ↓
响应式重渲染
```

---

## 性能优化

### 1. 虚拟滚动

```typescript
// 使用虚拟滚动前
// 1000 个任务 = 1000 个 DOM 节点

// 使用虚拟滚动后
// 1000 个任务 = 约 20 个可见 DOM 节点
```

### 2. 依赖关系计算优化

```typescript
// 旧版：O(n²) 算法
function findCriticalPathOld(tasks) {
  // 嵌套循环
  for (let i = 0; i < tasks.length; i++) {
    for (let j = 0; j < tasks.length; j++) {
      // ...
    }
  }
}

// 新版：O(n + m) 算法
function findCriticalPathNew(tasks) {
  const graph = new DependencyGraph()
  graph.buildGraph(tasks)  // O(n + m)
  return graph.calculateCriticalPath()  // O(n + m)
}
```

### 3. 拖拽节流

```typescript
// 使用 RAF 节流，目标 60fps
const throttledUpdate = rafThrottle((dayOffset) => {
  updatePreviewTask(dayOffset)
}, 16) // ~16ms per frame
```

---

## 设计模式

### 1. Composable 模式

将可复用逻辑提取到组合式函数中：

```typescript
// useGanttDrag.ts
export function useGanttDrag(options) {
  // 封装拖拽逻辑
  return {
    isDragging,
    startDrag,
    cancelDrag
  }
}
```

### 2. Provider/Inject 模式

通过 Store 提供全局状态：

```typescript
// ganttStore.js
export const ganttStore = reactive({
  state: {},
  getters: {},
  actions: {}
})

// 组件中使用
import { ganttStore } from '@/stores/ganttStore'
const { state } = ganttStore
```

### 3. 观察者模式

事件总线实现组件间通信：

```typescript
// eventBus.js
export const eventBus = mitt()

// 发送事件
eventBus.emit('task:updated', task)

// 监听事件
eventBus.on('task:updated', (task) => {
  // ...
})
```

### 4. 工厂模式

创建依赖关系图实例：

```typescript
export function createDependencyGraph(tasks) {
  const graph = new DependencyGraph()
  graph.buildGraph(tasks)
  return graph
}
```

### 5. 策略模式

不同断点使用不同渲染策略：

```typescript
const renderingStrategy = computed(() => {
  if (isMobile.value) {
    return 'card'      // 卡片视图
  } else if (tasks.length > 50) {
    return 'virtual'   // 虚拟滚动
  } else {
    return 'standard'   // 标准渲染
  }
})
```

---

## 扩展性设计

### 插槽系统

组件提供多个插槽供扩展：

```vue
<GanttChartRefactored>
  <template #actions>
    <CustomAction />
  </template>
  <template #tools>
    <CustomTool />
  </template>
</GanttChartRefactored>
```

### 主题系统

通过 CSS 变量扩展主题：

```css
/* 自定义主题 */
.gantt-chart.my-theme {
  --color-primary: #custom-color;
  --gantt-row-height: 80px;
}
```

### 插件系统（规划中）

```typescript
// 未来支持插件系统
GanttChart.use(PluginA)
GanttChart.use(PluginB)
```

---

## 安全性

### XSS 防护

```vue
<!-- 使用 v-text 防止 XSS -->
<div v-text="userInput"></div>

<!-- 避免使用 v-html -->
<div v-html="sanitizedHtml"></div>
```

### 事件验证

```typescript
function handleKeydown(event) {
  // 验证事件目标
  if (event.target.tagName === 'INPUT') {
    return
  }
  // 处理快捷键
}
```

---

## 可维护性

### 代码规范

- 单文件行数限制
- 函数复杂度控制
- TypeScript 类型覆盖率 > 90%
- ESLint 零警告

### 文档化

- API 文档完整
- 代码注释充分
- JSDoc 类型注释

### 测试覆盖

- 单元测试覆盖率 > 80%
- E2E 测试覆盖关键路径

---

## 未来规划

### Phase 6: 高级特性（规划中）

- 插件系统
- 数据导入/导出增强
- 实时协作
- 离线支持

### Phase 7: 性能优化（规划中）

- Web Worker 计算
- IndexedDB 缓存
- Service Worker

---

**更新日期**: 2026-01-31
**版本**: 2.0.0
