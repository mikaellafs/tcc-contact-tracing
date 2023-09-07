package interfaces

import (
	"contacttracing/pkg/models/db"
	"context"
)

type ReportRepository interface {
	Migrate(ctx context.Context) error
	Create(ctx context.Context, report db.Report) (*db.Report, error)
	GetById(ctx context.Context, reportId int64) (*db.Report, error)
}
