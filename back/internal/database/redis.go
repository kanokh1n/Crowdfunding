package database

import (
	"log"
	"context"

	"github.com/redis/go-redis/v9"
)

func NewRedis(addr string) *redis.Client {
	opt, err := redis.ParseURL("redis://" + addr)
	if err != nil {
		// addr might already be a full URL
		opt, err = redis.ParseURL(addr)
		if err != nil {
			log.Fatalf("failed to parse redis url: %v", err)
		}
	}

	rdb := redis.NewClient(opt)

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	return rdb
}
