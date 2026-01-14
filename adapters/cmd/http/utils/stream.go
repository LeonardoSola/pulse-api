package utils

import (
	"fmt"
	"net/http"
)

func StreamSSE(
	w http.ResponseWriter,
	r *http.Request,
	messages <-chan string,
) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		return
	}

	ctx := r.Context()

	for {
		select {
		case <-ctx.Done():
			return

		case msg, ok := <-messages:
			if !ok {
				return
			}

			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		}
	}

}
