package timer

import "github.com/maxence-charriere/go-app/v10/pkg/app"

func (t *Timer) OnMount(ctx app.Context) {
	t.ctx = ctx
}
