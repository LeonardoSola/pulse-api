package messages

import (
	"net/http"

	"github.com/LeonardoSola/wagon-api/adapters/bootstrap"
	"github.com/LeonardoSola/wagon-api/adapters/cmd/http/utils"
	message_entity "github.com/LeonardoSola/wagon-api/internal/entities/message"
	"github.com/go-chi/chi/v5"
)

func InitMessages(router *chi.Mux) {
	router.Post("/messages", createMessage)
}

func createMessage(w http.ResponseWriter, r *http.Request) {
	message := message_entity.Message{}
	err := utils.DecodeBody(r, &message)
	if err != nil {
		utils.ErrorResponse(err, w)
		return
	}

	err = bootstrap.Services.MessageService.CreateMessage(r.Context(), message)
	if err != nil {
		utils.ErrorResponse(err, w)
		return
	}

	utils.SendResponse(message, http.StatusCreated, w)
}
