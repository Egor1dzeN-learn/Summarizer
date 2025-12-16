package services_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"summarizer/backend/http/middleware"
	"summarizer/backend/http/models"
	"summarizer/backend/http/routes"
	services "summarizer/backend/http/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func setupFullRouter(authSvc services.AuthService, chatSvc services.ChatService) *gin.Engine {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(gin.Recovery())

	store := cookie.NewStore([]byte("test-secret"))
	r.Use(sessions.Sessions("summarizer", store))

	public := r.Group("/api")
	routes.NewAuthHandler(authSvc).Bind(public)

	protected := r.Group("/api")
	protected.Use(middleware.AuthRequired())
	routes.NewChatHandler(chatSvc).Bind(protected)

	return r
}

func loginAndGetCookie(t *testing.T, r *gin.Engine, token string) *http.Cookie {
	t.Helper()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(fmt.Sprintf(`{"token":"%s"}`, token)))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("login failed with status %d", w.Code)
	}
	cookies := w.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatalf("expected login to set cookie")
	}
	return cookies[0]
}

func TestChatHandlerHistoryEmptyForNewAccount(t *testing.T) {
	db := newTestDB(t)
	authSvc := services.NewAuthService(db)
	chatSvc := services.NewChatService(db, newStubWorker("", nil))
	user := authSvc.CreateUser(1, "User")
	token := authSvc.IssueLoginToken(user)

	r := setupFullRouter(authSvc, chatSvc)
	cookie := loginAndGetCookie(t, r, token)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/chats", nil)
	req.AddCookie(cookie)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 for empty history, got %d", w.Code)
	}

	var chats []models.Chat
	if err := json.Unmarshal(w.Body.Bytes(), &chats); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(chats) != 0 {
		t.Fatalf("expected no chats for new user, got %d", len(chats))
	}
}

func TestChatHandlerCreateAndSummarizeFlow(t *testing.T) {
	db := newTestDB(t)
	worker := newStubWorker("computed", nil)
	authSvc := services.NewAuthService(db)
	chatSvc := services.NewChatService(db, worker)
	user := authSvc.CreateUser(2, "User")
	token := authSvc.IssueLoginToken(user)

	r := setupFullRouter(authSvc, chatSvc)
	cookie := loginAndGetCookie(t, r, token)

	createReq := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/chats", strings.NewReader(`{"prompt":"some prompt"}`))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie)
	r.ServeHTTP(createReq, req)

	if createReq.Code != http.StatusOK {
		t.Fatalf("expected chat creation to succeed, got %d", createReq.Code)
	}

	var created models.Chat
	if err := json.Unmarshal(createReq.Body.Bytes(), &created); err != nil {
		t.Fatalf("failed to decode created chat: %v", err)
	}
	if created.ID == 0 {
		t.Fatalf("expected created chat to have ID")
	}

	entryReq := httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/chats/%d", created.ID), strings.NewReader(`{"msg":"question"}`))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie)
	r.ServeHTTP(entryReq, req)

	if entryReq.Code != http.StatusOK {
		t.Fatalf("expected message creation to succeed, got %d", entryReq.Code)
	}

	<-worker.done

	var storedEntry models.ChatEntry
	if err := db.Last(&storedEntry).Error; err != nil {
		t.Fatalf("failed to load stored entry: %v", err)
	}
	if storedEntry.Answer == nil || *storedEntry.Answer != "computed" {
		t.Fatalf("expected worker answer to be persisted, got %#v", storedEntry.Answer)
	}
}

func TestChatHandlerDeleteChat(t *testing.T) {
	db := newTestDB(t)
	authSvc := services.NewAuthService(db)
	chatSvc := services.NewChatService(db, newStubWorker("", nil))
	user := authSvc.CreateUser(3, "User")
	token := authSvc.IssueLoginToken(user)

	chat := chatSvc.NewChat(user.ID, "prompt")

	r := setupFullRouter(authSvc, chatSvc)
	cookie := loginAndGetCookie(t, r, token)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/chats/%d", chat.ID), nil)
	req.AddCookie(cookie)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected delete to respond with 200, got %d", w.Code)
	}

	var count int64
	db.Model(&models.Chat{}).Count(&count)
	if count != 0 {
		t.Fatalf("expected chat to be removed from DB, count=%d", count)
	}
}
