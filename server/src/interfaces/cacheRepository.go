package interfaces

import (
	"time"
)

type CacheRepository interface {
	SaveReport(userId, date string)
	GetReportsFrom(userId string) []time.Time
}
