package server

import (
	"net/http"

	"github.com/bruceesmith/logger"
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
		logger.Error("error on static file response Write", "error", err)
	}
}
