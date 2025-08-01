package wrserver

import (
	"context"
	"fmt"

	"github.com/bruceesmith/logger"
	"github.com/bruceesmith/terminator"
	"github.com/urfave/cli/v3"
)

// The following variables assist with testing
var (
	newServer         = NewServer
	terminateShutDown = terminator.ShutDown
	terminateWait     = terminator.Wait
)

func Daemon(ctx context.Context, cmd *cli.Command) error {
	var (
		err error
		svr *Server
	)

	svr, err = newServer(cmd.String("port"), cmd.String("static"))
	if err != nil {
		logger.Error("initialisation error", "error", err.Error())
		err = fmt.Errorf("initialisation error: [%w]", err)
		return err
	}

	logger.Info("wr server starting")
	go svr.Serve()

	// Wait for SIGTERM
	<-terminateShutDown()

	// Wait for all goroutines to stop
	terminateWait()
	logger.Info("wr server exiting")
	return nil
}
