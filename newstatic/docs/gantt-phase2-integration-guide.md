# Phase 2 Integration Guide

This guide provides step-by-step instructions for integrating the Phase 2 components into the existing Gantt chart system.

## Step 1: API Client Setup

Add to `src/api/index.js`:

```javascript
// ==================== Constraints API ====================
export const getTaskConstraints = (taskId) => {
  return request.get(`/progress/tasks/${taskId}/constraints`)
}

export const createConstraint = (taskId, data) => {
  return request.post(`/progress/tasks/${taskId}/constraints`, data)
}

export const updateConstraint = (constraintId, data) => {
  return request.put(`/progress/constraints/${constraintId}`, data)
}

export const deleteConstraint = (constraintId) => {
  return request.delete(`/progress/constraints/${constraintId}`)
}

export const applyConstraintToTask = (taskId, constraint) => {
  return request.post(`/progress/tasks/${taskId}/apply-constraint`, constraint)
}

// ==================== Resource Leveling API ====================
export const getResourceConflicts = (projectId) => {
  return request.get(`/progress/project/${projectId}/resource-conflicts`)
}

export const levelResources = (projectId, options) => {
  return request.post(`/progress/project/${projectId}/level-resources`, options)
}

// ==================== Calendar API ====================
export const getCalendars = (projectId) => {
  return request.get(`/progress/project/${projectId}/calendars`)
}

export const createCalendar = (projectId, data) => {
  return request.post(`/progress/project/${projectId}/calendars`, data)
}

export const updateCalendar = (calendarId, data) => {
  return request.put(`/progress/calendars/${calendarId}`, data)
}

export const deleteCalendar = (calendarId) => {
  return request.delete(`/progress/calendars/${calendarId}`)
}

export const getProjectCalendar = (projectId) => {
  return request.get(`/progress/project/${projectId}/calendar`)
}

export const setProjectCalendar = (projectId, calendarId) => {
  return request.put(`/progress/project/${projectId}/calendar`, { calendar_id: calendarId })
}

export const assignTaskCalendar = (taskId, calendarId) => {
  return request.post(`/progress/tasks/${taskId}/calendar`, { calendar_id: calendarId })
}

export const removeTaskCalendar = (taskId) => {
  return request.delete(`/progress/tasks/${taskId}/calendar`)
}

export const addHoliday = (calendarId, holiday) => {
  return request.post(`/progress/calendars/${calendarId}/holidays`, holiday)
}

export const deleteHoliday = (holidayId) => {
  return request.delete(`/progress/holidays/${holidayId}`)
}

export const addException = (calendarId, exception) => {
  return request.post(`/progress/calendars/${calendarId}/exceptions`, exception)
}

export const deleteException = (exceptionId) => {
  return request.delete(`/progress/exceptions/${exceptionId}`)
}
```

## Step 2: Update GanttStore

Add to `src/stores/ganttStore.js`:

```javascript
// Add imports at top
import { validateConstraint, calculateConstraintImpact } from '@/utils/ganttConstraints'
import { detectResourceConflicts, applyResourceLeveling } from '@/utils/resourceLeveling'
import { calculateCriticalPath, isTaskCritical } from '@/utils/criticalPath'
import { useCalendarStore } from '@/stores/calendarStore'

// Add to state
const state = reactive({
  // ... existing state ...

  // Phase 2 additions
  constraints: new Map(),
  resourceConflicts: [],
  criticalPathData: null,
  calendarStore: useCalendarStore(),

  // Dialog visibility
  constraintDialogVisible: false,
  resourceLevelingDialogVisible: false,
  calendarDialogVisible: false
})

// Add computed properties
const criticalTasks = computed(() => {
  if (!state.criticalPathData) return []
  return state.tasks.filter(task =>
    isTaskCritical(task, state.criticalPathData)
  )
})

// Add methods
const actions = {
  // ... existing methods ...

  // ==================== Constraint Methods ====================
  async updateTaskConstraint(taskId, constraint) {
    try {
      const task = state.tasks.find(t => t.id === taskId)
      if (!task) return

      const updatedTask = validateConstraint(task, constraint.type, constraint.date)
      if (!updatedTask.valid) {
        throw new Error(updatedTask.message)
      }

      // Apply constraint
      const result = applyConstraint(task, constraint.type, constraint.date)

      // Update in backend
      await progressApi.updateTask(taskId, {
        start_date: result.start_date,
        end_date: result.end_date,
        constraint: constraint
      })

      // Update local state
      const index = state.tasks.findIndex(t => t.id === taskId)
      if (index !== -1) {
        state.tasks[index] = result
      }

      state.constraints.set(taskId, constraint)

      ElMessage.success('Constraint applied successfully')
    } catch (error) {
      console.error('Failed to apply constraint:', error)
      ElMessage.error(error.message || 'Failed to apply constraint')
      throw error
    }
  },

  // ==================== Resource Leveling Methods ====================
  async checkResourceConflicts() {
    try {
      const conflicts = await detectResourceConflicts(
        state.tasks,
        state.resources
      )
      state.resourceConflicts = conflicts
      return conflicts
    } catch (error) {
      console.error('Failed to check resource conflicts:', error)
      ElMessage.error('Failed to check resource conflicts')
      return []
    }
  },

  async applyResourceLeveling(options) {
    try {
      const leveledTasks = await applyResourceLeveling(
        state.tasks,
        state.resources,
        options
      )

      // Update backend
      // TODO: Batch update tasks in backend

      // Update local state
      state.tasks = leveledTasks

      // Recalculate conflicts
      await this.checkResourceConflicts()

      ElMessage.success('Resources leveled successfully')
      return leveledTasks
    } catch (error) {
      console.error('Failed to level resources:', error)
      ElMessage.error('Failed to level resources')
      throw error
    }
  },

  // ==================== Critical Path Methods ====================
  async calculateCriticalPath() {
    try {
      const cpmResult = calculateCriticalPath(state.tasks)
      state.criticalPathData = cpmResult
      return cpmResult
    } catch (error) {
      console.error('Failed to calculate critical path:', error)
      ElMessage.error('Failed to calculate critical path')
      return null
    }
  },

  // ==================== Dialog Methods ====================
  openConstraintDialog(taskId) {
    state.selectedTaskId = taskId
    state.constraintDialogVisible = true
  },

  openResourceLevelingDialog() {
    state.resourceLevelingDialogVisible = true
  },

  openCalendarDialog(calendarId = null) {
    state.calendarDialogVisible = true
  }
}
```

## Step 3: Update GanttToolbar

Add to `src/components/gantt/core/GanttToolbar.vue`:

```vue
<template>
  <div class="gantt-toolbar">
    <!-- Existing toolbar buttons -->

    <!-- Phase 2 additions -->
    <el-divider direction="vertical" />

    <el-tooltip :content="t('gantt.toolbar.constraints')" placement="bottom">
      <el-button
        :disabled="!selectedTask"
        @click="handleConstraints"
      >
        <el-icon><Lock /></el-icon>
      </el-button>
    </el-tooltip>

    <el-tooltip :content="t('gantt.toolbar.resourceLeveling')" placement="bottom">
      <el-button @click="handleResourceLeveling">
        <el-icon><Grid /></el-icon>
      </el-button>
    </el-tooltip>

    <el-tooltip :content="t('gantt.toolbar.calendar')" placement="bottom">
      <el-button @click="handleCalendar">
        <el-icon><Calendar /></el-icon>
      </el-button>
    </el-tooltip>

    <!-- Phase 2 Dialogs -->
    <ConstraintEditDialog
      v-model="ganttStore.constraintDialogVisible"
      :task="selectedTask"
      @constraint-updated="handleConstraintUpdated"
    />

    <ResourceLevelingDialog
      v-model="ganttStore.resourceLevelingDialogVisible"
      :tasks="ganttStore.tasks"
      :resources="ganttStore.resources"
      @resources-leveled="handleResourcesLeveled"
    />

    <CalendarDialog
      v-model="ganttStore.calendarDialogVisible"
      :project-id="ganttStore.projectId"
      @calendar-updated="handleCalendarUpdated"
    />
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useGanttStore } from '@/stores/ganttStore'
import { useI18n } from 'vue-i18n'
import ConstraintEditDialog from '@/components/gantt/dialogs/ConstraintEditDialog.vue'
import ResourceLevelingDialog from '@/components/gantt/dialogs/ResourceLevelingDialog.vue'
import CalendarDialog from '@/components/gantt/dialogs/CalendarDialog.vue'

const { t } = useI18n()
const ganttStore = useGanttStore()

const selectedTask = computed(() =>
  ganttStore.tasks.find(t => t.id === ganttStore.selectedTaskId)
)

const handleConstraints = () => {
  ganttStore.openConstraintDialog(ganttStore.selectedTaskId)
}

const handleResourceLeveling = () => {
  ganttStore.openResourceLevelingDialog()
}

const handleCalendar = () => {
  ganttStore.openCalendarDialog()
}

const handleConstraintUpdated = ({ task, constraint }) => {
  console.log('Constraint updated:', task, constraint)
  // Refresh task data
}

const handleResourcesLeveled = (leveledTasks) => {
  console.log('Resources leveled:', leveledTasks)
  // Refresh view
}

const handleCalendarUpdated = (calendar) => {
  console.log('Calendar updated:', calendar)
  // Refresh calendar data
}
</script>
```

## Step 4: Add i18n Translations

Add to `src/locales/en-US.js`:

```javascript
export default {
  gantt: {
    toolbar: {
      constraints: 'Constraints',
      resourceLeveling: 'Resource Leveling',
      calendar: 'Calendar'
    },
    constraints: {
      editTitle: 'Edit Task Constraint',
      type: 'Constraint Type',
      selectType: 'Select constraint type',
      date: 'Constraint Date',
      selectDate: 'Select date',
      preview: 'Constraint Preview',
      currentSchedule: 'Current Schedule',
      constraintDate: 'Constraint Date',
      wouldShift: 'Task would shift by {days} days {direction}',
      earlier: 'earlier',
      later: 'later',
      noEffect: 'No schedule change required',
      descMSO: 'Task must start exactly on this date',
      descMFO: 'Task must finish exactly on this date',
      descSNET: 'Task cannot start before this date',
      descSNLT: 'Task cannot start after this date',
      descFNET: 'Task cannot finish before this date',
      descFNLT: 'Task cannot finish after this date',
      saveSuccess: 'Constraint saved successfully',
      saveError: 'Failed to save constraint',
      undoDescription: 'Constraint {constraint} on task {task}'
    },
    resourceLeveling: {
      title: 'Resource Leveling',
      conflictsFound: '{count} resource conflicts found',
      noConflicts: 'No resource conflicts detected',
      levelingOptions: 'Leveling Options',
      manual: 'Manual Leveling',
      manualDesc: 'Manually resolve resource conflicts',
      auto: 'Automatic Leveling',
      autoDesc: 'Automatically level resources using heuristics',
      priority: 'Leveling Priority',
      taskPriority: 'Task Priority',
      taskDuration: 'Task Duration',
      totalSlack: 'Total Slack',
      range: 'Leveling Range',
      allTasks: 'All Tasks',
      selectedTasks: 'Selected Tasks',
      splitTasks: 'Allow Task Splitting',
      adjustDependencies: 'Adjust Dependencies',
      conflictDetails: 'Resource Conflicts',
      resource: 'Resource',
      date: 'Date',
      assigned: 'Assigned',
      conflictingTasks: 'Conflicting Tasks',
      resolve: 'Resolve',
      preview: 'Before/After Preview',
      before: 'Before',
      after: 'After',
      tasksDelayed: 'Tasks Delayed',
      maxDelay: 'Maximum Delay',
      projectExtension: 'Project Extension',
      levelResources: 'Level Resources',
      loadConflictsError: 'Failed to load resource conflicts',
      previewError: 'Failed to preview leveling',
      levelingError: 'Failed to level resources',
      success: 'Resources leveled successfully'
    },
    resourceHistogram: {
      title: 'Resource Allocation',
      selectResource: 'Select resource',
      daily: 'Daily',
      weekly: 'Weekly',
      allocation: 'Allocation',
      capacity: 'Capacity',
      overallocation: 'Overallocation',
      assignedTasks: 'Assigned Tasks',
      noResourceSelected: 'No resource selected'
    },
    calendar: {
      editTitle: 'Edit Calendar',
      calendarType: 'Calendar Type',
      selectPreset: 'Select calendar preset',
      workingDays: 'Working Days',
      workingHours: 'Working Hours',
      startTime: 'Start Time',
      endTime: 'End Time',
      to: 'to',
      totalHours: 'Total Hours',
      holidays: 'Holidays',
      addHoliday: 'Add Holiday',
      noHolidays: 'No holidays defined',
      date: 'Date',
      name: 'Name',
      recurring: 'Recurring',
      yearly: 'Yearly',
      exceptions: 'Exceptions',
      addException: 'Add Exception',
      noExceptions: 'No exceptions defined',
      type: 'Type',
      working: 'Working',
      nonWorking: 'Non-Working',
      description: 'Description',
      preview: 'Preview',
      fillRequiredFields: 'Please fill all required fields',
      saveSuccess: 'Calendar saved successfully',
      saveError: 'Failed to save calendar',
      sunday: 'Sun',
      monday: 'Mon',
      tuesday: 'Tue',
      wednesday: 'Wed',
      thursday: 'Thu',
      friday: 'Fri',
      saturday: 'Sat',
      sun: 'S',
      mon: 'M',
      tue: 'T',
      wed: 'W',
      thu: 'T',
      fri: 'F',
      sat: 'S'
    }
  }
}
```

Add to `src/locales/zh-CN.js` (Chinese translations):

```javascript
// Similar structure with Chinese translations
export default {
  gantt: {
    toolbar: {
      constraints: '约束',
      resourceLeveling: '资源平衡',
      calendar: '日历'
    },
    constraints: {
      editTitle: '编辑任务约束',
      type: '约束类型',
      // ... etc
    }
    // ... etc
  }
}
```

## Step 5: Backend Router Setup

Add to `internal/api/progress/routes.go` (create if doesn't exist):

```go
package progress

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterRoutes registers all progress-related routes
func RegisterRoutes(r *gin.RouterGroup, db *gorm.DB) {
	// Initialize repositories
	constraintRepo := NewConstraintRepository(db)
	calendarRepo := NewCalendarRepository(db)

	// Initialize handlers
	constraintHandler := NewConstraintHandler(constraintRepo)

	// Progress task routes (existing)
	tasks := r.Group("/tasks")
	{
		// Existing task routes...

		// Constraint routes
		tasks("/:id/constraints", constraintHandler.GetConstraints)
		tasks("/:id/constraints", constraintHandler.CreateConstraint)
		tasks("/:id/apply-constraint", constraintHandler.ApplyConstraintToTask)
		tasks("/:id/calendar", setTaskCalendar)
	}

	// Constraint routes
	constraints := r.Group("/constraints")
	{
		constraints("/:id", constraintHandler.UpdateConstraint)
		constraints("/:id", constraintHandler.DeleteConstraint)
	}

	// Resource leveling routes
	project := r.Group("/project/:id")
	{
		project.GET("/resource-conflicts", getResourceConflicts)
		project.POST("/level-resources", levelResources)
		project.GET("/calendar", getProjectCalendar)
		project.PUT("/calendar", setProjectCalendar)
		project.GET("/calendars", getCalendars)
		project.POST("/calendars", createCalendar)
	}

	// Calendar routes
	calendars := r.Group("/calendars")
	{
		calendars("/:id", updateCalendar)
		calendars("/:id", deleteCalendar)
		calendars("/:id/holidays", addHoliday)
		calendars("/:id/exceptions", addException)
	}

	// Holiday routes
	holidays := r.Group("/holidays")
	{
		holidays("/:id", deleteHoliday)
	}

	// Exception routes
	exceptions := r.Group("/exceptions")
	{
		exceptions("/:id", deleteException)
	}
}

// Handler functions to implement
func getResourceConflicts(c *gin.Context) {
	// TODO: Implement
}

func levelResources(c *gin.Context) {
	// TODO: Implement
}

func getProjectCalendar(c *gin.Context) {
	// TODO: Implement
}

func setProjectCalendar(c *gin.Context) {
	// TODO: Implement
}

func getCalendars(c *gin.Context) {
	// TODO: Implement
}

func createCalendar(c *gin.Context) {
	// TODO: Implement
}

func updateCalendar(c *gin.Context) {
	// TODO: Implement
}

func deleteCalendar(c *gin.Context) {
	// TODO: Implement
}

func setTaskCalendar(c *gin.Context) {
	// TODO: Implement
}

func addHoliday(c *gin.Context) {
	// TODO: Implement
}

func addException(c *gin.Context) {
	// TODO: Implement
}
```

## Step 6: Testing

Create integration test file: `src/components/gantt/__tests__/phase2-integration.spec.js`:

```javascript
import { describe, it, expect, vi } from 'vitest'
import { validateConstraint } from '@/utils/ganttConstraints'
import { detectResourceConflicts } from '@/utils/resourceLeveling'
import { calculateCriticalPath } from '@/utils/criticalPath'

describe('Phase 2 Integration', () => {
  describe('Constraint System', () => {
    it('should validate SNET constraint', () => {
      const task = {
        id: 1,
        name: 'Test Task',
        start_date: '2026-02-20',
        end_date: '2026-02-25'
      }

      const result = validateConstraint(task, 'SNET', '2026-02-22')

      expect(result.valid).toBe(true)
      expect(result.message).toContain('delayed')
    })

    it('should calculate constraint impact', () => {
      const task = {
        id: 1,
        start_date: '2026-02-20',
        end_date: '2026-02-25'
      }

      const impact = calculateConstraintImpact(task, 'MSO', '2026-02-22')

      expect(impact.wouldShift).toBe(true)
      expect(impact.shiftDays).toBe(2)
    })
  })

  describe('Resource Leveling', () => {
    it('should detect resource conflicts', async () => {
      const tasks = [
        { id: 1, start_date: '2026-02-20', end_date: '2026-02-22', resources: [{ resource_id: 1, units: 100 }] },
        { id: 2, start_date: '2026-02-20', end_date: '2026-02-22', resources: [{ resource_id: 1, units: 100 }] }
      ]

      const resources = [
        { id: 1, name: 'Resource 1', capacity: 100 }
      ]

      const conflicts = await detectResourceConflicts(tasks, resources)

      expect(conflicts.length).toBe(1)
      expect(conflicts[0].overallocated).toBe(100)
    })
  })

  describe('Critical Path', () => {
    it('should calculate critical path', () => {
      const tasks = [
        { id: 1, start_date: '2026-02-20', end_date: '2026-02-22', dependencies: [] },
        { id: 2, start_date: '2026-02-22', end_date: '2026-02-25', dependencies: [1] }
      ]

      const result = calculateCriticalPath(tasks)

      expect(result.criticalTasks.has(1)).toBe(true)
      expect(result.criticalTasks.has(2)).toBe(true)
      expect(result.projectDuration).toBe(5)
    })
  })
})
```

## Step 7: Build and Deploy

```bash
# Frontend
cd newstatic
npm run build

# Backend
cd ../
go build -o bin/server cmd/server/main.go

# Run tests
npm run test
go test ./internal/api/progress/...
```

## Verification Checklist

- [ ] API client methods added
- [ ] GanttStore updated with new methods
- [ ] GanttToolbar includes new buttons
- [ ] Dialogs imported and configured
- [ ] i18n translations added
- [ ] Backend routes registered
- [ ] Database migrations run (if needed)
- [ ] Components render without errors
- [ ] API calls work correctly
- [ ] Undo/redo integration tested
- [ ] Calendar calculations accurate
- [ ] Resource conflicts detected
- [ ] Critical path calculated correctly

## Troubleshooting

### Issues and Solutions

**Issue**: Constraint dialog doesn't open
- **Solution**: Check that constraintDialogVisible is properly set in store

**Issue**: Resource histogram shows no data
- **Solution**: Ensure resources are loaded and tasks have resource assignments

**Issue**: Calendar calculations incorrect
- **Solution**: Verify workingDays is parsed correctly from JSON

**Issue**: Circular dependency detection slow
- **Solution**: Implement memoization or limit graph traversal depth

**Issue**: Backend routes not found
- **Solution**: Verify routes are registered in main router configuration

## Next Steps

1. Add more comprehensive error handling
2. Implement loading skeletons
3. Add accessibility features (ARIA labels)
4. Create user documentation
5. Add unit tests for utilities
6. Add E2E tests for workflows
7. Performance optimization for large projects
8. Add export/import functionality
