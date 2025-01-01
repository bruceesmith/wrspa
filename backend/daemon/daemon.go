package daemon

import (
	"github.com/bruceesmith/echidna/logger"
	"github.com/bruceesmith/echidna/terminator"
	"github.com/bruceesmith/echidna/vpr"
	"github.com/bruceesmith/go-wikiracing/backend/server"
	"github.com/bruceesmith/go-wikiracing/frontend/game"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// Daemon is the controller for the server-side execution
type daemon struct {
	server *server.Server
}

type config struct{}

var cfg config

func New() (d *daemon, err error) {
	if app.IsServer {
		err = vpr.Init(program, version, &cfg, "", nil)
		if err != nil {
			return
		}
		d = &daemon{}
		d.server, err = server.New()
		if err != nil {
			logger.Error(err.Error())
		}
	}
	return
}

func (d *daemon) Start() {
	logger.Info(program + " starting")
	app.Route(
		"/",
		func() app.Composer {
			return game.New()
		},
	)
	app.RunWhenOnBrowser()

	go d.server.Serve()
}

// Terminate shuts down the daemon
func (d *daemon) Terminate() {
	// Wait for all the independent goroutines to stop
	terminator.Wait()
	logger.Info(program + " exiting")
}
