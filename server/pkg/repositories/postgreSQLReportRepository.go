package repositories

import (
	"contacttracing/pkg/models/db"
	"context"
	"database/sql"
	"log"
)

const (
	reportRepositoryLog = "[Report Repository]"
)

type PostgreSQLReportRepository struct {
	db *sql.DB
}

func NewPostGreSQLReportRepository(db *sql.DB) *PostgreSQLReportRepository {
	repo := new(PostgreSQLReportRepository)
	repo.db = db

	return repo
}

func (r *PostgreSQLReportRepository) Migrate(ctx context.Context) error {
	log.Println(reportRepositoryLog, "Create reports table")

	query := `
    CREATE TABLE IF NOT EXISTS reports(
		id SERIAL PRIMARY KEY,
        userId VARCHAR(36) NOT NULL,
        dateStart TIMESTAMP WITH TIME ZONE NOT NULL,
        dateDiagnostic TIMESTAMP WITH TIME ZONE NOT NULL,
		dateReport TIMESTAMP WITH TIME ZONE NOT NULL,
		UNIQUE(userId, dateDiagnostic)
    );
    `

	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *PostgreSQLReportRepository) Create(ctx context.Context, report db.Report) (*db.Report, error) {
	log.Println(reportRepositoryLog, "Create new report: ", report)

	var id int64
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO reports(userId, dateStart, dateDiagnostic, dateReport) VALUES($1, $2, $3, $4)
		RETURNING id`,
		report.UserId, report.DateStart, report.DateDiagnostic, report.DateReport).Scan(&id)

	if err != nil {
		return nil, parsePostgreSQLError(err)
	}

	report.ID = id
	return &report, err
}

func (r *PostgreSQLReportRepository) GetById(ctx context.Context, reportId int64) (*db.Report, error) {
	query := `SELECT * FROM reports WHERE id = $1`

	report := db.Report{}
	err := r.db.QueryRowContext(ctx, query, reportId).Scan(
		&report.ID, &report.UserId, &report.DateStart, &report.DateDiagnostic, &report.DateReport)

	if err != nil {
		return nil, err
	}

	return &report, nil
}
