package repository

import (
	"context"
	"database/sql"

	"github.com/mitayush/visa-txn/internal/model"
)

type TransactionRepository interface {
	InsertTransactionWithConn(ctx context.Context, conn DBTX, tx *model.Transaction) (*model.Transaction, error)
	GetByIdempotencyKey(ctx context.Context, idempotencyKey string) (*model.Transaction, error)
}

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) InsertTransactionWithConn(ctx context.Context, conn DBTX, tx *model.Transaction) (*model.Transaction, error) {
	query := `
	INSERT INTO transactions (idempotency_key, account_id, operation_type_id, amount, event_date, created_at)
	VALUES (?, ?, ?, ?, ?, ?)
	RETURNING id, idempotency_key, account_id, operation_type_id, amount, event_date, created_at
	`
	row := conn.QueryRowContext(ctx, query, tx.IdempotencyKey, tx.AccountID, tx.OperationTypeID, tx.Amount, tx.EventDate, tx.CreatedAt)
	var transaction model.Transaction
	err := row.Scan(&transaction.ID, &transaction.IdempotencyKey, &transaction.AccountID, &transaction.OperationTypeID, &transaction.Amount, &transaction.EventDate, &transaction.CreatedAt)
	return &transaction, err
}

func (r *transactionRepository) GetByIdempotencyKey(ctx context.Context, idempotencyKey string) (*model.Transaction, error) {
	query := `
	SELECT id, idempotency_key, account_id, operation_type_id, amount, event_date, created_at
	FROM transactions WHERE idempotency_key = ? LIMIT 1
	`
	row := r.db.QueryRowContext(ctx, query, idempotencyKey)
	var transaction model.Transaction
	err := row.Scan(
		&transaction.ID,
		&transaction.IdempotencyKey,
		&transaction.AccountID,
		&transaction.OperationTypeID,
		&transaction.Amount,
		&transaction.EventDate,
		&transaction.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}
