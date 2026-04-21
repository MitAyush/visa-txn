package service_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/mitayush/visa-txn/internal/apperrors"
	"github.com/mitayush/visa-txn/internal/model"
	repomocks "github.com/mitayush/visa-txn/internal/repository/mocks"
	"github.com/mitayush/visa-txn/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 4, 18, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		in      *model.Account
		setup   func(*repomocks.MockAccountRepository)
		want    *model.Account
		wantErr error
	}{
		{
			name: "rejects_duplicate_document",
			in:   &model.Account{DocumentNumber: "dup"},
			setup: func(m *repomocks.MockAccountRepository) {
				m.On("ExistsByDocumentNumber", ctx, "dup").Return(true, nil)
			},
			wantErr: apperrors.ErrAccountExists,
		},
		{
			name: "persists_when_document_free",
			in:   &model.Account{DocumentNumber: "new"},
			setup: func(m *repomocks.MockAccountRepository) {
				m.On("ExistsByDocumentNumber", ctx, "new").Return(false, nil)
				created := &model.Account{AccountID: 7, DocumentNumber: "new", CreatedAt: now}
				m.On("Create", ctx, mock.AnythingOfType("*model.Account")).Return(created, nil)
			},
			want: &model.Account{AccountID: 7, DocumentNumber: "new", CreatedAt: now},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := repomocks.NewMockAccountRepository(t)
			if tt.setup != nil {
				tt.setup(m)
			}
			svc := service.NewAccountService(m)
			got, err := svc.CreateAccount(ctx, tt.in)
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

func TestGetAccount(t *testing.T) {
	ctx := context.Background()
	now := time.Date(2026, 4, 18, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		id      int64
		setup   func(*repomocks.MockAccountRepository)
		want    *model.Account
		wantErr error
	}{
		{
			name:    "rejects_non_positive_id",
			id:      0,
			wantErr: apperrors.ErrInvalidAccountID,
		},
		{
			name: "maps_missing_row_to_not_found",
			id:   99,
			setup: func(m *repomocks.MockAccountRepository) {
				m.On("GetByAccountID", ctx, int64(99)).Return((*model.Account)(nil), sql.ErrNoRows)
			},
			wantErr: apperrors.ErrAccountNotFound,
		},
		{
			name: "returns_account",
			id:   5,
			setup: func(m *repomocks.MockAccountRepository) {
				acc := &model.Account{AccountID: 5, DocumentNumber: "ok", CreatedAt: now}
				m.On("GetByAccountID", ctx, int64(5)).Return(acc, nil)
			},
			want: &model.Account{AccountID: 5, DocumentNumber: "ok", CreatedAt: now},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := repomocks.NewMockAccountRepository(t)
			if tt.setup != nil {
				tt.setup(m)
			}
			svc := service.NewAccountService(m)
			got, err := svc.GetAccount(ctx, tt.id)
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
