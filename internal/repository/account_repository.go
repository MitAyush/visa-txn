package repository

import (
	"context"
	"database/sql"

	"github.com/mitayush/visa-txn/internal/model"
)

type AccountRepository interface {
	Create(ctx context.Context, account *model.Account) (*model.Account, error)
	GetByAccountID(ctx context.Context, accountID int64) (*model.Account, error)
	ExistsByDocumentNumber(ctx context.Context, documentNumber string) (bool, error)
}

type accountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(ctx context.Context, account *model.Account) (*model.Account, error) {
	query := `
	INSERT INTO accounts (document_number) VALUES (?)
	RETURNING account_id, document_number, created_at
	`

	row := r.db.QueryRowContext(ctx, query, account.DocumentNumber)
	err := row.Scan(&account.AccountID, &account.DocumentNumber, &account.CreatedAt)
	return account, err
}

func (r *accountRepository) GetByAccountID(ctx context.Context, accountID int64) (*model.Account, error) {
	query := `
	SELECT account_id, document_number, created_at FROM accounts WHERE account_id = ? LIMIT 1
	`
	row := r.db.QueryRowContext(ctx, query, accountID)
	var account model.Account
	err := row.Scan(&account.AccountID, &account.DocumentNumber, &account.CreatedAt)
	return &account, err
}

func (r *accountRepository) ExistsByDocumentNumber(ctx context.Context, documentNumber string) (bool, error) {
	query := `
	SELECT COUNT(*) FROM accounts WHERE document_number = ? LIMIT 1
	`
	var count int
	err := r.db.QueryRowContext(ctx, query, documentNumber).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
