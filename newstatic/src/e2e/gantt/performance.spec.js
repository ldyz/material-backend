import { test, expect } from '@playwright/test'

test.describe('Gantt Chart Performance', () => {
  test('should render 100 tasks efficiently', async ({ page }) => {
    // Setup: Load test data
    await page.goto('/gantt?loadTestData=100')
    await page.waitForLoadState('networkidle')

    const startTime = Date.now()

    // Measure initial render
    await page.waitForSelector('[data-testid="task-item"]')
    const renderTime = Date.now() - startTime

    // Verify all tasks rendered
    const taskCount = await page.locator('[data-testid="task-item"]').count()
    expect(taskCount).toBe(100)

    // Render should be fast (< 2 seconds)
    expect(renderTime).toBeLessThan(2000)
  })

  test('should render 500 tasks efficiently', async ({ page }) => {
    await page.goto('/gantt?loadTestData=500')
    await page.waitForLoadState('networkidle')

    const startTime = Date.now()
    await page.waitForSelector('[data-testid="task-item"]')
    const renderTime = Date.now() - startTime

    const taskCount = await page.locator('[data-testid="task-item"]').count()
    expect(taskCount).toBe(500)
    expect(renderTime).toBeLessThan(3000)
  })

  test('should render 1000 tasks efficiently', async ({ page }) => {
    await page.goto('/gantt?loadTestData=1000')
    await page.waitForLoadState('networkidle')

    const startTime = Date.now()
    await page.waitForSelector('[data-testid="task-item"]')
    const renderTime = Date.now() - startTime

    const taskCount = await page.locator('[data-testid="task-item"]').count()
    expect(taskCount).toBe(1000)
    expect(renderTime).toBeLessThan(5000)
  })

  test('should handle smooth scrolling with many tasks', async ({ page }) => {
    await page.goto('/gantt?loadTestData=500')
    await page.waitForLoadState('networkidle')

    // Measure scroll performance
    const startTime = Date.now()

    await page.mouse.wheel(0, 1000)
    await page.waitForTimeout(100)

    const scrollTime = Date.now() - startTime

    // Scroll should be responsive
    expect(scrollTime).toBeLessThan(500)
  })

  test('should drag task quickly with many tasks', async ({ page }) => {
    await page.goto('/gantt?loadTestData=500')
    await page.waitForLoadState('networkidle')

    const task = page.locator('[data-testid="gantt-bar"]').first()

    const startTime = Date.now()

    await task.dragTo(page.locator('[data-testid="timeline-container"]'), {
      targetPosition: { x: 200, y: 0 },
    })

    const dragTime = Date.now() - startTime

    // Drag should be responsive
    expect(dragTime).toBeLessThan(1000)
  })

  test('should calculate dependencies quickly', async ({ page }) => {
    await page.goto('/gantt?loadTestData=200&withDependencies=true')
    await page.waitForLoadState('networkidle')

    const startTime = Date.now()

    // Trigger critical path calculation
    await page.click('[data-testid="calculate-critical-path-button"]')
    await page.waitForSelector('[data-testid="critical-path-indicator"]')

    const calcTime = Date.now() - startTime

    // Calculation should be fast
    expect(calcTime).toBeLessThan(2000)
  })

  test('should filter tasks quickly', async ({ page }) => {
    await page.goto('/gantt?loadTestData=500')
    await page.waitForLoadState('networkidle')

    const startTime = Date.now()

    await page.selectOption('[data-testid="task-status-filter"]', 'in-progress')
    await page.waitForTimeout(100)

    const filterTime = Date.now() - startTime

    // Filter should be fast
    expect(filterTime).toBeLessThan(500)
  })

  test('should search tasks quickly', async ({ page }) => {
    await page.goto('/gantt?loadTestData=500')
    await page.waitForLoadState('networkidle')

    const startTime = Date.now()

    await page.fill('[data-testid="task-search-input"]', 'test')
    await page.waitForTimeout(200)

    const searchTime = Date.now() - startTime

    // Search should be responsive
    expect(searchTime).toBeLessThan(500)
  })

  test('should handle zoom without lag', async ({ page }) => {
    await page.goto('/gantt?loadTestData=200')
    await page.waitForLoadState('networkidle')

    const startTime = Date.now()

    // Zoom in
    await page.click('[data-testid="zoom-in-button"]')
    await page.waitForTimeout(100)

    const zoomTime = Date.now() - startTime

    expect(zoomTime).toBeLessThan(500)
  })

  test('should measure frame rate during interaction', async ({ page }) => {
    await page.goto('/gantt?loadTestData=200')
    await page.waitForLoadState('networkidle')

    // Start FPS monitoring
    const fps = await page.evaluate(() => {
      return new Promise((resolve) => {
        let frames = 0
        let startTime = performance.now()

        function countFrames() {
          frames++
          if (performance.now() - startTime < 1000) {
            requestAnimationFrame(countFrames)
          } else {
            resolve(frames)
          }
        }

        // Start interaction
        const container = document.querySelector('[data-testid="timeline-container"]')
        container.dispatchEvent(new WheelEvent('wheel', { deltaY: 100 }))

        requestAnimationFrame(countFrames)
      })
    })

    // Should maintain decent FPS
    expect(fps).toBeGreaterThan(30)
  })
})

test.describe('Memory Management', () => {
  test('should not leak memory when creating/deleting tasks', async ({ page }) => {
    await page.goto('/gantt')
    await page.waitForLoadState('networkidle')

    // Get initial memory
    const initialMemory = await page.evaluate(() => {
      return (performance as any).memory?.usedJSHeapSize || 0
    })

    // Create and delete tasks multiple times
    for (let i = 0; i < 10; i++) {
      await page.click('[data-testid="add-task-button"]')
      await page.fill('[data-testid="task-name-input"]', `Task ${i}`)
      await page.click('[data-testid="save-task-button"]')
    }

    await page.click('[data-testid="task-item"]')
    await page.click('[data-testid="bulk-delete-button"]')
    await page.click('[data-testid="confirm-bulk-delete-button"]')

    await page.waitForTimeout(1000)

    // Get final memory
    const finalMemory = await page.evaluate(() => {
      return (performance as any).memory?.usedJSHeapSize || 0
    })

    // Memory increase should be reasonable (< 10MB)
    const memoryIncrease = (finalMemory - initialMemory) / 1024 / 1024
    expect(memoryIncrease).toBeLessThan(10)
  })

  test('should cleanup virtual scroll nodes', async ({ page }) => {
    await page.goto('/gantt?loadTestData=1000')
    await page.waitForLoadState('networkidle')

    // Get initial DOM node count
    const initialNodes = await page.evaluate(() => {
      return document.querySelectorAll('[data-testid="task-item"]').length
    })

    // Scroll to bottom
    await page.mouse.wheel(0, 10000)
    await page.waitForTimeout(200)

    // Get node count after scroll
    const scrolledNodes = await page.evaluate(() => {
      return document.querySelectorAll('[data-testid="task-item"]').length
    })

    // Virtual scroll should keep node count constant
    expect(scrolledNodes).toBeCloseTo(initialNodes, 50)
  })
})
