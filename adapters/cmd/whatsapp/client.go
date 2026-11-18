package whatsapp

import (
	"context"
	"sync"

	"github.com/LeonardoSola/pulse-api/adapters/cmd/whatsapp/events"
	"go.mau.fi/whatsmeow"
	waLog "go.mau.fi/whatsmeow/util/log"
)

// clientCache keeps one connected whatsmeow client per device JID.
type clientCache struct {
	mu      sync.Mutex
	clients map[string]*whatsmeow.Client
}

var cachedClients = &clientCache{clients: make(map[string]*whatsmeow.Client)}

func (c *clientCache) get(jid string) (*whatsmeow.Client, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	client, ok := c.clients[jid]
	return client, ok
}

func (c *clientCache) set(jid string, client *whatsmeow.Client) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.clients[jid] = client
}

func (c *clientCache) remove(jid string) (*whatsmeow.Client, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	client, ok := c.clients[jid]
	delete(c.clients, jid)
	return client, ok
}

func newClient(device *whatsmeowDevice) *whatsmeow.Client {
	return whatsmeow.NewClient(device, waLog.Stdout("Client", "ERROR", true))
}

// registerClient attaches the incoming-event handler and caches the client.
func registerClient(deviceJID string, client *whatsmeow.Client) {
	client.AddEventHandler(events.EventHandler(deviceJID))
	cachedClients.set(deviceJID, client)
}

// getClient returns a connected client for the device, creating it if needed.
func (s *WhatsappService) getClient(ctx context.Context, deviceJID string) (*whatsmeow.Client, error) {
	if client, ok := cachedClients.get(deviceJID); ok {
		if client.IsConnected() {
			return client, nil
		}

		if err := client.Connect(); err != nil {
			return nil, err
		}

		return client, nil
	}

	device, err := s.getDevice(ctx, deviceJID)
	if err != nil {
		return nil, err
	}

	client := newClient(device)
	if err := client.Connect(); err != nil {
		return nil, err
	}

	registerClient(deviceJID, client)
	return client, nil
}

// ConnectDevice makes sure a paired device has a connected client listening for events.
func (s *WhatsappService) ConnectDevice(ctx context.Context, deviceJID string) error {
	_, err := s.getClient(ctx, deviceJID)
	return err
}

// DisconnectAll disconnects every cached client (used on shutdown).
func DisconnectAll() {
	cachedClients.mu.Lock()
	defer cachedClients.mu.Unlock()

	for _, client := range cachedClients.clients {
		client.Disconnect()
	}
}
