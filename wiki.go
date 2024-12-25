package main

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type wiki struct {
	app.Compo
	ctx app.Context
}

func (w *wiki) Render() app.UI {
	return app.Div().Body(
		app.Div().Text("I am the wiki"),
		app.Raw(`<a href="https://www.w3schools.com" target="_blank" onclick="wikiAnchorClick(event, this)">Visit W3Schools.com!</a>`),
	)
}

func (w *wiki) OnMount(ctx app.Context) {
	w.ctx = ctx
	app.Window().Set("wikiUrlClicked", app.FuncOf(w.urlclick))
}

func (w *wiki) urlclick(this app.Value, args []app.Value) (x any) {
	if len(args) == 0 {
		fmt.Println("onclick handler received no arguments")
		return
	}
	object := args[0]
	if object.Type() != app.TypeObject {
		fmt.Printf("onclick handler received %s rather than an object\n", object.Type().String())
		return
	}
	if object.Get("localName").String() != "a" {
		fmt.Printf("onclick handler called from element type %s\n", object.Get("localName").String())
		return
	}
	// fmt.Printf("onclick handler called, href %s\n", object.Get("href").String())

	w.ctx.Update()
	return
}
