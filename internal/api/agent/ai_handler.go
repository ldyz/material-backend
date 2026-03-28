package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	openai "github.com/yourorg/material-backend/backend/pkg/openai"
	"gorm.io/gorm"
)

// AIHandler handles AI-powered agent operations
type AIHandler struct {
	db            *gorm.DB
	aiService     *openai.AIService
	client        *openai.Client      // DeepSeek client for chat (OpenAI 兼容)
	baiduClient   *openai.BaiduClient // 百度千帆客户端 (可选)
	whisperClient *openai.Client      // OpenAI client for voice transcription (optional)
	asrServiceURL string              // Local ASR service URL
	service       *Service
	baseURL       string // API基础URL，用于内部调用
	useBaidu      bool   // 是否使用百度千帆

	// 支持多模型切换
	providers     map[string]*openai.Client // 所有可用的提供者
	currentProvider string                   // 当前使用的提供者
	providerConfigs map[string]ProviderConfig // 提供者配置（用于API返回）
}

// ProviderConfig 提供者配置
type ProviderConfig struct {
	Name    string `json:"name"`
	Model   string `json:"model"`
	BaseURL string `json:"base_url"`
}

// NewAIHandler creates a new AI handler
// deepSeekKey: DeepSeek API key for chat
// deepSeekModel: DeepSeek model name (default: deepseek-chat)
// deepSeekBaseURL: DeepSeek API base URL (default: https://api.deepseek.com/v1)
// openAIKey: OpenAI API key for voice transcription (optional)
// asrServiceURL: Local ASR service URL (optional, e.g., http://localhost:8089)
func NewAIHandler(db *gorm.DB, deepSeekKey, deepSeekModel, deepSeekBaseURL, openAIKey, asrServiceURL string) *AIHandler {
	var aiService *openai.AIService
	var client *openai.Client
	var whisperClient *openai.Client

	// Initialize AI client for chat (supports DeepSeek, Baidu Qianfan, etc.)
	if deepSeekKey != "" {
		aiService = openai.NewAIServiceWithBaseURL(deepSeekKey, deepSeekModel, deepSeekBaseURL)
		client = openai.NewClientWithBaseURL(deepSeekKey, deepSeekModel, deepSeekBaseURL)
		log.Printf("[AIHandler] AI客户端初始化成功, BaseURL: %s, Model: %s", deepSeekBaseURL, deepSeekModel)
	}

	// Initialize OpenAI client for voice transcription (optional)
	if openAIKey != "" {
		whisperClient = openai.NewClient(openAIKey, "whisper-1")
	}

	// 检测是否是百度千帆
	isBaidu := strings.Contains(deepSeekBaseURL, "qianfan.baidubce.com")

	return &AIHandler{
		db:            db,
		aiService:     aiService,
		client:        client,
		whisperClient: whisperClient,
		asrServiceURL: asrServiceURL,
		service:       NewService(db),
		baseURL:       "http://127.0.0.1:8088", // 默认内部调用地址（后端服务端口）
		useBaidu:      isBaidu,
	}
}

// NewAIHandlerWithBaidu creates a new AI handler with Baidu Qianfan support
func NewAIHandlerWithBaidu(db *gorm.DB, baiduAPIKey, baiduSecretKey, baiduModel, baiduBaseURL, openAIKey, asrServiceURL string) *AIHandler {
	var baiduClient *openai.BaiduClient
	var whisperClient *openai.Client

	// Initialize Baidu Qianfan client for chat
	if baiduAPIKey != "" && baiduSecretKey != "" {
		baiduClient = openai.NewBaiduClient(baiduAPIKey, baiduSecretKey, baiduModel, baiduBaseURL)
		log.Printf("[AIHandler] 百度千帆客户端初始化成功，模型: %s", baiduModel)
	}

	// Initialize OpenAI client for voice transcription (optional)
	if openAIKey != "" {
		whisperClient = openai.NewClient(openAIKey, "whisper-1")
	}

	return &AIHandler{
		db:            db,
		baiduClient:   baiduClient,
		whisperClient: whisperClient,
		asrServiceURL: asrServiceURL,
		service:       NewService(db),
		baseURL:       "http://127.0.0.1:8088",
		useBaidu:      true,
	}
}

// NewMultiModelAIHandler creates a new AI handler with multiple model support
func NewMultiModelAIHandler(db *gorm.DB, baiduAPIKey, baiduModel, baiduBaseURL, deepSeekKey, deepSeekModel, deepSeekBaseURL, openAIKey, asrServiceURL, defaultProvider string) *AIHandler {
	providers := make(map[string]*openai.Client)
	providerConfigs := make(map[string]ProviderConfig)
	var whisperClient *openai.Client

	// Initialize Baidu Qianfan client (Anthropic 兼容 API)
	if baiduAPIKey != "" && baiduBaseURL != "" {
		client := openai.NewClientWithBaseURL(baiduAPIKey, baiduModel, baiduBaseURL)
		providers["baidu"] = client
		providerConfigs["baidu"] = ProviderConfig{
			Name:    "百度千帆",
			Model:   baiduModel,
			BaseURL: baiduBaseURL,
		}
		log.Printf("[AIHandler] 百度千帆客户端初始化成功: %s, 模型: %s", baiduBaseURL, baiduModel)
	}

	// Initialize DeepSeek client
	if deepSeekKey != "" {
		client := openai.NewClientWithBaseURL(deepSeekKey, deepSeekModel, deepSeekBaseURL)
		providers["deepseek"] = client
		providerConfigs["deepseek"] = ProviderConfig{
			Name:    "DeepSeek",
			Model:   deepSeekModel,
			BaseURL: deepSeekBaseURL,
		}
		log.Printf("[AIHandler] DeepSeek客户端初始化成功: %s, 模型: %s", deepSeekBaseURL, deepSeekModel)
	}

	// Initialize OpenAI client for voice transcription (optional)
	if openAIKey != "" {
		whisperClient = openai.NewClient(openAIKey, "whisper-1")
	}

	// Determine default provider
	currentProvider := defaultProvider
	if currentProvider == "" {
		if _, ok := providers["baidu"]; ok {
			currentProvider = "baidu"
		} else if _, ok := providers["deepseek"]; ok {
			currentProvider = "deepseek"
		}
	}

	// Get current client
	var currentClient *openai.Client
	var useBaidu bool
	if currentProvider == "baidu" {
		currentClient = providers["baidu"]
		useBaidu = true
	} else {
		currentClient = providers["deepseek"]
		useBaidu = false
	}

	return &AIHandler{
		db:              db,
		client:          currentClient,
		whisperClient:   whisperClient,
		asrServiceURL:   asrServiceURL,
		service:         NewService(db),
		baseURL:         "http://127.0.0.1:8088",
		useBaidu:        useBaidu,
		providers:       providers,
		currentProvider: currentProvider,
		providerConfigs: providerConfigs,
	}
}

// SetUseBaidu 设置是否使用百度千帆
func (h *AIHandler) SetUseBaidu(useBaidu bool) {
	h.useBaidu = useBaidu
}

// SwitchProvider 切换模型提供者
func (h *AIHandler) SwitchProvider(provider string) error {
	client, ok := h.providers[provider]
	if !ok {
		return fmt.Errorf("未知的提供者: %s", provider)
	}
	h.client = client
	h.currentProvider = provider
	h.useBaidu = (provider == "baidu")
	log.Printf("[AIHandler] 切换到模型提供者: %s", provider)
	return nil
}

// GetCurrentProvider 获取当前提供者
func (h *AIHandler) GetCurrentProvider() string {
	return h.currentProvider
}

// GetProviders 获取所有可用提供者
func (h *AIHandler) GetProviders() map[string]ProviderConfig {
	return h.providerConfigs
}

// SetBaseURL 设置API基础URL
func (h *AIHandler) SetBaseURL(url string) {
	h.baseURL = url
}

// TranscribeAudio transcribes audio to text using local ASR service or Whisper API
func (h *AIHandler) TranscribeAudio(ctx context.Context, audioFile io.Reader, filename string) (string, error) {
	// 优先使用本地 ASR 服务
	if h.asrServiceURL != "" {
		return h.transcribeWithLocalASR(ctx, audioFile, filename)
	}

	// 回退到 OpenAI Whisper API
	if h.whisperClient != nil {
		return h.whisperClient.TranscribeAudio(audioFile, filename)
	}

	return "", fmt.Errorf("语音转文字功能未配置，请设置 ASR_SERVICE_URL 或 OPENAI_API_KEY 环境变量")
}

// transcribeWithLocalASR 使用本地 ASR 服务进行语音识别
func (h *AIHandler) transcribeWithLocalASR(ctx context.Context, audioFile io.Reader, filename string) (string, error) {
	log.Printf("[ASR] Starting local ASR transcription, filename: %s, service URL: %s", filename, h.asrServiceURL)

	// 读取音频数据
	audioData, err := io.ReadAll(audioFile)
	if err != nil {
		return "", fmt.Errorf("读取音频数据失败: %w", err)
	}
	log.Printf("[ASR] Audio data read, size: %d bytes", len(audioData))

	// 创建 multipart form
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// 添加音频文件
	part, err := writer.CreateFormFile("audio", filename)
	if err != nil {
		return "", fmt.Errorf("创建 form file 失败: %w", err)
	}
	if _, err := part.Write(audioData); err != nil {
		return "", fmt.Errorf("写入音频数据失败: %w", err)
	}
	writer.Close()

	// 发送请求到本地 ASR 服务
	req, err := http.NewRequestWithContext(ctx, "POST", h.asrServiceURL+"/transcribe", &body)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	log.Printf("[ASR] Sending request to ASR service: %s", h.asrServiceURL+"/transcribe")
	startTime := time.Now()

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		elapsed := time.Since(startTime)
		log.Printf("[ASR] Request failed after %v: %v", elapsed, err)
		// Check if context was canceled
		if ctx.Err() != nil {
			return "", fmt.Errorf("ASR 请求被取消 (耗时 %v): %w", elapsed, ctx.Err())
		}
		return "", fmt.Errorf("ASR 服务请求失败: %w", err)
	}
	defer resp.Body.Close()

	elapsed := time.Since(startTime)
	log.Printf("[ASR] Response received in %v, status: %d", elapsed, resp.StatusCode)

	// 解析响应
	var result struct {
		Success bool   `json:"success"`
		Text    string `json:"text"`
		Error   string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	if !result.Success {
		return "", fmt.Errorf("语音识别失败: %s", result.Error)
	}

	log.Printf("[ASR] Transcription successful: %s", result.Text)
	return result.Text, nil
}

// HandleAIChat handles AI chat requests with tool calling
func (h *AIHandler) HandleAIChat(ctx context.Context, req *AIChatRequest) (*AIChatResponse, error) {
	if h.aiService == nil {
		return nil, fmt.Errorf("AI service not configured")
	}

	// Create executor
	executor := &toolExecutor{
		db:        h.db,
		service:   h.service,
		userID:    req.UserID,
		apiExec:   NewAPIToolExecutor(h.db, uint(req.UserID), ""),
	}

	// Call AI with tools
	resp, err := h.aiService.ChatWithTools(ctx, &openai.AIChatRequest{
		Message:             req.Message,
		ConversationHistory: req.ConversationHistory,
		Context:             req.Context,
		UserID:              req.UserID,
	}, executor)

	if err != nil {
		return nil, err
	}

	// Convert response
	return &AIChatResponse{
		Message:      resp.Message,
		ToolCalls:    convertToolCalls(resp.ToolCalls),
		NeedsAction:  resp.NeedsAction,
		ActionResult: resp.ActionResult,
		Conversation: resp.Conversation,
	}, nil
}

// HandleAIChatStream handles AI chat requests with streaming support
func (h *AIHandler) HandleAIChatStream(ctx context.Context, req *AIChatRequest, onChunk func(chunk string)) (*AIChatResponse, error) {
	// 检查是否使用百度千帆
	if h.useBaidu && h.baiduClient != nil {
		return h.handleBaiduChatStream(ctx, req, onChunk)
	}

	// 使用 DeepSeek/OpenAI 兼容客户端
	if h.client == nil {
		return nil, fmt.Errorf("AI service not configured")
	}

	// 设置总体超时为 60 秒
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// Build system prompt
	systemPrompt := h.buildSystemPrompt(req.Context)

	// Build messages
	messages := req.ConversationHistory
	if len(messages) == 0 {
		messages = []openai.Message{}
	}
	messages = append(messages, openai.Message{Role: "user", Content: req.Message})

	// Get tools
	tools := openai.GetMaterialManagementTools()

	// Create executor - 使用增强的API执行器
	executor := &toolExecutor{
		db:      h.db,
		service: h.service,
		userID:  req.UserID,
		apiExec: NewAPIToolExecutor(h.db, uint(req.UserID), ""),
	}

	// Track all tool calls
	allToolCallInfos := []ToolCallInfo{}

	// Use a loop to handle multiple rounds of tool calls
	maxIterations := 5 // 减少到5次，避免长时间等待
	for iteration := 0; iteration < maxIterations; iteration++ {
		// 检查 context 是否已取消
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("请求超时或已取消")
		default:
		}

		log.Printf("[AIHandler] Starting iteration %d", iteration+1)

		// Call AI with streaming
		resp, err := h.client.ChatWithToolsStream(systemPrompt, messages, tools, onChunk)
		if err != nil {
			return nil, fmt.Errorf("streaming chat failed: %w", err)
		}

		if len(resp.Choices) == 0 {
			return nil, fmt.Errorf("no response from AI")
		}

		choice := resp.Choices[0]

		// 添加调试日志
		log.Printf("[AIHandler] FinishReason: %s, ToolCalls count: %d, Content length: %d",
			choice.FinishReason, len(choice.Message.ToolCalls), len(choice.Message.Content))

		for i, tc := range choice.Message.ToolCalls {
			log.Printf("[AIHandler] ToolCall[%d]: ID=%s, Name=%s, Args=%s", i, tc.ID, tc.Function.Name, tc.Function.Arguments)
		}

		// Check if AI wants to call tools
		if len(choice.Message.ToolCalls) > 0 {
			// Add assistant message with tool calls to conversation
			// 确保消息有 content 字段（DeepSeek API 要求）
			assistantMsg := choice.Message
			if assistantMsg.Content == "" && len(assistantMsg.ToolCalls) > 0 {
				assistantMsg.Content = "" // 显式设置空字符串，避免缺失字段
			}
			messages = append(messages, assistantMsg)

			// Execute each tool call
			for _, toolCall := range choice.Message.ToolCalls {
				// Parse arguments
				var args map[string]interface{}
				if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
					log.Printf("Failed to parse tool arguments: %v", err)
					messages = append(messages, openai.FormatToolResult(toolCall.ID, "Invalid arguments", true))
					continue
				}

				// Execute tool - 优先使用API执行器
				result, err := executor.ExecuteToolCall(ctx, toolCall.Function.Name, args, req.UserID)
				if err != nil {
					resultJSON, _ := json.Marshal(map[string]interface{}{"error": err.Error()})
					messages = append(messages, openai.FormatToolResult(toolCall.ID, string(resultJSON), true))
				} else {
					resultJSON, _ := json.Marshal(result)
					log.Printf("[AIHandler] Tool result for %s: %s", toolCall.Function.Name, string(resultJSON))
					messages = append(messages, openai.FormatToolResult(toolCall.ID, string(resultJSON), false))
				}

				allToolCallInfos = append(allToolCallInfos, ToolCallInfo{
					ID:        toolCall.ID,
					Name:      toolCall.Function.Name,
					Arguments: args,
					Result:    result,
				})
			}

			// Continue the loop to get the next response
			continue
		}

		// No tool calls, this is the final response
		messages = append(messages, choice.Message)
		return &AIChatResponse{
			Message:      choice.Message.Content,
			ToolCalls:    allToolCallInfos,
			NeedsAction:  false,
			Conversation: messages,
		}, nil
	}

	// If we reach here, we've exceeded max iterations
	return nil, fmt.Errorf("exceeded maximum tool call iterations")
}

// handleBaiduChatStream 使用百度千帆处理聊天请求
func (h *AIHandler) handleBaiduChatStream(ctx context.Context, req *AIChatRequest, onChunk func(chunk string)) (*AIChatResponse, error) {
	if h.baiduClient == nil {
		return nil, fmt.Errorf("百度千帆服务未配置")
	}

	// 设置总体超时为 90 秒
	ctx, cancel := context.WithTimeout(ctx, 90*time.Second)
	defer cancel()

	// Build system prompt
	systemPrompt := h.buildSystemPrompt(req.Context)

	// Build messages
	messages := req.ConversationHistory
	if len(messages) == 0 {
		messages = []openai.Message{}
	}
	messages = append(messages, openai.Message{Role: "user", Content: req.Message})

	log.Printf("[AIHandler] Using Baidu Qianfan for chat, model: coding-plan")

	// Call Baidu API with streaming
	resp, err := h.baiduClient.ChatCompletionStream(&openai.ChatRequest{
		Messages:    append([]openai.Message{{Role: "system", Content: systemPrompt}}, messages...),
		Temperature: 0.7,
	}, onChunk)

	if err != nil {
		return nil, fmt.Errorf("百度千帆请求失败: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("百度千帆返回空响应")
	}

	choice := resp.Choices[0]
	messages = append(messages, choice.Message)

	log.Printf("[AIHandler] Baidu response completed, content length: %d", len(choice.Message.Content))

	return &AIChatResponse{
		Message:      choice.Message.Content,
		ToolCalls:    nil, // 百度千帆暂时不支持 Function Calling
		NeedsAction:  false,
		Conversation: messages,
	}, nil
}

// buildSystemPrompt builds the system prompt based on context
func (h *AIHandler) buildSystemPrompt(context map[string]interface{}) string {
	var sb strings.Builder

	// 获取当前日期时间
	now := time.Now()
	today := now.Format("2006-01-02")
	tomorrow := now.AddDate(0, 0, 1).Format("2006-01-02")
	yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")
	weekday := map[time.Weekday]string{
		time.Sunday: "周日", time.Monday: "周一", time.Tuesday: "周二",
		time.Wednesday: "周三", time.Thursday: "周四", time.Friday: "周五", time.Saturday: "周六",
	}

	sb.WriteString(`你是一个材料和施工管理系统的AI助手，主要通过语音交互为用户提供服务。

## 重要规则
1. 必须使用简体中文回复，不要使用繁体中文
2. 回复要简洁明了，适合语音播报
3. 不要编造不存在的功能或数据

## 当前日期时间
- 今天: ` + today + ` (` + weekday[now.Weekday()] + `)
- 明天: ` + tomorrow + `
- 昨天: ` + yesterday + `

## 系统能力概览

本系统是一个材料和施工管理平台，主要功能包括：

### 1. 施工预约管理
- 创建、查询、修改、取消施工预约
- 分配作业人员、审批流程
- 日历视图、人员排班
- 加急任务处理

**预约状态说明**：
- draft: 草稿 - 新建但未提交
- pending: 待审批 - 已提交等待审批
- scheduled: 已排期 - 审批通过，已分配人员
- in_progress: 进行中 - 任务已开始
- completed: 已完成 - 任务已完成
- cancelled: 已取消 - 任务已取消
- rejected: 已驳回 - 审批未通过

**时间段说明**：
- morning: 上午 (08:00-12:00)
- noon: 中午 (12:00-14:00)
- afternoon: 下午 (14:00-18:00)
- fullday: 全天

**作业人员日历锁定规则**（重要）：
- 创建预约时指定了作业人员，**立即锁定**该人员的时间段，其他人无法再预约
- 更换作业人员时，释放旧人员的锁定，锁定新人员
- 删除或取消预约时，释放作业人员的锁定
- 完成预约时，释放作业人员的锁定
- 撤回工作流（回到草稿状态）时，**保持锁定**不释放
- 如果提示"作业人员在指定时间段不可用"，说明该人员已被其他预约占用

### 2. 库存管理
- 库存查询、入库、出库
- 库存预警、安全库存管理
- 仓库管理

### 3. 物资计划
- 物资需求计划创建与审批
- 物资主数据管理
- 领用单/出库单管理

### 4. 项目管理
- 项目信息查询、进度跟踪

### 5. 考勤管理
- 查询打卡记录（按人员、日期、类型筛选）
- 考勤统计数据（月度汇总、加班统计）
- 打卡状态查询（待确认、已确认）

**打卡类型说明**：
- morning: 上午打卡
- afternoon: 下午打卡
- noon_overtime: 中午加班
- night_overtime: 晚上加班

### 6. 施工日志
- 查询施工日志记录
- 按项目、日期筛选日志

## 工具使用指南

### 查询任务/预约
用户说："明天的任务"、"今日安排"、"施工计划"
→ 使用 query_appointments 工具，设置 date 参数

用户说："本周有什么安排"
→ 使用 query_appointments 工具，设置 start_date 和 end_date 参数

### 创建预约
用户说："创建预约"、"安排施工"
→ 先用 query_projects 获取项目信息
→ 再用 query_workers 获取可用人员
→ 最后用 create_appointment 创建预约

### 库存相关
用户说："库存有多少"、"库存够吗"
→ 使用 query_stock 工具

用户说："低库存"、"库存预警"
→ 使用 query_stock_alerts 工具

### 审批相关
用户说："待审批"、"需要我审批"
→ 使用 query_pending_approvals 工具

用户说："审批通过"、"同意审批"
→ 使用 approve_appointment 工具，action 设为 "approve"

用户说："驳回"、"拒绝"
→ 使用 approve_appointment 工具，action 设为 "reject"

### 考勤相关
用户说："打卡记录"、"考勤记录"
→ 使用 query_attendance 工具

用户说："今天的考勤"、"某人出勤情况"
→ 使用 query_attendance 工具，设置日期或用户参数

用户说："考勤统计"、"出勤率"
→ 使用 query_attendance_stats 工具

### 施工日志相关
用户说："施工日志"、"施工记录"
→ 使用 query_construction_logs 工具

用户说："今天的日志"、"某项目的日志"
→ 使用 query_construction_logs 工具，设置日期或项目参数

用户说："日志详情"、"查看这条日志"
→ 使用 query_construction_log_detail 工具

### 工作流相关
用户说："待办任务"、"我的待办"
→ 使用 query_pending_workflow_tasks 工具

用户说："审批通过这个任务"、"同意这个工作流"
→ 使用 approve_workflow_task 工具

用户说："驳回这个任务"、"拒绝这个工作流"
→ 使用 reject_workflow_task 工具，需要提供驳回原因

### 通知相关
用户说："通知列表"、"消息列表"、"我的消息"
→ 使用 query_notifications 工具

用户说："未读消息数量"、"有多少未读"
→ 使用 query_unread_notification_count 工具

## 回复原则（语音交互优化）

1. **简洁明了**：回复要简短直接，适合语音播报，避免冗长
2. **关键信息优先**：把最重要的信息放在开头
3. **自然对话**：使用口语化表达，像真人对话一样自然
4. **确认机制**：执行重要操作前，简要确认用户意图
5. **错误友好**：如果无法完成操作，给出明确的建议

## 示例对话

用户："帮我查一下明天的任务"
助手："明天（` + tomorrow + `）有3个施工预约。第一个是水电改造，上午在A栋进行，张三负责。需要我详细说明其他预约吗？"

用户："给李四安排后天的施工"
助手："好的，后天是" + now.AddDate(0, 0, 2).Format("2006-01-02") + "。请告诉我：1. 具体时间？2. 施工内容？3. 哪个项目？"

用户："库存够吗"
助手："目前有2种物资库存偏低：水泥剩50袋，安全库存100袋；钢筋剩8吨，安全库存15吨。需要我帮您创建采购计划吗？"

用户："审批通过第一个"
助手："好的，已将编号APT001的施工预约审批通过。这是XX项目的安装任务，明天上午进行。"

用户："今天的考勤情况"
助手："今天有5人打卡。上午4人打卡，下午3人打卡。有2人中午加班共3小时。需要我详细说明每个人的打卡时间吗？"

用户："查看某项目的施工日志"
助手："请告诉我项目名称或编号，以及要查看的日期范围。"

用户："帮张三安排明天上午的施工"
助手："抱歉，张三在明天上午的时间段已被其他预约占用。请选择其他时间或更换作业人员。需要我查看张三的空闲时间吗？"

`)

	// Add user context
	if context != nil {
		if userID, ok := context["user_id"].(int); ok {
			sb.WriteString(fmt.Sprintf("\n当前用户ID: %d\n", userID))
		}
		if isAdmin, ok := context["is_admin"].(bool); ok && isAdmin {
			sb.WriteString("当前用户角色: 管理员\n")
		}
	}

	return sb.String()
}

// toolExecutor implements ToolExecutor interface
type toolExecutor struct {
	db      *gorm.DB
	service *Service
	userID  int
	apiExec *APIToolExecutor // API工具执行器
}

// ExecuteToolCall executes a tool call
func (e *toolExecutor) ExecuteToolCall(ctx context.Context, name string, arguments map[string]interface{}, userID int) (interface{}, error) {
	log.Printf("[AI Agent] Executing tool: %s with arguments: %+v", name, arguments)

	// 优先使用API执行器
	if e.apiExec != nil {
		// 尝试使用API执行器
		if _, ok := ToolDefinitions[name]; ok {
			result, err := e.apiExec.ExecuteToolCall(ctx, name, arguments)
			if err != nil {
				log.Printf("[AI Agent] API executor failed for %s: %v, falling back to DB", name, err)
			} else {
				return result, nil
			}
		}
	}

	// 回退到数据库直接操作
	switch name {
	case "query_materials":
		return e.queryMaterials(arguments)
	case "query_stock":
		return e.queryStock(arguments)
	case "query_stock_alerts":
		return e.queryStockAlerts(arguments)
	case "query_material_plans":
		return e.queryMaterialPlans(arguments)
	case "create_material_plan":
		return e.createMaterialPlan(arguments, userID)
	case "update_stock":
		return e.updateStock(arguments, userID)
	case "query_projects":
		return e.queryProjects(arguments)
	case "approve_workflow":
		return e.approveWorkflow(arguments, userID)
	case "generate_report":
		return e.generateReport(arguments)
	case "analyze_data":
		return e.analyzeData(arguments)
	case "query_appointments":
		return e.queryAppointments(arguments)
	case "create_appointment":
		return e.createAppointment(arguments, userID)
	case "query_worker_calendar":
		return e.queryWorkerCalendar(arguments)
	case "query_pending_approvals":
		return e.queryPendingApprovals(arguments, userID)
	case "query_workers":
		return e.queryWorkers()
	case "approve_appointment":
		return e.approveAppointment(arguments, userID)
	case "cancel_appointment":
		return e.cancelAppointment(arguments, userID)
	case "query_attendance":
		return e.queryAttendance(arguments)
	case "query_attendance_stats":
		return e.queryAttendanceStats(arguments)
	case "query_construction_logs":
		return e.queryConstructionLogs(arguments)
	default:
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
}

func (e *toolExecutor) queryMaterials(args map[string]interface{}) (interface{}, error) {
	search, _ := args["search"].(string)
	category, _ := args["category"].(string)
	limit := 20
	if l, ok := args["page_size"].(float64); ok {
		limit = int(l)
	}

	query := e.db.Table("material_master").
		Select("id, name, specification, unit, category, description").
		Limit(limit)

	if search != "" {
		query = query.Where("name ILIKE ? OR specification ILIKE ?",
			"%"+search+"%", "%"+search+"%")
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}

	var results []map[string]interface{}
	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success": true,
		"data":    results,
		"count":   len(results),
	}, nil
}

func (e *toolExecutor) queryStock(args map[string]interface{}) (interface{}, error) {
	materialID, _ := args["material_id"].(float64)
	lowStockOnly, _ := args["low_stock"].(bool)
	projectID, _ := args["project_id"].(float64)

	query := e.db.Table("stocks s").
		Select("s.id, s.material_id, m.name as material_name, m.specification, s.quantity, s.safety_stock, s.warehouse").
		Joins("LEFT JOIN material_master m ON s.material_id = m.id")

	if materialID > 0 {
		query = query.Where("s.material_id = ?", int(materialID))
	}
	if lowStockOnly {
		query = query.Where("s.quantity < s.safety_stock")
	}
	if projectID > 0 {
		query = query.Where("s.project_id = ?", int(projectID))
	}

	var results []map[string]interface{}
	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}

	// Calculate low stock count
	lowStockCount := 0
	for _, r := range results {
		qty, _ := r["quantity"].(float64)
		safety, _ := r["safety_stock"].(float64)
		if qty < safety {
			lowStockCount++
		}
	}

	return map[string]interface{}{
		"success":         true,
		"data":            results,
		"count":           len(results),
		"low_stock_count": lowStockCount,
	}, nil
}

func (e *toolExecutor) queryStockAlerts(args map[string]interface{}) (interface{}, error) {
	query := e.db.Table("stocks s").
		Select("s.id, s.material_id, m.name as material_name, m.specification, s.quantity, s.safety_stock, s.warehouse").
		Joins("LEFT JOIN material_master m ON s.material_id = m.id").
		Where("s.quantity < s.safety_stock")

	var results []map[string]interface{}
	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success":  true,
		"data":     results,
		"count":    len(results),
		"message":  fmt.Sprintf("发现 %d 种物资库存偏低", len(results)),
	}, nil
}

func (e *toolExecutor) queryMaterialPlans(args map[string]interface{}) (interface{}, error) {
	projectID, _ := args["project_id"].(float64)
	status, _ := args["status"].(string)
	limit := 20
	if l, ok := args["page_size"].(float64); ok {
		limit = int(l)
	}

	query := e.db.Table("material_plans").
		Select("id, plan_no, project_id, status, total_budget, created_at").
		Limit(limit)

	if projectID > 0 {
		query = query.Where("project_id = ?", int(projectID))
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var results []map[string]interface{}
	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success": true,
		"data":    results,
		"count":   len(results),
	}, nil
}

func (e *toolExecutor) createMaterialPlan(args map[string]interface{}, userID int) (interface{}, error) {
	projectID, _ := args["project_id"].(float64)
	items, _ := args["items"].([]interface{})
	remark, _ := args["remark"].(string)

	if projectID == 0 {
		return nil, fmt.Errorf("project_id is required")
	}

	// Use existing service
	op := &AgentOperation{
		Operation: OpCreatePlan,
		Resource:  "material_plan",
		Parameters: map[string]interface{}{
			"project_id": int(projectID),
			"items":      items,
			"remark":     remark,
		},
		Context: map[string]interface{}{
			"user_id": userID,
		},
		Reasoning: "AI initiated material plan creation",
	}

	result, err := e.service.HandleOperation(op, userID, "ai-assistant")
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success":  true,
		"plan_id":  result.Result["plan_id"],
		"plan_no":  result.Result["plan_no"],
		"message":  "物资计划创建成功",
	}, nil
}

func (e *toolExecutor) updateStock(args map[string]interface{}, userID int) (interface{}, error) {
	stockID, _ := args["stock_id"].(float64)
	quantity, _ := args["quantity"].(float64)
	operation, _ := args["operation"].(string)
	remark, _ := args["remark"].(string)

	if operation == "" {
		operation = "set"
	}

	if stockID == 0 {
		return nil, fmt.Errorf("stock_id is required")
	}

	// Get current stock
	var currentStock map[string]interface{}
	if err := e.db.Table("stocks").Where("id = ?", int(stockID)).First(&currentStock).Error; err != nil {
		return nil, fmt.Errorf("stock not found")
	}

	currentQty, _ := currentStock["quantity"].(float64)
	newQty := quantity

	switch operation {
	case "add":
		newQty = currentQty + quantity
	case "subtract":
		newQty = currentQty - quantity
	}

	// Update stock
	if err := e.db.Table("stocks").Where("id = ?", int(stockID)).Update("quantity", newQty).Error; err != nil {
		return nil, err
	}

	// Create log
	e.db.Table("stock_op_logs").Create(map[string]interface{}{
		"stock_id":  int(stockID),
		"quantity":  newQty,
		"operation": operation,
		"remark":    remark,
		"user_id":   userID,
	})

	return map[string]interface{}{
		"success":      true,
		"stock_id":     int(stockID),
		"old_quantity": currentQty,
		"new_quantity": newQty,
		"operation":    operation,
		"message":      "库存更新成功",
	}, nil
}

func (e *toolExecutor) queryProjects(args map[string]interface{}) (interface{}, error) {
	search, _ := args["search"].(string)
	status, _ := args["status"].(string)
	limit := 20
	if l, ok := args["page_size"].(float64); ok {
		limit = int(l)
	}

	query := e.db.Table("projects").
		Select("id, name, code, status, manager, start_date, end_date").
		Limit(limit)

	if search != "" {
		query = query.Where("name ILIKE ? OR code ILIKE ?",
			"%"+search+"%", "%"+search+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var results []map[string]interface{}
	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success": true,
		"data":    results,
		"count":   len(results),
	}, nil
}

func (e *toolExecutor) approveWorkflow(args map[string]interface{}, userID int) (interface{}, error) {
	taskID, _ := args["task_id"].(float64)
	action, _ := args["action"].(string)
	remark, _ := args["remark"].(string)

	if taskID == 0 {
		return nil, fmt.Errorf("task_id is required")
	}
	if action == "" {
		action = "approve"
	}

	// Use existing service
	op := &AgentOperation{
		Operation: OpApproveWorkflow,
		Resource:  "workflow",
		Parameters: map[string]interface{}{
			"task_id": int64(taskID),
			"action":  action,
			"remark":  remark,
		},
		Context: map[string]interface{}{
			"user_id": userID,
		},
		Reasoning: "AI initiated workflow approval",
	}

	result, err := e.service.HandleOperation(op, userID, "ai-assistant")
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success":  true,
		"task_id":  int(taskID),
		"action":   action,
		"message":  "工作流操作成功",
		"result":   result.Result,
	}, nil
}

func (e *toolExecutor) generateReport(args map[string]interface{}) (interface{}, error) {
	reportType, _ := args["report_type"].(string)
	projectID, _ := args["project_id"].(float64)

	if reportType == "" {
		return nil, fmt.Errorf("report_type is required")
	}

	op := &AgentOperation{
		Operation: OpGenerateReport,
		Resource:  "report",
		Parameters: map[string]interface{}{
			"report_type": reportType,
			"project_id":  int(projectID),
		},
	}

	result, err := e.service.HandleOperation(op, e.userID, "ai-assistant")
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success":     true,
		"report_type": reportType,
		"data":        result.Result,
		"message":     "报告生成成功",
	}, nil
}

func (e *toolExecutor) analyzeData(args map[string]interface{}) (interface{}, error) {
	analysisType, _ := args["analysis_type"].(string)
	question, _ := args["question"].(string)

	if analysisType == "" {
		return nil, fmt.Errorf("analysis_type is required")
	}

	op := &AgentOperation{
		Operation: OpAnalyze,
		Resource:  analysisType,
		Parameters: map[string]interface{}{
			"question": question,
		},
	}

	result, err := e.service.HandleOperation(op, e.userID, "ai-assistant")
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success":       true,
		"analysis_type": analysisType,
		"data":          result.Result,
		"message":       "数据分析完成",
	}, nil
}

func (e *toolExecutor) queryAppointments(args map[string]interface{}) (interface{}, error) {
	date, _ := args["date"].(string)
	projectID, _ := args["project_id"].(float64)
	status, _ := args["status"].(string)
	workerID, _ := args["worker_id"].(float64)
	startDate, _ := args["start_date"].(string)
	endDate, _ := args["end_date"].(string)
	limit := 20
	if l, ok := args["page_size"].(float64); ok {
		limit = int(l)
	}

	// 处理特殊日期值
	if date == "today" {
		date = time.Now().Format("2006-01-02")
	} else if date == "tomorrow" {
		date = time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	}

	query := e.db.Table("construction_appointments").
		Select(`id, appointment_no, work_date, time_slot, status, work_type,
			work_content, work_location, applicant_name, assigned_worker_names`).
		Limit(limit)

	if date != "" {
		query = query.Where("work_date = ?", date)
	}
	if projectID > 0 {
		query = query.Where("project_id = ?", int(projectID))
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if workerID > 0 {
		query = query.Where("assigned_worker_id = ? OR assigned_worker_ids::jsonb @ ?::jsonb",
			int(workerID), fmt.Sprintf("[%d]", int(workerID)))
	}
	if startDate != "" {
		query = query.Where("work_date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("work_date <= ?", endDate)
	}

	var results []map[string]interface{}
	if err := query.Order("work_date, time_slot").Find(&results).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success": true,
		"data":    results,
		"count":   len(results),
		"date":    date,
	}, nil
}

func (e *toolExecutor) createAppointment(args map[string]interface{}, userID int) (interface{}, error) {
	projectID, _ := args["project_id"].(float64)
	workDate, _ := args["work_date"].(string)
	timeSlot, _ := args["time_slot"].(string)
	workContent, _ := args["work_content"].(string)
	workLocation, _ := args["work_location"].(string)
	workType, _ := args["work_type"].(string)
	contactPerson, _ := args["contact_person"].(string)
	contactPhone, _ := args["contact_phone"].(string)
	isUrgent, _ := args["is_urgent"].(bool)
	priority, _ := args["priority"].(float64)
	urgentReason, _ := args["urgent_reason"].(string)

	if projectID == 0 {
		return nil, fmt.Errorf("project_id is required")
	}
	if workDate == "" {
		return nil, fmt.Errorf("work_date is required")
	}
	if timeSlot == "" {
		return nil, fmt.Errorf("time_slot is required")
	}
	if workContent == "" {
		return nil, fmt.Errorf("work_content is required")
	}

	// 获取用户名
	var userName string
	e.db.Table("users").Where("id = ?", userID).Pluck("full_name", &userName)

	// 创建预约
	aptData := map[string]interface{}{
		"project_id":       int(projectID),
		"work_date":        workDate,
		"time_slot":        timeSlot,
		"work_content":     workContent,
		"work_location":    workLocation,
		"work_type":        workType,
		"contact_person":   contactPerson,
		"contact_phone":    contactPhone,
		"is_urgent":        isUrgent,
		"priority":         int(priority),
		"urgent_reason":    urgentReason,
		"status":           "draft",
		"applicant_id":     userID,
		"applicant_name":   userName,
	}

	if err := e.db.Table("construction_appointments").Create(&aptData).Error; err != nil {
		return nil, err
	}

	aptID := aptData["id"]
	aptNo := aptData["appointment_no"]

	return map[string]interface{}{
		"success":        true,
		"appointment_id": aptID,
		"appointment_no": aptNo,
		"message":        "预约创建成功",
	}, nil
}

func (e *toolExecutor) queryWorkerCalendar(args map[string]interface{}) (interface{}, error) {
	workerID, _ := args["worker_id"].(float64)
	startDate, _ := args["start_date"].(string)
	endDate, _ := args["end_date"].(string)

	if startDate == "" || endDate == "" {
		return nil, fmt.Errorf("start_date and end_date are required")
	}

	query := e.db.Table("construction_appointments").
		Select(`id, appointment_no, work_date, time_slot, status,
			work_type, work_content, work_location, assigned_worker_names`).
		Where("work_date >= ? AND work_date <= ?", startDate, endDate).
		Where("status IN ?", []string{"scheduled", "in_progress", "pending"})

	if workerID > 0 {
		query = query.Where("assigned_worker_id = ? OR assigned_worker_ids::jsonb @ ?::jsonb",
			int(workerID), fmt.Sprintf("[%d]", int(workerID)))
	}

	var results []map[string]interface{}
	if err := query.Order("work_date, time_slot").Find(&results).Error; err != nil {
		return nil, err
	}

	// Group by date
	calendar := make(map[string][]map[string]interface{})
	for _, apt := range results {
		date, ok := apt["work_date"].(time.Time)
		if !ok {
			continue
		}
		dateStr := date.Format("2006-01-02")
		calendar[dateStr] = append(calendar[dateStr], apt)
	}

	return map[string]interface{}{
		"success":  true,
		"calendar": calendar,
		"count":    len(results),
	}, nil
}

func (e *toolExecutor) queryPendingApprovals(args map[string]interface{}, userID int) (interface{}, error) {
	limit := 20
	if l, ok := args["page_size"].(float64); ok {
		limit = int(l)
	}

	query := e.db.Table("construction_appointments").
		Select(`id, appointment_no, work_date, time_slot, status, work_type,
			work_content, work_location, applicant_name, priority, is_urgent`).
		Where("status = ?", "pending").
		Limit(limit).
		Order("is_urgent DESC, priority DESC, created_at ASC")

	var results []map[string]interface{}
	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success": true,
		"data":    results,
		"count":   len(results),
	}, nil
}

func (e *toolExecutor) queryWorkers() (interface{}, error) {
	query := e.db.Table("users").
		Select("id, full_name, username, email").
		Where("is_active = ? AND (role = ? OR role = ? OR role = ? OR role = ?)",
			true,
			"worker",
			"作业人员",
			"team_leader",
			"班组长",
		)

	var results []map[string]interface{}
	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success": true,
		"data":    results,
		"count":   len(results),
	}, nil
}

func (e *toolExecutor) approveAppointment(args map[string]interface{}, userID int) (interface{}, error) {
	id, _ := args["id"].(float64)
	action, _ := args["action"].(string)
	comment, _ := args["comment"].(string)

	if id == 0 {
		return nil, fmt.Errorf("id is required")
	}
	if action == "" {
		action = "approve"
	}

	// 获取预约单
	var appointment map[string]interface{}
	if err := e.db.Table("construction_appointments").Where("id = ?", int(id)).First(&appointment).Error; err != nil {
		return nil, fmt.Errorf("预约单不存在")
	}

	// 更新状态
	newStatus := "scheduled"
	if action == "reject" {
		newStatus = "rejected"
	}

	if err := e.db.Table("construction_appointments").Where("id = ?", int(id)).Update("status", newStatus).Error; err != nil {
		return nil, err
	}

	// 记录审批历史
	var userName string
	e.db.Table("users").Where("id = ?", userID).Pluck("full_name", &userName)

	e.db.Table("appointment_approval_history").Create(map[string]interface{}{
		"appointment_id": int(id),
		"user_id":        userID,
		"user_name":      userName,
		"action":         action,
		"comment":        comment,
	})

	return map[string]interface{}{
		"success":        true,
		"appointment_id": int(id),
		"action":         action,
		"new_status":     newStatus,
		"message":        fmt.Sprintf("预约单已%s", map[string]string{"approve": "审批通过", "reject": "驳回"}[action]),
	}, nil
}

func (e *toolExecutor) cancelAppointment(args map[string]interface{}, userID int) (interface{}, error) {
	id, _ := args["id"].(float64)
	reason, _ := args["reason"].(string)

	if id == 0 {
		return nil, fmt.Errorf("id is required")
	}

	// 获取预约单
	var appointment map[string]interface{}
	if err := e.db.Table("construction_appointments").Where("id = ?", int(id)).First(&appointment).Error; err != nil {
		return nil, fmt.Errorf("预约单不存在")
	}

	// 更新状态
	if err := e.db.Table("construction_appointments").Where("id = ?", int(id)).Updates(map[string]interface{}{
		"status":        "cancelled",
		"cancel_reason": reason,
	}).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success":        true,
		"appointment_id": int(id),
		"message":        "预约已取消",
	}, nil
}

func (e *toolExecutor) queryAttendance(args map[string]interface{}) (interface{}, error) {
	userID, _ := args["user_id"].(float64)
	date, _ := args["date"].(string)
	startDate, _ := args["start_date"].(string)
	endDate, _ := args["end_date"].(string)
	attendanceType, _ := args["attendance_type"].(string)
	status, _ := args["status"].(string)
	limit := 20
	if l, ok := args["page_size"].(float64); ok {
		limit = int(l)
	}

	// 处理特殊日期值
	if date == "today" {
		date = time.Now().Format("2006-01-02")
	} else if date == "yesterday" {
		date = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	}

	query := e.db.Table("attendance_records ar").
		Select(`ar.id, ar.user_id, u.full_name as user_name, ar.attendance_type, ar.clock_in_time,
			ar.clock_in_location, ar.overtime_hours, ar.status, ar.work_content,
			ar.appointment_id, ca.appointment_no`).
		Joins("LEFT JOIN users u ON ar.user_id = u.id").
		Joins("LEFT JOIN construction_appointments ca ON ar.appointment_id = ca.id").
		Limit(limit)

	if userID > 0 {
		query = query.Where("ar.user_id = ?", int(userID))
	}
	if date != "" {
		query = query.Where("DATE(ar.clock_in_time) = ?", date)
	}
	if startDate != "" {
		query = query.Where("DATE(ar.clock_in_time) >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("DATE(ar.clock_in_time) <= ?", endDate)
	}
	if attendanceType != "" {
		query = query.Where("ar.attendance_type = ?", attendanceType)
	}
	if status != "" {
		query = query.Where("ar.status = ?", status)
	}

	var results []map[string]interface{}
	if err := query.Order("ar.clock_in_time DESC").Find(&results).Error; err != nil {
		return nil, err
	}

	// 统计各类型数量
	morningCount := 0
	afternoonCount := 0
	noonOvertimeHours := 0.0
	nightOvertimeHours := 0.0

	for _, r := range results {
		switch r["attendance_type"] {
		case "morning":
			morningCount++
		case "afternoon":
			afternoonCount++
		case "noon_overtime":
			if hrs, ok := r["overtime_hours"].(float64); ok {
				noonOvertimeHours += hrs
			}
		case "night_overtime":
			if hrs, ok := r["overtime_hours"].(float64); ok {
				nightOvertimeHours += hrs
			}
		}
	}

	return map[string]interface{}{
		"success":              true,
		"data":                 results,
		"count":                len(results),
		"morning_count":        morningCount,
		"afternoon_count":      afternoonCount,
		"noon_overtime_hours":  noonOvertimeHours,
		"night_overtime_hours": nightOvertimeHours,
		"date":                 date,
	}, nil
}

func (e *toolExecutor) queryAttendanceStats(args map[string]interface{}) (interface{}, error) {
	userID, _ := args["user_id"].(float64)
	year, _ := args["year"].(float64)
	month, _ := args["month"].(float64)

	if year == 0 || month == 0 {
		return nil, fmt.Errorf("year and month are required")
	}

	query := e.db.Table("monthly_attendance_summary mas").
		Select(`mas.id, mas.user_id, u.full_name as user_name, mas.year, mas.month,
			mas.morning_count, mas.afternoon_count, mas.noon_overtime_hours,
			mas.night_overtime_hours, mas.total_work_days, mas.total_overtime_hours,
			mas.status`).
		Joins("LEFT JOIN users u ON mas.user_id = u.id").
		Where("mas.year = ? AND mas.month = ?", int(year), int(month))

	if userID > 0 {
		query = query.Where("mas.user_id = ?", int(userID))
	}

	var results []map[string]interface{}
	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success": true,
		"data":    results,
		"count":   len(results),
		"year":    int(year),
		"month":   int(month),
	}, nil
}

func (e *toolExecutor) queryConstructionLogs(args map[string]interface{}) (interface{}, error) {
	projectID, _ := args["project_id"].(float64)
	logDate, _ := args["log_date"].(string)
	startDate, _ := args["start_date"].(string)
	endDate, _ := args["end_date"].(string)
	limit := 20
	if l, ok := args["page_size"].(float64); ok {
		limit = int(l)
	}

	// 处理特殊日期值
	if logDate == "today" {
		logDate = time.Now().Format("2006-01-02")
	} else if logDate == "yesterday" {
		logDate = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	}

	query := e.db.Table("construction_log cl").
		Select(`cl.id, cl.title, cl.content, cl.images, cl.weather, cl.temperature,
			cl.progress, cl.issues, cl.log_date, cl.project_id, p.name as project_name,
			cl.creator_id, u.full_name as creator_name, cl.created_at`).
		Joins("LEFT JOIN projects p ON cl.project_id = p.id").
		Joins("LEFT JOIN users u ON cl.creator_id = u.id").
		Limit(limit)

	if projectID > 0 {
		query = query.Where("cl.project_id = ?", int(projectID))
	}
	if logDate != "" {
		query = query.Where("cl.log_date = ?", logDate)
	}
	if startDate != "" {
		query = query.Where("cl.log_date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("cl.log_date <= ?", endDate)
	}

	var results []map[string]interface{}
	if err := query.Order("cl.log_date DESC, cl.created_at DESC").Find(&results).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success": true,
		"data":    results,
		"count":   len(results),
	}, nil
}

// Helper functions

func convertToolCalls(tcs []openai.ToolCallInfo) []ToolCallInfo {
	result := make([]ToolCallInfo, len(tcs))
	for i, tc := range tcs {
		result[i] = ToolCallInfo{
			ID:        tc.ID,
			Name:      tc.Name,
			Arguments: tc.Arguments,
			Result:    tc.Result,
		}
	}
	return result
}

// AIChatRequest represents an AI chat request
type AIChatRequest struct {
	Message             string                   `json:"message"`
	ConversationHistory []openai.Message         `json:"conversation_history,omitempty"`
	Context             map[string]interface{}   `json:"context,omitempty"`
	UserID              int                      `json:"user_id"`
}

// AIChatResponse represents an AI chat response
type AIChatResponse struct {
	Message      string                   `json:"message"`
	ToolCalls    []ToolCallInfo           `json:"tool_calls,omitempty"`
	NeedsAction  bool                     `json:"needs_action"`
	ActionResult map[string]interface{}   `json:"action_result,omitempty"`
	Conversation []openai.Message         `json:"conversation"`
}

// ToolCallInfo represents tool call information
type ToolCallInfo struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
	Result    interface{}            `json:"result,omitempty"`
}
