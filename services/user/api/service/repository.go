package service

import (
	"context"
	"fmt"
	"os"

	"github.com/Creative-genius001/Stacklo/services/user/model"
	errors "github.com/Creative-genius001/Stacklo/services/user/utils/error"
	"github.com/Creative-genius001/Stacklo/services/user/utils/logger"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type Repository interface {
	GetUser(ctx context.Context, id string) (*model.User, error)
	CreateUser(ctx context.Context, user model.User) (*model.User, error)
	UpdateUser(ctx context.Context, user model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
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

func (r *postgresRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT 
		u.*
		FROM users u
		WHERE u.email = $1
		LIMIT 1;
	`
	var user model.User
	err := r.db.QueryRow(ctx, query, email).Scan(&email)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		logger.Logger.Warn("Error getting user", zap.Error(err))
		return nil, errors.Wrap(errors.TypeInternal, "Error getting user", err)
	}

	return &user, nil
}

func (r *postgresRepository) CreateUser(ctx context.Context, user model.User) (*model.User, error) {
	query := `
		INSERT into users (email, password_hash, first_name, last_name, phone_number, country, kyc_status, created_at, update_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
	`
	var newUser model.User
	err := r.db.QueryRow(
		ctx,
		query,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.PhoneNumber,
		user.Country,
		user.KycStatus,
	).Scan(
		&newUser.ID,
		&newUser.Email,
		&newUser.PasswordHash,
		&newUser.FirstName,
		&newUser.LastName,
		&newUser.PhoneNumber,
		&newUser.Country,
		&newUser.KycStatus,
		&newUser.CreatedAt,
		&newUser.UpdatedAt,
	)

	if err != nil {
		logger.Logger.Error("Error creating user in databse", zap.Error(err))
		return nil, errors.Wrap(errors.TypeInternal, "Error creating user in databse", err)
	}

	return &newUser, nil
}

// GetUser implements Repository.
func (r *postgresRepository) GetUser(ctx context.Context, id string) (*model.User, error) {
	panic("unimplemented")
}

// UpdateUser implements Repository.
func (r *postgresRepository) UpdateUser(ctx context.Context, w model.User) error {
	panic("unimplemented")
}

func (r *postgresRepository) Close() {
	r.db.Close(context.Background())
}
