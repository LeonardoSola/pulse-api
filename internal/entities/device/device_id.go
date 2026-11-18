package device_entity

import "github.com/google/uuid"

// ValidateID returns ErrDeviceNotFound when id is not a valid UUID, so that
// malformed IDs never reach the database.
func ValidateID(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return ErrDeviceNotFound
	}

	return nil
}
