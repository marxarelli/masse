package state

import (
	"fmt"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/solver/pb"
)

type DiffType llb.DiffType

const (
	DiffNone     DiffType = pb.AttrLocalDifferNone
	DiffMetadata          = pb.AttrLocalDifferMetadata
)

type Local struct {
	Name    string       `json:"local"`
	Options LocalOptions `json:"optionsValue"`
}

func (l *Local) Description() string {
	return fmt.Sprintf(
		"[%s]",
		l.Name,
	)
}

func (l *Local) CompileSource(_ ChainStates, constraints ...llb.ConstraintsOpt) (llb.State, error) {
	return llb.Local(
		l.Name,
		append(constraintsTo[llb.LocalOption](constraints), l.Options)...,
	), nil
}

type LocalOptions []*LocalOption

func (opts LocalOptions) SetLocalOption(info *llb.LocalInfo) {
	for _, opt := range opts {
		opt.SetLocalOption(info)
	}
}

type LocalOption struct {
	*Include
	*Exclude
	*FollowPaths
	*SharedKeyHint
	*Differ
	*Constraint
}

func (opt *LocalOption) SetLocalOption(info *llb.LocalInfo) {
	llbOpt, ok := oneof[llb.LocalOption](opt)
	if ok {
		llbOpt.SetLocalOption(info)
	}
}

type FollowPaths struct {
	FollowPaths []string
}

func (fp *FollowPaths) SetLocalOption(info *llb.LocalInfo) {
	llb.FollowPaths(fp.FollowPaths).SetLocalOption(info)
}

type SharedKeyHint struct {
	SharedKeyHint string
}

func (skh *SharedKeyHint) SetLocalOption(info *llb.LocalInfo) {
	llb.SharedKeyHint(skh.SharedKeyHint).SetLocalOption(info)
}

type Differ struct {
	Differ  DiffType
	Require bool
}

func (diff *Differ) SetLocalOption(info *llb.LocalInfo) {
	llb.Differ(llb.DiffType(diff.Differ), diff.Require).SetLocalOption(info)
}
