package schema

import (
	"embed"

	"cuelang.org/go/mod/modfile"
	"cuelang.org/go/mod/module"
)

var (
	//go:embed **/*.cue
	FS      embed.FS
	Tag     string
	ModFile = MustModFile()

	ModuleVersion = module.MustNewVersion(ModFile.QualifiedModule(), Version())
)

func MustModFile() *modfile.File {
	data, err := FS.ReadFile("cue.mod/module.cue")

	if err != nil {
		panic("could not read masse embedded `cue.mod/module.cue`")
	}

	file, err := modfile.Parse(data, "cue.mod/module.cue")
	if err != nil {
		panic(err)
	}

	return file
}

func Version() string {
	if Tag == "" {
		return ModFile.MajorVersion() + ".999.999-dev"
	}

	return Tag
}
