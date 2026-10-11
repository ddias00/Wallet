package integration

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wallet/internal/dto"
	"github.com/wallet/internal/handler"
	"github.com/wallet/internal/model"
	"github.com/wallet/internal/repository"
	"github.com/wallet/internal/routes"
	"github.com/wallet/internal/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupAPI(t *testing.T) *gin.Engine {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	if err := db.AutoMigrate(&model.User{}, &model.Wallet{}, &model.Transaction{}, &model.PixKey{}, &model.Boleto{}, &model.RefreshToken{}); err != nil {
		t.Fatal(err)
	}

	repo := repository.New(db)
	svc := service.New(repo, "secret", 15*time.Minute, 24*time.Hour)
	h := handler.New(svc)

	r := gin.Default()
	routes.Setup(r, h, "secret")

	return r
}

func TestAPIRegisterLogin(t *testing.T) {
	r := setupAPI(t)

	// Register
	regBody, _ := json.Marshal(dto.RegisterRequest{
		Name: "Test", CPF: "12345678909", Email: "test@example.com", Phone: "11987654321", Password: "password123",
	})
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(regBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 201 {
		t.Errorf("Register status = %d, want 201", w.Code)
	}

	// Login
	loginBody, _ := json.Marshal(dto.LoginRequest{Email: "test@example.com", Password: "password123"})
	req = httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Login status = %d, want 200", w.Code)
	}

	var loginResp dto.Response
	json.Unmarshal(w.Body.Bytes(), &loginResp)

	tokenData := loginResp.Data.(map[string]interface{})
	token := tokenData["access_token"].(string)

	if token == "" {
		t.Error("No access token returned")
	}

	// Get Balance
	req = httptest.NewRequest("GET", "/api/v1/wallet/balance", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Balance status = %d, want 200", w.Code)
	}
}

func TestAPIHealth(t *testing.T) {
	r := setupAPI(t)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Health status = %d, want 200", w.Code)
	}

	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["status"] != "ok" {
		t.Errorf("Health status = %s, want ok", response["status"])
	}
}

func TestAPIUnauthorized(t *testing.T) {
	r := setupAPI(t)

	req := httptest.NewRequest("GET", "/api/v1/wallet/balance", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("Unauthorized status = %d, want 401", w.Code)
	}
}

func TestAPIInvalidToken(t *testing.T) {
	r := setupAPI(t)

	req := httptest.NewRequest("GET", "/api/v1/wallet/balance", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("Invalid token status = %d, want 401", w.Code)
	}
}
