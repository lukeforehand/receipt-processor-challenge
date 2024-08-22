/*
database.go contains methods for accessing stored receipts
*/
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

// ReceiptDatabase is a simple in memory key/value store implementation
// for id/Receipt
type ReceiptDatabase struct {
	// context
	ctx context.Context
	// namespace
	namespace string
	// pointer to a redis client
	client *redis.Client
}

// NewReceiptDatabase initializes ReceiptDatabase with defaults
func NewReceiptDatabase(namespace string) ReceiptDatabase {
	return ReceiptDatabase{
		ctx:       context.Background(),
		namespace: namespace,
		client: redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:6379", os.Getenv("REDIS_HOST")),
		}),
	}
}

// Get retrieves a Receipt from the ReceiptDatabase
// id: the uuid string associated with a Receipt
// Returns: the Receipt for the given id
func (r ReceiptDatabase) Get(id string) (Receipt, error) {
	key := fmt.Sprintf("%s:%s", r.namespace, id)
	data, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return Receipt{}, err
	}
	var receipt Receipt
	err = json.Unmarshal([]byte(data), &receipt)
	return receipt, err
}

// Put stores a Receipt in the ReceiptDatabase
// id: the uuid string associated with the given Receipt
// receipt: the Receipt to store
func (r ReceiptDatabase) Set(id string, receipt Receipt) error {
	value, _ := json.Marshal(&receipt)
	key := fmt.Sprintf("%s:%s", r.namespace, id)
	return r.client.Set(r.ctx, key, value, 0).Err()
}
