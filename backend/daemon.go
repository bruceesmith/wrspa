package wrserver

import (
	"context"

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

func Daemon(ctx context.Context, cmd *cli.Command) error {
	svr, err := newServerAdapter(cmd.String("port"), cmd.String("static"), newClientAdapter())
	if err != nil {
		return err
	}
	return daemon(svr, terminator.New())
}
