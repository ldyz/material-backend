# Gantt Chart Testing Suite Documentation

## Overview

This document describes the comprehensive testing suite for the Gantt chart system, including unit tests, component tests, E2E tests, and backend tests.

## Test Structure

```
newstatic/
├── src/
│   ├── __tests__/
│   │   ├── setup.ts                 # Test configuration and mocks
│   │   ├── utils/
│   │   │   └── testHelpers.ts       # Helper functions for tests
│   │   └── mocks/
│   │       └── ganttMockData.js     # Mock data for tests
│   ├── utils/__tests__/             # Unit tests for utilities
│   ├── stores/__tests__/            # Unit tests for stores
│   └── e2e/gantt/                   # E2E tests
└── vitest.config.js                 # Vitest configuration
```

## Running Tests

### Unit Tests (Vitest)

```bash
# Install dependencies
npm install --save-dev vitest @vitest/ui jsdom @vue/test-utils

# Run all unit tests
npm run test

# Run tests in watch mode
npm run test:watch

# Run tests with UI
npm run test:ui

# Run tests with coverage
npm run test:coverage

# Run specific test file
npx vitest src/utils/__tests__/ganttConstraints.test.js
```

### E2E Tests (Playwright)

```bash
# Install Playwright
npm install --save-dev @playwright/test
npx playwright install

# Run all E2E tests
npm run test:e2e

# Run E2E tests in headed mode
npm run test:e2e:headed

# Run E2E tests for specific file
npx playwright test src/e2e/gantt/basic-workflow.spec.js

# Debug E2E tests
npx playwright test --debug
```

### Backend Tests (Go)

```bash
# Run all tests
go test ./internal/api/progress/...

# Run with coverage
go test -cover ./internal/api/progress/...

# Run specific test
go test -v ./internal/api/progress/ -run TestConstraintValidation

# Run benchmarks
go test -bench=. ./internal/api/progress/...
```

## Unit Tests

### 1. Gantt Constraints Tests (`ganttConstraints.test.js`)

**Purpose**: Test constraint validation, application, and conflict detection.

**Test Cases**:
- Validate all constraint types (must-start-on, must-finish-by, etc.)
- Apply constraints to tasks
- Calculate constraint impact
- Detect constraint conflicts
- Generate constraint suggestions
- Handle edge cases (holidays, invalid dates, etc.)

**Run**:
```bash
npx vitest src/utils/__tests__/ganttConstraints.test.js
```

### 2. Resource Leveling Tests (`resourceLeveling.test.js`)

**Purpose**: Test resource conflict detection and leveling algorithms.

**Test Cases**:
- Detect resource overallocation
- Calculate resource allocation
- Level resources with different strategies
- Find overallocated resources
- Calculate resource utilization
- Suggest leveling actions
- Handle material resources

**Run**:
```bash
npx vitest src/utils/__tests__/resourceLeveling.test.js
```

### 3. Dependency Validator Tests (`dependencyValidator.test.js`)

**Purpose**: Test dependency validation and critical path calculation.

**Test Cases**:
- Validate all dependency types (FS, SS, FF, SF)
- Detect circular dependencies
- Validate lag/lead times
- Calculate critical path
- Analyze dependency paths
- Handle complex dependency networks

**Run**:
```bash
npx vitest src/utils/__tests__/dependencyValidator.test.js
```

### 4. Critical Path Tests (`criticalPath.test.js`)

**Purpose**: Test critical path method calculations.

**Test Cases**:
- Calculate critical path for various network structures
- Calculate slack (float) for tasks
- Handle all 4 dependency types
- Calculate early/late start and finish dates
- Identify multiple critical paths
- Performance testing with large datasets

**Run**:
```bash
npx vitest src/utils/__tests__/criticalPath.test.js
```

### 5. Change Tracker Tests (`changeTracker.test.js`)

**Purpose**: Test change tracking and diff generation.

**Test Cases**:
- Track changes between task states
- Generate diffs (text, HTML, markdown, JSON)
- Export/import changes
- Compress changes
- Calculate change impact
- Handle change history

**Run**:
```bash
npx vitest src/utils/__tests__/changeTracker.test.js
```

### 6. Store Tests

#### Undo/Redo Store (`undoRedoStore.test.js`)

**Purpose**: Test undo/redo functionality.

**Test Cases**:
- Execute commands
- Undo/redo operations
- Macro commands
- Batch operations
- Command stack limits
- Clear operations

**Run**:
```bash
npx vitest src/stores/__tests__/undoRedoStore.test.js
```

#### Calendar Store (`calendarStore.test.js`)

**Purpose**: Test calendar and working time calculations.

**Test Cases**:
- Working days/hours calculations
- Holiday management
- Multiple calendars
- Calendar validation
- Date utilities

**Run**:
```bash
npx vitest src/stores/__tests__/calendarStore.test.js
```

## E2E Tests

### 1. Basic Workflow Tests (`basic-workflow.spec.js`)

**Purpose**: Test core Gantt chart functionality.

**Test Cases**:
- Create, edit, delete tasks
- Create dependencies
- Undo/redo operations
- Update progress
- Change assignees
- Filter and search tasks
- Keyboard shortcuts

**Run**:
```bash
npx playwright test src/e2e/gantt/basic-workflow.spec.js
```

### 2. Drag and Drop Tests (`drag-drop.spec.js`)

**Purpose**: Test interactive drag-and-drop features.

**Test Cases**:
- Drag tasks to change dates
- Resize tasks
- Create dependencies by dragging
- Bulk operations
- Snap to grid
- Locked tasks

**Run**:
```bash
npx playwright test src/e2e/gantt/drag-drop.spec.js
```

### 3. Performance Tests (`performance.spec.js`)

**Purpose**: Test performance with large datasets.

**Test Cases**:
- Render 100, 500, 1000 tasks
- Scroll performance
- Drag response time
- Filter/search speed
- Memory management
- Frame rate during interactions

**Run**:
```bash
npx playwright test src/e2e/gantt/performance.spec.js
```

## Backend Tests (Go)

### 1. Constraint Tests (`constraint_test.go`)

**Purpose**: Test constraint CRUD and validation.

**Test Cases**:
- Validate constraint types
- Apply constraints
- Detect conflicts
- Calculate impact
- CRUD operations

**Run**:
```bash
go test -v ./internal/api/progress/ -run TestConstraint
```

### 2. Resource Leveling Tests (`resource_leveling_test.go`)

**Purpose**: Test resource management.

**Test Cases**:
- Detect conflicts
- Calculate allocation
- Level resources
- Calculate utilization
- Handle dependencies and constraints

**Run**:
```bash
go test -v ./internal/api/progress/ -run TestResource
```

### 3. Comment Tests (`comment_test.go`)

**Purpose**: Test comment system.

**Test Cases**:
- CRUD operations
- Threading
- Mentions
- Attachments
- Pagination
- Search
- Permissions

**Run**:
```bash
go test -v ./internal/api/progress/ -run TestComment
```

## Test Coverage

### Coverage Goals

- **Statements**: 70%
- **Branches**: 65%
- **Functions**: 70%
- **Lines**: 70%

### Generate Coverage Report

```bash
# Vitest
npx vitest --coverage

# View HTML coverage report
open coverage/index.html
```

### Coverage by Module

- `ganttConstraints.js`: 85%+
- `resourceLeveling.js`: 80%+
- `dependencyValidator.js`: 85%+
- `criticalPath.js`: 80%+
- `changeTracker.js`: 75%+
- `calendarStore.js`: 80%+
- `undoRedoStore.js`: 85%+

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Tests

on: [push, pull_request]

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '18'
      - run: npm ci
      - run: npm run test:coverage
      - uses: codecov/codecov-action@v3

  e2e-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
      - run: npm ci
      - run: npx playwright install --with-deps
      - run: npm run test:e2e

  backend-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - run: go test -cover ./internal/api/progress/...
```

## Best Practices

### Unit Tests

1. **Keep tests isolated**: Each test should be independent
2. **Use descriptive names**: Test names should describe what they test
3. **Arrange-Act-Assert**: Structure tests clearly
4. **Mock external dependencies**: Use test helpers and mocks
5. **Test edge cases**: Don't just test happy paths

### E2E Tests

1. **Use data-testid**: Select elements by test IDs, not CSS selectors
2. **Wait for elements**: Use explicit waits, not `sleep()`
3. **Clean up**: Reset state between tests
4. **Use page objects**: Organize test code
5. **Avoid flakiness**: Make tests reliable and deterministic

### Backend Tests

1. **Use test databases**: Don't use production data
2. **Clean up resources**: Use defer for cleanup
3. **Test concurrency**: Test race conditions
4. **Use table-driven tests**: For multiple test cases
5. **Mock external services**: Don't make real network calls

## Test Data

### Mock Data Location

- Frontend: `/home/julei/backend/newstatic/src/__tests__/mocks/ganttMockData.js`
- Backend: Created programmatically in test setup

### Sample Test Data

- 6 sample tasks with dependencies
- 4 sample resources
- 2 sample calendars
- 3 sample constraints
- 4 sample users
- Large datasets for performance tests (100-1000 tasks)

## Debugging Tests

### Vitest Debugging

```bash
# Run with debug flag
npx vitest --inspect-brk --no-coverage

# Or use VS Code debugger
# Add launch configuration:
{
  "type": "node",
  "request": "launch",
  "name": "Debug Vitest",
  "program": "${workspaceFolder}/node_modules/.bin/vitest",
  "args": ["run", "--reporter=verbose"]
}
```

### Playwright Debugging

```bash
# Run in debug mode
npx playwright test --debug

# Run in headed mode
npx playwright test --headed

# Run with trace
npx playwright test --trace on

# View trace
npx playwright show-trace trace.zip
```

### Go Test Debugging

```bash
# Run with race detection
go test -race ./internal/api/progress/...

# Run with verbose output
go test -v ./internal/api/progress/

# Run specific test with debug
go test -v ./internal/api/progress/ -run TestConstraintValidation -test.v
```

## Performance Benchmarks

### Backend Benchmarks

```bash
# Run all benchmarks
go test -bench=. -benchmem ./internal/api/progress/...

# Run specific benchmark
go test -bench=BenchmarkResourceConflictDetection -benchmem ./internal/api/progress/
```

### Frontend Benchmarks

- Initial render: < 2s for 500 tasks
- Scroll performance: > 30 FPS
- Drag response: < 100ms
- Filter/search: < 500ms

## Troubleshooting

### Common Issues

1. **Test timeout**: Increase timeout in vitest.config.js
2. **Flaky E2E tests**: Add explicit waits, use data-testid
3. **Database locked**: Ensure proper cleanup in backend tests
4. **Memory leaks**: Check for event listener cleanup

### Test Isolation

- Each test should clean up after itself
- Use `beforeEach` and `afterEach` hooks
- Reset database state between tests
- Clear localStorage/sessionStorage

## Reporting

### HTML Reports

- Vitest: `coverage/index.html`
- Playwright: `playwright-report/index.html`

### JUnit Reports

```bash
# Vitest
npx vitest --reporter=junit --outputFile=test-results/junit.xml

# Playwright (configured in playwright.config.js)
# Output: test-results/junit.xml
```

## Maintenance

### Regular Tasks

1. Update test data when schema changes
2. Review and update flaky tests
3. Add tests for new features
4. Monitor coverage metrics
5. Update documentation

### Test Reviews

- Review failed tests before merging
- Ensure new code has tests
- Update tests when refactoring
- Keep test data realistic

## Resources

- [Vitest Documentation](https://vitest.dev/)
- [Playwright Documentation](https://playwright.dev/)
- [Vue Test Utils](https://test-utils.vuejs.org/)
- [Go Testing Guide](https://golang.org/doc/tutorial/add-a-test)

## Support

For questions or issues with tests:
1. Check this documentation
2. Review test code for examples
3. Check framework documentation
4. Ask team members for help
