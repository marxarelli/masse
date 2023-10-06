package layout

import "gitlab.wikimedia.org/dduvall/phyton/state"

type Root struct {
	Parameters Parameters   `json:"parameters"`
	Chains     state.Chains `json:"chains"`
	Layouts    Layouts      `json:"layouts"`
}
