package service

import (
	"context"
	"fmt"

	"matchmaking-service/internal/grpc/chatpb"
	"matchmaking-service/internal/repository"
)

// MatchmakingService - сервис поиска пары
type MatchmakingService struct {
	redisRepo *repository.RedisRepository
	chatSvc   chatpb.ChatServiceClient
}

// NewMatchmakingService - конструктор сервиса
func NewMatchmakingService(redisRepo *repository.RedisRepository, chatSvc chatpb.ChatServiceClient) *MatchmakingService {
	return &MatchmakingService{redisRepo: redisRepo, chatSvc: chatSvc}
}

// FindMatch - поиск собеседника и передача чата в chat-service
func (s *MatchmakingService) FindMatch(ctx context.Context, userID int64) error {
	// Получаем собеседника из очереди
	matchUserID, err := s.redisRepo.GetMatchingUser(ctx)
	if err != nil {
		return err
	}

	// Если нет собеседника, добавляем текущего пользователя в очередь
	if matchUserID == 0 {
		return s.redisRepo.AddUserToQueue(ctx, userID)
	}

	// Отправляем gRPC-запрос в chat-service
	chatReq := &chatpb.CreateChatRequest{
		User1Id: userID,
		User2Id: matchUserID,
	}
	_, err = s.chatSvc.CreateChat(ctx, chatReq)
	if err != nil {
		return fmt.Errorf("ошибка создания чата через gRPC: %w", err)
	}

	return nil
}
