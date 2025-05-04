package models

import "time"

type WsMessage struct {
	SenderID  int64     `json:"senderId"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}
