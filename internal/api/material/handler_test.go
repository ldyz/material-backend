package material

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	jwtpkg "github.com/yourorg/material-backend/backend/pkg/jwt"
	"github.com/yourorg/material-backend/backend/internal/api/auth"
	"github.com/yourorg/material-backend/backend/internal/api/project"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupMaterialRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil { t.Fatalf("open db fail: %v", err) }
	db.Migrator().DropTable(&auth.User{}, &auth.Role{}, &project.Project{}, &Material{})
	db.AutoMigrate(&auth.User{}, &auth.Role{}, &project.Project{}, &Material{})
	
	r := auth.Role{Name: "admin", Permissions: "admin"}
	db.Create(&r)
	u := auth.User{Username: "admin", Email: "admin@example.com", Role: "admin", IsActive: true}
	u.SetPassword("admin")
	db.Create(&u)

	g := gin.Default()
	rg := g.Group("/api")
	auth.RegisterRoutes(rg, db)
	project.RegisterRoutes(rg, db)
	RegisterRoutes(rg, db)
	return g, db
}

func TestMaterialCRUD(t *testing.T) {
	g, db := setupMaterialRouter(t)
	admin := auth.User{}
	db.Where("username = ?", "admin").First(&admin)
	tok, _ := jwtpkg.GenerateToken(admin.ID, admin.Username)

	// create project first
	p := project.Project{Name: "Test Project", Code: "TP-001"}
	db.Create(&p)

	// create material
	body := map[string]any{"name": "Steel Rod", "code": "MAT-001", "category": "metal", "project_id": strconv.Itoa(int(p.ID)), "quantity": 100, "price": 50.5}
	bb, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/material/materials", bytes.NewReader(bb))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tok)
	r := httptest.NewRecorder()
	g.ServeHTTP(r, req)
	if r.Code != 201 { t.Fatalf("create failed: %d %s", r.Code, r.Body.String()) }
	var cres map[string]any
	json.Unmarshal(r.Body.Bytes(), &cres)
	data := cres["data"].(map[string]any)
	id := data["id"]

	// list
	req2 := httptest.NewRequest("GET", "/api/material/materials", nil)
	req2.Header.Set("Authorization", "Bearer "+tok)
	r2 := httptest.NewRecorder()
	g.ServeHTTP(r2, req2)
	if r2.Code != 200 { t.Fatalf("list failed: %d %s", r2.Code, r2.Body.String()) }
	var lres map[string]any
	json.Unmarshal(r2.Body.Bytes(), &lres)
	materials := lres["data"].([]any)
	if len(materials) != 1 { t.Fatalf("expected 1 material, got %d", len(materials)) }

	// get by id
	url := "/api/material/materials/" + jsonNumberToString(id)
	req3 := httptest.NewRequest("GET", url, nil)
	req3.Header.Set("Authorization", "Bearer "+tok)
	r3 := httptest.NewRecorder()
	g.ServeHTTP(r3, req3)
	if r3.Code != 200 { t.Fatalf("get failed: %d %s", r3.Code, r3.Body.String()) }

	// update
	upbody := map[string]any{"quantity": 150}
	upbb, _ := json.Marshal(upbody)
	req4 := httptest.NewRequest("PUT", url, bytes.NewReader(upbb))
	req4.Header.Set("Content-Type", "application/json")
	req4.Header.Set("Authorization", "Bearer "+tok)
	r4 := httptest.NewRecorder()
	g.ServeHTTP(r4, req4)
	if r4.Code != 200 { t.Fatalf("update failed: %d %s", r4.Code, r4.Body.String()) }

	// delete
	req5 := httptest.NewRequest("DELETE", url, nil)
	req5.Header.Set("Authorization", "Bearer "+tok)
	r5 := httptest.NewRecorder()
	g.ServeHTTP(r5, req5)
	if r5.Code != 200 { t.Fatalf("delete failed: %d %s", r5.Code, r5.Body.String()) }
}

func TestMaterialFilters(t *testing.T) {
	g, db := setupMaterialRouter(t)
	admin := auth.User{}
	db.Where("username = ?", "admin").First(&admin)
	tok, _ := jwtpkg.GenerateToken(admin.ID, admin.Username)

	// create materials
	mats := []Material{
		{Name: "Steel Rod", Category: "metal", Quantity: 100},
		{Name: "Copper Wire", Category: "metal", Quantity: 50},
		{Name: "Concrete", Category: "cement", Quantity: 200},
	}
	for _, m := range mats { db.Create(&m) }

	// filter by category=metal
	req := httptest.NewRequest("GET", "/api/material/materials?category=metal", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	r := httptest.NewRecorder()
	g.ServeHTTP(r, req)
	if r.Code != 200 { t.Fatalf("expected 200, got %d", r.Code) }
	var resp map[string]any
	json.Unmarshal(r.Body.Bytes(), &resp)
	materials := resp["data"].([]any)
	if len(materials) != 2 { t.Fatalf("expected 2 metal materials, got %d", len(materials)) }

	// search by name
	req2 := httptest.NewRequest("GET", "/api/material/materials?search=Steel", nil)
	req2.Header.Set("Authorization", "Bearer "+tok)
	r2 := httptest.NewRecorder()
	g.ServeHTTP(r2, req2)
	if r2.Code != 200 { t.Fatalf("search failed: %d", r2.Code) }
	var resp2 map[string]any
	json.Unmarshal(r2.Body.Bytes(), &resp2)
	materials2 := resp2["data"].([]any)
	if len(materials2) != 1 { t.Fatalf("expected 1 Steel material, got %d", len(materials2)) }
}

// Helper to convert JSON numbers to string
func jsonNumberToString(v any) string {
	switch x := v.(type) {
	case float64:
		return strconv.FormatInt(int64(x), 10)
	default:
		return "0"
	}
}
