package wrserver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/bruceesmith/logger"
)

// apiHandler handles REST requests to the various /api/ endpoints
type apiHandler struct {
	bodyRe *regexp.Regexp
}

const (
	expectedMatches = 3                              // Expect 3 matches from FindSubmatch
	bodyRegex       = "(?ms)<body (.+?)>(.+)</body>" // Regex to extract the Wiki page's body
)

func NewAPIHandler() apiHandler {
	ah := apiHandler{
		bodyRe: regexp.MustCompile(bodyRegex),
	}
	return ah
}

// ServeHTTP is the request handler
func (a apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	function := EndPoint(strings.ToLower(strings.TrimPrefix(r.URL.Path, "/api/")))
	switch {
	case r.Method == http.MethodGet && function == Settings:
		a.Settings(w, r)
		return
	case r.Method == http.MethodGet && function == SpecialRandom:
		a.SpecialRandom(w, r)
		return
	case r.Method == http.MethodPost && function == WikiPage:
		a.WikiPage(w, r)
		return
	}
}

// Settings is the handler for the /api/settings REST endpoint
func (a apiHandler) Settings(w http.ResponseWriter, r *http.Request) {
	// Package up a JSON response
	// Get the current log level and trace IDs
	response := SettingsResponse{
		LogLevel: logger.Level(),
		TraceIDs: logger.TraceIDs(),
	}
	jason, err := json.Marshal(response)
	if err != nil {
		w.Write([]byte(marshalFailure("settings", err, response)))
	} else {
		w.Write(jason)
	}
}

// SpecialRandom is the handler for the /api/specialrandom REST endpoint
func (a apiHandler) SpecialRandom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := SpecialRandomResponse{
		Start: getRandom(),
		Goal:  getRandom(),
	}
	jason, err := json.Marshal(response)
	if err != nil {
		w.Write([]byte(marshalFailure("specialrandom", err, response)))
	} else {
		w.Write(jason)
	}
}

// WikiPage is the handler for the /api/wikipage REST endpoint
func (a apiHandler) WikiPage(w http.ResponseWriter, r *http.Request) {
	// Extract the subject from the POST requst
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("wikipage request failure", "error", err.Error())
		w.Write([]byte(marshalFailure("wikipage", err, body)))
		return
	}
	var request WikiPageRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		logger.Error("wikipage request failure", "error", err.Error())
		w.Write([]byte(marshalFailure("wikipage", err, body)))
		return
	}
	// Fetch the wiki page for the requested aubject
	page, err := getString(request.Subject)
	if err != nil {
		logger.Error("wikipage fetch failure", "error", err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// Split out the page body
	matches := a.bodyRe.FindSubmatch([]byte(page))
	if len(matches) != expectedMatches {
		logger.Error("wikipage fetch failure", "error", fmt.Sprint("expected ", expectedMatches, " regex matches, got ", len(matches)))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	page = "<div " + string(matches[1]) + ">" + string(matches[2]) + "</div"
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(page))
}
