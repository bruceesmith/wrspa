/*
Package setup is the page component for handling game initialisation, specifically
the choice of game type (custom or random) and the selection of endpoints (start
and goal)
*/
package setup

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/bruceesmith/wrspa/go-app/backend/api"
	"github.com/bruceesmith/logger"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// ---------------------------------------------------------------------------
//
// Model
//
// ---------------------------------------------------------------------------

type gametype string

const (
	custom gametype = "custom"
	random gametype = "random"
	unset  gametype = "unset"
)

func (gt gametype) String() string {
	return string(gt)
}

type Setup struct {
	app.Compo
	Tipe           gametype
	selector       typeSelector
	customSelected customSelected
	randomSelected randomSelected
}

var (
	randomStart, randomGoal string
	Default                 Setup
)

func init() {
	Default = Setup{
		Tipe: unset,
	}
}

// ---------------------------------------------------------------------------
//
// View
//
// ---------------------------------------------------------------------------

func (s *Setup) Render() app.UI {
	logger.Trace("setup.Render called")
	components := []app.UI{
		app.P().Text("Wiki Racing").
			Class("gwr-title"),
	}
	components = append(
		components,
		s.selector.view()...,
	)
	if s.Tipe == custom {
		components = append(
			components,
			s.customSelected.view()...,
		)
	} else if s.Tipe == random {
		components = append(
			components,
			s.randomSelected.view()...,
		)
	}
	return app.Div().
		Body(components...).
		Class("gwr-selector")
}

// ---------------------------------------------------------------------------
//
// Controller
//
// ---------------------------------------------------------------------------

func (s *Setup) OnMount(ctx app.Context) {
	ctx.ObserveState("gameTypeSelected", &s.Tipe)
	ctx.Async(
		func() {
			resp, err := http.Get("/api/SpecialRandom")
			if err != nil {
				logger.Error("Setup.OnMount error fetching SpecialRandom", "error", err.Error())
				return
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				logger.Error("Setup.OnMount error reading SpecialRandom response", "error", err.Error())
				return
			}
			var response api.SpecialRandomResponse
			err = json.Unmarshal(body, &response)
			if err != nil {
				logger.Error("Setup.OnMount error unmarshaling SpecialRandom response", "error", err.Error())
				return
			}
			randomStart = response.Start
			randomGoal = response.Goal
			logger.Trace("setup.OnMount random points fetched")
		},
	)
}
