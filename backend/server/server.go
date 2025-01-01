package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/bruceesmith/echidna/logger"
	"github.com/bruceesmith/echidna/terminator"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// Server is the HTTP server for this program
type Server struct {
	server *http.Server
}

// New returns a Server
func New() (s *Server, err error) {
	s = &Server{
		server: &http.Server{
			Addr: ":8000",
		},
	}
	mux := http.NewServeMux()
	mux.Handle("/api/", apiHandler{})
	s.server.Handler = s.apiHandlerFunc(mux)
	return
}

// apiHandlerFunc allows both go-app and REST HTTP(S) calls to co-exist
func (s *Server) apiHandlerFunc(restHandler http.Handler) http.Handler {
	h := &app.Handler{
		Name:        "WikiRacing",
		Description: "A wiki racing game",
		Styles: []string{
			"/web/game.css",
		},
		RawHeaders: []string{
			`<script>
				function wikiAnchorClick(a) {
					console.log("wikiAnchorClick called");
  					a.preventDefault();
					a.stopImmediatePropagation();
					wikiUrlClicked(a.target.href);
				};
				history.pushState("minus-2","");
        		history.pushState("minus-1","");
			</script>
			<link rel="stylesheet" href="https://en.wikipedia.com/w/load.php?lang=en&amp;modules=ext.cite.styles%7Cext.tmh.player.styles%7Cext.uls.interlanguage%7Cext.visualEditor.desktopArticleTarget.noscript%7Cext.wikimediaBadges%7Cext.wikimediamessages.styles%7Cjquery.makeCollapsible.styles%7Cmediawiki.page.gallery.styles%7Cskins.vector.icons%2Cstyles%7Cskins.vector.search.codex.styles%7Cwikibase.client.init&amp;only=styles&amp;skin=vector-2022">
			<link rel="stylesheet" href="https://en.wikipedia.com/w/load.php?lang=en&amp;modules=site.styles&amp;only=styles&amp;skin=vector-2022">`,
		},
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			restHandler.ServeHTTP(w, r)
		} else {
			h.ServeHTTP(w, r)
		}
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
	logger.Debug("API server starting on :8000")
	if err := s.server.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			logger.Error("Server.Serve() ListenAndServe error", "error", err.Error())
		}
	}
	logger.Debug("API server exiting")
}
