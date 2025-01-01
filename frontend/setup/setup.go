package setup

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/bruceesmith/echidna/logger"
	"github.com/bruceesmith/go-wikiracing/backend/api"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// ---------------------------------------------------------------------------
//
// Model
//
// ---------------------------------------------------------------------------

type state int

const (
	initial state = iota
	typeSelected
	queryingEndpoints
)

type Setup struct {
	app.Compo
	ctx   app.Context
	state state
	tipe  gametype
}

var (
	start, goal             string
	randomStart, randomGoal string
	Default                 Setup
)

func init() {
	Default = Setup{
		state: initial,
		tipe:  unset,
	}
}

// ---------------------------------------------------------------------------
//
// View
//
// ---------------------------------------------------------------------------

func (s *Setup) Render() (ui app.UI) {
	switch s.state {
	case initial:
		ui = app.Div().
			Body(
				app.Text("Wiki Racing"),
				&typeSelector{},
			)
	case typeSelected:
		ui = app.Div().
			Body(
				app.Text("Wiki Racing"),
				app.Br(),
				app.If(
					s.tipe == random,
					func() app.UI {
						return &randomSelected{}
					},
				).Else(
					func() app.UI {
						return &customSelected{}
					},
				),
			)
	default:
		ui = app.Text("unknown Setup state")
	}
	return
}

// ---------------------------------------------------------------------------
//
// Controller
//
// ---------------------------------------------------------------------------

func (s *Setup) OnMount(ctx app.Context) {
	s.ctx = ctx
	ctx.Handle("gameTypeSelected", s.typeSelected)
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
		},
	)
}

func (s *Setup) typeSelected(ctx app.Context, a app.Action) {
	s.state = typeSelected
	s.tipe, _ = a.Value.(gametype)
}
