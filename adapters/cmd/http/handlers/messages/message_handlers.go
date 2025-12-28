package message_handler

import (
	"net/http"
	"time"

	"github.com/LeonardoSola/wagon-api/adapters/bootstrap"
	"github.com/LeonardoSola/wagon-api/adapters/cmd/http/utils"
	message_entity "github.com/LeonardoSola/wagon-api/internal/entities/message"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func InitMessages(router *chi.Mux) {
	router.Post("/devices/{deviceId}/messages", createMessage)
}

func createMessage(w http.ResponseWriter, r *http.Request) {
	deviceID := chi.URLParam(r, "deviceId")

	var payload CreateMessageDTO
	err := utils.DecodeBody(r, &payload)
	if err != nil {
		utils.ErrorResponse(err, w)
		return
	}

	err = bootstrap.Services.MessageService.CreateMessage(r.Context(), message_entity.Message{
		Recipient: payload.Recipient,
		Text:      payload.Text,
		From:      deviceID,
		DeviceID:  deviceID,
		Timestamp: time.Now(),
		ID:        uuid.New().String(),
	})
	if err != nil {
		utils.ErrorResponse(err, w)
		return
	}

	utils.SendResponse("success", http.StatusCreated, w)
}
