package cuemod

import (
	"embed"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
)

//go:embed *.cue
var FS embed.FS

func Module() string {
	ctx := cuecontext.New()

	data, err := FS.ReadFile("module.cue")

	if err != nil {
		panic("could not read phyton embedded module.cue")
	}

	module, err := ctx.CompileBytes(data).LookupPath(cue.ParsePath("module")).String()

	if err != nil {
		panic("error reading phyton embedded module.cue")
	}

	return module
}
