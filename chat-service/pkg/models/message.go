package models

import "time"

// Message - структура для хранения сообщений в чате
type Message struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ChatID    int64     `gorm:"not null;index" json:"chat_id"`
	SenderID  int64     `gorm:"not null" json:"sender_id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
