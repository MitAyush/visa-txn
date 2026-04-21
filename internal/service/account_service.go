package service

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/mitayush/visa-txn/internal/apperrors"
	"github.com/mitayush/visa-txn/internal/model"
	"github.com/mitayush/visa-txn/internal/repository"
)

type AccountService struct {
	repo repository.AccountRepository
}

func NewAccountService(repo repository.AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

func (s *AccountService) CreateAccount(ctx context.Context, account *model.Account) (*model.Account, error) {
	// TODO: its not safe in distributed system, we need distributed lock to avoid race conditions
	exists, err := s.repo.ExistsByDocumentNumber(ctx, account.DocumentNumber)
	if err != nil {
		log.Println("error checking if account exists", err)
		return nil, apperrors.ErrInternalServerError
	}
	if exists {
		return nil, apperrors.ErrAccountExists
	}
	account, err = s.repo.Create(ctx, account)
	if err != nil {
		log.Println("error creating account", err)
		return nil, apperrors.ErrInternalServerError
	}
	return account, nil
}

func (s *AccountService) GetAccount(ctx context.Context, accountID int64) (*model.Account, error) {
	if accountID <= 0 {
		return nil, apperrors.ErrInvalidAccountID
	}
	account, err := s.repo.GetByAccountID(ctx, accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrAccountNotFound
		}
		log.Println("error getting account", err)
		return nil, apperrors.ErrInternalServerError
	}
	return account, nil
}
