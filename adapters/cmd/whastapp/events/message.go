package events

import (
	"go.mau.fi/whatsmeow/types/events"
)

func MessageHandler(evt *events.Message) {
	// message := models.Message{
	// 	ID:        evt.Info.ID,
	// 	Text:      evt.Message.GetConversation(),
	// 	From:      fmt.Sprintf("%s@%s", evt.Info.Sender.User, evt.Info.Sender.Server),
	// 	Recipient: evt.Info.Chat.String(),
	// 	Timestamp: evt.Info.Timestamp,
	// 	Platform:  constants.Whatsapp.Name,
	// }

	// platform := whastapp.NewWhatsappPlatform()
	// platform.SetEvent(*evt)

	// handlers.MessageHandler(platform, message)
}
