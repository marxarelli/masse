package state

import "github.com/moby/buildkit/client/llb"

type Scratch struct {
	Scratch bool `json:"scratch"`
}

func (scratch *Scratch) Compile(_ llb.State, _ ChainStates) (llb.State, error) {
	return llb.Scratch(), nil
}
