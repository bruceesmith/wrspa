package main

import (
	"errors"
	"os"

	"github.com/bruceesmith/echidna"
	"github.com/bruceesmith/echidna/logger"
	"github.com/bruceesmith/echidna/terminator"
	"github.com/bruceesmith/go-wikiracing/backend/daemon"
	"github.com/spf13/pflag"
)

func main() {
	daemon, err := daemon.New()
	if err != nil {
		if errors.Is(err, pflag.ErrHelp) || errors.Is(err, echidna.ErrVersion) {
			os.Exit(0)
		}
		logger.Error("initialisation error", "error", err.Error())
		os.Exit(1)

	}
	daemon.Start()
	<-terminator.ShutDown()
	daemon.Terminate()
}
