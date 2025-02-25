package wiki

import (
	"github.com/bruceesmith/go-wikiracing/frontend/observables"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// ---------------------------------------------------------------------------
//
// Model
//
// ---------------------------------------------------------------------------

type Topbar struct {
	app.Compo
	State state
}

var (
	theTopbar = Topbar{}
)

// ---------------------------------------------------------------------------
//
// View
//
// ---------------------------------------------------------------------------

func (t *Topbar) Render() app.UI {
	// https://www.streamlinehq.com
	return app.Div().
		Class("gwr-wiki-topbar").
		Body(
			app.Button().
				Disabled(true).
				Class("gwr-wiki-topbar-back").
				Body(
					app.Img().
						Src("/web/Line-Start-Arrow-Notch-Fill--Streamline-Outlined-Fill-Material.svg"),
					app.Text("Back"),
				),
			app.Button().
				Disabled(true).
				Class("gwr-wiki-topbar-forward").
				Body(
					app.Text("Forward"),
					app.Img().
						Src("/web/Line-End-Arrow-Notch-Fill--Streamline-Outlined-Fill-Material.svg"),
				),
			app.If(
				Default.State == ready || Default.State == paused,
				func() app.UI {
					return app.Button().
						OnClick(t.play).
						Class("gwr-wiki-topbar-playpause").
						AutoFocus(true).
						Body(
							app.Text("Play"),
							app.Img().
								Src("/web/Play-Arrow-Fill--Streamline-Outlined-Fill-Material.svg"),
						)
				},
			).ElseIf(
				Default.State == playing,
				func() app.UI {
					return app.Button().
						OnClick(t.pause).
						Class("gwr-wiki-topbar-playpause").
						Body(
							app.Text("Pause"),
							app.Img().
								Src("/web/Pause-Fill--Streamline-Outlined-Fill-Material.svg"),
						)
				},
			).Else(
				func() app.UI {
					return app.Button().
						OnClick(t.play).
						Class("gwr-wiki-topbar-playpause").
						Disabled(true).
						Body(
							app.Text("Play"),
							app.Img().
								Src("/web/Play-Arrow-Fill--Streamline-Outlined-Fill-Material.svg"),
						)
				},
			),
		)
}

// ---------------------------------------------------------------------------
//
// Controller
//
// ---------------------------------------------------------------------------

func (t *Topbar) OnMount(ctx app.Context) {
	ctx.ObserveState(observables.WikiState, &t.State)
}

func (t *Topbar) pause(ctx app.Context, e app.Event) {
	ctx.SetState(observables.WikiState, paused)
}

func (t *Topbar) play(ctx app.Context, e app.Event) {
	ctx.SetState(observables.WikiState, playing)
}
