package device_service

import (
	"context"
	"errors"
	"fmt"

	device_entity "github.com/LeonardoSola/wagon-api/internal/entities/device"
	"github.com/LeonardoSola/wagon-api/internal/ports"
	"github.com/jackc/pgx/v5"
)

type DeviceService struct {
	devicesRepository ports.DevicesRepository
	whatsappService   ports.WhatsappService
}

func NewDeviceService(devicesRepository ports.DevicesRepository, whatsappService ports.WhatsappService) *DeviceService {
	return &DeviceService{devicesRepository, whatsappService}
}

func (s *DeviceService) CreateDevice(ctx context.Context, payload CreateDeviceDTO) (device_entity.Device, error) {
	dev := device_entity.NewDevice(payload.Name)

	err := dev.Validate()
	if err != nil {
		return device_entity.Device{}, err
	}

	err = s.devicesRepository.Create(ctx, dev)
	if err != nil {
		return device_entity.Device{}, err
	}

	return dev, nil
}

func (s *DeviceService) GetDevices(ctx context.Context) ([]device_entity.Device, error) {
	devices, err := s.devicesRepository.GetAll(ctx)
	if err != nil {
		return []device_entity.Device{}, err
	}

	return devices, nil
}

func (s *DeviceService) GetDevice(ctx context.Context, deviceID string) (device_entity.Device, error) {
	device, err := s.devicesRepository.Get(ctx, deviceID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return device_entity.Device{}, device_entity.ErrDeviceNotFound
		}

		return device_entity.Device{}, device_entity.ErrDeviceSelect.SetError(err)
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
		ok, err := s.whatsappService.CheckDevice(ctx, *device.JID)
		if err == nil && ok {
			device.IsActive = true
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
