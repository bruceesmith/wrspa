package wrserver

import (
	"net/http"

	"github.com/bruceesmith/logger"
)

type staticHandler struct{}

// ServeHTTP is the request handler for PNG and SVG files from wikipedia.org
func (s staticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := get(r.URL.Path)
	if err != nil {
		logger.Error("static request failure", "error", err.Error(), "path", r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Write(body)
}
