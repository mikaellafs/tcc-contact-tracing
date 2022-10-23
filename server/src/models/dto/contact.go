package dto

import "time"

type Contact struct {
	User            string
	DateLastContact time.Time
	AnotherUser     string
	Duration        time.Duration
}
