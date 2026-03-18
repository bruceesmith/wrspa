/*
Package daemon is the heart of Go WikiRacing. On the server, it controls the REST server. In the
browser it hosts the game itself.
*/
package daemon

import (
	"context"
	"fmt"

	"github.com/bruceesmith/wrspa/go-app/backend/server"
	"github.com/bruceesmith/wrspa/go-app/frontend/game"
	"github.com/bruceesmith/logger"
	"github.com/bruceesmith/terminator"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/urfave/cli/v3"
)

func Daemon(ctx context.Context, cmd *cli.Command) error {
	var (
		err error
		svr *server.Server
	)
	if app.IsServer {
		svr, err = server.New(cmd.String("port"))
		if err != nil {
			logger.Error("initialisation error", "error", err.Error())
			err = fmt.Errorf("initialisation error: [%w]", err)
			return err
		}
	}

	app.Route(
		"/",
		func() app.Composer {
			return game.New()
		},
	)
	app.RunWhenOnBrowser()

	// Following code is only executed on the server, never in the browser

	logger.Info("gwr server starting")
	go svr.Serve()

	// Wait for SIGTERM
	<-terminator.ShutDown()

	// Wait for all goroutines to stop
	terminator.Wait()
	logger.Info("gwr server exiting")
	return nil
}
