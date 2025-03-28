package server

import "net/http"

type staticHandler struct{}

// ServeHTTP is the request handler for PNG and SVG files from wikipedia.org
func (s staticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := get(r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Write(body)
}
