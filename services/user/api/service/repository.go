package service

import (
	"context"
	"fmt"
	"os"

	"github.com/Creative-genius001/Stacklo/services/user/model"
	"github.com/Creative-genius001/Stacklo/services/user/utils/logger"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	GetUser(ctx context.Context, id string) (*model.User, error)
	CreateUser(ctx context.Context, w model.User) (*model.User, error)
	UpdateUser(ctx context.Context, w model.User) error
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

// CreateUser implements Repository.
func (r *postgresRepository) CreateUser(ctx context.Context, w model.User) (*model.User, error) {
	panic("unimplemented")
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
