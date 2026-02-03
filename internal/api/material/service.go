package material

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

// Service handles business logic for materials
type Service struct {
	db *gorm.DB
}

// NewService creates a new material service
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// ListQueryParams holds parameters for listing materials from plan items
type ListQueryParams struct {
	Page          int
	PageSize      int
	Search        string
	Name          string
	Specification string
	Category      string
	ProjectID     uint
	PlanID        uint
}

// PlanMaterialItem holds material item from plans
type PlanMaterialItem struct {
	ID               uint     `json:"id"`
	PlanID           uint     `json:"plan_id"`
	PlanNo           string   `json:"plan_no"`
	PlanName         string   `json:"plan_name"`
	ProjectID        uint     `json:"project_id"`
	ProjectName      string   `json:"project_name"`
	MaterialID       uint     `json:"material_id"`
	MaterialCode     string   `json:"material_code"`
	MaterialName     string   `json:"material_name"`
	Specification    string   `json:"specification"`
	Unit             string   `json:"unit"`
	Category         string   `json:"category"`
	PlannedQuantity  float64  `json:"planned_quantity"`
	UnitPrice        float64  `json:"unit_price"`
	ArrivedQuantity  float64  `json:"arrived_quantity"`
	RemainingQty     float64  `json:"remaining_quantity"`
	ArrivalPercent   float64  `json:"arrival_percent"`
	Priority         string   `json:"priority"`
	Status           string   `json:"status"`
	RequiredDate     *string  `json:"required_date"`
	Remark           string   `json:"remark"`
}

// ListMaterials retrieves plan material items with filters and pagination
func (s *Service) ListMaterials(params ListQueryParams) ([]PlanMaterialItem, int64, error) {
	var results []PlanMaterialItem
	var total int64

	// Build WHERE conditions
	conditions := []string{}
	args := []interface{}{}

	// Search filter (name, specification, code, category)
	if params.Search != "" {
		conditions = append(conditions, "(mm.name LIKE ? OR mm.specification LIKE ? OR mm.code LIKE ? OR mm.category LIKE ?)")
		args = append(args, "%"+params.Search+"%", "%"+params.Search+"%", "%"+params.Search+"%", "%"+params.Search+"%")
	}

	// Name filter
	if params.Name != "" {
		conditions = append(conditions, "mm.name LIKE ?")
		args = append(args, "%"+params.Name+"%")
	}

	// Specification filter
	if params.Specification != "" {
		conditions = append(conditions, "mm.specification LIKE ?")
		args = append(args, "%"+params.Specification+"%")
	}

	// Category filter
	if params.Category != "" {
		conditions = append(conditions, "mm.category LIKE ?")
		args = append(args, "%"+params.Category+"%")
	}

	// Project filter
	if params.ProjectID > 0 {
		conditions = append(conditions, "mp.project_id = ?")
		args = append(args, params.ProjectID)
	}

	// Plan filter
	if params.PlanID > 0 {
		conditions = append(conditions, "mpi.plan_id = ?")
		args = append(args, params.PlanID)
	}

	// Build WHERE clause
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + conditions[0]
		for i := 1; i < len(conditions); i++ {
			whereClause += " AND " + conditions[i]
		}
	}

	// Get total count
	countQuery := `
		SELECT COUNT(*)
		FROM material_plan_items mpi
		INNER JOIN material_plans mp ON mpi.plan_id = mp.id
		INNER JOIN projects p ON mp.project_id = p.id
		INNER JOIN material_master mm ON mpi.material_id = mm.id
	` + whereClause
	s.db.Raw(countQuery, args...).Scan(&total)

	// Main query
	offset := (params.Page - 1) * params.PageSize
	query := `
		SELECT
			mpi.id,
			mpi.plan_id,
			mp.plan_no,
			mp.plan_name,
			mp.project_id,
			p.name as project_name,
			mpi.material_id,
			mm.code as material_code,
			mm.name as material_name,
			mm.specification,
			mm.unit,
			mm.category,
			mpi.planned_quantity,
			COALESCE(mpi.unit_price, 0) as unit_price,
			mpi.priority,
			mpi.status,
			mpi.required_date,
			mpi.remark
		FROM material_plan_items mpi
		INNER JOIN material_plans mp ON mpi.plan_id = mp.id
		INNER JOIN projects p ON mp.project_id = p.id
		INNER JOIN material_master mm ON mpi.material_id = mm.id
	` + whereClause + `
		ORDER BY mpi.id DESC
		LIMIT ? OFFSET ?
	`
	queryArgs := append(args, params.PageSize, offset)
	rows, err := s.db.Raw(query, queryArgs...).Rows()
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Scan results
	for rows.Next() {
		var item PlanMaterialItem
		var requiredDate sql.NullString
		var unitPrice sql.NullFloat64

		err := rows.Scan(
			&item.ID,
			&item.PlanID,
			&item.PlanNo,
			&item.PlanName,
			&item.ProjectID,
			&item.ProjectName,
			&item.MaterialID,
			&item.MaterialCode,
			&item.MaterialName,
			&item.Specification,
			&item.Unit,
			&item.Category,
			&item.PlannedQuantity,
			&unitPrice,
			&item.Priority,
			&item.Status,
			&requiredDate,
			&item.Remark,
		)
		if err != nil {
			return nil, 0, err
		}

		if unitPrice.Valid {
			item.UnitPrice = unitPrice.Float64
		}
		if requiredDate.Valid {
			item.RequiredDate = &requiredDate.String
		}

		// Calculate arrived quantity from inbound_items
		var arrivedQty float64
		s.db.Raw(`
			SELECT COALESCE(SUM(ii.quantity), 0)
			FROM inbound_items ii
			INNER JOIN inbound_orders io ON ii.inbound_order_id = io.id
			WHERE ii.material_id = ?
			AND io.project_id = ?
			AND io.status IN ('approved', 'completed')
		`, item.MaterialID, item.ProjectID).Scan(&arrivedQty)

		item.ArrivedQuantity = arrivedQty
		item.RemainingQty = item.PlannedQuantity - arrivedQty
		if item.PlannedQuantity > 0 {
			item.ArrivalPercent = (arrivedQty / item.PlannedQuantity) * 100
		}

		results = append(results, item)
	}

	return results, total, nil
}

// GetMaterial retrieves a material by ID
func (s *Service) GetMaterial(id uint) (*Material, error) {
	var m Material
	if err := s.db.First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("物资不存在")
		}
		return nil, err
	}
	return &m, nil
}

// CreateMaterialRequest holds parameters for creating a material
type CreateMaterialRequest struct {
	Code          string
	Name          string
	Specification string
	Unit          string
	Description   string
	Category      string
	// These fields are not in material_master table but kept for API compatibility
	Price         float64
	Quantity      int
	ProjectID     uint
	Material      string
	Spec          string
}

// CreateMaterial creates a new material
func (s *Service) CreateMaterial(req *CreateMaterialRequest) (*Material, error) {
	// Validate name
	if req.Name == "" {
		return nil, errors.New("物资名称不能为空")
	}

	// Check for duplicate code
	var code *string
	if req.Code != "" {
		code = &req.Code
		var existing Material
		if s.db.Where("code = ?", req.Code).First(&existing).Error == nil {
			return nil, errors.New("物资编码已存在")
		}
	}

	m := &Material{
		Code:          code,
		Name:          req.Name,
		Specification: req.Specification,
		Unit:          req.Unit,
		Description:   req.Description,
		Category:      req.Category,
	}

	if err := s.db.Create(m).Error; err != nil {
		return nil, fmt.Errorf("创建物资失败: %w", err)
	}

	return m, nil
}

// UpdateMaterial updates an existing material
func (s *Service) UpdateMaterial(id uint, updates map[string]any) (*Material, error) {
	var m Material
	if err := s.db.First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("物资不存在")
		}
		return nil, err
	}

	// Handle code update
	if v, ok := updates["code"].(string); ok {
		var codePtr *string
		if v != "" {
			codePtr = &v
		}
		if (m.Code == nil && codePtr != nil) || (m.Code != nil && codePtr == nil) || (m.Code != nil && codePtr != nil && *m.Code != *codePtr) {
			var existing Material
			if s.db.Where("code = ?", v).First(&existing).Error == nil && existing.ID != m.ID {
				return nil, errors.New("物资编码已存在")
			}
			m.Code = codePtr
		}
		delete(updates, "code")
	}

	// Handle project_id update
	if v, ok := updates["project_id"]; ok {
		var projectID uint
		switch val := v.(type) {
		case string:
			if val != "" {
				pid, err := strconv.ParseUint(val, 10, 64)
				if err != nil || pid == 0 {
					return nil, errors.New("项目ID格式无效")
				}
				projectID = uint(pid)
			}
		case float64:
			if val > 0 {
				projectID = uint(val)
			}
		}

		if projectID > 0 {
			var projectExists int64
			if s.db.Table("projects").Where("id = ?", projectID).Count(&projectExists); projectExists == 0 {
				return nil, errors.New("指定的项目不存在")
			}
			m.ProjectID = &projectID
		} else {
			m.ProjectID = nil
		}
		delete(updates, "project_id")
	}

	// Apply other updates
	if err := s.db.Model(&m).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("更新物资失败: %w", err)
	}

	// Reload to get updated values
	s.db.First(&m, id)

	return &m, nil
}

// DeleteMaterial deletes a material
func (s *Service) DeleteMaterial(id uint) error {
	var m Material
	if err := s.db.First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("物资不存在")
		}
		return err
	}

	if err := s.db.Delete(&m).Error; err != nil {
		return fmt.Errorf("删除物资失败: %w", err)
	}

	return nil
}

// GetMaterialLogs retrieves operation logs for a material
func (s *Service) GetMaterialLogs(id uint) ([]map[string]any, error) {
	// Verify material exists
	var m Material
	if err := s.db.First(&m, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("物资不存在")
		}
		return nil, err
	}

	type Log struct {
		ID          uint   `json:"id"`
		Action      string `json:"action"`
		Description string `json:"description"`
		CreatedAt   string `json:"created_at"`
		CreatedBy   string `json:"created_by"`
	}

	var logs []Log
	s.db.Raw(`
		SELECT l.id, l.operation_type as action, l.remark as description,
		       l.created_at, u.username as created_by
		FROM stock_logs l
		LEFT JOIN stocks s ON l.stock_id = s.id
		LEFT JOIN users u ON l.created_by = u.id
		WHERE s.material_id = ?
		ORDER BY l.created_at DESC
		LIMIT 100
	`, id).Scan(&logs)

	result := make([]map[string]any, len(logs))
	for i, log := range logs {
		result[i] = map[string]any{
			"id":          log.ID,
			"action":      log.Action,
			"description": log.Description,
			"created_at":  log.CreatedAt,
			"created_by":  log.CreatedBy,
		}
	}

	return result, nil
}

// ImportMaterial holds data for importing a material
type ImportMaterial struct {
	Code          string
	Name          string
	Category      string
	Specification string
	Unit          string
	Price         float64
	Quantity      int
	Description   string
	ProjectID     uint
	Material      string
	Spec          string
}

// ImportResult holds the result of a batch import
type ImportResult struct {
	SuccessCount int      `json:"success_count"`
	ErrorCount   int      `json:"error_count"`
	Errors       []string `json:"errors"`
}

// ImportMaterials imports materials in batch
func (s *Service) ImportMaterials(projectID uint, materials []ImportMaterial) (*ImportResult, error) {
	// Validate project
	if projectID == 0 {
		return nil, errors.New("请指定所属项目")
	}
	var projectExists int64
	if s.db.Table("projects").Where("id = ?", projectID).Count(&projectExists); projectExists == 0 {
		return nil, errors.New("指定的项目不存在")
	}

	result := &ImportResult{
		Errors: []string{},
	}

	for i, item := range materials {
		name := item.Name
		if name == "" {
			result.ErrorCount++
			result.Errors = append(result.Errors, fmt.Sprintf("第%d行: 物资名称不能为空", i+1))
			continue
		}

		m := &Material{
			Name:      name,
			ProjectID: &projectID,
		}

		if item.Code != "" {
			m.Code = &item.Code
			var existing Material
			if s.db.Where("code = ?", item.Code).First(&existing).Error == nil {
				result.ErrorCount++
				result.Errors = append(result.Errors, fmt.Sprintf("第%d行: 物资编码已存在 - %s", i+1, item.Code))
				continue
			}
		}

		m.Specification = item.Specification
		m.Unit = item.Unit
		m.Price = item.Price
		m.Description = item.Description
		m.Category = item.Category
		m.Quantity = item.Quantity
		m.Material = item.Material
		m.Spec = item.Spec

		if err := s.db.Create(m).Error; err != nil {
			result.ErrorCount++
			result.Errors = append(result.Errors, fmt.Sprintf("第%d行: %s", i+1, err.Error()))
		} else {
			result.SuccessCount++
		}
	}

	return result, nil
}

// BatchCreateMaterial holds data for batch creating materials
type BatchCreateMaterial struct {
	Name          string
	Code          string
	Specification string
	Category      string
	Unit          string
	Price         *float64
	Quantity      *float64
	ProjectID     *uint
}

// BatchCreateResult holds the result of batch creation
type BatchCreateResult struct {
	Total     int                    `json:"total"`
	Success   int                    `json:"success"`
	Failed    int                    `json:"failed"`
	Materials []map[string]any       `json:"materials"`
	FailedItems []map[string]any     `json:"failed_items,omitempty"`
}

// BatchCreateMaterials creates or finds materials in batch
func (s *Service) BatchCreateMaterials(materials []BatchCreateMaterial) (*BatchCreateResult, error) {
	result := &BatchCreateResult{
		Materials: []map[string]any{},
		FailedItems: []map[string]any{},
	}

	for _, mat := range materials {
		// Check if material exists by name and specification
		var existingMaterial Material
		err := s.db.Where("name = ? AND (specification = '' OR specification IS NULL OR specification = ?)",
			mat.Name, mat.Specification).
			First(&existingMaterial).Error

		var materialID uint
		var isNew bool

		if err == nil {
			materialID = existingMaterial.ID
			isNew = false
		} else {
			price := 0.0
			if mat.Price != nil {
				price = *mat.Price
			}

			quantity := 0
			if mat.Quantity != nil {
				quantity = int(*mat.Quantity)
			}

			var codePtr *string
			if mat.Code != "" {
				codePtr = &mat.Code
			}

			material := Material{
				Code:          codePtr,
				Name:          mat.Name,
				Category:      mat.Category,
				Specification: mat.Specification,
				Unit:          mat.Unit,
				Price:         price,
				Quantity:      quantity,
				ProjectID:     mat.ProjectID,
			}

			if err := s.db.Create(&material).Error; err != nil {
				result.FailedItems = append(result.FailedItems, map[string]any{
					"name":  mat.Name,
					"error": err.Error(),
				})
				result.Failed++
				continue
			}

			materialID = material.ID
			isNew = true
		}

		result.Materials = append(result.Materials, map[string]any{
			"id":           materialID,
			"name":         mat.Name,
			"specification": mat.Specification,
			"is_new":       isNew,
		})
		result.Success++
	}

	result.Total = len(materials)

	return result, nil
}

// PlanInfo holds plan information for a material
type PlanInfo struct {
	MaterialID      uint
	PlannedQuantity int
	ArrivedQuantity int
}

// GetChildProjectIDs recursively gets all child project IDs
func (s *Service) GetChildProjectIDs(parentID uint) []uint {
	var children []struct {
		ID uint
	}
	s.db.Table("projects").Where("parent_id = ?", parentID).Find(&children)

	var ids []uint
	for _, child := range children {
		ids = append(ids, child.ID)
		childIDs := s.GetChildProjectIDs(child.ID)
		ids = append(ids, childIDs...)
	}

	return ids
}

// ExportMaterials retrieves materials for export
func (s *Service) ExportMaterials(search string, projectIDs []uint) ([]Material, error) {
	var materials []Material
	query := s.db.Model(&Material{})

	if search != "" {
		query = query.Where(
			"name LIKE ? OR material LIKE ? OR spec LIKE ? OR specification LIKE ? OR category LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%",
		)
	}

	if len(projectIDs) > 0 {
		query = query.Where("project_id IN ?", projectIDs)
	}

	if err := query.Order("id DESC").Find(&materials).Error; err != nil {
		return nil, err
	}

	return materials, nil
}
