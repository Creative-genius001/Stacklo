package service

import (
	"context"
	er "errors"
	"fmt"
	"os"
	"strings"

	"github.com/Creative-genius001/Stacklo/services/user/model"
	"github.com/Creative-genius001/Stacklo/services/user/types"
	errors "github.com/Creative-genius001/Stacklo/services/user/utils/error"
	"github.com/Creative-genius001/Stacklo/services/user/utils/logger"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type Repository interface {
	FindByID(ctx context.Context, id string) (*model.User, error)
	CreateUser(ctx context.Context, user model.User) (*model.User, error)
	UpdateUser(ctx context.Context, id string, data types.UpdateUser) error
	UpdateEmail(ctx context.Context, id string, email string) error
	UpdatePhone(ctx context.Context, id string, phone string) error
	FindByPhoneOrEmail(ctx context.Context, email string, phone string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateVerificationStatus(ctx context.Context, email string, status bool) error
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

func (r *postgresRepository) FindByPhoneOrEmail(ctx context.Context, email string, phone string) (*model.User, error) {
	query := `
		SELECT 
		u.*
		FROM users u
		WHERE u.email = $1 OR u.phone_number = $2
		LIMIT 1;
	`
	var user model.User
	err := r.db.QueryRow(ctx, query, email, phone).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
		&user.Country,
		&user.IsVerified,
		&user.KycStatus,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if er.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		logger.Logger.Warn("Error getting user", zap.Error(err))
		return nil, errors.Wrap(errors.TypeInternal, "Scan error", err)
	}

	return &user, nil
}

func (r *postgresRepository) CreateUser(ctx context.Context, user model.User) (*model.User, error) {
	query := `
		INSERT into users (email, password_hash, first_name, last_name, phone_number, country, isVerified, kyc_status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		RETURNING id, email, first_name, last_name;
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
		user.IsVerified,
		user.KycStatus,
	).Scan(
		&newUser.ID,
		&newUser.Email,
		&newUser.FirstName,
		&newUser.LastName,
	)

	if err != nil {
		logger.Logger.Error("Error creating user in databse", zap.Error(err))
		return nil, errors.Wrap(errors.TypeInternal, "Scan error:", err)
	}

	return &newUser, nil
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
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
		&user.Country,
		&user.IsVerified,
		&user.KycStatus,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if er.Is(err, pgx.ErrNoRows) {
		logger.Logger.Warn("user not found")
		return nil, nil
	}
	if err != nil {
		logger.Logger.Warn("Error getting user", zap.Error(err))
		return nil, errors.Wrap(errors.TypeInternal, "Error getting user", err)
	}

	return &user, nil
}

func (r *postgresRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	query := `
		SELECT
		u.*
		FROM users u
		WHERE u.id = $1
		LIMIT 1;
	`
	var user model.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.PhoneNumber,
		&user.Country,
		&user.IsVerified,
		&user.KycStatus,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if er.Is(err, pgx.ErrNoRows) {
		logger.Logger.Warn("user not found")
		return nil, nil
	}
	if err != nil {
		logger.Logger.Warn("Error getting user", zap.Error(err))
		return nil, errors.Wrap(errors.TypeInternal, "Error getting user", err)
	}

	return &user, nil
}

func (r *postgresRepository) UpdateUser(ctx context.Context, id string, data types.UpdateUser) error {
	columns := []string{}
	values := []interface{}{}
	argPos := 1

	if data.FirstName != nil {
		columns = append(columns, fmt.Sprintf("first_name = $%d", argPos))
		values = append(values, *data.FirstName)
		argPos++
	}

	if data.LastName != nil {
		columns = append(columns, fmt.Sprintf("last_name = $%d", argPos))
		values = append(values, *data.LastName)
		argPos++
	}

	if data.Country != nil {
		columns = append(columns, fmt.Sprintf("country = $%d", argPos))
		values = append(values, *data.Country)
		argPos++
	}
	if data.KycStatus != nil {
		columns = append(columns, fmt.Sprintf("kyc_status = $%d", argPos))
		values = append(values, *data.KycStatus)
		argPos++
	}

	if len(columns) == 0 {
		return nil
	}

	query := fmt.Sprintf("UPDATE users SET %s, updated_at = NOW() WHERE id = $%d", strings.Join(columns, ", "), argPos)
	values = append(values, id)

	cmdTag, err := r.db.Exec(ctx, query, values...)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		logger.Logger.Error("No user found", zap.String("user-ID", id))
		return errors.Wrap(errors.TypeNotFound, "no user found ", er.New("user not found"))
	}
	return nil
}

func (r *postgresRepository) UpdateEmail(ctx context.Context, id string, email string) error {
	query := `
		UPDATE users
		SET email = $1, updated_at = NOW()
		WHERE id = $2;	
	`

	cmdTag, err := r.db.Exec(ctx, query, email, id)
	if err != nil {
		logger.Logger.Error("failed to update user email", zap.String("user-ID", id), zap.Error(err))
		return errors.Wrap(errors.TypeInternal, "failed to update user email ", err)
	}
	if cmdTag.RowsAffected() == 0 {
		logger.Logger.Error("No user found", zap.String("user-ID", id))
		return errors.Wrap(errors.TypeNotFound, "no user found ", er.New("user not found"))
	}
	return nil
}

func (r *postgresRepository) UpdatePhone(ctx context.Context, id string, phone string) error {
	query := `
		UPDATE users
		SET phone_number = $1, updated_at = NOW()
		WHERE id = $2;	
	`

	cmdTag, err := r.db.Exec(ctx, query, phone, id)
	if err != nil {
		logger.Logger.Error("failed to update user phone number", zap.String("user-ID", id), zap.Error(err))
		return errors.Wrap(errors.TypeInternal, "failed to update user phone number ", err)
	}
	if cmdTag.RowsAffected() == 0 {
		logger.Logger.Error("No user found", zap.String("user-ID", id))
		return errors.Wrap(errors.TypeNotFound, "no user found ", er.New("user not found"))
	}
	return nil
}

func (r *postgresRepository) UpdateVerificationStatus(ctx context.Context, email string, status bool) error {
	query := `
		UPDATE users
		SET isVerified = $1, updated_at = NOW()
		WHERE email = $2;	
	`

	cmdTag, err := r.db.Exec(ctx, query, status, email)
	if err != nil {
		logger.Logger.Error("failed to update user verified status", zap.String("email", email), zap.Error(err))
		return errors.Wrap(errors.TypeInternal, "failed to update user verified status ", err)
	}
	if cmdTag.RowsAffected() == 0 {
		logger.Logger.Error("No user found", zap.String("user-ID", email))
		return errors.Wrap(errors.TypeNotFound, "no user found ", er.New("user not found"))
	}
	return nil
}

func (r *postgresRepository) Close() {
	r.db.Close(context.Background())
}
