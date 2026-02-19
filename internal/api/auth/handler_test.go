package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	jwtpkg "github.com/yourorg/material-backend/backend/pkg/jwt"
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
	db.AutoMigrate(&User{}, &Role{})
	// create admin role and user
	r := Role{Name: "admin", Permissions: "admin"}
	db.Create(&r)
	u := User{Username: "admin", Email: "admin@example.com", Role: "admin", IsActive: true}
	u.SetPassword("admin")
	db.Create(&u)

	// setup gin
	g := gin.Default()
	rg := g.Group("/api")
	RegisterRoutes(rg, db)
	return g, db
}

func TestLoginAndMeFlow(t *testing.T) {
	g, _ := setupTestRouter(t)
	// login
	body := map[string]string{"username": "admin", "password": "admin"}
	bb, _ := json.Marshal(body)
	req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(bb))
	req.Header.Set("Content-Type", "application/json")
	r := httptest.NewRecorder()
	g.ServeHTTP(r, req)
	if r.Code != http.StatusOK {
		t.Fatalf("login failed, status: %d, body: %s", r.Code, r.Body.String())
	}
	var resp map[string]any
	json.Unmarshal(r.Body.Bytes(), &resp)
	// Token is in meta.token
	meta, ok := resp["meta"].(map[string]any)
	if !ok {
		t.Fatalf("no meta in response: %v", resp)
	}
	token, ok := meta["token"].(string)
	if !ok || token == "" {
		t.Fatalf("no token returned, got: %v", meta)
	}
	// call /me
	req2 := httptest.NewRequest("GET", "/api/auth/me", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	r2 := httptest.NewRecorder()
	g.ServeHTTP(r2, req2)
	if r2.Code != http.StatusOK {
		t.Fatalf("me failed: %d %s", r2.Code, r2.Body.String())
	}
}

func TestPermissionDenied(t *testing.T) {
	g, db := setupTestRouter(t)
	// create user with no roles/permissions
	u := User{Username: "limited", Email: "l@example.com", Role: "user", IsActive: true}
	u.SetPassword("pass123")
	db.Create(&u)
	// generate token
	tok, _ := jwtpkg.GenerateToken(u.ID, u.Username)
	// try to access users list (requires user_view)
	req := httptest.NewRequest("GET", "/api/auth/users", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	r := httptest.NewRecorder()
	g.ServeHTTP(r, req)
	if r.Code != http.StatusForbidden {
		t.Fatalf("expected 403 forbidden, got %d body:%s", r.Code, r.Body.String())
	}
}
