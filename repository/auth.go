package repository

import (
	"time"

	"github.com/go-redis/redis"

	"github.com/acoshift/acourse/entity"
)

// StoreMagicLink stores magic link to redis
func StoreMagicLink(db *redis.Client, prefix string, linkID string, userID string) error {
	err := db.Set(prefix+"magic:"+linkID, userID, time.Hour/time.Second).Err()
	if err != nil {
		return err
	}
	return nil
}

// FindMagicLink finds magic link from redis
func FindMagicLink(db *redis.Client, prefix string, linkID string) (string, error) {
	key := prefix + "magic:" + linkID
	userID, err := db.Get(key).Result()
	if err == redis.Nil {
		return "", entity.ErrNotFound
	}
	if err != nil {
		return "", err
	}
	db.Del(key)
	return userID, nil
}

// CanAcquireMagicLink checks rate limit to acquire magic link
func CanAcquireMagicLink(db *redis.Client, prefix string, email string) (bool, error) {
	key := prefix + "magic-rate:" + email
	current, err := db.Incr(key).Result()
	if err != nil {
		return false, err
	}
	if current > 1 {
		return false, nil
	}
	err = db.Expire(key, 5*time.Minute/time.Second).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}
