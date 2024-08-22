package backend

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

// ReceiptProcessor pops a UUID and Receipt from the task broker
type ReceiptProcessor struct {
	// pointer to a redis client used as a backend broker
	redis *redis.Client
	// TODO :Database
}

func NewReceiptProcessor() ReceiptProcessor {
	return ReceiptProcessor{
		redis: NewRedisClient(),
	}
}

func (t *ReceiptProcessor) Start() {
	for {
		tasks, err := t.redis.BRPop(ctx, 0, RECEIPT_TASK_QUEUE).Result()
		if err != nil {
			log.Fatalf("Error retrieving tasks: %v", err)
		}
		err = t.processTask(tasks[1])
		if err != nil {
			log.Fatalf("Failed to process task: %v", err)
		}
	}
}

func (t *ReceiptProcessor) processTask(task string) error {
	// TODO: store receipt in database
	return nil
}
