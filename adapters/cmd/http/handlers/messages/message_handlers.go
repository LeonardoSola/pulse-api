package message_handler

import (
	"net/http"
	"strings"

	"github.com/LeonardoSola/pulse-api/adapters/bootstrap"
	"github.com/LeonardoSola/pulse-api/adapters/cmd/http/utils"
	message_entity "github.com/LeonardoSola/pulse-api/internal/entities/message"
	"github.com/go-chi/chi/v5"
)

func InitMessages(router chi.Router) {
	router.Post("/devices/{deviceId}/messages/sendText", createTextMessage)
	router.Post("/devices/{deviceId}/messages/sendImage", createImageMessage)
}

func createTextMessage(w http.ResponseWriter, r *http.Request) {
	sendMessage(w, r, false)
}

func createImageMessage(w http.ResponseWriter, r *http.Request) {
	sendMessage(w, r, true)
}

func sendMessage(w http.ResponseWriter, r *http.Request, requireImage bool) {
	var payload CreateMessageTextDTO
	if err := utils.DecodeBody(r, &payload); err != nil {
		utils.ErrorResponse(err, w)
		return
	}

	message, err := payload.toMessage(chi.URLParam(r, "deviceId"))
	if err != nil {
		utils.ErrorResponse(err, w)
		return
	}

	if requireImage && message.ImageURL == "" {
		utils.ErrorResponse(message_entity.ErrImageURLRequired, w)
		return
	}

	if err := bootstrap.Services.MessageService.CreateMessage(r.Context(), message); err != nil {
		utils.ErrorResponse(err, w)
		return
	}

	utils.SendResponse("success", http.StatusCreated, w)
}

// toMessage builds the outgoing message. The recipient may be a bare number /
// group ID (with recipient_type) or a full JID such as "5511999999999@s.whatsapp.net".
func (p CreateMessageTextDTO) toMessage(deviceID string) (message_entity.Message, error) {
	message := message_entity.NewMessageWithImage(deviceID, p.Text, p.Recipient, p.RecipientType, p.ImageURL)

	if strings.Contains(p.Recipient, "@") {
		if err := message.SetRecipient(p.Recipient); err != nil {
			return message_entity.Message{}, err
		}
	}

	return message, nil
}
