package whastapp

import (
	"context"
	"fmt"

	"github.com/LeonardoSola/wagon-api/adapters/cmd/whastapp/events"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func (s *WhatsappService) CreateDevice(ctx context.Context) (qrChan <-chan string, jidChan <-chan string, err error) {
	device := s.container.NewDevice()
	if device == nil {
		return nil, nil, fmt.Errorf("failed to create device")
	}

	client := whatsmeow.NewClient(device, waLog.Stdout("Client", "ERROR", true))
	client.AddEventHandler(events.EventHandler)

	qrChan, jidChan, err = s.getQrCodeChannel(ctx, client)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect device: %w", err)
	}

	return qrChan, jidChan, nil
}

func (s *WhatsappService) GetDevice(ctx context.Context, deviceID string) (string, error) {
	jid, err := types.ParseJID(deviceID)
	if err != nil {
		return "", err
	}

	device, err := s.container.GetDevice(ctx, jid)
	if err != nil {
		return "", err
	}

	if device == nil {
		return "", fmt.Errorf("device not found")
	}

	return jid.String(), nil
}

func (s *WhatsappService) DeleteDevice(ctx context.Context, deviceID string) error {
	jid, err := types.ParseJID(deviceID)
	if err != nil {
		return err
	}

	device, err := s.container.GetDevice(ctx, jid)
	if err != nil {
		return err
	}

	if device == nil {
		return fmt.Errorf("device not found")
	}

	err = device.Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *WhatsappService) getQrCodeChannel(ctx context.Context, client *whatsmeow.Client) (qrCodeChan <-chan string, jid <-chan string, err error) {
	whatsQrChan, err := client.GetQRChannel(ctx)
	if err != nil {
		return nil, nil, err
	}

	err = client.Connect()
	if err != nil {
		return nil, nil, err
	}

	// transform the whatsQrChan to a channel of strings
	qrChan := make(chan string)
	jidChan := make(chan string)
	go func() {
	loop:
		for evt := range whatsQrChan {
			switch evt.Event {
			case "code":
				qrCode := evt.Code
				qrChan <- qrCode
			case "success":
				qrChan <- "success"
				break loop
			case "timeout":
				qrChan <- "timeout"
				break loop
			}
		}

		close(qrChan)

		jidChan <- client.Store.GetJID().ADString()
		close(jidChan)
	}()

	return qrChan, jidChan, nil
}
