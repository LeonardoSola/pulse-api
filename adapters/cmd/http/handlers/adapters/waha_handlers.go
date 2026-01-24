package adapters

import (
	"net/http"

	"github.com/LeonardoSola/wagon-api/adapters/bootstrap"
	"github.com/LeonardoSola/wagon-api/adapters/cmd/http/middleware"
	"github.com/LeonardoSola/wagon-api/adapters/cmd/http/utils"
	message_entity "github.com/LeonardoSola/wagon-api/internal/entities/message"
	"github.com/go-chi/chi/v5"
)

func initWaha(router chi.Router) {
	router.Use(middleware.AuthAPIKey)
	router.Post("/api/sendText", sendText)
}

type sendTextPayload struct {
	Text      string `json:"text"`
	DeviceID  string `json:"session"`
	Recipient string `json:"chatId"`
}

func sendText(w http.ResponseWriter, r *http.Request) {
	var payload sendTextPayload
	err := utils.DecodeBody(r, &payload)
	if err != nil {
		utils.ErrorResponse(err, w)
		return
	}

	message := message_entity.NewMessage(payload.DeviceID, payload.Text, payload.Recipient, message_entity.RecipientTypeUser, "")
	if err != nil {
		utils.ErrorResponse(err, w)
		return
	}

	err = message.SetRecipient(payload.Recipient)
	if err != nil {
		utils.ErrorResponse(err, w)
		return
	}

	err = bootstrap.Services.MessageService.CreateMessage(
		r.Context(),
		message,
	)
	if err != nil {
		utils.ErrorResponse(err, w)
		return
	}

	utils.SendResponse("success", http.StatusOK, w)
}
