# Gantt Chart Editor - Main Integration Guide

**Version:** 1.0.0
**Last Updated:** 2026-02-19
**Author:** Material Management System Team

---

## Table of Contents

1. [Overview](#overview)
2. [Prerequisites](#prerequisites)
3. [Frontend Integration](#frontend-integration)
4. [Backend Integration](#backend-integration)
5. [Component Registration](#component-registration)
6. [Configuration](#configuration)
7. [Testing](#testing)
8. [Deployment](#deployment)

---

## Overview

This guide provides step-by-step instructions for integrating the complete Gantt Chart Editor system into your Material Management application. The editor is built across 5 phases, providing advanced project management capabilities including:

- **Phase 1:** Core Editor (Virtual timeline, task list, dependencies)
- **Phase 2:** Views & Modes (Kanban, Calendar, Dashboard, Mobile)
- **Phase 3:** Panels & Dialogs (Comments, History, Templates, Bulk Edit)
- **Phase 4:** Advanced Features (Resource leveling, Critical path, Reports)
- **Phase 5:** UX & Automation (Guided tour, AI suggestions, Workflow automation)

---

## Prerequisites

### System Requirements

- **Node.js:** v18.0.0 or higher
- **Vue:** 3.4.0 or higher
- **Element Plus:** 2.5.0 or higher
- **Pinia:** 2.1.0 or higher
- **Browser:** Chrome 90+, Firefox 88+, Safari 14+, Edge 90+

### Required Dependencies

All dependencies are already installed in `package.json`:

```json
{
  "dependencies": {
    "@vueuse/core": "^11.0.0",
    "date-fns": "^3.0.0",
    "element-plus": "^2.5.0",
    "lodash-es": "^4.17.21",
    "pinia": "^2.1.0",
    "socket.io-client": "^4.8.3",
    "vue-tour": "^2.0.0",
    "vue-virtual-scroller": "^2.0.0-beta.8",
    "chart.js": "^4.4.0",
    "vue-chartjs": "^5.3.0"
  }
}
```

---

## Frontend Integration

### Step 1: Update main.js

**Location:** `/home/julei/backend/newstatic/src/main.js`

The following registrations are already complete. Verify your main.js includes:

```javascript
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

// Virtual Scroller (CRITICAL for performance)
import 'vue-virtual-scroller/dist/vue-virtual-scroller.css'
import { RecycleScroller } from 'vue-virtual-scroller'

// Guided Tour
import VueTour from 'vue-tour'
import 'vue-tour/dist/vue-tour.css'

import App from './App.vue'
import router from './router'
import '@/assets/css/main.css'

// Theme system
import '@/styles/themes/variables.css'
import '@/styles/themes/dark.css'

const app = createApp(App)

// Register all Element Plus icons
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

// Register vue-virtual-scroller component (CRITICAL)
app.component('RecycleScroller', RecycleScroller)

// Register Vue Tour
app.use(VueTour)

app.use(createPinia())
app.use(router)
app.use(ElementPlus)

app.mount('#app')
```

### Step 2: Create Component Barrel Exports

**Location:** `/home/julei/backend/newstatic/src/components/gantt/index.js`

Create a centralized export file for all Gantt components:

```javascript
/**
 * Gantt Chart Components - Complete Export
 * All components organized by phase and category
 */

// ============ CORE COMPONENTS ============
export { default as GanttEditor } from './core/GanttEditor.vue'
export { default as GanttToolbar } from './core/GanttToolbar.vue'
export { default as GanttStatusBar } from './core/GanttStatusBar.vue'
export { default as GanttBody } from './core/GanttBody.vue'

// ============ TIMELINE ============
export { default as VirtualTimeline } from './timeline/VirtualTimeline.vue'
export { default as TimelineBackground } from './timeline/TimelineBackground.vue'
export { default as TimelineGrid } from './timeline/TimelineGrid.vue'
export { default as TaskBar } from './timeline/TaskBar.vue'
export { default as TaskBarSvg } from './timeline/TaskBarSvg.vue'
export { default as DependencyLines } from './timeline/DependencyLines.vue'

// ============ TABLE ============
export { default as VirtualTaskList } from './table/VirtualTaskList.vue'
export { default as TaskRow } from './table/TaskRow.vue'
export { default as EditableCell } from './table/EditableCell.vue'

// ============ VIRTUAL ============
export { default as VirtualTimelineWrapper } from './virtual/VirtualTimeline.vue'

// ============ VIEWS ============
export { default as DashboardView } from './views/DashboardView.vue'
export { default as KanbanView } from './views/KanbanView.vue'
export { default as KanbanColumn } from './views/KanbanColumn.vue'
export { default as KanbanCard } from './views/KanbanCard.vue'
export { default as CalendarView } from './views/CalendarView.vue'
export { default as CalendarPreview } from './views/CalendarPreview.vue'
export { default as ResourceHistogram } from './views/ResourceHistogram.vue'
export { default as GanttMiniView } from './views/GanttMiniView.vue'

// ============ MOBILE ============
export { default as MobileGanttView } from './mobile/MobileGanttView.vue'
export { default as TaskListView } from './mobile/TaskListView.vue'
export { default as MobileTaskRow } from './mobile/TaskRow.vue'
export { default as TimelineSwipeView } from './mobile/TimelineSwipeView.vue'

// ============ DASHBOARD ============
export { default as BurndownChart } from './dashboard/BurndownChart.vue'
export { default as EarnedValueChart } from './dashboard/EarnedValueChart.vue'
export { default as MilestoneTracker } from './dashboard/MilestoneTracker.vue'
export { default as ResourceUtilization } from './dashboard/ResourceUtilization.vue'
export { default as StatCard } from './dashboard/StatCard.vue'

// ============ PANELS ============
export { default as CommentsPanel } from './panels/CommentsPanel.vue'
export { default as CommentItem } from './panels/CommentItem.vue'
export { default as HistoryPanel } from './panels/HistoryPanel.vue'
export { default as HistoryItem } from './panels/HistoryItem.vue'
export { default as DiffViewer } from './panels/DiffViewer.vue'
export { default as DiffField } from './panels/DiffField.vue'
export { default as SmartSuggestionsPanel } from './panels/SmartSuggestionsPanel.vue'
export { default as SuggestionCard } from './panels/SuggestionCard.vue'

// ============ DIALOGS ============
export { default as BulkEditDialog } from './dialogs/BulkEditDialog.vue'
export { default as CalendarDialog } from './dialogs/CalendarDialog.vue'
export { default as ConstraintEditDialog } from './dialogs/ConstraintEditDialog.vue'
export { default as ReportBuilderDialog } from './dialogs/ReportBuilderDialog.vue'
export { default as ResourceLevelingDialog } from './dialogs/ResourceLevelingDialog.vue'
export { default as TaskTemplatesDialog } from './dialogs/TaskTemplatesDialog.vue'
export { default as TemplateManagerDialog } from './dialogs/TemplateManagerDialog.vue'

// ============ OVERLAYS ============
export { default as ContextMenu } from './overlays/ContextMenu.vue'
export { default as GuidedTour } from './overlays/GuidedTour.vue'
export { default as Minimap } from './overlays/Minimap.vue'

// ============ STORES ============
export { default as useGanttStore } from '@/stores/ganttStore'
export { default as useUndoRedoStore } from '@/stores/undoRedoStore'
export { default as useCalendarStore } from '@/stores/calendarStore'
export { default as useCollaborationStore } from '@/stores/collaborationStore'
export { default as useTemplateStore } from '@/stores/templateStore'

// ============ UTILITIES ============
export {
  getTourSteps,
  tourConfig,
  tourTypes,
  getCompletedTours,
  markTourCompleted,
  isTourCompleted,
  resetTourProgress
} from '@/utils/tourSteps'

export {
  WorkflowAutomationEngine,
  getWorkflowEngine,
  executeAutomation,
  RuleTypes,
  TriggerTypes,
  defaultRules
} from '@/utils/workflowAutomation'

export {
  analyzeSchedule,
  predictDelays,
  optimizeSchedule,
  createSuggestion
} from '@/utils/aiOptimizer'

export {
  formatDate,
  parseDate,
  addDays,
  subDays,
  differenceInDays,
  formatDuration
} from '@/utils/dateUtils'

export {
  exportToCSV,
  exportToJSON,
  exportToPDF,
  exportToICS
} from '@/utils/exportUtils'

export {
  calculateCriticalPath,
  calculateFloat,
  calculateES,
  calculateEF,
  calculateLS,
  calculateLF,
  isTaskOnCriticalPath
} from '@/utils/criticalPath'

// ============ CONSTANTS ============
export {
  ViewMode,
  TaskStatus,
  DependencyType,
  ConstraintType,
  Priority,
  ZoomLevel,
  ExportFormat
} from '@/constants/ganttConstants'
```

### Step 3: Update Router Configuration

**Location:** `/home/julei/backend/newstatic/src/router/index.js`

Add routes for Gantt views. Insert these routes in the children array:

```javascript
{
  path: 'progress/gantt',
  name: 'ProgressGantt',
  component: () => import('@/views/ProgressGantt.vue'),
  meta: {
    requiresAuth: true,
    title: '甘特图编辑器',
    permissions: ['progress_view', 'progress_edit']
  }
},
{
  path: 'progress/kanban',
  name: 'ProgressKanban',
  component: () => import('@/views/ProgressKanban.vue'),
  meta: {
    requiresAuth: true,
    title: '看板视图',
    permissions: ['progress_view']
  }
},
{
  path: 'progress/calendar',
  name: 'ProgressCalendar',
  component: () => import('@/views/ProgressCalendar.vue'),
  meta: {
    requiresAuth: true,
    title: '日历视图',
    permissions: ['progress_view']
  }
},
{
  path: 'progress/dashboard',
  name: 'ProgressDashboard',
  component: () => import('@/views/ProgressDashboard.vue'),
  meta: {
    requiresAuth: true,
    title: '进度仪表板',
    permissions: ['progress_view']
  }
}
```

### Step 4: Initialize Stores in App Setup

**Location:** `/home/julei/backend/newstatic/src/App.vue`

Ensure stores are initialized when the app starts:

```vue
<script setup>
import { onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useGanttStore } from '@/stores/ganttStore'
import { useUndoRedoStore } from '@/stores/undoRedoStore'
import { useCalendarStore } from '@/stores/calendarStore'
import { useCollaborationStore } from '@/stores/collaborationStore'
import { useTemplateStore } from '@/stores/templateStore'

const authStore = useAuthStore()
const ganttStore = useGanttStore()
const undoRedoStore = useUndoRedoStore()
const calendarStore = useCalendarStore()
const collaborationStore = useCollaborationStore()
const templateStore = useTemplateStore()

onMounted(async () => {
  // Initialize auth store
  await authStore.checkAuth()

  // Initialize Gantt stores if user is authenticated
  if (authStore.isAuthenticated) {
    await ganttStore.initialize()
    await calendarStore.loadCalendars()
    await templateStore.loadTemplates()

    // Initialize collaboration store
    collaborationStore.connect()
  }
})
</script>
```

### Step 5: Update API Client

**Location:** `/home/julei/backend/newstatic/src/api/gantt.js`

Create API client for Gantt endpoints:

```javascript
import axios from 'axios'

const BASE_URL = '/api/gantt'

export const ganttApi = {
  // ============ TASKS ============
  async getTasks(projectId) {
    const { data } = await axios.get(`${BASE_URL}/projects/${projectId}/tasks`)
    return data
  },

  async createTask(projectId, task) {
    const { data } = await axios.post(`${BASE_URL}/projects/${projectId}/tasks`, task)
    return data
  },

  async updateTask(projectId, taskId, updates) {
    const { data } = await axios.patch(
      `${BASE_URL}/projects/${projectId}/tasks/${taskId}`,
      updates
    )
    return data
  },

  async deleteTask(projectId, taskId) {
    await axios.delete(`${BASE_URL}/projects/${projectId}/tasks/${taskId}`)
  },

  async bulkUpdateTasks(projectId, updates) {
    const { data } = await axios.patch(
      `${BASE_URL}/projects/${projectId}/tasks/bulk`,
      { updates }
    )
    return data
  },

  // ============ DEPENDENCIES ============
  async getDependencies(projectId) {
    const { data } = await axios.get(`${BASE_URL}/projects/${projectId}/dependencies`)
    return data
  },

  async createDependency(projectId, dependency) {
    const { data } = await axios.post(
      `${BASE_URL}/projects/${projectId}/dependencies`,
      dependency
    )
    return data
  },

  async deleteDependency(projectId, dependencyId) {
    await axios.delete(`${BASE_URL}/projects/${projectId}/dependencies/${dependencyId}`)
  },

  // ============ CONSTRAINTS ============
  async getConstraints(projectId) {
    const { data } = await axios.get(`${BASE_URL}/projects/${projectId}/constraints`)
    return data
  },

  async updateConstraint(projectId, taskId, constraint) {
    const { data } = await axios.put(
      `${BASE_URL}/projects/${projectId}/tasks/${taskId}/constraint`,
      constraint
    )
    return data
  },

  // ============ COMMENTS ============
  async getComments(projectId, taskId) {
    const { data } = await axios.get(
      `${BASE_URL}/projects/${projectId}/tasks/${taskId}/comments`
    )
    return data
  },

  async addComment(projectId, taskId, comment) {
    const { data } = await axios.post(
      `${BASE_URL}/projects/${projectId}/tasks/${taskId}/comments`,
      comment
    )
    return data
  },

  // ============ HISTORY ============
  async getHistory(projectId, taskId) {
    const { data } = await axios.get(
      `${BASE_URL}/projects/${projectId}/tasks/${taskId}/history`
    )
    return data
  },

  // ============ TEMPLATES ============
  async getTemplates() {
    const { data } = await axios.get(`${BASE_URL}/templates`)
    return data
  },

  async createTemplate(template) {
    const { data } = await axios.post(`${BASE_URL}/templates`, template)
    return data
  },

  async applyTemplate(projectId, templateId) {
    const { data } = await axios.post(
      `${BASE_URL}/projects/${projectId}/apply-template/${templateId}`
    )
    return data
  },

  // ============ REPORTS ============
  async generateReport(projectId, config) {
    const { data } = await axios.post(
      `${BASE_URL}/projects/${projectId}/reports`,
      config
    )
    return data
  },

  // ============ RESOURCE LEVELING ============
  async levelResources(projectId, options) {
    const { data } = await axios.post(
      `${BASE_URL}/projects/${projectId}/level-resources`,
      options
    )
    return data
  },

  // ============ AI SUGGESTIONS ============
  async getAISuggestions(projectId) {
    const { data } = await axios.get(`${BASE_URL}/projects/${projectId}/ai-suggestions`)
    return data
  },

  async acceptSuggestion(projectId, suggestionId) {
    const { data } = await axios.post(
      `${BASE_URL}/projects/${projectId}/ai-suggestions/${suggestionId}/accept`
    )
    return data
  }
}
```

### Step 6: Add i18n Translations

**Location:** `/home/julei/backend/newstatic/src/locales/en.js` and `/home/julei/backend/newstatic/src/locales/zh.js`

Add translations for Gantt features. Example for English:

```javascript
export default {
  gantt: {
    title: 'Gantt Chart Editor',
    toolbar: {
      undo: 'Undo',
      redo: 'Redo',
      zoomIn: 'Zoom In',
      zoomOut: 'Zoom Out',
      fitToScreen: 'Fit to Screen',
      export: 'Export',
      import: 'Import',
      templates: 'Templates',
      bulkEdit: 'Bulk Edit',
      fullscreen: 'Fullscreen',
      settings: 'Settings'
    },
    views: {
      gantt: 'Gantt View',
      kanban: 'Kanban View',
      calendar: 'Calendar View',
      dashboard: 'Dashboard'
    },
    tasks: {
      add: 'Add Task',
      edit: 'Edit Task',
      delete: 'Delete Task',
      indent: 'Indent Task',
      outdent: 'Outdent Task',
      moveUp: 'Move Up',
      moveDown: 'Move Down'
    },
    // ... more translations
  }
}
```

---

## Backend Integration

### Step 1: Add New API Routes

**Location:** `/home/julei/backend/cmd/server/main.go`

Register the Gantt API routes:

```go
package main

import (
    "github.com/gin-gonic/gin"
    "your-project/internal/api/gantt"
)

func main() {
    r := gin.Default()

    // API v1 group
    v1 := r.Group("/api/v1")
    {
        // Gantt routes
        ganttGroup := v1.Group("/gantt")
        {
            ganttGroup.GET("/projects/:projectId/tasks", ganttHandler.GetTasks)
            ganttGroup.POST("/projects/:projectId/tasks", ganttHandler.CreateTask)
            ganttGroup.PATCH("/projects/:projectId/tasks/:taskId", ganttHandler.UpdateTask)
            ganttGroup.DELETE("/projects/:projectId/tasks/:taskId", ganttHandler.DeleteTask)
            ganttGroup.PATCH("/projects/:projectId/tasks/bulk", ganttHandler.BulkUpdateTasks)

            ganttGroup.GET("/projects/:projectId/dependencies", ganttHandler.GetDependencies)
            ganttGroup.POST("/projects/:projectId/dependencies", ganttHandler.CreateDependency)
            ganttGroup.DELETE("/projects/:projectId/dependencies/:dependencyId", ganttHandler.DeleteDependency)

            ganttGroup.GET("/projects/:projectId/constraints", ganttHandler.GetConstraints)
            ganttGroup.PUT("/projects/:projectId/tasks/:taskId/constraint", ganttHandler.UpdateConstraint)

            ganttGroup.GET("/projects/:projectId/tasks/:taskId/comments", ganttHandler.GetComments)
            ganttGroup.POST("/projects/:projectId/tasks/:taskId/comments", ganttHandler.AddComment)

            ganttGroup.GET("/projects/:projectId/tasks/:taskId/history", ganttHandler.GetHistory)

            ganttGroup.GET("/templates", ganttHandler.GetTemplates)
            ganttGroup.POST("/templates", ganttHandler.CreateTemplate)
            ganttGroup.POST("/projects/:projectId/apply-template/:templateId", ganttHandler.ApplyTemplate)

            ganttGroup.POST("/projects/:projectId/reports", ganttHandler.GenerateReport)
            ganttGroup.POST("/projects/:projectId/level-resources", ganttHandler.LevelResources)
            ganttGroup.GET("/projects/:projectId/ai-suggestions", ganttHandler.GetAISuggestions)
            ganttGroup.POST("/projects/:projectId/ai-suggestions/:suggestionId/accept", ganttHandler.AcceptSuggestion)
        }

        // WebSocket endpoint for real-time collaboration
        v1.GET("/ws/gantt/:projectId", ganttHandler.HandleWebSocket)
    }

    r.Run(":8080")
}
```

### Step 2: Run Database Migrations

Execute the SQL migration file to create required tables:

```bash
cd /home/julei/backend
psql -U your_username -d your_database -f docs/DATABASE_MIGRATIONS.sql
```

### Step 3: Register WebSocket Hub

**Location:** `/home/julei/backend/internal/api/gantt/websocket.go`

Create WebSocket hub for real-time collaboration:

```go
package gantt

import (
    "encoding/json"
    "log"
    "sync"

    "github.com/gorilla/websocket"
)

type WebSocketHub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
    mutex      sync.RWMutex
}

type Client struct {
    hub      *WebSocketHub
    conn     *websocket.Conn
    send     chan []byte
    projectID string
    userID    string
}

type WSMessage struct {
    Type    string      `json:"type"`
    Data    interface{} `json:"data"`
    UserID  string      `json:"userId"`
    Time    int64       `json:"time"`
}

var Hub = &WebSocketHub{
    clients:    make(map[*Client]bool),
    broadcast:  make(chan []byte),
    register:   make(chan *Client),
    unregister: make(chan *Client),
}

func (h *WebSocketHub) Run() {
    for {
        select {
        case client := <-h.register:
            h.mutex.Lock()
            h.clients[client] = true
            h.mutex.Unlock()
            log.Printf("Client connected to project %s", client.projectID)

        case client := <-h.unregister:
            if _, ok := h.clients[client]; ok {
                h.mutex.Lock()
                delete(h.clients, client)
                close(client.send)
                h.mutex.Unlock()
                log.Printf("Client disconnected from project %s", client.projectID)
            }

        case message := <-h.broadcast:
            h.mutex.RLock()
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
            h.mutex.RUnlock()
        }
    }
}

func (h *WebSocketHub) BroadcastToProject(projectID string, messageType string, data interface{}) {
    message := WSMessage{
        Type: messageType,
        Data: data,
        Time: time.Now().Unix(),
    }

    jsonData, err := json.Marshal(message)
    if err != nil {
        log.Printf("Error marshaling WebSocket message: %v", err)
        return
    }

    h.mutex.RLock()
    for client := range h.clients {
        if client.projectID == projectID {
            select {
            case client.send <- jsonData:
            default:
                close(client.send)
                delete(h.clients, client)
            }
        }
    }
    h.mutex.RUnlock()
}
```

### Step 4: Add CORS Configuration for WebSocket

**Location:** `/home/julei/backend/cmd/server/main.go`

Configure CORS to allow WebSocket connections:

```go
import (
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    // CORS configuration
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

    // ... rest of setup
}
```

### Step 5: Update Middleware

**Location:** `/home/julei/backend/internal/middleware/auth.go`

Add authentication for new Gantt endpoints:

```go
package middleware

import (
    "github.com/gin-gonic/gin"
    "your-project/internal/auth"
)

func GanttAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Check if user has required permissions
        user := auth.GetCurrentUser(c)
        if user == nil {
            c.JSON(401, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }

        // Check specific Gantt permissions
        projectID := c.Param("projectId")
        if !auth.CanAccessProject(user, projectID) {
            c.JSON(403, gin.H{"error": "Forbidden"})
            c.Abort()
            return
        }

        c.Next()
    }
}
```

---

## Component Registration

### Global Registration

**Location:** `/home/julei/backend/newstatic/src/main.js`

All critical components are already registered globally:

```javascript
// Element Plus icons (all registered)
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

// Virtual Scroller (CRITICAL for performance)
app.component('RecycleScroller', RecycleScroller)

// Vue Tour
app.use(VueTour)
```

### Local Registration

For views that use Gantt components, import them locally:

```vue
<script setup>
import { GanttEditor, KanbanView, CalendarView } from '@/components/gantt'
import { useGanttStore, useUndoRedoStore } from '@/components/gantt'
import { ganttApi } from '@/api/gantt'

const ganttStore = useGanttStore()
const undoRedoStore = useUndoRedoStore()

// Load data
const loadGanttData = async () => {
  const tasks = await ganttApi.getTasks(props.projectId)
  ganttStore.setTasks(tasks)
}
</script>
```

### Add Global Directives

**Location:** `/home/julei/backend/newstatic/src/directives/index.js`

Create utility directives:

```javascript
import { createApp } from 'vue'

// Click outside directive
export const clickOutside = {
  beforeMount(el, binding) {
    el.clickOutsideEvent = function(event) {
      if (!(el == event.target || el.contains(event.target))) {
        binding.value(event)
      }
    }
    document.body.addEventListener('click', el.clickOutsideEvent)
  },
  unmounted(el) {
    document.body.removeEventListener('click', el.clickOutsideEvent)
  }
}

// Long press directive
export const longPress = {
  beforeMount(el, binding) {
    let timer
    const duration = binding.value.duration || 500

    el.start = () => {
      timer = setTimeout(() => {
        binding.value.handler()
      }, duration)
    }

    el.cancel = () => {
      clearTimeout(timer)
    }

    el.addEventListener('mousedown', el.start)
    el.addEventListener('touchstart', el.start)
    el.addEventListener('mouseup', el.cancel)
    el.addEventListener('mouseleave', el.cancel)
    el.addEventListener('touchend', el.cancel)
  },
  unmounted(el) {
    el.removeEventListener('mousedown', el.start)
    el.removeEventListener('touchstart', el.start)
    el.removeEventListener('mouseup', el.cancel)
    el.removeEventListener('mouseleave', el.cancel)
    el.removeEventListener('touchend', el.cancel)
  }
}

// Register directives
export default function registerDirectives(app) {
  app.directive('click-outside', clickOutside)
  app.directive('long-press', longPress)
}
```

Register in main.js:

```javascript
import registerDirectives from './directives'
registerDirectives(app)
```

---

## Configuration

### Virtual Scroller Configuration

**Location:** `/home/julei/backend/newstatic/src/components/gantt/config/virtualScroller.js`

```javascript
export const virtualScrollerConfig = {
  // Buffer size (number of items to render outside viewport)
  buffer: 200,

  // Item size
  itemSize: 50, // Row height in pixels

  // Threshold for rendering more items
  threshold: 100,

  // Enable recycling
  recycle: true,

  // Key field for efficient updates
  keyField: 'id'
}
```

### Timeline Configuration

```javascript
export const timelineConfig = {
  // Day width in pixels
  minDayWidth: 20,
  maxDayWidth: 100,
  defaultDayWidth: 40,

  // Row height
  rowHeight: 50,

  // Timeline start buffer (days before first task)
  startBuffer: 7,

  // Timeline end buffer (days after last task)
  endBuffer: 30
}
```

### WebSocket Configuration

```javascript
export const websocketConfig = {
  // WebSocket server URL
  url: process.env.VUE_APP_WS_URL || 'ws://localhost:8080/api/v1/ws/gantt',

  // Reconnection settings
  reconnect: true,
  reconnectInterval: 3000,
  maxReconnectAttempts: 10,

  // Heartbeat
  heartbeatInterval: 30000,
  heartbeatTimeout: 5000
}
```

---

## Testing

### Unit Tests

Run unit tests for components and stores:

```bash
# Run all tests
npm run test

# Run tests in watch mode
npm run test:watch

# Run tests with UI
npm run test:ui

# Run tests with coverage
npm run test:coverage
```

### E2E Tests

Run end-to-end tests with Playwright:

```bash
# Run all E2E tests
npm run test:e2e

# Run E2E tests in headed mode
npm run test:e2e:headed

# Debug E2E tests
npm run test:e2e:debug
```

### Manual Testing Checklist

- [ ] Verify Gantt editor loads correctly
- [ ] Test task creation, editing, deletion
- [ ] Test dependency creation and deletion
- [ ] Test zoom in/out functionality
- [ ] Test view switching (Gantt, Kanban, Calendar, Dashboard)
- [ ] Test virtual scrolling with 1000+ tasks
- [ ] Test WebSocket collaboration
- [ ] Test export functionality (CSV, JSON, PDF, ICS)
- [ ] Test mobile responsiveness
- [ ] Test guided tour
- [ ] Test AI suggestions
- [ ] Test template application

---

## Deployment

### Production Build

```bash
# Build frontend
cd /home/julei/backend/newstatic
npm run build

# Build backend
cd /home/julei/backend
go build -o bin/server cmd/server/main.go
```

### Environment Variables

Create `.env.production`:

```env
# API
VUE_APP_API_URL=https://api.yourdomain.com/api/v1
VUE_APP_WS_URL=wss://api.yourdomain.com/api/v1/ws

# Feature flags
VUE_APP_ENABLE_GANTT=true
VUE_APP_ENABLE_COLLABORATION=true
VUE_APP_ENABLE_AI_SUGGESTIONS=true

# Performance
VUE_APP_VIRTUAL_SCROLL_BUFFER=200
VUE_APP_MAX_TASKS_RENDER=1000
```

### Nginx Configuration

For WebSocket support:

```nginx
server {
    listen 443 ssl;
    server_name yourdomain.com;

    # Frontend
    location / {
        root /var/www/newstatic/dist;
        try_files $uri $uri/ /index.html;
    }

    # API
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # WebSocket
    location /ws/ {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_read_timeout 86400;
    }
}
```

### Database Optimization

```sql
-- Create indexes for performance
CREATE INDEX idx_gantt_tasks_project_id ON gantt_tasks(project_id);
CREATE INDEX idx_gantt_tasks_parent_id ON gantt_tasks(parent_id);
CREATE INDEX idx_gantt_dependencies_project_id ON gantt_dependencies(project_id);
CREATE INDEX idx_gantt_dependencies_from_task ON gantt_dependencies(from_task_id);
CREATE INDEX idx_gantt_dependencies_to_task ON gantt_dependencies(to_task_id);
CREATE INDEX idx_gantt_history_task_id ON gantt_task_history(task_id);
CREATE INDEX idx_gantt_comments_task_id ON gantt_comments(task_id);
```

---

## Verification Steps

After completing the integration:

1. **Check Dependencies:**
   ```bash
   npm list vue-virtual-scroller vue-tour socket.io-client
   ```

2. **Check Component Imports:**
   - Verify main.js includes all registrations
   - Verify router includes Gantt routes
   - Verify API client has all endpoints

3. **Check Backend:**
   - Verify database tables exist
   - Verify API routes are registered
   - Verify WebSocket hub is running

4. **Check Frontend:**
   - Start dev server: `npm run dev`
   - Navigate to Gantt view
   - Open browser console and check for errors
   - Test basic functionality

5. **Performance Check:**
   - Load 1000+ tasks
   - Verify virtual scrolling works
   - Check memory usage in DevTools
   - Verify frame rate remains smooth (60fps)

---

## Troubleshooting

For common issues and solutions, refer to [TROUBLESHOOTING.md](./TROUBLESHOOTING.md)

For performance optimization, refer to [PERFORMANCE_GUIDE.md](./PERFORMANCE_GUIDE.md)

---

## Support

For questions or issues:
- Check the [COMPONENT_REFERENCE.md](./COMPONENT_REFERENCE.md) for detailed component documentation
- Review [API_ENDPOINTS.md](./API_ENDPOINTS.md) for API documentation
- See [MIGRATION_GUIDE.md](./MIGRATION_GUIDE.md) for migration from old implementation

---

**Document Version:** 1.0.0
**Last Updated:** 2026-02-19
