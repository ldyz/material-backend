/**
 * 甘特图 TypeScript 类型定义
 * 提供完整的类型安全支持
 */

// ==================== 枚举类型 ====================

/**
 * 任务状态枚举
 */
export enum TaskStatus {
  NOT_STARTED = 'not_started',
  IN_PROGRESS = 'in_progress',
  COMPLETED = 'completed',
  DELAYED = 'delayed'
}

/**
 * 任务优先级枚举
 */
export enum TaskPriority {
  LOW = 'low',
  MEDIUM = 'medium',
  HIGH = 'high'
}

/**
 * 视图模式枚举
 */
export enum ViewMode {
  DAY = 'day',
  WEEK = 'week',
  MONTH = 'month',
  QUARTER = 'quarter'
}

/**
 * 分组模式枚举
 */
export enum GroupMode {
  NONE = '',
  STATUS = 'status',
  PRIORITY = 'priority'
}

/**
 * 拖拽模式枚举
 */
export enum DragMode {
  NONE = 'none',
  MOVE = 'move',
  RESIZE_LEFT = 'resize_left',
  RESIZE_RIGHT = 'resize_right'
}

/**
 * 依赖关系类型枚举
 */
export enum DependencyType {
  FINISH_TO_START = 'finish_to_start',  // 完成-开始 (FS)
  START_TO_START = 'start_to_start',    // 开始-开始 (SS)
  FINISH_TO_FINISH = 'finish_to_finish',// 完成-完成 (FF)
  START_TO_FINISH = 'start_to_finish'   // 开始-完成 (SF)
}

// ==================== 基础类型 ====================

/**
 * 任务依赖关系
 */
export interface GanttTaskDependency {
  id: number | string
  predecessor_id: number | string
  successor_id: number | string
  type: DependencyType
  lag?: number  // 延滞天数（负数表示提前）
}

/**
 * 资源分配
 */
export interface GanttResource {
  id: number | string
  name: string
  type: 'human' | 'material' | 'equipment'
  unit?: string
  cost_per_unit?: number
  allocated_units?: number
}

/**
 * 基准计划（用于对比）
 */
export interface GanttBaseline {
  start: string
  end: string
  duration: number
}

// ==================== 核心任务类型 ====================

/**
 * 甘特图任务（完整版）
 */
export interface GanttTask {
  id: number | string
  name: string
  start: string           // ISO 日期字符串 YYYY-MM-DD
  end: string             // ISO 日期字符串 YYYY-MM-DD
  startDate: Date         // JavaScript Date 对象（运行时计算）
  endDate: Date           // JavaScript Date 对象（运行时计算）
  duration: number        // 工期（天数）
  progress: number        // 进度百分比 0-100
  status: TaskStatus
  priority: TaskPriority
  is_critical: boolean    // 是否关键路径
  is_milestone: boolean   // 是否里程碑

  // 层级关系
  parent_id?: number | string | null
  sort_order: number
  level?: number          // 树形层级深度

  // 依赖关系
  predecessors: GanttTaskDependency[]
  successors: GanttTaskDependency[]

  // 资源分配
  resources: GanttResource[]

  // 基准计划
  baseline?: GanttBaseline

  // UI 状态（运行时）
  is_visible?: boolean
  is_expanded?: boolean   // 是否展开（树形结构）
  is_selected?: boolean
  is_highlighted?: boolean

  // 额外数据
  description?: string
  assignee?: string
  tags?: string[]
  color?: string
}

// ==================== 时间轴相关类型 ====================

/**
 * 时间轴日期
 */
export interface TimelineDay {
  date: string           // ISO 日期字符串
  day: number            // 日期数字 (1-31)
  weekday: string        // 星期几的中文
  isToday: boolean       // 是否今天
  isWeekend: boolean     // 是否周末
  position: number       // X轴位置（像素）
}

/**
 * 时间轴周
 */
export interface TimelineWeek {
  weekNumber: number     // 周数
  startDate: string      // 开始日期
  endDate: string        // 结束日期
  position: number       // X轴位置
  width: number          // 宽度（像素）
}

/**
 * 时间轴月
 */
export interface TimelineMonth {
  month: number          // 月份 (1-12)
  year: number           // 年份
  name: string           // 中文名称
  startDate: string      // 开始日期
  endDate: string        // 结束日期
  position: number       // X轴位置
  width: number          // 宽度（像素）
}

/**
 * 时间轴季度
 */
export interface TimelineQuarter {
  quarter: number        // 季度 (1-4)
  year: number           // 年份
  name: string           // 中文名称
  startDate: string      // 开始日期
  endDate: string        // 结束日期
  position: number       // X轴位置
  width: number          // 宽度（像素）
}

// ==================== 视图配置类型 ====================

/**
 * 视图模式配置
 */
export interface ViewModeConfig {
  minWidth: number       // 最小日宽度
  maxWidth: number       // 最大日宽度
  default: number        // 默认日宽度
  label: string          // 显示标签
}

/**
 * 视图配置集合
 */
export interface ViewConfigs {
  [key: string]: ViewModeConfig
}

// ==================== 分组相关类型 ====================

/**
 * 任务分组
 */
export interface TaskGroup {
  key: string            // 分组键
  label: string          // 显示名称
  tasks: GanttTask[]     // 组内任务
  count: number          // 任务数量
  is_collapsed: boolean  // 是否折叠
  color?: string         // 分组颜色
}

/**
 * 分组任务集合
 */
export interface GroupedTasks {
  [key: string]: TaskGroup
}

// ==================== 拖拽相关类型 ====================

/**
 * 拖拽预览任务
 */
export interface DragPreviewTask {
  id: number | string
  start: string
  end: string
  duration: number
}

/**
 * 拖拽提示框位置
 */
export interface TooltipPosition {
  x: number
  y: number
}

// ==================== 虚拟滚动相关类型 ====================

/**
 * 虚拟滚动可见范围
 */
export interface VisibleRange {
  start: number
  end: number
}

/**
 * 虚拟滚动选项
 */
export interface VirtualScrollOptions {
  itemHeight: number        // 每项高度
  containerHeight: number   // 容器高度
  overscan?: number         // 额外渲染数量
}

// ==================== 统计相关类型 ====================

/**
 * 任务统计信息
 */
export interface TaskStats {
  total: number            // 总任务数
  completed: number        // 已完成
  inProgress: number       // 进行中
  notStarted: number       // 未开始
  delayed: number          // 延期
  critical: number         // 关键路径
  progressRate: number     // 平均进度
  onTime: number           // 按时完成
  overdue: number          // 已逾期
}

// ==================== 容器尺寸类型 ====================

/**
 * 容器尺寸
 */
export interface ContainerSize {
  width: number | null
  height: number | null
}

// ==================== 上下文菜单类型 ====================

/**
 * 上下文菜单位置
 */
export interface ContextMenuPosition {
  x: number
  y: number
}

// ==================== 关键路径相关类型 ====================

/**
 * 关键路径节点
 */
export interface CriticalPathNode {
  taskId: number | string
  earliestStart: number    // 最早开始时间（天数）
  earliestFinish: number   // 最早完成时间
  latestStart: number      // 最晚开始时间
  latestFinish: number     // 最晚完成时间
  slack: number            // 总时差
  isCritical: boolean      // 是否在关键路径上
}

// ==================== 资源管理类型 ====================

/**
 * 资源库
 */
export interface ResourceLibrary {
  [id: string]: GanttResource
}

/**
 * 资源分配请求
 */
export interface ResourceAllocationRequest {
  taskId: number | string
  resourceId: number | string
  units: number
}

// ==================== 过滤和搜索类型 ====================

/**
 * 任务筛选条件
 */
export interface TaskFilterOptions {
  status?: TaskStatus[]
  priority?: TaskPriority[]
  assignee?: string[]
  dateRange?: {
    start: string
    end: string
  }
  hasKeyword?: string
}

// ==================== 导出相关类型 ====================

/**
 * 导出选项
 */
export interface ExportOptions {
  format: 'png' | 'pdf' | 'excel'
  includeLegend?: boolean
  includeStats?: boolean
  filename?: string
  quality?: number
}

// ==================== 主题相关类型 ====================

/**
 * 主题模式
 */
export type ThemeMode = 'light' | 'dark' | 'auto'

/**
 * 主题配置
 */
export interface ThemeConfig {
  mode: ThemeMode
  primaryColor: string
  fontSize: 'small' | 'medium' | 'large'
}

// ==================== 响应式断点类型 ====================

/**
 * 断点名称
 */
export type Breakpoint = 'xs' | 'sm' | 'md' | 'lg' | 'xl' | '2xl'

/**
 * 断点配置
 */
export interface BreakpointConfig {
  [key: string]: string
}

// ==================== 触摸手势类型 ====================

/**
 * 触摸手势类型
 */
export type TouchGestureType = 'swipe_left' | 'swipe_right' | 'swipe_up' | 'swipe_down' | 'pinch' | 'tap' | 'long_press'

/**
 * 触摸点
 */
export interface TouchPoint {
  x: number
  y: number
}

/**
 * 触摸手势处理器
 */
export interface TouchGestureHandlers {
  onSwipeLeft?: () => void
  onSwipeRight?: () => void
  onSwipeUp?: () => void
  onSwipeDown?: () => void
  onPinch?: (scale: number) => void
  onTap?: () => void
  onLongPress?: () => void
}

// ==================== API 响应类型 ====================

/**
 * 进度数据响应
 */
export interface ScheduleDataResponse {
  activities: { [id: string]: any }
  nodes: { [id: string]: any }
}

/**
 * 任务更新请求
 */
export interface TaskUpdateRequest {
  id: number | string
  name?: string
  start?: string
  end?: string
  progress?: number
  status?: TaskStatus
  priority?: TaskPriority
  parent_id?: number | string | null
  description?: string
  assignee?: string
}

/**
 * 依赖关系创建请求
 */
export interface DependencyCreateRequest {
  predecessor_id: number | string
  successor_id: number | string
  type: DependencyType
  lag?: number
}

// ==================== 事件载荷类型 ====================

/**
 * 任务点击事件载荷
 */
export interface TaskClickPayload {
  task: GanttTask
  event: MouseEvent
}

/**
 * 任务拖拽事件载荷
 */
export interface TaskDragPayload {
  task: GanttTask
  originalTask: GanttTask
  dragMode: DragMode
}

/**
 * 依赖关系创建事件载荷
 */
export interface DependencyCreatePayload {
  sourceTaskId: number | string
  targetTaskId: number | string
  type: DependencyType
}

// ==================== 错误类型 ====================

/**
 * 甘特图错误类型
 */
export enum GanttErrorType {
  INVALID_DATE_RANGE = 'invalid_date_range',
  CIRCULAR_DEPENDENCY = 'circular_dependency',
  TASK_NOT_FOUND = 'task_not_found',
  INVALID_DEPENDENCY = 'invalid_dependency'
}

/**
 * 甘特图错误
 */
export class GanttError extends Error {
  constructor(
    public type: GanttErrorType,
    message: string,
    public details?: any
  ) {
    super(message)
    this.name = 'GanttError'
  }
}
