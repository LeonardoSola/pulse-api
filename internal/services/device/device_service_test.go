package device_service

import (
	"context"
	"errors"
	"testing"

	"github.com/LeonardoSola/pulse-api/internal/custom_error"
	device_entity "github.com/LeonardoSola/pulse-api/internal/entities/device"
	message_entity "github.com/LeonardoSola/pulse-api/internal/entities/message"
)

// fakeDeviceRepository is an in-memory ports.DeviceRepository.
type fakeDeviceRepository struct {
	devices map[string]device_entity.Device
	getErr  error
	deleted []string
}

func newFakeRepo(devices ...device_entity.Device) *fakeDeviceRepository {
	repo := &fakeDeviceRepository{devices: map[string]device_entity.Device{}}
	for _, d := range devices {
		repo.devices[d.ID] = d
	}
	return repo
}

func (r *fakeDeviceRepository) GetAll(context.Context) ([]device_entity.Device, error) {
	var out []device_entity.Device
	for _, d := range r.devices {
		out = append(out, d)
	}
	return out, nil
}

func (r *fakeDeviceRepository) Get(_ context.Context, id string) (device_entity.Device, error) {
	if r.getErr != nil {
		return device_entity.Device{}, r.getErr
	}
	d, ok := r.devices[id]
	if !ok {
		return device_entity.Device{}, device_entity.ErrDeviceNotFound
	}
	return d, nil
}

func (r *fakeDeviceRepository) GetByJID(context.Context, string) (device_entity.Device, error) {
	return device_entity.Device{}, device_entity.ErrDeviceNotFound
}

func (r *fakeDeviceRepository) Create(_ context.Context, d device_entity.Device) error {
	r.devices[d.ID] = d
	return nil
}

func (r *fakeDeviceRepository) Update(_ context.Context, d device_entity.Device) error {
	r.devices[d.ID] = d
	return nil
}

func (r *fakeDeviceRepository) Delete(_ context.Context, id string) error {
	r.deleted = append(r.deleted, id)
	delete(r.devices, id)
	return nil
}

func (r *fakeDeviceRepository) UpdateJID(_ context.Context, id string, jid *string) error {
	d := r.devices[id]
	d.JID = jid
	r.devices[id] = d
	return nil
}

// fakeWhatsapp is a ports.WhatsappService that records calls.
type fakeWhatsapp struct {
	paired     bool
	connectErr map[string]error
	connected  []string
	deleted    []string
	created    int
}

func (w *fakeWhatsapp) CreateDevice(context.Context) (<-chan string, <-chan string, error) {
	w.created++
	qr := make(chan string, 1)
	qr <- "qr-code"
	close(qr)
	jid := make(chan string)
	close(jid)
	return qr, jid, nil
}

func (w *fakeWhatsapp) DeleteDevice(_ context.Context, jid string) error {
	w.deleted = append(w.deleted, jid)
	return nil
}

func (w *fakeWhatsapp) CheckDevice(context.Context, string) (bool, error) { return w.paired, nil }

func (w *fakeWhatsapp) SendTextMessage(context.Context, message_entity.Message) (string, error) {
	return "", nil
}

func (w *fakeWhatsapp) SendImageMessage(context.Context, message_entity.Message) (string, error) {
	return "", nil
}

func (w *fakeWhatsapp) ConnectDevice(_ context.Context, jid string) error {
	w.connected = append(w.connected, jid)
	return w.connectErr[jid]
}

func pairedDevice(name, jid string) device_entity.Device {
	d := device_entity.NewDevice(name)
	d.JID = &jid
	return d
}

func TestGetDevice(t *testing.T) {
	existing := device_entity.NewDevice("phone")
	svc := NewDeviceService(newFakeRepo(existing), &fakeWhatsapp{})

	if _, err := svc.GetDevice(context.Background(), existing.ID); err != nil {
		t.Fatalf("existing device: %v", err)
	}

	for _, id := range []string{"not-a-uuid", device_entity.NewDevice("other").ID} {
		if _, err := svc.GetDevice(context.Background(), id); !errors.Is(err, device_entity.ErrDeviceNotFound) {
			t.Errorf("GetDevice(%q) = %v, want ErrDeviceNotFound", id, err)
		}
	}
}

func TestGetDeviceWrapsUnexpectedErrors(t *testing.T) {
	repo := newFakeRepo()
	repo.getErr = errors.New("connection refused")
	svc := NewDeviceService(repo, &fakeWhatsapp{})

	_, err := svc.GetDevice(context.Background(), device_entity.NewDevice("phone").ID)

	var customErr custom_error.CustomError
	if !errors.As(err, &customErr) || customErr.Message != device_entity.ErrDeviceSelect.Message {
		t.Fatalf("expected ErrDeviceSelect, got %v", err)
	}
	if !errors.Is(err, repo.getErr) {
		t.Errorf("the repository error should be wrapped, got %v", err)
	}
}

func TestCreateDeviceValidatesName(t *testing.T) {
	svc := NewDeviceService(newFakeRepo(), &fakeWhatsapp{})

	if _, err := svc.CreateDevice(context.Background(), CreateDeviceDTO{Name: "x"}); !errors.Is(err, device_entity.ErrDeviceInvalidName) {
		t.Fatalf("expected ErrDeviceInvalidName, got %v", err)
	}

	device, err := svc.CreateDevice(context.Background(), CreateDeviceDTO{Name: "Sales phone"})
	if err != nil {
		t.Fatal(err)
	}
	if device.Name != "Sales phone" || device.JID != nil {
		t.Errorf("unexpected device %+v", device)
	}
}

func TestConnectDeviceAlreadyPaired(t *testing.T) {
	device := pairedDevice("phone", "5511999999999.0:1@s.whatsapp.net")
	wa := &fakeWhatsapp{paired: true}
	svc := NewDeviceService(newFakeRepo(device), wa)

	qr, err := svc.ConnectDevice(context.Background(), device.ID)
	if err != nil {
		t.Fatal(err)
	}

	var events []string
	for evt := range qr {
		events = append(events, evt)
	}

	if len(events) != 1 || events[0] != QRConnected {
		t.Errorf("events = %v, want [%s]", events, QRConnected)
	}
	if wa.created != 0 {
		t.Error("a new pairing session must not be started for a paired device")
	}
}

func TestConnectDeviceStartsPairing(t *testing.T) {
	device := device_entity.NewDevice("phone")
	wa := &fakeWhatsapp{}
	svc := NewDeviceService(newFakeRepo(device), wa)

	qr, err := svc.ConnectDevice(context.Background(), device.ID)
	if err != nil {
		t.Fatal(err)
	}

	if got := <-qr; got != "qr-code" {
		t.Errorf("first event = %q, want qr-code", got)
	}
	if wa.created != 1 {
		t.Errorf("CreateDevice called %d times, want 1", wa.created)
	}
}

func TestDeleteDeviceRemovesWhatsappSession(t *testing.T) {
	jid := "5511999999999.0:1@s.whatsapp.net"
	device := pairedDevice("phone", jid)
	repo := newFakeRepo(device)
	wa := &fakeWhatsapp{}
	svc := NewDeviceService(repo, wa)

	if err := svc.DeleteDevice(context.Background(), device.ID); err != nil {
		t.Fatal(err)
	}

	if len(wa.deleted) != 1 || wa.deleted[0] != jid {
		t.Errorf("whatsapp sessions deleted = %v", wa.deleted)
	}
	if len(repo.deleted) != 1 {
		t.Errorf("repository deletes = %v", repo.deleted)
	}

	if err := svc.DeleteDevice(context.Background(), device.ID); !errors.Is(err, device_entity.ErrDeviceNotFound) {
		t.Errorf("second delete = %v, want ErrDeviceNotFound", err)
	}
}

func TestStartDevicesSkipsUnpairedAndContinuesOnError(t *testing.T) {
	broken := pairedDevice("broken", "111.0:1@s.whatsapp.net")
	ok := pairedDevice("ok", "222.0:1@s.whatsapp.net")
	unpaired := device_entity.NewDevice("unpaired")

	wa := &fakeWhatsapp{connectErr: map[string]error{*broken.JID: errors.New("boom")}}
	svc := NewDeviceService(newFakeRepo(broken, ok, unpaired), wa)

	if err := svc.StartDevices(context.Background()); err != nil {
		t.Fatalf("StartDevices: %v", err)
	}

	if len(wa.connected) != 2 {
		t.Errorf("connected = %v, want both paired devices", wa.connected)
	}
}
