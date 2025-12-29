package utils

import (
	"encoding/json"
	"net/http"

	"github.com/LeonardoSola/wagon-api/internal/custom_error"
)

func DecodeBody(r *http.Request, v any) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		return custom_error.ErrInvalidBody.SetError(err)
	}

	return nil
}
