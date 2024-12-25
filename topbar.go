package main

import "github.com/maxence-charriere/go-app/v10/pkg/app"

type topbar struct {
	app.Compo
	ctx app.Context
}

func (t *topbar) Render() app.UI {
	return app.Div().Body(
		app.Div().Text("I am the top bar"),
	)
}

func (t *topbar) OnMount(ctx app.Context) {
	t.ctx = ctx
}
