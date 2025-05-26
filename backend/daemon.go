package wrserver

import (
	"context"
	"fmt"

	"github.com/bruceesmith/logger"
	"github.com/bruceesmith/terminator"
	"github.com/urfave/cli/v3"
)

func Daemon(ctx context.Context, cmd *cli.Command) error {
	var (
		err error
		svr *Server
	)

	svr, err = NewServer(cmd.String("port"), cmd.String("static"))
	if err != nil {
		logger.Error("initialisation error", "error", err.Error())
		err = fmt.Errorf("initialisation error: [%w]", err)
		return err
	}

	logger.Info("wr server starting")
	go svr.Serve()

	// Wait for SIGTERM
	<-terminator.ShutDown()

	// Wait for all goroutines to stop
	terminator.Wait()
	logger.Info("wr server exiting")
	return nil
}
