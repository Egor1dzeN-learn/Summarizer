package services_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"summarizer/backend/http/models"
	"summarizer/backend/http/routes"
	services "summarizer/backend/http/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func setupAuthRouter(authSvc services.AuthService) *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(gin.Recovery())

	store := cookie.NewStore([]byte("test-secret"))
	r.Use(sessions.Sessions("summarizer", store))

	public := r.Group("/api")
	routes.NewAuthHandler(authSvc).Bind(public)

	return r
}

func TestAuthHandlerLoginAndMe(t *testing.T) {
	db := newTestDB(t)
	authSvc := services.NewAuthService(db)

	user := authSvc.CreateUser(777, "Tester")
	token := authSvc.IssueLoginToken(user)

	r := setupAuthRouter(authSvc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(fmt.Sprintf(`{"token":"%s"}`, token)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected login to succeed, got status %d", w.Code)
	}

	cookies := w.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatalf("expected session cookie to be set")
	}

	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/me", nil)
	req.AddCookie(cookies[0])
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected /api/me to succeed with session, got %d", w.Code)
	}

	var payload models.User
	if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to decode /api/me response: %v", err)
	}
	if payload.ID != user.ID {
		t.Fatalf("expected to receive authenticated user, got %+v", payload)
	}
}

func TestAuthHandlerRejectsInvalidToken(t *testing.T) {
	db := newTestDB(t)
	authSvc := services.NewAuthService(db)
	authSvc.CreateUser(888, "Tester")

	r := setupAuthRouter(authSvc)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(`{"token":"invalid"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected unauthorized for invalid token, got %d", w.Code)
	}
}
