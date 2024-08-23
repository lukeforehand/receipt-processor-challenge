/*
queue.go contains methods for queuing tasks
*/
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
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
func NewReceiptQueue(name string, port string) ReceiptQueue {
	return ReceiptQueue{
		ctx:  context.Background(),
		name: name,
		client: redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), port),
		}),
	}
}

// ReceiptMessage for the queue
type ReceiptMessage struct {
	// UUID
	Id uuid.UUID
	// Receipt
	Receipt Receipt
}

// Enqueue puts a id/Receipt to the back of the queue
// id: the uuid associated with a Receipt
// Returns: the Receipt for the given id
func (r ReceiptQueue) Enqueue(id uuid.UUID, receipt Receipt) error {
	value, _ := json.Marshal(ReceiptMessage{
		Id:      id,
		Receipt: receipt,
	})
	return r.client.LPush(r.ctx, r.name, value).Err()
}

// Dequeue gets a id/Receipt from the front of the queue
// Returns: id, receipt, error
func (r ReceiptQueue) Dequeue() (uuid.UUID, Receipt, error) {
	data, err := r.client.BRPop(r.ctx, 0, r.name).Result()
	if err != nil {
		return uuid.UUID{}, Receipt{}, err
	}
	var message ReceiptMessage
	err = json.Unmarshal([]byte(data[1]), &message)
	return message.Id, message.Receipt, err
}
