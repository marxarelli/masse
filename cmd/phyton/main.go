package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "phyton",
		Usage:   "Phyton container image builder",
		Version: "0.0.1",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Value:   "layout.cue",
			},
		},
		Commands: []*cli.Command{
			drawCommand,
			solveCommand,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
