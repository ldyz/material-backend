<template>
  <div class="comments-panel">
    <!-- Header -->
    <div class="comments-header">
      <h3>Comments</h3>
      <el-badge :value="unresolvedCount" :hidden="unresolvedCount === 0" type="danger">
        <el-button
          :icon="Message"
          circle
          size="small"
          @click="scrollToNewComment"
        />
      </el-badge>
    </div>

    <!-- Comments List -->
    <el-scrollbar ref="scrollbarRef" class="comments-list">
      <div v-if="loading" class="comments-loading">
        <el-skeleton :rows="3" animated />
      </div>

      <div v-else-if="comments.length === 0" class="comments-empty">
        <el-empty description="No comments yet" :image-size="80">
          <p class="empty-hint">Be the first to comment on this task</p>
        </el-empty>
      </div>

      <div v-else class="comments-container">
        <CommentItem
          v-for="comment in sortedComments"
          :key="comment.id"
          :comment="comment"
          :user="getUser(comment.user_id)"
          :is-replying="replyingTo === comment.id"
          @reply="onReply"
          @resolve="onResolve"
          @delete="onDelete"
          @edit="onEdit"
        />
      </div>
    </el-scrollbar>

    <!-- Comment Input -->
    <div class="comments-input">
      <div v-if="replyingTo" class="replying-to">
        <span>Replying to {{ getReplyingToUser().name }}</span>
        <el-button
          :icon="Close"
          size="small"
          text
          @click="cancelReply"
        />
      </div>

      <RichTextEditor
        ref="editorRef"
        v-model="newComment"
        :placeholder="replyingTo ? 'Write a reply...' : 'Write a comment...'"
        :mentions="availableMentions"
        :min-height="80"
        :max-height="200"
      />

      <div class="comments-actions">
        <div class="mentions-info">
          <el-tag
            v-for="mention in newCommentMentions"
            :key="mention.id"
            size="small"
            closable
            @close="removeMention(mention.id)"
          >
            @{{ mention.name }}
          </el-tag>
        </div>

        <el-button
          type="primary"
          :loading="submitting"
          :disabled="!canSubmit"
          @click="submitComment"
        >
          {{ replyingTo ? 'Reply' : 'Comment' }}
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
/**
 * CommentsPanel.vue
 *
 * Task comments panel with threading, @mentions, and rich text editing.
 * Integrates with backend API for real-time comment management.
 */

import { ref, computed, watch, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Message, Close } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { progressApi } from '@/api'
import RichTextEditor from '@/components/common/RichTextEditor.vue'
import CommentItem from './CommentItem.vue'

// ==================== Props ====================

const props = defineProps({
  taskId: {
    type: Number,
    required: true
  },
  projectUsers: {
    type: Array,
    default: () => []
  }
})

// ==================== Emits ====================

const emit = defineEmits(['comment-count-change'])

// ==================== State ====================

const authStore = useAuthStore()
const scrollbarRef = ref(null)
const editorRef = ref(null)

const loading = ref(false)
const submitting = ref(false)
const comments = ref([])
const replyingTo = ref(null)
const newComment = ref('')

// ==================== Computed ====================

/**
 * Sort comments by date (newest first for root, replies by thread)
 */
const sortedComments = computed(() => {
  const roots = comments.value.filter(c => !c.parent_id)
  const replies = comments.value.filter(c => c.parent_id)

  // Attach replies to their parents
  return roots.map(root => {
    const threadReplies = replies
      .filter(r => r.parent_id === root.id)
      .sort((a, b) => new Date(a.created_at) - new Date(b.created_at))

    return {
      ...root,
      replies: threadReplies
    }
  }).sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
})

/**
 * Count unresolved comments
 */
const unresolvedCount = computed(() => {
  return comments.value.filter(c => !c.is_resolved && !c.parent_id).length
})

/**
 * Get available users for @mentions
 */
const availableMentions = computed(() => {
  return props.projectUsers.map(user => ({
    id: user.id,
    name: user.name,
    avatar: user.avatar
  }))
})

/**
 * Extract mentions from new comment
 */
const newCommentMentions = computed(() => {
  const mentionRegex = /@(\w+)/g
  const mentions = []
  let match

  while ((match = mentionRegex.exec(newComment.value)) !== null) {
    const user = availableMentions.value.find(
      u => u.name.toLowerCase() === match[1].toLowerCase()
    )
    if (user && !mentions.find(m => m.id === user.id)) {
      mentions.push(user)
    }
  }

  return mentions
})

/**
 * Check if comment can be submitted
 */
const canSubmit = computed(() => {
  return newComment.value.trim().length > 0 && !submitting.value
})

// ==================== Methods ====================

/**
 * Get user by ID
 */
function getUser(userId) {
  return props.projectUsers.find(u => u.id === userId) || {
    name: 'Unknown User',
    avatar: null
  }
}

/**
 * Get user we're replying to
 */
function getReplyingToUser() {
  if (!replyingTo.value) return null
  return getUser(comments.value.find(c => c.id === replyingTo.value)?.user_id)
}

/**
 * Load comments for task
 */
async function loadComments() {
  if (!props.taskId) return

  loading.value = true
  try {
    const response = await progressApi.getComments(props.taskId)
    comments.value = response.data || []
    emit('comment-count-change', comments.value.length)
  } catch (error) {
    console.error('[CommentsPanel] Error loading comments:', error)
    ElMessage.error('Failed to load comments')
  } finally {
    loading.value = false
  }
}

/**
 * Submit new comment or reply
 */
async function submitComment() {
  if (!canSubmit.value) return

  submitting.value = true
  try {
    const commentData = {
      content: newComment.value,
      mentions: newCommentMentions.value.map(u => u.id),
      parent_id: replyingTo.value || undefined
    }

    const response = await progressApi.createComment(props.taskId, commentData)

    // Add to local list
    comments.value.push({
      ...response.data,
      user: authStore.user,
      created_at: new Date().toISOString()
    })

    // Clear input
    newComment.value = ''
    replyingTo.value = null

    // Scroll to bottom
    nextTick(() => {
      scrollToBottom()
    })

    ElMessage.success('Comment added')

    // Broadcast via WebSocket if available
    if (window.collaborationStore?.isConnected) {
      window.collaborationStore.emit('comment:create', {
        taskId: props.taskId,
        commentId: response.data.id
      })
    }
  } catch (error) {
    console.error('[CommentsPanel] Error submitting comment:', error)
    ElMessage.error('Failed to submit comment')
  } finally {
    submitting.value = false
  }
}

/**
 * Handle reply action
 */
function onReply(commentId) {
  replyingTo.value = commentId
  editorRef.value?.focus()
}

/**
 * Cancel reply
 */
function cancelReply() {
  replyingTo.value = null
}

/**
 * Handle resolve action
 */
async function onResolve(commentId) {
  try {
    await progressApi.resolveComment(commentId)

    const comment = comments.value.find(c => c.id === commentId)
    if (comment) {
      comment.is_resolved = true
    }

    ElMessage.success('Comment marked as resolved')
  } catch (error) {
    console.error('[CommentsPanel] Error resolving comment:', error)
    ElMessage.error('Failed to resolve comment')
  }
}

/**
 * Handle delete action
 */
async function onDelete(commentId) {
  try {
    await ElMessageBox.confirm(
      'Are you sure you want to delete this comment?',
      'Delete Comment',
      {
        type: 'warning'
      }
    )

    await progressApi.deleteComment(commentId)

    const index = comments.value.findIndex(c => c.id === commentId)
    if (index !== -1) {
      comments.value.splice(index, 1)
    }

    ElMessage.success('Comment deleted')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('[CommentsPanel] Error deleting comment:', error)
      ElMessage.error('Failed to delete comment')
    }
  }
}

/**
 * Handle edit action
 */
async function onEdit(commentId, newContent) {
  try {
    await progressApi.updateComment(commentId, { content: newContent })

    const comment = comments.value.find(c => c.id === commentId)
    if (comment) {
      comment.content = newContent
    }

    ElMessage.success('Comment updated')
  } catch (error) {
    console.error('[CommentsPanel] Error updating comment:', error)
    ElMessage.error('Failed to update comment')
  }
}

/**
 * Remove mention from tag
 */
function removeMention(userId) {
  const user = availableMentions.value.find(u => u.id === userId)
  if (user) {
    newComment.value = newComment.value.replace(
      new RegExp(`@${user.name}`, 'g'),
      ''
    )
  }
}

/**
 * Scroll to bottom of comments
 */
function scrollToBottom() {
  if (scrollbarRef.value) {
    scrollbarRef.value.setScrollTop(10000)
  }
}

/**
 * Scroll to new comment
 */
function scrollToNewComment() {
  scrollToBottom()
}

// ==================== Lifecycle ====================

onMounted(() => {
  loadComments()
})

// Watch for task changes
watch(() => props.taskId, () => {
  loadComments()
})

// Listen for real-time updates
if (window.collaborationStore) {
  window.collaborationStore.on('comment:create', (data) => {
    if (data.taskId === props.taskId) {
      loadComments()
    }
  })
}
</script>

<style scoped lang="scss">
.comments-panel {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: #fff;
  border-radius: 4px;
}

.comments-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  border-bottom: 1px solid #e4e7ed;

  h3 {
    margin: 0;
    font-size: 16px;
    font-weight: 600;
    color: #303133;
  }
}

.comments-list {
  flex: 1;
  padding: 16px;
}

.comments-loading {
  padding: 16px;
}

.comments-empty {
  padding: 40px 16px;
  text-align: center;

  .empty-hint {
    margin-top: 8px;
    font-size: 12px;
    color: #909399;
  }
}

.comments-container {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.comments-input {
  border-top: 1px solid #e4e7ed;
  padding: 16px;
}

.replying-to {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  margin-bottom: 12px;
  background: #f5f7fa;
  border-radius: 4px;
  font-size: 12px;
  color: #606266;

  .el-button {
    margin-left: 8px;
  }
}

.comments-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 12px;
}

.mentions-info {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}
</style>
