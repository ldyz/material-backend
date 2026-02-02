package system

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/material-backend/backend/internal/api/auth"
	"github.com/yourorg/material-backend/backend/internal/api/response"
	jwtpkg "github.com/yourorg/material-backend/backend/pkg/jwt"
	"gorm.io/gorm"
)

// round 四舍五入到指定小数位
func round(val float64, precision int) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

// DeepSeekConfig DeepSeek配置
type DeepSeekConfig struct {
	Model    string
	BaseURL  string
	APIKey   string
	Timeout  time.Duration
	MaxTokens int
}

// DeepSeekMessage DeepSeek消息格式
type DeepSeekMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// DeepSeekRequest DeepSeek API请求
type DeepSeekRequest struct {
	Model       string            `json:"model"`
	Messages    []DeepSeekMessage `json:"messages"`
	MaxTokens   int               `json:"max_tokens,omitempty"`
	Temperature float64           `json:"temperature,omitempty"`
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

// TableSchemaInfo 表结构信息
type TableSchemaInfo struct {
	Name        string
	Description string
	Columns     []string
	KeyInfo     string
}

// AIAnalyzer AI分析器
type AIAnalyzer struct {
	db     *gorm.DB
	config DeepSeekConfig
}

// NewAIAnalyzer 创建AI分析器
func NewAIAnalyzer(db *gorm.DB) *AIAnalyzer {
	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	if apiKey == "" {
		apiKey = "sk-93923beb1ba543709a514775d7dc2777" // 默认key
	}

	return &AIAnalyzer{
		db: db,
		config: DeepSeekConfig{
			Model:     "deepseek-chat",
			BaseURL:   "https://api.deepseek.com/v1",
			APIKey:    apiKey,
			Timeout:   30 * time.Second,
			MaxTokens: 800,
		},
	}
}

// RegisterAIRoutes 注册AI相关路由
func RegisterAIRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	r := rg.Group("/ai")
	r.Use(jwtpkg.TokenMiddleware())

	analyzer := NewAIAnalyzer(db)

	// AI分析接口
	r.POST("/analyze", auth.PermissionMiddleware(db, "system_statistics"), analyzer.handleAnalyze)
	
	// 获取建议问题
	r.GET("/suggestions", auth.PermissionMiddleware(db, "system_statistics"), analyzer.handleSuggestions)
	
	// 获取AI洞察
	r.GET("/insights", auth.PermissionMiddleware(db, "system_statistics"), analyzer.handleInsights)
	
	// 获取历史记录
	r.GET("/history", auth.PermissionMiddleware(db, "system_statistics"), analyzer.handleHistory)

	// 获取单条历史记录详情
	r.GET("/history/:id", auth.PermissionMiddleware(db, "system_statistics"), analyzer.handleHistoryDetail)

	// 删除历史记录
	r.DELETE("/history/:id", auth.PermissionMiddleware(db, "system_statistics"), analyzer.handleDeleteHistory)

	// 获取统计信息
	r.GET("/stats", auth.PermissionMiddleware(db, "system_statistics"), analyzer.handleStats)
	
	// 获取AI配置
	r.GET("/config", auth.PermissionMiddleware(db, "system_config"), analyzer.handleGetConfig)
	
	// 更新AI配置
	r.POST("/config", auth.PermissionMiddleware(db, "system_config"), analyzer.handleUpdateConfig)
	
	// 检查AI状态
	r.GET("/status", analyzer.handleStatus)
}

// handleAnalyze 处理AI分析请求
func (a *AIAnalyzer) handleAnalyze(c *gin.Context) {
	var req struct {
		Question             string   `json:"question"`
		ConversationMode     bool     `json:"conversation_mode"`
		ConversationID       string   `json:"conversation_id"`
		ConversationHistory  []map[string]string `json:"conversation_history"`
		MaxIterations        int      `json:"max_iterations"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}

	req.Question = strings.TrimSpace(req.Question)
	if req.Question == "" {
		response.BadRequest(c, "问题不能为空")
		return
	}

	if len(req.Question) > 500 {
		response.BadRequest(c, "问题长度不能超过500字符")
		return
	}

	if req.MaxIterations < 1 || req.MaxIterations > 5 {
		req.MaxIterations = 3
	}

	// 获取用户信息
	userID, _ := c.Get("user_id")
	ipAddress := c.ClientIP()

	// 创建分析日志
	var uid uint
	if userID != nil {
		uid = userID.(uint)
	}

	logEntry := &AIAnalysisLog{
		Question:  req.Question,
		UserID:    uid,
		IPAddress: ipAddress,
		Status:    "processing",
	}
	a.db.Create(logEntry)

	startTime := time.Now()

	var result map[string]any
	var err error

	// 根据是否为对话模式选择不同的处理逻辑
	if req.ConversationMode {
		result, err = a.analyzeWithConversation(req.Question, req.ConversationHistory, req.MaxIterations)
	} else {
		result, err = a.analyzeQuestion(req.Question)
	}

	processingTime := time.Since(startTime).Seconds()

	if err != nil {
		logEntry.Status = "failed"
		logEntry.ErrorMessage = err.Error()
		logEntry.ProcessingTime = processingTime
		a.db.Save(logEntry)

		response.InternalError(c, fmt.Sprintf("AI分析失败: %v", err))
		return
	}

	// 更新日志
	logEntry.Answer = result["answer"].(string)
	logEntry.QueryUsed = result["query_used"].(string)
	if dataSummary, ok := result["data_summary"].(string); ok {
		logEntry.DataSummary = dataSummary
	}
	logEntry.ProcessingTime = processingTime
	logEntry.Status = "completed"
	logEntry.ModelUsed = "deepseek-chat"
	a.db.Save(logEntry)

	result["success"] = true
	result["processing_time"] = processingTime

	// 如果是对话模式，返回对话相关信息
	if req.ConversationMode {
		result["conversation_id"] = req.ConversationID
		result["conversation_mode"] = true
		if iterations, ok := result["iterations_used"]; ok {
			result["iterations_used"] = iterations
		}
	}

	response.Success(c, result)
}

// analyzeWithConversation 对话模式分析
func (a *AIAnalyzer) analyzeWithConversation(question string, conversationHistory []map[string]string, maxIterations int) (map[string]any, error) {
	// 构建对话上下文
	var contextSummary string
	if len(conversationHistory) > 0 {
		contextSummary = "\n\n对话历史（用于理解上下文）：\n"
		for i, entry := range conversationHistory {
			if i >= 3 { // 只使用最近3轮对话作为上下文
				break
			}
			if q, ok := entry["question"]; ok {
				contextSummary += fmt.Sprintf("Q: %s\n", q)
			}
			if a, ok := entry["answer"]; ok {
				// 截取答案的关键部分
				if len(a) > 100 {
					contextSummary += fmt.Sprintf("A: %s...\n", a[:100])
				} else {
					contextSummary += fmt.Sprintf("A: %s\n", a)
				}
			}
		}
	}

	// 1. 生成SQL（结合上下文）
	schemas := a.getTableSchemas()
	fullQuestion := question + contextSummary
	sqlQuery, err := a.generateSQL(fullQuestion, schemas)
	if err != nil {
		return nil, fmt.Errorf("生成SQL失败: %w", err)
	}

	// 2. 执行SQL
	data, err := a.executeSQL(sqlQuery)
	if err != nil {
		return nil, fmt.Errorf("执行SQL失败: %w", err)
	}

	// 3. 生成答案（结合上下文）
	answer, err := a.generateAnswerWithContext(question, sqlQuery, data, conversationHistory)
	if err != nil {
		return nil, fmt.Errorf("生成答案失败: %w", err)
	}

	return map[string]any{
		"answer":         answer,
		"query_used":     sqlQuery,
		"data_summary":   fmt.Sprintf("查询返回 %d 条数据", len(data)),
		"data_count":     len(data),
		"iterations_used": 1,
	}, nil
}

// generateAnswerWithContext 生成带上下文的AI分析答案
func (a *AIAnalyzer) generateAnswerWithContext(question, sqlQuery string, data []map[string]any, conversationHistory []map[string]string) (string, error) {
	dataSample := data
	if len(data) > 20 {
		dataSample = data[:20]
	}

	dataJSON, _ := json.Marshal(dataSample)

	// 构建上下文信息
	contextInfo := ""
	if len(conversationHistory) > 0 {
		contextInfo = "\n\n【对话上下文参考】\n用户之前询问过相关问题，回答时可以适当关联之前的信息，但重点应放在当前问题上。"
	}

	prompt := fmt.Sprintf(`【当前问题】%s
【SQL】%s
【样本数据】%s
【总数】%d 条%s

【输出要求】
- 只输出简明扼要的分析结论，避免重复和废话
- 结构化分段：1. 关键统计 2. 趋势/变化 3. 洞察/建议
- 如无数据，直接给出合理建议
- 如果有对话历史，适当关联但以当前问题为主
`, question, sqlQuery, string(dataJSON), len(data), contextInfo)

	messages := []DeepSeekMessage{
		{Role: "system", Content: "提供专业、易懂的中文分析，即使无数据也要给出建议。注意结合对话历史但不被其限制。"},
		{Role: "user", Content: prompt},
	}

	return a.callDeepSeekAPI(messages, 800, 0.3)
}

// analyzeQuestion 分析问题
func (a *AIAnalyzer) analyzeQuestion(question string) (map[string]any, error) {
	// 1. 生成SQL
	schemas := a.getTableSchemas()
	sqlQuery, err := a.generateSQL(question, schemas)
	if err != nil {
		return nil, fmt.Errorf("生成SQL失败: %w", err)
	}

	// 2. 执行SQL
	data, err := a.executeSQL(sqlQuery)
	if err != nil {
		return nil, fmt.Errorf("执行SQL失败: %w", err)
	}

	// 3. 生成答案
	answer, err := a.generateAnswer(question, sqlQuery, data)
	if err != nil {
		return nil, fmt.Errorf("生成答案失败: %w", err)
	}

	return map[string]any{
		"answer":       answer,
		"query_used":   sqlQuery,
		"data_summary": fmt.Sprintf("查询返回 %d 条数据", len(data)),
		"data_count":   len(data),
	}, nil
}

// getTableSchemas 获取表结构信息
func (a *AIAnalyzer) getTableSchemas() []TableSchemaInfo {
	return []TableSchemaInfo{
		{
			Name:        "materials",
			Description: "物资表，存储所有物资信息",
			Columns:     []string{"id", "code", "name", "specification", "spec", "unit", "price", "description", "category", "quantity", "project_id", "material", "created_at", "updated_at"},
			KeyInfo:     "id:物资ID, name:物资名称, code:物资编码, specification:规格说明, spec:规格型号, unit:单位, category:分类, price:单价, quantity:数量, project_id:项目ID(关联projects.id)",
		},
		{
			Name:        "stocks",
			Description: "库存表，记录物资库存情况。通过material_id关联materials表获取物资详细信息",
			Columns:     []string{"id", "material_id", "quantity", "safety_stock", "location", "unit", "created_at", "updated_at"},
			KeyInfo:     "id:库存ID, material_id:物资ID(关联materials.id), quantity:当前库存数量, safety_stock:安全库存阈值, location:库位, unit:单位。注意：需要JOIN materials表获取物资name、code、specification等信息",
		},
		{
			Name:        "requisitions",
			Description: "领料单表，记录物资领用申请。需要JOIN projects表获取项目名称",
			Columns:     []string{"id", "requisition_no", "project_id", "applicant", "department", "status", "created_at", "remark", "approved_at", "approved_by", "urgent", "purpose", "issued_by", "issued_at"},
			KeyInfo:     "id:领料单ID, requisition_no:领料单号, project_id:项目ID(关联projects.id), applicant:申请人, department:申请部门, status:状态(pending待审批/approved已批准/rejected已拒绝/completed已完成), urgent:紧急程度(0普通1紧急), purpose:用途, approved_at:审批时间, issued_at:发放时间。注意：需要JOIN projects表获取项目名称",
		},
		{
			Name:        "inbound_orders",
			Description: "入库单表，记录物资入库情况",
			Columns:     []string{"id", "order_no", "supplier", "contact", "project_id", "creator_id", "creator_name", "status", "notes", "remark", "total_amount", "created_at", "updated_at"},
			KeyInfo:     "id:入库单ID, order_no:入库单号, supplier:供应商名称, contact:联系方式, project_id:项目ID(关联projects.id), creator_id:创建人ID(关联users.id), creator_name:创建人姓名, status:状态(pending待入库/completed已入库), total_amount:总金额(分为单位), notes:备注, remark:说明",
		},
		{
			Name:        "projects",
			Description: "项目表，存储所有项目信息",
			Columns:     []string{"id", "name", "code", "location", "start_date", "end_date", "description", "manager", "contact", "budget", "status", "created_at", "updated_at"},
			KeyInfo:     "id:项目ID, name:项目名称, code:项目编码, location:项目地点, start_date:开始日期, end_date:结束日期, description:项目描述, manager:项目经理, contact:联系方式, budget:预算, status:状态(planning策划中/active进行中/completed已完工)",
		},
		{
			Name:        "users",
			Description: "用户表，存储系统用户信息",
			Columns:     []string{"id", "username", "email", "full_name", "role", `"group"`, "is_active", "last_login", "created_at"},
			KeyInfo:     "id:用户ID, username:用户名(唯一), email:邮箱, full_name:真实姓名, role:角色(admin管理员/user普通用户/manager经理), \"group\":用户组(需要双引号转义字段名), is_active:是否激活(true/false), last_login:最后登录时间",
		},
	}
}

// generateSQL 生成SQL查询
func (a *AIAnalyzer) generateSQL(question string, schemas []TableSchemaInfo) (string, error) {
	schemaInfo := ""
	for _, schema := range schemas {
		schemaInfo += fmt.Sprintf("表名: %s\n描述: %s\n字段: %s\n重要信息: %s\n\n",
			schema.Name, schema.Description, strings.Join(schema.Columns, ", "), schema.KeyInfo)
	}

	prompt := fmt.Sprintf(`数据库表结构：
%s

用户问题：%s

规则：
1. 生成PostgreSQL兼容的SQL
2. 结果限制为 500 条
3. 使用中文别名
4. 默认查询最近30天
生成SQL：`, schemaInfo, question)

	messages := []DeepSeekMessage{
		{Role: "system", Content: "生成准确的PostgreSQL SQL查询，只返回SQL语句"},
		{Role: "user", Content: prompt},
	}

	response, err := a.callDeepSeekAPI(messages, 500, 0.1)
	if err != nil {
		return "", err
	}

	return a.cleanSQL(response), nil
}

// cleanSQL 清理和修正SQL
func (a *AIAnalyzer) cleanSQL(sql string) string {
	// 移除代码块标记
	sql = regexp.MustCompile("```sql|```").ReplaceAllString(sql, "")
	sql = strings.TrimSpace(sql)

	// 移除注释
	lines := strings.Split(sql, "\n")
	cleanedLines := []string{}
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "--") && !strings.HasPrefix(line, "#") {
			cleanedLines = append(cleanedLines, line)
		}
	}
	sql = strings.Join(cleanedLines, " ")

	// 字段名修正
	corrections := map[string]string{
		"min_quantity":                            "safety_stock",
		"real_name":                               "full_name",
		"DATE_SUB(CURRENT_DATE(), INTERVAL 30 DAY)": "CURRENT_DATE - INTERVAL '30 days'",
	}
	for old, new := range corrections {
		sql = strings.ReplaceAll(sql, old, new)
	}

	// 转义 group 字段（但不转义 GROUP BY）
	sql = regexp.MustCompile(`([\w\"]+)\.group(\s+AS|\s*,|\s*\)|\s*=|\s+|$)`).ReplaceAllString(sql, `$1."group"$2`)
	// 只转义非 GROUP BY 的 group（排除 GROUP 关键字）
	sql = regexp.MustCompile(`(?i)(?:^|[^\w\.])group(\s+AS|\s*,|\s*\)|\s*=)`).ReplaceAllString(sql, `"group"$1`)

	return sql
}

// executeSQL 执行SQL查询
func (a *AIAnalyzer) executeSQL(sqlQuery string) ([]map[string]any, error) {
	// 获取第一条SQL语句
	mainSQL := strings.Split(sqlQuery, ";")[0]
	mainSQL = strings.TrimSpace(mainSQL)
	if mainSQL == "" {
		return []map[string]any{}, nil
	}

	// 验证SQL安全性
	if err := a.validateSQL(mainSQL); err != nil {
		return nil, err
	}

	// 执行查询
	var results []map[string]any
	rows, err := a.db.Raw(mainSQL).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		values := make([]any, len(columns))
		valuePtrs := make([]any, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		row := make(map[string]any)
		for i, col := range columns {
			val := values[i]
			// 将字节切片转换为字符串
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		results = append(results, row)
	}

	return results, nil
}

// validateSQL 验证SQL安全性
func (a *AIAnalyzer) validateSQL(sql string) error {
	dangerousKeywords := []string{"DROP", "DELETE", "UPDATE", "INSERT", "ALTER", "CREATE", "TRUNCATE"}
	upperSQL := strings.ToUpper(sql)

	for _, keyword := range dangerousKeywords {
		pattern := `\b` + keyword + `\b`
		matched, _ := regexp.MatchString(pattern, upperSQL)
		if matched {
			return fmt.Errorf("SQL包含危险操作: %s", keyword)
		}
	}

	return nil
}

// generateAnswer 生成AI分析答案
func (a *AIAnalyzer) generateAnswer(question, sqlQuery string, data []map[string]any) (string, error) {
	dataSample := data
	if len(data) > 20 {
		dataSample = data[:20]
	}

	dataJSON, _ := json.Marshal(dataSample)
	prompt := fmt.Sprintf(`【问题】%s
【SQL】%s
【样本数据】%s
【总数】%d 条

【输出要求】
- 只输出简明扼要的分析结论，避免重复和废话
- 结构化分段：1. 关键统计 2. 趋势/变化 3. 洞察/建议
- 如无数据，直接给出合理建议
`, question, sqlQuery, string(dataJSON), len(data))

	messages := []DeepSeekMessage{
		{Role: "system", Content: "提供专业、易懂的中文分析，即使无数据也要给出建议"},
		{Role: "user", Content: prompt},
	}

	return a.callDeepSeekAPI(messages, 800, 0.3)
}

// callDeepSeekAPI 调用DeepSeek API
func (a *AIAnalyzer) callDeepSeekAPI(messages []DeepSeekMessage, maxTokens int, temperature float64) (string, error) {
	reqBody := DeepSeekRequest{
		Model:       a.config.Model,
		Messages:    messages,
		MaxTokens:   maxTokens,
		Temperature: temperature,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: a.config.Timeout}
	req, err := http.NewRequest("POST", a.config.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+a.config.APIKey)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API调用失败 (状态码: %d): %s", resp.StatusCode, string(body))
	}

	var deepseekResp DeepSeekResponse
	if err := json.NewDecoder(resp.Body).Decode(&deepseekResp); err != nil {
		return "", err
	}

	if len(deepseekResp.Choices) == 0 {
		return "", fmt.Errorf("API返回空响应")
	}

	return deepseekResp.Choices[0].Message.Content, nil
}

// handleSuggestions 获取建议问题
func (a *AIAnalyzer) handleSuggestions(c *gin.Context) {
	suggestionType := c.DefaultQuery("type", "question")

	switch suggestionType {
	case "questions":
		// 返回快捷问题列表
		questions := []string{
			"最近30天有多少条领料记录？",
			"库存量低于安全库存的物资有哪些？",
			"本月的入库单总金额是多少？",
			"哪些项目最近活跃？",
			"用户部门分布情况如何？",
			"最常用的物资是哪些？",
			"哪些物资需要补货？",
			"本周的领料趋势如何？",
		}
		response.SuccessWithMeta(c, questions, map[string]any{"suggestions": questions})

	case "recommendations":
		// 返回智能推荐建议（真实数据）
		recommendations := a.generateRealRecommendations()
		response.SuccessWithMeta(c, recommendations, map[string]any{"recommendations": recommendations})

	default:
		// 默认返回问题列表（向后兼容）
		defaultQuestions := []string{
			"最近30天有多少条领料记录？",
			"库存量低于安全库存的物资有哪些？",
			"本月的入库单总金额是多少？",
		}
		response.SuccessWithMeta(c, defaultQuestions, map[string]any{"suggestions": defaultQuestions})
	}
}

// handleInsights 获取AI洞察
func (a *AIAnalyzer) handleInsights(c *gin.Context) {
	dataType := c.DefaultQuery("type", "all")

	// 处理dashboard类型（前端需要的统一格式）
	if dataType == "dashboard" {
		a.handleDashboardInsights(c)
		return
	}

	insights := make(map[string]any)

	// 库存洞察
	if dataType == "all" || dataType == "inventory" {
		var totalItems, lowStockItems, outOfStockItems int64
		var avgQuantity float64

		// 总物资种类
		a.db.Raw(`
			SELECT COUNT(DISTINCT material_id) FROM stocks
		`).Scan(&totalItems)

		// 低库存物资
		a.db.Raw(`
			SELECT COUNT(*) FROM stocks
			WHERE quantity < safety_stock AND safety_stock > 0
		`).Scan(&lowStockItems)

		// 缺货物资
		a.db.Raw(`
			SELECT COUNT(*) FROM stocks
			WHERE quantity <= 0
		`).Scan(&outOfStockItems)

		// 平均库存量
		a.db.Raw(`
			SELECT AVG(quantity) FROM stocks
		`).Scan(&avgQuantity)

		insights["inventory"] = map[string]any{
			"total_items":        totalItems,
			"low_stock_items":    lowStockItems,
			"out_of_stock_items": outOfStockItems,
			"avg_quantity":       avgQuantity,
			"health_score":       calculateHealthScore(totalItems, lowStockItems, outOfStockItems),
		}
	}

	// 领料单洞察
	if dataType == "all" || dataType == "requisitions" {
		var reqStats []struct {
			Status string
			Count  int64
		}

		a.db.Raw(`
			SELECT status, COUNT(*) as count
			FROM requisitions
			WHERE created_at >= CURRENT_DATE - INTERVAL '30 days'
			GROUP BY status
		`).Scan(&reqStats)

		totalReqs := int64(0)
		statusDistribution := make([]map[string]any, 0)
		for _, stat := range reqStats {
			totalReqs += stat.Count
			statusDistribution = append(statusDistribution, map[string]any{
				"status":     stat.Status,
				"count":      stat.Count,
				"percentage": 0.0, // Will calculate below
			})
		}

		// Calculate percentages
		for _, stat := range statusDistribution {
			if totalReqs > 0 {
				stat["percentage"] = float64(stat["count"].(int64)) / float64(totalReqs) * 100
			}
		}

		insights["requisitions"] = statusDistribution
	}

	// 用户洞察
	if dataType == "all" || dataType == "users" {
		var totalUsers, activeUsers, recentActiveUsers int64
		thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

		// 总用户数
		a.db.Model(&struct{ ID uint }{}).Table("users").Count(&totalUsers)

		// 活跃用户（最近登录）
		a.db.Raw(`
			SELECT COUNT(*) FROM users
			WHERE last_login >= ?
		`, thirtyDaysAgo).Scan(&activeUsers)

		// 近期活跃（7天内）
		sevenDaysAgo := time.Now().AddDate(0, 0, -7)
		a.db.Raw(`
			SELECT COUNT(*) FROM users
			WHERE last_login >= ?
		`, sevenDaysAgo).Scan(&recentActiveUsers)

		insights["users"] = map[string]any{
			"total_users":       totalUsers,
			"active_users":      activeUsers,
			"recent_active_users": recentActiveUsers,
		}
	}

	// 项目洞察
	if dataType == "all" || dataType == "projects" {
		var activeProjects, totalProjects int64

		a.db.Model(&struct{ ID uint }{}).Table("projects").Count(&totalProjects)

		a.db.Raw(`
			SELECT COUNT(*) FROM projects
			WHERE status = 'active' OR status = '进行中'
		`).Scan(&activeProjects)

		insights["projects"] = map[string]any{
			"total_projects":  totalProjects,
			"active_projects": activeProjects,
		}
	}

	// 入库洞察
	if dataType == "all" || dataType == "inbound" {
		var totalInbound, totalAmount float64
		monthAgo := time.Now().AddDate(0, 0, -30)

		a.db.Raw(`
			SELECT COUNT(*) FROM inbound_orders
			WHERE created_at >= ?
		`, monthAgo).Scan(&totalInbound)

		a.db.Raw(`
			SELECT COALESCE(SUM(total_amount), 0) FROM inbound_orders
			WHERE created_at >= ?
		`, monthAgo).Scan(&totalAmount)

		insights["inbound"] = map[string]any{
			"total_orders":      totalInbound,
			"total_amount_cents": totalAmount,
			"total_amount":      totalAmount / 100.0, // Convert cents to yuan
		}
	}

	response.SuccessWithMeta(c, insights, map[string]any{"type": dataType})
}

// calculateHealthScore 计算库存健康分数 (0-100)
func calculateHealthScore(total, lowStock, outOfStock int64) float64 {
	if total == 0 {
		return 100
	}

	// Deduct points for issues
	score := 100.0
	if lowStock > 0 {
		score -= float64(lowStock) / float64(total) * 30
	}
	if outOfStock > 0 {
		score -= float64(outOfStock) / float64(total) * 50
	}

	if score < 0 {
		score = 0
	}
	return score
}

// handleDashboardInsights 处理dashboard类型的数据洞察
func (a *AIAnalyzer) handleDashboardInsights(c *gin.Context) {
	now := time.Now()
	thisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	lastMonth := thisMonth.AddDate(0, -1, 0)

	// 1. 本月入库金额（分转元）
	var thisMonthInbound, lastMonthInbound float64
	a.db.Raw(`
		SELECT COALESCE(SUM(total_amount), 0) FROM inbound_orders
		WHERE created_at >= ? AND status = 'completed'
	`, thisMonth).Scan(&thisMonthInbound)
	a.db.Raw(`
		SELECT COALESCE(SUM(total_amount), 0) FROM inbound_orders
		WHERE created_at >= ? AND created_at < ? AND status = 'completed'
	`, lastMonth, thisMonth).Scan(&lastMonthInbound)

	// 2. 本月出库金额（从stock_logs计算，关联materials获取价格）
	var thisMonthOutbound, lastMonthOutbound float64
	a.db.Raw(`
		SELECT COALESCE(SUM(sl.quantity * m.price), 0)
		FROM stock_logs sl
		JOIN stocks s ON s.id = sl.stock_id
		JOIN materials m ON m.id = s.material_id
		WHERE sl.time >= ? AND sl.type = 'out'
	`, thisMonth).Scan(&thisMonthOutbound)
	a.db.Raw(`
		SELECT COALESCE(SUM(sl.quantity * m.price), 0)
		FROM stock_logs sl
		JOIN stocks s ON s.id = sl.stock_id
		JOIN materials m ON m.id = s.material_id
		WHERE sl.time >= ? AND sl.time < ? AND sl.type = 'out'
	`, lastMonth, thisMonth).Scan(&lastMonthOutbound)

	// 3. 库存预警数量（库存量 <= 安全库存）
	var alertCount int64
	a.db.Raw(`
		SELECT COUNT(*) FROM stocks
		WHERE quantity <= safety_stock AND safety_stock > 0
	`).Scan(&alertCount)

	// 4. 待审批单据数量
	var pendingInbound, pendingRequisition int64
	a.db.Model(&struct{ ID uint }{}).Table("inbound_orders").
		Where("status = ?", "pending").Count(&pendingInbound)
	a.db.Model(&struct{ ID uint }{}).Table("requisitions").
		Where("status = ?", "pending").Count(&pendingRequisition)
	pendingCount := pendingInbound + pendingRequisition

	// 5. 计算增长率
	inboundGrowth := 0.0
	if lastMonthInbound > 0 {
		inboundGrowth = ((thisMonthInbound - lastMonthInbound) / lastMonthInbound) * 100
	}

	outboundGrowth := 0.0
	if lastMonthOutbound > 0 {
		outboundGrowth = ((thisMonthOutbound - lastMonthOutbound) / lastMonthOutbound) * 100
	}

	result := map[string]any{
		"total_inbound":   thisMonthInbound / 100.0, // 分转元
		"total_outbound":  thisMonthOutbound,
		"inbound_growth":  round(inboundGrowth, 2),
		"outbound_growth": round(outboundGrowth, 2),
		"alert_count":     alertCount,
		"pending_count":   pendingCount,
	}

	response.Success(c, result)
}

// handleHistory 获取历史记录
func (a *AIAnalyzer) handleHistory(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")
	myOnly := c.DefaultQuery("my_only", "false")

	query := a.db.Model(&AIAnalysisLog{})

	if myOnly == "true" {
		if userID, exists := c.Get("user_id"); exists {
			query = query.Where("user_id = ?", userID)
		}
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var logs []AIAnalysisLog
	query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&logs)

	history := make([]map[string]any, 0, len(logs))
	for _, log := range logs {
		history = append(history, map[string]any{
			"id":              log.ID,
			"question":        log.Question,
			"answer":          log.Answer,
			"data_summary":    log.DataSummary,
			"processing_time": log.ProcessingTime,
			"status":          log.Status,
			"error_message":   log.ErrorMessage,
			"model_used":      log.ModelUsed,
			"created_at":      log.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	response.SuccessWithPagination(c, history, int64(page), int64(pageSize), total)
}

// handleHistoryDetail 获取单条历史记录详情
func (a *AIAnalyzer) handleHistoryDetail(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "记录ID不能为空")
		return
	}

	var log AIAnalysisLog
	if err := a.db.Where("id = ?", id).First(&log).Error; err != nil {
		response.NotFound(c, "记录不存在")
		return
	}

	record := map[string]any{
		"id":              log.ID,
		"question":        log.Question,
		"answer":          log.Answer,
		"query_used":      log.QueryUsed,
		"data_summary":    log.DataSummary,
		"processing_time": log.ProcessingTime,
		"status":          log.Status,
		"error_message":   log.ErrorMessage,
		"model_used":      log.ModelUsed,
		"tokens_used":     log.TokensUsed,
		"user_id":         log.UserID,
		"ip_address":      log.IPAddress,
		"created_at":      log.CreatedAt.Format("2006-01-02 15:04:05"),
	}
	response.SuccessWithMeta(c, record, map[string]any{"record": record})
}

// handleDeleteHistory 删除历史记录
func (a *AIAnalyzer) handleDeleteHistory(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "记录ID不能为空")
		return
	}

	// 检查记录是否存在
	var log AIAnalysisLog
	if err := a.db.Where("id = ?", id).First(&log).Error; err != nil {
		response.NotFound(c, "记录不存在")
		return
	}

	// 删除记录
	if err := a.db.Delete(&log).Error; err != nil {
		response.InternalError(c, "删除失败")
		return
	}

	response.SuccessOnlyMessage(c, "删除成功")
}

// handleStats 获取统计信息
func (a *AIAnalyzer) handleStats(c *gin.Context) {
	today := time.Now().Truncate(24 * time.Hour)
	weekAgo := today.AddDate(0, 0, -7)
	monthAgo := today.AddDate(0, 0, -30)

	var totalQuestions, todayQuestions, weekQuestions, monthQuestions, successful int64

	a.db.Model(&AIAnalysisLog{}).Count(&totalQuestions)
	a.db.Model(&AIAnalysisLog{}).Where("created_at >= ?", today).Count(&todayQuestions)
	a.db.Model(&AIAnalysisLog{}).Where("created_at >= ?", weekAgo).Count(&weekQuestions)
	a.db.Model(&AIAnalysisLog{}).Where("created_at >= ?", monthAgo).Count(&monthQuestions)
	a.db.Model(&AIAnalysisLog{}).Where("status = ?", "completed").Count(&successful)

	successRate := 0.0
	if totalQuestions > 0 {
		successRate = float64(successful) / float64(totalQuestions) * 100
	}

	var avgTime float64
	a.db.Model(&AIAnalysisLog{}).
		Where("status = ?", "completed").
		Where("processing_time IS NOT NULL").
		Select("AVG(processing_time)").
		Scan(&avgTime)

	// Count by model
	var modelCounts []struct {
		ModelUsed string
		Count     int64
	}
	a.db.Model(&AIAnalysisLog{}).
		Select("model_used, COUNT(*) as count").
		Where("model_used IS NOT NULL").
		Group("model_used").
		Scan(&modelCounts)

	modelUsage := make(map[string]int64)
	for _, mc := range modelCounts {
		modelUsage[mc.ModelUsed] = mc.Count
	}

	stats := map[string]any{
		"total_questions":     totalQuestions,
		"today_questions":     todayQuestions,
		"week_questions":      weekQuestions,
		"month_questions":     monthQuestions,
		"success_rate":        successRate,
		"avg_processing_time": avgTime,
		"model_usage":         modelUsage,
	}

	response.SuccessWithMeta(c, stats, map[string]any{"stats": stats})
}

// handleGetConfig 获取AI配置
func (a *AIAnalyzer) handleGetConfig(c *gin.Context) {
	var configs []SystemConfig
	a.db.Where("key LIKE ?", "ai_%").Find(&configs)

	configMap := map[string]string{
		"ai_enabled":     "true",
		"ai_model":       "deepseek-chat",
		"max_iterations": "3",
		"timeout":        "30",
	}

	for _, cfg := range configs {
		configMap[cfg.Key] = cfg.Value
	}

	response.SuccessWithMeta(c, configMap, map[string]any{"config": configMap})
}

// handleUpdateConfig 更新AI配置
func (a *AIAnalyzer) handleUpdateConfig(c *gin.Context) {
	var data map[string]string
	if err := c.ShouldBindJSON(&data); err != nil {
		response.BadRequest(c, "配置数据格式错误")
		return
	}

	for key, value := range data {
		if strings.HasPrefix(key, "ai_") {
			var config SystemConfig
			result := a.db.Where("key = ?", key).First(&config)

			if result.Error == gorm.ErrRecordNotFound {
				config = SystemConfig{
					Key:         key,
					Value:       value,
					Type:        "string",
					Description: fmt.Sprintf("AI配置: %s", key),
				}
				a.db.Create(&config)
			} else {
				config.Value = value
				config.UpdatedAt = time.Now()
				a.db.Save(&config)
			}
		}
	}

	response.SuccessOnlyMessage(c, "AI配置已更新")
}

// generateRealRecommendations 生成真实的推荐建议
func (a *AIAnalyzer) generateRealRecommendations() []map[string]any {
	recommendations := []map[string]any{}

	// 1. 库存预警推荐
	var lowStockCount int64
	a.db.Raw(`
		SELECT COUNT(*) FROM stocks
		WHERE quantity > 0 AND quantity <= safety_stock AND safety_stock > 0
	`).Scan(&lowStockCount)

	var outOfStockCount int64
	a.db.Raw(`
		SELECT COUNT(*) FROM stocks
		WHERE quantity <= 0
	`).Scan(&outOfStockCount)

	if lowStockCount > 0 || outOfStockCount > 0 {
		desc := fmt.Sprintf("有%d种物资库存低于安全库存，%d种物资已缺货，建议及时补货",
			lowStockCount, outOfStockCount)
		recommendations = append(recommendations, map[string]any{
			"type":        "alert",
			"title":       "处理库存预警",
			"description": desc,
			"action":      "/stock",
			"priority":    "high",
		})
	}

	// 2. 待审批推荐
	var pendingInbound int64
	a.db.Model(&struct{ ID uint }{}).Table("inbound_orders").
		Where("status = ?", "pending").Count(&pendingInbound)

	var pendingRequisition int64
	a.db.Model(&struct{ ID uint }{}).Table("requisitions").
		Where("status = ?", "pending").Count(&pendingRequisition)

	totalPending := pendingInbound + pendingRequisition
	if totalPending > 0 {
		desc := fmt.Sprintf("有%d个入库单和%d个领料单待审批，请及时处理",
			pendingInbound, pendingRequisition)
		recommendations = append(recommendations, map[string]any{
			"type":        "approval",
			"title":       "审批待办",
			"description": desc,
			"action":      "/approvals",
			"priority":    "medium",
		})
	}

	// 3. 高价值物资补货建议（基于历史出库）
	var highValueMaterials []struct {
		Name     string
		Quantity float64
		Outbound float64
	}
	a.db.Raw(`
		SELECT m.name, s.quantity, COALESCE(SUM(sl.quantity), 0) as outbound
		FROM materials m
		JOIN stocks s ON s.material_id = m.id
		LEFT JOIN stock_logs sl ON sl.stock_id = s.id AND sl.type = 'out'
			AND sl.time >= CURRENT_DATE - INTERVAL '30 days'
		WHERE m.price > 1000 AND s.quantity < s.safety_stock * 1.5
		GROUP BY m.id, m.name, s.quantity
		ORDER BY m.price DESC
		LIMIT 3
	`).Scan(&highValueMaterials)

	if len(highValueMaterials) > 0 {
		materials := ""
		for i, mat := range highValueMaterials {
			if i > 0 {
				materials += "、"
			}
			materials += mat.Name
		}
		desc := fmt.Sprintf("根据历史出库数据，建议优先关注高价值物资：%s", materials)
		recommendations = append(recommendations, map[string]any{
			"type":        "insight",
			"title":       "采购建议",
			"description": desc,
			"action":      "",
			"priority":    "low",
		})
	}

	return recommendations
}

// handleStatus 检查AI状态
func (a *AIAnalyzer) handleStatus(c *gin.Context) {
	messages := []DeepSeekMessage{
		{Role: "system", Content: "你是一个AI助手。"},
		{Role: "user", Content: "简单回答：你好"},
	}

	_, err := a.callDeepSeekAPI(messages, 10, 0.1)
	if err != nil {
		response.SuccessWithMeta(c, map[string]any{
			"online":  false,
			"status":  "offline",
			"message": fmt.Sprintf("AI服务连接失败: %v", err),
		}, map[string]any{"success": false})
		return
	}

	response.SuccessWithMeta(c, map[string]any{
		"online":        true,
		"status":        "online",
		"message":       "AI服务正常运行",
		"model":         a.config.Model,
		"test_response": true,
	}, nil)
}
