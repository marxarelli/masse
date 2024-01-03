package cuemod

import (
	"embed"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
)

//go:embed *.cue
var FS embed.FS

// Module returns the module name defined by cue.mod/module.cue
func Module() string {
	ctx := cuecontext.New()

	data, err := FS.ReadFile("module.cue")

	if err != nil {
		panic("could not read masse embedded module.cue")
	}

	module, err := ctx.CompileBytes(data).LookupPath(cue.ParsePath("module")).String()

	if err != nil {
		panic("error reading masse embedded module.cue")
	}

	return module
}
