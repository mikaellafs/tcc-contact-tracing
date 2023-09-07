package interfaces

import (
	"contacttracing/pkg/models/db"
	"context"
)

type NotificationRepository interface {
	Migrate(ctx context.Context) error
	Create(ctx context.Context, notification db.Notification) (*db.Notification, error)
}
