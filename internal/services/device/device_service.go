package device_service

import (
	"context"
	"errors"
	"log/slog"

	device_entity "github.com/LeonardoSola/pulse-api/internal/entities/device"
	"github.com/LeonardoSola/pulse-api/internal/ports"
)

// QRConnected is sent on the QR channel when the device is already paired.
const QRConnected = "connected"

type DeviceService struct {
	DeviceRepository ports.DeviceRepository
	whatsappService  ports.WhatsappService
}

func NewDeviceService(DeviceRepository ports.DeviceRepository, whatsappService ports.WhatsappService) *DeviceService {
	return &DeviceService{DeviceRepository, whatsappService}
}

func (s *DeviceService) CreateDevice(ctx context.Context, payload CreateDeviceDTO) (device_entity.Device, error) {
	dev := device_entity.NewDevice(payload.Name)

	err := dev.Validate()
	if err != nil {
		return device_entity.Device{}, err
	}

	err = s.DeviceRepository.Create(ctx, dev)
	if err != nil {
		return device_entity.Device{}, err
	}

	return dev, nil
}

func (s *DeviceService) GetDevices(ctx context.Context) ([]device_entity.Device, error) {
	devices, err := s.DeviceRepository.GetAll(ctx)
	if err != nil {
		return []device_entity.Device{}, err
	}

	return devices, nil
}

func (s *DeviceService) GetDevice(ctx context.Context, deviceID string) (device_entity.Device, error) {
	if err := device_entity.ValidateID(deviceID); err != nil {
		return device_entity.Device{}, err
	}

	device, err := s.DeviceRepository.Get(ctx, deviceID)
	if err != nil {
		if errors.Is(err, device_entity.ErrDeviceNotFound) {
			return device_entity.Device{}, device_entity.ErrDeviceNotFound
		}

		return device_entity.Device{}, device_entity.ErrDeviceSelect.SetError(err)
	}

	return device, nil
}

func (s *DeviceService) UpdateDevice(ctx context.Context, deviceID string, payload UpdateDeviceDTO) error {
	device, err := s.GetDevice(ctx, deviceID)
	if err != nil {
		return err
	}

	device = device.WithName(payload.Name)

	err = device.Validate()
	if err != nil {
		return err
	}

	return s.DeviceRepository.Update(ctx, device)
}

// DeleteDevice removes the device and, if it was paired, its WhatsApp session.
func (s *DeviceService) DeleteDevice(ctx context.Context, deviceID string) error {
	device, err := s.GetDevice(ctx, deviceID)
	if err != nil {
		return err
	}

	if device.JID != nil && *device.JID != "" {
		if err := s.whatsappService.DeleteDevice(ctx, *device.JID); err != nil {
			slog.Warn("failed to delete WhatsApp session", "device_id", deviceID, "error", err)
		}
	}

	return s.DeviceRepository.Delete(ctx, deviceID)
}

// ConnectDevice returns a channel of QR codes to pair the device. If the device
// is already paired, the channel yields QRConnected and is closed.
func (s *DeviceService) ConnectDevice(ctx context.Context, deviceID string) (<-chan string, error) {
	device, err := s.GetDevice(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	if device.JID != nil && *device.JID != "" {
		ok, err := s.whatsappService.CheckDevice(ctx, *device.JID)
		if err == nil && ok {
			connected := make(chan string, 1)
			connected <- QRConnected
			close(connected)

			return connected, nil
		}
	}

	qrChan, jidChan, err := s.whatsappService.CreateDevice(ctx)
	if err != nil {
		return nil, err
	}

	go func() {
		for jid := range jidChan {
			if jid == "" {
				continue
			}

			if err := s.DeviceRepository.UpdateJID(context.Background(), deviceID, &jid); err != nil {
				slog.Error("failed to store device JID", "device_id", deviceID, "error", err)
			}
		}
	}()

	return qrChan, nil
}

// StartDevices reconnects every paired device. A device that fails to
// connect is logged and skipped so it does not block the others.
func (s *DeviceService) StartDevices(ctx context.Context) error {
	devices, err := s.DeviceRepository.GetAll(ctx)
	if err != nil {
		return err
	}

	for _, device := range devices {
		if device.JID == nil || *device.JID == "" {
			continue
		}

		if err := s.whatsappService.ConnectDevice(ctx, *device.JID); err != nil {
			slog.Error("failed to connect device", "device_id", device.ID, "error", err)
		}
	}

	return nil
}
