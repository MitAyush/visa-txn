package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/mitayush/visa-txn/internal/apperrors"
	"github.com/mitayush/visa-txn/internal/config"
	"github.com/mitayush/visa-txn/internal/model"
	"github.com/mitayush/visa-txn/internal/repository"
)

type TransactionService struct {
	db          *sql.DB
	txRepo      repository.TransactionRepository
	accountRepo repository.AccountRepository
	txAuditRepo repository.TransactionAuditRepository
	cfg         *config.Config
}

func NewTransactionService(
	db *sql.DB,
	txRepo repository.TransactionRepository,
	accountRepo repository.AccountRepository,
	txAuditRepo repository.TransactionAuditRepository,
	cfg *config.Config,
) *TransactionService {
	return &TransactionService{
		db:          db,
		txRepo:      txRepo,
		accountRepo: accountRepo,
		txAuditRepo: txAuditRepo,
		cfg:         cfg,
	}
}

func (s *TransactionService) CreateTransaction(ctx context.Context, tx *model.Transaction) (*model.Transaction, error) {
	if _, ok := s.cfg.OperationTypes[tx.OperationTypeID]; !ok {
		return nil, apperrors.ErrInvalidOperationType
	}
	// TODO: its not safe in distributed system, we need distributed lock to avoid race conditions
	existingTxn, err := s.txRepo.GetByIdempotencyKey(ctx, tx.IdempotencyKey)
	if err != nil {
		log.Println("error checking if transaction exists", err)
		return nil, apperrors.ErrInternalServerError
	}
	if existingTxn != nil {
		log.Println("transaction already exists with idempotency key", tx.IdempotencyKey)
		return nil, apperrors.ErrTransactionAlreadyExists
	}

	_, err = s.accountRepo.GetByAccountID(ctx, tx.AccountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrAccountDoesNotExist
		}
		log.Println("error verifying account", err)
		return nil, apperrors.ErrInternalServerError
	}

	now := time.Now()
	if tx.EventDate.IsZero() {
		tx.EventDate = now
	}

	tx.Amount = s.normalizeAmount(tx.OperationTypeID, tx.Amount)
	tx.CreatedAt = now

	request, err := json.Marshal(map[string]any{
		"account_id":        tx.AccountID,
		"operation_type_id": tx.OperationTypeID,
		"amount":            tx.Amount,
	})
	if err != nil {
		log.Println("error marshalling request", err)
		return nil, apperrors.ErrInternalServerError
	}

	var newTxn *model.Transaction
	err = repository.RunInTransaction(ctx, s.db, func(sqlTx *sql.Tx) error {
		var opErr error
		newTxn, opErr = s.txRepo.InsertTransactionWithConn(ctx, sqlTx, tx)
		if opErr != nil {
			return opErr
		}
		var response []byte
		response, opErr = json.Marshal(map[string]any{
			"transaction_id":    newTxn.ID,
			"amount":            newTxn.Amount,
			"event_date":        newTxn.EventDate,
			"created_at":        newTxn.CreatedAt,
			"idempotency_key":   newTxn.IdempotencyKey,
			"account_id":        newTxn.AccountID,
			"operation_type_id": newTxn.OperationTypeID,
		})
		if opErr != nil {
			return opErr
		}
		audit := &model.TransactionAudit{
			IdempotencyKey: tx.IdempotencyKey,
			Request:        string(request),
			Response:       string(response),
			Status:         "success",
			CreatedAt:      now,
		}
		return s.txAuditRepo.InsertTransactionAuditWithConn(ctx, sqlTx, audit)
	})
	if err != nil {
		log.Println("error persisting transaction", err)
		return nil, apperrors.ErrInternalServerError
	}
	return newTxn, nil
}

func (s *TransactionService) normalizeAmount(operationTypeID int, amount float64) float64 {
	if s.cfg.OperationTypes[operationTypeID].Sign == config.SignDebit {
		return -amount
	}
	return amount
}
