import { createPinia, setActivePinia } from 'pinia'
import { config } from '@vue/test-utils'
import { vi } from 'vitest'

/**
 * Setup Pinia store for testing
 */
export function setupPinia() {
  setActivePinia(createPinia())
}

/**
 * Create a mock store with state
 */
export function createMockStore(storeClass, initialState = {}) {
  setupPinia()
  const store = storeClass()
  Object.assign(store, initialState)
  return store
}

/**
 * Wait for Vue's next tick
 */
export async function flushPromises() {
  return new Promise(resolve => setTimeout(resolve, 0))
}

/**
 * Mock Axios response
 */
export function mockAxiosResponse(data, status = 200) {
  return {
    data,
    status,
    statusText: 'OK',
    headers: {},
    config: {},
  }
}

/**
 * Mock Axios error
 */
export function mockAxiosError(message, status = 400) {
  const error = new Error(message) as any
  error.response = {
    data: { error: message },
    status,
    statusText: 'Error',
    headers: {},
    config: {},
  }
  return error
}

/**
 * Create mock task
 */
export function createMockTask(overrides = {}) {
  return {
    id: `task-${Date.now()}`,
    name: 'Test Task',
    start: new Date().toISOString(),
    end: new Date(Date.now() + 86400000).toISOString(),
    progress: 0,
    type: 'task',
    ...overrides,
  }
}

/**
 * Create mock dependency
 */
export function createMockDependency(overrides = {}) {
  return {
    id: `dep-${Date.now()}`,
    from: 'task-1',
    to: 'task-2',
    type: 'finish-to-start',
    lag: 0,
    ...overrides,
  }
}

/**
 * Create mock resource
 */
export function createMockResource(overrides = {}) {
  return {
    id: `res-${Date.now()}`,
    name: 'Test Resource',
    type: 'work',
    capacity: 8,
    unit: 'hours',
    ...overrides,
  }
}

/**
 * Create mock calendar
 */
export function createMockCalendar(overrides = {}) {
  return {
    id: `cal-${Date.now()}`,
    name: 'Standard',
    workingDays: [1, 2, 3, 4, 5],
    workingHours: { start: '08:00', end: '17:00' },
    holidays: [],
    ...overrides,
  }
}

/**
 * Create mock project
 */
export function createMockProject(overrides = {}) {
  return {
    id: `proj-${Date.now()}`,
    name: 'Test Project',
    start: new Date().toISOString(),
    end: new Date(Date.now() + 30 * 86400000).toISOString(),
    calendarId: 'cal-1',
    ...overrides,
  }
}

/**
 * Mock WebSocket
 */
export function createMockWebSocket() {
  const ws = {
    send: vi.fn(),
    close: vi.fn(),
    onopen: null,
    onmessage: null,
    onerror: null,
    onclose: null,
    readyState: 1, // OPEN
  }

  // Simulate receiving a message
  ws.simulateMessage = (data) => {
    if (ws.onmessage) {
      ws.onmessage({ data: JSON.stringify(data) })
    }
  }

  // Simulate connection opening
  ws.simulateOpen = () => {
    ws.readyState = 1
    if (ws.onopen) ws.onopen({})
  }

  // Simulate connection closing
  ws.simulateClose = () => {
    ws.readyState = 3
    if (ws.onclose) ws.onclose({})
  }

  // Simulate error
  ws.simulateError = (error) => {
    if (ws.onerror) ws.onerror(error)
  }

  return ws
}

/**
 * Generate array of tasks
 */
export function generateTasks(count: number, overrides = {}) {
  return Array.from({ length: count }, (_, i) =>
    createMockTask({
      id: `task-${i}`,
      name: `Task ${i + 1}`,
      start: new Date(Date.now() + i * 86400000).toISOString(),
      end: new Date(Date.now() + (i + 1) * 86400000).toISOString(),
      ...overrides,
    })
  )
}

/**
 * Generate array of dependencies
 */
export function generateDependencies(count: number, taskCount: number) {
  return Array.from({ length: count }, (_, i) => ({
    id: `dep-${i}`,
    from: `task-${i % taskCount}`,
    to: `task-${(i + 1) % taskCount}`,
    type: 'finish-to-start',
    lag: 0,
  }))
}

/**
 * Mock performance API
 */
export function mockPerformance() {
  return {
    now: vi.fn(() => Date.now()),
    mark: vi.fn(),
    measure: vi.fn(),
  }
}

/**
 * Create mock element with dimensions
 */
export function createMockElement(dimensions = { width: 100, height: 100 }) {
  return {
    getBoundingClientRect: vi.fn(() => ({
      width: dimensions.width,
      height: dimensions.height,
      top: 0,
      left: 0,
      right: dimensions.width,
      bottom: dimensions.height,
      x: 0,
      y: 0,
      toJSON: () => ({}),
    })),
    offsetWidth: dimensions.width,
    offsetHeight: dimensions.height,
    clientWidth: dimensions.width,
    clientHeight: dimensions.height,
    scrollWidth: dimensions.width,
    scrollHeight: dimensions.height,
  }
}

/**
 * Mock date-fns functions
 */
export function mockDateUtils() {
  const mockDate = new Date('2024-01-01T00:00:00Z')
  vi.spyOn(global, 'Date').mockImplementation(() => mockDate as any)
  return mockDate
}

/**
 * Test helper to check if two arrays are equal regardless of order
 */
export function arraysEqualIgnoreOrder(arr1: any[], arr2: any[]) {
  if (arr1.length !== arr2.length) return false
  const sorted1 = [...arr1].sort()
  const sorted2 = [...arr2].sort()
  return sorted1.every((val, index) => val === sorted2[index])
}

/**
 * Mock router
 */
export function createMockRouter() {
  return {
    push: vi.fn(),
    replace: vi.fn(),
    go: vi.fn(),
    back: vi.fn(),
    forward: vi.fn(),
    currentRoute: {
      value: {
        path: '/',
        name: 'home',
        params: {},
        query: {},
        meta: {},
      },
    },
  }
}

/**
 * Mock route
 */
export function createMockRoute(overrides = {}) {
  return {
    path: '/',
    name: 'home',
    params: {},
    query: {},
    meta: {},
    ...overrides,
  }
}
