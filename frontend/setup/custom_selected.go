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
	start, goal string
}

// ---------------------------------------------------------------------------
//
// View
//
// ---------------------------------------------------------------------------

func (c *customSelected) nextstep() app.UI {
	next := app.Input().
		Type("submit").
		Value("Next").
		OnClick(c.next).
		Style("background", "#ADC6FF").
		Style("border-radius", "10px").
		Style("color", "#002E68").
		Style("font-size", "20px").
		Style("margin-top", "10px")
	if len(c.start) == 0 || len(c.goal) == 0 {
		return next.Disabled(true)
	} else {
		return next
	}
}

func (c *customSelected) view() []app.UI {
	return []app.UI{
		app.Hr(),
		app.P().
			Text("Enter the custom endpoints for this game:").
			Style("justify-self", "center"),
		app.P().
			Text("(only Wikipedia subjects, not URLs)").
			Style("justify-self", "center").
			Style("margin-bottom", "10px"),
		app.Label().
			For("start").
			Text("Start"),
		app.Input().
			Type("text").
			ID("start").
			Name("start").
			OnChange(c.ValueTo(&c.start)),
		app.Label().
			For("goal").
			Text("Goal"),
		app.Input().
			Type("text").
			ID("goal").
			Name("goal").
			OnChange(c.ValueTo(&c.goal)),
		c.nextstep(),
	}
}

// ---------------------------------------------------------------------------
//
// Controller
//
// ---------------------------------------------------------------------------

func (c *customSelected) next(ctx app.Context, e app.Event) {
	if len(c.start) > 0 && len(c.goal) > 0 {
		tags := app.Tags{}
		tags.Set("start", c.start)
		tags.Set("goal", c.goal)
		ctx.SetState("gameSelected", tags)
	} else {
		logger.Info("customSelected one or both fields not filled")
	}
}
