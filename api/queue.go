package api

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

type ReceiptQueue struct {
	// context
	ctx context.Context
	// name of the queue
	name string
	// pointer to a redis client used as a broker
	client *redis.Client
}

func (r ReceiptQueue) Enqueue(values ...string) error {
	value, _ := json.Marshal(values)
	return r.client.LPush(r.ctx, r.name, value).Err()
}

func (r ReceiptQueue) Dequeue() ([]string, error) {
	data, err := r.client.BRPop(r.ctx, 0, r.name).Result()
	if err != nil {
		return []string{}, err
	}
	var values []string
	err = json.Unmarshal([]byte(data[1]), &values)
	return values, err
}

func NewReceiptQueue(name string) ReceiptQueue {
	return ReceiptQueue{
		ctx:  context.Background(),
		name: name,
		client: redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:6379", os.Getenv("REDIS_HOST")),
		}),
	}
}
