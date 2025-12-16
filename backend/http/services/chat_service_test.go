package services

import (
	"context"
	"errors"
	"testing"

	"summarizer/backend/http/models"
)

func TestChatServiceHistoryEmptyForNewUser(t *testing.T) {
	db := newTestDB(t)
	svc := NewChatService(db, newStubWorker("", nil))

	chats := svc.GetChats(999)
	if len(chats) != 0 {
		t.Fatalf("expected empty history for new user, got %d chats", len(chats))
	}
}

func TestChatServiceNewChatTruncatesLongTitle(t *testing.T) {
	db := newTestDB(t)
	svc := NewChatService(db, newStubWorker("", nil))

	longPrompt := "This is a very long piece of text that should be trimmed for the title preview"
	chat := svc.NewChat(1, longPrompt)
	if chat.Title == longPrompt {
		t.Fatalf("expected chat title to be truncated for long prompt")
	}
	if got := len([]rune(chat.Title)); got != 23 {
		t.Fatalf("expected truncated title length 23, got %d (%q)", got, chat.Title)
	}
}

func TestChatServiceSummarizeCreatesEntryAndPersistsAnswer(t *testing.T) {
	db := newTestDB(t)
	worker := newStubWorker("short answer", nil)
	svc := NewChatService(db, worker)

	user := &models.User{TelegramID: 1, Name: "Alice"}
	db.Create(user)

	chat := svc.NewChat(user.ID, "original prompt")

	var onCompleteCalled bool
	entry := svc.Summarize(chat, "question text", func(result string) {
		onCompleteCalled = true
	})

	<-worker.done

	var stored models.ChatEntry
	if err := db.First(&stored, entry.ID).Error; err != nil {
		t.Fatalf("expected entry to be stored: %v", err)
	}
	if stored.Answer == nil || *stored.Answer != "short answer" {
		t.Fatalf("expected answer to be persisted, got %#v", stored.Answer)
	}
	if !onCompleteCalled {
		t.Fatalf("expected onComplete callback to be invoked")
	}

	calls := worker.Calls()
	if len(calls) != 1 {
		t.Fatalf("expected worker to be called once, got %d", len(calls))
	}
	if calls[0].Text != chat.Text || calls[0].Prompt != "question text" {
		t.Fatalf("unexpected worker call: %+v", calls[0])
	}
}

func TestChatServiceSummarizeStoresFailureOnWorkerError(t *testing.T) {
	db := newTestDB(t)
	worker := newStubWorker("", errors.New("worker failed"))
	svc := NewChatService(db, worker)

	user := &models.User{TelegramID: 2, Name: "Bob"}
	db.Create(user)
	chat := svc.NewChat(user.ID, "prompt")

	entry := svc.Summarize(chat, "question", func(string) {})
	<-worker.done

	var stored models.ChatEntry
	if err := db.First(&stored, entry.ID).Error; err != nil {
		t.Fatalf("expected entry to be stored: %v", err)
	}
	if stored.Answer == nil || *stored.Answer != "Failure: worker failed" {
		t.Fatalf("expected failure message to be stored, got %#v", stored.Answer)
	}
}

func TestChatServiceSummarizeStoresTimeoutAsFailure(t *testing.T) {
	db := newTestDB(t)
	worker := newStubWorker("", context.DeadlineExceeded)
	svc := NewChatService(db, worker)

	user := &models.User{TelegramID: 3, Name: "Carol"}
	db.Create(user)
	chat := svc.NewChat(user.ID, "prompt")

	entry := svc.Summarize(chat, "question", func(string) {})
	<-worker.done

	var stored models.ChatEntry
	if err := db.First(&stored, entry.ID).Error; err != nil {
		t.Fatalf("expected entry to be stored: %v", err)
	}
	if stored.Answer == nil || *stored.Answer == "" {
		t.Fatalf("expected timeout failure to be stored, got %#v", stored.Answer)
	}
	if expected := "Failure: context deadline exceeded"; *stored.Answer != expected {
		t.Fatalf("expected %q, got %q", expected, *stored.Answer)
	}
}

func TestChatServiceDeleteChatRemovesHistory(t *testing.T) {
	db := newTestDB(t)
	svc := NewChatService(db, newStubWorker("", nil))

	user := &models.User{TelegramID: 4, Name: "Dave"}
	db.Create(user)

	chat := svc.NewChat(user.ID, "prompt")
	entry := &models.ChatEntry{ChatID: chat.ID, Question: "q1"}
	db.Create(entry)

	svc.DeleteChat(chat)

	var chatCount int64
	db.Model(&models.Chat{}).Count(&chatCount)
	if chatCount != 0 {
		t.Fatalf("expected chat to be deleted, count=%d", chatCount)
	}

	var entryCount int64
	db.Model(&models.ChatEntry{}).Count(&entryCount)
	if entryCount != 0 {
		t.Fatalf("expected chat entries to be removed with chat, count=%d", entryCount)
	}
}
