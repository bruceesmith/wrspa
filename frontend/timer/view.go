package timer

import "github.com/maxence-charriere/go-app/v10/pkg/app"

func (t *Timer) Render() app.UI {
	return app.Div().Body(
		app.Div().Text("I am the timer"),
	)
}
