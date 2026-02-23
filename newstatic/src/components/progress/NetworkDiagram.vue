<template>
  <div class="network-diagram" :class="{ 'fullscreen': isFullscreen }" ref="containerRef" @keydown="handleKeydown" @keyup="handleKeyup" tabindex="0">
    <!-- 工具栏 -->
    <div class="network-toolbar">
      <!-- 工具模式切换 -->
      <el-button-group size="small" title="工具模式" style="margin-right: 12px">
        <!-- 选择工具 -->
        <el-button
          @click="setToolMode('select')"
          :type="toolMode === 'select' ? 'primary' : 'default'"
          title="选择模式"
        >
          <el-icon>
            <svg viewBox="0 0 1024 1024" width="14" height="14">
              <path d="M896 448H560V112c0-17.7-14.3-32-32-32s-32 14.3-32 32v336H160c-17.7 0-32 14.3-32 32s14.3 32 32 32h336v336c0 17.7 14.3 32 32 32s32-14.3 32-32V512h336c17.7 0 32-14.3 32-32s-14.3-32-32-32z" fill="currentColor"></path>
            </svg>
          </el-icon>
        </el-button>
        <!-- 连线工具 -->
        <el-button
          @click="setToolMode('connect')"
          :type="toolMode === 'connect' ? 'primary' : 'default'"
          title="连线工具（创建任务）"
        >
          <el-icon>
            <svg viewBox="0 0 1024 1024" width="14" height="14">
              <path d="M720 304v-64c0-17.7-14.3-32-32-32s-32 14.3-32 32v64c0 17.7 14.3 32 32 32s32-14.3 32-32zm-96 96v-64c0-17.7-14.3-32-32-32s-32 14.3-32 32v64c0 17.7 14.3 32 32 32s32-14.3 32-32zm-96 96v-64c0-17.7-14.3-32-32-32s-32 14.3-32 32v64c0 17.7 14.3 32 32 32s32-14.3 32-32zm-96 96v-64c0-17.7-14.3-32-32-32s-32 14.3-32 32v64c0 17.7 14.3 32 32 32s32-14.3 32-32zm-96 96v-64c0-17.7-14.3-32-32-32s-32 14.3-32 32v64c0 17.7 14.3 32 32 32s32-14.3 32-32zm384-416v-64c0-17.7-14.3-32-32-32s-32 14.3-32 32v64c0 17.7 14.3 32 32 32s32-14.3 32-32zM224 752v64c0 17.7 14.3 32 32 32s32-14.3 32-32v-64c0-17.7-14.3-32-32-32s-32 14.3-32 32zm96-96v64c0 17.7 14.3 32 32 32s32-14.3 32-32v-64c0-17.7-14.3-32-32-32s-32 14.3-32 32zm96-96v64c0 17.7 14.3 32 32 32s32-14.3 32-32v-64c0-17.7-14.3-32-32-32s-32 14.3-32 32zm96-96v64c0 17.7 14.3 32 32 32s32-14.3 32-32v-64c0-17.7-14.3-32-32-32s-32 14.3-32 32zm-288 96h384c17.7 0 32-14.3 32-32s-14.3-32-32-32H320c-17.7 0-32 14.3-32 32s14.3 32 32 32z" fill="currentColor"></path>
            </svg>
          </el-icon>
        </el-button>
        <!-- 手形平移工具 -->
        <el-button
          @click="setToolMode('pan')"
          :type="toolMode === 'pan' ? 'primary' : 'default'"
          title="平移工具"
        >
          <el-icon>
            <svg viewBox="0 0 1024 1024" width="14" height="14">
              <path d="M768 256h-64V128c0-17.7-14.3-32-32-32s-32 14.3-32 32v128H384V128c0-17.7-14.3-32-32-32s-32 14.3-32 32v128h-64c-17.7 0-32 14.3-32 32s14.3 32 32 32h64v416c0 17.7 14.3 32 32 32h256v128c0 17.7 14.3 32 32 32s32-14.3 32-32V768h64c17.7 0 32-14.3 32-32s-14.3-32-32-32h-64V320h64c17.7 0 32-14.3 32-32s-14.3-32-32-32z" fill="currentColor"></path>
            </svg>
          </el-icon>
        </el-button>
      </el-button-group>

      <!-- 缩放控制 -->
      <el-button-group size="small">
        <el-button @click="zoomOut" title="缩小">
          <el-icon><ZoomOut /></el-icon>
        </el-button>
        <el-button @click="resetZoom" title="重置">
          {{ Math.round(zoomLevel * 100) }}%
        </el-button>
        <el-button @click="zoomIn" title="放大">
          <el-icon><ZoomIn /></el-icon>
        </el-button>
      </el-button-group>

      <!-- 视图选项 -->
      <div style="margin-left: auto; display: flex; align-items: center; gap: 12px;">
        <el-button type="primary" size="small" @click="handleAddTask">
          <el-icon style="margin-right: 4px;"><Plus /></el-icon>
          添加任务
        </el-button>
        <el-checkbox v-model="showCriticalPath" @change="render" size="small">
          关键路径
        </el-checkbox>
        <el-checkbox v-model="focusCriticalPath" @change="render" size="small">
          聚焦关键路径
        </el-checkbox>
        <el-checkbox v-model="showTimeParams" @change="render" size="small">
          时间参数
        </el-checkbox>
        <el-checkbox v-model="showTaskNames" @change="render" size="small">
          任务名称
        </el-checkbox>
        <el-checkbox v-model="showSlack" @change="render" size="small">
          时差信息
        </el-checkbox>

        <!-- 布局方式 -->
        <el-select
          v-model="layoutMode"
          size="small"
          style="width: 120px"
          @change="handleLayoutChange"
        >
          <el-option label="自动布局" value="auto" />
          <el-option label="从左到右" value="left-right" />
          <el-option label="从上到下" value="top-down" />
        </el-select>

        <!-- 更多操作 -->
        <el-dropdown trigger="click" @command="handleMoreAction">
          <el-button size="small">
            更多
            <el-icon style="margin-left: 4px;"><ArrowDown /></el-icon>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="fit-view">
                <el-icon><FullScreen /></el-icon>
                适应视图
              </el-dropdown-item>
              <el-dropdown-item command="center-view">
                <el-icon><Aim /></el-icon>
                居中视图
              </el-dropdown-item>
              <el-dropdown-item command="export" divided>
                <el-icon><Download /></el-icon>
                导出图片
              </el-dropdown-item>
              <el-dropdown-item command="toggle-fullscreen">
                <el-icon v-if="!isFullscreen"><Crop /></el-icon>
                <el-icon v-else><Close /></el-icon>
                {{ isFullscreen ? '退出全屏' : '全屏' }}
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>

    <!-- 统计信息 -->
    <div class="network-stats" v-if="stats">
      <div class="stat-item">
        <span class="stat-label">事件节点</span>
        <span class="stat-value">{{ stats.nodes }}</span>
      </div>
      <div class="stat-item">
        <span class="stat-label">活动(任务)</span>
        <span class="stat-value">{{ stats.activities }}</span>
      </div>
      <div class="stat-item">
        <span class="stat-label">共享节点</span>
        <span class="stat-value">{{ stats.sharedNodes || 0 }}</span>
      </div>
      <div class="stat-item">
        <span class="stat-label">虚活动</span>
        <span class="stat-value">{{ stats.dummyActivities || 0 }}</span>
      </div>
      <div class="stat-item">
        <span class="stat-label">关键路径</span>
        <span class="stat-value critical">{{ stats.criticalActivities }}/{{ stats.activities }}</span>
      </div>
      <div class="stat-item">
        <span class="stat-label">关键节点</span>
        <span class="stat-value critical">{{ stats.criticalNodes || 0 }}</span>
      </div>
      <div class="stat-item">
        <span class="stat-label">总工期</span>
        <span class="stat-value">{{ Math.round(stats.totalDuration * 10) / 10 }}天</span>
      </div>
    </div>

    <!-- 时间轴（参考甘特图实现，改进显示方式） -->
    <div class="network-timeline-header" v-if="nodes.length > 0">
      <div class="timeline-content" ref="timelineHeaderRef" :style="{ transform: `translateX(${panX}px)`, width: timelineWidth + 'px' }">
        <!-- 月份行（上层） -->
        <div class="timeline-months-row">
          <div
            v-for="month in timelineMonths"
            :key="month.key"
            class="timeline-month-cell"
            :style="{ left: month.position + 'px', width: month.width + 'px' }"
          >
            <div class="month-label">{{ month.label }}</div>
          </div>
        </div>
        <!-- 日期行（下层） -->
        <div class="timeline-days-row">
          <div
            v-for="day in timelineDays"
            :key="day.date"
            class="timeline-day-cell"
            :class="{ 'is-today': day.isToday, 'is-weekend': day.isWeekend, 'is-odd': day.isOdd }"
            :style="{ left: day.position + 'px', width: dayWidth + 'px' }"
          >
            <div class="day-number">{{ day.day }}</div>
          </div>
        </div>
        <!-- 今天标记线 -->
        <div
          v-if="todayPosition !== null"
          class="today-line"
          :style="{ left: todayPosition + 'px' }"
        ></div>
      </div>
    </div>

    <!-- SVG 绘图区域 -->
    <div class="network-canvas-container" ref="canvasContainerRef">
      <svg
        ref="svgRef"
        class="network-svg"
        :width="svgWidth"
        :height="svgHeight"
        @wheel.prevent="handleWheel"
        @mousedown="handleMouseDown"
        @mousemove="handleMouseMove"
        @mouseup="handleMouseUp"
        @mouseleave="handleMouseUp"
      >
        <defs>
          <!-- 箭头标记 - refX设置为箭头长度(9)，使箭头尖端刚好在路径终点 -->
          <marker
            id="arrowhead-normal"
            markerWidth="10"
            markerHeight="7"
            refX="9"
            refY="3.5"
            orient="auto"
            markerUnits="userSpaceOnUse"
          >
            <path d="M0,0 L0,7 L9,3.5 z" fill="#64B5F6" />
          </marker>
          <marker
            id="arrowhead-critical"
            markerWidth="10"
            markerHeight="7"
            refX="9"
            refY="3.5"
            orient="auto"
            markerUnits="userSpaceOnUse"
          >
            <path d="M0,0 L0,7 L9,3.5 z" fill="#FF8A65" />
          </marker>
          <!-- 节点阴影 -->
          <filter id="nodeShadow" x="-50%" y="-50%" width="200%" height="200%">
            <feDropShadow dx="0" dy="2" stdDeviation="3" flood-opacity="0.2"/>
          </filter>
        </defs>

        <!-- 网格背景 -->
        <g class="grid-background" v-if="showGrid">
          <pattern
            id="grid"
            :width="gridSize * zoomLevel"
            :height="gridSize * zoomLevel"
            patternUnits="userSpaceOnUse"
          >
            <path
              :d="`M ${gridSize * zoomLevel} 0 L 0 0 0 ${gridSize * zoomLevel}`"
              fill="none"
              stroke="#f0f0f0"
              stroke-width="1"
            />
          </pattern>
          <rect
            :width="svgWidth"
            :height="svgHeight"
            fill="url(#grid)"
          />
        </g>

        <!-- 变换组 (只平移，不缩放) - 缩放通过动态计算节点位置实现 -->
        <g :transform="`translate(${panX}, ${panY})`">
          <!-- 连线过程中的临时连线 -->
          <line
            v-if="isConnecting && connectStartNode"
            :x1="connectStartNode.x"
            :y1="connectStartNode.y"
            :x2="connectEndPosition.x"
            :y2="connectEndPosition.y"
            stroke="#64B5F6"
            stroke-width="2"
            stroke-dasharray="5 5"
            opacity="0.6"
          />

          <!-- 连接线 -->
          <g class="edges">
            <path
              v-for="edge in edges"
              :key="edge.id"
              :d="edge.path"
              :stroke="edge.isDummy ? '#B0BEC5' : (edge.isCritical ? '#FF8A65' : (edge.isTaskEdge ? '#64B5F6' : '#CFD8DC'))"
              :stroke-width="edge.isDummy ? 1 : (edge.isTaskEdge ? (edge.isCritical ? 1.5 : 1.2) : 0.8)"
              :stroke-dasharray="edge.isDummy ? '3 3' : (edge.isTaskEdge ? 'none' : '4 2')"
              :opacity="focusCriticalPath && !edge.isCritical && !edge.isDummy ? 0.2 : 1"
              fill="none"
              :marker-end="edge.isTaskEdge ? (edge.isCritical ? 'url(#arrowhead-critical)' : 'url(#arrowhead-normal)') : 'none'"
              :class="{ 'edge-critical': edge.isCritical, 'edge-selected': selectedEdge?.id === edge.id, 'edge-task': edge.isTaskEdge, 'edge-dependency': !edge.isTaskEdge, 'edge-dummy': edge.isDummy, 'edge-dimmed': focusCriticalPath && !edge.isCritical && !edge.isDummy }"
              @click="handleEdgeClick(edge)"
              @dblclick="handleEdgeDblClick(edge)"
              @contextmenu.prevent="handleContextMenu($event, 'edge', edge)"
              style="cursor: pointer"
            />
            <!-- 活动标签（只在任务边上显示） -->
            <text
              v-for="edge in edges"
              v-show="showTaskNames && edge.isTaskEdge"
              :key="'label-' + edge.id"
              :x="edge.labelX"
              :y="edge.labelY"
              text-anchor="middle"
              class="edge-label"
              :class="{ 'edge-label-task': edge.isTaskEdge }"
              font-size="12"
              fill="#303133"
              font-weight="500"
            >
              {{ edge.label }}
            </text>
            <!-- 时差标签（只在任务边上显示） -->
            <text
              v-for="edge in edges"
              v-show="showSlack && edge.isTaskEdge && (edge.totalSlack !== undefined)"
              :key="'slack-' + edge.id"
              :x="edge.labelX"
              :y="edge.labelY + 14"
              text-anchor="middle"
              class="edge-slack-label"
              font-size="10"
              :fill="(edge.totalSlack || 0) < 0.1 ? '#FF8A65' : '#90A4AE'"
            >
              TS: {{ Math.round((edge.totalSlack || 0) * 10) / 10 }}
            </text>
          </g>

          <!-- 节点 -->
          <g class="nodes">
            <g
              v-for="node in nodes"
              :key="node.id"
              :transform="`translate(${node.x}, ${node.y})`"
              :opacity="focusCriticalPath && !node.isCritical ? 0.3 : 1"
              class="node-group"
              :class="{
                'node-critical': node.isCritical,
                'node-selected': selectedNode?.id === node.id,
                'node-connecting': connectMode && connectStartNode?.id === node.id,
                'node-connect-hover': isConnecting && connectHoverNode?.id === node.id,
                'node-drag-hover': isDragging && dragHoverNode?.id === node.id,
                'node-start': node.nodeType === 'start',
                'node-end': node.nodeType === 'end',
                'node-shared': node.nodeType === 'shared',
                'node-dimmed': focusCriticalPath && !node.isCritical
              }"
              @click="handleNodeClick(node)"
              @dblclick="handleNodeDblClick(node)"
              @mousedown="handleNodeMouseDown($event, node)"
              @mouseenter="handleNodeMouseEnter($event, node)"
              @mouseleave="handleNodeMouseLeave(node)"
              @contextmenu.prevent="handleContextMenu($event, 'node', node)"
              :style="{ cursor: connectMode ? 'crosshair' : (panMode ? 'grab' : 'pointer') }"
            >
              <!-- 节点圆形 -->
              <circle
                :r="nodeRadius"
                :fill="node.nodeType === 'start' ? '#5AB9A6' : node.nodeType === 'end' ? '#E57373' : node.nodeType === 'shared' ? '#64B5F6' : '#90A4AE'"
                filter="url(#nodeShadow)"
                stroke="white"
                stroke-width="2"
              />
              <!-- 节点编号 -->
              <text
                y="5"
                text-anchor="middle"
                fill="white"
                font-weight="bold"
                font-size="14"
              >{{ node.number }}</text>
              <!-- 节点类型标记 (S=起点, E=终点, Sh=共享) -->
              <text
                :y="node.nodeType === 'start' ? -32 : (node.nodeType === 'shared' ? -32 : 42)"
                text-anchor="middle"
                fill="#606266"
                font-weight="bold"
                font-size="10"
              >{{ node.nodeType === 'start' ? 'S' : (node.nodeType === 'shared' ? 'Sh' : 'E') }}</text>
              <!-- 时间参数 -->
              <g v-if="showTimeParams" class="time-params">
                <text
                  y="-40"
                  text-anchor="middle"
                  class="time-param-text"
                  font-size="11"
                  fill="#5AB9A6"
                >ES: {{ node.ES }}</text>
                <text
                  y="50"
                  text-anchor="middle"
                  class="time-param-text"
                  font-size="11"
                  fill="#64B5F6"
                >LS: {{ node.LS }}</text>
              </g>
              <!-- 时差信息 -->
              <g v-if="showSlack && (node.totalSlack !== undefined)" class="slack-params">
                <text
                  y="-52"
                  text-anchor="middle"
                  class="slack-param-text"
                  font-size="10"
                  :fill="(node.totalSlack || 0) < 0.1 ? '#FF8A65' : '#90A4AE'"
                >TS: {{ Math.round((node.totalSlack || 0) * 10) / 10 }}</text>
                <text
                  v-if="node.freeSlack !== undefined"
                  y="62"
                  text-anchor="middle"
                  class="slack-param-text"
                  font-size="10"
                  :fill="(node.freeSlack || 0) < 0.1 ? '#64B5F6' : '#90A4AE'"
                >FS: {{ Math.round((node.freeSlack || 0) * 10) / 10 }}</text>
              </g>
            </g>
          </g>
        </g>
      </svg>

      <!-- 加载状态 -->
      <div v-if="loading" class="network-loading">
        <el-icon class="is-loading" :size="40" />
        <p>加载中...</p>
      </div>

      <!-- 空状态 -->
      <div v-if="!loading && nodes.length === 0" class="network-empty">
        <el-empty description="暂无网络图数据" />
      </div>
    </div>

    <!-- 图例 -->
    <div class="network-legend">
      <div class="legend-item">
        <span class="legend-color start"></span>
        <span>任务起点(S)</span>
      </div>
      <div class="legend-item">
        <span class="legend-color shared"></span>
        <span>共享节点(Sh)</span>
      </div>
      <div class="legend-item">
        <span class="legend-color end"></span>
        <span>任务终点(E)</span>
      </div>
      <div class="legend-item">
        <span class="legend-edge task-edge"></span>
        <span>任务连线</span>
      </div>
      <div class="legend-item">
        <span class="legend-edge dummy-edge"></span>
        <span>虚活动</span>
      </div>
      <div class="legend-item">
        <span class="legend-color critical"></span>
        <span>关键路径</span>
      </div>
    </div>

    <!-- 右键菜单 -->
    <div
      v-if="contextMenuVisible"
      class="context-menu"
      :style="{ left: contextMenuPosition.x + 'px', top: contextMenuPosition.y + 'px' }"
      @click.stop
    >
      <div class="context-menu-header" v-if="contextMenuTarget">
        <template v-if="contextMenuTarget.type === 'node'">
          <span class="menu-title">节点 {{ contextMenuTarget.data.number }}</span>
          <span class="menu-subtitle">{{ contextMenuTarget.data.name || '未命名' }}</span>
        </template>
        <template v-else-if="contextMenuTarget.type === 'edge'">
          <span class="menu-title">任务连线</span>
          <span class="menu-subtitle">{{ contextMenuTarget.data.label || '未命名' }}</span>
        </template>
      </div>
      <el-divider style="margin: 8px 0" />
      <div class="context-menu-items">
        <!-- 节点菜单项 -->
        <template v-if="contextMenuTarget?.type === 'node'">
          <div class="menu-item" @click="handleContextMenuAction('edit-node')">
            <el-icon><Edit /></el-icon>
            <span>编辑节点</span>
          </div>
          <div class="menu-item" @click="handleContextMenuAction('view-details')">
            <el-icon><View /></el-icon>
            <span>查看详情</span>
          </div>
          <el-divider style="margin: 8px 0" />
          <div class="menu-item danger" @click="handleContextMenuAction('delete-node')">
            <el-icon><Delete /></el-icon>
            <span>删除节点</span>
          </div>
        </template>
        <!-- 边菜单项 -->
        <template v-else-if="contextMenuTarget?.type === 'edge'">
          <div class="menu-item" @click="handleContextMenuAction('edit-edge')">
            <el-icon><Edit /></el-icon>
            <span>编辑连线</span>
          </div>
          <div class="menu-item" @click="handleContextMenuAction('view-details')">
            <el-icon><View /></el-icon>
            <span>查看详情</span>
          </div>
          <el-divider style="margin: 8px 0" />
          <div class="menu-item danger" @click="handleContextMenuAction('delete-edge')">
            <el-icon><Delete /></el-icon>
            <span>删除连线</span>
          </div>
        </template>
      </div>
    </div>

    <!-- 节点详情对话框 -->
    <el-drawer
      v-model="nodeDetailVisible"
      title="节点详情"
      direction="rtl"
      size="500px"
      :close-on-click-modal="false"
    >
      <div v-if="selectedNode" class="node-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="节点编号" :span="2">
            <el-tag type="primary" size="large">节点 {{ selectedNode.number }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="关联任务" :span="2" v-if="selectedNode.taskId">
            {{ selectedNode.name }}
          </el-descriptions-item>
          <el-descriptions-item label="节点类型" :span="2">
            <el-tag v-if="selectedNode.nodeType === 'start'" type="success" size="small">
              任务起点 (S)
            </el-tag>
            <el-tag v-else-if="selectedNode.nodeType === 'end'" type="danger" size="small">
              任务终点 (E)
            </el-tag>
            <el-tag v-else-if="selectedNode.nodeType === 'shared'" type="warning" size="small">
              共享节点 (Sh)
            </el-tag>
            <el-tag v-else type="primary" size="small">普通节点</el-tag>
          </el-descriptions-item>
        </el-descriptions>

        <el-divider content-position="left">PERT 时间参数</el-divider>

        <el-descriptions :column="2" border>
          <el-descriptions-item label="最早时间(ET)">
            {{ Math.round(selectedNode.ES * 10) / 10 }} 天
          </el-descriptions-item>
          <el-descriptions-item label="最迟时间(LT)">
            {{ Math.round(selectedNode.LS * 10) / 10 }} 天
          </el-descriptions-item>
          <el-descriptions-item label="最早开始(ES)">
            {{ Math.round(selectedNode.ES * 10) / 10 }} 天
          </el-descriptions-item>
          <el-descriptions-item label="最早结束(EF)">
            {{ Math.round(selectedNode.EF * 10) / 10 }} 天
          </el-descriptions-item>
          <el-descriptions-item label="最迟开始(LS)">
            {{ Math.round(selectedNode.LS * 10) / 10 }} 天
          </el-descriptions-item>
          <el-descriptions-item label="最迟结束(LF)">
            {{ Math.round(selectedNode.LF * 10) / 10 }} 天
          </el-descriptions-item>
          <el-descriptions-item label="总时差" :span="2">
            <el-tag :type="(selectedNode.totalSlack || 0) < 0.1 ? 'danger' : 'success'" size="small">
              {{ Math.round((selectedNode.totalSlack || 0) * 10) / 10 }} 天
            </el-tag>
            <span style="margin-left: 8px; color: #909399; font-size: 12px">
              LT - ET
            </span>
          </el-descriptions-item>
          <el-descriptions-item label="自由时差" :span="2" v-if="selectedNode.freeSlack !== undefined">
            <el-tag :type="(selectedNode.freeSlack || 0) < 0.1 ? 'warning' : 'info'" size="small">
              {{ Math.round((selectedNode.freeSlack || 0) * 10) / 10 }} 天
            </el-tag>
            <span style="margin-left: 8px; color: #909399; font-size: 12px">
              不影响后续任务的最早开始
            </span>
          </el-descriptions-item>
          <el-descriptions-item label="关键节点" :span="2">
            <el-tag :type="selectedNode.isCritical ? 'danger' : 'success'" size="small">
              {{ selectedNode.isCritical ? '是 (总时差=0)' : '否' }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>

        <el-alert
          v-if="selectedNode.nodeType === 'shared'"
          type="info"
          :closable="false"
          show-icon
          style="margin-top: 16px"
        >
          此节点既是前置任务的终点，也是后置任务的起点
        </el-alert>

        <el-divider />

        <!-- 前置节点 -->
        <div v-if="selectedNode.predecessors && selectedNode.predecessors.length > 0">
          <h4>前置节点 ({{ selectedNode.predecessors.length }})</h4>
          <el-space wrap>
            <el-tag
              v-for="pred in selectedNode.predecessors"
              :key="pred.id"
              :type="pred.nodeType === 'start' ? 'success' : (pred.nodeType === 'shared' ? 'warning' : 'danger')"
              size="small"
              style="margin: 4px"
            >
              节点 {{ pred.number }} {{ pred.nodeType === 'start' ? '(起点S)' : (pred.nodeType === 'shared' ? '(共享Sh)' : '(终点E)') }}
            </el-tag>
          </el-space>
        </div>
        <div v-else style="color: #909399; font-size: 13px;">
          无前置节点
        </div>

        <!-- 后置节点 -->
        <div v-if="selectedNode.successors && selectedNode.successors.length > 0" style="margin-top: 16px">
          <h4>后置节点 ({{ selectedNode.successors.length }})</h4>
          <el-space wrap>
            <el-tag
              v-for="succ in selectedNode.successors"
              :key="succ.id"
              :type="succ.nodeType === 'start' ? 'success' : (succ.nodeType === 'shared' ? 'warning' : 'danger')"
              size="small"
              style="margin: 4px"
            >
              节点 {{ succ.number }} {{ succ.nodeType === 'start' ? '(起点S)' : (succ.nodeType === 'shared' ? '(共享Sh)' : '(终点E)') }}
            </el-tag>
          </el-space>
        </div>
        <div v-else style="color: #909399; font-size: 13px; margin-top: 16px;">
          无后置节点
        </div>
      </div>
    </el-drawer>

    <!-- 活动详情对话框 -->
    <el-drawer
      v-model="edgeDetailVisible"
      title="活动详情"
      direction="rtl"
      size="500px"
      :close-on-click-modal="false"
    >
      <div v-if="selectedEdge" class="edge-detail">
        <el-alert
          v-if="!selectedEdge.isTaskEdge"
          type="info"
          :closable="false"
          show-icon
          style="margin-bottom: 16px"
        >
          这是前置依赖关系，不是实际任务
        </el-alert>
        <el-descriptions :column="2" border v-if="selectedEdge.isTaskEdge">
          <el-descriptions-item label="活动名称" :span="2">
            {{ selectedEdge.label }}
          </el-descriptions-item>
          <el-descriptions-item label="起点">
            节点 {{ selectedEdge.from }} {{ getNodeTypeLabel(selectedEdge.from) }}
          </el-descriptions-item>
          <el-descriptions-item label="终点">
            节点 {{ selectedEdge.to }} {{ getNodeTypeLabel(selectedEdge.to) }}
          </el-descriptions-item>
          <el-descriptions-item label="工期">
            {{ selectedEdge.duration }} 天
          </el-descriptions-item>
          <el-descriptions-item label="关键活动">
            <el-tag :type="selectedEdge.isCritical ? 'danger' : 'primary'" size="small">
              {{ selectedEdge.isCritical ? '是' : '否' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="总时差" :span="2" v-if="selectedEdge.totalSlack !== undefined">
            <el-tag :type="(selectedEdge.totalSlack || 0) < 0.1 ? 'danger' : 'success'" size="small">
              {{ Math.round((selectedEdge.totalSlack || 0) * 10) / 10 }} 天
            </el-tag>
            <span style="margin-left: 8px; color: #909399; font-size: 12px">
              可以延迟而不影响项目工期
            </span>
          </el-descriptions-item>
          <el-descriptions-item label="自由时差" :span="2" v-if="selectedEdge.freeSlack !== undefined">
            <el-tag :type="(selectedEdge.freeSlack || 0) < 0.1 ? 'warning' : 'info'" size="small">
              {{ Math.round((selectedEdge.freeSlack || 0) * 10) / 10 }} 天
            </el-tag>
            <span style="margin-left: 8px; color: #909399; font-size: 12px">
              可用不影响后续活动
            </span>
          </el-descriptions-item>
        </el-descriptions>
        <el-descriptions :column="1" border v-else>
          <el-descriptions-item label="类型">
            <el-tag type="info" size="small">前置依赖关系</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="起点">
            节点 {{ selectedEdge.from }} {{ getNodeTypeLabel(selectedEdge.from) }}
          </el-descriptions-item>
          <el-descriptions-item label="终点">
            节点 {{ selectedEdge.to }} {{ getNodeTypeLabel(selectedEdge.to) }}
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </el-drawer>

    <!-- 任务编辑对话框 -->
    <el-dialog
      v-model="editDialogVisible"
      :title="editingTask ? '编辑任务' : '新建任务'"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form :model="taskForm" :rules="taskFormRules" ref="taskFormRef" label-width="100px">
        <el-form-item label="任务名称" prop="name">
          <el-input v-model="taskForm.name" placeholder="请输入任务名称" />
        </el-form-item>
        <el-form-item label="开始日期" prop="start">
          <el-date-picker
            v-model="taskForm.start"
            type="date"
            placeholder="选择开始日期"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="结束日期" prop="end">
          <el-date-picker
            v-model="taskForm.end"
            type="date"
            placeholder="选择结束日期"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="进度" prop="progress">
          <el-slider v-model="taskForm.progress" :marks="{ 0: '0%', 50: '50%', 100: '100%' }" />
        </el-form-item>
        <el-form-item label="优先级" prop="priority">
          <el-radio-group v-model="taskForm.priority">
            <el-radio-button label="urgent">紧急</el-radio-button>
            <el-radio-button label="high">高</el-radio-button>
            <el-radio-button label="medium">中</el-radio-button>
            <el-radio-button label="low">低</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="taskForm.notes" type="textarea" :rows="3" placeholder="任务备注" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveTask" :loading="saving">保存</el-button>
      </template>
    </el-dialog>

    <!-- 新建任务对话框（从连线创建） -->
    <el-dialog
      v-model="createTaskDialogVisible"
      title="新建任务"
      width="600px"
      :close-on-click-modal="false"
    >
      <el-form :model="createTaskForm" :rules="createTaskFormRules" ref="createTaskFormRef" label-width="100px">
        <el-form-item label="任务名称" prop="name">
          <el-input v-model="createTaskForm.name" placeholder="请输入任务名称" />
        </el-form-item>
        <el-form-item label="任务类型">
          <el-radio-group v-model="createTaskForm.isDummy">
            <el-radio :label="false">实际任务</el-radio>
            <el-radio :label="true">虚任务</el-radio>
          </el-radio-group>
          <div style="color: #909399; font-size: 12px; margin-top: 4px;">
            实际任务：有工期的任务，使用实线表示
            <br>
            虚任务：仅表示依赖关系，无实际工期，使用虚线表示
          </div>
        </el-form-item>
        <el-form-item label="开始日期" prop="start_date">
          <el-date-picker
            v-model="createTaskForm.start_date"
            type="date"
            placeholder="选择开始日期"
            style="width: 100%"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item label="结束日期" prop="end_date">
          <el-date-picker
            v-model="createTaskForm.end_date"
            type="date"
            placeholder="选择结束日期"
            style="width: 100%"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item label="进度" prop="progress">
          <el-slider v-model="createTaskForm.progress" :marks="{ 0: '0%', 50: '50%', 100: '100%' }" />
        </el-form-item>
        <el-form-item label="优先级" prop="priority">
          <el-radio-group v-model="createTaskForm.priority">
            <el-radio-button label="urgent">紧急</el-radio-button>
            <el-radio-button label="high">高</el-radio-button>
            <el-radio-button label="medium">中</el-radio-button>
            <el-radio-button label="low">低</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="createTaskForm.description" type="textarea" :rows="3" placeholder="任务描述" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="cancelCreateTask">取消</el-button>
        <el-button type="primary" @click="submitCreateTask" :loading="saving">创建</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted, nextTick, toRaw } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  ZoomIn,
  ZoomOut,
  ArrowDown,
  FullScreen,
  Close,
  Download,
  Crop,
  Aim,
  Plus
} from '@element-plus/icons-vue'
import html2canvas from 'html2canvas'
import { progressApi } from '@/api'
import { findOrthogonalPath } from '@/utils/ganttHelpers'

const props = defineProps({
  projectId: {
    type: [Number, String],
    required: true
  },
  scheduleData: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['node-selected', 'position-updated', 'task-updated'])

// 视图状态
const zoomLevel = ref(1)
const panX = ref(0)
const panY = ref(0)
const isFullscreen = ref(false)
const loading = ref(false)
const showGrid = ref(true)
const gridSize = ref(20)
const nodeRadius = ref(18)  // 减小节点半径，从25改为18
const toolMode = ref('select')  // 工具模式: 'select', 'connect', 'pan'
const panMode = computed(() => toolMode.value === 'pan')  // 平移模式
const connectMode = computed(() => toolMode.value === 'connect')  // 连线模式

// 显示选项
const showCriticalPath = ref(true)
const focusCriticalPath = ref(false)
const showTimeParams = ref(false)  // 默认不显示时间参数
const showTaskNames = ref(true)
const showSlack = ref(false)  // 默认不显示时差信息
const layoutMode = ref('auto')

// 画布尺寸
const svgWidth = computed(() => Math.max(2000, timelineWidth.value))
const svgHeight = ref(1200)
const headerHeight = 40 // 顶部时间刻度的高度
const dayWidth = ref(40) // 每天的像素宽度（用于时间轴显示和节点位置计算）

// 拖拽状态
const isDragging = ref(false)
const isPanning = ref(false)
const dragStartPos = ref({ x: 0, y: 0 })
const panStartPos = ref({ x: 0, y: 0 })
const draggedNode = ref(null)
const dragHoverNode = ref(null)  // 拖动时悬停的目标节点（用于节点合并）

// 连线状态
const isConnecting = ref(false)
const connectStartNode = ref(null)
const connectEndPosition = ref({ x: 0, y: 0 })
const connectHoverNode = ref(null)  // 连线时悬停的目标节点
const isCtrlPressed = ref(false)  // Ctrl键状态

// 右键菜单
const contextMenuVisible = ref(false)
const contextMenuPosition = ref({ x: 0, y: 0 })
const contextMenuTarget = ref(null) // { type: 'node' | 'edge', data: node | edge }

// 选中项
const selectedNode = ref(null)
const selectedEdge = ref(null)
const nodeDetailVisible = ref(false)
const edgeDetailVisible = ref(false)

// 任务编辑
const editDialogVisible = ref(false)
const editingTask = ref(null)
const saving = ref(false)
const taskFormRef = ref(null)
const taskForm = ref({
  name: '',
  start: null,
  end: null,
  progress: 0,
  priority: 'medium',
  notes: ''
})

// 新建任务对话框（从连线创建）
const createTaskDialogVisible = ref(false)
const createTaskForm = ref({
  name: '',
  isDummy: false,  // 是否为虚任务
  start_date: null,
  end_date: null,
  progress: 0,
  priority: 'medium',
  description: ''
})
const createTaskFormRef = ref(null)

const taskFormRules = {
  name: [{ required: true, message: '请输入任务名称', trigger: 'blur' }],
  start: [{ required: true, message: '请选择开始日期', trigger: 'change' }],
  end: [{ required: true, message: '请选择结束日期', trigger: 'change' }]
}

// 新建任务表单验证规则
const createTaskFormRules = {
  name: [{ required: true, message: '请输入任务名称', trigger: 'blur' }],
  start_date: [{ required: true, message: '请选择开始日期', trigger: 'change' }],
  end_date: [
    { required: true, message: '请选择结束日期', trigger: 'change' },
    {
      validator: (rule, value, callback) => {
        if (!value || !createTaskForm.value.start_date) {
          callback()
          return
        }
        const startDate = new Date(createTaskForm.value.start_date)
        const endDate = new Date(value)
        if (endDate < startDate) {
          callback(new Error('结束日期不能早于开始日期'))
        } else {
          callback()
        }
      },
      trigger: 'change'
    }
  ]
}

// Refs
const containerRef = ref(null)
const canvasContainerRef = ref(null)
const svgRef = ref(null)
const timelineHeaderRef = ref(null)

// 节点和边数据
const nodes = ref([])
const edges = ref([])

// 关键路径序列 (计算完整的关键路径任务序列)
const criticalPathSequence = computed(() => {
  const criticalEdges = edges.value.filter(e => e.isCritical && e.isTaskEdge)

  // 按开始时间排序关键路径任务
  const sorted = criticalEdges.sort((a, b) => {
    const aNode = nodes.value.find(n => n.id === a.from)
    const bNode = nodes.value.find(n => n.id === b.from)
    return (aNode?.ES || 0) - (bNode?.ES || 0)
  })

  return sorted.map(e => ({
    id: e.id,
    label: e.label,
    from: e.from,
    to: e.to,
    duration: e.duration
  }))
})

// 统计信息
const stats = computed(() => {
  // 计算任务数量（任务边数量）
  const taskEdges = edges.value.filter(e => e.isTaskEdge)
  const taskCount = taskEdges.length

  // 计算节点数量（所有节点）
  const nodeCount = nodes.value.length

  // 计算关键路径上的任务数量
  const criticalTaskCount = taskEdges.filter(e => e.isCritical).length

  // 计算关键路径节点数量
  const criticalNodeCount = nodes.value.filter(n => n.isCritical).length

  // 计算共享节点数量
  const sharedNodeCount = nodes.value.filter(n => n.nodeType === 'shared').length

  // 计算虚活动数量
  const dummyActivityCount = edges.value.filter(e => e.isDummy).length

  // 计算总工期
  const totalDuration = calculateTotalDuration()

  // 计算项目总时差（关键路径上时差为0）
  const criticalPathLength = criticalTaskCount

  return {
    nodes: nodeCount,
    activities: taskCount,
    criticalActivities: criticalTaskCount,
    criticalNodes: criticalNodeCount,
    sharedNodes: sharedNodeCount,
    dummyActivities: dummyActivityCount,
    totalDuration: totalDuration,
    criticalPathLength: criticalPathLength
  }
})

// 时间轴数据（参考甘特图实现，根据项目起始时间初始化）
const timelineDays = computed(() => {
  if (nodes.value.length === 0 || !props.scheduleData?.activities) return []

  const activities = Object.values(props.scheduleData.activities)
  if (activities.length === 0) return []

  // 过滤掉无效的活动数据：earliest_start 必须大于 0
  const validActivities = activities.filter(a => {
    const isValid = a.earliest_start && a.earliest_start > 0 && a.earliest_finish && a.earliest_finish > 0
    if (!isValid) {
      console.log('NetworkDiagram - Filtering out invalid activity:', a.name || a.id, 'earliest_start:', a.earliest_start, 'earliest_finish:', a.earliest_finish)
    }
    return isValid
  })

  if (validActivities.length === 0) return []

  // 将activities的时间戳转换为Date对象，找到最小和最大日期
  const taskDates = validActivities.map(a => {
    const startDate = new Date(a.earliest_start * 1000)
    const endDate = new Date(a.earliest_finish * 1000)
    return { startDate, endDate }
  })

  const minDate = new Date(Math.min(...taskDates.map(d => d.startDate.getTime())))
  const maxDate = new Date(Math.max(...taskDates.map(d => d.endDate.getTime())))

  const days = []
  const today = new Date()
  today.setHours(0, 0, 0, 0)

  // 添加缓冲：从项目起始日期前2天开始，结束后5天
  const bufferDays = 2
  const endBufferDays = 5

  const startDate = new Date(minDate)
  startDate.setDate(startDate.getDate() - bufferDays)

  const endDate = new Date(maxDate)
  endDate.setDate(endDate.getDate() + endBufferDays)

  const currentDate = new Date(startDate)
  let dayIndex = 0

  while (currentDate <= endDate) {
    const dateStr = formatDate(currentDate, 'YYYY-MM-DD')
    const weekday = formatDate(currentDate, 'ddd')
    const dayNum = currentDate.getDate()

    days.push({
      date: dateStr,
      day: dayNum,
      weekday: weekday,
      position: dayIndex * dayWidth.value,  // 位置从0开始
      isToday: currentDate.toDateString() === today.toDateString(),
      isWeekend: currentDate.getDay() === 0 || currentDate.getDay() === 6,
      isOdd: dayNum % 2 === 1  // 奇数日期标记，用于只显示奇数日期
    })

    currentDate.setDate(currentDate.getDate() + 1)
    dayIndex++
  }

  console.log('timelineDays computed:', days.length, 'days, start:', days[0]?.date, 'end:', days[days.length - 1]?.date)
  return days
})

// 时间轴月份（用于上层显示）
const timelineMonths = computed(() => {
  if (timelineDays.value.length === 0) return []

  const months = []
  let currentMonth = null
  let monthStartIndex = 0

  timelineDays.value.forEach((day, index) => {
    const date = new Date(day.date)
    const year = date.getFullYear()
    const month = date.getMonth() + 1
    const monthKey = `${year}-${month}`

    if (monthKey !== currentMonth) {
      // 保存上一个月
      if (currentMonth !== null) {
        months[months.length - 1].endIndex = index - 1
        months[months.length - 1].width = (index - monthStartIndex) * dayWidth.value
      }

      // 开始新月份
      currentMonth = monthKey
      monthStartIndex = index
      months.push({
        key: monthKey,
        label: `${year}年${month}月`,
        startIndex: index,
        position: day.position,
        width: 0  // 稍后计算
      })
    }
  })

  // 处理最后一个月
  if (months.length > 0) {
    const lastMonth = months[months.length - 1]
    lastMonth.endIndex = timelineDays.value.length - 1
    lastMonth.width = (timelineDays.value.length - lastMonth.startIndex) * dayWidth.value
  }

  return months
})

// 今天的位置
const todayPosition = computed(() => {
  if (nodes.value.length === 0 || !props.scheduleData?.activities) return null

  const activities = Object.values(props.scheduleData.activities)
  if (activities.length === 0) return null

  // 过滤掉无效的活动数据
  const validActivities = activities.filter(a => {
    return a.earliest_start && a.earliest_start > 0 && a.earliest_finish && a.earliest_finish > 0
  })

  if (validActivities.length === 0) return null

  // 计算时间轴的起始日期
  const taskDates = validActivities.map(a => {
    const startDate = new Date(a.earliest_start * 1000)
    const endDate = new Date(a.earliest_finish * 1000)
    return { startDate, endDate }
  })

  const minDate = new Date(Math.min(...taskDates.map(d => d.startDate.getTime())))
  const maxDate = new Date(Math.max(...taskDates.map(d => d.endDate.getTime())))

  const bufferDays = 2
  const endBufferDays = 5
  const timelineStartDate = new Date(minDate)
  timelineStartDate.setDate(timelineStartDate.getDate() - bufferDays)

  const timelineEndDate = new Date(maxDate)
  timelineEndDate.setDate(timelineEndDate.getDate() + endBufferDays)

  // 计算今天的位置
  const today = new Date()
  today.setHours(0, 0, 0, 0)

  // 检查今天是否在时间轴范围内
  if (today >= timelineStartDate && today <= timelineEndDate) {
    const daysDiff = Math.round((today.getTime() - timelineStartDate.getTime()) / (1000 * 60 * 60 * 24))
    return daysDiff * dayWidth.value
  }

  return null
})

// 时间轴总宽度
const timelineWidth = computed(() => {
  if (timelineDays.value.length === 0) return 800
  return Math.max(800, timelineDays.value.length * dayWidth.value)
})

// 日期格式化工具函数
function formatDate(date, format) {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const weekdays = ['日', '一', '二', '三', '四', '五', '六']
  const weekday = weekdays[date.getDay()]

  if (format === 'YYYY-MM-DD') {
    return `${year}-${month}-${day}`
  }

  return `${month}/${day}`
}

// 计算总工期
const calculateTotalDuration = () => {
  if (nodes.value.length === 0) return 0
  const endNode = nodes.value.find(n => n.isEnd)
  return endNode ? endNode.EF : 0
}

// 根据节点编号获取节点类型标签
const getNodeTypeLabel = (nodeNumber) => {
  const node = nodes.value.find(n => n.number === nodeNumber)
  if (node) {
    if (node.nodeType === 'start') return '(起点S)'
    if (node.nodeType === 'end') return '(终点E)'
    if (node.nodeType === 'shared') return '(共享Sh)'
  }
  return ''
}

// PERT 图时间参数计算
const calculatePERTTimeParams = (activities, taskStartNodeMap, taskEndNodeMap) => {
  // 正向计算：计算最早时间 (ET - Earliest Time)
  // ET(i) = max{ET(j) + duration(j,i)} for all predecessors j of i

  // 反向计算：计算最迟时间 (LT - Latest Time)
  // LT(i) = min{LT(j) - duration(i,j)} for all successors j of i

  const nodeET = new Map() // 节点最早时间
  const nodeLT = new Map() // 节点最迟时间

  // 第一阶段：正向计算最早时间
  // 按照拓扑顺序处理节点
  const processed = new Set()
  let changed = true
  let iterations = 0
  const maxIterations = Object.keys(activities).length * 2

  while (changed && iterations < maxIterations) {
    changed = false
    iterations++

    for (const [key, activity] of Object.entries(activities)) {
      if (activity.is_dummy) continue

      const taskId = activity.task_id || activity.id
      const startNode = taskStartNodeMap.get(taskId)
      const endNode = taskEndNodeMap.get(taskId)
      if (!startNode || !endNode) continue

      // 计算起点节点的最早时间
      if (!processed.has(startNode.id)) {
        let maxET = 0
        if (activity.predecessors && activity.predecessors.length > 0) {
          activity.predecessors.forEach(predId => {
            const predEndNode = taskEndNodeMap.get(predId)
            if (predEndNode && nodeET.has(predEndNode.id)) {
              const predET = nodeET.get(predEndNode.id)
              const predActivity = Object.values(activities).find(a => (a.task_id || a.id) === predId)
              const duration = predActivity?.duration || 0
              maxET = Math.max(maxET, predET + duration)
            }
          })
        }
        nodeET.set(startNode.id, maxET)
        startNode.ET = maxET
        startNode.ES = maxET // ES = ET (对于事件节点)
        startNode.EF = maxET // EF = ES (对于起点事件)
        processed.add(startNode.id)
        changed = true
      }

      // 计算终点节点的最早时间
      if (!processed.has(endNode.id) && nodeET.has(startNode.id)) {
        const duration = activity.duration || 0
        const et = nodeET.get(startNode.id) + duration
        nodeET.set(endNode.id, et)
        endNode.ET = et
        endNode.ES = et
        endNode.EF = et
        processed.add(endNode.id)
        changed = true
      }
    }
  }

  // 第二阶段：反向计算最迟时间
  // 找到项目结束时间（所有终点的最大ET）
  let projectEndET = 0
  for (const [key, activity] of Object.entries(activities)) {
    if (activity.is_dummy) continue
    const taskId = activity.task_id || activity.id
    const endNode = taskEndNodeMap.get(taskId)
    if (endNode && nodeET.has(endNode.id)) {
      projectEndET = Math.max(projectEndET, nodeET.get(endNode.id))
    }
  }

  // 初始化所有终点的LT = projectEndET
  for (const [key, activity] of Object.entries(activities)) {
    if (activity.is_dummy) continue
    const taskId = activity.task_id || activity.id
    const endNode = taskEndNodeMap.get(taskId)
    if (endNode && nodeET.has(endNode.id) && activity.successors?.length === 0) {
      nodeLT.set(endNode.id, projectEndET)
      endNode.LT = projectEndET
      endNode.LF = projectEndET
      endNode.LS = projectEndET - (activity.duration || 0)
    }
  }

  // 反向传播最迟时间
  changed = true
  iterations = 0

  while (changed && iterations < maxIterations) {
    changed = false
    iterations++

    for (const [key, activity] of Object.entries(activities)) {
      if (activity.is_dummy) continue

      const taskId = activity.task_id || activity.id
      const startNode = taskStartNodeMap.get(taskId)
      const endNode = taskEndNodeMap.get(taskId)
      if (!startNode || !endNode) continue

      // 如果终点已经计算了LT，则计算起点的LT
      if (nodeLT.has(endNode.id) && !nodeLT.has(startNode.id)) {
        const duration = activity.duration || 0
        const lt = nodeLT.get(endNode.id) - duration
        nodeLT.set(startNode.id, lt)
        startNode.LT = lt
        startNode.LF = lt
        startNode.LS = lt
        changed = true
      }
    }
  }

  // 第三阶段：计算时差
  // 总时差 = LT - ET
  // 自由时差 = min(后置任务的ET) - 当前任务的EF
  for (const [key, activity] of Object.entries(activities)) {
    if (activity.is_dummy) continue

    const taskId = activity.task_id || activity.id
    const startNode = taskStartNodeMap.get(taskId)
    const endNode = taskEndNodeMap.get(taskId)
    if (!startNode || !endNode) continue

    // 总时差
    const totalSlack = nodeLT.has(endNode.id) && nodeET.has(endNode.id)
      ? nodeLT.get(endNode.id) - nodeET.get(endNode.id)
      : 0

    // 自由时差 = min(所有后置任务的ES) - 当前任务的EF
    let freeSlack = 0
    if (activity.successors && activity.successors.length > 0) {
      const minSuccessorES = Math.min(...activity.successors.map(succId => {
        const succStartNode = taskStartNodeMap.get(succId)
        return succStartNode && nodeET.has(succStartNode.id) ? nodeET.get(succStartNode.id) : Infinity
      }))
      freeSlack = minSuccessorES - nodeET.get(endNode.id)
    } else {
      freeSlack = totalSlack
    }

    endNode.totalSlack = Math.round(totalSlack * 10) / 10
    endNode.freeSlack = Math.round(freeSlack * 10) / 10
    endNode.slack = Math.round(totalSlack * 10) / 10

    // 判断是否在关键路径上（总时差为0）
    const isCritical = Math.abs(totalSlack) < 0.1
    endNode.isCritical = isCritical
    startNode.isCritical = isCritical
  }
}

// 解析进度数据并构建网络图 (PERT格式：共享节点)
const buildNetworkDiagram = () => {
  console.log('NetworkDiagram - buildNetworkDiagram called')
  console.log('NetworkDiagram - scheduleData:', props.scheduleData)
  console.log('NetworkDiagram - activities:', props.scheduleData?.activities)

  // 输出前几个活动的时间戳用于调试
  const activityArray = Object.entries(props.scheduleData?.activities || {})
  if (activityArray.length > 0) {
    console.log('NetworkDiagram - Sample activity timestamps:')
    activityArray.slice(0, 3).forEach(([key, activity]) => {
      const startDate = new Date(activity.earliest_start * 1000)
      const endDate = new Date(activity.earliest_finish * 1000)
      console.log(`  ${activity.name || key}: earliest_start=${activity.earliest_start} (${startDate.toISOString().split('T')[0]}), earliest_finish=${activity.earliest_finish} (${endDate.toISOString().split('T')[0]})`)
    })
  }

  if (!props.scheduleData || !props.scheduleData.activities) {
    nodes.value = []
    edges.value = []
    console.log('NetworkDiagram - No scheduleData or activities')
    return
  }

  const activities = props.scheduleData.activities || {}
  const nodeList = []
  const edgeList = []
  const taskStartNodeMap = new Map() // taskId -> start node
  const taskEndNodeMap = new Map() // taskId -> end node
  // 事件节点映射：key是任务ID组合，value是节点对象
  // 例如：{"1-end,2-start"} 表示任务1的终点和任务2的起点共享的节点
  const eventNodes = new Map()

  console.log('NetworkDiagram - Processing activities, count:', Object.keys(activities).length)

  // 节点编号计数器
  let nodeNumber = 1

  // 第一阶段：收集所有需要创建的事件节点
  // 每个任务的起点和终点都是一个事件节点
  // 如果任务A的终点 = 任务B的起点（A是B的前置），则共享节点

  // 首先创建所有任务的起点和终点节点
  for (const [key, activity] of Object.entries(activities)) {
    const taskId = activity.task_id || activity.id
    console.log('NetworkDiagram - Processing activity:', key, 'taskId:', taskId, 'name:', activity.name)

    // 跳过虚拟活动
    if (activity.is_dummy) {
      console.log('NetworkDiagram - Skipping dummy activity:', key)
      continue
    }

    // 跳过无效的活动数据（时间戳为 0 或负值）
    if (!activity.earliest_start || activity.earliest_start <= 0 ||
        !activity.earliest_finish || activity.earliest_finish <= 0) {
      console.log('NetworkDiagram - Skipping invalid activity (bad timestamps):', key, 'name:', activity.name,
                  'earliest_start:', activity.earliest_start, 'earliest_finish:', activity.earliest_finish)
      continue
    }

    // 创建任务的起点节点（暂时不与其他任务共享）
    const startNodeKey = `${taskId}-start`
    const startNode = {
      id: startNodeKey,
      taskId: taskId,
      name: activity.name,
      nodeType: 'start',
      number: nodeNumber,
      // ES存储为天数（从1970-01-01开始的天数）
      ES: Math.floor(activity.earliest_start / 86400),
      EF: Math.floor(activity.earliest_start / 86400),
      LS: Math.floor(activity.latest_start / 86400),
      LF: Math.floor(activity.latest_start / 86400),
      // 存储原始时间戳（秒），用于调试和日期计算
      earliest_start_ts: activity.earliest_start,
      slack: 0,
      totalSlack: 0,
      freeSlack: 0,
      isCritical: activity.is_critical,
      isStart: activity.predecessors?.length === 0,
      isEnd: false,
      isDummy: false,
      predecessors: [],
      successors: []
    }
    taskStartNodeMap.set(taskId, startNode)
    eventNodes.set(startNodeKey, startNode)
    nodeNumber++

    // 创建任务的终点节点
    const endNodeKey = `${taskId}-end`
    const endNode = {
      id: endNodeKey,
      taskId: taskId,
      name: activity.name,
      nodeType: 'end',
      number: nodeNumber,
      ES: Math.floor(activity.earliest_finish / 86400),
      EF: Math.floor(activity.earliest_finish / 86400),
      LS: Math.floor(activity.latest_finish / 86400),
      LF: Math.floor(activity.latest_finish / 86400),
      // 存储原始时间戳（秒），用于调试和日期计算
      earliest_finish_ts: activity.earliest_finish,
      slack: Math.floor((activity.latest_finish - activity.earliest_finish) / 86400),
      totalSlack: 0,
      freeSlack: 0,
      isCritical: activity.is_critical,
      isStart: false,
      isEnd: activity.successors?.length === 0,
      isDummy: false,
      predecessors: [],
      successors: []
    }
    taskEndNodeMap.set(taskId, endNode)
    eventNodes.set(endNodeKey, endNode)
    nodeNumber++
  }

  // 第二阶段：处理依赖关系，合并共享的节点
  // 如果任务B依赖任务A，那么A的终点和B的起点应该是同一个节点
  for (const [key, activity] of Object.entries(activities)) {
    const taskId = activity.task_id || activity.id

    if (activity.is_dummy) continue

    if (activity.predecessors && activity.predecessors.length > 0) {
      const currentTaskStartNode = taskStartNodeMap.get(taskId)

      activity.predecessors.forEach(predId => {
        const predTaskEndNode = taskEndNodeMap.get(predId)
        if (predTaskEndNode && currentTaskStartNode) {
          // 合并节点：前置任务的终点 = 当前任务的起点
          // 使用前置任务的终点作为共享节点
          const sharedNode = predTaskEndNode
          // 更新当前任务的起点为共享节点
          taskStartNodeMap.set(taskId, sharedNode)
          // 删除原来独立的起点节点
          eventNodes.delete(`${taskId}-start`)
          // 更新共享节点的类型和编号
          sharedNode.nodeType = 'shared'
          sharedNode.id = `shared-${predId}-end-${taskId}-start`
        }
      })
    }
  }

  // 将所有唯一的节点添加到节点列表
  const uniqueNodes = new Set()
  for (const [key, node] of eventNodes) {
    if (!uniqueNodes.has(node.number)) {
      uniqueNodes.add(node.number)
      nodeList.push(node)
    }
  }

  console.log('NetworkDiagram - Total unique nodes created:', nodeList.length)

  // 第三阶段：计算 PERT 时间参数
  calculatePERTTimeParams(activities, taskStartNodeMap, taskEndNodeMap)

  // 第四阶段：创建任务边和依赖边
  for (const [key, activity] of Object.entries(activities)) {
    const taskId = activity.task_id || activity.id

    // 处理虚活动（Dummy Activity）- 只显示连接线，不显示任务边
    if (activity.is_dummy) {
      // 虚活动：从前置任务的终点连接到当前任务的起点
      if (activity.predecessors && activity.predecessors.length > 0) {
        const currentTaskStartNode = taskStartNodeMap.get(taskId)
        activity.predecessors.forEach(predId => {
          const predTaskEndNode = taskEndNodeMap.get(predId)
          if (predTaskEndNode && currentTaskStartNode) {
            edgeList.push({
              id: `dummy-edge-${predId}-${taskId}`,
              from: predTaskEndNode.number,
              to: currentTaskStartNode.number,
              x1: 0,
              y1: 0,
              x2: 0,
              y2: 0,
              label: '', // 虚活动不显示标签
              duration: 0,
              isCritical: false, // 虚活动不在关键路径上
              isDummy: true,
              isTaskEdge: false // 标记为非任务边
            })
          }
        })
      }
      continue
    }

    const startNode = taskStartNodeMap.get(taskId)
    const endNode = taskEndNodeMap.get(taskId)

    if (startNode && endNode) {
      // 创建任务本身的边（从起点到终点）
      edgeList.push({
        id: `task-edge-${taskId}`,
        taskId: taskId,  // 添加 taskId 属性，用于编辑和删除
        from: startNode.number,
        to: endNode.number,
        x1: 0,
        y1: 0,
        x2: 0,
        y2: 0,
        label: activity.name || `任务${taskId}`,
        duration: activity.duration,
        isCritical: endNode.isCritical,
        isDummy: false,
        isTaskEdge: true,
        totalSlack: endNode.totalSlack,
        freeSlack: endNode.freeSlack
      })

      // 更新节点的前后置关系
      if (!startNode.successors.includes(endNode)) {
        startNode.successors.push(endNode)
      }
      if (!endNode.predecessors.includes(startNode)) {
        endNode.predecessors.push(startNode)
      }
    }
  }

  // 自动布局
  layoutNodes(nodeList, edgeList)

  console.log('NetworkDiagram - Before assignment - nodeList:', nodeList.length, 'edgeList:', edgeList.length)
  console.log('NetworkDiagram - nodeList sample:', nodeList.slice(0, 3))
  console.log('NetworkDiagram - edgeList sample:', edgeList.slice(0, 3))

  nodes.value = nodeList
  edges.value = edgeList

  console.log('NetworkDiagram - After assignment - nodes.value:', nodes.value.length, 'edges.value:', edges.value.length)
}

// 节点布局算法
const layoutNodes = (nodeList, edgeList) => {
  console.log('layoutNodes - called with', nodeList.length, 'nodes,', edgeList.length, 'edges')

  const xScale = dayWidth.value // 使用动态宽度，与时间轴保持一致
  const ySpacing = 60  // 减小垂直间距，从80改为60
  const nodeRadius = 18  // 减小节点半径，从25改为18

  // 计算时间轴的起始日期
  if (!props.scheduleData?.activities || nodeList.length === 0) {
    console.warn('layoutNodes - No activities or nodes')
    return
  }

  const activities = Object.values(props.scheduleData.activities)

  // 过滤掉无效的活动数据
  const validActivities = activities.filter(a => {
    const isValid = a.earliest_start && a.earliest_start > 0 && a.earliest_finish && a.earliest_finish > 0
    if (!isValid) {
      console.log('layoutNodes - Filtering out invalid activity:', a.name || a.id, 'earliest_start:', a.earliest_start, 'earliest_finish:', a.earliest_finish)
    }
    return isValid
  })

  if (validActivities.length === 0) {
    console.warn('layoutNodes - No valid activities found')
    return
  }

  const taskDates = validActivities.map(a => {
    const startDate = new Date(a.earliest_start * 1000)
    const endDate = new Date(a.earliest_finish * 1000)
    return { startDate, endDate }
  })

  const minDate = new Date(Math.min(...taskDates.map(d => d.startDate.getTime())))
  const bufferDays = 2
  const timelineStartDate = new Date(minDate)
  timelineStartDate.setDate(timelineStartDate.getDate() - bufferDays)

  console.log('layoutNodes - Timeline start date:', timelineStartDate.toISOString().split('T')[0])

  // 添加调试：检查时间轴日期范围
  console.log('layoutNodes - Time range will be from', timelineStartDate.toISOString().split('T')[0], 'to',
    new Date(timelineStartDate.getTime() + (timelineDays.value.length * 86400000)).toISOString().split('T')[0])

  // 计算每个节点的依赖深度（拓扑排序）
  const nodeDepth = new Map()
  const calculateDepth = (nodeNumber, visited = new Set()) => {
    if (visited.has(nodeNumber)) return 0
    visited.add(nodeNumber)

    const node = nodeList.find(n => n.number === nodeNumber)
    if (!node) return 0

    const incomingEdges = edgeList.filter(e => e.to === nodeNumber)
    if (incomingEdges.length === 0) {
      return 0
    }

    let maxDepth = 0
    for (const edge of incomingEdges) {
      const depth = calculateDepth(edge.from, visited)
      if (depth > maxDepth) {
        maxDepth = depth
      }
    }
    return maxDepth + 1
  }

  nodeList.forEach(node => {
    const depth = calculateDepth(node.number)
    nodeDepth.set(node.number, depth)
  })

  // 分离有连接和无连接的节点
  const nodesWithConnections = nodeList.filter(n =>
    n.predecessors.length > 0 || n.successors.length > 0
  )
  const nodesWithoutConnections = nodeList.filter(n =>
    n.predecessors.length === 0 && n.successors.length === 0
  )

  console.log('nodesWithConnections:', nodesWithConnections.length)
  console.log('nodesWithoutConnections:', nodesWithoutConnections.length)

  // 第一阶段：布局有连接的节点（基于实际日期和依赖深度）
  // 关键改进：同一个任务的起点和终点使用相同的Y坐标

  // 按任务ID分组节点（每个任务有起点和终点）
  const taskNodesMap = new Map()  // taskId -> { startNode, endNode }
  const allTaskIds = new Set()

  nodesWithConnections.forEach(node => {
    if (node.taskId) {
      allTaskIds.add(node.taskId)
      if (!taskNodesMap.has(node.taskId)) {
        taskNodesMap.set(node.taskId, { startNode: null, endNode: null })
      }
      const taskNodes = taskNodesMap.get(node.taskId)
      if (node.nodeType === 'start' || node.nodeType === 'shared') {
        // 对于共享节点，它既是某个任务的终点也是另一个任务的起点
        // 我们用它作为终点的任务的Y坐标
        if (!taskNodes.startNode) {
          taskNodes.startNode = node
        }
      }
      if (node.nodeType === 'end' || node.nodeType === 'shared') {
        taskNodes.endNode = node
      }
    }
  })

  const depthYMap = new Map()
  const usedYPositions = new Set() // 跟踪已使用的Y位置

  // 使用 Map 存储新位置，避免直接修改原对象
  const nodePositionMap = new Map()

  // 为每个任务分配Y坐标
  let taskIndex = 0
  allTaskIds.forEach(taskId => {
    const taskNodes = taskNodesMap.get(taskId)
    if (!taskNodes) return

    const depth = taskNodes.startNode ? (nodeDepth.get(taskNodes.startNode.number) || 0) : 0
    const endDepth = taskNodes.endNode ? (nodeDepth.get(taskNodes.endNode.number) || 0) : 0
    const maxDepth = Math.max(depth, endDepth)

    // 为这个任务分配一个Y坐标
    let yPos
    let attempts = 0
    do {
      yPos = headerHeight + 60 + (taskIndex + attempts) * ySpacing
      attempts++
    } while (usedYPositions.has(yPos) && attempts < 10)

    usedYPositions.add(yPos)

    // 计算起点和终点的X坐标
    const startNode = taskNodes.startNode
    const endNode = taskNodes.endNode

    if (startNode) {
      // X坐标：基于节点的实际日期
      let nodeDate
      if (startNode.earliest_start_ts) {
        nodeDate = new Date(startNode.earliest_start_ts * 1000)
      } else {
        nodeDate = new Date('1970-01-01')
        nodeDate.setTime(nodeDate.getTime() + startNode.ES * 86400 * 1000)
      }
      const daysDiff = Math.round((nodeDate.getTime() - timelineStartDate.getTime()) / (1000 * 60 * 60 * 24))
      const xPos = Math.max(0, daysDiff * xScale)

      nodePositionMap.set(startNode.number, { x: xPos, y: yPos })

      console.log(`Task ${taskId} Start Node: ES=${startNode.ES}, date=${nodeDate.toISOString().split('T')[0]}, x=${xPos}, y=${yPos}`)
    }

    if (endNode && endNode !== startNode) {
      // X坐标：基于节点的实际日期
      let nodeDate
      if (endNode.earliest_finish_ts) {
        nodeDate = new Date(endNode.earliest_finish_ts * 1000)
      } else if (endNode.earliest_start_ts) {
        nodeDate = new Date(endNode.earliest_start_ts * 1000)
      } else {
        nodeDate = new Date('1970-01-01')
        nodeDate.setTime(nodeDate.getTime() + endNode.ES * 86400 * 1000)
      }
      const daysDiff = Math.round((nodeDate.getTime() - timelineStartDate.getTime()) / (1000 * 60 * 60 * 24))
      const xPos = Math.max(0, daysDiff * xScale)

      nodePositionMap.set(endNode.number, { x: xPos, y: yPos })  // 相同的Y坐标

      console.log(`Task ${taskId} End Node: ES=${endNode.ES}, date=${nodeDate.toISOString().split('T')[0]}, x=${xPos}, y=${yPos}`)
    }

    taskIndex++
  })

  // 应用位置更新到节点对象（创建新对象）
  nodesWithConnections.forEach(node => {
    const pos = nodePositionMap.get(node.number)
    if (pos) {
      node.x = pos.x
      node.y = pos.y
    } else {
      // 如果节点没有被分配位置（可能是不完整的任务数据），给它分配一个位置
      let nodeDate
      if (node.earliest_finish_ts) {
        nodeDate = new Date(node.earliest_finish_ts * 1000)
      } else if (node.earliest_start_ts) {
        nodeDate = new Date(node.earliest_start_ts * 1000)
      } else {
        nodeDate = new Date('1970-01-01')
        nodeDate.setTime(nodeDate.getTime() + node.ES * 86400 * 1000)
      }
      const daysDiff = Math.round((nodeDate.getTime() - timelineStartDate.getTime()) / (1000 * 60 * 60 * 24))
      const xPos = Math.max(0, daysDiff * xScale)
      const yPos = headerHeight + 60 + taskIndex * ySpacing
      nodePositionMap.set(node.number, { x: xPos, y: yPos })
      node.x = xPos
      node.y = yPos
      taskIndex++
    }
  })

  // 第二阶段：布局无连接的节点（垂直平铺在最前面）
  const startX = 50 // 最前面的X位置
  nodesWithoutConnections.forEach((node, index) => {
    node.x = startX
    node.y = headerHeight + 60 + taskIndex * ySpacing
    nodePositionMap.set(node.number, { x: startX, y: headerHeight + 60 + taskIndex * ySpacing })
    taskIndex++
    console.log(`Isolated Node ${node.taskId}: x=${startX}, y=${headerHeight + 60 + (taskIndex - 1) * ySpacing}`)
  })

  // 更新边的坐标和计算正交路径（创建新对象以触发响应式更新）
  for (let i = 0; i < edgeList.length; i++) {
    const edge = edgeList[i]
    const fromNode = nodeList.find(n => n.number === edge.from)
    const toNode = nodeList.find(n => n.number === edge.to)
    if (fromNode && toNode) {
      // 创建新的边对象以触发 Vue 响应式更新
      edgeList[i] = {
        ...edge,
        x1: fromNode.x,
        y1: fromNode.y,
        x2: toNode.x,
        y2: toNode.y,
        path: calculateOrthogonalPath(
          fromNode.x,
          fromNode.y,
          toNode.x,
          toNode.y,
          nodeRadius,
          nodeList,
          fromNode.id,
          toNode.id
        ),
        labelX: (fromNode.x + toNode.x) / 2,
        labelY: (fromNode.y + toNode.y) / 2 - 8
      }
      console.log(`Edge ${edge.id}: path calculated`)
    } else {
      console.warn(`Edge ${edge.id} missing nodes: from=${edge.from}, to=${edge.to}`)
    }
  }

  // 输出节点位置摘要，用于调试
  const nodeXPositions = nodeList.map(n => n.x).filter(x => x !== undefined)
  if (nodeXPositions.length > 0) {
    const minX = Math.min(...nodeXPositions)
    const maxX = Math.max(...nodeXPositions)
    console.log(`layoutNodes - Node X position range: min=${minX}, max=${maxX}, span=${maxX - minX}px`)
    console.log(`layoutNodes - Timeline day width: ${xScale}px, total days: ${timelineDays.value.length}`)
  }

  console.log('layoutNodes - edgeList after positioning:', edgeList)
}

// 计算正交连线路径（使用A*算法避开节点障碍物）
const calculateOrthogonalPath = (x1, y1, x2, y2, radius, allNodes = [], fromNodeId = null, toNodeId = null) => {
  // 从源节点圆周上的点开始
  // 到目标节点圆周外的点结束（留出箭头空间）

  // 箭头长度
  const arrowLength = 9

  // 计算起点（从源节点右侧出发）
  const startX = x1 + radius
  const startY = y1

  // 计算终点（到目标节点左侧，留出箭头空间）
  // 路径终点在节点左侧，箭头会从路径终点延伸到节点
  const endX = x2 - radius - 1  // 留1px间隙确保不穿过节点
  const endY = y2

  // 构建障碍物数组（排除源节点和目标节点）
  const obstacles = buildNodeObstacles(allNodes, radius, fromNodeId, toNodeId)

  // 使用A*算法计算路径
  try {
    return findOrthogonalPath(startX, startY, endX, endY, obstacles)
  } catch (error) {
    console.warn('Path finding failed, using fallback:', error)
    // 降级方案：使用简单的直角折线
    if (endX > startX) {
      const midX = (startX + endX) / 2
      return `M ${startX} ${startY} L ${midX} ${startY} L ${midX} ${endY} L ${endX} ${endY}`
    }
    const verticalDistance = Math.abs(endY - startY)
    const bendOffset = Math.max(30, verticalDistance / 2 + 20)
    const midY = startY < endY ? startY - bendOffset : startY + bendOffset
    const rightMostX = Math.max(startX, endX) + 30
    return `M ${startX} ${startY} L ${rightMostX} ${startY} L ${rightMostX} ${midY} L ${rightMostX} ${endY} L ${endX} ${endY}`
  }
}

// 构建节点障碍物数组（用于路径规划）
const buildNodeObstacles = (nodes, radius, excludeFromId, excludeToId) => {
  if (!nodes || nodes.length === 0) return []

  const obstacles = []
  const padding = radius + 10  // 节点半径 + 额外填充

  for (const node of nodes) {
    // 排除源节点和目标节点
    if (node.id === excludeFromId || node.id === excludeToId) continue
    // 排除没有坐标的节点
    if (node.x === undefined || node.y === undefined) continue

    // 将圆形节点转换为矩形障碍物（扩大检测范围）
    obstacles.push({
      minX: node.x - padding,
      maxX: node.x + padding,
      minY: node.y - padding,
      maxY: node.y + padding
    })
  }

  return obstacles
}

// 渲染
const render = () => {
  buildNetworkDiagram()
}

// 处理布局变化
const handleLayoutChange = () => {
  render()
}

// 缩放控制
const zoomIn = () => {
  const newZoom = Math.min(zoomLevel.value + 0.1, 3)
  zoomLevel.value = newZoom
  // 同步更新dayWidth，时间标尺也会相应缩放
  dayWidth.value = 40 * newZoom
}

const zoomOut = () => {
  const newZoom = Math.max(zoomLevel.value - 0.1, 0.3)
  zoomLevel.value = newZoom
  // 同步更新dayWidth，时间标尺也会相应缩放
  dayWidth.value = 40 * newZoom
}

const resetZoom = () => {
  zoomLevel.value = 1
  dayWidth.value = 40  // 重置为默认值
  centerView()
}

// 适应视图
const fitView = () => {
  zoomLevel.value = 0.8
  dayWidth.value = 40 * 0.8  // 同步更新dayWidth
  centerView()
  ElMessage.success('视图已适应')
}

// 居中视图
const centerView = () => {
  console.log('centerView called')
  if (!containerRef.value || !canvasContainerRef.value) {
    console.log('container refs not ready')
    return
  }

  const containerRect = canvasContainerRef.value.getBoundingClientRect()
  console.log('containerRect:', containerRect.width, 'x', containerRect.height)

  // 简单策略：直接设置 panX 和 panY 为 0
  // 这样节点会从 SVG 原点开始绘制
  panX.value = 0
  panY.value = 0

  console.log('set panX=0, panY=0')
}

// 鼠标滚轮缩放
const handleWheel = (event) => {
  const delta = event.deltaY > 0 ? -0.1 : 0.1
  const newZoom = Math.max(0.3, Math.min(3, zoomLevel.value + delta))
  zoomLevel.value = newZoom
  // 同步更新dayWidth，时间标尺也会相应缩放
  dayWidth.value = 40 * newZoom
}

// 设置工具模式
const setToolMode = (mode) => {
  toolMode.value = mode
  // 清理连线状态
  if (mode !== 'connect') {
    isConnecting.value = false
    connectStartNode.value = null
  }
  // 更新容器类
  if (containerRef.value) {
    containerRef.value.classList.remove('pan-mode', 'connect-mode')
    if (mode === 'pan') {
      containerRef.value.classList.add('pan-mode')
    } else if (mode === 'connect') {
      containerRef.value.classList.add('connect-mode')
    }
  }
}

// 节点鼠标进入（用于连线模式）
const handleNodeMouseEnter = (event, node) => {
  if (connectMode.value && isConnecting.value) {
    // 使用 toRaw 确保获取原始坐标
    const rawNode = toRaw(node)
    // 防御性检查：确保节点坐标有效
    if (rawNode.x !== undefined && rawNode.y !== undefined && rawNode.x >= 0 && rawNode.y >= 0) {
      // 只在位置真正改变时才更新
      if (connectEndPosition.value.x !== rawNode.x || connectEndPosition.value.y !== rawNode.y) {
        connectEndPosition.value = { x: rawNode.x, y: rawNode.y }
      }
    }
  }
}

// 节点鼠标离开
const handleNodeMouseLeave = (node) => {
  // 连线模式下，鼠标离开节点时不做特殊处理
}

// 鼠标拖拽
const handleMouseDown = (event) => {
  // 在平移模式下，任何地方都可以开始平移
  // 在非平移模式下，只在空白区域才允许平移
  if (panMode.value || event.target.tagName === 'svg') {
    isPanning.value = true
    // 添加类以禁用过渡，提高拖拽性能
    containerRef.value?.classList.add('is-panning')
    panStartPos.value = {
      x: event.clientX - panX.value,
      y: event.clientY - panY.value
    }
  }
}

const handleMouseMove = (event) => {
  if (isPanning.value) {
    panX.value = event.clientX - panStartPos.value.x
    panY.value = event.clientY - panStartPos.value.y
  }

  // 连线模式下更新连线终点
  if ((connectMode.value || isCtrlPressed.value) && isConnecting.value) {
    // 将鼠标坐标转换为SVG坐标
    const svgRect = svgRef.value.getBoundingClientRect()
    const svgX = event.clientX - svgRect.left - panX.value
    const svgY = event.clientY - svgRect.top - panY.value

    // 检测鼠标下的节点
    const hoverRadius = nodeRadius.value * 1.5  // 增加感应范围
    const hoveredNode = nodes.value.find(node => {
      // 使用 toRaw 确保获取原始坐标，避免响应式代理问题
      const rawNode = toRaw(node)
      // 防御性检查：确保坐标有效
      if (rawNode.x === undefined || rawNode.y === undefined || rawNode.x < 0 || rawNode.y < 0) {
        return false
      }
      const dx = rawNode.x - svgX
      const dy = rawNode.y - svgY
      return Math.sqrt(dx * dx + dy * dy) < hoverRadius
    })

    if (hoveredNode && hoveredNode.id !== connectStartNode.value?.id) {
      // 悬停在目标节点上，连线终点指向节点中心
      const rawHoveredNode = toRaw(hoveredNode)
      // 防御性检查：确保坐标有效
      if (rawHoveredNode.x !== undefined && rawHoveredNode.y !== undefined &&
          rawHoveredNode.x >= 0 && rawHoveredNode.y >= 0) {
        // 只在位置真正改变时才更新
        if (connectEndPosition.value.x !== rawHoveredNode.x || connectEndPosition.value.y !== rawHoveredNode.y) {
          connectHoverNode.value = hoveredNode
          connectEndPosition.value = { x: rawHoveredNode.x, y: rawHoveredNode.y }
        }
      }
    } else {
      // 没有悬停在节点上，连线终点跟随鼠标
      if (connectHoverNode.value !== null) {
        connectHoverNode.value = null
      }
      // 只在位置真正改变时才更新
      if (connectEndPosition.value.x !== svgX || connectEndPosition.value.y !== svgY) {
        connectEndPosition.value = { x: svgX, y: svgY }
      }
    }
  } else {
    connectHoverNode.value = null
  }

  if (isDragging.value && draggedNode.value) {
    const node = draggedNode.value
    const scale = zoomLevel.value

    // 计算鼠标移动的差值
    const deltaX = (event.clientX - dragStartPos.value.x)
    const deltaY = (event.clientY - dragStartPos.value.y)

    // 在选择模式下，允许沿X轴和Y轴拖动节点
    if (!panMode.value && !connectMode.value) {
      // 允许水平和垂直拖动
      const newX = node.originalX + deltaX
      const newY = node.originalY + deltaY

      // 限制拖动范围，不能拖到负数位置
      if (newX >= 0) {
        node.x = newX
      }
      // Y轴允许负值，但可以设置最小值
      if (newY >= -50) {
        node.y = newY
      }

      // 检测鼠标下的其他节点（用于节点合并）
      const svgRect = svgRef.value.getBoundingClientRect()
      const svgX = event.clientX - svgRect.left - panX.value
      const svgY = event.clientY - svgRect.top - panY.value
      const hoverRadius = nodeRadius.value * 2  // 拖动时使用更大的感应范围

      const hoveredNode = nodes.value.find(n => {
        if (n.id === node.id) return false  // 排除自己
        const rawNode = toRaw(n)
        if (rawNode.x === undefined || rawNode.y === undefined || rawNode.x < 0 || rawNode.y < 0) {
          return false
        }
        const dx = rawNode.x - svgX
        const dy = rawNode.y - svgY
        return Math.sqrt(dx * dx + dy * dy) < hoverRadius
      })

      if (hoveredNode && hoveredNode.id !== node.id) {
        dragHoverNode.value = hoveredNode
      } else {
        dragHoverNode.value = null
      }

      // 更新相关边的坐标和路径
      const nodeRadiusVal = nodeRadius.value
      edges.value.forEach(edge => {
        if (edge.from === node.number) {
          edge.x1 = node.x
          edge.y1 = node.y
          const toNode = nodes.value.find(n => n.number === edge.to)
          if (toNode) {
            edge.path = calculateOrthogonalPath(
              node.x,
              node.y,
              toNode.x,
              toNode.y,
              nodeRadiusVal,
              nodes.value,
              node.id,
              toNode.id
            )
            edge.labelX = (node.x + toNode.x) / 2
            edge.labelY = (node.y + toNode.y) / 2 - 8
          }
        }
        if (edge.to === node.number) {
          edge.x2 = node.x
          edge.y2 = node.y
          const fromNode = nodes.value.find(n => n.number === edge.from)
          if (fromNode) {
            edge.path = calculateOrthogonalPath(
              fromNode.x,
              fromNode.y,
              node.x,
              node.y,
              nodeRadiusVal,
              nodes.value,
              fromNode.id,
              node.id
            )
            edge.labelX = (fromNode.x + node.x) / 2
            edge.labelY = (fromNode.y + node.y) / 2 - 8
          }
        }
      })
    }
  }
}

const handleMouseUp = async () => {
  // 连线模式下松开鼠标，完成连线
  if (isConnecting.value && connectStartNode.value && connectHoverNode.value) {
    if (connectStartNode.value.id !== connectHoverNode.value.id) {
      // 完成连线，弹出新建任务对话框
      openCreateTaskDialog(connectStartNode.value, connectHoverNode.value)
    }
    // 重置连线状态
    isConnecting.value = false
    connectStartNode.value = null
    connectHoverNode.value = null
    connectEndPosition.value = { x: 0, y: 0 }
    return
  }

  if (isDragging.value && draggedNode.value) {
    const node = draggedNode.value
    const hasXChanged = node.originalX !== undefined && node.x !== node.originalX
    const hasYChanged = node.originalY !== undefined && node.y !== node.originalY

    // 检查是否释放到另一个节点上（节点合并）
    if (dragHoverNode.value && dragHoverNode.value.id !== node.id) {
      await mergeNodes(node, dragHoverNode.value)
      // 重置拖动状态
      isDragging.value = false
      draggedNode.value = null
      dragHoverNode.value = null
      containerRef.value?.classList.remove('is-panning')
      return
    }

    // 如果 X 坐标发生了变化（水平拖动），尝试更新任务日期
    if (hasXChanged) {
      await updateTaskDateFromNodePosition(node)

      // 更新任务日期后，节点保持在当前位置
      // 不触发整个网络图的重新计算，这样同一个任务的另一个节点不会受到影响
      // 下次刷新时会基于新日期重新计算位置
    } else if (hasYChanged) {
      // 如果只是 Y 坐标变化（垂直拖动），保存新的位置
      emit('position-updated', {
        nodes: nodes.value.map(n => ({
          id: n.id,
          x: n.x,
          y: n.y
        }))
      })
    }
  }
  isDragging.value = false
  isPanning.value = false
  draggedNode.value = null
  dragHoverNode.value = null
  // 移除类以恢复平滑过渡
  containerRef.value?.classList.remove('is-panning')
}

// 根据节点位置更新任务日期
const updateTaskDateFromNodePosition = async (node) => {
  try {
    // 拖动节点会影响所有连接到该节点的任务边
    // 找到所有与该节点相关的任务边
    const relatedEdges = edges.value.filter(edge =>
      edge.isTaskEdge && (edge.from === node.number || edge.to === node.number)
    )

    if (relatedEdges.length === 0) {
      console.warn('No task edges found for node:', node.number)
      return
    }

    console.log(`Updating ${relatedEdges.length} task(s) for node ${node.number} (type: ${node.nodeType})`)

    // 计算时间轴的起始日期
    const activities = props.scheduleData?.activities
    if (!activities) return

    const validActivities = Object.values(activities).filter(a => {
      return a.earliest_start && a.earliest_start > 0 && a.earliest_finish && a.earliest_finish > 0
    })

    if (validActivities.length === 0) return

    const taskDates = validActivities.map(a => {
      const startDate = new Date(a.earliest_start * 1000)
      const endDate = new Date(a.earliest_finish * 1000)
      return { startDate, endDate }
    })

    const minDate = new Date(Math.min(...taskDates.map(d => d.startDate.getTime())))
    const bufferDays = 2
    const timelineStartDate = new Date(minDate)
    timelineStartDate.setDate(timelineStartDate.getDate() - bufferDays)

    // 计算节点拖动的偏移量（天数）
    const deltaX = node.x - node.originalX
    const daysDelta = Math.round(deltaX / dayWidth.value)

    console.log(`Node ${node.number}: originalX=${node.originalX}, newX=${node.x}, deltaX=${deltaX}, daysDelta=${daysDelta}`)

    // 更新所有相关任务的日期
    for (const edge of relatedEdges) {
      if (!edge.taskId) continue

      // 找到对应的活动数据
      const activity = validActivities.find(a => a.task_id === edge.taskId || a.id === edge.taskId)
      if (!activity) continue

      // 计算原始日期和持续时间
      const originalStartDate = new Date(activity.earliest_start * 1000)
      const originalEndDate = new Date(activity.earliest_finish * 1000)
      const originalDuration = Math.round((originalEndDate.getTime() - originalStartDate.getTime()) / (1000 * 60 * 60 * 24))

      let newStartDate, newEndDate

      // 根据节点在任务边中的角色和节点类型来确定如何更新日期
      if (node.nodeType === 'start' || (node.nodeType === 'shared' && edge.to === node.number)) {
        // 节点作为任务的起点：拖动时只更新开始日期，结束日期保持不变
        // 这样任务的持续时间会改变，但终点节点不会移动
        newStartDate = new Date(originalStartDate.getTime() + daysDelta * 86400 * 1000)
        newEndDate = originalEndDate  // 保持结束日期不变
      } else if (node.nodeType === 'end' || (node.nodeType === 'shared' && edge.from === node.number)) {
        // 节点作为任务的终点：拖动时只更新结束日期，开始日期保持不变
        // 这样任务的持续时间会改变，但起点节点不会移动
        newEndDate = new Date(originalEndDate.getTime() + daysDelta * 86400 * 1000)
        newStartDate = originalStartDate  // 保持开始日期不变
      } else {
        // 默认：从新位置计算开始日期，结束日期相应延长
        newStartDate = new Date(originalStartDate.getTime() + daysDelta * 86400 * 1000)
        newEndDate = new Date(newStartDate.getTime() + originalDuration * 86400 * 1000)
      }

      // 更新任务
      const updateData = {
        start_date: formatDate(newStartDate, 'YYYY-MM-DD'),
        end_date: formatDate(newEndDate, 'YYYY-MM-DD')
      }

      console.log(`Updating task ${edge.taskId}:`, updateData)

      await progressApi.update(edge.taskId, updateData)

      // 触发数据重新加载，让前端显示最新的任务数据
      emit('task-updated', { silent: true })
    }

    ElMessage.success(`已更新 ${relatedEdges.length} 个任务的日期`)
  } catch (error) {
    console.error('更新任务日期失败:', error)
    ElMessage.error('更新失败：' + (error.response?.data?.message || error.message || '未知错误'))

    // 出错时恢复原位置（包括 X 和 Y 坐标）
    if (node.originalX !== undefined) {
      node.x = node.originalX
    }
    if (node.originalY !== undefined) {
      node.y = node.originalY
    }
  }
}

// 合并节点
const mergeNodes = async (sourceNode, targetNode) => {
  try {
    // 显示确认对话框
    await ElMessageBox.confirm(
      `确定要将节点 ${sourceNode.number} 合并到节点 ${targetNode.number} 吗？\n\n这将创建任务依赖关系并合并两个节点。`,
      '合并节点',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    console.log(`Merging node ${sourceNode.number} to ${targetNode.number}`)

    // 找到源节点和目标节点相关的所有任务边
    const sourceNodeEdges = edges.value.filter(edge =>
      edge.isTaskEdge && (edge.from === sourceNode.number || edge.to === sourceNode.number)
    )
    const targetNodeEdges = edges.value.filter(edge =>
      edge.isTaskEdge && (edge.from === targetNode.number || edge.to === targetNode.number)
    )

    console.log(`Source node has ${sourceNodeEdges.length} edges, target node has ${targetNodeEdges.length} edges`)

    // 创建依赖关系：源节点结束的任务 -> 目标节点开始的任务
    let dependenciesCreated = 0
    for (const sourceEdge of sourceNodeEdges) {
      if (!sourceEdge.taskId) continue
      // 源边中，节点作为终点的任务（需要创建到目标节点的依赖）
      if (sourceEdge.to === sourceNode.number) {
        for (const targetEdge of targetNodeEdges) {
          if (!targetEdge.taskId) continue
          // 目标边中，节点作为起点的任务（作为依赖的目标）
          if (targetEdge.from === targetNode.number) {
            try {
              await progressApi.createDependencyVisual(sourceEdge.taskId, targetEdge.taskId, { type: 'FS', lag: 0 })
              dependenciesCreated++
              console.log(`Created dependency: task ${sourceEdge.taskId} -> task ${targetEdge.taskId}`)
            } catch (error) {
              console.error(`Failed to create dependency between ${sourceEdge.taskId} and ${targetEdge.taskId}:`, error)
            }
          }
        }
      }
    }

    ElMessage.success(`已合并节点并创建 ${dependenciesCreated} 个依赖关系`)

    // 触发网络图重新计算
    emit('task-updated', { reloadOnly: true })
  } catch (error) {
    if (error !== 'cancel') {
      console.error('合并节点失败:', error)
      ElMessage.error('合并失败：' + (error.response?.data?.message || error.message || '未知错误'))
    }
    // 用户取消或出错，恢复源节点位置
    sourceNode.x = sourceNode.originalX
    sourceNode.y = sourceNode.originalY
  }
}

// 节点拖拽
const handleNodeMouseDown = (event, node) => {
  event.stopPropagation()

  // 连线模式或按住Ctrl键：开始连线
  if (connectMode.value || isCtrlPressed.value) {
    if (!isConnecting.value) {
      // 使用 toRaw 确保获取原始节点
      const rawNode = toRaw(node)
      // 防御性检查：确保节点坐标有效
      if (rawNode.x !== undefined && rawNode.y !== undefined && rawNode.x >= 0 && rawNode.y >= 0) {
        // 开始连线
        isConnecting.value = true
        connectStartNode.value = node
        connectEndPosition.value = { x: rawNode.x, y: rawNode.y }
        selectedNode.value = node  // 同时选中起始节点
      } else {
        console.error('handleNodeMouseDown - Invalid start node coordinates:', {
          nodeId: node.id,
          nodeX: rawNode.x,
          nodeY: rawNode.y
        })
      }
    }
    return
  }

  // 平移模式不处理节点点击
  if (panMode.value) {
    return
  }

  // 选择模式下允许拖动节点
  isDragging.value = true
  draggedNode.value = node

  // 保存拖拽开始时的鼠标位置和节点初始位置
  dragStartPos.value = {
    x: event.clientX,
    y: event.clientY
  }

  // 保存节点的初始位置
  node.originalX = node.x
  node.originalY = node.y
}

// 节点点击 - 选中节点或完成连线
const handleNodeClick = (node) => {
  // 连线模式：完成连线
  if ((connectMode.value || isCtrlPressed.value) && isConnecting.value) {
    if (connectStartNode.value && connectStartNode.value.id !== node.id) {
      // 完成连线，弹出新建任务对话框
      openCreateTaskDialog(connectStartNode.value, node)
    }
    // 重置连线状态
    isConnecting.value = false
    connectStartNode.value = null
    connectHoverNode.value = null
    return
  }

  // 选择模式：选中节点
  if (!panMode.value) {
    selectedNode.value = node
    emit('node-selected', node)
  }
}

// 右键菜单处理
const handleContextMenu = (event, type, data) => {
  event.preventDefault()
  event.stopPropagation()

  // 关闭之前的菜单
  contextMenuVisible.value = false

  // 使用 nextTick 确保 DOM 更新后再显示菜单
  nextTick(() => {
    contextMenuTarget.value = { type, data }
    contextMenuPosition.value = {
      x: event.clientX,
      y: event.clientY
    }
    contextMenuVisible.value = true
  })

  // 选中对应的节点或边
  if (type === 'node') {
    selectedNode.value = data
  } else if (type === 'edge') {
    selectedEdge.value = data
  }
}

// 处理右键菜单操作
const handleContextMenuAction = async (action) => {
  contextMenuVisible.value = false

  const target = contextMenuTarget.value
  if (!target) return

  try {
    switch (action) {
      case 'edit-node':
        await handleEditNode(target.data)
        break
      case 'view-details':
        if (target.type === 'node') {
          nodeDetailVisible.value = true
        } else if (target.type === 'edge') {
          edgeDetailVisible.value = true
        }
        break
      case 'delete-node':
        await handleDeleteNode(target.data)
        break
      case 'edit-edge':
        await handleEditEdge(target.data)
        break
      case 'delete-edge':
        await handleDeleteEdge(target.data)
        break
    }
  } catch (error) {
    console.error('操作失败:', error)
    ElMessage.error('操作失败：' + (error.response?.data?.message || error.message || '未知错误'))
  }
}

// 编辑节点
const handleEditNode = async (node) => {
  if (!node.taskId) {
    ElMessage.warning('该节点没有关联任务，无法编辑')
    return
  }

  // 从 scheduleData 中查找任务数据
  const activities = Object.values(props.scheduleData?.activities || {})
  const activity = activities.find(a => a.task_id === node.taskId || a.id === node.taskId)

  if (!activity) {
    ElMessage.error('未找到关联的任务数据')
    return
  }

  // 填充表单
  const startDate = new Date(activity.earliest_start * 1000)
  const endDate = new Date(activity.earliest_finish * 1000)

  taskForm.value = {
    name: activity.name || '',
    start: startDate,
    end: endDate,
    progress: activity.progress || 0,
    priority: activity.priority || 'medium',
    notes: activity.description || ''
  }

  // 保存正在编辑的任务 ID
  editingTask.value = { ...activity, id: activity.task_id || activity.id }

  editDialogVisible.value = true
}

// 删除节点
const handleDeleteNode = async (node) => {
  if (!node.taskId) {
    ElMessage.warning('该节点没有关联任务，无法删除')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要删除节点 ${node.number} 及其关联的任务吗？`,
      '确认删除',
      {
        confirmButtonText: '删除',
        cancelButtonText: '取消',
        type: 'warning',
        confirmButtonClass: 'el-button--danger'
      }
    )

    ElMessage.info('正在删除...')

    // 调用API删除任务
    await progressApi.delete(node.taskId)

    ElMessage.success('删除成功')

    // 触发刷新
    emit('task-updated', { deleted: true, id: node.taskId })
  } catch (error) {
    if (error !== 'cancel') {
      throw error
    }
  }
}

// 编辑连线
const handleEditEdge = async (edge) => {
  if (!edge.taskId) {
    ElMessage.warning('该连线没有关联任务，无法编辑')
    return
  }

  // 从 scheduleData 中查找任务数据
  const activities = Object.values(props.scheduleData?.activities || {})
  const activity = activities.find(a => a.task_id === edge.taskId || a.id === edge.taskId)

  if (!activity) {
    ElMessage.error('未找到关联的任务数据')
    return
  }

  // 填充表单
  const startDate = new Date(activity.earliest_start * 1000)
  const endDate = new Date(activity.earliest_finish * 1000)

  taskForm.value = {
    name: activity.name || '',
    start: startDate,
    end: endDate,
    progress: activity.progress || 0,
    priority: activity.priority || 'medium',
    notes: activity.description || ''
  }

  editingTask.value = { ...activity, id: activity.task_id || activity.id }

  editDialogVisible.value = true
}

// 删除连线
const handleDeleteEdge = async (edge) => {
  if (!edge.taskId) {
    ElMessage.warning('该连线没有关联任务，无法删除')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要删除连线 "${edge.label}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '删除',
        cancelButtonText: '取消',
        type: 'warning',
        confirmButtonClass: 'el-button--danger'
      }
    )

    ElMessage.info('正在删除...')

    // 调用API删除任务
    await progressApi.delete(edge.taskId)

    ElMessage.success('删除成功')

    // 触发刷新
    emit('task-updated', { deleted: true, id: edge.taskId })
  } catch (error) {
    if (error !== 'cancel') {
      throw error
    }
  }
}

// 点击其他地方关闭右键菜单
const handleClickOutside = () => {
  contextMenuVisible.value = false
}

// 通过连线创建任务
const createTaskFromConnection = async (startNode, endNode) => {
  try {
    // 计算任务的日期
    const startDate = new Date(startNode.earliest_start_ts * 1000 || startNode.ES * 86400 * 1000)
    const endDate = new Date(endNode.earliest_finish_ts * 1000 || endNode.ES * 86400 * 1000)

    // 计算工期（天数）
    const duration = Math.round((endDate.getTime() - startDate.getTime()) / (1000 * 60 * 60 * 24))

    const taskData = {
      project_id: props.projectId,
      name: `新任务 (${startNode.number}→${endNode.number})`,
      start_date: formatDate(startDate, 'YYYY-MM-DD'),
      end_date: formatDate(endDate, 'YYYY-MM-DD'),
      progress: 0,
      priority: 'medium',
      description: `从节点${startNode.number}到节点${endNode.number}的任务`
    }

    ElMessage.info('正在创建任务...')

    // 调用API创建任务
    await progressApi.create(taskData)

    ElMessage.success('任务创建成功')

    // 触发刷新事件，让父组件重新加载数据
    emit('task-updated', taskData)

    // 不需要调用render()，父组件接收到task-updated事件后会重新加载数据
    // watch(scheduleData) 会被触发并自动调用render()
  } catch (error) {
    console.error('创建任务失败:', error)
    ElMessage.error('创建任务失败：' + (error.response?.data?.message || error.message || '未知错误'))
  }
}

// 打开新建任务对话框（从连线创建）
const openCreateTaskDialog = (startNode, endNode) => {
  // 从 scheduleData 中查找对应的活动数据来获取正确的日期
  let startDate, endDate

  // 查找起始节点的活动
  if (startNode.taskId && props.scheduleData?.activities) {
    const activities = Object.values(props.scheduleData.activities)
    const startActivity = activities.find(a => a.task_id === startNode.taskId || a.id === startNode.taskId)
    if (startActivity?.earliest_start) {
      startDate = new Date(startActivity.earliest_start * 1000)
    }
  }

  // 查找结束节点的活动
  if (endNode.taskId && props.scheduleData?.activities) {
    const activities = Object.values(props.scheduleData.activities)
    const endActivity = activities.find(a => a.task_id === endNode.taskId || a.id === endNode.taskId)
    if (endActivity?.earliest_finish) {
      endDate = new Date(endActivity.earliest_finish * 1000)
    }
  }

  // 如果没有找到活动数据，使用节点的时间戳或ES值
  if (!startDate) {
    if (startNode.earliest_start_ts) {
      startDate = new Date(startNode.earliest_start_ts * 1000)
    } else if (startNode.ES && startNode.ES > 0) {
      startDate = new Date('1970-01-01')
      startDate.setTime(startDate.getTime() + startNode.ES * 86400 * 1000)
    } else {
      // 默认使用今天
      startDate = new Date()
    }
  }

  if (!endDate) {
    if (endNode.earliest_finish_ts) {
      endDate = new Date(endNode.earliest_finish_ts * 1000)
    } else if (endNode.EF && endNode.EF > 0) {
      endDate = new Date('1970-01-01')
      endDate.setTime(endDate.getTime() + endNode.EF * 86400 * 1000)
    } else {
      // 默认为开始日期+1天
      endDate = new Date(startDate.getTime() + 86400 * 1000)
    }
  }

  // 验证日期有效性：确保结束日期不早于开始日期
  if (endDate < startDate) {
    console.warn('Invalid date range: end date is before start date, adjusting...', {
      start: startDate.toISOString(),
      end: endDate.toISOString()
    })
    // 将结束日期设置为开始日期+1天
    endDate = new Date(startDate.getTime() + 86400 * 1000)
  }

  createTaskForm.value = {
    name: `新任务 (${startNode.number}→${endNode.number})`,
    isDummy: false,  // 默认创建实际任务
    start_date: formatDate(startDate, 'YYYY-MM-DD'),
    end_date: formatDate(endDate, 'YYYY-MM-DD'),
    progress: 0,
    priority: 'medium',
    description: `从节点${startNode.number}到节点${endNode.number}的任务`
  }

  // 保存连线信息，用于提交后显示连线
  createTaskForm.value._startNode = startNode
  createTaskForm.value._endNode = endNode

  createTaskDialogVisible.value = true
}

// 提交新建任务
const submitCreateTask = async () => {
  if (!createTaskFormRef.value) return

  try {
    await createTaskFormRef.value.validate()

    saving.value = true

    const taskData = {
      project_id: props.projectId,
      name: createTaskForm.value.name,
      start_date: createTaskForm.value.start_date,
      end_date: createTaskForm.value.end_date,
      progress: createTaskForm.value.progress,
      priority: createTaskForm.value.priority,
      description: createTaskForm.value.description,
      is_dummy: createTaskForm.value.isDummy  // 标记是否为虚任务
    }

    ElMessage.info('正在创建任务...')

    // 调用API创建任务
    const response = await progressApi.create(taskData)
    const newTaskId = response.data.id || response.data.task_id

    ElMessage.success('任务创建成功')

    // 创建依赖关系
    const startNode = createTaskForm.value._startNode
    const endNode = createTaskForm.value._endNode

    if (startNode && startNode.taskId) {
      try {
        await progressApi.createDependencyVisual(startNode.taskId, newTaskId, { type: 'FS', lag: 0 })
        console.log(`Created dependency: ${startNode.taskId} -> ${newTaskId}`)
      } catch (depError) {
        console.error('Failed to create dependency from start node:', depError)
        ElMessage.warning('任务创建成功，但前置依赖关系创建失败')
      }
    }

    if (endNode && endNode.taskId) {
      try {
        await progressApi.createDependencyVisual(newTaskId, endNode.taskId, { type: 'FS', lag: 0 })
        console.log(`Created dependency: ${newTaskId} -> ${endNode.taskId}`)
      } catch (depError) {
        console.error('Failed to create dependency to end node:', depError)
        ElMessage.warning('任务创建成功，但后置依赖关系创建失败')
      }
    }

    // 关闭对话框
    createTaskDialogVisible.value = false

    // 触发刷新事件，让父组件重新加载数据
    emit('task-updated', taskData)

    // 重置表单
    createTaskForm.value = {
      name: '',
      isDummy: false,
      start_date: null,
      end_date: null,
      progress: 0,
      priority: 'medium',
      description: ''
    }
  } catch (error) {
    if (error !== 'cancel') {  // 表单验证失败时会抛出'cancel'
      console.error('创建任务失败:', error)
      ElMessage.error('创建任务失败：' + (error.response?.data?.message || error.message || '未知错误'))
    }
  } finally {
    saving.value = false
  }
}

// 取消新建任务
const cancelCreateTask = () => {
  createTaskDialogVisible.value = false
  createTaskForm.value = {
    name: '',
    isDummy: false,
    start_date: null,
    end_date: null,
    progress: 0,
    priority: 'medium',
    description: ''
  }
}

// 节点双击 - 显示详情
const handleNodeDblClick = (node) => {
  selectedNode.value = node
  nodeDetailVisible.value = true
}

// 边点击 - 选中边
const handleEdgeClick = (edge) => {
  selectedEdge.value = edge
}

// 边双击 - 编辑任务详情（显示前置任务详情）
const handleEdgeDblClick = (edge) => {
  // 选中边并显示详情
  selectedEdge.value = edge
  edgeDetailVisible.value = true
}

// 更多操作
const handleMoreAction = (command) => {
  switch (command) {
    case 'fit-view':
      fitView()
      break
    case 'center-view':
      centerView()
      ElMessage.success('视图已居中')
      break
    case 'export':
      handleExportImage()
      break
    case 'toggle-fullscreen':
      toggleFullscreen()
      break
  }
}

// 导出图片
const handleExportImage = async () => {
  if (!svgRef.value) return

  try {
    const canvas = await html2canvas(svgRef.value.parentElement, {
      backgroundColor: '#ffffff',
      logging: false,
      useCORS: true,
      scale: 2
    })

    canvas.toBlob((blob) => {
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `网络图_${props.projectId}_${new Date().getTime()}.png`
      a.click()
      URL.revokeObjectURL(url)
      ElMessage.success('导出成功')
    })
  } catch (error) {
    console.error('导出失败:', error)
    ElMessage.error('导出失败，请重试')
  }
}

// 全屏切换
const toggleFullscreen = () => {
  isFullscreen.value = !isFullscreen.value

  if (isFullscreen.value) {
    document.body.classList.add('network-fullscreen')
  } else {
    document.body.classList.remove('network-fullscreen')
  }

  nextTick(() => {
    centerView()
  })
}

// 键盘事件处理
const handleKeydown = (event) => {
  // ESC键退出全屏或连线模式
  if (event.key === 'Escape') {
    if (isFullscreen.value) {
      isFullscreen.value = false
      document.body.classList.remove('network-fullscreen')
      nextTick(() => {
        centerView()
      })
    } else if (isConnecting.value) {
      // 退出连线模式
      isConnecting.value = false
      connectStartNode.value = null
      connectHoverNode.value = null
      connectEndPosition.value = { x: 0, y: 0 }
      // 切换回选择模式
      if (toolMode.value === 'connect') {
        toolMode.value = 'select'
      }
    }
  }

  // Ctrl键进入连线模式
  if (event.key === 'Control' && !isCtrlPressed.value) {
    isCtrlPressed.value = true
    // 如果已经在某个节点上，准备开始连线
    if (selectedNode.value && !isConnecting.value) {
      isConnecting.value = true
      connectStartNode.value = selectedNode.value
      connectEndPosition.value = { x: selectedNode.value.x, y: selectedNode.value.y }
    }
  }
}

// 键盘松开事件
const handleKeyup = (event) => {
  if (event.key === 'Control') {
    isCtrlPressed.value = false
    // 如果正在连线中，取消连线
    if (isConnecting.value) {
      isConnecting.value = false
      connectStartNode.value = null
      connectHoverNode.value = null
    }
  }
}

// 监听数据变化
watch(
  () => props.scheduleData,
  (newData, oldData) => {
    console.log('NetworkDiagram - scheduleData changed:', {
      hasNew: !!newData,
      hasOld: !!oldData,
      activitiesCount: newData?.activities ? Object.keys(newData.activities).length : 0,
      timestamp: Date.now()
    })
    render()
  },
  { deep: true, immediate: true }
)

// 监听dayWidth变化，重新计算节点位置和路径
watch(dayWidth, (newDayWidth, oldDayWidth) => {
  console.log('[watch dayWidth] Triggered:', { oldDayWidth, newDayWidth, nodesCount: nodes.value.length })

  if (nodes.value.length > 0 && props.scheduleData?.activities) {
    // 计算时间轴的起始日期
    const activities = Object.values(props.scheduleData.activities)
    const taskDates = activities.map(a => {
      const startDate = new Date(a.earliest_start * 1000)
      const endDate = new Date(a.earliest_finish * 1000)
      return { startDate, endDate }
    })

    const minDate = new Date(Math.min(...taskDates.map(d => d.startDate.getTime())))
    const bufferDays = 2
    const timelineStartDate = new Date(minDate)
    timelineStartDate.setDate(timelineStartDate.getDate() - bufferDays)

    const nodeRadiusVal = nodeRadius.value

    console.log('[watch dayWidth] Timeline start date:', timelineStartDate.toISOString().split('T')[0])
    console.log('[watch dayWidth] Before recalc - sample nodes:', nodes.value.slice(0, 2).map(n => ({ id: n.id, x: n.x, ES: n.ES })))

    // 重新计算节点位置（使用 dayWidth，与时间轴保持同步）
    // 关键修复：创建新数组以触发 Vue 的响应式更新
    const newNodes = nodes.value.map(node => {
      // 优先使用存储的时间戳（与 layoutNodes 保持一致）
      let nodeDate
      if (node.earliest_finish_ts) {
        nodeDate = new Date(node.earliest_finish_ts * 1000)
      } else if (node.earliest_start_ts) {
        nodeDate = new Date(node.earliest_start_ts * 1000)
      } else {
        // 回退到ES值转换
        nodeDate = new Date('1970-01-01')
        nodeDate.setTime(nodeDate.getTime() + node.ES * 86400 * 1000)
      }

      const daysDiff = Math.round((nodeDate.getTime() - timelineStartDate.getTime()) / (1000 * 60 * 60 * 24))
      const newX = Math.max(0, daysDiff * newDayWidth) // 使用 newDayWidth，与时间轴同步

      console.log(`[watch dayWidth] Node ${node.taskId} (${node.nodeType}): date=${nodeDate.toISOString().split('T')[0]}, daysDiff=${daysDiff}, oldX=${node.x}, newX=${newX}`)

      // 返回新对象以触发响应式更新
      return {
        ...node,
        x: newX
      }
    })

    nodes.value = newNodes

    // 重新计算边的坐标和正交路径
    const newEdges = edges.value.map(edge => {
      const fromNode = newNodes.find(n => n.number === edge.from)
      const toNode = newNodes.find(n => n.number === edge.to)
      if (fromNode && toNode) {
        const path = calculateOrthogonalPath(
          fromNode.x,
          fromNode.y,
          toNode.x,
          toNode.y,
          nodeRadiusVal,
          newNodes,
          fromNode.id,
          toNode.id
        )

        return {
          ...edge,
          x1: fromNode.x,
          y1: fromNode.y,
          x2: toNode.x,
          y2: toNode.y,
          path,
          labelX: (fromNode.x + toNode.x) / 2,
          labelY: (fromNode.y + toNode.y) / 2 - 8
        }
      }
      return edge
    })

    edges.value = newEdges

    console.log('[watch dayWidth] After recalc - sample nodes:', nodes.value.slice(0, 2).map(n => ({ id: n.id, x: n.x, ES: n.ES })))
    console.log('[watch dayWidth] Recalculation complete')
  }
})

onMounted(() => {
  render()
  centerView()
  // 添加全局点击监听，用于关闭右键菜单
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.body.classList.remove('network-fullscreen')
  // 移除全局点击监听
  document.removeEventListener('click', handleClickOutside)
})

// 添加任务
const handleAddTask = () => {
  editingTask.value = null
  taskForm.value = {
    name: '',
    start: new Date(),
    end: new Date(),
    progress: 0,
    priority: 'medium',
    notes: ''
  }
  editDialogVisible.value = true
}

// 保存任务
const handleSaveTask = async () => {
  if (!taskFormRef.value) return

  try {
    await taskFormRef.value.validate()

    saving.value = true

    const taskData = {
      project_id: props.projectId,
      name: taskForm.value.name,
      start_date: formatDate(taskForm.value.start),
      end_date: formatDate(taskForm.value.end),
      progress: taskForm.value.progress,
      priority: taskForm.value.priority,
      description: taskForm.value.notes
    }

    if (editingTask.value) {
      // 更新现有任务
      await progressApi.update(editingTask.value.id, taskData)
      ElMessage.success('任务更新成功')
      emit('task-updated', { ...editingTask.value, ...taskData })
    } else {
      // 创建新任务
      await progressApi.create(taskData)
      ElMessage.success('任务创建成功')
      emit('task-updated', taskData)
    }

    editDialogVisible.value = false
    editingTask.value = null
  } catch (error) {
    if (error !== false) {
      console.error('保存任务失败:', error)
      ElMessage.error('保存任务失败')
    }
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.network-diagram {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #fff;
  outline: none;
}

/* 平移模式下的光标样式 */
.network-diagram.pan-mode .network-svg {
  cursor: grab;
}

.network-diagram.pan-mode.is-panning .network-svg {
  cursor: grabbing;
}

/* 连线模式下的光标样式 */
.network-diagram.connect-mode .network-svg {
  cursor: crosshair;
}

.network-diagram:not(.pan-mode):not(.connect-mode) .network-svg {
  cursor: default;
}

/* 连线模式下节点的样式 */
.node-group.node-connecting {
  filter: drop-shadow(0 0 8px rgba(100, 181, 246, 0.8));
}

.node-group.node-connecting circle {
  stroke: #64B5F6;
  stroke-width: 3;
}

/* 连线时目标节点的hover效果 */
.node-group.node-connect-hover circle {
  stroke: #5AB9A6;
  stroke-width: 4;
  filter: drop-shadow(0 0 12px rgba(90, 185, 166, 0.8));
  transform-box: fill-box;
  transform-origin: center;
  animation: nodePulse 0.5s ease-in-out infinite alternate;
}

/* 拖动节点到目标节点上时的效果 */
.node-group.node-drag-hover circle {
  stroke: #FF9800;
  stroke-width: 5;
  filter: drop-shadow(0 0 15px rgba(255, 152, 0, 0.9));
  transform-box: fill-box;
  transform-origin: center;
  animation: nodeMergePulse 0.4s ease-in-out infinite alternate;
}

@keyframes nodeMergePulse {
  from {
    transform: scale(1);
  }
  to {
    transform: scale(1.2);
  }
}

@keyframes nodePulse {
  from {
    transform: scale(1);
  }
  to {
    transform: scale(1.15);
  }
}

.network-diagram.fullscreen {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 9999;
  background: #fff;
}

.network-toolbar {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  background: #fff;
  border-bottom: 1px solid #dcdfe6;
  flex-wrap: wrap;
  gap: 8px;
}

.network-stats {
  display: flex;
  align-items: center;
  gap: 24px;
  padding: 12px 16px;
  background: #f5f7fa;
  border-bottom: 1px solid #dcdfe6;
  flex-wrap: wrap;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
}

.stat-label {
  color: #909399;
}

.stat-value {
  font-weight: bold;
  color: #303133;
}

.stat-value.critical {
  color: #f56c6c;
}

.network-canvas-container {
  flex: 1;
  position: relative;
  overflow: hidden;
  background: #fafafa;
}

.network-svg {
  display: block;
  background: white;
}

.network-loading,
.network-empty {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  color: #909399;
}

.node-group {
  transition: all 0.3s;
}

.node-group:hover {
  filter: brightness(1.1);
}

.node-group.node-selected circle {
  stroke: #64B5F6;
  stroke-width: 4;
}

.node-start circle {
  stroke: #5AB9A6;
  stroke-width: 3;
}

.node-end circle {
  stroke: #E57373;
  stroke-width: 3;
}

.node-shared circle {
  stroke: #64B5F6;
  stroke-width: 3;
}

.node-critical circle {
  stroke: #FF8A65;
  stroke-width: 3;
}

.node-start.node-critical circle {
  stroke: #5AB9A6;
  stroke-width: 3;
}

.node-end.node-critical circle {
  stroke: #E57373;
  stroke-width: 3;
}

.node-shared.node-critical circle {
  stroke: #64B5F6;
  stroke-width: 3;
}

/* 聚焦模式样式 */
.node-dimmed {
  transition: opacity 0.3s ease;
}

.edge-dimmed {
  transition: opacity 0.3s ease;
}

/* 关键路径高亮效果 */
.node-critical {
  filter: drop-shadow(0 0 6px rgba(255, 138, 101, 0.4));
}

.edge-critical {
  filter: drop-shadow(0 0 4px rgba(255, 138, 101, 0.3));
}

.time-param-text {
  font-weight: 500;
}

.slack-param-text {
  font-weight: 500;
  opacity: 0.85;
}

.edge-critical {
  stroke-width: 1.5 !important;
}

.edge-selected {
  stroke: #409eff !important;
  stroke-width: 1.5 !important;
  stroke-dasharray: none;
}

.edge-task {
  transition: all 0.3s;
}

.edge-task:hover {
  stroke-width: 1.5 !important;
  filter: brightness(1.1);
}

.edge-dependency {
  opacity: 0.6;
}

.edge-dependency:hover {
  opacity: 1;
}

.edge-dummy {
  opacity: 0.7;
  transition: all 0.3s;
}

.edge-dummy:hover {
  opacity: 1;
  stroke-width: 1.5 !important;
}

.edge-label {
  pointer-events: none;
  font-weight: 500;
}

.edge-label-task {
  font-size: 13px;
  font-weight: bold;
}

.edge-slack-label {
  pointer-events: none;
  font-weight: 500;
  opacity: 0.85;
}

/* 时间轴样式（两行显示：上月份下日期） */
.network-timeline-header {
  position: relative;
  height: 60px;
  background: #f5f7fa;
  border-bottom: 1px solid #dcdfe6;
  flex-shrink: 0;
}

.timeline-content {
  position: relative;
  height: 100%;
  min-width: 800px;
  will-change: transform;
  /* 平滑过渡，避免拖拽时延迟 */
  transition: transform 0.1s ease-out;
}

/* 拖拽时禁用过渡以获得即时响应 */
.network-diagram.is-panning .timeline-content {
  transition: none;
}

/* 月份行（上层） */
.timeline-months-row {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 28px;
  border-bottom: 1px solid #dcdfe6;
}

.timeline-month-cell {
  position: absolute;
  top: 0;
  height: 100%;
  border-right: 1px solid #e4e7ed;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: bold;
  color: #303133;
  background: #fff;
}

.month-label {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  padding: 0 4px;
}

/* 日期行（下层） */
.timeline-days-row {
  position: absolute;
  top: 28px;
  left: 0;
  right: 0;
  bottom: 0;
}

.timeline-day-cell {
  position: absolute;
  top: 0;
  height: 100%;
  border-right: 1px solid #e4e7ed;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  color: #606266;
  transition: background 0.2s;
}

.timeline-day-cell.is-today {
  background: #fff3e0;
}

.timeline-day-cell.is-weekend {
  background: #fafafa;
}

/* 只显示奇数日期时，偶数日期淡化 */
.timeline-day-cell:not(.is-odd) {
  opacity: 0.3;
}

.day-number {
  font-size: 13px;
  font-weight: 500;
  color: #303133;
}

.timeline-day-cell.is-today .day-number {
  color: #e6a23c;
  font-weight: bold;
}

.network-timeline-header .today-line {
  position: absolute;
  top: 0;
  bottom: 0;
  width: 2px;
  background: #FF8A65;
  z-index: 5;
}

.network-timeline-header .today-line::before {
  content: '今天';
  position: absolute;
  top: -20px;
  left: 50%;
  transform: translateX(-50%);
  font-size: 11px;
  color: #FF8A65;
  white-space: nowrap;
}

.network-legend {
  display: flex;
  align-items: center;
  gap: 24px;
  padding: 12px 16px;
  background: #f5f7fa;
  border-top: 1px solid #dcdfe6;
  flex-wrap: wrap;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #606266;
}

.legend-color {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  border: 2px solid;
}

.legend-color.start {
  background: #5AB9A6;
  border-color: #5AB9A6;
}

.legend-color.shared {
  background: #64B5F6;
  border-color: #64B5F6;
}

.legend-color.end {
  background: #E57373;
  border-color: #E57373;
}

.legend-color.critical {
  background: #FF8A65;
  border-color: #FF8A65;
}

.legend-edge {
  width: 30px;
  height: 3px;
  border-radius: 2px;
}

.legend-edge.task-edge {
  background: #64B5F6;
  height: 2px;
}

.legend-edge.dummy-edge {
  background: #B0BEC5;
  height: 2px;
  border-top: 2px dashed #B0BEC5;
  background: none;
}

.legend-edge.dep-edge {
  background: #CFD8DC;
  height: 1px;
  border-top: 1px dashed #CFD8DC;
  background: none;
}

/* 节点详情 */
.node-detail,
.edge-detail {
  padding: 0 20px;
}

/* 全屏样式 */
:deep(.network-fullscreen) {
  overflow: hidden;
}

/* 右键菜单样式 */
.context-menu {
  position: fixed;
  z-index: 10000;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
  min-width: 200px;
  max-width: 280px;
  overflow: hidden;
  animation: contextMenuFadeIn 0.15s ease-out;
}

@keyframes contextMenuFadeIn {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.context-menu-header {
  padding: 12px 16px 8px;
  background: #f5f7fa;
}

.menu-title {
  display: block;
  font-size: 14px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 2px;
}

.menu-subtitle {
  display: block;
  font-size: 12px;
  color: #909399;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.context-menu-items {
  padding: 4px 0;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 16px;
  font-size: 14px;
  color: #606266;
  cursor: pointer;
  transition: all 0.2s;
  user-select: none;
}

.menu-item:hover {
  background: #f5f7fa;
  color: #409eff;
}

.menu-item.danger {
  color: #f56c6c;
}

.menu-item.danger:hover {
  background: #fef0f0;
  color: #f56c6c;
}

.menu-item .el-icon {
  font-size: 16px;
}


/* 打印样式 */
@media print {
  .network-toolbar,
  .network-stats,
  .network-legend {
    display: none !important;
  }

  .network-diagram {
    position: static !important;
    height: auto !important;
  }

  .network-canvas-container {
    overflow: visible !important;
  }
}
</style>
