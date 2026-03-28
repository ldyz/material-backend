package openai

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

// Client is an OpenAI API client
type Client struct {
	APIKey      string
	Model       string
	HTTPClient  *http.Client
	BaseURL     string
	UseXAPIKey  bool // 是否使用 x-api-key 头（百度千帆 Anthropic 兼容 API）
}

// NewClient creates a new OpenAI-compatible client
func NewClient(apiKey, model string) *Client {
	if model == "" {
		model = "gpt-4o"
	}
	return &Client{
		APIKey: apiKey,
		Model:  model,
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		BaseURL:    "https://api.openai.com/v1",
		UseXAPIKey: false,
	}
}

// NewClientWithBaseURL creates a client with custom base URL (for DeepSeek, etc.)
func NewClientWithBaseURL(apiKey, model, baseURL string) *Client {
	if model == "" {
		model = "deepseek-chat"
	}
	if baseURL == "" {
		baseURL = "https://api.deepseek.com/v1"
	}

	// 检测是否是 Anthropic 兼容 API (百度千帆、智谱AI等)
	// 这些 API 使用 x-api-key 头和 /v1/chat/completions 路径
	useXAPIKey := strings.Contains(baseURL, "qianfan.baidubce.com") ||
		strings.Contains(baseURL, "open.bigmodel.cn/api/anthropic")

	// 百度千帆 API 需要更长的超时时间
	timeout := 60 * time.Second
	if strings.Contains(baseURL, "qianfan.baidubce.com") {
		timeout = 600 * time.Second // 10分钟
	}

	return &Client{
		APIKey:     apiKey,
		Model:      model,
		HTTPClient: &http.Client{Timeout: timeout},
		BaseURL:    baseURL,
		UseXAPIKey: useXAPIKey,
	}
}

// Message represents a chat message
type Message struct {
	Role      string     `json:"role"`
	Content   string     `json:"content,omitempty"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string    `json:"tool_call_id,omitempty"`
}

// AnthropicRequestMessage represents a message in Anthropic request format (content always included)
type AnthropicRequestMessage struct {
	Role      string     `json:"role"`
	Content   string     `json:"content"` // 不使用 omitempty，确保总是序列化
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

// ToolCall represents a tool call in a message
type ToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

// Tool represents a function tool
type Tool struct {
	Type     string `json:"type"`
	Function struct {
		Name        string                 `json:"name"`
		Description string                 `json:"description"`
		Parameters  map[string]interface{} `json:"parameters"`
	} `json:"function"`
}

// AnthropicTool represents Anthropic format tool
type AnthropicTool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	InputSchema map[string]interface{} `json:"input_schema"`
}

// AnthropicChatRequest represents Anthropic format chat request
type AnthropicChatRequest struct {
	Model       string                   `json:"model"`
	System      string                   `json:"system,omitempty"`
	Messages    []AnthropicRequestMessage `json:"messages"`
	Tools       []AnthropicTool          `json:"tools,omitempty"`
	MaxTokens   int                      `json:"max_tokens"`
	Temperature float64                  `json:"temperature,omitempty"`
	Stream      bool                     `json:"stream,omitempty"`
}

// ChatRequest represents a chat completion request
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Tools       []Tool    `json:"tools,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
}

// ChatResponse represents a chat completion response
type ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int      `json:"index"`
		Message      Message  `json:"message"`
		FinishReason string   `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
}

// convertToAnthropicFormat converts OpenAI format request to Anthropic format
func convertToAnthropicFormat(req *ChatRequest, model string) *AnthropicChatRequest {
	// Extract system message and filter it from messages
	var systemPrompt string
	messages := make([]AnthropicRequestMessage, 0, len(req.Messages))

	for _, msg := range req.Messages {
		if msg.Role == "system" {
			systemPrompt = msg.Content
		} else if msg.Role == "tool" {
			// Anthropic 不支持 "tool" 角色，将工具结果转换为 user 消息
			// 格式: user 角色包含工具调用结果
			messages = append(messages, AnthropicRequestMessage{
				Role:    "user",
				Content: fmt.Sprintf("Tool result for %s:\n%s", msg.ToolCallID, msg.Content),
			})
		} else if msg.Role == "assistant" {
			// assistant 消息，确保 content 字段存在
			messages = append(messages, AnthropicRequestMessage{
				Role:      "assistant",
				Content:   msg.Content, // 即使为空也会被序列化
				ToolCalls: msg.ToolCalls,
			})
		} else {
			messages = append(messages, AnthropicRequestMessage{
				Role:    msg.Role,
				Content: msg.Content,
			})
		}
	}

	// Convert tools to Anthropic format
	tools := make([]AnthropicTool, 0, len(req.Tools))
	for _, tool := range req.Tools {
		tools = append(tools, AnthropicTool{
			Name:        tool.Function.Name,
			Description: tool.Function.Description,
			InputSchema: tool.Function.Parameters,
		})
	}

	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = 4096
	}

	return &AnthropicChatRequest{
		Model:       model,
		System:      systemPrompt,
		Messages:    messages,
		Tools:       tools,
		MaxTokens:   maxTokens,
		Temperature: req.Temperature,
		Stream:      req.Stream,
	}
}

// ChatCompletion sends a chat completion request
func (c *Client) ChatCompletion(req *ChatRequest) (*ChatResponse, error) {
	req.Model = c.Model

	var body []byte
	var err error

	if c.UseXAPIKey {
		// 转换为 Anthropic 格式
		anthropicReq := convertToAnthropicFormat(req, c.Model)
		body, err = json.Marshal(anthropicReq)
	} else {
		body, err = json.Marshal(req)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// 构建请求 URL
	reqURL := c.BaseURL
	if c.UseXAPIKey {
		// 百度千帆 Anthropic 兼容 API: BaseURL 已经包含完整路径
		// 例如: https://qianfan.baidubce.com/anthropic/coding
		// 需要添加 /v1/messages (Anthropic 风格路径)
		reqURL = c.BaseURL + "/v1/messages"
	} else {
		// OpenAI/DeepSeek: BaseURL + /chat/completions
		reqURL = c.BaseURL + "/chat/completions"
	}

	log.Printf("[Client] Request URL: %s", reqURL)

	httpReq, err := http.NewRequest("POST", reqURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// 根据配置选择认证方式
	if c.UseXAPIKey {
		// 百度千帆 Anthropic 兼容 API 使用 x-api-key 头
		httpReq.Header.Set("x-api-key", c.APIKey)
	} else {
		// OpenAI/DeepSeek 使用 Authorization: Bearer 头
		httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)
	}

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if chatResp.Error != nil {
		return nil, fmt.Errorf("OpenAI API error: %s", chatResp.Error.Message)
	}

	return &chatResp, nil
}

// Chat sends a simple chat message
func (c *Client) Chat(systemPrompt, userMessage string) (string, error) {
	req := &ChatRequest{
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userMessage},
		},
		Temperature: 0.7,
	}

	resp, err := c.ChatCompletion(req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response choices")
	}

	return resp.Choices[0].Message.Content, nil
}

// ChatWithTools sends a chat message with function calling support
func (c *Client) ChatWithTools(systemPrompt string, messages []Message, tools []Tool) (*ChatResponse, error) {
	req := &ChatRequest{
		Messages:    append([]Message{{Role: "system", Content: systemPrompt}}, messages...),
		Tools:       tools,
		Temperature: 0.7,
	}

	return c.ChatCompletion(req)
}

// TranscriptionResponse represents a transcription response from Whisper API
type TranscriptionResponse struct {
	Text string `json:"text"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
}

// TranscribeAudio calls OpenAI Whisper API to transcribe audio to text
func (c *Client) TranscribeAudio(audioFile io.Reader, filename string) (string, error) {
	// Create a pipe for multipart form data
	body, contentType, err := createMultipartForm(audioFile, filename)
	if err != nil {
		return "", fmt.Errorf("failed to create multipart form: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.BaseURL+"/audio/transcriptions", body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", contentType)
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)

	// Increase timeout for audio transcription
	c.HTTPClient.Timeout = 120 * time.Second

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var transcription TranscriptionResponse
	if err := json.Unmarshal(respBody, &transcription); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if transcription.Error != nil {
		return "", fmt.Errorf("OpenAI API error: %s", transcription.Error.Message)
	}

	return transcription.Text, nil
}

// createMultipartForm creates a multipart form body for audio transcription
func createMultipartForm(audioFile io.Reader, filename string) (*bytes.Buffer, string, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Add model field
	if err := writer.WriteField("model", "whisper-1"); err != nil {
		return nil, "", fmt.Errorf("failed to write model field: %w", err)
	}

	// Add language field for Chinese
	if err := writer.WriteField("language", "zh"); err != nil {
		return nil, "", fmt.Errorf("failed to write language field: %w", err)
	}

	// Add audio file
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(part, audioFile); err != nil {
		return nil, "", fmt.Errorf("failed to copy audio file: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, "", fmt.Errorf("failed to close writer: %w", err)
	}

	return &body, writer.FormDataContentType(), nil
}

// StreamChoice represents a choice in a streaming response (with delta instead of message)
type StreamChoice struct {
	Index        int          `json:"index"`
	Delta        Message      `json:"delta"`        // Streaming uses delta instead of message
	Message      Message      `json:"message"`      // Some APIs use message
	FinishReason string       `json:"finish_reason"`
}

// StreamChatResponse represents a streaming chat completion response
type StreamChatResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []StreamChoice `json:"choices"`
	Error   *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
}

// Anthropic streaming response types
type AnthropicStreamEvent struct {
	Type         string                   `json:"type"`
	Message      *AnthropicMessage        `json:"message,omitempty"`
	Index        int                      `json:"index,omitempty"`
	ContentBlock *AnthropicContentBlock   `json:"content_block,omitempty"`
	Delta        *AnthropicDelta          `json:"delta,omitempty"`
	Usage        *AnthropicUsage          `json:"usage,omitempty"`
}

type AnthropicMessage struct {
	ID      string                  `json:"id"`
	Type    string                  `json:"type"`
	Role    string                  `json:"role"`
	Content []AnthropicContentBlock `json:"content"`
	Model   string                  `json:"model"`
	Usage   AnthropicUsage          `json:"usage"`
	StopReason string               `json:"stop_reason,omitempty"`
}

type AnthropicContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Input map[string]interface{} `json:"input,omitempty"`
}

type AnthropicDelta struct {
	Type string `json:"type,omitempty"`
	Text string `json:"text,omitempty"`
	StopReason string `json:"stop_reason,omitempty"`
	PartialJSON string `json:"partial_json,omitempty"`
}

type AnthropicUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// ChatCompletionStream sends a streaming chat completion request
// The onChunk callback is called for each chunk of the response
func (c *Client) ChatCompletionStream(req *ChatRequest, onChunk func(chunk string)) (*ChatResponse, error) {
	req.Model = c.Model
	req.Stream = true // Enable streaming

	var body []byte
	var err error

	if c.UseXAPIKey {
		// 转换为 Anthropic 格式
		anthropicReq := convertToAnthropicFormat(req, c.Model)
		anthropicReq.Stream = true
		body, err = json.Marshal(anthropicReq)
	} else {
		body, err = json.Marshal(req)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// 构建请求 URL
	reqURL := c.BaseURL
	if c.UseXAPIKey {
		// 百度千帆 Anthropic 兼容 API: BaseURL 已经包含完整路径
		// 例如: https://qianfan.baidubce.com/anthropic/coding
		// 需要添加 /v1/messages (Anthropic 风格路径)
		reqURL = c.BaseURL + "/v1/messages"
	} else {
		// OpenAI/DeepSeek: BaseURL + /chat/completions
		reqURL = c.BaseURL + "/chat/completions"
	}

	log.Printf("[Client] Stream Request URL: %s", reqURL)

	httpReq, err := http.NewRequest("POST", reqURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")

	// 根据配置选择认证方式
	if c.UseXAPIKey {
		// 百度千帆 Anthropic 兼容 API 使用 x-api-key 头
		httpReq.Header.Set("x-api-key", c.APIKey)
	} else {
		// OpenAI/DeepSeek 使用 Authorization: Bearer 头
		httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)
	}

	// Increase timeout for streaming
	// 百度千帆 API 需要更长的超时时间
	if c.UseXAPIKey {
		c.HTTPClient.Timeout = 600 * time.Second // 10分钟
	} else {
		c.HTTPClient.Timeout = 120 * time.Second
	}

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(respBody))
	}

	// 根据API类型选择不同的解析方式
	if c.UseXAPIKey {
		return c.parseAnthropicStream(resp.Body, onChunk)
	}
	return c.parseOpenAIStream(resp.Body, onChunk)
}

// parseAnthropicStream 解析 Anthropic 格式的流式响应
func (c *Client) parseAnthropicStream(body io.Reader, onChunk func(chunk string)) (*ChatResponse, error) {
	var fullContent string
	var chatResp ChatResponse
	var currentToolCalls []ToolCall
	var currentToolCallIndex int = -1

	scanner := bufio.NewScanner(body)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines
		if line == "" {
			continue
		}

		// Check for data prefix
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")

		// Parse the JSON chunk
		var event AnthropicStreamEvent
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			log.Printf("[Anthropic] Failed to parse event: %v", err)
			continue
		}

		switch event.Type {
		case "message_start":
			// Initialize message
			if event.Message != nil {
				chatResp.ID = event.Message.ID
				chatResp.Model = event.Message.Model
			}

		case "content_block_start":
			// Start of a content block (text or tool_use)
			if event.ContentBlock != nil {
				if event.ContentBlock.Type == "tool_use" {
					tc := ToolCall{
						ID:   event.ContentBlock.ID,
						Type: "function",
					}
					tc.Function.Name = event.ContentBlock.Name
					currentToolCalls = append(currentToolCalls, tc)
					currentToolCallIndex = len(currentToolCalls) - 1
				}
			}

		case "content_block_delta":
			// Text or tool input delta
			if event.Delta != nil {
				if event.Delta.Type == "text_delta" && event.Delta.Text != "" {
					fullContent += event.Delta.Text
					if onChunk != nil {
						onChunk(event.Delta.Text)
					}
				} else if event.Delta.Type == "input_json_delta" && event.Delta.PartialJSON != "" {
					// Tool input delta
					if currentToolCallIndex >= 0 && currentToolCallIndex < len(currentToolCalls) {
						currentToolCalls[currentToolCallIndex].Function.Arguments += event.Delta.PartialJSON
					}
				}
			}

		case "content_block_stop":
			// End of content block
			currentToolCallIndex = -1

		case "message_delta":
			// Message-level changes
			if event.Delta != nil && event.Delta.StopReason != "" {
				chatResp.Choices = []struct {
					Index        int     `json:"index"`
					Message      Message `json:"message"`
					FinishReason string  `json:"finish_reason"`
				}{
					{
						Index:        0,
						FinishReason: event.Delta.StopReason,
					},
				}
			}

		case "message_stop":
			// End of message
			// Finalize the response
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading stream: %w", err)
	}

	// Build final response
	if len(chatResp.Choices) == 0 {
		chatResp.Choices = []struct {
			Index        int     `json:"index"`
			Message      Message `json:"message"`
			FinishReason string  `json:"finish_reason"`
		}{
			{
				Index:        0,
				FinishReason: "end_turn",
			},
		}
	}

	chatResp.Choices[0].Message = Message{
		Role:      "assistant",
		Content:   fullContent,
		ToolCalls: currentToolCalls,
	}

	log.Printf("[Anthropic] Stream completed. Content length: %d, ToolCalls: %d", len(fullContent), len(currentToolCalls))

	return &chatResp, nil
}

// parseOpenAIStream 解析 OpenAI 格式的流式响应
func (c *Client) parseOpenAIStream(body io.Reader, onChunk func(chunk string)) (*ChatResponse, error) {
	// Parse SSE stream
	var fullContent string
	var chatResp ChatResponse
	var toolCalls []ToolCall // 累积工具调用
	chunkCount := 0

	scanner := bufio.NewScanner(body)
	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines
		if line == "" {
			continue
		}

		// Check for data prefix
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")

		// Check for stream end
		if data == "[DONE]" {
			break
		}

		// Parse the JSON chunk using StreamChatResponse for proper delta handling
		var streamResp StreamChatResponse
		if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
			continue
		}

		if streamResp.Error != nil {
			return nil, fmt.Errorf("OpenAI API error: %s", streamResp.Error.Message)
		}

		if len(streamResp.Choices) > 0 {
			choice := streamResp.Choices[0]
			// In streaming mode, content is in Delta.Content, not Message.Content
			delta := choice.Delta.Content
			if delta != "" && onChunk != nil {
				onChunk(delta)
				chunkCount++
			}
			fullContent += delta

			// 处理工具调用 - 累积 Delta.ToolCalls
			if len(choice.Delta.ToolCalls) > 0 {
				for _, tc := range choice.Delta.ToolCalls {
					// 查找是否已存在该 tool call（通过 ID 或索引匹配）
					found := false
					for i, existing := range toolCalls {
						// 如果 ID 匹配，或者 ID 为空时通过索引匹配
						if (tc.ID != "" && existing.ID == tc.ID) || (tc.ID == "" && i == len(toolCalls)-1) {
							// 累积 arguments
							toolCalls[i].Function.Arguments += tc.Function.Arguments
							// 如果有新的 ID，更新它
							if tc.ID != "" {
								toolCalls[i].ID = tc.ID
							}
							// 如果有新的 function name，更新它
							if tc.Function.Name != "" {
								toolCalls[i].Function.Name = tc.Function.Name
							}
							// 如果有新的 type，更新它
							if tc.Type != "" {
								toolCalls[i].Type = tc.Type
							}
							found = true
							break
						}
					}
					if !found {
						// 新的 tool call
						toolCalls = append(toolCalls, tc)
					}
				}
			}

			// Store the final response structure
			if chatResp.ID == "" {
				chatResp.ID = streamResp.ID
				chatResp.Model = streamResp.Model
				chatResp.Created = streamResp.Created
				chatResp.Choices = make([]struct {
					Index        int     `json:"index"`
					Message      Message `json:"message"`
					FinishReason string  `json:"finish_reason"`
				}, 1)
			}
			chatResp.Choices[0].FinishReason = choice.FinishReason
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading stream: %w", err)
	}

	// Set the full content in the response
	if len(chatResp.Choices) > 0 {
		chatResp.Choices[0].Message.Content = fullContent
		chatResp.Choices[0].Message.Role = "assistant"
		chatResp.Choices[0].Message.ToolCalls = toolCalls
	} else {
		// Create a basic response if none was created from stream
		chatResp = ChatResponse{
			Choices: []struct {
				Index        int     `json:"index"`
				Message      Message `json:"message"`
				FinishReason string  `json:"finish_reason"`
			}{
				{
					Index: 0,
					Message: Message{
						Role:      "assistant",
						Content:   fullContent,
						ToolCalls: toolCalls,
					},
					FinishReason: "stop",
				},
			},
		}
	}

	return &chatResp, nil
}

// ChatWithToolsStream sends a streaming chat message with function calling support
func (c *Client) ChatWithToolsStream(systemPrompt string, messages []Message, tools []Tool, onChunk func(chunk string)) (*ChatResponse, error) {
	req := &ChatRequest{
		Messages:    append([]Message{{Role: "system", Content: systemPrompt}}, messages...),
		Tools:       tools,
		Temperature: 0.7,
	}

	return c.ChatCompletionStream(req, onChunk)
}
