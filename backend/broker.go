package backend

import (
	"github.com/redis/go-redis/v9"
)

const RECEIPT_TASK_QUEUE = "receipt_task_queue"

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}
