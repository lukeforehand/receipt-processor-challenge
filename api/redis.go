package api

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

type RedisQueue struct {
	// context
	ctx context.Context
	// name of the queue
	name string
	// pointer to a redis client used as a broker
	client *redis.Client
}

func (r RedisQueue) Enqueue(values ...interface{}) error {
	return r.client.LPush(r.ctx, r.name, values).Err()
}

func (r RedisQueue) Dequeue() ([]string, error) {
	return r.client.BRPop(r.ctx, 0, r.name).Result()
}

func NewRedisQueue(name string) RedisQueue {
	return RedisQueue{
		ctx:    context.Background(),
		name:   name,
		client: NewRedisClient(),
	}
}

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:6379", os.Getenv("REDIS_HOST")),
	})
}
