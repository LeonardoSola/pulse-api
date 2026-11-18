package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LeonardoSola/pulse-api/internal/custom_error"
)

func decodeError(t *testing.T, rec *httptest.ResponseRecorder) Error {
	t.Helper()

	var body Error
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("invalid JSON body %q: %v", rec.Body.String(), err)
	}
	return body
}

func TestErrorResponseCustomError(t *testing.T) {
	rec := httptest.NewRecorder()
	notFound := custom_error.CustomError{Status: http.StatusNotFound, Message: "Device not found"}

	ErrorResponse(notFound, rec)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q", ct)
	}
	if body := decodeError(t, rec); body.Message != "Device not found" || body.Status != http.StatusNotFound {
		t.Errorf("body = %+v", body)
	}
}

func TestErrorResponseHidesInternalErrors(t *testing.T) {
	rec := httptest.NewRecorder()

	// config.Config is nil in tests, which is treated as non-development.
	ErrorResponse(errors.New("pq: password authentication failed"), rec)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d", rec.Code)
	}
	body := decodeError(t, rec)
	if body.Message != "Internal server error" || body.Info != "" {
		t.Errorf("internal error details leaked: %+v", body)
	}
}

func TestSendResponse(t *testing.T) {
	rec := httptest.NewRecorder()

	SendResponse(map[string]string{"status": "ok"}, http.StatusCreated, rec)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d", rec.Code)
	}
	if got := rec.Body.String(); got != `{"data":{"status":"ok"}}` {
		t.Errorf("body = %s", got)
	}
}

func TestStreamSSE(t *testing.T) {
	messages := make(chan string, 2)
	messages <- "code-1"
	messages <- "success"
	close(messages)

	rec := httptest.NewRecorder()
	StreamSSE(rec, httptest.NewRequest(http.MethodGet, "/devices/x/qr", nil), messages)

	if ct := rec.Header().Get("Content-Type"); ct != "text/event-stream" {
		t.Errorf("Content-Type = %q", ct)
	}
	if got, want := rec.Body.String(), "data: code-1\n\ndata: success\n\n"; got != want {
		t.Errorf("body = %q, want %q", got, want)
	}
}
