package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	prefixKey = "sms:code:"
	codeTTL   = 5 * time.Minute
)

var ErrCodeNotFound = errors.New("code not found")

type SMSStorage struct {
	client *redis.Client
}

func NewSMSStorage(client *redis.Client) *SMSStorage {
	return &SMSStorage{client: client}
}

func (s *SMSStorage) Save(ctx context.Context, phone, code string) error {
	key := prefixKey + phone
	return s.client.Set(ctx, key, code, codeTTL).Err()
}

func (s *SMSStorage) Validate(ctx context.Context, phone, code string) (bool, error) {
	key := prefixKey + phone

	storedCode, err := s.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	if storedCode != code {
		return false, nil
	}

	s.client.Del(ctx, key)

	return true, nil
}
