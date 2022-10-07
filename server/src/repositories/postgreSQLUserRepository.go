package repositories

import (
	"contacttracing/src/models/db"
	"context"
	"database/sql"
	"errors"
)

type PostGreSQLUserRepository struct {
	db *sql.DB
}

func NewPostGreSQLUserRepository(db *sql.DB) *PostGreSQLUserRepository {
	repo := new(PostGreSQLUserRepository)
	repo.db = db
	return repo
}

func (r *PostGreSQLUserRepository) Migrate(ctx context.Context) error {
	query := `
    CREATE TABLE IF NOT EXISTS users(
        id SERIAL PRIMARY KEY,
        userId TEXT NOT NULL UNIQUE,
        pk TEXT NOT NULL,
        password TEXT NOT NULL
    );
    `

	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *PostGreSQLUserRepository) Create(ctx context.Context, user db.User) (*db.User, error) {
	var id int64
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO users(userId, pk, password) values($1, $2, $3) RETURNING id",
		user.UserId, user.Pk, user.Password).Scan(&id)

	user.ID = id
	return &user, err
}

func (r *PostGreSQLUserRepository) GetByUserId(ctx context.Context, userId string) (*db.User, error) {
	row := r.db.QueryRowContext(ctx, "SELECT * FROM users WHERE userId = $1", userId)

	var user db.User
	if err := row.Scan(&user.ID, &user.UserId, &user.Pk, &user.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExist
		}
		return nil, err
	}
	return &user, nil
}

func (r *PostGreSQLUserRepository) Update(ctx context.Context, user db.User) (*db.User, error) {
	return nil, nil
}

func (r *PostGreSQLUserRepository) Delete(ctx context.Context, id int64) error {
	return nil
}