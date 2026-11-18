package whatsapp

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	message_entity "github.com/LeonardoSola/pulse-api/internal/entities/message"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

func (s *WhatsappService) SendImageMessage(ctx context.Context, message message_entity.Message) (messageID string, err error) {
	if message.ImageURL == "" {
		return "", fmt.Errorf("image URL is required")
	}

	client, err := s.getClient(ctx, message.DeviceJID)
	if err != nil {
		return "", err
	}

	recipient, err := types.ParseJID(message.GetRecipient())
	if err != nil {
		return "", err
	}

	imageData, mimeType, err := downloadImage(ctx, message.ImageURL)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %w", err)
	}

	uploadResp, err := client.Upload(ctx, imageData, whatsmeow.MediaImage)
	if err != nil {
		return "", fmt.Errorf("failed to upload image: %w", err)
	}

	imageMsg := &waE2E.ImageMessage{
		Caption:       proto.String(message.Text),
		Mimetype:      proto.String(mimeType),
		URL:           &uploadResp.URL,
		DirectPath:    &uploadResp.DirectPath,
		MediaKey:      uploadResp.MediaKey,
		FileEncSHA256: uploadResp.FileEncSHA256,
		FileSHA256:    uploadResp.FileSHA256,
		FileLength:    &uploadResp.FileLength,
	}

	sentMessage, err := client.SendMessage(
		ctx, recipient, &waE2E.Message{
			ImageMessage: imageMsg,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to send image message: %w", err)
	}

	return sentMessage.ID, nil
}

func downloadImage(ctx context.Context, url string) ([]byte, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("failed to download image: status code %d", resp.StatusCode)
	}

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	mimeType := resp.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = inferMimeTypeFromURL(url)
	}

	if mimeType == "" {
		mimeType = "image/jpeg"
	}

	return imageData, mimeType, nil
}

func inferMimeTypeFromURL(url string) string {
	urlLower := strings.ToLower(url)
	switch {
	case strings.HasSuffix(urlLower, ".jpg") || strings.HasSuffix(urlLower, ".jpeg"):
		return "image/jpeg"
	case strings.HasSuffix(urlLower, ".png"):
		return "image/png"
	case strings.HasSuffix(urlLower, ".gif"):
		return "image/gif"
	case strings.HasSuffix(urlLower, ".webp"):
		return "image/webp"
	default:
		return ""
	}
}
