package whastapp

import (
	"context"
	"fmt"

	message_entity "github.com/LeonardoSola/wagon-api/internal/entities/message"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
)

func (s *WhatsappService) SendTextMessage(ctx context.Context, message message_entity.Message) error {
	fmt.Println("SendTextMessage", message)
	client, err := s.getClient(ctx, message.DeviceID)
	if err != nil {
		return err
	}

	recipient, err := types.ParseJID(message.Recipient)
	if err != nil {
		return err
	}

	resp, err := client.SendMessage(
		ctx, recipient, &waE2E.Message{
			Conversation: &message.Text,
		},
	)
	if err != nil {
		return err
	}

	fmt.Println("SendTextMessage", resp)

	return nil
}
