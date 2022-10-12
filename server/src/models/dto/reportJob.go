package dto

import "time"

type ReportJob struct {
	DBID     int64
	UserId   string
	Date     time.Time
	Attempts int
}
