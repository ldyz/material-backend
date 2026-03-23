package requisition

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

// StartRequisitionWorkflow 启动领料单工作流
func (wi *WorkflowIntegration) StartRequisitionWorkflow(requisition *Requisition, creatorID uint, creatorName string) error {
	// 获取领料单的工作流定义
	wf, err := wi.engine.GetWorkflowByModule("requisition")
	if err != nil {
		return fmt.Errorf("获取工作流失败: %w", err)
	}

	// 启动工作流实例
	instance, err := wi.engine.StartWorkflow(
		wf.ID,
		"requisition",
		requisition.ID,
		requisition.RequisitionNo,
		creatorID,
		creatorName,
		nil, // requisition 不关联项目
	)

	if err != nil {
		return fmt.Errorf("启动工作流失败: %w", err)
	}

	// 如果工作流不需要审批（直接完成），则执行发料逻辑
	if instance.Status == workflow.InstanceStatusApproved {
		return wi.executeRequisitionApproval(requisition, creatorID, creatorName, "", nil)
	}

	return nil
}

// ProcessRequisitionApproval 处理领料单审批
func (wi *WorkflowIntegration) ProcessRequisitionApproval(requisitionID uint, approverID uint, approverName string, action string, remark string, items []RequisitionApprovalItem) error {
	// 获取工作流实例
	instance, err := wi.engine.GetInstanceByBusiness("requisition", requisitionID)
	if err != nil {
		return errors.New("未找到工作流实例")
	}

	// 构建额外数据
	extraData := map[string]any{
		"items": items,
	}

	// 处理审批
	err = wi.engine.ProcessApproval(instance.ID, approverID, approverName, action, remark, extraData)
	if err != nil {
		return err
	}

	// 重新加载实例状态
	wi.db.Preload("Workflow").First(&instance, instance.ID)

	// 如果工作流已完全通过，执行发料逻辑
	if instance.Status == workflow.InstanceStatusApproved {
		var requisition Requisition
		if err := wi.db.Preload("Items").First(&requisition, requisitionID).Error; err != nil {
			return err
		}

		// 执行实际的发料操作
		if err := wi.executeRequisitionApproval(&requisition, approverID, approverName, remark, items); err != nil {
			return err
		}
	}

	return nil
}

// executeRequisitionApproval 执行领料单的实际发料逻辑
func (wi *WorkflowIntegration) executeRequisitionApproval(requisition *Requisition, approverID uint, approverName string, remark string, items []RequisitionApprovalItem) error {
	// 添加备注
	if remark != "" {
		if requisition.Remark != "" {
			requisition.Remark = requisition.Remark + "\n审批备注：" + remark
		} else {
			requisition.Remark = "审批备注：" + remark
		}
	}

	// 创建审批数量映射
	approvedQtyMap := make(map[uint]float64)
	if len(items) > 0 {
		for _, item := range items {
			approvedQtyMap[item.ID] = float64(item.ApprovedQuantity)
		}
	} else {
		// 默认全部通过
		for _, requisitionItem := range requisition.Items {
			approvedQtyMap[requisitionItem.ID] = requisitionItem.RequestedQuantity
		}
	}

	// 同步到库存
	for _, item := range requisition.Items {
		// 获取审批数量，默认为原始数量
		approvedQty := item.ApprovedQuantity
		if qty, ok := approvedQtyMap[item.ID]; ok {
			approvedQty = qty
		}

		// 查找库存记录（使用 StockID）
		var stock struct {
			ID          uint
			MaterialID  uint
			Quantity    float64
			SafetyStock float64
			Location    string
			Unit        string
		}

		// 查找现有库存
		if item.StockID != 0 {
			if err := wi.db.Where("id = ?", item.StockID).First(&stock).Error; err != nil {
				return fmt.Errorf("材料 %d 库存不存在", item.StockID)
			}
		} else {
			// 如果没有 StockID，使用 MaterialID 查找
			if err := wi.db.Where("material_id = ? AND project_id = ?",
				item.MaterialID, requisition.ProjectID).First(&stock).Error; err != nil {
				return fmt.Errorf("材料ID %d 库存不存在", item.MaterialID)
			}
		}

		// 检查库存是否充足
		if stock.Quantity < approvedQty {
			return fmt.Errorf("材料ID %d 库存不足，当前库存: %.2f，需要: %.2f",
				item.MaterialID, stock.Quantity, approvedQty)
		}

		// 更新库存数量
		stock.Quantity -= approvedQty
		wi.db.Save(&stock)

		// 记录库存操作日志
		detail := fmt.Sprintf("出库 %.2f，备注：领料单 %s", approvedQty, requisition.RequisitionNo)

		// 获取项目ID
		var stockProject struct {
			ProjectID uint
		}
		wi.db.Table("stocks").Where("id = ?", stock.ID).Select("project_id").First(&stockProject)
		projectID := stockProject.ProjectID
		if projectID == 0 {
			projectID = requisition.ProjectID // 回退到领料单的项目ID
		}

		// 创建库存日志（使用正确的字段结构）
		stockLog := map[string]interface{}{
			"stock_id":        stock.ID,
			"type":            "out",
			"quantity":        approvedQty,
			"quantity_before": stock.Quantity + approvedQty,
			"quantity_after":  stock.Quantity,
			"source_type":     "requisition",
			"source_id":       requisition.ID,
			"source_no":       requisition.RequisitionNo,
			"project_id":      projectID,
			"material_id":     item.MaterialID,
			"user_id":         approverID,
			"remark":          detail,
			"created_at":      time.Now(),
		}
		wi.db.Table("stock_logs").Create(&stockLog)

		// 创建库存操作日志
		opLog := map[string]interface{}{
			"stock_id": stock.ID,
			"op_type":  "out",
			"detail":   detail,
			"user_id":  approverID,
			"time":     time.Now(),
		}
		wi.db.Table("stock_op_logs").Create(&opLog)
	}

	// 更新领料单状态为已发料
	now := time.Now()
	requisition.Status = "issued"
	requisition.IssuedAt = &now
	requisition.IssuedBy = approverName
	return wi.db.Save(requisition).Error
}

// RequisitionApprovalItem 领料单审批项
type RequisitionApprovalItem struct {
	ID               uint `json:"id"`
	ApprovedQuantity int  `json:"approved_quantity"`
}

// GetRequisitionWorkflowStatus 获取领料单的工作流状态
func (wi *WorkflowIntegration) GetRequisitionWorkflowStatus(requisitionID uint) (*workflow.WorkflowInstance, error) {
	instance, err := wi.engine.GetInstanceByBusiness("requisition", requisitionID)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

// GetRequisitionPendingTasks 获取领料单的待办任务
func (wi *WorkflowIntegration) GetRequisitionPendingTasks(approverID uint) ([]workflow.WorkflowPendingTask, error) {
	return wi.engine.GetPendingTasksByBusiness(approverID, "requisition")
}

// ResubmitRequisition 重新提交领料单（退回后）
func (wi *WorkflowIntegration) ResubmitRequisition(requisitionID uint, submitterID uint, submitterName string) error {
	instance, err := wi.engine.GetInstanceByBusiness("requisition", requisitionID)
	if err != nil {
		return errors.New("未找到工作流实例")
	}

	return wi.engine.Resubmit(instance.ID, submitterID, submitterName)
}

// GetRequisitionWorkflowApprovals 获取领料单的审批记录
func (wi *WorkflowIntegration) GetRequisitionWorkflowApprovals(requisitionID uint) ([]workflow.WorkflowApproval, error) {
	instance, err := wi.engine.GetInstanceByBusiness("requisition", requisitionID)
	if err != nil {
		return nil, err
	}

	return wi.engine.GetInstanceApprovals(instance.ID)
}
