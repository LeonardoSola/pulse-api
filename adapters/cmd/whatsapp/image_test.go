package whatsapp

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// pngHeader is enough for http.DetectContentType to recognise a PNG.
var pngHeader = []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR")

func TestInferMimeTypeFromURL(t *testing.T) {
	tests := map[string]string{
		"https://x.test/a.jpg":            "image/jpeg",
		"https://x.test/a.JPEG":           "image/jpeg",
		"https://x.test/a.png?size=large": "image/png",
		"https://x.test/a.gif#frag":       "image/gif",
		"https://x.test/a.webp":           "image/webp",
		"https://x.test/a.txt":            "",
	}

	for url, want := range tests {
		if got := inferMimeTypeFromURL(url); got != want {
			t.Errorf("inferMimeTypeFromURL(%q) = %q, want %q", url, got, want)
		}
	}
}

func TestDetectImageMimeType(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		url         string
		data        []byte
		want        string
		wantErr     bool
	}{
		{name: "header wins", contentType: "image/webp", url: "https://x.test/a", data: pngHeader, want: "image/webp"},
		{name: "header parameters stripped", contentType: "image/png; charset=binary", data: pngHeader, want: "image/png"},
		{name: "sniffed when header is generic", contentType: "application/octet-stream", data: pngHeader, want: "image/png"},
		{name: "url extension as last resort", contentType: "application/octet-stream", url: "https://x.test/a.gif", data: []byte("??"), want: "image/gif"},
		{name: "html page rejected", contentType: "text/html", url: "https://x.test/page", data: []byte("<html></html>"), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := detectImageMimeType(tt.contentType, tt.url, tt.data)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected an error, got %q", got)
				}
				return
			}
			if err != nil || got != tt.want {
				t.Fatalf("got (%q, %v), want %q", got, err, tt.want)
			}
		})
	}
}

func TestDownloadImage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/ok.png":
			w.Header().Set("Content-Type", "image/png")
			_, _ = w.Write(pngHeader)
		case "/too-big":
			w.Header().Set("Content-Type", "image/jpeg")
			_, _ = w.Write(bytes.Repeat([]byte{0}, maxImageSize+1))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	data, mimeType, err := downloadImage(context.Background(), server.URL+"/ok.png")
	if err != nil {
		t.Fatalf("download: %v", err)
	}
	if mimeType != "image/png" || !bytes.Equal(data, pngHeader) {
		t.Errorf("got %q with %d bytes", mimeType, len(data))
	}

	if _, _, err := downloadImage(context.Background(), server.URL+"/missing"); err == nil {
		t.Error("expected an error for a 404")
	}

	if _, _, err := downloadImage(context.Background(), server.URL+"/too-big"); err == nil {
		t.Error("expected an error for an image over the size limit")
	}
}
