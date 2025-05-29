package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"matchmaking-service/internal/grpc/chatpb"
	"matchmaking-service/internal/repository"
)

type MatchmakingService struct {
	redisRepo   *repository.RedisRepository
	chatSvc     chatpb.ChatServiceClient
	subscribers map[int64]chan int64
	mu          sync.Mutex
}

func NewMatchmakingService(
	redisRepo *repository.RedisRepository,
	chatSvc chatpb.ChatServiceClient,
) *MatchmakingService {
	return &MatchmakingService{
		redisRepo:   redisRepo,
		chatSvc:     chatSvc,
		subscribers: make(map[int64]chan int64),
	}
}

func (s *MatchmakingService) FindMatch(ctx context.Context, userID int64) (<-chan int64, error) {
	// Проверяем, находится ли пользователь уже в очереди или чате
	inQueue, err := s.redisRepo.IsUserInQueue(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки статуса пользователя: %w", err)
	}
	if inQueue {
		return nil, fmt.Errorf("пользователь %d уже находится в очереди", userID)
	}

	ch := make(chan int64, 1)
	s.mu.Lock()
	s.subscribers[userID] = ch
	s.mu.Unlock()

	// Проверяем, есть ли подходящий партнер
	partnerID, err := s.redisRepo.GetMatchingUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска партнера: %w", err)
	}

	if partnerID == 0 {
		// Нет подходящего партнера, добавляем в очередь с таймаутом
		if err := s.redisRepo.AddUserToQueue(ctx, userID); err != nil {
			s.mu.Lock()
			delete(s.subscribers, userID)
			close(ch)
			s.mu.Unlock()
			return nil, fmt.Errorf("ошибка добавления в очередь: %w", err)
		}

		// Запускаем таймер для таймаута ожидания
		go func() {
			select {
			case <-time.After(30 * time.Second): // Таймаут 30 секунд
				s.mu.Lock()
				if ch, ok := s.subscribers[userID]; ok {
					close(ch)
					delete(s.subscribers, userID)
				}
				s.mu.Unlock()
				s.redisRepo.RemoveUserFromQueue(ctx, userID)
			case <-ctx.Done():
				s.mu.Lock()
				if ch, ok := s.subscribers[userID]; ok {
					close(ch)
					delete(s.subscribers, userID)
				}
				s.mu.Unlock()
				s.redisRepo.RemoveUserFromQueue(ctx, userID)
			}
		}()
		return ch, nil
	}

	// Создаем чат через gRPC
	resp, err := s.chatSvc.CreateChat(ctx, &chatpb.CreateChatRequest{
		User1Id: userID,
		User2Id: partnerID,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка создания чата через gRPC: %w", err)
	}
	chatID := resp.GetChatId()

	// Уведомляем обоих пользователей
	s.mu.Lock()
	if otherCh, ok := s.subscribers[partnerID]; ok {
		otherCh <- chatID
		close(otherCh)
		delete(s.subscribers, partnerID)
	}
	ch <- chatID
	close(ch)
	delete(s.subscribers, userID)
	s.mu.Unlock()

	return ch, nil
}
