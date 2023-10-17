package state

import (
	"fmt"

	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/phyton/common"
)

type Copy struct {
	Source      []common.Glob `json:"source"`
	From        ChainRef      `json:"from"`
	Destination string        `json:"destination"`
	Options     CopyOptions   `json:"optionsValue"`
}

func (cp *Copy) Description() string {
	return fmt.Sprintf(
		"{%s}%+v -> %+v",
		cp.From, cp.Source, cp.Destination,
	)
}

func (cp *Copy) CompileSource(secondary ChainStates, constraints ...llb.ConstraintsOpt) (llb.State, error) {
	return cp.Compile(llb.Scratch(), secondary, constraints...)
}

func (cp *Copy) Compile(primary llb.State, secondary ChainStates, constraints ...llb.ConstraintsOpt) (llb.State, error) {
	from, err := secondary.Resolve(cp.From)
	if err != nil {
		return primary, err
	}

	var fa *llb.FileAction

	copyOpts := []llb.CopyOption{
		&llb.CopyInfo{
			FollowSymlinks:      true,
			CopyDirContentsOnly: false,
			AttemptUnpack:       false,
			CreateDestPath:      true,
			AllowWildcard:       true,
			AllowEmptyWildcard:  true,
		},
		cp.Options,
	}

	dest := qualifyStatePath(primary, cp.Destination)
	for _, srcGlob := range cp.Source {
		src := qualifyStatePath(from, srcGlob.String())
		if fa == nil {
			fa = llb.Copy(from, src, dest, copyOpts...)
		} else {
			fa = fa.Copy(from, src, dest, copyOpts...)
		}
	}

	return primary.File(fa, constraints...), nil
}

func (cp *Copy) ChainRefs() []ChainRef {
	return []ChainRef{cp.From}
}

type CopyOptions []*CopyOption

func (opts CopyOptions) SetCopyOption(info *llb.CopyInfo) {
	for _, opt := range opts {
		opt.SetCopyOption(info)
	}
}

type CopyOption struct {
	*Creation             `json:",inline"`
	*User                 `json:",inline"`
	*Group                `json:",inline"`
	*Mode                 `json:",inline"`
	*Include              `json:",inline"`
	*Exclude              `json:",inline"`
	*FollowSymlinks       `json:",inline"`
	*CopyDirectoryContent `json:",inline"`
}

func (opt *CopyOption) SetCopyOption(info *llb.CopyInfo) {
	llbOpt, ok := oneof[llb.CopyOption](opt)
	if ok {
		llbOpt.SetCopyOption(info)
	}
}

type FollowSymlinks struct {
	FollowSymlinks bool `json:"followSymlinks"`
}

func (fs *FollowSymlinks) SetCopyOption(info *llb.CopyInfo) {
	info.FollowSymlinks = fs.FollowSymlinks
}

type CopyDirectoryContent struct {
	CopyDirectoryContent bool `json:"copyDirectoryContent"`
}

func (cdc *CopyDirectoryContent) SetCopyOption(info *llb.CopyInfo) {
	info.CopyDirContentsOnly = cdc.CopyDirectoryContent
}
