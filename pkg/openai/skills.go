package openai

import (
	"fmt"
	"time"
)

// Skill 预设技能定义
type Skill struct {
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	Tools         []string               `json:"tools"`          // 需要使用的工具列表
	DefaultParams map[string]interface{} `json:"default_params"` // 默认参数
	Examples      []SkillExample         `json:"examples"`       // 使用示例
	Priority      int                    `json:"priority"`       // 匹配优先级
}

// SkillExample 技能使用示例
type SkillExample struct {
	UserQuery     string              `json:"user_query"`     // 用户提问示例
	ToolSequence  []SkillToolCall     `json:"tool_sequence"`  // 工具调用序列
	ExpectedReply string              `json:"expected_reply"` // 预期回复
}

// SkillToolCall 技能工具调用定义
type SkillToolCall struct {
	ToolName string                 `json:"tool_name"`
	Params   map[string]interface{} `json:"params"`
}

// Skills 预设技能映射
var Skills = map[string]Skill{
	// ==================== 查询类技能 ====================
	"check_today_tasks": {
		Name:        "查看今日任务",
		Description: "查询今天的施工预约/任务安排",
		Tools:       []string{"query_appointments"},
		DefaultParams: map[string]interface{}{
			"date": "today",
		},
		Priority: 10,
		Examples: []SkillExample{
			{
				UserQuery: "今天有什么任务",
				ToolSequence: []SkillToolCall{
					{ToolName: "query_appointments", Params: map[string]interface{}{"date": "today"}},
				},
				ExpectedReply: "今天有X个任务，分别是...",
			},
		},
	},

	"check_tomorrow_tasks": {
		Name:        "查看明日任务",
		Description: "查询明天的施工预约/任务安排",
		Tools:       []string{"query_appointments"},
		DefaultParams: map[string]interface{}{
			"date": "tomorrow",
		},
		Priority: 10,
		Examples: []SkillExample{
			{
				UserQuery: "明天的任务安排",
				ToolSequence: []SkillToolCall{
					{ToolName: "query_appointments", Params: map[string]interface{}{"date": "tomorrow"}},
				},
				ExpectedReply: "明天（X月X日）有X个任务...",
			},
		},
	},

	"check_week_tasks": {
		Name:        "查看本周任务",
		Description: "查询本周的施工预约/任务安排",
		Tools:       []string{"query_appointments"},
		DefaultParams: map[string]interface{}{
			"start_date": "week_start",
			"end_date":   "week_end",
		},
		Priority: 9,
	},

	"check_low_stock": {
		Name:        "库存预警检查",
		Description: "检查库存不足的物资，返回低于安全库存的物资列表",
		Tools:       []string{"query_stock_alerts"},
		Priority:    8,
		Examples: []SkillExample{
			{
				UserQuery: "库存够不够",
				ToolSequence: []SkillToolCall{
					{ToolName: "query_stock_alerts", Params: map[string]interface{}{}},
				},
				ExpectedReply: "目前有X种物资库存偏低...",
			},
		},
	},

	"check_stock_status": {
		Name:        "查看库存情况",
		Description: "查询当前库存状态",
		Tools:       []string{"query_stock"},
		DefaultParams: map[string]interface{}{
			"page_size": 20,
		},
		Priority: 7,
	},

	"check_pending_approvals": {
		Name:        "查看待审批任务",
		Description: "查询需要当前用户审批的预约/任务",
		Tools:       []string{"query_pending_approvals"},
		Priority:    8,
		Examples: []SkillExample{
			{
				UserQuery: "有什么需要审批的",
				ToolSequence: []SkillToolCall{
					{ToolName: "query_pending_approvals", Params: map[string]interface{}{}},
				},
				ExpectedReply: "您有X个待审批的任务...",
			},
		},
	},

	"check_available_workers": {
		Name:        "查看可用人员",
		Description: "查询可安排的作业人员列表",
		Tools:       []string{"query_workers"},
		Priority:    6,
	},

	// ==================== 操作类技能 ====================
	"quick_approval": {
		Name:        "快速审批",
		Description: "审批待处理的预约单（需要先查询待审批列表）",
		Tools:       []string{"query_pending_approvals", "approve_appointment"},
		Priority:    7,
		Examples: []SkillExample{
			{
				UserQuery: "帮我审批通过第一个",
				ToolSequence: []SkillToolCall{
					{ToolName: "query_pending_approvals", Params: map[string]interface{}{}},
					{ToolName: "approve_appointment", Params: map[string]interface{}{"action": "approve"}},
				},
			},
		},
	},

	"create_task": {
		Name:        "创建施工任务",
		Description: "创建新的施工预约/任务（需要先获取项目和人员信息）",
		Tools:       []string{"query_projects", "query_workers", "create_appointment"},
		Priority:    5,
		Examples: []SkillExample{
			{
				UserQuery: "帮我创建一个施工预约",
				ToolSequence: []SkillToolCall{
					{ToolName: "query_projects", Params: map[string]interface{}{}},
					{ToolName: "query_workers", Params: map[string]interface{}{}},
					{ToolName: "create_appointment", Params: map[string]interface{}{}},
				},
				ExpectedReply: "已为您创建预约，编号XXX...",
			},
		},
	},

	"create_urgent_task": {
		Name:        "创建加急任务",
		Description: "创建加急的施工预约",
		Tools:       []string{"query_projects", "query_workers", "create_appointment"},
		DefaultParams: map[string]interface{}{
			"is_urgent": true,
			"priority":  8,
		},
		Priority: 6,
	},

	"cancel_my_task": {
		Name:        "取消我的任务",
		Description: "取消已创建的预约/任务",
		Tools:       []string{"query_appointments", "cancel_appointment"},
		Priority:    6,
	},

	// ==================== 分析类技能 ====================
	"inventory_report": {
		Name:        "生成库存报告",
		Description: "生成库存汇总分析报告",
		Tools:       []string{"generate_report"},
		DefaultParams: map[string]interface{}{
			"report_type": "inventory_summary",
		},
		Priority: 4,
	},

	"low_stock_report": {
		Name:        "生成库存预警报告",
		Description: "生成低库存预警报告",
		Tools:       []string{"generate_report"},
		DefaultParams: map[string]interface{}{
			"report_type": "low_stock_alert",
		},
		Priority: 5,
	},
}

// GetSkill 获取技能定义
func GetSkill(name string) (Skill, bool) {
	skill, ok := Skills[name]
	return skill, ok
}

// GetAllSkills 获取所有技能列表
func GetAllSkills() []Skill {
	result := make([]Skill, 0, len(Skills))
	for _, skill := range Skills {
		result = append(result, skill)
	}
	return result
}

// ResolveSkillParams 解析技能参数，替换特殊值为实际值
func ResolveSkillParams(params map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	now := time.Now()

	for k, v := range params {
		switch val := v.(type) {
		case string:
			switch val {
			case "today":
				result[k] = now.Format("2006-01-02")
			case "tomorrow":
				result[k] = now.AddDate(0, 0, 1).Format("2006-01-02")
			case "yesterday":
				result[k] = now.AddDate(0, 0, -1).Format("2006-01-02")
			case "week_start":
				// 本周一
				weekday := int(now.Weekday())
				if weekday == 0 {
					weekday = 7
				}
				weekStart := now.AddDate(0, 0, 1-weekday)
				result[k] = weekStart.Format("2006-01-02")
			case "week_end":
				// 本周日
				weekday := int(now.Weekday())
				if weekday == 0 {
					weekday = 7
				}
				weekEnd := now.AddDate(0, 0, 7-weekday)
				result[k] = weekEnd.Format("2006-01-02")
			default:
				result[k] = val
			}
		default:
			result[k] = v
		}
	}

	return result
}

// MatchSkill 根据用户意图匹配合适的技能
func MatchSkill(userMessage string) *Skill {
	// 简单的关键词匹配逻辑
	keywords := map[string]string{
		"今天":   "check_today_tasks",
		"今日":   "check_today_tasks",
		"明天":   "check_tomorrow_tasks",
		"明日":   "check_tomorrow_tasks",
		"本周":   "check_week_tasks",
		"这周":   "check_week_tasks",
		"库存预警": "check_low_stock",
		"低库存":  "check_low_stock",
		"库存不足": "check_low_stock",
		"库存":   "check_stock_status",
		"待审批":  "check_pending_approvals",
		"审批任务": "check_pending_approvals",
		"作业人员": "check_available_workers",
		"工人":   "check_available_workers",
		"创建":   "create_task",
		"新建":   "create_task",
		"加急":   "create_urgent_task",
	}

	for keyword, skillName := range keywords {
		if contains(userMessage, keyword) {
			if skill, ok := Skills[skillName]; ok {
				return &skill
			}
		}
	}

	return nil
}

// contains 检查字符串是否包含子串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// GetSkillPrompt 生成技能相关的系统提示
func GetSkillPrompt() string {
	var prompt string
	prompt += "\n## 预设技能\n\n"
	prompt += "以下是一些常用的操作技能，可以根据用户意图自动识别并执行：\n\n"

	for name, skill := range Skills {
		prompt += fmt.Sprintf("### %s (%s)\n", skill.Name, name)
		prompt += fmt.Sprintf("- 描述: %s\n", skill.Description)
		prompt += fmt.Sprintf("- 使用工具: %v\n\n", skill.Tools)
	}

	return prompt
}
