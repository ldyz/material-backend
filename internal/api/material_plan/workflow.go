package material_plan

import (
	"errors"
	"fmt"
	"time"

	"github.com/yourorg/material-backend/backend/internal/api/workflow"
	"gorm.io/gorm"
)

// WorkflowIntegration 工作流集成辅助类
type WorkflowIntegration struct {
	db     *gorm.DB
	engine *workflow.Engine
}

// NewWorkflowIntegration 创建工作流集成实例
func NewWorkflowIntegration(db *gorm.DB) *WorkflowIntegration {
	return &WorkflowIntegration{
		db:     db,
		engine: workflow.NewEngine(db),
	}
}

// StartPlanWorkflow 启动物资计划工作流
func (wi *WorkflowIntegration) StartPlanWorkflow(plan *MaterialPlan, creatorID uint, creatorName string) error {
	// 获取物资计划的工作流定义
	wf, err := wi.engine.GetWorkflowByModule("material_plan")
	if err != nil {
		return fmt.Errorf("获取工作流失败: %w", err)
	}

	// 启动工作流实例
	instance, err := wi.engine.StartWorkflow(
		wf.ID,
		"material_plan",
		plan.ID,
		plan.PlanNo,
		creatorID,
		creatorName,
		nil, // material_plan 不关联项目
	)

	if err != nil {
		return fmt.Errorf("启动工作流失败: %w", err)
	}

	// 保存工作流实例ID到计划，并更新状态为待审批
	plan.WorkflowInstanceID = &instance.ID
	plan.Status = PlanStatusPending
	if err := wi.db.Save(plan).Error; err != nil {
		return fmt.Errorf("更新计划状态失败: %w", err)
	}

	// 如果工作流不需要审批（直接完成），则自动通过
	if instance.Status == workflow.InstanceStatusApproved {
		return wi.executePlanApproval(plan, creatorID, creatorName, "")
	}

	return nil
}

// ProcessPlanApproval 处理物资计划审批
func (wi *WorkflowIntegration) ProcessPlanApproval(planID uint, approverID uint, approverName string, action string, remark string) error {
	// 获取工作流实例
	instance, err := wi.engine.GetInstanceByBusiness("material_plan", planID)
	if err != nil {
		return errors.New("未找到工作流实例")
	}

	// 处理审批
	err = wi.engine.ProcessApproval(instance.ID, approverID, approverName, action, remark, nil)
	if err != nil {
		return err
	}

	// 重新加载实例状态
	wi.db.Preload("Workflow").First(&instance, instance.ID)

	// 如果工作流已完全通过，更新计划状态
	if instance.Status == workflow.InstanceStatusApproved {
		var plan MaterialPlan
		if err := wi.db.First(&plan, planID).Error; err != nil {
			return err
		}

		// 执行批准后的逻辑
		if err := wi.executePlanApproval(&plan, approverID, approverName, remark); err != nil {
			return err
		}
	} else if instance.Status == workflow.InstanceStatusRejected {
		// 工作流被拒绝
		var plan MaterialPlan
		if err := wi.db.First(&plan, planID).Error; err != nil {
			return err
		}

		now := time.Now()
		plan.Status = PlanStatusRejected
		plan.ApproverID = &approverID
		plan.ApproverName = approverName
		plan.ApprovedAt = &now

		if remark != "" {
			if plan.Remark != "" {
				plan.Remark = plan.Remark + "\n审批备注：" + remark
			} else {
				plan.Remark = "审批备注：" + remark
			}
		}

		if err := wi.db.Save(&plan).Error; err != nil {
			return err
		}
	}

	return nil
}

// executePlanApproval 执行计划批准后的逻辑
func (wi *WorkflowIntegration) executePlanApproval(plan *MaterialPlan, approverID uint, approverName string, remark string) error {
	// 添加备注
	if remark != "" {
		if plan.Remark != "" {
			plan.Remark = plan.Remark + "\n审批备注：" + remark
		} else {
			plan.Remark = "审批备注：" + remark
		}
	}

	now := time.Now()
	plan.Status = PlanStatusApproved
	plan.ApproverID = &approverID
	plan.ApproverName = approverName
	plan.ApprovedAt = &now

	// 保存计划状态
	if err := wi.db.Save(plan).Error; err != nil {
		return fmt.Errorf("保存计划状态失败: %w", err)
	}

	// 自动将计划物资添加到物资库
	if err := wi.syncPlanItemsToMaterials(plan); err != nil {
		return fmt.Errorf("同步物资到物资库失败: %w", err)
	}

	return nil
}

// syncPlanItemsToMaterials 将计划物资同步到物资库
func (wi *WorkflowIntegration) syncPlanItemsToMaterials(plan *MaterialPlan) error {
	// 获取计划的所有项
	var items []MaterialPlanItem
	if err := wi.db.Where("plan_id = ?", plan.ID).Find(&items).Error; err != nil {
		return err
	}

	// 遍历每一项，处理material_id和stock
	for _, item := range items {
		// 新结构中 MaterialID 已经是必需的字段（uint 类型，不是指针）
		// 所以不需要创建新的物资记录，只需要验证物资存在即可

		// 验证物资在 material_master 表中存在
		var material struct {
			ID   uint
			Name string
		}
		if err := wi.db.Table("material_master").Where("id = ?", item.MaterialID).
			First(&material).Error; err != nil {
			return fmt.Errorf("计划项 ID %d 的 material_id %d 在 material_master 中不存在", item.ID, item.MaterialID)
		}

		// 物资已验证，不需要额外操作
		// item.MaterialID 已经正确设置，直接继续处理下一项
	}

	return nil
}

// GetPlanWorkflowStatus 获取计划的工作流状态
func (wi *WorkflowIntegration) GetPlanWorkflowStatus(planID uint) (*workflow.WorkflowInstance, error) {
	instance, err := wi.engine.GetInstanceByBusiness("material_plan", planID)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

// GetPlanPendingTasks 获取计划的待办任务
func (wi *WorkflowIntegration) GetPlanPendingTasks(approverID uint) ([]workflow.WorkflowPendingTask, error) {
	return wi.engine.GetPendingTasksByBusiness(approverID, "material_plan")
}

// ResubmitPlan 重新提交计划（退回后）
func (wi *WorkflowIntegration) ResubmitPlan(planID uint, submitterID uint, submitterName string) error {
	instance, err := wi.engine.GetInstanceByBusiness("material_plan", planID)
	if err != nil {
		return errors.New("未找到工作流实例")
	}

	// 检查计划状态
	var plan MaterialPlan
	if err := wi.db.First(&plan, planID).Error; err != nil {
		return errors.New("计划不存在")
	}

	if plan.Status != PlanStatusRejected {
		return errors.New("只有被拒绝的计划才能重新提交")
	}

	// 重新提交工作流
	if err := wi.engine.Resubmit(instance.ID, submitterID, submitterName); err != nil {
		return err
	}

	// 更新计划状态为待审批
	plan.Status = PlanStatusPending
	plan.Remark = ""
	if err := wi.db.Save(&plan).Error; err != nil {
		return err
	}

	return nil
}

// GetPlanWorkflowApprovals 获取计划的审批记录
func (wi *WorkflowIntegration) GetPlanWorkflowApprovals(planID uint) ([]workflow.WorkflowApproval, error) {
	instance, err := wi.engine.GetInstanceByBusiness("material_plan", planID)
	if err != nil {
		return nil, err
	}

	return wi.engine.GetInstanceApprovals(instance.ID)
}

// CancelPlanWorkflow 取消计划的工作流
func (wi *WorkflowIntegration) CancelPlanWorkflow(planID uint, userID uint, userName string, reason string) error {
	_, err := wi.engine.GetInstanceByBusiness("material_plan", planID)
	if err != nil {
		return errors.New("未找到工作流实例")
	}

	// 更新计划状态
	var plan MaterialPlan
	if err := wi.db.First(&plan, planID).Error; err != nil {
		return errors.New("计划不存在")
	}

	// 只有草稿或待审批状态的计划可以取消
	if plan.Status != PlanStatusDraft && plan.Status != PlanStatusPending {
		return errors.New("只有草稿或待审批状态的计划可以取消")
	}

	plan.Status = PlanStatusCancelled
	if reason != "" {
		if plan.Remark != "" {
			plan.Remark = plan.Remark + "\n取消原因：" + reason
		} else {
			plan.Remark = "取消原因：" + reason
		}
	}

	if err := wi.db.Save(&plan).Error; err != nil {
		return err
	}

	return nil
}
