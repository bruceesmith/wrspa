package wrserver

import (
	"net/http"

	"github.com/bruceesmith/terminator"
)

type clientAdapter struct {
	client ClientInterface
}

func newClientAdapter() (c *clientAdapter) {
	return &clientAdapter{client: NewClient("")}
}

func (ca *clientAdapter) Get(path string) (body []byte, err error) {
	return ca.client.Get(path)
}

func (ca *clientAdapter) GetRandom() (path string) {
	return ca.client.GetRandom()
}

type serverAdapter struct {
	server ServerInterface
}

func newServerAdapter(port, static string, client ClientInterface) (s *serverAdapter, err error) {
	svr, err := NewServer(port, static, client)
	return &serverAdapter{server: svr}, err
}

func (sa *serverAdapter) API(w http.ResponseWriter, r *http.Request) {
	sa.server.API(w, r)
}

func (sa *serverAdapter) MarshalFailure(function string, err error, response any) string {
	return sa.server.MarshalFailure(function, err, response)
}

func (sa *serverAdapter) Serve(t *terminator.Terminator) {
	sa.server.Serve(t)
}

func (sa *serverAdapter) Settings(w http.ResponseWriter, r *http.Request) {
	sa.server.Settings(w, r)
}

func (sa *serverAdapter) SPAFile(w http.ResponseWriter, r *http.Request) {
	sa.server.SPAFile(w, r)
}

func (sa *serverAdapter) SpecialRandom(w http.ResponseWriter, r *http.Request) {
	sa.server.SpecialRandom(w, r)
}

func (sa *serverAdapter) WikiPage(w http.ResponseWriter, r *http.Request) {
	sa.server.WikiPage(w, r)
}

func (sa *serverAdapter) WikipediaFile(w http.ResponseWriter, r *http.Request) {
	sa.server.WikipediaFile(w, r)
}
