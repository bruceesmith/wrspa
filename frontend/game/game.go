/*
Package game is the top level component in Go-WikiRacing. Its controller determines which
page is displayed depending on whether game endpoints (start and goal) have yet been set
*/
package game

import (
	"github.com/bruceesmith/go-wikiracing/frontend/observables"
	"github.com/bruceesmith/go-wikiracing/frontend/setup"
	"github.com/bruceesmith/go-wikiracing/frontend/wiki"
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
}
