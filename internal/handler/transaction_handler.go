package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/mitayush/visa-txn/internal/apperrors"
	"github.com/mitayush/visa-txn/internal/dto"
	"github.com/mitayush/visa-txn/internal/service"
)

var transactionCreateErrHTTPStatus = map[error]int{
	apperrors.ErrInvalidOperationType:     http.StatusBadRequest,
	apperrors.ErrAccountDoesNotExist:      http.StatusNotFound,
	apperrors.ErrTransactionAlreadyExists: http.StatusConflict,
}

type TransactionHandler struct {
	svc *service.TransactionService
}

func NewTransactionHandler(svc *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{svc: svc}
}

func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	tx := dto.CreateTransactionRequestToModel(&req)
	tx.IdempotencyKey = getIdempotencyKey(r)

	newTxn, err := h.svc.CreateTransaction(r.Context(), tx)
	if err != nil {
		writeJSON(w, transactionCreateErrHTTPStatus[err], dto.ErrorResponse{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, dto.TransactionToResponse(newTxn))
}

func getIdempotencyKey(r *http.Request) string {
	key := r.Header.Get("X-Idempotency-Key")
	if key == "" {
		key = uuid.NewString()
	}
	return key
}
