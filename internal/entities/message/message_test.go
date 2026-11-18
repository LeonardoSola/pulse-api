package message_entity

import (
	"errors"
	"testing"
)

func TestSetRecipient(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantUser  string
		wantType  RecipientType
		wantError bool
	}{
		{name: "user jid", input: "5511999999999@s.whatsapp.net", wantUser: "5511999999999", wantType: RecipientTypeUser},
		{name: "waha c.us jid", input: "5511999999999@c.us", wantUser: "5511999999999", wantType: RecipientTypeUser},
		{name: "group jid", input: "120363000000000000@g.us", wantUser: "120363000000000000", wantType: RecipientTypeGroup},
		{name: "missing server", input: "5511999999999", wantError: true},
		{name: "empty user", input: "@s.whatsapp.net", wantError: true},
		{name: "empty server", input: "5511999999999@", wantError: true},
		{name: "two at signs", input: "a@b@c", wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m Message
			err := m.SetRecipient(tt.input)

			if tt.wantError {
				if !errors.Is(err, ErrInvalidRecipient) {
					t.Fatalf("expected ErrInvalidRecipient, got %v", err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if m.Recipient != tt.wantUser || m.RecipientType != tt.wantType {
				t.Fatalf("got (%q, %q), want (%q, %q)", m.Recipient, m.RecipientType, tt.wantUser, tt.wantType)
			}
		})
	}
}

func TestGetRecipient(t *testing.T) {
	user := Message{Recipient: "5511999999999", RecipientType: RecipientTypeUser}
	if got := user.GetRecipient(); got != "5511999999999@s.whatsapp.net" {
		t.Errorf("user recipient = %q", got)
	}

	group := Message{Recipient: "120363000000000000", RecipientType: RecipientTypeGroup}
	if got := group.GetRecipient(); got != "120363000000000000@g.us" {
		t.Errorf("group recipient = %q", got)
	}
}

func TestSetRecipientRoundTrip(t *testing.T) {
	for _, jid := range []string{"5511999999999@s.whatsapp.net", "120363000000000000@g.us"} {
		var m Message
		if err := m.SetRecipient(jid); err != nil {
			t.Fatalf("SetRecipient(%q): %v", jid, err)
		}
		if got := m.GetRecipient(); got != jid {
			t.Errorf("round trip of %q gave %q", jid, got)
		}
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		message Message
		wantErr error
	}{
		{name: "valid text", message: Message{Recipient: "5511999999999", RecipientType: RecipientTypeUser, Text: "hi"}},
		{name: "valid image without caption", message: Message{Recipient: "5511999999999", RecipientType: RecipientTypeUser, ImageURL: "https://example.com/a.png"}},
		{name: "invalid type", message: Message{Recipient: "5511999999999", RecipientType: "individual", Text: "hi"}, wantErr: ErrInvalidRecipientType},
		{name: "empty recipient", message: Message{RecipientType: RecipientTypeUser, Text: "hi"}, wantErr: ErrInvalidRecipient},
		{name: "recipient still a jid", message: Message{Recipient: "5511@s.whatsapp.net", RecipientType: RecipientTypeUser, Text: "hi"}, wantErr: ErrInvalidRecipient},
		{name: "no content", message: Message{Recipient: "5511999999999", RecipientType: RecipientTypeUser, Text: "  "}, wantErr: ErrEmptyMessage},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.message.Validate()
			if tt.wantErr == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestNewMessageWithImage(t *testing.T) {
	m := NewMessageWithImage("device-id", "caption", "5511999999999", "", "https://example.com/a.png")

	if m.RecipientType != RecipientTypeUser {
		t.Errorf("default recipient type = %q, want %q", m.RecipientType, RecipientTypeUser)
	}
	if m.ID == "" {
		t.Error("ID was not generated")
	}
	if m.CreatedAt.IsZero() || m.UpdatedAt.IsZero() {
		t.Error("timestamps were not set")
	}
	if m.ImageURL != "https://example.com/a.png" {
		t.Errorf("ImageURL = %q", m.ImageURL)
	}
}
