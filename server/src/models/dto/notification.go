package dto

import "time"

type NotificationJob struct {
	DateNotified    time.Time
	DateLastContact time.Time
	ForUser         string
	FromReport      int64
	Duration        time.Duration
	Attempts        int
}

type Notification struct {
	DateNotified time.Time
	ForUser      string
	ReportId     int64
}
