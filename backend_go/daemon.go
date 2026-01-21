package wrserver

import (
	"context"
	"fmt"

	"github.com/bruceesmith/logger"
	"github.com/bruceesmith/terminator"
	"github.com/urfave/cli/v3"
)

func daemon(svr ServerInterface, t *terminator.Terminator) error {
	logger.Info("wr server starting")
	go svr.Serve(t)

	// Wait for SIGTERM
	<-t.ShutDown()

	// Wait for all goroutines to stop
	t.Wait()
	logger.Info("wr server exiting")
	return nil

}

const (
	portFlag   = "port"
	staticFlag = "static"
)

func Daemon(ctx context.Context, cmd *cli.Command) error {
	svr, err := newServerAdapter(cmd.String(portFlag), cmd.String(staticFlag), newClientAdapter())
	if err != nil {
		return fmt.Errorf("failed to create server adapter: %w", err)
	}
	return daemon(svr, terminator.New())
}
