package progress

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResourceLeveling(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	handler := NewProgressHandler(db)
	ctx := context.Background()

	t.Run("detect resource conflicts", func(t *testing.T) {
		// Create two overlapping tasks with same resource
		now := time.Now()
		task1 := createTestTaskWithResource(db, "Task 1", now, now.Add(24*time.Hour), "res-1")
		task2 := createTestTaskWithResource(db, "Task 2", now, now.Add(24*time.Hour), "res-1")

		conflicts, err := handler.DetectResourceConflicts(ctx, []int64{task1.ID, task2.ID})
		require.NoError(t, err)
		assert.NotEmpty(t, conflicts)
		assert.Contains(t, conflicts, "res-1")
	})

	t.Run("calculate resource allocation", func(t *testing.T) {
		now := time.Now()
		task := createTestTaskWithResource(db, "Task 1", now, now.Add(8*time.Hour), "res-1")

		allocation, err := handler.CalculateResourceAllocation(ctx, "res-1", now, now.Add(24*time.Hour))
		require.NoError(t, err)
		assert.Greater(t, allocation.Hours, float64(0))
		assert.LessOrEqual(t, allocation.Percentage, 100.0)
	})

	t.Run("level overallocated resources", func(t *testing.T) {
		// Create three overlapping tasks with same resource
		now := time.Now()
		task1 := createTestTaskWithResource(db, "Task 1", now, now.Add(24*time.Hour), "res-1")
		task2 := createTestTaskWithResource(db, "Task 2", now, now.Add(24*time.Hour), "res-1")
		task3 := createTestTaskWithResource(db, "Task 3", now, now.Add(24*time.Hour), "res-1")

		// Set priorities
		db.ExecContext(ctx, "UPDATE tasks SET priority = ? WHERE id = ?", "high", task1.ID)
		db.ExecContext(ctx, "UPDATE tasks SET priority = ? WHERE id = ?", "medium", task2.ID)
		db.ExecContext(ctx, "UPDATE tasks SET priority = ? WHERE id = ?", "low", task3.ID)

		result, err := handler.LevelResources(ctx, []int64{task1.ID, task2.ID, task3.ID})
		require.NoError(t, err)
		assert.NotEmpty(t, result.AdjustedTasks)
		assert.Equal(t, task1.StartTime, result.AdjustedTasks[0].StartTime) // High priority task not moved
	})

	t.Run("find overallocated resources", func(t *testing.T) {
		now := time.Now()
		createTestTaskWithResource(db, "Task 1", now, now.Add(24*time.Hour), "res-1")
		createTestTaskWithResource(db, "Task 2", now, now.Add(24*time.Hour), "res-1")

		overallocated, err := handler.FindOverallocatedResources(ctx, now, now.Add(24*time.Hour))
		require.NoError(t, err)
		assert.NotEmpty(t, overallocated)
	})

	t.Run("calculate resource utilization", func(t *testing.T) {
		now := time.Now()
		createTestTaskWithResource(db, "Task 1", now, now.Add(8*time.Hour), "res-1")

		utilization, err := handler.CalculateResourceUtilization(ctx, now, now.Add(24*time.Hour))
		require.NoError(t, err)
		assert.GreaterOrEqual(t, utilization.Overall, float64(0))
		assert.LessOrEqual(t, utilization.Overall, float64(100))
	})

	t.Run("suggest leveling actions", func(t *testing.T) {
		now := time.Now()
		task1 := createTestTaskWithResource(db, "Task 1", now, now.Add(24*time.Hour), "res-1")
		task2 := createTestTaskWithResource(db, "Task 2", now, now.Add(24*time.Hour), "res-1")

		actions, err := handler.SuggestLevelingActions(ctx, []int64{task1.ID, task2.ID})
		require.NoError(t, err)
		assert.NotEmpty(t, actions)
	})
}

func TestResourceAllocationCalculation(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	handler := NewProgressHandler(db)
	ctx := context.Background()

	t.Run("calculate daily allocation", func(t *testing.T) {
		now := time.Now()
		task := createTestTaskWithResource(db, "Task 1", now, now.Add(24*time.Hour), "res-1")

		daily, err := handler.CalculateDailyAllocation(ctx, "res-1", now, now.Add(7*24*time.Hour))
		require.NoError(t, err)
		assert.Len(t, daily, 7)
	})

	t.Run("calculate peak usage", func(t *testing.T) {
		now := time.Now()
		createTestTaskWithResource(db, "Task 1", now, now.Add(24*time.Hour), "res-1")
		createTestTaskWithResource(db, "Task 2", now, now.Add(24*time.Hour), "res-1")

		peak, err := handler.CalculatePeakUsage(ctx, "res-1", now, now.Add(24*time.Hour))
		require.NoError(t, err)
		assert.Greater(t, peak.Hours, 8.0) // Over capacity
	})

	t.Run("identify underutilized resources", func(t *testing.T) {
		now := time.Now()
		createTestTaskWithResource(db, "Task 1", now, now.Add(4*time.Hour), "res-1") // Only 4 hours

		underutilized, err := handler.IdentifyUnderutilizedResources(ctx, now, now.Add(24*time.Hour), 50.0)
		require.NoError(t, err)
		assert.Contains(t, underutilized, "res-1")
	})
}

func TestResourceLevelingWithDependencies(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	handler := NewProgressHandler(db)
	ctx := context.Background()

	t.Run("respect dependencies during leveling", func(t *testing.T) {
		now := time.Now()
		task1 := createTestTaskWithResource(db, "Task 1", now, now.Add(24*time.Hour), "res-1")
		task2 := createTestTaskWithResource(db, "Task 2", now.Add(48*time.Hour), now.Add(72*time.Hour), "res-1")

		// Create dependency
		_, err := db.ExecContext(ctx, `
			INSERT INTO task_dependencies (from_task_id, to_task_id, type, lag)
			VALUES (?, ?, 'finish-to-start', 0)
		`, task1.ID, task2.ID)
		require.NoError(t, err)

		// Level resources
		result, err := handler.LevelResources(ctx, []int64{task1.ID, task2.ID})
		require.NoError(t, err)

		// Verify dependency is preserved
		task1End := result.GetAdjustedTask(task1.ID).EndTime
		task2Start := result.GetAdjustedTask(task2.ID).StartTime
		assert.True(t, task2Start.After(task1End) || task2Start.Equal(task1End))
	})
}

func TestResourceLevelingWithConstraints(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	handler := NewProgressHandler(db)
	ctx := context.Background()

	t.Run("respect must-start-on constraint", func(t *testing.T) {
		now := time.Now()
		fixedStart := now.Add(48 * time.Hour)
		task := createTestTaskWithResource(db, "Task 1", now, now.Add(24*time.Hour), "res-1")

		// Add constraint
		_, err := db.ExecContext(ctx, `
			INSERT INTO task_constraints (task_id, type, constraint_date, applied)
			VALUES (?, 'must-start-on', ?, true)
		`, task.ID, fixedStart)
		require.NoError(t, err)

		// Try to level (should not move constrained task)
		result, err := handler.LevelResources(ctx, []int64{task.ID})
		require.NoError(t, err)

		adjustedTask := result.GetAdjustedTask(task.ID)
		assert.Equal(t, fixedStart.Truncate(time.Second), adjustedTask.StartTime.Truncate(time.Second))
	})
}

func TestMaterialResources(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	handler := NewProgressHandler(db)
	ctx := context.Background()

	t.Run("handle material resource conflicts", func(t *testing.T) {
		now := time.Now()
		// Create two tasks using same meeting room
		task1 := createTestTaskWithResource(db, "Meeting 1", now, now.Add(2*time.Hour), "room-1")
		task2 := createTestTaskWithResource(db, "Meeting 2", now, now.Add(2*time.Hour), "room-1")

		conflicts, err := handler.DetectResourceConflicts(ctx, []int64{task1.ID, task2.ID})
		require.NoError(t, err)
		assert.Contains(t, conflicts, "room-1")
	})

	t.Run("handle resource with capacity > 1", func(t *testing.T) {
		now := time.Now()
		// Resource with capacity of 5
		_, err := db.ExecContext(ctx, `
			INSERT INTO resources (id, name, type, capacity, unit)
			VALUES ('team-1', 'Team A', 'team', 5, 'people')
		`)
		require.NoError(t, err)

		// Create 3 tasks using the team
		for i := 1; i <= 3; i++ {
			createTestTaskWithResource(db, fmt.Sprintf("Task %d", i), now, now.Add(24*time.Hour), "team-1")
		}

		conflicts, err := handler.DetectResourceConflicts(ctx, []int64{1, 2, 3})
		require.NoError(t, err)
		assert.Empty(t, conflicts) // No conflict, within capacity
	})
}

func BenchmarkResourceConflictDetection(b *testing.B) {
	db := setupTestDB(b)
	defer teardownTestDB(b, db)

	handler := NewProgressHandler(db)
	ctx := context.Background()

	// Create 1000 tasks
	now := time.Now()
	var taskIDs []int64
	for i := 0; i < 1000; i++ {
		task := createTestTaskWithResource(db, fmt.Sprintf("Task %d", i), now, now.Add(24*time.Hour), "res-1")
		taskIDs = append(taskIDs, task.ID)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = handler.DetectResourceConflicts(ctx, taskIDs)
	}
}

// Helper functions

func createTestTaskWithResource(db *DB, name string, start, end time.Time, resourceID string) *Task {
	ctx := context.Background()
	task := createTestTask(db, name, start, end)

	// Assign resource
	_, err := db.ExecContext(ctx, `
		INSERT INTO task_resources (task_id, resource_id, allocation)
		VALUES (?, ?, 100)
	`, task.ID, resourceID)
	if err != nil {
		panic(err)
	}

	return task
}
