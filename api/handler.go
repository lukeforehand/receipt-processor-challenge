package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ReceiptHandler struct {
	Database      Database
	RuleProcessor RuleProcessor
}

func NewReceiptHandler() ReceiptHandler {
	return ReceiptHandler{
		Database:      Database{},
		RuleProcessor: NewRuleProcessor(),
	}
}

func (h *ReceiptHandler) PostReceiptsProcess(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
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

func (h *ReceiptHandler) GetReceiptsIdPoints(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	receipt, err := h.Database.GetReceipt(id)
	if err != nil {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}
	response := GetReceiptsIdPointsResponse{
		Points: h.RuleProcessor.TotalPoints(receipt),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
