/*
Package wrserver provides an HTTPS server for the WikiRacing game.
It serves the game web app and the API for the game.
It also serves static files from Wikipedia.org, such as PNG and SVG files.
*/
package wrserver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/bruceesmith/logger"
	"github.com/bruceesmith/terminator"
)

// Server is the HTTP server for this program
type Server struct {
	bodyRe *regexp.Regexp
	client ClientInterface
	port   string
	root   string
	server *http.Server
}

const (
	expectedMatches = 3                              // Expect 3 matches from FindSubmatch
	bodyRegex       = "(?ms)<body (.+?)>(.+)</body>" // Regex to extract the Wiki page's body
)

// NewServer returns a Server
func NewServer(port, static string, client ClientInterface) (svr ServerInterface, err error) {
	p, err := strconv.ParseInt(port, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("unable to parse port %s: %w", port, err)
	}
	if p <= 1024 || p > 65535 {
		return nil, fmt.Errorf("invalid port %s", port)
	}
	s := &Server{
		bodyRe: regexp.MustCompile(bodyRegex),
		client: client,
		port:   port,
		root:   static,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.SPAFile)
	mux.HandleFunc("/api/", s.API)
	mux.HandleFunc("/static/", s.WikipediaFile)
	mux.HandleFunc("/w/", s.WikipediaFile)

	s.server = &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	svr = s
	// mux := http.NewServeMux()
	// mux.Handle("/api/", NewAPIHandler())
	// mux.Handle("/static/", staticHandler{})
	// mux.Handle("/w/", staticHandler{})
	// mux.Handle("/", fileHandler{root: static})
	// svr.server.Handler = svr.MultiHandler(mux)
	return
}

// API provides the REST interface for the SPA
func (s *Server) API(w http.ResponseWriter, r *http.Request) {
	function := EndPoint(strings.ToLower(strings.TrimPrefix(r.URL.Path, "/api/")))
	switch {
	case r.Method == http.MethodGet && function == Settings:
		s.Settings(w, r)
		return
	case r.Method == http.MethodGet && function == SpecialRandom:
		s.SpecialRandom(w, r)
		return
	case r.Method == http.MethodPost && function == WikiPage:
		s.WikiPage(w, r)
		return
	}
}

// MarshalFailure creates a sensible JSON-format error message
func (s *Server) MarshalFailure(function string, err error, response any) string {
	return `{"msg": "unable to marshal API response", ` +
		`"function": "` + function + `", ` +
		`"error": "` + err.Error() + `", ` +
		`"response": "` + fmt.Sprintf("%+v", response) + `"}`
}

// Serve handles all HTTP(S) requests
func (s *Server) Serve(t *terminator.Terminator) {
	t.Add(1)
	go func() {
		<-t.ShutDown()
		s.server.Shutdown(context.Background())
		t.Done()
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

// Settings is the handler for the /api/settings REST endpoint
func (s *Server) Settings(w http.ResponseWriter, r *http.Request) {
	// Package up a JSON response
	// Get the current log level and trace IDs
	response := SettingsResponse{
		LogLevel: logger.Level(),
		TraceIDs: logger.TraceIDs(),
	}
	jason, err := json.Marshal(response)
	if err != nil {
		w.Write([]byte(s.MarshalFailure("settings", err, response)))
	} else {
		w.Write(jason)
	}
}

// SPAFile serves static files for the SPA (index.html, JavaScript, CSS, etc.)
func (s *Server) SPAFile(w http.ResponseWriter, r *http.Request) {
	file := s.root
	if r.URL.Path == "/" {
		file += "/index.html"
	} else {
		file += r.URL.Path
	}
	logger.TraceID("server", "fileHandler", "URL", r.URL.String(), "file", file)
	http.ServeFile(w, r, file)
}

// SpecialRandom is the handler for the /api/specialrandom REST endpoint
func (s *Server) SpecialRandom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := SpecialRandomResponse{
		Start: s.client.GetRandom(),
		Goal:  s.client.GetRandom(),
	}
	jason, err := json.Marshal(response)
	if err != nil {
		w.Write([]byte(s.MarshalFailure("specialrandom", err, response)))
	} else {
		w.Write(jason)
	}
}

// WikiPage is the handler for the /api/wikipage REST endpoint
func (s *Server) WikiPage(w http.ResponseWriter, r *http.Request) {
	// Extract the subject from the POST requst
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("wikipage request failure", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(s.MarshalFailure("wikipage", err, body)))
		return
	}
	var request WikiPageRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		logger.Error("wikipage request failure", "error", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(s.MarshalFailure("wikipage", err, body)))
		return
	}
	// Fetch the wiki page for the requested aubject
	pg, err := s.client.Get("/" + request.Subject)
	if err != nil {
		logger.Error("wikipage fetch failure", "error", err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// Split out the page body
	matches := s.bodyRe.FindSubmatch(pg)
	if len(matches) != expectedMatches {
		logger.Error("wikipage regex failure", "error", fmt.Sprint("expected ", expectedMatches, " regex matches, got ", len(matches)))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	page := "<div " + string(matches[1]) + ">" + string(matches[2]) + "</div"
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(page))
}

// WikipediaFile serves static files from Wikipedia
func (s *Server) WikipediaFile(w http.ResponseWriter, r *http.Request) {
	body, err := s.client.Get(r.URL.Path)
	if err != nil {
		logger.Error("static request failure", "error", err.Error(), "path", r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Write(body)
}
