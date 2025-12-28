package message_handler

type CreateMessageDTO struct {
	Recipient string `json:"recipient"`
	Text      string `json:"text"`
}
