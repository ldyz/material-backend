# Gantt Editing Components - Quick Reference Guide

## Quick Import Syntax

```javascript
// Import all components
import EditableCell from '@/components/gantt/table/EditableCell.vue'
import TaskTemplatesDialog from '@/components/gantt/dialogs/TaskTemplatesDialog.vue'
import BulkEditDialog from '@/components/gantt/dialogs/BulkEditDialog.vue'

// Import store
import { useUndoRedoStore } from '@/stores/undoRedoStore'
import { CreateTaskCommand, UpdateTaskCommand, BatchUpdateCommand } from '@/stores/undoRedoStore'
```

## Component Quick Setup

### 1. EditableCell - Minimal Setup

```vue
<EditableCell
  v-model="task.name"
  type="text"
  field="task_name"
  :task-id="task.id"
  :original-data="task"
/>
```

### 2. TaskTemplatesDialog - Minimal Setup

```vue
<TaskTemplatesDialog
  v-model="showTemplates"
  :project-id="projectId"
  @created="onTaskCreated"
/>
```

### 3. BulkEditDialog - Minimal Setup

```vue
<BulkEditDialog
  v-model="showBulkEdit"
  :tasks="selectedTasks"
  :project-id="projectId"
  @updated="onBulkUpdated"
/>
```

## Common Props Reference

### EditableCell Types

```javascript
// Text input
<EditableCell type="text" ... />

// Number input
<EditableCell
  type="number"
  :min="0"
  :max="100"
  :precision="0"
  ...
/>

// Date picker
<EditableCell
  type="date"
  date-format="YYYY-MM-DD"
  value-format="YYYY-MM-DD"
  ...
/>

// Select dropdown
<EditableCell
  type="select"
  :options="[
    { value: 'high', label: '高' },
    { value: 'medium', label: '中' },
    { value: 'low', label: '低' }
  ]"
  ...
/>
```

### Validation Examples

```javascript
// Required field
const rules = [
  { required: true, message: '必填项', trigger: 'blur' }
]

// Min/Max length
const rules = [
  { min: 3, max: 100, message: '长度3-100字符', trigger: 'blur' }
]

// Custom validator
const rules = [
  {
    validator: (value) => {
      if (!value) return true
      return value.length <= 100 || '最多100字符'
    },
    trigger: 'blur'
  }
]

// Multiple rules
const rules = [
  { required: true, message: '必填项', trigger: 'blur' },
  { min: 3, max: 100, message: '长度3-100字符', trigger: 'blur' }
]

// Use in component
<EditableCell :rules="rules" ... />
```

## Common Events

### EditableCell Events

```vue
<EditableCell
  @update:modelValue="newValue => task.name = newValue"
  @change="({ field, value, taskId }) => handleChange(field, value, taskId)"
  @edit="field => console.log('Editing:', field)"
  @cancel="field => console.log('Cancelled:', field)"
  ...
/>
```

### TaskTemplatesDialog Events

```vue
<TaskTemplatesDialog
  @created="task => console.log('Created:', task)"
  @template-selected="template => console.log('Selected:', template)"
  @update:modelValue="visible => console.log('Visible:', visible)"
  ...
/>
```

### BulkEditDialog Events

```vue
<BulkEditDialog
  @updated="({ count, fields, changes }) => {
    console.log(`Updated ${count} tasks`)
    console.log('Fields:', fields)
    console.log('Changes:', changes)
  }"
  @update:modelValue="visible => console.log('Visible:', visible)"
  ...
/>
```

## Display Format Examples

```javascript
// Status display
<EditableCell
  :display-format="value => {
    const map = {
      not_started: '未开始',
      in_progress: '进行中',
      completed: '已完成',
      delayed: '已延期'
    }
    return map[value] || value
  }"
  ...
/>

// Progress with percentage
<EditableCell
  :display-format="value => `${value}%`"
  ...
/>

// Date formatting
<EditableCell
  :display-format="value => {
    if (!value) return '-'
    const date = new Date(value)
    return `${date.getMonth() + 1}/${date.getDate()}`
  }"
  ...
/>

// Duration with unit
<EditableCell
  :display-format="value => `${value} 天`"
  ...
/>

// User name lookup
<EditableCell
  :display-format="value => {
    const user = users.find(u => u.id === value)
    return user ? user.name : '未分配'
  }"
  ...
/>
```

## Undo/Redo Integration

### Direct Command Usage

```javascript
import { useUndoRedoStore } from '@/stores/undoRedoStore'
import { CreateTaskCommand, UpdateTaskCommand, BatchUpdateCommand } from '@/stores/undoRedoStore'

const undoRedoStore = useUndoRedoStore()

// Create task
const createCommand = new CreateTaskCommand(
  progressApi.create.bind(progressApi),
  taskData,
  null
)
await undoRedoStore.executeCommand(createCommand)

// Update task
const updateCommand = new UpdateTaskCommand(
  taskId,
  { status: 'completed' },
  originalTaskData
)
await undoRedoStore.executeCommand(updateCommand)

// Batch update
const batchCommand = new BatchUpdateCommand([
  { taskId: 1, updates: { status: 'completed' }, originalData: task1Original },
  { taskId: 2, updates: { status: 'completed' }, originalData: task2Original }
])
await undoRedoStore.executeCommand(batchCommand)

// Undo
await undoRedoStore.undo()

// Redo
await undoRedoStore.redo()

// Clear history
undoRedoStore.clearHistory()
```

### Keyboard Shortcuts

```javascript
// In your component
function handleKeydown(event) {
  // Undo: Ctrl+Z or Cmd+Z
  if ((event.ctrlKey || event.metaKey) && event.key === 'z') {
    event.preventDefault()
    undoRedoStore.undo()
  }

  // Redo: Ctrl+Y or Cmd+Shift+Z
  if ((event.ctrlKey || event.metaKey) && event.key === 'y') {
    event.preventDefault()
    undoRedoStore.redo()
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleQuitkeydown)
})
```

## Common Patterns

### 1. Auto-Save on Change

```vue
<template>
  <EditableCell
    @change="handleCellChange"
    ...
  />
</template>

<script setup>
async function handleCellChange({ field, value, taskId }) {
  // Component already handles API call via undo/redo
  // Just need to refresh data if needed
  await refreshTasks()
}
</script>
```

### 2. Batch Operations

```vue
<template>
  <div>
    <el-button @click="selectAll">全选</el-button>
    <el-button @click="openBulkEdit">批量编辑</el-button>

    <BulkEditDialog
      v-model="bulkEditVisible"
      :tasks="selectedTaskObjects"
      :project-id="projectId"
      @updated="handleBulkUpdated"
    />
  </div>
</template>

<script setup>
const selectedTasks = ref([])
const bulkEditVisible = ref(false)

const selectedTaskObjects = computed(() =>
  tasks.value.filter(t => selectedTasks.value.includes(t.id))
)

function selectAll() {
  selectedTasks.value = tasks.value.map(t => t.id)
}

function openBulkEdit() {
  if (selectedTasks.value.length === 0) {
    ElMessage.warning('请先选择任务')
    return
  }
  bulkEditVisible.value = true
}

function handleBulkUpdated({ count }) {
  ElMessage.success(`已更新 ${count} 个任务`)
  selectedTasks.value = []
  refreshTasks()
}
</script>
```

### 3. Template-Based Creation

```vue
<template>
  <el-button @click="showTemplates = true">
    从模板创建
  </el-button>

  <TaskTemplatesDialog
    v-model="showTemplates"
    :project-id="projectId"
    :start-date="defaultStartDate"
    @created="handleTaskCreated"
  />
</template>

<script setup>
const showTemplates = ref(false)
const defaultStartDate = computed(() => new Date().toISOString().split('T')[0])

function handleTaskCreated(task) {
  ElMessage.success('任务已创建')
  refreshTasks()
}
</script>
```

## Error Handling

```javascript
// Try-catch with component operations
try {
  await undoRedoStore.undo()
  ElMessage.success('已撤销')
} catch (error) {
  console.error('Undo failed:', error)
  ElMessage.error(error.message || '操作失败')
}

// Validation error handling
<EditableCell
  :rules="[
    {
      validator: (value) => {
        if (!value) return true
        if (value.length > 100) {
          return '最多100个字符'
        }
        return true
      },
      trigger: 'blur'
    }
  ]"
  ...
/>
```

## Styling Customization

```css
/* Override CSS variables */
:root {
  --row-hover-bg: #f5f7fa;
  --row-selected-bg: #ecf5ff;
  --color-primary: #409eff;
}

/* Custom cell styling */
.editable-cell :deep(.status-completed) {
  color: #67c23a;
  font-weight: 500;
}

.editable-cell :deep(.priority-high) {
  color: #f56c6c;
  font-weight: 500;
}
```

## Performance Tips

1. **Lazy Loading**: Load user options on dialog open
2. **Virtual Scrolling**: Use for large task lists
3. **Debouncing**: Add for rapid edits if needed
4. **Computed Properties**: Use for filtered/transformed data

## Troubleshooting

### Issue: Edit mode not starting
- Check if `readonly` prop is false
- Ensure field is not disabled
- Check for z-index conflicts

### Issue: Validation not working
- Verify rules array format
- Check trigger value ('blur' or 'change')
- Ensure validator returns true or error message

### Issue: Undo/Redo not working
- Verify store is properly initialized
- Check if command canUndo() returns true
- Ensure async/await is used correctly

### Issue: Bulk edit not applying
- Check if tasks array has data
- Verify selectedFields array
- Ensure at least one field is selected

## Additional Resources

- Full documentation: `/docs/gantt-editing-components.md`
- Integration example: `/docs/gantt-integration-example.vue`
- Store documentation: `/src/stores/undoRedoStore.js`
- Event bus: `/src/utils/eventBus.js`

## Support

For questions or issues, contact the development team.
