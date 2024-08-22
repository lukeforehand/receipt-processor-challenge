/*
handler.go contains methods for handling receipt requests
*/
package api

import (
	"context"
	"encoding/json"
	"net/http"
	"receiptprocessor/backend"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

// ReceiptHandler handles the access of stored receipts, and the calculation
// of points earned for a Receipt
type ReceiptHandler struct {
	// pointer to a redis client used as a backend broker
	redis *redis.Client
	// database is the receipt storage
	database Database
	// ruleProcessor processes a receipt to determine points earned
	ruleProcessor RuleProcessor
}

// NewReceiptHandler initializes ReceiptHandler with defaults
func NewReceiptHandler() ReceiptHandler {
	return ReceiptHandler{
		redis:         backend.NewRedisClient(),
		database:      Database{},
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
	id := uuid.New()

	// queue task
	err := h.redis.LPush(ctx, backend.RECEIPT_TASK_QUEUE, id, receipt).Err()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := PostReceiptsProcessResponse{
		Id: id.String(),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
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
