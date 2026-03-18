package setup

import (
	"github.com/bruceesmith/go-wikiracing/frontend/observables"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// ---------------------------------------------------------------------------
//
// Model
//
// ---------------------------------------------------------------------------

type randomSelected struct{}

// ---------------------------------------------------------------------------
//
// View
//
// ---------------------------------------------------------------------------

func (r *randomSelected) view() []app.UI {
	return []app.UI{
		app.Text("The randomly selected endpoints are:"),
		app.Br(),
		app.Text("Start: " + randomStart),
		app.Br(),
		app.Text("Goal: " + randomGoal),
		app.Br(),
		app.Button().Text("Next").
			OnClick(r.next).
			Class("gwr-custom-next-step"),
	}
}

// ---------------------------------------------------------------------------
//
// Controller
//
// ---------------------------------------------------------------------------

func (r *randomSelected) next(ctx app.Context, e app.Event) {
	tags := app.Tags{}
	tags.Set("start", randomStart)
	tags.Set("goal", randomGoal)
	ctx.SetState(observables.GameSelected, tags)

}
