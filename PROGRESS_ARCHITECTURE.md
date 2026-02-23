# 进度管理系统架构说明文档

## 目录
1. [系统概述](#系统概述)
2. [技术栈](#技术栈)
3. [架构设计](#架构设计)
4. [前端架构](#前端架构)
5. [后端架构](#后端架构)
6. [前后端交互](#前后端交互)
7. [数据流向](#数据流向)
8. [代码清理建议](#代码清理建议)

---

## 系统概述

进度管理系统是一个企业级项目管理工具，提供以下核心功能：

- **甘特图可视化**：直观展示项目进度和任务关系
- **网络图视图**：基于节点和边的任务依赖关系可视化
- **任务管理**：创建、编辑、删除、拖拽任务
- **依赖关系管理**：支持FS/FF/SS/SF四种依赖类型，使用A*算法自动规划路径
- **资源管理**：资源分配、资源平衡、资源利用率分析
- **AI辅助**：AI自动生成项目进度计划
- **进度跟踪**：进度统计、挣值分析、燃尽图
- **多视图支持**：甘特图、网络图、日历、看板等

---

## 技术栈

### 前端技术栈
- **框架**: Vue 3 (Composition API)
- **UI组件库**: Element Plus
- **状态管理**: Pinia Store (自定义响应式Store)
- **图表**: 自定义SVG组件 + ECharts
- **构建工具**: Vite
- **语言**: JavaScript + TypeScript

### 后端技术栈
- **框架**: Go + Gin
- **ORM**: GORM
- **数据库**: PostgreSQL/MySQL
- **认证**: JWT
- **AI集成**: Claude API

---

## 架构设计

### 整体架构图

```
┌─────────────────────────────────────────────────────────────┐
│                        前端 (Vue 3)                         │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │  视图层      │  │  组件层      │  │  Store层     │     │
│  │  Progress.vue│  │  GanttChart  │  │  ganttStore  │     │
│  │  NetworkDiagram│ │  TaskTable   │  │  undoRedoStore│    │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
│         │                  │                  │              │
│         └──────────────────┼──────────────────┘              │
│                            │                                 │
│  ┌─────────────────────────┼─────────────────────────┐     │
│  │              API层 (api/index.js)                 │     │
│  │         progressApi (定义所有API调用)              │     │
│  └─────────────────────────┼─────────────────────────┘     │
└────────────────────────────┼───────────────────────────────┘
                             │ HTTP/REST
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                    后端 (Go + Gin)                          │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │  路由层      │  │  处理器层    │  │  数据层      │     │
│  │  routes.go   │  │  handler.go  │  │  GORM        │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
│                            │                                 │
│  ┌─────────────────────────┼─────────────────────────┐     │
│  │              服务层 (Business Logic)               │     │
│  │   - 任务管理   - 依赖管理   - 资源管理             │     │
│  │   - AI生成     - 数据验证   - 权限控制             │     │
│  └─────────────────────────┼─────────────────────────┘     │
└────────────────────────────┼───────────────────────────────┘
                             │
                             ▼
                    ┌───────────────┐
                    │   数据库      │
                    │ PostgreSQL/   │
                    │   MySQL       │
                    └───────────────┘
```

---

## 前端架构

### 目录结构

```
newstatic/src/
├── components/
│   ├── progress/                    # 进度管理组件
│   │   ├── Progress.vue             # 主容器
│   │   ├── GanttChart.vue           # 甘特图组件（正在使用）
│   │   ├── NetworkDiagram.vue       # 网络图组件
│   │   ├── TaskTable.vue            # 任务表格
│   │   ├── TaskEditDialog.vue       # 任务编辑对话框
│   │   ├── ResourceManagementDialog.vue
│   │   ├── CreateScheduleDialog.vue
│   │   └── ...
│   └── gantt/                       # 甘特图组件库
│       ├── core/                    # 核心组件
│       ├── timeline/                # 时间轴组件
│       ├── table/                   # 表格组件
│       ├── dialogs/                 # 对话框组件
│       ├── views/                   # 视图组件
│       └── ...
├── stores/
│   └── ganttStore.js                # 甘特图状态管理
│   └── undoRedoStore.js             # 撤销重做状态管理
├── composables/
│   ├── useGanttDrag.js              # 甘特图拖拽逻辑
│   ├── usePermission.js             # 权限管理
│   └── useUndoRedo.js               # 撤销重做逻辑
├── utils/
│   └── ganttHelpers.js              # 甘特图工具函数
└── api/
    └── index.js                     # API定义
```

### 核心组件职责

#### 1. Progress.vue (主容器)
**职责**:
- 项目和进度计划管理
- 视图切换（列表/甘特图/网络图）
- 数据加载和初始化
- 对话框管理

**关键方法**:
- `fetchProjects()`: 获取项目列表
- `loadAllProjectSchedules()`: 加载所有进度计划
- `handleCreateSchedule()`: 创建进度计划
- `handleDeleteSchedule()`: 删除进度计划
- `switchView()`: 切换视图

**引用的子组件**:
- `GanttChart`: 甘特图视图
- `NetworkDiagram`: 网络图视图
- `ProjectScheduleList`: 进度计划列表
- `CreateScheduleDialog`: 创建对话框

#### 2. GanttChart.vue (甘特图组件)
**职责**:
- 渲染甘特图时间轴和任务条
- 处理任务拖拽（移动、调整大小）
- 显示依赖关系线
- 缩放和滚动控制

**关键功能**:
- 时间轴渲染（日/周/月/季度）
- 任务条渲染和交互
- 依赖关系可视化（使用A*路径规划）
- 任务拖拽和调整

**使用的composables**:
- `useGanttDrag`: 拖拽逻辑
- `useGanttStore`: 状态管理

#### 3. NetworkDiagram.vue (网络图组件)
**职责**:
- 渲染节点和边
- 节点拖拽调整任务时间
- 显示依赖关系箭头
- 自动布局算法

**关键功能**:
- SVG节点渲染
- 节点拖拽和位置保存
- 箭头路径计算（A*算法）
- 自动布局和层级管理

#### 4. TaskTable.vue / TaskList.vue
**职责**:
- 表格形式展示任务列表
- 任务编辑和删除
- 进度更新

### 状态管理 (ganttStore.js)

#### State 结构
```javascript
{
  // 项目信息
  projectId: null,
  projectName: '',

  // 任务数据
  scheduleData: { activities: {}, nodes: {} },
  tasks: [],
  filteredTasks: [],

  // 视图配置
  viewMode: 'day',              // day, week, month, quarter
  dayWidth: 40,
  rowHeight: 50,

  // 显示选项
  showDependencies: true,
  showCriticalPath: true,

  // 交互状态
  selectedTaskId: null,
  isDragging: false,
  dragMode: 'none',             // none, move, resize_left, resize_right

  // 自动保存
  hasUnsavedChanges: false,
  pendingTaskUpdates: new Map(),
}
```

#### 主要 Actions
```javascript
// 数据管理
setProject(projectId)
loadData()
formatTasks(rawData)

// 任务操作
selectTask(taskId)
startDrag(task, mode)
updateDragPreview(dayOffset)
endDrag(newTask, originalTask)

// 依赖关系
startDependencyCreation(task)
completeDependencyCreation(targetTask)
cancelDependencyCreation()

// 视图控制
setViewMode(mode)
zoomIn()
zoomOut()
autoFit()

// 保存
saveAll()
```

#### Getters
```javascript
// 时间轴相关
timelineDays      // 按天计算的时间轴
timelineWeeks     // 按周计算的时间轴
timelineMonths    // 按月计算的时间轴
timelineWidth     // 时间轴总宽度

// 任务相关
taskStats         // 任务统计信息
groupedTasks      // 分组后的任务
selectedTask      // 当前选中的任务

// 依赖关系
visibleDependencies  // 可见的依赖关系
```

### Composables

#### useGanttDrag.js
**职责**: 封装甘特图拖拽逻辑

**功能**:
- 检测拖拽模式（移动/左边缘/右边缘）
- 处理鼠标移动和位置计算
- 调整任务开始/结束日期
- 触发拖拽开始/结束回调

**导出**:
```javascript
{
  isDragging,      // 是否正在拖拽
  draggedTask,     // 被拖拽的任务
  tooltipPosition, // 提示框位置
  tooltipVisible,  // 提示框可见性
  startDrag,       // 开始拖拽
  cancelDrag,      // 取消拖拽
  detectDragMode,  // 检测拖拽模式
}
```

#### useDependencyCreation.js
**职责**: 处理依赖关系创建流程

**状态**: ⚠️ **未使用，可删除**

### 工具函数 (ganttHelpers.js)

#### 依赖关系计算
- `calculateDependencyPath()`: 计算依赖路径
- `findOrthogonalPath()`: A*路径规划算法
- `buildObstacleMap()`: 构建障碍物地图
- `tryDirectPath()`: 尝试直接路径
- `tryRouteAbove()`: 尝试上方路径
- `tryRouteBelow()`: 尝试下方路径
- `tryRouteRight()`: 尝试右侧路径
- `tryZigzagPath()`: 尝试锯齿形路径

#### 任务分组和排序
- `groupTasksByStatus()`: 按状态分组
- `groupTasksByPriority()`: 按优先级分组
- `sortTasksByPriority()`: 按优先级排序
- `getPriorityWeight()`: 获取优先级权重

#### 日期处理
- `getWeekStart()`: 获取周开始
- `getWeekEnd()`: 获取周结束
- `getWeekNumber()`: 获取周数
- `calculateWorkingDays()`: 计算工作日
- `isDateRangeOverlapping()`: 检查日期重叠

#### 任务验证
- `isMilestone()`: 判断是否为里程碑
- `validateTaskDates()`: 验证任务日期
- `detectDependencyCycle()`: 检测依赖循环

---

## 后端架构

### 目录结构
```
internal/api/progress/
├── handler.go          # 主要请求处理器
├── model.go            # 数据模型
├── routes.go           # 路由定义
├── ai_handler.go       # AI生成处理
├── calendar.go         # 日历功能
├── change_log_handler.go  # 变更日志
├── constraint.go       # 约束管理
├── report.go           # 报告生成
├── resource_leveling.go  # 资源平衡
└── websocket.go        # WebSocket支持
```

### 数据模型

#### Task (任务模型)
```go
type Task struct {
    ID              uint
    ProjectID       uint
    ParentID        *uint           // 父任务ID
    Name            string
    Duration        *float64        // 工期（天数）
    StartDate       *time.Time      // 开始日期
    EndDate         *time.Time      // 结束日期
    Progress        float64         // 进度百分比
    Priority        string          // 优先级（high/medium/low）
    Status          string          // 状态
    PositionX       *float64        // 网络图X坐标
    PositionY       *float64        // 网络图Y坐标
    TaskType        string          // 任务类型
    IsMilestone     bool            // 是否为里程碑
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

#### TaskDependency (任务依赖模型)
```go
type TaskDependency struct {
    ID        uint
    TaskID    uint            // 后置任务ID
    DependsOn uint            // 前置任务ID
    Type      string          // 依赖类型：FS/FF/SS/SF
    Lag       int             // 延迟天数
    CreatedAt time.Time
}
```

#### ProjectSchedule (进度计划模型)
```go
type ProjectSchedule struct {
    ID        uint
    ProjectID uint
    Data      ScheduleData    // JSON格式
    CreatedAt time.Time
    UpdatedAt time.Time
}

type ScheduleData struct {
    Activities map[string]Activity  // 任务活动
    Nodes      map[string]Node      // 节点数据
}
```

### API 路由

#### 项目进度管理
```
GET    /progress/project/:id                    # 获取项目进度
PUT    /progress/project/:id                    # 更新项目进度
DELETE /progress/project/:id/schedule           # 删除项目进度
GET    /progress/project/:id/exists             # 检查进度是否存在
GET    /progress/project-schedules              # 获取所有项目进度
POST   /progress/project/:id/generate-plan      # AI生成进度计划
POST   /progress/project/:id/aggregate-plan     # 聚合子项目计划
```

#### 任务管理
```
GET    /progress/project/:id/tasks              # 获取项目任务列表
POST   /progress/project/:id/tasks              # 创建任务
PUT    /progress/tasks/:id                      # 更新任务
DELETE /progress/tasks/:id                      # 删除任务
GET    /progress/tasks/:id/dependencies         # 获取任务依赖
PUT    /progress/tasks/:id/position             # 更新任务位置
POST   /progress/tasks/:id/calculate-parent-progress  # 计算父任务进度
POST   /progress/tasks/:id/update-parent-progress    # 更新父任务进度
```

#### 依赖关系管理
```
POST   /progress/tasks/:id/dependencies         # 添加依赖
DELETE /progress/dependencies/:id               # 删除依赖
POST   /progress/dependencies/visual/:from/:to  # 可视化创建依赖
```

#### 资源管理
```
GET    /progress/project/:id/resources          # 获取项目资源
POST   /progress/project/:id/resources          # 创建资源
PUT    /progress/project/:id/resources/:id      # 更新资源
DELETE /progress/project/:id/resources/:id      # 删除资源
GET    /progress/tasks/:id/resources            # 获取任务资源
POST   /progress/tasks/:id/resources            # 分配资源
DELETE /progress/tasks/:id/resources/:id        # 移除任务资源
```

### 核心处理器 (handler.go)

#### 任务管理
- `GetProgressList()`: 获取进度列表
- `GetTasks()`: 获取项目任务
- `CreateTask()`: 创建新任务
- `UpdateTask()`: 更新任务（支持依赖关系条件更新）
- `DeleteTask()`: 删除任务

#### 进度计划管理
- `GetProjectSchedule()`: 获取项目进度计划
- `UpdateProjectSchedule()`: 更新进度计划
- `DeleteProjectSchedule()`: 删除进度计划（修复了类型转换问题）
- `CheckScheduleExists()`: 检查进度计划是否存在

#### 依赖关系管理
- `GetDependencies()`: 获取任务依赖
- `AddDependency()`: 添加依赖关系
- `RemoveDependency()`: 删除依赖关系
- `CreateDependencyVisual()`: 可视化创建依赖

#### 资源管理
- `GetProjectResources()`: 获取项目资源
- `CreateResource()`: 创建资源
- `UpdateResource()`: 更新资源
- `DeleteResource()`: 删除资源
- `AllocateTaskResource()`: 分配资源给任务

#### 高级功能
- `GeneratePlanWithAI()`: AI生成进度计划（ai_handler.go）
- `AggregateChildPlans()`: 聚合子项目计划
- `ResourceLeveling()`: 资源平衡（resource_leveling.go）

---

## 前后端交互

### API 调用流程

#### 1. 加载项目进度计划
```
Progress.vue
  └─> fetchProjects()
       └─> progressApi.getAllProjectSchedules()
            └─> GET /progress/project-schedules
                 └─> Handler.GetProgressList()
                      └─> GORM查询
                           └─> 返回JSON
```

#### 2. 创建任务
```
TaskEditDialog.vue
  └─> emit('create-task', data)
       └─> GanttChart.handleCreateTask()
            └─> actions.createTask(data)
                 └─> progressApi.createTask(projectId, data)
                      └─> POST /progress/project/:id/tasks
                           └─> Handler.CreateTask()
                                ├─> 数据验证
                                ├─> 创建任务记录
                                ├─> 处理依赖关系
                                └─> 返回新任务
```

#### 3. 拖拽任务更新（甘特图）
```
GanttChart.vue
  └─> useGanttDrag.startDrag()
       └─> onDragEnd callback
            └─> actions.endDrag()
                 ├─> 更新本地状态
                 ├─> 计算依赖约束
                 └─> progressApi.update(taskId, data)
                      └─> PUT /progress/tasks/:id
                           └─> Handler.UpdateTask()
                                ├─> 只在需要时更新依赖关系
                                └─> 保存任务更新
```

#### 4. 拖拽节点更新（网络图）
```
NetworkDiagram.vue
  └─> onNodeDragEnd()
       ├─> progressApi.update(taskId, data)
       │    └─> PUT /progress/tasks/:id
       │         └─> Handler.UpdateTask()
       └─> emit('task-updated', { silent: true })
            └─> Progress.vue.handleTaskUpdated()
                 └─> actions.loadData()
```

#### 5. 创建依赖关系
```
GanttChart.vue / NetworkDiagram.vue
  └─> startDependencyCreation()
       └─> 点击目标任务
            └─> completeDependencyCreation()
                 └─> progressApi.createDependencyVisual(fromId, toId, data)
                      └─> POST /progress/dependencies/visual/:from/:to
                           └─> Handler.CreateDependencyVisual()
                                ├─> 验证循环依赖
                                ├─> 创建依赖记录
                                └─> 返回依赖关系
```

### API 响应格式

#### 成功响应
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "name": "任务名称",
    "start_date": "2024-01-01",
    "end_date": "2024-01-10"
  }
}
```

#### 错误响应
```json
{
  "code": 400,
  "message": "错误信息",
  "error": "详细错误描述"
}
```

---

## 数据流向

### 1. 初始化数据流
```
用户打开进度页面
  ↓
Progress.vue created()
  ↓
fetchProjects() → 获取项目列表
  ↓
loadAllProjectSchedules() → 获取进度计划状态
  ↓
用户选择项目/进度
  ↓
ganttStore.setProject(projectId)
  ↓
ganttStore.loadData()
  ↓
progressApi.getProjectSchedule(projectId)
  ↓
格式化任务数据
  ↓
更新视图
```

### 2. 任务更新数据流
```
用户拖拽任务条
  ↓
useGanttDrag 捕获拖拽事件
  ↓
计算新的日期
  ↓
更新预览（本地状态）
  ↓
用户释放鼠标
  ↓
onDragEnd callback
  ↓
progressApi.update(taskId, data)
  ↓
后端保存到数据库
  ↓
actions.loadData() → 重新加载数据
  ↓
视图更新
```

### 3. 依赖关系创建数据流
```
用户点击"添加前置任务"
  ↓
startDependencyCreation(taskId)
  ↓
设置 isCreating = true
  ↓
用户点击目标任务
  ↓
completeDependencyCreation(targetTask)
  ↓
progressApi.createDependencyVisual(fromId, toId)
  ↓
后端验证并创建依赖
  ↓
actions.loadData() → 重新加载数据
  ↓
重新计算依赖路径（A*算法）
  ↓
更新视图显示新依赖线
```

---

## 代码清理建议

### 可以安全删除的文件

#### 1. 未使用的组件
```
❌ newstatic/src/components/progress/GanttChartRefactored.vue
   原因: Progress.vue 使用的是 GanttChart.vue，此文件未被引用

❌ newstatic/src/components/progress/SimpleGantt.vue
   原因: 显示"开发中"的空页面，未被使用

❌ newstatic/src/composables/useDependencyCreation.js
   原因: 整个项目中未被导入或使用
```

#### 2. 功能冗余的组件
```
⚠️ newstatic/src/components/progress/ResourceView.vue
   原因: ResourceManagementDialog 已覆盖其功能
   建议: 测试后删除
```

### 需要重构的代码

#### 1. API 废弃标记
以下API方法已标记为 `@deprecated` 但仍在使用，建议添加警告：

```javascript
// api/index.js
getList()    // @deprecated → 应使用 getTasks()
getDetail()  // @deprecated → 应使用 getTasks()
create()     // @deprecated → 应使用 createTask()
update()     // @deprecated → 应使用 updateTask()
delete()     // @deprecated → 应使用 deleteTask()
```

**建议**: 添加控制台警告，提示开发者使用新API

#### 2. Progress.vue 中的注释代码
```vue
<!-- 第158-203行: 旧的甘特图视图 -->
<!-- 可以删除，已被 GanttChart 组件替代 -->
```

### 重复组件处理

#### BulkEditDialog.vue
存在于两个位置：
- `newstatic/src/components/gantt/dialogs/BulkEditDialog.vue`
- `newstatic/src/components/progress/BulkEditDialog.vue` (如果存在)

**建议**: 保留 `gantt/dialogs` 版本，删除其他重复版本

### 未使用的辅助函数检查

以下函数需要验证是否在模板中使用：

```javascript
// Progress.vue
getTaskBarStyle()      // 仅在旧甘特图视图中使用
getStatusTagType()     // 检查模板引用
getPriorityTagType()   // 检查模板引用
canEdit()              // 检查模板引用
canDelete()            // 检查模板引用
```

**建议**: 使用 IDE 的"查找引用"功能验证，删除未使用的函数

### 清理步骤

1. **第一阶段（无风险）**:
   - 删除 `GanttChartRefactored.vue`
   - 删除 `SimpleGantt.vue`
   - 删除 `useDependencyCreation.js`

2. **第二阶段（需测试）**:
   - 删除 `ResourceView.vue`
   - 清理 Progress.vue 中的注释代码

3. **第三阶段（重构）**:
   - 统一 API 方法，移除废弃的别名
   - 清理未使用的辅助函数
   - 合并重复的对话框组件

---

## 最近修复的问题

### 1. 依赖关系删除问题
**问题**: 更改任务日期时，依赖关系被意外删除
**原因**: `UpdateTask` 处理器无条件删除依赖关系
**修复**: 只在 `predecessor_ids` 或 `successor_ids` 存在时才删除依赖
**位置**: `internal/api/progress/handler.go:542-560`

### 2. 进度计划删除问题
**问题**: 删除操作返回成功但实际未删除
**原因**: `projectID` 参数类型不匹配（string vs uint）
**修复**: 使用 `strconv.ParseUint()` 转换类型
**位置**: `internal/api/progress/handler.go:1241-1301`

### 3. 网络图箭头显示问题
**问题**: 箭头未接触节点，存在间隙
**原因**: SVG marker refX 值过大
**修复**: 调整 refX 从 20 到 9，路径端点到 `x2 - radius - 1`
**位置**: `newstatic/src/components/progress/NetworkDiagram.vue:205-226`

### 4. 任务删除不刷新
**问题**: 删除任务后需要手动刷新页面
**修复**: 删除后调用 `actions.loadData()`
**位置**: `newstatic/src/components/progress/GanttChartRefactored.vue:1063`

### 5. 网络图拖拽不保存
**问题**: 拖拽节点后不保存到后端
**修复**: 添加 `emit('task-updated', { silent: true })`
**位置**: `newstatic/src/components/progress/NetworkDiagram.vue:2352`

### 6. 甘特图拖拽不保存
**问题**: 拖拽任务条后只提示"记得保存"，不自动保存
**修复**: 修改 `handleDragEnd` 自动调用 API 并重新加载数据
**位置**: `newstatic/src/components/progress/GanttChartRefactored.vue:470-500`

### 7. 依赖关系线类型
**问题**: 使用曲线而不是直角折线
**修复**: 实现A*路径规划算法，生成直角折线路径
**位置**: `newstatic/src/utils/ganttHelpers.js:133-420`

---

## 最佳实践

### 前端开发

1. **组件职责单一化**: 每个组件只负责一个功能
2. **使用 Composables**: 复用逻辑抽取到 composables
3. **状态集中管理**: 复杂状态使用 Store 管理
4. **API 调用统一**: 所有API调用通过 `progressApi`
5. **错误处理**: 统一使用 try-catch 和 ElMessage 提示

### 后端开发

1. **RESTful 设计**: 遵循 REST API 设计规范
2. **类型安全**: 注意类型转换（string vs uint）
3. **事务处理**: 复杂操作使用数据库事务
4. **错误日志**: 添加详细的错误日志
5. **权限验证**: 所有API端点验证用户权限

### 代码维护

1. **定期清理**: 删除未使用的代码和组件
2. **文档更新**: 修改代码时同步更新文档
3. **版本控制**: 使用 Git 管理代码版本
4. **代码审查**: 重要修改进行代码审查
5. **单元测试**: 核心功能编写单元测试

---

## 文档维护

本文档应随代码演进定期更新，包括：

- 新增组件和功能
- API 端点变更
- 架构调整
- 已知问题和解决方案
- 性能优化记录

**最后更新**: 2026-02-22
**维护者**: 开发团队
