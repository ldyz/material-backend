package material_plan

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Service handles business logic for material plans
type Service struct {
	db *gorm.DB
}

// NewService creates a new material plan service
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// GeneratePlanNo generates a unique plan number
func (s *Service) GeneratePlanNo() (string, error) {
	now := time.Now()
	prefix := fmt.Sprintf("MP%s%02d%02d",
		strings.ToUpper(now.Format("2006")[2:]), now.Month(), now.Day())

	// 查找今天最大的纯数字序号计划（格式：MP2602010001）
	// 使用 PostgreSQL 正则表达式匹配纯数字后缀
	var lastPlan MaterialPlan
	err := s.db.Where("plan_no ~ ?", fmt.Sprintf("^%s[0-9]{4}$", prefix)).
		Order("plan_no DESC").
		First(&lastPlan).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Sprintf("%s%04d", prefix, 1), nil
		}
		return "", err
	}

	// Extract sequence number from last plan number
	seqStr := strings.TrimPrefix(lastPlan.PlanNo, prefix)
	seq := 1
	if seqStr != "" {
		_, err := fmt.Sscanf(seqStr, "%d", &seq)
		if err == nil {
			seq++
		}
	}

	return fmt.Sprintf("%s%04d", prefix, seq), nil
}

// mergeDuplicateItems merges items that reference the same material
// Returns a deduplicated list with quantities summed
func (s *Service) mergeDuplicateItems(items *[]CreateMaterialPlanItemRequest) []CreateMaterialPlanItemRequest {
	// Use a map to track unique materials by key
	// Key is either material_id (if provided) or material_name+specification
	type materialKey struct {
		MaterialID    uint
		MaterialName  string
		Specification string
	}

	mergedMap := make(map[materialKey]*CreateMaterialPlanItemRequest)

	for _, item := range *items {
		var key materialKey

		if item.MaterialID > 0 {
			// Use material_id as the key
			key = materialKey{
				MaterialID:    item.MaterialID,
				MaterialName:  "",
				Specification: "",
			}
		} else {
			// Use material_name+specification as the key
			key = materialKey{
				MaterialID:    0,
				MaterialName:  item.MaterialName,
				Specification: item.Specification,
			}
		}

		if existing, found := mergedMap[key]; found {
			// Merge: sum quantities and keep the first item's other fields
			existing.PlannedQuantity += item.PlannedQuantity
			// Use the higher unit price if different
			if item.UnitPrice > existing.UnitPrice {
				existing.UnitPrice = item.UnitPrice
			}
			// Combine remarks
			if item.Remark != "" && existing.Remark != "" {
				existing.Remark = existing.Remark + "; " + item.Remark
			} else if item.Remark != "" {
				existing.Remark = item.Remark
			}
		} else {
			// Add new item to the map
			itemCopy := item
			mergedMap[key] = &itemCopy
		}
	}

	// Convert map back to slice
	result := make([]CreateMaterialPlanItemRequest, 0, len(mergedMap))
	for _, item := range mergedMap {
		result = append(result, *item)
	}

	return result
}

// CreatePlan creates a new material plan
func (s *Service) CreatePlan(req *CreateMaterialPlanRequest, creatorID uint, creatorName string) (*MaterialPlan, error) {
	// Generate plan number
	planNo, err := s.GeneratePlanNo()
	if err != nil {
		return nil, fmt.Errorf("failed to generate plan number: %w", err)
	}

	// Validate project exists
	var project struct {
		ID uint
	}
	if err := s.db.Table("projects").Where("id = ?", req.ProjectID).First(&project).Error; err != nil {
		return nil, errors.New("project not found")
	}

	// Parse dates
	var plannedStartDate, plannedEndDate *time.Time
	if req.PlannedStartDate != "" {
		t, err := time.Parse("2006-01-02", req.PlannedStartDate)
		if err != nil {
			return nil, errors.New("invalid planned_start_date format, use YYYY-MM-DD")
		}
		plannedStartDate = &t
	}
	if req.PlannedEndDate != "" {
		t, err := time.Parse("2006-01-02", req.PlannedEndDate)
		if err != nil {
			return nil, errors.New("invalid planned_end_date format, use YYYY-MM-DD")
		}
		plannedEndDate = &t
	}

	// Validate date range
	if plannedStartDate != nil && plannedEndDate != nil && plannedEndDate.Before(*plannedStartDate) {
		return nil, errors.New("planned_end_date must be after planned_start_date")
	}

	// Create plan
	plan := &MaterialPlan{
		PlanNo:           planNo,
		PlanName:         req.PlanName,
		ProjectID:        req.ProjectID,
		PlanType:         req.PlanType,
		Status:           PlanStatusDraft,
		Priority:         req.Priority,
		PlannedStartDate: plannedStartDate,
		PlannedEndDate:   plannedEndDate,
		TotalBudget:      int(req.TotalBudget * 100),
		Description:      req.Description,
		Remark:           req.Remark,
		CreatorID:        creatorID,
		CreatorName:      creatorName,
	}

	// Set default values
	if plan.PlanType == "" {
		plan.PlanType = PlanTypeProcurement
	}
	if plan.Priority == "" {
		plan.Priority = PriorityNormal
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create plan
	if err := tx.Create(plan).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create plan: %w", err)
	}

	// Create items
	if len(req.Items) > 0 {
		// Merge duplicate items by material_id or material_name+specification
		mergedItems := s.mergeDuplicateItems(&req.Items)

		for i, itemReq := range mergedItems {
			item, err := s.CreatePlanItem(tx, plan.ID, &itemReq)
			if err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to create item %d: %w", i+1, err)
			}
			plan.Items = append(plan.Items, *item)
		}

		// Recalculate total budget from items if not provided
		if req.TotalBudget == 0 {
			plan.CalculateTotalBudget()
			tx.Save(plan)
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return plan, nil
}

// CreatePlanItem creates a plan item
func (s *Service) CreatePlanItem(tx *gorm.DB, planID uint, req *CreateMaterialPlanItemRequest) (*MaterialPlanItem, error) {
	var materialID uint

	// Debug log
	fmt.Printf("[DEBUG] CreatePlanItem: material_name=%q, material=%q, specification=%q, material_id=%d\n",
		req.MaterialName, req.Material, req.Specification, req.MaterialID)

	// Strategy 1: If material_id is provided, use it directly
	if req.MaterialID > 0 {
		// Validate material exists in material_master
		var material struct {
			ID uint
		}
		if err := tx.Table("material_master").Where("id = ?", req.MaterialID).First(&material).Error; err != nil {
			return nil, errors.New("material not found in material_master")
		}
		materialID = req.MaterialID
	} else if req.MaterialName != "" {
		// Strategy 2: Try to find material by name, specification AND material
		var material struct {
			ID       uint
			Material string
		}
		query := tx.Table("material_master").Where("name = ?", req.MaterialName)

		// If specification is provided, also match by specification
		if req.Specification != "" {
			query = query.Where("specification = ?", req.Specification)
		}

		// Also match by material - only exact match on all three fields
		if req.Material != "" {
			query = query.Where("material = ?", req.Material)
		} else {
			// If material is empty, only match records where material is also empty
			query = query.Where("(material = '' OR material IS NULL)")
		}

		if err := query.First(&material).Error; err == nil {
			// Found exact match (name + specification + material all match)
			materialID = material.ID
			fmt.Printf("[DEBUG] Found exact match material_master ID=%d (name=%q, spec=%q, material=%q)\n",
				materialID, req.MaterialName, req.Specification, req.Material)
		} else {
			// No exact match found - create new record
			fmt.Printf("[DEBUG] No exact match found (name=%q, spec=%q, material=%q). Creating new record.\n",
				req.MaterialName, req.Specification, req.Material)
			materialID = 0
		}

		// If no exact match, create new record
		if materialID == 0 {
			// Strategy 3: Auto-create material_master record
			// Generate a unique code if not provided
			code := req.MaterialCode
			if code == "" {
				// Get next ID from sequence to generate unique code
				var nextID int64
				if err := tx.Raw("SELECT nextval('material_master_id_seq')").Scan(&nextID).Error; err != nil {
					return nil, fmt.Errorf("failed to get next sequence value: %w", err)
				}
				code = fmt.Sprintf("AUTO%d", nextID)
			}

			newMaterial := struct {
				ID            uint
				Code          string
				Name          string
				Specification string
				Material      string
				Category      string
				Unit          string
			}{
				Code:          code,
				Name:          req.MaterialName,
				Specification: req.Specification,
				Material:      req.Material,
				Category:      req.Category,
				Unit:          req.Unit,
			}

			fmt.Printf("[DEBUG] Auto-creating material_master: code=%q, name=%q, material=%q, specification=%q\n",
				newMaterial.Code, newMaterial.Name, newMaterial.Material, newMaterial.Specification)

			if err := tx.Table("material_master").Create(&newMaterial).Error; err != nil {
				return nil, fmt.Errorf("failed to auto-create material_master: %w", err)
			}

			fmt.Printf("[DEBUG] Created material_master with ID=%d, material=%q\n", newMaterial.ID, newMaterial.Material)

			// Use the ID from the newly created material directly
			materialID = newMaterial.ID
		}
	} else {
		return nil, errors.New("either material_id or material_name is required")
	}

	// Parse required date
	var requiredDate *time.Time
	if req.RequiredDate != "" {
		t, err := time.Parse("2006-01-02", req.RequiredDate)
		if err != nil {
			return nil, errors.New("invalid required_date format, use YYYY-MM-DD")
		}
		requiredDate = &t
	}

	item := &MaterialPlanItem{
		PlanID:          planID,
		MaterialID:      materialID,
		PlannedQuantity: req.PlannedQuantity,
		UnitPrice:       req.UnitPrice,
		RequiredDate:    requiredDate,
		Priority:        req.Priority,
		Remark:          req.Remark,
		Status:          ItemStatusPending,
	}

	// Set default values
	if item.Priority == "" {
		item.Priority = PriorityNormal
	}

	if err := tx.Create(item).Error; err != nil {
		return nil, err
	}

	return item, nil
}

// UpdatePlan updates an existing material plan
func (s *Service) UpdatePlan(id uint, req *UpdateMaterialPlanRequest) (*MaterialPlan, error) {
	var plan MaterialPlan
	if err := s.db.Preload("Items").First(&plan, id).Error; err != nil {
		return nil, errors.New("plan not found")
	}

	// Check if plan can be edited
	if plan.Status != PlanStatusDraft {
		return nil, errors.New("only draft plans can be edited")
	}

	// Parse dates
	var plannedStartDate, plannedEndDate *time.Time
	if req.PlannedStartDate != "" {
		t, err := time.Parse("2006-01-02", req.PlannedStartDate)
		if err != nil {
			return nil, errors.New("invalid planned_start_date format, use YYYY-MM-DD")
		}
		plannedStartDate = &t
	}
	if req.PlannedEndDate != "" {
		t, err := time.Parse("2006-01-02", req.PlannedEndDate)
		if err != nil {
			return nil, errors.New("invalid planned_end_date format, use YYYY-MM-DD")
		}
		plannedEndDate = &t
	}

	// Update fields
	updates := map[string]any{}
	if req.PlanName != "" {
		updates["plan_name"] = req.PlanName
	}
	if req.PlanType != "" {
		updates["plan_type"] = req.PlanType
	}
	if req.Priority != "" {
		updates["priority"] = req.Priority
	}
	if plannedStartDate != nil {
		updates["planned_start_date"] = plannedStartDate
	}
	if plannedEndDate != nil {
		updates["planned_end_date"] = plannedEndDate
	}
	if req.TotalBudget > 0 {
		updates["total_budget"] = int(req.TotalBudget * 100)
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update plan
	if len(updates) > 0 {
		if err := tx.Model(&plan).Updates(updates).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update plan: %w", err)
		}
	}

	// Update items if provided
	if req.Items != nil {
		// Delete existing items
		if err := tx.Where("plan_id = ?", plan.ID).Delete(&MaterialPlanItem{}).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to delete old items: %w", err)
		}

		// Create new items
		plan.Items = nil
		for _, itemReq := range req.Items {
			item, err := s.CreatePlanItem(tx, plan.ID, &itemReq)
			if err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to create item: %w", err)
			}
			plan.Items = append(plan.Items, *item)
		}

		// Recalculate total budget
		plan.CalculateTotalBudget()
		tx.Save(plan)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &plan, nil
}

// GetPlan retrieves a plan by ID
func (s *Service) GetPlan(id uint) (*MaterialPlan, error) {
	var plan MaterialPlan
	if err := s.db.Preload("Items").First(&plan, id).Error; err != nil {
		return nil, errors.New("plan not found")
	}
	return &plan, nil
}

// ListPlans retrieves plans with filters
func (s *Service) ListPlans(projectID uint, status, priority string, page, pageSize int) ([]MaterialPlan, int64, error) {
	var plans []MaterialPlan
	var total int64

	query := s.db.Model(&MaterialPlan{})

	// Apply filters
	if projectID > 0 {
		query = query.Where("project_id = ?", projectID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if priority != "" {
		query = query.Where("priority = ?", priority)
	}

	// Count total
	query.Count(&total)

	// Paginate
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&plans).Error; err != nil {
		return nil, 0, err
	}

	return plans, total, nil
}

// DeletePlan deletes a plan
func (s *Service) DeletePlan(id uint) error {
	// Check if plan exists
	var plan MaterialPlan
	if err := s.db.First(&plan, id).Error; err != nil {
		return errors.New("plan not found")
	}

	// Check if plan can be deleted
	if plan.Status != PlanStatusDraft && plan.Status != PlanStatusCancelled {
		return errors.New("only draft or cancelled plans can be deleted")
	}

	// Delete plan (items will be cascade deleted)
	if err := s.db.Delete(&plan).Error; err != nil {
		return fmt.Errorf("failed to delete plan: %w", err)
	}

	return nil
}

// UpdatePlanStatus updates the status of a plan
func (s *Service) UpdatePlanStatus(id uint, status string, approverID *uint, approverName string) error {
	var plan MaterialPlan
	if err := s.db.First(&plan, id).Error; err != nil {
		return errors.New("plan not found")
	}

	// Validate status transition
	if !isValidStatusTransition(plan.Status, status) {
		return errors.New("invalid status transition")
	}

	updates := map[string]any{
		"status": status,
	}

	if approverID != nil {
		updates["approver_id"] = *approverID
		updates["approver_name"] = approverName
		now := time.Now()
		updates["approved_at"] = &now
	}

	if err := s.db.Model(&plan).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update plan status: %w", err)
	}

	return nil
}

// isValidStatusTransition checks if a status transition is valid
func isValidStatusTransition(from, to string) bool {
	validTransitions := map[string][]string{
		PlanStatusDraft:    {PlanStatusPending, PlanStatusCancelled},
		PlanStatusPending:  {PlanStatusApproved, PlanStatusRejected, PlanStatusCancelled},
		PlanStatusApproved: {PlanStatusActive, PlanStatusCancelled},
		PlanStatusActive:   {PlanStatusCompleted, PlanStatusCancelled},
	}

	allowed, exists := validTransitions[from]
	if !exists {
		return false
	}

	for _, status := range allowed {
		if status == to {
			return true
		}
	}

	return false
}

// GetPlanItemsProgress retrieves progress of items in a plan
func (s *Service) GetPlanItemsProgress(planID uint) ([]map[string]any, error) {
	// Get plan items
	var items []MaterialPlanItem
	if err := s.db.Where("plan_id = ?", planID).Find(&items).Error; err != nil {
		return nil, err
	}

	results := make([]map[string]any, 0, len(items))

	for _, item := range items {
		// Get material details
		var material struct {
			Code          string
			Name          string
			Specification string
			Unit          string
		}
		s.db.Table("material_master").Where("id = ?", item.MaterialID).
			Select("code, name, specification, unit").First(&material)

		// Get issued quantity from stock_logs (领料/发货)
		var issuedQty float64
		s.db.Table("stock_logs").
			Joins("JOIN stocks ON stocks.id = stock_logs.stock_id").
			Where("stocks.project_id IN (SELECT project_id FROM material_plans WHERE id = ?)", planID).
			Where("stock_logs.material_id = ?", item.MaterialID).
			Where("stock_logs.type = ?", "out").
			Where("stock_logs.source_type = ?", "requisition").
			Select("COALESCE(SUM(stock_logs.quantity), 0)").
			Scan(&issuedQty)

		// Get received quantity from stock_logs (入库/到货)
		var receivedQty float64
		s.db.Table("stock_logs").
			Joins("JOIN stocks ON stocks.id = stock_logs.stock_id").
			Where("stocks.project_id IN (SELECT project_id FROM material_plans WHERE id = ?)", planID).
			Where("stock_logs.material_id = ?", item.MaterialID).
			Where("stock_logs.type = ?", "in").
			Where("stock_logs.source_type = ?", "inbound").
			Select("COALESCE(SUM(stock_logs.quantity), 0)").
			Scan(&receivedQty)

		result := map[string]any{
			"id":                 item.ID,
			"plan_id":            item.PlanID,
			"material_id":        item.MaterialID,
			"material_code":      material.Code,
			"material_name":      material.Name,
			"specification":      material.Specification,
			"unit":               material.Unit,
			"planned_quantity":   item.PlannedQuantity,
			"received_quantity":  receivedQty,
			"issued_quantity":    issuedQty,
			"remaining_quantity": item.PlannedQuantity - issuedQty,
			"unit_price":         item.UnitPrice,
			"status":             item.Status,
			"progress":           0.0,
		}

		if item.PlannedQuantity > 0 {
			result["progress"] = (issuedQty / item.PlannedQuantity) * 100
		}

		results = append(results, result)
	}

	return results, nil
}
