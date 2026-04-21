package service_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/mitayush/visa-txn/internal/apperrors"
	"github.com/mitayush/visa-txn/internal/config"
	"github.com/mitayush/visa-txn/internal/model"
	"github.com/mitayush/visa-txn/internal/repository"
	repomocks "github.com/mitayush/visa-txn/internal/repository/mocks"
	"github.com/mitayush/visa-txn/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func transactionTestConfig() *config.Config {
	return &config.Config{
		OperationTypes: map[int]config.OperationType{
			1: {Description: "Normal Purchase", Sign: config.SignDebit},
			4: {Description: "Credit Voucher", Sign: config.SignCredit},
		},
	}
}

func openTxTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db := repository.NewDB(":memory:")
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func TestCreateTransaction(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 4, 18, 12, 0, 0, 0, time.UTC)
	cfg := transactionTestConfig()

	type setupFn func(
		t *testing.T,
		tx *repomocks.MockTransactionRepository,
		acc *repomocks.MockAccountRepository,
		aud *repomocks.MockTransactionAuditRepository,
	)

	tests := []struct {
		name    string
		in      *model.Transaction
		needsDB bool
		setup   setupFn
		want    *model.Transaction
		wantErr error
	}{
		{
			name: "rejects_duplicate_idempotency_key",
			in: &model.Transaction{
				IdempotencyKey: "dup-key", AccountID: 1, OperationTypeID: 1,
				Amount: 10, EventDate: now,
			},
			setup: func(_ *testing.T, tx *repomocks.MockTransactionRepository, _ *repomocks.MockAccountRepository, _ *repomocks.MockTransactionAuditRepository) {
				tx.On("GetByIdempotencyKey", ctx, "dup-key").Return(&model.Transaction{ID: 99}, nil)
			},
			wantErr: apperrors.ErrTransactionAlreadyExists,
		},
		{
			name: "rejects_missing_account",
			in: &model.Transaction{
				IdempotencyKey: "k1", AccountID: 404, OperationTypeID: 1,
				Amount: 10, EventDate: now,
			},
			setup: func(_ *testing.T, tx *repomocks.MockTransactionRepository, acc *repomocks.MockAccountRepository, _ *repomocks.MockTransactionAuditRepository) {
				tx.On("GetByIdempotencyKey", ctx, "k1").Return(nil, nil)
				acc.On("GetByAccountID", ctx, int64(404)).Return((*model.Account)(nil), sql.ErrNoRows)
			},
			wantErr: apperrors.ErrAccountDoesNotExist,
		},
		{
			name: "rejects_unknown_operation_type",
			in: &model.Transaction{
				IdempotencyKey: "k2", AccountID: 1, OperationTypeID: 999,
				Amount: 10, EventDate: now,
			},
			setup:   nil,
			wantErr: apperrors.ErrInvalidOperationType,
		},
		{
			name:    "persists_debit_with_negative_amount",
			needsDB: true,
			in: &model.Transaction{
				IdempotencyKey: "k-debit", AccountID: 1, OperationTypeID: 1,
				Amount: 50, EventDate: now,
			},
			setup: func(_ *testing.T, tx *repomocks.MockTransactionRepository, acc *repomocks.MockAccountRepository, aud *repomocks.MockTransactionAuditRepository) {
				tx.On("GetByIdempotencyKey", ctx, "k-debit").Return(nil, nil)
				acc.On("GetByAccountID", ctx, int64(1)).Return(&model.Account{AccountID: 1}, nil)
				tx.On("InsertTransactionWithConn", mock.Anything, mock.Anything, mock.MatchedBy(func(tr *model.Transaction) bool {
					return tr.IdempotencyKey == "k-debit" && tr.AccountID == 1 && tr.OperationTypeID == 1 && tr.Amount == -50 && tr.EventDate.Equal(now)
				})).Return(&model.Transaction{
					ID: 100, IdempotencyKey: "k-debit", AccountID: 1, OperationTypeID: 1,
					Amount: -50, EventDate: now, CreatedAt: now,
				}, nil)
				aud.On("InsertTransactionAuditWithConn", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			want: &model.Transaction{
				ID: 100, IdempotencyKey: "k-debit", AccountID: 1, OperationTypeID: 1,
				Amount: -50, EventDate: now, CreatedAt: now,
			},
		},
		{
			name:    "persists_credit_with_positive_amount",
			needsDB: true,
			in: &model.Transaction{
				IdempotencyKey: "k-credit", AccountID: 2, OperationTypeID: 4,
				Amount: 25, EventDate: now,
			},
			setup: func(_ *testing.T, tx *repomocks.MockTransactionRepository, acc *repomocks.MockAccountRepository, aud *repomocks.MockTransactionAuditRepository) {
				tx.On("GetByIdempotencyKey", ctx, "k-credit").Return(nil, nil)
				acc.On("GetByAccountID", ctx, int64(2)).Return(&model.Account{AccountID: 2}, nil)
				tx.On("InsertTransactionWithConn", mock.Anything, mock.Anything, mock.MatchedBy(func(tr *model.Transaction) bool {
					return tr.IdempotencyKey == "k-credit" && tr.AccountID == 2 && tr.OperationTypeID == 4 && tr.Amount == 25 && tr.EventDate.Equal(now)
				})).Return(&model.Transaction{
					ID: 200, IdempotencyKey: "k-credit", AccountID: 2, OperationTypeID: 4,
					Amount: 25, EventDate: now, CreatedAt: now,
				}, nil)
				aud.On("InsertTransactionAuditWithConn", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			want: &model.Transaction{
				ID: 200, IdempotencyKey: "k-credit", AccountID: 2, OperationTypeID: 4,
				Amount: 25, EventDate: now, CreatedAt: now,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var db *sql.DB
			if tt.needsDB {
				db = openTxTestDB(t)
			}
			txRepo := repomocks.NewMockTransactionRepository(t)
			accRepo := repomocks.NewMockAccountRepository(t)
			audRepo := repomocks.NewMockTransactionAuditRepository(t)
			if tt.setup != nil {
				tt.setup(t, txRepo, accRepo, audRepo)
			}
			svc := service.NewTransactionService(db, txRepo, accRepo, audRepo, cfg)
			inCopy := *tt.in
			got, err := svc.CreateTransaction(ctx, &inCopy)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, got)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
