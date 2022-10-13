package interfaces

import (
	"contacttracing/src/models/dto"
	"time"
)

type CacheRepository interface {
	SaveReport(userId, date string)
	GetReportsFrom(userId string) []time.Time

	SaveNotification(userId string, fromReport int64, date string)
	GetNotificationFrom(userId string, reportId int64) *dto.Notification
}
