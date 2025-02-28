/*
Package server provides the static file handler for the WASM side of Go-WikiRacing plus
the REST API handler for the WASM game
*/
package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/bruceesmith/logger"
	"github.com/bruceesmith/terminator"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

const (
	font = `<link rel="stylesheet"
  href="https://fonts.googleapis.com/css?family=Roboto:i,300,400,500,700&display=swap" />`
	materialHeader = `<script type="importmap">
	{
		"imports": {
			"@material/web/": "https://esm.run/@material/web/"
		}
	}
</script>
<script type="module">
	import '@material/web/all.js';
	import {styles as typescaleStyles} from '@material/web/typography/md-typescale-styles.js';
	document.adoptedStyleSheets.push(typescaleStyles.styleSheet);
</script>`

	wikiAnchorClick = `<script>
	function wikiAnchorClick(a) {
		console.log("wikiAnchorClick called");
  		a.preventDefault();
		a.stopImmediatePropagation();
		wikiUrlClicked(a.target.href);
	};
	history.pushState("minus-2","");
	history.pushState("minus-1","");
</script>`

	wikiLinks = `<link rel="stylesheet" href="https://en.wikipedia.com/w/load.php?lang=en&amp;modules=ext.cite.styles%7Cext.tmh.player.styles%7Cext.uls.interlanguage%7Cext.visualEditor.desktopArticleTarget.noscript%7Cext.wikimediaBadges%7Cext.wikimediamessages.styles%7Cjquery.makeCollapsible.styles%7Cmediawiki.page.gallery.styles%7Cskins.vector.icons%2Cstyles%7Cskins.vector.search.codex.styles%7Cwikibase.client.init&amp;only=styles&amp;skin=vector-2022">
<link rel="stylesheet" href="https://en.wikipedia.com/w/load.php?lang=en&amp;modules=site.styles&amp;only=styles&amp;skin=vector-2022">`
)

// Server is the HTTP server for this program
type Server struct {
	server *http.Server
	port   string
}

// New returns a Server
func New(port string) (s *Server, err error) {
	s = &Server{
		server: &http.Server{
			Addr: ":" + port,
		},
		port: port,
	}
	mux := http.NewServeMux()
	mux.Handle("/api/", apiHandler{})
	mux.Handle("/static/", staticHandler{})
	s.server.Handler = s.multiHandler(mux)
	return
}

// multiHandler allows go-app, REST HTTP(S) and Wikipediea static file calls to co-exist
func (s *Server) multiHandler(mux http.Handler) http.Handler {
	h := &app.Handler{
		Name:        "WikiRacing",
		Description: "A wiki racing game",
		Styles: []string{
			"/web/game.css",
		},
		RawHeaders: []string{
			font,
			wikiAnchorClick,
			wikiLinks,
			materialHeader,
		},
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.TraceID("server", "multiHandler", "path", r.URL.Path)
		if strings.HasPrefix(r.URL.Path, "/api/") || strings.HasPrefix(r.URL.Path, "/static/") {
			mux.ServeHTTP(w, r)
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
	logger.Debug("API server starting on :" + s.port)
	if err := s.server.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			logger.Error("Server.Serve() ListenAndServe error", "error", err.Error())
		}
	}
	logger.Debug("API server exiting")
}
