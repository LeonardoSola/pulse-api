package utils

import (
	"encoding/json"
	"net/http"

	"github.com/LeonardoSola/wagon-api/internal/config"
	"github.com/LeonardoSola/wagon-api/internal/custom_error"
)

type Success struct {
	Data any `json:"data"`
}

func SendResponse(data any, status int, w http.ResponseWriter) {
	w.WriteHeader(status)

	response := Success{
		Data: data,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(jsonData)
}

type Error struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Info    string `json:"info,omitempty"`
	Level   int    `json:"level"`
}

func ErrorResponse(data error, w http.ResponseWriter) {
	var response Error
	if e, ok := data.(custom_error.CustomError); ok {
		response = Error{
			Message: e.Message,
			Status:  e.Status,
		}

		if e.Err != nil && config.Config.Env == "development" {
			response.Info = e.Err.Error()
			response.Level = int(e.Level)
		}
	} else {
		response = Error{
			Message: data.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(jsonData)
}
