package main

import (
	"context"
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/bruceesmith/echidna"
	wrserver "github.com/bruceesmith/wrspa/backend/wrserver"
	"github.com/urfave/cli/v3"
)

func main() {
	var cmd = &cli.Command{
		Name:        "wrserver",
		Action:      wrserver.Daemon,
		Description: "Wiki Racing Server",
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
			&cli.StringFlag{
				Name:  "static",
				Usage: "path to the static SPA flags (CSS, MJS, ...)",
				Validator: func(p string) error {
					if len(p) == 0 {
						return fmt.Errorf("static file path cannot be empty")
					}
					return nil
				},
			},
		},
		Usage:   "Server for Wiki Racing",
		Version: "1.0",
	}

	echidna.Run(
		context.Background(),
		cmd,
	)
}
