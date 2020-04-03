package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// Build provides building of the arguments
func Build(args []string) {
	app := &cli.App{
		Name:  "gocker",
		Usage: "Tiny implementation of docker",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "project",
				Value: "bin",
				Usage: "type of the project. It might be bin or lib",
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "run",
				Usage:  "running of the image",
				Action: run,
			},
			{
				Name:   "images",
				Usage:  "list of images",
				Action: run,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		return
	}
}

// run provides running of the app
func run(c *cli.Context) error {
	return nil
}
