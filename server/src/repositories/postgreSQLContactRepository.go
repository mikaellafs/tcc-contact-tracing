package repositories

import (
	"contacttracing/src/models/db"
	"context"
	"database/sql"
	"time"
)

const (
	contactRepositoryLog = "Contact Repository: "
)

type PostgreSQLContactRepository struct {
	db *sql.DB
}

func NewPostgreSQLContactRepository(db *sql.DB) *PostgreSQLContactRepository {
	repo := new(PostgreSQLContactRepository)
	repo.db = db
	return repo
}

func (r *PostgreSQLContactRepository) Migrate(ctx context.Context) error {
	return nil
}

func (r *PostgreSQLContactRepository) Create(ctx context.Context, contact db.Contact) (*db.Contact, error) {
	return nil, nil
}

func (r *PostgreSQLContactRepository) GetContactsWithin(ctx context.Context, days int, from time.Time, userId string) ([]db.Contact, error) {
	return nil, nil
}
