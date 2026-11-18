package message_handler

import (
	"errors"
	"testing"

	message_entity "github.com/LeonardoSola/pulse-api/internal/entities/message"
)

func TestToMessage(t *testing.T) {
	tests := []struct {
		name     string
		payload  CreateMessageTextDTO
		wantUser string
		wantType message_entity.RecipientType
		wantErr  error
	}{
		{
			name:     "bare number defaults to user",
			payload:  CreateMessageTextDTO{Recipient: "5511999999999", Text: "hi"},
			wantUser: "5511999999999", wantType: message_entity.RecipientTypeUser,
		},
		{
			name:     "bare group id with explicit type",
			payload:  CreateMessageTextDTO{Recipient: "120363000000000000", RecipientType: message_entity.RecipientTypeGroup, Text: "hi"},
			wantUser: "120363000000000000", wantType: message_entity.RecipientTypeGroup,
		},
		{
			name:     "full group jid overrides type",
			payload:  CreateMessageTextDTO{Recipient: "120363000000000000@g.us", Text: "hi"},
			wantUser: "120363000000000000", wantType: message_entity.RecipientTypeGroup,
		},
		{
			name:    "malformed jid",
			payload: CreateMessageTextDTO{Recipient: "@s.whatsapp.net", Text: "hi"},
			wantErr: message_entity.ErrInvalidRecipient,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := tt.payload.toMessage("device-id")
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if msg.Recipient != tt.wantUser || msg.RecipientType != tt.wantType || msg.DeviceID != "device-id" {
				t.Errorf("got %+v", msg)
			}
		})
	}
}
