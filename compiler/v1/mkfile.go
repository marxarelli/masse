package v1

import (
	"cuelang.org/go/cue"
	"github.com/moby/buildkit/client/llb"
	"gitlab.wikimedia.org/dduvall/masse/internal/lookup"
)

func (c *compiler) compileMkfile(kind FileActionKind, state llb.State, fa *llb.FileAction, v cue.Value) (*llb.FileAction, error) {
	path, err := lookup.String(v, string(kind))
	if err != nil {
		return fa, vError(v, err)
	}

	options, err := lookup.DecodeOptions[MkfileOptions](v)
	if err != nil {
		return fa, vError(v, err)
	}

	mode := DefaultDirMode()
	if kind == MkfileKind {
		mode = DefaultFileMode()
	}

	for _, opt := range options {
		if opt.Mode != nil {
			mode = opt.Mode
		}
	}

	path = c.absPath(state, path)
	fileMode := mode.FileMode()

	if kind == MkfileKind {
		content, err := lookup.Bytes(v, "content")
		if err != nil {
			return fa, vError(v, err)
		}

		return fa.Mkfile(path, fileMode, content, options), nil
	}

	return fa.Mkdir(path, fileMode, options), nil
}

type MkfileOptions []*MkfileOption

func (opts MkfileOptions) SetMkfileOption(info *llb.MkfileInfo) {
	for _, opt := range opts {
		opt.SetMkfileOption(info)
	}
}

func (opts MkfileOptions) SetMkdirOption(info *llb.MkdirInfo) {
	for _, opt := range opts {
		opt.SetMkdirOption(info)
	}
}

type MkfileOption struct {
	*CreateParents
	*Creation
	*User
	*Group
	*Mode
}

func (opt *MkfileOption) SetMkfileOption(info *llb.MkfileInfo) {
	withOneOf(opt, func(mfo llb.MkfileOption) { mfo.SetMkfileOption(info) })
}

func (opt *MkfileOption) SetMkdirOption(info *llb.MkdirInfo) {
	withOneOf(opt, func(mfo llb.MkdirOption) { mfo.SetMkdirOption(info) })
}

type CreateParents struct {
	CreateParents bool `json:"createParents"`
}

func (opt *CreateParents) SetMkdirOption(info *llb.MkdirInfo) {
	llb.WithParents(opt.CreateParents).SetMkdirOption(info)
}
