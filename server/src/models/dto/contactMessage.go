package dto

type ContactMessage struct {
	User      string             `json:"user"`
	Contact   ContactFromMessage `json:"contact"`
	Signature string             `json:"signature"`
}

type ContactFromMessage struct {
	User                  string  `json:"token"`
	FirstContactTimestamp int64   `json:"firstContactTimestamp"` // milliseconds
	LastContactTimestamp  int64   `json:"lastContactTimestamp"`  // milliseconds
	Distance              float32 `json:"distance"`
	RSSI                  float32 `json:"rssi"`
	BatteryLevel          float32 `json:"batteryLevel"`
}
