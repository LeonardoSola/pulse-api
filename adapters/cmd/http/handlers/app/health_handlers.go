package app_handler

import (
	"net/http"
	"time"

	"github.com/LeonardoSola/pulse-api/adapters/cmd/http/utils"
	chi "github.com/go-chi/chi/v5"
)

func InitApp(router chi.Router) {
	router.Get("/health", health)
}

var startTime = time.Now()

func health(w http.ResponseWriter, r *http.Request) {
	type HealthResponse struct {
		Status string `json:"status"`
		Uptime string `json:"uptime"`
	}

	response := HealthResponse{
		Status: "ok",
		Uptime: time.Since(startTime).String(),
	}

	utils.SendResponse(response, http.StatusOK, w)
}
