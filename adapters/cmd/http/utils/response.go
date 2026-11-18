package utils

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/LeonardoSola/pulse-api/internal/config"
	"github.com/LeonardoSola/pulse-api/internal/custom_error"
)

type Success struct {
	Data any `json:"data"`
}

func SendResponse(data any, status int, w http.ResponseWriter) {
	writeJSON(w, status, Success{Data: data})
}

type Error struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Info    string `json:"info,omitempty"`
	Level   int    `json:"level"`
}

// ErrorResponse writes err as a JSON error body. CustomError values keep their
// status and message; any other error is reported as a 500 without leaking its
// details outside development mode.
func ErrorResponse(err error, w http.ResponseWriter) {
	var response Error

	var customErr custom_error.CustomError
	if errors.As(err, &customErr) {
		response = Error{
			Message: customErr.Message,
			Status:  customErr.Status,
		}

		if customErr.Level != custom_error.LevelInfo {
			slog.Log(context.Background(), slog.Level(customErr.Level), customErr.Message, "error", err.Error())
		}

		if customErr.Err != nil && isDevelopment() {
			response.Info = customErr.Err.Error()
			response.Level = int(customErr.Level)
		}
	} else {
		slog.Error("unhandled error", "error", err.Error())

		response = Error{
			Message: "Internal server error",
			Status:  http.StatusInternalServerError,
		}

		if isDevelopment() {
			response.Info = err.Error()
		}
	}

	writeJSON(w, response.Status, response)
}

func isDevelopment() bool {
	return config.Config != nil && config.Config.Env == "development"
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(jsonData)
}
