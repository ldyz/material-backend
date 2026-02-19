<template>
  <div class="comment-item" :class="{ 'is-resolved': comment.is_resolved, 'is-reply': isReply }">
    <!-- Comment Header -->
    <div class="comment-header">
      <div class="comment-author">
        <el-avatar :size="32" :src="user.avatar">
          {{ user.name?.charAt(0) || '?' }}
        </el-avatar>
        <div class="author-info">
          <span class="author-name">{{ user.name }}</span>
          <span class="comment-date">{{ formatDate(comment.created_at) }}</span>
        </div>
      </div>

      <el-dropdown trigger="click" @command="handleCommand">
        <el-button :icon="MoreFilled" circle size="small" text />
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="reply" :icon="ChatLineSquare">
              Reply
            </el-dropdown-item>
            <el-dropdown-item v-if="canEdit" command="edit" :icon="Edit">
              Edit
            </el-dropdown-item>
            <el-dropdown-item
              v-if="!comment.is_resolved"
              command="resolve"
              :icon="Select"
            >
              Mark Resolved
            </el-dropdown-item>
            <el-dropdown-item
              v-if="comment.is_resolved"
              command="unresolve"
              :icon="CircleClose"
            >
              Mark Unresolved
            </el-dropdown-item>
            <el-dropdown-item
              v-if="canDelete"
              command="delete"
              :icon="Delete"
              divided
              style="color: #f56c6c"
            >
              Delete
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>

    <!-- Comment Content -->
    <div class="comment-content">
      <div
        v-if="!isEditing"
        class="content-html"
        v-html="comment.content"
      />
      <el-input
        v-else
        v-model="editContent"
        type="textarea"
        :rows="4"
        placeholder="Edit your comment..."
      />
    </div>

    <!-- Mentions -->
    <div v-if="comment.mentions?.length" class="comment-mentions">
      <el-tag
        v-for="mentionId in comment.mentions"
        :key="mentionId"
        size="small"
        type="info"
      >
        @{{ getMentionName(mentionId) }}
      </el-tag>
    </div>

    <!-- Edit Actions -->
    <div v-if="isEditing" class="edit-actions">
      <el-button size="small" @click="cancelEdit">Cancel</el-button>
      <el-button
        type="primary"
        size="small"
        :loading="saving"
        @click="saveEdit"
      >
        Save
      </el-button>
    </div>

    <!-- Reply Actions -->
    <div v-if="showReplyInput" class="reply-input">
      <el-input
        v-model="replyText"
        type="textarea"
        :rows="3"
        placeholder="Write a reply..."
      />
      <div class="reply-actions">
        <el-button size="small" @click="cancelReply">Cancel</el-button>
        <el-button
          type="primary"
          size="small"
          :loading="saving"
          @click="submitReply"
        >
          Reply
        </el-button>
      </div>
    </div>

    <!-- Thread Replies -->
    <div v-if="comment.replies?.length" class="comment-replies">
      <CommentItem
        v-for="reply in comment.replies"
        :key="reply.id"
        :comment="reply"
        :user="getUser(reply.user_id)"
        :is-reply="true"
        @reply="$emit('reply', $event)"
        @resolve="$emit('resolve', $event)"
        @delete="$emit('delete', $event)"
        @edit="$emit('edit', $event)"
      />
    </div>

    <!-- Resolved Badge -->
    <div v-if="comment.is_resolved" class="resolved-badge">
      <el-tag type="success" size="small">Resolved</el-tag>
    </div>
  </div>
</template>

<script setup>
/**
 * CommentItem.vue
 *
 * Individual comment component with editing, resolving, and threading support.
 */

import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import {
  MoreFilled,
  ChatLineSquare,
  Edit,
  Delete,
  Select,
  CircleClose
} from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { formatDistanceToNow } from 'date-fns'

// ==================== Props ====================

const props = defineProps({
  comment: {
    type: Object,
    required: true
  },
  user: {
    type: Object,
    required: true
  },
  isReply: {
    type: Boolean,
    default: false
  },
  isReplying: {
    type: Boolean,
    default: false
  },
  projectUsers: {
    type: Array,
    default: () => []
  }
})

// ==================== Emits ====================

const emit = defineEmits(['reply', 'resolve', 'delete', 'edit'])

// ==================== State ====================

const authStore = useAuthStore()
const isEditing = ref(false)
const showReplyInput = ref(props.isReplying)
const editContent = ref('')
const replyText = ref('')
const saving = ref(false)

// ==================== Computed ====================

/**
 * Check if current user can edit this comment
 */
const canEdit = computed(() => {
  return props.comment.user_id === authStore.user?.id
})

/**
 * Check if current user can delete this comment
 */
const canDelete = computed(() => {
  return props.comment.user_id === authStore.user?.id || authStore.isAdmin
})

// ==================== Methods ====================

/**
 * Format date relative to now
 */
function formatDate(dateString) {
  try {
    return formatDistanceToNow(new Date(dateString), { addSuffix: true })
  } catch {
    return 'Unknown date'
  }
}

/**
 * Get mention name by ID
 */
function getMentionName(userId) {
  const user = props.projectUsers.find(u => u.id === userId)
  return user?.name || 'Unknown'
}

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
 * Handle dropdown command
 */
function handleCommand(command) {
  switch (command) {
    case 'reply':
      showReplyInput.value = true
      break
    case 'edit':
      startEditing()
      break
    case 'resolve':
      emit('resolve', props.comment.id)
      break
    case 'unresolve':
      emit('resolve', props.comment.id)
      break
    case 'delete':
      emit('delete', props.comment.id)
      break
  }
}

/**
 * Start editing comment
 */
function startEditing() {
  isEditing.value = true
  editContent.value = stripHtml(props.comment.content)
}

/**
 * Cancel editing
 */
function cancelEdit() {
  isEditing.value = false
  editContent.value = ''
}

/**
 * Save edit
 */
function saveEdit() {
  if (!editContent.value.trim()) {
    ElMessage.warning('Comment cannot be empty')
    return
  }

  emit('edit', props.comment.id, editContent.value)
  isEditing.value = false
}

/**
 * Cancel reply
 */
function cancelReply() {
  showReplyInput.value = false
  replyText.value = ''
}

/**
 * Submit reply
 */
async function submitReply() {
  if (!replyText.value.trim()) {
    ElMessage.warning('Reply cannot be empty')
    return
  }

  saving.value = true
  try {
    // Emit reply event to parent component
    emit('reply', props.comment.id)
    replyText.value = ''
    showReplyInput.value = false
  } finally {
    saving.value = false
  }
}

/**
 * Strip HTML tags from string
 */
function stripHtml(html) {
  const tmp = document.createElement('div')
  tmp.innerHTML = html
  return tmp.textContent || tmp.innerText || ''
}
</script>

<style scoped lang="scss">
.comment-item {
  position: relative;
  padding: 12px;
  background: #f5f7fa;
  border-radius: 8px;
  border-left: 3px solid #409eff;

  &.is-resolved {
    opacity: 0.7;
    border-left-color: #67c23a;
  }

  &.is-reply {
    margin-left: 48px;
    background: #fafafa;
    border-left-color: #909399;
  }
}

.comment-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.comment-author {
  display: flex;
  align-items: center;
  gap: 8px;
}

.author-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.author-name {
  font-size: 14px;
  font-weight: 600;
  color: #303133;
}

.comment-date {
  font-size: 12px;
  color: #909399;
}

.comment-content {
  margin-bottom: 8px;
  color: #606266;
  line-height: 1.6;
}

.content-html {
  :deep(p) {
    margin: 0 0 8px;

    &:last-child {
      margin-bottom: 0;
    }
  }

  :deep(a) {
    color: #409eff;
    text-decoration: none;

    &:hover {
      text-decoration: underline;
    }
  }
}

.comment-mentions {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
  margin-bottom: 8px;
}

.edit-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 12px;
}

.reply-input {
  margin-top: 12px;
  padding: 12px;
  background: #fff;
  border-radius: 4px;
}

.reply-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 8px;
}

.comment-replies {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid #e4e7ed;
}

.resolved-badge {
  position: absolute;
  top: 12px;
  right: 48px;
}
</style>
