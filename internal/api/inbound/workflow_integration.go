package inbound

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

// StartInboundWorkflow 启动入库单工作流
func (wi *WorkflowIntegration) StartInboundWorkflow(order *InboundOrder, creatorID uint, creatorName string) error {
	// 获取入库单的工作流定义
	wf, err := wi.engine.GetWorkflowByModule("inbound")
	if err != nil {
		return fmt.Errorf("获取工作流失败: %w", err)
	}

	// 启动工作流实例
	instance, err := wi.engine.StartWorkflow(
		wf.ID,
		"inbound_order",
		order.ID,
		order.OrderNo,
		creatorID,
		creatorName,
	)

	if err != nil {
		return fmt.Errorf("启动工作流失败: %w", err)
	}

	// 如果工作流不需要审批（直接完成），则执行入库逻辑
	if instance.Status == workflow.InstanceStatusApproved {
		return wi.executeInboundApproval(order, creatorID, creatorName, "", nil)
	}

	return nil
}

// ProcessInboundApproval 处理入库单审批
func (wi *WorkflowIntegration) ProcessInboundApproval(orderID uint, approverID uint, approverName string, action string, remark string, items []InboundApprovalItem) error {
	// 获取工作流实例
	instance, err := wi.engine.GetInstanceByBusiness("inbound_order", orderID)
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

	// 如果工作流已完全通过，执行入库逻辑
	if instance.Status == workflow.InstanceStatusApproved {
		var order InboundOrder
		if err := wi.db.Preload("Items").First(&order, orderID).Error; err != nil {
			return err
		}

		// 执行实际的入库操作
		if err := wi.executeInboundApproval(&order, approverID, approverName, remark, items); err != nil {
			return err
		}
	}

	return nil
}

// executeInboundApproval 执行入库单的实际入库逻辑
func (wi *WorkflowIntegration) executeInboundApproval(order *InboundOrder, approverID uint, approverName string, remark string, items []InboundApprovalItem) error {
	// 添加备注
	if remark != "" {
		if order.Remark != "" {
			order.Remark = order.Remark + "\n审批备注：" + remark
		} else {
			order.Remark = "审批备注：" + remark
		}
	}

	order.UpdatedAt = time.Now()

	// 创建审批数量映射
	approvedQtyMap := make(map[uint]float64)
	if len(items) > 0 {
		for _, item := range items {
			approvedQtyMap[item.ID] = float64(item.ApprovedQuantity)
		}
	} else {
		// 默认全部通过
		for _, orderItem := range order.Items {
			approvedQtyMap[orderItem.ID] = orderItem.Quantity
		}
	}

	// 同步到库存
	for _, item := range order.Items {
		// 获取审批数量，默认为原始数量
		approvedQty := item.Quantity
		if qty, ok := approvedQtyMap[item.ID]; ok {
			approvedQty = qty
		}

		// 获取单位信息
		var materialInfo struct {
			Unit string
		}
		wi.db.Table("material_master").Where("id = ?", item.MaterialID).
			Select("unit").First(&materialInfo)

		// 尝试查找现有库存
		var existingStock struct {
			ID       uint
			Quantity float64
		}
		var stockID uint
		var quantityBefore float64

		err := wi.db.Table("stocks").
			Where("material_id = ?", item.MaterialID).
			First(&existingStock).
			Error

		if err != nil {
			// 没有找到，创建新库存记录
			currentTime := time.Now()
			newStock := map[string]interface{}{
				"project_id":   order.ProjectID,
				"material_id":  item.MaterialID,
				"quantity":     float64(approvedQty),
				"safety_stock": 0,
				"location":     "",
				"created_at":   currentTime,
				"updated_at":   currentTime,
			}

			if err := wi.db.Table("stocks").Create(&newStock).Error; err != nil {
				return fmt.Errorf("创建库存记录失败: %w", err)
			}

			// 获取新创建的库存ID
			wi.db.Table("stocks").Where("material_id = ?", item.MaterialID).
				Order("id DESC").First(&existingStock)
			stockID = existingStock.ID
			quantityBefore = 0
		} else {
			// 找到了，更新现有库存数量
			stockID = existingStock.ID
			quantityBefore = existingStock.Quantity

			wi.db.Table("stocks").
				Where("id = ?", stockID).
				Update("quantity", gorm.Expr("quantity + ?", float64(approvedQty))).
				Update("updated_at", time.Now())
		}

		// 记录库存操作日志
		detail := fmt.Sprintf("入库 %.2f %s，备注：入库单 %s", float64(approvedQty), materialInfo.Unit, order.OrderNo)

		// 创建库存日志
		stockLog := map[string]interface{}{
			"stock_id":        stockID,
			"type":            "in",
			"quantity":        float64(approvedQty),
			"quantity_before": quantityBefore,
			"quantity_after":   quantityBefore + float64(approvedQty),
			"time":            time.Now(),
			"remark":          detail,
			"project_id":      order.ProjectID,
			"user_id":         approverID,
			"requisition_id":  nil,
			"inbound_code":    order.OrderNo,
		}
		wi.db.Table("stock_logs").Create(&stockLog)
	}



	// 更新订单状态为已完成
	order.Status = StatusCompleted
	return wi.db.Save(order).Error
}

// getMaterialIDsFromItems 从入库单项中提取物资ID列表
func getMaterialIDsFromItems(items []InboundOrderItem) []uint {
	materialIDs := make([]uint, 0, len(items))
	for _, item := range items {
		materialIDs = append(materialIDs, item.MaterialID)
	}
	return materialIDs
}

// InboundApprovalItem 入库单审批项
type InboundApprovalItem struct {
	ID               uint `json:"id"`
	ApprovedQuantity int  `json:"approved_quantity"`
}

// GetInboundWorkflowStatus 获取入库单的工作流状态
func (wi *WorkflowIntegration) GetInboundWorkflowStatus(orderID uint) (*workflow.WorkflowInstance, error) {
	instance, err := wi.engine.GetInstanceByBusiness("inbound_order", orderID)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

// GetInboundPendingTasks 获取入库单的待办任务
func (wi *WorkflowIntegration) GetInboundPendingTasks(approverID uint) ([]workflow.WorkflowPendingTask, error) {
	return wi.engine.GetPendingTasksByBusiness(approverID, "inbound_order")
}

// ResubmitInboundOrder 重新提交入库单（退回后）
func (wi *WorkflowIntegration) ResubmitInboundOrder(orderID uint, submitterID uint, submitterName string) error {
	instance, err := wi.engine.GetInstanceByBusiness("inbound_order", orderID)
	if err != nil {
		return errors.New("未找到工作流实例")
	}

	return wi.engine.Resubmit(instance.ID, submitterID, submitterName)
}

// GetInboundWorkflowApprovals 获取入库单的审批记录
func (wi *WorkflowIntegration) GetInboundWorkflowApprovals(orderID uint) ([]workflow.WorkflowApproval, error) {
	instance, err := wi.engine.GetInstanceByBusiness("inbound_order", orderID)
	if err != nil {
		return nil, err
	}

	return wi.engine.GetInstanceApprovals(instance.ID)
}
