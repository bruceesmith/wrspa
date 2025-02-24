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

func (t *typeSelector) button(label string, value gametype) app.HTMLButton {
	return app.Button().Text(label).
		Class("gwr-ts-text-2").
		Value(value).
		OnClick(t.selectType(value))
}

func (t *typeSelector) selectors() app.UI {
	return app.Div().Body(
		t.button("Custom", custom),
		t.button("Random", random),
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
