package dto

import "time"

type NotificationMessage struct {
	Risk         bool      `json:"risk"`
	Message      string    `json:"message"`
	Date         time.Time `json:"date"`
	AmountPeople int       `json:"amount"`
}
