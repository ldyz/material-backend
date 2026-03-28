import { vi } from 'vitest'
import { config } from '@vue/test-utils'
import ElementPlus from 'element-plus'

// Global setup for Vue Test Utils
config.global.plugins = [ElementPlus]

// Mock window.matchMedia
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation(query => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
})

// Mock ResizeObserver
global.ResizeObserver = vi.fn().mockImplementation(() => ({
  observe: vi.fn(),
  unobserve: vi.fn(),
  disconnect: vi.fn(),
}))

// Mock IntersectionObserver
global.IntersectionObserver = vi.fn().mockImplementation(() => ({
  observe: vi.fn(),
  unobserve: vi.fn(),
  disconnect: vi.fn(),
}))

// Mock localStorage
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
  length: 0,
  key: vi.fn(),
}
global.localStorage = localStorageMock

// Mock sessionStorage
const sessionStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
  length: 0,
  key: vi.fn(),
}
global.sessionStorage = sessionStorageMock

// Mock console methods in tests
global.console = {
  ...console,
  warn: vi.fn(),
  error: vi.fn(),
}

// Create calendarApi mock
const createCalendarApiMock = () => ({
  create: vi.fn().mockResolvedValue({ data: { id: 'cal-1', name: 'Test Calendar' } }),
  update: vi.fn().mockResolvedValue({ data: { id: 'cal-1', name: 'Updated Calendar' } }),
  delete: vi.fn().mockResolvedValue({}),
  addHoliday: vi.fn().mockResolvedValue({ data: { date: '2024-01-01', name: 'Holiday' } }),
  removeHoliday: vi.fn().mockResolvedValue({}),
  list: vi.fn().mockResolvedValue({ data: [] }),
  setProjectCalendar: vi.fn().mockResolvedValue({}),
  assignTaskCalendar: vi.fn().mockResolvedValue({}),
  removeTaskCalendar: vi.fn().mockResolvedValue({}),
  addException: vi.fn().mockResolvedValue({}),
  removeException: vi.fn().mockResolvedValue({})
})

// Mock API modules
vi.mock('@/api/index.js', () => ({
  authApi: {},
  userApi: {},
  materialApi: {},
  projectApi: {},
  stockApi: {},
  inboundApi: {},
  calendarApi: createCalendarApiMock()
}))

// Also mock for @/api alias
vi.mock('@/api', () => ({
  authApi: {},
  userApi: {},
  materialApi: {},
  projectApi: {},
  stockApi: {},
  inboundApi: {},
  calendarApi: createCalendarApiMock(),
  default: {
    authApi: {},
    userApi: {},
    materialApi: {},
    projectApi: {},
    stockApi: {},
    inboundApi: {},
    calendarApi: createCalendarApiMock()
  }
}))

// Reset mocks before each test
beforeEach(() => {
  vi.clearAllMocks()
  localStorageMock.clear()
  sessionStorageMock.clear()
})
