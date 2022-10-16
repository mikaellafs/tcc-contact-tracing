package dto

import (
	"encoding/json"
	"log"
)

type ContactMessage struct {
	User      string `json:"user"`
	Contact   []byte `json:"contact"`
	Signature []byte `json:"signature"`
}

type ContactFromMessage struct {
	User                  string  `json:"token"`
	FirstContactTimestamp int64   `json:"firstContactTimestamp"` // milliseconds
	LastContactTimestamp  int64   `json:"lastContactTimestamp"`  // milliseconds
	Distance              float32 `json:"distance"`
	RSSI                  float32 `json:"rssi"`
	BatteryLevel          float32 `json:"batteryLevel"`
}

func (m ContactMessage) ParseContact() *ContactFromMessage {
	var contact ContactFromMessage
	err := json.Unmarshal(m.Contact, &contact)

	if err != nil {
		log.Println("Parsing contact from message error:", err.Error())
		return nil
	}

	return &contact
}
