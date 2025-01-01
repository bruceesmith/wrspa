package setup

import "github.com/maxence-charriere/go-app/v10/pkg/app"

// ---------------------------------------------------------------------------
//
// Model
//
// ---------------------------------------------------------------------------

type randomSelected struct {
	app.Compo
}

// ---------------------------------------------------------------------------
//
// View
//
// ---------------------------------------------------------------------------

func (r *randomSelected) Render() (ui app.UI) {
	return app.Div().Body(
		app.Text("The randomly selected endpoints are:"),
		app.Br(),
		app.Text("Start: "+randomStart),
		app.Br(),
		app.Text("Goal: "+randomGoal),
		app.Br(),
		app.Button().
			Text("Next").
			AutoFocus(true).
			OnClick(r.next),
	)
}

// ---------------------------------------------------------------------------
//
// Controller
//
// ---------------------------------------------------------------------------

func (r *randomSelected) next(ctx app.Context, e app.Event) {
	start = randomStart
	goal = randomGoal
	ctx.NewActionWithValue("setupComplete", nil, app.T("start", start), app.T("goal", goal))
}
