package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/GdAyo19/FundsFlow/utils"
)

func TestHealthCheck(t *testing.T) {
	// Placeholder test - will expand with integration tests
	t.Log("Controller package loaded successfully")
}

func TestGenerateToken(t *testing.T) {
	token, err := utils.GenerateToken(1)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	if token == "" {
		t.Fatal("Expected non-empty token")
	}

	claims, err := utils.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}
	if claims.UserID != 1 {
		t.Fatalf("Expected UserID 1, got %d", claims.UserID)
	}
}

func TestInvalidToken(t *testing.T) {
	_, err := utils.ValidateToken("invalid-token")
	if err == nil {
		t.Fatal("Expected error for invalid token")
	}
}

func TestExpiredToken(t *testing.T) {
	claims := utils.Claims{
		UserID: 1,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("fundsflow-my-secret-key"))

	_, err := utils.ValidateToken(tokenString)
	if err == nil {
		t.Fatal("Expected error for expired token")
	}
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return r
}

func performRequest(r http.HandlerFunc, method, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func parseJSONResponse(t *testing.T, w *httptest.ResponseRecorder, v interface{}) {
	t.Helper()
	if err := json.Unmarshal(w.Body.Bytes(), v); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}
}
