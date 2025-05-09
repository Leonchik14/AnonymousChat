package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
}

// Конструктор для Redis-репозитория
func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{client: client}
}

// SetSession - сохраняет Refresh Token в Redis с TTL
func (r *RedisRepository) SetSession(ctx context.Context, sessionID string, userID uint, expiration time.Duration) error {
	return r.client.Set(ctx, "session:"+sessionID, userID, expiration).Err()
}

func (r *RedisRepository) SetEmailVerification(ctx context.Context, email string, token string, expiration time.Duration) error {
	return r.client.Set(ctx, email, token, expiration).Err()
}

func (r *RedisRepository) GetEmailVerification(ctx context.Context, email string) (time.Duration, error) {
	return r.client.TTL(ctx, email).Result()
}

// GetSessionTTL - получает оставшееся время жизни Refresh Token
func (r *RedisRepository) GetSessionTTL(ctx context.Context, sessionID string) (time.Duration, error) {
	return r.client.TTL(ctx, "session:"+sessionID).Result()
}

func (r *RedisRepository) GetUserIDBySession(ctx context.Context, sessionID string) (int64, error) {
	userID, err := r.client.Get(ctx, "session:"+sessionID).Int64()
	if err != nil {
		return 0, err
	}
	return userID, nil
}

// Удаление сессии пользователя (logout)
func (r *RedisRepository) DeleteSession(ctx context.Context, sessionID string) error {
	return r.client.Del(ctx, "session:"+sessionID).Err()
}
