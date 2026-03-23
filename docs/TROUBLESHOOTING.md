# Gantt Chart Editor - Troubleshooting Guide

**Version:** 1.0.0
**Last Updated:** 2026-02-19
**Author:** Material Management System Team

---

## Table of Contents

1. [Build Errors](#build-errors)
2. [Runtime Errors](#runtime-errors)
3. [Performance Issues](#performance-issues)
4. [WebSocket Connection Issues](#websocket-connection-issues)
5. [Testing Failures](#testing-failures)
6. [Browser-Specific Issues](#browser-specific-issues)
7. [Data Issues](#data-issues)
8. [UI/UX Issues](#uiux-issues)

---

## Build Errors

### Error: Module not found: Error: Can't resolve 'vue-virtual-scroller'

**Cause:** vue-virtual-scroller is not installed.

**Solution:**

```bash
cd /home/julei/backend/newstatic
npm install vue-virtual-scroller@2.0.0-beta.8
```

**Verify installation:**

```bash
npm list vue-virtual-scroller
```

---

### Error: Failed to resolve import: @/components/gantt

**Cause:** Import path alias is not configured.

**Solution:**

Check `vite.config.js`:

```javascript
import { defineConfig } from 'vite'
import path from 'path'

export default defineConfig({
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src')
    }
  }
})
```

---

### Error: Unexpected token <template>

**Cause:** Vue files are not being processed correctly.

**Solution:**

Check `vite.config.js` has Vue plugin:

```javascript
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()]
})
```

---

### Error: Cannot find module './components/gantt/index.js'

**Cause:** Component barrel export file doesn't exist.

**Solution:**

Create the file at `/src/components/gantt/index.js`:

```javascript
export { default as GanttEditor } from './core/GanttEditor.vue'
export { default as GanttToolbar } from './core/GanttToolbar.vue'
// ... more exports
```

---

### Error: TypeError: Cannot read property 'install' of undefined

**Cause:** Plugin is not imported correctly in main.js.

**Solution:**

Check main.js imports:

```javascript
// WRONG
import VueTour from 'vue-tour/dist/vue-tour'

// CORRECT
import VueTour from 'vue-tour'
```

---

## Runtime Errors

### Error: Maximum recursive updates exceeded

**Cause:** Watcher is triggering its own update, creating infinite loop.

**Example Problem:**

```javascript
// BAD: Causes infinite loop
watch(tasks, (newTasks) => {
  tasks.value = sortTasks(newTasks) // Triggers watch again
})
```

**Solution:**

```javascript
// GOOD: Use computed instead
const sortedTasks = computed(() => sortTasks(tasks.value))

// OR: Stop watching before updating
watch(tasks, (newTasks) => {
  stop() // Stop watching
  tasks.value = sortTasks(newTasks)
}, { flush: 'sync' })
```

---

### Error: Cannot read property 'id' of undefined

**Cause:** Task is undefined when accessing properties.

**Solution:**

Add null checks:

```javascript
// BAD
const taskName = task.name

// GOOD
const taskName = task?.name || 'Untitled'

// OR
const taskName = task && task.name
```

Use optional chaining in templates:

```vue
<template>
  <div>{{ task?.name || 'Unnamed Task' }}</div>
</template>
```

---

### Error: ResizeObserver loop limit exceeded

**Cause:** ResizeObserver is triggering more resize events than browser can handle.

**Solution:**

Debounce resize handler:

```javascript
import { debounce } from 'lodash-es'

const resizeObserver = new ResizeObserver(
  debounce((entries) => {
    // Handle resize
  }, 100)
)
```

---

### Error: Failed to execute 'appendChild' on 'Node'

**Cause:** Trying to append element that's already in DOM.

**Solution:**

Check if element exists before appending:

```javascript
if (!element.parentNode) {
  container.appendChild(element)
}
```

---

### Error: Unexpected token < in JSON at position 0

**Cause:** API returned HTML instead of JSON (likely error page).

**Solution:**

Check API response:

```javascript
try {
  const response = await fetch('/api/tasks')
  const data = await response.json()

  if (!response.ok) {
    throw new Error(data.message || 'API error')
  }
} catch (error) {
  console.error('API Error:', error)
  // Check if response is HTML
  if (error.message.includes('<')) {
    console.error('Server returned HTML instead of JSON')
  }
}
```

---

## Performance Issues

### Issue: Gantt chart is slow/laggy with 100+ tasks

**Symptoms:**
- Low FPS when scrolling
- Delayed UI response
- Browser uses 100% CPU

**Diagnosis:**

```javascript
// Check if virtual scrolling is working
console.log('Rendered task bars:', document.querySelectorAll('.task-bar').length)
console.log('Total tasks:', tasks.length)

// If numbers are close, virtual scrolling is not working
```

**Solution 1: Enable Virtual Scrolling**

```vue
<template>
  <RecycleScroller
    :items="tasks"
    :item-size="50"
    key-field="id"
  >
    <template #default="{ item }">
      <TaskRow :task="item" />
    </template>
  </RecycleScroller>
</template>
```

**Solution 2: Reduce Complexity**

```javascript
// BAD: Expensive computation on every render
const filteredTasks = tasks.filter(t => {
  return complexExpensiveCheck(t)
})

// GOOD: Memoize with computed
const filteredTasks = computed(() => {
  return tasks.value.filter(t => simpleCheck(t))
})
```

**Solution 3: Use pagination for 1000+ tasks**

```javascript
const loadTasksPaginated = async (page = 1) => {
  const response = await api.getTasks(projectId, {
    page,
    limit: 100
  })
  return response.data
}
```

---

### Issue: Memory usage keeps increasing

**Symptoms:**
- Browser tab uses > 500MB memory
- Memory keeps growing over time
- Eventually crashes

**Diagnosis:**

```javascript
// Monitor memory
setInterval(() => {
  if (performance.memory) {
    console.log('Memory:', {
      used: (performance.memory.usedJSHeapSize / 1048576).toFixed(2) + 'MB',
      total: (performance.memory.totalJSHeapSize / 1048576).toFixed(2) + 'MB'
    })
  }
}, 10000)
```

**Common Causes and Solutions:**

#### 1. Event Listeners Not Cleaned Up

```javascript
// BAD
onMounted(() => {
  window.addEventListener('resize', handleResize)
})

// GOOD
onMounted(() => {
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
})
```

#### 2. Timers Not Cleared

```javascript
// BAD
setInterval(() => {
  saveData()
}, 5000)

// GOOD
const timer = setInterval(() => {
  saveData()
}, 5000)

onUnmounted(() => {
  clearInterval(timer)
})
```

#### 3. Large Object References

```javascript
// BAD: Keeps reference to entire array
const tasks = ref([])
const archivedTasks = ref([])

function archiveTask(task) {
  archivedTasks.value.push(task) // Still in tasks array
  tasks.value = tasks.value.filter(t => t.id !== task.id)
}

// GOOD: Create new object
function archiveTask(task) {
  archivedTasks.value.push({ ...task })
  tasks.value = tasks.value.filter(t => t.id !== task.id)
}
```

---

### Issue: Initial load takes > 5 seconds

**Symptoms:**
- Blank screen for several seconds
- White screen before content appears
- Loader shows for too long

**Solution 1: Add Loading State**

```vue
<template>
  <div v-if="loading" class="loading">
    <el-icon class="is-loading"><Loading /></el-icon>
    <span>Loading tasks...</span>
  </div>
  <GanttEditor v-else :tasks="tasks" />
</template>

<script setup>
import { ref, onMounted } from 'vue'

const loading = ref(true)
const tasks = ref([])

onMounted(async () => {
  try {
    tasks.value = await loadTasks()
  } finally {
    loading.value = false
  }
})
</script>
```

**Solution 2: Progressive Loading**

```javascript
// Load critical data first
const [tasks, dependencies] = await Promise.all([
  loadTasks(),
  loadDependencies()
])

// Render initial UI
renderGantt(tasks, dependencies)

// Load optional data in background
setTimeout(async () => {
  const comments = await loadComments()
  const history = await loadHistory()
  // Update UI with additional data
}, 100)
```

**Solution 3: Code Splitting**

```javascript
// Instead of importing all at once
import { GanttEditor, KanbanView, CalendarView } from '@/components/gantt'

// Load on demand
const GanttEditor = defineAsyncComponent(() =>
  import('@/components/gantt/core/GanttEditor.vue')
)
```

---

## WebSocket Connection Issues

### Issue: WebSocket connection fails

**Symptoms:**
- "WebSocket connection failed" error
- Real-time updates not working
- Cursor positions not showing

**Diagnosis:**

```javascript
const socket = io('ws://localhost:8080')

socket.on('connect_error', (error) => {
  console.error('WebSocket error:', error)
})

socket.on('disconnect', (reason) => {
  console.log('WebSocket disconnected:', reason)
})
```

**Solution 1: Check Server is Running**

```bash
# Check if WebSocket server is listening
netstat -an | grep 8080
```

**Solution 2: Check CORS Configuration**

```javascript
// Server-side
const io = require('socket.io')(server, {
  cors: {
    origin: "http://localhost:5173",
    methods: ["GET", "POST"]
  }
})
```

**Solution 3: Check Firewall/Proxy**

Ensure WebSocket traffic is allowed:

```nginx
# Nginx configuration
location /ws/ {
    proxy_pass http://localhost:8080;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "upgrade";
}
```

**Solution 4: Implement Auto-Reconnection**

```javascript
const socket = io('ws://localhost:8080', {
  reconnection: true,
  reconnectionDelay: 1000,
  reconnectionAttempts: 10,
  reconnectionDelayMax: 5000
})

socket.on('reconnect', (attemptNumber) => {
  console.log('Reconnected after', attemptNumber, 'attempts')
})
```

---

### Issue: WebSocket messages not received

**Symptoms:**
- Connection established but no messages
- Other users' updates not showing

**Solution 1: Check Event Names**

```javascript
// Client
socket.on('task_updated', (data) => {
  console.log('Task updated:', data)
})

// Server
io.emit('task_updated', { taskId: 'task-1' })

// Make sure event names match exactly
```

**Solution 2: Check Room Subscription**

```javascript
// Client
socket.emit('join', { projectId: 'proj-123' })

// Server
socket.on('join', (data) => {
  socket.join(`project:${data.projectId}`)
})

// Emit to room
io.to(`project:${projectId}`).emit('task_updated', data)
```

**Solution 3: Add Logging**

```javascript
// Client-side logging
socket.on('message', (data) => {
  console.log('Received:', data)
})

// Server-side logging
io.on('connection', (socket) => {
  console.log('Client connected:', socket.id)

  socket.on('message', (data) => {
    console.log('Message from client:', socket.id, data)
  })
})
```

---

## Testing Failures

### Issue: Component tests fail with "Cannot find module"

**Cause:** Vitest doesn't handle Vue components correctly.

**Solution:**

Update `vitest.config.js`:

```javascript
import { defineConfig } from 'vitest/config'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  test: {
    environment: 'jsdom',
    globals: true
  },
  resolve: {
    alias: {
      '@': '/src'
    }
  }
})
```

---

### Issue: Test timeout - component takes too long to mount

**Cause:** Component has async operations or heavy computations.

**Solution:**

Increase timeout or use `flush-promises`:

```javascript
import { flushPromises } from '@vue/test-utils'

test('renders tasks', async () => {
  const wrapper = mount(GanttEditor, {
    props: { tasks: mockTasks }
  })

  // Wait for async operations
  await flushPromises()
  await nextTick()

  expect(wrapper.find('.task-bar').exists()).toBe(true)
})
```

---

### Issue: E2E test fails - element not found

**Cause:** Element not yet rendered when test looks for it.

**Solution:**

Add explicit wait:

```javascript
test('should display task bars', async ({ page }) => {
  await page.goto('/progress/gantt?projectId=proj-123')

  // Wait for selector
  await page.waitForSelector('.task-bar', { timeout: 5000 })

  const count = await page.locator('.task-bar').count()
  expect(count).toBeGreaterThan(0)
})
```

---

### Issue: Tests pass locally but fail in CI

**Cause:** CI environment has different timing or configuration.

**Solution:**

Add retries and longer timeouts:

```javascript
// playwright.config.js
export default defineConfig({
  retries: process.env.CI ? 2 : 0,
  timeout: 60000,
  use: {
    actionTimeout: 10000
  }
})
```

---

## Browser-Specific Issues

### Issue: Safari shows blank screen

**Cause:** Safari doesn't support some CSS features.

**Solution:**

Add vendor prefixes:

```css
.task-bar {
  -webkit-transform: translate(0, 0);
  -moz-transform: translate(0, 0);
  -ms-transform: translate(0, 0);
  transform: translate(0, 0);
}
```

---

### Issue: Firefox shows horizontal scrollbar

**Cause:** Firefox calculates box model differently.

**Solution:**

```css
.gantt-container {
  box-sizing: border-box;
  width: 100%;
  overflow-x: hidden;
}
```

---

### Issue: Mobile Safari has touch issues

**Cause:** Touch events not handled correctly.

**Solution:**

```javascript
// Add touch event support
const handleTouchStart = (e) => {
  // Prevent default zoom
  if (e.touches.length > 1) {
    e.preventDefault()
  }
}

const container = ref(null)
onMounted(() => {
  container.value?.addEventListener('touchstart', handleTouchStart, {
    passive: false
  })
})
```

---

## Data Issues

### Issue: Tasks appear in wrong order

**Cause:** Sort order not applied or position field incorrect.

**Solution:**

```javascript
// Ensure tasks are sorted
const sortedTasks = computed(() => {
  return [...tasks.value].sort((a, b) => a.position - b.position)
})
```

---

### Issue: Dependencies show circular reference error

**Cause:** Task A depends on B, B depends on A.

**Solution:**

```javascript
function detectCircularDependency(dependencies, taskId, visited = new Set()) {
  if (visited.has(taskId)) {
    throw new Error('Circular dependency detected')
  }

  visited.add(taskId)

  const dependentTasks = dependencies
    .filter(d => d.fromTaskId === taskId)
    .map(d => d.toTaskId)

  for (const depId of dependentTasks) {
    detectCircularDependency(dependencies, depId, visited)
  }
}

// Check before adding dependency
try {
  detectCircularDependency(dependencies, fromTaskId)
  addDependency(fromTaskId, toTaskId)
} catch (error) {
  showError(error.message)
}
```

---

### Issue: Dates are off by one day

**Cause:** Timezone conversion issue.

**Solution:**

```javascript
// Always work in UTC
function parseDate(dateStr) {
  const date = new Date(dateStr)
  return new Date(Date.UTC(
    date.getFullYear(),
    date.getMonth(),
    date.getDate()
  ))
}

// Or use date-fns with timezone
import { utcToZonedTime, zonedTimeToUtc } from 'date-fns-tz'

const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone
const zonedDate = utcToZonedTime(date, timezone)
```

---

## UI/UX Issues

### Issue: Task bars overlap

**Cause:** Tasks have same position/row.

**Solution:**

```javascript
function calculateTaskPositions(tasks) {
  const levels = []

  tasks.forEach(task => {
    let level = 0
    let found = false

    while (!found) {
      if (!levels[level]) {
        levels[level] = []
      }

      const canFit = levels[level].every(existingTask =>
        !doTasksOverlap(task, existingTask)
      )

      if (canFit) {
        levels[level].push(task)
        task.row = level
        found = true
      } else {
        level++
      }
    }
  })

  return tasks
}
```

---

### Issue: Text is cut off in task bars

**Cause:** Task bar width is too small for text.

**Solution:**

```vue
<template>
  <div
    class="task-bar"
    :style="{ width: taskWidth + 'px' }"
  >
    <span
      v-if="taskWidth > 50"
      class="task-name"
    >
      {{ task.name }}
    </span>
    <span v-else class="task-name-short">
      {{ task.name.charAt(0) }}
    </span>
  </div>
</template>
```

---

### Issue: Tooltip flickers

**Cause:** Mouse events triggering rapid show/hide.

**Solution:**

```javascript
const tooltipVisible = ref(false)
const tooltipDelay = ref(null)

const showTooltip = () => {
  clearTimeout(tooltipDelay.value)
  tooltipDelay.value = setTimeout(() => {
    tooltipVisible.value = true
  }, 200) // Wait 200ms before showing
}

const hideTooltip = () => {
  clearTimeout(tooltipDelay.value)
  tooltipVisible.value = false
}
```

---

## Getting Help

If you can't resolve your issue:

1. **Check Documentation:**
   - [MAIN_INTEGRATION_GUIDE.md](./MAIN_INTEGRATION_GUIDE.md)
   - [COMPONENT_REFERENCE.md](./COMPONENT_REFERENCE.md)
   - [API_ENDPOINTS.md](./API_ENDPOINTS.md)

2. **Search Existing Issues:**
   - GitHub Issues
   - Stack Overflow
   - Vue Forums

3. **Enable Debug Mode:**

```javascript
// main.js
app.config.devtools = true
app.config.performance = true

// Enable Vue warnings
if (import.meta.env.DEV) {
  const app = createApp(App)
  app.config.warnHandler = (msg, instance, trace) => {
    console.warn('Vue warning:', msg, trace)
  }
}
```

4. **Create Minimal Reproduction:**

```javascript
// Isolate the problem
import { mount } from '@vue/test-utils'

const wrapper = mount(Component, {
  props: { ...minimalProps }
})

// Test specific behavior
```

5. **Collect Diagnostic Information:**

```javascript
// Browser info
console.log('User Agent:', navigator.userAgent)
console.log('Screen:', screen.width, 'x', screen.height)

// Vue info
console.log('Vue Version:', version)
console.log('App Config:', app.config)

// Performance
console.log('Memory:', performance.memory)
console.log('Navigation:', performance.getEntriesByType('navigation'))
```

---

## Quick Fixes Checklist

Before diving deep into debugging, try these quick fixes:

- [ ] Clear browser cache and reload (Ctrl+Shift+R / Cmd+Shift+R)
- [ ] Restart dev server (`npm run dev`)
- [ ] Delete `node_modules` and reinstall (`rm -rf node_modules && npm install`)
- [ ] Check browser console for errors (F12 → Console)
- [ ] Check Network tab for failed requests (F12 → Network)
- [ ] Try in incognito/private mode
- [ ] Try different browser
- [ ] Update dependencies (`npm update`)
- [ ] Check if backend is running
- [ ] Check database connection
- [ ] Verify API endpoints are accessible
- [ ] Check CORS configuration
- [ ] Review recent code changes
- [ ] Rollback to last working version if needed

---

**Document Version:** 1.0.0
**Last Updated:** 2026-02-19
