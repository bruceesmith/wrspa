/*
Package wrserver provides an HTTPS server for the WikiRacing game.
It serves the game web app and the API for the game.
It also serves static files from Wikipedia.org, such as PNG and SVG files.
*/
package wrserver

import (
	"context"
	"net/http"

	"github.com/bruceesmith/logger"
	"github.com/bruceesmith/terminator"
)

// Server is the HTTP server for this program
type Server struct {
	server *http.Server
	port   string
}

// NewServer returns a Server
func NewServer(port string, static string) (s *Server, err error) {
	s = &Server{
		server: &http.Server{
			Addr: ":" + port,
		},
		port: port,
	}
	mux := http.NewServeMux()
	mux.Handle("/api/", NewAPIHandler())
	mux.Handle("/static/", staticHandler{})
	mux.Handle("/w/", staticHandler{})
	mux.Handle("/", fileHandler{root: static})
	s.server.Handler = s.multiHandler(mux)
	return
}

// multiHandler allows go-app, REST HTTP(S) and Wikipediea static file calls to co-exist
func (s *Server) multiHandler(mux http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.TraceID("server", "multiHandler", "path", r.URL.Path)
		mux.ServeHTTP(w, r)
	})
}

// Serve handles all HTTP(S) requests
func (s *Server) Serve() {
	terminator.Add(1)
	go func() {
		<-terminator.ShutDown()
		s.server.Shutdown(context.Background())
		terminator.Done()
	}()

	// Start the server
	logger.Debug("API server starting on :" + s.port)
	if err := s.server.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			logger.Error("Server.Serve() ListenAndServe error", "error", err.Error())
		}
	}
	logger.Debug("API server exiting")
}
