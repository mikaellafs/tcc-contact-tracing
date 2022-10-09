package dto

import "time"

type ReportJob struct {
	UserId   string
	Date     time.Time
	Attempts int
}
