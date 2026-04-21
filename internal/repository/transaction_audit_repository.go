package repository

import (
	"context"
	"database/sql"

	"github.com/mitayush/visa-txn/internal/model"
)

type TransactionAuditRepository interface {
	InsertTransactionAuditWithConn(ctx context.Context, conn DBTX, audit *model.TransactionAudit) error
}

type transactionAuditRepository struct {
	db *sql.DB
}

func NewTransactionAuditRepository(db *sql.DB) TransactionAuditRepository {
	return &transactionAuditRepository{db: db}
}

func (r *transactionAuditRepository) InsertTransactionAuditWithConn(ctx context.Context, conn DBTX, audit *model.TransactionAudit) error {
	query := `
	INSERT INTO transaction_audit (idempotency_key, request, response, status, created_at) VALUES (?, ?, ?, ?, ?)
	`
	_, err := conn.ExecContext(ctx, query, audit.IdempotencyKey, audit.Request, audit.Response, audit.Status, audit.CreatedAt)
	return err
}
