package state

import (
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/solver/pb"
	"gitlab.wikimedia.org/dduvall/phyton/common"
)

type DiffType llb.DiffType

const (
	DiffNone     DiffType = pb.AttrLocalDifferNone
	DiffMetadata          = pb.AttrLocalDifferMetadata
)

type Local struct {
	Name    string `json:"local"`
	Options []*LocalOption
}

type Include common.Include
type Exclude common.Exclude

type LocalOption struct {
	*Include
	*Exclude
	*FollowPaths
	*SharedKeyHint
	*Differ
}

type FollowPaths struct {
	FollowPaths []string
}

type SharedKeyHint struct {
	SharedKeyHint string
}

type Differ struct {
	Differ  DiffType
	Require bool
}
