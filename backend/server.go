/*
Package wrserver provides an HTTPS server for the WikiRacing game.
It serves the game web app and the API for the game.
It also serves static files from Wikipedia.org, such as PNG and SVG files.
*/
package wrserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/bruceesmith/logger"
	"github.com/bruceesmith/terminator"
	"golang.org/x/net/html"
)

// Server is the HTTP server for this program
type Server struct {
	client ClientInterface
	port   string
	root   string
	server *http.Server
}

// NewServer returns a Server
func NewServer(port, static string, client ClientInterface) (svr ServerInterface, err error) {
	p, err := strconv.ParseInt(port, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("unable to parse port %s: %w", port, err)
	}
	if p <= 1024 || p > 65535 {
		return nil, fmt.Errorf("invalid port %s", port)
	}

	info, err := os.Stat(static)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("static folder '%s' not found", static)
		}
		return nil, fmt.Errorf("error accessing static folder '%s': %w", static, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("static path '%s' is not a directory", static)
	}

	s := &Server{
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

// handleError is a helper function to handle errors in a consistent way
func (s *Server) handleError(w http.ResponseWriter, function string, err error, statusCode int, details any) {
	logger.Error(function+" failure", "error", err.Error())
	w.WriteHeader(statusCode)
	w.Write([]byte(s.MarshalFailure(function, err, details)))
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
	w.Header().Set("Content-Type", "application/json")
	// Package up a JSON response
	// Get the current log level and trace IDs
	response := SettingsResponse{
		LogLevel: logger.Level(),
		TraceIDs: logger.TraceIDs(),
	}
	jason, err := json.Marshal(response)
	if err != nil {
		s.handleError(w, "settings", err, http.StatusInternalServerError, response)
		return
	}
	w.Write(jason)
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
		s.handleError(w, "specialrandom", err, http.StatusInternalServerError, response)
		return
	}
	w.Write(jason)
}

// WikiPage is the handler for the /api/wikipage REST endpoint
func (s *Server) WikiPage(w http.ResponseWriter, r *http.Request) {
	// Extract the subject from the POST requst
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.handleError(w, "wikipage", err, http.StatusInternalServerError, string(body))
		return
	}
	var request WikiPageRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		s.handleError(w, "wikipage", err, http.StatusBadRequest, string(body))
		return
	}
	// All wiki page subjects must being with "/wiki/"
	if !strings.HasPrefix(request.Subject, "/wiki/") {
		err := fmt.Errorf("invalid subject: %s", request.Subject)
		s.handleError(w, "wikipage", err, http.StatusBadRequest, request.Subject)
		return
	}

	// Fetch the wiki page for the requested aubject
	pg, err := s.client.Get(request.Subject)
	if err != nil {
		s.handleError(w, "wikipage", err, http.StatusNotFound, request.Subject)
		return
	}

	// Extract the page body
	page, err := s.extractBody(pg)
	if err != nil {
		s.handleError(w, "wikipage", err, http.StatusInternalServerError, request.Subject)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(page)
}

// extractBody extracts the HTML from the <body> of a page
func (s *Server) extractBody(page []byte) ([]byte, error) {
	// Quick check for the presence of a body tag. This is not foolproof but
	// catches simple cases where a page is not a full HTML document.
	if !bytes.Contains(bytes.ToLower(page), []byte("<body")) {
		return nil, fmt.Errorf("no <body> tag found in html")
	}

	doc, err := html.Parse(bytes.NewReader(page))
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %w", err)
	}

	var body *html.Node
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "body" {
			body = n
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
			if body != nil {
				return
			}
		}
	}
	f(doc)

	if body == nil {
		// This is unlikely to be reached because html.Parse adds a body tag,
		// but we'll keep it for safety.
		return nil, fmt.Errorf("no <body> tag found in html (after parsing)")
	}

	var buf bytes.Buffer
	for c := body.FirstChild; c != nil; c = c.NextSibling {
		err := html.Render(&buf, c)
		if err != nil {
			return nil, fmt.Errorf("failed to render html: %w", err)
		}
	}
	return buf.Bytes(), nil
}

// WikipediaFile serves static files from Wikipedia
func (s *Server) WikipediaFile(w http.ResponseWriter, r *http.Request) {
	body, err := s.client.Get(r.URL.Path)
	if err != nil {
		s.handleError(w, "static", err, http.StatusNotFound, r.URL.Path)
		return
	}
	w.Write(body)
}
