package message_entity

import "time"

type Message struct {
	ID        string    `json:"id" db:"id"`
	DeviceID  string    `json:"device_id" db:"device_id"`
	Text      string    `json:"text" db:"text"`
	From      string    `json:"from" db:"from"`
	Recipient string    `json:"recipient" db:"recipient"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}
