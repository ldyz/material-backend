package notification

import (
	"bytes"
	"context"
	"log"
	"sync"

	"github.com/yourorg/material-backend/backend/internal/api/agent"
	openai "github.com/yourorg/material-backend/backend/pkg/openai"
	"gorm.io/gorm"
)

// VoiceHandler handles voice and AI chat messages over WebSocket
type VoiceHandler struct {
	db                  *gorm.DB
	aiHandler           *agent.AIHandler
	asrServiceURL       string
	conversationRepo    *agent.ConversationRepository
	// 流式识别会话存储（使用 clientID 作为 key，支持同一用户多终端独立会话）
	streamSessions sync.Map // map[clientID]*streamSession
}

// streamSession 存储流式识别会话
type streamSession struct {
	chunks    [][]byte
	mimeType  string
	lastText  string
	mu        sync.Mutex
}

// NewVoiceHandler creates a new voice handler
func NewVoiceHandler(db *gorm.DB, aiHandler *agent.AIHandler, asrServiceURL string) *VoiceHandler {
	return &VoiceHandler{
		db:               db,
		aiHandler:        aiHandler,
		asrServiceURL:    asrServiceURL,
		conversationRepo: agent.NewConversationRepository(db),
	}
}

// HandleVoiceMessage handles incoming voice messages
func (h *VoiceHandler) HandleVoiceMessage(ctx context.Context, client *Client, audioData []byte, mimeType string, userID int) {
	h.HandleVoiceMessageWithHistory(ctx, client, audioData, mimeType, userID, nil)
}

// HandleVoiceMessageWithHistory handles incoming voice messages with conversation history
func (h *VoiceHandler) HandleVoiceMessageWithHistory(ctx context.Context, client *Client, audioData []byte, mimeType string, userID int, history []map[string]interface{}) {
	log.Printf("[VoiceHandler] Processing voice message from user %d, size: %d bytes, mimeType: %s", userID, len(audioData), mimeType)

	// 1. Send processing status
	client.Send(map[string]interface{}{
		"type":    "voice_processing",
		"message": "正在识别语音...",
	})

	// 2. Transcribe audio
	filename := "recording.webm"
	if mimeType == "audio/mp4" || mimeType == "audio/m4a" {
		filename = "recording.m4a"
	}

	log.Printf("[VoiceHandler] Starting transcription for user %d, filename: %s", userID, filename)
	transcript, err := h.aiHandler.TranscribeAudio(ctx, bytes.NewReader(audioData), filename)
	if err != nil {
		log.Printf("[VoiceHandler] Transcription failed for user %d: %v", userID, err)
		// Check if context was canceled
		if ctx.Err() == context.Canceled {
			log.Printf("[VoiceHandler] Context was canceled for user %d", userID)
			client.Send(map[string]interface{}{
				"type":    "error",
				"message": "语音识别超时，请稍后重试",
			})
		} else if ctx.Err() == context.DeadlineExceeded {
			log.Printf("[VoiceHandler] Context deadline exceeded for user %d", userID)
			client.Send(map[string]interface{}{
				"type":    "error",
				"message": "语音识别超时，请稍后重试",
			})
		} else {
			client.Send(map[string]interface{}{
				"type":    "error",
				"message": "语音识别失败: " + err.Error(),
			})
		}
		return
	}

	log.Printf("[VoiceHandler] Transcription result: %s", transcript)

	// 3. Send transcript to client
	log.Printf("[VoiceHandler] Sending voice_transcript to user %d", userID)
	client.Send(map[string]interface{}{
		"type": "voice_transcript",
		"text": transcript,
	})

	// 4. Load conversation history from database
	convHistory := h.loadConversationHistory(userID)

	// 5. Save user message to database
	if h.conversationRepo != nil {
		if err := h.conversationRepo.SaveMessage(int64(userID), "user", transcript, nil); err != nil {
			log.Printf("[VoiceHandler] Failed to save user message: %v", err)
		}
	}

	// 6. Process with AI (streaming)
	h.processAIResponse(ctx, client, transcript, convHistory, userID)
}

// HandleChatMessage handles incoming text chat messages
func (h *VoiceHandler) HandleChatMessage(ctx context.Context, client *Client, message string, history []map[string]interface{}, userID int) {
	log.Printf("[ChatHandler] Processing chat message from user %d: %s", userID, message)

	// Load conversation history from database (ignore client history)
	convHistory := h.loadConversationHistory(userID)

	// Save user message to database
	if h.conversationRepo != nil {
		if err := h.conversationRepo.SaveMessage(int64(userID), "user", message, nil); err != nil {
			log.Printf("[ChatHandler] Failed to save user message: %v", err)
		}
	}

	// Process with AI (streaming)
	h.processAIResponse(ctx, client, message, convHistory, userID)
}

// processAIResponse processes AI response with streaming
func (h *VoiceHandler) processAIResponse(ctx context.Context, client *Client, message string, history []openai.Message, userID int) {
	// Send AI response start
	// 不生成 session_id，让前端自己管理
	log.Printf("[VoiceHandler] Sending ai_response_start to user %d", userID)
	client.Send(map[string]interface{}{
		"type": "ai_response_start",
	})

	// Build request
	req := &agent.AIChatRequest{
		Message:             message,
		ConversationHistory: history,
		UserID:              userID,
		Context: map[string]interface{}{
			"user_id": userID,
		},
	}

	// Check if streaming is supported
	var fullResponse string
	var resp *agent.AIChatResponse
	var err error

	// Try streaming first
	if h.aiHandler != nil {
		resp, err = h.aiHandler.HandleAIChatStream(ctx, req, func(chunk string) {
			fullResponse += chunk
			client.Send(map[string]interface{}{
				"type":    "ai_response_chunk",
				"content": chunk,
				"done":    false,
			})
		})
	} else {
		client.Send(map[string]interface{}{
			"type":    "error",
			"message": "AI 服务未配置",
		})
		return
	}

	// 检查 context 是否已取消或超时
	if ctx.Err() == context.DeadlineExceeded {
		log.Printf("[AIHandler] AI response timeout for user %d", userID)
		client.Send(map[string]interface{}{
			"type":    "error",
			"message": "AI 处理超时，请稍后重试",
		})
		return
	}

	if ctx.Err() == context.Canceled {
		log.Printf("[AIHandler] AI response canceled for user %d", userID)
		client.Send(map[string]interface{}{
			"type":    "error",
			"message": "AI 处理已取消",
		})
		return
	}

	if err != nil {
		log.Printf("[AIHandler] AI response failed: %v", err)
		client.Send(map[string]interface{}{
			"type":    "error",
			"message": "AI 处理失败: " + err.Error(),
		})
		return
	}

	// Send AI response done
	response := fullResponse
	if resp != nil && resp.Message != "" {
		response = resp.Message
	}

	// Save assistant message to database
	if h.conversationRepo != nil && response != "" {
		var toolCallsBytes []byte
		if resp != nil && len(resp.ToolCalls) > 0 {
			// ToolCalls already in JSON format, no need to marshal
		}
		if err := h.conversationRepo.SaveMessage(int64(userID), "assistant", response, toolCallsBytes); err != nil {
			log.Printf("[VoiceHandler] Failed to save assistant message: %v", err)
		}
	}

	log.Printf("[VoiceHandler] Sending ai_response_done to user %d, response length=%d", userID, len(response))
	client.Send(map[string]interface{}{
		"type":       "ai_response_done",
		"message":    response,
		"tool_calls": resp.ToolCalls,
	})

	log.Printf("[AIHandler] AI response completed for user %d", userID)
}

// HandleVoiceStreamChunk handles streaming voice chunk
func (h *VoiceHandler) HandleVoiceStreamChunk(ctx context.Context, client *Client, audioData []byte, mimeType string, userID int) {
	// 使用 clientID 作为 key，支持同一用户多终端独立会话
	sessionI, _ := h.streamSessions.LoadOrStore(client.clientID, &streamSession{
		chunks:   make([][]byte, 0),
		mimeType: mimeType,
	})
	session := sessionI.(*streamSession)

	session.mu.Lock()
	defer session.mu.Unlock()

	// 累积音频
	session.chunks = append(session.chunks, audioData)
	session.mimeType = mimeType

	// 合并所有音频进行识别
	if len(session.chunks) > 0 {
		combined := bytes.Join(session.chunks, nil)

		// 识别
		transcript, err := h.aiHandler.TranscribeAudio(ctx, bytes.NewReader(combined), "recording.webm")
		if err != nil {
			log.Printf("[VoiceHandler] Stream chunk transcription failed for client %s: %v", client.clientID, err)
			return
		}

		// 发送部分识别结果
		if transcript != "" && transcript != session.lastText {
			session.lastText = transcript
			client.Send(map[string]interface{}{
				"type": "voice_transcript_partial",
				"text": transcript,
			})
		}
	}
}

// HandleVoiceStreamEnd handles end of streaming voice
func (h *VoiceHandler) HandleVoiceStreamEnd(ctx context.Context, client *Client, userID int) {
	// 使用 clientID 作为 key，支持同一用户多终端独立会话
	sessionI, ok := h.streamSessions.Load(client.clientID)
	if !ok {
		return
	}
	session := sessionI.(*streamSession)

	session.mu.Lock()
	chunks := session.chunks
	_ = session.mimeType // mimeType not needed for final processing
	session.chunks = nil
	session.lastText = ""
	session.mu.Unlock()

	// 删除会话
	h.streamSessions.Delete(client.clientID)

	// 合并所有音频
	if len(chunks) == 0 {
		return
	}

	combined := bytes.Join(chunks, nil)

	log.Printf("[VoiceHandler] Stream end for client %s (user %d), total size: %d bytes", client.clientID, userID, len(combined))

	// 最终识别
	transcript, err := h.aiHandler.TranscribeAudio(ctx, bytes.NewReader(combined), "recording.webm")
	if err != nil {
		log.Printf("[VoiceHandler] Final transcription failed for client %s: %v", client.clientID, err)
		client.Send(map[string]interface{}{
			"type":    "error",
			"message": "语音识别失败: " + err.Error(),
		})
		return
	}

	// 发送最终识别结果
	client.Send(map[string]interface{}{
		"type": "voice_transcript",
		"text": transcript,
	})

	// Save user message to database
	if h.conversationRepo != nil {
		if err := h.conversationRepo.SaveMessage(int64(userID), "user", transcript, nil); err != nil {
			log.Printf("[VoiceHandler] Failed to save user message: %v", err)
		}
	}

	// Load conversation history from database
	convHistory := h.loadConversationHistory(userID)

	// 处理 AI 回复
	h.processAIResponse(ctx, client, transcript, convHistory, userID)
}

// loadConversationHistory 从数据库加载用户的对话历史
// 始终从数据库加载，确保历史记录按用户隔离
func (h *VoiceHandler) loadConversationHistory(userID int) []openai.Message {
	if h.conversationRepo == nil {
		return nil
	}

	messages, err := h.conversationRepo.GetRecentHistory(int64(userID), 20)
	if err != nil {
		log.Printf("[VoiceHandler] Failed to load conversation history: %v", err)
		return nil
	}

	convHistory := make([]openai.Message, 0, len(messages))
	for _, msg := range messages {
		convHistory = append(convHistory, openai.Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	log.Printf("[VoiceHandler] Loaded %d messages from database for user %d", len(convHistory), userID)
	return convHistory
}
