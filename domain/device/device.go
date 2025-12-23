package device

import (
	"time"

	"github.com/google/uuid"
)

type Device struct {
	ID        string     `json:"id" db:"id"`
	JID       *string    `json:"jid" db:"jid"`
	Name      string     `json:"name" db:"name"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
	IsActive  bool       `json:"is_active" db:"is_active"`
}

func NewDevice(name string) Device {
	return Device{
		ID:        uuid.New().String(),
		JID:       nil,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: nil,
		IsActive:  false,
	}
}

func (d Device) Validate() error {
	if len(d.Name) < NameMinLength || len(d.Name) > NameMaxLength {
		return ErrDeviceInvalidName
	}

	return nil
}

func (d Device) WithName(name *string) Device {
	if name != nil {
		d.Name = *name
	}

	return d
}

func (d Device) WithJID(jid *string) Device {
	if jid != nil {
		d.JID = jid
	}

	return d
}
