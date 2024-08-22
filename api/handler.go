/*
handler.go contains methods for handling receipt requests
*/
package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// ReceiptHandler handles the access of stored receipts, and the calculation
// of points earned for a Receipt
type ReceiptHandler struct {
	// queue is the receipt task queue
	queue RedisQueue
	// database is the receipt storage
	database Database
	// ruleProcessor processes a receipt to determine points earned
	ruleProcessor RuleProcessor
}

// NewReceiptHandler initializes ReceiptHandler with defaults
func NewReceiptHandler() ReceiptHandler {
	return ReceiptHandler{
		queue:         NewRedisQueue(os.Getenv("RECEIPT_QUEUE")),
		database:      Database{},
		ruleProcessor: NewRuleProcessor(),
	}
}

// PostReceiptsProcess handles POST requests to process a Receipt,
// storing the Receipt along with an associated UUID
// Response example: {"id":"7d4d837b-ef5e-47c0-89a9-889657b66eb9"}
func (h *ReceiptHandler) PostReceiptsProcess(w http.ResponseWriter, r *http.Request) {
	id := uuid.New()
	receipt, _ := io.ReadAll(r.Body)
	// queue task
	err := h.queue.Enqueue(id, receipt)
	if err != nil {
		log.Fatalf("Failed to queue task: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(PostReceiptsProcessResponse{
		Id: id.String(),
	})
}

// GetReceiptsIdPoints handles GET requests to get points earned for a Receipt
// Response example: {"points":31}
func (h *ReceiptHandler) GetReceiptsIdPoints(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	receipt, err := h.database.GetReceipt(id)
	if err != nil {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}
	response := GetReceiptsIdPointsResponse{
		Points: h.ruleProcessor.Points(receipt),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
