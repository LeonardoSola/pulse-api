package message_service

import (
	"context"
	"errors"
	"testing"

	device_entity "github.com/LeonardoSola/pulse-api/internal/entities/device"
	message_entity "github.com/LeonardoSola/pulse-api/internal/entities/message"
)

type fakeMessageRepository struct {
	created []message_entity.Message
}

func (r *fakeMessageRepository) Create(_ context.Context, m message_entity.Message) error {
	r.created = append(r.created, m)
	return nil
}
func (r *fakeMessageRepository) Get(context.Context, string) (message_entity.Message, error) {
	return message_entity.Message{}, nil
}
func (r *fakeMessageRepository) GetAll(context.Context) ([]message_entity.Message, error) {
	return nil, nil
}
func (r *fakeMessageRepository) Update(context.Context, message_entity.Message) error { return nil }
func (r *fakeMessageRepository) Delete(context.Context, string) error                 { return nil }

type fakeDeviceRepository struct {
	device device_entity.Device
}

func (r *fakeDeviceRepository) GetAll(context.Context) ([]device_entity.Device, error) {
	return []device_entity.Device{r.device}, nil
}
func (r *fakeDeviceRepository) Get(_ context.Context, id string) (device_entity.Device, error) {
	if id != r.device.ID {
		return device_entity.Device{}, device_entity.ErrDeviceNotFound
	}
	return r.device, nil
}
func (r *fakeDeviceRepository) GetByJID(_ context.Context, jid string) (device_entity.Device, error) {
	if r.device.JID == nil || *r.device.JID != jid {
		return device_entity.Device{}, device_entity.ErrDeviceNotFound
	}
	return r.device, nil
}
func (r *fakeDeviceRepository) Create(context.Context, device_entity.Device) error { return nil }
func (r *fakeDeviceRepository) Update(context.Context, device_entity.Device) error { return nil }
func (r *fakeDeviceRepository) Delete(context.Context, string) error               { return nil }
func (r *fakeDeviceRepository) UpdateJID(context.Context, string, *string) error   { return nil }

type fakeWhatsapp struct {
	text  []message_entity.Message
	image []message_entity.Message
}

func (w *fakeWhatsapp) CreateDevice(context.Context) (<-chan string, <-chan string, error) {
	return nil, nil, nil
}
func (w *fakeWhatsapp) DeleteDevice(context.Context, string) error        { return nil }
func (w *fakeWhatsapp) CheckDevice(context.Context, string) (bool, error) { return true, nil }
func (w *fakeWhatsapp) ConnectDevice(context.Context, string) error       { return nil }
func (w *fakeWhatsapp) SendTextMessage(_ context.Context, m message_entity.Message) (string, error) {
	w.text = append(w.text, m)
	return "TEXT-ID", nil
}
func (w *fakeWhatsapp) SendImageMessage(_ context.Context, m message_entity.Message) (string, error) {
	w.image = append(w.image, m)
	return "IMAGE-ID", nil
}

type published struct {
	queue   string
	message any
}

type fakeMessaging struct {
	published []published
}

func (m *fakeMessaging) Publish(_ context.Context, queue string, message any) error {
	m.published = append(m.published, published{queue, message})
	return nil
}

const deviceJID = "5511000000000.0:1@s.whatsapp.net"

func setup(paired bool) (*MessageService, *fakeMessageRepository, *fakeWhatsapp, *fakeMessaging, device_entity.Device) {
	device := device_entity.NewDevice("phone")
	if paired {
		jid := deviceJID
		device.JID = &jid
	}

	messages := &fakeMessageRepository{}
	wa := &fakeWhatsapp{}
	queue := &fakeMessaging{}
	svc := NewMessageService(messages, &fakeDeviceRepository{device: device}, wa, queue)

	return svc, messages, wa, queue, device
}

func TestCreateMessageSendsTextAndStoresExternalID(t *testing.T) {
	svc, messages, wa, _, device := setup(true)

	msg := message_entity.NewMessageWithImage(device.ID, "hello", "5511999999999", "", "")
	if err := svc.CreateMessage(context.Background(), msg); err != nil {
		t.Fatal(err)
	}

	if len(wa.text) != 1 || len(wa.image) != 0 {
		t.Fatalf("text sends = %d, image sends = %d", len(wa.text), len(wa.image))
	}
	if wa.text[0].DeviceJID != deviceJID {
		t.Errorf("DeviceJID = %q", wa.text[0].DeviceJID)
	}
	if len(messages.created) != 1 || messages.created[0].ExternalID != "TEXT-ID" {
		t.Errorf("stored messages = %+v", messages.created)
	}
}

func TestCreateMessageWithImageURLSendsImage(t *testing.T) {
	svc, _, wa, _, device := setup(true)

	msg := message_entity.NewMessageWithImage(device.ID, "caption", "5511999999999", "", "https://example.com/a.png")
	if err := svc.CreateMessage(context.Background(), msg); err != nil {
		t.Fatal(err)
	}

	if len(wa.image) != 1 || len(wa.text) != 0 {
		t.Fatalf("text sends = %d, image sends = %d", len(wa.text), len(wa.image))
	}
}

func TestCreateMessageErrors(t *testing.T) {
	t.Run("unpaired device", func(t *testing.T) {
		svc, _, wa, _, device := setup(false)
		msg := message_entity.NewMessageWithImage(device.ID, "hello", "5511999999999", "", "")

		if err := svc.CreateMessage(context.Background(), msg); !errors.Is(err, device_entity.ErrDeviceNotLoggedIn) {
			t.Fatalf("expected ErrDeviceNotLoggedIn, got %v", err)
		}
		if len(wa.text) != 0 {
			t.Error("message must not be sent")
		}
	})

	t.Run("malformed device id", func(t *testing.T) {
		svc, _, _, _, _ := setup(true)
		msg := message_entity.NewMessageWithImage("nope", "hello", "5511999999999", "", "")

		if err := svc.CreateMessage(context.Background(), msg); !errors.Is(err, device_entity.ErrDeviceNotFound) {
			t.Fatalf("expected ErrDeviceNotFound, got %v", err)
		}
	})

	t.Run("invalid payload", func(t *testing.T) {
		svc, _, wa, _, device := setup(true)
		msg := message_entity.NewMessageWithImage(device.ID, "", "5511999999999", "", "")

		if err := svc.CreateMessage(context.Background(), msg); !errors.Is(err, message_entity.ErrEmptyMessage) {
			t.Fatalf("expected ErrEmptyMessage, got %v", err)
		}
		if len(wa.text) != 0 {
			t.Error("message must not be sent")
		}
	})
}

func TestReceiveMessageStoresAndPublishes(t *testing.T) {
	svc, messages, _, queue, device := setup(true)

	msg := message_entity.NewMessage("", "hi there", "5511999999999", message_entity.RecipientTypeUser, "5511999999999@s.whatsapp.net")
	msg.DeviceJID = deviceJID

	if err := svc.ReceiveMessage(context.Background(), msg); err != nil {
		t.Fatal(err)
	}

	if len(messages.created) != 1 || messages.created[0].DeviceID != device.ID {
		t.Fatalf("stored messages = %+v", messages.created)
	}
	if len(queue.published) != 1 || queue.published[0].queue != "messages.received" {
		t.Fatalf("published = %+v", queue.published)
	}
}
