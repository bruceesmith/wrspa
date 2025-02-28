package server

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/bruceesmith/go-wikiracing/backend/api"
	"github.com/bruceesmith/logger"
)

// apiHandler handles REST requests to the various /api/ endpoints
type apiHandler struct {
}

// ServeHTTP is the request handler
func (a apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	function := api.EndPoint(strings.ToLower(strings.TrimPrefix(r.URL.Path, "/api/")))
	switch {
	case r.Method == http.MethodGet && function == api.SpecialRandom:
		a.SpecialRandom(w, r)
		return
	case r.Method == http.MethodPost && function == api.WikiPage:
		a.WikiPage(w, r)
		return
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
	var request api.WikiPageRequest
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
		w.Write([]byte(marshalFailure("wikipage", err, response)))
	} else {
		w.Write(jason)
	}
}
