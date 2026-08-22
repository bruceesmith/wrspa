package server

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/bruceesmith/logger"
	"github.com/bruceesmith/wrspa/go-app/backend/api"
)

// apiHandler handles REST requests to the various /api/ endpoints
type apiHandler struct {
}

// ServeHTTP is the request handler
func (a apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	function := api.EndPoint(strings.ToLower(strings.TrimPrefix(r.URL.Path, "/api/")))
	switch {
	case r.Method == http.MethodGet && function == api.Settings:
		a.Settings(w, r)
		return
	case r.Method == http.MethodGet && function == api.SpecialRandom:
		a.SpecialRandom(w, r)
		return
	case r.Method == http.MethodPost && function == api.WikiPage:
		a.WikiPage(w, r)
		return
	}
}

// Settings is the handler for the /api/settings REST endpoint
func (a apiHandler) Settings(w http.ResponseWriter, r *http.Request) {
	// Package up a JSON response
	// Get the current log level and trace IDs
	response := api.SettingsResponse{
		LogLevel: logger.Level(),
		TraceIDs: logger.TraceIDs(),
	}
	jason, err := json.Marshal(response)
	if err != nil {
		if _, err = w.Write([]byte(marshalFailure("settings", err, response))); err != nil {
			logger.Error("error on Settings response Write", "error", err)
		}
	} else {
		if _, err = w.Write(jason); err != nil {
			logger.Error("error on Settings response Write", "error", err)
		}
	}
}

// SpecialRandom is the handler for the /api/specialrandom REST endpoint
func (a apiHandler) SpecialRandom(w http.ResponseWriter, r *http.Request) {
	response := api.SpecialRandomResponse{
		Start: getRandom(),
		Goal:  getRandom(),
	}
	jason, err := json.Marshal(response)
	if err != nil {
		if _, err = w.Write([]byte(marshalFailure("specialrandom", err, response))); err != nil {
			logger.Error("error on SpecialRandom response Write", "error", err)
		}
	} else {
		if _, err = w.Write(jason); err != nil {
			logger.Error("error on SpecialRandom response Write", "error", err)
		}
	}
}

// WikiPage is the handler for the /api/wikipage REST endpoint
func (a apiHandler) WikiPage(w http.ResponseWriter, r *http.Request) {
	// Extract the subject from the POST requst
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("wikipage request failure", "error", err.Error())
		if _, err = w.Write([]byte(marshalFailure("wikipage", err, body))); err != nil {
			logger.Error("error on WikiPage response Write", "error", err)
		}
		return
	}
	var request api.WikiPageRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		logger.Error("wikipage request failure", "error", err.Error())
		if _, err = w.Write([]byte(marshalFailure("wikipage", err, body))); err != nil {
			logger.Error("error on WikiPage response Write", "error", err)
		}
		return
	}
	// Fetch the wiki page for the requested aubject
	page, err := getString(request.Subject)
	if err != nil {
		logger.Error("wikipage fetch failure", "error", err.Error())
	}
	// Package up a JSON response
	response := api.WikiPageResponse{
		Page: page,
	}
	if err != nil {
		response.Error = err.Error()
	}
	jason, err := json.Marshal(response)
	if err != nil {
		if _, err = w.Write([]byte(marshalFailure("wikipage", err, response))); err != nil {
			logger.Error("error on WikiPage response Write", "error", err)
		}
	} else {
		if _, err = w.Write(jason); err != nil {
			logger.Error("error on WikiPage response Write", "error", err)
		}
	}
}
