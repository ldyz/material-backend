package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// BaiduClient 百度千帆 API 客户端
type BaiduClient struct {
	APIKey      string
	SecretKey   string
	Model       string
	HTTPClient  *http.Client
	BaseURL     string
	accessToken string       // 缓存的 access token
	tokenExpire time.Time    // token 过期时间
}

// NewBaiduClient 创建百度千帆客户端
func NewBaiduClient(apiKey, secretKey, model, baseURL string) *BaiduClient {
	if model == "" {
		model = "coding-plan"
	}
	if baseURL == "" {
		baseURL = "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop/chat"
	}
	return &BaiduClient{
		APIKey:     apiKey,
		SecretKey:  secretKey,
		Model:      model,
		HTTPClient: &http.Client{Timeout: 120 * time.Second},
		BaseURL:    baseURL,
	}
}

// BaiduAccessTokenResponse 获取 access token 的响应
type BaiduAccessTokenResponse struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshToken     string `json:"refresh_token"`
	Scope            string `json:"scope"`
	SessionKey       string `json:"session_key"`
	SessionSecret    string `json:"session_secret"`
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
}

// getAccessToken 获取 access token
func (c *BaiduClient) getAccessToken() (string, error) {
	// 检查缓存的 token 是否有效
	if c.accessToken != "" && time.Now().Before(c.tokenExpire) {
		return c.accessToken, nil
	}

	// 构建请求 URL
	tokenURL := fmt.Sprintf("https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=%s&client_secret=%s",
		c.APIKey, c.SecretKey)

	resp, err := c.HTTPClient.Post(tokenURL, "application/json", nil)
	if err != nil {
		return "", fmt.Errorf("获取 access token 失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	var tokenResp BaiduAccessTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if tokenResp.Error != "" {
		return "", fmt.Errorf("获取 access token 失败: %s - %s", tokenResp.Error, tokenResp.ErrorDescription)
	}

	// 缓存 token（提前 5 分钟过期，确保 token 有效）
	c.accessToken = tokenResp.AccessToken
	c.tokenExpire = time.Now().Add(time.Duration(tokenResp.ExpiresIn-300) * time.Second)

	log.Printf("[BaiduClient] 获取 access token 成功，有效期: %d 秒", tokenResp.ExpiresIn)
	return c.accessToken, nil
}

// BaiduChatRequest 百度千帆聊天请求
type BaiduChatRequest struct {
	Messages    []BaiduMessage `json:"messages"`
	Temperature float64        `json:"temperature,omitempty"`
	TopP        float64        `json:"top_p,omitempty"`
	Stream      bool           `json:"stream,omitempty"`
	MaxOutput   int            `json:"max_output_tokens,omitempty"`
}

// BaiduMessage 百度千帆消息格式
type BaiduMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// BaiduChatResponse 百度千帆聊天响应
type BaiduChatResponse struct {
	ID               string `json:"id"`
	Object           string `json:"object"`
	Created          int64  `json:"created"`
	Result           string `json:"result"`             // 百度格式的响应内容
	IsTruncated      bool   `json:"is_truncated"`
	NeedClearHistory bool   `json:"need_clear_history"`
	Usage            struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// ChatCompletion 发送聊天请求
func (c *BaiduClient) ChatCompletion(req *ChatRequest) (*ChatResponse, error) {
	// 获取 access token
	accessToken, err := c.getAccessToken()
	if err != nil {
		return nil, err
	}

	// 转换消息格式
	baiduMessages := make([]BaiduMessage, len(req.Messages))
	for i, msg := range req.Messages {
		baiduMessages[i] = BaiduMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// 构建百度请求
	baiduReq := BaiduChatRequest{
		Messages:    baiduMessages,
		Temperature: req.Temperature,
		Stream:      false,
	}

	body, err := json.Marshal(baiduReq)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	// 构建请求 URL (包含 access token 和模型)
	reqURL := fmt.Sprintf("%s/%s?access_token=%s", c.BaseURL, c.Model, accessToken)

	httpReq, err := http.NewRequest("POST", reqURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	log.Printf("[BaiduClient] Response status: %d, body: %s", resp.StatusCode, string(respBody)[:min(500, len(respBody))])

	var baiduResp BaiduChatResponse
	if err := json.Unmarshal(respBody, &baiduResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w, body: %s", err, string(respBody))
	}

	if baiduResp.Error != nil {
		return nil, fmt.Errorf("百度 API 错误: [%d] %s", baiduResp.Error.Code, baiduResp.Error.Message)
	}

	// 转换为标准响应格式
	return &ChatResponse{
		ID:      baiduResp.ID,
		Object:  baiduResp.Object,
		Created: baiduResp.Created,
		Model:   c.Model,
		Choices: []struct {
			Index        int     `json:"index"`
			Message      Message `json:"message"`
			FinishReason string  `json:"finish_reason"`
		}{
			{
				Index: 0,
				Message: Message{
					Role:    "assistant",
					Content: baiduResp.Result,
				},
				FinishReason: "stop",
			},
		},
		Usage: struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		}{
			PromptTokens:     baiduResp.Usage.PromptTokens,
			CompletionTokens: baiduResp.Usage.CompletionTokens,
			TotalTokens:      baiduResp.Usage.TotalTokens,
		},
	}, nil
}

// Chat 发送简单聊天消息
func (c *BaiduClient) Chat(systemPrompt, userMessage string) (string, error) {
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

// ChatWithTools 发送带工具调用的聊天请求
// 百度千帆的 Function Calling 格式与 OpenAI 略有不同
func (c *BaiduClient) ChatWithTools(systemPrompt string, messages []Message, tools []Tool) (*ChatResponse, error) {
	// 构建完整消息列表
	allMessages := make([]Message, 0, len(messages)+1)
	allMessages = append(allMessages, Message{Role: "system", Content: systemPrompt})
	allMessages = append(allMessages, messages...)

	req := &ChatRequest{
		Messages:    allMessages,
		Tools:       tools,
		Temperature: 0.7,
	}

	return c.ChatCompletion(req)
}

// ChatCompletionStream 发送流式聊天请求
func (c *BaiduClient) ChatCompletionStream(req *ChatRequest, onChunk func(chunk string)) (*ChatResponse, error) {
	// 获取 access token
	accessToken, err := c.getAccessToken()
	if err != nil {
		return nil, err
	}

	// 转换消息格式
	baiduMessages := make([]BaiduMessage, len(req.Messages))
	for i, msg := range req.Messages {
		baiduMessages[i] = BaiduMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	// 构建百度请求
	baiduReq := BaiduChatRequest{
		Messages:    baiduMessages,
		Temperature: req.Temperature,
		Stream:      true,
	}

	body, err := json.Marshal(baiduReq)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	// 构建请求 URL
	reqURL := fmt.Sprintf("%s/%s?access_token=%s", c.BaseURL, c.Model, accessToken)

	httpReq, err := http.NewRequest("POST", reqURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API 错误: %s - %s", resp.Status, string(respBody))
	}

	// 解析 SSE 流
	var fullContent string
	var chatResp ChatResponse

	buf := make([]byte, 4096)
	for {
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			break
		}

		data := string(buf[:n])
		lines := strings.Split(data, "\n")

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || !strings.HasPrefix(line, "data: ") {
				continue
			}

			jsonData := strings.TrimPrefix(line, "data: ")
			if jsonData == "[DONE]" {
				continue
			}

			// 解析百度流式响应
			var streamResp struct {
				Result string `json:"result"`
				IsEnd  bool   `json:"is_end"`
			}
			if err := json.Unmarshal([]byte(jsonData), &streamResp); err != nil {
				continue
			}

			if streamResp.Result != "" {
				fullContent += streamResp.Result
				if onChunk != nil {
					onChunk(streamResp.Result)
				}
			}
		}

		if err == io.EOF {
			break
		}
	}

	// 构建最终响应
	chatResp = ChatResponse{
		Model: c.Model,
		Choices: []struct {
			Index        int     `json:"index"`
			Message      Message `json:"message"`
			FinishReason string  `json:"finish_reason"`
		}{
			{
				Index: 0,
				Message: Message{
					Role:    "assistant",
					Content: fullContent,
				},
				FinishReason: "stop",
			},
		},
	}

	return &chatResp, nil
}

// ChatWithToolsStream 发送带工具调用的流式聊天请求
func (c *BaiduClient) ChatWithToolsStream(systemPrompt string, messages []Message, tools []Tool, onChunk func(chunk string)) (*ChatResponse, error) {
	// 构建完整消息列表
	allMessages := make([]Message, 0, len(messages)+1)
	allMessages = append(allMessages, Message{Role: "system", Content: systemPrompt})
	allMessages = append(allMessages, messages...)

	req := &ChatRequest{
		Messages:    allMessages,
		Tools:       tools,
		Temperature: 0.7,
	}

	return c.ChatCompletionStream(req, onChunk)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// BuildRequestURL 构建带参数的请求 URL
func (c *BaiduClient) BuildRequestURL(endpoint string) string {
	return fmt.Sprintf("%s%s?access_token=%s", c.BaseURL, endpoint, url.QueryEscape(c.accessToken))
}
