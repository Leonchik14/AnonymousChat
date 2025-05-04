package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{client: client}
}

func (r *RedisRepository) AddUserToQueue(ctx context.Context, userID int64) error {
	err := r.client.RPush(ctx, "matchmaking_queue", strconv.FormatInt(userID, 10)).Err()
	if err != nil {
		return fmt.Errorf("ошибка добавления в очередь: %w", err)
	}
	log.Printf("🔹 Пользователь %d добавлен в очередь", userID)
	return nil
}

func (r *RedisRepository) GetMatchingUser(ctx context.Context) (int64, error) {
	existingUser, err := r.client.LPop(ctx, "matchmaking_queue").Result()
	if errors.Is(err, redis.Nil) {
		return 0, nil
	} else if err != nil {
		return 0, fmt.Errorf("ошибка работы с Redis: %w", err)
	}

	matchUserID, _ := strconv.ParseInt(existingUser, 10, 64)
	log.Printf("✅ Найден пользователь из очереди: %d", matchUserID)
	return matchUserID, nil
}
