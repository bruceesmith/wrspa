package setup

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// ---------------------------------------------------------------------------
//
// Model
//
// ---------------------------------------------------------------------------

type typeSelector struct{}

// ---------------------------------------------------------------------------
//
// View
//
// ---------------------------------------------------------------------------

func (t *typeSelector) view() []app.UI {
	return []app.UI{
		app.P().Text("Choose the type of game").
			Class("gwr-ts-text-1"),
		t.selectors(),
	}
}

func (t *typeSelector) button(label string, value gametype) (u app.UI) {
	switch value {
	case custom:
		u = app.Button().Text(label).
			Class("gwr-ts-text-2").
			Value(value).
			OnClick(t.selectType(value))
	case random:
		u = app.Button().Text(label).
			Class("gwr-ts-text-2").
			Value(value).
			OnClick(t.selectType(value))
	}
	return u
}

func (t *typeSelector) selectors() app.UI {
	return app.Div().Body(
		t.button("Custom", custom), // Custom should be a Filled Button
		t.button("Random", random), // Random should be an Outlined Button
	).
		Class("gwr-ts-selector")
}

// ---------------------------------------------------------------------------
//
// Controller
//
// ---------------------------------------------------------------------------

func (t *typeSelector) selectType(tipe gametype) func(ctx app.Context, e app.Event) {
	return func(ctx app.Context, e app.Event) {
		ctx.SetState("gameTypeSelected", tipe)
	}
}
