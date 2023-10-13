package layout

import (
	"gitlab.wikimedia.org/dduvall/phyton/common"
	"gitlab.wikimedia.org/dduvall/phyton/state"
)

type Layouts map[string]*Layout

type Layout struct {
	Comprises     []state.ChainRef
	Authors       []*Author
	Platforms     []*common.Platform
	Parameters    *Parameters
	Configuration *ImageConfig
}

// Merge returns a new [state.Merge] for the layout's comprised chains.
func (layout *Layout) Merge() *state.Merge {
	return &state.Merge{Merge: layout.Comprises}
}
