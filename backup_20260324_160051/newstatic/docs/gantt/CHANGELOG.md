# 变更日志

本文档记录甘特图组件的所有重要变更。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)

---

## [2.0.0] - 2026-01-31

### 🎉 重大重构 - 现代化甘特图组件

本次更新是完全重构，采用现代化的软件设计理念，全面升级用户体验和开发体验。

### ✨ 新增功能

#### 响应式设计
- 📱 移动端完全支持
- 📱 触摸手势支持（滑动、缩放、长按）
- 💻 平板和桌面端自适应
- 📐 响应式断点：320px / 375px / 768px / 1024px / 1280px

#### 性能优化
- ⚡ 虚拟滚动支持（自动启用阈值：50+ 任务）
- 🚀 依赖关系计算优化（O(n²) → O(n+m)）
- 🎯 拖拽性能优化（RAF 节流，目标 60fps）
- 🔍 搜索防抖优化

#### 主题系统
- 🎨 亮色/暗色/自动主题
- 🎨 主题持久化（localStorage）
- 🎨 CSS 变量系统（可定制）
- 🎨 动态主题切换

#### 键盘快捷键
- ⌨️ 完整键盘导航支持
- ⌨️ 50+ 快捷键
- ⌨️ 快捷键帮助面板（`?` 或 `Ctrl+/`）
- ⌨️ 自定义快捷键支持

#### 无障碍访问
- ♿ WCAG 2.1 AA 级别合规
- ♿ ARIA 标签和角色
- ♿ 屏幕阅读器支持
- ♿ 焦点管理
- ♿ 跳过导航链接
- ♿ 高对比度模式支持
- ♿ 减少动画模式支持

#### TypeScript 支持
- 📝 完整类型定义（30+ 类型）
- 📝 泛型支持
- 📝 类型推导
- 📝 JSDoc 注释

#### 新组件
- `GanttBody.vue` - 内容区容器
- `MobileGanttView.vue` - 移动端视图
- `TaskListView.vue` - 移动端任务列表
- `TimelineSwipeView.vue` - 移动端时间轴
- `TaskBar.vue` - 任务条组件
- `VirtualList.vue` - 虚拟列表
- `VirtualTimeline.vue` - 虚拟时间轴
- `ThemeToggle.vue` - 主题切换
- `ShortcutHelpPanel.vue` - 快捷键帮助
- `A11yAnnouncer.vue` - 屏幕阅读通知
- `FocusTrap.vue` - 焦点陷阱
- `SkipLink.vue` - 跳过链接

#### 新 Composables
- `useTheme.ts` - 主题管理
- `useBreakpoint.ts` - 响应式检测
- `useVirtualScroll.ts` - 虚拟滚动
- `useTouchGestures.ts` - 触摸手势
- `useGanttKeyboard.ts` - 键盘快捷键
- `useGanttTooltip.ts` - 工具提示
- `useGanttSelection.ts` - 选择管理
- `useA11y.ts` - 无障碍功能

#### 新工具函数
- `dependencyGraph.ts` - 依赖图优化
- `breakpoints.ts` - 断点配置

### 🔄 破坏性变更

#### 组件导入
```typescript
// 旧版本
import GanttChart from '@/components/progress/GanttChart.vue'

// 新版本
import GanttChartRefactored from '@/components/progress/GanttChartRefactored.vue'
```

#### 类型定义
```typescript
// Props 现在需要类型
interface Props {
  projectId: number | string
  projectName?: string
  scheduleData?: ScheduleData
}
```

#### 事件处理
```typescript
// 旧版本
this.$emit('task-updated', task)

// 新版本
const emit = defineEmits<{
  taskUpdated: [task: GanttTask]
}>()
emit('taskUpdated', task)
```

### 🐛 Bug 修复

- 修复大量任务时拖拽卡顿问题
- 修复依赖关系计算性能问题
- 修复移动端显示问题
- 修复暗色主题显示问题

### 📚 文档更新

- 添加 [README.md](./README.md)
- 添加 [API.md](./API.md)
- 添加 [MIGRATION.md](./MIGRATION.md)
- 添加 [ARCHITECTURE.md](./ARCHITECTURE.md)
- 添加 [CONTRIBUTING.md](./CONTRIBUTING.md)
- 添加 [CHANGELOG.md](./CHANGELOG.md)

### 📦 依赖更新

```json
{
  "vue": "^3.4.0",
  "typescript": "^5.0.0",
  "element-plus": "^2.5.0",
  "lodash-es": "^4.17.0"
}
```

### 🏗️ 架构变更

- 组件模块化（单文件 1300+ 行 → 50+ 文件）
- 业务逻辑提取到 Composables
- 样式系统模块化
- 类型系统集中管理

---

## [1.0.0] - 2025-01-15

### 初始版本

#### 功能
- ✅ 基础甘特图功能
- ✅ 任务管理（CRUD）
- ✅ 拖拽调整时间
- ✅ 依赖关系管理
- ✅ 资源分配
- ✅ 多视图切换（日/周/月/季度）
- ✅ 统计信息展示

#### 组件
- `GanttChart.vue` - 主组件
- `GanttToolbar.vue` - 工具栏
- `TaskTable.vue` - 任务表格
- `TaskTimeline.vue` - 时间轴
- `GanttHeader.vue` - 表头
- 等其他支持组件

#### 状态管理
- `ganttStore.js` - Pinia Store

#### 文档
- 基础使用说明
- API 文档
- 示例代码

---

## [Unreleased]

### 计划中

- [ ] 插件系统
- [ ] 数据导入/导出增强
- [ ] 实时协作
- [ ] 离线支持
- [ ] 更多主题预设

---

## 版本说明

- **[Unreleased]**: 开发中，尚未发布
- **[2.0.0]**: 重大重构版本（2026-01-31）
- **[1.0.0]**: 初始版本（2025-01-15）

---

**更新日期**: 2026-01-31
