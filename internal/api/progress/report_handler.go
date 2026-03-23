package progress

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ReportHandler handles report-related HTTP requests
type ReportHandler struct {
	db *gorm.DB
}

// NewReportHandler creates a new report handler
func NewReportHandler(db *gorm.DB) *ReportHandler {
	return &ReportHandler{db: db}
}

// GenerateReport generates a report based on the request parameters
// @Summary Generate Report
// @Tags reports
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param report body ReportRequest true "Report configuration"
// @Success 200 {object} ReportData
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /progress/project/:id/reports/generate [post]
func (h *ReportHandler) GenerateReport(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req ReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ProjectID = uint(projectID)

	// Set default date range if not provided
	if req.StartDate.IsZero() {
		req.StartDate = time.Now().AddDate(0, -1, 0) // 1 month ago
	}
	if req.EndDate.IsZero() {
		req.EndDate = time.Now()
	}

	// Set default sort order
	if req.SortOrder == "" {
		req.SortOrder = "asc"
	}

	// Generate report based on type
	reportData, err := h.generateReportData(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reportData)
}

// GetReports retrieves all saved reports for a project
// @Summary Get Project Reports
// @Tags reports
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {array} Report
// @Failure 500 {object} map[string]interface{}
// @Router /progress/project/:id/reports [get]
func (h *ReportHandler) GetReports(c *gin.Context) {
	projectID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var reports []Report
	err = h.db.Where("config_project_id = ?", uint(projectID)).
		Preload("User").
		Preload("Project").
		Find(&reports).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reports)
}

// GetReport retrieves a specific report by ID
// @Summary Get Report by ID
// @Tags reports
// @Produce json
// @Param id path int true "Report ID"
// @Success 200 {object} Report
// @Failure 404 {object} map[string]interface{}
// @Router /progress/reports/:id [get]
func (h *ReportHandler) GetReport(c *gin.Context) {
	reportID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}

	var report Report
	err = h.db.Preload("User").
		Preload("Project").
		First(&report, uint(reportID)).Error

	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Report not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, report)
}

// SaveReport saves a report configuration
// @Summary Save Report
// @Tags reports
// @Accept json
// @Produce json
// @Param report body ReportSaveRequest true "Report data"
// @Success 201 {object} Report
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /progress/reports [post]
func (h *ReportHandler) SaveReport(c *gin.Context) {
	var req ReportSaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	report := Report{
		Name:      req.Name,
		Config:    req.Config,
		CreatedBy: c.GetUint("user_id"), // Assuming user_id is set by auth middleware
	}

	if err := h.db.Create(&report).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Load associations
	h.db.Preload("User").Preload("Project").First(&report, report.ID)

	c.JSON(http.StatusCreated, report)
}

// UpdateReport updates an existing report
// @Summary Update Report
// @Tags reports
// @Accept json
// @Produce json
// @Param id path int true "Report ID"
// @Param report body ReportSaveRequest true "Report data"
// @Success 200 {object} Report
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /progress/reports/:id [put]
func (h *ReportHandler) UpdateReport(c *gin.Context) {
	reportID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}

	var req ReportSaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var report Report
	err = h.db.First(&report, uint(reportID)).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Report not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	report.Name = req.Name
	report.Config = req.Config

	if err := h.db.Save(&report).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Load associations
	h.db.Preload("User").Preload("Project").First(&report, report.ID)

	c.JSON(http.StatusOK, report)
}

// DeleteReport deletes a report
// @Summary Delete Report
// @Tags reports
// @Produce json
// @Param id path int true "Report ID"
// @Success 204
// @Failure 404 {object} map[string]interface{}
// @Router /progress/reports/:id [delete]
func (h *ReportHandler) DeleteReport(c *gin.Context) {
	reportID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}

	err = h.db.Delete(&Report{}, uint(reportID)).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ExportPDF exports a report as PDF
// @Summary Export Report to PDF
// @Tags reports
// @Accept json
// @Produce application/pdf
// @Param id path int true "Report ID"
// @Param export body ReportExportRequest true "Export options"
// @Success 200 {file} file
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /progress/reports/:id/export-pdf [post]
func (h *ReportHandler) ExportPDF(c *gin.Context) {
	reportID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}

	// Load report
	var report Report
	err = h.db.Preload("User").Preload("Project").First(&report, uint(reportID)).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Report not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generate PDF
	pdfData, err := h.generatePDF(&report)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename="+report.Name+".pdf")
	c.Data(http.StatusOK, "application/pdf", pdfData)
}

// ExportExcel exports a report as Excel
// @Summary Export Report to Excel
// @Tags reports
// @Accept json
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param id path int true "Report ID"
// @Param export body ReportExportRequest true "Export options"
// @Success 200 {file} file
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /progress/reports/:id/export-excel [post]
func (h *ReportHandler) ExportExcel(c *gin.Context) {
	reportID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}

	// Load report
	var report Report
	err = h.db.Preload("User").Preload("Project").First(&report, uint(reportID)).Error
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Report not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generate Excel
	excelData, err := h.generateExcel(&report)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename="+report.Name+".xlsx")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelData)
}

// GetTemplates retrieves all report templates
// @Summary Get Report Templates
// @Tags reports
// @Produce json
// @Success 200 {array} ReportTemplate
// @Router /progress/reports/templates [get]
func (h *ReportHandler) GetTemplates(c *gin.Context) {
	var templates []ReportTemplate
	err := h.db.Find(&templates).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, templates)
}

// SaveTemplate saves a report template
// @Summary Save Report Template
// @Tags reports
// @Accept json
// @Produce json
// @Param template body ReportTemplate true "Template data"
// @Success 201 {object} ReportTemplate
// @Router /progress/reports/templates [post]
func (h *ReportHandler) SaveTemplate(c *gin.Context) {
	var template ReportTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	template.CreatedBy = c.GetUint("user_id")

	if err := h.db.Create(&template).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, template)
}

// generateReportData generates the actual report data
func (h *ReportHandler) generateReportData(req *ReportRequest) (*ReportData, error) {
	reportData := &ReportData{
		Type:      req.Type,
		DateRange: []time.Time{req.StartDate, req.EndDate},
		Metadata: ReportMetadata{
			GeneratedAt: time.Now(),
			ProjectID:   req.ProjectID,
		},
	}

	// Set columns based on report type
	reportData.Columns = h.getColumnsForType(req.Type)

	// Fetch data based on report type
	var err error
	switch req.Type {
	case "task":
		reportData.Data, reportData.Summary, err = h.fetchTaskReportData(req)
	case "resource":
		reportData.Data, reportData.Summary, err = h.fetchResourceReportData(req)
	case "milestone":
		reportData.Data, reportData.Summary, err = h.fetchMilestoneReportData(req)
	case "progress":
		reportData.Data, reportData.Summary, err = h.fetchProgressReportData(req)
	default:
		return nil, fmt.Errorf("invalid report type: %s", req.Type)
	}

	if err != nil {
		return nil, err
	}

	// Apply grouping if specified
	if req.GroupBy != "" {
		reportData.Groups = h.groupData(reportData.Data, req.GroupBy)
	}

	// Generate title
	reportData.Title = h.generateReportTitle(req)

	return reportData, nil
}

// Helper methods for fetching report data
func (h *ReportHandler) fetchTaskReportData(req *ReportRequest) ([]map[string]interface{}, SummaryStats, error) {
	// Implementation for task report
	// Query tasks based on filters, date range, etc.
	return []map[string]interface{}{}, SummaryStats{}, nil
}

func (h *ReportHandler) fetchResourceReportData(req *ReportRequest) ([]map[string]interface{}, SummaryStats, error) {
	// Implementation for resource report
	return []map[string]interface{}{}, SummaryStats{}, nil
}

func (h *ReportHandler) fetchMilestoneReportData(req *ReportRequest) ([]map[string]interface{}, SummaryStats, error) {
	// Implementation for milestone report
	return []map[string]interface{}{}, SummaryStats{}, nil
}

func (h *ReportHandler) fetchProgressReportData(req *ReportRequest) ([]map[string]interface{}, SummaryStats, error) {
	// Implementation for progress report
	return []map[string]interface{}{}, SummaryStats{}, nil
}

func (h *ReportHandler) getColumnsForType(reportType string) []ColumnDef {
	// Return column definitions based on report type
	return []ColumnDef{}
}

func (h *ReportHandler) groupData(data []map[string]interface{}, groupBy string) []GroupData {
	// Implementation for grouping data
	return []GroupData{}
}

func (h *ReportHandler) generateReportTitle(req *ReportRequest) string {
	// Generate human-readable title
	return "Report"
}

func (h *ReportHandler) generatePDF(report *Report) ([]byte, error) {
	// Generate PDF using a library like github.com/jung-kurt/gofpdf
	// This is a placeholder - actual implementation would use a PDF library
	return []byte{}, nil
}

func (h *ReportHandler) generateExcel(report *Report) ([]byte, error) {
	// Generate Excel using a library like github.com/xuri/excelize
	// This is a placeholder - actual implementation would use an Excel library
	return []byte{}, nil
}
