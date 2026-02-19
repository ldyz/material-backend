package progress

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConstraintValidation(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	handler := NewProgressHandler(db)

	ctx := context.Background()

	t.Run("valid must-start-on constraint", func(t *testing.T) {
		task := createTestTask(db, "Task 1", time.Now(), time.Now().Add(24*time.Hour))

		constraint := &TaskConstraint{
			TaskID: task.ID,
			Type:   "must-start-on",
			Date:   task.StartTime,
		}

		valid, errs := handler.ValidateConstraint(ctx, constraint, task)
		assert.True(t, valid)
		assert.Empty(t, errs)
	})

	t.Run("invalid constraint type", func(t *testing.T) {
		task := createTestTask(db, "Task 2", time.Now(), time.Now().Add(24*time.Hour))

		constraint := &TaskConstraint{
			TaskID: task.ID,
			Type:   "invalid-type",
			Date:   task.StartTime,
		}

		valid, errs := handler.ValidateConstraint(ctx, constraint, task)
		assert.False(t, valid)
		assert.NotEmpty(t, errs)
	})

	t.Run("constraint conflicts with dependencies", func(t *testing.T) {
		task1 := createTestTask(db, "Task 1", time.Now(), time.Now().Add(24*time.Hour))
		task2 := createTestTask(db, "Task 2", time.Now().Add(48*time.Hour), time.Now().Add(72*time.Hour))

		// Create dependency
		_, err := db.ExecContext(ctx, `
			INSERT INTO task_dependencies (from_task_id, to_task_id, type, lag)
			VALUES (?, ?, 'finish-to-start', 0)
		`, task1.ID, task2.ID)
		require.NoError(t, err)

		// Constraint that violates dependency
		constraint := &TaskConstraint{
			TaskID: task2.ID,
			Type:   "must-start-on",
			Date:   task1.StartTime, // Before task1 finishes
		}

		valid, errs := handler.ValidateConstraint(ctx, constraint, task2)
		assert.False(t, valid)
		assert.NotEmpty(t, errs)
	})
}

func TestConstraintApplication(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	handler := NewProgressHandler(db)
	ctx := context.Background()

	t.Run("apply must-start-on constraint", func(t *testing.T) {
		task := createTestTask(db, "Task 1", time.Now(), time.Now().Add(24*time.Hour))
		newStart := time.Now().Add(48 * time.Hour)

		constraint := &TaskConstraint{
			TaskID: task.ID,
			Type:   "must-start-on",
			Date:   newStart,
		}

		updated, err := handler.ApplyConstraint(ctx, task, constraint)
		require.NoError(t, err)
		assert.Equal(t, newStart.Truncate(time.Second), updated.StartTime.Truncate(time.Second))
		assert.Equal(t, newStart.Add(24*time.Hour).Truncate(time.Second), updated.EndTime.Truncate(time.Second))
	})

	t.Run("apply must-finish-by constraint", func(t *testing.T) {
		task := createTestTask(db, "Task 1", time.Now(), time.Now().Add(24*time.Hour))
		newEnd := time.Now().Add(48 * time.Hour)

		constraint := &TaskConstraint{
			TaskID: task.ID,
			Type:   "must-finish-by",
			Date:   newEnd,
		}

		updated, err := handler.ApplyConstraint(ctx, task, constraint)
		require.NoError(t, err)
		assert.Equal(t, newEnd.Truncate(time.Second), updated.EndTime.Truncate(time.Second))
	})

	t.Run("apply start-no-earlier-than constraint", func(t *testing.T) {
		originalStart := time.Now().Add(12 * time.Hour)
		task := createTestTask(db, "Task 1", originalStart, originalStart.Add(24*time.Hour))
		earliestStart := time.Now().Add(48 * time.Hour)

		constraint := &TaskConstraint{
			TaskID: task.ID,
			Type:   "start-no-earlier-than",
			Date:   earliestStart,
		}

		updated, err := handler.ApplyConstraint(ctx, task, constraint)
		require.NoError(t, err)
		assert.Equal(t, earliestStart.Truncate(time.Second), updated.StartTime.Truncate(time.Second))
	})
}

func TestConstraintConflictDetection(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	handler := NewProgressHandler(db)
	ctx := context.Background()

	t.Run("detect conflicting constraints on same task", func(t *testing.T) {
		task := createTestTask(db, "Task 1", time.Now(), time.Now().Add(24*time.Hour))

		// Create two conflicting must-start-on constraints
		constraint1 := &TaskConstraint{
			TaskID: task.ID,
			Type:   "must-start-on",
			Date:   time.Now(),
		}
		constraint2 := &TaskConstraint{
			TaskID: task.ID,
			Type:   "must-start-on",
			Date:   time.Now().Add(48 * time.Hour),
		}

		_, err := db.ExecContext(ctx, `
			INSERT INTO task_constraints (task_id, type, constraint_date, applied)
			VALUES (?, ?, ?, true)
		`, task.ID, constraint1.Type, constraint1.Date)

		conflicts, err := handler.CheckConstraintConflicts(ctx, task.ID)
		require.NoError(t, err)
		assert.NotEmpty(t, conflicts)
	})

	t.Run("no conflicts with compatible constraints", func(t *testing.T) {
		task := createTestTask(db, "Task 1", time.Now(), time.Now().Add(24*time.Hour))

		constraint := &TaskConstraint{
			TaskID: task.ID,
			Type:   "start-no-earlier-than",
			Date:   time.Now(),
		}

		_, err := db.ExecContext(ctx, `
			INSERT INTO task_constraints (task_id, type, constraint_date, applied)
			VALUES (?, ?, ?, true)
		`, task.ID, constraint.Type, constraint.Date)

		conflicts, err := handler.CheckConstraintConflicts(ctx, task.ID)
		require.NoError(t, err)
		assert.Empty(t, conflicts)
	})
}

func TestConstraintCRUD(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	handler := NewProgressHandler(db)
	ctx := context.Background()

	t.Run("create constraint", func(t *testing.T) {
		task := createTestTask(db, "Task 1", time.Now(), time.Now().Add(24*time.Hour))

		constraint := &TaskConstraint{
			TaskID: task.ID,
			Type:   "must-start-on",
			Date:   time.Now().Add(48 * time.Hour),
		}

		err := handler.CreateConstraint(ctx, constraint)
		require.NoError(t, err)
		assert.NotZero(t, constraint.ID)
	})

	t.Run("get constraint by ID", func(t *testing.T) {
		task := createTestTask(db, "Task 1", time.Now(), time.Now().Add(24*time.Hour))

		constraint := &TaskConstraint{
			TaskID: task.ID,
			Type:   "must-start-on",
			Date:   time.Now().Add(48 * time.Hour),
		}

		err := handler.CreateConstraint(ctx, constraint)
		require.NoError(t, err)

		retrieved, err := handler.GetConstraint(ctx, constraint.ID)
		require.NoError(t, err)
		assert.Equal(t, constraint.Type, retrieved.Type)
		assert.Equal(t, constraint.TaskID, retrieved.TaskID)
	})

	t.Run("update constraint", func(t *testing.T) {
		task := createTestTask(db, "Task 1", time.Now(), time.Now().Add(24*time.Hour))

		constraint := &TaskConstraint{
			TaskID: task.ID,
			Type:   "must-start-on",
			Date:   time.Now().Add(48 * time.Hour),
		}

		err := handler.CreateConstraint(ctx, constraint)
		require.NoError(t, err)

		constraint.Type = "must-finish-by"
		err = handler.UpdateConstraint(ctx, constraint)
		require.NoError(t, err)

		retrieved, err := handler.GetConstraint(ctx, constraint.ID)
		require.NoError(t, err)
		assert.Equal(t, "must-finish-by", retrieved.Type)
	})

	t.Run("delete constraint", func(t *testing.T) {
		task := createTestTask(db, "Task 1", time.Now(), time.Now().Add(24*time.Hour))

		constraint := &TaskConstraint{
			TaskID: task.ID,
			Type:   "must-start-on",
			Date:   time.Now().Add(48 * time.Hour),
		}

		err := handler.CreateConstraint(ctx, constraint)
		require.NoError(t, err)

		err = handler.DeleteConstraint(ctx, constraint.ID)
		require.NoError(t, err)

		_, err = handler.GetConstraint(ctx, constraint.ID)
		assert.Error(t, err)
	})

	t.Run("get constraints by task", func(t *testing.T) {
		task := createTestTask(db, "Task 1", time.Now(), time.Now().Add(24*time.Hour))

		constraint1 := &TaskConstraint{
			TaskID: task.ID,
			Type:   "must-start-on",
			Date:   time.Now(),
		}
		constraint2 := &TaskConstraint{
			TaskID: task.ID,
			Type:   "must-finish-by",
			Date:   time.Now().Add(48 * time.Hour),
		}

		err := handler.CreateConstraint(ctx, constraint1)
		require.NoError(t, err)
		err = handler.CreateConstraint(ctx, constraint2)
		require.NoError(t, err)

		constraints, err := handler.GetConstraintsByTask(ctx, task.ID)
		require.NoError(t, err)
		assert.Len(t, constraints, 2)
	})
}

func TestConstraintImpact(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	handler := NewProgressHandler(db)
	ctx := context.Background()

	t.Run("calculate constraint impact", func(t *testing.T) {
		task1 := createTestTask(db, "Task 1", time.Now(), time.Now().Add(24*time.Hour))
		task2 := createTestTask(db, "Task 2", time.Now().Add(48*time.Hour), time.Now().Add(72*time.Hour))

		// Create dependency
		_, err := db.ExecContext(ctx, `
			INSERT INTO task_dependencies (from_task_id, to_task_id, type, lag)
			VALUES (?, ?, 'finish-to-start', 0)
		`, task1.ID, task2.ID)
		require.NoError(t, err)

		constraint := &TaskConstraint{
			TaskID: task1.ID,
			Type:   "must-start-on",
			Date:   time.Now().Add(96 * time.Hour),
		}

		impact, err := handler.CalculateConstraintImpact(ctx, constraint)
		require.NoError(t, err)

		assert.NotEmpty(t, impact.AffectedTasks)
		assert.Contains(t, impact.AffectedTasks, task2.ID)
		assert.Greater(t, impact.StartDateShift, time.Duration(0))
	})

	t.Run("get constraint suggestions", func(t *testing.T) {
		task := createTestTask(db, "Task 1", time.Now(), time.Now().Add(24*time.Hour))

		suggestions, err := handler.GetConstraintSuggestions(ctx, task)
		require.NoError(t, err)
		assert.NotEmpty(t, suggestions)
	})
}

// Helper functions

func createTestTask(db *DB, name string, start, end time.Time) *Task {
	ctx := context.Background()
	task := &Task{
		ProjectID:  1,
		Name:       name,
		StartTime:  start,
		EndTime:    end,
		Progress:   0,
		Status:     "pending",
		CreatedBy:  1,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err := db.QueryRowContext(ctx, `
		INSERT INTO tasks (project_id, name, start_time, end_time, progress, status, created_by, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id
	`, task.ProjectID, task.Name, task.StartTime, task.EndTime, task.Progress, task.Status,
		task.CreatedBy, task.CreatedAt, task.UpdatedAt).Scan(&task.ID)

	if err != nil {
		panic(err)
	}

	return task
}

func setupTestDB(t *testing.T) *DB {
	// Setup test database connection
	// This would typically use an in-memory SQLite or test PostgreSQL database
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	// Run migrations
	_, err = db.Exec(`
		CREATE TABLE tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			project_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			start_time TIMESTAMP NOT NULL,
			end_time TIMESTAMP NOT NULL,
			progress INTEGER DEFAULT 0,
			status TEXT NOT NULL,
			created_by INTEGER NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);

		CREATE TABLE task_dependencies (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			from_task_id INTEGER NOT NULL,
			to_task_id INTEGER NOT NULL,
			type TEXT NOT NULL,
			lag INTEGER DEFAULT 0
		);

		CREATE TABLE task_constraints (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			task_id INTEGER NOT NULL,
			type TEXT NOT NULL,
			constraint_date TIMESTAMP NOT NULL,
			applied BOOLEAN DEFAULT false
		);
	`)
	require.NoError(t, err)

	return &DB{db}
}

func teardownTestDB(t *testing.T, db *DB) {
	err := db.Close()
	require.NoError(t, err)
}
