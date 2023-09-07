package repositories

import (
	"contacttracing/pkg/models/db"
	"context"
	"database/sql"
	"errors"
	"log"
)

const (
	userRepositoryLog = "[User Repository]"
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
	log.Println(userRepositoryLog, "Create users table")

	query := `
    CREATE TABLE IF NOT EXISTS users(
        deviceId VARCHAR(100) PRIMARY KEY,
		userId VARCHAR(36) UNIQUE,
        pk TEXT NOT NULL
    );
    `

	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *PostGreSQLUserRepository) Create(ctx context.Context, user db.User) (*db.User, error) {
	log.Println(userRepositoryLog, "Create user: ", user)

	var id string
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO users(deviceId, userId, pk) VALUES($1, $2, $3)
		ON CONFLICT(deviceId)
		DO UPDATE SET
			pk = $3
		RETURNING userId`,
		user.DeviceId, user.Id, user.Pk).Scan(&id)

	if err != nil {
		err = parsePostgreSQLError(err)
	}

	user.Id = id
	return &user, err
}

func (r *PostGreSQLUserRepository) GetByUserId(ctx context.Context, userId string) (*db.User, error) {
	log.Println(userRepositoryLog, "Get user by id = ", userId)

	row := r.db.QueryRowContext(ctx, "SELECT * FROM users WHERE userId = $1", userId)

	var user db.User
	if err := row.Scan(&user.DeviceId, &user.Id, &user.Pk); err != nil {
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
