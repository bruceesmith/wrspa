package setup

import "github.com/maxence-charriere/go-app/v10/pkg/app"

// ---------------------------------------------------------------------------
//
// Model
//
// ---------------------------------------------------------------------------

type gametype string

const (
	custom gametype = "custom"
	random gametype = "random"
	unset  gametype = "unset"
)

type typeSelector struct {
	app.Compo
}

// ---------------------------------------------------------------------------
//
// View
//
// ---------------------------------------------------------------------------

func (t *typeSelector) Render() app.UI {
	return app.Div().Body(
		app.P().
			Body(
				app.Text("Choose the type of game"),
			),
		app.Input().
			Type("radio").
			ID("custom").
			Name("game_type").
			Value(string(custom)).
			OnChange(t.selectType),
		app.Label().
			For("custom").
			Text("Custom"),
		app.Input().
			Type("radio").
			ID("random").
			Name("game_type").
			Value(string(random)).
			OnChange(t.selectType),
		app.Label().
			For("random").
			Text("Random"),
	)
}

// ---------------------------------------------------------------------------
//
// Controller
//
// ---------------------------------------------------------------------------

func (t *typeSelector) selectType(ctx app.Context, e app.Event) {
	ctx.NewActionWithValue("gameTypeSelected", gametype(ctx.JSSrc().Get("value").String()))
}
