package whastapp

import (
	"context"
	"fmt"

	"github.com/LeonardoSola/wagon-api/domain/message"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func (s *WhatsappService) SendTextMessage(ctx context.Context, message message.Message) error {
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

	fmt.Println(resp)

	return err
}
