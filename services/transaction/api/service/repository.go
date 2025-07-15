package service

import (
	"context"
	er "errors"
	"fmt"
	"os"
	"strings"

	"github.com/Creative-genius001/Stacklo/services/transaction/model"
	errors "github.com/Creative-genius001/Stacklo/services/transaction/utils/error"
	"github.com/Creative-genius001/Stacklo/services/transaction/utils/logger"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type Repository interface {
	GetAllTransactions(ctx context.Context, id string) ([]*model.Transaction, error)
	CreateTransaction(ctx context.Context, w model.Transaction) error
	GetSingleTransaction(ctx context.Context, id string) (*model.Transaction, error)
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

func (r *postgresRepository) GetAllTransactions(ctx context.Context, userID string) ([]*model.Transaction, error) {
	query := `
		SELECT 
		t.*
        FROM transactions t 
		WHERE t.user_id = $1
		ORDER BY t.created_at DESC
		`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		logger.Logger.Error("Database query failed", zap.Error(err), zap.String("user-id", userID))
		return nil, errors.Wrap(errors.TypeInternal, "database query error", err)
	}
	defer rows.Close()
	var transactions []*model.Transaction
	for rows.Next() {
		var t model.Transaction
		err := rows.Scan(
			&t.ID,
			&t.UserId,
			&t.WalletId,
			&t.Currency,
			&t.Amount,
			&t.Reason,
			&t.EntryType,
			&t.Status,
			&t.TransactionType,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			logger.Logger.Error("Could not scan rows", zap.Error(err), zap.String("user-id", userID))
			return nil, errors.Wrap(errors.TypeInternal, "failed to retireve all transactions", err)
		}
		transactions = append(transactions, &t)
	}

	return transactions, nil
}

func (r *postgresRepository) CreateTransaction(ctx context.Context, t model.Transaction) error {

	query := `
		INSERT INTO transactions (
			user_id, wallet_id, currency, amount, status, reason, transaction_type, entry_type, created_at, updated_at
		)
		VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8,NOW(),NOW()
		)
		RETURNING id;
	`
	fQuery := `
		INSERT INTO fiat_transaction (
			id, reference_id, transaction_number, bank_name, account_name, account_number, fee, net_amount
		)
		VALUES (
			$1,$2,$3,$4,$5,$6,$7,$8
		)
	`
	cQuery := `
		INSERT INTO crypto_transaction (
			id, exchange_order_id, network, network_fee, price_at_transaction, quote_currency_amount
		)
		VALUES (
			$1,$2,$3,$4,$5,$6
		)
	`

	tx, err := r.db.Begin(ctx)
	if err != nil {
		logger.Logger.Error("Error starting db transaction", zap.Error(err))
		return errors.Wrap(errors.TypeInternal, "Error starting db transaction", err)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback(ctx)
			logger.Logger.Error("Panic recovered during fiat transaction creation, transaction rolled back", zap.Any("panic_value", r))
			panic(r)
		} else if err != nil {
			tx.Rollback(ctx)
			logger.Logger.Error("Error occurred during fiat transaction creation, transaction rolled back", zap.Error(err))
		}
	}()

	var transactionID string
	err = tx.QueryRow(
		ctx,
		query,
		t.UserId,
		t.WalletId,
		t.Currency,
		t.Amount,
		t.Status,
		t.Reason,
		t.TransactionType,
		t.EntryType,
	).Scan(&transactionID)
	if err != nil {
		logger.Logger.Error("Failed to insert transaction in table", zap.Error(err))
		return errors.Wrap(errors.TypeInternal, "Failed to insert transaction in table", err)
	}

	switch strings.ToUpper(t.TransactionType) {
	case "FIAT":
		if t.FiatDetails == nil {
			return errors.Wrap(errors.TypeInternal, "missing FiatDetails for fiat transaction", err)
		}
		_, err = tx.Exec(ctx, fQuery,
			transactionID,
			t.FiatDetails.ReferenceID,
			t.FiatDetails.TransactionNumber,
			t.FiatDetails.BankName,
			t.FiatDetails.AccountName,
			t.FiatDetails.AccountNumber,
			t.FiatDetails.Fee,
			t.FiatDetails.NetAmount,
		)
		if err != nil {
			return errors.Wrap(errors.TypeInternal, "failed to insert fiat details", err)
		}

	case "CRYPTO":
		if t.CryptoDetails == nil {
			return errors.Wrap(errors.TypeInternal, "missing crypto details for crpto transaction", err)
		}
		_, err = tx.Exec(ctx, cQuery,
			transactionID,
			t.CryptoDetails.ExchangeOrderID,
			t.CryptoDetails.Network,
			t.CryptoDetails.NetworkFee,
			t.CryptoDetails.PriceAtTransaction,
			t.CryptoDetails.QuoteCurrencyAmount,
		)
		if err != nil {
			return errors.Wrap(errors.TypeInternal, "failed to insert crypto details", err)
		}

	default:
		err := fmt.Sprintf("transaction type %s is invalid", t.TransactionType)
		return errors.Wrap(errors.TypeInternal, "invalid transaction type", er.New(err))
	}

	return tx.Commit(ctx)
}

func (r *postgresRepository) GetSingleTransaction(ctx context.Context, transactionID string) (*model.Transaction, error) {
	query := `
		SELECT 
		t.*,
		f.reference_id, 
		f.transaction_number, 
		f.bank_name, 
		f.account_name, 
		f.account_number, 
		f.fee, 
		f.net_amount,
		c.exchange_order_id,
		c.network,
		c.network_fee,
		c.price_at_transaction,
		c.quote_currency_amount
        FROM transactions t 
		LEFT JOIN fiat_transaction f ON f.id = t.id AND t.transaction_type = 'FIAT'
		LEFT JOIN crypto_transaction c ON c.id = t.id AND t.transaction_type = 'CRYPTO'
		WHERE t.id = $1
		LIMIT 1;
	`
	rows, err := r.db.Query(ctx, query, transactionID)
	if err != nil {
		logger.Logger.Error("Database query failed", zap.Error(err), zap.String("transaction-id", transactionID))
		return nil, errors.Wrap(errors.TypeInternal, "database query error", err)
	}
	defer rows.Close()
	var t model.Transaction
	var f model.FiatTransaction
	var c model.CryptoTransaction
	for rows.Next() {

		err := rows.Scan(
			&t.ID,
			&t.UserId,
			&t.WalletId,
			&t.Currency,
			&t.Amount,
			&t.Reason,
			&t.EntryType,
			&t.Status,
			&t.TransactionType,
			&t.CreatedAt,
			&t.UpdatedAt,
			&f.ReferenceID,
			&f.TransactionNumber,
			&f.BankName,
			&f.AccountName,
			&f.AccountNumber,
			&f.Fee,
			&f.NetAmount,
			&c.ExchangeOrderID,
			&c.Network,
			&c.NetworkFee,
			&c.PriceAtTransaction,
			&c.QuoteCurrencyAmount,
		)
		if err != nil {
			logger.Logger.Error("Could not scan rows", zap.Error(err), zap.String("transaction-id", transactionID))
			return nil, errors.Wrap(errors.TypeInternal, "failed to retireve transactions", err)
		}
		if f.ReferenceID != nil {
			t.FiatDetails = &f
		}

		if c.ExchangeOrderID != nil {
			t.CryptoDetails = &c
		}
	}

	return &t, nil
}
