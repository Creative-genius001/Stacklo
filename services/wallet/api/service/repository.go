package service

import (
	"context"
	er "errors"
	"fmt"
	"os"

	"github.com/Creative-genius001/Stacklo/services/wallet/model"
	errors "github.com/Creative-genius001/Stacklo/services/wallet/utils/error"
	"github.com/Creative-genius001/Stacklo/services/wallet/utils/logger"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type Repository interface {
	GetFiatWallet(ctx context.Context, id string) (*model.Wallet, error)
	GetAllWallets(ctx context.Context, id string) ([]*model.Wallet, error)
	CreateFiatWallet(ctx context.Context, w model.Wallet) (*model.Wallet, error)
	CreateCryptoWallet(ctx context.Context, w model.Wallet) error
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

func (r *postgresRepository) GetFiatWallet(ctx context.Context, id string) (*model.Wallet, error) {

	currency := "NGN" //can only support NGN for now
	var w model.Wallet
	query := `SELECT
				w.*,
				f.virtual_account_name,
				f.virtual_account_number,
				f.virtual_bank_name
			FROM wallets w
			LEFT JOIN fiat_wallet_metadata f ON w.id = f.wallet_id
			WHERE w.user_id = $1
			AND w.currency = $2
			LIMIT 1;
		`
	row := r.db.QueryRow(ctx, query, id, currency)
	err := row.Scan(
		&w.ID,
		&w.UserId,
		&w.Currency,
		&w.Balance,
		&w.WalletType,
		&w.Active,
		&w.CreatedAt,
		&w.UpdatedAt,
		&w.VirtualAccountName,
		&w.VirtualAccountNumber,
		&w.VirtualBankName,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			logger.Logger.Error("Wallet not found", zap.Error(err), zap.String("wallet-id", id))
			return nil, err
		}
		logger.Logger.Error("Database query failed", zap.Error(err), zap.String("wallet-id", id))
		return nil, errors.Wrap(errors.TypeInternal, "failed to retrieve wallet", err)
	}
	return &w, nil
}

func (r *postgresRepository) GetAllWallets(ctx context.Context, id string) ([]*model.Wallet, error) {
	query := `SELECT
			w.*,
			f.virtual_account_name,
			f.virtual_account_number,
			f.virtual_bank_name
			FROM wallets w
			LEFT JOIN fiat_wallet_metadata f ON w.id = f.wallet_id 
			WHERE user_id = $1
			LIMIT 1;
		`
	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		logger.Logger.Error("Database query failed", zap.Error(err), zap.String("wallet-id", id))
		return nil, errors.Wrap(errors.TypeInternal, "database query error", err)
	}
	defer rows.Close()
	var wallets []*model.Wallet
	for rows.Next() {
		var w model.Wallet
		err := rows.Scan(
			&w.ID,
			&w.UserId,
			&w.Currency,
			&w.Balance,
			&w.WalletType,
			&w.Active,
			&w.CreatedAt,
			&w.UpdatedAt,
			&w.VirtualAccountName,
			&w.VirtualAccountNumber,
			&w.VirtualBankName,
		)
		if err != nil {
			logger.Logger.Error("Database query failed", zap.Error(err), zap.String("wallet-id", id))
			return nil, errors.Wrap(errors.TypeInternal, "failed to retireve wallet", err)
		}
		wallets = append(wallets, &w)
	}

	return wallets, nil
}

func (r *postgresRepository) CreateFiatWallet(ctx context.Context, w model.Wallet) (*model.Wallet, error) {

	var exists bool
	query := `SELECT EXISTS (
		SELECT 1 FROM wallets WHERE user_id=$1 AND currency=$2
	)`

	err := r.db.QueryRow(ctx, query, w.UserId, w.Currency).Scan(&exists)
	if err != nil {
		logger.Logger.Error("Server error. Retry again", zap.Error(err))
		return nil, errors.Wrap(errors.TypeInternal, "Server error. Retry again", err)
	}

	if exists {
		logger.Logger.Warn("Attempt to create duplicate fiat wallet", zap.String("user_id", w.UserId), zap.String("currency", w.Currency), zap.Error(err))
		msg := fmt.Sprintf("Wallet for user %s with currency %s already exists", w.UserId, w.Currency)
		return nil, errors.Wrap(errors.TypeConflict, "Cannot create duplicate wallet", er.New(msg))
	} else {
		var newWallet model.Wallet
		query1 := `
		INSERT INTO wallets (
			user_id, currency, balance, active, wallet_type, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, NOW(),NOW()) 
		RETURNING id
	`
		query2 := `
		INSERT INTO fiat_wallet_metadata (
			wallet_id, virtual_account_name, virtual_account_number, virtual_bank_name
		) VALUES (
			$1, $2, $3, $4
		)
	`
		tx, err := r.db.Begin(ctx)
		if err != nil {
			logger.Logger.Error("Error starting transaction", zap.Error(err))
			return nil, errors.Wrap(errors.TypeInternal, "Error starting transaction", err)
		}

		defer func() {
			if r := recover(); r != nil {
				tx.Rollback(ctx)
				logger.Logger.Error("Panic recovered during fiat wallet creation, transaction rolled back", zap.Any("panic_value", r))
				panic(r)
			} else if err != nil {
				tx.Rollback(ctx)
				logger.Logger.Error("Error occurred during fiat wallet creation, transaction rolled back", zap.Error(err))
			}
		}()

		err = tx.QueryRow(ctx,
			query1,
			w.UserId,
			w.Currency,
			w.Balance,
			w.Active,
			w.WalletType,
		).Scan(
			&newWallet.ID,
		)

		if err != nil {
			logger.Logger.Error("Error creating fiat wallet", zap.Error(err))
			return nil, errors.Wrap(errors.TypeInternal, "Error creating wallet", err)
		}

		_, err = tx.Exec(
			ctx,
			query2,
			newWallet.ID,
			w.VirtualAccountName,
			w.VirtualAccountNumber,
			w.VirtualBankName,
		)

		if err != nil {
			logger.Logger.Error("Error creating fiat-wallet_metadata", zap.Error(err))
			return nil, errors.Wrap(errors.TypeInternal, "Error creating wallet", err)
		}

		if err := tx.Commit(ctx); err != nil {
			logger.Logger.Error("commit transaction error", zap.Error(err))
			return nil, errors.Wrap(errors.TypeInternal, "commit transaction error", err)
		}

		newWallet = model.Wallet{
			ID:                   newWallet.ID,
			UserId:               w.UserId,
			Balance:              w.Balance,
			Active:               w.Active,
			WalletType:           w.WalletType,
			Currency:             w.Currency,
			VirtualAccountName:   w.VirtualAccountName,
			VirtualAccountNumber: w.VirtualAccountNumber,
			VirtualBankName:      w.VirtualBankName,
		}
		return &newWallet, nil
	}
}

func (r *postgresRepository) CreateCryptoWallet(ctx context.Context, w model.Wallet) error {
	query := `
		INSERT INTO wallets (
			id, user_id, currency, balance, active, wallet_type, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, NOW(), NOW()
		)
	`
	_, err := r.db.Exec(ctx,
		query,
		w.ID,
		w.UserId,
		w.Currency,
		w.Balance,
		w.Active,
		w.WalletType,
	)

	if err != nil {
		logger.Logger.Error("Error creating crypto wallet", zap.Error(err))
		return errors.Wrap(errors.TypeInternal, "Error creating wallet", err)
	}

	return nil
}

func (r *postgresRepository) Deposit(ctx context.Context, amount string) error {
	// TODO: actual logic
	return nil
}

func (r *postgresRepository) Withdraw(ctx context.Context, amount string) error {
	// TODO: actual logic
	return nil
}
