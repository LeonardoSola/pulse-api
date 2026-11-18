package message_entity

import (
	"net/http"

	"github.com/LeonardoSola/pulse-api/internal/custom_error"
)

var (
	ErrInvalidRecipientType = custom_error.CustomError{
		Status:  http.StatusBadRequest,
		Message: "Invalid recipient type",
		Info:    "invalid recipient type",
		Level:   custom_error.LevelInfo,
	}

	ErrInvalidRecipient = custom_error.CustomError{
		Status:  http.StatusBadRequest,
		Message: "Invalid recipient",
		Info:    "recipient must be a phone number, a group ID or a full JID",
		Level:   custom_error.LevelInfo,
	}

	ErrEmptyMessage = custom_error.CustomError{
		Status:  http.StatusBadRequest,
		Message: "Message must have text or an image",
		Info:    "text and image_url are both empty",
		Level:   custom_error.LevelInfo,
	}

	ErrImageURLRequired = custom_error.CustomError{
		Status:  http.StatusBadRequest,
		Message: "image_url is required",
		Info:    "image_url is required",
		Level:   custom_error.LevelInfo,
	}
)

var ErrImageDownload = custom_error.CustomError{
	Status:  http.StatusUnprocessableEntity,
	Message: "Could not download image from image_url",
	Info:    "image download failed",
	Level:   custom_error.LevelInfo,
}
