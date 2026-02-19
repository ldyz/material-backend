import { test, expect } from '@playwright/test'

test.describe('Gantt Chart Basic Workflow', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to Gantt chart page
    await page.goto('/gantt')
    await page.waitForLoadState('networkidle')
  })

  test('should create a new task', async ({ page }) => {
    // Click "Add Task" button
    await page.click('[data-testid="add-task-button"]')

    // Fill in task details
    await page.fill('[data-testid="task-name-input"]', 'New Test Task')
    await page.fill('[data-testid="task-start-input"]', '2024-02-01')
    await page.fill('[data-testid="task-end-input"]', '2024-02-05')

    // Save task
    await page.click('[data-testid="save-task-button"]')

    // Verify task was created
    await expect(page.locator('text=New Test Task')).toBeVisible()
    await expect(page.locator('[data-testid="task-list"]')).toContainText('New Test Task')
  })

  test('should edit an existing task', async ({ page }) => {
    // Select first task
    await page.click('[data-testid="task-item"]:first-child')

    // Click edit button
    await page.click('[data-testid="edit-task-button"]')

    // Modify task name
    await page.fill('[data-testid="task-name-input"]', 'Updated Task Name')
    await page.click('[data-testid="save-task-button"]')

    // Verify task was updated
    await expect(page.locator('text=Updated Task Name')).toBeVisible()
  })

  test('should delete a task', async ({ page }) => {
    // Get initial task count
    const initialCount = await page.locator('[data-testid="task-item"]').count()

    // Select first task
    await page.click('[data-testid="task-item"]:first-child')

    // Click delete button
    await page.click('[data-testid="delete-task-button"]')

    // Confirm deletion
    await page.click('[data-testid="confirm-delete-button"]')

    // Verify task was deleted
    const newCount = await page.locator('[data-testid="task-item"]').count()
    expect(newCount).toBe(initialCount - 1)
  })

  test('should create a dependency between tasks', async ({ page }) => {
    // Get two tasks
    const task1 = page.locator('[data-testid="task-item"]').nth(0)
    const task2 = page.locator('[data-testid="task-item"]').nth(1)

    // Drag from task1 to task2 to create dependency
    await task1.dragTo(task2)

    // Verify dependency was created
    await expect(page.locator('[data-testid="dependency-line"]')).toBeVisible()
  })

  test('should undo task creation', async ({ page }) => {
    const initialCount = await page.locator('[data-testid="task-item"]').count()

    // Create a task
    await page.click('[data-testid="add-task-button"]')
    await page.fill('[data-testid="task-name-input"]', 'Task to Undo')
    await page.click('[data-testid="save-task-button"]')

    // Undo
    await page.click('[data-testid="undo-button"]')

    // Verify task was removed
    const newCount = await page.locator('[data-testid="task-item"]').count()
    expect(newCount).toBe(initialCount)
  })

  test('should redo undone action', async ({ page }) => {
    // Create and undo
    await page.click('[data-testid="add-task-button"]')
    await page.fill('[data-testid="task-name-input"]', 'Task to Redo')
    await page.click('[data-testid="save-task-button"]')
    await page.click('[data-testid="undo-button"]')

    // Redo
    await page.click('[data-testid="redo-button"]')

    // Verify task was restored
    await expect(page.locator('text=Task to Redo')).toBeVisible()
  })

  test('should update task progress', async ({ page }) => {
    // Select a task
    await page.click('[data-testid="task-item"]:first-child')

    // Update progress slider
    await page.fill('[data-testid="task-progress-input"]', '50')

    // Save
    await page.click('[data-testid="save-task-button"]')

    // Verify progress was updated
    await expect(page.locator('[data-testid="task-progress-bar"]:first-child')).toHaveAttribute(
      'style',
      /width: 50%/
    )
  })

  test('should change task assignee', async ({ page }) => {
    // Select a task
    await page.click('[data-testid="task-item"]:first-child')

    // Open assignee dropdown
    await page.click('[data-testid="task-assignee-select"]')

    // Select assignee
    await page.click('text=John Doe')

    // Save
    await page.click('[data-testid="save-task-button"]')

    // Verify assignee was changed
    await expect(page.locator('[data-testid="task-item"]:first-child')).toContainText('John Doe')
  })

  test('should filter tasks by status', async ({ page }) => {
    // Apply filter
    await page.selectOption('[data-testid="task-status-filter"]', 'in-progress')

    // Verify only in-progress tasks are shown
    const visibleTasks = await page.locator('[data-testid="task-item"]').allTextContents()
    visibleTasks.forEach(task => {
      expect(task).toMatch(/In Progress/i)
    })
  })

  test('should search for tasks', async ({ page }) => {
    // Enter search term
    await page.fill('[data-testid="task-search-input"]', 'Planning')

    // Verify search results
    await expect(page.locator('[data-testid="task-item"]')).toHaveCount(expect.any(Number))
    const firstTask = await page.locator('[data-testid="task-item"]:first-child').textContent()
    expect(firstTask).toContain('Planning')
  })
})

test.describe('Gantt Chart Keyboard Shortcuts', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/gantt')
    await page.waitForLoadState('networkidle')
  })

  test('should create task with Ctrl+N', async ({ page }) => {
    await page.keyboard.press('Control+N')

    await expect(page.locator('[data-testid="task-name-input"]')).toBeVisible()
  })

  test('should undo with Ctrl+Z', async ({ page }) => {
    const initialCount = await page.locator('[data-testid="task-item"]').count()

    await page.click('[data-testid="add-task-button"]')
    await page.fill('[data-testid="task-name-input"]', 'Keyboard Test')
    await page.click('[data-testid="save-task-button"]')

    await page.keyboard.press('Control+Z')

    const newCount = await page.locator('[data-testid="task-item"]').count()
    expect(newCount).toBe(initialCount)
  })

  test('should redo with Ctrl+Shift+Z', async ({ page }) => {
    await page.click('[data-testid="add-task-button"]')
    await page.fill('[data-testid="task-name-input"]', 'Redo Test')
    await page.click('[data-testid="save-task-button"]')
    await page.keyboard.press('Control+Z')
    await page.keyboard.press('Control+Shift+Z')

    await expect(page.locator('text=Redo Test')).toBeVisible()
  })

  test('should delete task with Delete key', async ({ page }) => {
    await page.click('[data-testid="task-item"]:first-child')
    await page.keyboard.press('Delete')
    await page.click('[data-testid="confirm-delete-button"]')

    await expect(page.locator('[data-testid="task-item"]:first-child')).not.toHaveText(
      await page.locator('[data-testid="task-item"]:first-child').textContent()
    )
  })
})
