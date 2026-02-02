package inbound

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/material-backend/backend/internal/api/auth"
	"github.com/yourorg/material-backend/backend/internal/db"
	"gorm.io/gorm"
)

func setupInboundRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
	d := db.New(":memory:")
	d.AutoMigrate(&auth.User{}, &auth.Role{}, &InboundOrder{}, &InboundOrderItem{})

	// Seed admin user
	admin := auth.User{
		Username: "admin",
		Email:    "admin@test.com",
	}
	admin.SetPassword("password123")
	d.Create(&admin)

	g := gin.New()
	auth.RegisterRoutes(g, d)
	RegisterRoutes(g, d)

	return g, d
}

func TestInboundOrderCRUD(t *testing.T) {
	g, db := setupInboundRouter(t)
	admin := auth.User{}
	db.Where("username = ?", "admin").First(&admin)
	tok, _ := auth.GenerateToken(admin.ID, admin.Username)

	// Create
	createReq := map[string]any{
		"supplier":     "Supplier A",
		"contact":      "Contact Person",
		"project_id":   1,
		"notes":        "Test notes",
		"total_amount": 1000.00,
		"items": []map[string]any{
			{
				"material_id": 1,
				"quantity":    100,
				"unit_price":  10.00,
				"remark":      "Item 1",
			},
		},
	}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/inbound/inbound-orders", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+tok)
	req.Header.Set("Content-Type", "application/json")
	r := httptest.NewRecorder()
	g.ServeHTTP(r, req)
	if r.Code != http.StatusCreated {
		t.Fatalf("expected %d, got %d", http.StatusCreated, r.Code)
	}
	var resp map[string]any
	json.Unmarshal(r.Body.Bytes(), &resp)
	order := resp["order"].(map[string]any)
	orderID := int(order["id"].(float64))

	// List
	req2 := httptest.NewRequest("GET", "/api/inbound/inbound-orders", nil)
	req2.Header.Set("Authorization", "Bearer "+tok)
	r2 := httptest.NewRecorder()
	g.ServeHTTP(r2, req2)
	if r2.Code != http.StatusOK {
		t.Fatalf("list failed: %d", r2.Code)
	}

	// Get
	req3 := httptest.NewRequest("GET", fmt.Sprintf("/api/inbound/inbound-orders/%d", orderID), nil)
	req3.Header.Set("Authorization", "Bearer "+tok)
	r3 := httptest.NewRecorder()
	g.ServeHTTP(r3, req3)
	if r3.Code != http.StatusOK {
		t.Fatalf("get failed: %d", r3.Code)
	}

	// Update
	updateReq := map[string]any{
		"supplier": "Supplier B",
		"contact":  "New Contact",
		"items": []map[string]any{
			{
				"material_id": 2,
				"quantity":    50,
				"unit_price":  20.00,
			},
		},
	}
	body, _ = json.Marshal(updateReq)
	req4 := httptest.NewRequest("PUT", fmt.Sprintf("/api/inbound/inbound-orders/%d", orderID), bytes.NewBuffer(body))
	req4.Header.Set("Authorization", "Bearer "+tok)
	req4.Header.Set("Content-Type", "application/json")
	r4 := httptest.NewRecorder()
	g.ServeHTTP(r4, req4)
	if r4.Code != http.StatusOK {
		t.Fatalf("update failed: %d", r4.Code)
	}

	// Delete
	req5 := httptest.NewRequest("DELETE", fmt.Sprintf("/api/inbound/inbound-orders/%d", orderID), nil)
	req5.Header.Set("Authorization", "Bearer "+tok)
	r5 := httptest.NewRecorder()
	g.ServeHTTP(r5, req5)
	if r5.Code != http.StatusOK {
		t.Fatalf("delete failed: %d", r5.Code)
	}
}

func TestInboundOrderApprove(t *testing.T) {
	g, db := setupInboundRouter(t)
	admin := auth.User{}
	db.Where("username = ?", "admin").First(&admin)
	tok, _ := auth.GenerateToken(admin.ID, admin.Username)

	// Create order
	createReq := map[string]any{
		"supplier": "Test Supplier",
		"items": []map[string]any{
			{
				"material_id": 1,
				"quantity":    100,
				"unit_price":  10.00,
			},
		},
	}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/inbound/inbound-orders", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+tok)
	req.Header.Set("Content-Type", "application/json")
	r := httptest.NewRecorder()
	g.ServeHTTP(r, req)
	var resp map[string]any
	json.Unmarshal(r.Body.Bytes(), &resp)
	order := resp["order"].(map[string]any)
	orderID := int(order["id"].(float64))

	// Approve
	approveReq := map[string]any{
		"approve": true,
		"remark":  "Approved",
	}
	body, _ = json.Marshal(approveReq)
	req2 := httptest.NewRequest("POST", fmt.Sprintf("/api/inbound/inbound-orders/%d/approve", orderID), bytes.NewBuffer(body))
	req2.Header.Set("Authorization", "Bearer "+tok)
	req2.Header.Set("Content-Type", "application/json")
	r2 := httptest.NewRecorder()
	g.ServeHTTP(r2, req2)
	if r2.Code != http.StatusOK {
		t.Fatalf("approve failed: %d", r2.Code)
	}

	// Verify status changed
	var updatedOrder InboundOrder
	db.First(&updatedOrder, orderID)
	if updatedOrder.Status != StatusApproved {
		t.Fatalf("expected status %s, got %s", StatusApproved, updatedOrder.Status)
	}
}

func TestInboundOrderFilters(t *testing.T) {
	g, db := setupInboundRouter(t)
	admin := auth.User{}
	db.Where("username = ?", "admin").First(&admin)
	tok, _ := auth.GenerateToken(admin.ID, admin.Username)

	// Create multiple orders
	orders := []InboundOrder{
		{
			OrderNo:     "RK001",
			Supplier:    "Supplier A",
			CreatorID:   admin.ID,
			CreatorName: "admin",
			Status:      StatusPending,
		},
		{
			OrderNo:     "RK002",
			Supplier:    "Supplier B",
			CreatorID:   admin.ID,
			CreatorName: "admin",
			Status:      StatusApproved,
		},
		{
			OrderNo:     "RK003",
			Supplier:    "Supplier A",
			CreatorID:   admin.ID,
			CreatorName: "admin",
			Status:      StatusPending,
		},
	}
	for _, o := range orders {
		db.Create(&o)
	}

	// Filter by supplier
	req := httptest.NewRequest("GET", "/api/inbound/inbound-orders?supplier=Supplier+A", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	r := httptest.NewRecorder()
	g.ServeHTTP(r, req)
	if r.Code != http.StatusOK {
		t.Fatalf("filter failed: %d", r.Code)
	}
	var resp map[string]any
	json.Unmarshal(r.Body.Bytes(), &resp)
	orders2 := resp["inbound_orders"].([]any)
	if len(orders2) != 2 {
		t.Fatalf("expected 2 orders, got %d", len(orders2))
	}

	// Filter by status
	req2 := httptest.NewRequest("GET", "/api/inbound/inbound-orders?status=pending", nil)
	req2.Header.Set("Authorization", "Bearer "+tok)
	r2 := httptest.NewRecorder()
	g.ServeHTTP(r2, req2)
	if r2.Code != http.StatusOK {
		t.Fatalf("status filter failed: %d", r2.Code)
	}
	json.Unmarshal(r2.Body.Bytes(), &resp)
	orders3 := resp["inbound_orders"].([]any)
	if len(orders3) != 2 {
		t.Fatalf("expected 2 pending orders, got %d", len(orders3))
	}
}
