/*
queue.go contains methods for queuing tasks
*/
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

// ReceiptQueue is a simple in memory key/value queue implementation
// for id/Receipt
type ReceiptQueue struct {
	// context
	ctx context.Context
	// name of the queue
	name string
	// pointer to a redis client used as a broker
	client *redis.Client
}

// NewReceiptQueue initializes ReceiptQueue with defaults
func NewReceiptQueue(name string) ReceiptQueue {
	return ReceiptQueue{
		ctx:  context.Background(),
		name: name,
		client: redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:6379", os.Getenv("REDIS_HOST")),
		}),
	}
}

// Enqueue puts a id/Receipt to the back of the queue
// id: the uuid string associated with a Receipt
// Returns: the Receipt for the given id
func (r ReceiptQueue) Enqueue(id string, receipt Receipt) error {
	value, _ := json.Marshal(&receipt)
	return r.client.LPush(r.ctx, r.name, value).Err()
}

// Dequeue gets a id/Receipt from the front of the queue
// Returns: slice of id/Receipt serialized as a string
func (r ReceiptQueue) Dequeue() ([]string, error) {
	data, err := r.client.BRPop(r.ctx, 0, r.name).Result()
	if err != nil {
		return []string{}, err
	}
	var values []string
	err = json.Unmarshal([]byte(data[1]), &values)
	return values, err
}
