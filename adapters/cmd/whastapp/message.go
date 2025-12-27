package whastapp

import (
	"context"

	message_entity "github.com/LeonardoSola/wagon-api/internal/entities/message"
	"github.com/LeonardoSola/wagon-api/internal/logger"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func (s *WhatsappService) SendTextMessage(ctx context.Context, message message_entity.Message) error {
	device, err := s.getDevice(ctx, message.DeviceID)
	if err != nil {
		return err
	}

	jid, err := types.ParseJID(message.Recipient)
	if err != nil {
		return err
	}

	client := whatsmeow.NewClient(device, waLog.Stdout("Client", "ERROR", true))
	resp, err := client.SendMessage(
		ctx, jid, &waE2E.Message{
			Conversation: &message.Text,
		},
	)

	logger.Log.Debug("Message sent", "response", resp)

	return err
}
