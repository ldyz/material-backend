package material_master

import (
	"fmt"

	"gorm.io/gorm"
)

// Service 物资主数据服务
type Service struct {
	db *gorm.DB
}

// NewService 创建物资主数据服务
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// Create 创建物资主数据
func (s *Service) Create(req *CreateMaterialMasterRequest) (*MaterialMaster, error) {
	// 检查编码是否已存在
	var count int64
	if err := s.db.Model(&MaterialMaster{}).Where("code = ?", req.Code).Count(&count).Error; err != nil {
		return nil, fmt.Errorf("检查编码失败: %w", err)
	}
	if count > 0 {
		return nil, fmt.Errorf("物资编码 %s 已存在", req.Code)
	}

	// 检查名称+规格是否已存在
	var existing MaterialMaster
	err := s.db.Where("name = ? AND COALESCE(specification, '') = ?", req.Name, req.Specification).
		First(&existing).Error
	if err == nil {
		return nil, fmt.Errorf("物资 %s (%s) 已存在", req.Name, req.Specification)
	}
	if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("检查物资失败: %w", err)
	}

	material := &MaterialMaster{
		Code:         req.Code,
		Name:         req.Name,
		Specification: req.Specification,
		Unit:         req.Unit,
		Category:     req.Category,
		SafetyStock:  req.SafetyStock,
		Description:  req.Description,
	}

	if err := s.db.Create(material).Error; err != nil {
		return nil, fmt.Errorf("创建物资主数据失败: %w", err)
	}

	return material, nil
}

// Update 更新物资主数据
func (s *Service) Update(id uint, req *UpdateMaterialMasterRequest) (*MaterialMaster, error) {
	var material MaterialMaster
	if err := s.db.First(&material, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("物资主数据不存在")
		}
		return nil, fmt.Errorf("查询物资主数据失败: %w", err)
	}

	// 检查编码是否与其他记录冲突
	var count int64
	if err := s.db.Model(&MaterialMaster{}).
		Where("code = ? AND id != ?", req.Code, id).
		Count(&count).Error; err != nil {
		return nil, fmt.Errorf("检查编码失败: %w", err)
	}
	if count > 0 {
		return nil, fmt.Errorf("物资编码 %s 已被其他物资使用", req.Code)
	}

	// 检查名称+规格是否与其他记录冲突
	var existing MaterialMaster
	err := s.db.Where("name = ? AND COALESCE(specification, '') = ? AND id != ?",
		req.Name, req.Specification, id).
		First(&existing).Error
	if err == nil {
		return nil, fmt.Errorf("物资 %s (%s) 已存在", req.Name, req.Specification)
	}
	if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("检查物资失败: %w", err)
	}

	// 更新字段
	material.Code = req.Code
	material.Name = req.Name
	material.Specification = req.Specification
	material.Unit = req.Unit
	material.Category = req.Category
	material.SafetyStock = req.SafetyStock
	material.Description = req.Description

	if err := s.db.Save(&material).Error; err != nil {
		return nil, fmt.Errorf("更新物资主数据失败: %w", err)
	}

	return &material, nil
}

// Delete 删除物资主数据
func (s *Service) Delete(id uint) error {
	// 检查是否有关联的库存记录
	var count int64
	if err := s.db.Table("stocks").Where("material_id = ?", id).Count(&count).Error; err != nil {
		return fmt.Errorf("检查库存失败: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("该物资存在库存记录，无法删除")
	}

	// 检查是否有关联的计划明细
	if err := s.db.Table("material_plan_items").Where("material_id = ?", id).Count(&count).Error; err != nil {
		return fmt.Errorf("检查计划明细失败: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("该物资存在计划明细，无法删除")
	}

	// 删除物资主数据
	if err := s.db.Delete(&MaterialMaster{}, id).Error; err != nil {
		return fmt.Errorf("删除物资主数据失败: %w", err)
	}

	return nil
}

// GetByID 根据 ID 获取物资主数据
func (s *Service) GetByID(id uint) (*MaterialMaster, error) {
	var material MaterialMaster
	if err := s.db.First(&material, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("物资主数据不存在")
		}
		return nil, fmt.Errorf("查询物资主数据失败: %w", err)
	}
	return &material, nil
}

// GetByCode 根据编码获取物资主数据
func (s *Service) GetByCode(code string) (*MaterialMaster, error) {
	var material MaterialMaster
	if err := s.db.Where("code = ?", code).First(&material).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("物资主数据不存在")
		}
		return nil, fmt.Errorf("查询物资主数据失败: %w", err)
	}
	return &material, nil
}

// List 获取物资主数据列表
func (s *Service) List(page, pageSize int, keyword, category string) ([]MaterialMaster, int64, error) {
	var materials []MaterialMaster
	var total int64

	query := s.db.Model(&MaterialMaster{})

	// 关键词搜索
	if keyword != "" {
		query = query.Where("code LIKE ? OR name LIKE ? OR specification LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 分类筛选
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计物资数量失败: %w", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Find(&materials).Error; err != nil {
		return nil, 0, fmt.Errorf("查询物资列表失败: %w", err)
	}

	return materials, total, nil
}

// ListWithProjectStock 获取指定项目的物资列表（带库存信息）
func (s *Service) ListWithProjectStock(projectID uint, page, pageSize int, keyword, category string) ([]MaterialMasterQueryDTO, int64, error) {
	var results []MaterialMasterQueryDTO
	var total int64

	// 构建查询：物资主数据 LEFT JOIN 库存表
	query := s.db.Table("material_master m").
		Select(`m.id, m.code, m.name, m.specification, m.unit, m.category,
			m.safety_stock, m.description, m.created_at, m.updated_at,
			? as project_id, s.id as stock_id, s.quantity, s.location, s.unit_cost`, projectID).
		Joins("LEFT JOIN stocks s ON s.material_id = m.id AND s.project_id = ?", projectID)

	// 关键词搜索
	if keyword != "" {
		query = query.Where("m.code LIKE ? OR m.name LIKE ? OR m.specification LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 分类筛选
	if category != "" {
		query = query.Where("m.category = ?", category)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计物资数量失败: %w", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).
		Order("m.created_at DESC").
		Find(&results).Error; err != nil {
		return nil, 0, fmt.Errorf("查询物资列表失败: %w", err)
	}

	return results, total, nil
}

// GetCategories 获取所有物资分类
func (s *Service) GetCategories() ([]string, error) {
	var categories []string

	if err := s.db.Model(&MaterialMaster{}).
		Distinct("category").
		Where("category != ''").
		Pluck("category", &categories).Error; err != nil {
		return nil, fmt.Errorf("获取物资分类失败: %w", err)
	}

	return categories, nil
}

// SyncOrCreate 同步或创建物资主数据
// 用于数据迁移或从外部系统同步物资数据
func (s *Service) SyncOrCreate(code, name, specification string) (*MaterialMaster, error) {
	var material MaterialMaster

	// 先尝试通过编码查找
	err := s.db.Where("code = ?", code).First(&material).Error
	if err == nil {
		return &material, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("查询物资主数据失败: %w", err)
	}

	// 不存在则创建
	material = MaterialMaster{
		Code:         code,
		Name:         name,
		Specification: specification,
	}

	if err := s.db.Create(&material).Error; err != nil {
		return nil, fmt.Errorf("创建物资主数据失败: %w", err)
	}

	return &material, nil
}
