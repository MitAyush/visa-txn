package dto

import (
	"errors"
	"time"

	"github.com/mitayush/visa-txn/internal/model"
)

type CreateTransactionRequest struct {
	AccountID       int64     `json:"account_id"`
	OperationTypeID int       `json:"operation_type_id"`
	Amount          float64   `json:"amount"`
	EventDate       time.Time `json:"event_date"`
}

type TransactionResponse struct {
	ID              int64     `json:"id,omitempty"`
	AccountID       int64     `json:"account_id"`
	OperationTypeID int       `json:"operation_type_id"`
	Amount          float64   `json:"amount"`
	EventDate       time.Time `json:"event_date"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
}

func CreateTransactionRequestToModel(req *CreateTransactionRequest) *model.Transaction {
	return &model.Transaction{
		AccountID:       req.AccountID,
		OperationTypeID: req.OperationTypeID,
		Amount:          req.Amount,
		EventDate:       req.EventDate,
	}
}

func TransactionToResponse(t *model.Transaction) TransactionResponse {
	return TransactionResponse{
		ID:              t.ID,
		AccountID:       t.AccountID,
		OperationTypeID: t.OperationTypeID,
		Amount:          t.Amount,
		EventDate:       t.EventDate,
		CreatedAt:       t.CreatedAt,
	}
}

func (r *CreateTransactionRequest) Validate() error {
	if r.AccountID <= 0 {
		return errors.New("account_id is required")
	}
	if r.OperationTypeID <= 0 {
		return errors.New("operation_type_id is required")
	}
	if r.Amount <= 0 {
		return errors.New("amount is required")
	}
	return nil
}
