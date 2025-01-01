package setup

import (
	"github.com/bruceesmith/echidna/logger"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// ---------------------------------------------------------------------------
//
// Model
//
// ---------------------------------------------------------------------------

type customSelected struct {
	app.Compo
}

// ---------------------------------------------------------------------------
//
// View
//
// ---------------------------------------------------------------------------

func (c *customSelected) Render() app.UI {
	return app.Div().Body(
		app.Text("Enter the custom endpoints for this game:"),
		app.Br(),
		app.Text("(only Wikipedia subjects, not URLs)"),
		app.Br(),
		app.Label().
			For("start").
			Text("Start"),
		app.Input().
			Type("text").
			ID("start").
			Name("start").
			OnChange(c.ValueTo(&start)),
		app.Label().
			For("goal").
			Text("Goal"),
		app.Input().
			Type("text").
			ID("goal").
			Name("goal").
			OnChange(c.ValueTo(&goal)),
		// app.Input().
		// 	Type("submit").
		// 	Value("Next").
		// 	Disabled(true).
		// 	OnClick(c.next),
		c.nextstep(),
	)
}

func (c *customSelected) nextstep() app.UI {
	if len(start) == 0 || len(goal) == 0 {
		return app.Input().
			Type("submit").
			Value("Next").
			Disabled(true).
			OnClick(c.next)
	} else {
		return app.Input().
			Type("submit").
			Value("Next").
			OnClick(c.next)
	}
}

// ---------------------------------------------------------------------------
//
// Controller
//
// ---------------------------------------------------------------------------

func (c *customSelected) next(ctx app.Context, e app.Event) {
	if len(start) > 0 && len(goal) > 0 {
		ctx.NewActionWithValue("setupComplete", nil, app.T("start", start), app.T("goal", goal))
	} else {
		logger.Info("customSelected one or both fields not filled")
	}
}
