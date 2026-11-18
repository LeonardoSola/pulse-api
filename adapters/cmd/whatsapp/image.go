package whatsapp

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

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
		return "", message_entity.ErrImageDownload.SetError(err)
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

// maxImageSize is WhatsApp's limit for image messages.
const maxImageSize = 16 << 20 // 16 MiB

var imageHTTPClient = &http.Client{Timeout: 30 * time.Second}

func downloadImage(ctx context.Context, url string) ([]byte, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, "", err
	}

	resp, err := imageHTTPClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("failed to download image: status code %d", resp.StatusCode)
	}

	imageData, err := io.ReadAll(io.LimitReader(resp.Body, maxImageSize+1))
	if err != nil {
		return nil, "", err
	}

	if len(imageData) > maxImageSize {
		return nil, "", fmt.Errorf("image is larger than %d bytes", maxImageSize)
	}

	mimeType, err := detectImageMimeType(resp.Header.Get("Content-Type"), url, imageData)
	if err != nil {
		return nil, "", err
	}

	return imageData, mimeType, nil
}

// detectImageMimeType picks the image MIME type from the response header,
// the file contents or the URL extension, in that order.
func detectImageMimeType(contentType string, url string, data []byte) (string, error) {
	candidates := []string{contentType, http.DetectContentType(data), inferMimeTypeFromURL(url)}
	for _, candidate := range candidates {
		mimeType, _, _ := strings.Cut(candidate, ";")
		mimeType = strings.TrimSpace(strings.ToLower(mimeType))
		if strings.HasPrefix(mimeType, "image/") {
			return mimeType, nil
		}
	}

	return "", fmt.Errorf("URL does not point to an image (content type %q)", contentType)
}

func inferMimeTypeFromURL(url string) string {
	urlLower := strings.ToLower(url)
	if i := strings.IndexAny(urlLower, "?#"); i >= 0 {
		urlLower = urlLower[:i]
	}
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
