package server

import (
	"log/slog"
	"net/http"
)

type staticHandler struct{}

// ServeHTTP is the request handler for PNG and SVG files from wikipedia.org
func (s staticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := get(r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if _, err = w.Write(body); err != nil {
		slog.Error("error on static file response Write", "error", err)
	}
}
