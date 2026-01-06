package message_entity

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID            string        `json:"id" db:"id"`
	ExternalID    string        `json:"external_id" db:"external_id"`
	DeviceID      string        `json:"device_id" db:"device_id"`
	DeviceJID     string        `json:"-" db:"-"`
	Text          string        `json:"text" db:"text"`
	ImageURL      string        `json:"image_url" db:"image_url"`
	From          string        `json:"from" db:"from"`
	Recipient     string        `json:"recipient" db:"recipient"`
	RecipientType RecipientType `json:"recipient_type" db:"recipient_type"`
	CreatedAt     time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at" db:"updated_at"`
	DeletedAt     *time.Time    `json:"deleted_at" db:"deleted_at"`
}

func NewMessage(deviceID string, text string, recipient string, recipientType RecipientType) Message {
	return NewMessageWithImage(deviceID, text, recipient, recipientType, "")
}

func NewMessageWithImage(deviceID string, text string, recipient string, recipientType RecipientType, imageURL string) Message {
	if recipientType == "" {
		recipientType = RecipientTypeUser
	}

	return Message{
		ID:            uuid.New().String(),
		DeviceID:      deviceID,
		Text:          text,
		ImageURL:      imageURL,
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

func (m *Message) SetRecipient(recipient string) error {
	recipentSplit := strings.Split(recipient, "@")
	if len(recipentSplit) != 2 {
		return ErrInvalidRecipientType
	}

	if len(recipentSplit) == 1 {
		m.Recipient = recipentSplit[0]
		m.RecipientType = RecipientTypeUser
	}

	if recipentSplit[1] == "g.us" {
		m.RecipientType = RecipientTypeGroup
	} else {
		m.RecipientType = RecipientTypeUser
	}

	m.Recipient = recipentSplit[0]

	return nil
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
