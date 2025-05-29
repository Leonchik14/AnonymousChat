package repository

import (
	"context"
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

func (r *RedisRepository) RemoveUserFromQueue(ctx context.Context, userID int64) error {
	// Удаляем пользователя из очереди
	err := r.client.LRem(ctx, "matchmaking_queue", 0, strconv.FormatInt(userID, 10)).Err()
	if err != nil {
		return fmt.Errorf("ошибка удаления из очереди: %w", err)
	}
	log.Printf("🔹 Пользователь %d удален из очереди", userID)
	return nil
}

func (r *RedisRepository) IsUserInQueue(ctx context.Context, userID int64) (bool, error) {
	// Проверяем, есть ли пользователь в очереди
	queue, err := r.client.LRange(ctx, "matchmaking_queue", 0, -1).Result()
	if err != nil {
		return false, fmt.Errorf("ошибка проверки очереди: %w", err)
	}
	userIDStr := strconv.FormatInt(userID, 10)
	for _, id := range queue {
		if id == userIDStr {
			return true, nil
		}
	}
	return false, nil
}

// Lua-скрипт для атомарного получения и удаления пользователя из очереди
const getMatchingUserScript = `
local queue_key = KEYS[1]
local user_id = ARGV[1]
local user = redis.call('LPOP', queue_key)
if user == false then
    return 0
end
if user == user_id then
    local next_user = redis.call('LPOP', queue_key)
    if next_user then
        redis.call('RPUSH', queue_key, user)
    end
    if next_user == false then
        return 0
    end
    return next_user
end
return user
`

func (r *RedisRepository) GetMatchingUser(ctx context.Context, userID int64) (int64, error) {
	// Используем Lua-скрипт для атомарного получения пользователя
	result, err := r.client.Eval(ctx, getMatchingUserScript, []string{"matchmaking_queue"}, strconv.FormatInt(userID, 10)).Result()
	if err != nil {
		return 0, fmt.Errorf("ошибка выполнения Lua-скрипта: %w", err)
	}

	switch v := result.(type) {
	case int64:
		if v == 0 {
			return 0, nil
		}
		return v, nil
	case string:
		matchUserID, _ := strconv.ParseInt(v, 10, 64)
		log.Printf("✅ Найден пользователь из очереди: %d", matchUserID)
		return matchUserID, nil
	default:
		return 0, fmt.Errorf("неожиданный тип результата Lua-скрипта: %T", v)
	}
}
