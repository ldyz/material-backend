# Gantt Chart Phase 2: Advanced Scheduling Components

This document describes all Phase 2 advanced scheduling components created for the Gantt chart system.

## Components Overview

### 1. Constraint System (Sprint 2.1)

#### Frontend Components

**File**: `newstatic/src/components/gantt/dialogs/ConstraintEditDialog.vue`
- Dialog for editing task constraints
- Supports 6 constraint types:
  - MSO (Must Start On) - 必须开始于
  - MFO (Must Finish On) - 必须完成于
  - SNET (Start No Earlier Than) - 不早于开始
  - SNLT (Start No Later Than) - 不晚于开始
  - FNET (Finish No Earlier Than) - 不早于完成
  - FNLT (Finish No Later Than) - 不晚于完成
- Features:
  - Date picker for constraint date
  - Visual preview of constraint effect
  - Constraint validation
  - Integration with undo/redo system
  - Bilingual labels (Chinese/English)

**File**: `newstatic/src/utils/ganttConstraints.js`
Utility functions for constraint management:
- `validateConstraint(task, type, date)` - Validates if constraint can be applied
- `calculateConstraintImpact(task, type, date)` - Calculates impact on task dates
- `applyConstraint(task, type, date)` - Applies constraint and returns updated task
- `checkConstraintSatisfied(task)` - Checks if task satisfies its constraint
- `getConstraintViolations(tasks)` - Gets all tasks violating constraints
- `removeConstraint(task)` - Removes constraint from task

#### Backend Components

**File**: `internal/api/progress/constraint.go`
Go models and repository:
- `Constraint` struct with validation
- `ConstraintRepository` for database operations
- Constraint type constants
- `ApplyToTask()` method to adjust dates based on constraint
- `IsSatisfied()` method to check constraint compliance

**File**: `internal/api/progress/constraint_handler.go`
HTTP handlers:
- `GET /progress/tasks/:id/constraints` - Get task constraints
- `POST /progress/tasks/:id/constraints` - Create constraint
- `PUT /progress/constraints/:id` - Update constraint
- `DELETE /progress/constraints/:id` - Delete constraint
- `POST /progress/tasks/:id/apply-constraint` - Apply constraint to task

---

### 2. Resource Leveling (Sprint 2.2)

#### Frontend Components

**File**: `newstatic/src/components/gantt/dialogs/ResourceLevelingDialog.vue`
Dialog for resource leveling operations:
- Resource conflict detection and display
- Manual and automatic leveling modes
- Leveling options:
  - Priority-based leveling (priority, duration, slack)
  - Range selection (all tasks or selected)
  - Task splitting option
  - Dependency adjustment option
- Before/after Gantt comparison preview
- Leveling statistics (tasks delayed, max delay, project extension)

**File**: `newstatic/src/utils/resourceLeveling.js`
Resource leveling algorithms:
- `detectResourceConflicts(tasks, resources)` - Detects overallocation
- `applyResourceLeveling(tasks, resources, options)` - Applies leveling
- `calculateResourceOverallocation(tasks, resources)` - Statistics
- `generateLevelingSuggestions(conflicts)` - Resolution suggestions
- `calculateLevelingStatistics(originalTasks, leveledTasks)` - Compare results
- `optimizeTaskSchedule(tasks, resources)` - Automatic optimization

**File**: `newstatic/src/components/gantt/views/ResourceHistogram.vue`
Visual resource allocation histogram:
- Resource selection dropdown
- Daily/weekly view modes
- Visual bar chart showing allocation vs capacity
- Overallocation highlighting
- Interactive tooltips with task details
- Capacity reference line
- Period click events for navigation

**File**: `newstatic/src/components/gantt/views/GanttMiniView.vue`
Compact Gantt preview component:
- Miniature timeline view
- Task bar visualization
- Conflict zone highlighting
- Used in leveling dialog for before/after comparison

#### Backend Components

**File**: `internal/api/progress/resource_leveling.go`
Resource leveling service:
- `ResourceAssignment` model
- `Resource` model with capacity
- `ResourceConflict` detection
- `ResourceLevelingService` with:
  - `DetectConflicts()` - Find overallocation
  - `LevelResources()` - Apply leveling algorithm
  - `calculateStatistics()` - Generate leveling stats

API endpoints (to be implemented in router):
- `POST /progress/project/:id/level-resources` - Level resources
- `GET /progress/project/:id/resource-conflicts` - Get conflicts

---

### 3. Advanced Dependencies (Sprint 2.3)

#### Frontend Components

**File**: `newstatic/src/utils/dependencyValidator.js`
Dependency validation utilities:
- `DependencyTypes` constants (FS, FF, SS, SF)
- `validateDependency(predecessor, successor, dependency)` - Validate single dependency
- `detectCircularDependencies(tasks)` - DFS-based circular dependency detection
- `validateAllDependencies(tasks)` - Validate all dependencies in project
- `analyzeDependencyPaths(taskId, tasks, direction)` - Trace dependency chains
- `calculatePathDuration(path)` - Calculate path length in days
- `getAffectedTasks(taskId, tasks)` - Find tasks affected by date change
- `calculateLagLimits(predecessor, successor, type)` - Min/max lag values

**File**: `newstatic/src/utils/criticalPath.js`
Enhanced Critical Path Method (CPM):
- `calculateCriticalPath(tasks)` - Main CPM calculation
  - Forward pass (Early Start/Finish)
  - Backward pass (Late Start/Finish)
  - Slack calculation
  - Critical path identification
- `calculateTaskSlack(task, cpmResult)` - Get task slack/float
- `isTaskCritical(task, cpmResult)` - Check if task is critical
- `calculateFreeSlack(task, tasks, cpmResult)` - Free slack calculation
- `calculateInterferingSlack(task, tasks, cpmResult)` - Interfering slack
- `getAllCriticalPaths(tasks)` - Find all critical paths (can be multiple)
- `generateCriticalPathReport(tasks)` - Generate comprehensive report
- `calculateDelayImpact(task, delayDays, tasks)` - Analyze delay impact

Features:
- Support for all 4 dependency types (FS, FF, SS, SF)
- Lag/lead time in calculations
- Multiple parallel critical paths detection
- Comprehensive slack analysis

---

### 4. Calendar System (Sprint 2.4)

#### Frontend Components

**File**: `newstatic/src/stores/calendarStore.js`
Pinia store for calendar management:
- State:
  - `calendars` - All project calendars
  - `projectCalendar` - Default project calendar
  - `taskCalendars` - Task-specific calendars
  - `standardPresets` - Predefined calendar templates
- Actions:
  - `fetchCalendars(projectId)` - Load calendars
  - `createCalendar(calendarData)` - Create new calendar
  - `updateCalendar(calendarId, updates)` - Update calendar
  - `deleteCalendar(calendarId)` - Delete calendar
  - `setProjectCalendar(calendarId)` - Set project default
  - `assignTaskCalendar(taskId, calendarId)` - Assign to task
  - `addHoliday(calendarId, holiday)` - Add holiday
  - `addException(calendarId, exception)` - Add exception
- Utilities:
  - `isWorkingDay(date, calendarId)` - Check if working day
  - `calculateWorkingDays(startDate, endDate, calendarId)` - Count working days
  - `addWorkingDays(startDate, days, calendarId)` - Add working days
  - `subtractWorkingDays(startDate, days, calendarId)` - Subtract working days
  - `getWorkingHours(date, calendarId)` - Get working hours for date

**File**: `newstatic/src/components/gantt/dialogs/CalendarDialog.vue`
Calendar editing dialog:
- Calendar type selection (Standard, 24/7, Custom)
- Working days selector (Sunday-Saturday checkboxes)
- Working hours configuration (start/end time)
- Holiday management:
  - Add/remove holidays
  - Recurring holiday support
  - Holiday list with dates and names
- Exception dates management:
  - Mark non-working days as working
  - Mark working days as non-working
  - Custom working hours for exceptions
- Calendar preview (30-day view)
- Bilingual interface

**File**: `newstatic/src/components/gantt/views/CalendarPreview.vue`
Calendar preview component:
- 30-day calendar grid visualization
- Color-coded days:
  - Working days (blue)
  - Non-working days (gray)
  - Holidays (red)
  - Exceptions (orange)
- Legend for day types
- Shows day numbers and day names

#### Backend Components

**File**: `internal/api/progress/calendar.go`
Calendar models and repository:
- `Calendar` model with working days/hours
- `Holiday` model for non-working days
- `CalendarException` model for exceptions
- `TaskCalendar` model for task assignments
- `CalendarRepository` for database operations
- `WorkingTimeCalculator` service:
  - `IsWorkingDay(date, calendar)` - Check working day
  - `CalculateWorkingDays(start, end, calendar)` - Count working days
  - `AddWorkingDays(date, days, calendar)` - Add working days

API endpoints (to be implemented in router):
- `GET /progress/project/:id/calendars` - List calendars
- `POST /progress/project/:id/calendars` - Create calendar
- `PUT /progress/calendars/:id` - Update calendar
- `DELETE /progress/calendars/:id` - Delete calendar
- `GET /progress/project/:id/calendar` - Get project calendar
- `PUT /progress/project/:id/calendar` - Set project calendar
- `POST /progress/tasks/:id/calendar` - Assign task calendar
- `DELETE /progress/tasks/:id/calendar` - Remove task calendar
- `POST /progress/calendars/:id/holidays` - Add holiday
- `DELETE /progress/holidays/:id` - Delete holiday
- `POST /progress/calendars/:id/exceptions` - Add exception
- `DELETE /progress/exceptions/:id` - Delete exception

---

## Integration Points

### 1. GanttStore Integration

Add to `ganttStore.js`:

```javascript
// Import new utilities
import { validateConstraint, calculateConstraintImpact } from '@/utils/ganttConstraints'
import { detectResourceConflicts, applyResourceLeveling } from '@/utils/resourceLeveling'
import { validateAllDependencies, calculateCriticalPath } from '@/utils/dependencyValidator'
import { useCalendarStore } from '@/stores/calendarStore'

// Add to state
constraints: ref(new Map()),
resourceConflicts: ref([]),
criticalPathData: ref(null),

// Add methods
async updateTaskConstraint(taskId, constraint) {
  // Apply constraint to task
},
async checkResourceConflicts() {
  // Check for overallocation
},
async calculateCriticalPath() {
  // Update critical path
},
async applyResourceLeveling(options) {
  // Level resources
}
```

### 2. API Integration

Create API client methods in `api/index.js`:

```javascript
// Constraints
export const getTaskConstraints = (taskId) => request.get(`/progress/tasks/${taskId}/constraints`)
export const createConstraint = (taskId, data) => request.post(`/progress/tasks/${taskId}/constraints`, data)
export const updateConstraint = (constraintId, data) => request.put(`/progress/constraints/${constraintId}`, data)
export const deleteConstraint = (constraintId) => request.delete(`/progress/constraints/${constraintId}`)

// Resource Leveling
export const getResourceConflicts = (projectId) => request.get(`/progress/project/${projectId}/resource-conflicts`)
export const levelResources = (projectId, options) => request.post(`/progress/project/${projectId}/level-resources`, options)

// Calendars
export const getCalendars = (projectId) => request.get(`/progress/project/${projectId}/calendars`)
export const createCalendar = (projectId, data) => request.post(`/progress/project/${projectId}/calendars`, data)
export const updateCalendar = (calendarId, data) => request.put(`/progress/calendars/${calendarId}`, data)
export const deleteCalendar = (calendarId) => request.delete(`/progress/calendars/${calendarId}`)
export const setProjectCalendar = (projectId, calendarId) => request.put(`/progress/project/${projectId}/calendar`, { calendar_id: calendarId })
```

### 3. Component Integration

Add to `GanttToolbar.vue`:

```vue
<template>
  <!-- Add new toolbar buttons -->
  <el-button @click="showConstraintsDialog">
    {{ t('gantt.toolbar.constraints') }}
  </el-button>
  <el-button @click="showResourceLevelingDialog">
    {{ t('gantt.toolbar.resourceLeveling') }}
  </el-button>
  <el-button @click="showCalendarDialog">
    {{ t('gantt.toolbar.calendar') }}
  </el-button>
</template>

<script setup>
import ConstraintEditDialog from '@/components/gantt/dialogs/ConstraintEditDialog.vue'
import ResourceLevelingDialog from '@/components/gantt/dialogs/ResourceLevelingDialog.vue'
import CalendarDialog from '@/components/gantt/dialogs/CalendarDialog.vue'
</script>
```

---

## Usage Examples

### 1. Using Constraints

```javascript
import { validateConstraint, applyConstraint } from '@/utils/ganttConstraints'

// Validate a constraint
const validation = validateConstraint(task, 'SNET', '2026-02-20')
if (!validation.valid) {
  console.error(validation.message)
}

// Apply constraint
const updatedTask = applyConstraint(task, 'SNET', '2026-02-20')
```

### 2. Resource Leveling

```javascript
import { detectResourceConflicts, applyResourceLeveling } from '@/utils/resourceLeveling'

// Detect conflicts
const conflicts = await detectResourceConflicts(tasks, resources)

// Level resources
const leveledTasks = await applyResourceLeveling(tasks, resources, {
  priority: 'priority',
  range: 'all',
  allowSplitting: false,
  adjustDependencies: true
})
```

### 3. Critical Path Analysis

```javascript
import { calculateCriticalPath, isTaskCritical } from '@/utils/criticalPath'

// Calculate critical path
const cpmResult = calculateCriticalPath(tasks)

// Check if task is critical
const critical = isTaskCritical(task, cpmResult)
```

### 4. Calendar Calculations

```javascript
import { useCalendarStore } from '@/stores/calendarStore'

const calendarStore = useCalendarStore()

// Check if working day
const isWorking = calendarStore.isWorkingDay(new Date(), calendarId)

// Calculate working days
const days = calendarStore.calculateWorkingDays(startDate, endDate, calendarId)

// Add working days
const newDate = calendarStore.addWorkingDays(startDate, 10, calendarId)
```

---

## Testing Considerations

### Unit Tests

- Constraint validation logic
- Resource conflict detection
- Circular dependency detection
- Critical path calculation
- Calendar working day calculations

### Integration Tests

- Constraint API endpoints
- Resource leveling workflow
- Calendar CRUD operations
- End-to-end scheduling scenarios

### E2E Tests

- Complete constraint application workflow
- Resource leveling with conflict resolution
- Calendar creation and task assignment
- Critical path visualization updates

---

## Performance Optimizations

1. **Memoization**: Cache calculation results for tasks
2. **Web Workers**: Run heavy calculations (critical path) in background
3. **Virtual Scrolling**: For large task lists in dialogs
4. **Debouncing**: For real-time validation and preview updates
5. **Lazy Loading**: Load resource conflicts on demand

---

## Future Enhancements

1. **Baseline Comparison**: Compare current schedule with baseline
2. **Portfolio Management**: Multiple project resource leveling
3. **What-If Analysis**: Scenario modeling
4. **Resource Pooling**: Shared resource pools across projects
5. **Advanced Constraints**: Constraint groups, constraint priorities
6. **Calendar Templates**: Reusable calendar templates
7. **Resource Skills**: Skill-based resource assignment
8. **Cost Calculations**: Resource cost estimation

---

## File Structure

```
newstatic/src/
├── components/gantt/
│   ├── dialogs/
│   │   ├── ConstraintEditDialog.vue
│   │   ├── ResourceLevelingDialog.vue
│   │   └── CalendarDialog.vue
│   └── views/
│       ├── ResourceHistogram.vue
│       ├── GanttMiniView.vue
│       └── CalendarPreview.vue
├── stores/
│   └── calendarStore.js
└── utils/
    ├── ganttConstraints.js
    ├── resourceLeveling.js
    ├── dependencyValidator.js
    └── criticalPath.js

internal/api/progress/
├── constraint.go
├── constraint_handler.go
├── resource_leveling.go
└── calendar.go
```

---

## Dependencies

### Frontend
- Vue 3
- Pinia (state management)
- Element Plus (UI components)
- Vue I18n (internationalization)

### Backend
- Go 1.19+
- Gin (web framework)
- GORM (ORM)

---

## License

These components are part of the project and follow the project's license.
