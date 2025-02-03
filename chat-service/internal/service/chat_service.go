package service

import (
	"context"
	"log"

	"chat-service/internal/grpc/chatpb"
	"chat-service/internal/repository"
)

// ChatService - сервис управления чатами
type ChatService struct {
	chatRepo *repository.ChatRepository
}

// NewChatService - конструктор сервиса
func NewChatService(chatRepo *repository.ChatRepository) *ChatService {
	return &ChatService{chatRepo: chatRepo}
}

// CreateChat - gRPC-метод создания чата
func (s *ChatService) CreateChat(ctx context.Context, req *chatpb.CreateChatRequest) (*chatpb.CreateChatResponse, error) {
	chatID, err := s.chatRepo.CreateChat(ctx, req.User1Id, req.User2Id)
	if err != nil {
		return nil, err
	}

	log.Printf("✅ Чат создан: %d (пользователи: %d, %d)", chatID, req.User1Id, req.User2Id)

	return &chatpb.CreateChatResponse{ChatId: chatID}, nil
}
