package config

import (
	"cuelang.org/go/cue"
	"gitlab.wikimedia.org/dduvall/masse/target"
)

type Root struct {
	Parameters cue.Value            `json:"parameters"`
	Targets    target.Targets       `json:"targets"`
	Chains     map[string]cue.Value `json:"chains"`
}
