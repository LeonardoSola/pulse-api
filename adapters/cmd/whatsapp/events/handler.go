package events

import (
	"go.mau.fi/whatsmeow/types/events"
)

func EventHandler(deviceJID string) func(evt interface{}) {
	return func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			MessageHandler(deviceJID, v)
		}
	}
}
