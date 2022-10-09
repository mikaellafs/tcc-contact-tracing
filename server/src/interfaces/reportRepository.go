package interfaces

import (
	"contacttracing/src/models/db"
	"context"
)

type ReportRepository interface {
	Migrate(ctx context.Context) error
	Create(ctx context.Context, report db.Report) (*db.Report, error)
	GetByUserId(ctx context.Context, reportId string) ([]*db.Report, error)
}
