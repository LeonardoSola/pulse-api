package whastapp

import (
	"context"

	"github.com/LeonardoSola/wagon-api/adapters/cmd/whastapp/events"
	"go.mau.fi/whatsmeow"
	waLog "go.mau.fi/whatsmeow/util/log"
)

var cachedClients map[string]*whatsmeow.Client = make(map[string]*whatsmeow.Client)

func (s *WhatsappService) getClient(ctx context.Context, deviceJID string) (*whatsmeow.Client, error) {
	if client, ok := cachedClients[deviceJID]; ok {
		if client.IsConnected() {
			return client, nil
		}

		err := client.Connect()
		if err != nil {
			return nil, err
		}

		return client, nil
	}

	device, err := s.getDevice(ctx, deviceJID)
	if err != nil {
		return nil, err
	}

	client := whatsmeow.NewClient(device, waLog.Stdout("Client", "ERROR", true))
	client.AddEventHandler(events.EventHandler(deviceJID))
	err = client.Connect()
	if err != nil {
		return nil, err
	}

	cachedClients[deviceJID] = client
	return client, nil
}

func (s *WhatsappService) ConnectDevice(ctx context.Context, deviceJID string) error {
	client, err := s.getClient(ctx, deviceJID)
	if err != nil {
		return err
	}

	if client.IsConnected() {
		return nil
	}

	client.AddEventHandler(events.EventHandler(deviceJID))

	err = client.Connect()
	if err != nil {
		return err
	}

	cachedClients[deviceJID] = client

	return nil
}
