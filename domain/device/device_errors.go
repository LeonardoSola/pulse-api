package device

import (
	"fmt"
	"net/http"

	"github.com/LeonardoSola/wagon-api/internal/custom_error"
)

var (
	ErrDeviceSelect = custom_error.CustomError{
		Status:  http.StatusInternalServerError,
		Message: "Error selecting device",
		Info:    "error selecting device",
		Level:   custom_error.Error,
	}

	ErrDeviceNotFound = custom_error.CustomError{
		Status:  http.StatusNotFound,
		Message: "Device not found",
		Info:    "device not found",
		Level:   custom_error.Info,
	}

	ErrDeviceInvalidName = custom_error.CustomError{
		Status:  http.StatusBadRequest,
		Message: fmt.Sprintf("Device invalid name, must be between %d and %d characters", NameMinLength, NameMaxLength),
		Info:    "name is invalid",
		Level:   custom_error.Info,
	}
)
