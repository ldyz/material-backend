package progress

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommentCRUD(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	handler := NewCommentHandler(db)
	ctx := context.Background()

	t.Run("create comment", func(t *testing.T) {
		task := createTestTask(db, "Task 1", time.Now(), time.Now().Add(24*time.Hour))

		comment := &Comment{
			TaskID:    task.ID,
			UserID:    1,
			Content:   "Test comment",
			CreatedAt: time.Now(),
		}

		err := handler.CreateComment(ctx, comment)
		require.NoError(t, err)
		assert.NotZero(t, comment.ID)
	})

	t.Run("get comment by ID", func(t *testing.T) {
		task := createTestTask(db, "Task 1", time.Now(), time.Now().Add(24*time.Hour))

		comment := &Comment{
			TaskID:    task.ID,
			UserID:    1,
			UserName:  "John Doe",
			Content:   "Test comment",
			CreatedAt: time.Now(),
		}

		err := handler.CreateComment(ctx, comment)
		require.NoError(t, err)

		retrieved, err := handler.GetComment(ctx, comment.ID)
		require.NoError(t, err)
		assert.Equal(t, comment.Content, retrieved.Content)
		assert.Equal(t, comment.UserName, retrieved.UserName)
	})

	t.Run("update comment", func(t *testing.T) {
		task := createTestTask(db, "Task 1", time.Now(), time.Now().Add(24*time.Hour))

		comment := &Comment{
			TaskID:    task.ID,
			UserID:    1,
			Content:   "Original comment",
			CreatedAt: time.Now(),
		}

		err := handler.CreateComment(ctx, comment)
		require.NoError(t, err)

		comment.Content = "Updated comment"
		err = handler.UpdateComment(ctx, comment)
		require.NoError(t, err)

		retrieved, err := handler.GetComment(ctx, comment.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated comment", retrieved.Content)
	})

	t.Run("delete comment", func(t *testing.T) {
		task := createTestTask(db, "Task 1", time.Now(), time.Now().Add(24*time.Hour))

		comment := &Comment{
			TaskID:    task.ID,
			UserID:    1,
			Content:   "Test comment",
			CreatedAt: time.Now(),
		}

		err := handler.CreateComment(ctx, comment)
		require.NoError(t, err)

		err = handler.DeleteComment(ctx, comment.ID)
		require.NoError(t, err)

		_, err = handler.GetComment(ctx, comment.ID)
		assert.Error(t, err)
	})

	t.Run("get comments by task", func(t *testing.T) {
		task := createTestTask(db, "Task 1", time.Now(), time.Now().Add(24*time.Hour))

		// Create multiple comments
		for i := 1; i <= 3; i++ {
			comment := &Comment{
				TaskID:    task.ID,
				UserID:    1,
				Content:   fmt.Sprintf("Comment %d", i),
				CreatedAt: time.Now(),
			}
			err := handler.CreateComment(ctx, comment)
			require.NoError(t, err)
		}

		comments, err := handler.GetCommentsByTask(ctx, task.ID)
		require.NoError(t, err)
		assert.Len(t, comments, 3)
	})
}

func TestCommentThreading(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	handler := NewCommentHandler(db)
	ctx := context.Background()

	t.Run("create threaded comments", func(t *testing.T) {
		task := createTestTask(db, "Task 1", time.Now(), time.Now().Add(24*time.Hour))

		// Create parent comment
		parent := &Comment{
			TaskID:    task.ID,
			UserID:    1,
			Content:   "Parent comment",
			CreatedAt: time.Now(),
		}

		err := handler.CreateComment(ctx, parent)
		require.NoError(t, err)

		// Create reply
		reply := &Comment{
			TaskID:    task.ID,
			UserID:    2,
			Content:   "Reply comment",
			ParentID:  parent.ID,
			CreatedAt: time.Now(),
		}

		err = handler.CreateComment(ctx, reply)
		require.NoError(t, err)

		// Verify threading
		comments, err := handler.GetCommentsByTask(ctx, task.ID)
		require.NoError(t, err)
		assert.Len(t, comments, 2)

		// Parent should have replies
		assert.Len(t, comments[0].Replies, 1)
	})

	t.Run("get comment thread", func(t *testing.T) {
		task := createTestTask(db, "Task 1", time.Now(), time.Now().Add(24*time.Hour))

		// Create thread
		parent := &Comment{
			TaskID:    task.ID,
			UserID:    1,
			Content:   "Parent comment",
			CreatedAt: time.Now(),
		}

		err := handler.CreateComment(ctx, parent)
		require.NoError(t, err)

		for i := 1; i <= 3; i++ {
			reply := &Comment{
				TaskID:    task.ID,
				UserID:    2,
				Content:   fmt.Sprintf("Reply %d", i),
				ParentID:  parent.ID,
				CreatedAt: time.Now(),
			}
			err = handler.CreateComment(ctx, reply)
			require.NoError(t, err)
		}

		// Get thread
		thread, err := handler.GetCommentThread(ctx, parent.ID)
		require.NoError(t, err)
		assert.Len(t, thread.Replies, 3)
	})
}

func TestCommentMentions(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	handler := NewCommentHandler(db)
	ctx := context.Background()

	t.Run("create comment with mentions", func(t *testing.T) {
		task := createTestTask(db, "Task 1", time.Now(), time.Now().Add(24*time.Hour))

		comment := &Comment{
			TaskID:    task.ID,
			UserID:    1,
			Content:   "@john please review this",
			CreatedAt: time.Now(),
			Mentions:  []string{"john"},
		}

		err := handler.CreateComment(ctx, comment)
		require.NoError(t, err)

		// Verify mentions stored
		retrieved, err := handler.GetComment(ctx, comment.ID)
		require.NoError(t, err)
		assert.Contains(t, retrieved.Mentions, "john")
	})

	t.Run("get mentions for user", func(t *testing.T) {
		task := createTestTask(db, "Task 1", time.Now(), time.Now().Add(24*time.Hour))

		// Create comments mentioning user
		for i := 1; i <= 3; i++ {
			comment := &Comment{
				TaskID:    task.ID,
				UserID:    1,
				Content:   fmt.Sprintf("Mentioning @jane in comment %d", i),
				CreatedAt: time.Now(),
				Mentions:  []string{"jane"},
			}
			err := handler.CreateComment(ctx, comment)
			require.NoError(t, err)
		}

		// Get mentions
		mentions, err := handler.GetMentionsForUser(ctx, "jane")
		require.NoError(t, err)
		assert.Len(t, mentions, 3)
	})
}

func TestCommentAttachments(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	handler := NewCommentHandler(db)
	ctx := context.Background()

	t.Run("create comment with attachments", func(t *testing.T) {
		task := createTestTask(db, "Task 1", time.Now(), time.Now().Add(24*time.Hour))

		comment := &Comment{
			TaskID:    task.ID,
			UserID:    1,
			Content:   "Please review attached files",
			CreatedAt: time.Now(),
			Attachments: []Attachment{
				{FileName: "document.pdf", FileSize: 1024, FileURL: "/files/document.pdf"},
				{FileName: "image.png", FileSize: 2048, FileURL: "/files/image.png"},
			},
		}

		err := handler.CreateComment(ctx, comment)
		require.NoError(t, err)

		// Verify attachments
		retrieved, err := handler.GetComment(ctx, comment.ID)
		require.NoError(t, err)
		assert.Len(t, retrieved.Attachments, 2)
	})

	t.Run("delete attachment", func(t *testing.T) {
		task := createTestTask(db, "Task 1", time.Now(), time.Now().Add(24*time.Hour))

		comment := &Comment{
			TaskID:    task.ID,
			UserID:    1,
			Content:   "Comment with attachment",
			CreatedAt: time.Now(),
			Attachments: []Attachment{
				{FileName: "file.pdf", FileSize: 1024, FileURL: "/files/file.pdf"},
			},
		}

		err := handler.CreateComment(ctx, comment)
		require.NoError(t, err)

		// Delete attachment
		err = handler.DeleteAttachment(ctx, comment.ID, comment.Attachments[0].ID)
		require.NoError(t, err)

		retrieved, err := handler.GetComment(ctx, comment.ID)
		require.NoError(t, err)
		assert.Len(t, retrieved.Attachments, 0)
	})
}

func TestCommentPagination(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	handler := NewCommentHandler(db)
	ctx := context.Background()

	t.Run("paginate comments", func(t *testing.T) {
		task := createTestTask(db, "Task 1", time.Now(), time.Now().Add(24*time.Hour))

		// Create 25 comments
		for i := 1; i <= 25; i++ {
			comment := &Comment{
				TaskID:    task.ID,
				UserID:    1,
				Content:   fmt.Sprintf("Comment %d", i),
				CreatedAt: time.Now(),
			}
			err := handler.CreateComment(ctx, comment)
			require.NoError(t, err)
		}

		// Get first page
		page1, err := handler.GetCommentsByTaskPaginated(ctx, task.ID, 1, 10)
		require.NoError(t, err)
		assert.Len(t, page1.Comments, 10)
		assert.True(t, page1.HasMore)

		// Get second page
		page2, err := handler.GetCommentsByTaskPaginated(ctx, task.ID, 2, 10)
		require.NoError(t, err)
		assert.Len(t, page2.Comments, 10)
		assert.True(t, page2.HasMore)

		// Get last page
		page3, err := handler.GetCommentsByTaskPaginated(ctx, task.ID, 3, 10)
		require.NoError(t, err)
		assert.Len(t, page3.Comments, 5)
		assert.False(t, page3.HasMore)
	})
}

func TestCommentSearch(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	handler := NewCommentHandler(db)
	ctx := context.Background()

	t.Run("search comments by content", func(t *testing.T) {
		task := createTestTask(db, "Task 1", time.Now(), time.Now().Add(24*time.Hour))

		comments := []*Comment{
			{TaskID: task.ID, UserID: 1, Content: "This is urgent", CreatedAt: time.Now()},
			{TaskID: task.ID, UserID: 1, Content: "Regular comment", CreatedAt: time.Now()},
			{TaskID: task.ID, UserID: 1, Content: "Another urgent matter", CreatedAt: time.Now()},
		}

		for _, comment := range comments {
			err := handler.CreateComment(ctx, comment)
			require.NoError(t, err)
		}

		// Search for "urgent"
		results, err := handler.SearchComments(ctx, task.ID, "urgent")
		require.NoError(t, err)
		assert.Len(t, results, 2)
	})

	t.Run("search comments by user", func(t *testing.T) {
		task := createTestTask(db, "Task 1", time.Now(), time.Now().Add(24*time.Hour))

		// Create comments from different users
		for i := 1; i <= 3; i++ {
			comment := &Comment{
				TaskID:    task.ID,
				UserID:    i,
				Content:   fmt.Sprintf("Comment from user %d", i),
				CreatedAt: time.Now(),
			}
			err := handler.CreateComment(ctx, comment)
			require.NoError(t, err)
		}

		// Search by user
		results, err := handler.GetCommentsByUser(ctx, task.ID, 1)
		require.NoError(t, err)
		assert.Len(t, results, 1)
	})
}

func TestCommentPermissions(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)

	handler := NewCommentHandler(db)
	ctx := context.Background()

	t.Run("check edit permission", func(t *testing.T) {
		task := createTestTask(db, "Task 1", time.Now(), time.Now().Add(24*time.Hour))

		comment := &Comment{
			TaskID:    task.ID,
			UserID:    1,
			Content:   "Test comment",
			CreatedAt: time.Now(),
		}

		err := handler.CreateComment(ctx, comment)
		require.NoError(t, err)

		// User can edit own comment
		canEdit, err := handler.CanEditComment(ctx, comment.ID, 1)
		require.NoError(t, err)
		assert.True(t, canEdit)

		// Other user cannot edit
		canEdit, err = handler.CanEditComment(ctx, comment.ID, 2)
		require.NoError(t, err)
		assert.False(t, canEdit)
	})

	t.Run("check delete permission", func(t *testing.T) {
		task := createTestTask(db, "Task 1", time.Now(), time.Now().Add(24*time.Hour))

		comment := &Comment{
			TaskID:    task.ID,
			UserID:    1,
			Content:   "Test comment",
			CreatedAt: time.Now(),
		}

		err := handler.CreateComment(ctx, comment)
		require.NoError(t, err)

		// User can delete own comment
		canDelete, err := handler.CanDeleteComment(ctx, comment.ID, 1)
		require.NoError(t, err)
		assert.True(t, canDelete)

		// Admin can delete any comment
		canDelete, err = handler.CanDeleteComment(ctx, comment.ID, 0) // Admin
		require.NoError(t, err)
		assert.True(t, canDelete)
	})
}
