package whastapp

import (
	"context"
	"fmt"
	"time"
	"whatsapp-shell/app/models"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

type whatsappPlatform struct {
	Client *whatsmeow.Client
	Event  events.Message
}

func NewWhatsappPlatform() *whatsappPlatform {
	client := GetClient()
	if client == nil {
		fmt.Println("Client is nil")
		return nil
	}

	return &whatsappPlatform{
		Client: client,
	}
}

func (p *whatsappPlatform) SetEvent(event events.Message) {
	p.Event = event
}

func (p *whatsappPlatform) SendMessage(message models.Message, to types.JID) {
	if p.Client == nil {
		fmt.Println("Client is nil, cannot send message")
		return
	}

	_, err := p.Client.SendMessage(context.Background(), to, &waE2E.Message{
		Conversation: proto.String(message.Text),
	})

	if err != nil {
		fmt.Println(err)
	}
}

func (p *whatsappPlatform) Reply(message models.Message, reply string) {
	if p.Client == nil {
		fmt.Println("Client is nil, cannot send message")
		return
	}

	err := p.Client.MarkRead([]types.MessageID{p.Event.Info.ID}, time.Now(), p.Event.Info.Chat, p.Event.Info.Chat)
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = p.Client.SendMessage(context.Background(), p.Event.Info.Chat, &waE2E.Message{
		Conversation: proto.String(reply),
	})

	if err != nil {
		fmt.Println(err)
	}
}
