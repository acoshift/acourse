package repository

import (
	"context"
	"time"

	"github.com/go-redis/redis"

	"github.com/acoshift/acourse/context/redisctx"
	"github.com/acoshift/acourse/entity"
)

// StoreMagicLink stores magic link to redis
func StoreMagicLink(ctx context.Context, linkID string, userID string) error {
	c := redisctx.GetClient(ctx)
	prefix := redisctx.GetPrefix(ctx)

	err := c.Set(prefix+"magic:"+linkID, userID, time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

// FindMagicLink finds magic link from redis
func FindMagicLink(ctx context.Context, linkID string) (string, error) {
	c := redisctx.GetClient(ctx)
	prefix := redisctx.GetPrefix(ctx)

	key := prefix + "magic:" + linkID
	userID, err := c.Get(key).Result()
	if err == redis.Nil {
		return "", entity.ErrNotFound
	}
	if err != nil {
		return "", err
	}
	c.Del(key)
	return userID, nil
}

// CanAcquireMagicLink checks rate limit to acquire magic link
func CanAcquireMagicLink(ctx context.Context, email string) (bool, error) {
	c := redisctx.GetClient(ctx)
	prefix := redisctx.GetPrefix(ctx)

	key := prefix + "magic-rate:" + email
	current, err := c.Incr(key).Result()
	if err != nil {
		return false, err
	}
	if current > 1 {
		return false, nil
	}
	err = c.Expire(key, 5*time.Minute).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}
