/*
Package daemon is the heart of Go WikiRacing. On the server, it controls the REST server. In the
browser it hosts the game itself.
*/
package daemon

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/bruceesmith/terminator"
	"github.com/bruceesmith/wrspa/go-app/backend/server"
	"github.com/bruceesmith/wrspa/go-app/frontend/game"
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
			slog.Error("initialisation error", "error", err.Error())
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

	slog.Info("gwr server starting")
	go svr.Serve()

	// Wait for SIGTERM
	<-terminator.ShutDown()

	// Wait for all goroutines to stop
	terminator.Wait()
	slog.Info("gwr server exiting")
	return nil
}
