package services

import (
	"fmt"
	"summarizer/backend/http/models"
	"summarizer/backend/http/resources"

	"github.com/jinzhu/gorm"
	"golang.org/x/exp/utf8string"
)

type ChatService interface {
	GetChats(userID uint) []models.Chat

	NewChat(userID uint, prompt string) *models.Chat
	FindChat(userID uint, id uint) *models.Chat
	DeleteChat(chat *models.Chat)

	SetUserTelegramActiveChat(user *models.User, chat *models.Chat)

	Summarize(chat *models.Chat, question string, onComplete func(result string)) *models.ChatEntry
}

type chatService struct {
	db   *gorm.DB
	node resources.WorkerNode
}

func NewChatService(db *gorm.DB, node resources.WorkerNode) ChatService {
	return &chatService{db: db, node: node}
}

func (s *chatService) GetChats(userID uint) []models.Chat {
	var chats []models.Chat

	err := s.db.Preload("Entries").
		Where("user_id = ?", userID).
		Find(&chats).Error
	if err != nil {
		panic(err)
	}

	return chats
}

func (s *chatService) NewChat(userID uint, prompt string) *models.Chat {
	truncLen := 20
	title := prompt
	if utf8Title := utf8string.NewString(prompt); utf8Title.RuneCount() > truncLen {
		title = utf8Title.Slice(0, truncLen) + "..."
	}
	chat := &models.Chat{
		UserID: userID,
		Title:  title,
		Text:   prompt,
	}
	if err := s.db.Create(chat).Error; err != nil {
		panic(err)
	}
	return chat
}

func (s *chatService) FindChat(userID uint, id uint) *models.Chat {
	var chat models.Chat
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&chat).Error; err != nil {
		panic(err)
	}
	return &chat
}

func (s *chatService) SetUserTelegramActiveChat(user *models.User, chat *models.Chat) {
	var id *uint = nil
	if chat != nil {
		id = &chat.ID
	}
	if err := s.db.Model(user).Update("current_telegram_chat_id", id).Error; err != nil {
		panic(err)
	}
}

func (s *chatService) DeleteChat(chat *models.Chat) {
	if err := s.db.Where("chat_id = ?", chat.ID).Delete(&models.ChatEntry{}).Error; err != nil {
		panic(err)
	}
	if err := s.db.Delete(&chat).Error; err != nil {
		panic(err)
	}
}

func (s *chatService) Summarize(chat *models.Chat, question string, onComplete func(result string)) *models.ChatEntry {
	entry := &models.ChatEntry{
		ChatID:   chat.ID,
		Question: question,
		Answer:   nil,
	}

	if err := s.db.Create(entry).Error; err != nil {
		panic(err)
	}

	s.node.Summarize(chat.Text, question, func(result string, err error) {
		var answer string
		if err != nil {
			answer = fmt.Sprintf("Failure: %s", err)
		} else {
			answer = result
		}
		if err := s.db.Model(entry).Update("answer", answer).Error; err != nil {
			panic(err)
		}
		onComplete(answer)
		// chat.Messages = append(chat.Messages, Message{
		// 	ID:   outMsgId,
		// 	From: "service",
		// 	Text: result,
		// })
	})

	return entry
}
