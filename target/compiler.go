package target

import (
	"context"

	"cuelang.org/go/cue"
	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/masse/common"
)

type Compiler interface {
	Compile(*Target) (CompilerResult, error)
	CompileChain(cue.Value) (llb.State, error)
	CompileState(llb.State, cue.Value) (llb.State, error)
	Error() error
	WithContext(context.Context) Compiler
}

type CompilerResult interface {
	ChainRef() string
	ChainState() llb.State
	Platform() common.Platform
	DependencyChainStates() map[string]llb.State
}
