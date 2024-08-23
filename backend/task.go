/*
task.go contains methods for processing receipts
*/
package backend

import (
	"log"
	"os"
	"receiptprocessor/api"

	"github.com/google/uuid"
)

// ReceiptProcessor pops a UUID and Receipt from the task broker
type ReceiptProcessor struct {
	// queue is the receipt task queue
	queue api.ReceiptQueue
	// database is the receipt storage
	database api.ReceiptDatabase
}

// NewReceiptProcessor initializes NewReceiptProcessor with defaults
func NewReceiptProcessor(port string) ReceiptProcessor {
	if port == "" {
		port = "6379"
	}
	return ReceiptProcessor{
		queue:    api.NewReceiptQueue(os.Getenv("RECEIPT_QUEUE"), port),
		database: api.NewReceiptDatabase(os.Getenv("RECEIPT_DATABASE"), port),
	}
}

// Start the ReceiptProcessor worker
func (r *ReceiptProcessor) Start() {
	for {
		id, receipt, err := r.queue.Dequeue()
		if err != nil {
			log.Fatalf("Error retrieving tasks: %v", err)
		}
		err = r.processTask(id, receipt)
		if err != nil {
			log.Fatalf("Failed to process task: %v", err)
		}
	}
}

// processTask processes each receipt in the queue
func (r *ReceiptProcessor) processTask(id uuid.UUID, receipt api.Receipt) error {
	log.Printf("Processing task: %s %s", id, receipt)
	return r.database.Set(id.String(), receipt)
}
