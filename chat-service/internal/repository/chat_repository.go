package repository

import (
	"context"
	"fmt"
	"log"

	"chat-service/pkg/models"
	"gorm.io/gorm"
)

// ChatRepository - репозиторий работы с БД
type ChatRepository struct {
	db *gorm.DB
}

// NewChatRepository - конструктор репозитория
func NewChatRepository(db *gorm.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

// CreateChat - создаёт новый чат между двумя пользователями
func (r *ChatRepository) CreateChat(ctx context.Context, user1ID, user2ID int64) (int64, error) {
	chat := models.Chat{User1ID: user1ID, User2ID: user2ID}

	result := r.db.WithContext(ctx).Create(&chat)
	if result.Error != nil {
		return 0, fmt.Errorf("ошибка создания чата: %w", result.Error)
	}

	log.Printf("✅ Чат создан: ID %d (пользователи: %d, %d)", chat.ID, user1ID, user2ID)
	return chat.ID, nil
}

// SaveMessage - сохраняет сообщение в базе данных
func (r *ChatRepository) SaveMessage(ctx context.Context, message *models.Message) error {
	result := r.db.WithContext(ctx).Create(message)
	if result.Error != nil {
		return fmt.Errorf("ошибка сохранения сообщения: %w", result.Error)
	}

	log.Printf("📩 Сообщение сохранено (чат %d): %s", message.ChatID, message.Content)
	return nil
}

// GetChatHistory - получает историю сообщений для указанного чата
func (r *ChatRepository) GetChatHistory(ctx context.Context, chatID int64) ([]models.Message, error) {
	var messages []models.Message

	result := r.db.WithContext(ctx).
		Where("chat_id = ?", chatID).
		Order("created_at ASC").
		Find(&messages)

	if result.Error != nil {
		return nil, fmt.Errorf("ошибка загрузки истории сообщений: %w", result.Error)
	}

	log.Printf("📜 Загружена история сообщений для чата %d", chatID)
	return messages, nil
}
