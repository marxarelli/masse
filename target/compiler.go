package target

import (
	"cuelang.org/go/cue"
	"github.com/moby/buildkit/client/llb"
)

type Compiler interface {
	Compile(*Target) (llb.State, error)
	CompileState(llb.State, cue.Value) (llb.State, error)
	Error() error
}
