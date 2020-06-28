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
	queue  string
	table  string
}

func NewRedisCache(address string, password string, db int, queue string, table string) *RedisCache {
	rq := &RedisCache{
		client: redis.NewClient(&redis.Options{
			Addr:     address,
			Password: password,
			DB:       db,
		}),
		ctx:   context.Background(),
		queue: queue,
		table: table,
	}

	_, err := rq.client.Ping(rq.ctx).Result()
	if err != nil {
		log.Fatalf("Error pinging redis instance %v ", err)
	}

	return rq
}

func (rq *RedisCache) Poll(key CacheKey, values chan string) {
	for {
		result, err := rq.client.BRPopLPush(rq.ctx, rq.queue, rq.table, 0).Result()

		if err != nil {
			continue
		}

		values <- string(result[0])
	}
}

func (rq *RedisCache) Get(key CacheKey) CacheValue {
	value, err := rq.client.Get(rq.ctx, string(key)).Result()

	if err != nil {
		return nil
	}

	return CacheValue(value)
}

func (rq *RedisCache) Push(key CacheKey, value CacheValue) {
	rq.client.Set(rq.ctx, string(key), string(value), 0)
	rq.client.LPush(rq.ctx, rq.queue, key)
}

func (rq *RedisCache) Keys() []CacheKey {
	values, err := rq.client.LRange(rq.ctx, rq.table, 0, -1).Result()
	if err != nil {
		return []CacheKey{}
	}

	keys := []CacheKey{}
	for _, value := range values {
		keys = append(keys, CacheKey(value))
	}

	return keys
}
