package models

import "time"

// Chat - структура для хранения информации о чате
type Chat struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	User1ID   int64     `gorm:"not null" json:"user1_id"`
	User2ID   int64     `gorm:"not null" json:"user2_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
