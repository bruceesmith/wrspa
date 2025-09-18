package wrserver

import (
	"net/http"

	"github.com/bruceesmith/terminator"
)

// ClientInterface is an interface for the Client struct
type ClientInterface interface {
	Get(path string) (body []byte, contentType string, err error)
	GetRandom() (path string)
}

// ServerInterface is an interface for the Server struct
type ServerInterface interface {
	API(w http.ResponseWriter, r *http.Request)
	MarshalFailure(function string, err error, response any) string
	Serve(t *terminator.Terminator)
	Settings(w http.ResponseWriter, r *http.Request)
	SPAFile(w http.ResponseWriter, r *http.Request)
	SpecialRandom(w http.ResponseWriter, r *http.Request)
	WikiPage(w http.ResponseWriter, r *http.Request)
	WikipediaFile(w http.ResponseWriter, r *http.Request)
}
