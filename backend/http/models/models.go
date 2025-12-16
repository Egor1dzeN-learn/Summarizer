package models

import "time"

type User struct {
	ID         uint   `gorm:"primaryKey"`
	TelegramID uint   `gorm:"uniqueIndex;not_null"`
	Name       string `gorm:"size:100;not_null" json:"name"`

	LoginToken           *string   `gorm:"uniqueIndex" json:"-"`
	LoginTokenExpiration time.Time `gorm:"" json:"-"`

	//
	CurrentTelegramChatID *uint
	CurrentTelegramChat   *Chat  `gorm:"foreignKey:UserID"`
	Chats                 []Chat `gorm:"foreignKey:UserID"`
}

type Chat struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	UserID uint   `gorm:"not null" json:"user_id"`
	Title  string `gorm:"type:text" json:"title"`
	Text   string `gorm:"type:text" json:"text"`

	//
	User    *User       `gorm:"constraint:OnDelete:CASCADE" json:"user"`
	Entries []ChatEntry `gorm:"foreignKey:ChatID" json:"entries"`
}

type ChatEntry struct {
	ID       uint    `gorm:"primaryKey" json:"id"`
	ChatID   uint    `gorm:"not null" json:"chat_id"`
	Question string  `gorm:"type:text;not null" json:"question"`
	Answer   *string `gorm:"type:text" json:"answer"`

	//
	Chat *Chat `gorm:"constraint:OnDelete:CASCADE"`
}
