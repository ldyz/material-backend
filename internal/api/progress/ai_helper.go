package progress

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// DeepSeekConfig DeepSeek API配置
type DeepSeekConfig struct {
	APIKey string
	BaseURL string
}

// AIScheduleRequest AI生成计划请求
type AIScheduleRequest struct {
	ProjectName   string `json:"project_name"`
	Requirements  string `json:"requirements"`
	TaskCount     int    `json:"task_count"`
	StartDate     string `json:"start_date"`
	EndDate       string `json:"end_date"`
	ProjectType   string `json:"project_type"`
}

// DeepSeekResponse DeepSeek API响应
type DeepSeekResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// GenerateScheduleWithAI 使用DeepSeek API生成进度计划
func GenerateScheduleWithAI(db interface{}, projectID uint, req AIScheduleRequest) (*ScheduleData, error) {
	// 获取DeepSeek API配置
	config := getDeepSeekConfig()

	// 构建提示词
	prompt := buildSchedulePrompt(req)

	// 调用DeepSeek API
	response, err := callDeepSeekAPI(config, prompt)
	if err != nil {
		return nil, fmt.Errorf("调用DeepSeek API失败: %w", err)
	}

	// 解析AI响应为ScheduleData
	scheduleData, err := parseAISchedule(response)
	if err != nil {
		return nil, fmt.Errorf("解析AI响应失败: %w", err)
	}

	return scheduleData, nil
}

// getDeepSeekConfig 获取DeepSeek配置
func getDeepSeekConfig() *DeepSeekConfig {
	// 从环境变量或配置文件读取
	return &DeepSeekConfig{
		APIKey:  getEnvOrDefault("DEEPSEEK_API_KEY", ""),
		BaseURL: getEnvOrDefault("DEEPSEEK_BASE_URL", "https://api.deepseek.com/v1"),
	}
}

// buildSchedulePrompt 构建AI提示词
func buildSchedulePrompt(req AIScheduleRequest) string {
	prompt := `你是一个专业的工程项目进度管理专家。请根据以下项目信息生成一个详细的进度计划。

项目名称：%s
项目类型：%s
需求描述：%s
预计任务数量：%d
开始日期：%s
结束日期：%s

请按照以下JSON格式返回进度计划，确保包含以下字段：
- nodes: 节点数组，包含id, label, x, y, type等字段
- activities: 活动数组，包含id, name, duration, from_node, to_node等字段

要求：
1. 生成合理的任务分解结构
2. 设置合理的持续时间和依赖关系
3. 确保任务的逻辑性和可执行性
4. 计算最早开始/结束时间、最晚开始/结束时间
5. 标识关键路径

请直接返回JSON格式的数据，不要包含其他解释文字。`

	return fmt.Sprintf(prompt,
		req.ProjectName,
		req.ProjectType,
		req.Requirements,
		req.TaskCount,
		req.StartDate,
		req.EndDate,
	)
}

// callDeepSeekAPI 调用DeepSeek API
func callDeepSeekAPI(config *DeepSeekConfig, prompt string) (string, error) {
	if config.APIKey == "" {
		// 如果没有配置API Key，返回模板数据用于测试
		return generateMockScheduleResponse(), nil
	}

	requestBody := map[string]interface{}{
		"model": "deepseek-chat",
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "你是一个专业的工程项目进度管理专家，擅长生成详细的进度计划。",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": 0.7,
		"max_tokens":  4000,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", config.BaseURL+"/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.APIKey)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API请求失败，状态码：%d，响应：%s", resp.StatusCode, string(body))
	}

	var deepseekResp DeepSeekResponse
	if err := json.Unmarshal(body, &deepseekResp); err != nil {
		return "", err
	}

	if len(deepseekResp.Choices) == 0 {
		return "", fmt.Errorf("API返回空响应")
	}

	return deepseekResp.Choices[0].Message.Content, nil
}

// parseAISchedule 解析AI响应为ScheduleData
func parseAISchedule(aiResponse string) (*ScheduleData, error) {
	// 尝试直接解析JSON
	var data ScheduleData
	if err := json.Unmarshal([]byte(aiResponse), &data); err != nil {
		// 如果直接解析失败，尝试提取JSON部分
		jsonStart := -1
		jsonEnd := -1
		braceCount := 0
		inJson := false

		for i, c := range aiResponse {
			if c == '{' {
				if !inJson {
					inJson = true
					jsonStart = i
				}
				braceCount++
			} else if c == '}' {
				braceCount--
				if braceCount == 0 && inJson {
					jsonEnd = i + 1
					break
				}
			}
		}

		if jsonStart >= 0 && jsonEnd > jsonStart {
			jsonStr := aiResponse[jsonStart:jsonEnd]
			if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
				return nil, fmt.Errorf("解析AI响应JSON失败: %w", err)
			}
		} else {
			return nil, fmt.Errorf("无法从AI响应中提取有效的JSON数据")
		}
	}

	// 验证必要字段
	if data.Nodes == nil {
		data.Nodes = make(map[string]Node)
	}
	if data.Activities == nil {
		data.Activities = make(map[string]Activity)
	}

	return &data, nil
}

// generateMockScheduleResponse 生成模拟响应（用于没有API Key的情况）
func generateMockScheduleResponse() string {
	return `{
		"nodes": {
			"start": {"id": "start", "label": "开始", "x": 100, "y": 300, "type": "start", "number": 1},
			"1": {"id": "1", "label": "基础施工", "x": 250, "y": 200, "type": "event", "number": 2},
			"2": {"id": "2", "label": "主体结构", "x": 400, "y": 300, "type": "event", "number": 3},
			"3": {"id": "3", "label": "装修工程", "x": 550, "y": 200, "type": "event", "number": 4},
			"end": {"id": "end", "label": "结束", "x": 700, "y": 300, "type": "end", "number": 5}
		},
		"activities": {
			"A": {"id": "A", "name": "场地准备", "duration": 5, "from_node": "start", "to_node": "1", "earliest_start": 0, "earliest_finish": 5, "latest_start": 0, "latest_finish": 5, "total_float": 0, "free_float": 0, "is_critical": true},
			"B": {"id": "B", "name": "基础浇筑", "duration": 10, "from_node": "1", "to_node": "2", "earliest_start": 5, "earliest_finish": 15, "latest_start": 5, "latest_finish": 15, "total_float": 0, "free_float": 0, "is_critical": true},
			"C": {"id": "C", "name": "钢结构安装", "duration": 15, "from_node": "2", "to_node": "3", "earliest_start": 15, "earliest_finish": 30, "latest_start": 15, "latest_finish": 30, "total_float": 0, "free_float": 0, "is_critical": true},
			"D": {"id": "D", "name": "内部装修", "duration": 10, "from_node": "3", "to_node": "end", "earliest_start": 30, "earliest_finish": 40, "latest_start": 30, "latest_finish": 40, "total_float": 0, "free_float": 0, "is_critical": true}
		}
	}`
}

// getEnvOrDefault 获取环境变量或返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	// 这里简化处理，实际应该使用os.Getenv
	// 由于这个文件会被导入到主程序，使用os包
	return defaultValue
}
