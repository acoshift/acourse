package app

import (
	"cloud.google.com/go/storage"
	"github.com/acoshift/go-firebase-admin"
	"github.com/go-redis/redis"

	"github.com/acoshift/acourse/file"
	"github.com/acoshift/acourse/image"
	"github.com/acoshift/acourse/notify"
)

// Config use to init app package
type Config struct {
	BaseURL            string
	RedisClient        *redis.Client
	RedisPrefix        string
	Auth               *firebase.Auth
	AdminNotifier      notify.AdminNotifier
	BucketHandle       *storage.BucketHandle
	BucketName         string
	ImageResizeEncoder image.JPEGResizeEncoder
	FileStorage        file.Storage
}
