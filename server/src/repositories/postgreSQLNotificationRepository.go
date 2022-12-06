package repositories

import (
	"context"
	"database/sql"
	"log"

	"contacttracing/src/models/db"
)

const (
	notificationRepositoryLog = "[Notification Repository]"
)

type PostgreSQLNotificationRepository struct {
	db *sql.DB
}

func NewPostgreSQLNotificationRepository(db *sql.DB) *PostgreSQLNotificationRepository {
	repo := new(PostgreSQLNotificationRepository)
	repo.db = db

	return repo
}

func (r *PostgreSQLNotificationRepository) Migrate(ctx context.Context) error {
	log.Println(notificationRepositoryLog, "Create notification table")
	query := `
    CREATE TABLE IF NOT EXISTS notifications(
		id SERIAL PRIMARY KEY,
		userId VARCHAR(36) NOT NULL,
		report INT NOT NULL,
        date TIMESTAMP WITH TIME ZONE NOT NULL,
		UNIQUE(userId,report),
		FOREIGN KEY (report) REFERENCES reports(id)
    );
    `
	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *PostgreSQLNotificationRepository) Create(ctx context.Context, notification db.Notification) (*db.Notification, error) {
	log.Println(notificationRepositoryLog, "Create new notification: ", notification)

	var id int64
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO notifications(userId, report, date) VALUES($1, $2, $3)
		RETURNING id`,
		notification.ForUser, notification.FromReport, notification.Date).Scan(&id)

	if err != nil {
		return nil, parsePostgreSQLError(err)
	}

	notification.ID = id
	return &notification, err
}
