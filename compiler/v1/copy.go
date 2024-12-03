package v1

import (
	"cuelang.org/go/cue"
	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/masse/internal/lookup"
)

func (c *compiler) compileCopy(state llb.State, fa *llb.FileAction, v cue.Value) (*llb.FileAction, error) {
	src, err := lookup.String(v, "copy")
	if err != nil {
		return fa, vError(v, err)
	}

	dest, err := lookup.String(v, "destination")
	if err != nil {
		return fa, vError(v, err)
	}

	from, err := lookup.Existing(v, "from")
	if err != nil {
		return fa, vError(v, err)
	}

	fromState, err := c.compileChainByRef(from)
	if err != nil {
		return fa, err
	}

	options, err := lookup.DecodeOptions[CopyOptions](v)
	if err != nil {
		return fa, vError(v, err)
	}

	return fa.Copy(
		fromState,
		c.absPath(fromState, src),
		c.absPath(state, dest),
		options,
	), nil
}

type CopyOptions []*CopyOption

func (opts CopyOptions) SetCopyOption(info *llb.CopyInfo) {
	for _, opt := range opts {
		opt.SetCopyOption(info)
	}
}

type CopyOption struct {
	*Creation
	*User
	*Group
	*Mode
	*Include
	*Exclude
	*FollowSymlinks
	*CopyDirectoryContents
}

func (opt *CopyOption) SetCopyOption(info *llb.CopyInfo) {
	withOneOf(opt, func(co llb.CopyOption) { co.SetCopyOption(info) })
}

type FollowSymlinks struct {
	FollowSymlinks bool `json:"followSymlinks"`
}

func (fs *FollowSymlinks) SetCopyOption(info *llb.CopyInfo) {
	info.FollowSymlinks = fs.FollowSymlinks
}

type CopyDirectoryContents struct {
	CopyDirectoryContents bool `json:"copyDirectoryContents"`
}

func (cdc *CopyDirectoryContents) SetCopyOption(info *llb.CopyInfo) {
	info.CopyDirContentsOnly = cdc.CopyDirectoryContents
}
