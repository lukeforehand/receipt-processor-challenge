/*
task.go contains methods for processing receipts
*/
package backend

import (
	"encoding/json"
	"log"
	"os"
	"receiptprocessor/api"
)

// ReceiptProcessor pops a UUID and Receipt from the task broker
type ReceiptProcessor struct {
	// queue is the receipt task queue
	queue api.ReceiptQueue
	// database is the receipt storage
	database api.ReceiptDatabase
}

// NewReceiptProcessor initializes NewReceiptProcessor with defaults
func NewReceiptProcessor() ReceiptProcessor {
	return ReceiptProcessor{
		queue:    api.NewReceiptQueue(os.Getenv("RECEIPT_QUEUE")),
		database: api.NewReceiptDatabase(os.Getenv("RECEIPT_DATABASE")),
	}
}

// Start starts the ReceiptProcessor worker
func (r *ReceiptProcessor) Start() {
	for {
		task, err := r.queue.Dequeue()
		if err != nil {
			log.Fatalf("Error retrieving tasks: %v", err)
		}
		err = r.processTask(task)
		if err != nil {
			log.Fatalf("Failed to process task: %v", err)
		}
	}
}

// processTask processes each receipt in the queue
func (r *ReceiptProcessor) processTask(task []string) error {
	log.Printf("Processing task: %v", task)
	id := task[0]
	var receipt api.Receipt
	json.Unmarshal([]byte(task[1]), &receipt)
	return r.database.Set(id, receipt)
}
