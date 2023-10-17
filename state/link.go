package state

import (
	"context"

	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/phyton/common"
)

type Link struct {
	Source      []common.Glob `json:"source"`
	From        ChainRef      `json:"from"`
	Destination string        `json:"destination"`
	Options     LinkOptions   `json:"optionsValue"`
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

func (ln *Link) Compile(primary llb.State, secondary ChainStates) (llb.State, error) {
	cwd, _ := primary.GetDir(context.TODO())

	state, err := ln.Copy().Compile(llb.Scratch().Dir(cwd), secondary)
	if err != nil {
		return primary, err
	}

	return llb.Merge([]llb.State{primary, state}), nil
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
