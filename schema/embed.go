package schema

import (
	"embed"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
)

//go:embed **/*.cue
var FS embed.FS

// Module returns the module name defined by cue.mod/module.cue
func Module() string {
	ctx := cuecontext.New()

	data, err := FS.ReadFile("cue.mod/module.cue")

	if err != nil {
		panic("could not read masse embedded `cue.mod/module.cue`")
	}

	module, err := ctx.CompileBytes(data).LookupPath(cue.ParsePath("module")).String()

	if err != nil {
		panic("error compiling masse embedded `cue.mod/module.cue`")
	}

	return module
}
