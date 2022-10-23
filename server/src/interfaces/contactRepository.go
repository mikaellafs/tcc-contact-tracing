package interfaces

import (
	"context"
	"time"

	"contacttracing/src/models/db"
	"contacttracing/src/models/dto"
)

type ContactRepository interface {
	Migrate(ctx context.Context) error
	Create(ctx context.Context, contact db.Contact) (*db.Contact, error)
	GetContactsWithin(ctx context.Context, days int, from time.Time, userId string) ([]dto.Contact, error)
	GetContactsBetweenUsers(ctx context.Context, user1, user2 string, from, to time.Time) ([]dto.Contact, error)
}
