package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

// AIService provides AI-powered operations
type AIService struct {
	client *Client
}

// NewAIService creates a new AI service
func NewAIService(apiKey, model string) *AIService {
	return &AIService{
		client: NewClient(apiKey, model),
	}
}

// NewAIServiceWithBaseURL creates a new AI service with custom base URL
func NewAIServiceWithBaseURL(apiKey, model, baseURL string) *AIService {
	return &AIService{
		client: NewClientWithBaseURL(apiKey, model, baseURL),
	}
}

// ChatRequest represents a chat request with context
type AIChatRequest struct {
	Message             string                   `json:"message"`
	ConversationHistory []Message                `json:"conversation_history,omitempty"`
	Context             map[string]interface{}   `json:"context,omitempty"`
	UserID              int                      `json:"user_id,omitempty"`
}

// AIChatResponse represents the AI response
type AIChatResponse struct {
	Message      string                   `json:"message"`
	ToolCalls    []ToolCallInfo           `json:"tool_calls,omitempty"`
	NeedsAction  bool                     `json:"needs_action"`
	ActionResult map[string]interface{}   `json:"action_result,omitempty"`
	Conversation []Message                `json:"conversation"`
}

// ToolCallInfo represents information about a tool call
type ToolCallInfo struct {
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
	Result    interface{}            `json:"result,omitempty"`
}

// ToolExecutor defines the interface for executing tools
type ToolExecutor interface {
	ExecuteToolCall(ctx context.Context, name string, arguments map[string]interface{}, userID int) (interface{}, error)
}

// ChatWithTools sends a message and handles tool calls
func (s *AIService) ChatWithTools(ctx context.Context, req *AIChatRequest, executor ToolExecutor) (*AIChatResponse, error) {
	// Build system prompt
	systemPrompt := s.buildSystemPrompt(req.Context)

	// Build messages
	messages := req.ConversationHistory
	if len(messages) == 0 {
		messages = []Message{}
	}
	messages = append(messages, Message{Role: "user", Content: req.Message})

	// Get tools
	tools := GetMaterialManagementTools()

	// Initial call
	resp, err := s.client.ChatWithTools(systemPrompt, messages, tools)
	if err != nil {
		return nil, fmt.Errorf("chat completion failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}

	choice := resp.Choices[0]

	// Check if AI wants to call tools
	if choice.FinishReason == "tool_calls" && len(choice.Message.ToolCalls) > 0 {
		// Add the assistant message to conversation
		messages = append(messages, choice.Message)

		// Execute each tool call
		toolCallInfos := []ToolCallInfo{}
		var lastResult interface{}

		for _, toolCall := range choice.Message.ToolCalls {
			// Parse arguments
			var args map[string]interface{}
			if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args); err != nil {
				log.Printf("Failed to parse tool arguments: %v", err)
				messages = append(messages, FormatToolResult(toolCall.ID, "Invalid arguments", true))
				continue
			}

			toolCallInfo := ToolCallInfo{
				ID:        toolCall.ID,
				Name:      toolCall.Function.Name,
				Arguments: args,
			}

			// Execute tool
			result, err := executor.ExecuteToolCall(ctx, toolCall.Function.Name, args, req.UserID)
			if err != nil {
				toolCallInfo.Result = map[string]interface{}{"error": err.Error()}
				messages = append(messages, FormatToolResult(toolCall.ID, err.Error(), true))
			} else {
				toolCallInfo.Result = result
				resultJSON, _ := json.Marshal(result)
				messages = append(messages, FormatToolResult(toolCall.ID, string(resultJSON), false))
			}

			toolCallInfos = append(toolCallInfos, toolCallInfo)
			lastResult = result
		}

		// Get final response after tool calls
		finalResp, err := s.client.ChatWithTools(systemPrompt, messages, tools)
		if err != nil {
			return nil, fmt.Errorf("final chat completion failed: %w", err)
		}

		if len(finalResp.Choices) > 0 {
			messages = append(messages, finalResp.Choices[0].Message)

			// 安全处理 ActionResult
			var actionResult map[string]interface{}
			if lastResult != nil {
				if m, ok := lastResult.(map[string]interface{}); ok {
					actionResult = m
				}
			}

			return &AIChatResponse{
				Message:      finalResp.Choices[0].Message.Content,
				ToolCalls:    toolCallInfos,
				NeedsAction:  false,
				ActionResult: actionResult,
				Conversation: messages,
			}, nil
		}
	}

	// No tool calls, return the message directly
	messages = append(messages, choice.Message)
	return &AIChatResponse{
		Message:     choice.Message.Content,
		NeedsAction: false,
		Conversation: messages,
	}, nil
}

// Chat sends a simple chat message without tools
func (s *AIService) Chat(ctx context.Context, message string, context map[string]interface{}) (string, error) {
	systemPrompt := s.buildSystemPrompt(context)

	req := &ChatRequest{
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: message},
		},
		Temperature: 0.7,
	}

	resp, err := s.client.ChatCompletion(req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	return resp.Choices[0].Message.Content, nil
}

// buildSystemPrompt builds the system prompt based on context
func (s *AIService) buildSystemPrompt(context map[string]interface{}) string {
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

## 当前日期时间
- 今天: ` + today + ` (` + weekday[now.Weekday()] + `)
- 明天: ` + tomorrow + `
- 昨天: ` + yesterday + `

## 你的核心能力

### 1. 施工预约管理
- 查询施工预约（按日期、项目、人员筛选）
- 创建施工预约（选择项目、日期、作业人员）
- 查询作业人员日历和空闲时间
- 预约状态跟踪、审批流程
- 加急任务处理

**预约状态说明**：
- draft: 草稿 - 新建但未提交
- pending: 待审批 - 已提交等待审批
- scheduled: 已排期 - 审批通过，已分配人员
- in_progress: 进行中 - 任务已开始
- completed: 已完成 - 任务已完成
- cancelled: 已取消 - 任务已取消

### 2. 库存管理
- 查询物资信息（名称、规格、分类、价格）
- 查询库存状态（数量、安全库存、低库存预警）
- 库存分析和预警
- 入库/出库记录查询

### 3. 物资计划
- 物资需求计划创建与审批
- 物资主数据管理
- 领用单/出库单管理

### 4. 项目管理
- 查询项目信息（状态、进度、成员）
- 项目进度分析

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

## 回复原则（语音交互优化）

1. **简洁明了**：回复要简短直接，适合语音播报，避免冗长
2. **关键信息优先**：把最重要的信息放在开头
3. **自然对话**：使用口语化表达，像真人对话一样自然
4. **确认机制**：执行重要操作前，简要确认用户意图
5. **错误友好**：如果无法完成操作，给出明确的建议

## 示例对话

用户："帮我查一下明天的预约"
助手："明天有3个施工预约。第一个是水电改造，上午9点在A栋进行，张三负责。需要我详细说明其他预约吗？"

用户："给李四安排后天的施工"
助手："好的，后天是3月26日。请告诉我：1. 具体时间？2. 施工内容？3. 哪个项目？"

用户："库存够吗"
助手："目前有2种物资库存偏低：水泥剩50袋，安全库存100袋；钢筋剩8吨，安全库存15吨。需要我帮您创建采购计划吗？"

用户："今天的考勤情况"
助手："今天有5人打卡。上午4人打卡，下午3人打卡。有2人中午加班共3小时。需要我详细说明每个人的打卡时间吗？"

用户："查看某人的施工日志"
助手："请告诉我：1. 查看哪个项目的日志？2. 查看哪一天的日志？"

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
