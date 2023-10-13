package main

import (
	"context"
	"os"

	"github.com/moby/buildkit/client/llb"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"gitlab.wikimedia.org/dduvall/phyton/layout"
	"gitlab.wikimedia.org/dduvall/phyton/state"
)

var solveCommand = &cli.Command{
	Name:    "solve",
	Aliases: []string{"s"},
	Usage:   "solve -t {target}",
	Action:  solveAction,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "target",
			Aliases:  []string{"t"},
			Required: true,
		},
	},
}

func solveAction(clicontext *cli.Context) error {
	file := clicontext.String("file")
	target := clicontext.String("target")

	data, err := os.ReadFile(file)
	if err != nil {
		return errors.Wrap(err, "failed to read layout file")
	}

	root, err := layout.Load(file, data)
	if err != nil {
		return errors.Wrap(err, "failed to load layout")
	}

	graph, err := root.LayoutGraph(target)
	if err != nil {
		return errors.Wrap(err, "failed to construct graph from layout")
	}

	st, err := state.Solve(graph)
	if err != nil {
		return errors.Wrap(err, "failed to solve layout graph")
	}

	def, err := st.Marshal(context.TODO())
	if err != nil {
		return errors.Wrap(err, "failed to marshal LLB state")
	}

	err = llb.WriteTo(def, os.Stdout)
	if err != nil {
		return errors.Wrap(err, "failed to output LLB definition")
	}

	return nil
}
