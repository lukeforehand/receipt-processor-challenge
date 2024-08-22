package backend

import (
	"log"
	"os"
	"receiptprocessor/api"
)

// ReceiptProcessor pops a UUID and Receipt from the task broker
type ReceiptProcessor struct {
	queue api.RedisQueue
	// TODO: Database
}

// NewReceiptProcessor initializes NewReceiptProcessor with defaults
func NewReceiptProcessor() ReceiptProcessor {
	return ReceiptProcessor{
		queue: api.NewRedisQueue(os.Getenv("RECEIPT_QUEUE")),
	}
}

// Start starts the ReceiptProcessor worker
func (p *ReceiptProcessor) Start() {
	for {
		tasks, err := p.queue.Dequeue()
		if len(tasks) > 1 {
			if err != nil {
				log.Fatalf("Error retrieving tasks: %v", err)
			}
			err = p.processTask(tasks[1])
			if err != nil {
				log.Fatalf("Failed to process task: %v", err)
			}
		}
	}
}

// processTask processes each receipt in the queue
func (t *ReceiptProcessor) processTask(task string) error {
	// TODO: store receipt in database
	log.Printf("Processing task: %v", task)
	return nil
}
