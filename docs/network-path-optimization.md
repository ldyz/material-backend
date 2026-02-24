# 网络图连线优化指南

## 核心目标

实现清晰、规范的双代号时标网络图连线，确保：
1. **最少弯折** - 减少连线的转角数量，使图形更清晰
2. **右侧出线优先** - 优先从节点右侧出线，符合阅读习惯
3. **箭头指向圆心** - 所有箭头都指向节点圆心，视觉统一
4. **正交路径** - 所有连线只使用水平和垂直线段，转角为90度
5. **任务行水平线段最长** - 优先在任务列表所在行保持水平，标签显示在这一行

## 实现策略

### 0. 线条重叠避免

**函数**: `nodeExitDirectionUsage` computed property

**动态优先级调整**：
```javascript
// 追踪每个节点的出线方向使用情况
const nodeExitDirectionUsage = computed(() => {
  const usage = {} // { nodeId: { right: count, top: count, bottom: count } }

  // 分析任务箭头的出线方向
  taskArrows.value.forEach(task => {
    // 解析路径，确定出线方向
    // 更新对应节点的出线方向计数
  })

  return usage
})

// 在路径计算中使用
const rightPenalty = exitUsage.right * 30  // 已使用越多，惩罚越大
const topPenalty = exitUsage.top * 30
const bottomPenalty = exitUsage.bottom * 30
```

**优先级计算**：
- 基础优先级 - 方向使用惩罚
- 每个已有线条降低30优先级
- 自动选择较少使用的方向

### 1. 任务连线（实工作）

**函数**: `calculateOptimizedPath()`

**路径选择逻辑**：
```javascript
// 策略0: 水平连接（0弯折）
if (dx > 0 && Math.abs(dy) < 10) {
  // 右侧出线 → 水平 → 圆心
  return path
}

// 策略1: 经过任务行（优先）
if (dx > 0) {
  // 计算方案1: 经过任务行的路径
  const pathViaTargetY = 右侧出线 → 垂直到任务行 → 水平 → 垂直到圆心

  // 计算方案2: 直接垂直连接
  const pathDirect = 右侧出线 → 水平 → 垂直到圆心

  // 如果任务行的水平线段足够长(>50px)，使用方案1
  if (horizontalLength > 50) {
    return pathViaTargetY  // 优先：经过任务行
  } else {
    return pathDirect  // 备选：直接连接
  }
}

// 策略2: 终点在上方
if (dy < -10) {
  // 右侧出线 → 水平 → 垂直到任务行 → 水平 → 垂直向上到圆心
  // 这样任务行有最长的水平线段
  return path
}

// 策略3: 终点在下方
if (dy > 10) {
  // 右侧出线 → 水平 → 垂直到任务行 → 水平 → 垂直向下到圆心
  // 这样任务行有最长的水平线段
  return path
}
```

**任务行水平线段优化**：
- 优先经过任务列表所在行（targetY）
- 如果任务行的水平线段长度 > 50px，使用经过任务行的路径
- 任务标签（名称、工期）显示在任务行的水平线段上

### 2. 虚工作连线（依赖关系）

**函数**: `calculateMultiDirectionPath()`

**路径方案优先级**：

| 优先级 | 弯折数 | 出线方向 | 方案描述 |
|--------|--------|----------|----------|
| 100 | 0 | 右侧 | 水平直接到圆心 |
| 95 | 1 | 上/下 | 垂直→水平到圆心 |
| 90 | 1 | 右侧 | 水平→垂直到圆心 |
| 80 | 2 | 右侧 | 向右回绕→垂直到圆心 |
| 70 | 2 | 上/下 | 垂直→水平→垂直到圆心 |

**选择算法**：
```javascript
pathOptions.sort((a, b) => {
  if (a.bends !== b.bends) {
    return a.bends - b.bends  // 弯折少的优先
  }
  if (a.priority !== b.priority) {
    return b.priority - a.priority  // 优先级高的优先
  }
  return a.length - b.length  // 路径短的优先
})
```

### 3. 箭头标记配置

**关键参数**: `refX = nodeRadius + arrowLength - 2`

```html
<marker
  :refX="nodeRadius + 7"  <!-- 任务箭头，箭头长度9px -->
  refY="3.5"
  orient="auto"
  markerUnits="userSpaceOnUse"
>
```

**工作原理**：
- 路径终点计算到圆心 `(toX, toY)`
- 箭头路径 `M0,0 L0,7 L9,3.5 z` 的尖端在 `x=9`
- `refX = nodeRadius + 7` 使得箭头尖端刚好接触圆边缘
- `nodeRadius + 7 - 9 = nodeRadius - 2`，箭头尖端距离圆心正好是半径
- 视觉效果：箭头指向圆心，尖端接触圆边缘

## 典型场景

### 场景1: 终点在右侧

```
最优路径（0弯折）：
[节点A] ─────────────────→ ●节点B
```

### 场景2: 终点在右侧但Y相差较大

```
最优路径（经过任务行，任务行水平线段最长）：
[节点A]
        │
        ├─────────────────┐
        │                 │
任务行 ├─────────────────→ ●节点B ← 标签显示在这里
        │                 │
        └─────────────────┘
```

### 场景3: 终点在左侧

```
最优路径（2弯折，右侧出线回绕）：
[节点A] ──┐
          │
          └──────→ ●节点B
```

### 场景4: 终点在上方

```
最优路径（2弯折，经过任务行）：
[节点A] ──┐
          │
          ├──────┐
          │      │
任务行     ├──────→ ●节点B ← 标签显示在这里
          │
    (向上绕行，经过任务行)
```

### 场景5: 任务行水平线段优化

```
任务A（行1）:
[节点A] ──────┐
             │
             ├───────── 任务行水平线段（最长）
             │             ↓
任务B（行5）:    ●节点B

标签"任务B"显示在任务行的水平线段上，而不是节点B的位置
```

## 代码示例

### 创建新的路径方案

```javascript
// 定义路径方案
pathOptions.push({
  name: '方案名称',
  bends: 弯折数,        // 0, 1, 2
  startX: 起点X,
  startY: 起点Y,
  endX: 圆心X,
  endY: 圆心Y,
  midPoints: [        // 中间点（如果有）
    { x: ..., y: ... }
  ],
  length: 路径长度,    // 用于优化
  priority: 优先级     // 200-50，越高越优先
})
```

### 提取任务行上的标签位置

```javascript
// 从路径中提取标签Y坐标，优先使用任务行
function extractLabelYFromPath(path, startX, endX, targetY) {
  const commands = path.match(/[MLH][^,]*/g)

  let targetYHorizontalLength = 0  // 任务行的水平线段长度
  let longestHorizontalY = null

  // 解析路径命令
  for (let i = 0; i < commands.length; i++) {
    const cmd = commands[i][0]
    const coords = commands[i].slice(1).split(/[, ]+/).map(Number)

    if (cmd === 'L') {
      const prevY = currentY
      const currentY = coords[1]

      // 检查是否是水平线段
      if (prevY === currentY) {
        const length = Math.abs(currentX - prevX)

        // 优先选择任务行上的水平线段
        if (Math.abs(prevY - targetY) < 5) {
          if (length > targetYHorizontalLength) {
            targetYHorizontalLength = length
            longestHorizontalY = prevY
          }
        }
      }
    }
  }

  return longestHorizontalY
}

// 使用示例
const labelY = extractLabelYFromPath(realPath, startX, endX, taskTargetY)
```

### 优先经过任务行的路径计算

```javascript
function calculateOptimizedPath(fromX, fromY, toX, toY, targetY, radius) {
  const dx = toX - fromX

  // 终点在右侧但Y相差较大
  if (dx > 0 && Math.abs(toY - fromY) >= 10) {
    const startX = fromX + radius

    // 方案1: 经过任务行的路径
    const pathViaTarget = `M ${startX} ${fromY} L ${startX} ${targetY} L ${toX} ${targetY} L ${toX} ${toY}`

    // 方案2: 直接垂直连接
    const bendX = fromX + dx * 0.6
    const pathDirect = `M ${startX} ${fromY} L ${bendX} ${fromY} L ${bendX} ${toY} L ${toX} ${toY}`

    // 计算水平线段长度
    const horizontalLength = toX - startX

    // 如果任务行的水平线段足够长，使用经过任务行的路径
    if (horizontalLength > 50) {
      return pathViaTarget  // 标签显示在任务行上
    } else {
      return pathDirect  // 直接连接
    }
  }
}
```

### 优先级设置指南

**任务连线**：
- 0弯折：无需设置priority（直接返回）
- 1弯折：优先使用，无需设置priority
- 2弯折：作为备选方案

**虚工作连线**：
- 右侧出线，0弯折：priority = 200
- 右侧出线，1弯折：priority = 190
- 右侧出线，2弯折：priority = 150
- 上/下出线，1弯折：priority = 100
- 上/下出线，2弯折：priority = 50

## 调试技巧

### 1. 检查路径弯折数

```javascript
console.log(`Path: ${path}`)
console.log(`Bends: ${selected.bends}`)
console.log(`Priority: ${selected.priority}`)
```

### 2. 验证正交约束

```javascript
// 解析路径，检查所有线段
const commands = path.match(/[MLH][^,]*/g)
commands.forEach(cmd => {
  const coords = cmd.slice(1).split(/[, ]+/).map(Number)
  // 验证只有水平或垂直线段
  if (cmd === 'L') {
    const isHorizontal = (coords[1] === currentY)
    const isVertical = (coords[0] === currentX)
    if (!isHorizontal && !isVertical) {
      console.error('非正交线段！')
    }
  }
})
```

### 3. 检查箭头方向

```javascript
// 验证箭头指向圆心
const lastCmd = commands[commands.length - 1]
const endX = lastCmd.coords[0]
const endY = lastCmd.coords[1]
console.log(`Arrow ends at: (${endX}, ${endY})`)
console.log(`Target center: (${toX}, ${toY})`)
```

### 4. 检查任务行水平线段

```javascript
// 检查路径是否经过任务行
function checkPathPassesThroughTargetRow(path, targetY) {
  const commands = path.match(/[MLH][^,]*/g)
  let hasTargetRowSegment = false

  for (let i = 0; i < commands.length; i++) {
    const cmd = commands[i][0]
    const coords = commands[i].slice(1).split(/[, ]+/).map(Number)

    if (cmd === 'L') {
      const prevY = parseFloat(commands[i-1].slice(1).split(/[, ]+/)[1])
      const currY = coords[1]

      // 检查是否在任务行上
      if (prevY === currY && Math.abs(prevY - targetY) < 5) {
        hasTargetRowSegment = true
        console.log('路径经过任务行:', prevY, '≈', targetY)
      }
    }
  }

  return hasTargetRowSegment
}

// 使用示例
const passesThroughTargetRow = checkPathPassesThroughTargetRow(realPath, taskTargetY)
console.log('经过任务行:', passesThroughTargetRow)
```

## 最佳实践

1. **优先使用右侧出线** - 符合从左到右的阅读习惯
2. **最少弯折原则** - 0弯折 > 1弯折 > 2弯折
3. **避免交叉** - 通过调整出线方向减少连线交叉
4. **保持一致性** - 相同场景使用相同路径策略
5. **测试边界情况** - 节点重叠、极近、极远等情况

## 常见问题

### Q: 为什么有些连线不使用右侧出线？

A: 当终点在起点左侧很近的位置时，右侧出线会产生很长的回绕路径，此时会选择上/下出线。

### Q: 如何调整连线优先级？

A: 修改路径方案的 `priority` 值，值越大优先级越高。各方向优先级应保持较平衡（70-100之间），避免强制使用右侧出线导致线条重叠。

### Q: 箭头被节点遮挡怎么办？

A: 检查 `refX` 设置，应为 `nodeRadius + 7`（任务箭头）或 `nodeRadius + 6`（依赖箭头）。确保路径终点在圆心 `(toX, toY)`。箭头路径尖端在x=9（或x=8），refX需要减去箭头长度来计算正确位置。

### Q: 连线转角不是90度？

A: 检查路径计算，确保只使用水平或垂直线段。所有线段应满足 `dx === 0 || dy === 0`。

### Q: 如何避免连线重叠？

A: 系统通过动态优先级调整自动避免重叠：
1. `nodeExitDirectionUsage` 追踪每个节点的出线方向使用情况
2. 路径计算时，对已使用的出线方向降低优先级（每条线降低30优先级）
3. 新路径会自动选择未使用或使用较少的出线方向
4. 这样可以有效减少线条重叠

A: 检查路径计算，确保只使用水平或垂直线段。所有线段应满足 `dx === 0 || dy === 0`。
