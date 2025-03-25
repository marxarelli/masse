package schema

import (
	"embed"

	"cuelang.org/go/mod/modfile"
	"cuelang.org/go/mod/module"
)

const embeddedProjectModuleName = "masse.local"

var (
	//go:embed *.cue **/*.cue
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

func ProjectModFile(module string) *modfile.File {
	return &modfile.File{
		Module: module,
		Language: &modfile.Language{
			Version: ModFile.Language.Version,
		},
		Source: &modfile.Source{
			Kind: "self",
		},
		Deps: map[string]*modfile.Dep{
			ModuleVersion.Path(): &modfile.Dep{
				Version: ModuleVersion.Version(),
				Default: true,
			},
		},
	}
}

func EmbeddedProjectModFile() *modfile.File {
	return ProjectModFile(embeddedProjectModuleName)
}
