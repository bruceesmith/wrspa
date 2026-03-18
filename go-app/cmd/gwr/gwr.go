package main

import (
	"context"
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/bruceesmith/echidna"
	"github.com/bruceesmith/wrspa/go-app/backend/daemon"
	"github.com/urfave/cli/v3"
)

func main() {
	var cmd = &cli.Command{
		Name:        "gwr",
		Action:      daemon.Daemon,
		Description: "Go Wiki Racing",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "port",
				Usage: "port where the server will listen",
				Validator: func(p string) error {
					if !govalidator.IsPort(p) {
						return fmt.Errorf("invalid port %s", p)
					}
					return nil
				},
				Value: "8080",
			},
		},
		Usage:   "Server for Go Wiki Racing",
		Version: "1.0",
	}

	echidna.Run(
		context.Background(),
		cmd,
	)
}
