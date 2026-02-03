package material

import (
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

// ListQueryParams holds parameters for listing materials
type ListQueryParams struct {
	Page            int
	PageSize        int
	Search          string
	Name            string
	Material        string
	Spec            string
	Specification   string
	Category        string
	ProjectID       uint
	ProjectIDs      []uint
	Filter          string // "unstored" to filter fully stored materials
	IncludeChildren bool
}

// MaterialWithEnrichment holds material with enriched data
type MaterialWithEnrichment struct {
	ID               uint     `gorm:"column:id"`
	Code             *string  `gorm:"column:code"`
	Name             string   `gorm:"column:name"`
	Specification    string   `gorm:"column:specification"`
	Unit             string   `gorm:"column:unit"`
	Price            float64  `gorm:"column:price"`
	Description      string   `gorm:"column:description"`
	Category         string   `gorm:"column:category"`
	Quantity         int      `gorm:"column:quantity"`
	StockQuantity    float64  `gorm:"column:stock_quantity"`
	ProjectID        *uint    `gorm:"column:project_id"`
	Material         string   `gorm:"column:material"`
	Spec             string   `gorm:"column:spec"`
	ProjectName      *string  `gorm:"column:project_name"`
	PlannedQuantity  int      `gorm:"column:planned_quantity"`
	ArrivedQuantity  int      `gorm:"column:arrived_quantity"`
}

// ListMaterials retrieves materials with filters and pagination
func (s *Service) ListMaterials(params ListQueryParams) ([]MaterialWithEnrichment, int64, error) {
	var results []MaterialWithEnrichment
	var total int64

	query := s.db.Model(&Material{})

	// Apply search filter
	if params.Search != "" {
		query = query.Where(
			"materials.name LIKE ? OR materials.material LIKE ? OR materials.spec LIKE ? OR materials.specification LIKE ? OR materials.category LIKE ?",
			"%"+params.Search+"%", "%"+params.Search+"%", "%"+params.Search+"%", "%"+params.Search+"%", "%"+params.Search+"%",
		)
	}

	// Apply field filters
	if params.Name != "" {
		query = query.Where("materials.name LIKE ?", "%"+params.Name+"%")
	}
	if params.Material != "" {
		query = query.Where("materials.material LIKE ?", "%"+params.Material+"%")
	}
	if params.Spec != "" {
		query = query.Where("materials.spec LIKE ? OR materials.specification LIKE ?", "%"+params.Spec+"%", "%"+params.Spec+"%")
	}
	if params.Specification != "" {
		query = query.Where("materials.specification LIKE ? OR materials.spec LIKE ?", "%"+params.Specification+"%", "%"+params.Specification+"%")
	}
	if params.Category != "" {
		query = query.Where("materials.category LIKE ?", "%"+params.Category+"%")
	}

	// Apply project filter
	if len(params.ProjectIDs) > 0 {
		query = query.Where("materials.project_id IN ?", params.ProjectIDs)
	} else if params.ProjectID > 0 {
		query = query.Where("materials.project_id = ?", params.ProjectID)
	}

	// Filter unstored materials
	if params.Filter == "unstored" {
		query = query.Where(`
			NOT EXISTS (
				SELECT 1
				FROM inbound_order_items ioi
				INNER JOIN inbound_orders io ON ioi.order_id = io.id
				WHERE ioi.material_id = materials.id
				AND io.status IN ('approved', 'completed')
				GROUP BY ioi.material_id
				HAVING SUM(ioi.quantity) >= materials.quantity
			)
		`)
	}

	// Get total count
	query.Count(&total)

	// Execute query with joins
	offset := (params.Page - 1) * params.PageSize
	s.db.Table("materials").
		Select("materials.*, projects.name as project_name, COALESCE(stocks.quantity, 0) as stock_quantity").
		Joins("LEFT JOIN projects ON materials.project_id = projects.id").
		Joins("LEFT JOIN stocks ON stocks.material_id = materials.id").
		Where(query).
		Offset(offset).Limit(params.PageSize).Order("materials.id DESC").
		Scan(&results)

	return results, total, nil
}

// EnrichWithPlanInfo enriches materials with plan information
func (s *Service) EnrichWithPlanInfo(materials []MaterialWithEnrichment) map[uint]PlanInfo {
	materialIDs := make([]uint, len(materials))
	for i, m := range materials {
		materialIDs[i] = m.ID
	}

	var planInfos []PlanInfo
	s.db.Raw(`
		SELECT
			material_id,
			COALESCE(SUM(planned_quantity), 0) as planned_quantity,
			COALESCE(SUM(arrived_quantity), 0) as arrived_quantity
		FROM material_plan_items
		WHERE material_id = ANY($1)
		GROUP BY material_id
	`, materialIDs).Scan(&planInfos)

	planMap := make(map[uint]PlanInfo)
	for _, info := range planInfos {
		planMap[info.MaterialID] = info
	}

	return planMap
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
	Price         float64
	Description   string
	Category      string
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

	// Validate project exists
	var projectExists int64
	if s.db.Table("projects").Where("id = ?", req.ProjectID).Count(&projectExists); projectExists == 0 {
		return nil, errors.New("指定的项目不存在")
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

	projectIDPtr := &req.ProjectID

	m := &Material{
		Code:          code,
		Name:          req.Name,
		Specification: req.Specification,
		Unit:          req.Unit,
		Price:         req.Price,
		Description:   req.Description,
		Category:      req.Category,
		Quantity:      req.Quantity,
		ProjectID:     projectIDPtr,
		Material:      req.Material,
		Spec:          req.Spec,
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
