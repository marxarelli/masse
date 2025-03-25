package v1

import (
	"fmt"

	"cuelang.org/go/cue"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/solver/pb"
	"gitlab.wikimedia.org/dduvall/masse/internal/lookup"
)

func (c *compiler) compileLocal(state llb.State, v cue.Value) (llb.State, error) {
	local, err := lookup.String(v, "local")
	if err != nil {
		return state, vError(v, err)
	}

	options, err := lookup.DecodeOptions[LocalOptions](v)
	if err != nil {
		return state, vError(v, err)
	}

	return llb.Local(local, c.sourceConstraints(), options), nil
}

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

type LocalOptions []LocalOption

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

func (opt LocalOption) SetLocalOption(info *llb.LocalInfo) {
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
