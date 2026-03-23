import { test, expect } from '@playwright/test'

test.describe('Gantt Chart Drag and Drop', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/gantt')
    await page.waitForLoadState('networkidle')
  })

  test('should drag task to change start date', async ({ page }) => {
    // Get initial position
    const task = page.locator('[data-testid="gantt-bar"]').first()
    const initialBox = await task.boundingBox()

    // Drag task to the right
    await task.dragTo(page.locator('[data-testid="timeline-container"]'), {
      targetPosition: { x: 100, y: 0 },
    })

    // Verify position changed
    const newBox = await task.boundingBox()
    expect(newBox.x).toBeGreaterThan(initialBox.x)
  })

  test('should resize task from right edge', async ({ page }) => {
    const task = page.locator('[data-testid="gantt-bar"]').first()
    const resizeHandle = task.locator('[data-testid="resize-handle-right"]')

    const initialWidth = (await task.boundingBox()).width

    // Drag resize handle
    await resizeHandle.dragTo(page.locator('[data-testid="timeline-container"]'), {
      targetPosition: { x: 50, y: 0 },
    })

    // Verify width changed
    const newWidth = (await task.boundingBox()).width
    expect(newWidth).toBeGreaterThan(initialWidth)
  })

  test('should resize task from left edge', async ({ page }) => {
    const task = page.locator('[data-testid="gantt-bar"]').first()
    const resizeHandle = task.locator('[data-testid="resize-handle-left"]')

    const initialBox = await task.boundingBox()

    // Drag left resize handle
    await resizeHandle.dragTo(page.locator('[data-testid="timeline-container"]'), {
      targetPosition: { x: -50, y: 0 },
    })

    // Verify position and width changed
    const newBox = await task.boundingBox()
    expect(newBox.x).toBeLessThan(initialBox.x)
  })

  test('should create dependency by dragging between tasks', async ({ page }) => {
    const task1 = page.locator('[data-testid="gantt-bar"]').nth(0)
    const task2 = page.locator('[data-testid="gantt-bar"]').nth(1)

    // Drag from task1 to task2
    await task1.dragTo(task2)

    // Verify dependency line appears
    await expect(page.locator('[data-testid="dependency-line"]')).toBeVisible()
  })

  test('should show dependency preview while dragging', async ({ page }) => {
    const task1 = page.locator('[data-testid="gantt-bar"]').nth(0)
    const task2 = page.locator('[data-testid="gantt-bar"]').nth(1)

    // Start dragging
    await task1.hover()
    await page.mouse.down()

    // Move to task2 (should show preview)
    await task2.hover()

    // Verify preview line is visible
    await expect(page.locator('[data-testid="dependency-preview"]')).toBeVisible()

    // Release to create dependency
    await page.mouse.up()
  })

  test('should prevent drag on locked tasks', async ({ page }) => {
    // Lock a task
    await page.click('[data-testid="task-item"]:first-child')
    await page.click('[data-testid="lock-task-button"]')

    const task = page.locator('[data-testid="gantt-bar"]').first()

    // Try to drag
    const initialBox = await task.boundingBox()
    await task.dragTo(page.locator('[data-testid="timeline-container"]'), {
      targetPosition: { x: 100, y: 0 },
    })

    // Verify position didn't change
    const newBox = await task.boundingBox()
    expect(newBox.x).toBe(initialBox.x)
  })

  test('should move multiple selected tasks', async ({ page }) => {
    // Select multiple tasks
    await page.click('[data-testid="task-item"]').nth(0, { modifiers: ['Control'] })
    await page.click('[data-testid="task-item"]').nth(1, { modifiers: ['Control'] })

    // Drag one task
    const task1 = page.locator('[data-testid="gantt-bar"]').nth(0)
    const initialBox1 = await task1.boundingBox()
    const task2 = page.locator('[data-testid="gantt-bar"]').nth(1)
    const initialBox2 = await task2.boundingBox()

    await task1.dragTo(page.locator('[data-testid="timeline-container"]'), {
      targetPosition: { x: 50, y: 0 },
    })

    // Verify both tasks moved
    const newBox1 = await task1.boundingBox()
    const newBox2 = await task2.boundingBox()
    expect(newBox1.x).not.toBe(initialBox1.x)
    expect(newBox2.x).not.toBe(initialBox2.x)
  })

  test('should snap to grid when dragging', async ({ page }) => {
    // Enable snap to grid
    await page.click('[data-testid="snap-to-grid-toggle"]')

    const task = page.locator('[data-testid="gantt-bar"]').first()
    await task.dragTo(page.locator('[data-testid="timeline-container"]'), {
      targetPosition: { x: 23, y: 0 }, // Non-grid position
    })

    // Verify position snapped to grid
    const box = await task.boundingBox()
    expect(box.x % 50).toBe(0) // Assuming 50px grid
  })

  test('should update dependencies when task moves', async ({ page }) => {
    // Create dependency
    const task1 = page.locator('[data-testid="gantt-bar"]').nth(0)
    const task2 = page.locator('[data-testid="gantt-bar"]').nth(1)
    await task1.dragTo(task2)

    // Move task1
    const initialBox = await task1.boundingBox()
    await task1.dragTo(page.locator('[data-testid="timeline-container"]'), {
      targetPosition: { x: 100, y: 0 },
    })

    // Verify dependency line updated
    await expect(page.locator('[data-testid="dependency-line"]')).toBeVisible()
  })

  test('should show tooltip with date while dragging', async ({ page }) => {
    const task = page.locator('[data-testid="gantt-bar"]').first()

    await task.hover()
    await page.mouse.down()
    await page.mouse.move(100, 0)

    // Verify tooltip shows new date
    await expect(page.locator('[data-testid="drag-tooltip"]')).toBeVisible()
    await expect(page.locator('[data-testid="drag-tooltip"]')).toContainText(/\d{4}-\d{2}-\d{2}/)
  })
})

test.describe('Gantt Chart Bulk Operations', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/gantt')
    await page.waitForLoadState('networkidle')
  })

  test('should select multiple tasks with Shift+Click', async ({ page }) => {
    await page.click('[data-testid="task-item"]').nth(0)
    await page.click('[data-testid="task-item"]').nth(2, { modifiers: ['Shift'] })

    // Verify 3 tasks selected (0, 1, 2)
    const selectedCount = await page.locator('[data-testid="task-item"].selected').count()
    expect(selectedCount).toBe(3)
  })

  test('should select all tasks with Ctrl+A', async ({ page }) => {
    await page.keyboard.press('Control+A')

    const allTasks = await page.locator('[data-testid="task-item"]').count()
    const selectedTasks = await page.locator('[data-testid="task-item"].selected').count()

    expect(selectedTasks).toBe(allTasks)
  })

  test('should bulk edit selected tasks', async ({ page }) => {
    // Select multiple tasks
    await page.click('[data-testid="task-item"]').nth(0)
    await page.click('[data-testid="task-item"]').nth(1, { modifiers: ['Control'] })

    // Open bulk edit
    await page.click('[data-testid="bulk-edit-button"]')

    // Change priority
    await page.selectOption('[data-testid="bulk-priority-select"]', 'high')

    // Apply
    await page.click('[data-testid="apply-bulk-edit-button"]')

    // Verify both tasks updated
    await expect(page.locator('[data-testid="task-item"]').nth(0)).toContainText('High')
    await expect(page.locator('[data-testid="task-item"]').nth(1)).toContainText('High')
  })

  test('should bulk delete selected tasks', async ({ page }) => {
    // Select multiple tasks
    await page.click('[data-testid="task-item"]').nth(0)
    await page.click('[data-testid="task-item"]').nth(1, { modifiers: ['Control'] })

    // Delete
    await page.click('[data-testid="bulk-delete-button"]')
    await page.click('[data-testid="confirm-bulk-delete-button"]')

    // Verify tasks deleted
    const remaining = await page.locator('[data-testid="task-item"]').count()
    expect(remaining).toBeLessThan(2)
  })
})
