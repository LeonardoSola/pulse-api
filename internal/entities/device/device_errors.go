package device_entity

import (
	"fmt"
	"net/http"

	"github.com/LeonardoSola/pulse-api/internal/custom_error"
)

var (
	ErrDeviceSelect = custom_error.CustomError{
		Status:  http.StatusInternalServerError,
		Message: "Error selecting device",
		Info:    "error selecting device",
		Level:   custom_error.LevelError,
	}

	ErrDeviceNotFound = custom_error.CustomError{
		Status:  http.StatusNotFound,
		Message: "Device not found",
		Info:    "device not found",
		Level:   custom_error.LevelInfo,
	}
	ErrDeviceNotLoggedIn = custom_error.CustomError{
		Status:  http.StatusConflict,
		Message: "Device not logged in",
		Info:    "device is not paired with WhatsApp yet, scan the QR code first",
		Level:   custom_error.LevelInfo,
	}

	ErrDeviceInvalidName = custom_error.CustomError{
		Status:  http.StatusBadRequest,
		Message: fmt.Sprintf("Device invalid name, must be between %d and %d characters", NameMinLength, NameMaxLength),
		Info:    "name is invalid",
		Level:   custom_error.LevelInfo,
	}
)
