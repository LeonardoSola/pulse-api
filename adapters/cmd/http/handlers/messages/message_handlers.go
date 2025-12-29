package message_handler

import (
	"net/http"

	"github.com/LeonardoSola/wagon-api/adapters/bootstrap"
	"github.com/LeonardoSola/wagon-api/adapters/cmd/http/utils"
	message_entity "github.com/LeonardoSola/wagon-api/internal/entities/message"
	"github.com/go-chi/chi/v5"
)

func InitMessages(router *chi.Mux) {
	router.Post("/devices/{deviceId}/messages/sendText", createMessage)
}

func createMessage(w http.ResponseWriter, r *http.Request) {
	deviceID := chi.URLParam(r, "deviceId")

	var payload CreateMessageTextDTO
	err := utils.DecodeBody(r, &payload)
	if err != nil {
		utils.ErrorResponse(err, w)
		return
	}

	err = bootstrap.Services.MessageService.CreateMessage(
		r.Context(),
		message_entity.NewMessage(
			deviceID,
			payload.Text,
			payload.Recipient,
			payload.RecipientType,
		),
	)
	if err != nil {
		utils.ErrorResponse(err, w)
		return
	}

	utils.SendResponse("success", http.StatusCreated, w)
}
