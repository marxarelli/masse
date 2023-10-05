package state

import "gitlab.wikimedia.org/dduvall/phyton/common"

type Image struct {
	Ref     string `json:"image"`
	Options []*ImageOption
}

type ImageOption struct {
	*Platform
	*LayerLimit
}

type Platform struct {
	Platform common.Platform `json:"value"`
}

type LayerLimit struct {
	LayerLimit uint32
}
