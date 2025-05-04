package service

import (
	"context"
	"fmt"
	"sync"

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
	// 1) Подготовить канал и сохранить его
	ch := make(chan int64, 1)
	s.mu.Lock()
	s.subscribers[userID] = ch
	s.mu.Unlock()

	// 2) Попробовать достать из очереди другого юзера
	partnerID, err := s.redisRepo.GetMatchingUser(ctx)
	if err != nil {
		return nil, err
	}

	// 3) Если никого нет или достали самих себя — кладём в очередь и ждём уведомления
	if partnerID == 0 || partnerID == userID {
		if err := s.redisRepo.AddUserToQueue(ctx, userID); err != nil {
			return nil, err
		}
		return ch, nil
	}

	// 4) Нашли настоящего партнёра — создаём чат по gRPC
	resp, err := s.chatSvc.CreateChat(ctx, &chatpb.CreateChatRequest{
		User1Id: userID,
		User2Id: partnerID,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка создания чата через gRPC: %w", err)
	}
	chatID := resp.GetChatId()

	// 5) Уведомляем обоих: сначала партнёра, потом себя
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
