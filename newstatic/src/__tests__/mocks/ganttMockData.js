/**
 * Mock data for Gantt chart tests
 */

export const mockTasks = [
  {
    id: 'task-1',
    name: 'Project Planning',
    start: '2024-01-01T08:00:00Z',
    end: '2024-01-03T17:00:00Z',
    progress: 100,
    type: 'task',
    assignee: 'user-1',
    resources: ['res-1'],
    priority: 'high',
    status: 'completed',
    color: '#67C23A',
    dependsOn: [],
  },
  {
    id: 'task-2',
    name: 'Requirements Gathering',
    start: '2024-01-04T08:00:00Z',
    end: '2024-01-08T17:00:00Z',
    progress: 80,
    type: 'task',
    assignee: 'user-2',
    resources: ['res-1', 'res-2'],
    priority: 'high',
    status: 'in-progress',
    color: '#409EFF',
    dependsOn: ['task-1'],
  },
  {
    id: 'task-3',
    name: 'Design Phase',
    start: '2024-01-09T08:00:00Z',
    end: '2024-01-15T17:00:00Z',
    progress: 50,
    type: 'task',
    assignee: 'user-3',
    resources: ['res-2'],
    priority: 'medium',
    status: 'in-progress',
    color: '#E6A23C',
    dependsOn: ['task-2'],
  },
  {
    id: 'task-4',
    name: 'Development',
    start: '2024-01-16T08:00:00Z',
    end: '2024-01-30T17:00:00Z',
    progress: 20,
    type: 'task',
    assignee: 'user-1',
    resources: ['res-1', 'res-2', 'res-3'],
    priority: 'high',
    status: 'in-progress',
    color: '#F56C6C',
    dependsOn: ['task-3'],
  },
  {
    id: 'task-5',
    name: 'Testing',
    start: '2024-01-31T08:00:00Z',
    end: '2024-02-05T17:00:00Z',
    progress: 0,
    type: 'task',
    assignee: 'user-4',
    resources: ['res-4'],
    priority: 'medium',
    status: 'pending',
    color: '#909399',
    dependsOn: ['task-4'],
  },
  {
    id: 'task-6',
    name: 'Deployment',
    start: '2024-02-06T08:00:00Z',
    end: '2024-02-07T17:00:00Z',
    progress: 0,
    type: 'milestone',
    assignee: 'user-2',
    resources: [],
    priority: 'high',
    status: 'pending',
    color: '#67C23A',
    dependsOn: ['task-5'],
  },
]

export const mockDependencies = [
  {
    id: 'dep-1',
    from: 'task-1',
    to: 'task-2',
    type: 'finish-to-start',
    lag: 0,
  },
  {
    id: 'dep-2',
    from: 'task-2',
    to: 'task-3',
    type: 'finish-to-start',
    lag: 0,
  },
  {
    id: 'dep-3',
    from: 'task-3',
    to: 'task-4',
    type: 'start-to-start',
    lag: 1,
  },
  {
    id: 'dep-4',
    from: 'task-4',
    to: 'task-5',
    type: 'finish-to-finish',
    lag: 0,
  },
  {
    id: 'dep-5',
    from: 'task-5',
    to: 'task-6',
    type: 'finish-to-start',
    lag: -1,
  },
]

export const mockResources = [
  {
    id: 'res-1',
    name: 'John Developer',
    type: 'work',
    capacity: 8,
    unit: 'hours',
    costPerHour: 50,
    available: true,
  },
  {
    id: 'res-2',
    name: 'Jane Designer',
    type: 'work',
    capacity: 8,
    unit: 'hours',
    costPerHour: 45,
    available: true,
  },
  {
    id: 'res-3',
    name: 'Bob Tester',
    type: 'work',
    capacity: 6,
    unit: 'hours',
    costPerHour: 40,
    available: true,
  },
  {
    id: 'res-4',
    name: 'Meeting Room A',
    type: 'material',
    capacity: 1,
    unit: 'rooms',
    costPerHour: 20,
    available: true,
  },
]

export const mockCalendars = [
  {
    id: 'cal-1',
    name: 'Standard',
    workingDays: [1, 2, 3, 4, 5],
    workingHours: { start: '08:00', end: '17:00' },
    holidays: [
      { date: '2024-01-01', name: 'New Year' },
      { date: '2024-01-15', name: 'MLK Day' },
    ],
  },
  {
    id: 'cal-2',
    name: 'Extended Hours',
    workingDays: [1, 2, 3, 4, 5, 6],
    workingHours: { start: '07:00', end: '19:00' },
    holidays: [],
  },
]

export const mockConstraints = [
  {
    id: 'con-1',
    taskId: 'task-4',
    type: 'must-start-on',
    date: '2024-01-16T08:00:00Z',
    applied: true,
  },
  {
    id: 'con-2',
    taskId: 'task-6',
    type: 'must-finish-by',
    date: '2024-02-07T17:00:00Z',
    applied: true,
  },
  {
    id: 'con-3',
    taskId: 'task-3',
    type: 'start-no-earlier-than',
    date: '2024-01-09T08:00:00Z',
    applied: true,
  },
]

export const mockProject = {
  id: 'proj-1',
  name: 'Sample Project',
  start: '2024-01-01T08:00:00Z',
  end: '2024-02-07T17:00:00Z',
  calendarId: 'cal-1',
  status: 'in-progress',
  progress: 45,
}

export const mockUsers = [
  {
    id: 'user-1',
    name: 'John Doe',
    email: 'john@example.com',
    avatar: 'avatar-1.jpg',
    role: 'developer',
  },
  {
    id: 'user-2',
    name: 'Jane Smith',
    email: 'jane@example.com',
    avatar: 'avatar-2.jpg',
    role: 'designer',
  },
  {
    id: 'user-3',
    name: 'Bob Johnson',
    email: 'bob@example.com',
    avatar: 'avatar-3.jpg',
    role: 'developer',
  },
  {
    id: 'user-4',
    name: 'Alice Williams',
    email: 'alice@example.com',
    avatar: 'avatar-4.jpg',
    role: 'tester',
  },
]

export const mockComments = [
  {
    id: 'comment-1',
    taskId: 'task-1',
    userId: 'user-1',
    userName: 'John Doe',
    content: 'Task completed successfully',
    timestamp: '2024-01-03T10:00:00Z',
    parentId: null,
    mentions: ['user-2'],
    attachments: [],
  },
  {
    id: 'comment-2',
    taskId: 'task-2',
    userId: 'user-2',
    userName: 'Jane Smith',
    content: 'Working on requirements',
    timestamp: '2024-01-05T14:30:00Z',
    parentId: null,
    mentions: [],
    attachments: [],
  },
]

export const mockChanges = [
  {
    id: 'change-1',
    timestamp: '2024-01-15T10:00:00Z',
    userId: 'user-1',
    userName: 'John Doe',
    action: 'update',
    entityType: 'task',
    entityId: 'task-1',
    changes: {
      progress: { from: 80, to: 100 },
      status: { from: 'in-progress', to: 'completed' },
    },
  },
  {
    id: 'change-2',
    timestamp: '2024-01-14T14:30:00Z',
    userId: 'user-2',
    userName: 'Jane Smith',
    action: 'create',
    entityType: 'dependency',
    entityId: 'dep-5',
    changes: {},
  },
]

export const largeTaskSet = Array.from({ length: 100 }, (_, i) => ({
  id: `task-${i}`,
  name: `Task ${i + 1}`,
  start: new Date(Date.now() + i * 86400000).toISOString(),
  end: new Date(Date.now() + (i + 1) * 86400000).toISOString(),
  progress: Math.floor(Math.random() * 100),
  type: i % 10 === 0 ? 'milestone' : 'task',
  assignee: `user-${(i % 4) + 1}`,
  resources: [`res-${(i % 4) + 1}`],
  priority: ['high', 'medium', 'low'][i % 3],
  status: ['pending', 'in-progress', 'completed'][i % 3],
  dependsOn: i > 0 ? [`task-${i - 1}`] : [],
}))

export const circularDependencyTasks = [
  {
    id: 'task-1',
    name: 'Task 1',
    start: '2024-01-01T08:00:00Z',
    end: '2024-01-03T17:00:00Z',
    progress: 0,
    type: 'task',
    dependsOn: ['task-3'],
  },
  {
    id: 'task-2',
    name: 'Task 2',
    start: '2024-01-04T08:00:00Z',
    end: '2024-01-06T17:00:00Z',
    progress: 0,
    type: 'task',
    dependsOn: ['task-1'],
  },
  {
    id: 'task-3',
    name: 'Task 3',
    start: '2024-01-07T08:00:00Z',
    end: '2024-01-09T17:00:00Z',
    progress: 0,
    type: 'task',
    dependsOn: ['task-2'],
  },
]

export const overAllocatedResources = [
  {
    id: 'task-1',
    name: 'Task 1',
    start: '2024-01-01T08:00:00Z',
    end: '2024-01-05T17:00:00Z',
    progress: 0,
    type: 'task',
    resources: ['res-1'],
    assignee: 'user-1',
  },
  {
    id: 'task-2',
    name: 'Task 2',
    start: '2024-01-01T08:00:00Z',
    end: '2024-01-05T17:00:00Z',
    progress: 0,
    type: 'task',
    resources: ['res-1'],
    assignee: 'user-1',
  },
  {
    id: 'task-3',
    name: 'Task 3',
    start: '2024-01-02T08:00:00Z',
    end: '2024-01-06T17:00:00Z',
    progress: 0,
    type: 'task',
    resources: ['res-1'],
    assignee: 'user-1',
  },
]

export const mockKanbanColumns = [
  {
    id: 'col-1',
    title: 'To Do',
    status: 'pending',
    wipLimit: 5,
    taskIds: ['task-5'],
    order: 0,
  },
  {
    id: 'col-2',
    title: 'In Progress',
    status: 'in-progress',
    wipLimit: 3,
    taskIds: ['task-2', 'task-3', 'task-4'],
    order: 1,
  },
  {
    id: 'col-3',
    title: 'Completed',
    status: 'completed',
    wipLimit: null,
    taskIds: ['task-1'],
    order: 2,
  },
]

export const mockTemplates = [
  {
    id: 'tpl-1',
    name: 'Software Development',
    description: 'Standard software development lifecycle',
    tasks: [
      {
        name: 'Requirements',
        duration: 5,
        type: 'task',
        defaultAssignee: 'user-2',
      },
      {
        name: 'Design',
        duration: 7,
        type: 'task',
        defaultAssignee: 'user-3',
      },
      {
        name: 'Development',
        duration: 14,
        type: 'task',
        defaultAssignee: 'user-1',
      },
      {
        name: 'Testing',
        duration: 5,
        type: 'task',
        defaultAssignee: 'user-4',
      },
      {
        name: 'Deployment',
        duration: 1,
        type: 'milestone',
        defaultAssignee: 'user-2',
      },
    ],
    dependencies: [
      { from: 0, to: 1, type: 'finish-to-start' },
      { from: 1, to: 2, type: 'finish-to-start' },
      { from: 2, to: 3, type: 'finish-to-start' },
      { from: 3, to: 4, type: 'finish-to-start' },
    ],
  },
]

export default {
  mockTasks,
  mockDependencies,
  mockResources,
  mockCalendars,
  mockConstraints,
  mockProject,
  mockUsers,
  mockComments,
  mockChanges,
  largeTaskSet,
  circularDependencyTasks,
  overAllocatedResources,
  mockKanbanColumns,
  mockTemplates,
}
