# 网络图节点和箭线逻辑归纳

## 📐 节点创建逻辑 (eventNodes computed)

### 第一步：收集任务依赖关系
```javascript
const taskPredecessors = new Map()
// key: taskId, value: predecessors[]
// 例如：Map { 1 => [2, 3], 2 => [], 3 => [1] }
```

### 第二步：为每个任务创建独立节点
每个任务创建2个节点（开始节点 + 结束节点）：

```javascript
任务A (2025-01-01 ~ 2025-01-10)
  → startNode_A (date: "2025-01-01", tasks.start: [A])
  → endNode_A (date: "2025-01-10", tasks.end: [A])

任务B (2025-01-10 ~ 2025-01-15)
  → startNode_B (date: "2025-01-10", tasks.start: [B])
  → endNode_B (date: "2025-01-15", tasks.end: [B])
```

### 第三步：智能合并节点

**合并规则：**
```javascript
条件1: 节点A.date === 节点B.date  // 同一天
条件2: 有直接依赖关系              // A依赖于B 或 B依赖于A
满足以上两个条件才合并
```

**合并示例：**

```javascript
// 场景1：应该合并
任务A结束于2025-01-10，任务B开始于2025-01-10
A.predecessors.includes(B.id)  // A依赖于B
→ endNode_A 和 startNode_B 合并为一个节点

// 场景2：不合并
任务C结束于2025-01-10，任务D开始于2025-01-10
C.predecessors 不包含 D.id
D.predecessors 不包含 C.id
→ endNode_C 和 startNode_D 保持独立
```

### 第四步：计算Y坐标

为同一天的多个节点添加垂直偏移，避免重叠：

```javascript
// 统计每个日期的节点数量
const dateGroupCount = new Map()
// 记录每个节点在日期组中的索引
const dateNodeIndex = new Map()

// 计算Y坐标
if (同一天有多个节点) {
  offsetInRange = nodeIndexInSameDate - (nodeCountInSameDate - 1) / 2
  baseTaskIndex += offsetInRange * 1.5  // 每个节点间隔1.5行
}
```

---

## 🔗 箭头绘制逻辑

### 1. 任务箭头 (taskArrows) - 实工作线

**作用：** 连接每个任务的开始节点和结束节点

**查找逻辑（已修复）：**
```javascript
// 查找开始节点：优先选择独立节点
let startNode = eventNodes.value.find(node =>
  node.date === taskStartDate &&
  node.tasks.start.includes(task.id) &&
  node.tasks.start.length === 1  // ✅ 优先选择只包含该任务的节点
)

// 如果找不到独立节点，选择合并节点
if (!startNode) {
  startNode = eventNodes.value.find(node =>
    node.date === taskStartDate &&
    node.tasks.start.includes(task.id)
  )
}

// 同样的逻辑用于结束节点
```

**为什么需要优先选择独立节点？**

当同一天有多个节点时：
```
节点1: { date: "2025-01-10", tasks.end: [任务A] }
节点2: { date: "2025-01-10", tasks.end: [任务B] }
```

如果不加筛选，`find` 会找到节点1或节点2（取决于数组顺序），可能导致：
- 任务A的箭头指向节点2（错误！）
- 任务B的箭头指向节点1（错误！）

加上 `node.tasks.end.length === 1` 筛选后：
- 任务A只匹配节点1 ✅
- 任务B只匹配节点2 ✅

### 2. 依赖关系箭头 (renderableDependencies) - 虚工作

**作用：** 连接前置任务的结束节点和后置任务的开始节点

**查找逻辑（已修复）：**
```javascript
// 前置任务：找到结束节点
let fromNode = eventNodes.value.find(node =>
  node.date === fromTaskEndDate &&
  node.tasks.end.includes(fromTask.id) &&
  node.tasks.end.length === 1  // ✅ 优先选择独立节点
)

// 后置任务：找到开始节点
let toNode = eventNodes.value.find(node =>
  node.date === toTaskStartDate &&
  node.tasks.start.includes(toTask.id) &&
  node.tasks.start.length === 1  // ✅ 优先选择独立节点
)
```

---

## 🎯 完整示例

### 数据示例：
```javascript
任务列表:
  - 任务1: 2025-01-01 ~ 2025-01-10, predecessors: []
  - 任务2: 2025-01-10 ~ 2025-01-15, predecessors: [1]
  - 任务3: 2025-01-10 ~ 2025-01-20, predecessors: []
```

### 节点创建过程：

**第二步：创建独立节点**
```
startNode_1 (date: "2025-01-01", tasks.start: [1])
endNode_1 (date: "2025-01-10", tasks.end: [1])
startNode_2 (date: "2025-01-10", tasks.start: [2])
endNode_2 (date: "2025-01-15", tasks.end: [2])
startNode_3 (date: "2025-01-10", tasks.start: [3])
endNode_3 (date: "2025-01-20", tasks.end: [3])
```

**第三步：合并节点**
```
endNode_1 和 startNode_2 在同一天(2025-01-10)
且 任务1.predecessors 不包含 任务2
但 任务2.predecessors 包含 任务1  // 任务2依赖于任务1
→ 合并为节点1-2: {date: "2025-01-10", tasks.end: [1], tasks.start: [2]}

startNode_3 和 endNode_1/startNode_2 在同一天
但 任务3 没有依赖关系
→ startNode_3 保持独立: {date: "2025-01-10", tasks.start: [3]}
```

**最终节点：**
```
节点1: {date: "2025-01-01", tasks.start: [1]}
节点1-2: {date: "2025-01-10", tasks.end: [1], tasks.start: [2]}  // 合并节点
节点3: {date: "2025-01-10", tasks.start: [3]}  // 独立节点
节点2: {date: "2025-01-15", tasks.end: [2]}
节点3: {date: "2025-01-20", tasks.end: [3]}
```

**第四步：Y坐标分配**
```
节点1-2 和 节点3 在同一天
→ 节点1-2: baseTaskIndex + (-0.5 * 1.5)  // 上移
→ 节点3:   baseTaskIndex + (0.5 * 1.5)   // 下移
```

### 箭头绘制：

**任务箭头：**
```
任务1: startNode_1 → 节点1-2 (使用 endNode_1 的位置)
任务2: 节点1-2 (使用 startNode_2 的位置) → endNode_2
任务3: startNode_3 (节点3) → endNode_3
```

**依赖箭头：**
```
任务2依赖于任务1:
  from: 节点1-2 (包含endNode_1, end.length=2)  // 合并节点
  to: 节点1-2 (包含startNode_2, start.length=2) // 同一个合并节点

  ✅ 优先选择独立节点...找不到合并节点，使用合并节点
```

---

## ⚠️ 边界情况处理

### 情况1：合并后的节点查找
```javascript
// 问题：合并节点包含多个任务，如何确保找到正确的节点？

解决方案：
// 优先选择独立节点 (length === 1)
// 如果找不到，再选择合并节点
// 因为合并节点本身就表示这些任务有关联关系
```

### 情况2：同一天多个独立节点
```javascript
// 问题：多个节点同一天但都不合并（无依赖关系）

解决方案：
// 使用 node.tasks.start/end.length === 1 筛选独立节点
// 加上任务ID匹配，确保每个任务找到自己的节点
```

### 情况3：Y坐标冲突
```javascript
// 问题：同一天多个节点Y坐标重叠

解决方案：
// 计算该日期的节点总数
// 为每个节点分配偏移：offsetInRange * 1.5
// 均匀分布在 baseTaskIndex 周围
```

---

## 📊 算法复杂度

- 节点创建：O(n²) - n为任务数，需要检查所有节点对
- 节点查找：O(m) - m为节点数，使用find查找
- 总体复杂度：O(n²) - 对于中小型项目（<1000任务）可接受

---

## 🔧 未来优化方向

1. **性能优化**
   - 使用 Map 代替数组 find 查找节点
   - 预先建立 taskId → node 的映射

2. **Y坐标优化**
   - 使用图布局算法自动计算最优Y坐标
   - 考虑任务依赖链的层级关系

3. **可视化优化**
   - 添加节点重叠检测和自动调整
   - 支持手动调整节点位置后保存

---

**文档版本：** v1.0
**更新日期：** 2026-02-24
**作者：** Claude AI Assistant
