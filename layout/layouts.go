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
