# Gantt Chart Editor - Migration Guide

**Version:** 1.0.0
**Last Updated:** 2026-02-19
**Author:** Material Management System Team

---

## Table of Contents

1. [Overview](#overview)
2. [What Changes](#what-changes)
3. [What Stays the Same](#what-stays-the-same)
4. [Data Migration](#data-migration)
5. [Feature Mapping](#feature-mapping)
6. [Breaking Changes](#breaking-changes)
7. [Rollback Plan](#rollback-plan)
8. [Step-by-Step Migration](#step-by-step-migration)
9. [Testing After Migration](#testing-after-migration)

---

## Overview

This guide helps you migrate from the legacy `GanttChart.vue` component to the new `GanttEditor.vue` component. The new component is built with performance, scalability, and modern UX in mind.

### Why Migrate?

- **10x Performance:** Virtual scrolling handles 1000+ tasks smoothly
- **Better UX:** Modern UI with guided tours, keyboard shortcuts, and mobile support
- **More Features:** AI suggestions, workflow automation, advanced reporting
- **Better Collaboration:** Real-time WebSocket collaboration
- **Maintainability:** Modular architecture with clear separation of concerns

### Migration Scope

- **Component:** `GanttChart.vue` → `GanttEditor.vue`
- **Location:** `/src/components/progress/GanttChart.vue` → `/src/components/gantt/core/GanttEditor.vue`
- **Affected Files:**
  - Views that import GanttChart
  - API calls to Gantt endpoints
  - Store usage
  - Router configuration

---

## What Changes

### 1. Component Structure

**Before (GanttChart.vue):**
```
GanttChart.vue (39,261 lines)
├── GanttToolbar.vue
├── GanttHeader.vue
├── GanttStats.vue
├── GanttStatusBar.vue
├── GanttLegend.vue
└── Internal components
```

**After (GanttEditor.vue):**
```
GanttEditor.vue (762 lines)
├── core/
│   ├── GanttEditor.vue
│   ├── GanttToolbar.vue
│   ├── GanttStatusBar.vue
│   └── GanttBody.vue
├── timeline/
│   ├── VirtualTimeline.vue
│   ├── TimelineBackground.vue
│   ├── TimelineGrid.vue
│   ├── TaskBar.vue
│   └── DependencyLines.vue
├── table/
│   ├── VirtualTaskList.vue
│   ├── TaskRow.vue
│   └── EditableCell.vue
└── [30+ more organized components]
```

### 2. Props Interface

**Before:**
```vue
<GanttChart
  :plan-id="planId"
  :plan-data="planData"
  :tasks="tasks"
  :read-only="readOnly"
  :height="containerHeight"
  @task-update="handleTaskUpdate"
  @task-delete="handleTaskDelete"
/>
```

**After:**
```vue
<GanttEditor
  :project-id="projectId"
  :tasks="tasks"
  :dependencies="dependencies"
  :read-only="readOnly"
  :height="containerHeight"
  :initial-view-mode="'gantt'"
  :enable-collaboration="true"
  @update:tasks="handleTasksUpdate"
  @update:dependencies="handleDependenciesUpdate"
  @task-select="handleTaskSelect"
/>
```

### 3. Data Model

**Before:**
```javascript
// Task object in old GanttChart
{
  id: 'task-1',
  planId: 'plan-1',
  name: 'Task 1',
  startDate: '2024-01-01',
  endDate: '2024-01-10',
  progress: 50,
  status: '进行中',
  // ... limited fields
}
```

**After:**
```javascript
// Task object in new GanttEditor
{
  id: 'task-1',
  projectId: 'proj-1',
  name: 'Task 1',
  description: 'Task description',
  startDate: new Date('2024-01-01'),
  endDate: new Date('2024-01-10'),
  duration: 10,
  progress: 50,
  status: 'in_progress',
  priority: 'medium',
  assignee: 'user-1',
  parentId: null,
  position: 0,
  color: '#409EFF',
  milestone: false,
  constraint: null,
  customFields: {},
  // ... many more fields
}
```

### 4. Event System

**Before:**
```javascript
// Events emitted by old GanttChart
@task-update      // Single task updated
@task-delete      // Task deleted
@add-dependency   // Dependency added
@delete-dependency // Dependency deleted
@save             // Save triggered
```

**After:**
```javascript
// Events emitted by new GanttEditor
@update:tasks         // All tasks updated
@update:dependencies  // All dependencies updated
@task-select         // Task selected
@task-update         // Task updated (with metadata)
@task-create         // Task created
@task-delete         // Task deleted
@dependency-create   // Dependency created
@dependency-delete   // Dependency deleted
@view-change         // View mode changed
@zoom-change         // Zoom level changed
@export             // Export triggered
@save              // Save triggered (with all data)
```

### 5. Store Usage

**Before:**
```javascript
import { usePlanStore } from '@/stores/planStore'

const planStore = usePlanStore()
planStore.loadPlanData(planId)
planStore.updateTask(taskId, updates)
```

**After:**
```javascript
import { useGanttStore } from '@/stores/ganttStore'
import { useUndoRedoStore } from '@/stores/undoRedoStore'
import { useCollaborationStore } from '@/stores/collaborationStore'

const ganttStore = useGanttStore()
const undoRedoStore = useUndoRedoStore()
const collaborationStore = useCollaborationStore()

ganttStore.initialize()
ganttStore.setTasks(tasks)
ganttStore.setDependencies(dependencies)
```

---

## What Stays the Same

### 1. Core Functionality

- ✅ Task display on timeline
- ✅ Task creation, editing, deletion
- ✅ Dependency management
- ✅ Zoom in/out
- ✅ Date navigation
- ✅ Export functionality
- ✅ Basic filtering and search

### 2. User Experience

- ✅ Visual appearance (with improvements)
- ✅ Keyboard shortcuts
- ✅ Context menus
- ✅ Drag-and-drop interactions
- ✅ Status indicators

### 3. API Compatibility

The backend API remains mostly compatible. New endpoints are added but old ones still work:

```javascript
// Old API (still works)
GET /api/plan/:planId/tasks
POST /api/plan/:planId/tasks
PUT /api/plan/:planId/tasks/:taskId

// New API (enhanced)
GET /api/gantt/projects/:projectId/tasks
POST /api/gantt/projects/:projectId/tasks
PATCH /api/gantt/projects/:projectId/tasks/:taskId
PATCH /api/gantt/projects/:projectId/tasks/bulk
```

---

## Data Migration

### Step 1: Update Task Data Format

Create a migration utility function:

**Location:** `/src/utils/ganttMigration.js`

```javascript
import { parseISO, format } from 'date-fns'

/**
 * Migrate task data from old format to new format
 */
export function migrateTaskData(oldTasks) {
  return oldTasks.map(oldTask => {
    // Convert date strings to Date objects
    const startDate = typeof oldTask.startDate === 'string'
      ? parseISO(oldTask.startDate)
      : oldTask.startDate

    const endDate = typeof oldTask.endDate === 'string'
      ? parseISO(oldTask.endDate)
      : oldTask.endDate

    // Map status values
    const statusMap = {
      '未开始': 'not_started',
      '进行中': 'in_progress',
      '已完成': 'completed',
      '已延期': 'delayed',
      '已阻塞': 'blocked'
    }

    // Map priority values
    const priorityMap = {
      '低': 'low',
      '中': 'medium',
      '高': 'high',
      '紧急': 'critical'
    }

    // Calculate duration if not present
    const duration = oldTask.duration || Math.ceil(
      (endDate - startDate) / (1000 * 60 * 60 * 24)
    )

    return {
      id: oldTask.id,
      projectId: oldTask.planId || oldTask.projectId,
      name: oldTask.name,
      description: oldTask.description || '',
      startDate,
      endDate,
      duration,
      progress: oldTask.progress || 0,
      status: statusMap[oldTask.status] || oldTask.status || 'not_started',
      priority: priorityMap[oldTask.priority] || oldTask.priority || 'medium',
      assignee: oldTask.assignee || oldTask.assigneeId,
      parentId: oldTask.parentId || null,
      position: oldTask.position || oldTask.sortOrder || 0,
      color: oldTask.color || null,
      milestone: oldTask.milestone || false,
      constraint: oldTask.constraint || null,
      customFields: oldTask.customFields || {},
      // Preserve any other existing fields
      ...oldTask
    }
  })
}

/**
 * Migrate dependency data from old format to new format
 */
export function migrateDependencyData(oldDependencies) {
  return oldDependencies.map(dep => ({
    id: dep.id,
    fromTaskId: dep.fromTaskId || dep.predecessorId,
    toTaskId: dep.toTaskId || dep.successorId,
    type: normalizeDependencyType(dep.type),
    lag: dep.lag || 0
  }))
}

function normalizeDependencyType(type) {
  const typeMap = {
    'FS': 'finish-to-start',
    'finish-to-start': 'finish-to-start',
    'SS': 'start-to-start',
    'start-to-start': 'start-to-start',
    'FF': 'finish-to-finish',
    'finish-to-finish': 'finish-to-finish',
    'SF': 'start-to-finish',
    'start-to-finish': 'start-to-finish'
  }
  return typeMap[type] || 'finish-to-start'
}
```

### Step 2: Create Migration Script

**Location:** `/scripts/migrateGanttData.js`

```javascript
import axios from 'axios'
import { migrateTaskData, migrateDependencyData } from '../src/utils/ganttMigration.js'

async function migrateProject(projectId) {
  try {
    console.log(`Migrating project ${projectId}...`)

    // Fetch old data
    const [oldTasks, oldDependencies] = await Promise.all([
      axios.get(`/api/plan/${projectId}/tasks`),
      axios.get(`/api/plan/${projectId}/dependencies`)
    ])

    // Migrate data
    const newTasks = migrateTaskData(oldTasks.data)
    const newDependencies = migrateDependencyData(oldDependencies.data)

    // Save new data
    await Promise.all([
      axios.post(`/api/gantt/projects/${projectId}/migrate`, {
        tasks: newTasks,
        dependencies: newDependencies
      })
    ])

    console.log(`✓ Migration complete for project ${projectId}`)
  } catch (error) {
    console.error(`✗ Migration failed for project ${projectId}:`, error)
    throw error
  }
}

// Run migration for all projects
async function migrateAllProjects() {
  const { data: projects } = await axios.get('/api/projects')

  console.log(`Found ${projects.length} projects to migrate`)

  for (const project of projects) {
    await migrateProject(project.id)
  }

  console.log('All migrations complete!')
}

// Execute
migrateAllProjects()
```

### Step 3: Backend Migration Endpoint

Add this endpoint to temporarily handle migration:

```go
// internal/api/gantt/migration.go
package gantt

func MigrateProjectData(c *gin.Context) {
    projectID := c.Param("projectId")

    var req struct {
        Tasks        []Task       `json:"tasks"`
        Dependencies []Dependency `json:"dependencies"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // Backup old data
    if err := backupOldData(projectID); err != nil {
        c.JSON(500, gin.H{"error": "Failed to backup data"})
        return
    }

    // Save new data
    if err := saveMigratedData(projectID, req.Tasks, req.Dependencies); err != nil {
        c.JSON(500, gin.H{"error": "Failed to save migrated data"})
        return
    }

    c.JSON(200, gin.H{"message": "Migration successful"})
}
```

---

## Feature Mapping

| Old Feature | New Feature | Notes |
|-------------|-------------|-------|
| Basic task display | Virtual task list | Performance improved with virtual scrolling |
| Simple toolbar | Enhanced toolbar | Added undo/redo, templates, bulk edit |
| Task properties dialog | Inline editing + EditableCell | Faster editing experience |
| Basic dependencies | Advanced dependencies | Added lag, different dependency types |
| Export PNG/PDF | Export PNG/PDF/CSV/JSON/ICS | More export formats |
| Statistics panel | Dashboard view | Rich analytics with charts |
| - | Kanban view | New: Agile-style board view |
| - | Calendar view | New: Calendar visualization |
| - | Comments panel | New: Task discussions |
| - | History panel | New: Change tracking with diff viewer |
| - | Resource leveling | New: Optimize resource allocation |
| - | Critical path analysis | New: Identify critical tasks |
| - | AI suggestions | New: Intelligent recommendations |
| - | Guided tour | New: Onboarding for new users |
| - | Real-time collaboration | New: Multi-user editing via WebSocket |
| - | Templates | New: Save and reuse task templates |
| - | Workflow automation | New: Automate repetitive tasks |
| - | Mobile view | New: Touch-optimized mobile interface |
| - | Minimap | New: Navigate large projects easily |

---

## Breaking Changes

### 1. Date Format

**Before:** String format (`'2024-01-01'`)

**After:** Date objects (`new Date('2024-01-01')`)

**Impact:** Any code that expects date strings will break.

**Fix:**
```javascript
// Before
const task = { startDate: '2024-01-01' }

// After
const task = { startDate: new Date('2024-01-01') }
```

### 2. Status Values

**Before:** Chinese values (`'未开始'`, `'进行中'`, etc.)

**After:** English enum values (`'not_started'`, `'in_progress'`, etc.)

**Impact:** Any hardcoded status checks will fail.

**Fix:**
```javascript
// Before
if (task.status === '进行中') { ... }

// After
if (task.status === 'in_progress') { ... }
// OR use constants
import { TaskStatus } from '@/constants/ganttConstants'
if (task.status === TaskStatus.IN_PROGRESS) { ... }
```

### 3. Component Imports

**Before:**
```javascript
import GanttChart from '@/components/progress/GanttChart.vue'
```

**After:**
```javascript
import { GanttEditor } from '@/components/gantt'
```

### 4. Store Structure

**Before:**
```javascript
planStore.tasks
planStore.dependencies
```

**After:**
```javascript
ganttStore.state.tasks
ganttStore.state.dependencies
```

### 5. API Endpoints

**Before:**
```javascript
/api/plan/:planId/tasks
```

**After:**
```javascript
/api/gantt/projects/:projectId/tasks
```

---

## Rollback Plan

### Preparation

Before migrating, create a backup:

```bash
# Backup database
pg_dump -U username -d database_name > backup_before_migration.sql

# Backup frontend code
cp -r /home/julei/backend/newstatic /home/julei/backend/newstatic_backup
```

### Rollback Procedure

If issues arise after migration:

#### 1. Database Rollback

```bash
# Restore database
psql -U username -d database_name < backup_before_migration.sql
```

#### 2. Frontend Rollback

```bash
# Restore frontend code
rm -rf /home/julei/backend/newstatic
mv /home/julei/backend/newstatic_backup /home/julei/backend/newstatic

# Rebuild
cd /home/julei/backend/newstatic
npm install
npm run build
```

#### 3. Gradual Rollback

If only some features have issues:

```vue
<!-- Use old component for specific views -->
<template>
  <div v-if="useLegacyGantt">
    <GanttChart
      :plan-id="planId"
      :tasks="tasks"
      @task-update="handleTaskUpdate"
    />
  </div>
  <div v-else>
    <GanttEditor
      :project-id="projectId"
      :tasks="migratedTasks"
      @update:tasks="handleTasksUpdate"
    />
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import GanttChart from '@/components/progress/GanttChart.vue'
import { GanttEditor } from '@/components/gantt'
import { migrateTaskData } from '@/utils/ganttMigration'

const useLegacyGatt = ref(false) // Toggle this flag to rollback

const migratedTasks = computed(() => {
  return useLegacyGantt.value ? tasks.value : migrateTaskData(tasks.value)
})
</script>
```

---

## Step-by-Step Migration

### Phase 1: Preparation (1-2 days)

1. **Backup everything**
   ```bash
   # Database backup
   pg_dump -U username -d database_name > backup.sql

   # Code backup
   git checkout -b backup-before-gantt-migration
   ```

2. **Create migration branch**
   ```bash
   git checkout -b feature/gantt-migration
   ```

3. **Install new dependencies**
   ```bash
   cd /home/julei/backend/newstatic
   npm install vue-virtual-scroller vue-tour
   ```

4. **Review and document current usage**
   ```bash
   # Find all files using GanttChart
   grep -r "GanttChart" src/
   ```

### Phase 2: Backend Updates (1-2 days)

1. **Add migration endpoint** (see Data Migration section)

2. **Run database migrations**
   ```bash
   psql -U username -d database_name -f docs/DATABASE_MIGRATIONS.sql
   ```

3. **Add new API routes** (see MAIN_INTEGRATION_GUIDE.md)

4. **Test backward compatibility**
   - Ensure old endpoints still work
   - Test new endpoints

### Phase 3: Frontend Updates (2-3 days)

1. **Update main.js** (see MAIN_INTEGRATION_GUIDE.md)

2. **Create migration utility**
   - Create `/src/utils/ganttMigration.js`
   - Test with sample data

3. **Update views one by one**

   **Example: Progress.vue**

   ```vue
   <!-- Before -->
   <template>
     <GanttChart
       :plan-id="planId"
       :plan-data="planData"
       :tasks="tasks"
       @task-update="handleTaskUpdate"
     />
   </template>

   <script setup>
   import GanttChart from '@/components/progress/GanttChart.vue'
   </script>

   <!-- After -->
   <template>
     <GanttEditor
       :project-id="projectId"
       :tasks="migratedTasks"
       :dependencies="migratedDependencies"
       @update:tasks="handleTasksUpdate"
       @update:dependencies="handleDependenciesUpdate"
     />
   </template>

   <script setup>
   import { computed } from 'vue'
   import { GanttEditor } from '@/components/gantt'
   import { migrateTaskData, migrateDependencyData } from '@/utils/ganttMigration'

   const migratedTasks = computed(() => migrateTaskData(tasks.value))
   const migratedDependencies = computed(() => migrateDependencyData(dependencies.value))
   </script>
   ```

4. **Update event handlers**

   ```javascript
   // Before
   const handleTaskUpdate = (task) => {
     const index = tasks.value.findIndex(t => t.id === task.id)
     if (index !== -1) {
       tasks.value[index] = task
     }
   }

   // After
   const handleTasksUpdate = (newTasks) => {
     tasks.value = newTasks
   }
   ```

### Phase 4: Data Migration (1 day)

1. **Run migration script for all projects**
   ```bash
   node scripts/migrateGanttData.js
   ```

2. **Verify migrated data**
   ```javascript
   // Check a sample project
   const project = await api.getProject('proj-123')
   console.log('Tasks:', project.tasks.length)
   console.log('Dependencies:', project.dependencies.length)
   ```

3. **Fix any data issues**

### Phase 5: Testing (2-3 days)

See "Testing After Migration" section below.

### Phase 6: Deployment (1 day)

1. **Build production version**
   ```bash
   npm run build
   ```

2. **Deploy to staging**
   - Test all functionality
   - Performance testing
   - Load testing

3. **Deploy to production**
   - Monitor logs
   - Be ready to rollback

---

## Testing After Migration

### Unit Tests

```javascript
// tests/unit/ganttMigration.test.js
import { describe, it, expect } from 'vitest'
import { migrateTaskData, migrateDependencyData } from '@/utils/ganttMigration'

describe('Gantt Migration', () => {
  it('should migrate task data correctly', () => {
    const oldTasks = [{
      id: 'task-1',
      planId: 'plan-1',
      name: 'Task 1',
      startDate: '2024-01-01',
      endDate: '2024-01-10',
      status: '进行中',
      priority: '高'
    }]

    const newTasks = migrateTaskData(oldTasks)

    expect(newTasks[0].status).toBe('in_progress')
    expect(newTasks[0].priority).toBe('high')
    expect(newTasks[0].startDate).toBeInstanceOf(Date)
  })

  it('should migrate dependency data correctly', () => {
    const oldDeps = [{
      id: 'dep-1',
      predecessorId: 'task-1',
      successorId: 'task-2',
      type: 'FS'
    }]

    const newDeps = migrateDependencyData(oldDeps)

    expect(newDeps[0].fromTaskId).toBe('task-1')
    expect(newDeps[0].toTaskId).toBe('task-2')
    expect(newDeps[0].type).toBe('finish-to-start')
  })
})
```

### Integration Tests

```javascript
// tests/integration/ganttEditor.test.js
import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { GanttEditor } from '@/components/gantt'

describe('GanttEditor Integration', () => {
  it('should render with migrated data', () => {
    const wrapper = mount(GanttEditor, {
      props: {
        projectId: 'proj-1',
        tasks: migratedTasks,
        dependencies: migratedDependencies
      }
    })

    expect(wrapper.find('.gantt-editor').exists()).toBe(true)
    expect(wrapper.findAll('.task-bar')).toHaveLength(migratedTasks.length)
  })

  it('should handle task updates', async () => {
    const wrapper = mount(GanttEditor, {
      props: {
        projectId: 'proj-1',
        tasks: migratedTasks
      }
    })

    await wrapper.vm.updateTask('task-1', { progress: 75 })

    expect(wrapper.emitted('update:tasks')).toBeTruthy()
  })
})
```

### E2E Tests

```javascript
// tests/e2e/ganttMigration.spec.js
import { test, expect } from '@playwright/test'

test.describe('Gantt Migration E2E', () => {
  test('should load migrated project correctly', async ({ page }) => {
    await page.goto('/progress/gantt?projectId=proj-1')

    // Wait for GanttEditor to load
    await page.waitForSelector('.gantt-editor')

    // Verify tasks are displayed
    const taskBars = await page.locator('.task-bar').count()
    expect(taskBars).toBeGreaterThan(0)

    // Verify interactions work
    await page.click('.task-bar:first-child')
    await expect(page.locator('.task-bar.is-selected')).toBeVisible()
  })

  test('should create new task with new format', async ({ page }) => {
    await page.goto('/progress/gantt?projectId=proj-1')

    // Click add task button
    await page.click('[data-testid="add-task-button"]')

    // Fill task form
    await page.fill('[data-testid="task-name-input"]', 'New Task')
    await page.fill('[data-testid="task-start-date"]', '2024-02-01')
    await page.fill('[data-testid="task-end-date"]', '2024-02-10')

    // Save
    await page.click('[data-testid="save-task-button"]')

    // Verify task was created
    await expect(page.locator('text=New Task')).toBeVisible()
  })
})
```

### Manual Testing Checklist

- [ ] Load project with 10 tasks
- [ ] Load project with 100 tasks
- [ ] Load project with 1000+ tasks (test virtual scrolling)
- [ ] Create new task
- [ ] Edit existing task
- [ ] Delete task
- [ ] Create dependency
- [ ] Delete dependency
- [ ] Zoom in/out
- [ ] Navigate dates
- [ ] Switch to Kanban view
- [ ] Switch to Calendar view
- [ ] Switch to Dashboard view
- [ ] Export to CSV
- [ ] Export to PDF
- [ ] Test mobile view
- [ ] Test keyboard shortcuts
- [ ] Test undo/redo
- [ ] Test search functionality
- [ ] Test bulk edit
- [ ] Test comments panel
- [ ] Test history panel

### Performance Testing

```javascript
// tests/performance/ganttPerformance.test.js
import { test, expect } from '@playwright/test'

test.describe('Gantt Performance', () => {
  test('should handle 1000 tasks smoothly', async ({ page }) => {
    // Load test project with 1000 tasks
    await page.goto('/progress/gantt?projectId=perf-test-1000')

    // Measure initial render time
    const renderStartTime = Date.now()
    await page.waitForSelector('.gantt-editor')
    const renderTime = Date.now() - renderStartTime

    expect(renderTime).toBeLessThan(3000) // Should render in < 3 seconds

    // Measure scroll performance
    const scrollStartTime = Date.now()
    await page.mouse.wheel(0, 1000)
    await page.waitForTimeout(500)
    const scrollTime = Date.now() - scrollStartTime

    expect(scrollTime).toBeLessThan(1000) // Should scroll smoothly

    // Check memory usage (via performance API)
    const memory = await page.evaluate(() => {
      return performance.memory
    })

    console.log('Memory usage:', memory)
  })
})
```

---

## Common Migration Issues

### Issue 1: Date Parsing Errors

**Problem:** Invalid date errors after migration

**Solution:**
```javascript
// Add date validation in migration utility
function safeParseDate(dateStr) {
  try {
    const date = parseISO(dateStr)
    return isNaN(date.getTime()) ? new Date() : date
  } catch {
    return new Date()
  }
}
```

### Issue 2: Missing Task Fields

**Problem:** Some tasks missing required fields

**Solution:**
```javascript
// Add defaults in migration
export function migrateTaskData(oldTasks) {
  return oldTasks.map(oldTask => ({
    ...DEFAULT_TASK_FIELDS,
    ...oldTask,
    // Apply migrations
    startDate: safeParseDate(oldTask.startDate),
    endDate: safeParseDate(oldTask.endDate),
    status: normalizeStatus(oldTask.status)
  }))
}

const DEFAULT_TASK_FIELDS = {
  progress: 0,
  priority: 'medium',
  status: 'not_started',
  milestone: false,
  position: 0
}
```

### Issue 3: Dependency Reference Errors

**Problem:** Dependencies referencing non-existent tasks

**Solution:**
```javascript
// Validate dependencies
export function migrateDependencyData(oldDependencies, validTaskIds) {
  return oldDependencies
    .filter(dep => validTaskIds.includes(dep.fromTaskId))
    .filter(dep => validTaskIds.includes(dep.toTaskId))
    .map(dep => ({
      ...dep,
      type: normalizeDependencyType(dep.type)
    }))
}
```

---

## Post-Migration Optimization

After successful migration:

1. **Remove old code** (after 1-2 weeks of stable operation)
   ```bash
   rm src/components/progress/GanttChart.vue
   rm src/components/progress/GanttChartRefactored.vue
   ```

2. **Update API to remove deprecated endpoints**

3. **Clean up unused imports**

4. **Update documentation**

5. **Train users on new features**

---

## Support

For migration issues:
- Check [MAIN_INTEGRATION_GUIDE.md](./MAIN_INTEGRATION_GUIDE.md) for detailed setup
- Review [COMPONENT_REFERENCE.md](./COMPONENT_REFERENCE.md) for API changes
- See [TROUBLESHOOTING.md](./TROUBLESHOOTING.md) for common issues
- Check [PERFORMANCE_GUIDE.md](./PERFORMANCE_GUIDE.md) for optimization tips

---

**Document Version:** 1.0.0
**Last Updated:** 2026-02-19
