package setup

import (
	"github.com/bruceesmith/wrspa/go-app/frontend/observables"
	"github.com/bruceesmith/logger"
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
	next := app.Button().Text("Next").
		OnClick(c.next).
		Class("gwr-custom-next-step")
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
			Text("Choose the custom endpoints for this game:").
			Class("gwr-custom-text-1"),
		app.P().
			Text("(only Wikipedia subjects, not URLs)").
			Class("gwr-custom-text-2"),
		app.Textarea().
			// Label("Start").
			Placeholder("Starting topic").
			Required(true).
			OnChange(c.ValueTo(&c.start)),
		app.Textarea().
			// Label("Goal").
			Placeholder("Target (goal) topic").
			Required(true).
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
		ctx.SetState(observables.GameSelected, tags)
	} else {
		logger.Info("customSelected one or both fields not filled")
	}
}
