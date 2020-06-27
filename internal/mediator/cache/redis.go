package cache

import (
	"context"
	"log"

	"github.com/go-redis/redis"
)

// Timeout is measured in nanoseconds
const QUEUE_TIMEOUT = 3.154e+16

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache(address string, password string, db int) *RedisCache {
	rq := &RedisCache{
		client: redis.NewClient(&redis.Options{
			Addr:     address,
			Password: password,
			DB:       db,
		}),
		ctx: context.Background(),
	}

	_, err := rq.client.Ping(rq.ctx).Result()
	if err != nil {
		log.Fatalf("Error pinging redis instance %v ", err)
	}

	return rq
}

func (rq *RedisCache) Poll(key CacheKey, values chan string) {
	for {
		val, err := rq.client.BLPop(rq.ctx, QUEUE_TIMEOUT, string(key)).Result()

		if err != nil {
			log.Fatalf(err)
		}

		c <- val
	}
}

// This is non-blocking.
func (rq *RedisCache) Get(key CacheKey) CacheValue {

}

func (rq *RedisCache) Push(key CacheKey, value CacheValue) {

}
