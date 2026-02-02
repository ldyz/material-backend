# 甘特图组件现代化重构

> 基于 Vue 3 + TypeScript 的现代化甘特图组件，支持响应式设计、虚拟滚动、主题切换和无障碍访问。

[![Vue 3](https://img.shields.io/badge/Vue-3.x-brightgreen)](https://vuejs.org/)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.x-blue)](https://www.typescriptlang.org/)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## ✨ 特性

### 🎯 核心功能
- **任务管理** - 创建、编辑、删除、复制任务
- **拖拽操作** - 拖拽调整任务时间和层级
- **依赖关系** - 可视化创建和管理任务依赖
- **关键路径** - 自动计算并高亮关键路径
- **资源管理** - 分配和管理项目资源
- **多视图** - 日/周/月/季度视图切换
- **分组显示** - 按状态、优先级分组

### 🎨 用户体验
- **主题切换** - 亮色/暗色/跟随系统主题
- **响应式设计** - 完美适配移动端和桌面端
- **键盘快捷键** - 完整的键盘操作支持
- **触摸手势** - 移动端滑动、缩放手势
- **无障碍访问** - WCAG 2.1 AA 级别合规

### ⚡ 性能优化
- **虚拟滚动** - 大数据量下流畅渲染
- **节流防抖** - 优化拖拽和搜索性能
- **按需加载** - 组件懒加载
- **计算优化** - O(n+m) 依赖关系计算

## 📦 安装

```bash
# 克隆项目
git clone https://github.com/your-org/gantt-chart.git

# 安装依赖
cd gantt-chart
npm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build
```

## 🚀 快速开始

### 基础用法

```vue
<template>
  <GanttChartRefactored
    :project-id="projectId"
    :project-name="projectName"
    :schedule-data="scheduleData"
    @task-updated="handleTaskUpdated"
  />
</template>

<script setup>
import { ref } from 'vue'
import { GanttChartRefactored } from '@/components/progress'

const projectId = ref(1)
const projectName = ref('我的项目')
const scheduleData = ref({ activities: {}, nodes: {} })

function handleTaskUpdated(task) {
  console.log('任务已更新:', task)
}
</script>
```

### 配置选项

```vue
<GanttChartRefactored
  :project-id="projectId"
  :project-name="projectName"
  :schedule-data="scheduleData"

  <!-- 可选配置 -->
  :default-view-mode="'day'"
  :default-theme="'light'"
  :enable-virtual-scroll="true"
  :virtual-scroll-threshold="50"

  @task-updated="handleUpdate"
  @task-selected="handleSelect"
/>
```

### Props

| 属性 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `project-id` | `Number \| String` | **必填** | 项目ID |
| `project-name` | `String` | `'未命名项目'` | 项目名称 |
| `schedule-data` | `Object` | `{}` | 进度数据对象 |
| `default-view-mode` | `String` | `'day'` | 默认视图模式 |
| `default-theme` | `String` | `'light'` | 默认主题 |

### Events

| 事件 | 参数 | 说明 |
|------|------|------|
| `@task-updated` | `(task: Object)` | 任务更新时触发 |
| `@task-selected` | `(task: Object)` | 任务选中时触发 |

## ⌨️ 键盘快捷键

### 导航
- `↑` - 选择上一个任务
- `↓` - 选择下一个任务
- `←` - 选择前驱任务
- `→` - 选择后继任务
- `Home` - 跳到第一个任务
- `End` - 跳到最后一个任务

### 编辑
- `Enter` - 编辑选中的任务
- `Delete` - 删除选中的任务
- `Escape` - 取消操作/关闭对话框
- `Ctrl+N` - 新建任务
- `Ctrl+D` - 复制任务
- `Ctrl+C/V` - 复制/粘贴
- `Ctrl+Z` - 撤销
- `Ctrl+S` - 保存

### 视图
- `Ctrl++` - 放大时间轴
- `Ctrl+-` - 缩小时间轴
- `Ctrl+0` - 重置缩放
- `Alt+D` - 切换依赖关系显示
- `Alt+P` - 切换关键路径显示

### 帮助
- `?` 或 `Ctrl+/` - 显示快捷键帮助

## 🎨 主题定制

### CSS 变量

在 `src/styles/themes/variables.css` 中定义的CSS变量可以自定义：

```css
:root {
  /* 主色调 */
  --color-primary: #409eff;

  /* 甘特图尺寸 */
  --gantt-row-height: 60px;
  --gantt-task-height: 32px;
  --gantt-day-width: 40px;

  /* 任务条颜色 */
  --task-bar-completed: #67c23a;
  --task-bar-in-progress: #409eff;
  --task-bar-delayed: #f56c6c;
}
```

### 切换主题

```typescript
import { useTheme } from '@/composables'

const { mode, setTheme } = useTheme()

// 切换到暗色主题
setTheme('dark')

// 切换到亮色主题
setTheme('light')

// 跟随系统
setTheme('auto')
```

## 📱 响应式设计

组件会根据屏幕尺寸自动调整：

| 断点 | 尺寸 | 布局 |
|------|------|------|
| `xs` | < 375px | 移动端卡片视图 |
| `sm` | 375px - 768px | 移动端卡片视图 |
| `md` | 768px - 1024px | 平板紧凑表格 |
| `lg` | 1024px - 1280px | 桌面标准视图 |
| `xl` | ≥ 1280px | 桌面宽屏视图 |

## ♿ 无障碍访问

组件完全符合 WCAG 2.1 AA 级别标准：

- ✅ 完整键盘导航
- ✅ 屏幕阅读器支持
- ✅ ARIA 标签和角色
- ✅ 焦点管理
- ✅ 跳过导航链接
- ✅ 高对比度模式支持
- ✅ 减少动画模式支持

## 🔧 开发

### 项目结构

```
src/
├── components/
│   ├── gantt/
│   │   ├── core/              # 核心组件
│   │   ├── table/             # 表格组件
│   │   ├── timeline/          # 时间轴组件
│   │   ├── virtual/           # 虚拟滚动组件
│   │   └── mobile/            # 移动端组件
│   ├── common/                # 通用组件
│   └── progress/              # 进度组件
├── composables/                # 组合式函数
├── stores/                     # 状态管理
├── types/                      # TypeScript 类型
├── utils/                      # 工具函数
└── styles/                     # 样式文件
```

### 构建

```bash
# 开发环境
npm run dev

# 类型检查
npm run type-check

# 代码检查
npm run lint

# 单元测试
npm run test:unit

# E2E 测试
npm run test:e2e

# 构建
npm run build
```

## 📖 文档

- [API 文档](./API.md) - 组件 API 详细说明
- [迁移指南](./MIGRATION.md) - 从旧版本迁移
- [架构文档](./ARCHITECTURE.md) - 系统架构设计
- [贡献指南](./CONTRIBUTING.md) - 如何贡献代码

## 🐛 问题反馈

请在 [Issues](https://github.com/your-org/gantt-chart/issues) 中提交问题。

## 📄 许可证

[MIT](LICENSE)

## 🙏 致谢

感谢所有贡献者的支持！

---

**更新日期**: 2026-01-31
**版本**: 2.0.0
