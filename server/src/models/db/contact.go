package db

import "time"

type Contact struct {
	ID                    int64
	User                  string
	AnotherUser           string
	FirstContactTimestamp time.Time
	LastContactTimestamp  time.Time
	Distance              float32
	RSSI                  float32
	BatteryLevel          float32
}
