package wrserver

import (
	"net/http"

	"github.com/bruceesmith/logger"
)

type fileHandler struct {
	root string
}

// ServeHTTP is the request handler for static local files, e.g. index.html and *.mjs
func (f fileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	file := f.root
	if r.URL.Path == "/" {
		file += "/index.html"
	} else {
		file += r.URL.Path
	}
	logger.TraceID("server", "fileHandler", "URL", r.URL.String(), "file", file)
	http.ServeFile(w, r, file)
}
