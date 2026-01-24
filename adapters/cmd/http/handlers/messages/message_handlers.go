package message_handler

import (
	"fmt"
	"net/http"

	"github.com/LeonardoSola/pulse-api/adapters/bootstrap"
	"github.com/LeonardoSola/pulse-api/adapters/cmd/http/utils"
	message_entity "github.com/LeonardoSola/pulse-api/internal/entities/message"
	"github.com/go-chi/chi/v5"
)

func InitMessages(router *chi.Mux) {
	router.Post("/devices/{deviceId}/messages/sendText", createTextMessage)
	router.Post("/devices/{deviceId}/messages/sendImage", createImageMessage)
}

func createTextMessage(w http.ResponseWriter, r *http.Request) {
	deviceID := chi.URLParam(r, "deviceId")

	var payload CreateMessageTextDTO
	err := utils.DecodeBody(r, &payload)
	if err != nil {
		utils.ErrorResponse(err, w)
		return
	}

	err = bootstrap.Services.MessageService.CreateMessage(
		r.Context(),
		message_entity.NewMessageWithImage(
			deviceID,
			payload.Text,
			payload.Recipient,
			payload.RecipientType,
			payload.ImageURL,
		),
	)
	if err != nil {
		utils.ErrorResponse(err, w)
		return
	}

	utils.SendResponse("success", http.StatusCreated, w)
}

func createImageMessage(w http.ResponseWriter, r *http.Request) {
	deviceID := chi.URLParam(r, "deviceId")

	var payload CreateMessageTextDTO
	err := utils.DecodeBody(r, &payload)
	if err != nil {
		utils.ErrorResponse(err, w)
		return
	}

	if payload.ImageURL == "" {
		utils.ErrorResponse(fmt.Errorf("image_url is required"), w)
		return
	}

	err = bootstrap.Services.MessageService.CreateMessage(
		r.Context(),
		message_entity.NewMessageWithImage(
			deviceID,
			payload.Text,
			payload.Recipient,
			payload.RecipientType,
			payload.ImageURL,
		),
	)
	if err != nil {
		utils.ErrorResponse(err, w)
		return
	}

	utils.SendResponse("success", http.StatusCreated, w)
}
