package state

import "gitlab.wikimedia.org/dduvall/phyton/common"

type Link struct {
	Source      []common.Glob `json:"source"`
	From        ChainRef      `json:"from"`
	Destination string        `json:"destination"`
	Options     []*LinkOption `json:"options"`
}

type LinkOption struct {
	*common.Creation      `json:",inline"`
	*common.User          `json:",inline"`
	*common.Group         `json:",inline"`
	*common.Mode          `json:",inline"`
	*common.Include       `json:",inline"`
	*common.Exclude       `json:",inline"`
	*CopyDirectoryContent `json:",inline"`
}
