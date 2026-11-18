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
	From          string        `json:"from" db:"-"`
	Recipient     string        `json:"recipient" db:"recipient"`
	RecipientType RecipientType `json:"recipient_type" db:"recipient_type"`
	CreatedAt     time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at" db:"updated_at"`
	DeletedAt     *time.Time    `json:"deleted_at" db:"deleted_at"`
}

func NewMessage(deviceID string, text string, recipient string, recipientType RecipientType, from string) Message {
	now := time.Now()

	return Message{
		ID:            uuid.New().String(),
		DeviceID:      deviceID,
		Text:          text,
		From:          from,
		Recipient:     recipient,
		RecipientType: recipientType,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// NewMessageWithImage builds an outgoing message. recipientType defaults to
// RecipientTypeUser when empty.
func NewMessageWithImage(deviceID string, text string, recipient string, recipientType RecipientType, imageURL string) Message {
	if recipientType == "" {
		recipientType = RecipientTypeUser
	}

	message := NewMessage(deviceID, text, recipient, recipientType, "")
	message.ImageURL = imageURL

	return message
}

// Validate checks the fields required to send a message.
func (m Message) Validate() error {
	if err := m.RecipientType.Validate(); err != nil {
		return err
	}

	if m.Recipient == "" || strings.Contains(m.Recipient, "@") {
		return ErrInvalidRecipient
	}

	if strings.TrimSpace(m.Text) == "" && m.ImageURL == "" {
		return ErrEmptyMessage
	}

	return nil
}

// GetRecipient returns the full WhatsApp JID of the recipient.
func (m Message) GetRecipient() string {
	if m.RecipientType == RecipientTypeGroup {
		return m.Recipient + "@g.us"
	}

	return m.Recipient + "@s.whatsapp.net"
}

// SetRecipient parses a full JID such as "5511999999999@s.whatsapp.net",
// "5511999999999@c.us" or "123456789@g.us" into Recipient and RecipientType.
// Any server other than g.us is treated as a user.
func (m *Message) SetRecipient(recipient string) error {
	user, server, found := strings.Cut(recipient, "@")
	if !found || user == "" || server == "" || strings.Contains(server, "@") {
		return ErrInvalidRecipient
	}

	if server == "g.us" {
		m.RecipientType = RecipientTypeGroup
	} else {
		m.RecipientType = RecipientTypeUser
	}

	m.Recipient = user

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
