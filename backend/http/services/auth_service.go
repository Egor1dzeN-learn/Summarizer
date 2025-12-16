package services

import (
	"errors"
	"summarizer/backend/http/models"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type AuthService interface {
	GetUser(id uint) *models.User
	GetUserByTelegramID(tid uint) *models.User
	CreateUser(tid uint, name string) *models.User
	FindUserByLoginToken(token string) *models.User
	IssueLoginToken(user *models.User) string
}

type authService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) AuthService {
	return &authService{db: db}
}

func (s *authService) GetUser(id uint) *models.User {
	var user models.User
	err := s.db.First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return &user
}

func (s *authService) GetUserByTelegramID(tid uint) *models.User {
	var user models.User
	err := s.db.Where("telegram_id = ?", tid).
		Preload("CurrentTelegramChat"). // todo: very slow, PoC
		Preload("CurrentTelegramChat.Entries").
		First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return &user
}

func (s *authService) CreateUser(tid uint, name string) *models.User {
	if existing := s.GetUserByTelegramID(tid); existing != nil {
		return existing
	}

	user := &models.User{
		TelegramID: tid,
		Name:       name,
	}
	if err := s.db.Create(user).Error; err != nil {
		panic(err)
	}
	return user
}

func (s *authService) FindUserByLoginToken(token string) *models.User {
	var user models.User
	// TODO: AND NOW() < timestamp+15m
	err := s.db.Where("login_token = ?", token).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err := s.db.Model(&user).Update("login_token", nil).Error; err != nil {
		panic(err)
	}
	return &user
}

func (s *authService) IssueLoginToken(user *models.User) string {
	token := uuid.New().String()
	if err := s.db.Model(user).Updates(map[string]interface{}{
		"login_token":            token,
		"login_token_expiration": time.Now().Add(15 * time.Minute),
	}).Error; err != nil {
		panic(err)
	}
	return token
}
