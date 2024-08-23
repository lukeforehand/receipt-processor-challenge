/*
handler.go contains methods for handling receipt requests
*/
package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// ReceiptHandler handles the access of stored receipts, and the calculation
// of points earned for a Receipt
type ReceiptHandler struct {
	// queue is the receipt task queue
	queue ReceiptQueue
	// database is the receipt storage
	database ReceiptDatabase
	// ruleProcessor processes a receipt to determine points earned
	ruleProcessor RuleProcessor
}

// NewReceiptHandler initializes ReceiptHandler with defaults
func NewReceiptHandler(port string) ReceiptHandler {
	if port == "" {
		port = "6379"
	}
	return ReceiptHandler{
		queue:         NewReceiptQueue(os.Getenv("RECEIPT_QUEUE"), port),
		database:      NewReceiptDatabase(os.Getenv("RECEIPT_DATABASE"), port),
		ruleProcessor: NewRuleProcessor(),
	}
}

// PostReceiptsProcess handles POST requests to process a Receipt,
// storing the Receipt along with an associated UUID
// Response example: {"id":"7d4d837b-ef5e-47c0-89a9-889657b66eb9"}
func (h *ReceiptHandler) PostReceiptsProcess(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// validate purchaseTime
	_, err := time.Parse("15:04", receipt.PurchaseTime)
	if err != nil {
		http.Error(w, "Invalid purchaseTime", http.StatusBadRequest)
		return
	}

	id := uuid.New()
	// queue task
	err = h.queue.Enqueue(id, receipt)
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
	receipt, err := h.database.Get(id)
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
