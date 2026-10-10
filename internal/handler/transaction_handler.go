package handler

import (
	"encoding/json"
	"net/http"

	"go-payment-gateway/internal/service"
)

// TransactionHandler handles HTTP requests for transaction operations.
type TransactionHandler struct {
	// Service performs transaction business operations.
	Service service.TransactionService
}

// NewTransactionHandler constructs a handler backed by s.
func NewTransactionHandler(s service.TransactionService) *TransactionHandler {
	return &TransactionHandler{Service: s}
}

// CreateTransaction decodes a transaction request, creates it, and writes the
// resulting transaction as JSON.
func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req service.CreateTransactionDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Format JSON payload tidak valid", http.StatusBadRequest)
		return
	}

	trx, err := h.Service.CreateTransaction(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Transaction created successfully",
		"data":    trx,
	})
}
