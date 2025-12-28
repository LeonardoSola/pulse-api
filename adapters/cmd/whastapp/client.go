package whastapp

import (
	"context"

	"go.mau.fi/whatsmeow"
	waLog "go.mau.fi/whatsmeow/util/log"
)

var cachedClients map[string]*whatsmeow.Client = make(map[string]*whatsmeow.Client)

func (s *WhatsappService) getClient(ctx context.Context, deviceID string) (*whatsmeow.Client, error) {
	if client, ok := cachedClients[deviceID]; ok {
		if client.IsConnected() {
			return client, nil
		}

		err := client.Connect()
		if err != nil {
			return nil, err
		}

		return client, nil
	}

	device, err := s.getDevice(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	client := whatsmeow.NewClient(device, waLog.Stdout("Client", "ERROR", true))
	err = client.Connect()
	if err != nil {
		return nil, err
	}

	cachedClients[deviceID] = client
	return client, nil
}
