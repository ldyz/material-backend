/**
 * 依赖关系图优化工具
 * 将 O(n²) 的依赖关系计算优化到 O(n + m)
 * 使用邻接表存储和拓扑排序算法
 */

import type { GanttTask, GanttTaskDependency, CriticalPathNode, DependencyType } from '@/types/gantt'

/**
 * 依赖关系图的边
 */
interface GraphEdge {
  from: string | number  // 前驱任务ID
  to: string | number    // 后继任务ID
  type: DependencyType   // 依赖类型
  lag: number            // 延滞天数
}

/**
 * 邻接表节点
 */
interface AdjacencyNode {
  taskId: string | number
  outgoing: GraphEdge[]  // 出边（后继）
  incoming: GraphEdge[]  // 入边（前驱）
}

/**
 * 关键路径计算结果
 */
interface CriticalPathResult {
  path: (string | number)[]  // 关键路径上的任务ID
  duration: number            // 总工期
  nodes: Map<string | number, CriticalPathNode>
}

/**
 * 依赖关系图类
 * 使用邻接表存储，提供高效的依赖关系查询和计算
 */
export class DependencyGraph {
  // 邻接表
  private adjacencyList: Map<string | number, AdjacencyNode> = new Map()

  // 任务ID到任务的映射
  private taskMap: Map<string | number, GanttTask> = new Map()

  /**
   * 构建依赖关系图
   * @param tasks 任务列表
   */
  buildGraph(tasks: GanttTask[]): void {
    this.adjacencyList.clear()
    this.taskMap.clear()

    // 初始化所有节点
    tasks.forEach(task => {
      this.taskMap.set(task.id, task)
      this.adjacencyList.set(task.id, {
        taskId: task.id,
        outgoing: [],
        incoming: []
      })
    })

    // 添加边
    tasks.forEach(task => {
      const node = this.adjacencyList.get(task.id)
      if (!node) return

      // 添加后继依赖（出边）
      task.successors.forEach(dep => {
        const edge: GraphEdge = {
          from: task.id,
          to: dep.successor_id,
          type: dep.type,
          lag: dep.lag || 0
        }
        node.outgoing.push(edge)

        // 在后继节点添加入边
        const successorNode = this.adjacencyList.get(dep.successor_id)
        if (successorNode) {
          successorNode.incoming.push(edge)
        }
      })

      // 添加前驱依赖（入边，用于双向查询）
      task.predecessors.forEach(dep => {
        const edge: GraphEdge = {
          from: dep.predecessor_id,
          to: task.id,
          type: dep.type,
          lag: dep.lag || 0
        }
        // 确保不重复添加
        if (!node.incoming.some(e => e.from === dep.predecessor_id && e.to === task.id)) {
          node.incoming.push(edge)
        }

        // 在前驱节点添加出边
        const predecessorNode = this.adjacencyList.get(dep.predecessor_id)
        if (predecessorNode) {
          if (!predecessorNode.outgoing.some(e => e.from === dep.predecessor_id && e.to === task.id)) {
            predecessorNode.outgoing.push(edge)
          }
        }
      })
    })
  }

  /**
   * 获取任务的所有后继任务
   * @param taskId 任务ID
   * @returns 后继任务ID数组
   */
  getSuccessors(taskId: string | number): (string | number)[] {
    const node = this.adjacencyList.get(taskId)
    return node?.outgoing.map(edge => edge.to) || []
  }

  /**
   * 获取任务的所有前驱任务
   * @param taskId 任务ID
   * @returns 前驱任务ID数组
   */
  getPredecessors(taskId: string | number): (string | number)[] {
    const node = this.adjacencyList.get(taskId)
    return node?.incoming.map(edge => edge.from) || []
  }

  /**
   * 检测循环依赖
   * 使用 DFS 深度优先搜索
   * @returns 是否存在循环依赖
   */
  detectCircularDependency(): boolean {
    const visited = new Set<string | number>()
    const recursionStack = new Set<string | number>()

    const dfs = (taskId: string | number): boolean => {
      visited.add(taskId)
      recursionStack.add(taskId)

      const node = this.adjacencyList.get(taskId)
      if (node) {
        for (const edge of node.outgoing) {
          if (!visited.has(edge.to)) {
            if (dfs(edge.to)) {
              return true
            }
          } else if (recursionStack.has(edge.to)) {
            // 找到环
            return true
          }
        }
      }

      recursionStack.delete(taskId)
      return false
    }

    for (const taskId of this.adjacencyList.keys()) {
      if (!visited.has(taskId)) {
        if (dfs(taskId)) {
          return true
        }
      }
    }

    return false
  }

  /**
   * 获取循环依赖的路径
   * @returns 循环路径数组（如果有）
   */
  getCircularPath(): (string | number)[] | null {
    const visited = new Set<string | number>()
    const recursionStack = new Set<string | number>()
    const parentMap = new Map<string | number, string | number>()

    const dfs = (taskId: string | number): (string | number)[] | null => {
      visited.add(taskId)
      recursionStack.add(taskId)

      const node = this.adjacencyList.get(taskId)
      if (node) {
        for (const edge of node.outgoing) {
          if (!visited.has(edge.to)) {
            parentMap.set(edge.to, taskId)
            const result = dfs(edge.to)
            if (result) {
              return result
            }
          } else if (recursionStack.has(edge.to)) {
            // 找到环，构建路径
            const path: (string | number)[] = [edge.to]
            let current = taskId
            while (current !== edge.to) {
              path.push(current)
              current = parentMap.get(current) as string | number
            }
            path.push(edge.to)
            return path.reverse()
          }
        }
      }

      recursionStack.delete(taskId)
      return null
    }

    for (const taskId of this.adjacencyList.keys()) {
      if (!visited.has(taskId)) {
        const result = dfs(taskId)
        if (result) {
          return result
        }
      }
    }

    return null
  }

  /**
   * 拓扑排序
   * @returns 拓扑排序后的任务ID数组
   */
  topologicalSort(): (string | number)[] {
    const inDegree = new Map<string | number, number>()
    const result: (string | number)[] = []
    const queue: (string | number)[] = []

    // 计算入度
    for (const [taskId, node] of this.adjacencyList) {
      inDegree.set(taskId, node.incoming.length)
      if (node.incoming.length === 0) {
        queue.push(taskId)
      }
    }

    // BFS 拓扑排序
    while (queue.length > 0) {
      const current = queue.shift()!
      result.push(current)

      const node = this.adjacencyList.get(current)
      if (node) {
        for (const edge of node.outgoing) {
          const newInDegree = inDegree.get(edge.to)! - 1
          inDegree.set(edge.to, newInDegree)
          if (newInDegree === 0) {
            queue.push(edge.to)
          }
        }
      }
    }

    return result
  }

  /**
   * 计算关键路径
   * 使用动态规划算法（CPM - 关键路径法）
   * @returns 关键路径计算结果
   */
  calculateCriticalPath(): CriticalPathResult {
    const sortedTasks = this.topologicalSort()
    const nodes = new Map<string | number, CriticalPathNode>()

    // 正向计算最早开始和最早完成时间
    for (const taskId of sortedTasks) {
      const task = this.taskMap.get(taskId)
      if (!task) continue

      const node = this.adjacencyList.get(taskId)
      const duration = this.calculateTaskDuration(task)

      // 最早开始时间 = max(所有前驱的最早完成时间 + 延滞)
      let earliestStart = 0
      if (node && node.incoming.length > 0) {
        for (const edge of node.incoming) {
          const predecessorNode = nodes.get(edge.from)
          if (predecessorNode) {
            const finishWithLag = predecessorNode.earliestFinish + edge.lag
            earliestStart = Math.max(earliestStart, finishWithLag)
          }
        }
      }

      const earliestFinish = earliestStart + duration

      nodes.set(taskId, {
        taskId,
        earliestStart,
        earliestFinish,
        latestStart: 0,
        latestFinish: 0,
        slack: 0,
        isCritical: false
      })
    }

    // 反向计算最晚开始和最晚完成时间
    const maxFinishTime = Math.max(
      ...Array.from(nodes.values()).map(n => n.earliestFinish)
    )

    for (let i = sortedTasks.length - 1; i >= 0; i--) {
      const taskId = sortedTasks[i]
      const node = nodes.get(taskId)
      const adjacencyNode = this.adjacencyList.get(taskId)

      if (!node) continue

      let latestFinish: number
      if (adjacencyNode && adjacencyNode.outgoing.length > 0) {
        // 最晚完成时间 = min(所有后继的最晚开始时间 - 延滞)
        latestFinish = Infinity
        for (const edge of adjacencyNode.outgoing) {
          const successorNode = nodes.get(edge.to)
          if (successorNode) {
            const startWithLag = successorNode.latestStart - edge.lag
            latestFinish = Math.min(latestFinish, startWithLag)
          }
        }
      } else {
        // 没有后继，使用项目总工期
        latestFinish = maxFinishTime
      }

      const task = this.taskMap.get(taskId)
      const duration = task ? this.calculateTaskDuration(task) : 0
      const latestStart = latestFinish - duration
      const slack = latestStart - node.earliestStart

      nodes.set(taskId, {
        ...node,
        latestFinish,
        latestStart,
        slack,
        isCritical: slack === 0
      })
    }

    // 回溯关键路径
    const criticalPath: (string | number)[] = []
    let currentTask: string | number | null = null

    // 找到最早开始时间为0且时差为0的任务作为起点
    for (const taskId of sortedTasks) {
      const node = nodes.get(taskId)
      if (node && node.earliestStart === 0 && node.isCritical) {
        const hasCriticalPredecessor = this.adjacencyList
          .get(taskId)
          ?.incoming.some(edge => nodes.get(edge.from)?.isCritical)

        if (!hasCriticalPredecessor) {
          currentTask = taskId
          break
        }
      }
    }

    // 沿着关键路径遍历
    while (currentTask) {
      criticalPath.push(currentTask)
      const node = this.adjacencyList.get(currentTask)

      // 找到下一个关键路径任务
      let nextTask: string | number | null = null
      if (node) {
        for (const edge of node.outgoing) {
          const successorNode = nodes.get(edge.to)
          if (successorNode && successorNode.isCritical) {
            nextTask = edge.to
            break
          }
        }
      }

      currentTask = nextTask
    }

    return {
      path: criticalPath,
      duration: maxFinishTime,
      nodes
    }
  }

  /**
   * 计算任务工期（天数）
   */
  private calculateTaskDuration(task: GanttTask): number {
    const start = new Date(task.start).getTime()
    const end = new Date(task.end).getTime()
    return Math.ceil((end - start) / (1000 * 60 * 60 * 24))
  }

  /**
   * 获取任务的所有依赖关系（前驱和后继）
   * @param taskId 任务ID
   * @returns 依赖关系数组
   */
  getTaskDependencies(taskId: string | number): GanttTaskDependency[] {
    const task = this.taskMap.get(taskId)
    return task?.predecessors || []
  }

  /**
   * 清空图
   */
  clear(): void {
    this.adjacencyList.clear()
    this.taskMap.clear()
  }

  /**
   * 获取图的统计信息
   */
  getStats() {
    let edgeCount = 0
    for (const node of this.adjacencyList.values()) {
      edgeCount += node.outgoing.length
    }

    return {
      nodeCount: this.adjacencyList.size,
      edgeCount,
      avgDegree: this.adjacencyList.size > 0 ? edgeCount / this.adjacencyList.size : 0
    }
  }
}

/**
 * 创建依赖关系图实例的工厂函数
 */
export function createDependencyGraph(tasks: GanttTask[]): DependencyGraph {
  const graph = new DependencyGraph()
  graph.buildGraph(tasks)
  return graph
}

/**
 * 计算任务的所有下游任务（传递闭包）
 * @param tasks 任务列表
 * @param taskId 起始任务ID
 * @returns 下游任务ID集合
 */
export function getDownstreamTasks(
  tasks: GanttTask[],
  taskId: string | number
): Set<string | number> {
  const graph = new DependencyGraph()
  graph.buildGraph(tasks)

  const visited = new Set<string | number>()
  const queue: (string | number)[] = [taskId]

  while (queue.length > 0) {
    const current = queue.shift()!
    if (visited.has(current)) continue

    visited.add(current)
    const successors = graph.getSuccessors(current)
    queue.push(...successors)
  }

  visited.delete(taskId) // 移除自身
  return visited
}

/**
 * 计算任务的所有上游任务（传递闭包）
 * @param tasks 任务列表
 * @param taskId 起始任务ID
 * @returns 上游任务ID集合
 */
export function getUpstreamTasks(
  tasks: GanttTask[],
  taskId: string | number
): Set<string | number> {
  const graph = new DependencyGraph()
  graph.buildGraph(tasks)

  const visited = new Set<string | number>()
  const queue: (string | number)[] = [taskId]

  while (queue.length > 0) {
    const current = queue.shift()!
    if (visited.has(current)) continue

    visited.add(current)
    const predecessors = graph.getPredecessors(current)
    queue.push(...predecessors)
  }

  visited.delete(taskId) // 移除自身
  return visited
}

/**
 * 验证依赖关系是否有效
 * @param tasks 任务列表
 * @returns 验证结果
 */
export function validateDependencies(tasks: GanttTask[]): {
  valid: boolean
  circularPath?: (string | number)[]
  error?: string
} {
  const graph = new DependencyGraph()
  graph.buildGraph(tasks)

  if (graph.detectCircularDependency()) {
    const path = graph.getCircularPath()
    return {
      valid: false,
      circularPath: path || undefined,
      error: '检测到循环依赖'
    }
  }

  return { valid: true }
}
