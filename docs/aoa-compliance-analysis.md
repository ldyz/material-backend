# AOA网络图规范符合性分析

## 当前实现 vs R11规范

### ❌ 不符合R11规范的部分

**R11规范要求**：
> 相同紧前条件 → 共享开始节点
> 相同紧后条件 → 共享结束节点

**当前实现**：
- 每个任务都有独立的开始和结束节点
- 节点ID格式：`node-${taskId}-start` 和 `node-${taskId}-end`
- 不共享节点，即使多个任务有相同的紧前/紧后条件

**示例**：
```javascript
// 场景：任务B和任务C都依赖任务A
// R11要求：任务B和任务C应该共享同一个开始节点
// 当前实现：
//   node-A-end (任务A的结束)
//   node-B-start (任务B的开始) ❌ 独立节点
//   node-C-start (任务C的开始) ❌ 独立节点

// R11规范应该是：
//   node-A-end
//   node-BC-start (任务B和C共享的开始节点) ✅
```

**为什么不共享节点？**
- 为了确保每个任务节点都能与任务列表行对齐
- 简化节点查找逻辑
- 提供更好的可视化效果

**权衡**：
- ✅ 优点：每个任务独立显示，易于理解
- ❌ 缺点：违反AOA规范，节点数量翻倍

---

## 甘特图逻辑关系在AOA中的表示

### 1. FS (Finish-to-Start) - 完成-开始

**定义**：前置任务完成后，后续任务才能开始

**AOA表示**：
```
任务A (前置) ─────●───> ●──── 任务B (后续)
                  ↑
            直接连接，无需虚工作
```

**当前实现**：✅ **已支持**
```javascript
// renderableDependencies computed
fromNodeId = `node-${fromTask.id}-end`  // 前置任务结束
toNodeId = `node-${toTask.id}-start`     // 后续任务开始
```

**是否需要虚工作**：❌ 不需要

---

### 2. SS (Start-to-Start) - 开始-开始

**定义**：前置任务开始后，后续任务才能开始

**AOA表示**：
```
任务A ──●────────────> ●─── 任务A结束
        │              │
        └──[虚工作]───> ●─── 任务B开始
                        │
                        └──> ●─── 任务B结束
```

**当前实现**：❌ **不支持**
```javascript
// 当前代码没有处理 dep.type === 'SS'
// 仍然连接 fromTask.end → toTask.start
```

**需要的改动**：
```javascript
if (dep.type === 'SS') {
  fromNodeId = `node-${fromTask.id}-start`  // 前置任务开始
  toNodeId = `node-${toTask.id}-start`      // 后续任务开始
  // 需要添加虚工作节点
}
```

**是否需要虚工作**：✅ 需要
- 添加虚工作节点：`dummy-${fromTask.id}-${toTask.id}-SS`

---

### 3. FF (Finish-to-Finish) - 完成-完成

**定义**：前置任务完成后，后续任务才能完成

**AOA表示**：
```
任务A ────────────> ●─── 任务A结束
                    │
                    └──[虚工作]───> ●─── 任务B结束
                                  ↑
任务B ──●──────────────────────────┘
       开始
```

**当前实现**：❌ **不支持**

**需要的改动**：
```javascript
if (dep.type === 'FF') {
  fromNodeId = `node-${fromTask.id}-end`  // 前置任务结束
  toNodeId = `node-${toTask.id}-end`      // 后续任务结束
  // 需要添加虚工作节点
}
```

**是否需要虚工作**：✅ 需要
- 添加虚工作节点：`dummy-${fromTask.id}-${toTask.id}-FF`

---

### 4. SF (Start-to-Finish) - 开始-完成

**定义**：前置任务开始后，后续任务才能完成

**AOA表示**：
```
任务A ──●──────────────────────────> ●─── 任务A结束
        │                            ↑
        └──[虚工作]───────────────────┘
                                    ●─── 任务B结束
                                    ↑
任务B ──────────────────────────────┘
       开始
```

**当前实现**：❌ **不支持**

**需要的改动**：
```javascript
if (dep.type === 'SF') {
  fromNodeId = `node-${fromTask.id}-start`  // 前置任务开始
  toNodeId = `node-${toTask.id}-end`        // 后续任务结束
  // 需要添加虚工作节点
}
```

**是否需要虚工作**：✅ 需要
- 添加虚工作节点：`dummy-${fromTask.id}-${toTask.id}-SF`

---

## 时间滞后 (Lag) 的处理

### 定义
时间滞后表示在依赖关系基础上额外的时间延迟。

### 示例
```javascript
// FS + 2天滞后
// 意味着：任务A完成后2天，任务B才能开始
dep = {
  depends_on: taskId_A,
  task_id: taskId_B,
  type: 'FS',
  lag: 2  // 2天滞后
}
```

### AOA表示
```
任务A ─────●───[虚工作:2天]───> ●─── 任务B
                            ↑
                        虚工作表示时间滞后
```

**当前实现**：❌ **不支持 lag 参数**
```javascript
// 当前代码没有使用 dep.lag
```

**需要的改动**：
```javascript
// 在计算虚工作节点位置时考虑 lag
dummyX = fromNode.x + dep.lag * props.dayWidth
```

---

## 实现方案建议

### 方案A：严格遵循AOA规范（R11）

**优点**：
- 符合工程规范
- 节点数量最少
- 理论上更优

**缺点**：
- 节点与任务列表不对齐
- 需要复杂的节点合并逻辑
- 可视化效果较差

### 方案B：当前方案（独立节点）+ 支持所有关系类型

**优点**：
- 每个任务独立显示
- 与任务列表完美对齐
- 实现相对简单

**缺点**：
- 违反R11规范
- 节点数量较多

### 方案C：混合方案

**核心思想**：
- 保留独立节点用于可视化
- 在内部逻辑上仍符合R11
- 渲染时为每个任务创建视觉副本

**实现**：
```javascript
// 逻辑层：符合R11的共享节点
const logicalNodes = createSharedNodeGraph() // R11规范

// 视觉层：为每个任务创建视觉节点
const visualNodes = logicalNodes.flatMap(node =>
  node.tasks.map(taskId => ({
    ...node,
    id: `visual-${taskId}`,
    taskId
  }))
)
```

---

## 总结表

| 关系类型 | 当前支持 | 需要虚工作 | 实现难度 | 优先级 |
|---------|---------|-----------|---------|--------|
| FS      | ✅ 是    | ❌ 否     | 简单    | -      |
| SS      | ❌ 否    | ✅ 是     | 中等    | 高    |
| FF      | ❌ 否    | ✅ 是     | 中等    | 高    |
| SF      | ❌ 否    | ✅ 是     | 复杂    | 低    |
| Lag     | ❌ 否    | ✅ 是     | 简单    | 中    |

---

## 推荐实现顺序

1. **Phase 1**: 添加SS支持（最常用）
2. **Phase 2**: 添加FF支持
3. **Phase 3**: 添加Lag支持
4. **Phase 4**: 添加SF支持（较少使用）
5. **Phase 5**: 考虑R11规范兼容性（可选）
