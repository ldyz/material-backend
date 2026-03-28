package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/yourorg/material-backend/backend/internal/api/auth"
)

// ToolDefinition 定义工具的API映射
type ToolDefinition struct {
	Name               string            `json:"name"`
	Description        string            `json:"description"`
	Method             string            `json:"method"`
	Path               string            `json:"path"`
	PathParams         []string          `json:"path_params,omitempty"`
	QueryParams        map[string]string `json:"query_params,omitempty"`
	BodyParams         []string          `json:"body_params,omitempty"`
	RequiresBody       bool              `json:"requires_body,omitempty"`
	ExampleResult      string            `json:"example_result,omitempty"`
	RequiredPermission string            `json:"required_permission,omitempty"` // 所需权限
}

// ToolDefinitions 全局工具定义映射
var ToolDefinitions = map[string]ToolDefinition{
	// ==================== 查询类工具 ====================
	"query_appointments": {
		Name:        "query_appointments",
		Description: `查询施工预约/施工任务/施工作业安排。

当用户询问以下内容时使用此工具：
- "明天的任务"、"今日安排"、"后天有什么"
- "施工计划"、"预约列表"
- "我的任务"、"待处理的预约"

参数说明：
- date: 指定日期 (YYYY-MM-DD格式)
- start_date/end_date: 日期范围查询
- status: 预约状态 (draft/pending/confirmed/in_progress/completed/cancelled)
- project_id: 项目ID筛选
- worker_id: 作业人员ID筛选`,
		Method:  "GET",
		Path:    "/api/appointments",
		QueryParams: map[string]string{
			"date":       "date",
			"start_date": "start_date",
			"end_date":   "end_date",
			"status":     "status",
			"project_id": "project_id",
			"worker_id":  "worker_id",
			"page":       "page",
			"page_size":  "page_size",
		},
		RequiredPermission: "appointment_view",
	},

	"query_appointment_detail": {
		Name:        "query_appointment_detail",
		Description: `查询单个预约/任务的详细信息。

当用户询问以下内容时使用：
- "预约详情"、"任务详情"
- "编号XXX的预约"
- "这个任务是什么内容"`,
		Method:            "GET",
		Path:              "/api/appointments/:id",
		PathParams:        []string{"id"},
		RequiredPermission: "appointment_view",
	},

	"query_stock": {
		Name:        "query_stock",
		Description: `查询库存信息，包括库存数量、安全库存、低库存预警。

当用户询问以下内容时使用：
- "库存有多少"、"库存情况"
- "物资库存"、"材料库存"
- "仓库里还有什么"`,
		Method:  "GET",
		Path:    "/api/stock/stocks",
		QueryParams: map[string]string{
			"material_id":  "material_id",
			"project_id":   "project_id",
			"warehouse":    "warehouse",
			"low_stock":    "low_stock",
			"page":         "page",
			"page_size":    "page_size",
		},
		RequiredPermission: "stock_view",
	},

	"query_stock_alerts": {
		Name:        "query_stock_alerts",
		Description: `查询库存预警信息，返回低于安全库存的物资列表。

当用户询问以下内容时使用：
- "库存预警"、"低库存"
- "库存不够"、"库存不足"
- "哪些物资需要补货"`,
		Method:            "GET",
		Path:              "/api/stock/stocks/alerts",
		RequiredPermission: "stock_view",
	},

	"query_materials": {
		Name:        "query_materials",
		Description: `查询物资主数据信息（物资档案）。

当用户询问以下内容时使用：
- "有哪些物资"、"物资信息"
- "物资名称"、"物资规格"
- "水泥/钢筋/电缆等具体物资"`,
		Method:  "GET",
		Path:    "/api/material-master",
		QueryParams: map[string]string{
			"search":    "search",
			"category":  "category",
			"page":      "page",
			"page_size": "page_size",
		},
		RequiredPermission: "material_view",
	},

	"query_material_plans": {
		Name:        "query_material_plans",
		Description: `查询物资计划列表。

当用户询问以下内容时使用：
- "物资计划"、"采购计划"
- "计划列表"、"需求计划"`,
		Method:  "GET",
		Path:    "/api/material-plan/plans",
		QueryParams: map[string]string{
			"project_id": "project_id",
			"status":     "status",
			"page":       "page",
			"page_size":  "page_size",
		},
		RequiredPermission: "material_plan_view",
	},

	"query_projects": {
		Name:        "query_projects",
		Description: `查询项目列表。

当用户询问以下内容时使用：
- "有哪些项目"、"项目列表"
- "项目信息"、"工程项目"`,
		Method:  "GET",
		Path:    "/api/project/projects",
		QueryParams: map[string]string{
			"search":    "search",
			"status":    "status",
			"page":      "page",
			"page_size": "page_size",
		},
		RequiredPermission: "project_view",
	},

	"query_requisitions": {
		Name:        "query_requisitions",
		Description: `查询领用单列表。

当用户询问以下内容时使用：
- "领用单"、"领料单"
- "出库记录"、"物资领用"`,
		Method:  "GET",
		Path:    "/api/requisitions",
		QueryParams: map[string]string{
			"project_id": "project_id",
			"status":     "status",
			"page":       "page",
			"page_size":  "page_size",
		},
		RequiredPermission: "requisition_view",
	},

	"query_inbounds": {
		Name:        "query_inbounds",
		Description: `查询入库单列表。

当用户询问以下内容时使用：
- "入库单"、"入库记录"
- "物资入库"、"进货记录"`,
		Method:  "GET",
		Path:    "/api/inbound/inbound-orders",
		QueryParams: map[string]string{
			"project_id": "project_id",
			"status":     "status",
			"page":       "page",
			"page_size":  "page_size",
		},
		RequiredPermission: "inbound_view",
	},

	"query_pending_approvals": {
		Name:        "query_pending_approvals",
		Description: `查询待审批的预约/任务列表。

当用户询问以下内容时使用：
- "待审批"、"审批任务"
- "需要我审批的"、"待处理审批"`,
		Method:  "GET",
		Path:    "/api/appointments/pending",
		QueryParams: map[string]string{
			"page":      "page",
			"page_size": "page_size",
		},
		RequiredPermission: "appointment_approve",
	},

	"query_workers": {
		Name:        "query_workers",
		Description: `查询作业人员列表。

当用户询问以下内容时使用：
- "作业人员"、"工人列表"
- "有哪些工人"、"人员名单"
- "谁可以安排"、"空闲人员"`,
		Method:            "GET",
		Path:              "/api/appointments/workers",
		RequiredPermission: "appointment_view",
	},

	"query_worker_calendar": {
		Name:        "query_worker_calendar",
		Description: `查询作业人员的日历/排班情况。

当用户询问以下内容时使用：
- "人员日历"、"排班情况"
- "谁有空"、"空闲时间"
- "某人的安排"`,
		Method:  "GET",
		Path:    "/api/appointments/calendar/view",
		QueryParams: map[string]string{
			"worker_id":  "worker_id",
			"start_date": "start_date",
			"end_date":   "end_date",
		},
		RequiredPermission: "appointment_view",
	},

	// ==================== 操作类工具 ====================
	"create_appointment": {
		Name:        "create_appointment",
		Description: `创建新的施工预约/任务。

当用户询问以下内容时使用：
- "创建预约"、"新建任务"
- "安排施工"、"添加预约"
- "帮我预约"、"创建一个任务"

必填参数：
- project_id: 项目ID
- work_date: 作业日期 (YYYY-MM-DD)
- time_slot: 时间段 (morning/noon/afternoon/fullday)
- work_content: 作业内容

可选参数：
- work_location: 作业地点
- contact_person: 联系人
- contact_phone: 联系电话
- is_urgent: 是否加急
- priority: 优先级(1-10)`,
		Method:       "POST",
		Path:         "/api/appointments",
		RequiresBody: true,
		BodyParams: []string{
			"project_id", "work_date", "time_slot", "work_content",
			"work_location", "work_type", "contact_person", "contact_phone",
			"is_urgent", "priority", "urgent_reason",
			"assigned_worker_id", "assigned_worker_ids", "assigned_worker_names",
		},
		RequiredPermission: "appointment_create",
	},

	"update_appointment": {
		Name:        "update_appointment",
		Description: `更新/修改已有的预约/任务。

当用户询问以下内容时使用：
- "修改预约"、"更新任务"
- "改一下时间"、"换个地点"`,
		Method:            "PUT",
		Path:              "/api/appointments/:id",
		PathParams:        []string{"id"},
		RequiresBody:      true,
		BodyParams: []string{
			"work_date", "time_slot", "work_content", "work_location",
			"contact_person", "contact_phone", "is_urgent", "priority",
		},
		RequiredPermission: "appointment_edit",
	},

	"submit_appointment": {
		Name:        "submit_appointment",
		Description: `提交预约/任务进入审批流程。

当用户询问以下内容时使用：
- "提交审批"、"提交预约"
- "发送审批"、"申请审批"`,
		Method:            "POST",
		Path:              "/api/appointments/:id/submit",
		PathParams:        []string{"id"},
		RequiredPermission: "appointment_edit",
	},

	"approve_appointment": {
		Name:        "approve_appointment",
		Description: `审批预约/任务。

当用户询问以下内容时使用：
- "审批通过"、"同意审批"
- "驳回"、"拒绝审批"
- "批准这个预约"`,
		Method:            "POST",
		Path:              "/api/appointments/:id/approve",
		PathParams:        []string{"id"},
		RequiresBody:      true,
		BodyParams:        []string{"action", "comment"},
		RequiredPermission: "appointment_approve",
	},

	"cancel_appointment": {
		Name:        "cancel_appointment",
		Description: `取消预约/任务。

当用户询问以下内容时使用：
- "取消预约"、"取消任务"
- "不要这个了"、"撤回预约"`,
		Method:            "POST",
		Path:              "/api/appointments/:id/cancel",
		PathParams:        []string{"id"},
		RequiresBody:      true,
		BodyParams:        []string{"reason"},
		RequiredPermission: "appointment_edit",
	},

	"stock_in": {
		Name:        "stock_in",
		Description: `物资入库操作。

当用户询问以下内容时使用：
- "入库"、"物资入库"
- "添加库存"、"进货"`,
		Method:            "POST",
		Path:              "/api/inbound",
		RequiresBody:      true,
		BodyParams: []string{
			"project_id", "supplier_id", "items",
			"warehouse", "remark",
		},
		RequiredPermission: "inbound_create",
	},

	"stock_out": {
		Name:        "stock_out",
		Description: `物资出库/领用操作。

当用户询问以下内容时使用：
- "出库"、"领料"
- "物资领用"、"取物资"`,
		Method:            "POST",
		Path:              "/api/requisition",
		RequiresBody:      true,
		BodyParams: []string{
			"project_id", "items", "applicant_id",
			"remark",
		},
		RequiredPermission: "requisition_create",
	},

	"create_material_plan": {
		Name:        "create_material_plan",
		Description: `创建物资计划。

当用户询问以下内容时使用：
- "创建物资计划"、"新建采购计划"
- "添加需求计划"`,
		Method:            "POST",
		Path:              "/api/material-plan/plans",
		RequiresBody:      true,
		BodyParams: []string{
			"project_id", "items", "remark",
		},
		RequiredPermission: "material_plan_create",
	},

	// ==================== 施工日志 ====================
	"query_construction_logs": {
		Name:        "query_construction_logs",
		Description: `查询施工日志。

当用户询问以下内容时使用：
- "施工日志"、"施工记录"
- "今天的日志"、"某项目的日志"`,
		Method:  "GET",
		Path:    "/api/construction-logs",
		QueryParams: map[string]string{
			"project_id": "project_id",
			"log_date":   "log_date",
			"start_date": "start_date",
			"end_date":   "end_date",
			"page":       "page",
			"page_size":  "page_size",
		},
		RequiredPermission: "constructionlog_view",
	},

	"query_construction_log_detail": {
		Name:        "query_construction_log_detail",
		Description: `查询单条施工日志的详细信息。

当用户询问以下内容时使用：
- "日志详情"、"施工日志详情"
- "查看这条日志的内容"`,
		Method:            "GET",
		Path:              "/api/construction-logs/:id",
		PathParams:        []string{"id"},
		RequiredPermission: "constructionlog_view",
	},

	// ==================== 工作流管理 ====================
	"query_pending_workflow_tasks": {
		Name:        "query_pending_workflow_tasks",
		Description: `查询待办工作流任务列表。

当用户询问以下内容时使用：
- "待办任务"、"我的待办"
- "需要处理的工作流"
- "待审批的工作流"`,
		Method:  "GET",
		Path:    "/api/workflow/tasks/pending",
		QueryParams: map[string]string{
			"task_type": "task_type",
			"page":      "page",
			"page_size": "page_size",
		},
		RequiredPermission: "appointment_approve",
	},

	"approve_workflow_task": {
		Name:        "approve_workflow_task",
		Description: `审批通过工作流任务。

当用户询问以下内容时使用：
- "审批通过这个任务"
- "同意这个工作流"`,
		Method:            "POST",
		Path:              "/api/workflow/tasks/:id/approve",
		PathParams:        []string{"id"},
		RequiresBody:      true,
		BodyParams:        []string{"comment"},
		RequiredPermission: "appointment_approve",
	},

	"reject_workflow_task": {
		Name:        "reject_workflow_task",
		Description: `驳回工作流任务。

当用户询问以下内容时使用：
- "驳回这个任务"
- "拒绝这个工作流"`,
		Method:            "POST",
		Path:              "/api/workflow/tasks/:id/reject",
		PathParams:        []string{"id"},
		RequiresBody:      true,
		BodyParams:        []string{"reason"},
		RequiredPermission: "appointment_approve",
	},

	// ==================== 通知管理 ====================
	"query_notifications": {
		Name:        "query_notifications",
		Description: `查询通知消息列表。

当用户询问以下内容时使用：
- "通知列表"、"消息列表"
- "我的消息"、"通知"`,
		Method:  "GET",
		Path:    "/api/notifications",
		QueryParams: map[string]string{
			"unread_only":       "unread_only",
			"notification_type": "notification_type",
			"page":              "page",
			"page_size":         "page_size",
		},
		RequiredPermission: "notification_view",
	},

	"query_unread_notification_count": {
		Name:        "query_unread_notification_count",
		Description: `查询未读通知数量。

当用户询问以下内容时使用：
- "未读消息数量"、"有多少未读"
- "未读通知"、"消息提醒"`,
		Method:            "GET",
		Path:              "/api/notifications/unread-count",
		RequiredPermission: "notification_view",
	},

	// ==================== 考勤管理 ====================
	"query_attendance": {
		Name:        "query_attendance",
		Description: `查询考勤打卡记录。

当用户询问以下内容时使用：
- "打卡记录"、"考勤记录"
- "今天的考勤"、"某人出勤情况"`,
		Method:  "GET",
		Path:    "/api/attendance/records",
		QueryParams: map[string]string{
			"user_id":         "user_id",
			"date":            "date",
			"start_date":      "start_date",
			"end_date":        "end_date",
			"attendance_type": "attendance_type",
			"status":          "status",
			"page":            "page",
			"page_size":       "page_size",
		},
		RequiredPermission: "attendance_view",
	},

	"query_attendance_stats": {
		Name:        "query_attendance_stats",
		Description: `查询考勤统计数据。

当用户询问以下内容时使用：
- "考勤统计"、"出勤率"
- "本月考勤情况"、"月度考勤"`,
		Method:  "GET",
		Path:    "/api/attendance/statistics",
		QueryParams: map[string]string{
			"user_id": "user_id",
			"year":    "year",
			"month":   "month",
		},
		RequiredPermission: "attendance_view",
	},
}

// APIToolExecutor 通过内部HTTP调用系统API执行工具
type APIToolExecutor struct {
	db          *gorm.DB
	baseURL     string // API基础URL
	userID      uint   // 当前用户ID
	userName    string // 当前用户名
	httpClient  *http.Client
}

// NewAPIToolExecutor 创建API工具执行器
func NewAPIToolExecutor(db *gorm.DB, userID uint, userName string) *APIToolExecutor {
	return &APIToolExecutor{
		db:       db,
		baseURL:  "http://127.0.0.1:8088", // 内部调用（后端服务端口）
		userID:   userID,
		userName: userName,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetBaseURL 设置API基础URL
func (e *APIToolExecutor) SetBaseURL(url string) {
	e.baseURL = url
}

// ExecuteToolCall 执行工具调用
func (e *APIToolExecutor) ExecuteToolCall(ctx context.Context, name string, args map[string]interface{}) (interface{}, error) {
	log.Printf("[APIToolExecutor] Executing tool: %s with args: %+v", name, args)

	// 获取工具定义
	toolDef, ok := ToolDefinitions[name]
	if !ok {
		return nil, fmt.Errorf("unknown tool: %s", name)
	}

	// 检查权限
	if toolDef.RequiredPermission != "" {
		if !e.hasPermission(toolDef.RequiredPermission) {
			log.Printf("[APIToolExecutor] Permission denied: user %d lacks permission %s", e.userID, toolDef.RequiredPermission)
			return map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("权限不足: 需要 %s 权限", toolDef.RequiredPermission),
			}, nil
		}
	}

	// 构建请求URL
	reqURL, err := e.buildRequestURL(toolDef, args)
	if err != nil {
		return nil, err
	}

	log.Printf("[APIToolExecutor] Request URL: %s %s", toolDef.Method, reqURL)

	// 构建请求体
	var bodyReader io.Reader
	if toolDef.RequiresBody {
		bodyData := e.buildRequestBody(toolDef, args)
		bodyBytes, err := json.Marshal(bodyData)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
		log.Printf("[APIToolExecutor] Request body: %s", string(bodyBytes))
	}

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, toolDef.Method, reqURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Internal-Call", "true") // 标记为内部调用
	req.Header.Set("X-User-ID", fmt.Sprintf("%d", e.userID))
	req.Header.Set("X-User-Name", e.userName)

	// 执行请求
	resp, err := e.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	log.Printf("[APIToolExecutor] Response status: %d, body length: %d", resp.StatusCode, len(respBody))

	// 检查HTTP状态码
	if resp.StatusCode >= 500 {
		// 服务器错误，返回错误信息
		return map[string]interface{}{
			"success": false,
			"error":   fmt.Sprintf("服务器错误 (HTTP %d)", resp.StatusCode),
			"status":  resp.StatusCode,
		}, nil
	}

	// 检查Content-Type是否为JSON
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		// 非JSON响应，可能是错误页面
		log.Printf("[APIToolExecutor] Non-JSON response, Content-Type: %s", contentType)
		return map[string]interface{}{
			"success": false,
			"error":   fmt.Sprintf("API返回非JSON响应 (Content-Type: %s)", contentType),
			"status":  resp.StatusCode,
		}, nil
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		// JSON解析失败，可能是格式问题
		log.Printf("[APIToolExecutor] Failed to parse JSON: %v, response: %s", err, string(respBody))
		return map[string]interface{}{
			"success": false,
			"error":   fmt.Sprintf("响应解析失败: %v", err),
		}, nil
	}

	// 检查API错误
	if resp.StatusCode >= 400 {
		errMsg := "API请求失败"
		if msg, ok := result["message"].(string); ok {
			errMsg = msg
		} else if msg, ok := result["error"].(string); ok {
			errMsg = msg
		}
		return map[string]interface{}{
			"success": false,
			"error":   errMsg,
			"status":  resp.StatusCode,
		}, nil
	}

	// 标准化返回结果
	return map[string]interface{}{
		"success": true,
		"data":    result,
	}, nil
}

// buildRequestURL 构建请求URL
func (e *APIToolExecutor) buildRequestURL(toolDef ToolDefinition, args map[string]interface{}) (string, error) {
	path := toolDef.Path

	// 替换路径参数
	for _, param := range toolDef.PathParams {
		value, ok := args[param]
		if !ok {
			return "", fmt.Errorf("missing path parameter: %s", param)
		}
		placeholder := ":" + param
		path = strings.Replace(path, placeholder, fmt.Sprintf("%v", value), 1)
	}

	// 构建查询参数
	if len(toolDef.QueryParams) > 0 {
		query := url.Values{}
		for argName, queryName := range toolDef.QueryParams {
			if value, ok := args[argName]; ok && value != nil && value != "" {
				query.Set(queryName, fmt.Sprintf("%v", value))
			}
		}
		if len(query) > 0 {
			path = path + "?" + query.Encode()
		}
	}

	return e.baseURL + path, nil
}

// buildRequestBody 构建请求体
func (e *APIToolExecutor) buildRequestBody(toolDef ToolDefinition, args map[string]interface{}) map[string]interface{} {
	body := make(map[string]interface{})

	for _, param := range toolDef.BodyParams {
		if value, ok := args[param]; ok {
			body[param] = value
		}
	}

	// 特殊处理：审批操作
	if toolDef.Name == "approve_appointment" {
		if action, ok := args["action"].(string); !ok || action == "" {
			body["action"] = "approve" // 默认通过
		}
	}

	return body
}

// GetToolList 获取工具列表（用于OpenAI Function Calling）
func GetToolList() []string {
	tools := make([]string, 0, len(ToolDefinitions))
	for name := range ToolDefinitions {
		tools = append(tools, name)
	}
	return tools
}

// GetToolDescription 获取工具描述
func GetToolDescription(name string) string {
	if tool, ok := ToolDefinitions[name]; ok {
		return tool.Description
	}
	return ""
}

// hasPermission 检查用户是否有指定权限
func (e *APIToolExecutor) hasPermission(permission string) bool {
	var user auth.User
	if err := e.db.Preload("Roles").First(&user, e.userID).Error; err != nil {
		log.Printf("[APIToolExecutor] Failed to load user %d: %v", e.userID, err)
		return false
	}

	// 管理员拥有所有权限
	if user.IsAdmin() {
		return true
	}

	// 检查用户角色是否包含所需权限
	return user.HasPermission(permission)
}
