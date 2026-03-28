package openai

import (
	"encoding/json"
	"fmt"
	"time"
)

// ToolBuilder helps build OpenAI function tools
type ToolBuilder struct {
	tools []Tool
}

// NewToolBuilder creates a new tool builder
func NewToolBuilder() *ToolBuilder {
	return &ToolBuilder{
		tools: make([]Tool, 0),
	}
}

// AddTool adds a function tool
func (tb *ToolBuilder) AddTool(name, description string, parameters map[string]interface{}) *ToolBuilder {
	tool := Tool{
		Type: "function",
	}
	tool.Function.Name = name
	tool.Function.Description = description
	tool.Function.Parameters = parameters
	tb.tools = append(tb.tools, tool)
	return tb
}

// Build returns the tools
func (tb *ToolBuilder) Build() []Tool {
	return tb.tools
}

// GetMaterialManagementTools returns all tools for the AI assistant
// 包含完整的系统能力工具定义
func GetMaterialManagementTools() []Tool {
	return NewToolBuilder().
		// ==================== 施工预约管理 ====================
		AddTool("query_appointments", `查询施工预约/施工任务/施工作业安排。

【使用场景】
当用户询问以下内容时使用此工具：
- "明天的任务"、"今日安排"、"后天有什么"
- "施工计划"、"预约列表"、"任务列表"
- "我的任务"、"待处理的预约"
- "本周有什么安排"

【参数说明】
- date: 指定查询日期，格式 YYYY-MM-DD
- start_date/end_date: 日期范围查询
- status: 预约状态筛选 (draft草稿/pending待审批/scheduled已排期/in_progress进行中/completed已完成/cancelled已取消)
- project_id: 按项目筛选
- worker_id: 按作业人员筛选

【返回信息】
返回预约列表，包含预约编号、日期、时间段、作业内容、地点、状态、作业人员等`, map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"date": map[string]interface{}{
					"type":        "string",
					"description": "查询日期，格式 YYYY-MM-DD，如 " + time.Now().Format("2006-01-02"),
				},
				"start_date": map[string]interface{}{
					"type":        "string",
					"description": "开始日期，用于日期范围查询",
				},
				"end_date": map[string]interface{}{
					"type":        "string",
					"description": "结束日期，用于日期范围查询",
				},
				"status": map[string]interface{}{
					"type":        "string",
					"description": "预约状态：draft(草稿), pending(待审批), scheduled(已排期), in_progress(进行中), completed(已完成), cancelled(已取消)",
					"enum":        []string{"draft", "pending", "scheduled", "in_progress", "completed", "cancelled"},
				},
				"project_id": map[string]interface{}{
					"type":        "integer",
					"description": "项目ID，用于筛选特定项目的预约",
				},
				"worker_id": map[string]interface{}{
					"type":        "integer",
					"description": "作业人员ID，用于筛选特定人员的预约",
				},
				"page": map[string]interface{}{
					"type":        "integer",
					"description": "页码，默认1",
					"default":     1,
				},
				"page_size": map[string]interface{}{
					"type":        "integer",
					"description": "每页数量，默认20",
					"default":     20,
				},
			},
		}).
		AddTool("query_appointment_detail", `查询单个预约/任务的详细信息。

【使用场景】
当用户询问以下内容时使用：
- "预约详情"、"任务详情"
- "编号XXX的预约是什么"
- "这个任务的具体内容"

【参数说明】
- id: 预约单ID（必填）

【返回信息】
返回预约的完整信息，包括申请人、作业人员、审批状态等详细内容`, map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type":        "integer",
					"description": "预约单ID（必填）",
				},
			},
			"required": []string{"id"},
		}).
		AddTool("create_appointment", `创建新的施工预约/任务。

【使用场景】
当用户询问以下内容时使用：
- "创建预约"、"新建任务"
- "安排施工"、"添加预约"
- "帮我预约明天"

【必填参数】
- project_id: 项目ID
- work_date: 作业日期 (YYYY-MM-DD)
- time_slot: 时间段
- work_content: 作业内容描述

【可选参数】
- work_location: 作业地点
- work_type: 作业类型 (construction施工/inspection检查/maintenance维护/other其他)
- contact_person: 联系人
- contact_phone: 联系电话
- is_urgent: 是否加急
- priority: 优先级(1-10，数字越大优先级越高)
- assigned_worker_id: 主作业人员ID
- assigned_worker_ids: 作业人员ID数组
- assigned_worker_names: 作业人员姓名（逗号分隔）

【重要提示】
- 指定作业人员后，会**立即锁定**该人员的时间段
- 如果作业人员在指定时间段不可用，会返回错误
- 锁定后其他人无法再预约该人员同一时间段`, map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"project_id": map[string]interface{}{
					"type":        "integer",
					"description": "项目ID（必填）",
				},
				"work_date": map[string]interface{}{
					"type":        "string",
					"description": "作业日期，格式 YYYY-MM-DD（必填）",
				},
				"time_slot": map[string]interface{}{
					"type":        "string",
					"description": "时间段：morning(上午08:00-12:00), noon(中午12:00-14:00), afternoon(下午14:00-18:00), fullday(全天)",
					"enum":        []string{"morning", "noon", "afternoon", "fullday"},
				},
				"work_content": map[string]interface{}{
					"type":        "string",
					"description": "作业内容描述（必填）",
				},
				"work_location": map[string]interface{}{
					"type":        "string",
					"description": "作业地点",
				},
				"work_type": map[string]interface{}{
					"type":        "string",
					"description": "作业类型：construction(施工), inspection(检查), maintenance(维护), other(其他)",
				},
				"contact_person": map[string]interface{}{
					"type":        "string",
					"description": "联系人姓名",
				},
				"contact_phone": map[string]interface{}{
					"type":        "string",
					"description": "联系电话",
				},
				"is_urgent": map[string]interface{}{
					"type":        "boolean",
					"description": "是否加急",
					"default":     false,
				},
				"priority": map[string]interface{}{
					"type":        "integer",
					"description": "优先级(1-10)，数字越大优先级越高",
					"default":     5,
				},
				"urgent_reason": map[string]interface{}{
					"type":        "string",
					"description": "加急原因（高优先级加急时必填）",
				},
				"assigned_worker_id": map[string]interface{}{
					"type":        "integer",
					"description": "主作业人员ID",
				},
				"assigned_worker_ids": map[string]interface{}{
					"type":        "array",
					"description": "作业人员ID数组",
					"items": map[string]interface{}{
						"type": "integer",
					},
				},
				"assigned_worker_names": map[string]interface{}{
					"type":        "string",
					"description": "作业人员姓名，多个用逗号分隔",
				},
			},
			"required": []string{"project_id", "work_date", "time_slot", "work_content"},
		}).
		AddTool("update_appointment", `更新/修改已有的预约/任务。

【使用场景】
当用户询问以下内容时使用：
- "修改预约"、"更新任务"
- "改一下时间"、"换个地点"

【参数说明】
- id: 预约单ID（必填）
- 其他字段：需要更新的字段`, map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type":        "integer",
					"description": "预约单ID（必填）",
				},
				"work_date": map[string]interface{}{
					"type":        "string",
					"description": "作业日期，格式 YYYY-MM-DD",
				},
				"time_slot": map[string]interface{}{
					"type":        "string",
					"description": "时间段",
				},
				"work_content": map[string]interface{}{
					"type":        "string",
					"description": "作业内容",
				},
				"work_location": map[string]interface{}{
					"type":        "string",
					"description": "作业地点",
				},
				"contact_person": map[string]interface{}{
					"type":        "string",
					"description": "联系人",
				},
				"contact_phone": map[string]interface{}{
					"type":        "string",
					"description": "联系电话",
				},
			},
			"required": []string{"id"},
		}).
		AddTool("submit_appointment", `提交预约/任务进入审批流程。

【使用场景】
当用户询问以下内容时使用：
- "提交审批"、"提交预约"
- "发送审批"、"申请审批"

【参数说明】
- id: 预约单ID（必填）

【返回信息】
返回提交后的预约状态`, map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type":        "integer",
					"description": "预约单ID（必填）",
				},
			},
			"required": []string{"id"},
		}).
		AddTool("approve_appointment", `审批预约/任务。

【使用场景】
当用户询问以下内容时使用：
- "审批通过"、"同意审批"
- "驳回"、"拒绝审批"
- "批准这个预约"

【参数说明】
- id: 预约单ID（必填）
- action: 审批动作 approve(通过)/reject(驳回)（必填）
- comment: 审批意见`, map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type":        "integer",
					"description": "预约单ID（必填）",
				},
				"action": map[string]interface{}{
					"type":        "string",
					"description": "审批动作：approve(通过) 或 reject(驳回)",
					"enum":        []string{"approve", "reject"},
				},
				"comment": map[string]interface{}{
					"type":        "string",
					"description": "审批意见",
				},
			},
			"required": []string{"id", "action"},
		}).
		AddTool("cancel_appointment", `取消预约/任务。

【使用场景】
当用户询问以下内容时使用：
- "取消预约"、"取消任务"
- "不要这个了"、"撤回预约"

【参数说明】
- id: 预约单ID（必填）
- reason: 取消原因

【重要提示】
- 取消预约会**释放**已锁定作业人员的日历，使其可以被其他预约使用`, map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type":        "integer",
					"description": "预约单ID（必填）",
				},
				"reason": map[string]interface{}{
					"type":        "string",
					"description": "取消原因",
				},
			},
			"required": []string{"id"},
		}).
		AddTool("query_pending_approvals", `查询待审批的预约/任务列表。

【使用场景】
当用户询问以下内容时使用：
- "待审批"、"审批任务"
- "需要我审批的"、"待处理审批"
- "有哪些待审批的预约"

【返回信息】
返回待当前用户审批的预约列表`, map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"page": map[string]interface{}{
					"type":        "integer",
					"description": "页码",
					"default":     1,
				},
				"page_size": map[string]interface{}{
					"type":        "integer",
					"description": "每页数量",
					"default":     20,
				},
			},
		}).
		AddTool("query_workers", `查询作业人员列表。

【使用场景】
当用户询问以下内容时使用：
- "作业人员"、"工人列表"
- "有哪些工人"、"人员名单"
- "谁可以安排"、"能派谁去"

【返回信息】
返回所有作业人员信息

【重要提示】
- 此工具返回所有作业人员，不包含日历锁定状态
- 要查看某人员某天是否有空，请使用 query_worker_calendar 工具
- 创建预约时指定人员会**立即锁定**其日历`, map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{},
		}).
		AddTool("query_worker_calendar", `查询作业人员的日历/排班情况。

【使用场景】
当用户询问以下内容时使用：
- "人员日历"、"排班情况"
- "谁有空"、"空闲人员"
- "某人明天的安排"

【参数说明】
- worker_id: 作业人员ID（可选，不填则查所有）
- start_date: 开始日期（必填）
- end_date: 结束日期（必填）

【返回信息】
按日期分组返回人员的预约安排`, map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"worker_id": map[string]interface{}{
					"type":        "integer",
					"description": "作业人员ID，不填则查询所有人员",
				},
				"start_date": map[string]interface{}{
					"type":        "string",
					"description": "开始日期，格式 YYYY-MM-DD（必填）",
				},
				"end_date": map[string]interface{}{
					"type":        "string",
					"description": "结束日期，格式 YYYY-MM-DD（必填）",
				},
			},
			"required": []string{"start_date", "end_date"},
		}).
		// ==================== 库存管理 ====================
		AddTool("query_stock", `查询库存信息，包括库存数量、安全库存、低库存预警。

【使用场景】
当用户询问以下内容时使用：
- "库存有多少"、"库存情况"
- "物资库存"、"材料库存"
- "仓库里还有什么"

【参数说明】
- material_id: 物资ID筛选
- project_id: 项目ID筛选
- warehouse: 仓库筛选
- low_stock: 是否只显示低库存

【返回信息】
返回库存列表，包含物资名称、数量、安全库存、仓库等`, map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"material_id": map[string]interface{}{
					"type":        "integer",
					"description": "物资ID筛选",
				},
				"project_id": map[string]interface{}{
					"type":        "integer",
					"description": "项目ID筛选",
				},
				"warehouse": map[string]interface{}{
					"type":        "string",
					"description": "仓库名称筛选",
				},
				"low_stock": map[string]interface{}{
					"type":        "boolean",
					"description": "是否只显示低库存物资",
					"default":     false,
				},
				"page": map[string]interface{}{
					"type":        "integer",
					"description": "页码",
					"default":     1,
				},
				"page_size": map[string]interface{}{
					"type":        "integer",
					"description": "每页数量",
					"default":     20,
				},
			},
		}).
		AddTool("query_stock_alerts", `查询库存预警信息，返回低于安全库存的物资列表。

【使用场景】
当用户询问以下内容时使用：
- "库存预警"、"低库存"
- "库存不够"、"库存不足"
- "哪些物资需要补货"

【返回信息】
返回所有低于安全库存的物资列表`, map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{},
		}).
		// ==================== 物资管理 ====================
		AddTool("query_materials", `查询物资主数据信息（物资档案）。

【使用场景】
当用户询问以下内容时使用：
- "有哪些物资"、"物资信息"
- "物资名称"、"物资规格"
- "水泥/钢筋/电缆等具体物资"

【参数说明】
- search: 搜索关键词（物资名称、规格等）
- category: 物资分类筛选

【返回信息】
返回物资列表，包含名称、规格、单位、分类等`, map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"search": map[string]interface{}{
					"type":        "string",
					"description": "搜索关键词，可以是物资名称、规格等",
				},
				"category": map[string]interface{}{
					"type":        "string",
					"description": "物资分类",
				},
				"page": map[string]interface{}{
					"type":        "integer",
					"description": "页码",
					"default":     1,
				},
				"page_size": map[string]interface{}{
					"type":        "integer",
					"description": "每页数量",
					"default":     20,
				},
			},
		}).
		AddTool("query_material_plans", `查询物资计划列表。

【使用场景】
当用户询问以下内容时使用：
- "物资计划"、"采购计划"
- "计划列表"、"需求计划"

【参数说明】
- project_id: 项目ID筛选
- status: 计划状态筛选

【返回信息】
返回物资计划列表`, map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"project_id": map[string]interface{}{
					"type":        "integer",
					"description": "项目ID",
				},
				"status": map[string]interface{}{
					"type":        "string",
					"description": "计划状态：draft(草稿)/pending(待审批)/approved(已批准)/completed(已完成)",
				},
				"page": map[string]interface{}{
					"type":        "integer",
					"description": "页码",
					"default":     1,
				},
				"page_size": map[string]interface{}{
					"type":        "integer",
					"description": "每页数量",
					"default":     20,
				},
			},
		}).
		AddTool("create_material_plan", `创建物资计划。

【使用场景】
当用户询问以下内容时使用：
- "创建物资计划"、"新建采购计划"
- "添加需求计划"

【参数说明】
- project_id: 项目ID（必填）
- items: 计划明细数组（必填）
- remark: 备注

【items格式】
每个item包含：material_id(物资ID), quantity(数量), unit_price(单价), planned_arrival_date(计划到货日期)`, map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"project_id": map[string]interface{}{
					"type":        "integer",
					"description": "项目ID（必填）",
				},
				"items": map[string]interface{}{
					"type":        "array",
					"description": "计划明细列表",
					"items": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"material_id":          map[string]interface{}{"type": "integer", "description": "物资ID"},
							"quantity":             map[string]interface{}{"type": "number", "description": "数量"},
							"unit_price":           map[string]interface{}{"type": "number", "description": "单价"},
							"planned_arrival_date": map[string]interface{}{"type": "string", "description": "计划到货日期"},
						},
					},
				},
				"remark": map[string]interface{}{
					"type":        "string",
					"description": "备注",
				},
			},
			"required": []string{"project_id", "items"},
		}).
		// ==================== 项目管理 ====================
		AddTool("query_projects", `查询项目列表。

【使用场景】
当用户询问以下内容时使用：
- "有哪些项目"、"项目列表"
- "项目信息"、"工程项目"

【参数说明】
- search: 搜索关键词
- status: 项目状态筛选

【返回信息】
返回项目列表，包含项目名称、编号、状态等`, map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"search": map[string]interface{}{
					"type":        "string",
					"description": "搜索关键词",
				},
				"status": map[string]interface{}{
					"type":        "string",
					"description": "项目状态：planning(规划中)/active(进行中)/suspended(暂停)/completed(已完成)",
				},
				"page": map[string]interface{}{
					"type":        "integer",
					"description": "页码",
					"default":     1,
				},
				"page_size": map[string]interface{}{
					"type":        "integer",
					"description": "每页数量",
					"default":     20,
				},
			},
		}).
		// ==================== 领用/入库管理 ====================
		AddTool("query_requisitions", `查询领用单列表。

【使用场景】
当用户询问以下内容时使用：
- "领用单"、"领料单"
- "出库记录"、"物资领用"

【返回信息】
返回领用单列表`, map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"project_id": map[string]interface{}{
					"type":        "integer",
					"description": "项目ID",
				},
				"status": map[string]interface{}{
					"type":        "string",
					"description": "状态筛选",
				},
				"page": map[string]interface{}{
					"type":        "integer",
					"description": "页码",
					"default":     1,
				},
				"page_size": map[string]interface{}{
					"type":        "integer",
					"description": "每页数量",
					"default":     20,
				},
			},
		}).
		AddTool("query_inbounds", `查询入库单列表。

【使用场景】
当用户询问以下内容时使用：
- "入库单"、"入库记录"
- "物资入库"、"进货记录"

【返回信息】
返回入库单列表`, map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"project_id": map[string]interface{}{
					"type":        "integer",
					"description": "项目ID",
				},
				"status": map[string]interface{}{
					"type":        "string",
					"description": "状态筛选",
				},
				"page": map[string]interface{}{
					"type":        "integer",
					"description": "页码",
					"default":     1,
				},
				"page_size": map[string]interface{}{
					"type":        "integer",
					"description": "每页数量",
					"default":     20,
				},
			},
		}).
		// ==================== 考勤管理 ====================
		AddTool("query_attendance", `查询考勤打卡记录。

【使用场景】
当用户询问以下内容时使用此工具：
- "打卡记录"、"考勤记录"
- "今天的考勤"、"某人出勤情况"
- "出勤记录"、"签到记录"

【参数说明】
- user_id: 用户ID（可选）
- date: 日期（可选，格式 YYYY-MM-DD）
- start_date/end_date: 日期范围查询
- attendance_type: 打卡类型 (morning上午/afternoon下午/noon_overtime中午加班/night_overtime晚上加班)
- status: 打卡状态 (pending待确认/confirmed已确认/rejected已驳回)

【返回信息】
返回打卡记录列表，包含用户名、打卡时间、打卡类型、状态等`, map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"user_id": map[string]interface{}{
					"type":        "integer",
					"description": "用户ID，查询特定用户的打卡记录",
				},
				"date": map[string]interface{}{
					"type":        "string",
					"description": "查询日期，格式 YYYY-MM-DD，如 " + time.Now().Format("2006-01-02"),
				},
				"start_date": map[string]interface{}{
					"type":        "string",
					"description": "开始日期，用于日期范围查询",
				},
				"end_date": map[string]interface{}{
					"type":        "string",
					"description": "结束日期，用于日期范围查询",
				},
				"attendance_type": map[string]interface{}{
					"type":        "string",
					"description": "打卡类型：morning(上午)/afternoon(下午)/noon_overtime(中午加班)/night_overtime(晚上加班)",
					"enum":        []string{"morning", "afternoon", "noon_overtime", "night_overtime"},
				},
				"status": map[string]interface{}{
					"type":        "string",
					"description": "打卡状态：pending(待确认)/confirmed(已确认)/rejected(已驳回)",
					"enum":        []string{"pending", "confirmed", "rejected"},
				},
				"page": map[string]interface{}{
					"type":        "integer",
					"description": "页码，默认1",
					"default":     1,
				},
				"page_size": map[string]interface{}{
					"type":        "integer",
					"description": "每页数量，默认20",
					"default":     20,
				},
			},
		}).
		AddTool("query_attendance_stats", `查询考勤统计数据。

【使用场景】
当用户询问以下内容时使用此工具：
- "考勤统计"、"出勤率"
- "本月考勤情况"、"月度考勤"
- "加班统计"、"加班小时"

【参数说明】
- user_id: 用户ID（可选，不填则查询全部）
- year: 年份（必填）
- month: 月份（必填，1-12）

【返回信息】
返回月度考勤汇总，包含工作天数、加班时长等`, map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"user_id": map[string]interface{}{
					"type":        "integer",
					"description": "用户ID，不填则查询全部用户",
				},
				"year": map[string]interface{}{
					"type":        "integer",
					"description": "年份，如 2026",
				},
				"month": map[string]interface{}{
					"type":        "integer",
					"description": "月份，1-12",
				},
			},
			"required": []string{"year", "month"},
		}).
		// ==================== 施工日志 ====================
		AddTool("query_construction_logs", `查询施工日志。

【使用场景】
当用户询问以下内容时使用此工具：
- "施工日志"、"施工记录"
- "今天的日志"、"某项目的日志"
- "工程日志"、"施工日记"

【参数说明】
- project_id: 项目ID（可选）
- log_date: 日志日期（可选，格式 YYYY-MM-DD）
- start_date/end_date: 日期范围查询

【返回信息】
返回施工日志列表，包含标题、内容、天气、进度等`, map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"project_id": map[string]interface{}{
					"type":        "integer",
					"description": "项目ID，筛选特定项目的日志",
				},
				"log_date": map[string]interface{}{
					"type":        "string",
					"description": "日志日期，格式 YYYY-MM-DD",
				},
				"start_date": map[string]interface{}{
					"type":        "string",
					"description": "开始日期，用于日期范围查询",
				},
				"end_date": map[string]interface{}{
					"type":        "string",
					"description": "结束日期，用于日期范围查询",
				},
				"page": map[string]interface{}{
					"type":        "integer",
					"description": "页码，默认1",
					"default":     1,
				},
				"page_size": map[string]interface{}{
					"type":        "integer",
					"description": "每页数量，默认20",
					"default":     20,
				},
			},
		}).
		AddTool("query_construction_log_detail", `查询单条施工日志的详细信息。

【使用场景】
当用户询问以下内容时使用此工具：
- "日志详情"、"施工日志详情"
- "查看这条日志的内容"
- "展开这个日志"

【参数说明】
- id: 日志ID（必填）

【返回信息】
返回日志的完整信息，包括详细内容、图片、问题记录等`, map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"id": map[string]interface{}{
					"type":        "integer",
					"description": "日志ID（必填）",
				},
			},
			"required": []string{"id"},
		}).
		// ==================== 工作流管理 ====================
		AddTool("query_pending_workflow_tasks", `查询待办工作流任务列表。

【使用场景】
当用户询问以下内容时使用此工具：
- "待办任务"、"我的待办"
- "需要处理的工作流"
- "待审批的工作流"
- "待办事项"

【参数说明】
- task_type: 任务类型筛选（可选）
- page/page_size: 分页参数

【返回信息】
返回当前用户待处理的工作流任务列表`, map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"task_type": map[string]interface{}{
					"type":        "string",
					"description": "任务类型筛选，如：appointment_approval(预约审批)/material_plan_approval(物资计划审批)/requisition_approval(领用审批)",
				},
				"page": map[string]interface{}{
					"type":        "integer",
					"description": "页码，默认1",
					"default":     1,
				},
				"page_size": map[string]interface{}{
					"type":        "integer",
					"description": "每页数量，默认20",
					"default":     20,
				},
			},
		}).
		AddTool("approve_workflow_task", `审批通过工作流任务。

【使用场景】
当用户询问以下内容时使用此工具：
- "审批通过这个任务"
- "同意这个工作流"
- "批准待办"

【参数说明】
- task_id: 任务ID（必填）
- comment: 审批意见（可选）

【返回信息】
返回审批结果`, map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"task_id": map[string]interface{}{
					"type":        "integer",
					"description": "任务ID（必填）",
				},
				"comment": map[string]interface{}{
					"type":        "string",
					"description": "审批意见",
				},
			},
			"required": []string{"task_id"},
		}).
		AddTool("reject_workflow_task", `驳回工作流任务。

【使用场景】
当用户询问以下内容时使用此工具：
- "驳回这个任务"
- "拒绝这个工作流"
- "不同意待办"

【参数说明】
- task_id: 任务ID（必填）
- reason: 驳回原因（必填）

【返回信息】
返回驳回结果`, map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"task_id": map[string]interface{}{
					"type":        "integer",
					"description": "任务ID（必填）",
				},
				"reason": map[string]interface{}{
					"type":        "string",
					"description": "驳回原因（必填）",
				},
			},
			"required": []string{"task_id", "reason"},
		}).
		// ==================== 通知管理 ====================
		AddTool("query_notifications", `查询通知消息列表。

【使用场景】
当用户询问以下内容时使用此工具：
- "通知列表"、"消息列表"
- "我的消息"、"通知"
- "系统通知"

【参数说明】
- unread_only: 是否只显示未读（可选，默认false）
- notification_type: 通知类型筛选（可选）
- page/page_size: 分页参数

【返回信息】
返回通知列表，包含标题、内容、类型、是否已读等`, map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"unread_only": map[string]interface{}{
					"type":        "boolean",
					"description": "是否只显示未读通知",
					"default":     false,
				},
				"notification_type": map[string]interface{}{
					"type":        "string",
					"description": "通知类型：system(系统)/approval(审批)/reminder(提醒)/announcement(公告)",
				},
				"page": map[string]interface{}{
					"type":        "integer",
					"description": "页码，默认1",
					"default":     1,
				},
				"page_size": map[string]interface{}{
					"type":        "integer",
					"description": "每页数量，默认20",
					"default":     20,
				},
			},
		}).
		AddTool("query_unread_notification_count", `查询未读通知数量。

【使用场景】
当用户询问以下内容时使用此工具：
- "未读消息数量"、"有多少未读"
- "未读通知"、"消息提醒"

【返回信息】
返回未读通知总数`, map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{},
		}).
		// ==================== 通用分析工具 ====================
		AddTool("generate_report", `生成数据报告。

【使用场景】
当用户询问以下内容时使用：
- "生成报告"、"数据分析"
- "库存报表"、"项目进度"

【参数说明】
- report_type: 报告类型
- project_id: 项目ID（可选）

【报告类型】
- inventory_summary: 库存汇总
- material_plan_summary: 物资计划汇总
- project_progress: 项目进度
- low_stock_alert: 低库存预警`, map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"report_type": map[string]interface{}{
					"type":        "string",
					"description": "报告类型：inventory_summary(库存汇总)/material_plan_summary(物资计划汇总)/project_progress(项目进度)/low_stock_alert(低库存预警)",
					"enum":        []string{"inventory_summary", "material_plan_summary", "project_progress", "low_stock_alert"},
				},
				"project_id": map[string]interface{}{
					"type":        "integer",
					"description": "项目ID，可选",
				},
			},
			"required": []string{"report_type"},
		}).
		Build()
}

// ParseToolCallArguments parses tool call arguments
func ParseToolCallArguments(args string, v interface{}) error {
	return json.Unmarshal([]byte(args), v)
}

// ToolCallResult represents the result of a tool call
type ToolCallResult struct {
	ToolCallID string
	Result     string
	IsError    bool
}

// FormatToolResult formats a tool result for the API
func FormatToolResult(toolCallID string, result string, isError bool) Message {
	content := result
	if isError {
		content = fmt.Sprintf("Error: %s", result)
	}
	return Message{
		Role:       "tool",
		ToolCallID: toolCallID,
		Content:    content,
	}
}
