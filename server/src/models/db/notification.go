package db

import "time"

type Notification struct {
	ID         int64
	ForUser    string
	FromReport int64
	Date       time.Time
}
