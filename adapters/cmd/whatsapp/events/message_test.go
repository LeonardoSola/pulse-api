package events

import (
	"testing"

	"go.mau.fi/whatsmeow/proto/waE2E"
	"google.golang.org/protobuf/proto"
)

func TestMessageText(t *testing.T) {
	tests := []struct {
		name string
		msg  *waE2E.Message
		want string
	}{
		{name: "nil message", msg: nil, want: ""},
		{name: "plain text", msg: &waE2E.Message{Conversation: proto.String("hello")}, want: "hello"},
		{name: "reply or link", msg: &waE2E.Message{ExtendedTextMessage: &waE2E.ExtendedTextMessage{Text: proto.String("see https://x.test")}}, want: "see https://x.test"},
		{name: "image caption", msg: &waE2E.Message{ImageMessage: &waE2E.ImageMessage{Caption: proto.String("photo")}}, want: "photo"},
		{name: "video caption", msg: &waE2E.Message{VideoMessage: &waE2E.VideoMessage{Caption: proto.String("clip")}}, want: "clip"},
		{name: "sticker has no text", msg: &waE2E.Message{StickerMessage: &waE2E.StickerMessage{}}, want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := messageText(tt.msg); got != tt.want {
				t.Errorf("messageText() = %q, want %q", got, tt.want)
			}
		})
	}
}
