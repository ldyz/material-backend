# Gantt Chart Editor - API Endpoints Reference

**Version:** 1.0.0
**Last Updated:** 2026-02-19
**Base URL:** `/api/v1/gantt`

---

## Table of Contents

1. [Overview](#overview)
2. [Authentication](#authentication)
3. [Task Endpoints](#task-endpoints)
4. [Dependency Endpoints](#dependency-endpoints)
5. [Constraint Endpoints](#constraint-endpoints)
6. [Comment Endpoints](#comment-endpoints)
7. [History Endpoints](#history-endpoints)
8. [Template Endpoints](#template-endpoints)
9. [Report Endpoints](#report-endpoints)
10. [Resource Leveling Endpoints](#resource-leveling-endpoints)
11. [AI Suggestion Endpoints](#ai-suggestion-endpoints)
12. [WebSocket Events](#websocket-events)
13. [Error Codes](#error-codes)
14. [Rate Limiting](#rate-limiting)

---

## Overview

The Gantt Chart API provides RESTful endpoints for managing project tasks, dependencies, constraints, and advanced features like resource leveling and AI-powered suggestions.

### Base URL Structure

```
Production: https://api.yourdomain.com/api/v1/gantt
Development: http://localhost:8080/api/v1/gantt
```

### Response Format

All responses follow this structure:

**Success Response:**
```json
{
  "success": true,
  "data": { ... },
  "message": "Operation successful"
}
```

**Error Response:**
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": { ... }
  }
}
```

---

## Authentication

All endpoints require authentication using Bearer token.

### Headers

```http
Authorization: Bearer <your-jwt-token>
Content-Type: application/json
```

### Example Request

```javascript
const response = await fetch('/api/v1/gantt/projects/proj-123/tasks', {
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  }
})
```

---

## Task Endpoints

### Get All Tasks

Retrieve all tasks for a project.

**Endpoint:** `GET /projects/:projectId/tasks`

**Path Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `projectId` | string | Yes | Project unique identifier |

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | number | 1 | Page number for pagination |
| `limit` | number | 100 | Items per page |
| `sort` | string | `position` | Sort field (`position`, `startDate`, `endDate`, `name`, `progress`) |
| `order` | string | `asc` | Sort order (`asc`, `desc`) |
| `status` | string | - | Filter by status (comma-separated) |
| `assignee` | string | - | Filter by assignee ID |
| `search` | string | - | Search in task name/description |

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "tasks": [
      {
        "id": "task-1",
        "projectId": "proj-123",
        "name": "Design System",
        "description": "Create design system components",
        "startDate": "2024-01-01T00:00:00Z",
        "endDate": "2024-01-10T00:00:00Z",
        "duration": 10,
        "progress": 75,
        "status": "in_progress",
        "priority": "high",
        "assignee": "user-1",
        "assigneeName": "John Doe",
        "parentId": null,
        "position": 0,
        "color": "#409EFF",
        "milestone": false,
        "constraint": null,
        "customFields": {},
        "createdAt": "2024-01-01T00:00:00Z",
        "updatedAt": "2024-01-05T00:00:00Z"
      }
    ],
    "total": 150,
    "page": 1,
    "limit": 100
  }
}
```

**Error Responses:**
- `401 Unauthorized` - Invalid or missing token
- `403 Forbidden` - No permission to access project
- `404 Not Found` - Project not found

---

### Create Task

Create a new task in a project.

**Endpoint:** `POST /projects/:projectId/tasks`

**Path Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `projectId` | string | Yes | Project unique identifier |

**Request Body:**

```json
{
  "name": "New Task",
  "description": "Task description",
  "startDate": "2024-01-15T00:00:00Z",
  "endDate": "2024-01-20T00:00:00Z",
  "duration": 5,
  "progress": 0,
  "status": "not_started",
  "priority": "medium",
  "assignee": "user-2",
  "parentId": null,
  "position": 1,
  "color": "#67C23A",
  "milestone": false,
  "customFields": {
    "storyPoints": 5,
    "tags": ["frontend", "ui"]
  }
}
```

**Response:** `201 Created`

```json
{
  "success": true,
  "data": {
    "task": {
      "id": "task-2",
      "projectId": "proj-123",
      "name": "New Task",
      "description": "Task description",
      "startDate": "2024-01-15T00:00:00Z",
      "endDate": "2024-01-20T00:00:00Z",
      "duration": 5,
      "progress": 0,
      "status": "not_started",
      "priority": "medium",
      "assignee": "user-2",
      "assigneeName": "Jane Smith",
      "parentId": null,
      "position": 1,
      "color": "#67C23A",
      "milestone": false,
      "constraint": null,
      "customFields": {
        "storyPoints": 5,
        "tags": ["frontend", "ui"]
      },
      "createdAt": "2024-01-10T00:00:00Z",
      "updatedAt": "2024-01-10T00:00:00Z"
    }
  },
  "message": "Task created successfully"
}
```

**Validation Errors:** `400 Bad Request`

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid task data",
    "details": {
      "name": ["Name is required"],
      "startDate": ["Start date must be before end date"]
    }
  }
}
```

---

### Update Task

Update an existing task.

**Endpoint:** `PATCH /projects/:projectId/tasks/:taskId`

**Path Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `projectId` | string | Yes | Project unique identifier |
| `taskId` | string | Yes | Task unique identifier |

**Request Body:** (all fields optional)

```json
{
  "name": "Updated Task Name",
  "progress": 50,
  "status": "in_progress",
  "endDate": "2024-01-25T00:00:00Z"
}
```

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "task": { ... },
    "changes": [
      {
        "field": "progress",
        "oldValue": 0,
        "newValue": 50
      },
      {
        "field": "status",
        "oldValue": "not_started",
        "newValue": "in_progress"
      }
    ]
  },
  "message": "Task updated successfully"
}
```

---

### Delete Task

Delete a task and its dependencies.

**Endpoint:** `DELETE /projects/:projectId/tasks/:taskId`

**Path Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `projectId` | string | Yes | Project unique identifier |
| `taskId` | string | Yes | Task unique identifier |

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `deleteChildren` | boolean | `false` | Also delete child tasks |
| `deleteDependencies` | boolean | `true` | Also delete associated dependencies |

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "deletedTaskId": "task-2",
    "deletedDependencies": 3,
    "deletedChildren": 0
  },
  "message": "Task deleted successfully"
}
```

---

### Bulk Update Tasks

Update multiple tasks at once.

**Endpoint:** `PATCH /projects/:projectId/tasks/bulk`

**Path Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `projectId` | string | Yes | Project unique identifier |

**Request Body:**

```json
{
  "updates": [
    {
      "taskId": "task-1",
      "progress": 100,
      "status": "completed"
    },
    {
      "taskId": "task-2",
      "progress": 50,
      "status": "in_progress"
    }
  ],
  "options": {
    "recalculateDates": true,
    "notifyAssignees": true
  }
}
```

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "updated": 2,
    "failed": 0,
    "tasks": [ ... ]
  },
  "message": "2 tasks updated successfully"
}
```

---

### Reorder Tasks

Change the order/position of tasks.

**Endpoint:** `POST /projects/:projectId/tasks/reorder`

**Request Body:**

```json
{
  "taskIds": ["task-3", "task-1", "task-2"],
  "parentId": null
}
```

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "updated": 3
  },
  "message": "Tasks reordered successfully"
}
```

---

## Dependency Endpoints

### Get All Dependencies

**Endpoint:** `GET /projects/:projectId/dependencies`

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "dependencies": [
      {
        "id": "dep-1",
        "projectId": "proj-123",
        "fromTaskId": "task-1",
        "toTaskId": "task-2",
        "type": "finish-to-start",
        "lag": 0,
        "createdAt": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

---

### Create Dependency

**Endpoint:** `POST /projects/:projectId/dependencies`

**Request Body:**

```json
{
  "fromTaskId": "task-1",
  "toTaskId": "task-2",
  "type": "finish-to-start",
  "lag": 0
}
```

**Dependency Types:**
- `finish-to-start` (FS) - Task must finish before next starts
- `start-to-start` (SS) - Tasks must start together
- `finish-to-finish` (FF) - Tasks must finish together
- `start-to-finish` (SF) - Task must start before next finishes

**Response:** `201 Created`

```json
{
  "success": true,
  "data": {
    "dependency": {
      "id": "dep-2",
      "fromTaskId": "task-1",
      "toTaskId": "task-2",
      "type": "finish-to-start",
      "lag": 0
    }
  }
}
```

---

### Delete Dependency

**Endpoint:** `DELETE /projects/:projectId/dependencies/:dependencyId`

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "deletedDependencyId": "dep-2"
  }
}
```

---

## Constraint Endpoints

### Get Constraints

**Endpoint:** `GET /projects/:projectId/constraints`

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "constraints": [
      {
        "taskId": "task-1",
        "type": "start-no-earlier-than",
        "date": "2024-01-15T00:00:00Z"
      }
    ]
  }
}
```

**Constraint Types:**
- `start-no-earlier-than` - Task cannot start before this date
- `finish-no-later-than` - Task must finish by this date
- `must-start-on` - Task must start exactly on this date
- `must-finish-on` - Task must finish exactly on this date

---

### Update Constraint

**Endpoint:** `PUT /projects/:projectId/tasks/:taskId/constraint`

**Request Body:**

```json
{
  "type": "start-no-earlier-than",
  "date": "2024-01-15T00:00:00Z"
}
```

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "constraint": {
      "taskId": "task-1",
      "type": "start-no-earlier-than",
      "date": "2024-01-15T00:00:00Z"
    }
  }
}
```

---

### Remove Constraint

**Endpoint:** `DELETE /projects/:projectId/tasks/:taskId/constraint`

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "removedConstraint": true
  }
}
```

---

## Comment Endpoints

### Get Comments

**Endpoint:** `GET /projects/:projectId/tasks/:taskId/comments`

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | number | 1 | Page number |
| `limit` | number | 50 | Items per page |

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "comments": [
      {
        "id": "comment-1",
        "taskId": "task-1",
        "userId": "user-1",
        "userName": "John Doe",
        "userAvatar": "/avatars/user-1.jpg",
        "content": "Task is progressing well",
        "createdAt": "2024-01-05T10:30:00Z",
        "updatedAt": "2024-01-05T10:30:00Z"
      }
    ],
    "total": 5
  }
}
```

---

### Add Comment

**Endpoint:** `POST /projects/:projectId/tasks/:taskId/comments`

**Request Body:**

```json
{
  "content": "This task needs more resources",
  "mentions": ["user-2", "user-3"]
}
```

**Response:** `201 Created`

```json
{
  "success": true,
  "data": {
    "comment": {
      "id": "comment-2",
      "taskId": "task-1",
      "userId": "user-1",
      "userName": "John Doe",
      "content": "This task needs more resources",
      "mentions": ["user-2", "user-3"],
      "createdAt": "2024-01-10T14:20:00Z"
    }
  }
}
```

---

### Delete Comment

**Endpoint:** `DELETE /projects/:projectId/tasks/:taskId/comments/:commentId`

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "deletedCommentId": "comment-2"
  }
}
```

---

## History Endpoints

### Get Task History

**Endpoint:** `GET /projects/:projectId/tasks/:taskId/history`

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | number | 1 | Page number |
| `limit` | number | 20 | Items per page |
| `fields` | string | - | Filter by fields (comma-separated) |

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "history": [
      {
        "id": "hist-1",
        "taskId": "task-1",
        "userId": "user-1",
        "userName": "John Doe",
        "timestamp": "2024-01-05T10:30:00Z",
        "changes": [
          {
            "field": "progress",
            "oldValue": 50,
            "newValue": 75
          },
          {
            "field": "status",
            "oldValue": "in_progress",
            "newValue": "completed"
          }
        ],
        "ipAddress": "192.168.1.100",
        "userAgent": "Mozilla/5.0..."
      }
    ],
    "total": 15
  }
}
```

---

### Restore from History

**Endpoint:** `POST /projects/:projectId/tasks/:taskId/history/:historyId/restore`

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "task": { ... },
    "restoredFromHistoryId": "hist-1"
  },
  "message": "Task restored from history"
}
```

---

## Template Endpoints

### Get Templates

**Endpoint:** `GET /templates`

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `category` | string | - | Filter by category |
| `search` | string | - | Search in name/description |

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "templates": [
      {
        "id": "tmpl-1",
        "name": "Software Development Sprint",
        "description": "Standard 2-week sprint template",
        "category": "software-development",
        "tasks": [
          {
            "name": "Planning",
            "duration": 2,
            "position": 0
          },
          {
            "name": "Development",
            "duration": 7,
            "position": 1
          },
          {
            "name": "Testing",
            "duration": 3,
            "position": 2
          },
          {
            "name": "Deployment",
            "duration": 1,
            "position": 3
          }
        ],
        "dependencies": [
          {
            "from": 0,
            "to": 1,
            "type": "finish-to-start"
          }
        ],
        "createdAt": "2024-01-01T00:00:00Z",
        "createdBy": "user-1"
      }
    ]
  }
}
```

---

### Create Template

**Endpoint:** `POST /templates`

**Request Body:**

```json
{
  "name": "Marketing Campaign",
  "description": "Standard marketing campaign template",
  "category": "marketing",
  "tasks": [ ... ],
  "dependencies": [ ... ]
}
```

**Response:** `201 Created`

```json
{
  "success": true,
  "data": {
    "template": {
      "id": "tmpl-2",
      "name": "Marketing Campaign",
      ...
    }
  }
}
```

---

### Apply Template

**Endpoint:** `POST /projects/:projectId/apply-template/:templateId`

**Request Body:**

```json
{
  "startDate": "2024-02-01T00:00:00Z",
  "assignee": "user-1",
  "options": {
    "includeDependencies": true,
    "adjustDates": true
  }
}
```

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "createdTasks": 5,
    "createdDependencies": 4,
    "tasks": [ ... ],
    "dependencies": [ ... ]
  },
  "message": "Template applied successfully"
}
```

---

### Delete Template

**Endpoint:** `DELETE /templates/:templateId`

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "deletedTemplateId": "tmpl-2"
  }
}
```

---

## Report Endpoints

### Generate Report

**Endpoint:** `POST /projects/:projectId/reports`

**Request Body:**

```json
{
  "type": "progress",
  "format": "pdf",
  "options": {
    "includeSummary": true,
    "includeCharts": true,
    "includeTaskList": true,
    "dateRange": {
      "start": "2024-01-01T00:00:00Z",
      "end": "2024-01-31T00:00:00Z"
    },
    "grouping": {
      "field": "assignee",
      "sort": "progress"
    },
    "filters": {
      "status": ["in_progress", "completed"],
      "priority": ["high", "critical"]
    }
  }
}
```

**Report Types:**
- `progress` - Overall progress report
- `resource` - Resource utilization report
- `milestone` - Milestone status report
- `critical_path` - Critical path analysis
- `burndown` - Burndown chart
- `custom` - Custom report

**Formats:**
- `pdf` - PDF document
- `xlsx` - Excel spreadsheet
- `csv` - CSV data
- `html` - HTML page

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "reportId": "rpt-1",
    "downloadUrl": "/api/v1/gantt/reports/rpt-1/download",
    "expiresAt": "2024-01-10T01:00:00Z",
    "size": 1024000
  }
}
```

---

### Download Report

**Endpoint:** `GET /reports/:reportId/download`

**Response:** `200 OK` (File download)

Returns the generated report file.

---

## Resource Leveling Endpoints

### Analyze Resources

**Endpoint:** `POST /projects/:projectId/level-resources/analyze`

**Request Body:**

```json
{
  "startDate": "2024-01-01T00:00:00Z",
  "endDate": "2024-01-31T00:00:00Z",
  "resources": ["user-1", "user-2", "user-3"]
}
```

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "overallocations": [
      {
        "resourceId": "user-1",
        "resourceName": "John Doe",
        "overallocatedDates": [
          {
            "date": "2024-01-15T00:00:00Z",
            "allocatedHours": 12,
            "capacityHours": 8,
            "tasks": ["task-1", "task-2", "task-3"]
          }
        ]
      }
    ],
    "recommendations": [
      {
        "type": "reschedule",
        "taskId": "task-2",
        "suggestedStartDate": "2024-01-16T00:00:00Z",
        "reason": "Reduce overallocation on 2024-01-15"
      }
    ]
  }
}
```

---

### Apply Resource Leveling

**Endpoint:** `POST /projects/:projectId/level-resources/apply`

**Request Body:**

```json
{
  "strategy": "delay",
  "priority": "high",
  "options": {
    "allowSplitting": false,
    "respectConstraints": true,
    "maxDelayDays": 7
  }
}
```

**Strategies:**
- `delay` - Delay tasks to resolve conflicts
- `split` - Split tasks into segments
- `reassign` - Reassign to different resources
- `extend` - Extend project duration

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "leveledTasks": 5,
    "changes": [
      {
        "taskId": "task-2",
        "oldStartDate": "2024-01-15T00:00:00Z",
        "newStartDate": "2024-01-16T00:00:00Z",
        "reason": "Resource overallocation"
      }
    ],
    "tasks": [ ... ]
  },
  "message": "Resource leveling applied successfully"
}
```

---

## AI Suggestion Endpoints

### Get AI Suggestions

**Endpoint:** `GET /projects/:projectId/ai-suggestions`

**Query Parameters:**
| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `type` | string | - | Filter by type (`schedule`, `resource`, `risk`, `optimization`) |
| `priority` | string | - | Filter by priority (`low`, `medium`, `high`) |

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "suggestions": [
      {
        "id": "sugg-1",
        "type": "schedule",
        "priority": "high",
        "title": "Delay Risk Detected",
        "description": "Task 'API Development' is at risk of delaying dependent tasks",
        "impact": "2-week delay to milestone",
        "effort": "Medium",
        "actions": [
          {
            "type": "reschedule",
            "taskId": "task-5",
            "suggestedDate": "2024-01-20T00:00:00Z"
          },
          {
            "type": "addResource",
            "taskId": "task-5",
            "resourceId": "user-3"
          }
        ],
        "confidence": 0.85,
        "createdAt": "2024-01-05T00:00:00Z"
      }
    ]
  }
}
```

---

### Accept Suggestion

**Endpoint:** `POST /projects/:projectId/ai-suggestions/:suggestionId/accept`

**Request Body:**

```json
{
  "actionIndex": 0,
  "note": "Rescheduling as suggested"
}
```

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "appliedAction": {
      "type": "reschedule",
      "taskId": "task-5",
      "newStartDate": "2024-01-20T00:00:00Z"
    },
    "updatedTasks": [ ... ]
  },
  "message": "Suggestion applied successfully"
}
```

---

### Dismiss Suggestion

**Endpoint:** `POST /projects/:projectId/ai-suggestions/:suggestionId/dismiss`

**Request Body:**

```json
{
  "reason": "Not applicable - timeline is acceptable"
}
```

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "dismissedSuggestionId": "sugg-1"
  }
}
```

---

### Refresh AI Analysis

**Endpoint:** `POST /projects/:projectId/ai-suggestions/refresh`

**Response:** `200 OK`

```json
{
  "success": true,
  "data": {
    "analysisInProgress": true,
    "estimatedTime": 30
  },
  "message": "AI analysis started"
}
```

---

## WebSocket Events

### Connection

**Endpoint:** `WS /api/v1/ws/gantt/:projectId`

**Connection URL:**
```
ws://localhost:8080/api/v1/ws/gantt/proj-123
```

**Query Parameters:**
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `token` | string | Yes | JWT authentication token |

---

### Client → Server Events

#### Subscribe to Project

```json
{
  "type": "subscribe",
  "projectId": "proj-123"
}
```

#### Unsubscribe

```json
{
  "type": "unsubscribe"
}
```

#### Broadcast Update

```json
{
  "type": "update",
  "data": {
    "tasks": [ ... ],
    "dependencies": [ ... ]
  }
}
```

#### Cursor Position (for real-time collaboration)

```json
{
  "type": "cursor",
  "data": {
    "userId": "user-1",
    "userName": "John Doe",
    "taskId": "task-5",
    "position": { "x": 100, "y": 50 }
  }
}
```

---

### Server → Client Events

#### Task Updated

```json
{
  "type": "task_updated",
  "data": {
    "taskId": "task-1",
    "changes": { ... },
    "updatedBy": {
      "userId": "user-2",
      "userName": "Jane Smith"
    },
    "timestamp": "2024-01-05T10:30:00Z"
  }
}
```

#### Task Created

```json
{
  "type": "task_created",
  "data": {
    "task": { ... },
    "createdBy": {
      "userId": "user-2",
      "userName": "Jane Smith"
    },
    "timestamp": "2024-01-05T10:30:00Z"
  }
}
```

#### Task Deleted

```json
{
  "type": "task_deleted",
  "data": {
    "taskId": "task-1",
    "deletedBy": {
      "userId": "user-2",
      "userName": "Jane Smith"
    },
    "timestamp": "2024-01-05T10:30:00Z"
  }
}
```

#### Dependency Created

```json
{
  "type": "dependency_created",
  "data": {
    "dependency": { ... },
    "createdBy": { ... }
  }
}
```

#### User Joined

```json
{
  "type": "user_joined",
  "data": {
    "userId": "user-3",
    "userName": "Bob Wilson",
    "timestamp": "2024-01-05T10:30:00Z"
  }
}
```

#### User Left

```json
{
  "type": "user_left",
  "data": {
    "userId": "user-3",
    "userName": "Bob Wilson",
    "timestamp": "2024-01-05T10:35:00Z"
  }
}
```

#### Cursor Update

```json
{
  "type": "cursor_update",
  "data": {
    "userId": "user-2",
    "userName": "Jane Smith",
    "cursor": {
      "taskId": "task-5",
      "position": { "x": 100, "y": 50 }
    }
  }
}
```

---

## Error Codes

### HTTP Status Codes

| Code | Description |
|------|-------------|
| `200` | Success |
| `201` | Created |
| `204` | No Content |
| `400` | Bad Request |
| `401` | Unauthorized |
| `403` | Forbidden |
| `404` | Not Found |
| `409` | Conflict (e.g., circular dependency) |
| `422` | Validation Error |
| `429` | Rate Limit Exceeded |
| `500` | Internal Server Error |
| `503` | Service Unavailable |

### Application Error Codes

| Error Code | Description |
|------------|-------------|
| `VALIDATION_ERROR` | Request validation failed |
| `AUTHENTICATION_REQUIRED` | No authentication token provided |
| `INVALID_TOKEN` | Authentication token is invalid or expired |
| `INSUFFICIENT_PERMISSIONS` | User lacks required permissions |
| `PROJECT_NOT_FOUND` | Project does not exist |
| `TASK_NOT_FOUND` | Task does not exist |
| `CIRCULAR_DEPENDENCY` | Creating circular dependency |
| `CONSTRAINT_VIOLATION` | Task constraint violated |
| `RESOURCE_OVERALLOCATED` | Resource exceeds capacity |
| `TEMPLATE_NOT_FOUND` | Template does not exist |
| `INVALID_DATE_RANGE` | Start date after end date |
| `DUPLICATE_TASK_NAME` | Task name already exists in project |
| `PARENT_TASK_NOT_FOUND` | Parent task does not exist |
| `RATE_LIMIT_EXCEEDED` | API rate limit exceeded |

---

## Rate Limiting

### Limits

| Tier | Requests per Minute | Requests per Hour |
|------|---------------------|-------------------|
| Free | 60 | 1000 |
| Pro | 300 | 5000 |
| Enterprise | Unlimited | Unlimited |

### Headers

Rate limit information is included in response headers:

```http
X-RateLimit-Limit: 300
X-RateLimit-Remaining: 250
X-RateLimit-Reset: 1704451200
```

### Retry-After

When rate limited:

```http
HTTP/1.1 429 Too Many Requests
Retry-After: 60
X-RateLimit-Remaining: 0
```

**Example Response:**

```json
{
  "success": false,
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "Rate limit exceeded. Please retry after 60 seconds.",
    "retryAfter": 60
  }
}
```

---

## SDK Examples

### JavaScript/TypeScript

```typescript
import axios from 'axios'

class GanttAPI {
  private baseURL = '/api/v1/gantt'
  private token: string

  constructor(token: string) {
    this.token = token
  }

  private get headers() {
    return {
      'Authorization': `Bearer ${this.token}`,
      'Content-Type': 'application/json'
    }
  }

  async getTasks(projectId: string, params?: any) {
    const response = await axios.get(
      `${this.baseURL}/projects/${projectId}/tasks`,
      { headers: this.headers, params }
    )
    return response.data
  }

  async createTask(projectId: string, task: Partial<Task>) {
    const response = await axios.post(
      `${this.baseURL}/projects/${projectId}/tasks`,
      task,
      { headers: this.headers }
    )
    return response.data
  }

  async updateTask(projectId: string, taskId: string, updates: Partial<Task>) {
    const response = await axios.patch(
      `${this.baseURL}/projects/${projectId}/tasks/${taskId}`,
      updates,
      { headers: this.headers }
    )
    return response.data
  }

  async deleteTask(projectId: string, taskId: string) {
    const response = await axios.delete(
      `${this.baseURL}/projects/${projectId}/tasks/${taskId}`,
      { headers: this.headers }
    )
    return response.data
  }
}

// Usage
const api = new GanttAPI('your-jwt-token')
const tasks = await api.getTasks('proj-123')
```

---

### WebSocket Client

```typescript
import { io, Socket } from 'socket.io-client'

class GanttWebSocket {
  private socket: Socket
  private projectId: string
  private token: string

  constructor(projectId: string, token: string) {
    this.projectId = projectId
    this.token = token
    this.socket = io(`/api/v1/ws/gantt/${projectId}`, {
      auth: { token },
      transports: ['websocket']
    })

    this.setupListeners()
  }

  private setupListeners() {
    this.socket.on('connect', () => {
      console.log('Connected to Gantt WebSocket')
      this.socket.emit('subscribe', { projectId: this.projectId })
    })

    this.socket.on('task_updated', (data) => {
      console.log('Task updated:', data)
      // Handle task update
    })

    this.socket.on('task_created', (data) => {
      console.log('Task created:', data)
      // Handle task creation
    })

    this.socket.on('user_joined', (data) => {
      console.log('User joined:', data.userName)
      // Show user joined notification
    })

    this.socket.on('disconnect', () => {
      console.log('Disconnected from Gantt WebSocket')
    })
  }

  broadcastUpdate(data: any) {
    this.socket.emit('update', data)
  }

  sendCursorPosition(position: any) {
    this.socket.emit('cursor', { data: position })
  }

  disconnect() {
    this.socket.disconnect()
  }
}

// Usage
const ws = new GanttWebSocket('proj-123', 'your-jwt-token')
```

---

**Document Version:** 1.0.0
**Last Updated:** 2026-02-19
