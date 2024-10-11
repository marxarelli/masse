package config

import (
	"cuelang.org/go/cue"
	"gitlab.wikimedia.org/dduvall/masse/internal/lookup"
	"gitlab.wikimedia.org/dduvall/masse/target"
)

type Root struct {
	Parameters cue.Value
	Targets    target.Targets `json:"targets"`
	Chains     map[string]cue.Value
}

// UnmarshalCUE parses the given cue.Value into the root config.
func (root *Root) UnmarshalCUE(rv cue.Value) error {
	params, err := lookup.Existing(rv, "parameters")
	if err != nil {
		return err
	}

	root.Parameters = params
	root.Targets = target.Targets{}
	root.Chains = map[string]cue.Value{}

	chains, err := lookup.Existing(rv, "chains")
	if err != nil {
		return err
	}

	err = chains.Decode(&root.Chains)
	if err != nil {
		return err
	}

	return lookup.EachField(
		rv,
		"targets",
		func(label string, tv cue.Value) error {
			target := &target.Target{}
			err := target.UnmarshalCUE(tv)
			if err != nil {
				return err
			}

			root.Targets[label] = target

			return nil
		},
	)
}
