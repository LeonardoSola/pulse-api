package whatsapp

import (
	"context"
	"errors"
	"fmt"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/types"
)

type whatsmeowDevice = store.Device

var errDeviceNotFound = errors.New("whatsapp device not found")

// CreateDevice starts a new pairing session. qrChan yields QR codes followed by
// "success" or "timeout"; jidChan yields the device JID once pairing succeeds.
func (s *WhatsappService) CreateDevice(ctx context.Context) (qrChan <-chan string, jidChan <-chan string, err error) {
	device := s.container.NewDevice()
	if device == nil {
		return nil, nil, fmt.Errorf("failed to create device")
	}

	client := newClient(device)

	qrChan, jidChan, err = s.getQrCodeChannel(ctx, client)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect device: %w", err)
	}

	return qrChan, jidChan, nil
}

func (s *WhatsappService) getDevice(ctx context.Context, deviceJID string) (*store.Device, error) {
	jid, err := types.ParseJID(deviceJID)
	if err != nil {
		return nil, err
	}

	device, err := s.container.GetDevice(ctx, jid)
	if err != nil {
		return nil, err
	}

	if device == nil {
		return nil, errDeviceNotFound
	}

	return device, nil
}

func (s *WhatsappService) CheckDevice(ctx context.Context, deviceJID string) (bool, error) {
	device, err := s.getDevice(ctx, deviceJID)
	if err != nil {
		return false, err
	}

	return device.Initialized, nil
}

// DeleteDevice logs the device out of WhatsApp (falling back to deleting the
// local session) and drops its cached client.
func (s *WhatsappService) DeleteDevice(ctx context.Context, deviceJID string) error {
	if client, ok := cachedClients.remove(deviceJID); ok {
		defer client.Disconnect()

		if client.IsLoggedIn() {
			if err := client.Logout(ctx); err == nil {
				return nil
			}
		}
	}

	device, err := s.getDevice(ctx, deviceJID)
	if err != nil {
		if errors.Is(err, errDeviceNotFound) {
			return nil
		}
		return err
	}

	return device.Delete(ctx)
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

	qrChan := make(chan string)
	jidChan := make(chan string, 1)

	// send forwards a value unless the caller has gone away.
	send := func(value string) bool {
		select {
		case qrChan <- value:
			return true
		case <-ctx.Done():
			return false
		}
	}

	go func() {
		defer close(jidChan)
		defer close(qrChan)

		for evt := range whatsQrChan {
			switch evt.Event {
			case whatsmeow.QRChannelEventCode:
				if !send(evt.Code) {
					// whatsmeow disconnects the client itself once ctx is done.
					return
				}
			case whatsmeow.QRChannelSuccess.Event:
				deviceJID := client.Store.GetJID().ADString()
				registerClient(deviceJID, client)
				jidChan <- deviceJID
				send("success")
				return
			default:
				// timeout or pairing error
				client.Disconnect()
				send(evt.Event)
				return
			}
		}
	}()

	return qrChan, jidChan, nil
}
