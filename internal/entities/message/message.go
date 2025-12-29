package message_entity

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID            string        `json:"id" db:"id"`
	ExternalID    string        `json:"external_id" db:"external_id"`
	DeviceID      string        `json:"device_id" db:"device_id"`
	Text          string        `json:"text" db:"text"`
	From          string        `json:"from" db:"from"`
	Recipient     string        `json:"recipient" db:"recipient"`
	RecipientType RecipientType `json:"recipient_type" db:"recipient_type"`
	CreatedAt     time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at" db:"updated_at"`
	DeletedAt     *time.Time    `json:"deleted_at" db:"deleted_at"`
}

func NewMessage(deviceID string, text string, recipient string, recipientType RecipientType) Message {
	if recipientType == "" {
		recipientType = RecipientTypeUser
	}

	return Message{
		ID:            uuid.New().String(),
		DeviceID:      deviceID,
		Text:          text,
		Recipient:     recipient,
		RecipientType: recipientType,
	}
}

func (m Message) Validate() error {
	if err := m.RecipientType.Validate(); err != nil {
		return err
	}

	return nil
}

func (m Message) GetRecipient() string {
	if m.RecipientType == RecipientTypeGroup {
		return m.Recipient + "@g.us"
	}

	return m.Recipient + "@s.whatsapp.net"
}

type RecipientType string

const (
	RecipientTypeGroup RecipientType = "group"
	RecipientTypeUser  RecipientType = "user"
)

func (r RecipientType) Validate() error {
	if r != RecipientTypeGroup && r != RecipientTypeUser {
		return ErrInvalidRecipientType
	}

	return nil
}
