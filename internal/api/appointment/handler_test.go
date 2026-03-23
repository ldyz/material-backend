package appointment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	jwtpkg "github.com/yourorg/material-backend/backend/pkg/jwt"
	"github.com/yourorg/material-backend/backend/internal/api/auth"
	"github.com/yourorg/material-backend/backend/internal/api/project"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupAppointmentRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db fail: %v", err)
	}
	db.AutoMigrate(&auth.User{}, &auth.Role{}, &project.Project{}, &ConstructionAppointment{}, &WorkerCalendar{})

	// Create admin role and user
	r := auth.Role{Name: "admin", Permissions: "admin"}
	db.Create(&r)
	u := auth.User{Username: "admin", Email: "admin@example.com", Role: "admin", IsActive: true}
	u.SetPassword("admin")
	db.Create(&u)

	// Create worker user
	worker := auth.User{Username: "worker1", Email: "worker1@example.com", Role: "user", FullName: "张三", IsActive: true}
	worker.SetPassword("worker1")
	db.Create(&worker)

	// Create project
	p := project.Project{Name: "测试项目", Code: "TEST-001"}
	db.Create(&p)

	g := gin.Default()
	rg := g.Group("/api")
	auth.RegisterRoutes(rg, db)
	RegisterRoutes(rg, db)
	return g, db
}

func TestAppointmentCRUD(t *testing.T) {
	g, db := setupAppointmentRouter(t)
	admin := auth.User{}
	db.Where("username = ?", "admin").First(&admin)
	tok, _ := jwtpkg.GenerateToken(admin.ID, admin.Username)

	// Get project
	var project project.Project
	db.Where("code = ?", "TEST-001").First(&project)

	// Get worker
	var worker auth.User
	db.Where("username = ?", "worker1").First(&worker)

	// Create appointment
	workDate := time.Now().Add(24 * time.Hour).Format("2006-01-02")
	// assigned_worker_ids 应该是 JSON 字符串格式
	workerIDsJSON, _ := json.Marshal([]uint{worker.ID})
	createReq := map[string]any{
		"project_id":          project.ID,
		"contact_person":      "张三",
		"contact_phone":       "13800138000",
		"work_date":           workDate,
		"time_slot":           "morning",
		"work_location":       "A区施工点",
		"work_content":        "浇筑混凝土",
		"work_type":           "施工",
		"assigned_worker_ids": string(workerIDsJSON),
	}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/appointments", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+tok)
	req.Header.Set("Content-Type", "application/json")
	r := httptest.NewRecorder()
	g.ServeHTTP(r, req)
	if r.Code != http.StatusOK && r.Code != http.StatusCreated {
		t.Fatalf("create failed: %d %s", r.Code, r.Body.String())
	}
	var resp map[string]any
	json.Unmarshal(r.Body.Bytes(), &resp)
	data := resp["data"].(map[string]any)
	appointmentID := int(data["id"].(float64))

	// List appointments
	req2 := httptest.NewRequest("GET", "/api/appointments", nil)
	req2.Header.Set("Authorization", "Bearer "+tok)
	r2 := httptest.NewRecorder()
	g.ServeHTTP(r2, req2)
	if r2.Code != http.StatusOK {
		t.Fatalf("list failed: %d", r2.Code)
	}

	// Get appointment by ID
	req3 := httptest.NewRequest("GET", "/api/appointments/"+jsonNumberToString(appointmentID), nil)
	req3.Header.Set("Authorization", "Bearer "+tok)
	r3 := httptest.NewRecorder()
	g.ServeHTTP(r3, req3)
	if r3.Code != http.StatusOK {
		t.Fatalf("get failed: %d", r3.Code)
	}

	// Update appointment
	updateReq := map[string]any{
		"work_content": "修改后的工作内容",
	}
	updateBody, _ := json.Marshal(updateReq)
	req4 := httptest.NewRequest("PUT", "/api/appointments/"+jsonNumberToString(appointmentID), bytes.NewReader(updateBody))
	req4.Header.Set("Authorization", "Bearer "+tok)
	req4.Header.Set("Content-Type", "application/json")
	r4 := httptest.NewRecorder()
	g.ServeHTTP(r4, req4)
	if r4.Code != http.StatusOK {
		t.Fatalf("update failed: %d %s", r4.Code, r4.Body.String())
	}

	// Delete appointment
	req5 := httptest.NewRequest("DELETE", "/api/appointments/"+jsonNumberToString(appointmentID), nil)
	req5.Header.Set("Authorization", "Bearer "+tok)
	r5 := httptest.NewRecorder()
	g.ServeHTTP(r5, req5)
	if r5.Code != http.StatusOK {
		t.Fatalf("delete failed: %d", r5.Code)
	}

	// Verify deleted
	req6 := httptest.NewRequest("GET", "/api/appointments/"+jsonNumberToString(appointmentID), nil)
	req6.Header.Set("Authorization", "Bearer "+tok)
	r6 := httptest.NewRecorder()
	g.ServeHTTP(r6, req6)
	if r6.Code != http.StatusNotFound {
		t.Fatalf("expected 404 after delete, got: %d", r6.Code)
	}
}

func TestAppointmentStats(t *testing.T) {
	g, db := setupAppointmentRouter(t)
	admin := auth.User{}
	db.Where("username = ?", "admin").First(&admin)
	tok, _ := jwtpkg.GenerateToken(admin.ID, admin.Username)

	// Get project
	var project project.Project
	db.Where("code = ?", "TEST-001").First(&project)

	// Create some test appointments
	workDate := time.Now().Add(24 * time.Hour)
	appointment1 := ConstructionAppointment{
		AppointmentNo:   "AP001",
		ProjectID:       &project.ID,
		ApplicantID:     admin.ID,
		ApplicantName:   "Admin",
		ContactPerson:   "张三",
		ContactPhone:    "13800138000",
		WorkDate:        workDate,
		TimeSlot:        "上午",
		WorkLocation:    "A区",
		WorkContent:     "测试内容1",
		Status:          "draft",
	}
	appointment2 := ConstructionAppointment{
		AppointmentNo:   "AP002",
		ProjectID:       &project.ID,
		ApplicantID:     admin.ID,
		ApplicantName:   "Admin",
		ContactPerson:   "李四",
		ContactPhone:    "13800138001",
		WorkDate:        workDate,
		TimeSlot:        "下午",
		WorkLocation:    "B区",
		WorkContent:     "测试内容2",
		Status:          "pending",
	}
	db.Create(&appointment1)
	db.Create(&appointment2)

	// Get stats
	req := httptest.NewRequest("GET", "/api/appointments/stats", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	r := httptest.NewRecorder()
	g.ServeHTTP(r, req)
	if r.Code != http.StatusOK {
		t.Fatalf("stats failed: %d %s", r.Code, r.Body.String())
	}
	var resp map[string]any
	json.Unmarshal(r.Body.Bytes(), &resp)
	if resp["data"] == nil {
		t.Fatalf("stats response missing data field")
	}
}

func TestAppointmentSearch(t *testing.T) {
	g, db := setupAppointmentRouter(t)
	admin := auth.User{}
	db.Where("username = ?", "admin").First(&admin)
	tok, _ := jwtpkg.GenerateToken(admin.ID, admin.Username)

	// Get project
	var project project.Project
	db.Where("code = ?", "TEST-001").First(&project)

	// Create test appointment
	workDate := time.Now().Add(24 * time.Hour)
	appointment := ConstructionAppointment{
		AppointmentNo:   "AP003",
		ProjectID:       &project.ID,
		ApplicantID:     admin.ID,
		ApplicantName:   "Admin",
		ContactPerson:   "王五",
		ContactPhone:    "13800138002",
		WorkDate:        workDate,
		TimeSlot:        "上午",
		WorkLocation:    "C区",
		WorkContent:     "特殊施工任务",
		Status:          "draft",
	}
	db.Create(&appointment)

	// Search by keyword
	req := httptest.NewRequest("GET", "/api/appointments/search?keyword=特殊", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	r := httptest.NewRecorder()
	g.ServeHTTP(r, req)
	if r.Code != http.StatusOK {
		t.Fatalf("search failed: %d %s", r.Code, r.Body.String())
	}
	var resp map[string]any
	json.Unmarshal(r.Body.Bytes(), &resp)
	data := resp["data"].([]any)
	if len(data) == 0 {
		t.Fatalf("search should return at least 1 result")
	}
}

func TestAppointmentFilters(t *testing.T) {
	g, db := setupAppointmentRouter(t)
	admin := auth.User{}
	db.Where("username = ?", "admin").First(&admin)
	tok, _ := jwtpkg.GenerateToken(admin.ID, admin.Username)

	// Get project
	var project project.Project
	db.Where("code = ?", "TEST-001").First(&project)

	// Delete all existing appointments from other tests
	db.Exec("DELETE FROM construction_appointments")

	// Create test appointments with different statuses
	workDate := time.Now().Add(24 * time.Hour)
	appointment1 := ConstructionAppointment{
		AppointmentNo:   "AP004",
		ProjectID:       &project.ID,
		ApplicantID:     admin.ID,
		ApplicantName:   "Admin",
		ContactPerson:   "测试人1",
		ContactPhone:    "13800138003",
		WorkDate:        workDate,
		TimeSlot:        "上午",
		WorkLocation:    "D区",
		WorkContent:     "测试内容",
		Status:          "draft",
	}
	appointment2 := ConstructionAppointment{
		AppointmentNo:   "AP005",
		ProjectID:       &project.ID,
		ApplicantID:     admin.ID,
		ApplicantName:   "Admin",
		ContactPerson:   "测试人2",
		ContactPhone:    "13800138004",
		WorkDate:        workDate,
		TimeSlot:        "下午",
		WorkLocation:    "E区",
		WorkContent:     "测试内容2",
		Status:          "pending",
	}
	db.Create(&appointment1)
	db.Create(&appointment2)

	// Filter by status
	req := httptest.NewRequest("GET", "/api/appointments?status=draft", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	r := httptest.NewRecorder()
	g.ServeHTTP(r, req)
	if r.Code != http.StatusOK {
		t.Fatalf("filter failed: %d - %s", r.Code, r.Body.String())
	}
	var resp map[string]any
	json.Unmarshal(r.Body.Bytes(), &resp)
	data := resp["data"].([]any)
	if len(data) != 1 {
		t.Fatalf("expected 1 draft appointment, got %d", len(data))
	}
}

// Test appointment validation
func TestAppointmentValidation(t *testing.T) {
	g, db := setupAppointmentRouter(t)
	admin := auth.User{}
	db.Where("username = ?", "admin").First(&admin)
	tok, _ := jwtpkg.GenerateToken(admin.ID, admin.Username)

	// Test missing required fields
	createReq := map[string]any{
		"contact_person": "张三",
		// Missing work_date, time_slot, work_location, work_content
	}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/appointments", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+tok)
	req.Header.Set("Content-Type", "application/json")
	r := httptest.NewRecorder()
	g.ServeHTTP(r, req)
	if r.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing fields, got: %d", r.Code)
	}
}

// Test my appointments endpoint
func TestMyAppointments(t *testing.T) {
	g, db := setupAppointmentRouter(t)
	admin := auth.User{}
	db.Where("username = ?", "admin").First(&admin)
	tok, _ := jwtpkg.GenerateToken(admin.ID, admin.Username)

	// Get project
	var project project.Project
	db.Where("code = ?", "TEST-001").First(&project)

	// Create appointment for current user
	workDate := time.Now().Add(24 * time.Hour)
	appointment := ConstructionAppointment{
		AppointmentNo:   "AP006",
		ProjectID:       &project.ID,
		ApplicantID:     admin.ID,
		ApplicantName:   "Admin",
		ContactPerson:   "我的预约",
		ContactPhone:    "13800138005",
		WorkDate:        workDate,
		TimeSlot:        "上午",
		WorkLocation:    "F区",
		WorkContent:     "我的预约内容",
		Status:          "draft",
	}
	db.Create(&appointment)

	// Get my appointments
	req := httptest.NewRequest("GET", "/api/appointments/my", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	r := httptest.NewRecorder()
	g.ServeHTTP(r, req)
	if r.Code != http.StatusOK {
		t.Fatalf("my appointments failed: %d %s", r.Code, r.Body.String())
	}
	var resp map[string]any
	json.Unmarshal(r.Body.Bytes(), &resp)
	data := resp["data"].([]any)
	if len(data) == 0 {
		t.Fatalf("expected at least 1 appointment")
	}
}

// Helper function
func jsonNumberToString(v any) string {
	switch x := v.(type) {
	case float64:
		return fmt.Sprintf("%.0f", x)
	case int:
		return fmt.Sprintf("%d", x)
	case int64:
		return fmt.Sprintf("%d", x)
	case uint:
		return fmt.Sprintf("%d", x)
	case uint64:
		return fmt.Sprintf("%d", x)
	default:
		return "0"
	}
}
