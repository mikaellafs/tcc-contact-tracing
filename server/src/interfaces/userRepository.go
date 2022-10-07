package interfaces

import (
	"contacttracing/src/models/db"
	"context"
)

type UserRepository interface {
	Migrate(ctx context.Context) error
	Create(ctx context.Context, user db.User) (*db.User, error)
	GetByUserId(ctx context.Context, userId string) (*db.User, error)
	Update(ctx context.Context, user db.User) (*db.User, error)
	Delete(ctx context.Context, id int64) error
}
