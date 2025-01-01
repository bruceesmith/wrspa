package game

import (
	"github.com/bruceesmith/go-wikiracing/frontend/setup"
	"github.com/bruceesmith/go-wikiracing/frontend/timer"
	"github.com/bruceesmith/go-wikiracing/frontend/wiki"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// ---------------------------------------------------------------------------
//
// Model
//
// ---------------------------------------------------------------------------

// GameState is the state of the WikiRacing finite state machine
type GameState int

const (
	preparing GameState = iota
	ready
	running
	paused
)

//go:generate stringer --type GameState

type Game struct {
	app.Compo
	ctx         app.Context
	state       GameState
	start, goal string
}

func New() (g *Game) {
	g = &Game{
		state: preparing,
	}
	return
}

// ---------------------------------------------------------------------------
//
// View
//
// ---------------------------------------------------------------------------

func (g *Game) Render() app.UI {
	switch g.state {
	case preparing:
		return &setup.Default
	case ready:
		wiki.Default.Targets(g.start, g.goal)
		return app.Div().Body(
			&wiki.Default,
			&timer.Timer{},
		)
	default:
		return app.Text("unknown game state")
	}
}

// ---------------------------------------------------------------------------
//
// Controller
//
// ---------------------------------------------------------------------------

func (g *Game) OnMount(ctx app.Context) {
	g.ctx = ctx
	ctx.Handle("setupComplete", g.setupComplete)
}

func (g *Game) setupComplete(ctx app.Context, a app.Action) {
	g.state = ready
	g.start = a.Tags.Get("start")
	g.goal = a.Tags.Get("goal")
	g.ctx.Update()
}
