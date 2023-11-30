package main

import (
	"context"
	"os"

	"github.com/moby/buildkit/client/llb"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"gitlab.wikimedia.org/dduvall/phyton/config"
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
		&cli.StringFlag{
			Name:    "platform",
			Aliases: []string{"p"},
			Value:   "linux/amd64",
		},
	},
}

func solveAction(clicontext *cli.Context) error {
	file := clicontext.String("file")
	targetName := clicontext.String("target")
	platform := clicontext.String("platform")

	data, err := os.ReadFile(file)
	if err != nil {
		return errors.Wrap(err, "failed to read config file")
	}

	root, err := config.Load(file, data)
	if err != nil {
		return errors.Wrap(err, "failed to load config")
	}

	target, ok := root.Targets[targetName]
	if !ok {
		return errors.Wrapf(err, "unknown target %q", targetName)
	}

	graph, err := target.Graph(root.Chains)
	if err != nil {
		return errors.Wrap(err, "failed to get target graph")
	}

	solver, err := target.ResolvePlatformSolver(platform)
	if err != nil {
		return err
	}

	st, err := solver.Solve(graph)
	if err != nil {
		return errors.Wrap(err, "failed to solve target graph")
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
