package interfaces

import (
	"contacttracing/src/models/dto"
	"time"
)

type CacheRepository interface {
	SaveReport(userId string, reportId int64, date time.Time)
	GetReportsFrom(userId string) []dto.Report
	UserHasReportedRecently(userId string) bool
	SaveUserHasReportedRecently(userId string, expiration time.Duration)

	SaveNotification(userId string, fromReport int64, date time.Time)
	GetNotificationFrom(userId string, reportId int64) *dto.Notification
	GetNotificationKeysFrom(userId string) []string
	RemoveNotificationsAfter(days time.Duration, maxDate time.Time) []string

	SavePotentialRiskJob(userId string, reportId int64)
	GetPotentialRiskJob(userId string, reportId int64) bool
	RemovePotentialRiskJob(userId string, reportId int64)
}
