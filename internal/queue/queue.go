package queue

import (
	"context"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

struct Queue {
	client *redis.Client
	ctx context.Context
}

func NewInstance(add string, password: string, db: string) {

}


func (q *Queue) IsConnected() bool {
	_, err := rdb.Ping(q.ctx).Result()

	if err != nil {
		return true
	}

	return false
}