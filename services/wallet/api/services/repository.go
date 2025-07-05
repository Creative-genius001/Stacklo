package services

import (
	"context"
	"fmt"
	"os"

	"github.com/Creative-genius001/Stacklo/services/wallet/types"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	GetWallet(ctx context.Context, id string) (*types.Wallet, error)
	CreateWallet(ctx context.Context, w *types.Wallet) error
	Deposit(ctx context.Context, amount string) error
	Withdraw(ctx context.Context, amount string) error
}

type postgresRepository struct {
	db *pgx.Conn
}

func NewPostgresRepository(url string) (*postgresRepository, error) {
	db, err := pgx.Connect(context.Background(), url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return &postgresRepository{db}, nil
}

func (r *postgresRepository) Close() {
	r.db.Close(context.Background())
}

func (r *postgresRepository) GetWallet(ctx context.Context, id string) (*types.Wallet, error) {
	var w types.Wallet
	query := `SELECT * FROM wallets WHERE id= $1 and LIMIT 1`
	row := r.db.QueryRow(ctx, query, id)
	err := row.Scan(&w.ID, &w.UserId, &w.Active, &w.VirtualAccountName, &w.VirtualAccountNumber, &w.VirtualBankCode, &w.VirtualBankName, &w.Currency, &w.Balance, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &w, nil
}
