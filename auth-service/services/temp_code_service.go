package services

import (
	"auth-server/storage"
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type tempCodeService struct {
	redisClient *storage.NamespacedRedis
}

func (s *tempCodeService) Save(code string, data string, expiresIn time.Duration) error {
	return s.redisClient.Set(context.TODO(), code, data, expiresIn)
}

func (s *tempCodeService) Get(code string) (string, error) {
	return s.redisClient.Get(context.TODO(), code)
}

func (s *tempCodeService) Delete(code string) error {
	return s.redisClient.Del(context.TODO(), code)
}

func NewTemporaryCodeService(redisClient *redis.Client, prefix string) TemporaryCodeService {
	return &tempCodeService{
		redisClient: storage.NewNamespacedRedis(redisClient, prefix),
	}
}
