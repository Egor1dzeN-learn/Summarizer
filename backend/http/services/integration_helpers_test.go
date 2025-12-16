package services_test

import (
	"sync"
	"testing"

	"summarizer/backend/http/models"
	"summarizer/backend/http/resources"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type workerCall struct {
	Text   string
	Prompt string
}

type stubWorkerNode struct {
	mu    sync.Mutex
	calls []workerCall

	result string
	err    error
	done   chan struct{}
}

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open sqlite in memory: %v", err)
	}

	db.Exec("PRAGMA foreign_keys = ON")

	if err := db.AutoMigrate(&models.User{}, &models.Chat{}, &models.ChatEntry{}).Error; err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}

	return db
}

func newStubWorker(result string, err error) *stubWorkerNode {
	return &stubWorkerNode{
		result: result,
		err:    err,
		done:   make(chan struct{}),
	}
}

func (s *stubWorkerNode) Summarize(text string, prompt string, accept resources.SummarizeResultConsumer) {
	s.mu.Lock()
	s.calls = append(s.calls, workerCall{Text: text, Prompt: prompt})
	res := s.result
	err := s.err
	done := s.done
	s.mu.Unlock()

	go func() {
		if err != nil {
			accept("", err)
		} else {
			if res == "" {
				res = "ok"
			}
			accept(res, nil)
		}

		if done != nil {
			select {
			case <-done:
			default:
				close(done)
			}
		}
	}()
}

func (s *stubWorkerNode) CloseConnection() {}

func (s *stubWorkerNode) Calls() []workerCall {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]workerCall, len(s.calls))
	copy(out, s.calls)
	return out
}
