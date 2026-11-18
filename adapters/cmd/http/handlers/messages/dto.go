package message_handler

import message_entity "github.com/LeonardoSola/pulse-api/internal/entities/message"

type CreateMessageTextDTO struct {
	Recipient     string                       `json:"recipient"`
	RecipientType message_entity.RecipientType `json:"recipient_type"`
	Text          string                       `json:"text"`
	ImageURL      string                       `json:"image_url"`
}
