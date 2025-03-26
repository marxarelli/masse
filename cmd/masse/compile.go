package main

import (
	"os"

	"github.com/moby/buildkit/client/llb"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"gitlab.wikimedia.org/dduvall/masse/common"
	v1compiler "gitlab.wikimedia.org/dduvall/masse/compiler/v1"
	"gitlab.wikimedia.org/dduvall/masse/config"
	"gitlab.wikimedia.org/dduvall/masse/load"
)

var compileCommand = &cli.Command{
	Name:    "compile",
	Aliases: []string{"s"},
	Usage:   "compile -t {target}",
	Action:  compileAction,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "no-cache",
			Aliases: []string{"n"},
		},
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

func compileAction(clicontext *cli.Context) error {
	file := clicontext.String("file")
	targetName := clicontext.String("target")
	platformName := clicontext.String("platform")
	noCache := clicontext.Bool("no-cache")

	platform, err := common.ParsePlatform(platformName)
	if err != nil {
		return errors.Wrapf(err, "failed to parse platform %q", platformName)
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return errors.Wrap(err, "failed to read config file")
	}

	root, err := config.Load(file, data, nil, load.WithNearestModFile())
	if err != nil {
		return errors.Wrapf(err, "failed to load config %q", file)
	}

	target, ok := root.Targets[targetName]
	if !ok {
		return errors.Wrapf(err, "unknown target %q", targetName)
	}

	compiler := v1compiler.New(
		root.Chains,
		v1compiler.WithPlatform(platform),
		v1compiler.WithContext(clicontext.Context),
		v1compiler.WithIgnoreCache(noCache),
	)

	result, err := compiler.Compile(target)
	if err != nil {
		return errors.Wrapf(err, "failed to compile target %q", targetName)
	}

	def, err := result.ChainState().Marshal(clicontext.Context)
	if err != nil {
		return errors.Wrap(err, "failed to marshal LLB state")
	}

	err = llb.WriteTo(def, os.Stdout)
	if err != nil {
		return errors.Wrap(err, "failed to output LLB definition")
	}

	return nil
}
