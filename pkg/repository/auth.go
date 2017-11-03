package repository

import (
	"context"
	"time"

	"github.com/garyburd/redigo/redis"

	"github.com/acoshift/acourse/pkg/app"
)

func (repo) StoreMagicLink(ctx context.Context, linkID string, userID string) error {
	pool, prefix := app.GetRedisPool(ctx)
	db := pool.Get()
	defer db.Close()

	_, err := db.Do("SETEX", prefix+"magic:"+linkID, int64(time.Hour/time.Second), userID)
	if err != nil {
		return err
	}
	return nil
}

func (repo) FindMagicLink(ctx context.Context, linkID string) (string, error) {
	pool, prefix := app.GetRedisPool(ctx)
	db := pool.Get()
	defer db.Close()

	key := prefix + "magic:" + linkID
	userID, err := redis.String(db.Do("GET", key))
	if err != nil {
		return "", err
	}
	db.Do("DEL", key)
	return userID, nil
}
