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
			Style("justify-self", "center"),
		t.selectors(),
	}
}

func (t *typeSelector) button(label string, value gametype) app.HTMLButton {
	return app.Button().Text(label).
		Style("background", "#44474E").
		Style("border-radius", "10px").
		Style("color", "#002E68").
		Style("font-size", "20px").
		Value(value).
		OnClick(t.selectType(value))
}

func (t *typeSelector) selectors() app.UI {
	return app.Div().Body(
		t.button("Custom", custom),
		t.button("Random", random),
	).
		Style("display", "grid").
		Style("grid-template-columns", "1fr 1fr").
		Style("gap", "5px")
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
