package service

import (
	"context"
	"fmt"
	"os"

	"github.com/Creative-genius001/Stacklo/services/wallet/types"
	errors "github.com/Creative-genius001/Stacklo/services/wallet/utils/error"
	"github.com/Creative-genius001/Stacklo/services/wallet/utils/logger"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type Repository interface {
	GetWallet(ctx context.Context, id string) (*types.Wallet, error)
	CreateWallet(ctx context.Context, w types.Wallet) (*types.Wallet, error)
	Deposit(ctx context.Context, amount string) error
	Withdraw(ctx context.Context, amount string) error
	Close()
}

type postgresRepository struct {
	db *pgx.Conn
}

func NewPostgresRepository(url string) (Repository, error) {
	db, err := pgx.Connect(context.Background(), url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	logger.Logger.Info("DB connection successful")
	return &postgresRepository{db}, nil
}

func (r *postgresRepository) Close() {
	r.db.Close(context.Background())
}

func (r *postgresRepository) GetWallet(ctx context.Context, id string) (*types.Wallet, error) {
	var w types.Wallet
	query := `SELECT id, user_id, currency, balance,virtual_account_name, virtual_account_number,  virtual_bank_name,active,  created_at, updated_at
          FROM wallets WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	err := row.Scan(&w.ID, &w.UserId, &w.Currency, &w.Balance, &w.VirtualAccountName, &w.VirtualAccountNumber, &w.VirtualBankName, &w.Active, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, err
		}
		logger.Logger.Error("Database query failed", zap.Error(err), zap.String("wallet-id", id))
		return nil, errors.Wrap(errors.TypeInternal, "database query error", err)
	}
	return &w, nil
}

func (r *postgresRepository) CreateWallet(ctx context.Context, w types.Wallet) (*types.Wallet, error) {

	query := `
		INSERT INTO wallets (
			id, user_id, currency, balance,
			virtual_account_name, virtual_account_number, virtual_bank_name, 
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		) RETURNING id, user_id, currency, balance,
					virtual_account_name, virtual_account_number, virtual_bank_name, active,
					created_at, updated_at
	`
	err := r.db.QueryRow(
		ctx,
		query,
		w.ID,
		w.UserId,
		w.Currency,
		w.Balance,
		w.VirtualAccountName,
		w.VirtualAccountNumber,
		w.VirtualBankName,
		w.Active,
		w.CreatedAt,
		w.UpdatedAt,
	).Scan(
		&w.ID,
		&w.UserId,
		&w.Currency,
		&w.Balance,
		&w.VirtualAccountName,
		&w.VirtualAccountNumber,
		&w.VirtualBankName,
		&w.Active,
		&w.CreatedAt,
		&w.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *postgresRepository) Deposit(ctx context.Context, amount string) error {
	// TODO: actual logic
	return nil
}

func (r *postgresRepository) Withdraw(ctx context.Context, amount string) error {
	// TODO: actual logic
	return nil
}
