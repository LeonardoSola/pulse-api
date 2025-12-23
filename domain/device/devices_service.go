package device

import (
	"context"
	"errors"
	"fmt"

	"github.com/LeonardoSola/wagon-api/internal/ports"
	"github.com/jackc/pgx/v5"
)

type DeviceService struct {
	devicesRepository DevicesRepository
	whatsappService   ports.WhatsappService
}

func NewDeviceService(devicesRepository DevicesRepository, whatsappService ports.WhatsappService) *DeviceService {
	return &DeviceService{devicesRepository, whatsappService}
}

func (s *DeviceService) CreateDevice(ctx context.Context, payload CreateDeviceDTO) (Device, error) {
	device := NewDevice(payload.Name)

	err := device.Validate()
	if err != nil {
		return Device{}, err
	}

	err = s.devicesRepository.Create(ctx, device)
	if err != nil {
		return Device{}, err
	}

	return device, nil
}

func (s *DeviceService) GetDevices(ctx context.Context) ([]Device, error) {
	devices, err := s.devicesRepository.GetAll(ctx)
	if err != nil {
		return []Device{}, err
	}

	return devices, nil
}

func (s *DeviceService) GetDevice(ctx context.Context, deviceID string) (Device, error) {
	device, err := s.devicesRepository.Get(ctx, deviceID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Device{}, ErrDeviceNotFound
		}

		return Device{}, ErrDeviceSelect.SetError(err)
	}

	return device, nil
}

func (s *DeviceService) UpdateDevice(ctx context.Context, deviceID string, payload UpdateDeviceDTO) error {
	device, err := s.devicesRepository.Get(ctx, deviceID)
	if err != nil {
		return err
	}

	device = device.WithName(payload.Name)

	err = device.Validate()
	if err != nil {
		return err
	}

	err = s.devicesRepository.Update(ctx, device)
	if err != nil {
		return err
	}

	return nil
}

func (s *DeviceService) DeleteDevice(ctx context.Context, deviceID string) error {
	err := s.devicesRepository.Delete(ctx, deviceID)
	if err != nil {
		return err
	}

	return nil
}

func (s *DeviceService) ConnectDevice(ctx context.Context, deviceID string) (qrChan <-chan string, err error) {
	device, err := s.devicesRepository.Get(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	if device.JID != nil {
		jid, err := s.whatsappService.GetDevice(ctx, *device.JID)
		if err == nil {
			device.JID = &jid
			err = s.devicesRepository.Update(ctx, device)
			if err != nil {
				return nil, err
			}

			return qrChan, nil
		}
	}

	qrChan, jidChan, err := s.whatsappService.CreateDevice(ctx)
	if err != nil {
		return nil, err
	}

	go func() {
		for jid := range jidChan {
			if jid != "" {
				err = s.devicesRepository.UpdateJID(context.Background(), deviceID, &jid)
				if err != nil {
					fmt.Println("error updating jid", err)
					continue
				}
			}
		}
	}()

	return qrChan, nil
}
