package state

import "github.com/moby/buildkit/client/llb"

type SourceMount struct {
	Target string `json:"mount"`
	From   ChainRef
	Source string
}

func (sm *SourceMount) LLBRunOptions(states ChainStates) ([]llb.RunOption, error) {
	from, err := states.Resolve(sm.From)
	if err != nil {
		return nil, err
	}

	return []llb.RunOption{
		llb.AddMount(
			sm.Target,
			from,
			llb.SourcePath(sm.Source),
		),
	}, nil
}
