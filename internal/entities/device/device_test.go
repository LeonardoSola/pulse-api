package device_entity

import (
	"errors"
	"strings"
	"testing"
)

func TestValidateName(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{name: "too short", value: "ab", wantErr: true},
		{name: "min length", value: "abc"},
		{name: "max length", value: strings.Repeat("a", NameMaxLength)},
		{name: "too long", value: strings.Repeat("a", NameMaxLength+1), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewDevice(tt.value).Validate()
			if tt.wantErr != errors.Is(err, ErrDeviceInvalidName) {
				t.Fatalf("Validate(%d chars) = %v", len(tt.value), err)
			}
		})
	}
}

func TestWithName(t *testing.T) {
	d := NewDevice("original")

	if got := d.WithName(nil).Name; got != "original" {
		t.Errorf("nil name changed the device name to %q", got)
	}

	name := "renamed"
	if got := d.WithName(&name).Name; got != "renamed" {
		t.Errorf("WithName = %q", got)
	}

	if d.Name != "original" {
		t.Error("WithName mutated the receiver")
	}
}

func TestValidateID(t *testing.T) {
	if err := ValidateID(NewDevice("abc").ID); err != nil {
		t.Errorf("generated ID rejected: %v", err)
	}

	for _, id := range []string{"", "not-a-uuid", "1"} {
		if err := ValidateID(id); !errors.Is(err, ErrDeviceNotFound) {
			t.Errorf("ValidateID(%q) = %v, want ErrDeviceNotFound", id, err)
		}
	}
}
