package queue

import (
	"context"
	"log"

	"github.com/go-redis/redis"
)

type RedisQueue struct {
	ctx    context.Context
	client *redis.Client
}

func NewRedisQueue(address string, password string, db int) *RedisQueue {
	rq := &RedisQueue{
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

func (rq *RedisQueue) Poll(key QueueKey) QueueValue {

}

// This is non-blocking.
func (rq *RedisQueue) Get(key QueueKey) QueueValue {

}

func (rq *RedisQueue) Push(key QueueKey, value QueueValue) {

}
