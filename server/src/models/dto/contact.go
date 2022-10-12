package dto

import "time"

type Contact struct {
	User        string
	AnotherUser string
	Duration    time.Duration
}
