package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/dominikbraun/graph/draw"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"gitlab.wikimedia.org/dduvall/phyton/config"
)

var drawCommand = &cli.Command{
	Name:    "draw",
	Aliases: []string{"d"},
	Usage:   "draw -t {target}",
	Action:  drawAction,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "target",
			Aliases:  []string{"t"},
			Required: true,
		},
		&cli.StringFlag{
			Name:    "colorscheme",
			Aliases: []string{"c"},
			Value:   "ylgn9",
		},
	},
}

func drawAction(clicontext *cli.Context) error {
	file := clicontext.String("file")
	target := clicontext.String("target")

	data, err := os.ReadFile(file)
	if err != nil {
		return errors.Wrap(err, "failed to read config file")
	}

	root, err := config.Load(file, data)
	if err != nil {
		return errors.Wrap(err, "failed to load config")
	}

	graph, err := root.TargetGraph(target)
	if err != nil {
		return errors.Wrap(err, "failed to construct target graph")
	}

	var buffer bytes.Buffer
	err = draw.DOT(graph.Graph, &buffer)
	if err != nil {
		return err
	}

	digraph := strings.SplitN(buffer.String(), "{", 2)

	io.WriteString(os.Stdout, digraph[0]+"{")
	io.WriteString(
		os.Stdout,
		fmt.Sprintf("\n\tnode [colorscheme=%s];\n", clicontext.String("colorscheme")),
	)
	io.WriteString(os.Stdout, digraph[1])

	return nil
}
