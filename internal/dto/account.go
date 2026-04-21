package dto

import (
	"errors"
	"strings"
	"time"

	"github.com/mitayush/visa-txn/internal/model"
)

type CreateAccountRequest struct {
	DocumentNumber string `json:"document_number"`
}

type GetAccountResponse struct {
	AccountID      int64     `json:"account_id"`
	DocumentNumber string    `json:"document_number"`
	CreatedAt      time.Time `json:"created_at"`
}

type CreateAccountResponse struct {
	AccountID      int64     `json:"account_id"`
	DocumentNumber string    `json:"document_number"`
	CreatedAt      time.Time `json:"created_at"`
}

func CreateAccountRequestToModel(req *CreateAccountRequest) *model.Account {
	return &model.Account{DocumentNumber: req.DocumentNumber}
}

func AccountToResponse(a *model.Account) GetAccountResponse {
	return GetAccountResponse{
		AccountID:      a.AccountID,
		DocumentNumber: a.DocumentNumber,
		CreatedAt:      a.CreatedAt,
	}
}

func (r *CreateAccountRequest) Validate() error {
	r.DocumentNumber = strings.TrimSpace(r.DocumentNumber)
	if r.DocumentNumber == "" {
		return errors.New("document_number is required")
	}
	if len(r.DocumentNumber) > 64 {
		return errors.New("document_number must be at most 64 characters")
	}
	return nil
}
