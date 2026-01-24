package whastapp

import (
	"context"

	message_entity "github.com/LeonardoSola/wagon-api/internal/entities/message"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
)

func (s *WhatsappService) SendTextMessage(ctx context.Context, message message_entity.Message) (messageID string, err error) {
	client, err := s.getClient(ctx, message.DeviceJID)
	if err != nil {
		return "", err
	}

	recipient, err := types.ParseJID(message.GetRecipient())
	if err != nil {
		return "", err
	}

	sentMessage, err := client.SendMessage(
		ctx, recipient, &waE2E.Message{
			Conversation: &message.Text,
		},
	)
	if err != nil {
		return "", err
	}

	return sentMessage.ID, nil
}
