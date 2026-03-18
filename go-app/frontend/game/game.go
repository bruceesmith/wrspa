/*
Package game is the top level component in Go-WikiRacing. Its controller determines which
page is displayed depending on whether game endpoints (start and goal) have yet been set
*/
package game

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/bruceesmith/wrspa/go-app/backend/api"
	"github.com/bruceesmith/wrspa/go-app/frontend/observables"
	"github.com/bruceesmith/wrspa/go-app/frontend/setup"
	"github.com/bruceesmith/wrspa/go-app/frontend/wiki"
	"github.com/bruceesmith/logger"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// ---------------------------------------------------------------------------
//
// Model
//
// ---------------------------------------------------------------------------

type Game struct {
	app.Compo
	EndPoints app.Tags
}

func New() (g *Game) {
	g = &Game{
		EndPoints: make(app.Tags),
	}
	return
}

// ---------------------------------------------------------------------------
//
// View
//
// ---------------------------------------------------------------------------

func (g *Game) Render() app.UI {
	var ui app.HTMLDiv
	switch {
	case g.EndPoints.Get("start") == "":
		ui = app.Div().
			Body(
				&setup.Default,
			)
	default:
		wiki.Default.Targets(g.EndPoints.Get("start"), g.EndPoints.Get("goal"))
		ui = app.Div().
			Body(
				&wiki.Default,
			)
	}
	return ui.Class("gwr-game")
}

// ---------------------------------------------------------------------------
//
// Controller
//
// ---------------------------------------------------------------------------

func (g *Game) OnMount(ctx app.Context) {
	ctx.ObserveState(observables.GameSelected, &g.EndPoints)
	// Fetch and apply some server settings
	g.settings()
}

// settings fetches certain settings from the server and applies them, This quirky process
// is necssary because the WASM program does not receive the command-line flags passed to
// the server.
func (g *Game) settings() {
	var err error
	resp, err := http.Get("/api/settings")
	if err != nil {
		logger.Error("Game.OnMount error fetching server settings", "error", err.Error())
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Game.OnMount error reading server settings", "error", err.Error())
		return
	}
	settings := api.SettingsResponse{}
	err = json.Unmarshal(body, &settings)
	if err != nil {
		logger.Error("Game.OnMount error unmarshaling server settings", "error", err.Error())
		return
	}
	// Apply the settings
	var ll logger.LogLevel
	err = (&ll).Set(settings.LogLevel)
	if err != nil {
		logger.Error("Game.OnMount error setting log level", "error", err.Error())
		return
	}
	logger.SetLevel(slog.Level(ll))
	logger.SetTraceIds(settings.TraceIDs...)
}
