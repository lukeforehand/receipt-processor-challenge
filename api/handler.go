/*
handler.go contains methods for handling receipt requests
*/
package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// ReceiptHandler handles the access of stored receipts, and the calculation
// of points earned for a Receipt
type ReceiptHandler struct {
	// Database is the receipt storage
	Database Database
	// RuleProcessor processes a receipt to determine points earned
	RuleProcessor RuleProcessor
}

// NewReceiptHandler initializes ReceiptHandler with defaults
func NewReceiptHandler() ReceiptHandler {
	return ReceiptHandler{
		Database:      Database{},
		RuleProcessor: NewRuleProcessor(),
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
	h.Database.PutReceipt(id.String(), receipt)
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
	receipt, err := h.Database.GetReceipt(id)
	if err != nil {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}
	response := GetReceiptsIdPointsResponse{
		Points: h.RuleProcessor.Points(receipt),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
