package service

import (
	"context"

	"github.com/Creative-genius001/Stacklo/services/transaction/model"
)

type Service interface {
	GetAllTransactions(ctx context.Context, id string) ([]*model.Transaction, error)
	CreateTransaction(ctx context.Context, w model.Transaction) error
	GetSingleTransaction(ctx context.Context, id string) (*model.Transaction, error)
	GetFilteredTransactions(ctx context.Context, f model.TransactionFilter) ([]model.Transaction, *string, error)
}

type transactionService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &transactionService{r}
}

func (t *transactionService) CreateTransaction(ctx context.Context, tr model.Transaction) error {
	err := t.repository.CreateTransaction(ctx, tr)
	if err != nil {
		return err
	}

	return nil
}

func (t *transactionService) GetAllTransactions(ctx context.Context, userID string) ([]*model.Transaction, error) {
	res, err := t.repository.GetAllTransactions(ctx, userID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *transactionService) GetSingleTransaction(ctx context.Context, transactionID string) (*model.Transaction, error) {
	res, err := t.repository.GetSingleTransaction(ctx, transactionID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t *transactionService) GetFilteredTransactions(ctx context.Context, f model.TransactionFilter) ([]model.Transaction, *string, error) {
	res, cursor, err := t.repository.GetFilteredTransactions(ctx, f)
	if err != nil {
		return nil, nil, err
	}

	return res, cursor, nil
}
