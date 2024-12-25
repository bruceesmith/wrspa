package main

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type timer struct {
	app.Compo
	ctx app.Context
}

func (t *timer) Render() app.UI {
	return app.Div().Body(
		app.Div().Text("I am the timer"),
	)
}

func (t *timer) OnMount(ctx app.Context) {
	t.ctx = ctx
}
