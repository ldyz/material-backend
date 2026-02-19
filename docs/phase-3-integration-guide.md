# Phase 3 Integration Guide

## Quick Start Integration

### Step 1: Add API Endpoints to Router

Update `/home/julei/backend/internal/api/progress/routes.go`:

```go
package progress

import (
    "github.com/gin-gonic/gin"
)

// RegisterRoutes registers all progress API routes
func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
    // Initialize repositories
    commentRepo := NewCommentRepository(db)
    changeLogRepo := NewChangeLogRepository(db)
    taskRepo := NewTaskRepository(db)

    // Initialize handlers
    commentHandler := NewCommentHandler(commentRepo)
    changeLogHandler := NewChangeLogHandler(changeLogRepo, taskRepo)

    // API group
    api := r.Group("/api/progress")
    {
        // Comment routes
        api.GET("/tasks/:id/comments", commentHandler.GetComments)
        api.POST("/tasks/:id/comments", commentHandler.CreateComment)
        api.PUT("/comments/:id", commentHandler.UpdateComment)
        api.DELETE("/comments/:id", commentHandler.DeleteComment)
        api.POST("/comments/:id/resolve", commentHandler.ResolveComment)
        api.POST("/comments/:id/unresolve", commentHandler.UnresolveComment)
        api.GET("/project/:id/comments", commentHandler.GetProjectComments)
        api.GET("/comments/my", commentHandler.GetUserComments)
        api.GET("/comments/mentions", commentHandler.GetMentions)
        api.GET("/tasks/:id/comments/count", commentHandler.GetUnresolvedCount)

        // Change log routes
        api.GET("/project/:id/change-log", changeLogHandler.GetProjectChangeLog)
        api.GET("/tasks/:id/history", changeLogHandler.GetTaskHistory)
        api.GET("/history/:entity_type/:entity_id", changeLogHandler.GetEntityHistory)
        api.POST("/changes/:id/rollback", changeLogHandler.RollbackChange)
        api.GET("/project/:id/statistics", changeLogHandler.GetStatistics)
        api.POST("/change-log/cleanup", changeLogHandler.DeleteOldChangeLogs)
        api.GET("/changes/my", changeLogHandler.GetUserChanges)

        // WebSocket route
        api.GET("/ws", HandleWebSocket)
    }

    // Initialize WebSocket hub
    InitWebSocket()
}
```

### Step 2: Create API Client Methods

Update `/home/julei/backend/newstatic/src/api/index.js`:

```javascript
// Comments API
export const getComments = (taskId) => request.get(`/progress/tasks/${taskId}/comments`)
export const createComment = (taskId, data) => request.post(`/progress/tasks/${taskId}/comments`, data)
export const updateComment = (commentId, data) => request.put(`/progress/comments/${commentId}`, data)
export const deleteComment = (commentId) => request.delete(`/progress/comments/${commentId}`)
export const resolveComment = (commentId) => request.post(`/progress/comments/${commentId}/resolve`)
export const unresolveComment = (commentId) => request.post(`/progress/comments/${commentId}/unresolve`)
export const getProjectComments = (projectId) => request.get(`/progress/project/${projectId}/comments`)
export const getMyComments = (params) => request.get('/progress/comments/my', { params })
export const getMentions = (params) => request.get('/progress/comments/mentions', { params })
export const getUnresolvedCount = (taskId) => request.get(`/progress/tasks/${taskId}/comments/count`)

// Change Log API
export const getProjectHistory = (projectId, params) =>
  request.get(`/progress/project/${projectId}/change-log`, { params })
export const getTaskHistory = (taskId, params) =>
  request.get(`/progress/tasks/${taskId}/history`, { params })
export const getEntityHistory = (entityType, entityId, params) =>
  request.get(`/progress/history/${entityType}/${entityId}`, { params })
export const rollbackChange = (changeId) =>
  request.post(`/progress/changes/${changeId}/rollback`)
export const getChangeStatistics = (projectId) =>
  request.get(`/progress/project/${projectId}/statistics`)
export const cleanupChangeLogs = (params) =>
  request.post('/progress/change-log/cleanup', null, { params })
export const getMyChanges = (params) =>
  request.get('/progress/changes/my', { params })
```

### Step 3: Initialize Collaboration in App

Update `/home/julei/backend/newstatic/src/main.js`:

```javascript
import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)
app.use(ElementPlus)

// Make collaboration store globally available
import { useCollaborationStore } from '@/stores/collaborationStore'
app.config.globalProperties.$collaboration = useCollaborationStore()

app.mount('#app')
```

### Step 4: Add Collaboration to Gantt View

Update your Gantt view component:

```vue
<template>
  <div class="gantt-view">
    <!-- Existing Gantt components -->

    <!-- Add collaboration toolbar -->
    <CollaborationToolbar
      v-if="showCollaboration"
      :project-id="projectId"
    />

    <!-- Add side panels -->
    <el-drawer
      v-model="commentsVisible"
      title="Comments"
      size="400px"
    >
      <CommentsPanel
        :task-id="selectedTaskId"
        :project-users="projectUsers"
      />
    </el-drawer>

    <el-drawer
      v-model="historyVisible"
      title="Change History"
      size="500px"
    >
      <HistoryPanel
        :task-id="selectedTaskId"
        :project-id="projectId"
        :project-users="projectUsers"
      />
    </el-drawer>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useCollaborationStore } from '@/stores/collaborationStore'
import CommentsPanel from '@/components/gantt/panels/CommentsPanel.vue'
import HistoryPanel from '@/components/gantt/panels/HistoryPanel.vue'
import CollaborationToolbar from '@/components/gantt/CollaborationToolbar.vue'

const collaborationStore = useCollaborationStore()
const projectId = ref(null)
const selectedTaskId = ref(null)
const commentsVisible = ref(false)
const historyVisible = ref(false)
const projectUsers = ref([])

onMounted(async () => {
  // Connect to WebSocket when project loads
  if (projectId.value) {
    await collaborationStore.connect(WS_URL, projectId.value)

    // Listen for updates
    collaborationStore.on('task:update', handleTaskUpdate)
  }

  // Load project users
  await loadProjectUsers()
})

onUnmounted(() => {
  collaborationStore.disconnect()
})

async function handleTaskUpdate(data) {
  // Refresh task data
  console.log('Task updated:', data)
}

async function loadProjectUsers() {
  // Fetch project members
  const response = await getProjectMembers(projectId.value)
  projectUsers.value = response.data
}
</script>
```

### Step 5: Add Database Migrations

Create migration file:

```sql
-- Comments table
CREATE TABLE IF NOT EXISTS task_comments (
    id SERIAL PRIMARY KEY,
    task_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    parent_id INTEGER,
    is_resolved BOOLEAN DEFAULT FALSE,
    mentions TEXT[] DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES task_comments(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_task_comments ON task_comments(task_id);
CREATE INDEX IF NOT EXISTS idx_comment_parent ON task_comments(parent_id);
CREATE INDEX IF NOT EXISTS idx_comment_user ON task_comments(user_id);

-- Change logs table
CREATE TABLE IF NOT EXISTS change_logs (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id INTEGER NOT NULL,
    action_type VARCHAR(20) NOT NULL,
    changes JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_project_changes ON change_logs(project_id);
CREATE INDEX IF NOT EXISTS idx_entity_changes ON change_logs(entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_user_changes ON change_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_change_dates ON change_logs(created_at);
```

### Step 6: Enable WebSocket in Backend

Update main server file:

```go
package main

import (
    "github.com/gin-gonic/gin"
    "your-project/internal/api/progress"
)

func main() {
    r := gin.Default()

    // CORS middleware for WebSocket
    r.Use(CORSMiddleware())

    // Register routes
    progress.RegisterRoutes(r, db)

    r.Run(":8080")
}

func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}
```

## Testing the Integration

### 1. Test WebSocket Connection

```javascript
// In browser console
const ws = new WebSocket('ws://localhost:8080/api/progress/ws?projectId=1&userId=1')

ws.onopen = () => console.log('Connected!')
ws.onmessage = (msg) => console.log('Received:', msg.data)
```

### 2. Test Comments

```javascript
// Create a comment
await createComment(taskId, {
  content: '<p>This is a test comment</p>',
  mentions: [2, 3]
})

// Get comments
const comments = await getComments(taskId)
console.log(comments)
```

### 3. Test Change Tracking

```javascript
// Get history
const history = await getTaskHistory(taskId)
console.log(history)

// Get statistics
const stats = await getChangeStatistics(projectId)
console.log(stats)
```

## Configuration

### Environment Variables

Add to `.env`:

```
# WebSocket
VITE_WS_URL=ws://localhost:8080
WS_PORT=8080

# Change Tracking
CHANGE_LOG_RETENTION_DAYS=90
```

### WebSocket Configuration

```javascript
// config/websocket.js
export default {
  url: import.meta.env.VITE_WS_URL || 'ws://localhost:8080',
  options: {
    timeout: 10000,
    reconnectAttempts: 5,
    reconnectDelay: 3000
  }
}
```

## Troubleshooting

### WebSocket Not Connecting

1. Check backend is running on correct port
2. Verify CORS configuration
3. Check firewall rules
4. Review browser console for errors

### Comments Not Appearing

1. Verify database tables exist
2. Check API responses in Network tab
3. Ensure task ID is valid
4. Review user permissions

### Change History Empty

1. Verify change tracking is enabled
2. Check database for change_logs table
3. Ensure user ID is being captured
4. Review repository queries

---

## Next Steps

1. Run database migrations
2. Restart backend server
3. Test WebSocket connection
4. Add comments to a task
5. Make changes and view history
6. Test rollback functionality

For full documentation, see `/home/julei/backend/docs/phase-3-collaboration.md`
