package custom_error

import "net/http"

var (
	ErrUnauthorized = CustomError{
		Status:  http.StatusUnauthorized,
		Message: "Unauthorized",
		Info:    "Unauthorized",
		Level:   Info,
	}
)
