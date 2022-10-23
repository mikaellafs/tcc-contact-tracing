package db

import "time"

type Report struct {
	ID             int64
	UserId         string
	DateStart      time.Time
	DateDiagnostic time.Time
	DateReport     time.Time
}
