package services

import (
	"testing"

	"summarizer/backend/http/models"
)

func TestAuthServiceIssueAndConsumeLoginToken(t *testing.T) {
	db := newTestDB(t)
	svc := NewAuthService(db)

	user := svc.CreateUser(123, "Alice")
	token := svc.IssueLoginToken(user)
	if token == "" {
		t.Fatalf("expected non-empty token")
	}

	found := svc.FindUserByLoginToken(token)
	if found == nil || found.ID != user.ID {
		t.Fatalf("expected to find user by token, got %#v", found)
	}

	var refreshed models.User
	if err := db.First(&refreshed, user.ID).Error; err != nil {
		t.Fatalf("failed to reload user: %v", err)
	}
	if refreshed.LoginToken != nil {
		t.Fatalf("expected login token to be cleared after consumption, got %v", *refreshed.LoginToken)
	}

	if reused := svc.FindUserByLoginToken(token); reused != nil {
		t.Fatalf("expected reused token to be invalid")
	}
}

func TestAuthServiceRejectsDuplicateTelegramID(t *testing.T) {
	db := newTestDB(t)
	svc := NewAuthService(db)

	first := svc.CreateUser(123, "Alice")
	second := svc.CreateUser(123, "Bob")

	var count int64
	db.Model(&models.User{}).Count(&count)
	if count != 1 {
		t.Fatalf("expected only one user in DB, got %d", count)
	}
	if second.ID != 0 && second.ID != first.ID {
		t.Fatalf("expected duplicate creation to reuse existing record or fail, got new id=%d", second.ID)
	}
}
