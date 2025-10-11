package services

import (
	"auth-server/storage"
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type oauthCodeService struct {
	redisClient *storage.NamespacedRedis
}

func (s *oauthCodeService) Save(code string, data string, expiresIn time.Duration) error {
	return s.redisClient.Set(context.TODO(), code, data, expiresIn)
}

func (s *oauthCodeService) Get(code string) (string, error) {
	return s.redisClient.Get(context.TODO(), code)
}

func (s *oauthCodeService) Delete(code string) error {
	return s.redisClient.Del(context.TODO(), code)
}

func NewOauthCodeService(redisClient *redis.Client) TemporaryCodeService {
	return &oauthCodeService{
		redisClient: storage.NewNamespacedRedis(redisClient, "oauth"),
	}
}
