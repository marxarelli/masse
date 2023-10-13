package state

import (
	"encoding/json"

	"github.com/moby/buildkit/client/llb"
	"github.com/pkg/errors"
)

type StateKind string

const (
	NilKind     StateKind = "nil"
	ScratchKind           = "scratch"
	GitKind               = "git"
	ImageKind             = "image"
	LocalKind             = "local"
	CopyKind              = "copy"
	DiffKind              = "diff"
	LinkKind              = "link"
	MergeKind             = "merge"
	ExtendKind            = "extend"
	RunKind               = "run"
	WithKind              = "with"
)

type Compilable interface {
	Compile(primary llb.State, secondary ChainStates) (llb.State, error)
}

type State struct {
	*Scratch `json:",inline"`
	*Git     `json:",inline"`
	*Image   `json:",inline"`
	*Local   `json:",inline"`
	*Copy    `json:",inline"`
	*Diff    `json:",inline"`
	*Link    `json:",inline"`
	*Merge   `json:",inline"`
	*Extend  `json:",inline"`
	*Run     `json:",inline"`
	*With    `json:",inline"`
}

func (state *State) UnmarshalJSON(data []byte) error {
	st := map[string]json.RawMessage{}
	err := json.Unmarshal(data, &st)
	if err != nil {
		return err
	}

	if _, ok := st["scratch"]; ok {
		state.Scratch = &Scratch{Scratch: true}
		return nil
	}

	if _, ok := st["git"]; ok {
		state.Git = &Git{}
		return json.Unmarshal(data, state.Git)
	}

	if _, ok := st["image"]; ok {
		state.Image = &Image{}
		return json.Unmarshal(data, state.Image)
	}

	if _, ok := st["local"]; ok {
		state.Local = &Local{}
		return json.Unmarshal(data, state.Local)
	}

	if _, ok := st["copy"]; ok {
		state.Copy = &Copy{}
		return json.Unmarshal(data, state.Copy)
	}

	if _, ok := st["diff"]; ok {
		state.Diff = &Diff{}
		return json.Unmarshal(data, state.Diff)
	}

	if _, ok := st["link"]; ok {
		state.Link = &Link{}
		return json.Unmarshal(data, state.Link)
	}

	if _, ok := st["merge"]; ok {
		state.Merge = &Merge{}
		return json.Unmarshal(data, state.Merge)
	}

	if _, ok := st["extend"]; ok {
		state.Extend = &Extend{}
		return json.Unmarshal(data, state.Extend)
	}

	if _, ok := st["run"]; ok {
		state.Run = &Run{}
		return json.Unmarshal(data, state.Run)
	}

	if _, ok := st["with"]; ok {
		state.With = &With{}
		return json.Unmarshal(data, state.With)
	}

	return nil
}

func (state *State) Kind() StateKind {
	field, ok := oneof[any](state)

	if ok {
		switch field.(type) {
		case *Scratch:
			return ScratchKind
		case *Git:
			return GitKind
		case *Image:
			return ImageKind
		case *Local:
			return LocalKind
		case *Copy:
			return CopyKind
		case *Diff:
			return DiffKind
		case *Link:
			return LinkKind
		case *Merge:
			return MergeKind
		case *Extend:
			return ExtendKind
		case *Run:
			return RunKind
		case *With:
			return WithKind
		}
	}

	return NilKind
}

func (state *State) AnonymousChains() ([]Chain, bool) {
	cd, ok := oneof[ChainDefiner](state)

	if !ok {
		return []Chain{}, false
	}

	return cd.AnonymousChains()
}

func (state *State) ChainRefs() []ChainRef {
	cr, ok := oneof[ChainReferencer](state)

	if !ok {
		return []ChainRef{}
	}

	return cr.ChainRefs()
}

func (state *State) Compile(primary llb.State, secondary ChainStates) (llb.State, error) {
	c, ok := oneof[Compilable](state)
	if !ok {
		return llb.State{}, errors.Errorf("no compilable state")
	}

	return c.Compile(primary, secondary)
}
