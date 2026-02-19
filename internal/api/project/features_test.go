package project

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
	"bytes"
	"strconv"
	
	"github.com/gin-gonic/gin"
	"github.com/yourorg/material-backend/backend/internal/api/auth"
	jwtpkg "github.com/yourorg/material-backend/backend/pkg/jwt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupFeatureRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil { t.Fatalf("open db fail: %v", err) }
	// ensure a clean slate in the shared in-memory DB between tests
	db.Migrator().DropTable(&auth.User{}, &auth.Role{}, &Project{})
	db.AutoMigrate(&auth.User{}, &auth.Role{}, &Project{})
	r := auth.Role{Name: "admin", Permissions: "admin"}
	db.Create(&r)
	u := auth.User{Username: "admin", Email: "admin@example.com", Role: "admin", IsActive: true}
	u.SetPassword("admin")
	db.Create(&u)

	g := gin.Default()
	rg := g.Group("/api")
	auth.RegisterRoutes(rg, db)
	RegisterRoutes(rg, db)
	return g, db
}

func TestProjectFilters(t *testing.T) {
	g, db := setupFeatureRouter(t)
	// admin token
	admin := auth.User{}
	db.Where("username = ?", "admin").First(&admin)
	tok, _ := jwtpkg.GenerateToken(admin.ID, admin.Username)

	// create projects
	projs := []Project{
		{Name: "A1", Status: "active", Manager: "Alice"},
		{Name: "B1", Status: "closed", Manager: "Bob"},
		{Name: "A2", Status: "active", Manager: "Alice"},
	}
	for _, p := range projs { db.Create(&p) }

	// filter by status=active
	req := httptest.NewRequest("GET", "/api/project/projects?status=active", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	r := httptest.NewRecorder()
	g.ServeHTTP(r, req)
	if r.Code != 200 { t.Fatalf("expected 200, got %d body:%s", r.Code, r.Body.String()) }
	var resp map[string]any
	json.Unmarshal(r.Body.Bytes(), &resp)
	projects, ok := resp["data"].([]any)
	if !ok { t.Fatalf("expected data array in response, got: %v", resp) }
	if len(projects) != 2 { t.Fatalf("expected 2 active projects, got %d", len(projects)) }

	// filter by manager=Bob
	req2 := httptest.NewRequest("GET", "/api/project/projects?manager=Bob", nil)
	req2.Header.Set("Authorization", "Bearer "+tok)
	r2 := httptest.NewRecorder()
	g.ServeHTTP(r2, req2)
	if r2.Code != 200 { t.Fatalf("expected 200, got %d body:%s", r2.Code, r2.Body.String()) }
	var resp2 map[string]any
	json.Unmarshal(r2.Body.Bytes(), &resp2)
	projects2, ok := resp2["data"].([]any)
	if !ok { t.Fatalf("expected data array in response, got: %v", resp2) }
	if len(projects2) != 1 { t.Fatalf("expected 1 Bob project, got %d", len(projects2)) }
}

func TestProjectCodeGeneration(t *testing.T) {
	g, db := setupFeatureRouter(t)
	admin := auth.User{}
	db.Where("username = ?", "admin").First(&admin)
	tok, _ := jwtpkg.GenerateToken(admin.ID, admin.Username)

	// create projects without code and ensure unique generated codes
	for i := 0; i < 5; i++ {
		body := map[string]any{"name": "GenProj" + strconv.Itoa(i)}
		bb, _ := json.Marshal(body)
		req := httptest.NewRequest("POST", "/api/project/projects", bytes.NewReader(bb))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tok)
		r := httptest.NewRecorder()
		g.ServeHTTP(r, req)
		if r.Code != 201 { t.Fatalf("create failed: %d %s", r.Code, r.Body.String()) }
		var cres map[string]any
		json.Unmarshal(r.Body.Bytes(), &cres)
		proj, ok := cres["data"].(map[string]any)
		if !ok { t.Fatalf("expected data in response: %v", cres) }
		if proj["code"] == nil || proj["code"] == "" { t.Fatalf("expected code generated, got empty") }
	}
	// ensure uniqueness among non-empty codes (ignore legacy empty codes created directly in tests)
	var codes []string
	db.Model(&Project{}).Pluck("code", &codes)
	m := map[string]bool{}
	for _, c := range codes {
		if c == "" { continue }
		if m[c] { t.Fatalf("duplicate code found: %s", c) }
		m[c] = true
	}
}
