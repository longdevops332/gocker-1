package cmd

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/saromanov/gocker/pkg/images"
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
				Action: list,
			},
			{
				Name:   "pull",
				Usage:  "pull of the image",
				Action: pull,
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

func list(c *cli.Context) error {
	images, err := images.List("./images")
	if err != nil {
		logrus.Fatalf("unable to get list of images: %v", err)
		return err
	}
	for _, img := range images {
		fmt.Println(img)
	}
	return nil
}

func pull(c *cli.Context) error {
	err := func(ctx *cli.Context) error {
		img := c.Args().Get(0)
		if img == "" {
			return errors.New("image name is not defined")
		}
		if err := images.NewPull(img).Do(); err != nil {
			return errors.Wrap(err, "unable to pull image")
		}
		return nil
	}(c)
	if err != nil {
		logrus.Fatalf("unable to apply pull of the image: %v", err)
	}
	return nil
}
