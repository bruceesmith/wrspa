package wiki

import (
	"fmt"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// ---------------------------------------------------------------------------
//
// Model
//
// ---------------------------------------------------------------------------

type Timer struct {
	app.Compo
	ticker  *time.Ticker
	done    chan bool
	elapsed time.Duration
	Value   string
}

var (
	tmr = Timer{
		Value: "00:00:00",
	}
)

// ---------------------------------------------------------------------------
//
// View
//
// ---------------------------------------------------------------------------

func (t *Timer) Render() app.UI {
	return app.Div().Body(
		app.Text(
			"time is "+t.Value,
		),
	).
		Style("display", "grid").
		Style("place-content", "center")
}

// ---------------------------------------------------------------------------
//
// Controller
//
// ---------------------------------------------------------------------------

func (t *Timer) OnMount(ctx app.Context) {
	tmr.ticker = time.NewTicker(time.Second)
	tmr.done = make(chan bool)
	ctx.ObserveState("time", &t.Value)
	ctx.Async(
		func() {
			defer t.ticker.Stop()
			for {
				select {
				case <-t.done:
					break
				case <-t.ticker.C:
					switch Default.State {
					case ready, paused, finished:
					case playing:
						t.elapsed += time.Second
						ctx.SetState(
							"time",
							fmt.Sprintf(
								"%02.0f:%02.0f:%02.0f",
								t.elapsed.Hours(),
								t.elapsed.Minutes(),
								t.elapsed.Seconds(),
							),
						)
					}
				}
			}
		},
	)
}

func (t *Timer) finished() {
	t.done <- true
}
