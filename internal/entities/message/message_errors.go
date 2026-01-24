package message_entity

import (
	"net/http"

	"github.com/LeonardoSola/wagon-api/internal/custom_error"
)

var ErrInvalidRecipientType = custom_error.CustomError{
	Status:  http.StatusBadRequest,
	Message: "Invalid recipient type",
	Info:    "invalid recipient type",
	Level:   custom_error.LevelInfo,
}
