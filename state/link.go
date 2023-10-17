package state

import (
	"context"
	"fmt"

	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/phyton/common"
)

type Link struct {
	Source      []common.Glob `json:"source"`
	From        ChainRef      `json:"from"`
	Destination string        `json:"destination"`
	Options     LinkOptions   `json:"optionsValue"`
}

func (ln *Link) Description() string {
	return fmt.Sprintf("∪ %s", ln.Copy().Description())
}

func (ln *Link) ChainRefs() []ChainRef {
	return []ChainRef{ln.From}
}

func (ln *Link) Copy() *Copy {
	return &Copy{
		Source:      ln.Source,
		From:        ln.From,
		Destination: ln.Destination,
		Options:     ln.Options.CopyOptions(),
	}
}

func (ln *Link) Compile(primary llb.State, secondary ChainStates, constraints ...llb.ConstraintsOpt) (llb.State, error) {
	cwd, _ := primary.GetDir(context.TODO())

	cp := ln.Copy()
	state, err := cp.Compile(
		llb.Scratch().Dir(cwd),
		secondary,
		append(constraints, llb.WithCustomName(
			(&State{Copy: cp}).Description(),
		))...,
	)

	if err != nil {
		return primary, err
	}

	return llb.Merge([]llb.State{primary, state}, constraints...), nil
}

type LinkOptions []*LinkOption

func (opts LinkOptions) CopyOptions() CopyOptions {
	copts := make(CopyOptions, len(opts))
	for i, opt := range opts {
		copts[i] = opt.CopyOption()
	}
	return copts
}

type LinkOption struct {
	*Creation             `json:",inline"`
	*User                 `json:",inline"`
	*Group                `json:",inline"`
	*Mode                 `json:",inline"`
	*Include              `json:",inline"`
	*Exclude              `json:",inline"`
	*CopyDirectoryContent `json:",inline"`
}

func (opt *LinkOption) CopyOption() *CopyOption {
	return &CopyOption{
		Creation:             opt.Creation,
		User:                 opt.User,
		Group:                opt.Group,
		Mode:                 opt.Mode,
		Include:              opt.Include,
		Exclude:              opt.Exclude,
		CopyDirectoryContent: opt.CopyDirectoryContent,
	}
}
