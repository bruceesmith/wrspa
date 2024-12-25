package main

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// GameState is the state of the WikiRacing finite state machine
type GameState int

const (
	loading GameState = iota
	ready
	launching
)

//go:generate stringer --type GameState

type game struct {
	app.Compo
	ctx   app.Context
	state GameState
	Name  string
}

func newGame() (h *game) {
	h = &game{
		state: loading,
		Name:  "fred",
	}
	return
}

func (g *game) Render() app.UI {
	return app.Div().Body(
		&topbar{},
		&wiki{},
		&timer{},
	)
}

func (g *game) OnMount(ctx app.Context) {
	g.ctx = ctx
}
