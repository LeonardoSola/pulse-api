package custom_error

import "net/http"

var (
	ErrInvalidBody = CustomError{
		Status:  http.StatusBadRequest,
		Message: "Invalid body",
		Info:    "Invalid body",
		Level:   Info,
	}
)
