package project

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"strconv"

	"github.com/gin-gonic/gin"
	jwtpkg "github.com/yourorg/material-backend/backend/pkg/jwt"
	"github.com/yourorg/material-backend/backend/internal/api/auth"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
	// use sqlite in-memory DB for tests
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db fail: %v", err)
	}
	// migrate
	db.AutoMigrate(&auth.User{}, &auth.Role{}, &Project{})
	// create admin role and user
	r := auth.Role{Name: "admin", Permissions: "admin"}
	db.Create(&r)
	u := auth.User{Username: "admin", Email: "admin@example.com", Role: "admin", IsActive: true}
	u.SetPassword("admin")
	db.Create(&u)

	// setup gin
	g := gin.Default()
	rg := g.Group("/api")
	auth.RegisterRoutes(rg, db)
	RegisterRoutes(rg, db)
	return g, db
}

func TestProjectCRUD(t *testing.T) {
	g, _ := setupTestRouter(t)
	// login to get token
	login := map[string]string{"username": "admin", "password": "admin"}
	bb, _ := json.Marshal(login)
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(bb))
	req.Header.Set("Content-Type", "application/json")
	r := httptest.NewRecorder()
	g.ServeHTTP(r, req)
	if r.Code != http.StatusOK { t.Fatalf("login failed: %d %s", r.Code, r.Body.String()) }
	var resp map[string]any
	json.Unmarshal(r.Body.Bytes(), &resp)
	meta := resp["meta"].(map[string]any)
	tok := meta["token"].(string)

	// create project
	p := map[string]any{"name": "Test Project", "description": "desc"}
	pb, _ := json.Marshal(p)
	req2 := httptest.NewRequest("POST", "/api/project/projects", bytes.NewReader(pb))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Authorization", "Bearer "+tok)
	r2 := httptest.NewRecorder()
	g.ServeHTTP(r2, req2)
	if r2.Code != http.StatusCreated { t.Fatalf("create failed: %d %s", r2.Code, r2.Body.String()) }
	var cres map[string]any
	json.Unmarshal(r2.Body.Bytes(), &cres)
	proj := cres["data"].(map[string]any)
	id := proj["id"]

	// list
	req3 := httptest.NewRequest("GET", "/api/project/projects", nil)
	req3.Header.Set("Authorization", "Bearer "+tok)
	r3 := httptest.NewRecorder()
	g.ServeHTTP(r3, req3)
	if r3.Code != http.StatusOK { t.Fatalf("list failed: %d %s", r3.Code, r3.Body.String()) }

	// get by id
	url := "/api/project/projects/" + (jsonNumberToString(id))
	req4 := httptest.NewRequest("GET", url, nil)
	req4.Header.Set("Authorization", "Bearer "+tok)
	r4 := httptest.NewRecorder()
	g.ServeHTTP(r4, req4)
	if r4.Code != http.StatusOK { t.Fatalf("get failed: %d %s", r4.Code, r4.Body.String()) }

	// update
	up := map[string]any{"name": "Updated Project"}
	ub, _ := json.Marshal(up)
	req5 := httptest.NewRequest("PUT", url, bytes.NewReader(ub))
	req5.Header.Set("Content-Type", "application/json")
	req5.Header.Set("Authorization", "Bearer "+tok)
	r5 := httptest.NewRecorder()
	g.ServeHTTP(r5, req5)
	if r5.Code != http.StatusOK { t.Fatalf("update failed: %d %s", r5.Code, r5.Body.String()) }

	// delete
	req6 := httptest.NewRequest("DELETE", url, nil)
	req6.Header.Set("Authorization", "Bearer "+tok)
	r6 := httptest.NewRecorder()
	g.ServeHTTP(r6, req6)
	if r6.Code != http.StatusOK { t.Fatalf("delete failed: %d %s", r6.Code, r6.Body.String()) }

	// get should 404
	req7 := httptest.NewRequest("GET", url, nil)
	req7.Header.Set("Authorization", "Bearer "+tok)
	r7 := httptest.NewRecorder()
	g.ServeHTTP(r7, req7)
	if r7.Code != http.StatusNotFound { t.Fatalf("expected 404 after delete, got: %d %s", r7.Code, r7.Body.String()) }
}

func TestProjectMembers(t *testing.T) {
	g, db := setupTestRouter(t)
	// admin token
	admin := auth.User{}
	db.Where("username = ?", "admin").First(&admin)
	tok, _ := jwtpkg.GenerateToken(admin.ID, admin.Username)

	// create user to add
	u := auth.User{Username: "bob", Email: "b@example.com", Role: "user", IsActive: true}
	u.SetPassword("pass")
	db.Create(&u)

	// create project
	p := Project{Name: "Member Project", Code: "MP-1"}
	db.Create(&p)

	// add member
	body := map[string]any{"user_ids": []uint{u.ID}}
	bb, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/project/projects/"+jsonNumberToString(p.ID)+"/members", bytes.NewReader(bb))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tok)
	r := httptest.NewRecorder()
	g.ServeHTTP(r, req)
	if r.Code != http.StatusOK { t.Fatalf("add member failed: %d %s", r.Code, r.Body.String()) }

	// list members
	req2 := httptest.NewRequest("GET", "/api/project/projects/"+jsonNumberToString(p.ID)+"/members", nil)
	req2.Header.Set("Authorization", "Bearer "+tok)
	r2 := httptest.NewRecorder()
	g.ServeHTTP(r2, req2)
	if r2.Code != http.StatusOK { t.Fatalf("list members failed: %d %s", r2.Code, r2.Body.String()) }
	var resp map[string]any
	json.Unmarshal(r2.Body.Bytes(), &resp)
	members := resp["data"].([]any)
	if len(members) != 1 { t.Fatalf("expected 1 member, got %d", len(members)) }

	// remove member
	req3 := httptest.NewRequest("DELETE", "/api/project/projects/"+jsonNumberToString(p.ID)+"/members/"+jsonNumberToString(u.ID), nil)
	req3.Header.Set("Authorization", "Bearer "+tok)
	r3 := httptest.NewRecorder()
	g.ServeHTTP(r3, req3)
	if r3.Code != http.StatusOK { t.Fatalf("remove member failed: %d %s", r3.Code, r3.Body.String()) }
}

// helper to turn numeric json number into string without float formatting
func jsonNumberToString(v any) string {
	switch x := v.(type) {
	case float64:
		return strconv.FormatInt(int64(x), 10)
	case int:
		return strconv.Itoa(x)
	case int64:
		return strconv.FormatInt(x, 10)
	case uint:
		return strconv.FormatUint(uint64(x), 10)
	case uint64:
		return strconv.FormatUint(x, 10)
	default:
		return "0"
	}
}
