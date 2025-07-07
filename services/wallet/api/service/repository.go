package service

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/Creative-genius001/Stacklo/services/wallet/types"
	"github.com/Creative-genius001/go-logger"
	"github.com/jackc/pgx/v5"
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

	logger.Info("db connection successful")
	return &postgresRepository{db}, nil
}

func (r *postgresRepository) Close() {
	r.db.Close(context.Background())
}

func (r *postgresRepository) GetWallet(ctx context.Context, id string) (*types.Wallet, error) {
	var w types.Wallet
	query := `SELECT * FROM wallets WHERE id = $1 LIMIT 1`
	row := r.db.QueryRow(ctx, query, id)
	err := row.Scan(&w.ID, &w.UserId, &w.Active, &w.VirtualAccountName, &w.VirtualAccountNumber, &w.VirtualBankName, &w.Currency, &w.Balance, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Error("Wallet not found due to error: ", err)
			return nil, errors.New("wallet not found")
		}
		logger.Error("Error retrieving wallet due to error: ", err)
		return nil, err
	}
	return &w, nil
}

func (r *postgresRepository) CreateWallet(ctx context.Context, w types.Wallet) (*types.Wallet, error) {

	query := `
		INSERT INTO wallets (
			id, user_id, currency, balance,
			virtual_account_name, virtual_account_number, virtual_bank_name, active,
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
		logger.Error("Failed to create wallet: ", err)
		return nil, errors.New("Failed to create wallet")
	}

	logger.Info("Successfully created wallet with ID: " + w.ID)
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
